package translator

import (
	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/mo"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// -----------------------------------------------------------------------------
// Translator - Public Constants and Package Variables
// -----------------------------------------------------------------------------

const KindGateway = gatewayapi.Kind("Gateway")

// -----------------------------------------------------------------------------
// Translator - Public Types
// -----------------------------------------------------------------------------

// FeatureFlags are used to control the behavior of the translator.
type FeatureFlags struct {
	// ReportConfiguredKubernetesObjects turns on object reporting for this translator:
	// each subsequent call to BuildKongConfig() will track the Kubernetes objects which
	// were successfully translated.
	ReportConfiguredKubernetesObjects bool

	// ExpressionRoutes indicates whether to translate Kubernetes objects to expression based Kong Routes.
	ExpressionRoutes bool

	// EnterpriseEdition indicates whether to translate objects that are only available in the Kong enterprise edition.
	EnterpriseEdition bool

	// FillIDs enables the translator to fill in the IDs fields of Kong entities - Services, Routes, and Consumers - based
	// on their names. It ensures that IDs remain stable across restarts of the controller.
	FillIDs bool

	// RewriteURIs enables the translator to translate the konghq.com/rewrite annotation to the proper set of Kong plugins.
	RewriteURIs bool

	// KongServiceFacade indicates whether we should support KongServiceFacades as Ingress backends.
	KongServiceFacade bool
}

func NewFeatureFlags(
	featureGates featuregates.FeatureGates,
	routerFlavor dpconf.RouterFlavor,
	updateStatusFlag bool,
	enterpriseEdition bool,
) FeatureFlags {
	return FeatureFlags{
		ReportConfiguredKubernetesObjects: updateStatusFlag,
		ExpressionRoutes:                  dpconf.ShouldEnableExpressionRoutes(routerFlavor),
		EnterpriseEdition:                 enterpriseEdition,
		FillIDs:                           featureGates.Enabled(featuregates.FillIDsFeature),
		RewriteURIs:                       featureGates.Enabled(featuregates.RewriteURIsFeature),
		KongServiceFacade:                 featureGates.Enabled(featuregates.KongServiceFacade),
	}
}

// LicenseGetter is an interface for getting the Kong Enterprise license.
type LicenseGetter interface {
	// GetLicense returns an optional license.
	GetLicense() mo.Option[kong.License]
}

// Translator translates Kubernetes objects and configurations into their
// equivalent Kong objects and configurations, producing a complete
// state configuration for the Kong Admin API.
type Translator struct {
	logger           logr.Logger
	storer           store.Storer
	workspace        string
	licenseGetter    LicenseGetter
	featureFlags     FeatureFlags
	ingressClassName string

	failuresCollector          *failures.ResourceFailuresCollector
	translatedObjectsCollector *ObjectsCollector
}

// NewTranslator produces a new Translator object provided a logging mechanism
// and a Kubernetes object store.
func NewTranslator(
	logger logr.Logger,
	storer store.Storer,
	workspace string,
	featureFlags FeatureFlags,
	ingressClassName string,
) (*Translator, error) {
	failuresCollector := failures.NewResourceFailuresCollector(logger)

	// If the feature flag is enabled, create a new collector for translated objects.
	var translatedObjectsCollector *ObjectsCollector
	if featureFlags.ReportConfiguredKubernetesObjects {
		translatedObjectsCollector = NewObjectsCollector()
	}

	return &Translator{
		logger:                     logger,
		storer:                     storer,
		workspace:                  workspace,
		featureFlags:               featureFlags,
		failuresCollector:          failuresCollector,
		translatedObjectsCollector: translatedObjectsCollector,
		ingressClassName:           ingressClassName,
	}, nil
}

// -----------------------------------------------------------------------------
// Translator - Public Methods
// -----------------------------------------------------------------------------

// KongConfigBuildingResult is a result of Translator.BuildKongConfig method.
type KongConfigBuildingResult struct {
	// KongState is the Kong configuration used to configure the Gateway(s).
	KongState *kongstate.KongState

	// TranslationFailures is a list of resource failures that occurred during parsing.
	// They should be used to provide users with feedback on Kubernetes objects validity.
	TranslationFailures []failures.ResourceFailure

	// ConfiguredKubernetesObjects is a list of Kubernetes objects that were successfully translated.
	ConfiguredKubernetesObjects []client.Object
}

func (t *Translator) UpdateStore(s store.Storer) {
	t.storer = s
}

// BuildKongConfig creates a Kong configuration from Ingress and Custom resources
// defined in Kubernetes.
func (t *Translator) BuildKongConfig() KongConfigBuildingResult {
	// Translate and merge all rules together from all Kubernetes API sources
	ingressRules := mergeIngressRules(
		t.ingressRulesFromIngressV1(),
		t.ingressRulesFromTCPIngressV1beta1(),
		t.ingressRulesFromUDPIngressV1beta1(),
		t.ingressRulesFromHTTPRoutes(),
		t.ingressRulesFromUDPRoutes(),
		t.ingressRulesFromTCPRoutes(),
		t.ingressRulesFromTLSRoutes(),
		t.ingressRulesFromGRPCRoutes(),
	)

	// populate any Kubernetes Service objects relevant objects and get the
	// services to be skipped because of annotations inconsistency
	servicesToBeSkipped := ingressRules.populateServices(t.logger, t.storer, t.failuresCollector, t.translatedObjectsCollector)

	// add the routes and services to the state
	var result kongstate.KongState

	// generate Upstreams and Targets from service defs
	// update ServiceNameToServices with resolved ports (translating any name references to their number, as Kong
	// services require a number)
	result.Upstreams, ingressRules.ServiceNameToServices = t.getUpstreams(ingressRules.ServiceNameToServices)

	for key, service := range ingressRules.ServiceNameToServices {
		// if the service doesn't need to be skipped, then add it to the
		// list of services.
		if _, ok := servicesToBeSkipped[key]; !ok {
			result.Services = append(result.Services, service)
		}
	}

	// merge KongIngress with Routes, Services and Upstream
	result.FillOverrides(t.logger, t.storer, t.failuresCollector)

	// generate consumers and credentials
	result.FillConsumersAndCredentials(t.logger, t.storer, t.failuresCollector)
	for i := range result.Consumers {
		t.registerSuccessfullyTranslatedObject(&result.Consumers[i].K8sKongConsumer)
	}

	// generate vaults
	result.FillVaults(t.logger, t.storer, t.failuresCollector)
	for i := range result.Vaults {
		t.registerSuccessfullyTranslatedObject(result.Vaults[i].K8sKongVault)
	}

	// process consumer groups
	result.FillConsumerGroups(t.logger, t.storer)
	for i := range result.ConsumerGroups {
		t.registerSuccessfullyTranslatedObject(&result.ConsumerGroups[i].K8sKongConsumerGroup)
	}

	// process annotation plugins
	result.FillPlugins(t.logger, t.storer, t.failuresCollector)
	for i := range result.Plugins {
		t.registerSuccessfullyTranslatedObject(result.Plugins[i].K8sParent)
	}

	// generate Certificates and SNIs
	ingressCerts := t.getCerts(ingressRules.SecretNameToSNIs)
	gatewayCerts := t.getGatewayCerts()
	// note that ingress-derived certificates will take precedence over gateway-derived certificates for SNI assignment
	result.Certificates = mergeCerts(t.logger, ingressCerts, gatewayCerts)

	// populate CA certificates in Kong
	result.CACertificates = t.getCACerts()

	if t.licenseGetter != nil && t.featureFlags.EnterpriseEdition {
		optionalLicense := t.licenseGetter.GetLicense()
		if l, ok := optionalLicense.Get(); ok {
			result.Licenses = append(result.Licenses, kongstate.License{License: l})
		}
	}

	if t.featureFlags.FillIDs {
		// generate IDs for Kong entities
		result.FillIDs(t.logger, t.workspace)
	}

	return KongConfigBuildingResult{
		KongState:                   &result,
		TranslationFailures:         t.popTranslationFailures(),
		ConfiguredKubernetesObjects: t.popConfiguredKubernetesObjects(),
	}
}

func (t *Translator) IngressClassName() string {
	return t.ingressClassName
}

// -----------------------------------------------------------------------------
// Translator - Public Methods - Other Optional Features
// -----------------------------------------------------------------------------

// InjectLicenseGetter sets a license getter to be used by the translator.
func (t *Translator) InjectLicenseGetter(licenseGetter LicenseGetter) {
	t.licenseGetter = licenseGetter
}

// -----------------------------------------------------------------------------
// Translator - Private Methods
// -----------------------------------------------------------------------------

// registerTranslationFailure should be called when any Kubernetes object translation failure is encountered.
func (t *Translator) registerTranslationFailure(reason string, causingObjects ...client.Object) {
	t.failuresCollector.PushResourceFailure(reason, causingObjects...)
}

func (t *Translator) popTranslationFailures() []failures.ResourceFailure {
	return t.failuresCollector.PopResourceFailures()
}

// registerSuccessfullyTranslatedObject should be called when any Kubernetes object is successfully translated.
// It collects the object for reporting purposes.
func (t *Translator) registerSuccessfullyTranslatedObject(obj client.Object) {
	t.translatedObjectsCollector.Add(obj)
}

// popConfiguredKubernetesObjects provides a list of all the Kubernetes objects
// that have been successfully translated as part of BuildKongConfig() call so far.
func (t *Translator) popConfiguredKubernetesObjects() []client.Object {
	return t.translatedObjectsCollector.Pop()
}
