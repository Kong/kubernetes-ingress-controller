package admission

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

const (
	KindKongPlugin        = "KongPlugin"
	KindKongClusterPlugin = "KongClusterPlugin"
)

// RequestHandler is an HTTP server that can validate Kong Ingress Controllers'
// Custom Resources using Kubernetes Admission Webhooks.
type RequestHandler struct {
	// Validator validates the entities that the k8s API-server asks
	// it the server to validate.
	Validator KongValidator
	// ReferenceIndexers gets the resources (KongPlugin and KongClusterPlugin)
	// referring the validated resource (Secret) to check the changes on
	// referred Secret will produce invalid configuration of the plugins.
	ReferenceIndexers ctrlref.CacheIndexers

	Logger logr.Logger
}

// ServeHTTP parses AdmissionReview requests and responds back
// with the validation result of the entity.
func (h RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.Logger.Error(nil, "Received request with empty body")
		http.Error(w, "Admission review object is missing",
			http.StatusBadRequest)
		return
	}

	review := admissionv1.AdmissionReview{}
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		h.Logger.Error(err, "Failed to decode admission review")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := h.handleValidation(r.Context(), *review.Request)
	if err != nil {
		h.Logger.Error(err, "Failed to run validation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	review.Response = response

	if err := json.NewEncoder(w).Encode(&review); err != nil {
		h.Logger.Error(err, "Failed to encode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var (
	consumerGVResource = metav1.GroupVersionResource{
		Group:    kongv1.SchemeGroupVersion.Group,
		Version:  kongv1.SchemeGroupVersion.Version,
		Resource: "kongconsumers",
	}
	consumerGroupGVResource = metav1.GroupVersionResource{
		Group:    kongv1beta1.SchemeGroupVersion.Group,
		Version:  kongv1beta1.SchemeGroupVersion.Version,
		Resource: "kongconsumergroups",
	}
	pluginGVResource = metav1.GroupVersionResource{
		Group:    kongv1.SchemeGroupVersion.Group,
		Version:  kongv1.SchemeGroupVersion.Version,
		Resource: "kongplugins",
	}
	clusterPluginGVResource = metav1.GroupVersionResource{
		Group:    kongv1.SchemeGroupVersion.Group,
		Version:  kongv1.SchemeGroupVersion.Version,
		Resource: "kongclusterplugins",
	}
	kongIngressGVResource = metav1.GroupVersionResource{
		Group:    kongv1.SchemeGroupVersion.Group,
		Version:  kongv1.SchemeGroupVersion.Version,
		Resource: "kongingresses",
	}
	kongVaultGVResource = metav1.GroupVersionResource{
		Group:    kongv1alpha1.SchemeGroupVersion.Group,
		Version:  kongv1alpha1.SchemeGroupVersion.Version,
		Resource: "kongvaults",
	}
	secretGVResource = metav1.GroupVersionResource{
		Group:    corev1.SchemeGroupVersion.Group,
		Version:  corev1.SchemeGroupVersion.Version,
		Resource: "secrets",
	}
	ingressGVResource = metav1.GroupVersionResource{
		Group:    netv1.SchemeGroupVersion.Group,
		Version:  netv1.SchemeGroupVersion.Version,
		Resource: "ingresses",
	}
	serviceGVResource = metav1.GroupVersionResource{
		Group:    corev1.SchemeGroupVersion.Group,
		Version:  corev1.SchemeGroupVersion.Version,
		Resource: "services",
	}
)

func (h RequestHandler) handleValidation(ctx context.Context, request admissionv1.AdmissionRequest) (
	*admissionv1.AdmissionResponse, error,
) {
	responseBuilder := NewResponseBuilder(request.UID)

	switch request.Resource {
	case consumerGVResource:
		return h.handleKongConsumer(ctx, request, responseBuilder)
	case consumerGroupGVResource:
		return h.handleKongConsumerGroup(ctx, request, responseBuilder)
	case pluginGVResource:
		return h.handleKongPlugin(ctx, request, responseBuilder)
	case clusterPluginGVResource:
		return h.handleKongClusterPlugin(ctx, request, responseBuilder)
	case secretGVResource:
		return h.handleSecret(ctx, request, responseBuilder)
	case gatewayapi.V1GatewayGVResource, gatewayapi.V1beta1GatewayGVResource:
		return h.handleGateway(ctx, request, responseBuilder)
	case gatewayapi.V1HTTPRouteGVResource, gatewayapi.V1beta1HTTPRouteGVResource:
		return h.handleHTTPRoute(ctx, request, responseBuilder)
	case kongIngressGVResource:
		return h.handleKongIngress(ctx, request, responseBuilder)
	case kongVaultGVResource:
		return h.handleKongVault(ctx, request, responseBuilder)
	case serviceGVResource:
		return h.handleService(ctx, request, responseBuilder)
	case ingressGVResource:
		return h.handleIngress(ctx, request, responseBuilder)
	default:
		return nil, fmt.Errorf("unknown resource type to validate: %s/%s %s",
			request.Resource.Group, request.Resource.Version,
			request.Resource.Resource)
	}
}

// +kubebuilder:webhook:verbs=create;update,groups=configuration.konghq.com,resources=kongconsumers,versions=v1,name=kongconsumers.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleKongConsumer(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	consumer := kongv1.KongConsumer{}
	deserializer := codecs.UniversalDeserializer()
	_, _, err := deserializer.Decode(request.Object.Raw, nil, &consumer)
	if err != nil {
		return nil, err
	}

	switch request.Operation {
	case admissionv1.Create:
		ok, msg, err := h.Validator.ValidateConsumer(ctx, consumer)
		if err != nil {
			return nil, err
		}
		return responseBuilder.Allowed(ok).WithMessage(msg).Build(), nil
	case admissionv1.Update:
		var oldConsumer kongv1.KongConsumer
		_, _, err = deserializer.Decode(request.OldObject.Raw, nil, &oldConsumer)
		if err != nil {
			return nil, err
		}
		// validate only if the username is being changed
		if consumer.Username == oldConsumer.Username {
			return responseBuilder.Allowed(true).Build(), nil
		}
		ok, message, err := h.Validator.ValidateConsumer(ctx, consumer)
		if err != nil {
			return nil, err
		}
		return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
	default:
		return nil, fmt.Errorf("unknown operation %q", string(request.Operation))
	}
}

// +kubebuilder:webhook:verbs=create;update,groups=configuration.konghq.com,resources=kongconsumergroups,versions=v1beta1,name=kongconsumergroups.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleKongConsumerGroup(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	var consumerGroup kongv1beta1.KongConsumerGroup
	if _, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &consumerGroup); err != nil {
		return nil, err
	}
	ok, message, err := h.Validator.ValidateConsumerGroup(ctx, consumerGroup)
	if err != nil {
		return nil, err
	}

	return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
}

