# Kong Ingress on Minikube

## Setup Minikube

1. Install [`minikube`](https://github.com/kubernetes/minikube)  
  
    Minikube is a tool that makes it easy to run Kubernetes locally.
    Minikube runs a single-node Kubernetes cluster inside a VM on your laptop
    for users looking to try out Kubernetes or develop with it day-to-day.

1. Start `minikube`

    ```bash
    minikube start
    ```

    It will take a few minutes to get all resources provisioned.

    ```bash
    kubectl get nodes
    ```

## Deploy Kong Ingress Controller

Deploy Kong Ingress Controller using `kubectl`:

```bash
$ curl -sL https://bit.ly/kong-ingress | kubectl create -f -
namespace/kong created
customresourcedefinition.apiextensions.k8s.io/kongplugins.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongconsumers.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongcredentials.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongingresses.configuration.konghq.com created
service/postgres created
statefulset.apps/postgres created
serviceaccount/kong-serviceaccount created
clusterrole.rbac.authorization.k8s.io/kong-ingress-clusterrole created
clusterrolebinding.rbac.authorization.k8s.io/kong-ingress-clusterrole-nisa-binding created
service/kong-ingress-controller created
deployment.extensions/kong-ingress-controller created
service/kong-proxy created
deployment.extensions/kong created
job.batch/kong-migrations created
```

Alternatively, you can use our helm chart as well.
Please ensure that you've Tiller working and then execute:

```bash
$ helm install stable/kong --set ingressController.enabled=true
```

> Note: this process could take up to five minutes the first time.

## Setup environment variables

Next, we will setup an environment variable with the IP address at which
Kong is accesssible. This will be used to actually send reqeusts into the
Kubernetes cluster.

```bash
$ export PROXY_IP=$(minikube service -n kong kong-proxy --url | head -1)
$ echo $PROXY_IP
http://192.168.99.100:32728
```

Once you've Kong Ingress Controlled installed, please follow our
[getting started](../guides/getting-started.md) tutorial to learn
about how to use the Ingress Controller.
