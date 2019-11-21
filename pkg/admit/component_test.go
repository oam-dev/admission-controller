package admit

import (
	"testing"

	"github.com/oam-dev/admission-controller/common"

	"github.com/oam-dev/admission-controller/pkg/apis/core.oam.dev/v1alpha1"
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
