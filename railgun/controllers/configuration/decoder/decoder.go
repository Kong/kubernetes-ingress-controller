package decoder

import (
	"fmt"
	"reflect"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StoreBuilder struct {
	objs store.FakeObjects
}

func (sb *StoreBuilder) Add(obj client.Object) error {
	switch obj := obj.(type) {
	case *networkingv1beta1.Ingress:
		sb.objs.IngressesV1beta1 = append(sb.objs.IngressesV1beta1, obj)
	case *networkingv1.Ingress:
		sb.objs.IngressesV1 = append(sb.objs.IngressesV1, obj)
	case *configurationv1beta1.TCPIngress:
		sb.objs.TCPIngresses = append(sb.objs.TCPIngresses, obj)
	case *v1alpha1.UDPIngress:
		sb.objs.UDPIngresses = append(sb.objs.UDPIngresses, obj)
	case *corev1.Service:
		sb.objs.Services = append(sb.objs.Services, obj)
	case *corev1.Endpoints:
		sb.objs.Endpoints = append(sb.objs.Endpoints, obj)
	case *corev1.Secret:
		sb.objs.Secrets = append(sb.objs.Secrets, obj)
	case *configurationv1.KongPlugin:
		sb.objs.KongPlugins = append(sb.objs.KongPlugins, obj)
	case *configurationv1.KongClusterPlugin:
		sb.objs.KongClusterPlugins = append(sb.objs.KongClusterPlugins, obj)
	case *configurationv1.KongIngress:
		sb.objs.KongIngresses = append(sb.objs.KongIngresses, obj)
	case *configurationv1.KongConsumer:
		sb.objs.KongConsumers = append(sb.objs.KongConsumers, obj)
	case *knative.Ingress:
		sb.objs.KnativeIngresses = append(sb.objs.KnativeIngresses, obj)
	default:
		return fmt.Errorf("unsupported type %q", reflect.TypeOf(obj))
	}
	return nil
}

func (sb *StoreBuilder) Build() (store.Storer, error) {
	return store.NewFakeStore(sb.objs)
}
