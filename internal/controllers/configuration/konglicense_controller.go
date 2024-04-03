package configuration

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/samber/mo"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/crds"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
)

// -----------------------------------------------------------------------------
// KongV1Alpha1 KongLicense - Reconciler
// -----------------------------------------------------------------------------

// KongV1Alpha1KongLicenseReconciler reconciles KongLicense resources.
type KongV1Alpha1KongLicenseReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	LicenseCache     cache.Store
	CacheSyncTimeout time.Duration
	StatusQueue      *status.Queue
	// ControllerName is the name of the controller to use in `controllerName` field in controller status item.
	ControllerName string

	chosenLicenseLock sync.RWMutex
	chosenLicense     *kongv1alpha1.KongLicense
}

const (
	// LicenseControllerType annotates the controller type.
	LicenseControllerType = "konghq.com/kong-ingress-controller"
)

const (
	ConditionTypeProgrammed = "Programmed"
	// ConditionReasonPickedAsLatest represents that the KongLicense being picked as the newest one.
	ConditionReasonPickedAsLatest = "PickedAsLatest"
	// ConditionReasonReplacedByNewer represents that the KongLicense is replaced by other one that is newer.
	ConditionReasonReplacedByNewer = "ReplacedByNewer"
	// maxConditionNum is the maximum number of condition items in the controller status.
	maxConditionNum = 8
)

var _ controllers.Reconciler = &KongV1Alpha1KongLicenseReconciler{}

func NewLicenseCache() cache.Store {
	return cache.NewStore(kongLicenseKeyFunc)
}

func kongLicenseKeyFunc(obj interface{}) (string, error) {
	l, ok := obj.(*kongv1alpha1.KongLicense)
	if !ok {
		return "", fmt.Errorf("object is type %T, not KongLicense", obj)
	}
	return l.Name, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongV1Alpha1KongLicenseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("KongV1Alpha1KongLicense", mgr, controller.Options{
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
				Group:   "configuration.konghq.com",
				Version: "v1alpha1",
				Kind:    "KongLicense",
			})},
			&handler.EnqueueRequestForObject{},
		); err != nil {
			return err
		}
	}
	return c.Watch(
		source.Kind(mgr.GetCache(), &kongv1alpha1.KongLicense{}),
		&handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(isKongLicenseEnabled),
	)
}

// SetLogger sets the logger.
func (r *KongV1Alpha1KongLicenseReconciler) SetLogger(l logr.Logger) {
	r.Log = l
}

//+kubebuilder:rbac:groups=configuration.konghq.com,resources=konglicenses,verbs=get;list;watch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=konglicenses/status,verbs=get;update;patch

