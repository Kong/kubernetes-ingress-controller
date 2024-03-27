package gateway

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// GatewayClass Controller - Vars & Consts
// -----------------------------------------------------------------------------

var (
	// _controllerName is the unique identifier for this controller and is used
	// within GatewayClass resources to indicate that this controller should
	// support connected Gateway resources.
	_controllerName gatewayv1.GatewayController = "konghq.com/kic-gateway-controller"

	// _controllerNameLock guards access to _controllerName.
	_controllerNameLock sync.RWMutex
)

func SetControllerName(name gatewayv1.GatewayController) {
	_controllerNameLock.Lock()
	defer _controllerNameLock.Unlock()
	_controllerName = name
}

func GetControllerName() gatewayv1.GatewayController {
	_controllerNameLock.RLock()
	defer _controllerNameLock.RUnlock()
	return _controllerName
}

// -----------------------------------------------------------------------------
// GatewayClass Controller - Reconciler
// -----------------------------------------------------------------------------

// GatewayClassReconciler reconciles a GatewayClass object.
type GatewayClassReconciler struct { //nolint:revive
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	CacheSyncTimeout time.Duration
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("V1Beta1GatewayClass", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	return c.Watch(
		source.Kind(mgr.GetCache(), &gatewayv1.GatewayClass{}),
		&handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(r.GatewayClassIsUnmanaged),
	)
}

// -----------------------------------------------------------------------------
// GatewayClass Controller - Watch Predicates
// -----------------------------------------------------------------------------

// GatewayClassIsUnmanaged is a watch predicate which filters out reconciliation events for
// gateway objects which aren't annotated as unmanaged.
func (r *GatewayClassReconciler) GatewayClassIsUnmanaged(obj client.Object) bool {
	gatewayClass, ok := obj.(*gatewayv1.GatewayClass)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"gatewayclass watch predicate received unexpected object type",
			"expected", "*gatewayv1.GatewayClass", "found", reflect.TypeOf(obj),
		)
		return false
	}

	return isGatewayClassControlledAndUnmanaged(gatewayClass)
}

// -----------------------------------------------------------------------------
// GatewayClass Controller - Reconciliation
// -----------------------------------------------------------------------------

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses/status,verbs=get;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *GatewayClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("GatewayV1Beta1GatewayClass", req.NamespacedName)

	gwc := new(gatewayv1.GatewayClass)
	if err := r.Client.Get(ctx, req.NamespacedName, gwc); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(util.DebugLevel).Info("object enqueued no longer exists, skipping", "name", req.Name)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("processing gatewayclass", "name", req.Name)

	if isGatewayClassControlledAndUnmanaged(gwc) {
		alreadyAccepted := util.CheckCondition(
			gwc.Status.Conditions,
			util.ConditionType(gatewayv1.GatewayClassConditionStatusAccepted),
			util.ConditionReason(gatewayv1.GatewayClassReasonAccepted),
			metav1.ConditionTrue,
			gwc.Generation,
		)

		if !alreadyAccepted {
			acceptedCondtion := metav1.Condition{
				Type:               string(gatewayv1.GatewayClassConditionStatusAccepted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: gwc.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1.GatewayClassReasonAccepted),
				Message:            "the gatewayclass has been accepted by the controller",
			}
			setGatewayClassCondition(gwc, acceptedCondtion)
			return ctrl.Result{}, r.Status().Update(ctx, pruneGatewayClassStatusConds(gwc))
		}
	}

	return ctrl.Result{}, nil
}

// SetLogger sets the logger.
func (r *GatewayClassReconciler) SetLogger(l logr.Logger) {
	r.Log = l
}

// -----------------------------------------------------------------------------
// GatewayClass Controller - Private
// -----------------------------------------------------------------------------

// pruneGatewayClassStatusConds cleans out old status conditions if the
// Gatewayclass currently has more status conditions set than the 8 maximum
// allowed by the Kubernetes API.
func pruneGatewayClassStatusConds(gwc *gatewayv1.GatewayClass) *gatewayv1.GatewayClass {
	if len(gwc.Status.Conditions) > maxConds {
		gwc.Status.Conditions = gwc.Status.Conditions[len(gwc.Status.Conditions)-maxConds:]
	}
	return gwc
}

// setGatewayClassCondition sets the condition with specified type in gatewayclass status
// to expected condition in newCondition.
// if the gatewayclass status does not contain a condition with that type, add one more condition.
// if the gatewayclass status contains condition(s) with the type, then replace with the new condition.
func setGatewayClassCondition(gwc *gatewayv1.GatewayClass, newCondition metav1.Condition) {
	newConditions := []metav1.Condition{}
	for _, condition := range gwc.Status.Conditions {
		if condition.Type != newCondition.Type {
			newConditions = append(newConditions, condition)
		}
	}
	newConditions = append(newConditions, newCondition)
	gwc.Status.Conditions = newConditions
}
