package store

import (
	"reflect"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	apiv1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/client-go/tools/cache"
)

func keyFunc(obj interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	name := v.FieldByName("Name")
	namespace := v.FieldByName("Namespace")
	return namespace.String() + "/" + name.String(), nil
}

// FakeObjects can be used to populate a fake Store.
type FakeObjects struct {
	Ingresses       []*networking.Ingress
	Services        []*apiv1.Service
	Endpoints       []*apiv1.Endpoints
	Secrets         []*apiv1.Secret
	KongPlugins     []*configurationv1.KongPlugin
	KongIngresses   []*configurationv1.KongIngress
	KongConsumers   []*configurationv1.KongConsumer
	KongCredentials []*configurationv1.KongCredential
}

// NewFakeStore creates a store backed by the objects passed in as arguments.
func NewFakeStore(
	objects FakeObjects) (Storer, error) {
	var s Storer

	ingressStore := cache.NewStore(keyFunc)
	for _, ingress := range objects.Ingresses {
		err := ingressStore.Add(ingress)
		if err != nil {
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
	kongCredentialsStore := cache.NewStore(keyFunc)
	for _, c := range objects.KongCredentials {
		err := kongCredentialsStore.Add(c)
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
	s = Store{
		stores: CacheStores{
			Ingress:  ingressStore,
			Service:  serviceStore,
			Endpoint: endpointStore,
			Secret:   secretsStore,

			Plugin:        kongPluginsStore,
			Consumer:      consumerStore,
			Credential:    kongCredentialsStore,
			Configuration: kongIngressStore,
		},
		isValidIngresClass: annotations.IngressClassValidatorFuncFromObjectMeta("kong"),
	}
	return s, nil
}
