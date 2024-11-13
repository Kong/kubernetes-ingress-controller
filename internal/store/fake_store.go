package store

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/go-logr/zapr"
	"github.com/samber/lo"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/tools/cache"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	"sigs.k8s.io/yaml"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// FakeObjects can be used to populate a fake Store.
type FakeObjects struct {
	IngressesV1                    []*netv1.Ingress
	IngressClassesV1               []*netv1.IngressClass
	HTTPRoutes                     []*gatewayapi.HTTPRoute
	UDPRoutes                      []*gatewayapi.UDPRoute
	TCPRoutes                      []*gatewayapi.TCPRoute
	TLSRoutes                      []*gatewayapi.TLSRoute
	GRPCRoutes                     []*gatewayapi.GRPCRoute
	ReferenceGrants                []*gatewayapi.ReferenceGrant
	Gateways                       []*gatewayapi.Gateway
	TCPIngresses                   []*kongv1beta1.TCPIngress
	UDPIngresses                   []*kongv1beta1.UDPIngress
	IngressClassParametersV1alpha1 []*kongv1alpha1.IngressClassParameters
	Services                       []*corev1.Service
	EndpointSlices                 []*discoveryv1.EndpointSlice
	Secrets                        []*corev1.Secret
	KongPlugins                    []*kongv1.KongPlugin
	KongClusterPlugins             []*kongv1.KongClusterPlugin
	KongIngresses                  []*kongv1.KongIngress
	KongConsumers                  []*kongv1.KongConsumer
	KongConsumerGroups             []*kongv1beta1.KongConsumerGroup
	KongUpstreamPolicies           []*kongv1beta1.KongUpstreamPolicy
	KongServiceFacades             []*incubatorv1alpha1.KongServiceFacade
	KongVaults                     []*kongv1alpha1.KongVault
	KongCustomEntities             []*kongv1alpha1.KongCustomEntity
}

