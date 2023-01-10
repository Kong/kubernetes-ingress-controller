package dataplane

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func TestUniqueObjects(t *testing.T) {
	t.Log("generating some objects to test the de-duplication of objects")
	ing1 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  corev1.NamespaceDefault,
			Name:       "test-ingress-1",
			Generation: 1,
		},
	}
	ing1.SetGroupVersionKind(ingGVK)
	ing2 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  corev1.NamespaceDefault,
			Name:       "test-ingress-2",
			Generation: 1,
		},
	}
	ing2.SetGroupVersionKind(ingGVK)
	ing3 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  "other-namespace",
			Name:       "test-ingress-1",
			Generation: 1,
		},
	}
	ing3.SetGroupVersionKind(ingGVK)
	ing4 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  "other-namespace",
			Name:       "test-ingress-2",
			Generation: 1,
		},
	}
	ing4.SetGroupVersionKind(ingGVK)

	testCases := []struct {
		name         string
		reportedObjs []client.Object
		failedObjs   [][]client.Object
		uniqueObjs   []client.Object
	}{
		{
			name:         "no failures",
			reportedObjs: []client.Object{ing1, ing2},
			uniqueObjs:   []client.Object{ing1, ing2},
		},
		{
			name:         "has failures",
			reportedObjs: []client.Object{ing1, ing3},
			failedObjs: [][]client.Object{
				{ing1},
				{ing4},
			},
			uniqueObjs: []client.Object{ing1, ing3, ing4},
		},
		{
			name:         "one object in multiple failures",
			reportedObjs: []client.Object{ing1, ing2},
			failedObjs: [][]client.Object{
				{ing3},
				{ing2, ing3},
			},
			uniqueObjs: []client.Object{ing1, ing2, ing3},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			translationFailures := []failures.ResourceFailure{}
			for _, failedObjs := range tc.failedObjs {
				translationFailure, err := failures.NewResourceFailure(
					"for test", failedObjs...,
				)
				require.NoError(t, err)
				translationFailures = append(translationFailures, translationFailure)
			}
			uniqueObjs := uniqueObjects(tc.reportedObjs, translationFailures)
			require.Len(t, uniqueObjs, len(tc.uniqueObjs))
			require.ElementsMatch(t, tc.uniqueObjs, uniqueObjs)
		})
	}
}

// initialized objects don't have GVK's, so we fake those for unit tests.
var (
	ingGVK = schema.GroupVersionKind{
		Group:   "networking.k8s.io",
		Version: "v1",
		Kind:    "Ingress",
	}
)

func TestClientAddressesNotifications(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		logger      = logrus.New()
		expected    = map[string]int{}
		serverCalls int32
	)

	const numberOfServers = 2

	createTestServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This test server serves as kong Admin API checking that we only get
			// as many calls as new clients requests.
			// That said: when we have 1 client with url1 and we receive a notification
			// with url1 and url2 we should only create the second client with
			// url2 and leave the existing one (for url1) in place and reuse it.

			atomic.AddInt32(&serverCalls, 1)
			n := int(atomic.LoadInt32(&serverCalls))

			if n > numberOfServers {
				t.Errorf("clients should only call out to the server %d times, but we received %d requests",
					numberOfServers, n,
				)
			}
		}))
	}

	srv := createTestServer()
	defer srv.Close()
	expected[srv.URL] = 0

	srv2 := createTestServer()
	defer srv2.Close()
	expected[srv2.URL] = 0

	client, err := NewKongClient(ctx, logger, time.Second, "", false, true, util.ConfigDumpDiagnostic{},
		sendconfig.New(ctx, logr.Discard(), []*kong.Client{}, semver.Version{}, "off", 10, []string{}),
		nil,
		"off",
		func(ctx context.Context, addr string) (*kong.Client, error) {
			num, ok := expected[addr]
			if !ok {
				return nil, fmt.Errorf("got %s which was unexpected", addr)
			}
			if num != 0 {
				return nil, fmt.Errorf("got %s more than once", addr)
			}
			expected[addr] = 1

			kongClient, err := kong.NewTestClient(lo.ToPtr(addr), &http.Client{})
			require.NoError(t, err)
			return kongClient, nil
		},
	)
	require.NoError(t, err)

	requireClientsCountEventually := func(t *testing.T, c *KongClient, n int, args ...any) {
		require.Eventually(t, func() bool {
			c.lock.RLock()
			defer c.lock.RUnlock()
			return len(c.kongConfig.Clients) == n
		}, 5*time.Second, 5*time.Millisecond, args...,
		)
	}

	requireClientsCountEventually(t, client, 0,
		"initially there should be 0 clients")

	client.Notify([]string{srv.URL})
	requireClientsCountEventually(t, client, 1,
		"after notifying about a new address we should get 1 client eventually")

	client.Notify([]string{srv.URL})
	requireClientsCountEventually(t, client, 1,
		"after notifying the same address there's no update in clients")

	client.Notify([]string{srv.URL, srv2.URL})
	requireClientsCountEventually(t, client, 2,
		"after notifying new address set including the old already existing one we get both the old and the new")

	client.Notify([]string{srv.URL, srv2.URL})
	requireClientsCountEventually(t, client, 2,
		"notifying again with the same set of URLs should not change the existing URLs")

	client.Notify([]string{srv.URL})
	requireClientsCountEventually(t, client, 1,
		"notifying again with just one URL should decrease the set of URLs to just this one")

	client.Notify([]string{})
	requireClientsCountEventually(t, client, 0)

	// We could test here notifying about srv.URL and srv2.URL again but there's
	// no data structure in the client that could notify us about a removal of
	// a client which we could use here.

	require.NoError(t, client.Shutdown(context.Background()))
	require.NoError(t, client.Shutdown(context.Background()), "closing second time shouldn't return an error")

	client.Notify([]string{})
}
