package admission

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	credsvalidation "github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation/consumers/credentials"
	gatewayvalidation "github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation/gateway"
	ingressvalidation "github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation/ingress"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation/kongplugin"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

// KongValidator validates Kong entities.
type KongValidator interface {
	ValidateConsumer(ctx context.Context, consumer kongv1.KongConsumer) (bool, string, error)
	ValidateConsumerGroup(ctx context.Context, consumerGroup kongv1beta1.KongConsumerGroup) (bool, string, error)
	ValidatePlugin(ctx context.Context, plugin kongv1.KongPlugin, overrideSecrets []*corev1.Secret) (bool, string, error)
	ValidateClusterPlugin(ctx context.Context, plugin kongv1.KongClusterPlugin, overrideSecrets []*corev1.Secret) (bool, string, error)
	ValidateVault(ctx context.Context, vault kongv1alpha1.KongVault) (bool, string, error)
	ValidateCredential(ctx context.Context, secret corev1.Secret) (bool, string)
	ValidateGateway(ctx context.Context, gateway gatewayapi.Gateway) (bool, string, error)
	ValidateHTTPRoute(ctx context.Context, httproute gatewayapi.HTTPRoute) (bool, string, error)
	ValidateIngress(ctx context.Context, ingress netv1.Ingress) (bool, string, error)
	ValidateService(ctx context.Context, ingress corev1.Service) (bool, string, error)
}

// AdminAPIServicesProvider provides KongHTTPValidator with Kong Admin API services that are needed to perform
// validation against entities stored by the Gateway.
type AdminAPIServicesProvider interface {
	GetConsumersService() (kong.AbstractConsumerService, bool)
	GetPluginsService() (kong.AbstractPluginService, bool)
	GetConsumerGroupsService() (kong.AbstractConsumerGroupService, bool)
	GetInfoService() (kong.AbstractInfoService, bool)
	GetRoutesService() (kong.AbstractRouteService, bool)
	GetVaultsService() (kong.AbstractVaultService, bool)
}

// ConsumerGetter is an interface for retrieving KongConsumers.
type ConsumerGetter interface {
	ListAllConsumers(ctx context.Context) ([]kongv1.KongConsumer, error)
}

// SecretGetterWithOverride returns the override secrets in the list if the namespace and name matches,
// or use the nested secretGetter to fetch the secret otherwise.
// Used for validating changes of secrets to override existing the one in cache with the one to be updated.
type SecretGetterWithOverride struct {
	overrideSecrets map[k8stypes.NamespacedName]*corev1.Secret
	secretGetter    kongstate.SecretGetter
}

var _ kongstate.SecretGetter = &SecretGetterWithOverride{}

func (s *SecretGetterWithOverride) GetSecret(namespace, name string) (*corev1.Secret, error) {
	nsName := k8stypes.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	overrideSecret, ok := s.overrideSecrets[nsName]
	if ok {
		return overrideSecret, nil
	}

	return s.secretGetter.GetSecret(namespace, name)
}

// NewSecretGetterWithOverride returns a secret getter with given override secrets.
func NewSecretGetterWithOverride(s kongstate.SecretGetter, overrideSecrets []*corev1.Secret) *SecretGetterWithOverride {
	overrideSecretMap := lo.SliceToMap(overrideSecrets, func(secret *corev1.Secret) (k8stypes.NamespacedName, *corev1.Secret) {
		return k8stypes.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, secret.DeepCopy()
	})
	return &SecretGetterWithOverride{
		overrideSecrets: overrideSecretMap,
		secretGetter:    s,
	}
}

// KongHTTPValidator implements KongValidator interface to validate Kong
// entities using the Admin API of Kong.
type KongHTTPValidator struct {
	Logger                   logr.Logger
	SecretGetter             kongstate.SecretGetter
	ConsumerGetter           ConsumerGetter
	Storer                   store.Storer
	ManagerClient            client.Client
	AdminAPIServicesProvider AdminAPIServicesProvider
	TranslatorFeatures       translator.FeatureFlags

	ingressClassMatcher   func(*metav1.ObjectMeta, string, annotations.ClassMatching) bool
	ingressV1ClassMatcher func(*netv1.Ingress, annotations.ClassMatching) bool
}