// Reconcile processes the watched objects.
func (r *KongV1Alpha1KongLicenseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("KongV1Alpha1KongLicense", req.NamespacedName)

	// get the relevant object
	obj := new(kongv1alpha1.KongLicense)

	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name
			_, objectExistsInCache, err := r.LicenseCache.Get(obj)
			if err != nil {
				return ctrl.Result{}, err
			}
			if objectExistsInCache {
				// Delete the object in the cache first.
				log.V(util.DebugLevel).Info("KongLicense deleted in cluster, delete it in cache")
				if err := r.LicenseCache.Delete(obj); err != nil {
					return ctrl.Result{}, err
				}

				// Then pick the effective license in KongLicenses remaining in cache.
				if err := r.repickLicenseOnDelete(ctx, obj); err != nil {
					return ctrl.Result{}, err
				}
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("Reconciling resource", "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("Resource is being deleted, its configuration will be removed", "type", "KongLicense")

		_, objectExistsInCache, err := r.LicenseCache.Get(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			// Delete the object in the cache first.
			if err := r.LicenseCache.Delete(obj); err != nil {
				return ctrl.Result{}, err
			}

			// Then pick the effective license in KongLicenses remaining in cache.
			if err := r.repickLicenseOnDelete(ctx, obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}

	// Add the new/updated KongLicense to cache.
	if err := r.LicenseCache.Add(obj); err != nil {
		return ctrl.Result{}, err
	}

	// Trigger a compare on stored KongLicenses in cache and pick the newest.
	chosenLicense := r.pickLicenseInCache()
	if chosenLicense.Name == obj.Name {
		log.V(util.DebugLevel).Info("Picked KongLicense being reconciled", "name", obj.Name)
		err := r.ensureControllerStatusProgrammedCondition(ctx, obj, metav1.ConditionTrue, ConditionReasonPickedAsLatest, "")
		if err != nil {
			return ctrl.Result{}, err
		}

		oldChosenLicense := r.getChosenLicense()
		if oldChosenLicense != nil && oldChosenLicense.Name != chosenLicense.Name {
			r.Log.V(util.DebugLevel).Info("Originally picked KongLicense replaced", "name", oldChosenLicense.Name)
			err := r.ensureControllerStatusProgrammedCondition(ctx, oldChosenLicense, metav1.ConditionFalse, ConditionReasonReplacedByNewer, "Replaced by newer created KongLicense")
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		r.setChosenLicense(chosenLicense)
	}

	return ctrl.Result{}, nil
}

// GetLicense is the interface to get the license in Kong configuration format to use in translator.
func (r *KongV1Alpha1KongLicenseReconciler) GetLicense() mo.Option[kong.License] {
	chosenLicense := r.getChosenLicense()
	if chosenLicense == nil {
		r.Log.V(util.DebugLevel).Info("No KongLicense available")
		return mo.None[kong.License]()
	}
	r.Log.V(util.DebugLevel).Info("Get license from KongLicense resource", "name", chosenLicense.Name)
	// TODO: Validate KongLicense on Kong gateway:
	// https://github.com/Kong/kubernetes-ingress-controller/issues/5566
	return mo.Some(kong.License{
		ID:      kong.String(uuid.NewSHA1(uuid.Nil, []byte("KongLicense:"+chosenLicense.Name)).String()),
		Payload: kong.String(chosenLicense.RawLicenseString),
	})
}

// -----------------------------------------------------------------------------
// KongV1Alpha1 KongLicense - Private Methods of Reconciler
// -----------------------------------------------------------------------------

func isKongLicenseEnabled(obj client.Object) bool {
	kongLicense, ok := obj.(*kongv1alpha1.KongLicense)
	if !ok {
		return false
	}
	return kongLicense.Enabled
}

// compareKongLicense returns true if license1 is newer than license2 (compared by metadata.creationTimestamp).
// If the creationTimestamp equals or not comparable, returns the one with lexical smaller name.
func compareKongLicense(license1, license2 *kongv1alpha1.KongLicense) bool {
	if license1.CreationTimestamp.After(license2.CreationTimestamp.Time) {
		return true
	}
	if license2.CreationTimestamp.After(license1.CreationTimestamp.Time) {
		return false
	}
	return license1.Name < license2.Name
}

// pickLicenseInCache picks the newest license in the cache.
func (r *KongV1Alpha1KongLicenseReconciler) pickLicenseInCache() *kongv1alpha1.KongLicense {
	licenseList := r.LicenseCache.List()
	var chosenLicense *kongv1alpha1.KongLicense
	for _, obj := range licenseList {
		license, ok := obj.(*kongv1alpha1.KongLicense)
		if !ok {
			continue
		}
		if chosenLicense == nil || compareKongLicense(license, chosenLicense) {
			chosenLicense = license
		}
	}
	return chosenLicense
}

// fullControllerName returns the full controllerName used in the controller status item
// combined with constant type and reconciler's own controller name.
func (r *KongV1Alpha1KongLicenseReconciler) fullControllerName() string {
	return LicenseControllerType + "/" + r.ControllerName
}

// setChosenLicense sets the chosen effective KongLicense copy in the cache.
func (r *KongV1Alpha1KongLicenseReconciler) setChosenLicense(l *kongv1alpha1.KongLicense) {
	r.chosenLicenseLock.Lock()
	defer r.chosenLicenseLock.Unlock()
	r.chosenLicense = l.DeepCopy()
}

// getChosenLicense fetches the the chosen effective KongLicense in the cache.
func (r *KongV1Alpha1KongLicenseReconciler) getChosenLicense() *kongv1alpha1.KongLicense {
	r.chosenLicenseLock.RLock()
	defer r.chosenLicenseLock.RUnlock()
	return r.chosenLicense
}

func (r *KongV1Alpha1KongLicenseReconciler) repickLicenseOnDelete(ctx context.Context, deletedLicense *kongv1alpha1.KongLicense) error {
	oldChosenLicense := r.getChosenLicense()
	// Trigger a repick of license if the originally chosen license is deleted.
	if oldChosenLicense.Name == deletedLicense.Name {
		chosenLicense := r.pickLicenseInCache()
		r.setChosenLicense(chosenLicense)
		if chosenLicense != nil {
			r.Log.V(util.DebugLevel).Info("Picked KongLicense remaining in cache", "name", chosenLicense.Name)
			return r.ensureControllerStatusProgrammedCondition(ctx, chosenLicense, metav1.ConditionTrue, ConditionReasonPickedAsLatest, "")
		}
	}
	return nil
}

// ensureControllerStatusProgrammedCondition updates the "programmed" condition
// in the controller status item managed by the reconciler if required.
func (r *KongV1Alpha1KongLicenseReconciler) ensureControllerStatusProgrammedCondition(
	ctx context.Context, license *kongv1alpha1.KongLicense,
	programmedStatus metav1.ConditionStatus,
	reason string, message string,
) error {
	// Get the latest status of target KongLicense.
	err := r.Client.Get(ctx, k8stypes.NamespacedName{Name: license.Name}, license)
	if err != nil {
		return fmt.Errorf("failed to get latest version of KongLicense %s: %w", license.Name, err)
	}

	fullControllerName := r.fullControllerName()
	// Find the managed controller status item and append new item when absent.
	controllerStatus, controllerIndex, found := lo.FindIndexOf(license.Status.KongLicenseControllerStatuses, func(controllerStatus kongv1alpha1.KongLicenseControllerStatus) bool {
		return controllerStatus.ControllerName == fullControllerName
	})
	if !found {
		wantedControllerStatus := kongv1alpha1.KongLicenseControllerStatus{
			ControllerName: fullControllerName,
			Conditions: []metav1.Condition{
				{
					Type:               ConditionTypeProgrammed,
					LastTransitionTime: metav1.Now(),
					Status:             programmedStatus,
					Reason:             reason,
				},
			},
		}
		license.Status.KongLicenseControllerStatuses = append(license.Status.KongLicenseControllerStatuses, wantedControllerStatus)
		return r.Client.Status().Update(ctx, license)
	}

	wantedCondition := metav1.Condition{
		Type:               ConditionTypeProgrammed,
		LastTransitionTime: metav1.Now(),
		Status:             programmedStatus,
		Reason:             reason,
		Message:            message,
	}
	// Find the "programmed" condition in the condition list.
	programmedCondition, conditionIndex, found := lo.FindIndexOf(controllerStatus.Conditions, func(condition metav1.Condition) bool {
		return condition.Type == ConditionTypeProgrammed
	})
	if !found {
		// Return error if the number of conditions already reached the limit.
		if len(license.Status.KongLicenseControllerStatuses) >= maxConditionNum {
			return fmt.Errorf("already %d condition items in controller status exceeding the limit %d, cannot add new item",
				len(license.Status.KongLicenseControllerStatuses), maxConditionNum)
		}
		license.Status.KongLicenseControllerStatuses[controllerIndex].Conditions = append(
			license.Status.KongLicenseControllerStatuses[controllerIndex].Conditions, wantedCondition,
		)
		return r.Client.Status().Update(ctx, license)
	}
	if programmedCondition.Status != programmedStatus {
		license.Status.KongLicenseControllerStatuses[controllerIndex].Conditions[conditionIndex] = wantedCondition
		return r.Client.Status().Update(ctx, license)
	}

	return nil
}

// WrapKongLicenseReconcilerToDynamicCRDController wraps KongLicenseReconciler to DynamicCRDController
// to watch precense of KongLicense CRD to avoid aborts if KongLicense is not installed when controller initialized.
func WrapKongLicenseReconcilerToDynamicCRDController(
	ctx context.Context, mgr ctrl.Manager, r *KongV1Alpha1KongLicenseReconciler,
) *crds.DynamicCRDController {
	return &crds.DynamicCRDController{
		Manager:          mgr,
		Log:              ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Dynamic/KongLicense"),
		CacheSyncTimeout: r.CacheSyncTimeout,
		RequiredCRDs: []schema.GroupVersionResource{
			{
				Group:    kongv1alpha1.GroupVersion.Group,
				Version:  kongv1alpha1.GroupVersion.Version,
				Resource: "konglicenses",
			},
		},
		Controller: r,
	}
}
