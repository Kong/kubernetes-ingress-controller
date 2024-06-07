package sendconfig

import (
	"github.com/blang/semver/v4"
)

// Config gathers parameters that are needed for sending configuration to Kong Admin APIs.
type Config struct {
	// Currently, this assumes that all underlying clients are using the same version
	// hence this shared field in here.
	Version semver.Version

	// InMemory tells whether a Kong Gateway Admin APIs should be communicated in DB-less mode.
	// It's not relevant for Konnect client.
	InMemory bool

	// Concurrency defines how many concurrent goroutines should be used when syncing configuration in DB-mode.
	Concurrency int

	// FilterTags are tags used to manage and filter entities in Kong.
	FilterTags []string

	// SkipCACertificates disables CA certificates, to avoid fighting over configuration in multi-workspace
	// environments. See https://github.com/Kong/deck/pull/617
	SkipCACertificates bool

	// EnableReverseSync indicates that reverse sync should be enabled for
	// updates to the data-plane.
	EnableReverseSync bool

	// ExpressionRoutes indicates whether to use Kong's expression routes.
	ExpressionRoutes bool

	// SanitizeKonnectConfigDumps indicates whether to sanitize Konnect config dumps.
	SanitizeKonnectConfigDumps bool

	// FallbackConfiguration indicates whether to generate fallback configuration in the case of entity
	// errors returned by the Kong Admin API.
	FallbackConfiguration bool

	// UseLastValidConfigForFallback indicates whether to use the last valid config cache to backfill broken objects
	// when recovering from a config push failure.
	UseLastValidConfigForFallback bool
}