// NewKongHTTPValidator provides a new KongHTTPValidator object provided
// a controller-runtime client which will be used to retrieve reference objects
// such as consumer credentials secrets. If you do not pass a cached client
// here, the performance of this validator can get very poor at high scales.
func NewKongHTTPValidator(
	logger logr.Logger,
	managerClient client.Client,
	ingressClass string,
	servicesProvider AdminAPIServicesProvider,
	translatorFeatures translator.FeatureFlags,
	storer store.Storer,
) KongHTTPValidator {
	return KongHTTPValidator{
		Logger:                   logger,
		SecretGetter:             &managerClientSecretGetter{managerClient: managerClient},
		ConsumerGetter:           &managerClientConsumerGetter{managerClient: managerClient},
		Storer:                   storer,
		ManagerClient:            managerClient,
		AdminAPIServicesProvider: servicesProvider,
		TranslatorFeatures:       translatorFeatures,

		ingressClassMatcher:   annotations.IngressClassValidatorFuncFromObjectMeta(ingressClass),
		ingressV1ClassMatcher: annotations.IngressClassValidatorFuncFromV1Ingress(ingressClass),
	}
}

// ValidateConsumer checks if consumer has a Username and a consumer with
// the same username doesn't exist in Kong.
// If an error occurs during validation, it is returned as the last argument.
// The first boolean communicates if the consumer is valid or not and string
// holds a message if the entity is not valid.
func (validator KongHTTPValidator) ValidateConsumer(
	ctx context.Context,
	consumer kongv1.KongConsumer,
) (bool, string, error) {
	// ignore consumers that are being managed by another controller
	if !validator.ingressClassMatcher(&consumer.ObjectMeta, annotations.IngressClassKey, annotations.ExactClassMatch) {
		return true, "", nil
	}
	if err := kongplugin.ValidatePluginUniquenessPerObject(validator.Storer, &consumer); err != nil {
		return false, fmt.Sprintf("KongConsumer has invalid KongPlugin annotation: %s", err), nil
	}

	errText, err := validator.ensureConsumerDoesNotExistInGateway(ctx, consumer.Username)
	if err != nil || errText != "" {
		return false, errText, err
	}

	// if there are no credentials for this consumer, there's no need to move on
	// to credentials validation.
	if len(consumer.Credentials) == 0 {
		return true, "", nil
	}

	// pull all the managed consumers in order to build a validation index of
	// credentials so that the consumers credentials references can be validated.
	managedConsumers, err := validator.listManagedConsumers(ctx)
	if err != nil {
		return false, fmt.Sprintf("failed to fetch managed KongConsumers from cache: %s", err), nil
	}

	// retrieve the consumer's credentials secrets to validate them with the index
	credentials := make([]*corev1.Secret, 0, len(consumer.Credentials))
	ignoredSecrets := make(map[string]map[string]struct{})
	for _, secretName := range consumer.Credentials {
		// retrieve the credentials secret
		secret, err := validator.SecretGetter.GetSecret(consumer.Namespace, secretName)
		if err != nil {
			if apierrors.IsNotFound(err) {
				return false, fmt.Sprintf("%s: %s", ErrTextConsumerCredentialSecretNotFound, err), nil
			}
			return false, ErrTextFailedToRetrieveSecret, err
		}

		// do the basic credentials validation
		if err := credsvalidation.ValidateCredentials(secret); err != nil {
			return false, fmt.Sprintf("%s: %s", ErrTextConsumerCredentialValidationFailed, err), nil
		}

		// if valid, store it so we can index it for upcoming constraints validation
		credentials = append(credentials, secret)

		// later we'll build a global index of all credentials which is needed to
		// validate unique key constraints. That index should omit the secrets that
		// are referenced by this consumer to avoid duplication.
		if _, ok := ignoredSecrets[consumer.Namespace]; !ok {
			ignoredSecrets[consumer.Namespace] = make(map[string]struct{}, len(consumer.Credentials))
		}
		ignoredSecrets[consumer.Namespace][secretName] = struct{}{}
	}

	// unique constraints on consumer credentials are global to all consumers
	// and credentials, so we must build an index based on all existing credentials.
	// we ignore the secrets referenced by this consumer so that the index is not
	// testing them against themselves.
	credentialsIndex, err := globalValidationIndexForCredentials(ctx, validator.ManagerClient, managedConsumers, ignoredSecrets)
	if err != nil {
		return false, fmt.Sprintf("%s: %s", ErrTextConsumerCredentialValidationFailed, err), nil
	}

	// validate the consumer's credentials against the index of all managed
	// credentials to ensure they're not in violation of any unique constraints.
	for _, secret := range credentials {
		// do the unique constraints validation of the credentials using the credentials index
		if err := credentialsIndex.ValidateCredentialsForUniqueKeyConstraints(secret); err != nil {
			return false, fmt.Sprintf("%s: %s", ErrTextConsumerCredentialValidationFailed, err), nil
		}
	}

	return true, "", nil
}

