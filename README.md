[![][kong-logo]][kong-url]
[![Build Status](https://github.com/kong/kubernetes-ingress-controller/workflows/Test/badge.svg)](https://github.com/kong/kubernetes-ingress-controller/actions?query=branch%3Amaster+event%3Apush)
[![codecov](https://codecov.io/gh/Kong/kubernetes-ingress-controller/branch/main/graph/badge.svg?token=S1aqcXiGEo)](https://codecov.io/gh/Kong/kubernetes-ingress-controller)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/Kong/kong/blob/master/LICENSE)
[![Twitter](https://img.shields.io/twitter/follow/thekonginc.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=thekonginc)

# Kong for Kubernetes

# Not supported, k8s 1.22

⚠️ **Due to Bintray image registries going out of service, we've moved our Docker images to [Docker Hub](https://hub.docker.com/r/kong/kubernetes-ingress-controller/tags).** ⚠️

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
[Minikube](https://kubernetes.io/docs/setup/minikube/)
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

## Documentation

All documentation around Kong Ingress Controller is present in this
repository inside the [docs][docs] directory.
Pull Requests are welcome for additions and corrections.

## Guides and Tutorials

Please browse through [guides][guides] to get started or understand how to configure
a specific setting with Kong Ingress Controller.

## main branch builds

Pre-release builds of the `main` branch are available from the
[kong-ingress-controller repository][docker-images] hosted on Docker Hub:

`main` contains unreleased new features for upcoming minor and major releases:

```
docker pull kong/kubernetes-ingress-controller:main
```

## Seeking help

Please search through the posts on
[Kong Nation](https://discuss.konghq.com/c/kubernetes) as it's
likely that another user has run into the same problem.
If you don't find an answer, please feel free to post a question.
For a feature request, please post in
[Feature Suggestions](https://discuss.konghq.com/c/feature-suggestions)
category.

You can also talk to the developers behind Kong in the
[#kong](https://kubernetes.slack.com/messages/kong) channel on the
Kubernetes Slack server.

If you've spotted a bug, please open an issue
on our [Github](https://github.com/kong/kubernetes-ingress-controller/issues).

### Community meetings

You can join monthly meetings hosted by the maintainers of the project
to ask questions, provide feedback or just come and say hello.
The meeting takes place on every second Tuesday of the month
at 10 am Pacific time.
Please submit your contact details on the
[online meetups](https://konghq.com/online-meetups/) page to receive
meeting invite and Zoom links to join the meeting.

[ingress]: https://kubernetes.io/docs/concepts/services-networking/ingress/
[kong]: https://konghq.com/kong-community-edition/
[kong-hub]: https://docs.konghq.com/hub/
[docs]: https://docs.konghq.com/kubernetes-ingress-controller/latest/introduction/
[deployment]: https://docs.konghq.com/kubernetes-ingress-controller/latest/deployment/overview/
[annotations]: https://docs.konghq.com/kubernetes-ingress-controller/latest/references/annotations/
[crds]: https://docs.konghq.com/kubernetes-ingress-controller/latest/references/custom-resources/
[faqs]: https://docs.konghq.com/kubernetes-ingress-controller/latest/faq/
[getting-started-guide]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/getting-started/
[badge-travis-image]: https://travis-ci.org/Kong/kubernetes-ingress-controller.svg?branch=master
[badge-travis-url]: https://travis-ci.org/Kong/kubernetes-ingress-controller
[docker-images]: https://hub.docker.com/r/kong/kubernetes-ingress-controller
[kong-url]: https://konghq.com/
[kong-logo]: https://konghq.com/wp-content/uploads/2018/05/kong-logo-github-readme.png
[k4k8s-enterprise-setup]: https://docs.konghq.com/kubernetes-ingress-controller/latest/deployment/k4k8s-enterprise/
[guides]: https://docs.konghq.com/kubernetes-ingress-controller/latest/guides/overview/
