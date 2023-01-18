package parser

import (
	"errors"
	"fmt"
	"sort"

	"github.com/kong/go-kong/kong"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func (p *Parser) ingressRulesFromKnativeIngress() ingressRules {
	result := newIngressRules()

	// IngressClass is not actually part of the Knative spec, and we are getting networking.k8s.io IngressClasses here,
	// not a resource specific to Knative. However, the reason we're using it (enabling the 2.x regex heuristic) is
	// Kong-specific, so in absence of a proper Knative IngressClass to attach our IngressClassParams to, we may as
	// well use the stock Kubernetes resource.
	icp, err := getIngressClassParametersOrDefault(p.storer)
	if err != nil {
		if !errors.As(err, &store.ErrNotFound{}) {
			// anything else is unexpected
			p.logger.Errorf("could not find IngressClassParameters, using defaults: %s", err)
		}
	}
	ingressList, err := p.storer.ListKnativeIngresses()
	if err != nil {
		p.logger.WithError(err).Error("failed to list Knative Ingresses")
		return result
	}

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	services := map[string]kongstate.Service{}
	secretToSNIs := newSecretNameToSNIs()

	for _, ingress := range ingressList {
		regexPrefix := translators.ControllerPathRegexPrefix
		if prefix, ok := ingress.ObjectMeta.Annotations[annotations.AnnotationPrefix+annotations.RegexPrefixKey]; ok {
			regexPrefix = prefix
		}
		ingressSpec := ingress.Spec

		secretToSNIs.addFromIngressV1TLS(knativeIngressToNetworkingTLS(ingress.Spec.TLS), ingress)

		var objectSuccessfullyParsed bool
		for i, rule := range ingressSpec.Rules {
			hosts := rule.Hosts
			if rule.HTTP == nil {
				continue
			}
			for j, rule := range rule.HTTP.Paths {
				path := rule.Path

				path = maybePrependRegexPrefix(path, regexPrefix, icp.EnableLegacyRegexDetection && p.flagEnabledRegexPathPrefix)
				if path == "" {
					path = "/"
				}
				r := kongstate.Route{
					Ingress: util.FromK8sObject(ingress),
					Route: kong.Route{
						Name:              kong.String(fmt.Sprintf("%s.%s.%d%d", ingress.Namespace, ingress.Name, i, j)),
						Paths:             kong.StringSlice(path),
						StripPath:         kong.Bool(false),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						RequestBuffering:  kong.Bool(true),
						ResponseBuffering: kong.Bool(true),
					},
				}
				r.Hosts = kong.StringSlice(hosts...)

				knativeBackend := knativeSelectSplit(rule.Splits)
				serviceName := fmt.Sprintf("%s.%s.%s", knativeBackend.ServiceNamespace, knativeBackend.ServiceName,
					knativeBackend.ServicePort.String())
				serviceHost := fmt.Sprintf("%s.%s.%s.svc", knativeBackend.ServiceName, knativeBackend.ServiceNamespace,
					knativeBackend.ServicePort.String())
				service, ok := services[serviceName]
				if !ok {
					var headers []string
					for key, value := range knativeBackend.AppendHeaders {
						headers = append(headers, key+":"+value)
					}
					for key, value := range rule.AppendHeaders {
						headers = append(headers, key+":"+value)
					}

					service = kongstate.Service{
						Service: kong.Service{
							Name:           kong.String(serviceName),
							Host:           kong.String(serviceHost),
							Port:           kong.Int(DefaultHTTPPort),
							Protocol:       kong.String("http"),
							Path:           kong.String("/"),
							ConnectTimeout: kong.Int(DefaultServiceTimeout),
							ReadTimeout:    kong.Int(DefaultServiceTimeout),
							WriteTimeout:   kong.Int(DefaultServiceTimeout),
							Retries:        kong.Int(DefaultRetries),
						},
						Namespace: ingress.Namespace,
						Backends: []kongstate.ServiceBackend{{
							Name:    knativeBackend.ServiceName,
							PortDef: translators.PortDefFromIntStr(knativeBackend.ServicePort),
						}},
						Parent: ingress,
					}
					if len(headers) > 0 {
						service.Plugins = append(service.Plugins, kong.Plugin{
							Name: kong.String("request-transformer"),
							Config: kong.Configuration{
								"add": map[string]interface{}{
									"headers": headers,
								},
							},
						})
					}
				}
				service.Routes = append(service.Routes, r)
				services[serviceName] = service
				objectSuccessfullyParsed = true
			}
		}

		if objectSuccessfullyParsed {
			p.ReportKubernetesObjectUpdate(ingress)
		}
	}

	result.ServiceNameToServices = services
	result.SecretNameToSNIs = secretToSNIs
	return result
}

func knativeSelectSplit(splits []knative.IngressBackendSplit) knative.IngressBackendSplit {
	if len(splits) == 0 {
		return knative.IngressBackendSplit{}
	}
	res := splits[0]
	maxPercentage := splits[0].Percent
	if len(splits) == 1 {
		return res
	}
	for i := 1; i < len(splits); i++ {
		if splits[i].Percent > maxPercentage {
			res = splits[i]
			maxPercentage = res.Percent
		}
	}
	return res
}
