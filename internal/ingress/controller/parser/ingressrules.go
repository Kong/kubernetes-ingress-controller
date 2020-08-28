package parser

import (
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	"github.com/sirupsen/logrus"
)

type ingressRules struct {
	SecretNameToSNIs      map[string][]string
	ServiceNameToServices map[string]Service
}

func newIngressRules() ingressRules {
	return ingressRules{
		SecretNameToSNIs:      make(map[string][]string),
		ServiceNameToServices: make(map[string]Service),
	}
}

func mergeIngressRules(objs ...*ingressRules) ingressRules {
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
