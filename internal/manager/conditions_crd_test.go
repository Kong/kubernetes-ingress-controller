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
		enabled        bool
		expectedResult bool
	}{
		{
			name:           "enabled_and_registered_resource",
			gvr:            knownGvr,
			enabled:        true,
			expectedResult: true,
		},
		{
			name:           "disabled_and_registered_resource",
			gvr:            knownGvr,
			enabled:        false,
			expectedResult: false,
		},
		{
			name:           "enabled_and_not_registered_resource",
			gvr:            unknownGVR,
			enabled:        true,
			expectedResult: false,
		},
		{
			name:           "disabled_and_not_registered_resource",
			gvr:            unknownGVR,
			enabled:        false,
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedResult,
				manager.NewCRDCondition(tc.gvr, tc.enabled, restMapper).Enabled(),
			)
		})
	}
}
