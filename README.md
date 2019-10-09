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
- [**License**](#license)

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

The following matrix lists supported versions of
Kong for every release of the Kong Ingress Controller:

| Kong Ingress Controller  | <= 0.0.4           | 0.0.5              | 0.1.x              | 0.2.x              | 0.3.x              | 0.4.x              | 0.5.x              | 0.6.x              |
|--------------------------|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|
| Kong 0.13.x              | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                | :x:                | :x:                | :x:                | :x:                |
| Kong 0.14.x              | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                | :x:                | :x:                |
| Kong 1.0.x               | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Kong 1.1.x               | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Kong 1.2.x               | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Kong 1.3.x               | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: |
| Kong Enterprise 0.32-x   | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                | :x:                | :x:                | :x:                |
| Kong Enterprise 0.33-x   | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                | :x:                | :x:                | :x:                |
| Kong Enterprise 0.34-x   | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                | :x:                | :x:                | :x:                |
| Kong Enterprise 0.35-x   | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Kong Enterprise 0.36-x   | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :white_check_mark: |

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
# using YAMLs
kubectl apply -f https://bit.ly/kong-ingress

# or using Helm
helm install stable/kong --set ingressController.enabled=true
```

You can also spin up Kong Ingress Controller without a database dependency:

```
# using YAMLs
kubectl apply -f https://bit.ly/kong-ingress-dbless
# or using Helm
helm install stable/kong --set ingressController.enabled=true \
  --set postgresql.enabled=false --set env.database=off
```

It takes a few minutes for all components to
spin up.
You now have set up Kong as your Ingress controller and
all Ingress resources in your Kubernetes Cluster will be satisfied.

Please refer to our [deployment documentation][deployment]
for a detailed introduction to Kong Ingress Controller
and Ingress spec.

## Seeking help

Please search through the posts on
[Kong Nation](https://discuss.konghq.com/c/kubernetes) as it's
likely that another user has run into the same problem.
If you don't find an answer, please feel free to post a question.
For a feature request, please post in
[Feature Suggestions](https://discuss.konghq.com/c/feature-suggestions)
category.

You can also talk to us in the
[#kong](https://kubernetes.slack.com/messages/kong) channel on the
Kubernetes Slack server.

If you've spotted a bug, please open an issue
on our [Github](https://github.com/kong/kubernetes-ingress-controller/issues).

## Documentation

All documentation around Kong Ingress Controller is present in this
repository inside the [docs][docs] directory.
Pull Requests are welcome for additions and corrections.

Following are some helpful link:

- [**Getting Started**](docs/guides/getting-started.md):
  Get Kubernetes Ingress setup up and running.
- [**Deployment**][deployment]:
  Deployment guides for Minikube, GKE
  and other types of clusters.
- [**Custom Resources Definitions (CRDs)**][crds]:
  Use custom resources
  to configure Kong in addition to the Ingress resource.
- [**Annotations**][annotations]:
  Associate plugins with your requests using annotations
- [**FAQs**][faqs]: Frequently Asked Questions.

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

[ingress]: https://kubernetes.io/docs/concepts/services-networking/ingress/
[kong]: https://konghq.com/kong-community-edition/
[kong-hub]: https://docs.konghq.com/hub/
[docs]: docs/
[deployment]: docs/deployment/
[annotations]: docs/references/annotations.md
[crds]: docs/references/custom-resources.md
[faqs]: docs/faq.md
