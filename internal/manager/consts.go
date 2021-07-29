package manager

// -----------------------------------------------------------------------------
// Controller Manager - Constants & Vars
// -----------------------------------------------------------------------------

// HealthzPort is the default port the manager's health service listens on.
// Changing this will result in a breaking change. Existing deployments may use the literal
// port number in their liveness and readiness probes, and upgrading to a controller version
// with a changed HealthzPort will result in crash loops until users update their probe config.
// Note that there are several stock manifests in this repo that also use the literal port number. If you
// update this value, search for the old port number and update the stock manifests also.
const HealthzPort = 10254

// MetricsPort is the default port the manager's metrics service listens on.
// Similar to HealthzPort, it may be used in existing user deployment configurations, and its
// literal value is used in several stock manifests, which must be updated along with this value.
const MetricsPort = 10255

// DiagnosticsPort is the default port of the manager's diagnostics service listens on.
const DiagnosticsPort = 10256
