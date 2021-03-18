---
title: KIC Kubebuilder Rearchitecture
status: implementable
---

# Kong Kubernetes Ingress Controller (KIC) Re-architecture using Kubebuilder

<!-- toc -->
- [Release Signoff Checklist](#release-signoff-checklist)
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
    - [Story 1](#story-1)
    - [Story 2](#story-2)
- [Design Details](#design-details)
  - [Test Plan](#test-plan)
  - [Graduation Criteria](#graduation-criteria)
- [Implementation History](#implementation-history)
<!-- /toc -->

## Summary

Historically the [Kong Kubernetes Ingress Controller (KIC)][kic] was built on older Kubernetes controller design patterns and maintained a lot of functionality for the runtime of the controller. We will modernize the KIC by re-architecting it using the [Kubebuilder SDK][kb] provided by [Kubernetes Special Interest Groups (SIGs)][sig] and make more extensive use of the [controller-runtime][cr] and other upstream libraries to maintain [controllers][ctrl].

[kic]:https://github.com/kong/kubernetes-ingress-controller
[kb]:https://kubebuilder.io/
[sig]:https://github.com/kubernetes-sigs/kubebuilder/
[cr]:https://github.com/kubernetes-sigs/controller-runtime
[ctrl]:https://kubernetes.io/docs/concepts/architecture/controller/

## Motivation

- increase the speed and efficiency for contributors adding new APIs
- decrease the complexity of interwoven dependencies between disparate APIs
- improve statuses and events associated with our APIs and improve the logging within their controller(s)
- provide separate controllers for each upstream API which we support (e.g. `netv1.Ingress`, `netv1beta1.Ingress`, e.t.c.)
- provide automation and tooling to generate and maintain multiple API implementations
- provide automation and tooling to maintain multiple API versions (e.g. `v1alpha1`, `v1beta1`, `v1`, e.t.c.)
- provide feature gates which limits access to newer and experimental Kubernetes APIs and controllers
- provide documentation for feature gates which communicates API maturity and longevity to end-users
- prepare ourselves for an eventual pivot to [Gateway APIs][gw]

[gw]:https://gateway-api.sigs.k8s.io/

### Goals

- re-architect KIC onto Kubebuilder and modern `controller-runtime`
- break the existing functionality related to each of our Kubernetes APIs (e.g. `KongIngress`, `TCPIngress`, e.t.c.) into their own [controller implementation][impl]
- produce `UDPIngress` to provide a starting point example of how to build and maintain APIs going forward

[impl]:https://kubebuilder.io/cronjob-tutorial/controller-overview.html

## Proposal

### User Stories

#### Story 1

As a maintainer of KIC I want functionality concerns regarding our APIs clearly delimited between API boundaries and their controllers for simpler maintainance.

#### Story 2

As an operator of the KIC I want logging to be clearly deliniated by the responsible APIs/Controllers for improved transparency and readability.

#### Story 3

As a contributor to the KIC I want to be able to quickly contribute new ideas and experimental features without making them immediately available in upcoming releases.

#### Story 4

As a user of KIC, I want to be able to inspect the intermediate objects produced by KIC (collected KongState, generated decK config) for debugging purposes, as inspired by [this review comment](https://github.com/Kong/kubernetes-ingress-controller/pull/991#pullrequestreview-570627606).

## Implementation History

- spike and prototype started to re-architect the KIC and move it to modern tooling: https://github.com/kong/railgun
- merging the railgun prototype (KIC on [controller-runtime][cr] v0.7.x+ and [kubebuilder][kb]) into upstream [KIC][kic] (PR [KIC#1032][rgpr])
- developed a prototype [Kubernetes Testing Framework (KTF)][ktf] for KIC now used in railgun
- KIC 2.0 signed off by product for prioritization
- [KIC 2.0 Milestone established][ms12]
- [UDPIngress][udp] supported added to `railgun/` POC and demoed

[cr]:https://github.com/kubernetes-sigs/controller-runtime
[kb]:https://github.com/kubernetes-sigs/kubebuilder
[kic]:https://github.com/kong/kubernetes-ingress-controller
[rgpr]:https://github.com/Kong/kubernetes-ingress-controller/pull/1032
[ktf]:https://github.com/kong/kubernetes-testing-framework
[ms12]:https://github.com/Kong/kubernetes-ingress-controller/milestone/12
[udp]:https://github.com/Kong/kubernetes-ingress-controller/milestone/14
