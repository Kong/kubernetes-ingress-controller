package kongstate

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configurationv1alpha1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1alpha1"
	configurationv1beta1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1beta1"
	"github.com/kong/kubernetes-configuration/v2/pkg/metadata"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation/consumers/credentials"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
)

// KongState holds the configuration that should be applied to Kong.
type KongState struct {
	Services       []Service
	Upstreams      []Upstream
	Certificates   []Certificate
	CACertificates []kong.CACertificate
	Licenses       []License
	Plugins        []Plugin
	Consumers      []Consumer
	ConsumerGroups []ConsumerGroup
	Vaults         []Vault

	CustomEntities map[string]*KongCustomEntityCollection
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (ks *KongState) SanitizedCopy(uuidGenerator util.UUIDGenerator) *KongState {
	return &KongState{
		Services:  ks.Services,
		Upstreams: ks.Upstreams,
		Certificates: func() []Certificate {
			if ks.Certificates == nil {
				return nil
			}

			return lo.Map(ks.Certificates, func(c Certificate, _ int) Certificate {
				return c.SanitizedCopy()
			})
		}(),
		CACertificates: ks.CACertificates,
		Plugins:        ks.Plugins,
		Consumers: func() []Consumer {
			if ks.Consumers == nil {
				return nil
			}

			return lo.Map(ks.Consumers, func(c Consumer, _ int) Consumer {
				return c.SanitizedCopy(uuidGenerator)
			})
		}(),
		Licenses: func() []License {
			if ks.Licenses == nil {
				return nil
			}
			return lo.Map(ks.Licenses, func(l License, _ int) License {
				return l.SanitizedCopy()
			})
		}(),
		ConsumerGroups: ks.ConsumerGroups,
		Vaults:         ks.Vaults,
		CustomEntities: ks.CustomEntities,
	}
}

func (ks *KongState) FillConsumersAndCredentials(
	_ logr.Logger,
	s store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
) {
	consumerIndex := make(map[string]Consumer)
	credentialsConflictDetector := NewCredentialConflictsDetector()

	// build consumer index
	for _, consumer := range s.ListKongConsumers() {
		var c Consumer
		// This is now enforce that the CRD level but we're keeping this just for those
		// rare cases where the CRD Validation Expressions are disabled.
		if consumer.Username == "" && consumer.CustomID == "" {
			failuresCollector.PushResourceFailure("no username or custom_id specified", consumer)
			continue
		}
		if consumer.Username != "" {
			c.Username = kong.String(consumer.Username)
		}
		if consumer.CustomID != "" {
			c.CustomID = kong.String(consumer.CustomID)
		}
		c.K8sKongConsumer = *consumer
		c.Tags = util.GenerateTagsForObject(consumer)

		// Get consumer groups
		for _, cgName := range consumer.ConsumerGroups {
			cg, err := s.GetKongConsumerGroup(consumer.Namespace, cgName)
			if err != nil {
				failuresCollector.PushResourceFailure(fmt.Sprintf("nonexistent consumer group: %q", err), consumer)
				continue
			}
			c.ConsumerGroups = append(c.ConsumerGroups, kong.ConsumerGroup{
				Name: &cg.Name,
			})
		}

		for _, cred := range consumer.Credentials {
			pushCredentialResourceFailures := func(message string) {
				failuresCollector.PushResourceFailure(fmt.Sprintf("credential %q failure: %s", cred, message), consumer)
			}
			secret, err := s.GetSecret(consumer.Namespace, cred)
			if err != nil {
				pushCredentialResourceFailures(fmt.Sprintf("Failed to fetch secret: %v", err))
				continue
			}
			credConfig := map[string]any{}
			// try the label first. if it's present, no need to check the field
			credType, err := util.ExtractKongCredentialType(secret)
			if err != nil {
				pushCredentialResourceFailures(fmt.Sprintf("could not load credential from Secret: %s", err))
			}
			if !credentials.SupportedTypes.Has(credType) {
				pushCredentialResourceFailures(
					fmt.Sprintf("failed to provision credential: unsupported credential type: %q", credType),
				)
				continue
			}
			for k, v := range secret.Data {
				// TODO populate these based on schema from Kong
				// and remove this workaround
				if k == "redirect_uris" {
					credConfig[k] = strings.Split(string(v), ",")
					continue
				}
				// TODO this is a credential type-agnostic mutation that should only apply to Oauth2 credentials.
				// However, the credential-specific code after deals only in interface{}s, and we can't fix individual
				// keys. To handle this properly we'd need to refactor the types used in all following code.
				if k == "hash_secret" {
					boolVal, err := strconv.ParseBool(string(v))
					if err != nil {
						// add a translation error here to tell that parsing hash_secret failed.
						pushCredentialResourceFailures(
							fmt.Sprintf("Failed to parse hash_secret to bool: %v. defaulting to false", err),
						)
						credConfig[k] = false
					} else {
						credConfig[k] = boolVal
					}
					continue
				}
				// ttl is a field that only appears in keyAuth credentials and has int type.
				// Same as above, we cannot fix individual keys after translated to credConfig.
				if k == "ttl" {
					intVal, err := strconv.Atoi(string(v))
					if err != nil {
						// add a translation error here to tell that parsing TTL failed.
						pushCredentialResourceFailures(
							fmt.Sprintf("Failed to parse ttl to int: %v, skipfilling the field", err),
						)
					} else {
						credConfig[k] = intVal
					}
					continue
				}
				credConfig[k] = string(v)
			}
			tags := util.GenerateTagsForObject(secret)
			cred, err := c.SetCredential(credType, credConfig, tags)
			if err != nil {
				pushCredentialResourceFailures(
					fmt.Sprintf("failed to provision credential: %v", err),
				)
				continue
			}
			credentialsConflictDetector.RegisterForConflictDetection(cred, secret, consumer)
		}

		consumerIndex[consumer.Namespace+"/"+consumer.Name] = c
	}

	// Detect credential conflicts, push failures for them, and remove the consumers that have conflicts.
	for _, conflict := range credentialsConflictDetector.DetectConflicts() {
		// Push the failure for the credential Secret first.
		failuresCollector.PushResourceFailure(
			fmt.Sprintf("Secret has conflicting credential: %s", conflict.Message),
			conflict.Credential.CredentialSecret,
		)

		// Check if the consumer is still in the index...
		consumerKey := client.ObjectKeyFromObject(conflict.Credential.Consumer).String()
		_, ok := consumerIndex[consumerKey]
		if !ok {
			// If it's not, it's already been removed, and we can skip this step.
			continue
		}
		// ...if it is, push the failure for the KongConsumer and remove it from the index.
		failuresCollector.PushResourceFailure(
			fmt.Sprintf("Consumer has conflicting credentials: %s", conflict.Message),
			conflict.Credential.Consumer,
		)
		delete(consumerIndex, consumerKey)
	}

	// populate the consumer in the state
	for _, c := range consumerIndex {
		ks.Consumers = append(ks.Consumers, c)
	}
}

func (ks *KongState) FillConsumerGroups(_ logr.Logger, s store.Storer) {
	for _, cg := range s.ListKongConsumerGroups() {
		ks.ConsumerGroups = append(ks.ConsumerGroups, ConsumerGroup{
			ConsumerGroup: kong.ConsumerGroup{
				Name: kong.String(cg.Name),
				Tags: util.GenerateTagsForObject(cg),
			},
			K8sKongConsumerGroup: *cg,
		})
	}
}

// servicesAsObjects returns a corev1.Service as a client.Object. It's used as a helper with lo.Map to return something
// acceptable to functions that accept multiple client.Objects, since simply expanding the Service slice results in a
// compiler error.
func servicesAsObjects(svc *corev1.Service, _ int) client.Object {
	return svc
}

func (ks *KongState) FillOverrides(
	logger logr.Logger,
	s store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
	kongVersion semver.Version,
) {
	for i := 0; i < len(ks.Services); i++ {
		// Services
		if err := ks.Services[i].override(); err != nil {
			servicesGroup := lo.Values(ks.Services[i].K8sServices)
			failuresCollector.PushResourceFailure(err.Error(), lo.Map(servicesGroup, servicesAsObjects)...)
		}

		// Routes
		for j := 0; j < len(ks.Services[i].Routes); j++ {
			// Routes, opposed to Services, are not validated in their override method therefore we have no error check here.
			// Routes nested under Services here do not include their original parent object info in kongstate. Translators
			// convert this into a util.K8sObjectInfo, which includes name/namespace/GVK as strings. Unfortunately we can't
			// really convert this into a client.Object for use with the failures collector, and plumbing the original object
			// down into the kongstate.Route copy looked a bit annoying. Protocol validation for routes instead lives in the
			// HTTPRoute and Ingress translators (these may override to ws/wss, whereas the others are expected to derive
			// their protcol from the resource type alone).
			ks.Services[i].Routes[j].Override(logger)
		}
	}

	ks.FillUpstreamOverrides(s, logger, failuresCollector, kongVersion)
}

func (ks *KongState) FillUpstreamOverrides(
	s store.Storer,
	logger logr.Logger,
	failuresCollector *failures.ResourceFailuresCollector,
	kongVersion semver.Version,
) {
	for i := 0; i < len(ks.Upstreams); i++ {
		servicesGroup := lo.Values(ks.Upstreams[i].Service.K8sServices)

		// In case `konghq.com/override` annotation is set on any of the services, we should log a deprecation error.
		maybeLogKongIngressDeprecationError(logger, servicesGroup)

		kongIngress, err := getKongIngressForServices(s, servicesGroup)
		if err != nil {
			failuresCollector.PushResourceFailure(err.Error(), lo.Map(servicesGroup, servicesAsObjects)...)
		} else {
			for _, svc := range servicesGroup {
				ks.Upstreams[i].override(kongIngress, svc)
			}
		}

		kongUpstreamPolicy, err := GetKongUpstreamPolicyForServices(s, servicesGroup)
		if err != nil {
			failuresCollector.PushResourceFailure(err.Error(), lo.Map(servicesGroup, servicesAsObjects)...)
		} else if kongUpstreamPolicy != nil {
			if kongUpstreamPolicy.Spec.Algorithm != nil && *kongUpstreamPolicy.Spec.Algorithm == "sticky-sessions" &&
				kongVersion.LT(versions.KongStickySessionsCutoff) {
				failuresCollector.PushResourceFailure(
					fmt.Sprintf(
						"sticky sessions algorithm specified in KongUpstreamPolicy '%s' is not supported with Kong Gateway versions < %s",
						kongUpstreamPolicy.Name, versions.KongStickySessionsCutoff,
					),
					lo.Map(servicesGroup, servicesAsObjects)...,
				)
				continue
			}
			ks.Upstreams[i].overrideByKongUpstreamPolicy(kongUpstreamPolicy)
		}
	}
}

// compareKongVault compares two `KongVault`s when they have the same `spec.prefix`.
// When 2 or more KongVaults have the same prefix, only one of them is translated.
// It returns true when v1 has higher priority then v2, by the following order:
// - The one created earlier (earlier `creationTimestamp`) takes precedence.
// - If the creationTimestamp equals, the one with smaller lexical order (`<` for strings) takes precedence.
func compareKongVault(v1, v2 *configurationv1alpha1.KongVault) bool {
	if v1.CreationTimestamp.Before(&v2.CreationTimestamp) {
		return true
	}
	if v2.CreationTimestamp.Before(&v1.CreationTimestamp) {
		return false
	}
	// None of them can be seen created before the other (equal or not comparable), compare by lexical order of name.
	return v1.Name < v2.Name
}

func (ks *KongState) FillVaults(
	logger logr.Logger,
	s store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
) {
	// List all vaults and reject the KongVaults with duplicate prefix to prevent invalid Kong configuration generated.
	allKongVaults := s.ListKongVaults()
	prefixToKongVault := map[string]*configurationv1alpha1.KongVault{}
	for _, vault := range allKongVaults {
		prefix := vault.Spec.Prefix
		existingVault, ok := prefixToKongVault[prefix]
		if !ok {
			prefixToKongVault[prefix] = vault
			continue
		}
		// ok == true, which means we have KongVaults with same spec.prefix
		if compareKongVault(existingVault, vault) {
			// the one already in the map has higher priority, the current KongVault is rejected.
			failuresCollector.PushResourceFailure(
				fmt.Sprintf("spec.prefix %q is duplicate", prefix), vault,
			)
		} else {
			// the current vault has the higher priority, the existing one is rejected and taken out of the map.
			prefixToKongVault[prefix] = vault
			failuresCollector.PushResourceFailure(
				fmt.Sprintf("spec.prefix %q is duplicate", prefix), existingVault,
			)
		}
	}

	for _, vault := range prefixToKongVault {
		config, err := RawConfigToConfiguration(vault.Spec.Config.Raw)
		if err != nil {
			logger.Error(err, "failed to parse configuration of vault to JSON", "name", vault.Name)
			failuresCollector.PushResourceFailure(
				fmt.Sprintf("failed to parse configuration of vault %q to JSON: %v", vault.Name, err),
				vault,
			)
			continue
		}
		kongVault := kong.Vault{
			Name:   kong.String(vault.Spec.Backend),
			Prefix: kong.String(vault.Spec.Prefix),
			Config: config,
			Tags:   util.GenerateTagsForObject(vault),
		}
		if len(vault.Spec.Description) > 0 {
			kongVault.Description = kong.String(vault.Spec.Description)
		}
		ks.Vaults = append(ks.Vaults, Vault{
			Vault:        kongVault,
			K8sKongVault: vault.DeepCopy(),
		})
	}
}

type NamespacedKongPlugin struct {
	Namespace string
	Name      string
}

func (ks *KongState) getPluginRelations(
	cacheStore store.Storer, log logr.Logger, failuresCollector *failures.ResourceFailuresCollector,
) map[string]rels.ForeignRelations {
	// KongPlugin key (KongPlugin's name:namespace) to corresponding associations
	pluginRels := map[string]rels.ForeignRelations{}

	type entityRelationType int
	const (
		ConsumerRelation      entityRelationType = iota
		ConsumerGroupRelation entityRelationType = iota
		RouteRelation         entityRelationType = iota
		ServiceRelation       entityRelationType = iota
	)
	addRelation := func(referrer client.Object, plugin k8stypes.NamespacedName, identifier string, t entityRelationType) {
		// There are 2 types of KongPlugin references: local and remote.
		// A local reference is one where the KongPlugin is in the same namespace as the referrer.
		// A remote reference is one where the KongPlugin is in a different namespace.
		// By default a KongPlugin is considered local.
		// If the plugin has a namespace specified, it is considered remote.
		//
		// The referrer is the entity that the KongPlugin is associated with.
		//
		// Code in buildPlugins() will combine plugin associations into
		// multi-entity plugins within the local namespace
		namespace, err := extractReferredPluginNamespaceIfAllowed(log, cacheStore, referrer, plugin)
		if err != nil {
			if errors.As(err, &referredPluginNotAllowedError{}) {
				failuresCollector.PushResourceFailure(err.Error(), referrer)
			} else {
				log.Error(err, "could not bind requested plugin", "plugin", plugin.Name, "namespace", plugin.Namespace)
			}
			return
		}

		pluginKey := namespace + ":" + plugin.Name
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = rels.ForeignRelations{}
		}
		fr := rels.FR{
			Identifier: identifier,
			Referer:    referrer,
		}
		switch t {
		case ConsumerRelation:
			relations.Consumer = append(relations.Consumer, fr)
		case ConsumerGroupRelation:
			relations.ConsumerGroup = append(relations.ConsumerGroup, fr)
		case RouteRelation:
			relations.Route = append(relations.Route, fr)
		case ServiceRelation:
			relations.Service = append(relations.Service, fr)
		}
		pluginRels[pluginKey] = relations
	}

	for i := range ks.Services {
		for _, svc := range ks.Services[i].K8sServices {
			for _, plugin := range metadata.ExtractPluginsNamespacedNames(svc) {
				addRelation(svc, plugin, *ks.Services[i].Name, ServiceRelation)
			}
		}

		for j := range ks.Services[i].Routes {
			ingress := ks.Services[i].Routes[j].Ingress
			for _, plugin := range metadata.ExtractPluginsNamespacedNames(ingress) {
				// pretend we have a full Ingress struct for reference checks.
				// REVIEW: we only need an object to carry type meta and object meta here, maybe we should create some other types of virtual object here?
				virtualIngress := netv1.Ingress{
					// Fill the actual type of the object for reference checks.
					TypeMeta: util.TypeMetaFromGVK(ingress.GroupVersionKind),
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ingress.Namespace,
						Name:      ingress.Name,
					},
				}
				addRelation(&virtualIngress, plugin, *ks.Services[i].Routes[j].Name, RouteRelation)
			}
		}
	}

	for _, c := range ks.Consumers {
		if c.Username == nil && c.CustomID == nil {
			// Either Custom ID or username has to be filled so this shouldn't
			// happen but let's log here just in case.
			log.Error(nil, "no username or custom_id specified", "consumer", c.K8sKongConsumer.Name)
			break
		}

		var identifier string
		switch {
		case c.Username != nil:
			identifier = *c.Username
		case c.CustomID != nil:
			identifier = *c.CustomID
		}

		for _, plugin := range metadata.ExtractPluginsNamespacedNames(&c.K8sKongConsumer) {
			addRelation(&c.K8sKongConsumer, plugin, identifier, ConsumerRelation)
		}
	}

	for _, cg := range ks.ConsumerGroups {
		for _, plugin := range metadata.ExtractPluginsNamespacedNames(&cg.K8sKongConsumerGroup) {
			addRelation(&cg.K8sKongConsumerGroup, plugin, *cg.Name, ConsumerGroupRelation)
		}
	}

	return pluginRels
}

