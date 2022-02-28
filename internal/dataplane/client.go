package dataplane

import (
	"context"
)

// -----------------------------------------------------------------------------
// Dataplane Client - Public Interface
// -----------------------------------------------------------------------------

type Client interface {
	// Update the data-plane by parsing the current configuring and applying
	// it to the backend API.
	Update(ctx context.Context) error
}
