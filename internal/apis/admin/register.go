package admin

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	adminv1 "github.com/kong/ingress-controller/internal/apis/admin/v1"
)

const GroupName = "konghq.com"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&adminv1.Route{},
		&adminv1.Service{},
	)
	return nil
}
