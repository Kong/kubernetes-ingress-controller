# Kong for Kubernetes

[![Build Status][badge-travis-image]][badge-travis-url]

Use [Kong][kong] for Kubernetes [Ingress][ingress].  
Configure [plugins][kong-hub], health checking,
load balancing and more in Kong
for Kubernetes Services, all using
Custom Resource Definitions(CRDs) and Kubernetes-native tooling.

## Tables of content

- [**Features**](#features)
- [**Get started**](#get-started)
- [**Documentation**](#documentation)
- [**Version support matrix**](#version-support-matrix)
- [**Master branch builds**](#master-branch-builds)
- [**Seeking help**](#seeking-help)
- [**License**](#license)

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

# Helm 2
$ helm install kong/kong

# Helm 3
$ helm install kong/kong --generate-name --set ingressController.installCRDs=false
```

If you are setting up Kong for Kubernetes Enterprise, please
follow along [this guide](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/deployment/k4k8s-enterprise.md).

Follow the [Getting Started guide][getting-started-guide] to start
using Ingress in Kubernetes.

## Documentation

All documentation around Kong Ingress Controller is present in this
repository inside the [docs][docs] directory.
Pull Requests are welcome for additions and corrections.

## Version support matrix

[Version compatibility doc](docs/references/version-compatibility.md)
details on compatibility between versions of the
controller and versions of Kong, Kong for Kubernetes Enterprise and
Kong Enterprise.

## Master branch builds

If you would like to use the latest and the greatest version of the controller,
you can use `latest` tag from the [master repository][bintray-master-builds]
hosted on Bintray:

```
docker pull kong-docker-kubernetes-ingress-controller.bintray.io/master:latest
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

## License

```text
Copyright 2018-2020 Kong Inc.

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
[getting-started-guide]: docs/guides/getting-started.md
[badge-travis-image]: https://travis-ci.org/Kong/kubernetes-ingress-controller.svg?branch=master
[badge-travis-url]: https://travis-ci.org/Kong/kubernetes-ingress-controller
[bintray-master-builds]: https://bintray.com/kong/kubernetes-ingress-controller/master
