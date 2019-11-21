package admit

import (
	"github.com/oam-dev/admission-controller/common"
	"github.com/oam-dev/admission-controller/pkg/apis/core.oam.dev/v1alpha1"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// validate Scope Spec here
func (a *Admit) ScopeSpec(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.Info("admitting Application Scope")
	scopeResource := metav1.GroupVersionResource{Group: common.AppConfigGroup, Version: common.AppConfigVersion, Resource: common.ScopeCRD}

	if ar.Request.Resource != scopeResource {
		klog.Errorf("expect resource to be %s", scopeResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	comp := v1alpha1.Scope{}
	deserializer := common.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &comp); err != nil {
		klog.Error(err)
		return common.ToErrorResponse(err)
	}

	// TODO validate scope spec

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}
