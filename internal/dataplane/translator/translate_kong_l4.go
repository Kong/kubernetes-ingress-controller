package translator

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/kong/go-kong/kong"
	netv1 "k8s.io/api/networking/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

func (t *Translator) ingressRulesFromTCPIngressV1beta1() ingressRules {
	result := newIngressRules()

	ingressList, err := t.storer.ListTCPIngresses()
	if err != nil {
		t.logger.Error(err, "Failed to list TCPIngresses")
		return result
	}

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		result.SecretNameToSNIs.addFromIngressV1TLS(tcpIngressToNetworkingTLS(ingress.Spec.TLS), ingress)

		var objectSuccessfullyTranslated bool
		for i, rule := range ingress.Spec.Rules {
			r := kongstate.Route{
				Ingress: util.FromK8sObject(ingress),
				Route: kong.Route{
					Name:      kong.String(ingress.Namespace + "." + ingress.Name + "." + strconv.Itoa(i)),
					Protocols: kong.StringSlice("tcp", "tls"),
					Destinations: []*kong.CIDRPort{
						{
							Port: kong.Int(rule.Port),
						},
					},
					Tags: util.GenerateTagsForObject(ingress),
				},
			}
			if host := rule.Host; host != "" {
				r.SNIs = kong.StringSlice(host)
			}

			serviceBackend, err := kongstate.NewServiceBackendForService(
				k8stypes.NamespacedName{
					Namespace: ingress.Namespace,
					Name:      rule.Backend.ServiceName,
				},
				kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(rule.Backend.ServicePort)},
			)
			if err != nil {
				t.logger.Error(err, "failed to create ServiceBackend for TCPIngress rule",
					"ingress_name", ingress.Name,
					"ingress_namespace", ingress.Namespace,
					"rule_idx", i,
				)
				continue
			}

			serviceName := fmt.Sprintf("%s.%s.%d", ingress.Namespace, rule.Backend.ServiceName, rule.Backend.ServicePort)
			service, ok := result.ServiceNameToServices[serviceName]
			if !ok {
				service = kongstate.Service{
					Service: kong.Service{
						Name: kong.String(serviceName),
						Host: kong.String(fmt.Sprintf("%s.%s.%d.svc", rule.Backend.ServiceName, ingress.Namespace,
							rule.Backend.ServicePort)),
						Port:           kong.Int(DefaultHTTPPort),
						Protocol:       kong.String("tcp"),
						ConnectTimeout: kong.Int(DefaultServiceTimeout),
						ReadTimeout:    kong.Int(DefaultServiceTimeout),
						WriteTimeout:   kong.Int(DefaultServiceTimeout),
						Retries:        kong.Int(DefaultRetries),
					},
					Namespace: ingress.Namespace,
					Backends:  []kongstate.ServiceBackend{serviceBackend},
					Parent:    ingress,
				}
			}
			service.Routes = append(service.Routes, r)
			result.ServiceNameToServices[serviceName] = service
			result.ServiceNameToParent[serviceName] = ingress
			objectSuccessfullyTranslated = true
		}

		if objectSuccessfullyTranslated {
			t.registerSuccessfullyTranslatedObject(ingress)
		}
	}

	if t.featureFlags.ExpressionRoutes {
		applyExpressionToIngressRules(&result)
	}

	return result
}

func (t *Translator) ingressRulesFromUDPIngressV1beta1() ingressRules {
	result := newIngressRules()

	ingressList, err := t.storer.ListUDPIngresses()
	if err != nil {
		t.logger.Error(err, "Failed to list UDPIngresses")
		return result
	}

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		var objectSuccessfullyTranslated bool
		for i, rule := range ingress.Spec.Rules {
			// generate the kong Route based on the listen port
			route := kongstate.Route{
				Ingress: util.FromK8sObject(ingress),
				Route: kong.Route{
					Name:         kong.String(ingress.Namespace + "." + ingress.Name + "." + strconv.Itoa(i) + ".udp"),
					Protocols:    kong.StringSlice("udp"),
					Destinations: []*kong.CIDRPort{{Port: kong.Int(rule.Port)}},
					Tags:         util.GenerateTagsForObject(ingress),
				},
			}

			serviceBackend, err := kongstate.NewServiceBackendForService(
				k8stypes.NamespacedName{
					Namespace: ingress.Namespace,
					Name:      rule.Backend.ServiceName,
				},
				kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(rule.Backend.ServicePort)},
			)
			if err != nil {
				t.logger.Error(err, "failed to create ServiceBackend for UDPIngress rule",
					"ingress_name", ingress.Name,
					"ingress_namespace", ingress.Namespace,
					"rule_idx", i,
				)
				continue
			}

			// generate the kong Service backend for the UDPIngress rules
			host := fmt.Sprintf("%s.%s.%d.svc", rule.Backend.ServiceName, ingress.Namespace, rule.Backend.ServicePort)
			serviceName := fmt.Sprintf("%s.%s.%d.udp", ingress.Namespace, rule.Backend.ServiceName, rule.Backend.ServicePort)
			service, ok := result.ServiceNameToServices[serviceName]
			if !ok {
				service = kongstate.Service{
					Namespace: ingress.Namespace,
					Service: kong.Service{
						Name:     kong.String(serviceName),
						Protocol: kong.String("udp"),
						Host:     kong.String(host),
						Port:     kong.Int(rule.Backend.ServicePort),
					},
					Backends: []kongstate.ServiceBackend{serviceBackend},
					Parent:   ingress,
				}
			}
			service.Routes = append(service.Routes, route)
			result.ServiceNameToServices[serviceName] = service
			result.ServiceNameToParent[serviceName] = ingress
			objectSuccessfullyTranslated = true
		}

		if objectSuccessfullyTranslated {
			t.registerSuccessfullyTranslatedObject(ingress)
		}
	}

	if t.featureFlags.ExpressionRoutes {
		applyExpressionToIngressRules(&result)
	}

	return result
}

func tcpIngressToNetworkingTLS(tls []kongv1beta1.IngressTLS) []netv1.IngressTLS {
	var result []netv1.IngressTLS

	for _, t := range tls {
		result = append(result, netv1.IngressTLS{
			Hosts:      t.Hosts,
			SecretName: t.SecretName,
		})
	}
	return result
}