// +kubebuilder:webhook:verbs=create;update,groups=configuration.konghq.com,resources=kongplugins,versions=v1,name=kongplugins.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleKongPlugin(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	plugin := kongv1.KongPlugin{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &plugin)
	if err != nil {
		return nil, err
	}

	ok, message, err := h.Validator.ValidatePlugin(ctx, plugin, nil)
	if err != nil {
		return nil, err
	}

	return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
}

// +kubebuilder:webhook:verbs=create;update,groups=configuration.konghq.com,resources=kongclusterplugins,versions=v1,name=kongclusterplugins.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleKongClusterPlugin(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	plugin := kongv1.KongClusterPlugin{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &plugin)
	if err != nil {
		return nil, err
	}

	ok, message, err := h.Validator.ValidateClusterPlugin(ctx, plugin, nil)
	if err != nil {
		return nil, err
	}

	return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
}

// NOTE this handler _does not_ use a kubebuilder directive, as our Secret handling requires webhook features
// kubebuilder does not support (objectSelector). Instead, config/webhook/additional_secret_hooks.yaml includes
// handwritten webhook rules that we patch into the webhook manifest.
// See https://github.com/kubernetes-sigs/controller-tools/issues/553 for further info.

func (h RequestHandler) handleSecret(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	secret := corev1.Secret{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &secret)
	if err != nil {
		return nil, err
	}

	switch request.Operation {
	case admissionv1.Update, admissionv1.Create:
		// credential secrets
		if credType, err := util.ExtractKongCredentialType(&secret); err == nil && credType != "" {
			ok, message := h.Validator.ValidateCredential(ctx, secret)
			if !ok {
				return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
			}
		}

		// TODO this check duplicates the objectSelector filter in the webhook definition, and will only check updates to
		// plugin configuration secrets if they have the new 3.2+ plugin configuration label. we could optionally remove
		// this check to allow users to remove the secret filter configuration from the webhook definition and check any
		// referenced secret, labeled or not.

		// plugin configuration secrets
		switch validate := secret.Labels[labels.ValidateLabel]; labels.ValidateType(validate) {
		case labels.PluginValidate:
			ok, message, err := h.checkReferrersOfSecret(ctx, &secret)
			if err != nil {
				return responseBuilder.Allowed(false).WithMessage(fmt.Sprintf("failed to validate other objects referencing the secret: %v", err)).Build(), err
			}
			if !ok {
				return responseBuilder.Allowed(false).WithMessage(message).Build(), nil
			}
		default:
			// TODO this duplicates the above plugin handling block. prior to 3.2, the admission webhook ingested all
			// Secrets and used this to validate updates to plugin configuration. this non-labeled case is retained
			// for environments that still use ingest all configuration.
			ok, message, err := h.checkReferrersOfSecret(ctx, &secret)
			if err != nil {
				return responseBuilder.Allowed(false).WithMessage(fmt.Sprintf("failed to validate other objects referencing the secret: %v", err)).Build(), err
			}
			if !ok {
				return responseBuilder.Allowed(false).WithMessage(message).Build(), nil
			}

			// no reference found in the blanket block, this is some random unrelated Secret and KIC should ignore it.
			return responseBuilder.Allowed(true).Build(), nil
		}

	default:
		return nil, fmt.Errorf("unknown operation %q", string(request.Operation))
	}
	// fallback allow. it should not be possible to hit this because of the defaults above, but the compiler wants it.
	// if a request somehow has reached this, we shouldn't touch it.
	return responseBuilder.Allowed(true).Build(), nil
}

