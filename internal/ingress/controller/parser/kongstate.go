package parser

import (
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser/consumer"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/sirupsen/logrus"
)

// KongState holds the configuration that should be applied to Kong.
type KongState struct {
	Services       []Service
	Upstreams      []Upstream
	Certificates   []Certificate
	CACertificates []kong.CACertificate
	Plugins        []Plugin
	Consumers      []consumer.Consumer
}

func (ks *KongState) fillConsumersAndCredentials(log logrus.FieldLogger, s store.Storer) {
	consumerIndex := make(map[string]consumer.Consumer)

	// build consumer index
	for _, kConsumer := range s.ListKongConsumers() {
		var c consumer.Consumer
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

	// legacy attach credentials
	credentials := s.ListKongCredentials()
	if len(credentials) > 0 {
		log.Warnf("deprecated KongCredential resource in use; " +
			"please use secret-based credentials, " +
			"KongCredential resource will be removed in future")
	}
	for _, credential := range credentials {
		log = log.WithFields(logrus.Fields{
			"kongcredential_name":      credential.Name,
			"kongcredential_namespace": credential.Namespace,
			"consumerRef":              credential.ConsumerRef,
		})
		cons, ok := consumerIndex[credential.Namespace+"/"+
			credential.ConsumerRef]
		if !ok {
			continue
		}
		if credential.Type == "" {
			log.Errorf("invalid KongCredential: no Type provided")
			continue
		}
		if !supportedCreds.Has(credential.Type) {
			log.Errorf("invalid KongCredential: invalid Type provided")
			continue
		}
		if credential.Config == nil {
			log.Errorf("invalid KongCredential: empty config")
			continue
		}
		err := cons.SetCredential(log, credential.Type, credential.Config)
		if err != nil {
			log.Errorf("failed to provision credential: %v", err)
			continue
		}
		consumerIndex[credential.Namespace+"/"+credential.ConsumerRef] = cons
	}

	// populate the consumer in the state
	for _, c := range consumerIndex {
		ks.Consumers = append(ks.Consumers, c)
	}
}

func (ks *KongState) fillOverrides(log logrus.FieldLogger, s store.Storer) {
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
		overrideService(&ks.Services[i], kongIngress, anns)

		// Routes
		for j := 0; j < len(ks.Services[i].Routes); j++ {
			var kongIngress *configurationv1.KongIngress
			var err error
			if ks.Services[i].Routes[j].IsTCP {
				kongIngress, err = getKongIngressFromTCPIngress(s,
					&ks.Services[i].Routes[j].TCPIngress)
				if err != nil {
					log.WithFields(logrus.Fields{
						"tcpingress_name":      ks.Services[i].Routes[j].TCPIngress.Name,
						"tcpingress_namespace": ks.Services[i].Routes[j].TCPIngress.Namespace,
					}).Errorf("failed to fetch KongIngress resource for Ingress: %v", err)
				}
			} else {
				kongIngress, err = getKongIngressFromIngress(s,
					&ks.Services[i].Routes[j].Ingress)
				if err != nil {
					log.WithFields(logrus.Fields{
						"ingress_name":      ks.Services[i].Routes[j].Ingress.Name,
						"ingress_namespace": ks.Services[i].Routes[j].Ingress.Namespace,
					}).Errorf("failed to fetch KongIngress resource for Ingress: %v", err)
				}
			}

			overrideRoute(log, &ks.Services[i].Routes[j], kongIngress)
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
		overrideUpstream(&ks.Upstreams[i], kongIngress, anns)
	}
}

func getPluginRelations(state KongState) map[string]foreignRelations {
	// KongPlugin key (KongPlugin's name:namespace) to corresponding associations
	pluginRels := map[string]foreignRelations{}
	addConsumerRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = foreignRelations{}
		}
		relations.Consumer = append(relations.Consumer, identifier)
		pluginRels[pluginKey] = relations
	}
	addRouteRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = foreignRelations{}
		}
		relations.Route = append(relations.Route, identifier)
		pluginRels[pluginKey] = relations
	}
	addServiceRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = foreignRelations{}
		}
		relations.Service = append(relations.Service, identifier)
		pluginRels[pluginKey] = relations
	}

	for i := range state.Services {
		// service
		svc := state.Services[i].K8sService
		pluginList := annotations.ExtractKongPluginsFromAnnotations(
			svc.GetAnnotations())
		for _, pluginName := range pluginList {
			addServiceRelation(svc.Namespace, pluginName,
				*state.Services[i].Name)
		}
		// route
		for j := range state.Services[i].Routes {
			ingress := state.Services[i].Routes[j].Ingress
			pluginList := annotations.ExtractKongPluginsFromAnnotations(ingress.GetAnnotations())
			for _, pluginName := range pluginList {
				addRouteRelation(ingress.Namespace, pluginName, *state.Services[i].Routes[j].Name)
			}
		}
	}
	// consumer
	for _, c := range state.Consumers {
		pluginList := annotations.ExtractKongPluginsFromAnnotations(c.K8sKongConsumer.GetAnnotations())
		for _, pluginName := range pluginList {
			addConsumerRelation(c.K8sKongConsumer.Namespace, pluginName, *c.Username)
		}
	}
	return pluginRels
}
