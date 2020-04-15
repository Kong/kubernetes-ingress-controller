# Kong Ingress on Elastic Kubernetes Service (EKS)

## Requirements

1. A fully functional EKS cluster.
   Please follow Amazon's Guide to
   [set up an EKS cluster](https://aws.amazon.com/getting-started/projects/deploy-kubernetes-app-amazon-eks/).
2. Basic understanding of Kubernetes
3. A working `kubectl`  linked to the EKS Kubernetes
   cluster we will work on. The above EKS setup guide will help
   you set this up.

## Deploy Kong Ingress Controller

Deploy Kong Ingress Controller using `kubectl`:

```bash
$ kubectl create -f https://bit.ly/k4k8s
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

It may take a few minutes for all containers to start and report
healthy statuses.

Alternatively, you can use our helm chart as well.
Please ensure that you have Tiller working and then execute:

```bash
$ helm repo add kong https://charts.konghq.com
$ helm repo update

# Helm 2
$ helm install kong/kong

# Helm 3
$ helm install kong/kong --generate-name --set ingressController.installCRDs=false
```

*Note:* this process could take up to five minutes the first time.

## Setup environment variables

Next, create an environment variable with the IP address at which
Kong is accesssible. This IP address sends requests to the
Kubernetes cluster.

Execute the following command to get the IP address at which Kong is accessible:

```bash
$ kubectl get services -n kong
NAME         TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                      AGE
kong-proxy   LoadBalancer   10.63.250.199   203.0.113.42   80:31929/TCP,443:31408/TCP   57d
```

Create an environment variable to hold the IP address:

```bash
$ export PROXY_IP=$(kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" service -n kong kong-proxy)
```

> Note: It may take some time for Amazon to actually associate the
IP address to the `kong-proxy` Service.

Once you've installed Kong Ingress Controller, please follow our
[getting started](../guides/getting-started.md) tutorial to learn
about how to use the Ingress Controller.

## TLS configuration

Versions of Kong prior to 2.0.0 default to using [the "modern" cipher suite
list](https://wiki.mozilla.org/Security/Server_Side_TLS). This is not
compatible with ELBs when the ELB terminates TLS at the edge and establishes a
new session with Kong. This error will appear in Kong's logs:

```
*7961 SSL_do_handshake() failed (SSL: error:1417A0C1:SSL routines:tls_post_process_client_hello:no shared cipher) while SSL handshaking
```

To correct this issue, set `KONG_SSL_CIPHER_SUITE=intermediate` in your
environment variables.
