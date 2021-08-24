---
title: Kong Gateway API
status: provisional
---

# Kong Gateway API

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
<!-- /toc -->

## Summary

[Gateway API][gateway] is the successor the the [Ingress][ingress] API upon which the [Kong Kubernetes Ingress Controller (KIC)][kic] was founded
and includes a variety of improvements over its predecessor for feature richness, protocol support, lifecycle management, automation and more. We will
implement Gateway APIs for the KIC in order to improve the operations, automation and integration characteristics of using the Kong Gateway on Kubernetes.

[gateway]:https://kubernetes-sigs.github.io/gateway-api
[kic]:https://github.com/kong/kubernetes-ingress-controller

## Motivation

- enable multi-gateway single KIC deployments
- enable lifecycle management of Kong Gateways on Kubernetes
- adhere to upstream standards to create a low barrier to entry for end-users

### Goals

- integrate our existing controller manager with `GatewayClass`
- develop a `Gateway` controller which is responsible for multi-proxy deployments and Kong Gateway lifecycle management
- expose as much Kong Gateway functionality as possible behind upstream standard APIs

### Non-Goals

- we wont retroactively support the `v0.1.x` Gateway API specification

## Proposal

Similar to how the KIC has historically had support for the [Ingress][ingress] resource from `networking.k8s.io/v1` the intention is to start conforming to the new Gateway API resources including (but not limited to):

- [Gateway][gateway] - the actual Gateway which routes traffic, in our case the [Kong Gateway][kong]
- [GatewayClass][gateway-class] - the controller which is responsible for a `Gateway`
- [HTTPRoute][httproute] - the quintessential HTTP traffic routing mechanism, with parallels to the `Ingress` API

For the purposes of this KEP we will be focused on implementing the [v1alpha2][v1alpha2] specification of Gateway APIs as it is (at the time of writing) the latest release and current iteration.

[ingress]:https://kubernetes.io/docs/concepts/services-networking/ingress/
[gateway]:https://gateway-api.sigs.k8s.io/v1alpha2/api-types/gateway/
[kong]:https://github.com/kong/kong
[gateway-class]:https://gateway-api.sigs.k8s.io/v1alpha2/api-types/gatewayclass/
[httproute]:https://gateway-api.sigs.k8s.io/v1alpha2/api-types/httproute/

### User Stories

#### Story 1

As a Kubernetes operator I want to use standard Kubernetes APIs to define ingress rules for communication with my services from outside the cluster.

#### Story 2

As a Kubernetes operator with an already existing gateway for ingress I want transitioning to Kong to be as seamless and require as minimal changes
to my existing deployments as possible.

#### Story 3

As a Kubernetes operator I want automated lifecycle management of Kong Gateways deployed to my clusters, rather than having to manage them by hand using Helm.

#### Story 4

As a Kubernetes operator I want the KIC to automate the deployment and manage the lifecycle of multiple Kong Gateways when my use case is high enough scale that a single gateway is not sufficient.

## Design Details

Similar to all other APIs we've implemented previously (such as `Ingress`, `TCPRoute`, `UDPRoute`, e.t.c.) we will need to create controllers
for the new resources:

- `internal/controllers/gateway/gatewayclass_controller.go`
- `internal/controllers/gateway/gateway_controller.go`
- `internal/controllers/gateway/httproute_controller.go`

These controllers will be responsible for their respective `GatewayClass`, `Gateway` and `HTTPRoute` types from `gateway.networking.k8s.io/v1alpha2`.

Each of these types will need handlers in the backend `Proxy` implementation so that Kong `Services`, `Routes`, e.t.c. will be generated for them
and posted to the Kong Admin API's `/config` endpoint, as is true of all other types currently.

### Controller Implementations

At the core we will need to support the following APIs:

- `GatewayClass`
- `Gateway`
- `HTTPRoute`

These will give us a meaningful implementation, though we will need to add support for other types in time.

#### GatewayClass Controller