func (validator KongHTTPValidator) ValidateConsumerGroup(
	ctx context.Context,
	consumerGroup kongv1beta1.KongConsumerGroup,
) (bool, string, error) {
	// Ignore ConsumerGroups that are being managed by another controller.
	if !validator.ingressClassMatcher(&consumerGroup.ObjectMeta, annotations.IngressClassKey, annotations.ExactClassMatch) {
		return true, "", nil
	}

	infoSvc, ok := validator.AdminAPIServicesProvider.GetInfoService()
	if !ok {
		return true, "", nil
	}
	info, err := infoSvc.Get(ctx)
	if err != nil {
		validator.Logger.V(util.DebugLevel).Info("Failed to fetch Kong info", "error", err)
		return false, ErrTextAdminAPIUnavailable, nil
	}
	version, err := kong.NewVersion(info.Version)
	if err != nil {
		validator.Logger.V(util.DebugLevel).Info("Failed to parse Kong version", "error", err)
	} else if !version.IsKongGatewayEnterprise() {
		return false, ErrTextConsumerGroupUnsupported, nil
	}

	cgs, ok := validator.AdminAPIServicesProvider.GetConsumerGroupsService()
	if !ok {
		return true, "", nil
	}
	// This check forbids consumer group creation if the license is invalid or missing.
	// There is no other way to robustly check the validity of a license than actually trying an enterprise feature.
	if _, _, err := cgs.List(ctx, &kong.ListOpt{Size: 0}); err != nil {
		switch {
		case kong.IsNotFoundErr(err):
			// This is the case when consumer group is not supported (Kong OSS) and previous version
			// check (if !version.IsKongGatewayEnterprise()) has been omitted due to a parsing error.
			return false, ErrTextConsumerGroupUnsupported, nil
		case kong.IsForbiddenErr(err):
			return false, ErrTextConsumerGroupUnlicensed, nil
		default:
			return false, fmt.Sprintf("%s: %s", ErrTextConsumerGroupUnexpected, err), nil
		}
	}
	return true, "", nil
}

// ValidateCredential checks if the secret contains a credential meant to
// be installed in Kong. If so, then it verifies if all the required fields
// are present in it or not. If valid, it returns true with an empty string,
// else it returns false with the error message. If an error happens during
// validation, error is returned.
func (validator KongHTTPValidator) ValidateCredential(ctx context.Context, secret corev1.Secret) (bool, string) {
	// If the secret doesn't specify a credential type it's not a credentials secret. We shouldn't actually reach this
	// codepath in practice because such secrets will be filtered out by the webhook secrets objectSelector and ignored.
	// However, installs could potentially use an outdated webhook definition. Prior to 3.2 we only filtered in code and
	// used a blanket selector.
	if _, err := util.ExtractKongCredentialType(&secret); err != nil {
		return true, ""
	}

	// If we know it's a credentials secret, we can ensure its base-level validity.
	if err := credsvalidation.ValidateCredentials(&secret); err != nil {
		return false, fmt.Sprintf("%s: %s", ErrTextConsumerCredentialValidationFailed, err)
	}

	// Credentials are validated further for unique key constraints only if they are referenced by a managed consumer
	// in the namespace, as such we pull a list of all consumers from the cached client to determine
	// if the credentials are referenced.
	managedConsumers, err := validator.listManagedConsumers(ctx)
	if err != nil {
		return false, fmt.Sprintf("failed to fetch managed KongConsumers from cache: %s", err)
	}

	// Verify whether this secret is referenced by any managed consumer.
	managedConsumersWithReferences := listManagedConsumersReferencingCredentialsSecret(secret, managedConsumers)
	if len(managedConsumersWithReferences) == 0 {
		// If no managed consumers reference this secret, its considered unmanaged, and we don't validate it
		// unless it becomes referenced by a managed consumer at a later time.
		return true, ""
	}

	// If base-level validation passed and the credential is referenced by a consumer,
	// we move on to create an index of all managed credentials so that we can verify that
	// the updates to this secret are not in violation of any unique key constraints.
	ignoreSecrets := map[string]map[string]struct{}{secret.Namespace: {secret.Name: {}}}
	credentialsIndex, err := globalValidationIndexForCredentials(ctx, validator.ManagerClient, managedConsumers, ignoreSecrets)
	if err != nil {
		return false, fmt.Sprintf("%s: %s", ErrTextConsumerCredentialValidationFailed, err)
	}

	// The index is built, now validate that the newly updated secret
	// is not in violation of any constraints.
	if err := credentialsIndex.ValidateCredentialsForUniqueKeyConstraints(&secret); err != nil {
		return false, fmt.Sprintf("%s: %s", ErrTextConsumerCredentialValidationFailed, err)
	}

	return true, ""
}

