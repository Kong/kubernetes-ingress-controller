package license

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	"github.com/samber/lo"
	"github.com/samber/mo"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/crds"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
)

// -----------------------------------------------------------------------------
// KongV1Alpha1 KongLicense - Reconciler
// -----------------------------------------------------------------------------

// ValidatorFunc is the function type used to validate the license string by KongV1Alpha1KongLicenseReconciler.
type ValidatorFunc func(rawLicenseString string) error

// NewKongV1Alpha1KongLicenseReconciler creates a new KongV1Alpha1KongLicenseReconciler.
// It can validate the license and set conditions accordingly when licenseValidator is provided.
// Based on whether it returns an error it sets
// `status.status.controllers[].controllerName.conditions[].type`: "LicenseValid" to "True" or "False"
// according with other fields of the condition. Field `message` is set to the returned error message.
// If not provided, the whole validation step is skipped.
func NewKongV1Alpha1KongLicenseReconciler(
	client client.Client,
	log logr.Logger,
	scheme *runtime.Scheme,
	licenseCache cache.Store,
	cacheSyncTimeout time.Duration,
	statusQueue *status.Queue,
	licenseControllerType string,
	electionID mo.Option[string],
	licenseValidator mo.Option[ValidatorFunc],
) *KongV1Alpha1KongLicenseReconciler {
	return &KongV1Alpha1KongLicenseReconciler{
		Client:                client,
		Log:                   log,
		Scheme:                scheme,
		LicenseCache:          licenseCache,
		CacheSyncTimeout:      cacheSyncTimeout,
		StatusQueue:           statusQueue,
		LicenseControllerType: licenseControllerType,
		ElectionID:            electionID,
		licenseValidator:      licenseValidator.OrEmpty(),
	}
}

// KongV1Alpha1KongLicenseReconciler reconciles KongLicense resources.
type KongV1Alpha1KongLicenseReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	LicenseCache     cache.Store
	CacheSyncTimeout time.Duration
	StatusQueue      *status.Queue
	// ControllerName is the part of field status.controllers[].controllerName: "LicenseControllerType/ControllerName".
	LicenseControllerType string
	// ElectionID is the part of field status.controllers[].controllerName: "LicenseControllerType/ElectionID".
	// This is unique identifier of the controller instance. When not specified, the field will be set to
	// status.controllers[].controllerName: "LicenseControllerType".
	ElectionID mo.Option[string]

	licenseValidator  func(rawLicenseString string) error
	chosenLicenseLock sync.RWMutex
	chosenLicense     *kongv1alpha1.KongLicense
}

const (
	// LicenseControllerTypeKIC annotates the controller type.
	LicenseControllerTypeKIC = "konghq.com/kong-ingress-controller"
)

