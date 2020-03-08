# Kong for Kubernetes with Kong Enterprise

This guide walks through setting up Kong Ingress Controller using Kong
Enterprise. This architecture is described in detail in [this doc](../concepts/k4k8s-with-kong-enterprise.md).

We assume that we start from scratch and you don't have Kong Enterprise
deployed. For the sake of simplicity, we will deploy Kong Enterprise and
it's database in Kubernetes itself. You can safely run them outside
Kubernetes as well.

## Table of content

- [Prerequisites](#prerequisites)
- [Install](#install)
- [Using Kong for Kubernetes](#using-kong-for-kubernetes-with-kong-enterprise)

## Prerequisites

Before we can deploy Kong Ingress Controller with Kong Enterprise,
we need to satisfy the following prerequisites:
- [Kong Enterprise License secret](#kong-enterprise-license-secret)
- [Kong Enterprise Docker registry access](#kong-enterprise-docker-registry-access)
- [Kong Enterprise bootstrap password](#kong-enterprise-bootstrap-password)

In order to create these secrets, let's provision the `kong`
namespace first:
```bash
$ kubectl create namespace kong
namespace/kong created
```

### Kong Enterprise License secret

Kong Enterprise requires a valid license to run.
As part of sign up for Kong Enterprise, you should have received a license file.
Save the license file temporarily to disk and execute the following:

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
$ kubectl create secret -n kong docker-registry kong-enterprise-edition-docker \
    --docker-server=kong-docker-kong-enterprise-edition-docker.bintray.io \
    --docker-username=<your-bintray-username@kong> \
    --docker-password=<your-bintray-api-key>
secret/kong-enterprise-edition-docker created
```

### Kong Enterprise bootstrap password

Next, we need to create a secret containing the password using which we can login into Kong Manager.
Please replace `cloudnative` with a random password of your choice and note it down.

```bash
kubectl create secret generic kong-enterprise-superuser-password  -n kong --from-literal=password=cloudnative

```

Once these are created, we are ready to deploy Kong Enterprise
Ingress Controller.

## Install

```bash
$ kubectl apply -f https://bit.ly/kong-ingress-enterprise
```

It takes a little while to bootstrap the database.
Once bootstrapped, you should see Kong Ingress Controller running with
Kong Enterprise as it's core:

```bash
$ kubectl get pods -n kong
NAME                            READY   STATUS      RESTARTS   AGE
ingress-kong-548b9cff98-n44zj   2/2     Running     0          21s
kong-migrations-pzrzz           0/1     Completed   0          4m3s
postgres-0                      1/1     Running     0          4m3s
```

You can also see the `kong-proxy` service:

```bash
$ kubectl get services -n kong
NAME                      TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                      AGE
kong-admin                LoadBalancer   10.63.255.85    34.83.95.105    80:30574/TCP                 4m35s
kong-manager              LoadBalancer   10.63.247.16    34.83.242.237   80:31045/TCP                 4m34s
kong-proxy                LoadBalancer   10.63.242.31    35.230.122.13   80:32006/TCP,443:32007/TCP   4m34s
kong-validation-webhook   ClusterIP      10.63.240.154   <none>          443/TCP                      4m34s
postgres                  ClusterIP      10.63.241.104   <none>          5432/TCP                     4m34s

```

> Note: Depending on the Kubernetes distribution you are using, you might or might
not see an external IP assigned to the three LoadBalancer type services. Please see
your provider's guide on obtaining an IP address for a Kubernetes Service of
type `LoadBalancer`. If you are running Minikube, you will not get an
external IP address.

### Setup Kong Manager

Next, if you browse to the IP address or host of the `kong-manager` service in your Browser,
which in our case is `http://34.83.242.237`.
Kong Manager should load in your browser.
Try logging in to the Manager with the username `kong_admin`
and the password you supplied in the prerequisite, it should fail.
The reason being we've not yet told Kong Manager where it can find the Admin API.

Let's set that up. We will take the External IP address of `kong-admin` service and
set the environment variable `KONG_ADMIN_API_URI`:

```bash
KONG_ADMIN_IP=$(kubectl get svc -n kong kong-admin --output=jsonpath='{.status.loadBalancer.ingress[0].ip}')
kubectl patch deployment -n kong ingress-kong -p "{\"spec\": { \"template\" : { \"spec\" : {\"containers\":[{\"name\":\"proxy\",\"env\": [{ \"name\" : \"KONG_ADMIN_API_URI\", \"value\": \"${KONG_ADMIN_IP}\" }]}]}}}}"
```

It will take a few minutes to roll out the updated deployment and once the new
`ingress-kong` pod is up and running, you should be able to log into the Kong Manager UI.

As you follow along with other guides on how to use your newly deployed Kong Ingress Controller,
you will be able to browse Kong Manager and see changes reflectded in the UI as Kong's
configuration changes.

## Using Kong for Kubernetes with Kong Enterprise

Let's setup an environment variable to hold the IP address of `kong-proxy` service:

```bash
$ export PROXY_IP=$(kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" service -n kong kong-proxy)
```

Once you've installed Kong for Kubernetes Enterprise, please follow our
[getting started](../guides/getting-started.md) tutorial to learn more.

## Customizing by use-case

The deployment in this guide is a point to start using Ingress Controller.
Based on your existing architecture, this deployment will require custom
work to make sure that it needs all of your requirements.

In this guide, there are three load-balancers deployed for each of
Kong Proxy, Kong Admin and Kong Manager services. It is possible and
recommended to instead have a single Load balancer and then use DNS names
and Ingress resources to expose the Admin and Manager services outside
the cluster.
