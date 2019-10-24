package builds

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	e2e "k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/framework/pod"

	buildv1 "github.com/openshift/api/build/v1"
	exutil "github.com/openshift/origin/test/extended/util"
	"github.com/openshift/origin/test/extended/util/jenkins"
)

var _ = g.Describe("[Feature:Jenkins][Slow]jenkins pipeline test openshift-client-plugin", func() {
	defer g.GinkgoRecover()

	var (
		jenkinsEphemeralTemplatePath  = exutil.FixturePath("..", "..", "examples", "jenkins", "jenkins-ephemeral-template.json")
		nodejsDeclarativePipelinePath = exutil.FixturePath("..", "..", "examples", "jenkins", "pipeline", "raw-commands-pipeline.yaml")
		oc                            = exutil.NewCLI("jenkins-pipeline", exutil.KubeConfigPath())
		ticker                        *time.Ticker
		j                             *jenkins.JenkinsRef

		cleanup = func(jenkinsTemplatePath string) {
			if g.CurrentGinkgoTestDescription().Failed {
				exutil.DumpPodStates(oc)
				exutil.DumpPodLogsStartingWith("", oc)
				exutil.DumpPersistentVolumeInfo(oc)
			}
			if os.Getenv(jenkins.EnableJenkinsMemoryStats) != "" {
				ticker.Stop()
			}
			client := oc.AsAdmin().KubeFramework().ClientSet
			g.By("removing jenkins")
			exutil.RemoveDeploymentConfigs(oc, "jenkins")			
		}
		setupJenkins = func(jenkinsTemplatePath string) {
			exutil.PreTestDump()
			// Deploy Jenkins
			// NOTE, we use these tests for both a) nightly regression runs against the latest openshift jenkins image on docker hub, and
			// b) PR testing for changes to the various openshift jenkins plugins we support.  With scenario b), a container image that extends
			// our jenkins image is built, where the proposed plugin change is injected, overwritting the current released version of the plugin.
			// Our test/PR jobs on ci.openshift create those images, as well as set env vars this test suite looks for.  When both the env var
			// and test image is present, a new image stream is created using the test image, and our jenkins template is instantiated with
			// an override to use that images stream and test image
			var licensePrefix, pluginName string
			useSnapshotImage := false

			err := oc.Run("create").Args("-n", oc.Namespace(), "-f", jenkinsTemplatePath).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())

			jenkinsTemplateName := "jenkins-ephemeral"
		
			// our pipeline jobs, between jenkins and oc invocations, need more mem than the default
			newAppArgs := []string{"--template", fmt.Sprintf("%s/%s", oc.Namespace(), jenkinsTemplateName), "-p", "MEMORY_LIMIT=2Gi", "-p", "DISABLE_ADMINISTRATIVE_MONITORS=true"}
			newAppArgs = jenkins.OverridePodTemplateImages(newAppArgs)
			clientPluginNewAppArgs, useClientPluginSnapshotImage := jenkins.SetupSnapshotImage(jenkins.UseLocalClientPluginSnapshotEnvVarName, localClientPluginSnapshotImage, localClientPluginSnapshotImageStream, newAppArgs, oc)
			syncPluginNewAppArgs, useSyncPluginSnapshotImage := jenkins.SetupSnapshotImage(jenkins.UseLocalSyncPluginSnapshotEnvVarName, localSyncPluginSnapshotImage, localSyncPluginSnapshotImageStream, newAppArgs, oc)

			switch {
			case useClientPluginSnapshotImage && useSyncPluginSnapshotImage:
				fmt.Fprintf(g.GinkgoWriter,
					"\nBOTH %s and %s for PR TESTING ARE SET.  WILL NOT CHOOSE BETWEEN THE TWO SO TESTING CURRENT PLUGIN VERSIONS IN LATEST OPENSHIFT JENKINS IMAGE ON DOCKER HUB.\n",
					jenkins.UseLocalClientPluginSnapshotEnvVarName, jenkins.UseLocalSyncPluginSnapshotEnvVarName)
			case useClientPluginSnapshotImage:
				fmt.Fprintf(g.GinkgoWriter, "\nTHE UPCOMING TESTS WILL LEVERAGE AN IMAGE THAT EXTENDS THE LATEST OPENSHIFT JENKINS IMAGE AND OVERRIDES THE OPENSHIFT CLIENT PLUGIN WITH A NEW VERSION BUILT FROM PROPOSED CHANGES TO THAT PLUGIN.\n")
				licensePrefix = clientLicenseText
				pluginName = clientPluginName
				useSnapshotImage = true
				newAppArgs = clientPluginNewAppArgs
			case useSyncPluginSnapshotImage:
				fmt.Fprintf(g.GinkgoWriter, "\nTHE UPCOMING TESTS WILL LEVERAGE AN IMAGE THAT EXTENDS THE LATEST OPENSHIFT JENKINS IMAGE AND OVERRIDES THE OPENSHIFT SYNC PLUGIN WITH A NEW VERSION BUILT FROM PROPOSED CHANGES TO THAT PLUGIN.\n")
				licensePrefix = syncLicenseText
				pluginName = syncPluginName
				useSnapshotImage = true
				newAppArgs = syncPluginNewAppArgs
			default:
				fmt.Fprintf(g.GinkgoWriter, "\nNO PR TEST ENV VARS SET SO TESTING CURRENT PLUGIN VERSIONS IN LATEST OPENSHIFT JENKINS IMAGE ON DOCKER HUB.\n")
			}

			g.By(fmt.Sprintf("calling oc new-app useSnapshotImage %v with license text %s and newAppArgs %#v", useSnapshotImage, licensePrefix, newAppArgs))
			err = oc.Run("new-app").Args(newAppArgs...).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())

	
			g.By("waiting for jenkins deployment")
			err = exutil.WaitForDeploymentConfig(oc.KubeClient(), oc.AppsClient().AppsV1(), oc.Namespace(), "jenkins", 1, false, oc)
			if err != nil {
				exutil.DumpApplicationPodLogs("jenkins", oc)
			}
			o.Expect(err).NotTo(o.HaveOccurred())

			j = jenkins.NewRef(oc)

			g.By("wait for jenkins to come up")
			_, err = j.WaitForContent("", 200, 5*time.Minute, "")

			if err != nil {
				exutil.DumpApplicationPodLogs("jenkins", oc)
			}

			o.Expect(err).NotTo(o.HaveOccurred())

			if useSnapshotImage {
				g.By("verifying the test image is being used")
				// for the test image, confirm that a snapshot version of the plugin is running in the jenkins image we'll test against
				_, err = j.WaitForContent(licensePrefix+` ([0-9\.]+)-SNAPSHOT`, 200, 10*time.Minute, "/pluginManager/plugin/"+pluginName+"/thirdPartyLicenses")
				o.Expect(err).NotTo(o.HaveOccurred())
			}

			// Start capturing logs from this deployment config.
			// This command will terminate if the Jenkins instance crashes. This
			// ensures that even if the Jenkins DC restarts, we should capture
			// logs from the crash.
			_, _, _, err = oc.Run("logs").Args("-f", "dc/jenkins").Background()
			o.Expect(err).NotTo(o.HaveOccurred())

			if os.Getenv(jenkins.EnableJenkinsMemoryStats) != "" {
				ticker = jenkins.StartJenkinsMemoryTracking(oc, oc.Namespace())
			}
		}

		debugAnyJenkinsFailure = func(br *exutil.BuildResult, name string, oc *exutil.CLI, dumpMaster bool) {
			if !br.BuildSuccess {
				br.LogDumper = jenkins.DumpLogs
				fmt.Fprintf(g.GinkgoWriter, "\n\n START debugAnyJenkinsFailure\n\n")
				j := jenkins.NewRef(oc)
				jobLog, err := j.GetJobConsoleLogsAndMatchViaBuildResult(br, "")
				if err == nil {
					fmt.Fprintf(g.GinkgoWriter, "\n %s job log:\n%s", name, jobLog)
				} else {
					fmt.Fprintf(g.GinkgoWriter, "\n error getting %s job log: %#v", name, err)
				}
				if dumpMaster {
					exutil.DumpApplicationPodLogs("jenkins", oc)
				}
				fmt.Fprintf(g.GinkgoWriter, "\n\n END debugAnyJenkinsFailure\n\n")
			}
		}
	)

	// these tests are isolated so that PR testing the jenkins-client-plugin can execute the extended
	// tests with a ginkgo focus that runs only the tests within this ginkgo context
	g.Context("Client Plugin tests", func() {
		g.It("using the ephemeral template", func() {
			defer cleanup(jenkinsEphemeralTemplatePath)
			setupJenkins(jenkinsEphemeralTemplatePath)
			g.By("Pipeline containing calls to openshift.raw() function")
			g.By("should perform oc calls and complete successfully", func() {

				// instantiate the template
				g.By(fmt.Sprintf("calling oc new-app -f %q", mavenSlavePipelinePath))
				err := oc.Run("new-app").Args("-f", mavenSlavePipelinePath).Execute()
				o.Expect(err).NotTo(o.HaveOccurred())

				// start the build
				g.By("starting the pipeline build and waiting for it to complete")
				br, err := exutil.StartBuildAndWait(oc, "openshift-jee-sample")
				if err != nil || !br.BuildSuccess {
					exutil.DumpBuilds(oc)
					exutil.DumpPodLogsStartingWith("maven", oc)
					exutil.DumpBuildLogs("openshift-jee-sample-docker", oc)
					exutil.DumpDeploymentLogs("openshift-jee-sample", 1, oc)
				}
				debugAnyJenkinsFailure(br, oc.Namespace()+"-openshift-jee-sample", oc, true)
				br.AssertSuccess()

				// wait for the service to be running
				g.By("expecting the openshift-jee-sample service to be deployed and running")
				_, err = exutil.GetEndpointAddress(oc, "openshift-jee-sample")
				o.Expect(err).NotTo(o.HaveOccurred())

				g.By("clean up openshift resources for next potential run")
				err = oc.Run("delete").Args("bc", "openshift-jee-sample").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
				err = oc.Run("delete").Args("bc", "openshift-jee-sample-docker").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
				err = oc.Run("delete").Args("dc", "openshift-jee-sample").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
				err = oc.Run("delete").Args("is", "openshift-jee-sample").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
				err = oc.Run("delete").Args("svc", "openshift-jee-sample").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
				err = oc.Run("delete").Args("route", "openshift-jee-sample").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
				err = oc.Run("delete").Args("is", "wildfly").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
			})

			g.By("Pipelines with declarative syntax")

			g.By("should build successfully", func() {
				// create the bc
				g.By(fmt.Sprintf("calling oc create -f %q", nodejsDeclarativePipelinePath))
				err := oc.Run("create").Args("-f", nodejsDeclarativePipelinePath).Execute()
				o.Expect(err).NotTo(o.HaveOccurred())

				// start the build
				g.By("starting the pipeline build and waiting for it to complete")
				br, err := exutil.StartBuildAndWait(oc, "nodejs-sample-pipeline")
				if err != nil || !br.BuildSuccess {
					exutil.DumpBuilds(oc)
					exutil.DumpPodLogsStartingWith("nodejs", oc)
					exutil.DumpBuildLogs("nodejs-mongodb-example", oc)
					exutil.DumpDeploymentLogs("mongodb", 1, oc)
					exutil.DumpDeploymentLogs("nodejs-mongodb-example", 1, oc)
				}
				debugAnyJenkinsFailure(br, oc.Namespace()+"-nodejs-sample-pipeline", oc, true)
				br.AssertSuccess()

				// wait for the service to be running
				g.By("expecting the nodejs-mongodb-example service to be deployed and running")
				_, err = exutil.GetEndpointAddress(oc, "nodejs-mongodb-example")
				o.Expect(err).NotTo(o.HaveOccurred())

				g.By("clean up openshift resources for next potential run")
				err = oc.Run("delete").Args("all", "-l", "app=nodejs-mongodb-example").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
				err = oc.Run("delete").Args("secret", "nodejs-mongodb-example").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
				err = oc.Run("delete").Args("bc", "nodejs-sample-pipeline").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
				err = oc.Run("delete").Args("is", "nodejs-mongodb-example-staging").Execute()
				o.Expect(err).NotTo(o.HaveOccurred())
			})
		})

	})

})
