package controllers

import (
	"context"

	"github.com/kong/go-kong/kong"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8sobj "github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object"
)

// DataPlane is a common interface that is used by reconcilers to interact
// with the dataplane.
//
// TODO: This can probably be used in other reconcilers as well.
// Related issue: https://github.com/Kong/kubernetes-ingress-controller/issues/3794
type DataPlane interface {
	DataPlaneClient

	Listeners(ctx context.Context) ([]kong.ProxyListener, []kong.StreamListener, error)
	AreKubernetesObjectReportsEnabled() bool
	KubernetesObjectConfigurationStatus(obj client.Object) k8sobj.ConfigurationStatus
}

// DataPlaneClient is a common client interface that is used by reconcilers to interact
// with the dataplane to perform CRUD operations on provided objects.
type DataPlaneClient interface {
	UpdateObject(obj client.Object) error
	DeleteObject(obj client.Object) error
}
