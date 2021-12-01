# Table of Contents

 - [2.1.0](#210)
 - [2.0.6](#206)
 - [2.0.5](#205)
 - [2.0.4](#204)
 - [2.0.3](#203)
 - [2.0.2](#202)
 - [2.0.1](#201)
 - [2.0.0](#200)
 - [1.3.3](#133)
 - [1.3.2](#132)
 - [1.3.1](#131)
 - [1.3.0](#130)
 - [1.2.0](#120)
 - [1.1.1](#111)
 - [1.1.0](#110)
 - [1.0.0](#100)
 - [0.10.0](#0100)
 - [0.9.1](#091)
 - [0.9.0](#090)
 - [0.8.1](#081)
 - [0.8.0](#080)
 - [0.7.1](#071)
 - [0.7.0](#070)
 - [0.6.2](#062)
 - [0.6.1](#061)
 - [0.6.0](#060)
 - [0.5.0](#050)
 - [0.4.0](#040)
 - [0.3.0](#030)
 - [0.2.2](#022)
 - [0.1.3](#013)
 - [0.2.1](#021)
 - [0.1.2](#012)
 - [0.1.1](#011)
 - [0.2.0](#020)
 - [0.1.0](#010)
 - [0.0.5](#005)
 - [0.0.4 and prior](#004-and-prior)

## [2.1.0]

> Release date: TBD

**Note:** the admission webhook updates originally released in [2.0.6](#206)
are _not_ applied automatically by the upgrade. If you set one up previously,
you should edit it (`kubectl edit validatingwebhookconfiguration
kong-validations`) and add `kongclusterplugins` under the `resources` block for
the `configuration.konghq.com` API group.

#### Added

- [Feature Gates][k8s-fg] have been added to the controller manager in order to
  enable alpha/beta/experimental features and provide documentation about those
  features and their maturity over time. For more information see the
  [KIC Feature Gates Documentation][kic-fg].
  [#1970](https://github.com/Kong/kubernetes-ingress-controller/pull/1970)
- a Gateway controller has been added in support of [Gateway APIs][gwapi].
  This controller is foundational and doesn't serve any end-user purpose alone.
  [#1945](https://github.com/Kong/kubernetes-ingress-controller/issues/1945)

[k8s-fg]:https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
[kic-fg]:https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md
[gwapi]:https://github.com/kubernetes-sigs/gateway-api

#### Fixed

- The validating webhook now validates that required fields data is not empty
  [#1993](https://github.com/Kong/kubernetes-ingress-controller/issues/1993)
- The validating webhook now validates unique key constraints for KongConsumer
  credentials secrets on update of secrets, and on create or update of
  KongConsumers.
  [#729](https://github.com/Kong/kubernetes-ingress-controller/issues/729)

## [2.0.6]

> Release date: 2021/11/19

**Note:** the admission webhook updates are _not_ applied automatically by the
upgrade. If you set one up previously, you should edit it (`kubectl edit
validatingwebhookconfiguration kong-validations`) and add `kongclusterplugins`
under the `resources` block for the `configuration.konghq.com` API group.

#### Fixed

- Fixed an issue where statuses would not update properly when a single service
  had multiple Ingress resources associated with it.
  [#2013](https://github.com/Kong/kubernetes-ingress-controller/pull/2013)
- Fixed an issue where statuses would not update for Ingress resources with
  periods in the name.
  [#2012](https://github.com/Kong/kubernetes-ingress-controller/issues/2012)
- The template admission webhook configuration now includes KongClusterPlugins.
  [#2000](https://github.com/Kong/kubernetes-ingress-controller/issues/2000)
- Failures to set up the status update subsystem no longer block the update
  loop.
  [#2005](https://github.com/Kong/kubernetes-ingress-controller/issues/2005)

#### Under the hood

- Updated several Go dependencies. See `go.mod` in the [diff][2.0.6] for details.

## [2.0.5]

> Release date: 2021/11/02

#### Fixed

- Fixed a bug where version reported for the controller manager was missing
  due to incorrect linker flags and missing build args in image builds.
  [#1943](https://github.com/Kong/kubernetes-ingress-controller/issues/1943)
- `hash_secret` strings in OAuth2 credentials now correctly convert to bools
  in the generated Kong configuration.
  [#1984](https://github.com/Kong/kubernetes-ingress-controller/issues/1984)
- Fixed an issue where the admission controller returned an incorrect
  status code for invalid plugin configuration.
  [#1980](https://github.com/Kong/kubernetes-ingress-controller/issues/1980)

## [2.0.4]

> Release date: 2021/10/22

#### Added

- Go Module V2 has been published so that APIs and Clients
  can be imported from external Golang projects.
  [#1936](https://github.com/Kong/kubernetes-ingress-controller/pull/1936)

#### Fixed

- Fixed a bug where the admission server's logger was missing, resulting in
  panics when the admission server tried logging.
  [#1954](https://github.com/Kong/kubernetes-ingress-controller/issues/1954)
- The admission controller now also validates KongClusterPlugin resources.
  [#1764](https://github.com/Kong/kubernetes-ingress-controller/issues/1764)
- Fixed a segfault when the version reporter failed to initialize.
  [#1961](https://github.com/Kong/kubernetes-ingress-controller/issues/1961)

## [2.0.3]

> Release date: 2021/10/19

#### Fixed

- Debug logging for resource status updates have been fixed to ensure that
  debug output isn't silently lost and to fix some formatting issues.
  [#1930](https://github.com/Kong/kubernetes-ingress-controller/pull/1930)
- Fixed a bug where Ingress resources would not be able to receive status
  updates containing relevant addresses in environments where LoadBalancer
  type services provision slowly.
  [#1931](https://github.com/Kong/kubernetes-ingress-controller/pull/1931)

## [2.0.2]

> Release date: 2021/10/14

#### Added

- Builds now produce Red Hat UBI-based images.

## [2.0.1]

> Release date: 2021/10/11
#### Added

- The ingress controller version now gets logged on startup.
  [#1911](https://github.com/Kong/kubernetes-ingress-controller/pull/1911)

#### Fixed

- Fixed an issue reading workspace information with RBAC permissions that
  only allow access to the specified workspace.
  [#1900](https://github.com/Kong/kubernetes-ingress-controller/issues/1900)

## [2.0.0]

> Release date: 2021/10/07

**NOTE**: This changelog entry was compiled from every changelog entry in the
  `alpha` and `beta` pre-releases of `2.0.0`. If you're looking for the interim
  changelog between `alpha` and/or `beta` versions prior to the release see
  the [historical changelog here][alpha-beta-changelog].

[alpha-beta-changelog]:https://github.com/Kong/kubernetes-ingress-controller/blob/3e9761c378d02eda1c4622d87b899ea9ea4c35b4/CHANGELOG.md

#### Breaking changes

While you're reviewing the breaking changes below we also recommend you check
out our [upgrade guide][upgrade-1-3-to-2-0] which covers upgrading from the
previous `v1.3.x` releases to this release.

- The admission webhook now requires clients that support TLS 1.2 or higher.
  [#1671](https://github.com/Kong/kubernetes-ingress-controller/issues/1671)
- autonegotiation of the Ingress API version (extensions v1beta1, networking
  v1beta1, networking v1) has been disabled. Instead, the user is expected to
  set **exactly** one of:
  `--controller-ingress-networkingv1`
  `--controller-ingress-networkingv1beta1`
  `--controller-ingress-extensionsv1beta1`
- several miscellaneous flags have been removed.
  The following flags are no longer present:
  - `--disable-ingress-extensionsv1beta1` (replaced by `--enable-controller-ingress-extensionsv1beta1=false`)
  - `--disable-ingress-networkingv1` (replaced by `--enable-controller-ingress-networkingv1=false`)
  - `--disable-ingress-networkingv1beta1` (replaced by `--enable-controller-ingress-networkingv1beta1=false`)
  - `--version`
  - `--alsologtostderr`
  - `--logtostderr`
  - `--v`
  - `--vmodule`
- support for "classless" ingress types has been removed.
  The following flags are no longer present:
  - `--process-classless-ingress-v1beta1`
  - `--process-classless-ingress-v1`
  - `--process-classless-kong-consumer`
- `--dump-config` (a diagnostic option) is now a boolean. `true` is equivalent
  to the old `enabled` value. `false` is equivalent to the old `disabled`
  value. `true` with the additional new `--dump-sensitive-config=true` flag is
  equivalent to the old `sensitive` value.
- The historical `--stderrthreshold` flag is now deprecated: it no longer has
  any effect when used and will be removed in a later release.
  [#1297](https://github.com/Kong/kubernetes-ingress-controller/issues/1297)
- The `--update-status-on-shutdown` flag which supplements the `--update-status`
  flag has been deprecated and will no longer have any effect, it will be removed
  in a later release.
  [#1304](https://github.com/Kong/kubernetes-ingress-controller/issues/1304)
- the `--sync-rate-limit` is now deprecated in favor of `--sync-time-seconds`.
  This functionality no longer blocks goroutines until the provided number of
  seconds has passed to enforce rate limiting, now instead it configures a
  non-blocking [time.Ticker][go-tick] that runs at the provided seconds
  interval. Input remains a float that indicates seconds.

[upgrade-1-3-to-2-0]:https://docs.konghq.com/kubernetes-ingress-controller/2.0.x/guides/upgrade/

#### Added

- Individual controllers can now be enabled or disabled at a granular level.
  For example you can disable the controller for `TCPIngress` with:
  `--enable-controller-tcpingress=false`
  To see the entire list of configurable controllers run the controller manager
  with `--help`.
  [#1638](https://github.com/Kong/kubernetes-ingress-controller/issues/1638)
- The `--watch-namespace` flag was added and supports watching a single
  specific namespace (e.g. `--watch-namespace namespaceA`) or multiple
  distinct namespaces using a comma-separated list (e.g.
  `--watch-namespace "namespaceA,namespaceB"`). If not provided the default
  behavior is to watch **all namespaces** as it was in previous releases.
  [#1317](https://github.com/Kong/kubernetes-ingress-controller/pull/1317)
- UDP support was added via the `v1beta1.UDPIngress` API.
  [#1454](https://github.com/Kong/kubernetes-ingress-controller/pull/1454)
  [UDP Blog Post][kong-udp]
- Renamed roles and bindings to reflect their association with Kong.
  [#1801](https://github.com/Kong/kubernetes-ingress-controller/issues/1801)
- Upgraded Kong Gateway from 2.4 to 2.5
  [#1684](https://github.com/Kong/kubernetes-ingress-controller/issues/1684)
- Decreased log level of some status update messages.
  [#1641](https://github.com/Kong/kubernetes-ingress-controller/issues/1641)
- Added metrics tracking whether configuration was successfully generated and
  applied and the time taken to sync configuration to Kong.
  [#1622](https://github.com/Kong/kubernetes-ingress-controller/issues/1622)
- Added a [Prometheus operator PodMonitor](https://github.com/Kong/kubernetes-ingress-controller/blob/v2.0.0-beta.1/config/prometheus/monitor.yaml)
  to scrape controller and Kong metrics. To use it:
  ```
  kubectl apply -f https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/main/config/prometheus/monitor.yaml
  ```
  [#1657](https://github.com/Kong/kubernetes-ingress-controller/issues/1657)
- Added controller functional metrics in 2.x
  [#705](https://github.com/Kong/kubernetes-ingress-controller/issues/705)
- Implemented Ingress status updates in 2.x.
  [#1451](https://github.com/Kong/kubernetes-ingress-controller/pull/1451)
- Added `--publish-status-address` and `--publish-service` flags to 2.x.
  [#1451](https://github.com/Kong/kubernetes-ingress-controller/pull/1451)
  [#1509](https://github.com/Kong/kubernetes-ingress-controller/pull/1509)
- Added scripts to generate 2.x manifests.
  [#1563](https://github.com/Kong/kubernetes-ingress-controller/pull/1563)
- Added support for --dump-config to 2.x.
  [#1589](https://github.com/Kong/kubernetes-ingress-controller/pull/1589)
- profiling using `pprof` is now a standalone HTTP server listening on port 10256.
- adds support for selector tags (filter) tags refractored work. KIC 1.x
  [#1415](https://github.com/Kong/kubernetes-ingress-controller/pull/1415)
- Profiling using `pprof` is now a standalone HTTP server listening on port 10256.
  [#1417](https://github.com/Kong/kubernetes-ingress-controller/pull/1417)
- Reduced 2.x RBAC permissions to match 1.x permissions and added a generated
  single-namespace Role that matches the ClusterRole.
  [#1457](https://github.com/Kong/kubernetes-ingress-controller/pull/1457)
- support for the `konghq.com/host-aliases` annotation.
  [#1016](https://github.com/Kong/kubernetes-ingress-controller/pull/1016/)
- Added `--proxy-timeout-seconds` flag to configure the kong client api timeout.
  [#1401](https://github.com/Kong/kubernetes-ingress-controller/pull/1401)

[kong-udp]:https://konghq.com/blog/kong-gateway-2-2-released/#UDP-Support

#### Fixed

- In DB-less mode, the controller only marks itself ready once it has
  successfully applied configuration at least once. This ensures that proxies
  do not start handling traffic until they are configured.
  [#1720](https://github.com/Kong/kubernetes-ingress-controller/issues/1720)
- Prometheus metrics were not exposed on the metrics endpoint in 2.0.0-beta.1 by default
  [#1497](https://github.com/Kong/kubernetes-ingress-controller/issues/1497)
- Resolved an issue where certain UDPIngress and TCPIngress configurations
  resulted in overlapping incompatible Kong configuration.
  [#1702](https://github.com/Kong/kubernetes-ingress-controller/issues/1702)
- Fixed a panic that would occur in the controller manager when a
  `KongConsumer` object with an empty name was submitted.
  Any `KongConsumer` resource created with an empty `UserName` will
  now throw an error in the controller manager logs (this wont stop
  other configurations from proceeding), but the object in question
  will thereafter otherwise be skipped for backend configuration
  until the resource has been corrected.
  [#1658](https://github.com/Kong/kubernetes-ingress-controller/issues/1658)
- The controller will now retry unsuccessful TCPIngress status updates.
  [#1641](https://github.com/Kong/kubernetes-ingress-controller/issues/1641)
- The controller now correctly disables Knative controllers automatically when
  Knative controllers are not installed.
  [#1585](https://github.com/Kong/kubernetes-ingress-controller/issues/1585)
- Corrected the old Ingress v1beta1 API group.
  [#1584](https://github.com/Kong/kubernetes-ingress-controller/pull/1584)
- Updated our Knative API support for more recent upstream releases.
  [#1148] (https://github.com/Kong/kubernetes-ingress-controller/pull/1396)

#### Under the hood

- Updated the compiler to [Go v1.17](https://golang.org/doc/go1.17)
  [#1714](https://github.com/Kong/kubernetes-ingress-controller/issues/1714)
- Code for the previous v1.x releases of the Kubernetes Ingress Controller
  have been removed. Maintenance of the v1.x era codebase lives on in the
  `1.3.x` and related branches going forward.
  [#1591](https://github.com/Kong/kubernetes-ingress-controller/issues/1591)
- Made assorted improvements to CI and test code.
  [#1646](https://github.com/Kong/kubernetes-ingress-controller/issues/1646)
  [#1664](https://github.com/Kong/kubernetes-ingress-controller/issues/1664)
  [#1669](https://github.com/Kong/kubernetes-ingress-controller/issues/1669)
  [#1672](https://github.com/Kong/kubernetes-ingress-controller/issues/1672)
- New `v1` versions of `CustomResourceDefinitions` introduced for KIC 2.0 are now
  backwards compatible with the previous `v1beta1` CRD definitions (i.e. `v1beta1 -> v1`
  upgrades of KIC's CustomResourceDefinitions now work fully automatically). In practice
  the upgrade process should be seamless for end-users (e.g. `kubectl apply -f <NEW CRDS>`).
  If you're interested in better understanding the differences and what's going on
  under the hood, please see the relevant PR which includes the user facing changes.
  [Kubernetes#79604](https://github.com/kubernetes/kubernetes/pull/79604)
  [#1133](https://github.com/Kong/kubernetes-ingress-controller/issues/1133)
- The uuid generation is now done by the same library in the whole project
  [#1604](https://github.com/Kong/kubernetes-ingress-controller/issues/1604)
- the controller manager will no longer log multiple entries for `nil` updates
  to the Kong Admin API. The result is that operators will no longer see multiple
  "no configuration change, skipping sync to kong" entries for any single update,
  instead it will only report this `nil` update scenario the first time it is
  encountered for any particular SHA derived from the configuration contents.
- project layout for contributions has been changed: this project now uses the
  [Kubebuilder SDK][kubebuilder] and there are layout changes and
  configurations specific to the new build environment.
- controller architecture has been changed: each API type now has an
  independent controller implementation and all controllers now utilize
  [controller-runtime][controller-runtime].
- full integration testing in [Golang][go] has been added for testing APIs and
  controllers on a fully featured Kubernetes cluster, this is now supported by
  the new [Kong Kubernetes Testing Framework (KTF)][ktf] project and now runs
  as part of CI.
- the mechanism for caching and resolving Kong Admin `/config` configurations
  when running in `DBLESS` mode has been reimplemented to enable fine-tuned
  configuration options in later iterations.
- contains the refactored admission webhook server. The server key and
  certificate flags have improved semantics: the default flag value is no
  longer the default path, but an empty string. When both key/cert value flags
  and key/cert file flags remain unset, KIC will read cert/key files from the
  default paths, as said in the flag descriptions. This change should not
  affect any existing configuration - in all configuration cases, behavior is
  expected to remain unchanged.
- taking configuration values from environment variables no longer uses Viper.

[go-tick]:https://golang.org/pkg/time/#Ticker
[kubebuilder]:https://github.com/kubernetes-sigs/kubebuilder
[controller-runtime]:https://github.com/kubernetes-sigs/controller-runtime
[go]:https://golang.org
[ktf]:https://github.com/kong/kubernetes-testing-framework

## [1.3.3]

> Release date: 2021/10/01

#### Fixed

- Fixed invalid plugin validation code in admission controller.
  [go-kong#81](https://github.com/Kong/go-kong/pull/81)
- Fixed a panic when sorting consumers.
  [#1658](https://github.com/Kong/kubernetes-ingress-controller/pull/1658)

## [1.3.2]

> Release date: 2021/08/12

#### Under the hood

- Updated Alpine image to 3.14.
  [#1691](https://github.com/Kong/kubernetes-ingress-controller/pull/1691/)
- Update Kong images to 2.5.

## [1.3.1]

> Release date: 2021/06/03

#### Fixed

- fixed a bug that now stops `v1.3.x` releases from advertising themselves as `v2` if manually built with default configurations.

#### Under the hood

- Upgraded CI dependencies
- Some cleanup iterations on RELEASE.md release process

## [1.3.0]

> Release date: 2021/05/27

#### Added

- support for the `konghq.com/host-aliases` annotation.
  [#1016](https://github.com/Kong/kubernetes-ingress-controller/pull/1016/)

#### Fixed

- Sort SNIs and certificates consistently to avoid an issue with unnecessary
  configuration re-syncs.
  [#1268](https://github.com/Kong/kubernetes-ingress-controller/pull/1268/)

#### Under the hood

- Upgraded various dependencies.

## [1.2.0]

> Release date: 2021/03/24

#### Added

- Ingresses now support `konghq.com/request-buffering` and
  `konghq.com/response-buffering` annotations, which set the
  `request-buffering` and `response-buffering` settings on associated Kong
  routes.
  [#1016](https://github.com/Kong/kubernetes-ingress-controller/pull/1016/)
- Added `--dump-config` flag to dump generated Kong configuration to a
  temporary file to debug issues where the controller generates unexpected
  configuration or unacceptable configuration. When set to `enabled` it redacts
  sensitive values (credentials and certificate keys), and when set to
  `sensitive`, it includes all configuration.
  [#991](https://github.com/Kong/kubernetes-ingress-controller/pull/991/)
- Added support for mtls-auth plugin credentials (requires Enterprise 2.3.2.0
  or newer).
  [#1078](https://github.com/Kong/kubernetes-ingress-controller/pull/1078/)
- The KongClusterPlugin CRD is now optional, for installation in clusters
  where KIC administrators do not have cluster-wide permissions.

#### Fixed

- The admission webhook can now validate KongPlugin configurations stored in a
  Secret.
  [#1036](https://github.com/Kong/kubernetes-ingress-controller/pull/1036/)

#### Under the hood

- Build configuration allows target architectures other than `amd64`. Note that
  other architectures are not officially supported.
  [#1046](https://github.com/Kong/kubernetes-ingress-controller/pull/1046/)
- Updated to Go 1.16. Make sure to update your Go version if you build your own
  controller binaries.
  [#1110](https://github.com/Kong/kubernetes-ingress-controller/pull/1110/)
- Refactored synchronization loop into more discrete components and created
  packages for them.
  [#1027](https://github.com/Kong/kubernetes-ingress-controller/pull/1027/)
  [#1029](https://github.com/Kong/kubernetes-ingress-controller/pull/1029/)
- Broad refactoring (with the purpose of exposing KIC's logic as libraries), in preparation for an architectural upgrade of KIC to a [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)-based
  implementation of the controller (expected to be released as KIC v2.0).
  [#1037](https://github.com/Kong/kubernetes-ingress-controller/pull/1037/)
- Added a Go-based integration test environment and initial set of tests.
  [#1102](https://github.com/Kong/kubernetes-ingress-controller/pull/1102/)
- CI improvements check test coverage for PRs and automaticall open PRs for
  dependency updates.
- Upgraded almost all Go library dependencies (from now on, using Dependabot to ensure that minor releases use the newest versions available).

## [1.1.1]

> Release date: 2021/01/07

#### Fixed

- Ingress controller now correctly sets ports for ExternalName services [#985](https://github.com/Kong/kubernetes-ingress-controller/pull/985)
- TCPIngress CRD: removed the duplicated subresource YAML key [#997](https://github.com/Kong/kubernetes-ingress-controller/pull/997)

#### Deprecated

- Removed Helm 2 installation instructions because Helm 2 is EOL. Use Helm 3 instead. [#993](https://github.com/Kong/kubernetes-ingress-controller/pull/993)

## [1.1.0]

> Release date: 2020/12/09

#### Breaking changes

- The controller no longer supports Cassandra-backed Kong clusters, following
  deprecation in 0.9.0. You must migrate to a Postgres-backed or DB-less
  cluster before upgrading to 1.1.0. The controller will restore configuration
  from Kubernetes resources (Ingresses, Services, KongPlugins, etc.) into the
  new datastore automatically. Kong Enterprise users with
  non-controller-managed configuration (Portal configuration, RBAC
  configuration, etc.) will need to migrate that configuration manually.
  [#974](https://github.com/Kong/kubernetes-ingress-controller/pull/974)

#### Added

- The default Kong version is now 2.2.x and the default Kong Enterprise version
  is now 2.2.0.0.
  [#932](https://github.com/Kong/kubernetes-ingress-controller/pull/932)
  [#965](https://github.com/Kong/kubernetes-ingress-controller/pull/965)
- The default worker count is now 2 instead of 1. This avoids request latency
  during blocking configuration changes.
  [#957](https://github.com/Kong/kubernetes-ingress-controller/pull/957)
- Knative Services now support `konghq.com/override` (for attaching
  KongIngress resources).
  [#908](https://github.com/Kong/kubernetes-ingress-controller/pull/908)
- Added the `konghq.com/snis` Ingress annotation. This populates SNI
  configuration on the routes derived from the annotated Ingress.
  [#863](https://github.com/Kong/kubernetes-ingress-controller/pull/863)

#### Fixed

- The controller now correctly prints the affected Service name when logging
  warnings about Services without any endpoints.
  [#915](https://github.com/Kong/kubernetes-ingress-controller/pull/915)
- Credentials that lack critical fields no longer result in a panic.
  [#944](https://github.com/Kong/kubernetes-ingress-controller/pull/944)

## [1.0.0]

> Release date: 2020/10/05

#### Breaking changes

- The controller no longer supports versions of Kong prior to 2.0.0.
  [#875](https://github.com/Kong/kubernetes-ingress-controller/pull/875)
- Deprecated 0.x.x flags are no longer supported. Please see [the documentation
  changes](https://github.com/Kong/kubernetes-ingress-controller/pull/866/files#diff-9a686fb3bf9c18ab81952a0933fb5c00)
  for a complete list of removed flags and their replacements. Note that this
  change applies to both flags and their equivalent environment variables, e.g.
  for `--admin-header`, if you set `CONTROLLER_ADMIN_HEADER`, you should now
  use `CONTROLLER_KONG_ADMIN_HEADER`.
  [#866](https://github.com/Kong/kubernetes-ingress-controller/pull/866)
- KongCredential custom resources are no longer supported. You should convert
  any KongCredential resources to [credential Secrets](https://docs.konghq.com/kubernetes-ingress-controller/1.0.x/guides/using-consumer-credential-resource/#provision-a-consumer)
  before upgrading to 1.0.0.
  [#862](https://github.com/Kong/kubernetes-ingress-controller/pull/862)
- Deprecated 0.x.x annotations are no longer supported. Please see [the
  documentation changes](https://github.com/Kong/kubernetes-ingress-controller/pull/873/files#diff-777e08783d63482620961c4f93a1e1f6)
  for a complete list of removed annotations and their replacements.
  [#873](https://github.com/Kong/kubernetes-ingress-controller/pull/873)

#### Added

- The controller Docker registry now has minor version tags. These always point
  to the latest patch release for a given minor version, e.g. if `1.0.3` is the
  latest patch release for the `1.0.x` series, the `1.0` Docker tag will point
  to `1.0.3`.
  [#747](https://github.com/Kong/kubernetes-ingress-controller/pull/747)
- Custom resources now all have a status field. For 1.0.0, this field is a
  placeholder, and does not contain any actual status information. Future
  versions will add status information that reflects whether the controller has
  created Kong configuration for that custom resource.
  [#824](https://github.com/Kong/kubernetes-ingress-controller/pull/824)
- Version compatibility documentation now includes [information about supported
  Kubernetes versions for a given controller version](https://docs.konghq.com/kubernetes-ingress-controller/1.0.x/references/version-compatibility/#kubernetes).
  [#820](https://github.com/Kong/kubernetes-ingress-controller/pull/820)

#### Fixed

- EKS documentation now uses hostnames rather than IP addresses.
  [#877](https://github.com/Kong/kubernetes-ingress-controller/pull/877)

## [0.10.0]

> Release date: 2020/09/15

#### Breaking changes

- Ingress resources now require `kubernetes.io/ingress.class` annotations by
  default. Kong recommends adding this annotation to Ingresses that previously
  did not have it, but you can override this change and instruct the controller
  to process Ingresses without this annotation if desired. See the [ingress
  class documentation](https://example.com/link-tbd) for details.
  [#767](https://github.com/Kong/kubernetes-ingress-controller/pull/767)
- KongConsumer resources now require `kubernetes.io/ingress.class` annotations
  by default. This change can also be overriden using a flag.
  [#767](https://github.com/Kong/kubernetes-ingress-controller/pull/767)
- TCPIngress resources now require `kubernetes.io/ingress.class` annotations.
  This change _cannot_ be overriden.
  [#767](https://github.com/Kong/kubernetes-ingress-controller/pull/767)
- CA certificate secrets now require `kubernetes.io/ingress.class` annotations.
  This change _cannot_ be overriden.
  [#815](https://github.com/Kong/kubernetes-ingress-controller/pull/815)
- Removed support for global KongPlugin resources. You must now use
  KongClusterPlugin resources for global plugins. You should run
  `kubectl get kongplugin -l global=true --all-namespaces` to list existing
  global KongPlugins to find and convert them before upgrading. The controller
  will also log a warning if it finds any global KongPlugins that are still in
  place.
  [#751](https://github.com/Kong/kubernetes-ingress-controller/pull/751)


#### Added

- Added support for [Ingress
  v1](https://github.com/kubernetes/enhancements/tree/master/keps/sig-network/1453-ingress-api#summary-of-the-proposed-changes).
  [#832](https://github.com/Kong/kubernetes-ingress-controller/pull/832).
  [#843](https://github.com/Kong/kubernetes-ingress-controller/pull/843).
- Added support for the port mapping functionality in Kong versions 2.1 and
  newer in example manifests. This feature [improves Kong's functionality when
  behind a load balancer that uses different ports than Kong's proxy
  listens](https://github.com/Kong/kong/pull/5861).
  [#753](https://github.com/Kong/kubernetes-ingress-controller/pull/753)
- Added support for the `ingress.kubernetes.io/force-ssl-redirect` annotation.
  [#745](https://github.com/Kong/kubernetes-ingress-controller/pull/745)
- Transitioned to structured logging.
  [#748](https://github.com/Kong/kubernetes-ingress-controller/pull/748)
- Added flags to enable processing of Ingress and KongConsumer resources
  without `ingress.class` annotations regardless of the controller class.
  Previously, this functionality was only available when using the default
  controller class, and could not be disabled.
  [#767](https://github.com/Kong/kubernetes-ingress-controller/pull/767)
- Added support for `admission.k8s.io/v1` validating webhooks.
  [#759](https://github.com/Kong/kubernetes-ingress-controller/pull/759)
- Migrated to Go 1.13-style error handling.
  [#765](https://github.com/Kong/kubernetes-ingress-controller/pull/765)
- Added documentation for using the controller along with Istio.
  [#798](https://github.com/Kong/kubernetes-ingress-controller/pull/798)
- Updated documentation to include information on Kong 2.1.

#### Fixed

- Removed `securityContext` from example deployments. Earlier Kong versions
  had to run as root to support some Enterprise features. This is no longer the
  case in modern Kong versions.
  [#672](https://github.com/Kong/kubernetes-ingress-controller/pull/672)
- Added missing documentation for `--enable-reverse-sync` flag.
  [#718](https://github.com/Kong/kubernetes-ingress-controller/pull/718)
- Fixed a bug where the controller did not track updates to resources that
  should not have required `ingress.class` unless that annotation was present.
  [#767](https://github.com/Kong/kubernetes-ingress-controller/pull/767)
- Clarified build instructions for pushing Docker artifacts.
  [#768](https://github.com/Kong/kubernetes-ingress-controller/pull/768)
- Improved controller startup behavior in scenarios where Kong was not
  available. The controller will now retry and exit with an error after a
  timeout, rather than hanging indefinitely.
  [#771](https://github.com/Kong/kubernetes-ingress-controller/pull/771)
  [#799](https://github.com/Kong/kubernetes-ingress-controller/pull/799)
- Addressed several documentation typos and incongruent examples.
  [#776](https://github.com/Kong/kubernetes-ingress-controller/pull/776)
  [#785](https://github.com/Kong/kubernetes-ingress-controller/pull/785)
  [#809](https://github.com/Kong/kubernetes-ingress-controller/pull/809)
- Corrected a Helm 3 example that still used deprecated Helm 2 flags.
  [#793](https://github.com/Kong/kubernetes-ingress-controller/pull/793)

#### Under the hood

- Improved tests by removing many hard-coded default values. The tests now
  reference variables that define the default value in a single location.
  [#815](https://github.com/Kong/kubernetes-ingress-controller/pull/815)
- Added CI warning when base and single-file example manifests diverge.
  [#797](https://github.com/Kong/kubernetes-ingress-controller/pull/797)
- Updated Kubernetes dependencies from v0.17.x to v0.19.0 and switched from
  `knative.dev/serving` to `knative.dev/networking`.
  [#813](https://github.com/Kong/kubernetes-ingress-controller/pull/813)
  [#817](https://github.com/Kong/kubernetes-ingress-controller/pull/817)
- Updated Go build configuration to use Go 1.15.
  [#816](https://github.com/Kong/kubernetes-ingress-controller/pull/816)

## [0.9.1]

> Release date: 2020/06/08

#### Fixed

- Parse TLS section of Knative Ingress resources
  [#721](https://github.com/Kong/kubernetes-ingress-controller/pull/721)

## [0.9.0]

> Release date: 2020/05/26

#### Breaking change

Health-check behavior of the default manifest has been changed to use
`status` interface of Kong instead of a simple Nginx server block.
The change is transparent and doesn't require any additional work.
[#634](https://github.com/Kong/kubernetes-ingress-controller/pull/634)

### Deprecations

Kong deployments backed by Cassandra are deprecated and will not be supported
in future. Cassandra deployments for Ingress Controller use cases are rare
and seldom make sense since the features that Cassandra brings are
provided by other means in such architectures.
[#617](https://github.com/Kong/kubernetes-ingress-controller/pull/617)

#### Added

- **Plugin configuration via Kubernetes Secrets**  Configuration of plugins
  can be stored in Kubernetes Secrets and then referenced in `KongPlugin`
  and `KongClusterPlugin` resources.
  [#618](https://github.com/Kong/kubernetes-ingress-controller/pull/618)
- **mTLS authentication**  The controller can configure CA Certificates
  in Kong and these can be used by `mtls-auth` plugin in Kong. The plugin
  is currently enterprise-only.
  [#616](https://github.com/Kong/kubernetes-ingress-controller/pull/616)
- **Kong Custom entities in DB-less mode** Custom entities used in
  custom plugins can now be configured for DB-less deployments of Kong.
  [#630](https://github.com/Kong/kubernetes-ingress-controller/pull/630)
- **Host-header manipulation**  Host header of a request destined to a
  Kubernetes Service can now be manipulated using the `konghq.com/host-header`
  annotation on the `Service` resource.
  [#597](https://github.com/Kong/kubernetes-ingress-controller/pull/597)
- **Method-based routing**  Method based routing can be performed using the
  Ingress resource. A new annotation `konghq.com/methods` can now be used to
  match HTTP method in addition to HTTP `host` and `path`. This was
  previously supported only via `KongIngress` Custom Resource.
  [#591](https://github.com/Kong/kubernetes-ingress-controller/pull/591)
- **New configuration options** Following new CLI flags and corresponding
  environment variables have been added:
  - `--admission-webhook-cert`, `--admission-webhook-key`
    and `--kong-admin-ca-cert`. These have been added to ease configuration
    by enabling users to supply sensitive values using `Secret`
    references inside `PodSpec`.
    [#628](https://github.com/Kong/kubernetes-ingress-controller/pull/628)
  - `--kong-custom-entities-secret` flag has been added to support
    custom entities in DB-less mode feature.

#### Fixed

- Some errors that were previously ignored are being caught and handled
  correctly
  [#635](https://github.com/Kong/kubernetes-ingress-controller/pull/635)
- Ingress rules with consecutive slashes (`//`) are now ignored
  [#663](https://github.com/Kong/kubernetes-ingress-controller/pull/663)

## [0.8.1]

> Release date: 2020/04/15

#### Added

- Added `--enable-reverse-sync` flag to enable checks from Kong to kubernetes
  state. This should be enabled only if a human has access to Kong's Admin API
  or Kong Manager (for Enterprise users). This flag will disable an optimization
  in the controller and result in an increase read activity on Kong's Admin
  API and database.
  [#559](https://github.com/Kong/kubernetes-ingress-controller/issues/559)

#### Fixed

- Fix certificate and SNI sync to avoid a deadlock due to a conflict when
  Kong is running with a database backend.
  [#524](https://github.com/Kong/kubernetes-ingress-controller/issues/524)
- Correctly set Knative Ingress Status
  [#600](https://github.com/Kong/kubernetes-ingress-controller/pull/600)

## [0.8.0]

> Release date: 2020/03/25

#### Breaking changes

- **`strip_path` disabled by default**
  The value of `strip_path` of routes in Kong is now set to `false`.
  If you are upgrading from a previous version, please carefully test the change
  before rolling it out as this change can possibly break the routing
  for your clusters.
  You can use `konghq.com/strip-path` annotation to set the value to `true`.

#### Deprecations

The following annotations are now deprecated and will be removed in a future
release:
- `configuration.konghq.com`
- `plugins.konghq.com`
- `configuration.konghq.com/protocols`
- `configuration.konghq.com/protocol`
- `configuration.konghq.com/client-cert`

Please read the annotations document for new annotations.

#### Added

- **Knative Ingress support**  The controller can now proxy traffic for
  serverless workloads running on top of Knative. Users can also select
  Kong plugins to execute on a per Knative workload/service basis.
  [#563](https://github.com/Kong/kubernetes-ingress-controller/pull/563)
- **TCP/TLS routing**  New Custom Resource TCPIngress has been introduced
  to support TCP proxy. SNI-based proxying is also supported for TLS encrypted
  TCP streams.
  [#527](https://github.com/Kong/kubernetes-ingress-controller/pull/527)
- **New Custom Resource KongClusterPlugin**  Plugin configuration can now
  be shared acrossed Kubernetes namespaces using `KongClusterPlugin`, a new
  cluster-level Custom Resource.
  [#520](https://github.com/Kong/kubernetes-ingress-controller/pull/520)
- **New annotation group `konghq.com`** A new annotations group has
  been introduced which should simplify configuration and reduce the need of
  `KongIngress` resource in most use-cases. The following new annotations
  have been introduced:
  - `konghq.com/plugins`
  - `konghq.com/override`
  - `konghq.com/client-cert`
  - `konghq.com/protocols`
  - `konghq.com/protocol`
  - `konghq.com/preserve-host`
  - `konghq.com/plugins`
  - `konghq.com/override`
  - `konghq.com/path`
  - `konghq.com/strip-path`
  - `konghq.com/https-redirect-status-code`

#### Fixed

- Admission webhook now checks for the correct fields for JWT credential
  type.
  [#556](https://github.com/Kong/kubernetes-ingress-controller/pull/556)

#### Under the hood

- decK has been upgraded to v1.0.3.
  [#576](https://github.com/Kong/kubernetes-ingress-controller/pull/576)
- Go has been upgraded to 1.14.
  [#579](https://github.com/Kong/kubernetes-ingress-controller/pull/579)
- Alpine docker image has been upgraded to 3.11.
  [#567](https://github.com/Kong/kubernetes-ingress-controller/pull/567)

## [0.7.1]

> Release date: 2020/01/31

#### Summary

This releases contains bug-fixes only. All users are advised to upgrade.

#### Fixed

- De-duplicate SNIs when the same SNI is associated with multiple secrets.
  [#510](https://github.com/Kong/kubernetes-ingress-controller/issues/510)
- `plugin.RunOn` is not injected when Kong version >= 2.0.0.
  [#521](https://github.com/Kong/kubernetes-ingress-controller/issues/521)
- Parse default backend in `Ingress` resource correctly.
  [#511](https://github.com/Kong/kubernetes-ingress-controller/issues/511)
- KongPlugin resources with `global: true` label are correctly processed
  to include `protocols` fields while rendering Kong's configuration.
  [#502](https://github.com/Kong/kubernetes-ingress-controller/issues/502)
- Admission Controller: correctly process updates to `KongConsumer` resource
  [#501](https://github.com/Kong/kubernetes-ingress-controller/issues/501)
- Do not send multiple update events for a single CRD update
  [#514](https://github.com/Kong/kubernetes-ingress-controller/issues/514)

## [0.7.0]

> Release date: 2020/01/06

#### Summary

This release adds secret-based credentials, gRPC routing, upstream mutual
authentication, DB-less deployment by default and performance improvements.

#### Breaking changes

- The default value of `--admission-webhook-listen` flag is now `off` to avoid
  an error in the logs when the cert and key pair is not provided. Users will
  have to explicitly set this flag to `:8080` to enable it. Please do note that
  it is recommended to always set up the Admission Controller.

#### Added

- **Multi-port services** Ingress rules forwarding traffic to multiple ports
  of the same services are now supported. The names of the services configured
  in Kong have been changed to include the port number/name for uniqueness.
  [#404](https://github.com/Kong/kubernetes-ingress-controller/pull/404)
- When using the controller with Kong Enterprise,
  Controller now attempts to create the workspace configured via
  `--kong-workspace`, if it does not exist.
  [#429](https://github.com/Kong/kubernetes-ingress-controller/pull/429)
- **Controller configuration revamped** Configuration of the controller itself
  can now be tweaked via environment flags and CLI flags, both. Environment
  variables and Secrets can be used to pass sensitive information to the
  controller.
  [#436](https://github.com/Kong/kubernetes-ingress-controller/pull/436)
- **Encrypted credentials via Secrets** Credentials can now be configured via
  `Secret` resource from the Kubernetes core API. These credentials are
  encrypted at rest by Kubernetes. The controller loads these secrets into
  Kong's memory or database from the Kubernetes data-store.
  [#430](https://github.com/Kong/kubernetes-ingress-controller/pull/430)
- **Multi-entity plugins** Plugins can now be configured for a combination of
  an Ingress rule(s) and KongConsumer or a combination of a Service
  and KongConsumer.
  [#386](https://github.com/Kong/kubernetes-ingress-controller/issues/386)
- **Mutual authentication using mTLS** Kong and the Kubernetes Service can
  mutually authenticate each other now. Use the new
  `configuration.konghq.com/client-cert` annotation on a Kubernetes Service
  to specify the cert-key pair Kong should use to authenticate itself.
  [#483](https://github.com/Kong/kubernetes-ingress-controller/pull/483)
- **gRPC routing** Kong Ingress Controller can now expose and proxy gRPC
  protocol based services, in addition to HTTP-based services. These can
  be configured using the core Ingress resource itself.
  [#454](https://github.com/Kong/kubernetes-ingress-controller/pull/454)
- **Performance improvement** Number of sync calls to Kong, in both DB and
  DB-less mode, should be reduced by an order of magnitude for most deployments.
  This will also improve Kong's performance.
  [#484](https://github.com/Kong/kubernetes-ingress-controller/pull/484)
- `credentials` property has been added to the `KongConsumer` Custom Resource.
  This property holds the references to the secrets containing the credentials.
  [#430](https://github.com/Kong/kubernetes-ingress-controller/pull/430)
- Flag `--kong-admin-filter-tag` has been added to change the tag used
  to filter and managed entity in Kong's database. This defaults to
  `managed-by-ingress-controller`.
  [#440](https://github.com/Kong/kubernetes-ingress-controller/pull/440)
- Flag `--kong-admin-concurrency` has been added to control the number of
  concurrent requests between the controller and Kong's Admin API.
  This defaults to `10`.
  [#481](https://github.com/Kong/kubernetes-ingress-controller/pull/481)
- Flag `--kong-admin-token` has been added to supply the RBAC token
  for the Admin API for Kong Enterprise deployments.
  [#489](https://github.com/Kong/kubernetes-ingress-controller/pull/489)
- Admission Controller now validates Secret-based credentials. It ensures that
  the required fields are set in the secret and the credential type is a
  valid one.
  [#446](https://github.com/Kong/kubernetes-ingress-controller/pull/446)
- `http2` is now enabled by default on the TLS port.
  [#456](https://github.com/Kong/kubernetes-ingress-controller/pull/456)
- DB-less or the in-memory mode is now the new default in the reference
  manifests. It is recommended to run Kong without a database for Ingress
  Controller deployments.
  [#456](https://github.com/Kong/kubernetes-ingress-controller/pull/456)
- `upstream.host_header` property has been added to the `KongIngress` Custom
  Resource. This property can be used to change the `host` header in every
  request that is sent to the upstream service.
  [#478](https://github.com/Kong/kubernetes-ingress-controller/pull/478)

#### Fixed

- Every event in the queue is not logged anymore as it can leak sensitive
  information in the logs. Thanks to [@goober](https://github.com/goober)
  for the report.
  [#439](https://github.com/Kong/kubernetes-ingress-controller/pull/439)
- For database deployments, `upstream` entity are now created with `round-robin`
  as default `algorithm` to avoid false positives during a sync operation.
  These false positives can have a negative impact on Kong's performance.
  [#480](https://github.com/Kong/kubernetes-ingress-controller/pull/480)


#### Deprecated

- `KongCredential` Custom Resource is now deprecated and will be remove in a
  future release. Instead, please use Secret-based credentials.
  [#430](https://github.com/Kong/kubernetes-ingress-controller/pull/430):
- Following flags have been deprecated and new ones have been added in place
  [#436](https://github.com/Kong/kubernetes-ingress-controller/pull/436):
  - `--kong-url`, instead use `--kong-admin-url`
  - `--admin-tls-skip-verify`, instead use `--kong-admin-tls-skip-verify`
  - `--admin-header`, instead use `--kong-admin-header`
  - `--admin-tls-server-name`, instead use `--kong-admin-tls-server-name`
  - `--admin-ca-cert-file`, instead use `--kong-admin-ca-cert-file`

#### Under the hood

- decK has been bumped up to v0.6.2.

## [0.6.2]

> Release date: 2019/11/13

#### Summary

This is a minor patch release to fix version parsing issue with new
Kong Enterprise packages.

## [0.6.1]

> Release date: 2019/10/09

#### Summary

This is a minor patch release to update Kong Ingress Controller's
Docker image to use a non-root by default.

## [0.6.0]

> Release date: 2019/09/17

#### Summary

This release introduces an Admission Controller for CRDs,
Istio compatibility, support for `networking/ingress`,
Kong 1.3 additions and enhancements to documentation and deployments.

#### Added

- **Service Mesh integration** Kong Ingress Controller can now be deployed
  alongside Service Mesh solutions like Kuma and Istio. In such a deployment,
  Kong handles all the external client facing routing and policies while the
  mesh takes care of these aspects for internal service-to-service traffic.
- **`ingress.kubernetes.io/service-upstream`**, a new annotation has
  been introduced.
  Adding this annotation to a Kubernetes service resource
  will result in Kong directly forwarding traffic to kube-proxy.
  In other words, Kong will not send traffic directly to the pods.
  [#365](https://github.com/Kong/kubernetes-ingress-controller/pull/365)
- Ingress resources created in the new `networking.k8s.io` API group are
  now be supported. The controller dynamically figures out the API group
  to use based on the metadata it receives from k8s API-server.
- **Kong Credential enhancements**
  - Kong Credentials are now live-synced as they are created and updated in
    DB-mode.
    [#230](https://github.com/Kong/kubernetes-ingress-controller/issues/#230)
  - A single Consumer can now contain multiple credentials of the same type
    and multiple ACL group associations.
    [#371](https://github.com/Kong/kubernetes-ingress-controller/pull/371)
- **Admission controller** Kong Ingress Controller now ships with an in-built
  admission controller for KongPlugin and KongConsumer entities. The validations
  stop users from misconfiguring the Ingress controller.
  [#372](https://github.com/Kong/kubernetes-ingress-controller/pull/372)
- **Kong 1.3 support**:
  - HTTP Header based routing is now supported using `KongIngress.Route.Headers`
    property.
  - The algorithm to use for load-balancing traffic sent upstream can be
    set using `KongIngress.Upstream.Algorithm` field.
- **Kustomize**: Users can now use `kustomize` to tweak the reference deployment
  as per their needs. Both, DB and DB-less modes are supported. Please have
  a look at `deploy/manifests` directory in the Github repository.
- **Documentation**: The documentation for the project has been revamped.
  Deployment guides, how-to guides, and reference docs have been added.
- **Deployment**: The deployment of Kong Ingress Controller in DB and DB-less
  modes has been simplified, and Kong Ingress Controller now always runs as a
  side-car to Kong in proxy mode. There is no dedicated deployment for Kong
  Ingress Controller that needs to be run.

#### Fixed

- SNIs and Certificates are now de-duplicated across namespaces.
  [#360](https://github.com/Kong/kubernetes-ingress-controller/issues/#360)
  [#327](https://github.com/Kong/kubernetes-ingress-controller/issues/#327)
- Empty TLS secret no longer stops the controller from syncing configuration
  [#321](https://github.com/Kong/kubernetes-ingress-controller/issues/#321)
- Fix a nil reference when empty Ingress rules are created
  [#365](https://github.com/Kong/kubernetes-ingress-controller/pull/365)

#### Under the hood

- Kubernetes client-go library has been updated to v1.15.3.
- Credentials sync has been moved into decK and decK has been bumped up
  to v0.5.1.

## [0.5.0]

> Release date: 2019/06/25

#### Summary

This release introduces automated TLS certificates, consumer-level plugins,
enabling deployments using controller and Kong's Admin API at the same time
and numerous bug-fixes and enhancements.

#### Breaking changes

- UUID of consumers in Kong are no longer associated with UID of KongConsumer
  custom resource.

#### Added

- Kong 1.2 is now supported, meaning wild-card hosts in TLS section of Ingress
  resources are allowed.
- **Automated TLS certificates using Let's Encrypt**: Use Kong's Ingress
  Controller and
  [cert-manager](https://docs.cert-manager.io/en/latest/tasks/issuing-certificates/ingress-shim.html)
  to automatically provision TLS certs and serve them.
- **Tagging support**: All entities managed by Kong Ingress Controller in Kong's
  database are now tagged and the controller manages only a subset of Kong's
  configuration. Any entity created via Kong's Admin API will not be
  automatically deleted by the Ingress Controller.
- **Consumer-level plugins** can now be configured by applying
  `plugins.konghq.com` annotation on KongConsumer custom resources.
  [#250](https://github.com/Kong/kubernetes-ingress-controller/issues/#250)
- **Kong Enterprise workspaces**: Ingress Controller can manage a specific
  workspace inside Kong Enterprise (previously, only default workspace).
- Avoid reloading configuration in Kong in db-less mode when there is no
  change in configuration.
  [#308](https://github.com/Kong/kubernetes-ingress-controller/pull/308)
- Service scoped plugins for Kong 1.1 are now configured correctly.
  [#289](https://github.com/Kong/kubernetes-ingress-controller/issues/#289)

#### Fixed

- Multiple certificates are now correctly populated in Kong.
  [#285](https://github.com/Kong/kubernetes-ingress-controller/issues/#285)
- Missing entities like certificate secrets, services or plugins in Kubernetes
  object store will not stop controller from syncing configuration to Kong.
- A Ingress rule with an empty path is correctly parsed and populated in Kong.
  [#98](https://github.com/Kong/kubernetes-ingress-controller/issues/#98)
- Plugins with a nested schema are now correctly configured.
  [#294](https://github.com/Kong/kubernetes-ingress-controller/issues/#294)

#### Under the hood

- Dependency management for the project is done using Go modules.
- Kubernetes client-go library has been updated to v1.14.1.
- Makefile and Dockerfiles have been simplified.

## [0.4.0]

> Release date: 2019/04/24

#### Summary

This release introduces support to run Kong as an Ingress Controller
without a database!
This release comes with major under the hood rewrites to fix numerous
bugs and design issues in the codebase. Most of the syncing logic has
now been ported over to [decK](http://github.com/hbagdi/deck).

This release comes with a number of breaking changes.
Please read the changelog and test in your environment.

#### Breaking Changes

- :warning: Annotation `<plugin-name>.plugin.konghq.com`
  (deprecated in 0.2.0) is no longer supported.
- :warning: `--default-backend-service` CLI flag is now removed. The default
  service will now be picked up from the default backend in the Ingress rules.
- :warning: Service and Upstream entity overrides via KongIngress CRD are now
  supported only with `configuration.konghq.com` annotation on Kubernetes
  services.
  Route level overrides work same as before,
  using the `configuration.konghq.com` annotation on Ingress resources.
- :warning: `strip_path` property of Routes in Kong is set to `true` by default.
- :warning: `preserve_host` property of Routes in Kong is set to
  `true` by default.
- Plugins created for a combination of Route and Consumer using `consumerRef`
  property in KongPlugin CRD are not supported anymore. This functionality
  will be added back in future
  via [#250](https://github.com/Kong/kubernetes-ingress-controller/issues/250).
- Service and upstream Host name have changed from
  `namespace.service-name.port` to `service-name.namespace.svc`.

#### Added

- Ingress Controller now supports a DB-less deployment mode using Kong 1.1.
  [#244](https://github.com/Kong/kubernetes-ingress-controller/issues/244)
- New `run_on` and `protocols` properties are added to KongPlugin CRD.
  These can be used to further tune behaviors of plugins
  in Service Mesh deployments.
- New fields are added to KongIngress CRD to support HTTPS Active health-checks.
- Ingress Controller is now built using Go 1.12.
- Default service, which handles all traffic that is not matched against
  any of the Ingress rules, is now configured using the default backend
  defined via the Ingress resources.

#### Fixed

- Logs to stdout and stderr will be much more quieter and helpful and won't
  be as verbose as before.
- Routes with same path but different methods can now be created.
  [#202](https://github.com/Kong/kubernetes-ingress-controller/issues/202)
- Removing a value in KongPlugin config will now correctly sync it to Kong.
  [#117](https://github.com/Kong/kubernetes-ingress-controller/issues/117)
- Setting `--update-state=false` no longer causes a panic and performs leader
  election correctly.
  [#232](https://github.com/Kong/kubernetes-ingress-controller/issues/232)
  Thanks to [@lijiaocn](https://github.com/lijiaocn) for the fix!!
- KongIngress will now correctly override properites of Upstream object
  in Kong.
  [#252](https://github.com/Kong/kubernetes-ingress-controller/issues/252)
- Removing a value from KongPlugin config will now correctly unset it in
  Kong's datastore.
  [#117](https://github.com/Kong/kubernetes-ingress-controller/issues/117)

#### Under the hood

- Translation of Ingress rules and CRDs to Kong entities is completey
  re-written.
  [#241](https://github.com/Kong/kubernetes-ingress-controller/issues/241)
- For database deployments, an external tool, decK is used to sync resources
  to Kong, fixing numerous bugs and making Ingress Controller code saner
  and easier to maintain.

## [0.3.0]

> Release date: 2019/01/08

#### Breaking Changes

 - :warning: Default Ingress class is now `kong`.
   If you were relying on the previous default of `nginx`, you will
   need to explicitly set the class using `--ingress-class` CLI flag.

#### Added

- **Support for Kong 1.0.x** Kong 1.0 introduces a number of breaking changes
  in the Admin API. Ingress controller is updated to make correct calls
  and parse responses correctly.
  [#213](https://github.com/Kong/kubernetes-ingress-controller/pull/213)
- **ingress.class annotation-based filtering on CRD** Multiple Kong clusters
  can be deployed and configured individually on the same Kubernetes Cluster.
  This feature allows configuring
  global Plugins, Consumers & credentials
  using a different `ingress.class` annotation for each Kong cluster.
  [#220](https://github.com/Kong/kubernetes-ingress-controller/pull/220)
- **TLS support for Ingress Controller <-> Kong communication**
  The ingress controller can now talk to Kong's Control-Plane using TLS with
  custom certificates. Following new CLI flags are introduces:
  - `--admin-tls-skip-verify`: to skip validation of a certificate; it
  shouldn't be used in production environments.
  - `--admin-tls-server-name`: use this if the FQDN of Kong's Control Plane
  doesn't match the CN in the certificate.
  - `--admin-ca-cert-file`: use this to specify a custom CA cert which is
  not part of the bundled CA certs.
  [#212](https://github.com/Kong/kubernetes-ingress-controller/pull/212)

#### Fixed

- Retries for services in Kong can be set to zero.
   [#211](https://github.com/Kong/kubernetes-ingress-controller/pull/211)


## [0.2.2]

> Release date: 2018/11/09

#### Fixed

 - Fix plugin config comparison logic to avoid unnecessary PATCH requests
   to Kong
   [#196](https://github.com/Kong/kubernetes-ingress-controller/pull/196)
 - Fix `strip_path` in Routes in Kong. It is now set to false by default
   as in all other versions of Ingress controller except 0.2.1.
   [#194](https://github.com/Kong/kubernetes-ingress-controller/pull/194)
 - Fix path-only based Ingress rule parsing and configuration where only a
   path based rule for a Kubernetes Service
   would not setup Routes and Service in Kong.
   [#190](https://github.com/Kong/kubernetes-ingress-controller/pull/190)
 - Fix a nil pointer reference when overriding Ingress resource with KongIngress
   [#188](https://github.com/Kong/kubernetes-ingress-controller/pull/188)


## [0.1.3]

> Release date: 2018/11/09

#### Fixed

 - Fix path-only based Ingress rule parsing and configuration where only a
   path based rule for a Kubernetes Service
   would not setup Routes and Service in Kong.
   [#190](https://github.com/Kong/kubernetes-ingress-controller/pull/190)
 - Fix plugin config comparison logic to avoid unnecessary PATCH requests
   to Kong
   [#196](https://github.com/Kong/kubernetes-ingress-controller/pull/196)


## [0.2.1]

> Release date: 2018/10/26

#### Added

 - **Header Injection in requests to Kong's Admin API** HTTP Headers
   can be set via CLI which will be injected in every request sent to
   Kong's Admin API, enabling the use of Ingress Controller when Kong's
   Control Plane is protected by Authentication/Authorization.
   [#172](https://github.com/Kong/kubernetes-ingress-controller/pull/172)
 - **Path only based routing** Path only Ingress rules (without a host)
   are now parsed and served correctly.
   [#142](https://github.com/Kong/kubernetes-ingress-controller/pull/142)
 - Under the hood, an external library is now used to talk to Kong's Admin
   API. Several other packages and dead code has been dropped. These changes
   don't have any user facing changes but are steps in direction to simplify
   code and make it more testable.
   [#150](https://github.com/Kong/kubernetes-ingress-controller/pull/150)
   [#154](https://github.com/Kong/kubernetes-ingress-controller/pull/154)
   [#179](https://github.com/Kong/kubernetes-ingress-controller/pull/179)

#### Fixed

 - Fixed KongIngress overrides to enable overriding hashing attributes in
   Upstream object in Kong.
   Thanks @jdevalk2 for the patch!
   [#139](https://github.com/Kong/kubernetes-ingress-controller/pull/139)
 - Remove and sync certificates correctly when TLS secret reference changes
   for a hostname in Ingress spec.
   [#169](https://github.com/Kong/kubernetes-ingress-controller/pull/169)
 - Migrations for Kong are run using 'Job' in Kubernetes to avoid any
   issues that might arise due to multiple Kong nodes running migrations.
   [#161](https://github.com/Kong/kubernetes-ingress-controller/pull/161)
 - Kong and Ingress controller now wait for Postgres to start and migrations
   to finish before attempting to start.
   [#168](https://github.com/Kong/kubernetes-ingress-controller/pull/168)


## [0.1.2]

> Release date: 2018/10/26

#### Deprecated

 - :warning: Configuring plugins in Kong using `<plugin-name>.plugin.konghq.com`
   annotation is now deprecated and will be removed in a future release.
   Please use `plugins.konghq.com` annotation instead.

#### Added

 - **Header Injection in requests to Kong's Admin API** HTTP Headers
   can be set via CLI which will be injected in every request sent to
   Kong's Admin API, enabling the use of Ingress Controller when Kong's
   Control Plane is protected by Authentication/Authorization.
   [#172](https://github.com/Kong/kubernetes-ingress-controller/pull/172)
 - **Path only based routing** Path only Ingress rules (without a host)
   are now parsed and served correctly.
   [#142](https://github.com/Kong/kubernetes-ingress-controller/pull/142)
 - **Global Plugins** Plugins can be configured to run globally in Kong
   using a "global" label on `KongPlugin` resource.
   [#112](https://github.com/Kong/kubernetes-ingress-controller/pull/112)
 - A new property `plugin` has been introduced in `KongPlugin` resource
   which ties the configuration to be used and the type of the plugin.
   [#122](https://github.com/Kong/kubernetes-ingress-controller/pull/122)
 - Multiple plugins can be configured for an Ingress or a Service in k8s
   using `plugins.konghq.com` annotation.
   [#124](https://github.com/Kong/kubernetes-ingress-controller/pull/124)
 - `KongPlugin` resources do not need to be duplicated any more.
   The same `KongPlugin` resource can be used across
   multiple Ingress/Service resources.
   [#121](https://github.com/Kong/kubernetes-ingress-controller/pull/121)

#### Fixed

 - Avoid issuing unnecessary PATCH requests on Services in Kong during the
   reconciliation loop, which lead to unnecessary Router rebuilds inside Kong.
   [#107](https://github.com/Kong/kubernetes-ingress-controller/pull/107)
 - Fixed the diffing logic for plugin configuration between KongPlugin
   resource in k8s and plugin config in Kong to avoid false positives.
   [#106](https://github.com/Kong/kubernetes-ingress-controller/pull/106)
 - Correctly format IPv6 address for Targets in Kong.
   Thanks @NixM0nk3y for the patch!
   [#118](https://github.com/Kong/kubernetes-ingress-controller/pull/118)
 - Fixed KongIngress overrides to enable overriding hashing attributes in
   Upstream object in Kong.
   Thanks @jdevalk2 for the patch!
   [#139](https://github.com/Kong/kubernetes-ingress-controller/pull/139)
 - Remove and sync certificates correctly when TLS secret reference changes
   for a hostname in Ingress spec.
   [#169](https://github.com/Kong/kubernetes-ingress-controller/pull/169)


## [0.1.1]

> Release date: 2018/09/26

#### Fixed

 - Fix version parsing for minor releases of Kong Enterprise (like 0.33-1).
   The dash(`-`) didn't go well with the semver parsing
   [#141](https://github.com/Kong/kubernetes-ingress-controller/pull/141)

## [0.2.0]

> Release date: 2018/09/21

#### Breaking Changes

 - :warning: Support for Kong 0.13.x has been dropped in favor of 0.14.x

#### Deprecated

 - :warning: Configuring plugins in Kong using `<plugin-name>.plugin.konghq.com`
   annotation is now deprecated and will be removed in a future release.
   Please use `plugins.konghq.com` annotation instead.

#### Added

 - **Support for Kong 0.14.x** The supported version of Kong 0.14.x
   has been introduced. Kong 0.14.x introduced breaking changes to a few
   Admin API endpoints which have been updated in the Ingress Controller.
   [#101](https://github.com/Kong/kubernetes-ingress-controller/pull/101)
 - **Global Plugins** Plugins can be configured to run globally in Kong
   using a "global" label on `KongPlugin` resource.
   [#112](https://github.com/Kong/kubernetes-ingress-controller/pull/112)
 - A new property `plugin` has been introduced in `KongPlugin` resource
   which ties the configuration to be used and the type of the plugin.
   [#122](https://github.com/Kong/kubernetes-ingress-controller/pull/122)
 - Multiple plugins can be configured for an Ingress or a Service in k8s
   using `plugins.konghq.com` annotation.
   [#124](https://github.com/Kong/kubernetes-ingress-controller/pull/124)
 - `KongPlugin` resources do not need to be duplicated any more.
   The same `KongPlugin` resource can be used across
   multiple Ingress/Service resources.
   [#121](https://github.com/Kong/kubernetes-ingress-controller/pull/121)
 - The custom resource definitions now have a short-name for all the
   CRDs, making it easy to interact with `kubectl`.
   [#120](https://github.com/Kong/kubernetes-ingress-controller/pull/120)

#### Fixed

 - Avoid issuing unnecessary PATCH requests on Services in Kong during the
   reconciliation loop, which lead to unnecessary Router rebuilds inside Kong.
   [#107](https://github.com/Kong/kubernetes-ingress-controller/pull/107)
 - Fixed the diffing logic for plugin configuration between KongPlugin
   resource in k8s and plugin config in Kong to avoid false positives.
   [#106](https://github.com/Kong/kubernetes-ingress-controller/pull/106)
 - Correctly format IPv6 address for Targets in Kong.
   Thanks @NixM0nk3y for the patch!
   [#118](https://github.com/Kong/kubernetes-ingress-controller/pull/118)


## [0.1.0]

> Release date: 2018/08/17

#### Breaking Changes

 - :warning: **Declarative Consumers in Kong** Kong consumers can be
   declaratively configured via `KongConsumer` custom resources. Any consumers
   created directly in Kong without a corresponding `KongConsumer` custom
   resource will be deleted by the ingress controller.
   [#81](https://github.com/Kong/kubernetes-ingress-controller/pull/81)

#### Added

 - **Support Upstream TLS** Service in Kong can be configured to use HTTPS
   via `KongIngress` custom resource.
   [#79](https://github.com/Kong/kubernetes-ingress-controller/pull/79)
 - Support for control over protocol(HTTP/HTTPS) to use for ingress traffic
   via `KongIngress` custom resource.
   [#64](https://github.com/Kong/kubernetes-ingress-controller/pull/64)

#### Fixed

 - Multiple SNIs are created in Kong if multiple hosts are specified in TLS
   section of an `Ingress` resource.
   [#76](https://github.com/Kong/kubernetes-ingress-controller/pull/76)
 - Updates to `KongIngress` resource associated with an Ingress
   now updates the corresponding routing properties in Kong.
   [#92](https://github.com/Kong/kubernetes-ingress-controller/pull/92)


## [v0.0.5]

> Release date: 2018/06/02

#### Added

 - Add support for Kong Enterprise Edition 0.32 and above

## [v0.0.4] and prior

 - The initial versions  were rapildy iterated to deliver
   a working ingress controller.

[2.1.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.6...v2.1.0
[2.0.6]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.5...v2.0.6
[2.0.5]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.4...v2.0.5
[2.0.4]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.3...v2.0.4
[2.0.3]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.2...v2.0.3
[2.0.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.0-beta.2...v2.0.0
[1.3.3]: https://github.com/kong/kubernetes-ingress-controller/compare/1.3.2...1.3.3
[1.3.2]: https://github.com/kong/kubernetes-ingress-controller/compare/1.3.1...1.3.2
[1.3.1]: https://github.com/kong/kubernetes-ingress-controller/compare/1.3.0...1.3.1
[1.3.0]: https://github.com/kong/kubernetes-ingress-controller/compare/1.2.0...1.3.0
[1.2.0]: https://github.com/kong/kubernetes-ingress-controller/compare/1.1.1...1.2.0
[1.1.1]: https://github.com/kong/kubernetes-ingress-controller/compare/1.1.0...1.1.1
[1.1.0]: https://github.com/kong/kubernetes-ingress-controller/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.10.0...1.0.0
[0.10.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.9.1...0.10.0
[0.9.1]: https://github.com/kong/kubernetes-ingress-controller/compare/0.9.0...0.9.1
[0.9.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.8.1...0.9.0
[0.8.1]: https://github.com/kong/kubernetes-ingress-controller/compare/0.8.0...0.8.1
[0.8.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.7.1...0.8.0
[0.7.1]: https://github.com/kong/kubernetes-ingress-controller/compare/0.7.0...0.7.1
[0.7.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.6.2...0.7.0
[0.6.2]: https://github.com/kong/kubernetes-ingress-controller/compare/0.6.1...0.6.2
[0.6.1]: https://github.com/kong/kubernetes-ingress-controller/compare/0.6.0...0.6.1
[0.6.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.5.0...0.6.0
[0.5.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.4.0...0.5.0
[0.4.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.3.0...0.4.0
[0.3.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.2.2...0.3.0
[0.2.2]: https://github.com/kong/kubernetes-ingress-controller/compare/0.2.1...0.2.2
[0.1.3]: https://github.com/kong/kubernetes-ingress-controller/compare/0.1.2...0.1.3
[0.2.1]: https://github.com/kong/kubernetes-ingress-controller/compare/0.2.0...0.2.1
[0.1.2]: https://github.com/kong/kubernetes-ingress-controller/compare/0.1.1...0.1.2
[0.1.1]: https://github.com/kong/kubernetes-ingress-controller/compare/0.1.0...0.1.1
[0.2.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v0.0.5...0.1.0
[v0.0.5]: https://github.com/kong/kubernetes-ingress-controller/compare/v0.0.4...v0.0.5
[v0.0.4]: https://github.com/kong/kubernetes-ingress-controller/compare/7866a27f268c32c5618fba546da2c73ba74d4a46...v0.0.4
