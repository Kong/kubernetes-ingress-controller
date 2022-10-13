package manager_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

func TestCRDControllerCondition(t *testing.T) {
	knownGvr := schema.GroupVersionResource{
		Group:    "group",
		Version:  "v1",
		Resource: "resources",
	}
	unknownGVR := schema.GroupVersionResource{
		Group:    "otherGroup",
		Version:  "v1",
		Resource: "resources",
	}

	restMapper := meta.NewDefaultRESTMapper(nil)
	restMapper.Add(schema.GroupVersionKind{
		Group:   knownGvr.Group,
		Version: knownGvr.Version,
		Kind:    "Resource",
	}, meta.RESTScopeRoot)

	testCases := []struct {
		name           string
		gvr            schema.GroupVersionResource
		expectedResult bool
	}{
		{
			name:           "registered_resource",
			gvr:            knownGvr,
			expectedResult: true,
		},
		{
			name:           "not_registered_resource",
			gvr:            unknownGVR,
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedResult,
				manager.NewCRDCondition(tc.gvr, restMapper).Enabled(),
			)
		})
	}
}
