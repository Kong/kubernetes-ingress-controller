package helpers

import (
	"fmt"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// AddIngressToCleaner adds a runtime.Object to the cleanup list if it is a supported version of Ingress. It panics if the
// runtime.Object is something else.
func AddIngressToCleaner(cleaner *clusters.Cleaner, obj runtime.Object) {
	switch i := obj.(type) {
	case *netv1.Ingress:
		cleaner.Add(i)
	case *netv1beta1.Ingress:
		cleaner.Add(i)
	default:
		panic(fmt.Sprintf("%s passed to addIngressToCleaner but is not an Ingress", obj.GetObjectKind()))
	}
}
