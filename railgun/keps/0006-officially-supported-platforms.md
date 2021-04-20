---
status: provisional
---

# Official Support for new Platforms and Integrations with KIC

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
    - [Story 1](#story-1)
    - [Story 2](#story-2)
<!-- /toc -->

## Summary

Provide the criteria for contributors to request and add new "official supported" platforms and integrations for the [Kong Kubernetes Ingress Controller (KIC)][kic].

[kic]:https://github.com/kong/kubernetes-ingress-controller

## Motivation

- we want to document for contributors and end-users the criteria for any platform or integration to be "officially supported"
- we want to demonstrate our documentation by adding new support for an existing feature request using the newly established criteria

### Goals

- establish preliminary acceptance criteria, including considerations about impact and maintenance burden
- establish testing and infrastructure criteria
- document the timeframes relevant to any official support, with clear end of life (EOL) criteria
- provide walkthrough documentation that covers beginning to end
- provide a repository issue template for this purpose which links to the relevant documentation
- exercise and demonstrate all of the above by implementing the [previously requested ARM64 support][issues-451]

[arm64]:https://en.wikipedia.org/wiki/ARM64
[upstream]:https://github.com/kong/kong
[issues-451]:https://github.com/Kong/kubernetes-ingress-controller/issues/451

## Proposal

### User Stories

#### Story 1

As an end-user and contributor I want ARM64 builds of the KIC available so that I can deploy on my ARM64 based Kubernetes clusters.

(real request: https://github.com/Kong/kubernetes-ingress-controller/issues/451)

#### Story 2

As a maintainer of the KIC considering whether to accept new platforms or integrations I want integration tests that provide protections against regressions.
