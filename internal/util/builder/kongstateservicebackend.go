package builder

import (
	"fmt"

	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
)

// KongstateServiceBackendBuilder is a builder for KongstateServiceBackend.
// Primarily used for testing.
type KongstateServiceBackendBuilder struct {
	name      string
	namespace string
	weight    *int32
	portDef   kongstate.PortDef
	t         kongstate.ServiceBackendType
}

func NewKongstateServiceBackend(name string) *KongstateServiceBackendBuilder {
	return &KongstateServiceBackendBuilder{
		name: name,
	}
}

func (b *KongstateServiceBackendBuilder) WithNamespace(namespace string) *KongstateServiceBackendBuilder {
	b.namespace = namespace
	return b
}

func (b *KongstateServiceBackendBuilder) WithWeight(weight int) *KongstateServiceBackendBuilder {
	b.weight = lo.ToPtr(int32(weight))
	return b
}

func (b *KongstateServiceBackendBuilder) WithPortNumber(port int) *KongstateServiceBackendBuilder {
	b.portDef = kongstate.PortDef{
		Number: int32(port),
		Mode:   kongstate.PortModeByNumber,
	}
	return b
}

func (b *KongstateServiceBackendBuilder) WithPortName(port string) *KongstateServiceBackendBuilder {
	b.portDef = kongstate.PortDef{
		Name: port,
		Mode: kongstate.PortModeByName,
	}
	return b
}

func (b *KongstateServiceBackendBuilder) WithType(t kongstate.ServiceBackendType) *KongstateServiceBackendBuilder {
	b.t = t
	return b
}

func (b *KongstateServiceBackendBuilder) MustBuild() kongstate.ServiceBackend {
	// Default to Kubernetes Service backend type if not specified.
	if b.t == "" {
		b.t = kongstate.ServiceBackendTypeKubernetesService
	}
	// Default to default namespace if not specified.
	if b.namespace == "" {
		b.namespace = metav1.NamespaceDefault
	}

	s, err := kongstate.NewServiceBackend(
		b.t,
		k8stypes.NamespacedName{
			Namespace: b.namespace,
			Name:      b.name,
		},
		b.portDef,
	)
	if err != nil {
		// This should never happen. If it does, it's a bug that will be discovered in tests as the builder
		// is used in tests only.
		panic(fmt.Errorf("failed to build service backend: %w", err))
	}
	if b.weight != nil {
		s.SetWeight(*b.weight)
	}
	return s
}
