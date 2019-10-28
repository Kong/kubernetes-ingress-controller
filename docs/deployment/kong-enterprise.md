# Kong Ingress Controller with Kong Enterprise

This guide walks through setting up Kong Ingress Controller using Kong
Enterprise.

Kong Ingress Controller is compatible with Kong Enterprise.
For version compatibility, please checkout the
[version matrix](../../README.md#version-support-matrix).

Before we can deploy the Ingress Controller, we need to satisfy two
prerequisites:

- [Kong Enterprise License secret](#kong-enterprise-license-secret)
- [Kong Enterprise Docker registry access](#kong-enterprise-docker-registry-access)

## Kong Enterprise License secret

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

## Kong Enterprise Docker registry access

Next, we need to setup Docker credentials in order to allow Kubernetes
nodes to pull down Kong Enterprise Docker image, which is hosted as a private
repository.
As part of your sign up for Kong Enterprise, you should have received credentials
for these as well.

```bash
$ kubectl create secret -n kong docker-registry kong-enterprise-docker \
    --docker-server=kong-docker-kong-enterprise-edition-docker.bintray.io \
    --docker-username=<your-username> \
    --docker-password=<your-password>
secret/kong-enterprise-docker created
```

Once these are created, we are ready to deploy Kong Enterprise
Ingress Controller.

## Deploy the Kong Ingress Controller

```bash
$ kubectl apply -f https://bit.ly/kong-ingress-enterprise
```

It takes a little while to bootstrap the database.
Once bootstrapped, you should see Kong Ingress Controller running with
Kong Enterprise as it's core:

```bash
$ kubectl get pods -n kong
NAME                            READY   STATUS      RESTARTS   AGE
ingress-kong-79784d576d-fbvs4   2/2     Running     1          10h
ingress-kong-79784d576d-sszcc   2/2     Running     2          10h
kong-migrations-6xvst           0/1     Completed   0          20h
postgres-0                      1/1     Running     0          20h
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

> Note: It may take a while for Google to actually associate the
IP address to the `kong-proxy` Service.

Once you've installed Kong Ingress Controller, please follow our
[getting started](../guides/getting-started.md) tutorial to learn
about how to use the Ingress Controller.
