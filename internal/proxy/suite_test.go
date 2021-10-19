package proxy

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	kongt "github.com/kong/kubernetes-testing-framework/pkg/utils/kong"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Test Vars
// -----------------------------------------------------------------------------

var (
	logger logrus.FieldLogger

	fakeKongConfig   sendconfig.Kong
	fakeKongAdminAPI *kongt.FakeAdminAPIServer

	fakeK8sClient client.Client
)

// -----------------------------------------------------------------------------
// Test Main
// -----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	setupVars()
	code := m.Run()
	os.Exit(code)
}

// -----------------------------------------------------------------------------
// Test Main - Helper Functions
// -----------------------------------------------------------------------------

func setupVars() {
	var err error

	// setup logging and other general configurations
	logger, err = util.MakeLogger("debug", "json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create test logger: %v\n", err)
		os.Exit(10)
	}

	// setup kubernetes related configurations
	fakeK8sClient = fake.NewClientBuilder().Build()

	// setup kong proxy related configurations
	fakeKongAdminAPI, err = kongt.NewFakeAdminAPIServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup Kong Proxy testing environment: %v\n", err)
		os.Exit(10)
	}
	fakeKongConfig.Client = fakeKongAdminAPI.KongClient
}

// this mock will stop us from needing a functioning Kong proxy and simply
// assumes that all updates to the Kong Admin API succeed. Test which want
// to test failure conditions will need to add their own mock, and then put
// this one back when they're done.
//
// NOTE: as these tests grow, we can use the kongt.FakeAdminAPIServer implementation to properly
//       mock out the requests (this is used above to ensure that tests can properly initialize a Proxy instance)
//       instead of the always-succeed functionality we use currently.
var mockKongAdmin KongUpdater = func(ctx context.Context,
	lastConfigSHA []byte,
	cache *store.CacheStores,
	ingressClassName string,
	deprecatedLogger logrus.FieldLogger,
	kongConfig sendconfig.Kong,
	enableReverseSync bool,
	diagnostic util.ConfigDumpDiagnostic,
	proxyRequestTimeout time.Duration,
	promMetrics *metrics.CtrlFuncMetrics,
) ([]byte, error) {
	fakeKongAdminUpdateCount(1)
	return lastConfigSHA, nil
}

// these globs are for threadsafety and tracking of the fakeKongAdminUpdateCount() function,
// use that function directly, don't use these vars.
var (
	countLock   = &sync.RWMutex{}
	updateCount int
)

// fakeKongAdminUpdateCount keeps track of the number of times the Kong Admin API has
// received updates from the proxy cache server during the runtime of the tests. This
// can be useful for simple checks to ensure updates ran, or to validate counts of the
// number of updates run over time when cache server resolution configuration options
// are tweaked to change throughput and test performance.
//
// If newcounts are provided the function appends the count with those provided counts first.
func fakeKongAdminUpdateCount(newcounts ...int) int {
	countLock.Lock()
	defer countLock.Unlock()
	if len(newcounts) < 1 {
		return updateCount
	}
	if len(newcounts) == 1 && newcounts[0] == 0 {
		updateCount = 0
		return 0
	}
	for _, count := range newcounts {
		updateCount = updateCount + count
	}
	return updateCount
}
