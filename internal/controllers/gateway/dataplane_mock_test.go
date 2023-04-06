package gateway

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	k8sobj "github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object"
)

type DataplaneMock struct {
	KubernetesObjectReportsEnabled bool
	// Mapping namespace to name to status
	// Note: this will come in useful when implementing
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3793
	// which requires the status to be reported for route objects.
	ObjectsStatuses map[string]map[string]k8sobj.ConfigurationStatus
}

func (d DataplaneMock) UpdateObject(_ context.Context, _ client.Object) error {
	return nil
}

func (d DataplaneMock) DeleteObject(_ context.Context, _ client.Object) error {
	return nil
}

func (d DataplaneMock) AreKubernetesObjectReportsEnabled() bool {
	return d.KubernetesObjectReportsEnabled
}

func (d DataplaneMock) KubernetesObjectConfigurationStatus(obj client.Object) k8sobj.ConfigurationStatus {
	return d.ObjectsStatuses[obj.GetNamespace()][obj.GetName()]
}