const (
	ConditionTypeProgrammed = "Programmed"
	// ConditionReasonPickedAsLatest represents that the KongLicense being picked as the newest one.
	ConditionReasonPickedAsLatest = "PickedAsLatest"
	// ConditionReasonReplacedByNewer represents that the KongLicense is replaced by other one that is newer.
	ConditionReasonReplacedByNewer = "ReplacedByNewer"

	// ConditionTypeLicenseValid is the type of condition for the license validation.
	ConditionTypeLicenseValid = "LicenseValid"
	// ConditionReasonLicenseValid represents that the license is validated successfully.
	ConditionReasonLicenseValid = "ValidatedSuccessfully"
	// ConditionReasonLicenseInvalid represents that the provided license is invalid.
	ConditionReasonLicenseInvalid = "InvalidLicenseProvided"

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
	blder := ctrl.NewControllerManagedBy(mgr).
		Named("KongV1Alpha1KongLicense").
		WithOptions(controller.Options{
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
		})

	// if configured, start the status updater controller
	if r.StatusQueue != nil {
		blder.WatchesRawSource(
			source.Channel(r.StatusQueue.Subscribe(schema.GroupVersionKind{
				Group:   "configuration.konghq.com",
				Version: "v1alpha1",
				Kind:    "KongLicense",
			}),
				&handler.EnqueueRequestForObject{},
			),
		)
	}
	return blder.Watches(&kongv1alpha1.KongLicense{},
		&handler.EnqueueRequestForObject{},
		builder.WithPredicates(predicate.NewPredicateFuncs(isKongLicenseEnabled)),
	).
		Complete(r)
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
				log.V(logging.DebugLevel).Info("KongLicense deleted in cluster, delete it in cache")
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
	log.V(logging.DebugLevel).Info("Reconciling resource", "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(logging.DebugLevel).Info("Resource is being deleted, its configuration will be removed", "type", "KongLicense")

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
		log.V(logging.DebugLevel).Info("Picked KongLicense being reconciled", "name", obj.Name)
		err := r.ensureControllerStatusConditions(ctx, obj, metav1.ConditionTrue, ConditionReasonPickedAsLatest, "")
		if err != nil {
			return ctrl.Result{}, err
		}

		oldChosenLicense := r.getChosenLicense()
		if oldChosenLicense != nil && oldChosenLicense.Name != chosenLicense.Name {
			r.Log.V(logging.DebugLevel).Info("Originally picked KongLicense replaced", "name", oldChosenLicense.Name)
			err := r.ensureControllerStatusConditions(ctx, oldChosenLicense, metav1.ConditionFalse, ConditionReasonReplacedByNewer, "Replaced by newer created KongLicense")
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		r.setChosenLicense(chosenLicense)
	}

	return ctrl.Result{}, nil
}

// License is a wrapper for kong.License that include the information about its validity.
// Field IsValid is set to true if the license is valid, otherwise it's set to false.
// If it's not set at all, it means that the license validation is not configured,
// thus the information about the license validity is not available.
type License struct {
	License kong.License
	IsValid mo.Option[bool]
}

// GetValidatedLicense is the interface to get the license in Kong configuration format and information
// about its validity when if validator has been configured.
func (r *KongV1Alpha1KongLicenseReconciler) GetValidatedLicense() mo.Option[License] {
	chosenLicense := r.getChosenLicense()
	if chosenLicense == nil {
		r.Log.V(logging.DebugLevel).Info("No KongLicense available")
		return mo.None[License]()
	}
	r.Log.V(logging.DebugLevel).Info("Get license from KongLicense resource", "name", chosenLicense.Name)
	isValid := mo.None[bool]()
	if r.licenseValidator != nil {
		isValid = mo.Some(r.licenseValidator(chosenLicense.RawLicenseString) == nil)
	}
	return mo.Some(License{
		License: kong.License{
			ID:      kong.String(uuid.NewSHA1(uuid.Nil, []byte("KongLicense:"+chosenLicense.Name)).String()),
			Payload: kong.String(chosenLicense.RawLicenseString),
		},
		IsValid: isValid,
	})
}

// GetLicense returns the license in Kong configuration format. It does not check the license validity.
// This method is provided to implement the translator.LicenseGetter interface (use case where information
// about validity is not required). When the info about validity is required, use method GetValidatedLicense.
func (r *KongV1Alpha1KongLicenseReconciler) GetLicense() mo.Option[kong.License] {
	unpackedL, ok := r.GetValidatedLicense().Get()
	if !ok {
		return mo.None[kong.License]()
	}
	return mo.Some(unpackedL.License)
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
	if cn, ok := r.ElectionID.Get(); ok {
		return r.LicenseControllerType + "/" + cn
	}
	return r.LicenseControllerType
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
			r.Log.V(logging.DebugLevel).Info("Picked KongLicense remaining in cache", "name", chosenLicense.Name)
			return r.ensureControllerStatusConditions(ctx, chosenLicense, metav1.ConditionTrue, ConditionReasonPickedAsLatest, "")
		}
	}
	return nil
}

