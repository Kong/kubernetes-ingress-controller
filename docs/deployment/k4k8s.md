# Kong for Kubernetes

Kong for Kubernetes is an Ingress Controller based on the
Open-Source Kong Gateway. It consists of two components:

- **Kong**: the Open-Source Gateway
- **Controller**: a daemon process that integrates with the
  Kubernetes platform and configures Kong.

## Table of content

- [Installers](#installers)
    - [YAML manifests](#yaml-manifests)
    - [Kustomize](#kustomize)
    - [Helm](#helm)
- [Using Kong for Kubernetes](#using-kong-for-kubernetes)

## Installers

Kong for Kubernetes can be installed using an installer of
your choice.

Once you've installed Kong for Kubernetes,
jump to the [next section](#using-kong-for-kubernetes)
on using it.

### YAML manifests

Please pick one of the following guides depending on your platform:

- [Minikube](minikube.md)
- [Google Kubernetes Engine(GKE) by Google](gke.md)
- [Elastic Kubernetes Service(EKS) by Amazon](eks.md)
- [Azure Kubernetes Service(AKS) by Microsoft](aks.md)

### Kustomize

Use Kustomize to install Kong for Kubernetes:

```
kustomize build github.com/kong/kubernetes-ingress-controller/deploy/manifests/base
```

You can use the above URL as a base kustomization and build on top of it
to make it suite better for your cluster and use-case.

Once installed, set an environment variable, $PROXY_IP with the External IP address of
the `kong-proxy` service in `kong` namespace:

```
export PROXY_IP=$(kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" service -n kong kong-proxy)
```

### Helm

You can use Helm to install Kong via the official Helm chart:

```
helm repo update
helm install stable/kong --name demo --namespace kong
```

Once installed, set an environment variable, $PROXY_IP with the External IP address of
the `demo-kong-proxy` service in `kong` namespace:

```
export PROXY_IP=$(kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" service -n kong demo-kong-proxy)
```

## Using Kong for Kubernetes

Once you've installed Kong for Kubernetes, please follow our
[getting started](../guides/getting-started.md) tutorial to learn more.
