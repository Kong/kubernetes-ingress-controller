package controllers

import (
	"context"

	"github.com/kong/go-kong/kong"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8sobj "github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object"
)

// DataPlane is a common interface that is used by reconcilers to interact
// with the Kong dataplane.
type DataPlane interface {
	DataPlaneClient
	DataPlaneStatusClient

	Listeners(ctx context.Context) ([]kong.ProxyListener, []kong.StreamListener, error)
}

type DataPlaneStatusClient interface {
	AreKubernetesObjectReportsEnabled() bool
	KubernetesObjectConfigurationStatus(obj client.Object) k8sobj.ConfigurationStatus
	KubernetesObjectIsConfigured(obj client.Object) bool
}

// DataPlaneClient is a common client interface that is used by reconcilers to interact
// with the Kong dataplane to perform CRUD operations on provided objects.
type DataPlaneClient interface {
	UpdateObject(obj client.Object) error
	DeleteObject(obj client.Object) error
	ObjectExists(obj client.Object) (bool, error)
}
