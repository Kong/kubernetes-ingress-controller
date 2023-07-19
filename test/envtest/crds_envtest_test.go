//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// TestMissingCRDsDontCrashTheManager ensures that in case of missing CRDs installation in the cluster, specific
// controllers are disabled, this fact is properly logged, and the manager does not crash.
func TestMissingCRDsDontCrashTheManager(t *testing.T) {
	emptyScheme := runtime.NewScheme()
	envcfg := Setup(t, emptyScheme)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	loggerHook := RunManager(ctx, t, envcfg, func(cfg *manager.Config) {
		// Reducing controllers' cache synchronisation timeout in order to trigger the possible sync timeout quicker.
		// It's a regression test for https://github.com/Kong/gateway-operator/issues/326.
		cfg.CacheSyncTimeout = time.Millisecond * 500
	})

	require.Eventually(t, func() bool {
		gvrs := []schema.GroupVersionResource{
			{
				Group:    kongv1beta1.GroupVersion.Group,
				Version:  kongv1beta1.GroupVersion.Version,
				Resource: "udpingresses",
			},
			{
				Group:    kongv1beta1.GroupVersion.Group,
				Version:  kongv1beta1.GroupVersion.Version,
				Resource: "tcpingresses",
			},
			{
				Group:    kongv1.GroupVersion.Group,
				Version:  kongv1.GroupVersion.Version,
				Resource: "kongingresses",
			},
			{
				Group:    kongv1alpha1.GroupVersion.Group,
				Version:  kongv1alpha1.GroupVersion.Version,
				Resource: "ingressclassparameterses",
			},
			{
				Group:    kongv1.GroupVersion.Group,
				Version:  kongv1.GroupVersion.Version,
				Resource: "kongplugins",
			},
			{
				Group:    kongv1.GroupVersion.Group,
				Version:  kongv1.GroupVersion.Version,
				Resource: "kongconsumers",
			},
			{
				Group:    kongv1beta1.GroupVersion.Group,
				Version:  kongv1beta1.GroupVersion.Version,
				Resource: "kongconsumergroups",
			},
			{
				Group:    kongv1.GroupVersion.Group,
				Version:  kongv1.GroupVersion.Version,
				Resource: "kongclusterplugins",
			},
		}

		for _, gvr := range gvrs {
			expectedLog := fmt.Sprintf("Disabling controller for Group=%s, Resource=%s due to missing CRD", gvr.GroupVersion(), gvr.Resource)
			if !lo.ContainsBy(loggerHook.AllEntries(), func(entry *logrus.Entry) bool {
				return strings.Contains(entry.Message, expectedLog)
			}) {
				t.Logf("expected log not found: %s", expectedLog)
				return false
			}
		}
		return true
	}, time.Minute, time.Millisecond*500)
}
