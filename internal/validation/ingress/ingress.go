package ingress

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	netv1 "k8s.io/api/networking/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
)

type routeValidator interface {
	Validate(context.Context, *kong.Route) (bool, string, error)
}

func ValidateIngress(
	ctx context.Context,
	routesValidator routeValidator,
	parserFeatures parser.FeatureFlags,
	kongVersion semver.Version,
	ingress *netv1.Ingress,
) (bool, string, error) {
	discardLogger := logrus.New()
	discardLogger.Out = io.Discard
	failuresCollector, err := failures.NewResourceFailuresCollector(discardLogger)
	if err != nil {
		return false, "", err
	}
	var icp kongv1alpha1.IngressClassParametersSpec
	if kongVersion.LT(versions.ExplicitRegexPathVersionCutoff) {
		icp.EnableLegacyRegexDetection = true
	}
	result := parser.IngressesV1ToKongServices(
		parserFeatures,
		[]*netv1.Ingress{ingress},
		icp,
		&parser.ObjectsCollector{}, // It's irrelevant for validation.
		failuresCollector,
	)

	var errMsgs []string
	for _, f := range failuresCollector.PopResourceFailures() {
		errMsgs = append(errMsgs, f.Message())
	}
	if len(errMsgs) > 0 {
		return false, validationMsg(errMsgs), nil
	}
	var kongRoutes []kong.Route
	for _, r := range result {
		for _, route := range r.Routes {
			kongRoutes = append(kongRoutes, route.Route)
		}
	}
	// Validate by using feature of Kong Gateway.
	for _, kg := range kongRoutes {
		kg := kg
		ok, msg, err := routesValidator.Validate(ctx, &kg)
		if err != nil {
			return false, fmt.Sprintf("unable to validate Ingress schema: %s", err.Error()), nil
		}
		if !ok {
			errMsgs = append(errMsgs, msg)
		}
	}
	if len(errMsgs) > 0 {
		return false, validationMsg(errMsgs), nil
	}
	return true, "", nil
}

func validationMsg(errMsgs []string) string {
	return fmt.Sprintf("Ingress failed schema validation: %s", strings.Join(errMsgs, ", "))
}
