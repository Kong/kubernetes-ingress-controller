---
title: KIC Kubebuilder Rearchitecture
status: implementable
---

# Notes

For reference see the milestones related to this proposal to check the progress of related efforts:

- Kong Kubernetes Testing Framework (KTF) `v0.1.0` - https://github.com/Kong/kubernetes-testing-framework/milestone/1
- KIC 2.0 Alpha - https://github.com/Kong/kubernetes-ingress-controller/milestone/12
- KIC 2.0 Testing Renaissance - https://github.com/Kong/kubernetes-ingress-controller/milestone/16
- KIC 2.0 Milestone - https://github.com/Kong/kubernetes-ingress-controller/milestone/15

# Kong Kubernetes Ingress Controller (KIC) Re-architecture using Kubebuilder

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
    - [Story 1](#story-1)
    - [Story 2](#story-2)
- [Design Details](#design-details)
  - [Test Plan](#test-plan)
  - [Graduation Criteria](#graduation-criteria)
- [Implementation History](#implementation-history)
- [Alternatives](#alternatives)
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

### Non-Goals

#### Configuration Monolith Re-architecture

When we first started experimenting and prototyping for this KEP we were motivated to deconstruct the existing monolithic controller such that individual controllers for APIs could be separate microservices and autonomous. Due to limitations of maintainer capacity and some desires for upstream changes that would not be able to occur logistically fast enough to support us (namely having a single upstream API to develop against rather than a separate API for DB vs DBLESS Kong instances) we've pulled this work out of scope. We are however still motivated to make this change and continue our re-architecture, just consider it out of scope for this KEP and it will become the subject of its own.

#### Kubebuilder Controller Management

Despite the motivation present in this KEP to automate some of our controller management, logistics and time constraints have led us to keeping the conversion our existing controller to `kubebuilder` managed controllers _out of scope_ for this KEP. For this iteration we will focus on using the API, CRD, and configuration management features of `kubebuilder`, but the controller management features will be considered as part of a later iteration to reduce the number of changes we make with a single release (we will however still convert to controller runtime and ultimately use the kubebuilder provided controller machinery to replace our historical machinery).

#### New APIs

We reference [UDPIngress][udpingress] in the implementation history (below) as it was used for demonstration, but completing new features and APIs is not in scope for this KEP, though the result of this KEP is that newer features are intended to be easier to contribute.

We can use the improvements made here to _demonstrate_ the ease of adding new features on the new architecture, but the full scope of GA for new features will need to be a follow up KEP.

[udpingress]:https://github.com/Kong/kubernetes-ingress-controller/milestone/14

## Proposal

The historical releases of the KIC (which we will refer to as `pre-v2`) were built on an older controller architecture forked from the [NGinx Ingress Controller][nginx-ingress-controller] some years prior.

This legacy served us well for the years to come, but at the point where this KEP was written (early 2021) it was becoming noticeably harder to continue maintaining and adding new features to the KIC as it fundamentally hadn't grown alongside much of the rest of the Kubernetes open source community.

Since the inception of KIC new Software Development Kits (SDK) have been created to support building and maintaining Kubernetes controllers:

- [Kubebuilder][kb]
- [OperatorSDK][osdk]

These SDKs simplify, automate and ultimately generate some of the code we had been historically maintaining ourselves including (but not limited to):

- API schemas
- Controller reconcilation machinery
- Custom Resource Defition (CRD) management
- Kustomize configurations
- RBAC security
- controller-manager CLI (and flags)

In short, using a Kubernetes SDK gets most of the actual machinery and scaffolding needed to start writing our explicit API reconcilation logic for "free" (paid for by the last few years of community contributions, which we are extremely grateful for).

This enhancement is about re-architecting the KIC onto [Kubebuilder][kb] (for reasons why we did not choose [Redhat's OperatorSDK][osdk] see the [alternatives section below](/#alternatives) and as a consequence putting ourselves on a modern version of [Kubernetes Controller Runtime][cr] with a multitude of new enhancements and features.

The result will be a large portion of our KIC maintainence is automated (and possibly for some things even automated via CI) making it easier and faster to contribute to the project so that we can focus harder on fixes and enhancements.

[nginx-ingress-controller]:https://docs.nginx.com/nginx-ingress-controller/
[kb]:https://kubebuilder.io
[osdk]:https://sdk.operatorframework.io/
[cr]:https://github.com/kubernetes-sigs/controller-runtime

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
- [Established KIC 2.0 Preview release criteria][ms15]
- KTF fully separated into it's [own repo][ktf]
- integration tests [added][legacy-tests] to test `v1.x` and railgun controllers on every PR from now until release
- First alpha release objectives defined in milestone: https://github.com/Kong/kubernetes-ingress-controller/milestone/15

[cr]:https://github.com/kubernetes-sigs/controller-runtime
[kb]:https://github.com/kubernetes-sigs/kubebuilder
[kic]:https://github.com/kong/kubernetes-ingress-controller
[rgpr]:https://github.com/Kong/kubernetes-ingress-controller/pull/1032
[ktf]:https://github.com/kong/kubernetes-testing-framework
[ms12]:https://github.com/Kong/kubernetes-ingress-controller/milestone/12
[udp]:https://github.com/Kong/kubernetes-ingress-controller/milestone/14
[ms15]:https://github.com/Kong/kubernetes-ingress-controller/milestone/15
[legacy-tests]:https://github.com/Kong/kubernetes-ingress-controller/issues/1040

## Alternatives

### CRD/Secret vs. In-Memory Cache

To help break apart the monolithic controller from 1.x we considered using a `Secret` or `CRD` as the caching location for resources as an interim solution between previous arch and the future arch we wanted to define, however the timing and logistics of that simply didn't work for this KEPs scope and limitations such as maximum object size for `Secrets` led us to stop on this and save it for a later iteration.

### OperatorSDK vs. Kubebuilder

The [OperatorSDK][osdk] from [Redhat][rhel] was considered for our new Kubernetes SDK, but ultimately decided against due to lack of familiarity and preferring a more generic and flexible toolkit.
