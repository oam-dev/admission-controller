package admit

import (
	"encoding/json"
	"fmt"

	"github.com/oam-dev/admission-controller/common"

	"github.com/oam-dev/admission-controller/pkg/apis/core.oam.dev/v1alpha1"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

var AppConfigResource = metav1.GroupVersionResource{Group: common.AppConfigGroup, Version: common.AppConfigVersion, Resource: common.AppConfigCRD}

// validate Application Configuration Spec here
func (a *Admit) AppConfigSpec(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.V(2).Info("admitting Application Configuration")

	if ar.Request.Resource != AppConfigResource {
		klog.Errorf("expect resource to be %s", AppConfigResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	appConfig := v1alpha1.ApplicationConfiguration{}
	deserializer := common.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &appConfig); err != nil {
		klog.Error(err)
		return common.ToErrorResponse(err)
	}

	//check application scope exist when they created by AppConfig
	if err := a.checkAppScope(&appConfig); err != nil {
		return common.ToErrorResponse(err)
	}
	// check component existence
	if err := a.checkComponent(&appConfig); err != nil {
		return common.ToErrorResponse(err)
	}
	// check trait existence
	if err := a.checkTrait(&appConfig); err != nil {
		return common.ToErrorResponse(err)
	}
	//check appscope instance existence
	if err := a.checkAppScopeInstance(&appConfig); err != nil {
		return common.ToErrorResponse(err)
	}
	// TODO check component/trait variables and properties in AppConfig

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}

// mutate Application Configuration Spec here
func (a *Admit) MutateAppConfigSpec(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.V(2).Info("mutating Application Configuration")
	if ar.Request.Resource != AppConfigResource {
		klog.Errorf("expect resource to be %s", AppConfigResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	appConfig := v1alpha1.ApplicationConfiguration{}
	deserializer := common.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &appConfig); err != nil {
		klog.Error(err)
		return common.ToErrorResponse(err)
	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	var patches []patchOperation
	patches = append(patches, mutateFromVariable(&appConfig)...)

	if len(patches) > 0 {
		patchBytes, err := json.Marshal(patches)
		if err != nil {
			klog.Error(err)
			return common.ToErrorResponse(err)
		}
		reviewResponse.Patch = patchBytes
	}
	pt := v1beta1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt

	return &reviewResponse
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func mutateFromVariable(appConfig *v1alpha1.ApplicationConfiguration) []patchOperation {
	// TODO generate patchOperation here
	return nil
}

func (a *Admit) checkComponent(appConf *v1alpha1.ApplicationConfiguration) error {
	for _, v := range appConf.Spec.Components {
		_, err := a.componentInformer.Lister().ComponentSchematics(appConf.Namespace).Get(v.ComponentName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Admit) checkAppScopeInstance(appConf *v1alpha1.ApplicationConfiguration) error {
	for _, v := range appConf.Spec.Components {
		for _, s := range v.ApplicationScopes {
			scopeInstance, err := a.appConfigInformer.Lister().ApplicationConfigurations(appConf.Namespace).Get(s)
			if err != nil {
				return err
			}
			if len(scopeInstance.Spec.Scopes) < 1 {
				return fmt.Errorf("%s doesn't have any scope binded", s)
			}

		}
	}
	return nil
}

func (a *Admit) checkTrait(appConf *v1alpha1.ApplicationConfiguration) error {
	for _, v := range appConf.Spec.Components {
		for _, t := range v.Traits {
			_, err := a.traitInformer.Lister().Traits(appConf.Namespace).Get(t.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Admit) checkAppScope(appConf *v1alpha1.ApplicationConfiguration) error {
	for _, v := range appConf.Spec.Scopes {
		scope, err := a.scopeInformer.Lister().Scopes(appConf.Namespace).Get(v.Name)
		if err != nil {
			return err
		}
		if scope.Spec.Type != v.Type {
			return fmt.Errorf("don't have type %s for scope %s but %s", v.Name, v.Type, scope.Spec.Type)
		}

	}
	return nil
}
