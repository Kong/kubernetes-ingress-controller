package manager

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/configuration"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/knative"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object/status"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	konghqcomv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
	konghqcomv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
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
	mgr manager.Manager,
	dataplaneClient *dataplane.KongClient,
	dataplaneAddressFinder *dataplane.AddressFinder,
	udpDataplaneAddressFinder *dataplane.AddressFinder,
	kubernetesStatusQueue *status.Queue,
	c *Config,
	featureGates map[string]bool,
	kongAdminAPIEndpointsNotifier configuration.EndpointsNotifier,
) ([]ControllerDef, error) {
	restMapper := mgr.GetClient().RESTMapper()

	// Choose the best API version of Ingress to inform which ingress controller to enable.
	ingressConditions, err := NewIngressControllersConditions(c, restMapper)
	if err != nil {
		return nil, fmt.Errorf("ingress version picker failed: %w", err)
	}

	referenceGrantsEnabled := featureGates[gatewayFeature] && ShouldEnableCRDController(
		schema.GroupVersionResource{
			Group:    gatewayv1beta1.GroupVersion.Group,
			Version:  gatewayv1beta1.GroupVersion.Version,
			Resource: "referencegrants",
		},
		restMapper,
	)

	referenceIndexers := ctrlref.NewCacheIndexers()

	controllers := []ControllerDef{
		// ---------------------------------------------------------------------------
		// Kong Gateway Admin API Service discovery
		// ---------------------------------------------------------------------------
		{
			Enabled: c.KongAdminSvc.Name != "",
			Controller: &configuration.KongAdminAPIServiceReconciler{
				Client:            mgr.GetClient(),
				ServiceNN:         c.KongAdminSvc,
				Log:               ctrl.Log.WithName("controllers").WithName("KongAdminAPIService"),
				CacheSyncTimeout:  c.CacheSyncTimeout,
				EndpointsNotifier: kongAdminAPIEndpointsNotifier,
			},
		},
		// ---------------------------------------------------------------------------
		// Core API Controllers
		// ---------------------------------------------------------------------------
		{
			Enabled: ingressConditions.IngressClassNetV1Enabled(),
			Controller: &configuration.NetV1IngressClassReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("IngressClass").WithName("netv1"),
				DataplaneClient:  dataplaneClient,
				Scheme:           mgr.GetScheme(),
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: ingressConditions.IngressNetV1Enabled(),
			Controller: &configuration.NetV1IngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1"),
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
			Enabled: ingressConditions.IngressNetV1beta1Enabled(),
			Controller: &configuration.NetV1Beta1IngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1beta1"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				IngressClassName: c.IngressClassName,
				// this and other resources that support class get an additional watch to account for the default
				// IngressClass even if the cluster uses an Ingress version other than networking/v1 (the only version
				// we support IngressClass for). we pass the v1 controller disable flag to them to avoid
				// https://github.com/Kong/kubernetes-ingress-controller/issues/2563
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				DataplaneAddressFinder:     dataplaneAddressFinder,
				CacheSyncTimeout:           c.CacheSyncTimeout,
				ReferenceIndexers:          referenceIndexers,
			},
		},
		{
			Enabled: ingressConditions.IngressExtV1beta1Enabled(),
			Controller: &configuration.ExtV1Beta1IngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("Ingress").WithName("extv1beta1"),
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
				Log:               ctrl.Log.WithName("controllers").WithName("Service"),
				Scheme:            mgr.GetScheme(),
				DataplaneClient:   dataplaneClient,
				CacheSyncTimeout:  c.CacheSyncTimeout,
				ReferenceIndexers: referenceIndexers,
			},
		},
		{
			Enabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1EndpointsReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Endpoints"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: true,
			Controller: &configuration.CoreV1SecretReconciler{
				Client:            mgr.GetClient(),
				Log:               ctrl.Log.WithName("controllers").WithName("Secrets"),
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
			Enabled: c.UDPIngressEnabled && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    konghqcomv1beta1.GroupVersion.Group,
					Version:  konghqcomv1beta1.GroupVersion.Version,
					Resource: "udpingresses",
				},
				restMapper,
			),
			Controller: &configuration.KongV1Beta1UDPIngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("UDPIngress"),
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
			Enabled: c.TCPIngressEnabled && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    konghqcomv1beta1.GroupVersion.Group,
					Version:  konghqcomv1beta1.GroupVersion.Version,
					Resource: "tcpingresses",
				},
				restMapper,
			),
			Controller: &configuration.KongV1Beta1TCPIngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("TCPIngress"),
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
			Enabled: c.KongIngressEnabled && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    konghqcomv1.GroupVersion.Group,
					Version:  konghqcomv1.GroupVersion.Version,
					Resource: "kongingresses",
				},
				restMapper,
			),
			Controller: &configuration.KongV1KongIngressReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("KongIngress"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: c.IngressClassParametersEnabled && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    konghqcomv1alpha1.GroupVersion.Group,
					Version:  konghqcomv1alpha1.GroupVersion.Version,
					Resource: "ingressclassparameterses",
				},
				restMapper,
			),
			Controller: &configuration.KongV1Alpha1IngressClassParametersReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("IngressClassParameters"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: c.KongPluginEnabled && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    konghqcomv1.GroupVersion.Group,
					Version:  konghqcomv1.GroupVersion.Version,
					Resource: "kongplugins",
				},
				restMapper,
			),
			Controller: &configuration.KongV1KongPluginReconciler{
				Client:            mgr.GetClient(),
				Log:               ctrl.Log.WithName("controllers").WithName("KongPlugin"),
				Scheme:            mgr.GetScheme(),
				DataplaneClient:   dataplaneClient,
				CacheSyncTimeout:  c.CacheSyncTimeout,
				ReferenceIndexers: referenceIndexers,
			},
		},
		{
			Enabled: c.KongConsumerEnabled && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    konghqcomv1.GroupVersion.Group,
					Version:  konghqcomv1.GroupVersion.Version,
					Resource: "kongconsumers",
				},
				restMapper,
			),
			Controller: &configuration.KongV1KongConsumerReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("KongConsumer"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				CacheSyncTimeout:           c.CacheSyncTimeout,
				ReferenceIndexers:          referenceIndexers,
			},
		},
		{
			Enabled: c.KongClusterPluginEnabled && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    konghqcomv1.GroupVersion.Group,
					Version:  konghqcomv1.GroupVersion.Version,
					Resource: "kongclusterplugins",
				},
				restMapper,
			),
			Controller: &configuration.KongV1KongClusterPluginReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				CacheSyncTimeout:           c.CacheSyncTimeout,
				ReferenceIndexers:          referenceIndexers,
			},
		},
		// ---------------------------------------------------------------------------
		// Other Controllers
		// ---------------------------------------------------------------------------
		{
			// knative is a special case because it existed before we added feature gates functionality
			// for this controller (only) the existing --enable-controller-knativeingress flag overrides
			// any feature gate configuration. See FEATURE_GATES.md for more information.
			Enabled: (featureGates[knativeFeature] || c.KnativeIngressEnabled) && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    knativev1alpha1.SchemeGroupVersion.Group,
					Version:  knativev1alpha1.SchemeGroupVersion.Version,
					Resource: "ingresses",
				},
				restMapper,
			),
			Controller: &knative.Knativev1alpha1IngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("Ingress").WithName("KnativeV1Alpha1"),
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
		// ---------------------------------------------------------------------------
		// Gateway API Controllers - Beta APIs
		// ---------------------------------------------------------------------------
		{
			Enabled: featureGates[gatewayFeature] && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    gatewayv1beta1.GroupVersion.Group,
					Version:  gatewayv1beta1.GroupVersion.Version,
					Resource: "gateways",
				},
				restMapper,
			),
			Controller: &gateway.GatewayReconciler{
				Client:               mgr.GetClient(),
				Log:                  ctrl.Log.WithName("controllers").WithName(gatewayFeature),
				Scheme:               mgr.GetScheme(),
				DataplaneClient:      dataplaneClient,
				PublishService:       c.PublishService.String(),
				WatchNamespaces:      c.WatchNamespaces,
				EnableReferenceGrant: referenceGrantsEnabled,
				CacheSyncTimeout:     c.CacheSyncTimeout,
				ReferenceIndexers:    referenceIndexers,
			},
		},
		{
			Enabled: featureGates[gatewayFeature] && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    gatewayv1beta1.GroupVersion.Group,
					Version:  gatewayv1beta1.GroupVersion.Version,
					Resource: "httproutes",
				},
				restMapper,
			),
			Controller: &gateway.HTTPRouteReconciler{
				Client:               mgr.GetClient(),
				Log:                  ctrl.Log.WithName("controllers").WithName("HTTPRoute"),
				Scheme:               mgr.GetScheme(),
				DataplaneClient:      dataplaneClient,
				EnableReferenceGrant: referenceGrantsEnabled,
				CacheSyncTimeout:     c.CacheSyncTimeout,
			},
		},
		// ---------------------------------------------------------------------------
		// Gateway API Controllers - Alpha APIs
		// ---------------------------------------------------------------------------
		{
			Enabled: referenceGrantsEnabled,
			Controller: &gateway.ReferenceGrantReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("ReferenceGrant"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: featureGates[gatewayAlphaFeature] && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    gatewayv1alpha2.GroupVersion.Group,
					Version:  gatewayv1alpha2.GroupVersion.Version,
					Resource: "udproutes",
				},
				restMapper,
			),
			Controller: &gateway.UDPRouteReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("UDPRoute"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: featureGates[gatewayAlphaFeature] && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    gatewayv1alpha2.GroupVersion.Group,
					Version:  gatewayv1alpha2.GroupVersion.Version,
					Resource: "tcproutes",
				},
				restMapper,
			),
			Controller: &gateway.TCPRouteReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("TCPRoute"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: featureGates[gatewayAlphaFeature] && ShouldEnableCRDController(
				schema.GroupVersionResource{
					Group:    gatewayv1alpha2.GroupVersion.Group,
					Version:  gatewayv1alpha2.GroupVersion.Version,
					Resource: "tlsroutes",
				},
				restMapper,
			),
			Controller: &gateway.TLSRouteReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("TLSRoute"),
				Scheme:           mgr.GetScheme(),
				DataplaneClient:  dataplaneClient,
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
	}

	return controllers, nil
}
