# Kong Ingress Controller Design

## Overview

Kong Ingress Controller is a dynamic and
highly available Ingress Controller which configures Kong
using Ingress resources created in your Kubernetes cluster.  
In addition, it can configure plugins, load balancing, health checking
on your services running in Kubernetes.

## Deployment

Ingress Controller is a Golang application that talks to Kubernetes API server
and [translates](#translation) Kubernetes resources into Kong.

Kong Ingress controller can be deployed in any Kubernetes cluster
which is running Kong.
Ingress Controller does not pose any limitations on how Kong is deployed
in your Kubernetes environment.
All it needs, is access to Kubernetes API server and
Admin API of Kong which it uses to configure Kong.

The Admin API of Kong can be running on
Data-Plane Kong nodes or Kong nodes running Control-Plane and Data-plane,
both at the same time.

![kong components](images/deployment.png "Kong Components")

In the above deployment figure, a Kong Control-Plane pod is
deployed alongside the Ingress Controller pod.
As mentioned above, this is only one way of deploying Kong with Ingress Controller
and it won't matter how we do it.

Kong's state is stored in Postgres (can be Cassandra) which should be deployed
as a StatefulSet and all Kong nodes in your Kubernetes cluster should be
able to connect to the database.

Please check out [Deployment Guides](deployment/) for more
details on how to deploy Kong with Ingress Controller.

## High Availability

Multiple instances of Kong Ingress Controller pod can be deployed.
The Ingress Controller uses a leader election protocol and elects a leader.
At any point, only one leader Controller pod will be configuring Kong and
other follower pods will be ready to take over as soon as the leader fails.

## Scaling Kong

If Kong is deployed in Control-plane and Data-plane mode, then
Kong proxy can be scaled independently.
By using this approach we can deploy and scale the data-plane
with the requirements of your applications,
i.e. using a DaemonSet, a deployment with affinity rules,
HorizontalPodAutoscaler etc.

## Translation

Kong Ingress Controller is a translator from
Kubernetes resources to resources in Kong.

Following figure shows how resources in Kubernetes are translated
into Kong's data model:

![translating k8s to kong](images/k8s-to-kong.png "Translating k8s resources to Kong")

The figure shows the translation of Ingress resource, services and
pods in Kubernetes to corresponding resources in Kong.

## Custom Resources

There are a few [Custom Resource Definitions(CRDs)][k8s-crd] available
that can be used to configure plugins, health checks, load balancing,
consumers, their credentials in Kong using the Ingress Controller.

Please read through our [Custom Resources][custom-resources]
documentation for details.

[k8s-deployment]: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
[k8s-initcontainer]: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
[k8s-crd]: https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/
[custom-resources]: custom-resources.md