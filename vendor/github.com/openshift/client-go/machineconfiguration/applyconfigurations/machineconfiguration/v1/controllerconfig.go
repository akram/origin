// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	machineconfigurationv1 "github.com/openshift/api/machineconfiguration/v1"
	internal "github.com/openshift/client-go/machineconfiguration/applyconfigurations/internal"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	managedfields "k8s.io/apimachinery/pkg/util/managedfields"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

// ControllerConfigApplyConfiguration represents a declarative configuration of the ControllerConfig type for use
// with apply.
type ControllerConfigApplyConfiguration struct {
	metav1.TypeMetaApplyConfiguration    `json:",inline"`
	*metav1.ObjectMetaApplyConfiguration `json:"metadata,omitempty"`
	Spec                                 *ControllerConfigSpecApplyConfiguration   `json:"spec,omitempty"`
	Status                               *ControllerConfigStatusApplyConfiguration `json:"status,omitempty"`
}

// ControllerConfig constructs a declarative configuration of the ControllerConfig type for use with
// apply.
func ControllerConfig(name string) *ControllerConfigApplyConfiguration {
	b := &ControllerConfigApplyConfiguration{}
	b.WithName(name)
	b.WithKind("ControllerConfig")
	b.WithAPIVersion("machineconfiguration.openshift.io/v1")
	return b
}

// ExtractControllerConfig extracts the applied configuration owned by fieldManager from
// controllerConfig. If no managedFields are found in controllerConfig for fieldManager, a
// ControllerConfigApplyConfiguration is returned with only the Name, Namespace (if applicable),
// APIVersion and Kind populated. It is possible that no managed fields were found for because other
// field managers have taken ownership of all the fields previously owned by fieldManager, or because
// the fieldManager never owned fields any fields.
// controllerConfig must be a unmodified ControllerConfig API object that was retrieved from the Kubernetes API.
// ExtractControllerConfig provides a way to perform a extract/modify-in-place/apply workflow.
// Note that an extracted apply configuration will contain fewer fields than what the fieldManager previously
// applied if another fieldManager has updated or force applied any of the previously applied fields.
// Experimental!
func ExtractControllerConfig(controllerConfig *machineconfigurationv1.ControllerConfig, fieldManager string) (*ControllerConfigApplyConfiguration, error) {
	return extractControllerConfig(controllerConfig, fieldManager, "")
}

// ExtractControllerConfigStatus is the same as ExtractControllerConfig except
// that it extracts the status subresource applied configuration.
// Experimental!
func ExtractControllerConfigStatus(controllerConfig *machineconfigurationv1.ControllerConfig, fieldManager string) (*ControllerConfigApplyConfiguration, error) {
	return extractControllerConfig(controllerConfig, fieldManager, "status")
}

func extractControllerConfig(controllerConfig *machineconfigurationv1.ControllerConfig, fieldManager string, subresource string) (*ControllerConfigApplyConfiguration, error) {
	b := &ControllerConfigApplyConfiguration{}
	err := managedfields.ExtractInto(controllerConfig, internal.Parser().Type("com.github.openshift.api.machineconfiguration.v1.ControllerConfig"), fieldManager, b, subresource)
	if err != nil {
		return nil, err
	}
	b.WithName(controllerConfig.Name)

	b.WithKind("ControllerConfig")
	b.WithAPIVersion("machineconfiguration.openshift.io/v1")
	return b, nil
}

// WithKind sets the Kind field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Kind field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithKind(value string) *ControllerConfigApplyConfiguration {
	b.TypeMetaApplyConfiguration.Kind = &value
	return b
}

// WithAPIVersion sets the APIVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the APIVersion field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithAPIVersion(value string) *ControllerConfigApplyConfiguration {
	b.TypeMetaApplyConfiguration.APIVersion = &value
	return b
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithName(value string) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ObjectMetaApplyConfiguration.Name = &value
	return b
}

// WithGenerateName sets the GenerateName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GenerateName field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithGenerateName(value string) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ObjectMetaApplyConfiguration.GenerateName = &value
	return b
}

