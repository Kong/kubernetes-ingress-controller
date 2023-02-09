package dataplane

import (
	"context"
)

// -----------------------------------------------------------------------------
// Dataplane Client - Public Vars & Consts
// -----------------------------------------------------------------------------

const (
	// DefaultTimeoutSeconds indicates the time.Duration allowed for responses to
	// come back from the backend data-plane API.
	//
	// NOTE: the current default is based on observed latency in a CI environment using
	// the GKE cloud provider with the Kong Admin API.
	DefaultTimeoutSeconds float32 = 30.0
)

// -----------------------------------------------------------------------------
// Dataplane Client - Public Interface
// -----------------------------------------------------------------------------

type Client interface {
	// DBMode informs the caller which DB mode the data-plane has employed
	// (e.g. "off" (dbless) or "postgres").
	DBMode() string

	// Update the data-plane by parsing the current configuring and applying
	// it to the backend API.
	Update(ctx context.Context) error
}
