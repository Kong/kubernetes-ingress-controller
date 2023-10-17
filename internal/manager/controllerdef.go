package manager

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/configuration"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/crds"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object/status"
)

// -----------------------------------------------------------------------------
// Controller Manager - Controller Definition Interfaces
// -----------------------------------------------------------------------------

// Controller is a Kubernetes controller that can be plugged into Manager.
type Controller interface {
	SetupWithManager(ctrl.Manager) error
}

// ControllerDef is a specification of a Controller that can be conditionally registered with Manager.
type ControllerDef struct {
	Enabled    bool
	Controller Controller
}

// Name returns a human-readable name of the controller.
func (c *ControllerDef) Name() string {
	return reflect.TypeOf(c.Controller).String()
}

// MaybeSetupWithManager runs SetupWithManager on the controller if it is enabled.
func (c *ControllerDef) MaybeSetupWithManager(mgr ctrl.Manager) error {
	if !c.Enabled {
		return nil
	}

	return c.Controller.SetupWithManager(mgr)
}

// -----------------------------------------------------------------------------
// Controller Manager - Controller Setup Functions
// -----------------------------------------------------------------------------

func setupControllers(
	ctx context.Context,
	mgr manager.Manager,
	dataplaneClient controllers.DataPlane,
	dataplaneAddressFinder *dataplane.AddressFinder,
	udpDataplaneAddressFinder *dataplane.AddressFinder,
	kubernetesStatusQueue *status.Queue,
	c *Config,
	featureGates map[string]bool,
	kongAdminAPIEndpointsNotifier configuration.EndpointsNotifier,
	adminAPIsDiscoverer configuration.AdminAPIsDiscoverer,
) []ControllerDef {
	referenceIndexers := ctrlref.NewCacheIndexers(ctrl.LoggerFrom(ctx).WithName("controllers").WithName("reference-indexers"))

	controllers := []ControllerDef{
		// ---------------------------------------------------------------------------
		// Kong Gateway Admin API Service discovery
		// ---------------------------------------------------------------------------
		{
			Enabled: c.KongAdminSvc.IsPresent(),
			Controller: &configuration.KongAdminAPIServiceReconciler{
				Client:              mgr.GetClient(),
				ServiceNN:           c.KongAdminSvc.OrEmpty(),
				Log:                 ctrl.LoggerFrom(ctx).WithName("controllers").WithName("KongAdminAPIService"),
				CacheSyncTimeout:    c.CacheSyncTimeout,
				EndpointsNotifier:   kongAdminAPIEndpointsNotifier,
				AdminAPIsDiscoverer: adminAPIsDiscoverer,
			},
		},
		// ---------------------------------------------------------------------------
		// Core API Controllers
		// ---------------------------------------------------------------------------
		{
			Enabled: c.IngressClassNetV1Enabled,
			Controller: &configuration.NetV1IngressClassReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("IngressClass").WithName("netv1"),
				DataplaneClient:  dataplaneClient,
				Scheme:           mgr.GetScheme(),
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: c.IngressNetV1Enabled,
			Controller: &configuration.NetV1IngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Ingress").WithName("netv1"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				DataplaneAddressFinder:     dataplaneAddressFinder,
				CacheSyncTimeout:           c.CacheSyncTimeout,
				ReferenceIndexers:          referenceIndexers,
			},
		},
		{
			Enabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1ServiceReconciler{
				Client:            mgr.GetClient(),
				Log:               ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Service"),
				Scheme:            mgr.GetScheme(),
				DataplaneClient:   dataplaneClient,
				CacheSyncTimeout:  c.CacheSyncTimeout,
				ReferenceIndexers: referenceIndexers,
			},
		},
		{
			Enabled: c.ServiceEnabled,
			Controller: &configuration.DiscoveryV1EndpointSliceReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("EndpointSlice"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: true,
			Controller: &configuration.CoreV1SecretReconciler{
				Client:            mgr.GetClient(),
				Log:               ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Secrets"),
				Scheme:            mgr.GetScheme(),
				DataplaneClient:   dataplaneClient,
				CacheSyncTimeout:  c.CacheSyncTimeout,
				ReferenceIndexers: referenceIndexers,
			},
		},
		// ---------------------------------------------------------------------------
		// Kong API Controllers
		// ---------------------------------------------------------------------------
		{
			Enabled: c.UDPIngressEnabled,
			Controller: &configuration.KongV1Beta1UDPIngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.LoggerFrom(ctx).WithName("controllers").WithName("UDPIngress"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				DataplaneAddressFinder:     udpDataplaneAddressFinder,
				CacheSyncTimeout:           c.CacheSyncTimeout,
			},
		},
		{
			Enabled: c.TCPIngressEnabled,
			Controller: &configuration.KongV1Beta1TCPIngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.LoggerFrom(ctx).WithName("controllers").WithName("TCPIngress"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				DataplaneAddressFinder:     dataplaneAddressFinder,
				CacheSyncTimeout:           c.CacheSyncTimeout,
				ReferenceIndexers:          referenceIndexers,
			},
		},
		{
			Enabled: c.KongIngressEnabled,
			Controller: &configuration.KongV1KongIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("KongIngress"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: c.IngressClassParametersEnabled,
			Controller: &configuration.KongV1Alpha1IngressClassParametersReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("IngressClassParameters"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: c.KongPluginEnabled,
			Controller: &configuration.KongV1KongPluginReconciler{
				Client:            mgr.GetClient(),
				Log:               ctrl.LoggerFrom(ctx).WithName("controllers").WithName("KongPlugin"),
				Scheme:            mgr.GetScheme(),
				DataplaneClient:   dataplaneClient,
				CacheSyncTimeout:  c.CacheSyncTimeout,
				ReferenceIndexers: referenceIndexers,
				// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/4578
				// StatusQueue:       kubernetesStatusQueue,
			},
		},
		{
			Enabled: c.KongConsumerEnabled,
			Controller: &configuration.KongV1KongConsumerReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.LoggerFrom(ctx).WithName("controllers").WithName("KongConsumer"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				CacheSyncTimeout:           c.CacheSyncTimeout,
				ReferenceIndexers:          referenceIndexers,
				StatusQueue:                kubernetesStatusQueue,
			},
		},
		{
			Enabled: c.KongConsumerEnabled,
			Controller: &configuration.KongV1Beta1KongConsumerGroupReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.LoggerFrom(ctx).WithName("controllers").WithName("KongConsumerGroup"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				CacheSyncTimeout:           c.CacheSyncTimeout,
				ReferenceIndexers:          referenceIndexers,
				StatusQueue:                kubernetesStatusQueue,
			},
		},
		{
			Enabled: c.KongClusterPluginEnabled,
			Controller: &configuration.KongV1KongClusterPluginReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.LoggerFrom(ctx).WithName("controllers").WithName("KongClusterPlugin"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				CacheSyncTimeout:           c.CacheSyncTimeout,
				ReferenceIndexers:          referenceIndexers,
				// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/4578
				// StatusQueue:       kubernetesStatusQueue,
			},
		},
		// ---------------------------------------------------------------------------
		// Gateway API Controllers - Beta APIs
		// ---------------------------------------------------------------------------
		{
			Enabled: featureGates[featuregates.GatewayFeature],
			Controller: &crds.DynamicCRDController{
				Manager:          mgr,
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Dynamic/Gateway"),
				CacheSyncTimeout: c.CacheSyncTimeout,
				RequiredCRDs:     baseGatewayCRDs(),
				Controller: &gateway.GatewayReconciler{
					Client:               mgr.GetClient(),
					Log:                  ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Gateway"),
					Scheme:               mgr.GetScheme(),
					DataplaneClient:      dataplaneClient,
					IngressServiceRef:    c.IngressService.OrEmpty(),
					IngressServiceUDPRef: c.IngressServiceUDP,
					WatchNamespaces:      c.WatchNamespaces,
					CacheSyncTimeout:     c.CacheSyncTimeout,
					ReferenceIndexers:    referenceIndexers,
				},
			},
		},
		{
			Enabled: featureGates[featuregates.GatewayFeature],
			Controller: &crds.DynamicCRDController{
				Manager:          mgr,
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Dynamic/HTTPRoute"),
				CacheSyncTimeout: c.CacheSyncTimeout,
				RequiredCRDs: append(baseGatewayCRDs(), schema.GroupVersionResource{
					Group:    gatewayv1beta1.GroupVersion.Group,
					Version:  gatewayv1beta1.GroupVersion.Version,
					Resource: "httproutes",
				}),
				Controller: &gateway.HTTPRouteReconciler{
					Client:           mgr.GetClient(),
					Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("HTTPRoute"),
					Scheme:           mgr.GetScheme(),
					DataplaneClient:  dataplaneClient,
					CacheSyncTimeout: c.CacheSyncTimeout,
					StatusQueue:      kubernetesStatusQueue,
				},
			},
		},
		{
			Enabled: featureGates[featuregates.GatewayFeature],
			Controller: &crds.DynamicCRDController{
				Manager:          mgr,
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Dynamic/ReferenceGrant"),
				CacheSyncTimeout: c.CacheSyncTimeout,
				RequiredCRDs: append(baseGatewayCRDs(), schema.GroupVersionResource{
					Group:    gatewayv1beta1.GroupVersion.Group,
					Version:  gatewayv1beta1.GroupVersion.Version,
					Resource: "referencegrants",
				}),
				Controller: &gateway.ReferenceGrantReconciler{
					Client:           mgr.GetClient(),
					Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("ReferenceGrant"),
					Scheme:           mgr.GetScheme(),
					DataplaneClient:  dataplaneClient,
					CacheSyncTimeout: c.CacheSyncTimeout,
				},
			},
		},
		// ---------------------------------------------------------------------------
		// Gateway API Controllers - Alpha APIs
		// ---------------------------------------------------------------------------
		{
			Enabled: featureGates[featuregates.GatewayAlphaFeature],
			Controller: &crds.DynamicCRDController{
				Manager:          mgr,
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Dynamic/UDPRoute"),
				CacheSyncTimeout: c.CacheSyncTimeout,
				RequiredCRDs: append(baseGatewayCRDs(), schema.GroupVersionResource{
					Group:    gatewayv1alpha2.GroupVersion.Group,
					Version:  gatewayv1alpha2.GroupVersion.Version,
					Resource: "udproutes",
				}),
				Controller: &gateway.UDPRouteReconciler{
					Client:           mgr.GetClient(),
					Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("UDPRoute"),
					Scheme:           mgr.GetScheme(),
					DataplaneClient:  dataplaneClient,
					CacheSyncTimeout: c.CacheSyncTimeout,
					StatusQueue:      kubernetesStatusQueue,
				},
			},
		},
		{
			Enabled: featureGates[featuregates.GatewayAlphaFeature],
			Controller: &crds.DynamicCRDController{
				Manager:          mgr,
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Dynamic/TCPRoute"),
				CacheSyncTimeout: c.CacheSyncTimeout,
				RequiredCRDs: append(baseGatewayCRDs(), schema.GroupVersionResource{
					Group:    gatewayv1alpha2.GroupVersion.Group,
					Version:  gatewayv1alpha2.GroupVersion.Version,
					Resource: "tcproutes",
				}),
				Controller: &gateway.TCPRouteReconciler{
					Client:           mgr.GetClient(),
					Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("TCPRoute"),
					Scheme:           mgr.GetScheme(),
					DataplaneClient:  dataplaneClient,
					CacheSyncTimeout: c.CacheSyncTimeout,
					StatusQueue:      kubernetesStatusQueue,
				},
			},
		},
		{
			Enabled: featureGates[featuregates.GatewayAlphaFeature],
			Controller: &crds.DynamicCRDController{
				Manager:          mgr,
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Dynamic/TLSRoute"),
				CacheSyncTimeout: c.CacheSyncTimeout,
				RequiredCRDs: append(baseGatewayCRDs(), schema.GroupVersionResource{
					Group:    gatewayv1alpha2.GroupVersion.Group,
					Version:  gatewayv1alpha2.GroupVersion.Version,
					Resource: "tlsroutes",
				}),
				Controller: &gateway.TLSRouteReconciler{
					Client:           mgr.GetClient(),
					Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("TLSRoute"),
					Scheme:           mgr.GetScheme(),
					DataplaneClient:  dataplaneClient,
					CacheSyncTimeout: c.CacheSyncTimeout,
					StatusQueue:      kubernetesStatusQueue,
				},
			},
		},
		{
			Enabled: featureGates[featuregates.GatewayAlphaFeature],
			Controller: &crds.DynamicCRDController{
				Manager:          mgr,
				Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Dynamic/GRPCRoute"),
				CacheSyncTimeout: c.CacheSyncTimeout,
				RequiredCRDs: append(baseGatewayCRDs(), schema.GroupVersionResource{
					Group:    gatewayv1alpha2.GroupVersion.Group,
					Version:  gatewayv1alpha2.GroupVersion.Version,
					Resource: "grpcroutes",
				}),
				Controller: &gateway.GRPCRouteReconciler{
					Client:           mgr.GetClient(),
					Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("GRPCRoute"),
					Scheme:           mgr.GetScheme(),
					DataplaneClient:  dataplaneClient,
					CacheSyncTimeout: c.CacheSyncTimeout,
					StatusQueue:      kubernetesStatusQueue,
				},
			},
		},
	}

	return controllers
}

// baseGatewayCRDs returns a slice of base CRDs required for running all the Gateway API controllers.
func baseGatewayCRDs() []schema.GroupVersionResource {
	return []schema.GroupVersionResource{
		{
			Group:    gatewayv1beta1.GroupVersion.Group,
			Version:  gatewayv1beta1.GroupVersion.Version,
			Resource: "gateways",
		},
		{
			Group:    gatewayv1beta1.GroupVersion.Group,
			Version:  gatewayv1beta1.GroupVersion.Version,
			Resource: "gatewayclasses",
		},
	}
}
