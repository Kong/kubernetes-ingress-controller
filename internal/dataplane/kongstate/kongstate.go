package kongstate

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/admission/validation/consumers/credentials"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (ks *KongState) SanitizedCopy() *KongState {
	return &KongState{
		Services:  ks.Services,
		Upstreams: ks.Upstreams,
		Certificates: func() (res []Certificate) {
			for _, v := range ks.Certificates {
				res = append(res, *v.SanitizedCopy())
			}
			return
		}(),
		CACertificates: ks.CACertificates,
		Plugins:        ks.Plugins,
		Consumers: func() (res []Consumer) {
			for _, v := range ks.Consumers {
				res = append(res, *v.SanitizedCopy())
			}
			return
		}(),
		Licenses: func() (res []License) {
			for _, v := range ks.Licenses {
				res = append(res, *v.SanitizedCopy())
			}
			return
		}(),
		ConsumerGroups: ks.ConsumerGroups,
	}
}

func (ks *KongState) FillConsumersAndCredentials(
	s store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
	kongVersion semver.Version,
) {
	consumerIndex := make(map[string]Consumer)

	// build consumer index
	for _, consumer := range s.ListKongConsumers() {
		var c Consumer
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
				pushCredentialResourceFailures(fmt.Sprintf("failed to fetch secret: %v", err))
				continue
			}
			credConfig := map[string]interface{}{}
			for k, v := range secret.Data {
				// TODO populate these based on schema from Kong
				// and remove this workaround
				if k == "redirect_uris" {
					credConfig[k] = strings.Split(string(v), ",")
					continue
				}
				// TODO this is a kongCredType-agnostic mutation that should only apply to Oauth2 credentials.
				// However, the credential-specific code after deals only in interface{}s, and we can't fix individual
				// keys. To handle this properly we'd need to refactor the types used in all following code.
				if k == "hash_secret" {
					boolVal, err := strconv.ParseBool(string(v))
					if err != nil {
						// add a translation error here to tell that parsing hash_secret failed.
						pushCredentialResourceFailures(
							fmt.Sprintf("failed to parse hash_secret to bool: %v. defaulting to false", err),
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
							fmt.Sprintf("faield to parse ttl to int: %v, skipfilling the fiedl", err),
						)
					} else {
						credConfig[k] = intVal
					}
					continue
				}
				credConfig[k] = string(v)
			}
			credType, ok := credConfig["kongCredType"].(string)
			if !ok {
				pushCredentialResourceFailures(
					fmt.Sprintf("failed to provision credential: invalid kongCredType: type '%T' not string", credType),
				)
				continue
			}
			if !credentials.SupportedTypes.Has(credType) {
				pushCredentialResourceFailures(
					fmt.Sprintf("failed to provision credential: unsupported kongCredType: %q", credType),
				)
				continue
			}
			if len(credConfig) <= 1 { // 1 key of credType itself
				pushCredentialResourceFailures(
					"failed to provision credential: empty secret",
				)
				continue
			}
			credTags := util.GenerateTagsForObject(secret)
			err = c.SetCredential(credType, credConfig, credTags, kongVersion)
			if err != nil {
				pushCredentialResourceFailures(
					fmt.Sprintf("failed to provision credential: %v", err),
				)
				continue
			}
		}

		consumerIndex[consumer.Namespace+"/"+consumer.Name] = c
	}

	// populate the consumer in the state
	for _, c := range consumerIndex {
		ks.Consumers = append(ks.Consumers, c)
	}
}

