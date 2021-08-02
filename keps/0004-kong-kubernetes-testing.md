---
status: implementable
---

**NOTE**: this will be considered `implemented` once [Milestone 1][m1] is completed.

[m1]:https://github.com/Kong/kubernetes-testing-framework/milestone/1

# Kong Kubernetes Testing Framework (KTF)

<!-- toc -->
- [Release Signoff Checklist](#release-signoff-checklist)
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories (Optional)](#user-stories-optional)
    - [Story 1](#story-1)
    - [Story 2](#story-2)
  - [Risks and Mitigations](#risks-and-mitigations)
- [Design Details](#design-details)
  - [Test Plan](#test-plan)
  - [Graduation Criteria](#graduation-criteria)
- [Implementation History](#implementation-history)
- [Drawbacks](#drawbacks)
- [Alternatives](#alternatives)
<!-- /toc -->

## Release Signoff Checklist

- [ ] [v0.1.0 Milestone][ms1]

[ms1]:https://github.com/Kong/kubernetes-testing-framework/milestone/1

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

[go]:https://golang.org
[team-repos]:https://github.com/orgs/Kong/teams/team-k8s/repositories

### Goals

- provide _provisioning functionality_ for testing clusters (e.g. `kind`, `minikube`, `GKE`, e.t.c.)
- provide _deployment functionality_ for Kubernetes components to create complete testing environments (e.g. `helm`, `metallb`, e.t.c.)
- provide _deployment functionality_ for Kong components (e.g. deploying Proxy only, deploying Proxy with KIC, version matrix, e.t.c.)
- provide _generators_ to quickly generate default objects commonly used in testing (e.g. `Service`, `Deployment`, e.t.c.)
- provide _mocking functionality_ for the Kong Admin API

### Non-Goals

Due to incongruencies with one of our most prominent upstream tools ([Kind][kind]) we're going to need to skip on creating complete container images for testing environments in favor of writing setup logic aftermarket for existing default images. While being able to move runtime problems to build time would be helpful, we'll potentially need to look at migrating to new tools in some future iteration to follow up.

For this iteration we're _not trying to solve the problem of automated cleanup of testing environments_, as this is expected to greatly increase scope. Instead tests are expected to use separate namespaces and cleanup after themselves for this iteration. There is a follow up item to take care of cleanup as a separate scope: https://github.com/Kong/kubernetes-testing-framework/issues/4.

## Proposal

The overall purpose of the [Kong Kubernetes Testing Framework (KTF)][ktf] is to provide environment provisioners for Kubernetes integration tests specific to Kong components with a minimal user interface.

Contributors to tests should be able expect an existing Kubernetes cluster, and that they will have access to a functional Kubernetes client object to start using in their tests.

[ktf]:https://github.com/kong/kubernetes-testing-framework

### User Stories

#### Story 1

As a developer and contributor to any Kubernetes controller, I want my tests to require _no function calls for test environment setup_ so that I can effectively start writing my tests right away and consider the testing environment _somewhat opaque_.

#### Story 2

As a developer and contributor to any Kubernetes controller, when feasible I want the way I write tests for any single controller to be the same way I would write unit tests - by testing the controller on its Kubernetes API interface.

### Risks and Mitigations

One of the main risks of this proposal is the increase in CI costs for the team. The main way in which we mitigate this is making sure our CI workflow scope is limited to the minimum necessary and keep our build triggers clean.

## Design Details

### Test Plan

All code which can feasibly be unit tested will be, the only code we expect will not be feasible to unit test will be code that makes shell calls to other programs.

Higher level integration tests will exist to cover:

- setup of a Kind Kubernetes cluster
- basic usage of a Kind Kubernetes cluster's client-go client
- teardown of a Kind Kubernetes cluster

Integration tests written in the testing framework will also serve the purpose of being examples for callers using the library to write tests, as such `go doc` documentation and heavily commented tests is a must.

### Graduation Criteria

- [X] Testing Framework Prototype
- [X] Testing Framework plugged into KIC
- [ ] KTF `v0.1.0` milestone completed & `v0.1.0` released

## Implementation History

- As part of [KEP 1][kep1] we created the [Kubernetes Testing Framework (KTF) Prototype][ktf]
- In service of the [KTF][ktf] we added a prototype [Kind Image Builder][kind-images] which creates container images used by KTF to bootstrap test clusters.
- The KTF was [updated to include Kubernetes API Object generators and cluster runbooks][ktf-pr3].
- A minimum test [was added to KIC][kic-pr1102] using the new KTF functionality, `v0.0.1` tagged.
- The runbook concept was removed in favor of factory-style cluster provisioning
- We remove the experimental image builder in favor of adding deployment tooling to the test framework, `v0.0.2` tagged.
- Admin API mocking was added in `pkg/kong/fake_admin_api.go` and is now in use by KIC integration tests

[kep1]:/keps/0001-single-kic-multi-gateway.md
[ktf]:https://github.com/kong/kubernetes-testing-framework
[kind-images]:https://github.com/kong/kind-images
[ktf-pr3]:https://github.com/Kong/kubernetes-testing-framework/pull/3
[kic-pr1102]:https://github.com/Kong/kubernetes-ingress-controller/pull/1102

## Alternatives

Using the existing open source [Kuttl][kuttl] testing tool was considered, but it is very limited in expressivity and the project is not well aligned to wanting to express tests in Go at the time of writing.

[kuttl]:https://github.com/kudobuilder/kuttl
