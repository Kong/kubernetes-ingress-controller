package dataplane

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

// clientFactoryWithExpected implements ClientFactory interface and can be used
// in tests to assert which clients have been created and signal failure if:
// - client for an unexpected address gets created
// - client which already got created was tried to be created second time.
type clientFactoryWithExpected struct {
	expected map[string]bool
	t        *testing.T
}

func (cf clientFactoryWithExpected) CreateAdminAPIClient(_ context.Context, address string) (*adminapi.Client, error) {
	stillExpecting, ok := cf.expected[address]
	if !ok {
		cf.t.Errorf("got %s which was unexpected", address)
		return nil, fmt.Errorf("got %s which was unexpected", address)
	}
	if !stillExpecting {
		cf.t.Errorf("got %s more than once", address)
		return nil, fmt.Errorf("got %s more than once", address)
	}
	cf.expected[address] = false

	return adminapi.NewTestClient(address)
}

func (cf clientFactoryWithExpected) AssertExpectedCalls() {
	for _, addr := range cf.ExpectedCallsLeft() {
		cf.t.Errorf("%s client expected to be called, but wasn't", addr)
	}
}

func (cf clientFactoryWithExpected) ExpectedCallsLeft() []string {
	var notCalled []string
	for addr, stillExpected := range cf.expected {
		if stillExpected {
			notCalled = append(notCalled, addr)
		}
	}
	return notCalled
}

func TestClientAddressesNotifications(t *testing.T) {
	var (
		logger      = logrus.New()
		expected    = map[string]bool{}
		serverCalls int32
	)

	const numberOfServers = 2

	createTestServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This test server serves as kong Admin API checking that we only get
			// as many calls as new clients requests.
			// That said: when we have 1 manager with url1 and we receive a notification
			// with url1 and url2 we should only create the second manager with
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initialClient, err := adminapi.NewTestClient("localhost:8083")
	require.NoError(t, err)
	manager, err := NewAdminAPIClientsManager(
		ctx,
		logger,
		[]*adminapi.Client{initialClient},
		testClientFactoryWithExpected,
	)
	require.NoError(t, err)
	require.NotNil(t, manager)
	manager.RunNotifyLoop()
	<-manager.Running()

	defer testClientFactoryWithExpected.AssertExpectedCalls()

	requireClientsCountEventually := func(t *testing.T, c *AdminAPIClientsManager, addresses []string, args ...any) {
		require.Eventually(t, func() bool {
			clientAddresses := lo.Map(c.AllClients(), func(cl *adminapi.Client, _ int) string {
				return cl.BaseRootURL()
			})
			return slices.Equal(addresses, clientAddresses)
		}, time.Second, time.Millisecond, args...,
		)
	}

	requireClientsCountEventually(t, manager, []string{"localhost:8083"},
		"initially there should be the initial client")

	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(srv.URL)})
	requireClientsCountEventually(t, manager, []string{srv.URL},
		"after notifying about a new address we should get 1 client eventually")

	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(srv.URL)})
	requireClientsCountEventually(t, manager, []string{srv.URL},
		"after notifying the same address there's no update in clients")

	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(srv.URL), testDiscoveredAdminAPI(srv2.URL)})
	requireClientsCountEventually(t, manager, []string{srv.URL, srv2.URL},
		"after notifying new address set including the old already existing one we get both the old and the new")

	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(srv.URL), testDiscoveredAdminAPI(srv2.URL)})
	requireClientsCountEventually(t, manager, []string{srv.URL, srv2.URL},
		"notifying again with the same set of URLs should not change the existing URLs")

	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(srv.URL)})
	requireClientsCountEventually(t, manager, []string{srv.URL},
		"notifying again with just one URL should decrease the set of URLs to just this one")

	manager.Notify([]adminapi.DiscoveredAdminAPI{})
	requireClientsCountEventually(t, manager, []string{})

	// We could test here notifying about srv.URL and srv2.URL again but there's
	// no data structure in the manager that could notify us about a removal of
	// a manager which we could use here.

	cancel()
	require.NotPanics(t, func() { manager.Notify([]adminapi.DiscoveredAdminAPI{}) }, "notifying about new clients after manager has been shut down shouldn't panic")
}