// ValidatePlugin checks if k8sPlugin is valid. It does so by performing
// an HTTP request to Kong's Admin API entity validation endpoints.
// If an error occurs during validation, it is returned as the last argument.
// The first boolean communicates if k8sPluign is valid or not and string
// holds a message if the entity is not valid.
func (validator KongHTTPValidator) ValidatePlugin(
	ctx context.Context,
	k8sPlugin kongv1.KongPlugin,
	overrideSecrets []*corev1.Secret,
) (bool, string, error) {
	var plugin kong.Plugin
	plugin.Name = kong.String(k8sPlugin.PluginName)
	var err error

	secretGetter := NewSecretGetterWithOverride(validator.SecretGetter, overrideSecrets)

	plugin.Config, err = kongstate.RawConfigurationWithPatchesToConfiguration(
		secretGetter,
		k8sPlugin.Namespace,
		k8sPlugin.Config,
		k8sPlugin.ConfigPatches,
	)
	if err != nil {
		return false, fmt.Sprintf("%s: %s", ErrTextPluginConfigInvalid, err), nil
	}
	if k8sPlugin.ConfigFrom != nil {
		config, err := kongstate.SecretToConfiguration(secretGetter, k8sPlugin.ConfigFrom.SecretValue, k8sPlugin.Namespace)
		if err != nil {
			return false, fmt.Sprintf("%s: %s", ErrTextPluginSecretConfigUnretrievable, err), nil
		}
		plugin.Config = config
	}
	if k8sPlugin.RunOn != "" {
		plugin.RunOn = kong.String(k8sPlugin.RunOn)
	}
	plugin.Ordering = k8sPlugin.Ordering
	plugin.Protocols = kong.StringSlice(kongv1.KongProtocolsToStrings(k8sPlugin.Protocols)...)

	errText, err := validator.validatePluginAgainstGatewaySchema(ctx, plugin)
	if err != nil || errText != "" {
		validator.Logger.Info("validate KongPlugin on Kong gateway failed",
			"plugin", fmt.Sprintf("%s/%s", k8sPlugin.Namespace, k8sPlugin.Name),
			"error", err,
		)
		return false, errText, err
	}

	return true, "", nil
}

// ValidateClusterPlugin transfers relevant fields from a KongClusterPlugin into a KongPlugin and then returns
// the result of ValidatePlugin for the derived KongPlugin.
func (validator KongHTTPValidator) ValidateClusterPlugin(
	ctx context.Context,
	k8sPlugin kongv1.KongClusterPlugin,
	overrideSecrets []*corev1.Secret,
) (bool, string, error) {
	var plugin kong.Plugin
	plugin.Name = kong.String(k8sPlugin.PluginName)
	var err error

	secretGetter := NewSecretGetterWithOverride(validator.SecretGetter, overrideSecrets)
	plugin.Config, err = kongstate.RawConfigurationWithNamespacedPatchesToConfiguration(
		secretGetter,
		k8sPlugin.Config,
		k8sPlugin.ConfigPatches,
	)
	if err != nil {
		return false, fmt.Sprintf("%s: %s", ErrTextPluginConfigInvalid, err), nil
	}

	if k8sPlugin.ConfigFrom != nil {
		config, err := kongstate.NamespacedSecretToConfiguration(secretGetter, k8sPlugin.ConfigFrom.SecretValue)
		if err != nil {
			return false, fmt.Sprintf("%s: %s", ErrTextPluginSecretConfigUnretrievable, err), nil
		}
		plugin.Config = config
	}
	if k8sPlugin.RunOn != "" {
		plugin.RunOn = kong.String(k8sPlugin.RunOn)
	}

	plugin.Ordering = k8sPlugin.Ordering
	plugin.Protocols = kong.StringSlice(kongv1.KongProtocolsToStrings(k8sPlugin.Protocols)...)

	errText, err := validator.validatePluginAgainstGatewaySchema(ctx, plugin)
	if err != nil || errText != "" {
		validator.Logger.Info("validate KongClusterPlugin on Kong gateway failed",
			"plugin", k8sPlugin.Name,
			"error", err,
		)
		return false, errText, err
	}

	return true, "", nil
}

