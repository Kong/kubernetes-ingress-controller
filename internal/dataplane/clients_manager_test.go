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

func (cf clientFactoryWithExpected) CreateAdminAPIClient(ctx context.Context, address string) (*adminapi.Client, error) {
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
			clientAddresses := lo.Map(c.Clients(), func(cl *adminapi.Client, _ int) string {
				return cl.BaseRootURL()
			})
			return slices.Equal(addresses, clientAddresses)
		}, time.Second, time.Millisecond, args...,
		)
	}

	requireClientsCountEventually(t, manager, []string{"localhost:8083"},
		"initially there should be the initial client")

	manager.Notify([]string{srv.URL})
	requireClientsCountEventually(t, manager, []string{srv.URL},
		"after notifying about a new address we should get 1 client eventually")

	manager.Notify([]string{srv.URL})
	requireClientsCountEventually(t, manager, []string{srv.URL},
		"after notifying the same address there's no update in clients")

	manager.Notify([]string{srv.URL, srv2.URL})
	requireClientsCountEventually(t, manager, []string{srv.URL, srv2.URL},
		"after notifying new address set including the old already existing one we get both the old and the new")

	manager.Notify([]string{srv.URL, srv2.URL})
	requireClientsCountEventually(t, manager, []string{srv.URL, srv2.URL},
		"notifying again with the same set of URLs should not change the existing URLs")

	manager.Notify([]string{srv.URL})
	requireClientsCountEventually(t, manager, []string{srv.URL},
		"notifying again with just one URL should decrease the set of URLs to just this one")

	manager.Notify([]string{})
	requireClientsCountEventually(t, manager, []string{})

	// We could test here notifying about srv.URL and srv2.URL again but there's
	// no data structure in the manager that could notify us about a removal of
	// a manager which we could use here.

	cancel()
	require.NotPanics(t, func() { manager.Notify([]string{}) }, "notifying about new clients after manager has been shut down shouldn't panic")
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

	clients := manager.Clients()
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
		manager.adjustKongClients([]string{"localhost:8080", "localhost:8081"})
		requireNoExpectedCallsLeftEventually(t)
	})

	t.Run("1 addresses, no new client", func(t *testing.T) {
		// Change expected addresses
		cf.expected = map[string]bool{}
		// there is address contained in the notification but a client for that
		// address already exists, client creator should not be called
		manager.adjustKongClients([]string{"localhost:8080"})
		requireNoExpectedCallsLeftEventually(t)
	})

	t.Run("2 addresses, 1 new client", func(t *testing.T) {
		// Change expected addresses
		cf.expected = map[string]bool{"localhost:8081": true}
		// there are 2 addresses contained in the notification but only 1 is new
		// hence the client creator should be called only once
		manager.adjustKongClients([]string{"localhost:8080", "localhost:8081"})
		requireNoExpectedCallsLeftEventually(t)
	})
}

func TestNewAdminAPIClientsManager_NoInitialClientsDisallowed(t *testing.T) {
	cf := &clientFactoryWithExpected{t: t}
	_, err := NewAdminAPIClientsManager(context.Background(), logrus.New(), nil, cf)
	require.Error(t, err)
}
