# Kong for Kubernetes Enterprise

Kong for Kubernetes Enterprise is an enhanced version of
the Open-Source Ingress Controller. It includes all
Enterprise plugins and comes with 24x7 support for worry-free
production deployment.
This is available to enterprise customers of Kong, Inc. only.

## Table of content

- [Prerequisites](#prerequisites)
- [Installers](#installers)
    - [YAML manifests](#yaml-manifests)
    - [Kustomize](#kustomize)
    - [Helm](#helm)
- [Using Kong for Kubernetes Enterprise](#using-kong-for-kubernetes-enterprise)

## Prerequisites

Before we can deploy Kong, we need to satisfy two
prerequisites:

- [Kong Enterprise License secret](#kong-enterprise-license-secret)
- [Kong Enterprise Docker registry access](#kong-enterprise-docker-registry-access)

In order to create these secrets, let's provision the `kong`
namespace first:

```bash
$ kubectl create namespace kong
namespace/kong created
```

### Kong Enterprise License secret

Enterprise version requires a valid license to run.  
As part of sign up for Kong Enterprise, you should have received a license file.
If you do not have one, please contact your sales representative.
Save the license file temporarily to disk with filename `license`
and execute the following:

```bash
$ kubectl create secret generic kong-enterprise-license --from-file=./license -n kong
secret/kong-enterprise-license created
```

Please note:

- There is no `.json` extension in the `--from-file` parameter.
- `-n kong` specifies the namespace in which you are deploying
  Kong Ingress Controller. If you are deploying in a different namespace,
  please change this value.

### Kong Enterprise Docker registry access

Next, we need to setup Docker credentials in order to allow Kubernetes
nodes to pull down Kong Enterprise Docker image, which is hosted as a private
repository.
As part of your sign up for Kong Enterprise, you should have received
credentials to access Enterprise Bintray repositories.
Your username is the same username you use
to log in to Bintray and password
is an API-key that can be provisioned via Bintray.

```bash
$ kubectl create secret -n kong docker-registry kong-enterprise-k8s-docker \
    --docker-server=kong-docker-kong-enterprise-k8s.bintray.io \
    --docker-username=<your-bintray-username@kong> \
    --docker-password=<your-bintray-api-key>
secret/kong-enterprise-k8s-docker created
```

Again, please take a note of the namespace `kong`.

## Installers

Once the secrets are in-place, we can proceed with installation.

Kong for Kubernetes can be installed using an installer of
your choice:

### YAML manifests

Execute the following to install Kong for Kubernetes Enteprise using YAML
manifests:

```bash
$ kubectl apply -f https://bit.ly/k4k8s-enterprise
```

It takes a few minutes the first time this setup is done.

```bash
$ kubectl get pods -n kong
NAME                            READY   STATUS    RESTARTS   AGE
ingress-kong-6ffcf8c447-5qv6z   2/2     Running   1          44m
```

You can also see the `kong-proxy` service:

```bash
$ kubectl get service kong-proxy -n kong
NAME         TYPE           CLUSTER-IP     EXTERNAL-IP     PORT(S)                      AGE
kong-proxy   LoadBalancer   10.63.254.78   35.233.198.16   80:32697/TCP,443:32365/TCP   22h
```

> Note: Depending on the Kubernetes distribution you are using, you might or might
not see an external IP address assigned to the service. Please see
your provider's guide on obtaining an IP address for a Kubernetes Service of
type `LoadBalancer`.

Let's setup an environment variable to hold the IP address:

```bash
$ export PROXY_IP=$(kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" service -n kong kong-proxy)
```

> Note: It may take a while for your cloud provider to actually associate the
IP address to the `kong-proxy` Service.

### Kustomize

Use Kustomize to install Kong for Kubernetes Enterprise:

```
kustomize build github.com/kong/kubernetes-ingress-controller/deploy/manifests/enterprise-k8s
```

You can use the above URL as a base kustomization and build on top of it
as well.

Once installed, set an environment variable, $PROXY_IP with the External IP address of
the `kong-proxy` service in `kong` namespace:

```
export PROXY_IP=$(kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" service -n kong kong-proxy)
```

### Helm

You can use Helm to install Kong via the official Helm chart:

```
helm repo add kong https://charts.konghq.com
helm repo update
helm install kong/kong --name demo --namespace kong --values https://fstlnk.in/k4k8s-enterprise-helm-values
```

Once installed, set an environment variable, $PROXY_IP with the External IP address of
the `demo-kong-proxy` service in `kong` namespace:

```
export PROXY_IP=$(kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" service -n kong demo-kong-proxy)
```

## Using Kong for Kubernetes Enterprise

Once you've installed Kong for Kubernetes Enterprise, please follow our
[getting started](../guides/getting-started.md) tutorial to learn more.