func (validator KongHTTPValidator) ValidateGateway(
	ctx context.Context, gateway gatewayapi.Gateway,
) (bool, string, error) {
	// validate the gatewayclass reference
	gwc := gatewayapi.GatewayClass{}
	if err := validator.ManagerClient.Get(ctx, client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, &gwc); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return true, "", nil // not managed by this controller
		}
		return false, ErrTextCantRetrieveGatewayClass, err
	}

	// validate whether the gatewayclass is a supported class, if not
	// then this gateway belongs to another controller.
	if gwc.Spec.ControllerName != gatewaycontroller.GetControllerName() {
		return true, "", nil
	}

	return true, "", nil
}

func (validator KongHTTPValidator) ValidateHTTPRoute(
	ctx context.Context, httproute gatewayapi.HTTPRoute,
) (bool, string, error) {
	var routeValidator routeValidator = noOpRoutesValidator{}
	if routesSvc, ok := validator.AdminAPIServicesProvider.GetRoutesService(); ok {
		routeValidator = routesSvc
	}
	return gatewayvalidation.ValidateHTTPRoute(
		ctx, routeValidator, validator.TranslatorFeatures, &httproute, validator.ManagerClient,
	)
}

func (validator KongHTTPValidator) ValidateIngress(
	ctx context.Context, ingress netv1.Ingress,
) (bool, string, error) {
	// Ignore Ingresses that are being managed by another controller.
	if !validator.ingressClassMatcher(&ingress.ObjectMeta, annotations.IngressClassKey, annotations.ExactClassMatch) &&
		!validator.ingressV1ClassMatcher(&ingress, annotations.ExactClassMatch) {
		return true, "", nil
	}

	var routeValidator routeValidator = noOpRoutesValidator{}
	if routesSvc, ok := validator.AdminAPIServicesProvider.GetRoutesService(); ok {
		routeValidator = routesSvc
	}
	return ingressvalidation.ValidateIngress(ctx, routeValidator, validator.TranslatorFeatures, &ingress, validator.Logger, validator.Storer)
}

func (validator KongHTTPValidator) ValidateService(
	_ context.Context, service corev1.Service,
) (bool, string, error) {
	if err := kongplugin.ValidatePluginUniquenessPerObject(validator.Storer, &service); err != nil {
		return false, fmt.Sprintf("Service has invalid KongPlugin annotation: %s", err), nil
	}
	return true, "", nil
}

type routeValidator interface {
	Validate(context.Context, *kong.Route) (bool, string, error)
}

type noOpRoutesValidator struct{}

func (noOpRoutesValidator) Validate(_ context.Context, _ *kong.Route) (bool, string, error) {
	return true, "", nil
}

func (validator KongHTTPValidator) ValidateVault(ctx context.Context, k8sKongVault kongv1alpha1.KongVault) (bool, string, error) {
	// Ignore KongVaults that are being managed by another controller.
	if !validator.ingressClassMatcher(&k8sKongVault.ObjectMeta, annotations.IngressClassKey, annotations.ExactClassMatch) {
		return true, "", nil
	}
	config, err := kongstate.RawConfigToConfiguration(k8sKongVault.Spec.Config.Raw)
	if err != nil {
		return false, fmt.Sprintf(ErrTextVaultConfigUnmarshalFailed, err), nil
	}

	// list existing KongVaults and reject if the spec.prefix is duplicate with another `KongVault`.
	existingKongVaults := validator.Storer.ListKongVaults()
	dupeVault, hasDupe := lo.Find(existingKongVaults, func(v *kongv1alpha1.KongVault) bool {
		return v.Spec.Prefix == k8sKongVault.Spec.Prefix && v.Name != k8sKongVault.Name
	})
	if hasDupe {
		return false, fmt.Sprintf("spec.prefix %q is duplicate with existing KongVault %q",
			k8sKongVault.Spec.Prefix, dupeVault.Name), nil
	}

	kongVault := kong.Vault{
		Name:   kong.String(k8sKongVault.Spec.Backend),
		Prefix: kong.String(k8sKongVault.Spec.Prefix),
		Config: config,
	}
	if len(k8sKongVault.Spec.Description) > 0 {
		kongVault.Description = kong.String(k8sKongVault.Spec.Description)
	}
	errText, err := validator.validateVaultAgainstGatewaySchema(ctx, kongVault)
	if err != nil || errText != "" {
		return false, errText, err
	}
	return true, "", nil
}

