package manager

import (
	"fmt"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"

	ctrlutils "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
)

type IngressAPI int

const (
	OtherAPI          IngressAPI = iota
	NetworkingV1      IngressAPI = iota
	NetworkingV1beta1 IngressAPI = iota
	ExtensionsV1beta1 IngressAPI = iota
)

// IngressControllerConditions negotiate the best Ingress API version supported by both KIC and the k8s apiserver and
// provide functions to determine if particular controllers should be enabled.
type IngressControllerConditions struct {
	chosenVersion IngressAPI
	cfg           *Config
}

func NewIngressControllersConditions(cfg *Config, mapper meta.RESTMapper) (*IngressControllerConditions, error) {
	chosenVersion, err := negotiateIngressAPI(cfg, mapper)
	if err != nil {
		return nil, err
	}

	return &IngressControllerConditions{chosenVersion: chosenVersion, cfg: cfg}, nil
}

// IngressExtV1beta1Enabled returns true if the chosen ingress API version is extensions/v1beta1 and it's enabled.
func (s *IngressControllerConditions) IngressExtV1beta1Enabled() bool {
	return s.chosenVersion == ExtensionsV1beta1 && s.cfg.IngressExtV1beta1Enabled
}

// IngressNetV1beta1Enabled returns true if the chosen ingress API version is networking.k8s.io/v1beta1 and it's enabled.
func (s *IngressControllerConditions) IngressNetV1beta1Enabled() bool {
	return s.chosenVersion == NetworkingV1beta1 && s.cfg.IngressNetV1beta1Enabled
}

// IngressNetV1Enabled returns true if the chosen ingress API version is networking.k8s.io/v1 and it's enabled.
func (s *IngressControllerConditions) IngressNetV1Enabled() bool {
	return s.chosenVersion == NetworkingV1 && s.cfg.IngressNetV1Enabled
}

// IngressClassNetV1Enabled returns true if the chosen ingress class API version is networking.k8s.io/v1 and it's enabled.
func (s *IngressControllerConditions) IngressClassNetV1Enabled() bool {
	return s.chosenVersion == NetworkingV1 && s.cfg.IngressClassNetV1Enabled
}

func negotiateIngressAPI(config *Config, mapper meta.RESTMapper) (IngressAPI, error) {
	var allowedAPIs []IngressAPI
	candidateAPIs := map[IngressAPI]schema.GroupVersionResource{
		NetworkingV1: {
			Group:    netv1.SchemeGroupVersion.Group,
			Version:  netv1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		},
		NetworkingV1beta1: {
			Group:    netv1beta1.SchemeGroupVersion.Group,
			Version:  netv1beta1.SchemeGroupVersion.Version,
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
		if ctrlutils.CRDExists(mapper, candidateAPIs[candidate]) {
			return candidate, nil
		}
	}
	return OtherAPI, fmt.Errorf("no suitable Ingress API found")
}
