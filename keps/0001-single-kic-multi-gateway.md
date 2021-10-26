---
title: Single KIC Multi Gateway (DBLESS)
status: declined
---

# NOTES

For the time being we are considering this KEP **declined** as there was not sufficient requirements to move forward. In the future we can re-open this in an implementable state if needed. The research done here provided influence for several other follow-up projects that spun off of it (see [KIC 2.0][kic2]).

[kic2]:https://github.com/Kong/kubernetes-ingress-controller/milestone/12

# Single KIC Multi Gateway (DBLESS)

<!-- toc -->
- [Release Signoff Checklist](#release-signoff-checklist)
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
    - [Story 1](#story-1)
    - [Story 2](#story-2)
  - [Notes/Constraints/Caveats](#notesconstraintscaveats)
  - [Risks and Mitigations](#risks-and-mitigations)
- [Design Details](#design-details)
  - [Test Plan](#test-plan)
  - [Graduation Criteria](#graduation-criteria)
  - [Upgrade / Downgrade Strategy](#upgrade--downgrade-strategy)
  - [Version Skew Strategy](#version-skew-strategy)
- [Production Readiness Review Questionnaire](#production-readiness-review-questionnaire)
  - [Feature Enablement and Rollback](#feature-enablement-and-rollback)
  - [Rollout, Upgrade and Rollback Planning](#rollout-upgrade-and-rollback-planning)
  - [Monitoring Requirements](#monitoring-requirements)
  - [Dependencies](#dependencies)
  - [Scalability](#scalability)
  - [Troubleshooting](#troubleshooting)
- [Implementation History](#implementation-history)
- [Drawbacks](#drawbacks)
- [Alternatives](#alternatives)
- [Infrastructure Needed](#infrastructure-needed)
<!-- /toc -->

## Summary

Historically the Kong Kubernetes Ingress Controller (KIC) [required a controller pod to be present per proxy pod in DBLESS mode][kic702]. This KEP aims to make it possible to deploy a single controller to many any number of proxy instances in DBLESS mode.

[kic702]:https://github.com/Kong/kubernetes-ingress-controller/issues/702

## Motivation

- reduce load on the K8s API in high scale Kong environments where large numbers of proxies are needed
- we want to break up the [Kubernetes Ingress Controller (KIC)][kic] reconciliation loop in to smaller, specialized parts, to allow for easier maintenance and more flexibility
- provide status information for each proxy instance (configuration status, mTLS status)
- aggregate and deduplicate analytics data across any number of proxy instances
- make operations (including deployment, teardown, logging, events and statuses) for an HA deployment easy to understand and use, and well documented.
- ensure that considerations are made (and documented) so that all resulting efforts try to align us better towards a future built on [Gateway API][gateway-api]

[kic]:https://github.com/kong/kubernetes-ingress-controller
[gateway-api]:https://kubernetes-sigs.github.io/gateway-api/

### Goals

- support single-controller mode for configuring any number of `proxy` pods
- provide service discovery and automatic configuration for new proxy instances
- enable controller mTLS authentication on Admin API connections between KIC and proxy containers for security
- develop tooling to aggregate and deduplicate analytics across proxy pods
- develop documentation for new components and the overall HA setup options
- develop and present support and engineering enablement
- provide a seamless upgrade path for existing KIC deployments

### Non-Goals

- [Daemonset style HA deployments][daemonset-mode] will be part of a separate scope to keep scope tight.

[daemonset-mode]:https://github.com/Kong/kubernetes-ingress-controller/issues/682

## Proposal

### User Stories

#### Story 1

As an operator of Kong on Kubernetes I want to have an HA proxy deployment so that if a node containing a proxy fails other nodes are available to take on requests, without the overhead of regenerating the whole Kong configuration for each proxy on each update.

#### Story 2

As an operator of Kong on Kubernetes I want to be able to see the status of all proxies deployed in an HA configuration at a glance, and know what state they are in. I Specifically want to be able to easily see what configuration is currently loaded on a proxy server, and if there's any problem with loading a new configuration (Kong Admin API Errors or KIC validation errors). I also specifically want to know at a glance whether the proxy has the up to date configuration and is running and ready for requests.

#### Story 3

As a developer and contributor to Kong Kubernetes projects I want the purpose of Kong Kubernetes controllers to be clearly defined and with a succint scope to make contributing and maintenance simpler.

#### Story 4

As an operator of kong for kubernetes I want a single controller to be able to manage my proxy instances to avoid undue stress on the Kubernetes API when scaling to large numbers of proxies in DBLESS mode.

### Risks and Mitigations

One of the main risks of developing the HA deployment strategy is the addition of complexity at many different layers.

In order to reduce the risk of complexity and maintenance burden while still allowing for new functionality and good user experience, we will maintain a focus on strong separation of concerns between existing and new components, and endeaver to adhere strongly to [modular programming philosophy][modular-programming] and Kubernetes best practices.

[modular-programming]:https://en.wikipedia.org/wiki/Module_(programming)

### Graduation Criteria

#### Prerequisite: KIC 2.0 GA

- [ ] [KIC 2.0][kic2] GA

[kic2]:https://github.com/Kong/kubernetes-ingress-controller/milestone/12

## Alternatives

### Using Secrets

"@hbagdi" brought up concerns about the usage of secrets for the mTLS configuration:

https://github.com/Kong/keps/pull/1#discussion_r566176503

we need to investigate alternative solutions, including potentially creating our own CRD for this purpose (with encrypted `storage` in the backend).

We need a resolution on this concern that "@hbagdi" and "@shaneutt" sign off on before moving this KEP to `implementable` state.

# ADDITIONAL NOTES & Follow Up

- we need to handle DB mode as well (not just DBLESS mode) "@hbagdi" & "@shaneutt"
- we need to consider what statuses we will provide on the `Pods` themselves "@hbagdi" & "@shaneutt"
- for the new controllers the intent is to provide a single controller manager for all our new controllers AND include the existing KIC, making operational and maintenance burden of multiple controllers near `nil`, and all controller logs can be reached from a single pod.
- pod anti-affinity for proxy pods!
- integration/e2e tests for new APIs!
- user documentation for new APIs: ProxyInstance, proxy config, mtls config, logging, statuses, events, e.t.c.
- existing users MUST be able to upgrade without ingress downtime
- `helm release upgrade` needs to be capable of performing this upgrade (cleanly)
- we must have E2E/integration tests which PROVE the safety of this upgrade
- the complexity of the current KIC can be quite different in DB vs DBLESS, so as far as keeping the KIC simple, or breaking it up more we need to investigate and research this motivation further and support it with user stories "@rainest"
