package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type TraitSpec struct {
	AppliesTo  []string `json:"appliesTo"`
	Properties string   `json:"properties"`
}

type TraitStatus struct {
}

// +genclient

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Trait struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TraitSpec   `json:"spec,omitempty"`
	Status TraitStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TraitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Trait `json:"items"`
}
