package common

import (
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (

	// AppConfigGroup is the Group of OAM
	AppConfigGroup = "core.oam.dev"
	// AppConfigVersion is the version of OAM spec this admission controller could validate
	AppConfigVersion = "v1alpha1"
	// AppConfigCRD is the resources of the Application Configuration
	AppConfigCRD = "applicationconfigurations"
	// ComponentCRD is the resources of the Component
	ComponentCRD = "componentschematics"
	// ScopeCRD is the resources of the Scope
	ScopeCRD = "applicationscopes"
	// TraitCRD is the resources of the Trait
	TraitCRD = "traits"
)

const (
	WorkloadServer          = "core.oam.dev/v1alpha1.Server"
	WorkloadSingletonServer = "core.oam.dev/v1alpha1.SingletonServer"
	WorkloadTask            = "core.oam.dev/v1alpha1.Task"
	WorkloadSingletonTask   = "core.oam.dev/v1alpha1.SingletonTask"
	WorkloadSingletonWorker = "core.oam.dev/v1alpha1.SingletonWorker"
	WorkloadWorker          = "core.oam.dev/v1alpha1.Worker"
)

// ToErrorResponse is a helper function to create an AdmissionResponse
// with an embedded error
func ToErrorResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}