type pluginReference struct {
	Referrer  client.Object
	Namespace string
	Name      string
}

func buildPlugins(
	logger logr.Logger,
	s store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
	pluginRels map[string]rels.ForeignRelations,
) []Plugin {
	var plugins []Plugin

	for pluginIdentifier, relations := range pluginRels {
		identifier := strings.Split(pluginIdentifier, ":")
		namespace, kongPluginName := identifier[0], identifier[1]
		k8sPlugin, k8sClusterPlugin, pluginFound, err := getKongPluginOrKongClusterPlugin(s, namespace, kongPluginName)
		if err != nil {
			logger.Error(err, "Failed to fetch KongPlugin resource",
				"kongplugin_name", kongPluginName,
				"kongplugin_namespace", namespace)
			continue
		}
		// For non-existing plugin push failure K8s Events for all K8s object that reference it.
		if !pluginFound {
			for _, relation := range relations.GetCombinations() {
				for _, ref := range relation.ToList() {
					if !ref.IsEmpty() {
						failuresCollector.PushResourceFailure(
							fmt.Sprintf("referenced KongPlugin or KongClusterPlugin %q does not exist", kongPluginName),
							ref.Referer,
						)
					}
				}
			}
			continue
		}
		var plugin Plugin
		if k8sPlugin != nil {
			plugin, err = kongPluginFromK8SPlugin(s, *k8sPlugin)
			if err != nil {
				failuresCollector.PushResourceFailure(err.Error(), k8sPlugin)
				continue
			}
		}
		if k8sClusterPlugin != nil {
			plugin, err = kongPluginFromK8SClusterPlugin(s, *k8sClusterPlugin)
			if err != nil {
				failuresCollector.PushResourceFailure(err.Error(), k8sClusterPlugin)
				continue
			}
		}

		usedInstanceNames := sets.New[string]()
		for _, rel := range relations.GetCombinations() {
			plugin := plugin.DeepCopy()

			var toHash bytes.Buffer

			// ID is populated because that is read by decK and in_memory translator too
			if !rel.Service.IsEmpty() {
				relSvc := rel.Service.Identifier
				plugin.Service = &kong.Service{ID: kong.String(relSvc)}
				toHash.WriteString("_service-" + relSvc)
			}
			if !rel.Route.IsEmpty() {
				relRoute := rel.Route.Identifier
				plugin.Route = &kong.Route{ID: kong.String(relRoute)}
				toHash.WriteString("_route-" + relRoute)
			}
			if !rel.Consumer.IsEmpty() {
				relConsumer := rel.Consumer.Identifier
				plugin.Consumer = &kong.Consumer{ID: kong.String(relConsumer)}
				toHash.WriteString("_consumer-" + relConsumer)
			}
			if !rel.ConsumerGroup.IsEmpty() {
				relConsumerGroup := rel.ConsumerGroup.Identifier
				plugin.ConsumerGroup = &kong.ConsumerGroup{ID: kong.String(relConsumerGroup)}
				toHash.WriteString("_group-" + relConsumerGroup)
			}

			// instance_name must be unique. Using the same KongPlugin on multiple resources will result in duplicates
			// unless we add some sort of suffix.
			if plugin.InstanceName != nil {
				sha := sha256.Sum256(toHash.Bytes())
				suffix := fmt.Sprintf("%x", sha)
				short := suffix[:9]
				suffixed := fmt.Sprintf("%s-%s", *plugin.InstanceName, short)
				if usedInstanceNames.Has(suffixed) {
					// in the unlikely event of a short hash collision, use the full one
					suffixed = fmt.Sprintf("%s-%s", *plugin.InstanceName, suffix)
				}
				usedInstanceNames.Insert(suffixed)
				plugin.InstanceName = &suffixed
			}
			plugins = append(plugins, plugin)
		}
	}

	gKCPs, err := globalKongClusterPlugins(logger, s)
	if err != nil {
		logger.Error(err, "Failed to fetch global plugins")
	}
	// global plugins have no instance_name transform as they can only be applied once
	plugins = append(plugins, gKCPs...)

	return plugins
}

