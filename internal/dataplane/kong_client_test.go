package dataplane

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
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

// clientFactoryWithExpected implements ClientFactory interface and can be used
// in tests to assert which clients have been created and signal failure if:
// - client for an unexpected address gets created
// - client which already got created was tried to be created second time.
type clientFactoryWithExpected struct {
	expected map[string]bool
	t        *testing.T
}

func (cf clientFactoryWithExpected) CreateAdminAPIClient(ctx context.Context, address string) (adminapi.Client, error) {
	stillExpecting, ok := cf.expected[address]
	if !ok {
		cf.t.Errorf("got %s which was unexpected", address)
		return adminapi.Client{}, fmt.Errorf("got %s which was unexpected", address)
	}
	if !stillExpecting {
		cf.t.Errorf("got %s more than once", address)
		return adminapi.Client{}, fmt.Errorf("got %s more than once", address)
	}
	cf.expected[address] = false

	kongClient, err := kong.NewTestClient(lo.ToPtr(address), &http.Client{})
	if err != nil {
		return adminapi.Client{}, err
	}

	return adminapi.NewClient(kongClient), nil
}

func (cf clientFactoryWithExpected) AssertExpectedCalls() {
	for addr, stillExpected := range cf.expected {
		if stillExpected {
			cf.t.Errorf("%s client expected to be called, but wasn't", addr)
		}
	}
}

func TestClientAddressesNotifications(t *testing.T) {
	var (
		ctx         = context.Background()
		logger      = logrus.New()
		expected    = map[string]bool{}
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
	expected[srv.URL] = true

	srv2 := createTestServer()
	defer srv2.Close()
	expected[srv2.URL] = true

	testClientFactoryWithExpected := clientFactoryWithExpected{
		expected: expected,
		t:        t,
	}
	client, err := NewKongClient(ctx, logger, time.Second, "", false, true, util.ConfigDumpDiagnostic{},
		sendconfig.New(ctx, logr.Discard(), []adminapi.Client{},
			sendconfig.Config{
				InMemory:    true,
				Concurrency: 10,
			},
		),
		nil,
		"off",
		testClientFactoryWithExpected,
	)
	require.NoError(t, err)
	defer testClientFactoryWithExpected.AssertExpectedCalls()

	requireClientsCountEventually := func(t *testing.T, c *KongClient, addresses []string, args ...any) {
		require.Eventually(t, func() bool {
			c.lock.RLock()
			defer c.lock.RUnlock()
			clientAddresses := lo.Map(c.kongConfig.Clients, func(cl adminapi.Client, _ int) string {
				return cl.BaseRootURL()
			})
			return slices.Equal(addresses, clientAddresses)
		}, time.Second, time.Millisecond, args...,
		)
	}

	requireClientsCountEventually(t, client, []string{},
		"initially there should be 0 clients")

	client.Notify([]string{srv.URL})
	requireClientsCountEventually(t, client, []string{srv.URL},
		"after notifying about a new address we should get 1 client eventually")

	client.Notify([]string{srv.URL})
	requireClientsCountEventually(t, client, []string{srv.URL},
		"after notifying the same address there's no update in clients")

	client.Notify([]string{srv.URL, srv2.URL})
	requireClientsCountEventually(t, client, []string{srv.URL, srv2.URL},
		"after notifying new address set including the old already existing one we get both the old and the new")

	client.Notify([]string{srv.URL, srv2.URL})
	requireClientsCountEventually(t, client, []string{srv.URL, srv2.URL},
		"notifying again with the same set of URLs should not change the existing URLs")

	client.Notify([]string{srv.URL})
	requireClientsCountEventually(t, client, []string{srv.URL},
		"notifying again with just one URL should decrease the set of URLs to just this one")

	client.Notify([]string{})
	requireClientsCountEventually(t, client, []string{})

	// We could test here notifying about srv.URL and srv2.URL again but there's
	// no data structure in the client that could notify us about a removal of
	// a client which we could use here.

	require.NoError(t, client.Shutdown(context.Background()), "closing shouldn't return an error")
	require.NoError(t, client.Shutdown(context.Background()), "closing second time shouldn't return an error")

	require.NotPanics(t, func() { client.Notify([]string{}) }, "notifying about new clients after client has been shut down shouldn't panic")
}

func TestClientAdjustInternalClientsAfterNotification(t *testing.T) {
	var (
		ctx    = context.Background()
		logger = logrus.New()
	)

	cf := &clientFactoryWithExpected{
		t: t,
	}
	defer cf.AssertExpectedCalls()
	client, err := NewKongClient(ctx, logger, time.Second, "", false, true, util.ConfigDumpDiagnostic{},
		sendconfig.New(ctx, logr.Discard(), []adminapi.Client{},
			sendconfig.Config{
				InMemory:    true,
				Concurrency: 10,
			},
		),
		nil,
		"off",
		cf,
	)
	require.NoError(t, err)
	require.NotNil(t, client)

	t.Run("2 new clients", func(t *testing.T) {
		// Change expected addresses
		cf.expected = map[string]bool{"localhost:8080": true, "localhost:8081": true}

		// there are 2 addresses contained in the notification of which 2 are new
		// and client creator should be called exactly 2 times
		client.adjustKongClients(ctx, []string{"localhost:8080", "localhost:8081"})
	})

	t.Run("1 addresses, no new client", func(t *testing.T) {
		// Change expected addresses
		cf.expected = map[string]bool{"localhost:8080": true}
		// there is address contained in the notification but a client for that
		// address already exists, client creator should not be called
		client.adjustKongClients(ctx, []string{"localhost:8080"})
	})

	t.Run("2 addresses, 1 new client", func(t *testing.T) {
		// Change expected addresses
		cf.expected = map[string]bool{"localhost:8080": true, "localhost:8081": true}
		// there are 2 addresses contained in the notification but only 1 is new
		// hence the client creator should be called only once
		client.adjustKongClients(ctx, []string{"localhost:8080", "localhost:8081"})
	})
}
