package parser

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func (p *Parser) ingressRulesFromTCPIngressV1beta1() ingressRules {
	result := newIngressRules()

	ingressList, err := p.storer.ListTCPIngresses()
	if err != nil {
		p.logger.WithError(err).Error("failed to list TCPIngresses")
		return result
	}

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		ingressSpec := ingress.Spec

		log := p.logger.WithFields(logrus.Fields{
			"tcpingress_namespace": ingress.Namespace,
			"tcpingress_name":      ingress.Name,
		})

		result.SecretNameToSNIs.addFromIngressV1beta1TLS(tcpIngressToNetworkingTLS(ingressSpec.TLS), ingress.Namespace)

		var objectSuccessfullyParsed bool
		for i, rule := range ingressSpec.Rules {
			if !util.IsValidPort(rule.Port) {
				log.Errorf("invalid TCPIngress: invalid port: %v", rule.Port)
				continue
			}
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
				},
			}
			host := rule.Host
			if host != "" {
				r.SNIs = kong.StringSlice(host)
			}
			if rule.Backend.ServiceName == "" {
				log.Errorf("invalid TCPIngress: empty serviceName")
				continue
			}
			if !util.IsValidPort(rule.Backend.ServicePort) {
				log.Errorf("invalid TCPIngress: invalid servicePort: %v", rule.Backend.ServicePort)
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
					Backends: []kongstate.ServiceBackend{{
						Name:    rule.Backend.ServiceName,
						PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(rule.Backend.ServicePort)},
					}},
				}
			}
			service.Routes = append(service.Routes, r)
			result.ServiceNameToServices[serviceName] = service
			objectSuccessfullyParsed = true
		}

		if objectSuccessfullyParsed {
			p.ReportKubernetesObjectUpdate(ingress)
		}
	}

	return result
}

func (p *Parser) ingressRulesFromUDPIngressV1beta1() ingressRules {
	result := newIngressRules()

	ingressList, err := p.storer.ListUDPIngresses()
	if err != nil {
		p.logger.WithError(err).Errorf("failed to list UDPIngresses")
		return result
	}

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(&ingressList[j].CreationTimestamp)
	})

	for _, ingress := range ingressList {
		ingressSpec := ingress.Spec

		log := p.logger.WithFields(logrus.Fields{
			"udpingress_namespace": ingress.Namespace,
			"udpingress_name":      ingress.Name,
		})

		var objectSuccessfullyParsed bool
		for i, rule := range ingressSpec.Rules {
			// validate the ports and servicenames for the rule
			if !util.IsValidPort(rule.Port) {
				log.Errorf("invalid UDPIngress: invalid port: %d", rule.Port)
				continue
			}
			if rule.Backend.ServiceName == "" {
				log.Errorf("invalid UDPIngress: empty serviceName")
				continue
			}
			if !util.IsValidPort(rule.Backend.ServicePort) {
				log.Errorf("invalid UDPIngress: invalid servicePort: %d", rule.Backend.ServicePort)
				continue
			}

			// generate the kong Route based on the listen port
			route := kongstate.Route{
				Ingress: util.FromK8sObject(ingress),
				Route: kong.Route{
					Name:         kong.String(ingress.Namespace + "." + ingress.Name + "." + strconv.Itoa(i) + ".udp"),
					Protocols:    kong.StringSlice("udp"),
					Destinations: []*kong.CIDRPort{{Port: kong.Int(rule.Port)}},
				},
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
					Backends: []kongstate.ServiceBackend{{
						Name:    rule.Backend.ServiceName,
						PortDef: kongstate.PortDef{Mode: kongstate.PortModeByNumber, Number: int32(rule.Backend.ServicePort)},
					}},
				}
			}
			service.Routes = append(service.Routes, route)
			result.ServiceNameToServices[serviceName] = service
			objectSuccessfullyParsed = true
		}

		if objectSuccessfullyParsed {
			p.ReportKubernetesObjectUpdate(ingress)
		}
	}

	return result
}
