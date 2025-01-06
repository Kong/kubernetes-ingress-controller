[![kong-logo]][kong-url]
[![Build Status](https://github.com/Kong/kubernetes-ingress-controller/actions/workflows/checks.yaml/badge.svg)](https://github.com/Kong/kubernetes-ingress-controller/actions/workflows/checks.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kong/kubernetes-ingress-controller/v3.svg)](https://pkg.go.dev/github.com/kong/kubernetes-ingress-controller/v3)
[![Codecov](https://codecov.io/gh/Kong/kubernetes-ingress-controller/branch/main/graph/badge.svg?token=S1aqcXiGEo)](https://codecov.io/gh/Kong/kubernetes-ingress-controller)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/Kong/kong/blob/master/LICENSE)
[![Twitter](https://img.shields.io/twitter/follow/thekonginc.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=thekonginc)
[![Conformance](https://img.shields.io/badge/Gateway%20API%20Conformance%20v1.0.0-Kong%20Ingress%20Controller%203.0-green)](https://github.com/kubernetes-sigs/gateway-api/blob/main/conformance/reports/v1.0.0/kong-kubernetes-ingress-controller.yaml)

# Kong Ingress Controller for Kubernetes (KIC)

Use [Kong][kong] for Kubernetes [Gateway API][gwapi] or [Ingress][ingress].
Configure [plugins][docs-konghq-hub], health checking,
load balancing and more, all using
Custom Resource Definitions (CRDs) and Kubernetes-native tooling.

[**Features**](#features) | [**Get started**](#get-started) | [**Documentation**](#documentation) | [**main branch builds**](#main-branch-builds) | [**Seeking help**](#seeking-help)

## Features

- **Gateway API support**
  Use [Gateway API][gwapi] resources (official successor of [Ingress][ingress] resources) to configure Kong.
  Native support for TCP, UDP, TLS, gRPC and HTTP/HTTPS traffic, reuse the same gateway for multiple protocols and namespaces.
- **Ingress support**
  Use [Ingress][ingress] resources to configure Kong.
- **Declarative configuration for Kong**
  Configure all of Kong features in declarative Kubernetes native way with CRDs.
- **Seamlessly operate Kong**
  Scale and manage multiple replicas of Kong Gateway automatically to ensure performance and high-availability.
- **Health checking and load-balancing**
  Load balance requests across your pods and supports active & passive health-checks.
- **Enhanced API management using plugins**
  Use a wide array of [plugins][docs-konghq-hub] for e.g.
  - authentication
  - request/response transformations
  - rate-limiting

## Get started (using Helm)

You can use [Minikube or Kind][k8s-io-tools] on your local machine or use
a hosted Kubernetes service like [GKE](https://cloud.google.com/kubernetes-engine/).

### Install the Gateway API CRDs

This command will install all resources that have graduated to GA or beta,
including `GatewayClass`, `Gateway`, `HTTPRoute`, and `ReferenceGrant`.

```shell
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/standard-install.yaml
```

Or, if you want to use experimental resources and fields such as `TCPRoute`s and `UDPRoute`s,
please run this command.

```shell
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/experimental-install.yaml
```

### Install the Kong Ingress Controller with Helm

```shell
helm install kong --namespace kong --create-namespace --repo https://charts.konghq.com ingress
```

To learn more details about Helm chart follow the [Helm chart documentation](https://charts.konghq.com/).

Once installed, please follow the [Getting Started guide][docs-konghq-getting-started-guide]
to start using Kong in your Kubernetes cluster.

> Note: Kong Enterprise users, please follow along with our
[enterprise guide][docs-konghq-k4k8s-enterprise-setup] to setup the enterprise version.

## Get started (using Operator) 

As an alternative to Helm, you can also install Kong Ingress Controller using the **Kong Gateway Operator** by following this [quick start guide][kgo-guide].

## Container images

### Release images

Release builds of Kong Ingress Controller can be found on Docker Hub in
[kong/kubernetes-ingress-controller repository][dockerhub-kic].

At the moment we're providing images for:

- Linux `amd64`
- Linux `arm64`

### `main` branch builds

Nightly pre-release builds of the `main` branch are available from the
[kong/nightly-ingress-controller repository][dockerhub-kic-nightly] hosted on Docker Hub:

`main` contains unreleased new features for upcoming minor and major releases:

```shell
docker pull kong/nightly-ingress-controller:nightly
```

## Documentation

All documentation for the Kong Ingress Controller is present in the [kong/docs.konghq.com](https://github.com/kong/docs.konghq.com) repository. Pull Requests are welcome for additions and corrections.

### Guides and Tutorials

Please browse through the [guides][docs-konghq-kic-guides] to get started and to learn specific ingress controller operations.

## Contributing

We ❤️ pull requests and we’re continually working hard to make it as easy as possible for developers to contribute.
Before beginning development with the Kong Ingress Controller, please familiarize yourself with the following developer resources:

- [TESTING](TESTING.md)
- [CONTRIBUTING](CONTRIBUTING.md)
- [CODE_OF_CONDUCT](CODE_OF_CONDUCT.md)
- [COPYRIGHT](https://github.com/Kong/kong/blob/master/COPYRIGHT)

## Seeking help

Please search through the [FAQs][docs-konghq-faqs], posts on the
[discussions page][github-kic-discussions] or the
[Kong Nation Forums](https://discuss.konghq.com/c/kubernetes)
as it's likely that another user has run into the same problem.
If you don't find an answer, please feel free to post a question.

If you've found a bug, please [open an issue][github-kic-issues].

For a feature request, please open an issue using the feature request template.

You can also talk to the developers behind Kong in the
[#kong][slack-kubernetes-kong] channel on the Kubernetes Slack server.

### Community meetings

You can join monthly meetups hosted by [Kong](https://konghq.com) to ask questions, provide feedback, or just to listen and hang out.
See the [Online Meetups Page](https://konghq.com/online-meetups/) to sign up and receive meeting invites and [Zoom](https://zoom.us) links.

## Preview and Experimental Features

At any time the KIC may include features or options that are considered
experimental and are not enabled by default, nor available in the [Kong
Documentation Site][docs-konghq].

To try out new features that are behind feature gates, please see the
preview features in [FEATURE_GATES.md][fgates] and documentation for these
preview features can be found in [FEATURE_PREVIEW_DOCUMENTATION.md][fpreview].

[fgates]:/FEATURE_GATES.md
[fpreview]:/FEATURE_PREVIEW_DOCUMENTATION.md
[ingress]: https://kubernetes.io/docs/concepts/services-networking/ingress/
[gwapi]: https://gateway-api.sigs.k8s.io/
[kong]: https://konghq.com/kong
[kong-url]: https://konghq.com/
[kong-logo]: https://konghq.com/wp-content/uploads/2018/05/kong-logo-github-readme.png
[k8s-io-tools]: https://kubernetes.io/docs/tasks/tools/
[slack-Kubernetes-kong]: https://kubernetes.slack.com/messages/kong

[dockerhub-kic]: https://hub.docker.com/r/kong/kubernetes-ingress-controller
[dockerhub-kic-nightly]: https://hub.docker.com/r/kong/nightly-ingress-controller

[github-kic-discussions]: https://github.com/Kong/kubernetes-ingress-controller/discussions
[github-kic-issues]: https://github.com/kong/kubernetes-ingress-controller/issues

[docs-konghq]:https://docs.konghq.com
[docs-konghq-hub]: https://docs.konghq.com/hub/
[docs-konghq-faqs]: https://docs.konghq.com/kubernetes-ingress-controller/latest/faq/
[docs-konghq-getting-started-guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/getting-started/
[docs-konghq-k4k8s-enterprise-setup]: https://docs.konghq.com/kubernetes-ingress-controller/latest/deployment/k4k8s-enterprise/
[docs-konghq-kic-guides]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/overview/

[kgo-guide]: https://docs.konghq.com/gateway-operator/latest/get-started/kic/install/
