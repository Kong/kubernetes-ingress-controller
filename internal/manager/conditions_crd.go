package manager

import (
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
)

// CRDControllerCondition determines if a resource controller for a given CRD should be enabled.
type CRDControllerCondition struct {
	gvr           schema.GroupVersionResource
	toggleEnabled bool
	restMapper    meta.RESTMapper
	log           logr.Logger
}

func NewCRDCondition(gvr schema.GroupVersionResource, enabled bool, restMapper meta.RESTMapper) CRDControllerCondition {
	return CRDControllerCondition{
		toggleEnabled: enabled,
		gvr:           gvr,
		restMapper:    restMapper,
		log: ctrl.Log.WithName("controllers").WithName("crdCondition").
			WithValues("group", gvr.Group, "version", gvr.Version, "resource", gvr.Resource),
	}
}

func (c CRDControllerCondition) Enabled() bool {
	if c.toggleEnabled {
		if !utils.CRDExists(c.restMapper, c.gvr) {
			c.log.Info(fmt.Sprintf("disabling the '%s' controller due to missing CRD installation", c.gvr.Resource))
			return false
		}

		return true
	}

	return false
}