func TestClientAdjustInternalClientsAfterNotification(t *testing.T) {
	var (
		ctx    = context.Background()
		logger = logrus.New()
	)

	cf := &clientFactoryWithExpected{
		t: t,
	}

	// Initial client is expected to be replaced later on.
	testClient, err := adminapi.NewTestClient("localhost:8083")
	require.NoError(t, err)
	manager, err := NewAdminAPIClientsManager(ctx, logger, []*adminapi.Client{testClient}, cf)
	require.NoError(t, err)
	require.NotNil(t, manager)
	manager.RunNotifyLoop()
	<-manager.Running()

	clients := manager.AllClients()
	require.Len(t, clients, 1)
	require.Equal(t, "localhost:8083", clients[0].BaseRootURL())

	requireNoExpectedCallsLeftEventually := func(t *testing.T) {
		require.Eventually(t, func() bool {
			return len(cf.ExpectedCallsLeft()) == 0
		}, time.Second, time.Millisecond)
	}

	t.Run("2 new clients", func(t *testing.T) {
		// Change expected addresses
		cf.expected = map[string]bool{"localhost:8080": true, "localhost:8081": true}
		// there are 2 addresses contained in the notification of which 2 are new
		// and client creator should be called exactly 2 times
		manager.adjustGatewayClients([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI("localhost:8080"), testDiscoveredAdminAPI("localhost:8081")})
		requireNoExpectedCallsLeftEventually(t)
	})

	t.Run("1 addresses, no new client", func(t *testing.T) {
		// Change expected addresses
		cf.expected = map[string]bool{}
		// there is address contained in the notification but a client for that
		// address already exists, client creator should not be called
		manager.adjustGatewayClients([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI("localhost:8080")})
		requireNoExpectedCallsLeftEventually(t)
	})

	t.Run("2 addresses, 1 new client", func(t *testing.T) {
		// Change expected addresses
		cf.expected = map[string]bool{"localhost:8081": true}
		// there are 2 addresses contained in the notification but only 1 is new
		// hence the client creator should be called only once
		manager.adjustGatewayClients([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI("localhost:8080"), testDiscoveredAdminAPI("localhost:8081")})
		requireNoExpectedCallsLeftEventually(t)
	})
}

func TestNewAdminAPIClientsManager_NoInitialClientsDisallowed(t *testing.T) {
	cf := &clientFactoryWithExpected{t: t}
	_, err := NewAdminAPIClientsManager(context.Background(), logrus.New(), nil, cf)
	require.Error(t, err)
}

func TestAdminAPIClientsManager_NotRunningNotifyLoop(t *testing.T) {
	t.Parallel()

	testClient, err := adminapi.NewTestClient("localhost:8080")
	require.NoError(t, err)
	m, err := NewAdminAPIClientsManager(
		context.Background(),
		logrus.New(),
		[]*adminapi.Client{testClient},
		&clientFactoryWithExpected{t: t},
	)
	require.NoError(t, err)

	select {
	case <-m.Running():
		t.Error("expected manager to not run without explicitly running it with RunNotifyLoop method")
	case <-time.After(time.Millisecond * 100):
	}
}

func TestAdminAPIClientsManager_Clients(t *testing.T) {
	t.Parallel()

	testClient, err := adminapi.NewTestClient("localhost:8080")
	require.NoError(t, err)
	m, err := NewAdminAPIClientsManager(
		context.Background(),
		logrus.New(),
		[]*adminapi.Client{testClient},
		&clientFactoryWithExpected{t: t},
	)
	require.NoError(t, err)
	require.Len(t, m.GatewayClients(), 1, "expecting one initial client")
	require.Equal(t, m.GatewayClientsCount(), 1, "expecting one initial client")
	require.Len(t, m.AllClients(), 1, "expecting one initial client")

	konnectTestClient, err := adminapi.NewTestClient("https://us.api.konghq.tech")
	require.NoError(t, err)
	m.SetKonnectClient(konnectTestClient)
	require.Len(t, m.GatewayClients(), 1, "konnect client should not be returned from GatewayClients")
	require.Equal(t, m.GatewayClientsCount(), 1, "konnect client should not be counted in GatewayClientsCount")
	require.Len(t, m.AllClients(), 2, "konnect client should be returned from AllClients")
}

func TestAdminAPIClientsManager_GatewayClientsFromNotificationsAreExpectedToHavePodRef(t *testing.T) {
	t.Parallel()

	cf := &clientFactoryWithExpected{t: t, expected: map[string]bool{"http://10.0.0.1:8080": true}}
	testClient, err := adminapi.NewTestClient("http://localhost:8080")
	require.NoError(t, err)
	m, err := NewAdminAPIClientsManager(
		context.Background(),
		logrus.New(),
		[]*adminapi.Client{testClient},
		cf,
	)
	require.NoError(t, err)
	m.RunNotifyLoop()

	m.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI("http://10.0.0.1:8080")})

	require.Eventually(t, func() bool {
		gwClients := m.GatewayClients()
		if len(gwClients) != 1 {
			t.Log("there's no gateway clients...")
			return false
		}
		_, ok := gwClients[0].PodReference()
		if !ok {
			t.Log("there's no pod reference attached")
		}
		return ok
	}, time.Second, time.Millisecond)
}

