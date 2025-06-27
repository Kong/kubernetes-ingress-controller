package types

import "github.com/kong/kubernetes-telemetry/pkg/types"

type (
	// Payload is an alias for types.ProviderReport.
	// It represents the structure of the data used in telemetry reporting.
	// This alias allows for easier reference and usage of the ProviderReport type
	// within the telemetry package.
	Payload = types.ProviderReport

	// PayloadKey is an alias for types.ProviderReportKey.
	// It represents the key type used in the ProviderReport map.
	// This alias simplifies the usage of ProviderReportKey within the telemetry package.
	PayloadKey = types.ProviderReportKey

	// PayloadCustomizer is a function type that takes a Payload as input
	// and returns a modified Payload. This allows for dynamic adjustments to the
	// payload data based on specific requirements or conditions.
	PayloadCustomizer = func(payload Payload) Payload
)
