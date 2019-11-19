package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type TraitBinding struct {
	Name string `json:"name"`

	// A properties object (for trait and scope configuration) is an object whose structure is determined by the trait or scope property schema. It may be a simple value, or it may be a complex object.
	// Properties are validated against the schema appropriate for the trait or scope.
	Properties string `json:"properties,omitempty"`
}

/// A value that is substituted into a parameter.
type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ParameterValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ComponentConfiguration struct {
	ComponentName     string           `json:"componentName"`
	InstanceName      string           `json:"instanceName"`
	ParameterValues   []ParameterValue `json:"parameterValues,omitempty"`
	Traits            []TraitBinding   `json:"traits,omitempty"`
	ApplicationScopes []string         `json:"applicationScopes,omitempty"`
}

type ScopeBinding struct {
	Name string `json:"name"`
	Type string `json:"type"`

	// A properties object (for trait and scope configuration) is an object whose structure is determined by the trait or scope property schema. It may be a simple value, or it may be a complex object.
	// Properties are validated against the schema appropriate for the trait or scope.
	Properties []map[string]intstr.IntOrString `json:"properties,omitempty"`
}

// ConfigurationSpec defines the desired state of ApplicationConfiguration
type ApplicationConfigurationSpec struct {
	Variables  []Variable               `json:"variables,omitempty"`
	Scopes     []ScopeBinding           `json:"scopes,omitempty"`
	Components []ComponentConfiguration `json:"components,omitempty"`
}

// ConfigurationStatus defines the observed state of ApplicationConfiguration
type ApplicationConfigurationStatus struct {
}

// +genclient

// ApplicationConfiguration is the Schema for the configurations API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApplicationConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationConfigurationSpec   `json:"spec,omitempty"`
	Status ApplicationConfigurationStatus `json:"status,omitempty"`
}

// ConfigurationList contains a list of ApplicationConfiguration
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApplicationConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationConfiguration `json:"items"`
}
