# Kubernetes Ingress Controller for Kong

[![Build Status](https://travis-ci.org/Kong/kubernetes-ingress-controller.svg?branch=master)](https://travis-ci.org/Kong/kubernetes-ingress-controller)

Use [Kong][kong] for your Kubernetes [Ingress][ingress]
and further configure [plugins][kong-hub], health checking,
load balancing and more in Kong
for your Kubernetes services, all using
Custom Resource Definitions(CRDs).

## Tables of content

- [**Features**](#features)
- [**Version support matrix**](#version-support-matrix)
- [**Get started**](#get-started)
- [**Documentation**](#documentation)
- [**Seeking help**](#seeking-help)
- [**Design**](#design)
- [**License**](#license)
- [**Roadmap**](#roadmap)

## Features

- **Ingress routing**: Use [Ingress][ingress] resources to configure Kong
- **Health checking and Load-balancing**: Load balance requests across
  your pods and supports active & passive health-checks.
- **Configure Plugins**: Execute custom code
  as a request is proxied to your service.
- **Request/response transformations**: Use plugins to
  modify your requests/responses on the fly.
- **Authentication**: Protect your services using authentication
  plugins.
- **Declarative configuration for Kong** Configure all of Kong
  using CRDs in Kubernetes and manage Kong declaratively.

## Version support matrix

The Ingress controller is tested on
Kubernetes version `1.8` through `1.10`.

The following matrix lists supported versions of
Kong for every release of the Kong Ingress Controller:

| Kong Ingress Controller  | <= 0.0.4           | 0.0.5              | 0.1.x              | 0.2.x              |
|--------------------------|:------------------:|:------------------:|:------------------:|:------------------:|
| Kong 0.13.x              | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                |
| Kong 0.14.x              | :x:                | :x:                | :x:                | :white_check_mark: |
| Kong Enterprise 0.32.x   | :x:                | :white_check_mark: | :white_check_mark: | :x:                |
| Kong Enterprise 0.33.x   | :x:                | :white_check_mark: | :white_check_mark: | :x:                |

## Get started

You can deploy Kong Ingress Controller on any
Kubernetes cluster which supports a Service of `type: LoadBalancer`.

You can use
[Minikube](https://kubernetes.io/docs/setup/minikube/)
on your local machine or use
a hosted k8s service like
[GKE](https://cloud.google.com/kubernetes-engine/).

To setup Kong Ingress Controller in your k8s cluster, execute:

```shell
kubectl apply -f https://bit.ly/kong-ingress
```

It takes a few minutes for all components to
spin up.
You now have set up Kong as your Ingress controller and
all Ingress resources in your Kubernetes Cluster will be satisfied.

Please refer our [deployment documentation][deployment-doc]
for a detailed  introduction to Kong Ingress Controller
and Ingress spec.

## Seeking help

Please search through the posts on
[Kong Nation](https://discuss.konghq.com/c/kubernetes) as it's
likely that another user has run into the same problem.
If you don't find an answer, please feel free to post a question.
If you have a feature request, please post in
[Feature Suggestions](https://discuss.konghq.com/c/feature-suggestions)
category.

If you've spotted a bug, please open an issue
on our [Github](https://github.com/kong/kubernetes-ingress-controller/issues).

## Documentation

All documentation around Kong Ingress Controller is present in this
repository. Pull Requests are welcome for additions and corrections.

Following are some helpful link:

- [**Getting Started**][deployment-doc]:
  Get Kubernetes Ingress setup up and running.
- [**Deployment**][deployment-doc]:
  Deployment guides for Minikube, GKE
  and other types of clusters.
- [**Custom Resources Definitions (CRDs**][crds]:
  Use custom resources
  to configure Kong in addition to the Ingress resource.
- [**Annotations**][annotations]:
  Associate plugins with your requests using annotations
- [**FAQs**][faqs]: Frequently Asked Questions.

## Design

Kong Ingress Controller is built to satisfy the [Ingress][ingress]
spec in Kubernetes.
Kong Ingress Controller is a [Go](https://golang.org/) app
that listens to events from the API-server of your Kubernetes cluster
and then sets up Kong to handle your configuration accordingly,
meaning you never have to configure Kong yourself manually.

The controller can configure any Kong cluster via a
Kong node running either in a control-plane mode
or running both, control and data planes.

For detailed design, please refer to our
[design][design] documentation.

## License

```text
Copyright 2018 Kong Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Roadmap

Please check the [roadmap][roadmap] document.

[ingress]: https://kubernetes.io/docs/concepts/services-networking/ingress/
[kong]: https://konghq.com/kong-community-edition/
[kong-hub]: https://docs.konghq.com/hub/
[deployment-doc]: deploy/README.md
[annotations]: docs/annotations.md
[crds]: docs/custom-types.md
[roadmap]: docs/roadmap.md
[design]: docs/design.md
[faqs]: docs/faq.md