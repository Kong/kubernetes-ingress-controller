package manager

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/configuration"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	ctrlutils "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object/status"
	konghqcomv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

// -----------------------------------------------------------------------------
// Controller Manager - Controller Definition Interfaces
// -----------------------------------------------------------------------------

// Controller is a Kubernetes controller that can be plugged into Manager.
type Controller interface {
	SetupWithManager(ctrl.Manager) error
}

// AutoHandler decides whether the specific controller shall be enabled (true) or disabled (false).
type AutoHandler func(client.Client) bool

// ControllerDef is a specification of a Controller that can be conditionally registered with Manager.
type ControllerDef struct {
	Enabled     bool
	AutoHandler AutoHandler
	Controller  Controller
}

// Name returns a human-readable name of the controller.
func (c *ControllerDef) Name() string {
	return reflect.TypeOf(c.Controller).String()
}

// MaybeSetupWithManager runs SetupWithManager on the controller if it is enabled
// and its AutoHandler (if any) indicates that it can load
func (c *ControllerDef) MaybeSetupWithManager(mgr ctrl.Manager) error {
	if !c.Enabled {
		return nil
	}

	if c.AutoHandler != nil {
		if enable := c.AutoHandler(mgr.GetClient()); !enable {
			return nil
		}
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
	kubernetesStatusQueue *status.Queue,
	c *Config,
	featureGates map[string]bool,
) ([]ControllerDef, error) {
	// Choose the best API version of Ingress to inform which ingress controller to enable.
	var ingressPicker ingressControllerStrategy
	if err := ingressPicker.Initialize(c, mgr.GetClient()); err != nil {
		return nil, fmt.Errorf("ingress version picker failed: %w", err)
	}

	controllers := []ControllerDef{
		// ---------------------------------------------------------------------------
		// Core API Controllers
		// ---------------------------------------------------------------------------
		{
			Enabled:     c.IngressClassNetV1Enabled,
			AutoHandler: ingressPicker.IsNetV1,
			Controller: &configuration.NetV1IngressClassReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("IngressClass").WithName("netv1"),
				DataplaneClient: dataplaneClient,
				Scheme:          mgr.GetScheme(),
			},
		},
		{
			Enabled:     c.IngressNetV1Enabled,
			AutoHandler: ingressPicker.IsNetV1,
			Controller: &configuration.NetV1IngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				DataplaneAddressFinder:     dataplaneAddressFinder,
			},
		},
		{
			Enabled:     c.IngressNetV1beta1Enabled,
			AutoHandler: ingressPicker.IsNetV1beta1,
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
			},
		},
		{
			Enabled:     c.IngressExtV1beta1Enabled,
			AutoHandler: ingressPicker.IsExtV1beta1,
			Controller: &configuration.ExtV1Beta1IngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("Ingress").WithName("extv1beta1"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				DataplaneAddressFinder:     dataplaneAddressFinder,
			},
		},
		{
			Enabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1ServiceReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("Service"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
		{
			Enabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1EndpointsReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("Endpoints"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
		{
			Enabled: true,
			Controller: &configuration.CoreV1SecretReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("Secrets"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
		// ---------------------------------------------------------------------------
		// Kong API Controllers
		// ---------------------------------------------------------------------------
		{
			Enabled: c.UDPIngressEnabled,
			Controller: &configuration.KongV1Beta1UDPIngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("UDPIngress"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				DataplaneAddressFinder:     dataplaneAddressFinder,
			},
		},
		{
			Enabled: c.TCPIngressEnabled,
			Controller: &configuration.KongV1Beta1TCPIngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("TCPIngress"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				DataplaneAddressFinder:     dataplaneAddressFinder,
			},
		},
		{
			Enabled: c.KongIngressEnabled,
			Controller: &configuration.KongV1KongIngressReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("KongIngress"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
		{
			Enabled: c.KongPluginEnabled,
			Controller: &configuration.KongV1KongPluginReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("KongPlugin"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
		{
			Enabled: c.KongConsumerEnabled,
			Controller: &configuration.KongV1KongConsumerReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("KongConsumer"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
			},
		},
		{
			Enabled: c.KongClusterPluginEnabled,
			AutoHandler: crdExistsChecker{GVR: schema.GroupVersionResource{
				Group:    konghqcomv1.SchemeGroupVersion.Group,
				Version:  konghqcomv1.SchemeGroupVersion.Version,
				Resource: "kongclusterplugins",
			}}.CRDExists,
			Controller: &configuration.KongV1KongClusterPluginReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("KongClusterPlugin"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
			},
		},
		// ---------------------------------------------------------------------------
		// Other Controllers
		// ---------------------------------------------------------------------------
		{
			// knative is a special case because it existed before we added feature gates functionality
			// for this controller (only) the existing --enable-controller-knativeingress flag overrides
			// any feature gate configuration. See FEATURE_GATES.md for more information.
			Enabled: featureGates[gatewayFeature] || c.KnativeIngressEnabled,
			AutoHandler: crdExistsChecker{GVR: schema.GroupVersionResource{
				Group:    knativev1alpha1.SchemeGroupVersion.Group,
				Version:  knativev1alpha1.SchemeGroupVersion.Version,
				Resource: "ingresses",
			}}.CRDExists,
			Controller: &configuration.Knativev1alpha1IngressReconciler{
				Client:                     mgr.GetClient(),
				Log:                        ctrl.Log.WithName("controllers").WithName("Ingress").WithName("KnativeV1Alpha1"),
				Scheme:                     mgr.GetScheme(),
				DataplaneClient:            dataplaneClient,
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				DataplaneAddressFinder:     dataplaneAddressFinder,
			},
		},
		// ---------------------------------------------------------------------------
		// GatewayAPI Controllers
		// ---------------------------------------------------------------------------
		{
			Enabled: featureGates[gatewayFeature],
			AutoHandler: crdExistsChecker{
				GVR: schema.GroupVersionResource{
					Group:    gatewayv1alpha2.SchemeGroupVersion.Group,
					Version:  gatewayv1alpha2.SchemeGroupVersion.Version,
					Resource: "gateways",
				},
			}.CRDExists,
			Controller: &gateway.GatewayReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName(gatewayFeature),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
				PublishService:  c.PublishService,
				WatchNamespaces: c.WatchNamespaces,
			},
		},
		{
			Enabled: featureGates[gatewayFeature],
			AutoHandler: crdExistsChecker{
				GVR: schema.GroupVersionResource{
					Group:    gatewayv1alpha2.SchemeGroupVersion.Group,
					Version:  gatewayv1alpha2.SchemeGroupVersion.Version,
					Resource: "referencepolicies",
				},
			}.CRDExists,
			Controller: &gateway.ReferencePolicyReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("ReferencePolicy"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
		{
			Enabled: featureGates[gatewayFeature],
			AutoHandler: crdExistsChecker{
				GVR: schema.GroupVersionResource{
					Group:    gatewayv1alpha2.SchemeGroupVersion.Group,
					Version:  gatewayv1alpha2.SchemeGroupVersion.Version,
					Resource: "httproutes",
				},
			}.CRDExists,
			Controller: &gateway.HTTPRouteReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("HTTPRoute"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
		{
			Enabled: featureGates[gatewayFeature],
			AutoHandler: crdExistsChecker{
				GVR: schema.GroupVersionResource{
					Group:    gatewayv1alpha2.SchemeGroupVersion.Group,
					Version:  gatewayv1alpha2.SchemeGroupVersion.Version,
					Resource: "udproutes",
				},
			}.CRDExists,
			Controller: &gateway.UDPRouteReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("UDPRoute"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
		{
			Enabled: featureGates[gatewayFeature],
			AutoHandler: crdExistsChecker{
				GVR: schema.GroupVersionResource{
					Group:    gatewayv1alpha2.SchemeGroupVersion.Group,
					Version:  gatewayv1alpha2.SchemeGroupVersion.Version,
					Resource: "tcproutes",
				},
			}.CRDExists,
			Controller: &gateway.TCPRouteReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("TCPRoute"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
		{
			Enabled: featureGates[gatewayFeature],
			AutoHandler: crdExistsChecker{
				GVR: schema.GroupVersionResource{
					Group:    gatewayv1alpha2.SchemeGroupVersion.Group,
					Version:  gatewayv1alpha2.SchemeGroupVersion.Version,
					Resource: "tlsroutes",
				},
			}.CRDExists,
			Controller: &gateway.TLSRouteReconciler{
				Client:          mgr.GetClient(),
				Log:             ctrl.Log.WithName("controllers").WithName("TLSRoute"),
				Scheme:          mgr.GetScheme(),
				DataplaneClient: dataplaneClient,
			},
		},
	}

	return controllers, nil
}

// crdExistsChecker verifies whether the resource type defined by GVR is supported by the k8s apiserver.
type crdExistsChecker struct {
	GVR schema.GroupVersionResource
}

// CRDExists returns true iff the apiserver supports the specified group/version/resource.
func (c crdExistsChecker) CRDExists(r client.Client) bool {
	return ctrlutils.CRDExists(r, c.GVR)
}

// ingressControllerStrategy picks the best Ingress API supported by k8s apiserver.
type ingressControllerStrategy struct {
	chosenVersion IngressAPI
}

// Initialize negotiates the best Ingress API version supported by both KIC and the k8s apiserver.
func (s *ingressControllerStrategy) Initialize(cfg *Config, cl client.Client) error {
	var err error
	s.chosenVersion, err = negotiateIngressAPI(cfg, cl)
	return err
}

// IsExtV1beta1 returns true iff the best supported API version is extensions/v1beta1.
func (s *ingressControllerStrategy) IsExtV1beta1(_ client.Client) bool {
	return s.chosenVersion == ExtensionsV1beta1
}

// IsExtV1beta1 returns true iff the best supported API version is networking.k8s.io/v1beta1.
func (s *ingressControllerStrategy) IsNetV1beta1(_ client.Client) bool {
	return s.chosenVersion == NetworkingV1beta1
}

// IsExtV1beta1 returns true iff the best supported API version is networking.k8s.io/v1.
func (s *ingressControllerStrategy) IsNetV1(_ client.Client) bool {
	return s.chosenVersion == NetworkingV1
}