func globalKongClusterPlugins(logger logr.Logger, s store.Storer) ([]Plugin, error) {
	res := make(map[string]Plugin)
	var duplicates []string // keep track of duplicate
	// TODO respect the oldest CRD
	// Current behavior is to skip creating the plugin but in case
	// of duplicate plugin definitions, we should respect the oldest one
	// This is important since if a user comes in to k8s and creates a new
	// CRD, the user now deleted an older plugin

	globalClusterPlugins, err := s.ListGlobalKongClusterPlugins()
	if err != nil {
		return nil, fmt.Errorf("error listing global KongClusterPlugins: %w", err)
	}
	for i := range globalClusterPlugins {
		k8sPlugin := *globalClusterPlugins[i]
		pluginName := k8sPlugin.PluginName
		// empty pluginName skip it
		if pluginName == "" {
			logger.Error(nil, "Invalid KongClusterPlugin: empty plugin property",
				"kongclusterplugin_name", k8sPlugin.Name)
			continue
		}
		if _, ok := res[pluginName]; ok {
			logger.Error(nil, "Multiple KongPlugin with 'global' label found, cannot apply",
				"kongplugin_name", pluginName)
			duplicates = append(duplicates, pluginName)
			continue
		}
		if plugin, err := kongPluginFromK8SClusterPlugin(s, k8sPlugin); err == nil {
			res[pluginName] = plugin
		} else {
			logger.Error(err, "Failed to generate configuration from KongClusterPlugin",
				"kongclusterplugin_name", k8sPlugin.Name)
		}
	}
	for _, plugin := range duplicates {
		delete(res, plugin)
	}
	return lo.Values(res), nil
}

