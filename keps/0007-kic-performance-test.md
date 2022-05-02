---
title: Performance Testing
status: declined
---

**NOTE**: this KEP is currently declined as we had other higher priorities
          at the time of writing. However this does not mean it's never
          going to happen, it's simply a question of when makes the most
          sense. If you're reading this and interested in this functionality
          please feel free to create a discussion or issue regarding your
          use case and we can re-assess whether it makes sense to work on
          this sooner.

## Summary

We want a performance testing solution for the [Kong Kubernetes Ingress
Controller (KIC)][kic] to collect CPU, Memory and I/O performance metrics from a
running KIC deployment under a variety of scenarios. This testing framework will
be used to evaluate the performance characteristics of KIC in specific
environments and will help us identify bottlenecks.

[kic]:https://github.com/kong/kubernetes-ingress-controller

## Motivation

- we currently have no data, documentation, or determined performance
  characteristics for the KIC in any environment
- we want end-users and customers to be able to run a portable copy of our
  performance tests in their own environments

### Goals

- create Golang performance tests which can be easily run with `go test`
- add a metrics report for performance test runs: CPU, Memory, I/O
- make the tests portable, provide a container image for the test suite
