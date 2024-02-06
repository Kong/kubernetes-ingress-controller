package envtest

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/configuration"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
)

func TestKongLicenseController(t *testing.T) {
	scheme := Scheme(t, WithKong)
	cfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reconciler := &configuration.KongV1Alpha1KongLicenseReconciler{
		Client:         ctrlClient,
		Log:            logr.Discard(),
		Scheme:         scheme,
		LicenseCache:   configuration.NewLicenseCache(),
		ControllerName: "test",
	}
	StartReconcilers(ctx, t, ctrlClient.Scheme(), cfg, reconciler)

	const (
		waitTime            = 3 * time.Second
		tickTime            = 100 * time.Millisecond
		fullControllerName  = configuration.LicenseControllerType + "/test"
		conditionProgrammed = "Programmed"
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