// NewFakeStore creates a store backed by the objects passed in as arguments.
func NewFakeStore(
	objects FakeObjects,
) (Storer, error) {
	var s Storer

	ingressV1Store := cache.NewStore(namespacedKeyFunc)
	for _, ingress := range objects.IngressesV1 {
		err := ingressV1Store.Add(ingress)
		if err != nil {
			return nil, err
		}
	}
	ingressClassV1Store := cache.NewStore(clusterWideKeyFunc)
	for _, ingress := range objects.IngressClassesV1 {
		err := ingressClassV1Store.Add(ingress)
		if err != nil {
			return nil, err
		}
	}
	IngressClassParametersV1alpha1Store := cache.NewStore(clusterWideKeyFunc)
	for _, IngressClassParametersV1alpha1 := range objects.IngressClassParametersV1alpha1 {
		err := IngressClassParametersV1alpha1Store.Add(IngressClassParametersV1alpha1)
		if err != nil {
			return nil, err
		}
	}
	httprouteStore := cache.NewStore(namespacedKeyFunc)
	for _, httproute := range objects.HTTPRoutes {
		if err := httprouteStore.Add(httproute); err != nil {
			return nil, err
		}
	}
	udprouteStore := cache.NewStore(namespacedKeyFunc)
	for _, udproute := range objects.UDPRoutes {
		if err := udprouteStore.Add(udproute); err != nil {
			return nil, err
		}
	}
	tcprouteStore := cache.NewStore(namespacedKeyFunc)
	for _, tcproute := range objects.TCPRoutes {
		if err := tcprouteStore.Add(tcproute); err != nil {
			return nil, err
		}
	}
	tlsrouteStore := cache.NewStore(namespacedKeyFunc)
	for _, tlsroute := range objects.TLSRoutes {
		if err := tlsrouteStore.Add(tlsroute); err != nil {
			return nil, err
		}
	}
	grpcrouteStore := cache.NewStore(namespacedKeyFunc)
	for _, grpcroute := range objects.GRPCRoutes {
		if err := grpcrouteStore.Add(grpcroute); err != nil {
			return nil, err
		}
	}
	referencegrantStore := cache.NewStore(namespacedKeyFunc)
	for _, referencegrant := range objects.ReferenceGrants {
		if err := referencegrantStore.Add(referencegrant); err != nil {
			return nil, err
		}
	}
	gatewayStore := cache.NewStore(namespacedKeyFunc)
	for _, gw := range objects.Gateways {
		if err := gatewayStore.Add(gw); err != nil {
			return nil, err
		}
	}
	tcpIngressStore := cache.NewStore(namespacedKeyFunc)
	for _, ingress := range objects.TCPIngresses {
		err := tcpIngressStore.Add(ingress)
		if err != nil {
			return nil, err
		}
	}
	udpIngressStore := cache.NewStore(namespacedKeyFunc)
	for _, ingress := range objects.UDPIngresses {
		if err := udpIngressStore.Add(ingress); err != nil {
			return nil, err
		}
	}
	serviceStore := cache.NewStore(namespacedKeyFunc)
	for _, s := range objects.Services {
		err := serviceStore.Add(s)
		if err != nil {
			return nil, err
		}
	}
	secretsStore := cache.NewStore(namespacedKeyFunc)
	for _, s := range objects.Secrets {
		err := secretsStore.Add(s)
		if err != nil {
			return nil, err
		}
	}
	endpointSliceStore := cache.NewStore(namespacedKeyFunc)
	for _, e := range objects.EndpointSlices {
		err := endpointSliceStore.Add(e)
		if err != nil {
			return nil, err
		}
	}
	kongIngressStore := cache.NewStore(namespacedKeyFunc)
	for _, k := range objects.KongIngresses {
		err := kongIngressStore.Add(k)
		if err != nil {
			return nil, err
		}
	}
	consumerStore := cache.NewStore(namespacedKeyFunc)
	for _, c := range objects.KongConsumers {
		err := consumerStore.Add(c)
		if err != nil {
			return nil, err
		}
	}
	consumerGroupStore := cache.NewStore(namespacedKeyFunc)
	for _, c := range objects.KongConsumerGroups {
		err := consumerGroupStore.Add(c)
		if err != nil {
			return nil, err
		}
	}
	kongPluginsStore := cache.NewStore(namespacedKeyFunc)
	for _, p := range objects.KongPlugins {
		err := kongPluginsStore.Add(p)
		if err != nil {
			return nil, err
		}
	}
	kongClusterPluginsStore := cache.NewStore(clusterWideKeyFunc)
	for _, p := range objects.KongClusterPlugins {
		err := kongClusterPluginsStore.Add(p)
		if err != nil {
			return nil, err
		}
	}
	kongUpstreamPolicyStore := cache.NewStore(namespacedKeyFunc)
	for _, p := range objects.KongUpstreamPolicies {
		err := kongUpstreamPolicyStore.Add(p)
		if err != nil {
			return nil, err
		}
	}
	kongServiceFacade := cache.NewStore(namespacedKeyFunc)
	for _, s := range objects.KongServiceFacades {
		err := kongServiceFacade.Add(s)
		if err != nil {
			return nil, err
		}
	}
	kongVaultStore := cache.NewStore(clusterWideKeyFunc)
	for _, v := range objects.KongVaults {
		err := kongVaultStore.Add(v)
		if err != nil {
			return nil, err
		}
	}
	kongCustomEntityStore := cache.NewStore(namespacedKeyFunc)
	for _, e := range objects.KongCustomEntities {
		if err := kongCustomEntityStore.Add(e); err != nil {
			return nil, err
		}
	}

	s = &Store{
		stores: CacheStores{
			IngressV1:                      ingressV1Store,
			IngressClassV1:                 ingressClassV1Store,
			HTTPRoute:                      httprouteStore,
			UDPRoute:                       udprouteStore,
			TCPRoute:                       tcprouteStore,
			TLSRoute:                       tlsrouteStore,
			GRPCRoute:                      grpcrouteStore,
			ReferenceGrant:                 referencegrantStore,
			Gateway:                        gatewayStore,
			TCPIngress:                     tcpIngressStore,
			UDPIngress:                     udpIngressStore,
			Service:                        serviceStore,
			EndpointSlice:                  endpointSliceStore,
			Secret:                         secretsStore,
			Plugin:                         kongPluginsStore,
			ClusterPlugin:                  kongClusterPluginsStore,
			Consumer:                       consumerStore,
			ConsumerGroup:                  consumerGroupStore,
			KongIngress:                    kongIngressStore,
			IngressClassParametersV1alpha1: IngressClassParametersV1alpha1Store,
			KongUpstreamPolicy:             kongUpstreamPolicyStore,
			KongServiceFacade:              kongServiceFacade,
			KongVault:                      kongVaultStore,
			KongCustomEntity:               kongCustomEntityStore,
		},
		ingressClass:          annotations.DefaultIngressClass,
		isValidIngressClass:   annotations.IngressClassValidatorFuncFromObjectMeta(annotations.DefaultIngressClass),
		isValidIngressV1Class: annotations.IngressClassValidatorFuncFromV1Ingress(annotations.DefaultIngressClass),
		ingressClassMatching:  annotations.ExactClassMatch,
		logger:                zapr.NewLogger(zap.NewNop()),
	}
	return s, nil
}