func (ks *KongState) FillPlugins(
	log logr.Logger,
	s store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
) {
	ks.Plugins = buildPlugins(log, s, failuresCollector, ks.getPluginRelations(s, log, failuresCollector))
}

// FillIDs iterates over the KongState and fills in the ID field for each entity
// that supports the FillID method (these are Service, Route, Consumer and Consumer
// Group). It makes their IDs deterministic, enabling their correct identification
// in external systems (e.g. Konnect Analytics).
// The workspace parameter is used for guarantee that the ID is unique across all workspaces,
// as required by Kong gateway.
func (ks *KongState) FillIDs(logger logr.Logger, workspace string) {
	for svcIndex, svc := range ks.Services {
		if err := svc.FillID(workspace); err != nil {
			logger.Error(err, "Failed to fill ID for service", "service_name", *svc.Name)
		} else {
			ks.Services[svcIndex] = svc
		}

		for routeIndex, route := range svc.Routes {
			if err := route.FillID(workspace); err != nil {
				logger.Error(err, "Failed to fill ID for route", "route_name", *route.Name)
			} else {
				ks.Services[svcIndex].Routes[routeIndex] = route
			}
		}
	}

	for consumerIndex, consumer := range ks.Consumers {
		if err := consumer.FillID(workspace); err != nil {
			logger.Error(err, "Failed to fill ID for consumer", "consumer_name", consumer.FriendlyName())
		} else {
			ks.Consumers[consumerIndex] = consumer
		}
	}

	for consumerGroupIndex, consumerGroup := range ks.ConsumerGroups {
		if err := consumerGroup.FillID(workspace); err != nil {
			logger.Error(err, "Failed to fill ID for consumer group", "consumer_group_name", *consumerGroup.Name)
		} else {
			ks.ConsumerGroups[consumerGroupIndex] = consumerGroup
		}
	}

	for valutIndex, vault := range ks.Vaults {
		if err := vault.FillID(workspace); err != nil {
			logger.Error(err, "Failed to fill ID for vault", "vault_name", vault.FriendlyName())
		} else {
			ks.Vaults[valutIndex] = vault
		}
	}

	for pluginIndex, plugin := range ks.Plugins {
		if err := plugin.FillID(workspace); err != nil {
			logger.Error(err, "failed to fill ID for plugin", "plugin_name", plugin.Name)
		} else {
			ks.Plugins[pluginIndex] = plugin
		}
	}
}

