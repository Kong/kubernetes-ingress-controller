package kongstate

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kong/kubernetes-ingress-controller/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/store"
	"github.com/kong/kubernetes-ingress-controller/internal/util"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
)

// KongState holds the configuration that should be applied to Kong.
type KongState struct {
	Services       []Service
	Upstreams      []Upstream
	Certificates   []Certificate
	CACertificates []kong.CACertificate
	Plugins        []Plugin
	Consumers      []Consumer
	Version        semver.Version
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
	}
}

func (ks *KongState) FillConsumersAndCredentials(log logrus.FieldLogger, s store.Storer) {
	consumerIndex := make(map[string]Consumer)

	// TODO: jwts in consumer credentials are GLOBALLY unique, need to deduplicate these
	//       for all consumer here
	// jwtDedup := make(map[string]credsValidationMeta)

	// build consumer index
	for _, consumer := range s.ListKongConsumers() {
		var c Consumer
		if consumer.Username == "" && consumer.CustomID == "" {
			continue
		}
		if consumer.Username != "" {
			c.Username = kong.String(consumer.Username)
		}
		if consumer.CustomID != "" {
			c.CustomID = kong.String(consumer.CustomID)
		}
		c.K8sKongConsumer = *consumer

		log = log.WithFields(logrus.Fields{
			"kongconsumer_name":      consumer.Name,
			"kongconsumer_namespace": consumer.Namespace,
		})

		// ---------------------------------------------------------------------------
		// KongConsumer Credentials - Processing & Validation
		// ---------------------------------------------------------------------------

		credsDedup := make(map[string][]credsValidationMeta)
		for _, cred := range consumer.Credentials {
			log = log.WithFields(logrus.Fields{
				"secret_name":      cred,
				"secret_namespace": consumer.Namespace,
			})

			// --------------------------------------------------------------------------
			// Credentials Retrieval
			// --------------------------------------------------------------------------

			secret, err := s.GetSecret(consumer.Namespace, cred)
			if err != nil {
				log.Errorf("failed to fetch secret: %v", err)
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
				credConfig[k] = string(v)
			}

			// --------------------------------------------------------------------------
			// Credentials Validation
			// --------------------------------------------------------------------------

			// if any validation fails the credentials will be marked as invalid and not
			// applied, however validation continues after each failure so as to get the
			// benefit of reporting ALL validation issues that were found with this
			// credential, so the user can make all the required changes and not update
			// one only to find there were more the whole time.
			invalid := false

			// if any failure occur when validating the credentials, the error will contain
			// this prefix in the message body to make it easier to search for this class
			// of errors (KongConsumer credentials validation errors) in the manager logs.
			failHeader := fmt.Sprintf("failed to provision credentials for consumer %s", consumer.Name)

			// validate that the required "kongCredType" key is present in the creds.
			credType, ok := credConfig[credTypeKey].(string)
			if !ok {
				log.Errorf("%s: secret %s provided no %s", failHeader, secret.Name, credTypeKey)
				invalid = true
			}
			if !supportedCreds.Has(credType) {
				supportedCredsString := fmt.Sprintf("valid options are: %s", strings.Join(supportedCreds.List(), ","))
				log.Errorf("%s: invalid credType: %s (%s)", failHeader, credType, supportedCredsString)
				invalid = true
			}

			// validate that the credentials actually includes configuration, not just the type
			if len(credConfig) <= 1 {
				log.Errorf("%s: secret %s has no data", failHeader, secret.Name)
				continue // if this is true, there's nothing left to validate anyway
			}

			// consumer credentials can technically be configured with two (different)
			// secrets that contain the same key. For some types there are unique
			// contraints on the key. Here we validate the unique constraints for
			// several types.
			for k, v := range credConfig {
				if k == credTypeKey {
					continue
				}

				// TODO: filter out non-unique keys, such as "username" for "basic-auth"
				//       See: https://github.com/kong/kong/blob/master/kong/plugins/hmac-auth/daos.lua

				// validate the integrity of the secret key value data
				value, ok := v.([]byte)
				if !ok {
					log.Errorf("%s: invalid credential %s: key value can't be %T", failHeader, consumer.Name, cred, v)
					invalid = true
				}

				// check if the key has ever been seen before, if it has we need to
				// process it to ensure deduplication.
				//
				// if all of the following are true (inclusively) then the key is
				// invalid and the credentials will be dropped:
				//
				//  - if the credential type is the same
				//  - if the key has a unique constraint
				//  - if the data is duplicated
				if metas, ok := credsDedup[k]; ok {
					for _, meta := range metas {
						if meta.credType == credType && bytes.Equal(meta.value, value) && supportedCredsWithUniqueConstraints.Has(credType) {
							log.Errorf("%s: unique constaint violated: more than 1 credential secret has defined %s", failHeader, k)
							invalid = true
						}
					}
				}

				// report to the deduplicator that we've seen this key now so that
				// when future credentials are processed this validation logic can check.
				credsDedup[k] = append(credsDedup[k], credsValidationMeta{value: value, credType: credType})
			}

			// --------------------------------------------------------------------------
			// Credentials Updates
			// --------------------------------------------------------------------------

			if !invalid {
				err = c.SetCredential(credType, credConfig, ks.Version)
				if err != nil {
					log.Errorf("%s: %w", failHeader, err)
					continue
				}
			}
		}

		consumerIndex[consumer.Namespace+"/"+consumer.Name] = c
	}

	// populate the consumer in the state
	for _, c := range consumerIndex {
		ks.Consumers = append(ks.Consumers, c)
	}
}