The `GatewayClass` API is a reference to the controller that is managing Gateways of this class. Our controller
will provide a tag (similar to what we currently call `ingress-class`) to indicate whether the controller is
responsible for a `GatewayClass`:

```yaml
kind: GatewayClass
apiVersion: gateway.networking.k8s.io/v1alpha2
metadata:
  name: acme-lb
spec:
  controller: konghq.com/gateway-controller
```

This will be instrumented by the controller manager as a flag, though the above will be the default case when
the flag isn't set:

```shell
manager --gateway-class konghq.com/gateway-controller
```

##### Additional Considerations

- It's up to the operator whether they want a single or multiple controllers responsible for their `GatewayClasses`, for some high end deployments an entirely separate manager can be spun up but will need to have a distinct `--gateway-class`
- If someone mutates the specification for a `GatewayClass` such as to drop the `controller: <tag>` which enables our support of it, the controller will stop managing related resources and will clean out the Kong Gateway's relevant configurations

#### Gateway Controller

Where the `GatewayClass` API maps to our controller manager, the `Gateway` API maps to the Kong Gateway (Proxy Server).

For each `Gateway` resource that is provided and configured for the relevant `GatewayClass` a separate proxy server will be initialized:

```yaml
kind: GatewayClass
apiVersion: gateway.networking.k8s.io/v1alpha2
metadata:
  name: default-match-example
spec:
  controller: acme.io/gateway-controller
---
kind: Gateway
apiVersion: gateway.networking.k8s.io/v1alpha2
metadata:
  name: project-1-ingress
spec:
  gatewayClassName: default-match-example
  listeners:
  - name: http
    protocol: HTTP
    port: 80
---
kind: Gateway
apiVersion: gateway.networking.k8s.io/v1alpha2
metadata:
  name: project-2-ingress
spec:
  gatewayClassName: default-match-example
  listeners:
  - name: http
    protocol: HTTP
    port: 80
```

In the above example there will be a separate Kong proxy container present for both `project-1-ingress` and `project-2-ingress`.

##### Additional Considerations

- https://github.com/Kong/kubernetes-ingress-controller/issues/702 is related to our single controller multi-proxy lifecycle management concerns
- for legacy reasons we may need to include a single proxy deployment mechanism where one proxy server can host for multiple gateways, as this is how the KIC has historically operated (see the current Helm chart)

#### HTTPRoute Controller

The HTTPRoute controller is fairly straightforward, which backend proxy is responsible is based on which Gateway its attached to in the `parentRefs`, but then the resource otherwise
gets parsed, converted into Kong types, and posted to `/config` like any other API resource that currently exists.

### Proxy Cache Implementation

While the controllers will need to handle `GatewayClass` and `Gateway` to determine which `HTTPRoutes` we need to support, once we need to support them
the backend caching, parsing, and Kong API updating logic will only need to be aware of `Gateway` and `HTTPRoute`.

Updates between `Gateway` objects and the `HTTPRoutes` which are attached to them will be synchronized and are tightly coupled in terms of the backend
Kong configurations they will produce as `Gateway` is responsible for the TLS configurations for listeners separately from the `HTTPRoute`, which with the
historical `Ingress` API this was all managed within that single resource.

We'll need to update `internal/store/store.go` to cache the objects that are provided from the aforementioned controllers via `proxy.UpdateObject(obj)`:

```go
type Storer interface {
    // ... existing methods
    ListGateways() ([]*gatewayv1alpha2.Gateway, error)
    ListHTTPRoutes() ([]*gatewayv1alpha2.HTTPRoute, error)
}

type CacheStores struct {
    // ... existing stores
    Gateway      cache.Store
    HTTPRoute    cache.Store
}
```

We will need to ensure all the relevant plumbing to cache the Gateway API objects therein is in place so that they can be fed to `internal/parser/`:

