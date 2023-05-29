package manager

import (
	"fmt"

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"

	ctrlutils "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
)

type IngressAPI int

const (
	NoIngressAPI IngressAPI = iota
	NetworkingV1
)

// IngressControllerConditions negotiates the best Ingress API version supported by both KIC and the k8s apiserver and
// provides functions to determine if particular controllers should be enabled.
// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/1666
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

// IngressNetV1Enabled returns true if the chosen ingress API version is networking.k8s.io/v1 and it's enabled.
func (s *IngressControllerConditions) IngressNetV1Enabled() bool {
	return s.chosenVersion == NetworkingV1 && s.cfg.IngressNetV1Enabled
}

// IngressClassNetV1Enabled returns true if the chosen ingress class API version is networking.k8s.io/v1 and it's enabled.
func (s *IngressControllerConditions) IngressClassNetV1Enabled() bool {
	return s.chosenVersion == NetworkingV1 && s.cfg.IngressClassNetV1Enabled
}

func negotiateIngressAPI(config *Config, mapper meta.RESTMapper) (IngressAPI, error) { //nolint:unparam
	var allowedAPIs []IngressAPI
	candidateAPIs := map[IngressAPI]schema.GroupVersionResource{
		NetworkingV1: {
			Group:    netv1.SchemeGroupVersion.Group,
			Version:  netv1.SchemeGroupVersion.Version,
			Resource: "ingresses",
		},
	}

	// Please note the order is not arbitrary - the most mature APIs will get picked first.
	if config.IngressNetV1Enabled {
		allowedAPIs = append(allowedAPIs, NetworkingV1)
	}

	for _, candidate := range allowedAPIs {
		if ctrlutils.CRDExists(mapper, candidateAPIs[candidate]) {
			return candidate, nil
		}
	}
	return NoIngressAPI, nil
}

func ShouldEnableCRDController(gvr schema.GroupVersionResource, restMapper meta.RESTMapper) bool {
	if !ctrlutils.CRDExists(restMapper, gvr) {
		ctrl.Log.WithName("controllers").WithName("crdCondition").
			Info(fmt.Sprintf("Disabling controller for Group=%s, Resource=%s due to missing CRD", gvr.GroupVersion(), gvr.Resource))
		return false
	}
	return true
}