func (ks *KongState) FillOverrides(log logrus.FieldLogger, s store.Storer) {
	for i := 0; i < len(ks.Services); i++ {
		// Services
		anns := ks.Services[i].K8sService.Annotations
		kongIngress, err := getKongIngressForService(s, ks.Services[i].K8sService)
		if err != nil {
			log.WithFields(logrus.Fields{
				"service_name":      ks.Services[i].K8sService.Name,
				"service_namespace": ks.Services[i].K8sService.Namespace,
			}).Errorf("failed to fetch KongIngress resource for Service: %v", err)
		}
		ks.Services[i].override(kongIngress, anns)

		// Routes
		for j := 0; j < len(ks.Services[i].Routes); j++ {
			var kongIngress *configurationv1.KongIngress
			var err error
			kongIngress, err = getKongIngressFromObjectMeta(s, &ks.Services[i].Routes[j].Ingress)
			if err != nil {
				log.WithFields(logrus.Fields{
					"resource_name":      ks.Services[i].Routes[j].Ingress.Name,
					"resource_namespace": ks.Services[i].Routes[j].Ingress.Namespace,
				}).Errorf("failed to fetch KongIngress resource: %v", err)
			}

			ks.Services[i].Routes[j].override(log, kongIngress)
		}
	}

	// Upstreams
	for i := 0; i < len(ks.Upstreams); i++ {
		kongIngress, err := getKongIngressForService(s,
			ks.Upstreams[i].Service.K8sService)
		anns := ks.Upstreams[i].Service.K8sService.Annotations
		if err != nil {
			log.WithFields(logrus.Fields{
				"service_name":      ks.Upstreams[i].Service.K8sService.Name,
				"service_namespace": ks.Upstreams[i].Service.K8sService.Namespace,
			}).Errorf("failed to fetch KongIngress resource for Service: %v", err)
			continue
		}
		ks.Upstreams[i].override(kongIngress, anns)
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
		svc := ks.Services[i].K8sService
		pluginList := annotations.ExtractKongPluginsFromAnnotations(
			svc.GetAnnotations())
		for _, pluginName := range pluginList {
			addServiceRelation(svc.Namespace, pluginName,
				*ks.Services[i].Name)
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
	return pluginRels
}

func buildPlugins(log logrus.FieldLogger, s store.Storer, pluginRels map[string]util.ForeignRelations) []Plugin {
	var plugins []Plugin

	for pluginIdentifier, relations := range pluginRels {
		identifier := strings.Split(pluginIdentifier, ":")
		namespace, kongPluginName := identifier[0], identifier[1]
		plugin, err := getPlugin(s, namespace, kongPluginName)
		if err != nil {
			log.WithFields(logrus.Fields{
				"kongplugin_name":      kongPluginName,
				"kongplugin_namespace": namespace,
			}).Errorf("failed to fetch KongPlugin: %v", err)
			continue
		}

		for _, rel := range relations.GetCombinations() {
			plugin := *plugin.DeepCopy()
			// ID is populated because that is read by decK and in_memory
			// translator too
			if rel.Service != "" {
				plugin.Service = &kong.Service{ID: kong.String(rel.Service)}
			}
			if rel.Route != "" {
				plugin.Route = &kong.Route{ID: kong.String(rel.Route)}
			}
			if rel.Consumer != "" {
				plugin.Consumer = &kong.Consumer{ID: kong.String(rel.Consumer)}
			}
			plugins = append(plugins, Plugin{plugin})
		}
	}

	globalPlugins, err := globalPlugins(log, s)
	if err != nil {
		log.Errorf("failed to fetch global plugins: %v", err)
	}
	plugins = append(plugins, globalPlugins...)

	return plugins
}

func globalPlugins(log logrus.FieldLogger, s store.Storer) ([]Plugin, error) {
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
		if plugin, err := kongPluginFromK8SClusterPlugin(s, k8sPlugin); err == nil {
			res[pluginName] = Plugin{
				Plugin: plugin,
			}
		} else {
			log.WithFields(logrus.Fields{
				"kongclusterplugin_name": k8sPlugin.Name,
			}).Errorf("failed to generate configuration from KongClusterPlugin: %v ", err)
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

func (ks *KongState) FillPlugins(log logrus.FieldLogger, s store.Storer) {
	ks.Plugins = buildPlugins(log, s, ks.getPluginRelations())
}

// -----------------------------------------------------------------------------
// KongConsumer - Credentials Helper Vars & Functions
// -----------------------------------------------------------------------------

var (
	// supportedCredsWithUniqueConstraints indicates all the "kongCredType"s which are supported for
	// KongConsumer credentials and have a unique constraint (multiple credentials can not define
	// these keys more than once for a single KongConsumer resource)
	supportedCredsWithUniqueConstraints = sets.NewString(
		"basic-auth",
		"hmac-auth",
		"jwt",
		"key-auth",
		"oauth2",
	)

	// supportedCredsWithoutUniqueConstraints indicates all the "kongCredType"s which are supported for
	// KongConsumer credentials and also do not have a unique constraint (multiple credentials can define
	// these keys more than once for a single KongConsumer resource)
	supportedCredsWithoutUniqueConstraints = sets.NewString(
		"acl",
	)

	// supportedCreds indicates all the "kongCredType"s which are supported for KongConsumer credentials.
	supportedCreds = supportedCredsWithUniqueConstraints.Union(supportedCredsWithoutUniqueConstraints)
)

const (
	// credTypeKey indicates the key in a consumer secret which identifies the type
	// of credential that is being provided for the consumer.
	credTypeKey = "kongCredType"

	// credPrimaryKey indicates the key "key" in a consumer secret for credentials.
	credPrimaryKey = "key"
)

// credsValidationMeta is a metadata struct to help validate the contents of
// consumer credentials, particularly unique constraints on the underlying data.
type credsValidationMeta struct {
	value    []byte
	credType string
}
