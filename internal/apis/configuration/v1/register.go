package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// TODO: move SchemeBuilder with zz_generated.deepcopy.go to k8s.io/api.
	// localSchemeBuilder and AddToScheme will stay in k8s.io/kubernetes.

	// SchemeBuilder is a schemeBuilder with all CRDs
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme points to SchemeBuilder's AddToScheme
	AddToScheme = SchemeBuilder.AddToScheme
	// SchemeGroupVersion is API Group and Version of
	// Kong Ingress Controller's API.
	SchemeGroupVersion = schema.GroupVersion{Group: "configuration.konghq.com", Version: "v1"}
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&KongIngress{},
		&KongIngressList{},
		&KongPlugin{},
		&KongPluginList{},
		&KongConsumer{},
		&KongConsumerList{},
		&KongCredential{},
		&KongCredentialList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
