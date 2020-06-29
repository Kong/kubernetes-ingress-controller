# Using Redis for rate-limiting

Kong can rate-limit your traffic without any external dependency.
In such a case, Kong stores the request counters in-memory
and each Kong node applies the rate-limiting policy independently.
There is no synchronization of information being done in this case.
But if Redis is available in your cluster, Kong
can take advantage of it and synchronize the rate-limit information
across multiple Kong nodes and enforce a slightly different rate-limiting
policy.

This guide walks through the steps of using Redis for rate-limiting in
a multi-node Kong deployment.

## Installation

Please follow the [deployment](../deployment) documentation to install
Kong Ingress Controller on your Kubernetes cluster.

## Testing Connectivity to Kong

This guide assumes that the `PROXY_IP` environment variable is
set to contain the IP address or URL pointing to Kong.
Please follow one of the
[deployment guides](../deployment) to configure this environment variable.

If everything is setup correctly, making a request to Kong should return
HTTP 404 Not Found.

```bash
$ curl -i $PROXY_IP
HTTP/1.1 404 Not Found
Date: Fri, 21 Jun 2019 17:01:07 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 48
Server: kong/1.2.1

{"message":"no Route matched with those values"}
```

This is expected as Kong does not yet know how to proxy the request.

## Setup a Sample Service

For the purpose of this guide, we will setup an [httpbin](https://httpbin.org)
service in the cluster and proxy it.

```bash
kubectl apply -f https://bit.ly/k8s-httpbin
service/httpbin created
deployment.apps/httpbin created
```

Create an Ingress rule to proxy the httpbin service we just created:

```bash
$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
  annotations:
    konghq.com/strip-path: "true"
spec:
  rules:
  - http:
      paths:
      - path: /foo
        backend:
          serviceName: httpbin
          servicePort: 80
' | kubectl apply -f -
ingress.extensions/demo created
```

Test the Ingress rule:

```bash
$ curl -i $PROXY_IP/foo/status/200
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Wed, 17 Jul 2019 19:25:32 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1
```

## Setup rate-limiting

We will start by creating a global rate-limiting policy:

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongClusterPlugin
metadata:
  name: global-rate-limit
  labels:
    global: \"true\"
config:
  minute: 5
  policy: local
plugin: rate-limiting
" | kubectl apply -f -
kongclusterplugin.configuration.konghq.com/global-rate-limit created
```

Here we are configuring Kong Ingress Controller to rate-limit traffic from
any client to 5 requests per minute, and we are applying this policy in a
global sense, meaning the rate-limit will apply across all services.

You can set this up for a specific Ingress or a specific service as well,
please follow [using KongPlugin resource](using-kongplugin-resource.md)
guide on steps for doing that.

Next, test the rate-limiting policy by executing the following command
multiple times and observe the rate-limit headers in the response:

```bash
$ curl -I $PROXY_IP/foo/headers
```

As there is a single Kong instance running, Kong correctly imposes
the rate-limit and you can make only 5 requests in a minute.

## Scale Kong Ingress Controller to multiple pods

Now, let's scale up the Kong Ingress controller deployment to 3 pods, for
scalability and redundancy:

```bash
$ kubectl scale --replicas 3 -n kong deployment ingress-kong
deployment.extensions/ingress-kong scaled
```

It will take a couple minutes for the new pods to start up.
Once the new pods are up and running, test the rate-limiting policy by 
executing the following command and observing the rate-limit headers:

```bash
$ curl -I $PROXY_IP/foo/headers
```

You will observe that the rate-limit is not consistent anymore
and you can make more than 5 requests in a minute.

To understand this behavior, we need to understand how we have configured Kong.
In the current policy, each Kong node is tracking a rate-limit in-memory
and it will allow 5 requests to go through for a client.
There is no synchronization of the rate-limit information across Kong nodes.
In use-cases where rate-limiting is used as a protection mechanism and to
avoid over-loading your services, each Kong node tracking its own counter
for requests is good enough as a malicious user will hit rate-limits on all
nodes eventually.
Or if the load-balancer in-front of Kong is performing some
sort of deterministic hashing of requests such that the same Kong node always
receives the requests from a client, then we won't have this problem at all.

In some cases, a synchronization of information that each Kong node maintains
in-memory is needed. For that purpose, Redis can be used.
Let's go ahead and set this up next.

## Deploy Redis to your Kubernetes cluster

First, we will deploy redis in our Kubernetes cluster:

```bash
$ kubectl apply -n kong -f https://bit.ly/k8s-redis
deployment.apps/redis created
service/redis created
```

Once this is deployed, let's update our KongClusterPlugin configuration to use
Redis as a datastore rather than each Kong node storing the counter information
in-memory:

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongClusterPlugin
metadata:
  name: global-rate-limit
  labels:
    global: \"true\"
config:
  minute: 5
  policy: redis
  redis_host: redis
plugin: rate-limiting
" | kubectl apply -f -
kongclusterplugin.configuration.konghq.com/global-rate-limit configured
```

Notice, how the `policy` is now set to `redis` and we have configured Kong
to talk to the `redis`  server available at `redis` DNS name, which is the
Redis node we deployed earlier.

## Test it

Now, if you go ahead and execute the following commands, you should be able
to make only 5 requests in a minute:

```bash
$ curl -I $PROXY_IP/foo/headers
```

This guide shows how to use Redis as a data-store for rate-limiting plugin,
but this can be used for other plugins which support Redis as a data-store
like proxy-cache.
