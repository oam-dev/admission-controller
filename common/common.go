package common

import "k8s.io/api/admission/v1beta1"
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	AppConfigGroup   = "core.oam.dev"
	AppConfigVersion = "v1alpha1"
	AppConfigCRD     = "applicationconfigurations"

	ComponentCRD = "componentschematics"
	ScopeCRD     = "scopes"
	TraitCRD     = "traits"
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
