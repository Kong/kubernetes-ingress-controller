[![][kong-logo]][kong-url]
[![Build Status](https://github.com/kong/kubernetes-ingress-controller/workflows/Test/badge.svg)](https://github.com/kong/kubernetes-ingress-controller/actions?query=branch%3Amaster+event%3Apush)
[![Go Reference](https://pkg.go.dev/badge/github.com/kong/kubernetes-ingress-controller/v2.svg)](https://pkg.go.dev/github.com/kong/kubernetes-ingress-controller/v2)
[![Codecov](https://codecov.io/gh/Kong/kubernetes-ingress-controller/branch/main/graph/badge.svg?token=S1aqcXiGEo)](https://codecov.io/gh/Kong/kubernetes-ingress-controller)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/Kong/kong/blob/master/LICENSE)
[![Twitter](https://img.shields.io/twitter/follow/thekonginc.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=thekonginc)

# Kong Ingress Controller for Kubernetes (KIC)

Use [Kong][kong] for Kubernetes [Ingress][ingress].
Configure [plugins][kong-hub], health checking,
load balancing and more in Kong
for Kubernetes Services, all using
Custom Resource Definitions(CRDs) and Kubernetes-native tooling.

[**Features**](#features) | [**Get started**](#get-started) | [**Documentation**](#documentation) | [**main branch builds**](#main-branch-builds) | [**Seeking help**](#seeking-help)

## Features

- **Ingress routing**
  Use [Ingress][ingress] resources to configure Kong
- **Enhanced API management using plugins**
  Use a wide-array of [plugins][kong-hub]
  to monitor, transform, protect your traffic.
- **Native gRPC support**
  Proxy gRPC traffic and gain visibility into it using
  Kong's plugin.
- **Health checking and Load-balancing**
  Load balance requests across your pods and supports active & passive health-checks.
- **Request/response transformations**
  Use plugins to
  modify your requests/responses on the fly.
- **Authentication**
  Protect your services using authentication methods
  of your choice.
- **Declarative configuration for Kong**
  Configure all of Kong
  using CRDs in Kubernetes and manage Kong declaratively.

## Get started

You can use
[Minikube, Kind](https://kubernetes.io/docs/tasks/tools/)
on your local machine or use
a hosted k8s service like
[GKE](https://cloud.google.com/kubernetes-engine/).

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

Once installed, please follow the [Getting Started guide][getting-started-guide]
to start using Ingress in your Kubernetes cluster.

> Note: Kong Enterprise users, please follow along our
[enterprise guide][k4k8s-enterprise-setup] to setup the enterprise version.

## Contributing

We ❤️ pull requests and we’re continually working hard to make it as easy as possible for developers to contribute. Before beginning development with the Kong Ingress Controller, please familiarize yourself with the following developer resources:
- [CONTRIBUTING](CONTRIBUTING.md)
- [CODE_OF_CONDUCT](CODE_OF_CONDUCT.md) and [COPYRIGHT](https://github.com/Kong/kong/blob/master/COPYRIGHT)

## Documentation

All documentation for the Kong Ingress Controller is present in the [kong/docs.konghq.com](https://github.com/kong/docs.konghq.com) repository. Pull Requests are welcome for additions and corrections.

## Guides and Tutorials

Please browse through the [guides][guides] to get started and to learn specific ingress controller operations.

## main branch builds

Nightly pre-release builds of the `main` branch are available from the
[kong/nightly-ingress-controller repository][nightly-images] hosted on Docker Hub:

`main` contains unreleased new features for upcoming minor and major releases:

```
docker pull kong/nightly-ingress-controller:nightly
```

## Seeking help

Please search through the posts on the [discussions
page](https://github.com/Kong/kubernetes-ingress-controller/discussions)
or the [Kong Nation Forums](https://discuss.konghq.com/c/kubernetes)
as it's likely that another user has run into the same problem.
If you don't find an answer, please feel free to post a question.

If you've found a bug please [open an issue](https://github.com/kong/kubernetes-ingress-controller/issues).

For a feature request please open an issue using the feature request template.

You can also talk to the developers behind Kong in the
[#kong](https://kubernetes.slack.com/messages/kong) channel on the
Kubernetes Slack server.

### Preview and Experimental Features

At any time the KIC may include features or options that are considered
experimental and are not enabled by default, nor available in the [Kong
Documentation Site][kongdocs].

To try out new features that are behind feature gates, please see the
preview features in [FEATURE_GATES.md][fgates] and documentation for these
preview features can be found in [FEATURE_PREVIEW_DOCUMENTATION.md][fpreview].

[kongdocs]:https://docs.konghq.com
[fgates]:/FEATURE_GATES.md
[fpreview]:/FEATURE_PREVIEW_DOCUMENTATION.md

### Community meetings

You can join monthly meetups hosted by [Kong](https://konghq.com) to ask questions, provide feedback or just to listen and hang out.
See the [Online Meetups Page](https://konghq.com/online-meetups/) to sign up and receive meeting invites and [Zoom](https://zoom.us) links.

[ingress]: https://kubernetes.io/docs/concepts/services-networking/ingress/
[kong]: https://konghq.com/kong-community-edition/
[kong-hub]: https://docs.konghq.com/hub/
[docs]: https://docs.konghq.com/kubernetes-ingress-controller/latest/introduction/
[deployment]: https://docs.konghq.com/kubernetes-ingress-controller/latest/deployment/overview/
[annotations]: https://docs.konghq.com/kubernetes-ingress-controller/latest/references/annotations/
[crds]: https://docs.konghq.com/kubernetes-ingress-controller/latest/references/custom-resources/
[faqs]: https://docs.konghq.com/kubernetes-ingress-controller/latest/faq/
[getting-started-guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/getting-started/
[docker-images]: https://hub.docker.com/r/kong/kubernetes-ingress-controller
[nightly-images]: https://hub.docker.com/r/kong/nightly-ingress-controller
[kong-url]: https://konghq.com/
[kong-logo]: https://konghq.com/wp-content/uploads/2018/05/kong-logo-github-readme.png
[k4k8s-enterprise-setup]: https://docs.konghq.com/kubernetes-ingress-controller/latest/deployment/k4k8s-enterprise/
[guides]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/overview/
