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
$ curl -sL https://bit.ly/k4k8s | kubectl create -f -
namespace/kong created
customresourcedefinition.apiextensions.k8s.io/kongplugins.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongconsumers.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongcredentials.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongingresses.configuration.konghq.com created
serviceaccount/kong-serviceaccount created
clusterrole.rbac.authorization.k8s.io/kong-ingress-clusterrole created
clusterrolebinding.rbac.authorization.k8s.io/kong-ingress-clusterrole-nisa-binding created
configmap/kong-server-blocks created
service/kong-proxy created
service/kong-validation-webhook created
deployment.extensions/kong created
```

Alternatively, you can use our helm chart as well.
Please ensure that you've Tiller working and then execute:

```bash
$ helm repo add kong https://charts.konghq.com
$ helm repo update
$ helm install kong/kong
```

> Note: this process could take up to five minutes the first time.

## Setup environment variables

Next, we will setup an environment variable with the IP address at which
Kong is accesssible. This will be used to actually send reqeusts into the
Kubernetes cluster.

```bash
$ export PROXY_IP=$(minikube service kong-kong-proxy --url | head -1)
# If installed by helm without specifying namespace, use the line below instead.
# $ export PROXY_IP=$(minikube service -n kong kong-kong-proxy --url | head -1)
$ echo $PROXY_IP
http://192.168.99.100:32728
```

Once you've installed Kong Ingress Controller, please follow our
[getting started](../guides/getting-started.md) tutorial to learn
about how to use the Ingress Controller.
