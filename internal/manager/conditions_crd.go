package manager

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
)

// CRDControllerCondition determines if a resource controller for a given CRD should be enabled.
type CRDControllerCondition struct {
	gvr        schema.GroupVersionResource
	enabled    bool
	restMapper meta.RESTMapper
	log        logr.Logger
}

func NewCRDCondition(gvr schema.GroupVersionResource, enabled bool, restMapper meta.RESTMapper) CRDControllerCondition {
	return CRDControllerCondition{
		enabled:    enabled,
		gvr:        gvr,
		restMapper: restMapper,
		log: ctrl.Log.WithName("crd_controller_condition").
			WithValues("group", gvr.Group, "version", gvr.Version, "resource", gvr.Resource),
	}
}

func (c CRDControllerCondition) Enabled() bool {
	if c.enabled {
		if !utils.CRDExists(c.restMapper, c.gvr) {
			c.log.Info("Disabling the resource controller due to missing CRD installation")
			return false
		}

		return true
	}

	return false
}