func (ks *KongState) FillConsumerGroups(_ logrus.FieldLogger, s store.Storer) {
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

func (ks *KongState) FillOverrides(log logrus.FieldLogger, s store.Storer) {
	for i := 0; i < len(ks.Services); i++ {
		// Services
		kongIngress, err := getKongIngressForServices(s, ks.Services[i].K8sServices)
		if err != nil {
			log.WithError(err).
				Errorf("failed to fetch KongIngress resource for Services %s",
					PrettyPrintServiceList(ks.Services[i].K8sServices),
				)
			continue
		}

		for _, svc := range ks.Services[i].K8sServices {
			ks.Services[i].override(log, kongIngress, svc)
		}

		// Routes
		for j := 0; j < len(ks.Services[i].Routes); j++ {
			kongIngress, err := getKongIngressFromObjectMeta(s, ks.Services[i].Routes[j].Ingress)
			if err != nil {
				log.WithFields(logrus.Fields{
					"resource_name":      ks.Services[i].Routes[j].Ingress.Name,
					"resource_namespace": ks.Services[i].Routes[j].Ingress.Namespace,
				}).WithError(err).Errorf("failed to fetch KongIngress resource")
			}

			ks.Services[i].Routes[j].override(log, kongIngress)
		}
	}

	// Upstreams
	for i := 0; i < len(ks.Upstreams); i++ {
		kongIngress, err := getKongIngressForServices(s, ks.Upstreams[i].Service.K8sServices)
		if err != nil {
			log.WithError(err).
				Errorf("failed to fetch KongIngress resource for Services %s",
					PrettyPrintServiceList(ks.Upstreams[i].Service.K8sServices),
				)
			continue
		}

		for _, svc := range ks.Upstreams[i].Service.K8sServices {
			ks.Upstreams[i].override(kongIngress, svc)
		}
	}
}

func (ks *KongState) getPluginRelations() map[string]util.ForeignRelations {
	// KongPlugin key (KongPlugin's name:namespace) to corresponding associations
	pluginRels := map[string]util.ForeignRelations{}
	addConsumerRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = util.ForeignRelations{}
		}
		relations.Consumer = append(relations.Consumer, identifier)
		pluginRels[pluginKey] = relations
	}
	addConsumerGroupRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = util.ForeignRelations{}
		}
		relations.ConsumerGroup = append(relations.ConsumerGroup, identifier)
		pluginRels[pluginKey] = relations
	}
	addRouteRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = util.ForeignRelations{}
		}
		relations.Route = append(relations.Route, identifier)
		pluginRels[pluginKey] = relations
	}
	addServiceRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = util.ForeignRelations{}
		}
		relations.Service = append(relations.Service, identifier)
		pluginRels[pluginKey] = relations
	}

	for i := range ks.Services {
		// service
		for _, svc := range ks.Services[i].K8sServices {
			pluginList := annotations.ExtractKongPluginsFromAnnotations(svc.GetAnnotations())
			for _, pluginName := range pluginList {
				addServiceRelation(svc.Namespace, pluginName, *ks.Services[i].Name)
			}
		}
		// route
		for j := range ks.Services[i].Routes {
			ingress := ks.Services[i].Routes[j].Ingress
			pluginList := annotations.ExtractKongPluginsFromAnnotations(ingress.Annotations)
			for _, pluginName := range pluginList {
				addRouteRelation(ingress.Namespace, pluginName, *ks.Services[i].Routes[j].Name)
			}
		}
	}
	// consumer
	for _, c := range ks.Consumers {
		pluginList := annotations.ExtractKongPluginsFromAnnotations(c.K8sKongConsumer.GetAnnotations())
		for _, pluginName := range pluginList {
			addConsumerRelation(c.K8sKongConsumer.Namespace, pluginName, *c.Username)
		}
	}
	// consumer group
	for _, cg := range ks.ConsumerGroups {
		pluginList := annotations.ExtractKongPluginsFromAnnotations(cg.K8sKongConsumerGroup.GetAnnotations())
		for _, pluginName := range pluginList {
			addConsumerGroupRelation(cg.K8sKongConsumerGroup.Namespace, pluginName, *cg.Name)
		}
	}

	return pluginRels
}