// MarshalToYAML marshals the contents of every object in the store as YAML, separated by "---".
// This is useful for debugging.
func (objects FakeObjects) MarshalToYAML() ([]byte, error) {
	// In many cases objects we'd like to dump do not have their GVK set, so we need to set it manually based on
	// their known type - otherwise the YAML dump will not work.
	typeToGVK := map[reflect.Type]schema.GroupVersionKind{
		reflect.TypeOf(&netv1.Ingress{}):                       netv1.SchemeGroupVersion.WithKind("Ingress"),
		reflect.TypeOf(&netv1.IngressClass{}):                  netv1.SchemeGroupVersion.WithKind("IngressClass"),
		reflect.TypeOf(&gatewayapi.HTTPRoute{}):                gatewayv1.SchemeGroupVersion.WithKind("HTTPRoute"),
		reflect.TypeOf(&gatewayapi.UDPRoute{}):                 gatewayv1alpha2.SchemeGroupVersion.WithKind("UDPRoute"),
		reflect.TypeOf(&gatewayapi.TCPRoute{}):                 gatewayv1alpha2.SchemeGroupVersion.WithKind("TCPRoute"),
		reflect.TypeOf(&gatewayapi.TLSRoute{}):                 gatewayv1alpha2.SchemeGroupVersion.WithKind("TLSRoute"),
		reflect.TypeOf(&gatewayapi.GRPCRoute{}):                gatewayv1.SchemeGroupVersion.WithKind("GRPCRoute"),
		reflect.TypeOf(&gatewayapi.ReferenceGrant{}):           gatewayv1beta1.SchemeGroupVersion.WithKind("ReferenceGrant"),
		reflect.TypeOf(&gatewayapi.Gateway{}):                  gatewayv1.SchemeGroupVersion.WithKind("Gateway"),
		reflect.TypeOf(&kongv1beta1.TCPIngress{}):              kongv1beta1.SchemeGroupVersion.WithKind("TCPIngress"),
		reflect.TypeOf(&kongv1beta1.UDPIngress{}):              kongv1beta1.SchemeGroupVersion.WithKind("UDPIngress"),
		reflect.TypeOf(&kongv1alpha1.IngressClassParameters{}): kongv1alpha1.SchemeGroupVersion.WithKind("IngressClassParameters"),
		reflect.TypeOf(&corev1.Service{}):                      corev1.SchemeGroupVersion.WithKind("Service"),
		reflect.TypeOf(&discoveryv1.EndpointSlice{}):           discoveryv1.SchemeGroupVersion.WithKind("EndpointSlice"),
		reflect.TypeOf(&corev1.Secret{}):                       corev1.SchemeGroupVersion.WithKind("Secret"),
		reflect.TypeOf(&kongv1.KongPlugin{}):                   kongv1.SchemeGroupVersion.WithKind("KongPlugin"),
		reflect.TypeOf(&kongv1.KongClusterPlugin{}):            kongv1.SchemeGroupVersion.WithKind("KongClusterPlugin"),
		reflect.TypeOf(&kongv1.KongIngress{}):                  kongv1.SchemeGroupVersion.WithKind("KongIngress"),
		reflect.TypeOf(&kongv1.KongConsumer{}):                 kongv1.SchemeGroupVersion.WithKind("KongConsumer"),
		reflect.TypeOf(&kongv1beta1.KongConsumerGroup{}):       kongv1beta1.SchemeGroupVersion.WithKind("KongConsumerGroup"),
		reflect.TypeOf(&kongv1alpha1.KongVault{}):              kongv1alpha1.SchemeGroupVersion.WithKind(kongv1alpha1.KongVaultKind),
		reflect.TypeOf(&kongv1alpha1.KongCustomEntity{}):       kongv1alpha1.SchemeGroupVersion.WithKind(kongv1alpha1.KongCustomEntityKind),
	}

	out := &bytes.Buffer{}

	// fillGVKAndAppendToBuffer is a helper function that sets the GVK of the given object and appends it to the output.
	fillGVKAndAppendToBuffer := func(obj runtime.Object) error {
		gvk, ok := typeToGVK[reflect.TypeOf(obj)]
		if !ok {
			return fmt.Errorf("unknown type: %T", obj)
		}
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		b, err := marshalObjToYAML(obj)
		if err != nil {
			return err
		}
		out.WriteString("---\n")
		out.Write(b)
		return nil
	}

	// Let's gather all objects in a single generic slice.
	var allObjects []any
	allObjects = append(allObjects, lo.ToAnySlice(objects.IngressesV1)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.IngressClassesV1)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.HTTPRoutes)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.UDPRoutes)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.TCPRoutes)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.TLSRoutes)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.GRPCRoutes)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.ReferenceGrants)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.Gateways)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.TCPIngresses)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.UDPIngresses)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.IngressClassParametersV1alpha1)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.Services)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.EndpointSlices)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.Secrets)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.KongPlugins)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.KongClusterPlugins)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.KongIngresses)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.KongConsumers)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.KongConsumerGroups)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.KongVaults)...)
	allObjects = append(allObjects, lo.ToAnySlice(objects.KongCustomEntities)...)

	for _, obj := range allObjects {
		if err := fillGVKAndAppendToBuffer(obj.(runtime.Object)); err != nil {
			return nil, err
		}
	}

	return out.Bytes(), nil
}

// marshalObjToYAML marshals the given object as YAML.
// It uses the JSON printer to dump the object as JSON, and then converts the JSON to YAML because some of the objects
// we want to dump do not have YAML tags.
func marshalObjToYAML(obj runtime.Object) ([]byte, error) {
	buff := bytes.Buffer{}
	printer := printers.JSONPrinter{}
	err := printer.PrintObj(obj, &buff)
	if err != nil {
		return nil, err
	}

	return yaml.JSONToYAML(buff.Bytes())
}