// maybeLogKongIngressDeprecationError iterates over services and logs a deprecation error if a service
// is annotated with `konghq.com/override` annotation.
func maybeLogKongIngressDeprecationError(logger logr.Logger, services []*corev1.Service) {
	for _, svc := range services {
		_, upstreamPolicyAnnotationSet := annotations.ExtractUpstreamPolicy(svc.Annotations)
		kongOverrideAnnotationSet := annotations.ExtractConfigurationName(svc.Annotations) != ""

		// If both `konghq.com/override` and `konghq.com/upstream-policy` are set, we should log a more specific error.
		if kongOverrideAnnotationSet && upstreamPolicyAnnotationSet {
			logger.Error(nil, fmt.Sprintf("Service uses both %s and %s annotations, should use only %s annotation. Settings "+
				"from %s will take precedence",
				annotations.AnnotationPrefix+annotations.ConfigurationKey,
				configurationv1beta1.KongUpstreamPolicyAnnotationKey,
				configurationv1beta1.KongUpstreamPolicyAnnotationKey,
				configurationv1beta1.KongUpstreamPolicyAnnotationKey),
				"namespace", svc.Namespace, "name", svc.Name,
			)
		}

		// In case it's just `konghq.com/override` set, we should log a deprecation error.
		if kongOverrideAnnotationSet {
			logger.Error(nil, fmt.Sprintf(
				"Service uses deprecated %s annotation and KongIngress, migrate to %s and KongUpstreamPolicy",
				annotations.AnnotationPrefix+annotations.ConfigurationKey,
				configurationv1beta1.KongUpstreamPolicyAnnotationKey),
				"namespace", svc.Namespace, "name", svc.Name,
			)
		}
	}
}

