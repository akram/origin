// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// CloudPrivateIPConfigSpecApplyConfiguration represents a declarative configuration of the CloudPrivateIPConfigSpec type for use
// with apply.
type CloudPrivateIPConfigSpecApplyConfiguration struct {
	Node *string `json:"node,omitempty"`
}

// CloudPrivateIPConfigSpecApplyConfiguration constructs a declarative configuration of the CloudPrivateIPConfigSpec type for use with
// apply.
func CloudPrivateIPConfigSpec() *CloudPrivateIPConfigSpecApplyConfiguration {
	return &CloudPrivateIPConfigSpecApplyConfiguration{}
}

// WithNode sets the Node field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Node field is set to the value of the last call.
func (b *CloudPrivateIPConfigSpecApplyConfiguration) WithNode(value string) *CloudPrivateIPConfigSpecApplyConfiguration {
	b.Node = &value
	return b
}
