package gateway

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// GatewayClass Controller - Vars & Consts
// -----------------------------------------------------------------------------

const (
	// ControllerName is the unique identifier for this controller and is used
	// within GatewayClass resources to indicate that this controller should
	// support connected Gateway resources.
	ControllerName gatewayv1alpha2.GatewayController = "konghq.com/kic-gateway-controller"
)

// -----------------------------------------------------------------------------
// GatewayClass Controller - Reconciler
// -----------------------------------------------------------------------------

// GatewayClassReconciler reconciles a GatewayClass object
type GatewayClassReconciler struct { //nolint:revive
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("V1Alpha2GatewayClass", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
	})
	if err != nil {
		return err
	}
	return c.Watch(
		&source.Kind{Type: &gatewayv1alpha2.GatewayClass{}},
		&handler.EnqueueRequestForObject{},
	)
}

// -----------------------------------------------------------------------------
// Gateway Controller - Reconciliation
// -----------------------------------------------------------------------------

//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses/status,verbs=get;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *GatewayClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("V1Alpha2GatewayClass", req.NamespacedName)

	gwc := new(gatewayv1alpha2.GatewayClass)
	if err := r.Client.Get(ctx, req.NamespacedName, gwc); err != nil {
		if errors.IsNotFound(err) {
			log.V(util.DebugLevel).Info("object enqueued no longer exists, skipping", "name", req.Name)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("processing gatewayclass", "name", req.Name)

	if gwc.Spec.ControllerName == ControllerName {
		alreadyAccepted := false
		for _, cond := range gwc.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayClassConditionStatusAccepted) {
				if cond.ObservedGeneration == gwc.Generation {
					alreadyAccepted = true
				}
			}
		}

		if !alreadyAccepted {
			gwc.Status.Conditions = append(gwc.Status.Conditions, metav1.Condition{
				Type:               string(gatewayv1alpha2.GatewayClassConditionStatusAccepted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: gwc.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1alpha2.GatewayClassReasonAccepted),
				Message:            "the gatewayclass has been accepted by the controller",
			})
			return ctrl.Result{}, r.Status().Update(ctx, pruneGatewayClassStatusConds(gwc))
		}
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// GatewayClass Controller - Private
// -----------------------------------------------------------------------------

// pruneGatewayClassStatusConds cleans out old status conditions if the
// Gatewayclass currently has more status conditions set than the 8 maximum
// allowed by the Kubernetes API.
func pruneGatewayClassStatusConds(gwc *gatewayv1alpha2.GatewayClass) *gatewayv1alpha2.GatewayClass {
	if len(gwc.Status.Conditions) > maxConds {
		gwc.Status.Conditions = gwc.Status.Conditions[len(gwc.Status.Conditions)-maxConds:]
	}
	return gwc
}
