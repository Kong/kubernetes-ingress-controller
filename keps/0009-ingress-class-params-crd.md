---
title: Introduce IngressClassParams CRD to control IngressClass behavior
status: provisional
---

# Introduce IngressClassParams CRD to control IngressClass behavior

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [Graduation Criteria](#graduation-criteria)
- [Implementation History](#implementation-history)
<!-- /toc -->

## Summary

As of now users can control the behavior of an `IngressClass` via e.g.
[`ingress.kubernetes.io/service-upstream` annotation][service-upstream-annotation].
This is handy in general but gets cumbersome when one needs to change this for
multiple services.

Controlling `IngressClass` parameters in one place would be much more desirable in
such situations hence the proposal to introduce a new object - `IngressClassParams` -
which would allow said customizations.

[service-upstream-annotation]: https://docs.konghq.com/kubernetes-ingress-controller/2.3.x/references/annotations/#ingresskubernetesioservice-upstream

## Motivation

- we want to allow `IngressClass` behavior customizations in one place
- we want to make configuration management related to `IngressClass` to not be
  a burden for the user even when IngressClass is assigned to many services

### Goals

- provide a definition of new object - `IngressClassParams`
- provide an _ability to configure_ `IngressClass` in one place via aforementioned
  `IngressClassParams` object
- provide `ServiceUpstream` field in `IngressClassParams` object, as means to control
  the behavior already provided by `ingress.kubernetes.io/service-upstream` annotation

## Proposal

### Graduation Criteria

- [ ] `IngressClassParams` object available for deployment for end users
- [ ] `ServiceUpstream` field implemented in `IngressClassParams` object to allow
  control of behavior of all services handled by an `IngressClass`

## Implementation History

- First issue proposing to enable `service-upstream` for an entire ingress class
  in [#774][774]
- [#1131][1131] was raised to discuss `IngressClass` parameters more generically
- An initial proposal has been made in [#1586][1586] (trying to address [#1131][1131])
  but the work did not conclude and eventually the PR didn't get merged
- Another proposal has been made in [#2535][2535] basing its work on [#1586][1586]

[774]: https://github.com/Kong/kubernetes-ingress-controller/pull/774
[1131]: https://github.com/Kong/kubernetes-ingress-controller/pull/1131
[1586]: https://github.com/Kong/kubernetes-ingress-controller/pull/1586
[2535]: https://github.com/Kong/kubernetes-ingress-controller/pull/2535