// getServiceIDFromPluginRels returns the ID of the services which a plugin refers to in RelatedEntitiesRef.
// It fills the IDs of services directly referred, and IDs of services where referred routes attaches to.
func getServiceIDFromPluginRels(
	log logr.Logger, rer RelatedEntitiesRef, routeAttachedService map[string]*Service, workspace string,
) []rels.FR {
	// Return IDs of directly referred services.
	if len(rer.Services) > 0 {
		return lo.FilterMap(rer.Services, func(s *Service, _ int) (rels.FR, bool) {
			if err := s.FillID(workspace); err != nil {
				log.Error(err, "failed to fill ID for service")
				return rels.FR{}, false
			}
			return rels.FR{Identifier: *s.ID, Referer: s.Parent}, true
		})
	}
	// Returns IDs of services where the referred routes attaches.
	serviceIDs := lo.FilterMap(
		rer.Routes,
		func(r *Route, _ int) (rels.FR, bool) {
			svc, ok := routeAttachedService[*r.Name]
			if !ok {
				return rels.FR{}, false
			}
			if err := svc.FillID(workspace); err != nil {
				log.Error(err, "failed to fill ID for service")
				return rels.FR{}, false
			}
			return rels.FR{
				Identifier: *svc.ID,
				Referer:    svc.Parent,
			}, true
		},
	)
	return lo.UniqBy(serviceIDs, func(fr rels.FR) string {
		return fr.Identifier
	})
}

