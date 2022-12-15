package admission

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	configuration "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

// RequestHandler is an HTTP server that can validate Kong Ingress Controllers'
// Custom Resources using Kubernetes Admission Webhooks.
type RequestHandler struct {
	// Validator validates the entities that the k8s API-server asks
	// it the server to validate.
	Validator KongValidator

	Logger logrus.FieldLogger
}

// ServeHTTP parses AdmissionReview requests and responds back
// with the validation result of the entity.
func (h RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.Logger.Error("received request with empty body")
		http.Error(w, "admission review object is missing",
			http.StatusBadRequest)
		return
	}

	review := admissionv1.AdmissionReview{}
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		h.Logger.WithError(err).Error("failed to decode admission review")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := h.handleValidation(r.Context(), *review.Request)
	if err != nil {
		h.Logger.WithError(err).Error("failed to run validation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	review.Response = response

	if err := json.NewEncoder(w).Encode(&review); err != nil {
		h.Logger.WithError(err).Error("failed to encode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var (
	consumerGVResource = metav1.GroupVersionResource{
		Group:    configuration.SchemeGroupVersion.Group,
		Version:  configuration.SchemeGroupVersion.Version,
		Resource: "kongconsumers",
	}
	pluginGVResource = metav1.GroupVersionResource{
		Group:    configuration.SchemeGroupVersion.Group,
		Version:  configuration.SchemeGroupVersion.Version,
		Resource: "kongplugins",
	}
	clusterPluginGVResource = metav1.GroupVersionResource{
		Group:    configuration.SchemeGroupVersion.Group,
		Version:  configuration.SchemeGroupVersion.Version,
		Resource: "kongclusterplugins",
	}
	kongIngressGVResource = metav1.GroupVersionResource{
		Group:    configuration.SchemeGroupVersion.Group,
		Version:  configuration.SchemeGroupVersion.Version,
		Resource: "kongingresses",
	}
	secretGVResource = metav1.GroupVersionResource{
		Group:    corev1.SchemeGroupVersion.Group,
		Version:  corev1.SchemeGroupVersion.Version,
		Resource: "secrets",
	}
	gatewayGVResource = metav1.GroupVersionResource{
		Group:    gatewayv1beta1.SchemeGroupVersion.Group,
		Version:  gatewayv1beta1.SchemeGroupVersion.Version,
		Resource: "gateways",
	}
	httprouteGVResource = metav1.GroupVersionResource{
		Group:    gatewayv1beta1.SchemeGroupVersion.Group,
		Version:  gatewayv1beta1.SchemeGroupVersion.Version,
		Resource: "httproutes",
	}
)

func (h RequestHandler) handleValidation(ctx context.Context, request admissionv1.AdmissionRequest) (
	*admissionv1.AdmissionResponse, error,
) {
	responseBuilder := NewResponseBuilder(request.UID)

	switch request.Resource {
	case consumerGVResource:
		return h.handleKongConsumer(ctx, request, responseBuilder)
	case pluginGVResource:
		return h.handleKongPlugin(ctx, request, responseBuilder)
	case clusterPluginGVResource:
		return h.handleKongClusterPlugin(ctx, request, responseBuilder)
	case secretGVResource:
		return h.handleSecret(ctx, request, responseBuilder)
	case gatewayGVResource:
		return h.handleGateway(ctx, request, responseBuilder)
	case httprouteGVResource:
		return h.handleHTTPRoute(ctx, request, responseBuilder)
	case kongIngressGVResource:
		return h.handleKongIngress(ctx, request, responseBuilder)
	default:
		return nil, fmt.Errorf("unknown resource type to validate: %s/%s %s",
			request.Resource.Group, request.Resource.Version,
			request.Resource.Resource)
	}
}

func (h RequestHandler) handleKongConsumer(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	consumer := configuration.KongConsumer{}
	deserializer := codecs.UniversalDeserializer()
	_, _, err := deserializer.Decode(request.Object.Raw, nil, &consumer)
	if err != nil {
		return nil, err
	}
	//nolint:exhaustive
	switch request.Operation {
	case admissionv1.Create:
		ok, msg, err := h.Validator.ValidateConsumer(ctx, consumer)
		if err != nil {
			return nil, err
		}
		return responseBuilder.Allowed(ok).WithMessage(msg).Build(), nil
	case admissionv1.Update:
		var oldConsumer configuration.KongConsumer
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

func (h RequestHandler) handleKongPlugin(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	plugin := configuration.KongPlugin{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &plugin)
	if err != nil {
		return nil, err
	}

	ok, message, err := h.Validator.ValidatePlugin(ctx, plugin)
	if err != nil {
		return nil, err
	}

	return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
}

func (h RequestHandler) handleKongClusterPlugin(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	plugin := configuration.KongClusterPlugin{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &plugin)
	if err != nil {
		return nil, err
	}

	ok, message, err := h.Validator.ValidateClusterPlugin(ctx, plugin)
	if err != nil {
		return nil, err
	}

	return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
}

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
	if _, ok := secret.Data["kongCredType"]; !ok {
		// secret does not look like a credential resource in Kong
		return responseBuilder.Allowed(true).Build(), nil
	}

	// secrets are only validated on update because they must be referenced by a
	// managed consumer in order for us to validate them, and because credentials
	// validation also happens at the consumer side of the reference so a
	// credentials secret can not be referenced without being validated.
	//nolint:exhaustive
	switch request.Operation {
	case admissionv1.Update:
		ok, message, err := h.Validator.ValidateCredential(ctx, secret)
		if err != nil {
			return nil, err
		}
		return responseBuilder.Allowed(ok).WithMessage(message).Build(), nil
	default:
		return nil, fmt.Errorf("unknown operation %q", string(request.Operation))
	}
}

func (h RequestHandler) handleGateway(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	gateway := gatewayv1beta1.Gateway{}
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

func (h RequestHandler) handleHTTPRoute(
	ctx context.Context,
	request admissionv1.AdmissionRequest,
	responseBuilder *ResponseBuilder,
) (*admissionv1.AdmissionResponse, error) {
	httproute := gatewayv1beta1.HTTPRoute{}
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

func (h RequestHandler) handleKongIngress(_ context.Context, request admissionv1.AdmissionRequest, responseBuilder *ResponseBuilder) (*admissionv1.AdmissionResponse, error) {
	kongIngress := configuration.KongIngress{}
	_, _, err := codecs.UniversalDeserializer().Decode(request.Object.Raw, nil, &kongIngress)
	if err != nil {
		return nil, err
	}

	// KongIngress is always allowed.
	responseBuilder = responseBuilder.Allowed(true)

	if kongIngress.Proxy != nil {
		const warning = "'proxy' is DEPRECATED. Use Service's annotations instead."
		responseBuilder = responseBuilder.WithWarning(warning)
	}

	if kongIngress.Route != nil {
		const warning = "'route' is DEPRECATED. Use Ingress' annotations instead."
		responseBuilder = responseBuilder.WithWarning(warning)
	}

	return responseBuilder.Build(), nil
}
