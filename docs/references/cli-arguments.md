# CLI Arguments

Various settings and configurations of the controller can be tweaked
using CLI flags.

## Environment variables

Each flag defined in the table below can also be configured using
an environment variable. The name of the environment variable is `CONTROLLER_`
string followed by the name of flag in uppercase.

For example, `--ingress-class` can be configured using the following
environment variable:

```
CONTROLLER_INGRESS_CLASS=kong-foobar
```

It is recommended that all the configuration is done via environment variables
and not CLI flags.

## Flags

Following table describes all the flags that are available:

|  Flag | Type | Default | Description |
|-------|------|---------|-------------|
| --admin-ca-cert-file                 |`string`   | none                            | DEPRECATED, use `--kong-admin-ca-cert-file`|
| --admin-header                       |`string`   | none                            | DEPRECATED, use `--kong-admin-header`|
| --admin-tls-server-name              |`string`   | none                            | DEPRECATED, use `--kong-admin-tls-server-name`|
| --admin-tls-skip-verify              |`boolean`  | none                            | DEPRECATED, use `--kong-admin-tls-skip-verify`|
| --admission-webhook-cert-file        |`string`   | `/admission-webhook/tls.crt`    | Path to the PEM-encoded certificate file for TLS handshake.|
| --admission-webhook-key-file         |`string`   | `/admission-webhook/tls.key`    | Path to the PEM-encoded private key file for TLS handshake.|
| --admission-webhook-cert             |`string`   | none                            | PEM-encoded certificate string for TLS handshake.|
| --admission-webhook-key              |`string`   | none                            | PEM-encoded private key string for TLS handshake.|
| --admission-webhook-listen           |`string`   | `off`                           | The address to start admission controller on (ip:port). Setting it to 'off' disables the admission controller.|
| --alsologtostderr                    |`boolean`  | `false`                         | Logs are written to standard error as well as to files.|
| --anonymous-reports                  |`string`   | `true`                          | Send anonymized usage data to help improve Kong.|
| --apiserver-host                     |`string`   | none                            | The address of the Kubernetes Apiserver to connect to in the format of protocol://address:port, e.g., "http://localhost:8080. If not specified, the assumption is that the binary runs inside a Kubernetes cluster and local discovery is attempted.|
| --disable-ingress-extensionsv1beta1  |`boolean`  | `false`                         | Disable processing Ingress resources with apiVersion `extensions/v1beta1`.|
| --disable-ingress-networkingv1beta1  |`boolean`  | `false`                         | Disable processing Ingress resources with apiVersion `networking/v1beta1`.|
| --disable-ingress-networkingv1       |`boolean`  | `false`                         | Disable processing Ingress resources with apiVersion `networking/v1`.|
| --election-id                        |`string`   | `ingress-controller-leader`     | The name of ConfigMap (in the same namespace) to use to facilitate leader-election between multiple instances of the controller.|
| --ingress-class                      |`string`   | `kong`                          | Ingress class name to use to filter Ingress and custom resources when multiple Ingress Controllers are running in the same Kubernetes cluster.|
| --kong-admin-ca-cert-file            |`string`   | none                            | Path to PEM-encoded CA certificate file to verify Kong's Admin SSL certificate.|
| --kong-admin-ca-cert                 |`string`   | none                            | PEM-encoded CA certificate string to verify Kong's Admin SSL certificate.|
| --kong-admin-concurrency             |`int`      | `10`                            | Max number of concurrent requests sent to Kong's Admin API.|
| --kong-admin-filter-tag              |`string`   | `managed-by-ingress-controller` | The tag used to manage entities in Kong.|
| --kong-admin-header                  |`string`   | none                            | Add a header (key:value) to every Admin API call, this flag can be used multiple times to specify multiple headers.|
| --kong-admin-token                   |`string`   | none                            | Set the Kong Enterprise RBAC token to be used by the controller.|
| --kong-admin-tls-server-name         |`string`   | none                            | SNI name to use to verify the certificate presented by Kong in TLS.|
| --kong-admin-tls-skip-verify         |`boolean`  | `false`                         | Disable verification of TLS certificate of Kong's Admin endpoint.|
| --kong-admin-url                     |`string`   | `http://localhost:8001`         | The address of the Kong Admin URL to connect to in the format of `protocol://address:port`.|
| --kong-url                           |`string`   | none                            | DEPRECATED, use `--kong-admin-url` |
| --kong-workspace                     |`string`   | `default`                       | Workspace in Kong Enterprise to be configured.|
| --kong-custom-entities-secret        |`string`   | none                            | Secret containing custom entities to be populated in DB-less mode, takes the form `namespace/name`.|
| --enable-reverse-sync                |`bool`     | `false`                         | Enable reverse checks from Kong to Kubernetes. Use this option only if a human has edit access to Kong's Admin API. |
| --kubeconfig                         |`string`   | none                            | Path to kubeconfig file with authorization and master location information.|
| --log_backtrace_at                   |`string`   | none                            | When set to a file and line number holding a logging statement, such as -log_backtrace_at=gopherflakes.go:234 a stack trace will be written to the Info log whenever execution hits that statement. (Unlike with -vmodule, the ".go" must be present.)|
| --log_dir                            |`string`   | none                            | If non-empty, write log files in this directory.|
| --logtostderr                        |`boolean`  | `true`                          | Logs to standard error instead of files.|
| --profiling                          |`boolean`  | `true`                          | Enable profiling via web interface `host:port/debug/pprof/`. |
| --publish-service                    |`string`   | none                            | The namespaces and name of the Kubernetes Service fronting Kong Ingress Controller in the form of namespace/name. The controller will set the status of the Ingress resouces to match the endpoints of this service. In reference deployments, this is kong/kong-proxy.|
| --publish-status-address             |`string`   | none                            | User customized address to be set in the status of ingress resources. The controller will set the endpoint records on the ingress using this address.|
| --skip-classless-ingress-v1beta1     |`boolean`  | `false`                         | Toggles whether the controller processes `extensions/v1beta1` and `networking/v1beta1` Ingress resources that have no `kubernetes.io/ingress.class` annotation.|
| --process-classless-ingress-v1       |`boolean`  | `false`                         | Toggles whether the controller processes  `networking/v1` Ingress resources that have no `kubernetes.io/ingress.class` annotation or class field.|
| --process-classless-kong-consumer    |`boolean`  | `false`                         | Toggles whether the controller processes KongConsumer resources that have no `kubernetes.io/ingress.class` annotation.|
| --stderrthreshold                    |`string`   | `2`                             | logs at or above this threshold go to stderr.|
| --sync-period                        |`duration` | `10m`                           | Relist and confirm cloud resources this often.|
| --sync-rate-limit                    |`float32`  | `0.3`                           | Define the sync frequency upper limit. |
| --update-status                      |`boolean`  | `true`                          | Indicates if the ingress controller should update the Ingress status IP/hostname.|
| --update-status-on-shutdown          |`boolean`  | `true`                          | Indicates if the ingress controller should update the Ingress status IP/hostname when the controller is being stopped.|
|  -v, --v Level                        `int`      | `0`                             | | Enable V-leveled logging at the specified level.|
| --version                            |`boolean`  | `false`                         | Shows release information about the Kong Ingress controller.|
| --vmodule moduleSpec                 |`string`   | none                            | The syntax of the argument is a comma-separated list of pattern=N, where pattern is a literal file name (minus the ".go" suffix) or "glob" pattern and N is a V level. For instance, -vmodule=gopher*=3 sets the V level to 3 in all Go files whose names begin "gopher".|
| --watch-namespace                    |`string`   | none                            | Namespace to watch for Ingress and custom resources. The default value of an empty string results in the controller watching for resources in all namespaces and configuring Kong accordingly.|
| --help                               |`boolean`  | `false`                         | Shows this documentation on the CLI and exit.|