// getPluginRelatedEntitiesRef gets services/routes/consumers referred by each plugin and returns the pointer to them.
// It is for the custom entities to fill the IDs of the referred entities into their "foreign" fields.
// It basically does the same thing as getPluginRelations but stores the pointers
// because we need to call the FillID method of the entities to fetch the ID,
// as Kong gateway requires IDs in the "foreign" fields (not other identifiers such as name.)
//
// TODO: refactor the building of plugin related entities and share the result between here and building plugins:
// https://github.com/Kong/kubernetes-ingress-controller/issues/6115
func (ks *KongState) getPluginRelatedEntitiesRef(
	cacheStore store.Storer, log logr.Logger, failuresCollector *failures.ResourceFailuresCollector,
) PluginRelatedEntitiesRefs {
	pluginRels := PluginRelatedEntitiesRefs{
		RelatedEntities:      map[string]RelatedEntitiesRef{},
		RouteAttachedService: map[string]*Service{},
	}
	addRelation := func(referrer client.Object, plugin k8stypes.NamespacedName, entity any) {
		namespace, err := extractReferredPluginNamespaceIfAllowed(log, cacheStore, referrer, plugin)
		if err != nil {
			if errors.As(err, &referredPluginNotAllowedError{}) {
				failuresCollector.PushResourceFailure(err.Error(), referrer)
			} else {
				log.Error(err, "could not bind requested plugin", "plugin", plugin.Name, "namespace", plugin.Namespace)
			}
			return
		}
		pluginKey := namespace + ":" + plugin.Name
		relations, ok := pluginRels.RelatedEntities[pluginKey]
		if !ok {
			relations = RelatedEntitiesRef{}
		}
		switch e := entity.(type) {
		case *Consumer:
			relations.Consumers = append(relations.Consumers, e)
		case *Route:
			relations.Routes = append(relations.Routes, e)
		case *Service:
			relations.Services = append(relations.Services, e)
		}
		pluginRels.RelatedEntities[pluginKey] = relations
	}

	for i := range ks.Services {
		for _, svc := range ks.Services[i].K8sServices {
			for _, plugin := range metadata.ExtractPluginsNamespacedNames(svc) {
				addRelation(svc, plugin, &ks.Services[i])
			}
		}

		for j, r := range ks.Services[i].Routes {
			ingress := ks.Services[i].Routes[j].Ingress
			for _, plugin := range metadata.ExtractPluginsNamespacedNames(ingress) {
				// Pretend we have a full Ingress struct for reference checks.
				virtualIngress := netv1.Ingress{
					// Fill the actual type of the object for reference checks.
					TypeMeta: util.TypeMetaFromGVK(ingress.GroupVersionKind),
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ingress.Namespace,
						Name:      ingress.Name,
					},
				}
				addRelation(&virtualIngress, plugin, &ks.Services[i].Routes[j])
				// For some entities, we need to find the referred service via route
				// but the `Service` field of translated routes are empty because we do not want the field to appear in final declarative config.
				// So we need to maintain a map from route name to service to find the service object.
				pluginRels.RouteAttachedService[*r.Name] = &ks.Services[i]
			}
		}
	}

	for i, c := range ks.Consumers {
		for _, plugin := range metadata.ExtractPluginsNamespacedNames(&c.K8sKongConsumer) {
			addRelation(&c.K8sKongConsumer, plugin, &ks.Consumers[i])
		}
	}
	return pluginRels
}

