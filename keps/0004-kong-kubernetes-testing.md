---
title: Kong Kubernetes Testing Framework (KTF)
status: implemented
---

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

## Summary

Historically the [Kong Kubernetes Ingress Controller (KIC)][kic] used bash
scripts and `curl` to perform integration tests. This KEP aims to provide a
standard [Golang][go] testing framework for Kong Kubernetes integrations to
make testing Kong [Kubernetes][k8s] controllers straightforward and
prescriptive.

[kic]:https://github.com/kong/kubernetes-ingress-controller
[go]:https://golang.org
[k8s]:https://github.com/kubernetes/kubernetes

## Motivation

- we want to have strong documentation of our _testing process_
- we want our integration tests to be in [Golang][go] to match our primary
  development language and provide high levels of expressibility for tests
- we want to be prescriptive how to test Kong Kubernetes components across
  all [@kong/team-k8s repos][team-repos] to limit to promote code re-use
  and limit test maintenance
- we want contributions of integration tests for Kong Kubernetes components to
  be a friendly experience

[go]:https://golang.org
[team-repos]:https://github.com/orgs/Kong/teams/team-k8s/repositories

### Goals

- provide _provisioning functionality_ for testing clusters (e.g. `kind`,
  `minikube`, `GKE`, e.t.c.)
- provide _deployment functionality_ for Kubernetes components (e.g. `helm`,
  `metallb`, e.t.c.)
- provide _deployment functionality_ for Kong components (e.g. deploying Proxy
  only, deploying Proxy with KIC, version matrix, e.t.c.)
- provide _generators_ to quickly generate default objects commonly used in
  testing (e.g. `Service`, `Deployment`, e.t.c.)

### Non-Goals

- provide _mocking functionality_ for the Kong Admin API

## Proposal

The main purpose of the [Kong Kubernetes Testing Framework (KTF)][ktf] is to
provision Kubernetes testing environments on a local system for `go test` to
run integration/e2e tests against. Secondarily the testing framework will
provide an assortment of helper functions for common functionality (e.g.
deploying addons like Kong).

[ktf]:https://github.com/kong/kubernetes-testing-framework

### User Stories

#### Story 1

As a developer I want the testing environment to already pre-populate a
Kubernetes cluster and provide me the client object to use the cluster
for my tests.

### Risks and Mitigations

One of the main risks of this proposal is the increase in CI costs for the
team. The main way in which we mitigate this is making sure our CI workflow
scope is limited to the minimum necessary and keep our build triggers clean.

## Design Details

### Test Plan

The testing framework itself needs to be tested: all code which can feasibly
be unit tested will be, and cluster and addon provisioning will have
integration and e2e tests. These tests will also help to serve as examples for
end-users of the testing framework.

### Graduation Criteria

- [X] Testing Framework Prototype
- [X] Testing Framework plugged into KIC
- [x] KTF `v0.1.0` milestone completed & `v0.1.0` released

## Implementation History

- As part of [KEP 1][kep1] we created the [Kubernetes Testing Framework (KTF)
  Prototype][ktf]
- In service of the [KTF][ktf] we added a prototype [Kind Image
  Builder][kind-images] which creates container images used by KTF to
  bootstrap test clusters.
- The KTF was [updated to include Kubernetes API Object generators and cluster
  runbooks][ktf-pr3].
- A minimum test [was added to KIC][kic-pr1102] using the new KTF functionality
- The runbook concept was removed in favor of factory-style cluster provisioning
- We remove the experimental image builder in favor of adding deployment tooling
  to the test framework.
- Admin API mocking was added in `pkg/kong/fake_admin_api.go` and is now in use by KIC integration tests

[kep1]:/keps/0001-single-kic-multi-gateway.md
[ktf]:https://github.com/kong/kubernetes-testing-framework
[kind-images]:https://github.com/kong/kind-images
[ktf-pr3]:https://github.com/Kong/kubernetes-testing-framework/pull/3
[kic-pr1102]:https://github.com/Kong/kubernetes-ingress-controller/pull/1102

## Alternatives

Using the existing open source [Kuttl][kuttl] testing tool was considered, but
it is very limited in expressivity and the project is not well aligned to
wanting to express tests in Go at the time of writing.

[kuttl]:https://github.com/kudobuilder/kuttl
