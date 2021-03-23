---
status: implementable
---

# Kong Kubernetes Testing Framework (KTF)

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories (Optional)](#user-stories-optional)
    - [Story 1](#story-1)
    - [Story 2](#story-2)
- [Implementation History](#implementation-history)
<!-- /toc -->

## Summary

Historically the [Kong Kubernetes Ingress Controller (KIC)][kic] used bash scripts and `curl` to perform integration tests. This KEP aims to provide a standard [Golang][go] testing framework for Kong Kubernetes integrations to make testing Kong [Kubernetes][k8s] controllers straightforward and prescriptive.

[kic]:https://github.com/kong/kubernetes-ingress-controller
[go]:https://golang.org
[k8s]:https://github.com/kubernetes/kubernetes

## Motivation

- we want to have strong documentation of our _testing process_
- we want our integration tests to be in [Golang][go] to match our primary development language and provide high levels of expressibility for tests
- we want to be prescriptive how to test Kong Kubernetes components across all [@kong/team-k8s repos][team-repos] to limit to promote code re-use and limit test maintenance
- we want contributions of integration tests for Kong Kubernetes components to be a friendly experience
- we want to be able to define cluster environments wholly or in part by a simple container image and tag as the definition

[go]:https://golang.org
[team-repos]:https://github.com/orgs/Kong/teams/team-k8s/repositories

### Goals

## Proposal

### User Stories

#### Story 1

As a developer and contributor to any Kubernetes controller, I want my tests to require _no function calls for test environment setup_ so that I can effectively start writing my tests right away and consider the testing environment _somewhat opaque_.

#### Story 2

As a developer and contributor to any Kubernetes controller, when feasible I want the way I write tests for any single controller to be the same way I would write unit tests - by testing the controller on its Kubernetes API interface.

#### Story 3

As a contributor testing KIC and contributing PRs to KIC I want to be able to define most of my testing environment using container images, with tags that can indicate variants and versions.

## Design Details

The overall purpose of the [Kong Kubernetes Testing Framework (KTF)][ktf] is to provide environment provisioners for Kubernetes integration tests specific to Kong components with a minimal user interface.

Contributors to tests should be able expect an existing Kubernetes cluster, and that they will have access to a functional Kubernetes client object to start using in their tests.

We rely heavily on [Kubernetes In Docker (Kind)][kind] as a major component in our provisioning logic for test environments, as as part of our efforts we provide [Kong Kind Images][kind-images] which build on the existing images provided by Kind with Kong addons loaded on top.

[ktf]:https://github.com/kong/kubernetes-testing-framework

## Implementation History

- As part of [KEP 1][kep1] we created the [Kubernetes Testing Framework (KTF) Prototype][ktf]
- In service of the [KTF][ktf] we added a prototype [Kind Image Builder][kind-images] which creates container images used by KTF to bootstrap test clusters.
- The KTF was [updated to include Kubernetes API Object generators and cluster runbooks][ktf-pr3].
- A minimum test [was added to KIC][kic-pr1102] using the new KTF functionality
- The runbook concept was removed in favor of factory-style cluster provisioning

[kep1]:/keps/0001-single-kic-multi-gateway.md
[ktf]:https://github.com/kong/kubernetes-testing-framework
[kind-images]:https://github.com/kong/kind-images
[ktf-pr3]:https://github.com/Kong/kubernetes-testing-framework/pull/3
[kic-pr1102]:https://github.com/Kong/kubernetes-ingress-controller/pull/1102

## Alternatives

Using the existing open source [Kuttl][kuttl] testing tool was considered, but it is very limited in expressivity and the project is not well aligned to wanting to express tests in Go at the time of writing.

TODO: we should investigate more open source and off the shelf solutions just to make sure we're not missing anything that would save us time.

[kuttl]:https://github.com/kudobuilder/kuttl
