package kongstate

import (
	"fmt"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

// KongState holds the configuration that should be applied to Kong.
type KongState struct {
	Services       []Service
	Upstreams      []Upstream
	Certificates   []Certificate
	CACertificates []kong.CACertificate
	Plugins        []Plugin
	Consumers      []Consumer
}

func (ks *KongState) FillConsumersAndCredentials(log logrus.FieldLogger, s store.Storer) {
	consumerIndex := make(map[string]Consumer)

	// build consumer index
	for _, kConsumer := range s.ListKongConsumers() {
		var c Consumer
		if kConsumer.Username == "" && kConsumer.CustomID == "" {
			continue
		}
		if kConsumer.Username != "" {
			c.Username = kong.String(kConsumer.Username)
		}
		if kConsumer.CustomID != "" {
			c.CustomID = kong.String(kConsumer.CustomID)
		}
		c.K8sKongConsumer = *kConsumer

		log = log.WithFields(logrus.Fields{
			"kongconsumer_name":      kConsumer.Name,
			"kongconsumer_namespace": kConsumer.Namespace,
		})
		for _, cred := range kConsumer.Credentials {
			log = log.WithFields(logrus.Fields{
				"secret_name":      cred,
				"secret_namespace": kConsumer.Namespace,
			})
			secret, err := s.GetSecret(kConsumer.Namespace, cred)
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
			credType, ok := credConfig["kongCredType"].(string)
			if !ok {
				log.Errorf("failed to provision credential: invalid credType: %v", credType)
			}
			if !supportedCreds.Has(credType) {
				log.Errorf("failed to provision credential: invalid credType: %v", credType)
				continue
			}
			if len(credConfig) <= 1 { // 1 key of credType itself
				log.Errorf("failed to provision credential: empty secret")
				continue
			}
			err = c.SetCredential(log, credType, credConfig)
			if err != nil {
				log.Errorf("failed to provision credential: %v", err)
				continue
			}
		}

		consumerIndex[kConsumer.Namespace+"/"+kConsumer.Name] = c
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

var supportedCreds = sets.NewString(
	"acl",
	"basic-auth",
	"hmac-auth",
	"jwt",
	"key-auth",
	"oauth2",
)
