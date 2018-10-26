# Table of Contents

 - [0.2.1](#021-20181026)
 - [0.1.1](#011-20180926)
 - [0.2.0](#020-20180921)
 - [0.1.0](#010-20180817)
 - [0.0.5](#005---20180602)
 - [0.0.4 and prior](#004-and-prior)

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

[0.2.1]: https://github.com/kong/kubernetes-ingress-controller/compare/0.2.0...0.2.1
[0.1.1]: https://github.com/kong/kubernetes-ingress-controller/compare/0.1.0...0.1.1
[0.2.0]: https://github.com/kong/kubernetes-ingress-controller/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/kong/kubernetes-ingress-controller/compare/v0.0.5...0.1.0
[v0.0.5]: https://github.com/kong/kubernetes-ingress-controller/compare/v0.0.4...v0.0.5
[v0.0.4]: https://github.com/kong/kubernetes-ingress-controller/compare/7866a27f268c32c5618fba546da2c73ba74d4a46...v0.0.4
