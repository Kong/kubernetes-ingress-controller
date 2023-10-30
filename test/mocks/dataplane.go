package mocks

import (
	"context"

	"github.com/kong/go-kong/kong"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8sobj "github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object"
)

type Dataplane struct {
	KubernetesObjectReportsEnabled bool
	// Mapping namespace to name to status
	// Note: this will come in useful when implementing
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3793
	// which requires the status to be reported for route objects.
	ObjectsStatuses map[string]map[string]k8sobj.ConfigurationStatus
}

func (d Dataplane) UpdateObject(_ client.Object) error {
	return nil
}

func (d Dataplane) DeleteObject(_ client.Object) error {
	return nil
}

func (d Dataplane) ObjectExists(obj client.Object) (bool, error) {
	_, exists := d.ObjectsStatuses[obj.GetNamespace()][obj.GetName()]
	return exists, nil
}

func (d Dataplane) Listeners(context.Context) ([]kong.ProxyListener, []kong.StreamListener, error) {
	return nil, nil, nil
}

func (d Dataplane) AreKubernetesObjectReportsEnabled() bool {
	return d.KubernetesObjectReportsEnabled
}

func (d Dataplane) KubernetesObjectConfigurationStatus(obj client.Object) k8sobj.ConfigurationStatus {
	return d.ObjectsStatuses[obj.GetNamespace()][obj.GetName()]
}

func (d Dataplane) KubernetesObjectIsConfigured(obj client.Object) bool {
	return d.ObjectsStatuses[obj.GetNamespace()][obj.GetName()] == k8sobj.ConfigurationStatusSucceeded
}