// -----------------------------------------------------------------------------
// KongHTTPValidator - Private Methods
// -----------------------------------------------------------------------------

func (validator KongHTTPValidator) listManagedConsumers(ctx context.Context) ([]*kongv1.KongConsumer, error) {
	// Gather a list of all consumers from the cached client.
	consumers, err := validator.ConsumerGetter.ListAllConsumers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumers: %w", err)
	}

	// Reduce the consumer set to consumers managed by this controller.
	managedConsumers := make([]*kongv1.KongConsumer, 0)
	for _, consumer := range consumers {
		consumer := consumer
		if !validator.ingressClassMatcher(&consumer.ObjectMeta, annotations.IngressClassKey,
			annotations.ExactClassMatch) {
			// ignore consumers (and subsequently secrets) that are managed by other controllers
			continue
		}
		consumerCopy := consumer
		managedConsumers = append(managedConsumers, &consumerCopy)
	}

	return managedConsumers, nil
}

func (validator KongHTTPValidator) ensureConsumerDoesNotExistInGateway(ctx context.Context, username string) (string, error) {
	if consumerSvc, hasClient := validator.AdminAPIServicesProvider.GetConsumersService(); hasClient {
		// verify that the consumer is not already present in the data-plane
		c, err := consumerSvc.Get(ctx, &username)
		if err != nil {
			if !kong.IsNotFoundErr(err) {
				validator.Logger.Error(err, "Failed to fetch consumer from kong")
				return ErrTextConsumerUnretrievable, err
			}
		}
		if c != nil {
			return ErrTextConsumerExists, nil
		}
	}

	// if there's no client, do not verify existence with data-plane as there's none available
	return "", nil
}

func (validator KongHTTPValidator) validatePluginAgainstGatewaySchema(ctx context.Context, plugin kong.Plugin) (string, error) {
	pluginService, hasClient := validator.AdminAPIServicesProvider.GetPluginsService()
	if hasClient {
		isValid, msg, err := pluginService.Validate(ctx, &plugin)
		if err != nil {
			return ErrTextPluginConfigValidationFailed, err
		}
		if !isValid {
			return fmt.Sprintf(ErrTextPluginConfigViolatesSchema, msg), nil
		}
	}

	// if there's no client, do not verify with data-plane as there's none available
	return "", nil
}

func (validator KongHTTPValidator) validateVaultAgainstGatewaySchema(ctx context.Context, vault kong.Vault) (string, error) {
	vaultService, hasClient := validator.AdminAPIServicesProvider.GetVaultsService()
	if !hasClient {
		return "", nil
	}
	isValid, msg, err := vaultService.Validate(ctx, &vault)
	if err != nil {
		return ErrTextVaultUnableToValidate, err
	}
	if !isValid {
		return fmt.Sprintf(ErrTextVaultConfigValidationResultInvalid, msg), nil
	}
	return "", nil
}

type managerClientSecretGetter struct {
	managerClient client.Client
}

func (m *managerClientSecretGetter) GetSecret(namespace, name string) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	return secret, m.managerClient.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, secret)
}

type managerClientConsumerGetter struct {
	managerClient client.Client
}

func (m *managerClientConsumerGetter) ListAllConsumers(ctx context.Context) ([]kongv1.KongConsumer, error) {
	consumers := &kongv1.KongConsumerList{}
	if err := m.managerClient.List(ctx, consumers, &client.ListOptions{
		Namespace: corev1.NamespaceAll,
	}); err != nil {
		return nil, err
	}
	return consumers.Items, nil
}
