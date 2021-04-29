package ctrlutils

import (
	"context"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/kong/kubernetes-ingress-controller/pkg/parser"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/mgrutils"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KongIngressFinalizer is the finalizer used to ensure Kong configuration cleanup for deleted resources.
const KongIngressFinalizer = "configuration.konghq.com/ingress"

// UpdateKongAdmin is a helper function to take the contents of a Kong config and update the Admin API with the parsed contents.
func UpdateKongAdmin(ctx context.Context, kongCFG *sendconfig.Kong, ingressClassName string) error {
	// build the kongstate object from the Kubernetes objects in the storer
	storer := store.New(*mgrutils.CacheStores, "kong", false, false, false, logrus.StandardLogger())
	kongstate, err := parser.Build(logrus.StandardLogger(), storer)
	if err != nil {
		return err
	}

	// generate the deck configuration to be applied to the admin API
	targetConfig := deckgen.ToDeckContent(ctx, logrus.StandardLogger(), kongstate, nil, nil)

	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = sendconfig.PerformUpdate(timedCtx, logrus.StandardLogger(), kongCFG, true, false, targetConfig, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// CleanupFinalizer ensures that a deleted resource is no longer present in the object cache.
func CleanupFinalizer(ctx context.Context, c client.Client, log logr.Logger, nsn types.NamespacedName, obj client.Object) (ctrl.Result, error) {
	if HasFinalizer(obj, KongIngressFinalizer) {
		log.Info("kong ingress finalizer needs to be removed from a resource which is deleting", "ingress", obj.GetName(), "finalizer", KongIngressFinalizer)
		finalizers := []string{}
		for _, finalizer := range obj.GetFinalizers() {
			if finalizer != KongIngressFinalizer {
				finalizers = append(finalizers, finalizer)
			}
		}
		obj.SetFinalizers(finalizers)
		if err := c.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("the kong ingress finalizer was removed from an a resource which is deleting", "ingress", obj.GetName(), "finalizer", KongIngressFinalizer)
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// HasFinalizer is a helper function to check whether a client.Object
// already has a specific finalizer set.
func HasFinalizer(obj client.Object, finalizer string) bool {
	hasFinalizer := false
	for _, foundFinalizer := range obj.GetFinalizers() {
		if foundFinalizer == finalizer {
			hasFinalizer = true
		}
	}
	return hasFinalizer
}

// IsAPIAvailable is a hack to short circuit controllers for APIs which aren't available on the cluster,
// enabling us to keep separate logic and logging for some legacy API versions.
func IsAPIAvailable(mgr ctrl.Manager, obj client.Object) (bool, error) {
	if err := mgr.GetAPIReader().Get(context.Background(), client.ObjectKey{Namespace: DefaultNamespace, Name: "non-existent"}, obj); err != nil {
		if strings.Contains(err.Error(), "no matches for kind") {
			return false, nil
		}
		if errors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	}

	return true, nil
}
