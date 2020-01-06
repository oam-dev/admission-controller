package admit

import (
	"github.com/oam-dev/admission-controller/common"
	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// validate Trait Spec here
func (a *Admit) TraitSpec(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.Info("admitting Trait")
	traitResource := metav1.GroupVersionResource{Group: common.AppConfigGroup, Version: common.AppConfigVersion, Resource: common.TraitCRD}

	if ar.Request.Resource != traitResource {
		klog.Errorf("expect resource to be %s", traitResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	comp := v1alpha1.Trait{}
	deserializer := common.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &comp); err != nil {
		klog.Error(err)
		return common.ToErrorResponse(err)
	}

	// TODO validate trait spec

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}
