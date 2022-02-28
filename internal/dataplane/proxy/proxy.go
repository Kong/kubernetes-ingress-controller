package proxy

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Proxy - Public Vars
// -----------------------------------------------------------------------------

const (
	// DefaultProxyTimeoutSeconds indicates the time.Duration allowed for responses to
	// come back from the backend proxy API.
	//
	// NOTE: the current default is based on observed latency in a CI environment using
	// the GKE cloud provider.
	DefaultProxyTimeoutSeconds float32 = 10.0

	// DefaultSyncSeconds indicates the time.Duration (minimum) that will occur between
	// updates to the Kong Proxy Admin API when using the NewProxy() constructor.
	// this 1s default was based on local testing wherein it appeared sub-second updates
	// to the Admin API could be problematic (or at least operate differently) based on
	// which storage backend was in use (i.e. "dbless", "postgres"). This is a workaround
	// for improvements we still need to investigate upstream.
	//
	// See Also: https://github.com/Kong/kubernetes-ingress-controller/issues/1398
	DefaultSyncSeconds float32 = 3.0
)

// -----------------------------------------------------------------------------
// Proxy - Public Types
// -----------------------------------------------------------------------------

// Proxy represents the Kong Proxy from the perspective of Kubernetes allowing
// callers to update and remove Kubernetes objects in the backend proxy without
// having to understand or be aware of Kong DSLs or how types are converted between
// Kubernetes and the Kong Admin API.
//
// NOTE: implementations of this interface are: threadsafe, non-blocking
type Proxy interface {
	// IsReady returns true if the proxy is considered ready.
	// A ready proxy has configuration available and can handle traffic.
	IsReady() bool

	manager.Runnable
	manager.LeaderElectionRunnable
}

// KongUpdater is a type of function that describes how to provide updates to the Kong Admin API
// and implementations will report the configuration SHA that results from any update performed.
type KongUpdater func(ctx context.Context,
	lastConfigSHA []byte,
	cache *store.CacheStores,
	ingressClassName string,
	deprecatedLogger logrus.FieldLogger,
	kongConfig sendconfig.Kong,
	enableReverseSync bool,
	diagnostic util.ConfigDumpDiagnostic,
	proxyRequestTimeout time.Duration,
	promMetrics *metrics.CtrlFuncMetrics) ([]byte, error)
