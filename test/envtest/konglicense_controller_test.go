package envtest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"

	ctrllicense "github.com/kong/kubernetes-ingress-controller/v3/controllers/license"
)

func TestKongLicenseController(t *testing.T) {
	scheme := Scheme(t, WithKong)
	cfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reconciler := ctrllicense.NewKongV1Alpha1KongLicenseReconciler(
		ctrlClient,
		logr.Discard(),
		scheme,
		ctrllicense.NewLicenseCache(),
		time.Second,
		nil,
		ctrllicense.LicenseControllerTypeKIC,
		mo.Some("test"),
		mo.None[ctrllicense.ValidatorFunc](),
	)

	StartReconcilers(ctx, t, ctrlClient.Scheme(), cfg, reconciler)

	const (
		fullControllerName  = ctrllicense.LicenseControllerTypeKIC + "/test"
		conditionProgrammed = ctrllicense.ConditionTypeProgrammed
	)

	t.Log("Create a KongLicense and verify that it is reconciled")
	kongLicense1 := &kongv1alpha1.KongLicense{
		ObjectMeta: metav1.ObjectMeta{
			Name: "license-1",
		},
		RawLicenseString: "test-license-1",
		Enabled:          true,
	}
	err := ctrlClient.Create(ctx, kongLicense1)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		l, ok := reconciler.GetLicense().Get()
		if !ok {
			return false
		}
		return l.Payload != nil && *l.Payload == "test-license-1"
	}, waitTime, tickTime, "Should return the license in GetLicense method")

	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{Name: kongLicense1.Name}, kongLicense1)
		require.NoError(t, err, "Should not have error in getting the latest status")
		return findConditionInControllerStatus(
			kongLicense1, fullControllerName, conditionProgrammed, metav1.ConditionTrue)
	}, waitTime, tickTime, "The Programmed condition in controller status should be true")

	t.Log("Wait for a second and create a new KongLicense, then verify it replaces the old one")
	// We're waiting specifically for 1 second because Kubernetes object's CreationTimestamp is stored with a seconds-level
	// precision. If we create the new KongLicense immediately after the old one, the CreationTimestamps will be the same
	// and the controller will not pick up the new one as the latest.
	// The controller should always pick up only a single KongLicense with the latest CreationTimestamp to be used for
	// configuring Gateways. That one KongLicense should have Programmed condition set to true, while all others should
	// have it set to false.
	// CreationTimestamp precision upstream issue: https://github.com/kubernetes/kubernetes/issues/81026
	time.Sleep(time.Second)
	kongLicense2 := &kongv1alpha1.KongLicense{
		ObjectMeta: metav1.ObjectMeta{
			Name: "license-2",
		},
		RawLicenseString: "test-license-2",
		Enabled:          true,
	}
	err = ctrlClient.Create(ctx, kongLicense2)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		l, ok := reconciler.GetLicense().Get()
		if !ok {
			return false
		}
		return l.Payload != nil && *l.Payload == "test-license-2"
	}, waitTime, tickTime, "Should return *NEW* license in GetLicense method")

	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{Name: kongLicense1.Name}, kongLicense1)
		require.NoError(t, err, "Should not have error in getting the latest status")
		return findConditionInControllerStatus(
			kongLicense1, fullControllerName, conditionProgrammed, metav1.ConditionFalse)
	}, waitTime, tickTime, "The Programmed condition in controller status of the *OLD* KongLicense should be false")

	t.Log("Delete the new KongLicense and verify that the old one get chosen")
	err = ctrlClient.Delete(ctx, kongLicense2)
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		l, ok := reconciler.GetLicense().Get()
		if !ok {
			return false
		}
		return l.Payload != nil && *l.Payload == "test-license-1"
	}, waitTime, tickTime, "Should return old license in GetLicense method after the new one deleted")
	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{Name: kongLicense1.Name}, kongLicense1)
		require.NoError(t, err, "Should not have error in getting the latest status")
		return findConditionInControllerStatus(
			kongLicense1, fullControllerName, conditionProgrammed, metav1.ConditionTrue)
	}, waitTime, tickTime, "The Programmed condition in controller status should be back to true after the new one deleted")
}

func TestKongLicenseControllerValidation(t *testing.T) {
	scheme := Scheme(t, WithKong)
	cfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const (
		controllerName = "test-controller"
		validLicense   = "valid-license-for-testing"
	)
	licenseValidator := ctrllicense.ValidatorFunc(
		func(licenseRaw string) error {
			if licenseRaw == validLicense {
				return nil
			}
			return errors.New("invalid signature")
		},
	)
	reconciler := ctrllicense.NewKongV1Alpha1KongLicenseReconciler(
		ctrlClient,
		logr.Discard(),
		scheme,
		ctrllicense.NewLicenseCache(),
		time.Second,
		nil,
		controllerName,
		mo.None[string](),
		mo.Some(licenseValidator),
	)
	StartReconcilers(ctx, t, ctrlClient.Scheme(), cfg, reconciler)

	t.Log("Create a KongLicense and verify that it is reconciled")
	kongLicense1 := &kongv1alpha1.KongLicense{
		ObjectMeta: metav1.ObjectMeta{
			Name: "license-1",
		},
		RawLicenseString: "invalid-license-for-testing",
		Enabled:          true,
	}
	err := ctrlClient.Create(ctx, kongLicense1)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{Name: kongLicense1.Name}, kongLicense1)
		require.NoError(t, err, "Should not have error in getting the latest status")
		return findConditionInControllerStatus(
			kongLicense1, controllerName, ctrllicense.ConditionTypeLicenseValid, metav1.ConditionFalse,
		)
	}, waitTime, tickTime, "The Programmed condition for LicenseValid in controller status should be set to False")
	require.Eventually(t, func() bool {
		l, ok := reconciler.GetValidatedLicense().Get()
		if !ok {
			return false
		}
		isValid, ok := l.IsValid.Get()
		return ok && !isValid
	}, waitTime, tickTime, "Should return that license is invalid in GetLicense method")

	kongLicense1.RawLicenseString = "valid-license-for-testing"
	err = ctrlClient.Update(ctx, kongLicense1)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{Name: kongLicense1.Name}, kongLicense1)
		require.NoError(t, err, "Should not have error in getting the latest status")
		return findConditionInControllerStatus(
			kongLicense1, controllerName, ctrllicense.ConditionTypeLicenseValid, metav1.ConditionTrue,
		)
	}, waitTime, tickTime, "The Programmed condition for LicenseValid in controller status should be set to True")
	require.Eventually(t, func() bool {
		l, ok := reconciler.GetValidatedLicense().Get()
		if !ok {
			return false
		}
		isValid, ok := l.IsValid.Get()
		return ok && isValid
	}, waitTime, tickTime, "Should return that license is valid in GetLicense method")
}

func findConditionInControllerStatus(
	l *kongv1alpha1.KongLicense, controllerName string, conditionType string, conditionStatus metav1.ConditionStatus,
) bool {
	controllerStatus, found := lo.Find(
		l.Status.KongLicenseControllerStatuses, func(c kongv1alpha1.KongLicenseControllerStatus) bool {
			return c.ControllerName == controllerName
		})
	if !found {
		return false
	}
	return lo.ContainsBy(controllerStatus.Conditions, func(c metav1.Condition) bool {
		return c.Type == conditionType && c.Status == conditionStatus
	})
}
