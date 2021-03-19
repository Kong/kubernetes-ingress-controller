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

[Gateway API][gateway] is the premier upstream Kubernetes API aiming to shape the future of ingress for Kubernetes clusters. This KEP aims to align the [Kong Kubernetes Ingress Controller (KIC)][kic] with [upstream direction][gapi] of [Kubernetes SIG Networking][signet] and provide an initial implementation of the [latest specification][gapi-releases].

[gateway]:https://kubernetes-sigs.github.io/gateway-api/api-overview/
[kic]:https://github.com/kong/kubernetes-ingress-controller
[gapi]:https://kubernetes-sigs.github.io/gateway-api/
[signet]:https://github.com/kubernetes/community/tree/master/sig-network
[gapi-releases]:https://github.com/kubernetes-sigs/gateway-api/releases

## Motivation

- stay up to date (and when possible ahead) of competing solutions in terms of Kubernetes functionality and UX
- conform to Kubernetes upstream implementations for gateway services to align ourselves and grow our influence in upstream direction
- make it possible to use and manage Kong in an entirely Kubernetes native way, lowering the barrier to entry for human operators

### Goals

- develop an initial `Gateway` implementation against `Gateway API v0.2.x`
- demo our prototype implementation internally and in an upcoming [SIG Networking Meeting][cal]

[cal]:https://kubernetes-sigs.github.io/gateway-api/community/
