package store

import (
	"reflect"

	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/client-go/tools/cache"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func keyFunc(obj interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	name := v.FieldByName("Name")
	namespace := v.FieldByName("Namespace")
	return namespace.String() + "/" + name.String(), nil
}

func clusterResourceKeyFunc(obj interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	return v.FieldByName("Name").String(), nil
}

// FakeObjects can be used to populate a fake Store.
type FakeObjects struct {
	IngressesV1beta1   []*networkingv1beta1.Ingress
	IngressesV1        []*networkingv1.Ingress
	TCPIngresses       []*configurationv1beta1.TCPIngress
	UDPIngresses       []*configurationv1beta1.UDPIngress
	Services           []*apiv1.Service
	Endpoints          []*apiv1.Endpoints
	Secrets            []*apiv1.Secret
	KongPlugins        []*configurationv1.KongPlugin
	KongClusterPlugins []*configurationv1.KongClusterPlugin
	KongIngresses      []*configurationv1.KongIngress
	KongConsumers      []*configurationv1.KongConsumer

	KnativeIngresses []*knative.Ingress
}

// NewFakeStore creates a store backed by the objects passed in as arguments.
func NewFakeStore(
	objects FakeObjects) (Storer, error) {
	var s Storer

	ingressV1beta1Store := cache.NewStore(keyFunc)
	for _, ingress := range objects.IngressesV1beta1 {
		err := ingressV1beta1Store.Add(ingress)
		if err != nil {
			return nil, err
		}
	}
	ingressV1Store := cache.NewStore(keyFunc)
	for _, ingress := range objects.IngressesV1 {
		err := ingressV1Store.Add(ingress)
		if err != nil {
			return nil, err
		}
	}
	tcpIngressStore := cache.NewStore(keyFunc)
	for _, ingress := range objects.TCPIngresses {
		err := tcpIngressStore.Add(ingress)
		if err != nil {
			return nil, err
		}
	}
	udpIngressStore := cache.NewStore(keyFunc)
	for _, ingress := range objects.UDPIngresses {
		if err := udpIngressStore.Add(ingress); err != nil {
			return nil, err
		}
	}
	serviceStore := cache.NewStore(keyFunc)
	for _, s := range objects.Services {
		err := serviceStore.Add(s)
		if err != nil {
			return nil, err
		}
	}
	secretsStore := cache.NewStore(keyFunc)
	for _, s := range objects.Secrets {
		err := secretsStore.Add(s)
		if err != nil {
			return nil, err
		}
	}
	endpointStore := cache.NewStore(keyFunc)
	for _, e := range objects.Endpoints {
		err := endpointStore.Add(e)
		if err != nil {
			return nil, err
		}
	}
	kongIngressStore := cache.NewStore(keyFunc)
	for _, k := range objects.KongIngresses {
		err := kongIngressStore.Add(k)
		if err != nil {
			return nil, err
		}
	}
	consumerStore := cache.NewStore(keyFunc)
	for _, c := range objects.KongConsumers {
		err := consumerStore.Add(c)
		if err != nil {
			return nil, err
		}
	}
	kongPluginsStore := cache.NewStore(keyFunc)
	for _, p := range objects.KongPlugins {
		err := kongPluginsStore.Add(p)
		if err != nil {
			return nil, err
		}
	}
	kongClusterPluginsStore := cache.NewStore(clusterResourceKeyFunc)
	for _, p := range objects.KongClusterPlugins {
		err := kongClusterPluginsStore.Add(p)
		if err != nil {
			return nil, err
		}
	}

	knativeIngressStore := cache.NewStore(keyFunc)
	for _, ingress := range objects.KnativeIngresses {
		err := knativeIngressStore.Add(ingress)
		if err != nil {
			return nil, err
		}
	}
	s = Store{
		stores: CacheStores{
			IngressV1beta1: ingressV1beta1Store,
			IngressV1:      ingressV1Store,
			TCPIngress:     tcpIngressStore,
			UDPIngress:     udpIngressStore,
			Service:        serviceStore,
			Endpoint:       endpointStore,
			Secret:         secretsStore,

			Plugin:        kongPluginsStore,
			ClusterPlugin: kongClusterPluginsStore,
			Consumer:      consumerStore,
			KongIngress:   kongIngressStore,

			KnativeIngress: knativeIngressStore,
		},
		ingressClass:                annotations.DefaultIngressClass,
		isValidIngressClass:         annotations.IngressClassValidatorFuncFromObjectMeta(annotations.DefaultIngressClass),
		isValidIngressV1Class:       annotations.IngressClassValidatorFuncFromV1Ingress(annotations.DefaultIngressClass),
		ingressV1Beta1ClassMatching: annotations.ExactClassMatch,
		ingressV1ClassMatching:      annotations.ExactClassMatch,
		kongConsumerClassMatching:   annotations.ExactClassMatch,
	}
	return s, nil
}
