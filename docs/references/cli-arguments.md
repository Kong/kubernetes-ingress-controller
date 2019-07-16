# CLI Arguments

Use the following flags to tweak the behavior of Kong Ingress Controller:

|  Flag | Type | Default | Description |
|-------|------|---------|-------------|
| Configuration |
| `--ingress-class` | `string` | `kong` | Ingress class name to use to filter Ingress and custom resources when multiple Ingress Controllers are running in the same Kubernetes cluster. |
| `--election-id` | `string` | `ingress-controller-leader` | The name of ConfigMap (in the same namespace) to use to facilitate leader-election between multiple instances of the controller. |
| `--watch-namespace` | `string` | none | Namespace to watch for Ingress and custom resources. The default value of an empty string results in the controller watching for resources in all namespaces and configuring Kong accordingly. |
| `--kong-workspace` | `string` | `default` | Name of the workspace to be configured via the Ingress Controller. The workspace must be already created. |
| `--kong-url` | `string` | `http://localhost:8001` | The address of the Kong Admin URL to connect to in the format of protocol://address:port. If Kong's Admin API is not co-located with the Ingress Controller, please update it using this flag. |
| `--apiserver-host` | `string` | none | The address of the Kubernetes Apiserver to connect to in the format of protocol://address:port. If not specified, the assumption is that the binary runs inside a Kubernetes cluster and local discovery is attempted. |
| `--kubeconfig` | `string` | none | Path to kubeconfig file with authorization and master location information. |
| `--publish-service` | `string` | none | The namespaces and name of the Kubernetes Service fronting Kong Ingress Controller in the form of `namespace/name`. The controller will set the status of the Ingress resouces to match the endpoints of this service. In reference deployments, this is `kong/kong-proxy`. |
| `--profiling` | `boolean` | `true` | Enable profiling via web interface at `/debug/pprof/` |
| `--publish-status-address` | `string` | none | User customized address to be set as the status of ingress resources. The controller will set the status of the Ingress resourde to this address.|
| `--sync-period` | `duration` | `10m` | Relist and resync all the configuration every so often. |
| `--sync-rate-limit` | `float32` | `0.3` | Define the maximum per second sync frequency. |
| `--update-status` | `boolean` | `true` | If true, the controller will update the status of the Ingress resource with the endpoints of the service in `--publish-service`. |
| `--update-status-on-shutdown` | `boolean` | `true`  | If true, the controller will update the status of the Ingress resource when it being stoppped. |
| `--version` | `boolean` | `false` | Shows release information about the Kong Ingress controller and exit |
| `--help` | `boolean` | `false` | Shows this documentation on the CLI and exit. |
| Authentication|
| `--admin-header` | `string` in the form of `key:value` | none | Add a header (key:value) to every HTTP request to Kong's Admin API; it can be used multiple times to inject multiple headers |
| `--admin-ca-cert-file` | `string` | none | Path to PEM-encoded CA certificate file to verify the certificate served on Kong's Admin API |
| `--admin-tls-server-name` | `string` | none | SNI name to use for verification of the certificate presented by Kong |
| `--admin-tls-skip-verify` | `boolean` | `false` | Disable verification of TLS certificate of Kong's Admin endpoint |
| Logging |
| `--alsologtostderr` | `boolean` | `false` | Logs are written to standard error as well as to files. |
| `--log_backtrace_at` | `file:N` | none | When set to a file and line number holding a logging statement, such as -log_backtrace_at=gopherflakes.go:234 a stack trace will be written to the Info log whenever execution hits that statement. (Unlike with -vmodule, the ".go" must be present.) |
| `--log_dir` | `string` | none | Log files will be written to this directory instead of the default temporary directory. |
| `--stderrthreshold` | `string` | `ERROR` | logs at or above this threshold go to stderr |
| `--logtostderr` | `boolean` | `false` | Logs are written to standard error instead of to files. |
| `-v or --v` | `int` | `0`| Enable V-leveled logging at the specified level |
| `--vmodule` | `string` | none | The syntax of the argument is a comma-separated list of pattern=N, where pattern is a literal file name (minus the ".go" suffix) or "glob" pattern and N is a V level. For instance, -vmodule=gopher*=3 sets the V level to 3 in all Go files whose names begin "gopher".|
