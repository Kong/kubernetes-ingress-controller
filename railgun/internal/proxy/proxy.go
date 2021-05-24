package proxy

import (
	"context"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
)

// -----------------------------------------------------------------------------
// Proxy - Public Vars
// -----------------------------------------------------------------------------

const (
	// DefaultSyncSeconds indicates the time.Duration (minimum) that will occur between
	// updates to the Kong Proxy Admin API when using the NewProxy() constructor.
	//
	// NOTE: this default was originally inherited from KIC v1.x.
	DefaultSyncSeconds float32 = 0.3

	// DefaultObjectBufferSize is the number of client.Objects that the server will buffer
	// before it starts rejecting new objects while it processes the originals.
	// If you get to the point that objects are rejected, you'll find that the
	// UpdateObject() and DeleteObject() methods will start throwing errors and you'll
	// need to retry queing the object at a later time.
	//
	// NOTE: implementations of the Proxy interface should error, not block on full buffer.
	//
	// TODO: the current default of 50 is based on a loose approximation to allow ~5mb
	//       of buffer space for client.Objects and assuming a throughput of ~50 API
	//       updates per second, but in the future we may want to make this configurable,
	//       provide metrics for it, and furthermore automate detecting good values for it.
	//       depending on configuration and/or available system memory and the amount of
	//       throughput (in Kubernetes object updates) that the API is meant to handle.
	DefaultObjectBufferSize = 500
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
	// UpdateObject accepts a Kubernetes controller-runtime client.Object and adds/updates that to the configuration cache.
	// It will be asynchronously converted into the upstream Kong DSL and applied to the Kong Admin API.
	// A status will later be added to the object whether the configuration update succeeds or fails.
	UpdateObject(obj client.Object) error

	// DeleteObject accepts a Kubernetes controller-runtime client.Object and removes it from the configuration cache.
	// The delete action will asynchronously be converted to Kong DSL and applied to the Kong Admin API.
	// A status will later be added to the object whether the configuration update succeeds or fails.
	DeleteObject(obj client.Object) error
}

// KongUpdater is a type of function that describes how to provide updates to the Kong Admin API
// and implementations will report the configuration SHA that results from any update performed.
type KongUpdater func(ctx context.Context,
	lastConfigSHA []byte,
	cache *store.CacheStores,
	ingressClassName string,
	deprecatedLogger logrus.FieldLogger,
	kongConfig sendconfig.Kong,
	enableReverseSync bool) ([]byte, error)