func buildPlugins(
	log logrus.FieldLogger,
	s store.Storer,
	kongVersion semver.Version,
	failuresCollector *failures.ResourceFailuresCollector,
	pluginRels map[string]util.ForeignRelations,
) []Plugin {
	var plugins []Plugin

	for pluginIdentifier, relations := range pluginRels {
		identifier := strings.Split(pluginIdentifier, ":")
		namespace, kongPluginName := identifier[0], identifier[1]
		k8sPlugin, k8sClusterPlugin, err := getKongPluginOrKongClusterPlugin(s, namespace, kongPluginName)
		if err != nil {
			log.WithFields(logrus.Fields{
				"kongplugin_name":      kongPluginName,
				"kongplugin_namespace": namespace,
			}).WithError(err).Errorf("failed to fetch KongPlugin resource")
			continue
		}

		var plugin Plugin
		if k8sPlugin != nil {
			plugin, err = kongPluginFromK8SPlugin(s, *k8sPlugin, kongVersion)
			if err != nil {
				failuresCollector.PushResourceFailure(err.Error(), k8sPlugin)
				continue
			}
		}
		if k8sClusterPlugin != nil {
			plugin, err = kongPluginFromK8SClusterPlugin(s, *k8sClusterPlugin, kongVersion)
			if err != nil {
				failuresCollector.PushResourceFailure(err.Error(), k8sClusterPlugin)
				continue
			}
		}

		usedInstanceNames := sets.New[string]()
		for _, rel := range relations.GetCombinations() {
			plugin := plugin.DeepCopy()
			var sha [32]byte
			// ID is populated because that is read by decK and in_memory
			// translator too
			if rel.Service != "" {
				plugin.Service = &kong.Service{ID: kong.String(rel.Service)}
				sha = sha256.Sum256([]byte("service-" + rel.Service))
			}
			if rel.Route != "" {
				plugin.Route = &kong.Route{ID: kong.String(rel.Route)}
				sha = sha256.Sum256([]byte("route-" + rel.Route))
			}
			if rel.Consumer != "" {
				plugin.Consumer = &kong.Consumer{ID: kong.String(rel.Consumer)}
				sha = sha256.Sum256([]byte("consumer-" + rel.Consumer))
			}
			if rel.ConsumerGroup != "" {
				plugin.ConsumerGroup = &kong.ConsumerGroup{ID: kong.String(rel.ConsumerGroup)}
				sha = sha256.Sum256([]byte("group-" + rel.ConsumerGroup))
			}
			// instance_name must be unique. Using the same KongPlugin on multiple resources will result in duplicates
			// unless we add some sort of suffix.
			if plugin.InstanceName != nil {
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

	globalPlugins, err := globalPlugins(log, s, kongVersion)
	if err != nil {
		log.WithError(err).Error("failed to fetch global plugins")
	}
	// global plugins have no instance_name transform as they can only be applied once
	plugins = append(plugins, globalPlugins...)

	return plugins
}

func globalPlugins(log logrus.FieldLogger, s store.Storer, kongVersion semver.Version) ([]Plugin, error) {
	// removed as of 0.10.0
	// only retrieved now to warn users
	globalPlugins, err := s.ListGlobalKongPlugins()
	if err != nil {
		return nil, fmt.Errorf("error listing global KongPlugins: %w", err)
	}
	if len(globalPlugins) > 0 {
		log.Warning("global KongPlugins found. These are no longer applied and",
			" must be replaced with KongClusterPlugins.",
			" Please run \"kubectl get kongplugin -l global=true --all-namespaces\" to list existing plugins")
	}
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
	for i := 0; i < len(globalClusterPlugins); i++ {
		k8sPlugin := *globalClusterPlugins[i]
		pluginName := k8sPlugin.PluginName
		// empty pluginName skip it
		if pluginName == "" {
			log.WithFields(logrus.Fields{
				"kongclusterplugin_name": k8sPlugin.Name,
			}).Errorf("invalid KongClusterPlugin: empty plugin property")
			continue
		}
		if _, ok := res[pluginName]; ok {
			log.Error("multiple KongPlugin definitions found with"+
				" 'global' label for '", pluginName,
				"', the plugin will not be applied")
			duplicates = append(duplicates, pluginName)
			continue
		}
		if plugin, err := kongPluginFromK8SClusterPlugin(s, k8sPlugin, kongVersion); err == nil {
			res[pluginName] = plugin
		} else {
			log.WithFields(logrus.Fields{
				"kongclusterplugin_name": k8sPlugin.Name,
			}).WithError(err).Error("failed to generate configuration from KongClusterPlugin")
		}
	}
	for _, plugin := range duplicates {
		delete(res, plugin)
	}
	var plugins []Plugin
	for _, p := range res {
		plugins = append(plugins, p)
	}
	return plugins, nil
}

func (ks *KongState) FillPlugins(
	log logrus.FieldLogger,
	s store.Storer,
	kongVersion semver.Version,
	failuresCollector *failures.ResourceFailuresCollector,
) {
	ks.Plugins = buildPlugins(log, s, kongVersion, failuresCollector, ks.getPluginRelations())
}

// FillIDs iterates over the KongState and fills in the ID field for each entity
// that supports the FillID method (these are Service, Route, Consumer and Consumer
// Group). It makes their IDs deterministic, enabling their correct identification
// in external systems (e.g. Konnect Analytics).
func (ks *KongState) FillIDs(logger logrus.FieldLogger) {
	for svcIndex, svc := range ks.Services {
		if err := svc.FillID(""); err != nil {
			logger.WithError(err).Errorf("failed to fill ID for service %s", *svc.Name)
		} else {
			ks.Services[svcIndex] = svc
		}

		for routeIndex, route := range svc.Routes {
			if err := route.FillID(""); err != nil {
				logger.WithError(err).Errorf("failed to fill ID for route %s", *route.Name)
			} else {
				ks.Services[svcIndex].Routes[routeIndex] = route
			}
		}
	}

	for consumerIndex, consumer := range ks.Consumers {
		if err := consumer.FillID(""); err != nil {
			logger.WithError(err).Errorf("failed to fill ID for consumer %s", *consumer.Username)
		} else {
			ks.Consumers[consumerIndex] = consumer
		}
	}

	for consumerGroupIndex, consumerGroup := range ks.ConsumerGroups {
		if err := consumerGroup.FillID(""); err != nil {
			logger.WithError(err).Errorf("failed to fill ID for consumer group %s", *consumerGroup.Name)
		} else {
			ks.ConsumerGroups[consumerGroupIndex] = consumerGroup
		}
	}
}
