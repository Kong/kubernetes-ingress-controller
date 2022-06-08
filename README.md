[![][kong-logo]][kong-url]
[![Build Status](https://github.com/kong/kubernetes-ingress-controller/workflows/Test/badge.svg)](https://github.com/kong/kubernetes-ingress-controller/actions?query=branch%3Amaster+event%3Apush)
[![Go Reference](https://pkg.go.dev/badge/github.com/kong/kubernetes-ingress-controller/v2.svg)](https://pkg.go.dev/github.com/kong/kubernetes-ingress-controller/v2)
[![Codecov](https://codecov.io/gh/Kong/kubernetes-ingress-controller/branch/main/graph/badge.svg?token=S1aqcXiGEo)](https://codecov.io/gh/Kong/kubernetes-ingress-controller)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/Kong/kong/blob/master/LICENSE)
[![Twitter](https://img.shields.io/twitter/follow/thekonginc.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=thekonginc)

# Kong Ingress Controller for Kubernetes (KIC)

Use [Kong][kong] for Kubernetes [Ingress][ingress].
Configure [plugins][docs-konghq-hub], health checking,
load balancing, and more in Kong
for Kubernetes Services, all using
Custom Resource Definitions (CRDs) and Kubernetes-native tooling.

[**Features**](#features) | [**Get started**](#get-started) | [**Documentation**](#documentation) | [**main branch builds**](#main-branch-builds) | [**Seeking help**](#seeking-help)

## Features

- **Ingress routing**
  Use [Ingress][ingress] resources to configure Kong.
- **Enhanced API management using plugins**
  Use a wide array of [plugins][docs-konghq-hub] to monitor, transform
  and protect your traffic.
- **Native gRPC support**
  Proxy gRPC traffic and gain visibility into it using Kong's plugins.
- **Health checking and Load-balancing**
  Load balance requests across your pods and supports active & passive health-checks.
- **Request/response transformations**
  Use plugins to modify your requests/responses on the fly.
- **Authentication**
  Protect your services using authentication methods of your choice.
- **Declarative configuration for Kong**
  Configure all of Kong using CRDs in Kubernetes and manage Kong declaratively.

## Get started

You can use [Minikube or Kind][k8s-io-tools] on your local machine or use
a hosted Kubernetes service like [GKE](https://cloud.google.com/kubernetes-engine/).

Setting up Kong for Kubernetes is as simple as:

```shell
# using YAMLs
$ kubectl apply -f https://bit.ly/k4k8s

# or using Helm
$ helm repo add kong https://charts.konghq.com
$ helm repo update

# Helm 3
$ helm install kong/kong --generate-name --set ingressController.installCRDs=false
```

Once installed, please follow the [Getting Started guide][docs-konghq-getting-started-guide]
to start using Ingress in your Kubernetes cluster.

> Note: Kong Enterprise users, please follow along with our
[enterprise guide][docs-konghq-k4k8s-enterprise-setup] to setup the enterprise version.

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

```
docker pull kong/nightly-ingress-controller:nightly
```

## Documentation

All documentation for the Kong Ingress Controller is present in the [kong/docs.konghq.com](https://github.com/kong/docs.konghq.com) repository. Pull Requests are welcome for additions and corrections.

### Guides and Tutorials

Please browse through the [guides][docs-konghq-kic-guides] to get started and to learn specific ingress controller operations.

## Contributing

We ❤️ pull requests and we’re continually working hard to make it as easy as possible for developers to contribute.
Before beginning development with the Kong Ingress Controller, please familiarize yourself with the following developer resources:

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
