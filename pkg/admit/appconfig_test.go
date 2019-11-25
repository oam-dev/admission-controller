package admit

import (
	"testing"

	"github.com/oam-dev/admission-controller/pkg/apis/core.oam.dev/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestParseFromVariable(t *testing.T) {
	var cases = []struct {
		Variables []v1alpha1.Variable
		Value     string
		expValue  string
		expChange bool
	}{
		{
			Variables: []v1alpha1.Variable{
				{
					Name:  "VAR_42",
					Value: "v1",
				},
			},
			Value:     "VAR_42",
			expValue:  "",
			expChange: false,
		},
		{
			Variables: []v1alpha1.Variable{
				{
					Name:  "VAR_42",
					Value: "v1",
				},
			},
			Value:     "[fromVariable (VAR_42)]",
			expValue:  "",
			expChange: false,
		},
		{
			Variables: []v1alpha1.Variable{
				{
					Name:  "VAR_42",
					Value: "v1",
				},
				{
					Name:  "VAR_",
					Value: "v2",
				},
			},
			Value:     "[fromVariable(VAR_42)]",
			expValue:  "v1",
			expChange: true,
		},
		{
			Variables: []v1alpha1.Variable{
				{
					Name:  "VAR_42",
					Value: "v1",
				},
				{
					Name:  "VAR_",
					Value: "v2",
				},
			},
			Value:     "[fromVariable(VAR_)]",
			expValue:  "v2",
			expChange: true,
		},
		{
			Variables: []v1alpha1.Variable{
				{
					Name:  "_",
					Value: "v_",
				},
			},
			Value:     "[fromVariable(_)]",
			expValue:  "v_",
			expChange: true,
		},
		{
			Variables: []v1alpha1.Variable{
				{
					Name:  "var",
					Value: "v",
				},
			},
			Value:     "[fromVariable()]",
			expValue:  "",
			expChange: false,
		},
		{
			Variables: []v1alpha1.Variable{
				{
					Name:  "42",
					Value: "v42",
				},
			},
			Value:     "[fromVariable(42)]",
			expValue:  "v42",
			expChange: true,
		},
		{
			Variables: []v1alpha1.Variable{},
			Value:     "[fromVariable(42)]",
			expValue:  "",
			expChange: true,
		},
		{
			Variables: []v1alpha1.Variable{
				{
					Name:  "VAR",
					Value: "vv",
				},
			},
			Value:     "[fromVariable(VAR)]",
			expValue:  "vv",
			expChange: true,
		},
	}
	for idx, tc := range cases {
		gotVar, gotChange := parseFromVariable(tc.Value, tc.Variables)
		assert.Equal(t, tc.expChange, gotChange, "case index %d", idx)
		assert.Equal(t, tc.expValue, gotVar, "case index %d", idx)
	}
}

func TestMutateFromVariable(t *testing.T) {
	patches := mutateFromVariable(&v1alpha1.ApplicationConfiguration{
		Spec: v1alpha1.ApplicationConfigurationSpec{
			Variables: []v1alpha1.Variable{
				{
					Name:  "host",
					Value: "www.example.com",
				},
				{
					Name:  "port",
					Value: "9091",
				},
			},
			Components: []v1alpha1.ComponentConfiguration{
				{
					ParameterValues: []v1alpha1.ParameterValue{
						{
							Name:  "hostname",
							Value: "[fromVariable(host)]",
						},
					},
					Traits: []v1alpha1.TraitBinding{
						{
							Name:       "service_port",
							Properties: "[fromVariable(port)]",
						},
					},
				},
			},
		},
	})
	assert.Equal(t, []patchOperation{
		{Op: "replace", Path: "/spec/components/0/parameterValues/0/value", Value: "www.example.com"},
		{Op: "replace", Path: "/spec/components/0/traits/0/properties", Value: "9091"},
	}, patches)

}
