---
title: Introduce IngressClassParameters CRD to control IngressClass behavior
status: provisional
---

# Introduce IngressClassParameters CRD to control IngressClass behavior

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
- [Proposal](#proposal)
  - [Graduation Criteria](#graduation-criteria)
- [Implementation History](#implementation-history)
<!-- /toc -->

## Summary

Currently users can control how traffic is load-balanced for an `Ingress`
resource using the [`ingress.kubernetes.io/service-upstream` annotation][service-upstream-annotation].
This can quickly get cumbersome, specifically for environments with large number
of `Ingress` resources.

Controlling `IngressClass` parameters in one place would be much more desirable in
such situations, hence the proposal to introduce a new object - `IngressClassParameters` -
which would allow said customizations.

[service-upstream-annotation]: https://docs.konghq.com/kubernetes-ingress-controller/2.3.x/references/annotations/#ingresskubernetesioservice-upstream

## Motivation

- we want to allow `IngressClass` behavior customizations in one place
- we want to make configuration management related to `IngressClass` to not be
  a burden for the user even when IngressClass is assigned to many services

### Goals

- provide a definition of new object - `IngressClassParameters`
- provide an _ability to configure_ `IngressClass` in one place via aforementioned
  `IngressClassParameters` object
- provide `ServiceUpstream` field in `IngressClassParameters` object, as means to control
  the behavior already provided by `ingress.kubernetes.io/service-upstream` annotation

## Proposal

### Graduation Criteria

- [ ] introduce `IngressClassParameters` object available for deployment for end users
- [ ] introduce a [feature gate][feature-gates] for the functionality described in this KEP
- [ ] `ServiceUpstream` field implemented in `IngressClassParameters` object to allow
  control of behavior of all services handled by an `Ingress` configured with an
  `IngressClass` that references said parameters object

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
[feature-gates]: ../FEATURE_GATES.md
