package manager

import (
	"fmt"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/ctrlutils"
)

type IngressAPI int

const (
	OtherAPI          IngressAPI = iota
	NetworkingV1      IngressAPI = iota
	NetworkingV1beta1 IngressAPI = iota
	ExtensionsV1beta1 IngressAPI = iota
)

func negotiateIngressAPI(config *Config, client client.Client) (IngressAPI, error) {
	var allowedAPIs []IngressAPI
	candidateAPIs := map[IngressAPI]schema.GroupVersionResource{
		NetworkingV1: {
			Group:    networkingv1.SchemeGroupVersion.Group,
			Version:  networkingv1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		},
		NetworkingV1beta1: {
			Group:    networkingv1beta1.SchemeGroupVersion.Group,
			Version:  networkingv1beta1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		},
		ExtensionsV1beta1: {
			Group:    extensionsv1beta1.SchemeGroupVersion.Group,
			Version:  extensionsv1beta1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		},
	}

	if config.IngressNetV1Enabled {
		allowedAPIs = append(allowedAPIs, NetworkingV1)
	}

	if config.IngressNetV1beta1Enabled {
		allowedAPIs = append(allowedAPIs, NetworkingV1beta1)
	}

	if config.IngressExtV1beta1Enabled {
		allowedAPIs = append(allowedAPIs, ExtensionsV1beta1)
	}

	for _, candidate := range allowedAPIs {
		if ctrlutils.CRDExists(client, candidateAPIs[candidate]) {
			return candidate, nil
		}
	}
	return OtherAPI, fmt.Errorf("no suitable Ingress API found")
}