// checkReferrersOfSecret validates all referrers (KongPlugins and KongClusterPlugins) of the secret
// and rejects the secret if it generates invalid configurations for any of the referrers.
func (h RequestHandler) checkReferrersOfSecret(ctx context.Context, secret *corev1.Secret) (bool, string, error) {
	referrers, err := h.ReferenceIndexers.ListReferrerObjectsByReferent(secret)
	if err != nil {
		return false, "", fmt.Errorf("failed to list referrers of secret: %w", err)
	}

	for _, obj := range referrers {
		gvk := obj.GetObjectKind().GroupVersionKind()
		if gvk.Group == kongv1.GroupVersion.Group && gvk.Version == kongv1.GroupVersion.Version && gvk.Kind == KindKongPlugin {
			plugin := obj.(*kongv1.KongPlugin)
			ok, message, err := h.Validator.ValidatePlugin(ctx, *plugin, []*corev1.Secret{secret})
			if err != nil {
				return false, "", fmt.Errorf("failed to run validation on KongPlugin %s/%s: %w",
					plugin.Namespace, plugin.Name, err,
				)
			}
			if !ok {
				return false,
					fmt.Sprintf("Change on secret will generate invalid configuration for KongPlugin %s/%s: %s",
						plugin.Namespace, plugin.Name, message,
					), nil
			}
		}
		if gvk.Group == kongv1.GroupVersion.Group && gvk.Version == kongv1.GroupVersion.Version && gvk.Kind == KindKongClusterPlugin {
			plugin := obj.(*kongv1.KongClusterPlugin)
			ok, message, err := h.Validator.ValidateClusterPlugin(ctx, *plugin, []*corev1.Secret{secret})
			if err != nil {
				return false, "", fmt.Errorf("failed to run validation on KongClusterPlugin %s: %w",
					plugin.Name, err,
				)
			}
			if !ok {
				return false, fmt.Sprintf("Change on secret will generate invalid configuration for KongClusterPlugin %s: %s",
					plugin.Name, message,
				), nil
			}
		}
	}
	return true, "", nil
}

// +kubebuilder:webhook:verbs=create;update,groups=gateway.networking.k8s.io,resources=gateways,versions=v1;v1beta1,name=gateways.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleGateway(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	gateway := gatewayapi.Gateway{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &gateway)
	if err != nil {
		return nil, err
	}
	ok, message, err := h.Validator.ValidateGateway(ctx, gateway)
	if err != nil {
		return nil, err
	}

	return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
}