type referredPluginNotAllowedError struct {
	pluginReference gatewayapi.PluginLabelReference
}

func (e referredPluginNotAllowedError) Error() string {
	return fmt.Sprintf("no grant found to referenced %q plugin in the requested remote KongPlugin bind", e.pluginReference)
}

func extractReferredPluginNamespaceIfAllowed(
	log logr.Logger, cacheStore store.Storer, referrer client.Object, plugin k8stypes.NamespacedName,
) (string, error) {
	// There are 2 types of KongPlugin references: local and remote.
	// A local reference is one where the KongPlugin is in the same namespace as the referrer.
	// A remote reference is one where the KongPlugin is in a different namespace.
	// By default a KongPlugin is considered local.
	// If the plugin has a namespace specified, it is considered remote.
	//
	// The referrer is the entity that the KongPlugin is associated with.
	if plugin.Namespace == "" {
		return referrer.GetNamespace(), nil
	}

	// remote KongPlugin, permitted if ReferenceGrant allows.
	err := isRemotePluginReferenceAllowed(
		log,
		cacheStore,
		pluginReference{
			Referrer:  referrer,
			Namespace: plugin.Namespace,
			Name:      plugin.Name,
		},
	)
	if err != nil {
		return "", err
	}
	return plugin.Namespace, nil
}

func isRemotePluginReferenceAllowed(log logr.Logger, s store.Storer, r pluginReference) error {
	// remote plugin. considered part of this namespace if a suitable ReferenceGrant exists
	grants, err := s.ListReferenceGrants()
	if err != nil {
		return fmt.Errorf("could not retrieve ReferenceGrants from store when building plugin relations map: %w", err)
	}
	allowed := gatewayapi.GetPermittedForReferenceGrantFrom(
		log,
		gatewayapi.ReferenceGrantFrom{
			Group:     gatewayapi.Group(r.Referrer.GetObjectKind().GroupVersionKind().Group),
			Kind:      gatewayapi.Kind(r.Referrer.GetObjectKind().GroupVersionKind().Kind),
			Namespace: gatewayapi.Namespace(r.Referrer.GetNamespace()),
		},
		grants,
	)

	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/6000
	// Our reference checking code wasn't originally designed to handle multiple object types and relationships and
	// wasn't designed with much guidance for future usage. It expects something akin to the BackendRef section of an
	// HTTPRoute or similar, since that's what it was originally designed for. A future simplified model should probably
	// work with source and target client.Object resources derived from the particular resources' reference specs, as
	// those reference formats aren't necessarily consistent. We should possibly push for a standard upstream utility as
	// part of https://github.com/kubernetes/enhancements/issues/3766 if it proceeds further, as this is can be difficult
	// to wrap your head around when deep in the code that's checking an individual use case. A standard model for "these
	// are the inputs and outputs for a reference grant check, derive them from your spec" should help avoid inconsistency
	// in check implementation.

	// we don't have a full plugin resource here for the grant checker, so we build a fake one with the correct
	// name and namespace.
	pluginReference := gatewayapi.PluginLabelReference{
		Namespace: &r.Namespace,
		Name:      r.Name,
	}

	// Because we should check whether it is allowed to refer FROM the referrer TO the plugin here,
	// we should put the referrer on the "target" and the plugin on the "backendRef".

	log.V(logging.TraceLevel).Info("requested grant to plugins",
		"from-namespace", r.Referrer.GetNamespace(),
		"from-group", r.Referrer.GetObjectKind().GroupVersionKind().Group,
		"from-kind", r.Referrer.GetObjectKind().GroupVersionKind().Kind,
		"to-namespace", r.Namespace,
		"to-name", r.Name,
	)

	if !gatewayapi.NewRefCheckerForKongPlugin(log, r.Referrer, pluginReference).IsRefAllowedByGrant(allowed) {
		return referredPluginNotAllowedError{
			pluginReference: pluginReference,
		}
	}
	return nil
}