func setLicenseValidityCondition(ks *[]metav1.Condition, licenseValidationError error) {
	_, index, found := lo.FindIndexOf(*ks, func(condition metav1.Condition) bool {
		return condition.Type == ConditionTypeLicenseValid
	})
	if !found {
		*ks = append(*ks, metav1.Condition{
			Type: ConditionTypeLicenseValid,
		})
		index = len(*ks) - 1
	}
	licenseCondition := &(*ks)[index]

	licenseCondition.LastTransitionTime = metav1.Now()
	setCondition := func(status metav1.ConditionStatus, reason string, message string) {
		licenseCondition.Status = status
		licenseCondition.Reason = reason
		licenseCondition.Message = message
	}
	if licenseValidationError != nil {
		setCondition(metav1.ConditionFalse, ConditionReasonLicenseInvalid, licenseValidationError.Error())
		return
	}
	setCondition(metav1.ConditionTrue, ConditionReasonLicenseValid, "")
}

// ensureControllerStatusConditions updates the "programmed" condition
// in the controller status item managed by the reconciler if required.
// If licenseValidator is provided, it also updates the "LicenseValid" condition.
func (r *KongV1Alpha1KongLicenseReconciler) ensureControllerStatusConditions(
	ctx context.Context, license *kongv1alpha1.KongLicense,
	programmedStatus metav1.ConditionStatus,
	reason string, message string,
) error {
	// Get the latest status of target KongLicense.
	if err := r.Client.Get(ctx, k8stypes.NamespacedName{Name: license.Name}, license); err != nil {
		return fmt.Errorf("failed to get latest version of KongLicense %s: %w", license.Name, err)
	}

	fullControllerName := r.fullControllerName()
	// Find the managed controller status item and append new item when absent.
	_, controllerIndex, found := lo.FindIndexOf(license.Status.KongLicenseControllerStatuses, func(controllerStatus kongv1alpha1.KongLicenseControllerStatus) bool {
		return controllerStatus.ControllerName == fullControllerName
	})
	if !found {
		wantedControllerStatus := kongv1alpha1.KongLicenseControllerStatus{
			ControllerName: fullControllerName,
		}
		license.Status.KongLicenseControllerStatuses = append(license.Status.KongLicenseControllerStatuses, wantedControllerStatus)
		controllerIndex = 0
	}
	ctrlManagedConditions := &license.Status.KongLicenseControllerStatuses[controllerIndex].Conditions
	if r.licenseValidator != nil {
		setLicenseValidityCondition(ctrlManagedConditions, r.licenseValidator(license.RawLicenseString))
	}

	wantedCondition := metav1.Condition{
		Type:               ConditionTypeProgrammed,
		LastTransitionTime: metav1.Now(),
		Status:             programmedStatus,
		Reason:             reason,
		Message:            message,
	}
	// Find the "programmed" condition in the condition list.
	programmedCondition, conditionIndex, found := lo.FindIndexOf(*ctrlManagedConditions, func(condition metav1.Condition) bool {
		return condition.Type == ConditionTypeProgrammed
	})
	if !found {
		// Return error if the number of conditions already reached the limit.
		if len(license.Status.KongLicenseControllerStatuses) >= maxConditionNum {
			return fmt.Errorf("already %d condition items in controller status exceeding the limit %d, cannot add new item",
				len(license.Status.KongLicenseControllerStatuses), maxConditionNum)
		}
		*ctrlManagedConditions = append(
			*ctrlManagedConditions, wantedCondition,
		)
		return r.Client.Status().Update(ctx, license)
	}
	if programmedCondition.Status != programmedStatus {
		(*ctrlManagedConditions)[conditionIndex] = wantedCondition
	}

	return r.Client.Status().Update(ctx, license)
}

// WrapKongLicenseReconcilerToDynamicCRDController wraps KongLicenseReconciler to DynamicCRDController
// to watch presence of KongLicense CRD to avoid aborts if KongLicense is not installed when controller initialized.
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
