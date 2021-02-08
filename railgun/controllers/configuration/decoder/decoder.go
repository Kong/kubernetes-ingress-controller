package decoder

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

const KeyDelimiter = "_"

func isGV(want schema.GroupVersion, gotGroup, gotVersion string) bool {
	return want.Group == gotGroup && want.Version == gotVersion
}

func DecodeObject(key string, value []byte) (client.Object, error) {
	const wantLen = 5
	keyElems := strings.SplitN(key, KeyDelimiter, wantLen+1)
	if len(keyElems) != 5 {
		return nil, fmt.Errorf("key had %d elements, expected 5", len(keyElems))
	}
	group, version, kind, namespace, name :=
		keyElems[0], keyElems[1], keyElems[2], keyElems[3], keyElems[4]

	var result client.Object

	switch {
	case isGV(corev1.SchemeGroupVersion, group, version) && kind == "Service":
		result = new(corev1.Service)
	case isGV(corev1.SchemeGroupVersion, group, version) && kind == "Endpoints":
		result = new(corev1.Endpoints)
	case isGV(corev1.SchemeGroupVersion, group, version) && kind == "Secret":
		result = new(corev1.Secret)
	case isGV(networkingv1beta1.SchemeGroupVersion, group, version) && kind == "Ingress":
		result = new(networkingv1beta1.Ingress)
	case isGV(networkingv1.SchemeGroupVersion, group, version) && kind == "Ingress":
		result = new(networkingv1.Ingress)
	case isGV(configurationv1beta1.SchemeGroupVersion, group, version) && kind == "TCPIngress":
		result = new(configurationv1beta1.TCPIngress)
	case isGV(configurationv1beta1.SchemeGroupVersion, group, version) && kind == "KongPlugin":
		result = new(configurationv1.KongPlugin)
	case isGV(configurationv1beta1.SchemeGroupVersion, group, version) && kind == "KongClusterPlugin":
		result = new(configurationv1.KongClusterPlugin)
	case isGV(configurationv1beta1.SchemeGroupVersion, group, version) && kind == "KongIngress":
		result = new(configurationv1.KongIngress)
	case isGV(configurationv1beta1.SchemeGroupVersion, group, version) && kind == "KongConsumer":
		result = new(configurationv1.KongConsumer)
	case isGV(knative.SchemeGroupVersion, group, version) && kind == "Ingress":
		result = new(knative.Ingress)
	}

	if err := json.Unmarshal(value, result); err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	if namespace != result.GetNamespace() || name != result.GetName() {
		return nil, fmt.Errorf("NS/name of the object (%q, %q) does not match key (%q, %q)",
			result.GetNamespace(), result.GetName(), namespace, name)
	}

	return result, nil
}

func DecodeSecret(s *corev1.Secret) (store.Storer, error) {
	return nil, nil
}