// WithNamespace sets the Namespace field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespace field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithNamespace(value string) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ObjectMetaApplyConfiguration.Namespace = &value
	return b
}

// WithUID sets the UID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UID field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithUID(value types.UID) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ObjectMetaApplyConfiguration.UID = &value
	return b
}

// WithResourceVersion sets the ResourceVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ResourceVersion field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithResourceVersion(value string) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ObjectMetaApplyConfiguration.ResourceVersion = &value
	return b
}

// WithGeneration sets the Generation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Generation field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithGeneration(value int64) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ObjectMetaApplyConfiguration.Generation = &value
	return b
}

// WithCreationTimestamp sets the CreationTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CreationTimestamp field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithCreationTimestamp(value apismetav1.Time) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ObjectMetaApplyConfiguration.CreationTimestamp = &value
	return b
}

// WithDeletionTimestamp sets the DeletionTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionTimestamp field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithDeletionTimestamp(value apismetav1.Time) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ObjectMetaApplyConfiguration.DeletionTimestamp = &value
	return b
}

// WithDeletionGracePeriodSeconds sets the DeletionGracePeriodSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionGracePeriodSeconds field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithDeletionGracePeriodSeconds(value int64) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ObjectMetaApplyConfiguration.DeletionGracePeriodSeconds = &value
	return b
}

// WithLabels puts the entries into the Labels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Labels field,
// overwriting an existing map entries in Labels field with the same key.
func (b *ControllerConfigApplyConfiguration) WithLabels(entries map[string]string) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.ObjectMetaApplyConfiguration.Labels == nil && len(entries) > 0 {
		b.ObjectMetaApplyConfiguration.Labels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.ObjectMetaApplyConfiguration.Labels[k] = v
	}
	return b
}

// WithAnnotations puts the entries into the Annotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Annotations field,
// overwriting an existing map entries in Annotations field with the same key.
func (b *ControllerConfigApplyConfiguration) WithAnnotations(entries map[string]string) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.ObjectMetaApplyConfiguration.Annotations == nil && len(entries) > 0 {
		b.ObjectMetaApplyConfiguration.Annotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.ObjectMetaApplyConfiguration.Annotations[k] = v
	}
	return b
}

// WithOwnerReferences adds the given value to the OwnerReferences field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the OwnerReferences field.
func (b *ControllerConfigApplyConfiguration) WithOwnerReferences(values ...*metav1.OwnerReferenceApplyConfiguration) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithOwnerReferences")
		}
		b.ObjectMetaApplyConfiguration.OwnerReferences = append(b.ObjectMetaApplyConfiguration.OwnerReferences, *values[i])
	}
	return b
}

// WithFinalizers adds the given value to the Finalizers field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Finalizers field.
func (b *ControllerConfigApplyConfiguration) WithFinalizers(values ...string) *ControllerConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		b.ObjectMetaApplyConfiguration.Finalizers = append(b.ObjectMetaApplyConfiguration.Finalizers, values[i])
	}
	return b
}

func (b *ControllerConfigApplyConfiguration) ensureObjectMetaApplyConfigurationExists() {
	if b.ObjectMetaApplyConfiguration == nil {
		b.ObjectMetaApplyConfiguration = &metav1.ObjectMetaApplyConfiguration{}
	}
}

// WithSpec sets the Spec field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Spec field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithSpec(value *ControllerConfigSpecApplyConfiguration) *ControllerConfigApplyConfiguration {
	b.Spec = value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *ControllerConfigApplyConfiguration) WithStatus(value *ControllerConfigStatusApplyConfiguration) *ControllerConfigApplyConfiguration {
	b.Status = value
	return b
}

// GetName retrieves the value of the Name field in the declarative configuration.
func (b *ControllerConfigApplyConfiguration) GetName() *string {
	b.ensureObjectMetaApplyConfigurationExists()
	return b.ObjectMetaApplyConfiguration.Name
}
