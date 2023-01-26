package knative

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	knativev1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	knativeApis "knative.dev/pkg/apis"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ctrlref "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/reference"
	ctrlutils "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object/status"
)

// -----------------------------------------------------------------------------
// Knativev1alpha1 Ingress - Reconciler
// -----------------------------------------------------------------------------

// Knativev1alpha1IngressReconciler reconciles Ingress resources.
type Knativev1alpha1IngressReconciler struct {
	client.Client

	Log             logr.Logger
	Scheme          *runtime.Scheme
	DataplaneClient *dataplane.KongClient

	DataplaneAddressFinder *dataplane.AddressFinder
	StatusQueue            *status.Queue

	IngressClassName           string
	DisableIngressClassLookups bool
	CacheSyncTimeout           time.Duration

	ReferenceIndexers ctrlref.CacheIndexers
}

// SetupWithManager sets up the controller with the Manager.
func (r *Knativev1alpha1IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("Knativev1alpha1Ingress", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	// if configured, start the status updater controller
	if r.StatusQueue != nil {
		if err := c.Watch(
			&source.Channel{Source: r.StatusQueue.Subscribe(schema.GroupVersionKind{
				Group:   "networking.internal.knative.dev",
				Version: "v1alpha1",
				Kind:    "Ingress",
			})},
			&handler.EnqueueRequestForObject{},
		); err != nil {
			return err
		}
	}
	if !r.DisableIngressClassLookups {
		err = c.Watch(
			&source.Kind{Type: &netv1.IngressClass{}},
			handler.EnqueueRequestsFromMapFunc(r.listClassless),
			predicate.NewPredicateFuncs(ctrlutils.IsDefaultIngressClass),
		)
		if err != nil {
			return err
		}
	}
	preds := ctrlutils.GeneratePredicateFuncsForIngressClassFilter(r.IngressClassName)
	return c.Watch(
		&source.Kind{Type: &knativev1alpha1.Ingress{}},
		&handler.EnqueueRequestForObject{},
		preds,
	)
}

// listClassless finds and reconciles all objects without ingress class information.
func (r *Knativev1alpha1IngressReconciler) listClassless(obj client.Object) []reconcile.Request {
	resourceList := &knativev1alpha1.IngressList{}
	if err := r.Client.List(context.Background(), resourceList); err != nil {
		r.Log.Error(err, "failed to list classless ingresses")
		return nil
	}
	var recs []reconcile.Request
	for i, resource := range resourceList.Items {
		if ctrlutils.IsIngressClassEmpty(&resourceList.Items[i]) {
			recs = append(recs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: resource.Namespace,
					Name:      resource.Name,
				},
			})
		}
	}
	return recs
}

//+kubebuilder:rbac:groups=networking.internal.knative.dev,resources=ingresses,verbs=get;list;watch
//+kubebuilder:rbac:groups=networking.internal.knative.dev,resources=ingresses/status,verbs=get;update;patch

// Reconcile processes the watched objects.
func (r *Knativev1alpha1IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Knativev1alpha1Ingress", req.NamespacedName)

	// get the relevant object
	obj := new(knativev1alpha1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name

			err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, obj)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(obj)
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("resource is being deleted, its configuration will be removed", "type", "Ingress", "namespace", req.Namespace, "name", req.Name)
		err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, obj)
		if err != nil {
			return ctrl.Result{}, err
		}

		objectExistsInCache, err := r.DataplaneClient.ObjectExists(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.DataplaneClient.DeleteObject(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}

	class := new(netv1.IngressClass)
	if !r.DisableIngressClassLookups {
		if err := r.Get(ctx, types.NamespacedName{Name: r.IngressClassName}, class); err != nil {
			// we log this without taking action to support legacy configurations that only set ingressClassName or
			// used the class annotation and did not create a corresponding IngressClass. We only need this to determine
			// if the IngressClass is default or to configure default settings, and can assume no/no additional defaults
			// if none exists.
			log.V(util.DebugLevel).Info("could not retrieve IngressClass", "ingressclass", r.IngressClassName)
		}
	}
	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClass(obj, r.IngressClassName, ctrlutils.IsDefaultIngressClass(class)) {
		log.V(util.DebugLevel).Info("object missing ingress class, ensuring it's removed from configuration", "namespace", req.Namespace, "name", req.Name)
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(obj)
	}

	// update reference records for secrets referred by the ingress
	referredSecretNames := make(map[types.NamespacedName]struct{}, len(obj.Spec.TLS))
	for _, tls := range obj.Spec.TLS {
		secretNamespace := tls.SecretNamespace
		if tls.SecretNamespace != "" {
			secretNamespace = tls.SecretNamespace
		}
		nsName := types.NamespacedName{
			Namespace: secretNamespace,
			Name:      tls.SecretName,
		}
		referredSecretNames[nsName] = struct{}{}
	}
	if err := ctrlref.UpdateReferencesToSecret(ctx, r.Client, r.ReferenceIndexers, r.DataplaneClient,
		obj, referredSecretNames); err != nil {
		return ctrl.Result{}, err
	}

	// update the kong Admin API with the changes
	if err := r.DataplaneClient.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}
	// if status updates are enabled report the status for the object
	if r.DataplaneClient.AreKubernetesObjectReportsEnabled() {
		log.V(util.DebugLevel).Info("determining whether data-plane configuration has succeeded", "namespace", req.Namespace, "name", req.Name)

		if !r.DataplaneClient.KubernetesObjectIsConfigured(obj) {
			log.V(util.DebugLevel).Error(fmt.Errorf("resource not yet configured"), "resource not yet configured in the data-plane", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{Requeue: true}, nil // requeue until the object has been properly configured
		}

		log.V(util.DebugLevel).Info("determining gateway addresses for object status updates", "namespace", req.Namespace, "name", req.Name)
		addrs, err := r.DataplaneAddressFinder.GetLoadBalancerAddresses(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}

		log.V(util.DebugLevel).Info("found addresses for data-plane updating object status", "namespace", req.Namespace, "name", req.Name)
		var knativeLBIngress []knativev1alpha1.LoadBalancerIngressStatus
		for _, addr := range addrs {
			knativeIng := knativev1alpha1.LoadBalancerIngressStatus{
				IP:     addr.IP,
				Domain: addr.Hostname,
			}
			knativeLBIngress = append(knativeLBIngress, knativeIng)
		}
		ingressCondSet := knativeApis.NewLivingConditionSet()
		if obj.Status.PublicLoadBalancer == nil || len(obj.Status.PublicLoadBalancer.Ingress) != len(addrs) || !reflect.DeepEqual(obj.Status.PublicLoadBalancer.Ingress, knativeLBIngress) {
			obj.Status.MarkLoadBalancerReady(knativeLBIngress, knativeLBIngress)
			ingressCondSet.Manage(&obj.Status).MarkTrue(knativev1alpha1.IngressConditionReady)
			ingressCondSet.Manage(&obj.Status).MarkTrue(knativev1alpha1.IngressConditionNetworkConfigured)
			obj.Status.ObservedGeneration = obj.Generation
			return ctrl.Result{}, r.Status().Update(ctx, obj)
		}
		log.V(util.DebugLevel).Info("status update not needed", "namespace", req.Namespace, "name", req.Name)
	}

	return ctrl.Result{}, nil
}