// +kubebuilder:webhook:verbs=create;update,groups=gateway.networking.k8s.io,resources=httproutes,versions=v1;v1beta1,name=httproutes.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleHTTPRoute(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	httproute := gatewayapi.HTTPRoute{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &httproute)
	if err != nil {
		return nil, err
	}
	ok, message, err := h.Validator.ValidateHTTPRoute(ctx, httproute)
	if err != nil {
		return nil, err
	}
	return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
}

const (
	proxyWarning    = "Support for 'proxy' was removed in 3.0. It will have no effect. Use Service's annotations instead."
	routeWarning    = "Support for 'route' was removed in 3.0. It will have no effect. Use Ingress' annotations instead."
	upstreamWarning = "'upstream' is DEPRECATED and will be removed in a future version. Use a KongUpstreamPolicy resource instead."
)

// +kubebuilder:webhook:verbs=create;update,groups=configuration.konghq.com,resources=kongingresses,versions=v1,name=kongingresses.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleKongIngress(_ context.Context, request admissionv1.AdmissionRequest, responseBuilder *ResponseBuilder) (*admissionv1.AdmissionResponse, error) {
	kongIngress := kongv1.KongIngress{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &kongIngress)
	if err != nil {
		return nil, err
	}

	// KongIngress is always allowed.
	responseBuilder = responseBuilder.Allowed(true)

	// Proxy and Route fields are now disallowed to be set with the use of CEL rules in the CRD definition.
	// We still warn about them here only just in case someone doesn't install new CRDs with CEL rules.
	if kongIngress.Proxy != nil {
		responseBuilder = responseBuilder.WithWarning(proxyWarning)
	}
	if kongIngress.Route != nil {
		responseBuilder = responseBuilder.WithWarning(routeWarning)
	}

	if kongIngress.Upstream != nil {
		responseBuilder = responseBuilder.WithWarning(upstreamWarning)
	}

	return responseBuilder.Build(), nil
}

const (
	serviceWarning = "%s is deprecated and will be removed in a future release. Use Service annotations " +
		"for the 'proxy' section and %s with a KongUpstreamPolicy resource instead."
)

// +kubebuilder:webhook:verbs=create;update,groups=core,resources=services,versions=v1,name=services.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleService(_ context.Context, request admissionv1.AdmissionRequest, responseBuilder *ResponseBuilder) (*admissionv1.AdmissionResponse, error) {
	service := corev1.Service{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &service)
	if err != nil {
		return nil, err
	}

	// Service is always allowed.
	responseBuilder = responseBuilder.Allowed(true)

	if annotations.ExtractConfigurationName(service.Annotations) != "" {
		warning := fmt.Sprintf(serviceWarning, annotations.AnnotationPrefix+annotations.ConfigurationKey,
			kongv1beta1.KongUpstreamPolicyAnnotationKey)

		responseBuilder = responseBuilder.WithWarning(warning)
	}

	return responseBuilder.Build(), nil
}

// +kubebuilder:webhook:verbs=create;update,groups=networking.k8s.io,resources=ingresses,versions=v1,name=ingresses.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleIngress(ctx context.Context, request admissionv1.AdmissionRequest, responseBuilder *ResponseBuilder) (*admissionv1.AdmissionResponse, error) {
	ingress := netv1.Ingress{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &ingress)
	if err != nil {
		return nil, err
	}
	ok, message, err := h.Validator.ValidateIngress(ctx, ingress)
	if err != nil {
		return nil, err
	}

	return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
}

// +kubebuilder:webhook:verbs=create;update,groups=configuration.konghq.com,resources=kongvaults,versions=v1alpha1,name=kongvaults.validation.ingress-controller.konghq.com,path=/,webhookVersions=v1,matchPolicy=equivalent,mutating=false,failurePolicy=fail,sideEffects=None,admissionReviewVersions=v1

func (h RequestHandler) handleKongVault(ctx context.Context, request admissionv1.AdmissionRequest, responseBuilder *ResponseBuilder) (*admissionv1.AdmissionResponse, error) {
	kongVault := kongv1alpha1.KongVault{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &kongVault)
	if err != nil {
		return nil, err
	}
	ok, message, err := h.Validator.ValidateVault(ctx, kongVault)
	if err != nil {
		return nil, err
	}

	return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
}
