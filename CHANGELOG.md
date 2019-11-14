# Table of Contents

 - [0.6.2](#062---20191113)
 - [0.6.1](#061---20191009)
 - [0.6.0](#060---20190917)
 - [0.5.0](#050---20190625)
 - [0.4.0](#040---20190424)
 - [0.3.0](#030---20190108)
 - [0.2.2](#022---20181109)
 - [0.1.3](#013---20181109)
 - [0.2.1](#021---20181026)
 - [0.1.2](#012---20181026)
 - [0.1.1](#011---20180926)
 - [0.2.0](#020---20180921)
 - [0.1.0](#010---20180817)
 - [0.0.5](#005---20180602)
 - [0.0.4 and prior](#004-and-prior)

## [0.6.2] - 2019/11/13

### Summary

This is a minor patch release to fix version parsing issue with new
Kong Enterprise packages.

## [0.6.1] - 2019/10/09

### Summary

This is a minor patch release to update Kong Ingress Controller's
Docker image to use a non-root by default.

## [0.6.0] - 2019/09/17

### Summary

This release brings introduces an Admission Controller for CRDs,
Istio compatibility, support for `networking/ingress`,
Kong 1.3 addtions and enhancements to documentation and deployments.

### Added

- **Service Mesh integration** Kong Ingress Controller can now be deployed
  alongside Service Mesh solutions like Kuma and Istio. In such a deployment,
  Kong handles all the external client facing routing and policies while the
  mesh takes care of these aspects for internal service-to-service traffic.
- **`ingress.kubernetes.io/service-upstream`**, a new annotation has
  been introduced.
  Adding this annotation to a kubernetes service resource
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
  stop users from mis-configuring the Ingress controller.
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

### Fixed

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

## [0.5.0] - 2019/06/25

#### Summary

This release introduces automated TLS certificates, consumer-level plugins,
enabling deployments using controller and Kong's Admin API at the same time
and numerous bug-fixes and enhancements.

#### Breaking changes

- UUID of consumers in Kong are no longer associated with UID of KongConsumer
  custom resource.

#### Added

- Kong 1.2 is now supported, meaning wildcard hosts in TLS section of Ingress
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

## [0.4.0] - 2019/04/24

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
  (deprecatd in 0.2.0) is no longer supported.
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
- New fields are added to KongIngress CRD to support HTTPS Active healthchecks.
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

## [0.3.0] - 2019/01/08

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


## [0.2.2] - 2018/11/09

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
 - Fix a nil pointer reference when overiding Ingress resource with KongIngress
   [#188](https://github.com/Kong/kubernetes-ingress-controller/pull/188)


## [0.1.3] - 2018/11/09

#### Fixed

 - Fix path-only based Ingress rule parsing and configuration where only a
   path based rule for a Kubernetes Service
   would not setup Routes and Service in Kong.
   [#190](https://github.com/Kong/kubernetes-ingress-controller/pull/190)
 - Fix plugin config comparison logic to avoid unnecessary PATCH requests
   to Kong
   [#196](https://github.com/Kong/kubernetes-ingress-controller/pull/196)


## [0.2.1] - 2018/10/26

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


## [0.1.2] - 2018/10/26

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
   reconcillation loop, which lead to unnecessary Router rebuilds inside Kong.
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


## [0.1.1] - 2018/09/26

#### Fixed

 - Fix version parsing for minor releases of Kong Enterprise (like 0.33-1).
   The dash(`-`) didn't go well with the semver parsing
   [#141](https://github.com/Kong/kubernetes-ingress-controller/pull/141)

## [0.2.0] - 2018/09/21

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
 - The custom resource definitions now have a shortname for all the
   CRDs, making it easy to interract with `kubectl`.
   [#120](https://github.com/Kong/kubernetes-ingress-controller/pull/120)

#### Fixed

 - Avoid issuing unnecessary PATCH requests on Services in Kong during the
   reconcillation loop, which lead to unnecessary Router rebuilds inside Kong.
   [#107](https://github.com/Kong/kubernetes-ingress-controller/pull/107)
 - Fixed the diffing logic for plugin configuration between KongPlugin
   resource in k8s and plugin config in Kong to avoid false positives.
   [#106](https://github.com/Kong/kubernetes-ingress-controller/pull/106)
 - Correctly format IPv6 address for Targets in Kong.
   Thanks @NixM0nk3y for the patch!
   [#118](https://github.com/Kong/kubernetes-ingress-controller/pull/118)


## [0.1.0] - 2018/08/17

#### Breaking Changes

 - :warning: **Declarative Consumers in Kong** Kong consumers can be
   declaritively configured via `KongConsumer` custom resources. Any consumers
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


## [v0.0.5] - 2018/06/02

#### Added

 - Add support for Kong Enterprise Edition 0.32 and above

## [v0.0.4] and prior

 - The initial versions rapidly were iterated delivering
   a working ingress controller.

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
