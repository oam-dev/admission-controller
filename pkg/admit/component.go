package admit

import (
	"github.com/oam-dev/admission-controller/common"
	"github.com/oam-dev/admission-controller/pkg/apis/core.oam.dev/v1alpha1"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// validate Component Spec here
func (a *Admit) ComponentSpec(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.Info("admitting Component Schematics")
	componentResource := metav1.GroupVersionResource{Group: common.AppConfigGroup, Version: common.AppConfigVersion, Resource: common.ComponentCRD}

	if ar.Request.Resource != componentResource {
		klog.Errorf("expect resource to be %s", componentResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	comp := v1alpha1.ComponentSchematic{}
	deserializer := common.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &comp); err != nil {
		klog.Error(err)
		return common.ToErrorResponse(err)
	}

	// TODO invalidate if a worker tries to bind to a port

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}
