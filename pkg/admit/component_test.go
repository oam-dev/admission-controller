package admit

import (
	"testing"

	"github.com/oam-dev/admission-controller/common"

	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestValidateWorker(t *testing.T) {
	err := validateWorker(&v1alpha1.ComponentSchematic{
		Spec: v1alpha1.ComponentSpec{
			WorkloadType: common.WorkloadServer,
			Containers: []v1alpha1.Container{
				{
					Ports: []v1alpha1.Port{{Name: "p1", ContainerPort: 9999}},
				},
			},
		},
	})
	assert.NoError(t, err)
	err = validateWorker(&v1alpha1.ComponentSchematic{
		Spec: v1alpha1.ComponentSpec{
			WorkloadType: common.WorkloadWorker,
			Containers: []v1alpha1.Container{
				{
					Ports: []v1alpha1.Port{{Name: "p1", ContainerPort: 9999}},
				},
			},
		},
	})
	assert.Error(t, err)
}

func TestValidateGVK(t *testing.T) {
	tests := []struct {
		WorkloadType string
		Group        string
		Version      string
		Kind         string
		Err          bool
	}{
		{
			WorkloadType: "core.oam.dev/v1alpha1.Singleton",
			Group:        "core.oam.dev",
			Version:      "v1alpha1",
			Kind:         "Singleton",
			Err:          false,
		},
		{
			WorkloadType: "core.oam.dev/v1alpha1",
			Group:        "core.oam.dev",
			Version:      "v1alpha1",
			Kind:         "Singleton",
			Err:          true,
		},
		{
			WorkloadType: "caching.oam.dev/v2.Redis",
			Group:        "caching.oam.dev",
			Version:      "v2",
			Kind:         "Redis",
			Err:          false,
		},
	}
	for _, ti := range tests {
		g, v, k, err := validateGVK(ti.WorkloadType)
		if !ti.Err {
			assert.NoError(t, err)
			assert.Equal(t, ti.Group, g)
			assert.Equal(t, ti.Version, v)
			assert.Equal(t, ti.Kind, k)
		} else {
			assert.Error(t, err)
		}
	}
}