```go
func parseAll(log logrus.FieldLogger, s store.Storer) ingressRules {
    // ... existing parsing logic
    gateways, err := s.ListGateways()
    parsedGateways := fromGatewayAPIsV1alpha2Gateway(log, gateways)
    httproutes, err := s.ListHTTPRoutes()
    parsedHTTPRoutes := fromGatewayAPIsV1alpha2HTTPRoute(log, httproutes)

    return mergeIngressRules(otherTypes..., parsedGateways, parsedHTTPRoutes)
```

The parsing logic above (e.g. `fromGatewayAPIsV1alpha2HTTPRoute()`) is the ultimate low-level logic we will need to compile `Gateways` and `HTTPRoutes`
into `kong.Services` and `kong.Routes`. Practically speaking almost all the prior art needed to implement this is covered by the API types we already
support, so the work should be predominantly integration work and little to no exploratory work needed.

### Test Plan

In order to consider this KEP `implemented` we must have the following present:

- `examples/` includes manifests for (at least) a working HTTP ingress with `GatewayClass`, `Gateway`, and `HTTPRoute` using Kong
- KIC [docs.konghq.com][docs] must include a `guide/` for Gateway APIs (flagged **experimental** until Gateway APIs are GA)
- `test/integration/` includes integration tests which cover deployment and configuration options of `HTTPRoute`, including TLS configurations
- our `.github/workflows/release.yaml` release CI must test our `HTTPRoute` implementation against a historical list of the Kubernetes versions we support

[docs]:https://docs.konghq.com/kubernetes-ingress-controller/

## Production Readiness

Production readiness of this feature is marked by the following requirements:

- Gateway APIs themselves have reached a GA status in upstream Kubernetes
- A version of Kubernetes which includes the Gateway APIs standard becomes GA
- We have support for most if not all of the [entire Gateway APIs spec][gateway-spec]
- If there's any features or implementations for Gateway APIs we explicitly choose not to support, we document that and the reasoning
- Our integration and E2E testing provides feature cover and strong regression protections

[gateway-spec]:https://gateway-api.sigs.k8s.io/v1alpha2/references/spec/

### Feature Enablement and Rollback

In order to significantly limit the exposure of our Gateway APIs implementations during the alpha/beta periods of their release
we will not provide the standard `--enable-controller-<name>` flags for the manager as these are [prominently documented][flags-docs]
and otherwise easily visible when running `manager --help` which ostensibly gives them a sense that they are a normal part of any included release.

We will provide a [Feature Gate][feature-gates] option for enabling Gateway APIs which will [conform to upstream standards][gate-alpha-beta] as to how
alpha and beta features can be enabled, which in turn will be strongly familiar way to manage these kinds of features to Kubernetes cluster operators.

Consequently we will need to add [documentation][docs-refs] a table of our feature gates statuses as upstream does for their controllers.

###### How can this feature be enabled / disabled in a live cluster?

Can be enabled with:

```shell
--feature-gates="...,GatewayAPIs=true"
```

Disabling can be a matter of not setting a value on the feature gate at all, or being explicitly if desired for documentative reasons:

```shell
--feature-gates="...,GatewayAPIs=false"
```

[flags-docs]:https://docs.konghq.com/kubernetes-ingress-controller/latest/references/cli-arguments/
[features-gates]:https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
[gate-alpha-beta]:https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/#feature-gates-for-alpha-or-beta-features
[docs-refs]:https://docs.konghq.com/kubernetes-ingress-controller/latest/references

## Drawbacks

Given the direction of the upstream Kubernetes community which appears to be conforming around Gateway,
the only drawback we've seen is the time cost we will need to pay in order to implement the API during
a time when the API isn't GA.

## Alternatives

The primary alternative we've considered is that we could continue to support the featuresets which Gateway APIs is bound for using our custom APIs,
however this will make us an outlier and less attractive an option to operators as many competing implementations already have and are planning long-term support of Gateway APIs.

## Infrastructure Needed

We can rely on the existing infrastructure available currently which tests `Ingress` and our CRD based APIs, no new infrastructure will be required.