func TestAdminAPIClientsManager_SubscribeToGatewayClientsChanges(t *testing.T) {
	t.Parallel()

	cf := &clientFactoryWithExpected{t: t, expected: map[string]bool{
		"http://10.0.0.1:8080": true,
		"http://10.0.0.2:8080": true,
	}}
	testClient, err := adminapi.NewTestClient("http://localhost:8080")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	m, err := NewAdminAPIClientsManager(ctx, logrus.New(), []*adminapi.Client{testClient}, cf)
	require.NoError(t, err)

	t.Run("no notify loop running should return false when subscribing", func(t *testing.T) {
		ch, ok := m.SubscribeToGatewayClientsChanges()
		require.Nil(t, ch)
		require.Falsef(t, ok, "expected no subscription to be created because no notification loop is running")
	})

	m.RunNotifyLoop()

	t.Run("when notification loop is running subscription should be created", func(t *testing.T) {
		ch, ok := m.SubscribeToGatewayClientsChanges()
		require.NotNil(t, ch)
		require.True(t, ok)

		m.Notify([]adminapi.DiscoveredAdminAPI{
			testDiscoveredAdminAPI("http://10.0.0.1:8080"),
			testDiscoveredAdminAPI("http://10.0.0.2:8080"),
		})

		select {
		case <-ch:
			require.Len(t, m.GatewayClients(), 2, "expected to get 2 clients after the update")
		case <-time.After(time.Second):
			t.Error("did not receive notification after gateway clients changes")
		}
	})

	t.Run("when multiple subscriptions are created, each of them should receive notifications", func(t *testing.T) {
		sub1, ok := m.SubscribeToGatewayClientsChanges()
		require.NotNil(t, sub1)
		require.True(t, ok)

		sub2, ok := m.SubscribeToGatewayClientsChanges()
		require.NotNil(t, sub2)
		require.True(t, ok)

		m.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI("http://10.0.0.1:8080")})

		select {
		case <-sub2:
			require.Len(t, m.GatewayClients(), 1, "expected to get 1 client after the update")
		case <-time.After(time.Second):
			t.Error("did not receive notification after gateway clients changes")
		}
		select {
		case <-sub1:
			require.Len(t, m.GatewayClients(), 1, "expected to get 1 client after the update")
		case <-time.After(time.Second):
			t.Error("did not receive notification after gateway clients changes")
		}
	})

	t.Run("when the context gets cancelled, subscriber channel gets closed", func(t *testing.T) {
		ch, ok := m.SubscribeToGatewayClientsChanges()
		require.NotNil(t, ch)
		require.True(t, ok)

		cancel()

		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Error("subscription channel wasn't closed after cancelling the context")
		}
	})

	t.Run("when the context is cancelled, subscriptions cannot be created", func(t *testing.T) {
		ch, ok := m.SubscribeToGatewayClientsChanges()
		require.Nil(t, ch)
		require.False(t, ok)
	})
}

func TestAdminAPIClientsManager_ConcurrentNotify(t *testing.T) {
	t.Parallel()

	cf := &clientFactoryWithExpected{t: t, expected: map[string]bool{
		"http://10.0.0.1:8080": true,
	}}
	testClient, err := adminapi.NewTestClient("http://10.0.0.1:8080")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m, err := NewAdminAPIClientsManager(ctx, logrus.New(), []*adminapi.Client{testClient}, cf)
	require.NoError(t, err)
	m.RunNotifyLoop()

	var receivedNotificationsCount atomic.Uint32
	ch, ok := m.SubscribeToGatewayClientsChanges()
	require.NotNil(t, ch)
	require.True(t, ok)

	// Run subscriber worker in a separate goroutine to consume notifications.
	go func() {
		for {
			select {
			case <-ch:
				// Call GatewayClients() here to make sure that we can access the clients safely
				// from the subscriber goroutine without causing a deadlock in the notify loop.
				require.Len(t, m.GatewayClients(), 1, "expected to get 1 client")
				receivedNotificationsCount.Add(1)
			case <-ctx.Done():
				t.Log("Test is done, stopping subscriber worker")
				return
			}
		}
	}()

	// Run multiple notifiers in parallel to make sure that Notify is safe for concurrent use.
	for i := 0; i < 10; i++ {
		go func() {
			m.Notify([]adminapi.DiscoveredAdminAPI{
				testDiscoveredAdminAPI("http://10.0.0.1:8080"),
			})
		}()
	}

	require.Eventually(t, func() bool {
		if receivedNotificationsCount.Load() != 10 {
			t.Logf("Received %d notifications, expected 10, waiting...", receivedNotificationsCount.Load())
			return false
		}
		return true
	}, time.Second, time.Millisecond, "expected to receive 10 notifications")
}

func testDiscoveredAdminAPI(address string) adminapi.DiscoveredAdminAPI {
	return adminapi.DiscoveredAdminAPI{
		Address: address,
		PodRef:  types.NamespacedName{Name: "pod-1", Namespace: "ns"},
	}
}
