package config

// MetricsAccessFilter defines the access filter function for the metrics endpoint.
type MetricsAccessFilter string

const (
	// MetricsAccessFilterOff disabled the access filter on metrics endpoint.
	MetricsAccessFilterOff MetricsAccessFilter = "off"
	// MetricsAccessFilterRBAC enables the access filter on metrics endpoint.
	// For more information consult:
	// https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/metrics/filters#WithAuthenticationAndAuthorization
	MetricsAccessFilterRBAC MetricsAccessFilter = "rbac"
)

// String returns the string representation of the MetricsFilter.
func (mf MetricsAccessFilter) String() string {
	return string(mf)
}
