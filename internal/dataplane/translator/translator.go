package translator

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
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

	// KongCustomEntity indicates whether we should support translating custom entities from KongCustomEntity CRs.
	KongCustomEntity bool
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
		KongCustomEntity:                  featureGates.Enabled(featuregates.KongCustomEntity),
	}
}

// SchemaServiceProvider returns a kong schema service required for translating custom entities.
type SchemaServiceProvider interface {
	GetSchemaService() kong.AbstractSchemaService
}

// Translator translates Kubernetes objects and configurations into their
// equivalent Kong objects and configurations, producing a complete
// state configuration for the Kong Admin API.
type Translator struct {
	logger        logr.Logger
	storer        store.Storer
	workspace     string
	licenseGetter license.Getter
	featureFlags  FeatureFlags

	// schemaServiceProvider provides the schema service required for fetching schemas of custom entities.
	schemaServiceProvider SchemaServiceProvider

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
	schemaServiceProvider SchemaServiceProvider,
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
		schemaServiceProvider:      schemaServiceProvider,
		failuresCollector:          failuresCollector,
		translatedObjectsCollector: translatedObjectsCollector,
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

// UpdateCache updates the store cache used by the translator.
// This method can be used to swap the cache with another one (e.g. the last valid snapshot).
func (t *Translator) UpdateCache(c store.CacheStores) {
	t.storer.UpdateCache(c)
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

	// process custom entities
	if t.featureFlags.KongCustomEntity {
		result.FillCustomEntities(t.logger, t.storer, t.failuresCollector, t.schemaServiceProvider.GetSchemaService(), t.workspace)
		// Register successcully translated KCEs to set the status of these KCEs.
		for _, collection := range result.CustomEntities {
			for i := range collection.Entities {
				t.registerSuccessfullyTranslatedObject(collection.Entities[i].K8sKongCustomEntity)
			}
		}
	}

	// generate Certificates and SNIs
	ingressCerts := t.getCerts(ingressRules.SecretNameToSNIs)
	gatewayCerts := t.getGatewayCerts()
	// note that ingress-derived certificates will take precedence over gateway-derived certificates for SNI assignment
	var certIDsSeen certIDToMergedCertID
	result.Certificates, certIDsSeen = mergeCerts(t.logger, ingressCerts, gatewayCerts)

	// re-fill client certificate IDs of services after certificates are merged.
	for i, s := range result.Services {
		if s.ClientCertificate != nil && s.ClientCertificate.ID != nil {
			certID := s.ClientCertificate.ID
			mergedCertID := certIDsSeen[*certID]
			result.Services[i].ClientCertificate = &kong.Certificate{
				ID: kong.String(mergedCertID),
			}
		}
	}

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

// -----------------------------------------------------------------------------
// Translator - Public Methods - Other Optional Features
// -----------------------------------------------------------------------------

// InjectLicenseGetter sets a license getter to be used by the translator.
func (t *Translator) InjectLicenseGetter(licenseGetter license.Getter) {
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

// UnavailableSchemaService is a fake schema service used when no gateway admin API clients available.
// It always returns error in its Get and Validate methods.
type UnavailableSchemaService struct{}

var _ kong.AbstractSchemaService = UnavailableSchemaService{}

func (s UnavailableSchemaService) Get(_ context.Context, _ string) (kong.Schema, error) {
	return nil, errors.New("schema service unavailable")
}

func (s UnavailableSchemaService) Validate(_ context.Context, _ kong.EntityType, _ any) (bool, string, error) {
	return false, "", errors.New("schema service unavailable")
}
