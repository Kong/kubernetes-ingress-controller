package ingress

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
)

type routeValidator interface {
	Validate(context.Context, *kong.Route) (bool, string, error)
}

func ValidateIngress(
	ctx context.Context,
	routesValidator routeValidator,
	translatorFeatures translator.FeatureFlags,
	ingress *netv1.Ingress,
	logger logr.Logger,
) (bool, string, error) {
	// Validate by using feature of Kong Gateway.
	var errMsgs []string
	for _, kg := range ingressToKongRoutesForValidation(translatorFeatures, ingress, logger) {
		kg := kg
		ok, msg, err := routesValidator.Validate(ctx, &kg)
		if err != nil {
			return false, fmt.Sprintf("Unable to validate Ingress schema: %s", err.Error()), nil
		}
		if !ok {
			errMsgs = append(errMsgs, msg)
		}
	}
	if len(errMsgs) > 0 {
		return false, fmt.Sprintf("Ingress failed schema validation: %s", strings.Join(errMsgs, ", ")), nil
	}
	return true, "", nil
}

// ingressToKongRoutesForValidation converts Ingress to Kong Routes that can be validated by Kong Gateway,
// discards everything else that is not needed for validation.
func ingressToKongRoutesForValidation(
	translatorFeatures translator.FeatureFlags,
	ingress *netv1.Ingress,
	logger logr.Logger,
) []kong.Route {
	kongServices := translator.IngressesV1ToKongServices(
		translatorFeatures,
		[]*netv1.Ingress{ingress},
		kongv1alpha1.IngressClassParametersSpec{EnableLegacyRegexDetection: true},
		&translator.ObjectsCollector{}, // It's irrelevant for validation.
		logger,
	)

	var kongRoutes []kong.Route
	for _, svc := range kongServices {
		for _, route := range svc.Routes {
			kongRoutes = append(kongRoutes, route.Route)
		}
	}
	return kongRoutes
}
