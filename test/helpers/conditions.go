package helpers

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/net"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// TODO: for now this can stay here but ideally we'd use a common package for this
// and github.com/kong/kubernetes-ingress-controller/v3/test/interla/helpers.
// At the moment we can't use the test/internal package in e.g. internal/controllers
// package because of how the internal packages work.
// This might require a separate PR that will reorder code and put it in a top
// level internal package instead of test/internal.

// HTTPRouteEventuallyContainsConditions returns a predicate function that can be
// used with assert.Eventually or require.Eventually in order to check - via the
// provided client - that the HTTPRoute with the NamespacedName as provided in
// the arguments, does indeed contain the provied conditions in the status.
func HTTPRouteEventuallyContainsConditions(ctx context.Context, t *testing.T, client ctrlclient.Client, nn k8stypes.NamespacedName, conds ...metav1.Condition) func() bool {
	return func() bool {
		t.Helper()

		var (
			ns    = nn.Namespace
			name  = nn.Name
			route = gatewayapi.HTTPRoute{}
		)

		err := client.Get(ctx, ctrlclient.ObjectKey{Namespace: ns, Name: name}, &route)
		if err != nil {
			// No point in continuing if connection is down.
			if net.IsConnectionRefused(err) {
				require.NoError(t, err)
				return false
			}
			t.Logf("Failed to get HTTPRoute: %v", err)
			return false
		}

		return lo.ContainsBy(route.Status.Parents, func(p gatewayapi.RouteParentStatus) bool {
			var count int
			for _, cond := range conds {
				contains := lo.ContainsBy(p.Conditions, func(c metav1.Condition) bool {
					return c.Type == cond.Type && c.Status == cond.Status && c.Reason == cond.Reason
				})
				if !contains {
					t.Logf("condition Type:%s, Status:%s, Reason:%s missing from route:%s/%s status",
						cond.Type, cond.Status, cond.Reason, ns, name,
					)
					return false
				}
				count++
			}
			return count == len(conds)
		})
	}
}

func HTTPRouteEventuallyNotContainsConditions(ctx context.Context, t *testing.T, client ctrlclient.Client, nn k8stypes.NamespacedName, conds ...metav1.Condition) func() bool {
	return func() bool {
		t.Helper()

		var (
			ns    = nn.Namespace
			name  = nn.Name
			route = gatewayapi.HTTPRoute{}
		)

		err := client.Get(ctx, ctrlclient.ObjectKey{Namespace: ns, Name: name}, &route)
		if err != nil {
			// No point in continuing if connection is down.
			if net.IsConnectionRefused(err) {
				require.NoError(t, err)
				return false
			}
			t.Logf("Failed to get HTTPRoute: %v", err)
			return false
		}

		return !lo.ContainsBy(route.Status.Parents, func(p gatewayapi.RouteParentStatus) bool {
			var count int
			for _, cond := range conds {
				contains := lo.ContainsBy(p.Conditions, func(c metav1.Condition) bool {
					return c.Type == cond.Type && c.Status == cond.Status && c.Reason == cond.Reason
				})
				if contains {
					t.Logf("condition Type:%s, Status:%s, Reason:%s present in route:%s/%s status",
						cond.Type, cond.Status, cond.Reason, ns, name,
					)
					return false
				}
				count++
			}
			return count == 0
		})
	}
}
