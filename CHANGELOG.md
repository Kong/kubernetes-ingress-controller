# Table of Contents

<!---
Adding a new version? You'll need three changes:
* Add the ToC link, like "[1.2.3](#123)".
* Add the section header, like "## [1.2.3]".
* Add the diff link, like "[2.7.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v1.2.2...v1.2.3".
  This is all the way at the bottom. It's the thing we always forget.
--->
 - [3.2.3](#323)
 - [3.2.2](#322)
 - [3.2.1](#321)
 - [3.2.0](#320)
 - [3.1.6](#316)
 - [3.1.5](#315)
 - [3.1.4](#314)
 - [3.1.3](#313)
 - [3.1.2](#312)
 - [3.1.1](#311)
 - [3.1.0](#310)
 - [3.0.2](#302)
 - [3.0.1](#301)
 - [3.0.0](#300)
 - [2.12.5](#2125)
 - [2.12.4](#2124)
 - [2.12.3](#2123)
 - [2.12.2](#2122)
 - [2.12.1](#2121)
 - [2.12.0](#2120)
 - [2.11.1](#2111)
 - [2.11.0](#2110)
 - [2.10.5](#2105)
 - [2.10.4](#2104)
 - [2.10.3](#2103)
 - [2.10.2](#2102)
 - [2.10.1](#2101)
 - [2.10.0](#2100)
 - [2.9.3](#293)
 - [2.9.2](#292)
 - [2.9.1](#291)
 - [2.9.0](#290)
 - [2.8.2](#282)
 - [2.8.1](#281)
 - [2.8.0](#280)
 - [2.7.0](#270)
 - [2.6.0](#260)
 - [2.5.0](#250)
 - [2.4.2](#242)
 - [2.4.1](#241)
 - [2.4.0](#240)
 - [2.3.1](#231)
 - [2.3.0](#230)
 - [2.2.1](#221)
 - [2.2.0](#220)
 - [2.1.1](#211)
 - [2.1.0](#210)
 - [2.0.7](#207)
 - [2.0.6](#206)
 - [2.0.5](#205)
 - [2.0.4](#204)
 - [2.0.3](#203)
 - [2.0.2](#202)
 - [2.0.1](#201)
 - [2.0.0](#200)
 - [1.3.4](#134)
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

## Unreleased

### Added

- `KongCustomEntity` is now supported by the `FallbackConfiguration` feature.
  [#6286](https://github.com/Kong/kubernetes-ingress-controller/pull/6286)
- It is now possible to disable synchronization of consumers to Konnect through the
  flag `--konnect-disable-consumers-sync`.
  [#6313](https://github.com/Kong/kubernetes-ingress-controller/pull/6313)
- Allow `KongCustomEntity` to refer to plugins in another namespace via 
  `spec.parentRef.namespace`. The reference is allowed only when there is a
  `ReferenceGrant` in the namespace of the `KongPlugin` to grant permissions
  to `KongCustomEntity` of referring to `KongPlugin`.
  [#6289](https://github.com/Kong/kubernetes-ingress-controller/pull/6289)
- Konnect configuration updates are now handled separately from gateway
  updates. This allows the controller to handle sync errors for the gateway and
  Konnect speparately, and avoids one blocking the other.
  [#6341](https://github.com/Kong/kubernetes-ingress-controller/pull/6341)
  [#6349](https://github.com/Kong/kubernetes-ingress-controller/pull/6349)
- Added `duration` field in logs after successfully sent configuration to Kong
  gateway or Konnect.
  [#6360](https://github.com/Kong/kubernetes-ingress-controller/pull/6360)
- `KongCustomEntity` is now included in last valid configuration retrieved from
  Kong gateways.
  [#6305](https://github.com/Kong/kubernetes-ingress-controller/pull/6305)
- Added `/debug/config/diff-report` diagnostic endpoint. This endpoint is
  available in DB mode when the `--dump-config` and `--dump-sensitive-config`
  are enabled. It returns the latest diff information for the controller's last
  configuration sync along with config hash and sync timestamp metadata. The
  controller maintains the last 5 diffs in cache. You can retrieve older diffs
  by appending a `?hash=<hash>` query string argument. Available config hashes
  and their timestamps are listed under the `available` section of the
  response.
  [#6131](https://github.com/Kong/kubernetes-ingress-controller/pull/6131)

### Fixed

- Services using `Secret`s containing the same certificate as client certificates
  by annotation `konghq.com/client-cert` can be correctly translated.
  [#6228](https://github.com/Kong/kubernetes-ingress-controller/pull/6228)
- Generate one entity for each attached foreign entity if a `KongCustomEntity`
  resource is attached to multiple foreign Kong entities.
  [#6280](https://github.com/Kong/kubernetes-ingress-controller/pull/6280)

### Changed

- Check Kong Gateway readiness concurrently. This greatly reduces the time which
  is required to check all Gateway instances readiness, especially when there's many
  of them. Increased individual readiness check timeout from 1s to 5s.
  [#6347](https://github.com/Kong/kubernetes-ingress-controller/pull/6347)
  [#6357](https://github.com/Kong/kubernetes-ingress-controller/pull/6357)

## 3.2.3

> Release date: 2024-07-23

### Fixed

- Fixed the reference checker in checking permission of remote plugins to use
  the correct namespace of `ReferenceGrant` required. Add trace logging to
  `ReferenceGrant` check functions.
  [#6295](https://github.com/Kong/kubernetes-ingress-controller/pull/6295)
  [#6302](https://github.com/Kong/kubernetes-ingress-controller/pull/6302)

## 3.2.2

> Release date: 2024-07-01

### Fixed

- Fixed an issue where new gateways were not being populated with the current configuration when
  `FallbackConfiguration` feature gate was turned on. Previously, configuration updates were skipped
  if the Kubernetes config cache did not change, leading to inconsistencies. Now, the system ensures
  that all gateways are populated with the latest configuration regardless of cache changes.
  [#6271](https://github.com/Kong/kubernetes-ingress-controller/pull/6271)

## 3.2.1

> Release date: 2024-06-28

### Fixed

- Do not try recovering from gateways synchronization errors with fallback configuration
  (either generated or the last valid one) when an unexpected error (e.g. 5xx or network issue) occurs.
  [#6237](https://github.com/Kong/kubernetes-ingress-controller/pull/6237)
- Admission webhook will accept multiple plugins of the same type associated with a single route-like,
  Service, KongConsumer, KongConsumerGroup object to allow plugins to be associated with combinations
  of those objects.
  [#6252](https://github.com/Kong/kubernetes-ingress-controller/pull/6252)

## 3.2.0

> Release date: 2024-06-12

### Highlights

- 🚀 **Fallback Configuration**: New `FallbackConfiguration` feature enables isolating configuration failure domains so
  that one broken object no longer prevents the entire configuration from being applied. See [Fallback Configuration guide]
  to learn more.
- 🏗️ **Custom Kong Entities**: New `KongCustomEntity` CRD allows defining Kong custom entities programmatically in KIC.
  See [Using Custom Entities guide] to learn more.
- 📨 **GRPCRoute v1 support**: Following the GA graduation in the Gateway API, KIC now supports the v1 version of the GRPCRoute.
  See [GRPCRoute reference] to learn more. _Requires upgrading the Gateway API's CRDs to v1.1._

[Fallback Configuration guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/high-availability/fallback-config/
[Using Custom Entities guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/services/custom-entity
[GRPCRoute reference]: https://gateway-api.sigs.k8s.io/api-types/grpcroute/

### Breaking changes

- Removed support for the deprecated `kongCredType` Secret field. If you have
  not previously [updated to the credential label](https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/migrate/credential-kongcredtype-label/)
  you must do so before upgrading to this version. This removal includes an
  update to the webhook configuration that checks only Secrets with
  `konghq.com/credential` or `konghq.com/validate` labels (for Secrets that
  contain plugin configuration). This filter improves performance and
  reliability by not checking Secrets the controller will never use. Users that
  wish to defer adding `konghq.com/validate` to Secrets with plugin
  configuration can set the `ingressController.admissionWebhook.filterSecrets`
  chart values.yaml key to `false`. Doing so does not benefit from the
  performance benefits, however, so labeling plugin configuration Secrets and
  enabling the filter is recommended as soon as is convenient.
  [#5856](https://github.com/Kong/kubernetes-ingress-controller/pull/5856)
- Dynamically set the proxy protocol of GRPCRoute to `grpc` or `grpcs` based on the port listened by Gateway.
  If you don't set the protocol for Service via `konghq.com/protocol` annotation, Kong will use `grpc` instead of `grpcs`.
  [#5776](https://github.com/Kong/kubernetes-ingress-controller/pull/5776)
- The `/debug/config/failed` and `/debug/config/successful` diagnostic
  endpoints now nest configuration dumps under a `config` key. These endpoints
  previously returned the configuration dump at the root. They now return
  additional metadata along with the configuration. This change should not
  impact normal usage, but if you scrape these endpoints, be aware that their
  output format has changed.
  [#6101](https://github.com/Kong/kubernetes-ingress-controller/pull/6101)

### Added

- Added `FallbackConfiguration` feature gate to enable the controller to generate a fallback configuration
  for Kong when it fails to apply the original one. The feature gate is disabled by default.
  [#5993](https://github.com/Kong/kubernetes-ingress-controller/pull/5993)
  [#6010](https://github.com/Kong/kubernetes-ingress-controller/pull/6010)
  [#6047](https://github.com/Kong/kubernetes-ingress-controller/pull/6047)
  [#6071](https://github.com/Kong/kubernetes-ingress-controller/pull/6071)
- Added `--use-last-valid-config-for-fallback` CLI flag to enable using the last valid configuration cache
  to backfill excluded broken objects when the `FallbackConfiguration` feature gate is enabled.
  [#6098](https://github.com/Kong/kubernetes-ingress-controller/pull/6098)
- Added `FallbackKongConfigurationSucceeded`, `FallbackKongConfigurationTranslationFailed` and
  `FallbackKongConfigurationApplyFailed` Kubernetes Events to report the status of the fallback configuration.
  [#6099](https://github.com/Kong/kubernetes-ingress-controller/pull/6099)
- Added Prometheus metrics covering `FallbackConfiguration` feature:
  - `ingress_controller_fallback_translation_count`
  - `ingress_controller_fallback_translation_broken_resource_count`
  - `ingress_controller_fallback_configuration_push_count`
  - `ingress_controller_fallback_configuration_push_last`
  - `ingress_controller_fallback_configuration_push_duration_milliseconds`
  - `ingress_controller_fallback_configuration_push_broken_resource_count`
  - `ingress_controller_fallback_cache_generating_duration_milliseconds`
  - `ingress_controller_processed_config_snapshot_cache_hit`
  - `ingress_controller_processed_config_snapshot_cache_miss`
    [#6105](https://github.com/Kong/kubernetes-ingress-controller/pull/6105)
- Added a `GET /debug/config/fallback` diagnostic endpoint to expose the fallback configuration
  details (currently broken, excluded and backfilled objects, as well as the overall status).
  [#6184](https://github.com/Kong/kubernetes-ingress-controller/pull/6184)
- Added the CRD `KongCustomEntity` to support custom Kong entities that are not
  defined in KIC yet. The current version only supports translating custom
  entities into declarative configuration in DBless mode, and cannot apply
  custom entities to DB backed Kong gateways.
  Feature gate `KongCustomEntity` is required to set to `true` to enabled the
  `KongCustomEntity` controller.
  **Note**: The IDs of Kong services, routes and consumers referred by custom
  entities via `foreign` type fields of custom entities are filled by the `FillID`
  method of the corresponding type because the IDs of these entities are required
  to fill the `foreign` fields of custom entities. So the `FillIDs` feature gate
  is also required when `KongCustomEntity` is enabled.
  [#5982](https://github.com/Kong/kubernetes-ingress-controller/pull/5982)
  [#6006](https://github.com/Kong/kubernetes-ingress-controller/pull/6006)
  [#6055](https://github.com/Kong/kubernetes-ingress-controller/pull/6055)
- Add support for Kubernetes Gateway API v1.1:
  - add a flag `--enable-controller-gwapi-grpcroute` to control whether enable or disable GRPCRoute controller.
  - add support for `GRPCRoute` v1, which requires users to upgrade the Gateway API's CRD to v1.1.
    [#5918](https://github.com/Kong/kubernetes-ingress-controller/pull/5918)
- Add a `/debug/config/raw-error` endpoint to the config dump diagnostic
  server. This endpoint outputs the original Kong `/config` endpoint error for
  failed configuration pushes in case error parsing fails. Attempt to log the
  `message` field of errors that KIC cannot fully parse.
  [#5773](https://github.com/Kong/kubernetes-ingress-controller/issues/5773), [#5846](https://github.com/Kong/kubernetes-ingress-controller/pull/5846)
- Add constraint to limit items in `Credentials` and `ConsumerGroups` in
  `KongConsumer`s to be unique by defining their `x-kubernetes-list-type` as `set`.
  Please note that if you're using `helm` as the installation method, upgrading alone
  won't make this change take effect until you manually update the CRD manifests in your
  cluster to the current version. See [Updates to CRDs] for more details.
  [#5894](https://github.com/Kong/kubernetes-ingress-controller/pull/5894)
- Add support in `HTTPRoute`s for `URLRewrite`:
  - `FullPathRewrite` [#5855](https://github.com/Kong/kubernetes-ingress-controller/pull/5855)
  - `ReplacePrefixMatch` for both router modes:
    - `traditional_compatible` [#5895](https://github.com/Kong/kubernetes-ingress-controller/pull/5895)
    - `expressions` [#5940](https://github.com/Kong/kubernetes-ingress-controller/pull/5940)
  - `Hostname` [#5952](https://github.com/Kong/kubernetes-ingress-controller/pull/5952)
- DB mode now supports Event reporting for resources that failed to apply.
  [#5785](https://github.com/Kong/kubernetes-ingress-controller/pull/5785)
- Improve validation - reject `Ingresses`, `Services`, `HTTPRoutes`, `KongConsumers` or `KongConsumerGroups`
  that have multiple instances of the same type plugin attached.
  [#5972](https://github.com/Kong/kubernetes-ingress-controller/pull/5972)
  [#5979](https://github.com/Kong/kubernetes-ingress-controller/pull/5979)
- Added support for `konghq.com/headers-separator` that allows to set custom separator (instead of default `,`)
  for headers specified with `konghq.com/headers.*` annotations. Moreover parsing a content of `konghq.com/headers.*`
  is more robust - leading and trailing whitespace characters are discarded.
  [#5977](https://github.com/Kong/kubernetes-ingress-controller/pull/5977)
- The `konghq.com/plugins` annotation supports a new `<namespace>:<name>`
  format. This format requests a KongPlugin from a remote namespace. Binding
  plugins across namespaces requires a ReferenceGrant from the requesting
  resource to KongPlugins in the target namespace. This approach is useful for
  some plugins bound to different types of entities, such as a set of
  rate-limiting plugins applied to a service and various consumers. The
  cross-namespace grant allows the service manager to define different limits
  for consumers managed by other users without requiring those users to create
  consumers in the Service's namespace.
  [#5965](https://github.com/Kong/kubernetes-ingress-controller/pull/5965)
- The last valid configuration no longer omits licenses and vaults.
  [#6048](https://github.com/Kong/kubernetes-ingress-controller/pull/6048)
- Add support for Gateway API GRPCRoute and pass related Gateway API conformance test.
  [#5776](https://github.com/Kong/kubernetes-ingress-controller/pull/5776)

### Fixed

- Set proper User-Agent for request made to Kong and Konnect.
  [#5753](https://github.com/Kong/kubernetes-ingress-controller/pull/5753)
- Reconcile Secrets with `konghq.com/credential` label instead of waiting for other
  object to contain a reference to that Secrets
  [#5816](https://github.com/Kong/kubernetes-ingress-controller/pull/5816)
- Support to apply licenses to DB backed Kong gateway from `KongLicense`.
  [#5648](https://github.com/Kong/kubernetes-ingress-controller/pull/5648)
- Do not generate invalid duplicate upstream targets when routes use multiple
  Services with the same endpoints.
  [#5817](https://github.com/Kong/kubernetes-ingress-controller/pull/5817)
- Remove the constraint of items of `parentRefs` can only be empty or 
  `gateway.network.k8s.io/Gateway` in validating `HTTPRoute`s. If an item in
  `parentRefs`'s group/kind is not `gateway.network.k8s.io/Gateway`, the item
  is seen as a parent other than the controller and ignored in parentRef check.
  [#5919](https://github.com/Kong/kubernetes-ingress-controller/pull/5919)
- Redacted values no longer cause collisions in configuration reported to Konnect.
  [#5964](https://github.com/Kong/kubernetes-ingress-controller/pull/5964)
- The `--dump-sensitive-config` flag is no longer backwards.
  [#6073](https://github.com/Kong/kubernetes-ingress-controller/pull/6073)
- Fixed KIC clearing Gateway API *Route status of routes that it shouldn't reconcilce, e.g.
  those attached to Gateways that do not belong to GatewayClass that KIC reconciles.
  [#6079](https://github.com/Kong/kubernetes-ingress-controller/pull/6079)
- Fixed KIC non leaders correctly getting up to date Admin API addresses by not
  requiring leader election for the related controller.
  [#6126](https://github.com/Kong/kubernetes-ingress-controller/pull/6126)
- Plugins attached to both a KongConsumerGroup and a route-like resource or
  Service now properly generate a plugin attached to both a Kong consumer group
  and route or service. Previously, these incorrectly generated plugins
  attached to the route or service only.
  [#6132](https://github.com/Kong/kubernetes-ingress-controller/pull/6132)
- KongPlugin's `config` field is no longer incorrectly sanitized.
  [#6138](https://github.com/Kong/kubernetes-ingress-controller/pull/6138)

### Changed

- Preallocate slices for Gateway API objects when listing in store.
  This yields a significant performance improvements in time spent, bytes allocated
  and allocations per list operation.
  [#5824](https://github.com/Kong/kubernetes-ingress-controller/pull/5824)

[Updates to CRDs]: https://github.com/Kong/charts/blob/main/charts/kong/UPGRADE.md#updates-to-crds

## [3.1.6]

> Release date: 2024-06-11

### Fixed

- Konnect instances report correct plugin configuration to Konnect.
  [#6138](https://github.com/Kong/kubernetes-ingress-controller/pull/6138)
- Plugins attached to both a KongConsumerGroup and a route-like resource or
  Service now properly generate a plugin attached to both a Kong consumer group
  and route or service. Previously, these incorrectly generated plugins
  attached to the route or service only.
  [#6132](https://github.com/Kong/kubernetes-ingress-controller/pull/6132)


## [3.1.5]

> Release date: 2024-05-17

### Fixed

- Support to apply licenses to DB backed Kong gateway from `KongLicense`.
  [#5648](https://github.com/Kong/kubernetes-ingress-controller/pull/5648)
- Redacted values no longer cause collisions in configuration reported to Konnect.
  [#5964](https://github.com/Kong/kubernetes-ingress-controller/pull/5964)
- Assign a default value for `weight` in Kong target if the `weight` is nil.
  [#5946](https://github.com/Kong/kubernetes-ingress-controller/pull/5946)

## [3.1.4]

> Release date: 2024-04-26

### Fixed

- Do not generate invalid duplicate upstream targets when routes use multiple
  Services with the same endpoints.
  [#5817](https://github.com/Kong/kubernetes-ingress-controller/pull/5817)
- Bump golang version to 1.21.9 to fix CVE [GO-2024-2687](https://pkg.go.dev/vuln/GO-2024-2687).
  [#5905](https://github.com/Kong/kubernetes-ingress-controller/pull/5905)

## [3.1.3]

> Release date: 2024-04-08

### Fixed

- Remove unnecessary tag support check that could incorrectly delete configuration if the check did not execute properly.
  [#5658](https://github.com/Kong/kubernetes-ingress-controller/issues/5658)
- Do not require `rsa_public_key` field in credential `Secret`s when working with jwt HMAC credentials.
  [#5737](https://github.com/Kong/kubernetes-ingress-controller/issues/5737)
- `KongUpstreamPolicy` controller no longer requires existence of `HTTPRoute` CRD
  to start.
  [#5780](https://github.com/Kong/kubernetes-ingress-controller/pull/5780)
- Do not require namespaces when parsing errors about cluster scoped resources
  [#5764](https://github.com/Kong/kubernetes-ingress-controller/issues/5764)

## [3.1.2]

> Release date: 2024-03-08

### Fixed

- When managed Kong gateways are OSS edition, KIC will not apply licenses to
  the Kong gateway instances to avoid invalid configurations.
  [#5640](https://github.com/Kong/kubernetes-ingress-controller/pull/5640)

## [3.1.1]

> Release date: 2024-02-29

### Added

- Managed `Gateway`s now get reconciled by the `Gateway` controller, but do not
  get their status updated, they only become part of the configuration published
  to Kong.
  [#5662](https://github.com/Kong/kubernetes-ingress-controller/pull/5662)

### Fixed

- Fixed an issue where single-`Gateway` mode did not actually filter out routes
  associated with other `Gateway`s in the controller class.
  [#5642](https://github.com/Kong/kubernetes-ingress-controller/pull/5642)

## [3.1.0]

> Release date: 2024-02-07

### Highlights

- 🔒 Kong Gateway's [secret vaults][kong-vault] now become a first-class citizen for Kubernetes users thanks to the new
  `KongVault` CRD. _See [Kong Vault guide][vault-guide] and [CRDs reference][crds-ref] for more details._
- 🔒 Providing an Enterprise license to KIC-managed Kong Gateways becomes much easier thanks to the new `KongLicense`
  CRD which is used to dynamically provision all the replicas with the latest license found in the cluster. _See
  [Enterprise License][license-guide] and [CRDs reference][crds-ref] for more details._
- ✨ Populating a single field of `KongPlugin`'s configuration with use of a Kubernetes Secret becomes possible thanks
  to the new `KongPlugin`'s `configPatches` field. _See [Using Kubernetes Secrets in Plugins][secrets-in-plugins-guide]
  and [CRDs reference][crds-ref] for more details._
- 🔒 All sensitive information stored in the cluster is now sanitized while sending configuration to Konnect.

[crds-ref]: https://docs.konghq.com/kubernetes-ingress-controller/latest/reference/custom-resources/
[vault-guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/security/kong-vault/
[license-guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/license/
[secrets-in-plugins-guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/security/plugin-secrets/

### Added

- New CRD `KongVault` to represent a custom Kong vault for storing sensitive
  data used in plugin configurations. Now users can create a `KongVault` to
  create a custom Kong vault and reference the values in configurations of
  plugins. Reference of using Kong vaults: [Kong vault](kong-vault). Since the prefix of
  Kong vault is restrained unique, the `spec.prefix` field is set to immutable,
  and only one of multiple `KongVault`s with the same `spec.prefix` will get
  translated. Translation failiure events will be recorded for others with
  duplicating `spec.prefix`.
  [#5354](https://github.com/Kong/kubernetes-ingress-controller/pull/5354)
  [#5384](https://github.com/Kong/kubernetes-ingress-controller/pull/5384)
  [#5435](https://github.com/Kong/kubernetes-ingress-controller/pull/5435)
  [#5412](https://github.com/Kong/kubernetes-ingress-controller/pull/5412)
- New CRD `KongLicense` to represent a Kong enterprise license to apply to
  managed Kong gateway enterprise instances. The `Enabled` field of `KongLicense`
  (set to `True` if not present) need to be set to true to get reconciled.
  If there are multiple `KongLicense`s in the cluster, the newest one
  (with latest `metadata.creationTimestamp`) is chosen. The `KongLicense`
  controller is disabled when synchoroniztion of license with Konnect is turned
  on. When sync license with Konnect is turned on, licenses from Konnect are used.  
  [#5487](https://github.com/Kong/kubernetes-ingress-controller/pull/5487)
  [#5514](https://github.com/Kong/kubernetes-ingress-controller/pull/5514)
- Added `configPatches` field to KongPlugin and KongClusterPlugin to
  support populating configuration fields from Secret values. An item in
  `configPatches` defines a JSON patch to add a field on the path in its `path`
  and value from the value in the secret given in `valueFrom`. The JSON patches
  are applied to the raw JSON in `config`. Can only be specified when
  `configFrom` is not present.
  [#5158](https://github.com/Kong/kubernetes-ingress-controller/pull/5158)
  [#5208](https://github.com/Kong/kubernetes-ingress-controller/pull/5208)
- Added `SanitizeKonnectConfigDumps` feature gate allowing to enable sanitizing
  sensitive data (like TLS private keys, Secret-sourced Plugins configuration, etc.)
  in Konnect configuration dumps. It's turned on by default.
  [#5489](https://github.com/Kong/kubernetes-ingress-controller/pull/5489)
  [#5573](https://github.com/Kong/kubernetes-ingress-controller/pull/5573)
- Kong Plugin's `config` field now is sanitized when it contains sensitive data
  sourced from a Secret (i.e. `configFrom` or `configPatches` is used).
  [#5495](https://github.com/Kong/kubernetes-ingress-controller/pull/5495)
- `KongServiceFacade` CRD allowing creating Kong Services directly from Kubernetes using
  Kubernetes Services as their backends. `KongServiceFacade` can be used as a backend in
  Kubernetes Ingress. This API is highly experimental and is not distributed by default.
  It can be installed with `kubectl kustomize "github.com/Kong/kubernetes-ingress-controller/config/crd/incubator/?ref=v3.1.0"`
  When installed, it has to be enabled with `ServiceFacade` feature gate.
  [#5220](https://github.com/Kong/kubernetes-ingress-controller/pull/5220)
  [#5234](https://github.com/Kong/kubernetes-ingress-controller/pull/5234)
  [#5290](https://github.com/Kong/kubernetes-ingress-controller/pull/5290)
  [#5282](https://github.com/Kong/kubernetes-ingress-controller/pull/5282)
  [#5298](https://github.com/Kong/kubernetes-ingress-controller/pull/5298)
  [#5302](https://github.com/Kong/kubernetes-ingress-controller/pull/5302)
- Added support for GRPC over HTTP (without TLS) in Gateway API.
  [#5128](https://github.com/Kong/kubernetes-ingress-controller/pull/5128)
  [#5283](https://github.com/Kong/kubernetes-ingress-controller/pull/5283)
- Added `--init-cache-sync-duration` CLI flag. This flag configures how long the
  controller waits for Kubernetes resources to populate at startup before
  generating the initial Kong configuration. It also fixes a bug that removed
  the default 5 second wait period.
  [#5238](https://github.com/Kong/kubernetes-ingress-controller/pull/5238)
- Added `--emit-kubernetes-events` CLI flag to disable the creation of events
  in translating and applying configurations to Kong.
  [#5296](https://github.com/Kong/kubernetes-ingress-controller/pull/5296)
  [#5299](https://github.com/Kong/kubernetes-ingress-controller/pull/5299)
- Added validation on `Secret`s to reject the change if it will generate
  invalid confiugration of plugins for `KongPlugin`s or `KongClusterPlugin`s
  referencing to the secret.
  [#5203](https://github.com/Kong/kubernetes-ingress-controller/pull/5203)
- Validate `HTTPRoute` in admission webhook and reject it if the spec uses
  the following features that we do not support:
  - `parentRefs` other than `gatewayapi.networking.k8s.io/Gateway`
  - using `timeouts` in rules
  - `URLRewrite`, `RequestMirror` filters
  - using filters in backendRefs of rules
  [#5312](https://github.com/Kong/kubernetes-ingress-controller/pull/5312)
- Added functionality to the `KongUpstreamPolicy` controller to properly set and
  enforce `KongUpstreamPolicy` status. 
  The controller will set an ancestor status in `KongUpstreamPolicy` status for 
  each of its ancestors (i.e. `Service` or `KongServiceFacade`) with the `Accepted`
  and `Programmed` condition.
  [#5185](https://github.com/Kong/kubernetes-ingress-controller/pull/5185)
  [#5428](https://github.com/Kong/kubernetes-ingress-controller/pull/5428)
  [#5444](https://github.com/Kong/kubernetes-ingress-controller/pull/5444)
- Added flag `--gateway-to-reconcile` to set KIC to only reconcile
  the specified Gateway resource in Kubernetes.
  [#5405](https://github.com/Kong/kubernetes-ingress-controller/pull/5405)
- Added support for `HTTPRouteTimeoutBackendRequest` in Gateway API.
  Due to only one field being available in the Gateway API to control this behavior,
  when users set `spec.rules[].timeouts` in HTTPRoute,
  KIC will set `connect_timeout`, `read_timeout` and `write_timeout` for the service to this value.
  It's only possible to set the same timeout for each rule in a single `HTTPRoute`. Other settings
  will be rejected by the admission webhook validation.
  [#5243](https://github.com/Kong/kubernetes-ingress-controller/pull/5243)
- Log the details in response from Konnect when failed to push configuration
  to Konnect.
  [#5453](https://github.com/Kong/kubernetes-ingress-controller/pull/5453)

### Fixed

- Validators of `KongPlugin` and `KongClusterPlugin` will not return `500` on
  failures to parse configurations and failures to retrieve secrets used for
  configuration. Instead, it will return `400` with message to tell the
  validation failures.
  [#5208](https://github.com/Kong/kubernetes-ingress-controller/pull/5208)
- Fixed an issue that prevented the controller from associating admin API
  errors with a GRPCRoute.
  [#5267](https://github.com/Kong/kubernetes-ingress-controller/pull/5267)
  [#5275](https://github.com/Kong/kubernetes-ingress-controller/pull/5275)
- Restore the diagnostics server functionality, which was accidentally disabled.
  [#5270](https://github.com/Kong/kubernetes-ingress-controller/pull/5270)
- Allow configuring a GRPCRoute without hostnames and matches that catch all
  requests.
  [#5303](https://github.com/Kong/kubernetes-ingress-controller/pull/5303)
- Allow the `ws` and `wss` Enterprise protocol values for protocol annotations.
  [#5399](https://github.com/Kong/kubernetes-ingress-controller/pull/5399)
- Add a `workspace` parameter in filling IDs of Kong entities to avoid
  duplicate IDs cross different workspaces.
  [#5401](https://github.com/Kong/kubernetes-ingress-controller/pull/5401)
- Support properly ConsumerGroups when fallback to the last valid configuration.
  [#5438](https://github.com/Kong/kubernetes-ingress-controller/pull/5438)
- When specifying Gateway API Routes' `backendRef`s with namespace specified, those
  refs are checked for existence and allowed if they exist.
  [#5392](https://github.com/Kong/kubernetes-ingress-controller/pull/5392)
- Unmanaged Gateway mode honors the `--publish-status-address(-udp)` flags. If
  set, the controller will set these addresses in the Gateway status addresses
  instead of the proxy service/publish service addresses. The controller _no
  longer sets addresses in the Gateway spec addresses_. Review of the official
  specification indicated that the spec addresses are for user address
  requests, and that implementations can and should set status addresses to a
  different set of addresses if they assign addresses other than the requested
  set.
  [#5445](https://github.com/Kong/kubernetes-ingress-controller/pull/5445)
- Fixed a potential race condition that could occur when fetching the last applied
  or failed-to-be-applied config from the diagnostics server. The race could occur
  if the config was being updated while the HTTP endpoint was being hit at the same
  time.
  [#5474](https://github.com/Kong/kubernetes-ingress-controller/pull/5474)
- Stale `HTTPRoute`'s parent statuses are now removed when the `HTTPRoute` no longer
  defines a parent `Gateway` in its `spec.parentRefs`.
  [#5477](https://github.com/Kong/kubernetes-ingress-controller/pull/5477)
- `expressions` router flavor can now successfully be used with Konnect synchronization
  turned on. The controller will no longer populate disallowed `regex_priority` and `path_handling`
  Kong Route's fields when the router flavor is `expressions` that were causing Konnect to reject
  the configuration.
  [#5581](https://github.com/Kong/kubernetes-ingress-controller/pull/5581)

### Changed

- `SecretKeyRef` of `ConfigFrom` field in `KongPlugin` and `KongClusterPlugin`
  are `Required`. When `ConfigFrom` is specified, the validation of there CRDs
  will require `SecretKeyRef` to be present.
  [#5103](https://github.com/Kong/kubernetes-ingress-controller/pull/5103)
- CRD Validation Expressions
  - `KongPlugin` and `KongClusterPlugin` now enforce only one of `config` and `configFrom`
    to be set.
    [#5119](https://github.com/Kong/kubernetes-ingress-controller/pull/5119)
  - `KongConsumer` now enforces that at least one of `username` or `custom_id` is provided.
    [#5137](https://github.com/Kong/kubernetes-ingress-controller/pull/5137)
  - `KongPlugin` and `KongClusterPlugin` now enforce `plugin` to be immutable.
    [#5142](https://github.com/Kong/kubernetes-ingress-controller/pull/5142)
- `HTTPRoute` does no longer get rejected by the admission webhook when:
  - There's no `Gateway`'s `Listener` with `AllowedRoutes` matching the `HTTPRoute`.
  - There's no `Gateway`'s `Listener` with `Protocol` matching the `HTTPRoute`.
  - There's no `Gateway`'s `Listener` matching `HTTPRoute`'s `ParentRef`'s `SectionName`.
  All of these are validated by the controller and the results are reported in a `HTTPRoute`'s
  `Accepted` condition reported for a `Gateway`.
  [#5469](https://github.com/Kong/kubernetes-ingress-controller/pull/5469)

[kong-vault]: https://docs.konghq.com/gateway/latest/kong-enterprise/secrets-management/

## [3.0.2]

> Release date: 2024-01-11

### Added

- Added `--emit-kubernetes-events` CLI flag to disable the creation of events
  in translating and applying configurations to Kong.
  [#5296](https://github.com/Kong/kubernetes-ingress-controller/pull/5296)
  [#5299](https://github.com/Kong/kubernetes-ingress-controller/pull/5299)
- Added `-init-cache-sync-duration` CLI flag. This flag configures how long the controller waits for Kubernetes resources to populate at startup before generating the initial Kong configuration. It also fixes a bug that removed the default 5 second wait period.
  [#5238](https://github.com/Kong/kubernetes-ingress-controller/pull/5238)

## [3.0.1]

> Release date: 2023-11-22

### Fixed

- Using an Ingress with annotation `konghq.com/rewrite` and another Ingress without it pointing to the same Service,
  will no longer cause synchronization loop and random request failures due to incorrect routing.
  [#5218](https://github.com/Kong/kubernetes-ingress-controller/pull/5218)
- Using the same Service in one Ingress as a target for ingress rule and default backend works without issues.
  [#5219](https://github.com/Kong/kubernetes-ingress-controller/pull/5219)

## [3.0.0]

> Release date: 2023-11-03

### Highlights

- 🚀 Support for [Gateway API](https://kubernetes.io/docs/concepts/services-networking/gateways/) is now GA!
  - You only need to install Gateway API CRDs to use GA features of Gateway API with KIC.
  - Check the [Ingress to Gateway migration guide] to learn how to start using Gateway API already.
- 📈 Gateway Discovery feature is enabled by default both in DB-less and DB mode, allowing for scaling
  your gateways independently of the controller.
- 📖 Brand-new docs: [The KIC docs] have been totally revamped to be Gateway API first, and every single guide
  is as easy as copying and pasting your way down the page.

[Ingress to Gateway migration guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/migrate/ingress-to-gateway/
[The KIC docs]: https://docs.konghq.com/kubernetes-ingress-controller/latest/

### Breaking changes

- Only Kong Gateway in version >= 3.4.1 is supported. The controller will refuse to start
  if the version is lower, also won't discover such Kong Gateways.
  [#4766](https://github.com/Kong/kubernetes-ingress-controller/pull/4766)
- Removed feature gates: 
  - `CombinedServices`: The feature is enabled and it can't be changed.
    [#4743](https://github.com/Kong/kubernetes-ingress-controller/pull/4743)
  - `CombinedRoutes`: The feature is enabled and it can't be changed.
    [#4749](https://github.com/Kong/kubernetes-ingress-controller/pull/4749)
  - `ExpressionRoutes`: The feature is enabled and it can't be changed.
    KIC now translates to expression based Kong routes when Kong's router flavor `expressions`.
    [#4892](https://github.com/Kong/kubernetes-ingress-controller/pull/4892)
- Removed Knative support.
  [#4748](https://github.com/Kong/kubernetes-ingress-controller/pull/4748)
- The "text" logging format has changed. "json" should be used for
  machine-parseable logs.
  [#4688](https://github.com/Kong/kubernetes-ingress-controller/pull/4688)
- The "warn", "fatal", and "panic" log levels are no longer available. "error"
  is now the highest log level. "warn" logs are now logged at "error" level.
  [#4688](https://github.com/Kong/kubernetes-ingress-controller/pull/4688)
- Removed support for deprecated `KongIngress` fields: `Proxy` and `Route`. Respective
  `Service` or `Ingress` annotations should be used instead. See [KIC Annotations reference].
  [#4760](https://github.com/Kong/kubernetes-ingress-controller/pull/4760)
- Removed previously deprecated CLI flags:
  - `sync-rate-limit`
  - `stderrthreshold`
  - `update-status-on-shutdown`
  - `kong-custom-entities-secret`
  - `leader-elect`
  - `enable-controller-ingress-extensionsv1beta1`
  - `enable-controller-ingress-networkingv1beta1`
    [#4770](https://github.com/Kong/kubernetes-ingress-controller/pull/4770)
  - `debug-log-reduce-redundancy`
    [#4688](https://github.com/Kong/kubernetes-ingress-controller/pull/4688)
- `--konnect-runtime-group-id` CLI flag is now deprecated. Please use `--konnect-control-plane-id`
  instead.
  [#4783](https://github.com/Kong/kubernetes-ingress-controller/pull/4783)
- All manifests from `deploy/single` are no longer supported as installation
  method and were removed, please use Helm chart or Kong Gateway Operator instead.
  [#4866](https://github.com/Kong/kubernetes-ingress-controller/pull/4866)
  [#4873](https://github.com/Kong/kubernetes-ingress-controller/pull/4873)
  [#4970](https://github.com/Kong/kubernetes-ingress-controller/pull/4970)
- Credentials now use a `konghq.com/credential` label to indicate
  credential type instead of the `kongCredType` field. This allows controller
  compontents to avoid caching unnecessary Secrets. The `kongCredType` field is
  still supported but is now deprecated.
  See the [Migrate Credential Type Labels] guide to see how to update your `Secrets` smoothly.
- `KongIngress` is now entirely deprecated and will be removed in a future release.
  Its fields that were previously deprecated (`proxy` and `route`) are now not allowed to be set.
  They must be migrated to annotations. `upstream` field is deprecated - it's recommended
  to migrate its settings to the new `KongUpstreamPolicy` resource.
  See the [KongIngress to KongUpstreamPolicy migration guide] for details.
  [#5022](https://github.com/Kong/kubernetes-ingress-controller/pull/5022)
- Fixed `HTTPRoute` and `KongConsumer` admission webhook validators to properly
  signal validation failures, resulting in returning responses with `AdmissionResponse`
  filled instead of 500 status codes. It will make them work as expected in cases where
  the `ValidatingWebhookConfiguration` has `failurePolicy: Ignore`.
  This will enable validations of `HTTPRoute` and `KongConsumer` that were previously only
  accidentally effective with `failurePolicy: Fail`, thus it can be considered a breaking change.
  [#5063](https://github.com/Kong/kubernetes-ingress-controller/pull/5063)

### Fixed

- No more "log.SetLogger(...) was never called..." log entry during shutdown of KIC
  [#4738](https://github.com/Kong/kubernetes-ingress-controller/pull/4738)
- Changes to referenced Secrets are now tracked independent of their referent.
  [#4758](https://github.com/Kong/kubernetes-ingress-controller/pull/4758)
- When Kong returns a flattened error related to a Kong entity, the entity's type and name
  will be included in the message reported in `KongConfigurationApplyFailed` Kubernetes event
  generated for it.
  [#4813](https://github.com/Kong/kubernetes-ingress-controller/pull/4813)
- Fixed an incorrect watch, set in UDPRoute controller watching UDProute status updates.
  [#4835](https://github.com/Kong/kubernetes-ingress-controller/pull/4835)
- Fixed setting proper destination port for TCPRoute and UDPRoute, now field `SectionName`
  for `TCPRoute` and `UDPRoute` works as expected. It **breaks** some configurations that
  relied on matching multiple Gateway's listener ports to ports of services automatically.
  [#4928](https://github.com/Kong/kubernetes-ingress-controller/pull/4928)
- Fixed a panic when receiving broken configuration from Kong Gateway.
  [#5003](https://github.com/Kong/kubernetes-ingress-controller/pull/5003)
- Use 46 bits in values of priorities of generated Kong routes when expression
  rotuer is enabled to limit the priorities to be less than `1e14`. This
  prevents them to be encoded into scientific notation when dumping 
  configurations from admin API that brings precision loss and type 
  inconsistency in decoding JSON/YAML data to `uint64`. 
  This change will limit number of `HTTPRoute`s that can be 
  deterministically sorted by their creation timestamps, names and internal
  rule orders to `2^12=4096` and number of `GRPCRoutes` can be sorted to `2^8=256`.
  [#5024](https://github.com/Kong/kubernetes-ingress-controller/pull/5024)
- Error logs emitted from Gateway Discovery readiness checker that should be
  logged at `debug` level are now logged at that level.
  [#5029](https://github.com/Kong/kubernetes-ingress-controller/pull/5029)

### Changed

- Update paths of Konnect APIs from `runtime_groups/*` to `control-planes/*`.
  [#4566](https://github.com/Kong/kubernetes-ingress-controller/pull/4566)
- Docker images now use UID and GID 1000 to match Kong images. This should have
  no user-facing effect.
  [#4911](https://github.com/Kong/kubernetes-ingress-controller/pull/4911)
- Bump version of gateway API to `1.0.0` and support `Gateway`, `GatewayClass`
  and `HTTPRoute` in API version `gateway.networking.k8s.io/v1`.
  [#4893](https://github.com/Kong/kubernetes-ingress-controller/pull/4893)
  [#4981](https://github.com/Kong/kubernetes-ingress-controller/pull/4981)
  [#5041](https://github.com/Kong/kubernetes-ingress-controller/pull/5041)
- Update `Gateway`s, `GatewayClass`es and `HTTPRoute`s in examples to API
  version `gateway.networking.k8s.io/v1`.
  [#4935](https://github.com/Kong/kubernetes-ingress-controller/pull/4935)
- Controller to admin API communications are exempted from mesh proxies when
  the controller resides in a separate Deployment. This allows the controller
  to manage its own mTLS negotiation.
  [#4942](https://github.com/Kong/kubernetes-ingress-controller/pull/4942)
- Remove `Gateway` feature flag for Gateway API.
  [#4968](https://github.com/Kong/kubernetes-ingress-controller/pull/4968)

  It was enabled by default since 2.6.0 so the default behavior doesn't change.
  If users want to disable related functionality, they still can by disabling
  related Gateway API controllers via setting the following flags to `false`:
  - `--enable-controller-gwapi-gateway`
  - `--enable-controller-gwapi-httproute`
  - `--enable-controller-gwapi-reference-grant`
- Count `HTTPRoute` to gateway's number of attached route if the gateway is
  present in its `status.parents`, even if the gateway has unresolved refs.
  [#4987](https://github.com/Kong/kubernetes-ingress-controller/pull/4987)
- The default value for `--kong-admin-svc-port-names` is now `"admin-tls,kong-admin-tls"`
  instead of `"admin,admin-tls,kong-admin,kong-admin-tls"`. HTTP port names
  have been removed as discovery does not support plaintext HTTP connections.
  Instances configured with both HTTP and HTTPS admin ports resulted in
  discovery unsuccessfully trying to use HTTPS to talk to HTTP ports.
  [#5043](https://github.com/Kong/kubernetes-ingress-controller/pull/5043)
- The log format has been standardized to start with uppercase letters.
  [#5033](https://github.com/Kong/kubernetes-ingress-controller/pull/5033)
  [#5035](https://github.com/Kong/kubernetes-ingress-controller/pull/5035)
  [#5037](https://github.com/Kong/kubernetes-ingress-controller/pull/5037)
  [#5038](https://github.com/Kong/kubernetes-ingress-controller/pull/5038)
  [#5049](https://github.com/Kong/kubernetes-ingress-controller/pull/5049)
  [#5050](https://github.com/Kong/kubernetes-ingress-controller/pull/5050)
  [#5065](https://github.com/Kong/kubernetes-ingress-controller/pull/5065)

### Added

- Added support for expression-based Kong routes for `TLSRoute`. This requires
  Kong installed with `KONG_ROUTER_FLAVOR=expressions` set in the environment.
  [#4574](https://github.com/Kong/kubernetes-ingress-controller/pull/4574).
- The `FillIDs` feature gate is now enabled by default.
  [#4746](https://github.com/Kong/kubernetes-ingress-controller/pull/4746)
- Get rid of deprecation warning in logs for unsupported label `global: true` for `KongPlugin`,
  it'll be treated as any other label without a special meaning.
  [#4737](https://github.com/Kong/kubernetes-ingress-controller/pull/4737)
- Telemetry now reports the router flavor.
  [#4762](https://github.com/Kong/kubernetes-ingress-controller/pull/4762)
- Support Query Parameter matching of `HTTPRoute` when expression router enabled.
  [#4780](https://github.com/Kong/kubernetes-ingress-controller/pull/4780)
- Support `ExtensionRef` HTTPRoute filter. It is now possibile to set a KongPlugin
  reference in the `HTTPRoute`s' `ExtensionRef` filter field.
  [#4838](https://github.com/Kong/kubernetes-ingress-controller/pull/4838)
- Added `--kong-admin-token-file` flag to provide the Kong admin token via a
  file. This is an alternative to the existing `--kong-admin-token` for users
  that prefer to mount a file over binding a Secret to an environment variable
  value. Only one of the two options can be used.
  [#4808](https://github.com/Kong/kubernetes-ingress-controller/pull/4808)
- New `KongUpstreamPolicy` CRD superseding `KongIngress.Upstream` was introduced.
  It allows overriding Kong Upstream settings generated for a specific `Service` used
  in an `Ingress` or Gateway API `Route`. A policy can be applied to a `Service` by
  setting `konghq.com/upstream-policy: <policy-name>` annotation on the `Service` object.
  Read more in [KIC CRDs reference].
  [#4880](https://github.com/Kong/kubernetes-ingress-controller/pull/4880)
  [#4943](https://github.com/Kong/kubernetes-ingress-controller/pull/4943)
  [#4955](https://github.com/Kong/kubernetes-ingress-controller/pull/4955)
  [#4957](https://github.com/Kong/kubernetes-ingress-controller/pull/4957)
  [#4969](https://github.com/Kong/kubernetes-ingress-controller/pull/4969)
  [#4979](https://github.com/Kong/kubernetes-ingress-controller/pull/4979)
  [#4989](https://github.com/Kong/kubernetes-ingress-controller/pull/4989)
- KIC now specifies its `UserAgent` when communicating with kube-apiserver
  as `kong-ingress-controller/${VERSION}` where version is the version of KIC.
  [#5019](https://github.com/Kong/kubernetes-ingress-controller/pull/5019)
- Allow Gateway Discovery with database backed Kong. KIC will send Kong 
  configurations to one of the backend pods of the service specified by the
  flag `--kong-admin-svc` if Kong gateway is DB backed.
  [#4828](https://github.com/Kong/kubernetes-ingress-controller/pull/4828)

[KIC Annotations reference]: https://docs.konghq.com/kubernetes-ingress-controller/latest/references/annotations/
[KIC CRDs reference]: https://docs.konghq.com/kubernetes-ingress-controller/latest/references/custom-resources/
[KongIngress to KongUpstreamPolicy migration guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/migrate/kongingress/
[Migrate Credential Type Labels]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/migrate/credential-kongcredtype-label/

## [2.12.5]

> Release date: 2024-06-25

### Fixed

- Services using `Secret`s containing the same certificate as client certificates
  by annotation `konghq.com/client-cert` can be correctly translated.
  [#6228](https://github.com/Kong/kubernetes-ingress-controller/pull/6228)

## [2.12.4]

> Release date: 2024-04-30

- Bump golang version to 1.21.9 to fix CVE [GO-2024-2687](https://pkg.go.dev/vuln/GO-2024-2687).
  [#5916](https://github.com/Kong/kubernetes-ingress-controller/pull/5916)
- Bump `golang.org/x/net` to `0.23.0` and `google.golang.org/protobuf` to `1.33.0`
  To fix [GO-2024-2687](https://pkg.go.dev/vuln/GO-2024-2687) and [GO-2024-2611](https://pkg.go.dev/vuln/GO-2024-2611).
  [#5947](https://github.com/Kong/kubernetes-ingress-controller/pull/5947)

## [2.12.3]

> Release date: 2023-12-19

### Fixed

- Fix(manager): set InitCacheSyncDuration to 5s by default and allow it to be configured via `--init-cache-sync-duration` CLI flag
  [#5238](https://github.com/Kong/kubernetes-ingress-controller/pull/5238)
- Don't set `instance_name` of plugin if Kong version is below 3.2.0.
  [#5250](https://github.com/Kong/kubernetes-ingress-controller/pull/5250)
- Added `--emit-kubernetes-events` CLI flag to disable the creation of events
  in translating and applying configurations to Kong.
  [#5296](https://github.com/Kong/kubernetes-ingress-controller/pull/5296)
  [#5299](https://github.com/Kong/kubernetes-ingress-controller/pull/5299)

## [2.12.2]

> Release date: 2023-11-22

### Fixed

- Using an Ingress with annotation `konghq.com/rewrite` and another Ingress without it pointing to the same Service,
  will no longer cause synchronization loop and random request failures due to incorrect routing.
  [#5215](https://github.com/Kong/kubernetes-ingress-controller/pull/5215)
- Using the same Service in one Ingress as a target for ingress rule and default backend works without issues.
  [#5217](https://github.com/Kong/kubernetes-ingress-controller/pull/5217)

### Known issues

- **Only when combined routes are not enabled**, generated Kong routes may have conflicting names, that leads to
  incorrect routing. In such case the descriptive error message is now provided. Use feature gate `CombinedRoutes=true`
  or update Kong Kubernetes Ingress Controller version to 3.0.0 or above (both remediation changes naming schema of Kong routes).
  [#5198](https://github.com/Kong/kubernetes-ingress-controller/issues/5198)

## [2.12.1]

> Release date: 2023-11-09

### Fixed

- Credentials Secrets that are not referenced by any KongConsumer but violate the KongConsumer
  basic level validation (invalid credential type or missing required fields) are now rejected
  by the admission webhook.
  [#4887](https://github.com/Kong/kubernetes-ingress-controller/pull/4887)
- Error logs emitted from Gateway Discovery readiness checker that should be
  logged at `debug` level are now logged at that level.
  [#5030](https://github.com/Kong/kubernetes-ingress-controller/pull/5030)
- Fix `panic` when last known configuration fetcher gets a `nil` Status when requesting
  `/status` from Kong Gateway.
  This happens when Gateway is responding with a 50x HTTP status code.
  [#5120](https://github.com/Kong/kubernetes-ingress-controller/pull/5120)
- Use 46 bits in values of priorities of generated Kong routes when expression
  rotuer is enabled to limit the priorities to be less than `1e14`. This
  prevents them to be encoded into scientific notation when dumping
  configurations from admin API that brings precision loss and type
  inconsistency in decoding JSON/YAML data to `uint64`.
  This change will limit number of `HTTPRoute`s that can be
  deterministically sorted by their creation timestamps, names and internal
  rule orders to `2^12=4096` and number of `GRPCRoutes` can be sorted to `2^8=256`.
  [#5124](https://github.com/Kong/kubernetes-ingress-controller/pull/5124)

## [2.12.0]

> Release date: 2023-09-25

### Deprecated

- Knative Ingress is deprecated and will be removed in KIC 3.0. [#2813](https://github.com/Kong/kubernetes-ingress-controller/issues/2813)
- `KongIngress` for `Service` and `Route` parameters has been deprecated since KIC 2.8 and will be removed in KIC 3.0.
    - We expect to eventually deprecate `KongIngress` also for `Upstream` parameters as described in [#3174](https://github.com/Kong/kubernetes-ingress-controller/issues/3174)
- Existing Kustomize (`deploy/manifests/`) and `deploy/single/` YAML manifests as a method of installing KIC.
    - The `deploy/single/` and `deploy/manifests/` directories will no longer work with KIC 3.0+. You should use the [Helm chart](https://docs.konghq.com/kubernetes-ingress-controller/latest/deployment/k4k8s/#helm) or [Kong Gateway Operator](https://docs.konghq.com/gateway-operator/latest/) instead.
- DB-less deployments of Kong running with KIC as a sidecar. The [Gateway Discovery](https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/using-gateway-discovery/) feature added in KIC 2.9 should be used instead.
    - The mode where Kong runs with a database (Postgres) is not affected by the migration to Gateway Discovery yet, but likely will in the future [#4751](https://github.com/Kong/kubernetes-ingress-controller/issues/4751)

### Added

- `konghq.com/rewrite` annotation has been introduced to manage URI rewriting.
  This feature requires enabling the `RewriteURIs` feature gate.
  [#4360](https://github.com/Kong/kubernetes-ingress-controller/pull/4360), [#4646](https://github.com/Kong/kubernetes-ingress-controller/pull/4646)
- Provide validation in admission webhook for `Ingress` paths (validate regex expressions).
  [#4647](https://github.com/Kong/kubernetes-ingress-controller/pull/4647)
  [#4360](https://github.com/Kong/kubernetes-ingress-controller/pull/4360)
- Added support for expression-based Kong routes for `TCPRoute`, `UDPRoute`,
  `TCPIngress`, and `UDPIngress`. This requires the `ExpressionRoutes` feature
  gate and a Kong 3.4+ install with `KONG_ROUTER_FLAVOR=expressions` set in the
  environment.
  [#4385](https://github.com/Kong/kubernetes-ingress-controller/pull/4385)
  [#4550](https://github.com/Kong/kubernetes-ingress-controller/pull/4550)
  [#4612](https://github.com/Kong/kubernetes-ingress-controller/pull/4612)
- `KongIngress` CRD now supports `latency` algorithm in its `upstream.algorithm`
  field. This can be used with Kong Gateway 3.2+.
  [#4703](https://github.com/Kong/kubernetes-ingress-controller/pull/4703)

### Changed

- Generate wildcard routes to match all `HTTP` or `GRPC` requests for rules
  in `HTTPRoute` or `GRPCRoute` if there are no matches in the rule and no
  hostnames in their parent objects.
  [#4526](https://github.com/Kong/kubernetes-ingress-controller/pull/4528)
- The Gateway API has been bumped to 0.8.1.
  [#4700](https://github.com/Kong/kubernetes-ingress-controller/pull/4700)

### Fixed

- Allow regex expressions in `HTTPRoute` configuration and provide validation in admission webhook.
  Before this change admission webhook used to reject entirely such configurations incorrectly as not supported yet.
  [#4608](https://github.com/Kong/kubernetes-ingress-controller/pull/4608)
- Do not parse error body when failed to get response from reloading declarative
  configurations to produce proper error log in such situations,
  [#4666](https://github.com/Kong/kubernetes-ingress-controller/pull/4666)
- Set type meta of objects when adding them to caches and reference indexers
  to ensure that indexes of objects in reference indexers have correct object
  kind. This ensures referece relations of objects are stored and indexed
  correctly.
  [#4663](https://github.com/Kong/kubernetes-ingress-controller/pull/4663)
- Display Service ports on generated Kong services, instead of a static default
  value. This change is cosmetic only.
  [#4503](https://github.com/Kong/kubernetes-ingress-controller/pull/4503)
- Create routes that match any service and method for `GRPCRoute` rules with no
  matches.
  [#4512](https://github.com/Kong/kubernetes-ingress-controller/issues/4512)
- KongPlugins used on multiple resources will no longer result in
  `instance_name` collisions.
  [#4588](https://github.com/Kong/kubernetes-ingress-controller/issues/4588)
- Fix `panic` when last known configuration fetcher gets a `nil` Status when requesting
  `/status` from Kong Gateway.
  This happens when Gateway is responding with a 50x HTTP status code.
  [#4627](https://github.com/Kong/kubernetes-ingress-controller/issues/4627)
- Ensure the API server is available at startup and do not disable CRD
  controllers if Kong CRDs are unavailable. Do not disable the Ingress
  controller if the Ingress API is unavailable. This avoids incorrectly
  deleting existing configuration during an API server restart.
  [#4641](https://github.com/Kong/kubernetes-ingress-controller/issues/4641)
  [#4643](https://github.com/Kong/kubernetes-ingress-controller/issues/4643)
- Fix `Licenses` and `ConsumerGroups` missing in sanitized copies of Kong configuration.
  [#4710](https://github.com/Kong/kubernetes-ingress-controller/pull/4710)

## [2.11.1]

> Release date: 2023-08-29

### Changed

- Bumped the default Kong version to 3.4 in example manifests.
  [#4534](https://github.com/Kong/kubernetes-ingress-controller/pull/4534)

### Fixed

- Disable KongPlugin and KongClusterPlugin Programmed statuses. These were
  introduced in 2.11.0 and caused excessively frequent status updates in
  clusters with multiple KIC instances installed.
  [#4584](https://github.com/Kong/kubernetes-ingress-controller/pull/4584)

## [2.11.0]

> Release date: 2023-08-09

### Added

- Introduce `KongConsumerGroup` CRD (supported by Kong Enterprise only)
  [#4325](https://github.com/Kong/kubernetes-ingress-controller/pull/4325)
  [#4387](https://github.com/Kong/kubernetes-ingress-controller/pull/4387)
  [#4419](https://github.com/Kong/kubernetes-ingress-controller/pull/4419)
  [#4437](https://github.com/Kong/kubernetes-ingress-controller/pull/4437)
  [#4452](https://github.com/Kong/kubernetes-ingress-controller/pull/4452)
- The ResponseHeaderModifier Gateway API filter is now supported and translated
  to the proper set of Kong plugins.
  [#4350](https://github.com/Kong/kubernetes-ingress-controller/pull/4350)
- The `CombinedServices` feature gate is now enabled by default.
  [#4138](https://github.com/Kong/kubernetes-ingress-controller/pull/4138)
- Plugin CRDs now support the `instance_name` field introduced in Kong 3.2.
  [#4174](https://github.com/Kong/kubernetes-ingress-controller/pull/4174)
- Gateway resources no longer use the _Ready_ condition following changes to
  the upstream Gateway API specification in version 0.7.
  [#4142](https://github.com/Kong/kubernetes-ingress-controller/pull/4142)
- Prometheus metrics now include counts of resources that the controller cannot
  send to the proxy instances and the last successful configuration push time.
  [#4181](https://github.com/Kong/kubernetes-ingress-controller/pull/4181)
- Store the last known good configuration. If Kong rejects the latest
  configuration, send the last good configuration to Kong instances with no
  configuration. This allows newly-started Kong instances to serve traffic even
  if a configuration error prevents the controller from sending the latest
  configuration.
  [#4205](https://github.com/Kong/kubernetes-ingress-controller/pull/4205)
- Telemetry reports now include the OpenShift version, if any.
  [#4211](https://github.com/Kong/kubernetes-ingress-controller/pull/4211)
- Assign priorities to routes translated from Ingresses when parser translate
  them to expression based Kong routes. The assigning method is basically the
  same as in Kong gateway's `traditional_compatible` router, except that
  `regex_priority` field in Kong traditional route is not supported. This
  method is adopted to keep the compatibility with traditional router on
  maximum effort.
  [#4240](https://github.com/Kong/kubernetes-ingress-controller/pull/4240)
- Assign priorities to routes translated from HTTPRoutes when parser translates
  them to expression based Kong routes. The assigning method follows the
  [specification on priorities of matches in `HTTPRoute`][httproute-specification].
  [#4296](https://github.com/Kong/kubernetes-ingress-controller/pull/4296)
  [#4434](https://github.com/Kong/kubernetes-ingress-controller/pull/4434)
- Assign priorities to routes translated from GRPCRoutes when the parser translates
  them to expression based Kong routes. The priority order follows the
  [specification on match priorities in GRPCRoute][grpcroute-specification].
  [#4364](https://github.com/Kong/kubernetes-ingress-controller/pull/4364)
- When a translated Kong configuration is empty in DB-less mode, the controller
  will now send the configuration with a single empty `Upstream`. This is to make
  Gateways using `/status/ready` as their health check ready after receiving the
  initial configuration (even if it's empty).
  [#4316](https://github.com/Kong/kubernetes-ingress-controller/pull/4316)
- Fetch the last known good configuration from existing proxy instances. If
  KIC restarts, it is now able to fetch the last good configuration from a running
  proxy instance and store it in its internal cache.
  [#4265](https://github.com/Kong/kubernetes-ingress-controller/pull/4265)
- Gateway Discovery feature was adapted to handle Gateways that are not ready yet
  in terms of accepting data-plane traffic, but are ready to accept configuration
  updates. The controller will now send configuration to such Gateways and will
  actively monitor their readiness for accepting configuration updates.
  [#4368](https://github.com/Kong/kubernetes-ingress-controller/pull/4368)
- `KongConsumer`, `KongConsumerGroup` `KongPlugin`, and `KongClusterPlugin` CRDs were extended with
  `Status.Conditions` field. It will contain the `Programmed` condition describing
  whether an object was successfully translated into Kong entities and sent to Kong.
  [#4409](https://github.com/Kong/kubernetes-ingress-controller/pull/4409)
  [#4412](https://github.com/Kong/kubernetes-ingress-controller/pull/4412)
  [#4423](https://github.com/Kong/kubernetes-ingress-controller/pull/4423)
- `KongConsumer`, `KongConsumerGroup`, `KongPlugin`, and `KongClusterPlugin`'s `additionalPrinterColumns`
  were extended with `Programmed` column. It will display the status of the
  `Programmed` condition of an object when `kubectl get` is used.
  [#4425](https://github.com/Kong/kubernetes-ingress-controller/pull/4425)
  [#4423](https://github.com/Kong/kubernetes-ingress-controller/pull/4423)
- Parser instead of logging errors for invalid `KongPlugin` or `KongClusterPlugin`
  configuration, will now propagate a translation failure that will result
  in the `Programmed` condition of the object being set to `False` and an
  event being emitted.
  [#4428](https://github.com/Kong/kubernetes-ingress-controller/pull/4428)

### Changed

- Log message `no active endpoints` is now logged at debug instead of
  warning level.
  [#4161](https://github.com/Kong/kubernetes-ingress-controller/pull/4161)
- Events and logs for inconsistent multi-Service backend annotations now list
  all involved Services, not just Services whose annotation does not match the
  first observed value, as that value is not necessarily the desired value.
  [#4171](https://github.com/Kong/kubernetes-ingress-controller/pull/4171)
- Use [`gojson`][gojson] for marshalling JSON when generating SHA for config.
  This should yield some performance benefits during config preparation and
  sending stage (we've observed around 35% reduced time in config marshalling
  time but be aware that your mileage may vary).
  [#4222](https://github.com/Kong/kubernetes-ingress-controller/pull/4222)
- Changed the Gateway's readiness probe in all-in-one manifests from `/status`
  to `/status/ready`. Gateways will be considered ready only after an initial
  configuration is applied by the controller.
  [#4368](https://github.com/Kong/kubernetes-ingress-controller/pull/4368)
- When translating to expression based Kong routes, annotations to specify
  protocols are translated to `protocols` field of the result Kong route,
  instead of putting the conditions to match protocols inside expressions.
  [#4422](https://github.com/Kong/kubernetes-ingress-controller/pull/4422)

### Fixed

- Correctly support multi-Service backends that have multiple Services sharing
  the same name in different namespaces.
  [#4375](https://github.com/Kong/kubernetes-ingress-controller/pull/4375)
- Properly construct targets for IPv6-only clusters.
  [#4391](https://github.com/Kong/kubernetes-ingress-controller/pull/4391)
- Attach kubernetes events to `KongConsumer`s when the parser fails to
  translate its credentials to Kong configuration, instead of logging thet
  error to reduce the redundant logs.
  [#4398](https://github.com/Kong/kubernetes-ingress-controller/pull/4398)
- `Gateway` can now correctly update `AttachedRoutes` even if there are more
  than 100 `HttpRoute`s.
  [#4458](https://github.com/Kong/kubernetes-ingress-controller/pull/4458)

[gojson]: https://github.com/goccy/go-json
[httproute-specification]: https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io/v1.HTTPRoute
[grpcroute-specification]:  https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io/v1alpha2.GRPCRouteRule

## [2.10.5]

> Release date: 2023-08-31

### Fixed

- Fixed compatibility with Kong 3.3+ when using chart versions 2.26+.
  [#4515](https://github.com/Kong/kubernetes-ingress-controller/pull/4515)

## [2.10.4]

> Release date: 2023-07-25

### Fixed

- Fixed a bug that prevented the controller from updating configuration in
  Konnect Runtime Groups API when an existing Certificate was updated.
  [#4384](https://github.com/Kong/kubernetes-ingress-controller/issues/4384)

## [2.10.3]

> Release date: 2023-07-13

### Fixed

- Nodes in Konnect Runtime Groups API are not updated every 3s anymore.
  This was caused by a bug in `NodeAgent` that was sending the updates
  despite the fact that the configuration status was not changed.
  [#4324](https://github.com/Kong/kubernetes-ingress-controller/pull/4324)

## [2.10.2]

> Release date: 2023-07-07

### Added

- Added `--update-status-queue-buffer-size` allowing configuring the size of
  the status update queue's underlying channels used to buffer updates to the
  status of Kubernetes resources.
  [#4267](https://github.com/Kong/kubernetes-ingress-controller/pull/4267)

### Fixed

- Translator of `GRPCRoute` generates paths without leading `~` when running
  with Kong gateway with version below 3.0.
  [#4238](https://github.com/Kong/kubernetes-ingress-controller/pull/4238)
- Fixed a bug where the controller sync loop would get stuck when a number of
  updates for one of Gateway API resources kinds (`HTTPRoute`, `TCPRoute`,
  `UDPRoute`, `TLSRoute`, `GRPCRoute`) exceeded 8192. This was caused by the
  fact that the controller was using a fixed-size buffer to store updates for
  each resource kind and there were no consumers for the updates. The sending
  was blocked after a buffer got full, resulting in a deadlock.
  [#4267](https://github.com/Kong/kubernetes-ingress-controller/pull/4267)

## [2.10.1]

> Release date: 2023-06-27

### Added

- `--konnect-initial-license-polling-period` and `--konnect-license-polling-period`
  CLI flags were added to allow configuring periods at which KIC polls license
  from Konnect. The initial period will be used until a valid license is retrieved.
  The default values are 1m and 12h respectively.
  [#4178](https://github.com/Kong/kubernetes-ingress-controller/pull/4178)

### Fixed

- Fix KIC crash which occurred when invalid config was applied in DB mode.
  [#4213](https://github.com/Kong/kubernetes-ingress-controller/pull/4213)

## [2.10.0]

> Release date: 2023-06-02

### Added

- Gateways now track UDP Listener status when `--publish-service-udp` is set.
  `UDPRoute`s that do not match a valid UDP Listener are excluded from
  configuration. Previously KIC added any UDPRoute that indicated an associated
  Gateway as its parent regardless of Listener configuration or status.
  [#3832](https://github.com/Kong/kubernetes-ingress-controller/pull/3832)
- Added license agent for Konnect-managed instances.
  [#3883](https://github.com/Kong/kubernetes-ingress-controller/pull/3883)
- `Service`, `Route`, and `Consumer` Kong entities now can get assigned
  deterministic IDs based on their unique properties (name, username, etc.)
  instead of random UUIDs. To enable this feature, set `FillIDs` feature gate
  to `true`.
  It's going to be useful in cases where stable IDs are needed across multiple
  Kong Gateways managed by KIC (e.g. for the integration with Konnect and
  reporting metrics that later can be aggregated across multiple instances based
  on the entity's ID).
  When `FillIDs` will be enabled, the controller will re-create all the existing
  entities (Services, Routes, and Consumers) with the new IDs assigned. That can
  potentially lead to temporary downtime between the deletion of the old entities
  and the creation of the new ones.
  Users should be cautious about enabling the feature if their existing DB-backed
  setup consists of a huge amount of entities for which the recreation can take
  significant time.
  [#3933](https://github.com/Kong/kubernetes-ingress-controller/pull/3933)
  [#4075](https://github.com/Kong/kubernetes-ingress-controller/pull/4075)
- Added translator to translate ingresses under `networking.k8s.io/v1` to
  expression based Kong routes. The translator is enabled when feature gate
  `ExpressionRoutes` is turned on and the managed Kong gateway runs in router
  flavor `expressions`. We assume `router_flavor` to be `traditional`
  for versions below 3.0. If it is not available in Kong's configuration,
  for versions 3.0 and above, abort with an error.
  [#3935](https://github.com/Kong/kubernetes-ingress-controller/pull/3935)
  [#4076](https://github.com/Kong/kubernetes-ingress-controller/pull/4076)
- Added `CombinedServices` feature gate that prevents the controller from
  creating a separate Kong `Service` for each `netv1.Ingress` that uses
  the same Kubernetes `Service` as its backend when `CombinedRoutes` feature
  is turned on. Instead, the controller will create a single Kong `Service`
  in such circumstances.
  The feature is disabled by default to allow users to adapt to the changed
  Kong Services naming scheme (`<namespace>.<service-name>.<port>` instead
  of `<namespace>.<ingress-name>.<service-name>.<port>`) - it could break
  existing monitoring rules that rely on the old naming scheme.
  It will become the default behavior in the next minor release with the possibility
  to opt-out.
  [#3963](https://github.com/Kong/kubernetes-ingress-controller/pull/3963)
- Added translator to translate `HTTPRoute` and `GRPCRoute` in gateway APIs to
  expression based kong routes. Similar to ingresses, this translator is only
  enabled when feature gate `ExpressionRoutes` is turned on and the managed
  Kong gateway runs in router flavor `expressions`.
  [#3956](https://github.com/Kong/kubernetes-ingress-controller/pull/3956)
  [#3988](https://github.com/Kong/kubernetes-ingress-controller/pull/3988)
- Configuration updates to Konnect Runtime Group's Admin API now respect a backoff
  strategy that prevents KIC from exceeding API calls limits.
  [#3989](https://github.com/Kong/kubernetes-ingress-controller/pull/3989)
  [#4015](https://github.com/Kong/kubernetes-ingress-controller/pull/4015)
- When Gateway API CRDs are not installed, the controllers of those are not started
  during the start-up phase. From now on, they will be dynamically started in runtime
  once their installation is detected, making restarting the process unnecessary.
  [#3996](https://github.com/Kong/kubernetes-ingress-controller/pull/3996)
- Disable translation of unsupported Kubernetes objects when translating to
  expression based routes is enabled (`ExpressionRoutes` feature enabled AND
  kong using router flavor `expressions`), and generate a translation failure
  event attached to each of the unsupported objects.
  [#4022](https://github.com/Kong/kubernetes-ingress-controller/pull/4022)
- Controller's configuration synchronization status reported to Konnect's Node API
  now accounts for potential failures in synchronizing configuration with Konnect's
  Runtime Group Admin API.
  [#4029](https://github.com/Kong/kubernetes-ingress-controller/pull/4029)
- Record an event attached to KIC pod after applying configuration to Kong. If
  the applying succeeded, a `Normal` event with `KongConfigurationSucceeded`
  reason is recorded. If the applying failed, a `Warning` event with
  `KongConfigurationApplyFailed` reason is recorded.
  [#4054](https://github.com/Kong/kubernetes-ingress-controller/pull/4054)
- Disable translation to expression routes when feature gate `ExpressionRoutes`
  is enabled but feature gate `CombinedRoutes` is not enabled.
  [#4057](https://github.com/Kong/kubernetes-ingress-controller/pull/4057)
- Added `--gateway-discovery-dns-strategy` flag which allows specifying which
  DNS strategy to use when generating Gateway's Admin API addresses.
  [#4071](https://github.com/Kong/kubernetes-ingress-controller/pull/4071)

  There are 3 options available
  - `ip` (default): which will make KIC create Admin API addresses built out of
    IP addresses.
  - `pod`: will make KIC build addresses using the following template:
    `pod-ip-address.my-namespace.pod`.
  - `service`: will make KIC build addresses using the following template:
    `pod-ip-address.service-name.my-namespace.svc`.
    This is known to not work on GKE because it uses `kube-dns` instead of `coredns`.
- Gateway's `AttachedRoutes` fields get updated with the number of routes referencing
  and using each listener.
  [#4052](https://github.com/Kong/kubernetes-ingress-controller/pull/4052)
- `all-in-one-postgres.yaml` and `all-in-one-postgres-enterprise.yaml` manifests'
  migrations job now works properly when running against an already bootstrapped
  database, allowing upgrades from one version of Kong Gateway to another without
  tearing down the database.
  [#4116](https://github.com/Kong/kubernetes-ingress-controller/pull/4116)
- Telemetry reports now include a count for every `gateway.networking.k8s.io` CRD.
  [#4058](https://github.com/Kong/kubernetes-ingress-controller/pull/4058)

### Changed

- Kong Ingress Controller no longer relies on `k8s.io.api.core.v1` `Endpoints`,
  and instead uses `discovery.k8s.io/v1` `EndpointSlice` to discover endpoints
  for Kubernetes `Service`s.
  [#3997](https://github.com/Kong/kubernetes-ingress-controller/pull/3997)
- Gateway Discovery now produces DNS names instead of IP addresses
  [#4044](https://github.com/Kong/kubernetes-ingress-controller/pull/4044)

### Fixed

- Fixed paging in `GetAdminAPIsForService` which might have caused the controller
  to only return the head of the list of Endpoints for Admin API service.
  [#3846](https://github.com/Kong/kubernetes-ingress-controller/pull/3846)
- Fixed a race condition in the version-specific feature system.
  [#3852](https://github.com/Kong/kubernetes-ingress-controller/pull/3852)
- Fixed a missing reconciliation behavior for Admin API `EndpointSlice` reconciler
  when the `EndpointSlice` that we receive a reconciliation request for is already
  missing
  [#3889](https://github.com/Kong/kubernetes-ingress-controller/pull/3889)
- Fixed leader election role manifest where `""` and `"coordination"` API groups
  together with the related manifest resources (`configmaps` and `leases`) might
  become mixed up when the manifest is unmarshalled.
  [#3932](https://github.com/Kong/kubernetes-ingress-controller/pull/3932)

### Deprecated

- Removed support for `extensions/v1beta1` `Ingress` which was removed in kubernetes 1.22.
  At the same time deprecate `--enable-controller-ingress-extensionsv1beta1` CLI flag.
  [#3710](https://github.com/Kong/kubernetes-ingress-controller/pull/3710)

- Removed support for `networking.k8s.io/v1beta1` `Ingress` which was removed in kubernetes 1.22.
  At the same time deprecate `--enable-controller-ingress-networkingv1beta1` CLI flag.
  [#3867](https://github.com/Kong/kubernetes-ingress-controller/pull/3867)

## [2.9.3]

> Release date: 2023-04-17

### Fixed

- Fixed a missing reconciliation behavior for Admin API EndpointSlice reconciler
  when the EndpointSlice that we receive a reconciliation request for is already
  missing
  [#3889](https://github.com/Kong/kubernetes-ingress-controller/pull/3889)
- Update enterprise manifests to use Kong Gateway 3.2
  [#3885](https://github.com/Kong/kubernetes-ingress-controller/pull/3885)

## [2.9.2]

> Release date: 2023-04-03

### Fixed

- Fixed a deadlock in `AdminAPIClientsManager` which could occur when Konnect integration
  was enabled, and multiple `Notify` calls were made in parallel (e.g. when scaling Gateway
  deployment up).
  [#3816](https://github.com/Kong/kubernetes-ingress-controller/pull/3816)

## [2.9.1]

> Release date: 2023-03-29

This release was intended to include a fix for a deadlock in `AdminAPIClientsManager`
([#3816](https://github.com/Kong/kubernetes-ingress-controller/pull/3816)), but it wasn't
backported properly. It is included in the next patch release.

## [2.9.0]

> Release date: 2023-03-27

### Added

- Konnect Runtime Group's nodes are reactively updated on each discovered Gateway clients
  change.
  [#3727](https://github.com/Kong/kubernetes-ingress-controller/pull/3727)
- Telemetry reports now include a number of discovered Gateways when the Gateway Discovery
  feature is turned on.
  [#3783](https://github.com/Kong/kubernetes-ingress-controller/pull/3783)
- Adding the `konghq.com/tags: csv,of,tags` annotation will add tags to
  generated resources.
  [#3778](https://github.com/Kong/kubernetes-ingress-controller/pull/3778)
- `HTTPRoute` reconciler now watches relevant `ReferenceGrant`s for changes.
  [#3759](https://github.com/Kong/kubernetes-ingress-controller/pull/3759)
- Bumped Kong version in manifests to 3.2.
  [#3804](https://github.com/Kong/kubernetes-ingress-controller/pull/3804)
- Store status of whether configuration succeeded or failed for Kubernetes
  objects in dataplane client and publish the events to let controllers know
  if the controlled objects succeeded or failed to be translated to Kong
  configuration.
  [#3359](https://github.com/Kong/kubernetes-ingress-controller/pull/3359)
- Added `version` command
  [#3379](https://github.com/Kong/kubernetes-ingress-controller/pull/3379)  
- Added `--publish-service-udp` to indicate the Service that handles inbound
  UDP traffic.
  [#3325](https://github.com/Kong/kubernetes-ingress-controller/pull/3325)
- Added possibility to configure multiple Kong Gateways through the
  `--kong-admin-url` CLI flag (which can be specified multiple times) or through
  a corresponding environment variable `CONTROLLER_KONG_ADMIN_URL` (which can
  specify multiple values separated by a comma).
  [#3268](https://github.com/Kong/kubernetes-ingress-controller/pull/3268)
- Added a new `dbless-konnect` configuration variant to the manifests. It can
  be used to deploy a DB-less variant of KIC that will also synchronise its
  data-plane configuration with Konnect cloud.
  [#3448](https://github.com/Kong/kubernetes-ingress-controller/pull/3448)
- The Gateway API has been bumped to 0.6.1. The `GatewayConditionScheduled` has
  been replaced by the `GatewayConditionAccepted`, and the `ListenerConditionDetached`
  condition has been replaced by the `ListenerConditionAccepted`.
  [#3496](https://github.com/Kong/kubernetes-ingress-controller/pull/3496)
  [#3524](https://github.com/Kong/kubernetes-ingress-controller/pull/3524)
- The `ReferenceGrant` has been promoted to beta.
  [#3507](https://github.com/Kong/kubernetes-ingress-controller/pull/3507)
- Enable `ReferenceGrant` if `Gateway` feature gate is turned on (default).
  [#3519](https://github.com/Kong/kubernetes-ingress-controller/pull/3519)
- Experimental `--konnect-sync-enabled` feature flag has been introduced. It
  enables the integration with Kong's Konnect cloud. It's turned off by default.
  When enabled, it allows to synchronise data-plane configuration with
  a Konnect Runtime Group specified by `--konnect-runtime-group-id`.
  It requires `--konnect-tls-client-*` set of flags to be set to provide
  Runtime Group's TLS client certificates for authentication.
  [#3455](https://github.com/Kong/kubernetes-ingress-controller/pull/3455)
- Added Konnect client to upload status of KIC instance to Konnect cloud if
  flag `--konnect-sync-enabled` is set to `true`.
  [#3469](https://github.com/Kong/kubernetes-ingress-controller/pull/3469)
- Added Gateway discovery using Kong Admin API service configured via `--kong-admin-svc`
  which accepts a namespaced name of a headless service which should have
  Admin API endpoints exposed under a named port called `admin`. Gateway 
  discovery is only allowed to run with dbless kong gateways.
  [#3421](https://github.com/Kong/kubernetes-ingress-controller/pull/3421)
  [#3642](https://github.com/Kong/kubernetes-ingress-controller/pull/3642)
- Added configurable port names for Gateway discovery through
  `--kong-admin-svc-port-names`. This flag accepts a list of port names that
  Admin API service ports will be matched against.
  [#3556](https://github.com/Kong/kubernetes-ingress-controller/pull/3556)
- Added `dataplane` metrics label for `ingress_controller_configuration_push_count`
  and `ingress_controller_configuration_push_duration_milliseconds`. This means
  that all time series for those metrics will get a new label designating the
  address of the dataplane that the configuration push has been targeted for.
  [#3521](https://github.com/Kong/kubernetes-ingress-controller/pull/3521)
- In DB-less mode, failure to push a config will now generate Kubernetes Events
  with the reason `KongConfigurationApplyFailed` and an `InvolvedObject`
  indicating which Kubernetes resource was responsible for the broken Kong
  configuration.
  [#3446](https://github.com/Kong/kubernetes-ingress-controller/pull/3446)
- Leader election is enabled by default when Kong Gateway discovery is enabled.
  [#3529](https://github.com/Kong/kubernetes-ingress-controller/pull/3529)
- Added flag `--konnect-refresh-node-period` to set the period of uploading 
  status of KIC instance to Konnect runtime group.
  [#3533](https://github.com/Kong/kubernetes-ingress-controller/pull/3533)
- Replaced service account's token static secret with a projected volume in 
  deployment manifests.
  [#3563](https://github.com/Kong/kubernetes-ingress-controller/pull/3563)
- Added `GRPCRoute` controller and implemented basic `GRPCRoute` functionality.
  [#3537](https://github.com/Kong/kubernetes-ingress-controller/pull/3537)
- Included Konnect sync and Gateway discovery features in telemetry reports.
  [#3588](https://github.com/Kong/kubernetes-ingress-controller/pull/3588)
- Upload the status of controlled Kong gateway nodes to Konnect when syncing with 
  Konnect is enabled by setting the flag `--konnect-sync-enabled` to true. 
  If gateway discovery is enabled via `--kong-admin-svc` flag, the hostname of a node 
  corresponding to each Kong gateway instance will use `<pod_namespace>/<pod_name>` 
  format, where `pod_namespace` and `pod_name` are the namespace and name of the Kong 
  gateway pod. If gateway discovery is disabled, the Kong gateway nodes will use `gateway_<address>` 
  as the hostname, where `address` is the Admin API address used by KIC.
  [#3587](https://github.com/Kong/kubernetes-ingress-controller/pull/3587)
- All all-in-one DB-less deployment manifests will now use separate deployments 
  for the controller and the proxy. This enables the proxy to be scaled independently
  of the controller. The old `all-in-one-dbless.yaml` manifest has been deprecated and 
  renamed to `all-in-one-dbless-legacy.yaml`. It will be removed in a future release.
  [#3629](https://github.com/Kong/kubernetes-ingress-controller/pull/3629)
- The RequestRedirect Gateway API filter is now supported and translated
  to the proper set of Kong plugins.
  [#3702](https://github.com/Kong/kubernetes-ingress-controller/pull/3702)

### Fixed

- Fixed the issue where the status of an ingress is not updated when `secretName` is
  not specified in `ingress.spec.tls`.
  [#3719](https://github.com/Kong/kubernetes-ingress-controller/pull/3719)
- Fixed incorrectly set parent status for Gateway API routes
  [#3732](https://github.com/Kong/kubernetes-ingress-controller/pull/3732)
- Disabled non-functioning mesh reporting when `--watch-namespaces` flag set.
  [#3336](https://github.com/Kong/kubernetes-ingress-controller/pull/3336)
- Fixed the duplicate update of status of `HTTPRoute` caused by incorrect check
  of whether status is changed.
  [#3346](https://github.com/Kong/kubernetes-ingress-controller/pull/3346)
- Change existing `resolvedRefs` condition in status `HTTPRoute` if there is
  already one to avoid multiple appearance of conditions with same type
  [#3386](https://github.com/Kong/kubernetes-ingress-controller/pull/3386)
- Event messages for invalid multi-Service backends now indicate their derived
  Kong resource name.
  [#3318](https://github.com/Kong/kubernetes-ingress-controller/pull/3318)
- Removed a duplicated status update of the HTTPRoute, which led to a potential
  status flickering.
  [#3451](https://github.com/Kong/kubernetes-ingress-controller/pull/3451)
- Made Admission Webhook fetch the latest list of Gateways to avoid calling
  outdated services set statically during the setup.
  [#3601](https://github.com/Kong/kubernetes-ingress-controller/pull/3601)
- Fixed the way configuration flags `KongAdminSvc` and `PublishService` are
  checked for being set. The old one was always evaluating to `true`.
  [#3602](https://github.com/Kong/kubernetes-ingress-controller/pull/3602)

### Under the hood

- Controller manager scheme is constructed based on the provided feature gates
  [#3539](https://github.com/Kong/kubernetes-ingress-controller/pull/3539)
- Updated the compiler to [Go v1.20](https://golang.org/doc/go1.20)
  [#3540](https://github.com/Kong/kubernetes-ingress-controller/issues/3540)
- Fixed an issue with the fake clientset using the wrong group name
  [#3517](https://github.com/Kong/kubernetes-ingress-controller/issues/3517)

### Deprecated

- `kong-custom-entities-secret` flag has been marked as deprecated and will be
  removed in 3.0.
  [#3262](https://github.com/Kong/kubernetes-ingress-controller/pull/3262)

## [2.8.2]

> Release date: 2022-03-30

### Under the hood

- Updated Golang from 1.19.4 to 1.20.1 to address several CVEs.
  [#3775](https://github.com/Kong/kubernetes-ingress-controller/pull/3775)


## [2.8.1]

> Release date: 2022-01-04

### Fixed

- When `CombinedRoutes` is turned on, translator will replace each occurrence of
  `*` in `Ingress`'s host to `_` in kong route names because `*` is not
  allowed in kong route names.
  [#3312](https://github.com/Kong/kubernetes-ingress-controller/pull/3312)
- Fix an issue with `CombinedRoutes`, which caused the controller to fail when
  creating config for Ingress when backend services specified only port names
  [#3313](https://github.com/Kong/kubernetes-ingress-controller/pull/3313)
- Parse `ttl` field of key-auth credentials in secrets to `int` type before filling
  to kong credentials to fix the invalid type error.
  [#3319](https://github.com/Kong/kubernetes-ingress-controller/pull/3319)

## [2.8.0]

> Release date: 2022-12-19

### Breaking changes

- The `CombinedRoutes` feature flag is enabled by default, and traditional
  route generation is deprecated. This reduces configuration size without
  affecting routing, but does change route names and IDs. Metrics
  monitors or other systems that track data by route name or ID will see a
  break in continuity. The [feature gates document](https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md#differences-between-traditional-and-combined-routes)
  covers changes in greater detail. Please [comment on the deprecation
  issue](https://github.com/Kong/kubernetes-ingress-controller/issues/3131)
  if you have questions or concerns about the transition.
  [#3132](https://github.com/Kong/kubernetes-ingress-controller/pull/3132)

### Known issues

- Processing of custom entities through `--kong-custom-entities-secret` flag or
  `CONTROLLER_KONG_CUSTOM_ENTITIES_SECRET` environment variable does not work.
  [#3278](https://github.com/Kong/kubernetes-ingress-controller/issues/3278)

### Deprecated

- KongIngress' `proxy` and `route` fields are now deprecated in favor of
  Service and Ingress annotations. The annotations will become the only
  means of configuring those settings in 3.0 release.
  [#3246](https://github.com/Kong/kubernetes-ingress-controller/pull/3246)

### Added

- Added `HTTPRoute` support for `CombinedRoutes` feature. When enabled,
  `HTTPRoute.HTTPRouteRule` objects with identical `backendRefs` generate a 
  single Kong service instead of a service per rule, and 
  `HTTPRouteRule.HTTPRouteMatche` objects using the same `backendRefs` can be 
  consolidated into a single Kong route instead of always creating a route per 
  match, reducing configuration size.
  The following limitations apply:
  - `HTTPRouteRule` objects cannot be consolidated into a single Kong Service 
    if they belong to different `HTTPRoute`.
  - `HTTPRouteRule` objects cannot be consolidated into a single Kong Service 
    if they have different `HTTPRouteRule.HTTPBackendRef[]` objects. The order
    of the backend references is not important.
  - `HTTPRouteMatch` objects cannot be consolidated into a single Kong Route
    if parent `HTTPRouteRule` objects cannot be consolidated into a single Kong Service.
  - `HTTPRouteMatch` objects cannot be consolidated into a single Kong Route
    if parent `HTTPRouteRule` objects have different `HTTPRouteRule.HTTPRouteFilter[]` filters.
  - `HTTPRouteMatch` objects cannot be consolidated into a single Kong Route 
    if they have different matching spec (`HTTPHeaderMatch.Headers`, `HTTPHeaderMatch.QueryParams`, 
    `HTTPHeaderMatch.Method`). Different `HTTPHeaderMatch.Path` paths between 
    `HTTPRouteMatch[]` objects does not prevent consolidation.
  This change does not functionally impact routing: requests that went to a given Service
  using the original method still go to the same Service when `CombinedRoutes` is enabled.
  [#3008](https://github.com/Kong/kubernetes-ingress-controller/pull/3008)
  [#3060]https://github.com/Kong/kubernetes-ingress-controller/pull/3060)
- Added `--cache-sync-timeout` flag allowing to change the default controllers' 
  cache synchronisation timeout. 
  [#3013](https://github.com/Kong/kubernetes-ingress-controller/pull/3013)
- Secrets validation introduced: CA certificates won't be synchronized
  to Kong if the certificate is expired.
  [#3063](https://github.com/Kong/kubernetes-ingress-controller/pull/3063)
- Changed the logic of storing secrets into object cache. Now only the secrets
  that are possibly used in Kong configuration are stored into cache, and the 
  irrelevant secrets (e.g: service account tokens) are not stored. This change
  is made to reduce memory usage of the cache.
  [#3047](https://github.com/Kong/kubernetes-ingress-controller/pull/3047)
- Services support annotations for connect, read, and write timeouts.
  [#3121](https://github.com/Kong/kubernetes-ingress-controller/pull/3121)
- Services support annotations for retries.
  [#3121](https://github.com/Kong/kubernetes-ingress-controller/pull/3121)
- Routes support annotations for headers. These use a special
  `konghq.com/headers.HEADERNAME` format. For example, adding
  `konghq.com/headers.x-example: green` to an Ingress will create routes that
  only match requests with an `x-example: green` request header.
  [#3121](https://github.com/Kong/kubernetes-ingress-controller/pull/3121)
  [#3155](https://github.com/Kong/kubernetes-ingress-controller/pull/3155)
- Routes support annotations for path handling.
  [#3121](https://github.com/Kong/kubernetes-ingress-controller/pull/3121)
- Warning Kubernetes API events with a `KongConfigurationTranslationFailed` 
  reason are recorded when:
  - CA secrets cannot be properly translated into Kong configuration
    [#3125](https://github.com/Kong/kubernetes-ingress-controller/pull/3125)
  - Annotations in services backing a single route do not match
    [#3130](https://github.com/Kong/kubernetes-ingress-controller/pull/3130)
  - A service's referred client-cert does not exist.
    [#3137](https://github.com/Kong/kubernetes-ingress-controller/pull/3137)
  - One of `netv1.Ingress` related issues occurs (e.g. backing Kubernetes service couldn't 
    be found, matching Kubernetes service port couldn't be found).
    [#3138](https://github.com/Kong/kubernetes-ingress-controller/pull/3138)
  - A Gateway Listener has more than one CertificateRef specified or refers to a Secret 
    that has no valid TLS key-pair.
    [#3147](https://github.com/Kong/kubernetes-ingress-controller/pull/3147)
  - An Ingress refers to a TLS secret that does not exist or
    has no valid TLS key-pair.
    [#3150](https://github.com/Kong/kubernetes-ingress-controller/pull/3150)
  - An HTTPRoute has no backendRefs specified.
    [#3167](https://github.com/Kong/kubernetes-ingress-controller/pull/3167)
- CRDs' validations improvements: `UDPIngressRule.Port`, `IngressRule.Port` and `IngressBackend.ServiceName`
  instead of being validated in the Parser, are validated by the Kubernetes API now.
  [#3136](https://github.com/Kong/kubernetes-ingress-controller/pull/3136)
- Gateway API: Implement port matching for HTTPRoute, TCPRoute and TLSRoute as defined in
  [GEP-957](https://gateway-api.sigs.k8s.io/geps/gep-957/)
  [#3129](https://github.com/Kong/kubernetes-ingress-controller/pull/3129)
  [#3226](https://github.com/Kong/kubernetes-ingress-controller/pull/3226)
- Gateway API: Matching routes by `Listener.AllowedRoutes`
  [#3181](https://github.com/Kong/kubernetes-ingress-controller/pull/3181)
- Admission webhook will warn if any of `KongIngress` deprecated fields gets populated.
  [#3261](https://github.com/Kong/kubernetes-ingress-controller/pull/3261)

### Fixed

- The controller now logs an error for and skips multi-Service rules that have
  inconsistent Service annotations. Previously this issue prevented the
  controller from applying configuration until corrected.
  [#2988](https://github.com/Kong/kubernetes-ingress-controller/pull/2988)
- Gateway API has been updated to 0.5.1. That version brought in some changes
  in the conformance tests logic. Now, when the TLS config of a listener
  references a non-existing secret, the listener ResolvedRefs condition reason
  is set to InvalidCertificateRef. In addition, if a TLS config references a
  secret in another namespace, and no ReferenceGrant allows that
  reference, the listener ResolvedRefs condition reason is set to
  RefNotPermitted.
  [#3024](https://github.com/Kong/kubernetes-ingress-controller/pull/3024)
- The `distroless` target is now the last target in the Dockerfile. This makes
  it the default target if `docker buildx build` is invoked without a target.
  While custom image build pipelines _should_ specify a target, this change
  makes the default the same target released as the standard
  `kong/kubernetes-ingress-controller:X.Y.Z` tags in the official repo.
  [#3043](https://github.com/Kong/kubernetes-ingress-controller/pull/3043)
- The controller will no longer crash in case of missing CRDs installation.
  Instead, an explicit message will be logged, informing that a given resource
  controller has been disabled.
  [#3013](https://github.com/Kong/kubernetes-ingress-controller/pull/3013)
- Improve signal handling and cancellation. With this change broken connection to
  Admin API and/or initial data plane sync can be cancelled properly.
  [#3076](https://github.com/Kong/kubernetes-ingress-controller/pull/3076)
- Admin and proxy listens in the deploy manifests now use the same parameters
  as the default upstream kong.conf.
  [#3165](https://github.com/Kong/kubernetes-ingress-controller/pull/3165)
- Fix the behavior of filtering hostnames in `HTTPRoute` when listeners 
  of parent gateways specified hostname.
  If an `HTTPRoute` does not specify hostnames, and one of its parent listeners
  has not specified hostname, the `HTTPRoute` matches any hostname. 
  If an `HTTPRoute` specifies hostnames, and no intersecting hostnames 
  could be found in its parent listners, it is not accepted.
  [#3180](https://github.com/Kong/kubernetes-ingress-controller/pull/3180)
- Matches `sectionName` in parentRefs of route objects in gateway API. Now 
  if a route specifies `sectionName` in parentRefs, and no listener can 
  match the specified name, the route is not accepted.
  [#3230](https://github.com/Kong/kubernetes-ingress-controller/pull/3230)
- If there's no matching Kong listener for a protocol specified in a Gateway's
  Listener, only one `Detached` condition is created in the Listener's status.
  [#3257](https://github.com/Kong/kubernetes-ingress-controller/pull/3257)

## [2.7.0]

> Release date: 2022-09-26

2.7 patches several bugs in 2.6.0. One of these required a breaking change. The
breaking change is not expected to affect most configurations, but does require
a minor version bump to comply with semver. If you have not already upgraded to
2.6, you should upgrade directly from 2.5 to 2.7, and follow the 2.6 upgrade
instructions and the [revised Kong 3.x upgrade instructions](https://docs.konghq.com/kubernetes-ingress-controller/2.7.x/guides/upgrade-kong-3x).

### Breaking changes

- Ingress paths that begin with `/~` are now treated as regular expressions,
  and are translated into a Kong route path that begins with `~` instead of
  `/~`. To preserve the existing translation, set `konghq.com/regex-prefix` to
  some value. For example, if you set `konghq.com/regex-prefix: /@`, paths
  beginning with `/~` will result in route paths beginning in `/~`, whereas
  paths beginning in `/@` will result in route paths beginning in `~`.
  [#2956](https://github.com/Kong/kubernetes-ingress-controller/pull/2956)

### Added

- The controller-specific `/~` prefix translates to the Kong `~` prefix, as
  Ingress does not allow paths that do not begin in `/`. The prefix can be
  overriden by setting a `konghq.com/regex-prefix` annotation, for routes that
  need their paths to actually begin with `/~`
  [#2956](https://github.com/Kong/kubernetes-ingress-controller/pull/2956)
- Prometheus metrics now highlight configuration push failures caused by
  conflicts. The `ingress_controller_configuration_push_count` Prometheus
  metric now reports `success="false"` with a `failure_reason="conflict|other"`
  label, distinguishing configuration conflicts from other errors (transient
  network errors, Kong offline, Kong reported non-conflict error, etc.).
  [#2965](https://github.com/Kong/kubernetes-ingress-controller/pull/2965)

### Fixed

- The legacy regex heuristic toggle on IngressClassParameters now works when
  the combined routes feature flag is enabled.
  [#2942](https://github.com/Kong/kubernetes-ingress-controller/pull/2942)
- Handles Kubernetes versions that do not support namespaced
  IngressClassParameters without panicking. Although the controller will run on
  clusters without the `IngressClassNamespacedParams` feature gate enabled
  (1.21) or without it available (<1.21), these clusters do not support the
  legacy regular expression heuristic IngressClassParameters option. These
  versions are EOL, and we advise users to upgrade to Kubernetes 1.22 or later
  before upgrading to KIC 2.6+ or Kong 3.0+.
  [#2970](https://github.com/Kong/kubernetes-ingress-controller/pull/2970)

## [2.6.0]

> Release date: 2022-09-14

### Breaking changes

- Kong 3.x changes regular expression configuration and the controller does not
  handle these changes automatically. You will need to enable compatibility
  features initially and then update Ingress configuration before disabling
  them. This procedure is covered in the [Kong 3.x upgrade guide for the
  controller](https://docs.konghq.com/kubernetes-ingress-controller/2.6.x/guides/upgrade-kong-3x).
- When using the `CombinedRoutes=true` feature gate, Ingress rules with no
  PathType now use ImplementationSpecific instead of Prefix. While Kong's
  ImplementationSpecific handling is similar to Prefix, it does not require
  that the prefix be a directory: an ImplementationSpecific `/foo` will match
  `/foo`, `/foo/`, and `/foo/.*`, whereas Prefix will only match the latter
  two. If you have rules with no PathType, use `CombinedRoutes=true`, and wish
  to preserve existing behavior, add `PathType=prefix` configuration to those
  rules.
  [#2883](https://github.com/Kong/kubernetes-ingress-controller/pull/2883)
- The GatewayClass objects now require the annotation
  "konghq.com/gatewayclass-unmanaged" to be reconciled by the controller.
  The annotation "konghq.com/gateway-unmanaged" is not considered anymore and
  doesn't need to be set on Gateways to be reconciled. Only the Gateways using
  an unmanaged GatewayClass are reconciled.
  [#2917](https://github.com/Kong/kubernetes-ingress-controller/pull/2917)

#### Added

- IngressClassParameters now supports a `enableLegacyRegexDetection` boolean
  field. Kong 3.x+ requires adding a `~` prefix to regular expression paths,
  whereas Kong 2.x and earlier attempted to detect regular expression paths
  using heuristics. By default, if you use regular expression paths and wish to
  migrate to Kong 3.x, you must update all Ingresses to use this prefix.
  Enabling this field will use the 2.x heuristic to detect if an Ingress path
  is a regular expression and add the prefix for you. You should update your
  Ingresses to include the new prefix as soon as possible after upgrading to
  Kong 3.x+, however, as the heuristic has known flaws that will not be fixed.
  [#2883](https://github.com/Kong/kubernetes-ingress-controller/pull/2883)
- Added support for plugin ordering (requires Kong Enterprise 3.0 or higher).
  [#2657](https://github.com/Kong/kubernetes-ingress-controller/pull/2657)
- The all-in-one manifests now use a separate ClusterRole for Gateway API
  resources, allowing non-admin users to apply these manifests (minus the
  Gateway API role) on clusters without Gateway API CRDs installed.
  [#2529](https://github.com/Kong/kubernetes-ingress-controller/issues/2529)
- Gateway API support which had previously been off by default behind a feature
  gate (`--feature-gates=Gateway=true`) is now **on by default** and covers beta
  stage APIs (`GatewayClass`, `Gateway`, and `HTTPRoute`). Alpha stage APIs
  (`TCPRoute`, `UDPRoute`, `TLSRoute`, `ReferenceGrant`) have been moved behind
  a different feature gate called `GatewayAlpha` and are off by default. When
  upgrading if you're using the alpha APIs, switch your feature gate flags to
  `--feature-gates=GatewayAlpha=true` to keep them enabled.
  [#2781](https://github.com/Kong/kubernetes-ingress-controller/pull/2781)
- Added all the Gateway-related conformance tests.
  [#2777](https://github.com/Kong/kubernetes-ingress-controller/issues/2777)
- Added all the HTTPRoute-related conformance tests.
  [#2776](https://github.com/Kong/kubernetes-ingress-controller/issues/2776)
- Added support for Kong 3.0 upstream `query_arg` and `uri_capture` hash
  configuration to KongIngress.
  [#2822](https://github.com/Kong/kubernetes-ingress-controller/issues/2822)
- Added support for Gateway API's `v1beta1` versions of: `GatewayClass`, `Gateway`
  and `HTTPRoute`.
  [#2889](https://github.com/Kong/kubernetes-ingress-controller/issues/2889)
  [#2894](https://github.com/Kong/kubernetes-ingress-controller/issues/2894)
  [#2900](https://github.com/Kong/kubernetes-ingress-controller/issues/2900)
- Manifests now use `/bin/bash` instead of `/bin/sh` and use bash-based
  connectivity checks for compatibility with the new Debian Kong images.
  [#2923](https://github.com/Kong/kubernetes-ingress-controller/issues/2923)

#### Fixed

- When `Endpoints` could not be found for a `Service` to add them as targets of
  a Kong `Upstream`, this would produce a log message at `error` and `warning`
  levels which was inaccurate because this condition is often expected when
  `Pods` are being provisioned. Those log entries now report at `info` level.
  [#2820](https://github.com/Kong/kubernetes-ingress-controller/issues/2820)
  [#2825](https://github.com/Kong/kubernetes-ingress-controller/pull/2825)
- Added `mtls-auth` to the admission webhook supported credential types list.
  [#2739](https://github.com/Kong/kubernetes-ingress-controller/pull/2739)
- Disabled additional IngressClass lookups in other reconcilers when the
  IngressClass reconciler is disabled.
  [#2724](https://github.com/Kong/kubernetes-ingress-controller/pull/2724)
- ReferencePolicy support has been dropped in favor of the newer ReferenceGrant API.
  [#2775](https://github.com/Kong/kubernetes-ingress-controller/pull/2772)
- Fixed a bug that caused the `Knative` feature gate to not be checked. Since our
  knative integration is on by default and because it gets very little usage
  this likely did not cause any troubles for anyone as all fixing this will do
  is make it possible to disable the knative controller using the feature gate.
  (it is also possible to control it via the `--enable-controller-knativeingress`
  which was working properly).
  [#2781](https://github.com/Kong/kubernetes-ingress-controller/pull/2781)
- Treat status conditions in `Gateway` and `GatewayClass` as snapshots, replace
  existing conditions with same type on setting conditions.
  [#2791](https://github.com/Kong/kubernetes-ingress-controller/pull/2791)
- Update Listener statuses whenever they change, not just on Gateway creation.
  [#2797](https://github.com/Kong/kubernetes-ingress-controller/pull/2797)
- StripPath for `HTTPRoute`s is now disabled by default to be conformant with the
  Gateway API requirements.
  #[#2737](https://github.com/Kong/kubernetes-ingress-controller/pull/2737)

#### Under the hood

- Updated the compiler to [Go v1.19](https://golang.org/doc/go1.19)
  [#2794](https://github.com/Kong/kubernetes-ingress-controller/issues/2794)

## [2.5.0]

> Release date: 2022-07-11

#### Breaking changes in Gateway API technical preview:

- The controller no longer overrides Gateway Listeners with a list of Listeners
  derived from Kong configuration. User-provided Listener lists are preserved
  as-is. Listener status information indicates if a requested Listener is not
  ready because of missing Kong listen configuration. This is necessary to
  properly support allowed routes and TLS configuration in Listeners, which
  would otherwise be wiped out by automatic updates. This has no immediate
  impact on existing Gateway resources used with previous versions: their
  automatically-set Listeners are now treated as user-defined Listeners and
  will not be modified by upgrading. This only affects new Gateway resources:
  you will need to populate the Listeners you want, and they will need to match
  Kong's listen configuration to become ready.
  [#2555](https://github.com/Kong/kubernetes-ingress-controller/pull/2555)

#### Added

- Updated Gateway API dependencies to [v0.5.0][gw-v0.5.0] and updated `examples`
  directory to use `v1beta1` versions of APIs where applicable.
  [#2691](https://github.com/Kong/kubernetes-ingress-controller/pull/2691)
- Added support for Gateway Listener TLS configuration, to enable full use of
  TLSRoute and HTTPS HTTPRoutes.
  [#2580](https://github.com/Kong/kubernetes-ingress-controller/pull/2580)
- Added information about service mesh deployment and distribution in telemetry data reported to Kong.
  [#2642](https://github.com/Kong/kubernetes-ingress-controller/pull/2642)

[gw-v0.5.0]:https://github.com/kubernetes-sigs/gateway-api/releases/tag/v0.5.0

#### Fixed

- Fixed the problem that logs from reporter does not appear in the pod log.
  [#2645](https://github.com/Kong/kubernetes-ingress-controller/pull/2645)

## [2.4.2]

> Release date: 2022-06-30

#### Fixed

- Fix an issue with ServiceAccount token mount.
  [#2620](https://github.com/Kong/kubernetes-ingress-controller/issues/2620)
  [#2626](https://github.com/Kong/kubernetes-ingress-controller/issues/2626)

## [2.4.1]

> Release date: 2022-06-22

#### Added

- Increased the default Kong admin API timeout from 10s to 30s and added a
  log mentioning the flag to increase it further.
  [#2594](https://github.com/Kong/kubernetes-ingress-controller/issues/2594)

#### Fixed

- Disabling the IngressClass controller now disables IngressClass watches in
  other controllers. This fixes a crash on Kubernetes versions that do not
  offer an IngressClass version that KIC can read.
  [#2577](https://github.com/Kong/kubernetes-ingress-controller/issues/2577)

## [2.4.0]

> Release date: 2022-06-14

#### Added

- A new gated feature called `CombinedRoutes` has been added. Historically
  a `kong.Route` would be created for _each path_ on an `Ingress` resource
  in the phase where Kubernetes resources are translated to Kong Admin API
  configuration. This new feature changes how `Ingress` resources are
  translated so that a single route can be created for any unique combination
  of ingress object, hostname, service and port which has multiple paths.
  This option is helpful for end-users who are making near constant changes
  to their configs (e.g. constantly adding, updating, and removing `Ingress`
  resources) at scale, and users that have enormous numbers of paths all
  pointing to a single Kubernetes `Service` as it can significantly reduce
  the overall size of the dataplane configuration that is pushed to the Kong
  Admin API. This feature is expected to be disruptive (routes may be dropped
  briefly in postgres mode when switching to this mode) so for the moment it
  is behind a feature gate while we continue to iterate on it and evaluate it
  and seek a point where it would become the default behavior. Enable it with
  the controller argument `--feature-gates=CombinedRoutes`.
  [#2490](https://github.com/Kong/kubernetes-ingress-controller/issues/2490)
- `UDPRoute` resources now support multiple backendRefs for load-balancing.
  [#2405](https://github.com/Kong/kubernetes-ingress-controller/issues/2405)
- `TCPRoute` resources now support multiple backendRefs for load-balancing.
  [#2405](https://github.com/Kong/kubernetes-ingress-controller/issues/2405)
- `TCPRoute` resources are now supported.
  [#2086](https://github.com/Kong/kubernetes-ingress-controller/issues/2086)
- `HTTPRoute` resources now support multiple `backendRefs` with a round-robin
  load-balancing strategy applied by default across the `Endpoints` or the
  `Services` (if the `ingress.kubernetes.io/service-upstream`
  annotation is set). They also now support weights to enable more
  fine-tuning of the load-balancing between those backend services.
  [#2166](https://github.com/Kong/kubernetes-ingress-controller/issues/2166)
- `Gateway` resources now honor [`listener.allowedRoutes.namespaces`
  filters](https://gateway-api.sigs.k8s.io/v1alpha2/references/spec/#gateway.networking.k8s.io/v1alpha2.RouteNamespaces).
  Note that the unmanaged Kong Gateway implementation populates listeners
  automatically based on the Kong Service and Deployment, and user-provided
  `allowedRoutes` filters are merged into generated listeners with the same
  protocol.
  [#2389](https://github.com/Kong/kubernetes-ingress-controller/issues/2389)
- Added `--skip-ca-certificates` flag to ignore CA certificate resources for
  [use with multi-workspace environments](https://github.com/Kong/deck/blob/main/CHANGELOG.md#v1120).
  [#2341](https://github.com/Kong/kubernetes-ingress-controller/issues/2341)
- Gateway API Route types now support cross-namespace BackendRefs if a
  [ReferencePolicy](https://gateway-api.sigs.k8s.io/v1alpha2/api-types/referencepolicy/)
  permits them.
  [#2451](https://github.com/Kong/kubernetes-ingress-controller/issues/2451)
- Added description of each field of `kongIngresses` CRD.
  [#1766](https://github.com/Kong/kubernetes-ingress-controller/issues/1766)
- Added support for `TLSRoute` resources.
  [#2476](https://github.com/Kong/kubernetes-ingress-controller/issues/2476)
- Added `--term-delay` flag to support setting a time delay before processing
  `SIGTERM` and `SIGINT` signals. This was added to specifically help in
  situations where the Kong Gateway has a load-balancer in front of it to help
  stagger and stabilize the shutdown procedure when the load-balancer is
  draining or otherwise needs to remove the Gateway from it's rotation.
  [#2494](https://github.com/Kong/kubernetes-ingress-controller/pull/2494)
- Added `kong-ingress-controller` category to CRDs
  [#2517](https://github.com/Kong/kubernetes-ingress-controller/pull/2517)
- Added `v1alpha1.IngressClassParameters` CRD and its first field `ServiceUpstream`
  to control the behavior of routing traffic via an upstream service for all services managed
  by an ingress class without the need of adding an annotation to every single one
  [#2535](https://github.com/Kong/kubernetes-ingress-controller/pull/2535)

#### Fixed

- Unmanaged-mode `Gateway` resources which reference a `LoadBalancer` type
  `Service` will now tolerate the IPs/Hosts for that `Service` not becoming
  provisioned (e.g. the `LoadBalancer` implementation is broken or otherwise
  and the `EXTERNAL-IP` is stuck in `<pending>`) and will still attempt to
  configure `Routes` for that `Gateway` as long as the dataplane API can be
  otherwise reached.
  [#2413](https://github.com/Kong/kubernetes-ingress-controller/issues/2413)
- Fixed a race condition in the newer Gateway route controllers which could
  trigger when an object's status was updated shortly after the object was
  cached in the dataplane client.
  [#2446](https://github.com/Kong/kubernetes-ingress-controller/issues/2446)
- Added a mechanism to retry the initial connection to the Kong
  Admin API on controller start to fix an issue where the controller
  pod could crash loop on start when waiting for Gateway readiness 
  (e.g. if the Gateway is waiting for its database to initialize). 
  The new retry mechanism can be manually configured using the 
  `--kong-admin-init-retries` and `--kong-admin-init-retry-delay` flags.
  [#2274](https://github.com/Kong/kubernetes-ingress-controller/issues/2274)
- diff logging now honors log level instead of printing at all log levels. It
  will only print at levels `debug` and `trace`.
  [#2422](https://github.com/Kong/kubernetes-ingress-controller/issues/2422)
- For KNative Ingress resources, KIC now reads both the new style annotation
  `networking.knative.dev/ingress-class` and the deprecated `networking.knative.dev/ingress.class` one
  to adapt to [what has already been done in knative](https://github.com/knative/networking/pull/522).
  [#2485](https://github.com/Kong/kubernetes-ingress-controller/issues/2485)
- Remove KongIngress support for Gateway API Route objects and Services referenced
  by those Routes. This disables an undocumented ability of customizing Gateway API
  `*Route` objects and `Service`s that are set as backendRefs for those `*Route`s
  via `konghq.com/override` annotations.
  [#2554](https://github.com/Kong/kubernetes-ingress-controller/issues/2554)
- Fixed a vulnerability that permission could be escalated by running custom lua
  scripts.
  [#2572](https://github.com/Kong/kubernetes-ingress-controller/pull/2572)

## [2.3.1]

> Release date: 2022-04-07

#### Fixed

- Fixed an issue where admission controllers configured without certificates
  would incorrectly detect invalid configuration and prevent the controller
  from starting.
  [#2403](https://github.com/Kong/kubernetes-ingress-controller/pull/2403)

## [2.3.0]

> Release date: 2022-04-05

#### Breaking changes

- HTTPRoute header matches no longer interpret CSV values as multiple match
  values, as this was not part of the HTTPRoute specification. Multiple values
  should use regular expressions instead.
  [#2302](https://github.com/Kong/kubernetes-ingress-controller/pull/2302)

#### Added

- `Gateway` resources which have a `LoadBalancer` address among their list of
  addresses will have those addresses listed on the top for convenience, and
  so that those addresses are made prominent in the `kubectl get gateways`
  short view.
  [#2339](https://github.com/Kong/kubernetes-ingress-controller/pull/2339)
- The controller manager can now be flagged with a client certificate to use
  for mTLS authentication with the Kong Admin API.
  [#1958](https://github.com/Kong/kubernetes-ingress-controller/issues/1958)
- Deployment manifests now include an IngressClass resource and permissions to
  read IngressClass resources.
  [#2292](https://github.com/Kong/kubernetes-ingress-controller/pull/2292)
- The controller now reads IngressClass resources to determine if its
  IngressClass is the default IngressClass. If so, the controller will ingest
  resources that require a class (Ingress, KongConsumer, KongClusterPlugin,
  etc.) but have none set.
  [#2313](https://github.com/Kong/kubernetes-ingress-controller/pull/2313)
- HTTPRoute header matches now support regular expressions.
  [#2302](https://github.com/Kong/kubernetes-ingress-controller/pull/2302)
- HTTPRoutes that define multiple matches for the same header are rejected to
  comply with the HTTPRoute specification.
  [#2302](https://github.com/Kong/kubernetes-ingress-controller/pull/2302)
- Admission webhook certificate files now track updates to the file, and will
  update when the corresponding Secret has changed.
  [#2258](https://github.com/Kong/kubernetes-ingress-controller/pull/2258)
- Added support for Gateway API [UDPRoute](https://gateway-api.sigs.k8s.io/v1alpha2/references/spec/#gateway.networking.k8s.io/v1alpha2.UDPRoute)
  resources.
  [#2363](https://github.com/Kong/kubernetes-ingress-controller/pull/2363)
- The controller can now detect whether a Kong container has crashed and needs
  a configuration push. Requires Kong 2.8+.
  [#2343](https://github.com/Kong/kubernetes-ingress-controller/pull/2343)

#### Fixed

- Fixed an issue where duplicated route names in `HTTPRoute` resources with
  multiple matches would cause the Kong Admin API to collide the routes into
  one, effectively dropping routes for services beyond the first.
  [#2345](https://github.com/Kong/kubernetes-ingress-controller/pull/2345)
- Status updates for `HTTPRoute` objects no longer mark the resource as
  `ConditionRouteAccepted` until the object has been successfully configured
  in Kong Gateway at least once, as long as `--update-status`
  is enabled (enabled by default).
  [#2339](https://github.com/Kong/kubernetes-ingress-controller/pull/2339)
- Status updates for `HTTPRoute` now properly use the `ConditionRouteAccepted`
  value for parent `Gateway` conditions when the route becomes configured in
  the `Gateway` rather than the previous random `"attached"` string.
  [#2339](https://github.com/Kong/kubernetes-ingress-controller/pull/2339)
- Fixed a minor issue where addresses on `Gateway` resources would be
  duplicated depending on how many listeners are configured.
  [#2339](https://github.com/Kong/kubernetes-ingress-controller/pull/2339)
- Unconfigured fields now use their default value according to the Kong proxy
  instance's reported schema. This addresses an issue where configuration
  updates would send unnecessary requests to clear a default value.
  [#2286](https://github.com/Kong/kubernetes-ingress-controller/pull/2286)
- Certificate selection for hostnames is no longer random if both certificate
  Secrets have the same creation timestamp, and no longer results in
  unnecessary configuration updates.
  [#2338](https://github.com/Kong/kubernetes-ingress-controller/pull/2338)

## [2.2.1]

> Release date: 2022/02/15

#### Fixed

- Added mitigation for an issue where controllers may briefly delete and
  recreate configuration upon gaining leadership while populating their
  Kubernetes object cache.
  [#2255](https://github.com/Kong/kubernetes-ingress-controller/pull/2255)

## [2.2.0]

> Release date: 2022/02/04

#### Added

- Support for Kubernetes [Gateway APIs][gwapis] is now available [by enabling 
  the `Gateway` feature gate](https://docs.konghq.com/kubernetes-ingress-controller/2.2.x/guides/using-gateway-api/).
  This is an alpha feature, with limited support for the `HTTPRoute` API.
  [Gateway Milestone 1][gwm1]
- Kubernetes client rate limiting can now be configured using `--apiserver-qps`
  (default 100) and `--apiserver-burst` (default 300) settings. Defaults have
  been increased to prevent ratelimiting under normal loads.
  [#2169](https://github.com/Kong/kubernetes-ingress-controller/issues/2169)
- The KIC Grafana dashboard [is now published on grafana.com](https://grafana.com/grafana/dashboards/15662).
  [#2235](https://github.com/Kong/kubernetes-ingress-controller/issues/2235)

[gwapis]:https://github.com/kubernetes-sigs/gateway-api
[gwm1]:https://github.com/Kong/kubernetes-ingress-controller/milestone/21

#### Fixed

- Fixed an issue where validation could fail for credentials secrets if the
  `value` for a unique constrained `key` were updated in place while linked
  to a managed `KongConsumer`.
  [#2190](https://github.com/Kong/kubernetes-ingress-controller/issues/2190)
- The controller now retries status updates if the publish service LoadBalancer
  has not yet provisioned. This fixes an issue where controllers would not
  update status until the first configuration change after the LoadBalancer
  became ready.

## [2.1.1]

> Release date: 2022/01/05

2.1.1 has no user-facing changes from 2.1.0. It updates a certificate used in
the test environment which expired during the 2.1.0 release process.
[#2133](https://github.com/Kong/kubernetes-ingress-controller/pull/2133)

## [2.1.0]

> Release date: 2022/01/05

**Note:** the admission webhook updates originally released in [2.0.6](#206)
are _not_ applied automatically by the upgrade. If you set one up previously,
you should edit it (`kubectl edit validatingwebhookconfiguration
kong-validations`) and add `kongclusterplugins` under the `resources` block for
the `configuration.konghq.com` API group.

#### Breaking changes

- The `--leader-elect` flag has been deprectated and will be removed in a
  future release. Leader election is now enabled or disabled automatically
  based on the database mode. The flag is no longer honored.
  [#2053](https://github.com/Kong/kubernetes-ingress-controller/issues/2053)
- You must upgrade to 2.0.x before upgrading to 2.1.x to properly handle the
  transition from apiextensions.k8s.io/v1beta1 CRDs to apiextensions.k8s.io/v1
  CRDSs. CRDs are now generated from their underlying Go structures to avoid
  accidental mismatches between implementation and Kubernetes configuration.
  KongIngresses previously included `healthchecks.passive.unhealthy.timeout`
  and `healthchecks.active.unhealthy.timeout` fields that did not match the
  corresponding Kong configuration and had no effect. These are now
  `healthchecks.passive.unhealthy.timeouts` and
  `healthchecks.active.unhealthy.timeouts`, respectively. If you use these
  fields, you must rename them in your KongIngresses before upgrading.
  [#1971](https://github.com/Kong/kubernetes-ingress-controller/pull/1971)

#### Added

- Added validation for `Gateway` objects in the admission webhook
  [#1946](https://github.com/Kong/kubernetes-ingress-controller/issues/1946)
- [Feature Gates][k8s-fg] have been added to the controller manager in order to
  enable alpha/beta/experimental features and provide documentation about those
  features and their maturity over time. For more information see the
  [KIC Feature Gates Documentation][kic-fg].
  [#1970](https://github.com/Kong/kubernetes-ingress-controller/pull/1970)
- a Gateway controller has been added in support of [Gateway APIs][gwapi].
  This controller is foundational and doesn't serve any end-user purpose alone.
  [#1945](https://github.com/Kong/kubernetes-ingress-controller/issues/1945)
- Anonymous reports now use TLS instead of UDP.
  [#2089](https://github.com/Kong/kubernetes-ingress-controller/pull/2089)
- The new `--election-namespace` flag sets the leader election namespace. This
  is normally only used if a controller is running outside a Kubernetes
  cluster.
  [#2053](https://github.com/Kong/kubernetes-ingress-controller/issues/2053)
- There is now a [Grafana dashboard](https://github.com/Kong/kubernetes-ingress-controller/blob/main/grafana.json)
  for the controller metrics.
  [#2035](https://github.com/Kong/kubernetes-ingress-controller/issues/2035)
- TCPIngresses now support TLS passthrough in Kong 2.7+, by setting a
  `konghq.com/protocols: tls_passthrough` annotation.
  [#2041](https://github.com/Kong/kubernetes-ingress-controller/issues/2041)

[k8s-fg]:https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
[kic-fg]:https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md
[gwapi]:https://github.com/kubernetes-sigs/gateway-api

#### Fixed

- Fixed an edge case which could theoretically remove data-plane config for
  objects which couldn't be retrieved from the manager's cached client.
  [#2057](https://github.com/Kong/kubernetes-ingress-controller/pull/2057)
- The validating webhook now validates that required fields data are not empty.
  [#1993](https://github.com/Kong/kubernetes-ingress-controller/issues/1993)
- The validating webhook now validates unique key constraints for KongConsumer
  credentials secrets on update of secrets, and on create or update of
  KongConsumers.
  [#729](https://github.com/Kong/kubernetes-ingress-controller/issues/729)
- Fixed a race condition where multiple actors may simultaneously attempt to
  create the configured Enterprise workspaces.
  [#2070](https://github.com/Kong/kubernetes-ingress-controller/pull/2070)
- Fixed incorrect leader election behavior. Previously, non-leader instances
  would still attempt to update Kong configuration, but would not scan for
  Kubernetes resources to translate into Kong configuration.
  [#2053](https://github.com/Kong/kubernetes-ingress-controller/issues/2053)
- Configuration updates that time out now correctly report a failure.
  [deck #529](https://github.com/Kong/deck/pull/529)
  [#2125](https://github.com/Kong/kubernetes-ingress-controller/pull/2125)

## [2.0.7]

> Release date: 2022/01/19

#### Under the hood

- Anonymous reports now use TLS instead of UDP.
  [#2089](https://github.com/Kong/kubernetes-ingress-controller/pull/2089)

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
- Per documentation and by design, KongClusterPlugin resources require an
  `kubernetes.io/ingress.class` annotation, but this was not fully enforced. In
  2.0, all KongClusterPlugin resources require this annotation set to the
  controller's ingress class. Check your resources to confirm they are annotated
  before upgrading.
  [#2090](https://github.com/Kong/kubernetes-ingress-controller/issues/2090)

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

## [1.3.4]

> Release date: 2022/01/19

#### Under the hood

- Anonymous reports now use TLS instead of UDP.
  [#2089](https://github.com/Kong/kubernetes-ingress-controller/pull/2089)

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


## [0.0.5]

> Release date: 2018/06/02

#### Added

 - Add support for Kong Enterprise Edition 0.32 and above

## [0.0.4] and prior

 - The initial versions  were rapildy iterated to deliver
   a working ingress controller.

[3.2.3]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.2.2...v3.2.3
[3.2.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.2.1...v3.2.2
[3.2.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.2.0...v3.2.1
[3.2.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.1.6...v3.2.0
[3.1.6]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.1.5...v3.1.6
[3.1.5]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.1.4...v3.1.5
[3.1.4]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.1.3...v3.1.4
[3.1.3]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.1.2...v3.1.3
[3.1.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.1.1...v3.1.2
[3.1.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.1.0...v3.1.1
[3.1.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.0.2...v3.1.0
[3.0.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.0.1...v3.0.2
[3.0.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v3.0.0...v3.0.1
[3.0.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.12.0...v3.0.0
[2.12.5]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.12.4...v2.12.5
[2.12.4]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.12.3...v2.12.4
[2.12.3]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.12.2...v2.12.3
[2.12.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.12.1...v2.12.2
[2.12.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.12.0...v2.12.1
[2.12.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.11.1...v2.12.0
[2.11.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.11.0...v2.11.1
[2.11.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.10.4...v2.11.0
[2.10.5]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.10.4...v2.10.5
[2.10.4]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.10.3...v2.10.4
[2.10.3]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.10.2...v2.10.3
[2.10.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.10.1...v2.10.2
[2.10.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.10.0...v2.10.1
[2.10.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.9.3...v2.10.0
[2.9.3]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.9.2...v2.9.3
[2.9.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.9.1...v2.9.2
[2.9.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.9.0...v2.9.1
[2.9.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.8.1...v2.9.0
[2.8.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.8.1...v2.8.2
[2.8.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.8.0...v2.8.1
[2.8.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.7.0...v2.8.0
[2.7.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.6.0...v2.7.0
[2.6.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.5.0...v2.6.0
[2.5.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.4.2...v2.5.0
[2.4.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.4.1...v2.4.2
[2.4.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.4.0...v2.4.1
[2.4.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.3.1...v2.4.0
[2.3.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.3.0...v2.3.1
[2.3.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.2.1...v2.3.0
[2.2.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.2.0...v2.2.1
[2.2.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.1.1...v2.2.0
[2.1.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.1.0...v2.1.1
[2.1.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.6...v2.1.0
[2.0.7]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.6...v2.0.7
[2.0.6]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.5...v2.0.6
[2.0.5]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.4...v2.0.5
[2.0.4]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.3...v2.0.4
[2.0.3]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.2...v2.0.3
[2.0.2]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v2.0.0-beta.2...v2.0.0
[1.3.4]: https://github.com/kong/kubernetes-ingress-controller/compare/1.3.3...1.3.4
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
