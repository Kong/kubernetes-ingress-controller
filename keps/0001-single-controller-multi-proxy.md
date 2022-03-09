---
title: Single Controller Multi Proxy DBLESS
status: provisional
---

# Decoupled KIC Components

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
    - [Story 1](#story-1)
- [Alternatives](#alternatives)
<!-- /toc -->

## Summary

Historically the [Kong Kubernetes Ingress Controller (KIC)][kic] has [required
a controller pod to be present per proxy pod in DBLESS mode][kic702] and
generally had all components coupled tightly together in a single `Pod`. This
KEP aims to make it possible to deploy a single controller to manage any number
of proxy instances in DBLESS mode.

[kic]:https://github.com/kong/kubernetes-ingress-controller
[kic702]:https://github.com/Kong/kubernetes-ingress-controller/issues/702

## Motivation

- reduce Kubernetes API overhead in high scale Kong environments where large
  numbers of proxies are needed in DBLESS mode
- better support an operator use case based on the `Gateway` API to enable a
  single KIC to manage multiple `Gateways`.

[gateway-api]:https://kubernetes-sigs.github.io/gateway-api/

### Goals

- enable TLS and authentication for DBLESS `proxy` pods
- enable multiple `proxy` instances to be managed by a single KIC
- provide service discovery and automatic configuration for new `proxy` pods

## Proposal

### User Stories

#### Story 1

As an operator with DBLESS deployed Kong I want a single controller to be able
to manage my proxy instances so that scaling out my gateway setup doesn't
linearly increase the Kubernetes Admin API load and I can safely allow KIC to talk
 to the Admin API between pods.

## Alternatives

**TODO**: Now with the advent of the [Gateway API][gwm1] we need to consider
          whether we want to invest in adding multi-proxy support in the
          traditional KIC ecosystem with some kind of independent service
          discovery mechanism, or whether we want to simply (and only) implement
          support on top of the already existing `Gateway` support.

[gwm1]:https://github.com/Kong/kubernetes-ingress-controller/milestone/21
