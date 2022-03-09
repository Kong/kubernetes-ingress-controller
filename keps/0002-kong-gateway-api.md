---
title: Kong Gateway API
status: implementable
---

# Kong Gateway API

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
- [Design Details](#design-details)
  - [Test Plan](#test-plan)
  - [Graduation Criteria](#graduation-criteria)
- [Production Readiness](#production-readiness)
  - [Feature Enablement and Rollback](#feature-enablement-and-rollback)
- [Drawbacks](#drawbacks)
- [Alternatives](#alternatives)
- [Infrastructure Needed](#infrastructure-needed)
<!-- /toc -->

## Summary

[Gateway API][gateway] is the [successor][ingv2] the the [Ingress][ingress] API
upon which the [Kong Kubernetes Ingress Controller (KIC)][kic] was founded and
includes a variety of improvements over its predecessor for feature richness,
protocol support, lifecycle management, automation and more. We will implement
Gateway APIs for the KIC in order to improve the operations, automation and
integration characteristics of using the Kong Gateway on Kubernetes.

[gateway]:https://kubernetes-sigs.github.io/gateway-api
[ingv2]:https://www.youtube.com/watch?v=Ne9UJL6irXY
[ingress]:https://kubernetes.io/docs/concepts/services-networking/ingress/
[kic]:https://github.com/kong/kubernetes-ingress-controller

## Motivation

- adhere to upstream standards to create a low barrier to entry for end-users

### Goals

- support `GatewayClass` and `Gateway` APIs
- support `HTTPRoute`, `TCPRoute`, `UDPRoute` and `TLSRoute` APIs

### Non-Goals

- we wont retroactively support the `v1alpha1` Gateway API specification
- we wont support the full lifecycle of `Gateway` resources

## Proposal

The KIC will include limited support for `Gateway` resources, and full
"core" level support (as defined in upstream) for L7 and L4 routes.

Defining "limited Gateway support": One notable connotation of `Gateway`
resources is the idea that an end-user might create a `Gateway` and this
results in an underlying `Deployment` and a gateway server (e.g. [Kong
Gateway][kong]) being deployed and the lifecycle of that `Gateway` is managed
by the controller. We'll call this "managed" mode and conversely when a gateway
already exists on the cluster and the `Gateway` resource is more or less just
metadata about that existing resource that will be referred to as "unmanaged"
mode. We will only be supporting unmanaged mode in the KIC, see the
alternatives section for more information about managed mode considerations.

[kong]:https://github.com/kong/kong

### User Stories

#### Story 1

As a Kubernetes operator I want to use standard Kubernetes APIs to define
ingress rules for communication with my services from outside the cluster.

#### Story 2

As a **developer** I want to use **standard Kubernetes APIs** so that my
manifests and deployments don't require (at least minimize) domain specific
APIs.

## Design Details

The implementation of `Gateway` support in KIC will not require any new CRDs
and will only require controllers which will configure routing APIs like
`HTTPRoute` in the data-plane.

At the core we will support the following APIs:

- `GatewayClass`
- `Gateway`
- `HTTPRoute`
- `TCPRoute`
- `UDPRoute`
- `TLSRoute`

We will add support for these resources and controllers to provision them.

### Validating Webhook

For the first iteration of our Gateway API support we will have some limits on
functionality that will be codified in validation:

- `Gateway` limitations (e.g. disallowing managed mode gateways)
- `HTTPRoute` limitations (e.g. query params matching)

Otherwise the basic validation for these API types is to be provided by the
[upstream validating webhook][gwhook].

[gwhook]:https://github.com/kubernetes-sigs/gateway-api/tree/master/cmd/admission

### Controllers

#### Gateway Controller

In our implementation the `Gateway` resource will be a reflection of the
Kubernetes `Service` object for the Kong Gateway, historically defined with:

```console
$ manager --publish-service namespace/name
```

Rather than the end-user directly defining the `Addresses` and `Listeners` of
the `Gateway` resource, they will be derived from the publish service and any
updates to the `Service` will cause reconcilation to update the `Gateway`
according to the `Service` spec and status.

This operational mode is called "unmanaged" mode and needs to be configured on
`Gateway` resources with an annotation:

```yaml
kind: Gateway
apiVersion: gateway.networking.k8s.io/v1alpha2
metadata:
  annotations:
    konghq.com/gateway-unmanaged: "true"
  name: project-1-ingress
spec:
  gatewayClassName: default-match-example
  listeners:
  - name: http
    protocol: HTTP
    port: 80
```

The KIC Gateway controller will be responsible for watching the Kubernetes
`Service` belonging to the Kong Gateway `Deployment` and derive the `Gateway`
listeners, ports, and protocols (which routes like `HTTPRoute` need to attach
to) from that `Service` object.

If the `Service` object is updated to include new available ports and/or
protocols (e.g. a human operator performs a `helm upgrade` which changes the
Kong Gateway configuration) the KIC Gateway Controller will reconfigure the
gateways listener specification accordingly.

Direct manipulation of the `Gateway` object in this mode will have no effect:
the KIC Gateway Controller will always update the resource to be derivative of
the `Service` object.

[chart]:https://github.com/kong/charts
[kong-admin-api]:https://docs.konghq.com/gateway-oss/latest/admin-api/

#### Routing Controller

For the `HTTPRoute`, `TCPRoute`, `UDPRoute` and `TLSRoute` APIs the controller
implementation will be very similar to the parallel `Ingress`, `TCPIngress` and
`UDPIngress` which already exist and will simply be a matter of translating
routing rules into Kong `Routes`, `Services` and `Upstreams`.

### Test Plan

Automated testing will be provided for the `main` branch and all PRs:

- unit tests for all Go packages using [go tests][gotest]
- integration tests using the [Kong Kubernetes Testing Framework (KTF)][ktf]
- e2e tests using [KTF][ktf]

All tests should be able to run locally using `go test` and integration and e2e
tests using a local system Kubernetes deployment like [Kubernetes in Docker
(KIND)][kind].

[gotest]:https://pkg.go.dev/testing
[ktf]:https://github.com/kong/kubernetes-testing-framework
[kind]:https://github.com/kubernetes-sigs/kind

### Graduation Criteria

The following highlights milestones along the path of Gateway support maturity.

The milestones correlate directly with [Github Milestones][ms] in the same
repository.

[ms]:https://github.com/Kong/kubernetes-ingress-controller/milestones

#### Milestone 1 - Initial L7 Support

- [x] our validating webhook for Gateway resources has been added and provides
      errors for unimplemented features
- [x] a `Gateway` controller implementation with basic support for "unmanaged"
      operational mode is introduced
- [x] an initial implementation of `HTTPRoute` is added which includes support
      for the most important options for basic usage

#### Milestone 2 - Core L7 Support

- [ ] all `HTTPRoute` features marked as "core" support have been added
- [ ] upstream [conformance tests][cft] have been added to our CI and pass
- [ ] optional: add as many "extended" and "custom" features as can already be
      supported by the [Kong Gateway][kong].

[cft]:https://github.com/kubernetes-sigs/gateway-api/tree/master/conformance
[kong]:https://github.com/kong/kong

#### Milestone 3 - Core L4 Support

- [ ] `TCPRoute` support is added with all "core" features supported
- [ ] `UDPRoute` support is added with all "core" features supported
- [ ] `TLSRoute` support is added with all "core" features supported
- [ ] upstream conformance tests are passing for all three APIs
- [ ] optional: add as many "extended" and "custom" features as can already be
      supported by the [Kong Gateway][kong].

## Production Readiness

Production readiness of this feature is marked by the following requirements:

- All milestones of the above `Graduation Criteria` have been completed
- Gateway APIs themselves have reached a GA status in upstream Kubernetes
- A version of Kubernetes which includes the Gateway APIs standard becomes GA
- Gateway API documentation is added to [our documentation][kongdocs]

[gateway-spec]:https://gateway-api.sigs.k8s.io/v1alpha2/references/spec/
[kongdocs]:https://docs.konghq.com

### Feature Enablement and Rollback

The Gateway API support will be disabled by default prior to GA and enabled by
using `--feature-gates=Gateway=true`.

## Drawbacks

Given the direction of the upstream Kubernetes community which appears to be
conforming around Gateway, the only drawback we've seen is the time cost we
will need to pay in order to implement the API during a time when the API isn't
GA. As such for the initial iteration we've decided to hold back on some of the
provisioning and lifecycle management features inherent to the Gateway APIs
project (e.g. provisioning `Gateway` resources by deploying and maintaining new
proxy `Deployments` for them) in order to spread time costs out.

## Alternatives

### Not Supporting Gateway

The main alternative we've considered is that we could continue to support
the featuresets which Gateway APIs is bound for using our custom APIs (e.g.
`TCPIngress`, `UDPIngress`, e.t.c.), however this will make us an outlier and
less attractive an option to operators as many competing implementations
already have and are planning long-term support of Gateway APIs.

### Managed Gateways

We considered supporting managed `Gateways` as a part of the KIC but this
lends itself better to an operator which works at a higher level: for now
we're going to assess managed gateway support separately possibly as part
of our [operator][kongop] project instead.

[kongop]:https://github.com/kong/kong-operator
