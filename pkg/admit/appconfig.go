package admit

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/oam-dev/admission-controller/common"

	"github.com/oam-dev/admission-controller/pkg/apis/core.oam.dev/v1alpha1"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

var AppConfigResource = metav1.GroupVersionResource{Group: common.AppConfigGroup, Version: common.AppConfigVersion, Resource: common.AppConfigCRD}

// validate Application Configuration Spec here
func (a *Admit) AppConfigSpec(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.Info("admitting Application Configuration")

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

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}

// mutate Application Configuration Spec here
func (a *Admit) MutateAppConfigSpec(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.Info("mutating Application Configuration")
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

	patches := mutateFromVariable(&appConfig)
	if len(patches) > 0 {
		patchBytes, err := json.Marshal(patches)
		if err != nil {
			klog.Error(err)
			return common.ToErrorResponse(err)
		}
		reviewResponse.Patch = patchBytes
		reviewResponse.Result = &metav1.Status{
			Status: "Success",
		}
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

// value of properties and parameters will contain [fromVariable(VARNAME)] function call
func mutateFromVariable(appConfig *v1alpha1.ApplicationConfiguration) []patchOperation {
	var patches []patchOperation
	for cid, comp := range appConfig.Spec.Components {
		for idx, v := range comp.ParameterValues {
			parsedValue, changed := parseFromVariable(v.Value, appConfig.Spec.Variables)
			if changed {
				patches = append(patches, patchOperation{
					Op:    "replace",
					Path:  fmt.Sprintf("/spec/components/%d/parameterValues/%d/value", cid, idx),
					Value: parsedValue,
				})
			}
		}
		for idx, v := range comp.Traits {
			parsedValue, changed := parseFromVariable(v.Properties, appConfig.Spec.Variables)
			if changed {
				patches = append(patches, patchOperation{
					Op:    "replace",
					Path:  fmt.Sprintf("/spec/components/%d/traits/%d/properties", cid, idx),
					Value: parsedValue,
				})
			}
		}
	}
	//TODO make sure the type of properties in scope binding, if they will contain fromVariable function call, mutate it
	return patches
}

var variableExp = regexp.MustCompile(`^\[fromVariable\((?P<var>[[:word:]]+)\)\]$`)

func parseFromVariable(value string, variables []v1alpha1.Variable) (string, bool) {
	res := variableExp.FindStringSubmatch(value)
	if len(res) < 2 {
		return "", false
	}
	for _, v := range variables {
		if v.Name == res[1] {
			return v.Value, true
		}
	}
	return "", true
}

func (a *Admit) checkComponent(appConf *v1alpha1.ApplicationConfiguration) error {
	for _, v := range appConf.Spec.Components {
		comp, err := a.componentInformer.Lister().ComponentSchematics(appConf.Namespace).Get(v.ComponentName)
		if err != nil {
			return err
		}
		for _, p := range v.ParameterValues {
			var found = false
			for _, cp := range comp.Spec.Parameters {
				if cp.Name == p.Name {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("component %s don't have this parameter %s", v.ComponentName, p.Name)
			}
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
			//TODO check trait properties in AppConfig really exist in that Trait spec
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
