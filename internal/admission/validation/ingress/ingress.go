package ingress

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation/kongplugin"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
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
	storer store.Storer,
	managerClient client.Client,
) (bool, string, error) {
	var (
		errMsgs           []string
		failuresCollector = failures.NewResourceFailuresCollector(logger)
	)

	if err := validation.ValidateRouteSourceAnnotations(ingress); err != nil {
		return false, fmt.Sprintf("Ingress has invalid Kong annotations: %s", err), nil
	}

	if err := kongplugin.ValidatePluginUniquenessPerObject(ctx, managerClient, ingress); err != nil {
		return false, fmt.Sprintf("Ingress has invalid KongPlugin annotation: %s", err), nil
	}

	for _, kg := range ingressToKongRoutesForValidation(translatorFeatures, ingress, failuresCollector, storer) {
		kg := kg
		// Validate by using feature of Kong Gateway.
		ok, msg, err := routesValidator.Validate(ctx, &kg)
		if err != nil {
			return false, fmt.Sprintf("Unable to validate Ingress schema: %s", err.Error()), nil
		}
		if !ok {
			errMsgs = append(errMsgs, msg)
		}
	}
	// Collect failures from the translation.
	for _, failure := range failuresCollector.PopResourceFailures() {
		errMsgs = append(errMsgs, failure.Message())
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
	failuresCollector subtranslator.FailuresCollector,
	storer store.Storer,
) []kong.Route {
	kongServices := subtranslator.TranslateIngresses(
		[]*netv1.Ingress{ingress},
		kongv1alpha1.IngressClassParametersSpec{EnableLegacyRegexDetection: true},
		subtranslator.TranslateIngressFeatureFlags{
			ExpressionRoutes:  translatorFeatures.ExpressionRoutes,
			KongServiceFacade: translatorFeatures.KongServiceFacade,
		},
		&translator.ObjectsCollector{}, // It's irrelevant for validation.
		failuresCollector,
		storer,
	)

	var kongRoutes []kong.Route
	for _, svc := range kongServices {
		for _, route := range svc.Routes {
			kongRoutes = append(kongRoutes, route.Route)
		}
	}
	return kongRoutes
}
