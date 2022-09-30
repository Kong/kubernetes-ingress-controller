package manager_test

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestCRDControllerCondition(t *testing.T) {
	gvr := schema.GroupVersionResource{
		Group:    "group",
		Version:  "version",
		Resource: "resources",
	}
	meta.NewDefaultRESTMapper()
	fake.ClientBuilder{}.WithRESTMapper()

	meta.RESTMapper()
}
