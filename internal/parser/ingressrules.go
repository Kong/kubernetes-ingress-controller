package parser

import (
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

type ingressRules struct {
	SecretNameToSNIs      SecretNameToSNIs
	ServiceNameToServices map[string]kongstate.Service
}

func newIngressRules() ingressRules {
	return ingressRules{
		SecretNameToSNIs:      newSecretNameToSNIs(),
		ServiceNameToServices: make(map[string]kongstate.Service),
	}
}

func mergeIngressRules(objs ...ingressRules) ingressRules {
	result := newIngressRules()

	for _, obj := range objs {
		for k, v := range obj.SecretNameToSNIs {
			result.SecretNameToSNIs[k] = append(result.SecretNameToSNIs[k], v...)
		}
		for k, v := range obj.ServiceNameToServices {
			result.ServiceNameToServices[k] = v
		}
	}
	return result
}

func (ir *ingressRules) populateServices(log logrus.FieldLogger, s store.Storer) {
	// populate Kubernetes Service
	for key, service := range ir.ServiceNameToServices {
		k8sSvc, err := s.GetService(service.Namespace, service.Backend.Name)
		if err != nil {
			log.WithFields(logrus.Fields{
				"service_name":      service.Backend.Name,
				"service_namespace": service.Namespace,
			}).Errorf("failed to fetch service: %v", err)
		}
		if k8sSvc != nil {
			service.K8sService = *k8sSvc
		}
		secretName := annotations.ExtractClientCertificate(
			service.K8sService.GetAnnotations())
		if secretName != "" {
			secret, err := s.GetSecret(service.K8sService.Namespace,
				secretName)
			secretKey := service.K8sService.Namespace + "/" + secretName
			// ensure that the cert is loaded into Kong
			if _, ok := ir.SecretNameToSNIs[secretKey]; !ok {
				ir.SecretNameToSNIs[secretKey] = []string{}
			}
			if err == nil {
				service.ClientCertificate = &kong.Certificate{
					ID: kong.String(string(secret.UID)),
				}
			} else {
				log.WithFields(logrus.Fields{
					"secret_name":      secretName,
					"secret_namespace": service.K8sService.Namespace,
				}).Errorf("failed to fetch secret: %v", err)
			}
		}
		ir.ServiceNameToServices[key] = service
	}
}

type SecretNameToSNIs map[string][]string

func newSecretNameToSNIs() SecretNameToSNIs {
	return SecretNameToSNIs(map[string][]string{})
}

func (m SecretNameToSNIs) addFromIngressV1beta1TLS(tlsSections []networkingv1beta1.IngressTLS, namespace string) {
	// Assume that v1beta1 and v1 tlsSections have identical semantics and field-wise content.
	var v1 []networkingv1.IngressTLS
	for _, item := range tlsSections {
		v1 = append(v1, networkingv1.IngressTLS{Hosts: item.Hosts, SecretName: item.SecretName})
	}
	m.addFromIngressV1TLS(v1, namespace)
}

func (m SecretNameToSNIs) addFromIngressV1TLS(tlsSections []networkingv1.IngressTLS, namespace string) {
	for _, tls := range tlsSections {
		if len(tls.Hosts) == 0 {
			continue
		}
		if tls.SecretName == "" {
			continue
		}
		hosts := tls.Hosts
		secretName := namespace + "/" + tls.SecretName
		hosts = m.filterHosts(hosts)
		if m[secretName] != nil {
			hosts = append(hosts, m[secretName]...)
		}
		m[secretName] = hosts
	}
}

func (m SecretNameToSNIs) filterHosts(hosts []string) []string {
	hostsToAdd := []string{}
	seenHosts := map[string]bool{}
	for _, hosts := range m {
		for _, host := range hosts {
			seenHosts[host] = true
		}
	}
	for _, host := range hosts {
		if !seenHosts[host] {
			hostsToAdd = append(hostsToAdd, host)
		}
	}
	return hostsToAdd
}
