package admit

import (
	"fmt"
	"strings"

	"github.com/oam-dev/admission-controller/common"
	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
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
	if err := validateWorker(&comp); err != nil {
		return common.ToErrorResponse(err)
	}

	if _, _, _, err := validateGVK(comp.Spec.WorkloadType); err != nil {
		return common.ToErrorResponse(err)
	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}

func validateGVK(gvk string) (g, v, k string, err error) {
	parts := strings.Split(gvk, "/")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("missing version and kind in workload type")
	}

	vk := strings.Split(parts[1], ".")
	if len(vk) != 2 {
		return "", "", "", fmt.Errorf("missing kind in workload type")
	}
	return parts[0], vk[0], vk[1], nil
}

func validateWorker(comp *v1alpha1.ComponentSchematic) error {
	if comp.Spec.WorkloadType != common.WorkloadSingletonWorker && comp.Spec.WorkloadType != common.WorkloadWorker {
		return nil
	}
	for _, v := range comp.Spec.Containers {
		if len(v.Ports) > 0 {
			return fmt.Errorf("worker container named %s has a port declared", v.Name)
		}
	}
	return nil
}
