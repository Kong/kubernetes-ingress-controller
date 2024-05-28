//go:build envtest

package envtest

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestSpecificGatewayNN(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 3 * time.Second
		tickTime = 500 * time.Millisecond
	)

	scheme := Scheme(t, WithGatewayAPI, WithKong)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		gw, gwc = deployGateway(ctx, t, ctrlClient)
		nn      = client.ObjectKeyFromObject(&gw)
		// We use the same GatewayClass here.
		gwIgnored = deployGatewayUsingGatewayClass(ctx, t, ctrlClient, gwc)
		nnIgnored = client.ObjectKeyFromObject(&gwIgnored)
	)

	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(gw.Namespace),
		WithGatewayFeatureEnabled,
		WithGatewayAPIControllers(),
		WithGatewayToReconcile(nn.String()),
	)

	const routeCount = 10
	routes := createHTTPRoutes(ctx, t, ctrlClient, gw, routeCount)
	ignoredRoutes := createHTTPRoutes(ctx, t, ctrlClient, gwIgnored, routeCount)

	t.Run("configured specific gateway gets its listener status filled", func(t *testing.T) {
		require.Eventually(t, func() bool {
			if err := ctrlClient.Get(ctx, nn, &gw); err != nil {
				t.Logf("Failed to get gateway %s: %v", nn, err)
				return false
			}
			httpListener, ok := lo.Find(gw.Status.Listeners, func(listener gatewayapi.ListenerStatus) bool {
				return listener.Name == "http"
			})
			if !ok {
				t.Logf("Failed to find http listener status in gateway %s", nn)
				return false
			}
			if httpListener.AttachedRoutes != int32(routeCount) {
				t.Logf("Expected %d routes to be attached to the http listener, got %d", routeCount, httpListener.AttachedRoutes)
				return false
			}

			return true
		}, waitTime, tickTime, "Failed to attach route to gateway")
	})
	t.Run("HTTPRoute attached to configured specific gateway gets its status parent filled", func(t *testing.T) {
		require.Eventually(t, func() bool {
			route := gatewayapi.HTTPRoute{}
			routeNN := client.ObjectKeyFromObject(routes[0])
			if err := ctrlClient.Get(ctx, routeNN, &route); err != nil {
				t.Logf("Failed to get route %s: %v", routeNN, err)
				return false
			}
			for _, p := range route.Status.Parents {
				if lo.ContainsBy(p.Conditions, func(c metav1.Condition) bool {
					return c.Type == string(gatewayapi.RouteConditionAccepted) && c.Status == metav1.ConditionTrue
				}) {
					return true
				}
			}

			return false
		}, waitTime, tickTime, "Failed to get accepted condition on HTTPRoute")
	})

	t.Run("not configured gateway does not gets its listener status filled and HTTPRoute attached to it doesn't get its status parent filled", func(t *testing.T) {
		require.Never(t, func() bool {
			t.Logf("Checking if Gateway %s is ignored (does not get status listeners filled)", nnIgnored)
			if err := ctrlClient.Get(ctx, nnIgnored, &gwIgnored); err != nil {
				t.Logf("Failed to get gateway %s: %v", nnIgnored, err)
				return true
			}

			if len(gwIgnored.Status.Listeners) != 0 {
				t.Logf("%s gateway should not have status listeners filled.", nnIgnored)
				return true
			}
			// Programmed and Accepted conditions are set by default to Unknown.
			if gatewayStatusContainsProgrammedOrAcceptedCondition(t, gwIgnored) {
				return true
			}

			t.Logf("Checking if HTTPRoute %s is ignored as expected", ignoredRoutes[0].Name)
			routeIgnored := gatewayapi.HTTPRoute{}
			routeIgnoredNN := client.ObjectKeyFromObject(ignoredRoutes[0])
			if err := ctrlClient.Get(ctx, routeIgnoredNN, &routeIgnored); err != nil {
				t.Logf("Failed to get ignored route %s: %v", routeIgnoredNN, err)
				return false
			}
			for _, p := range routeIgnored.Status.Parents {
				return lo.ContainsBy(p.Conditions, func(c metav1.Condition) bool {
					return c.Type == string(gatewayapi.RouteConditionAccepted) &&
						c.Status == metav1.ConditionTrue
				})
			}

			return false
		}, waitTime, tickTime, "Non configured gateway should not be processed and HTTPRoute should not be accepted")
	})

	t.Run("changes to gatewayclass used by not configured gateway do not get gateway's listener status filled", func(t *testing.T) {
		require.Never(t, func() bool {
			gwcOld := gwc.DeepCopy()
			gwc.Annotations = map[string]string{"foo": strconv.Itoa(time.Now().Nanosecond())}
			if err := ctrlClient.Patch(ctx, &gwc, client.MergeFrom(gwcOld)); err != nil {
				t.Logf("Failed patching gatewayclass %s: %v", client.ObjectKeyFromObject(&gwc), err)
				return true
			}

			t.Logf("Checking if Gateway %s is ignored (does not get status listeners filled)", nnIgnored)
			if err := ctrlClient.Get(ctx, nnIgnored, &gwIgnored); err != nil {
				t.Logf("Failed to get gateway %s: %v", nnIgnored, err)
				return true
			}

			if len(gwIgnored.Status.Listeners) != 0 {
				t.Logf("%s gateway should not be processed.", nnIgnored)
				return true
			}

			// Programmed and Accepted conditions are set by default to Unknown.
			if gatewayStatusContainsProgrammedOrAcceptedCondition(t, gwIgnored) {
				return true
			}

			return false
		}, waitTime, tickTime)
	})

	t.Run("changes to httproute attached to ignored gateway do not get gateway's listener status filled", func(t *testing.T) {
		require.Never(t, func() bool {
			routeIgnored := ignoredRoutes[0].DeepCopy()
			routeIgnoredOld := routeIgnored.DeepCopy()
			routeIgnored.Annotations = map[string]string{"foo": strconv.Itoa(time.Now().Nanosecond())}
			if err := ctrlClient.Patch(ctx, routeIgnored, client.MergeFrom(routeIgnoredOld)); err != nil {
				t.Logf("Failed patching gatewayclass %s: %v", client.ObjectKeyFromObject(routeIgnored), err)
				return true
			}

			t.Logf("Checking if Gateway %s is ignored (does not get status listeners filled)", gwIgnored.Name)
			if err := ctrlClient.Get(ctx, nnIgnored, &gwIgnored); err != nil {
				t.Logf("Failed to get gateway %s: %v", nnIgnored, err)
				return true
			}

			if len(gwIgnored.Status.Listeners) != 0 {
				t.Logf("%s gateway should not have status listeners filled.", nnIgnored)
				return true
			}

			// Programmed and Accepted conditions are set by default to Unknown.
			if gatewayStatusContainsProgrammedOrAcceptedCondition(t, gwIgnored) {
				return true
			}
			return false
		}, waitTime, tickTime)
	})

	t.Run("changes to referencegrant used by not configured gateway do not get gateway's listener status filled", func(t *testing.T) {
		refGrant := gatewayapi.ReferenceGrant{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "refgrant-1",
				Namespace: gwIgnored.Namespace,
			},
			Spec: gatewayapi.ReferenceGrantSpec{
				From: []gatewayapi.ReferenceGrantFrom{
					{
						Group:     gatewayapi.V1Group,
						Kind:      "Gateway",
						Namespace: gatewayapi.Namespace(gwIgnored.Namespace),
					},
				},
				To: []gatewayapi.ReferenceGrantTo{
					{
						Group: gatewayapi.V1Group,
						Kind:  "HTTPRoute",
					},
				},
			},
		}
		require.NoError(t, ctrlClient.Create(ctx, &refGrant))

		require.Never(t, func() bool {
			refGrantOld := refGrant.DeepCopy()
			refGrant.Annotations = map[string]string{"foo": strconv.Itoa(time.Now().Nanosecond())}
			if err := ctrlClient.Patch(ctx, &refGrant, client.MergeFrom(refGrantOld)); err != nil {
				t.Logf("Failed patching referenceGrant %s: %v", client.ObjectKeyFromObject(&refGrant), err)
				return true
			}

			t.Logf("Checking if Gateway %s is ignored (does not get status listeners filled)", nnIgnored)
			if err := ctrlClient.Get(ctx, nnIgnored, &gwIgnored); err != nil {
				t.Logf("Failed to get gateway %s: %v", nnIgnored, err)
				return true
			}

			if len(gwIgnored.Status.Listeners) != 0 {
				t.Logf("%s gateway should not be processed.", nnIgnored)
				return true
			}

			// Programmed and Accepted conditions are set by default to Unknown.
			if gatewayStatusContainsProgrammedOrAcceptedCondition(t, gwIgnored) {
				return true
			}

			return false
		}, waitTime, tickTime)
	})
}

func gatewayStatusContainsProgrammedOrAcceptedCondition(t *testing.T, gw gatewayapi.Gateway) bool {
	nn := client.ObjectKeyFromObject(&gw)
	for _, c := range gw.Status.Conditions {
		if c.Type == string(gatewayapi.GatewayConditionProgrammed) && c.Status != metav1.ConditionUnknown {
			t.Logf("%s gateway should not have Programmed status condition set to something else than Unknown.", nn)
			return true
		}
		if c.Type == string(gatewayapi.GatewayConditionAccepted) && c.Status != metav1.ConditionUnknown {
			t.Logf("%s gateway should not have Accepted status condition set to something else than Unknown.", nn)
			return true
		}
	}
	return false
}
