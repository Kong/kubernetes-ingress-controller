# Setting up Active and Passive health checks

In this guide, we will go through steps necessary to setup active and passive
health checking using Kong Ingress Controller. This configuration allows
Kong to automatically short-circuit requests to specific Pods that are
mis-behaving in your Kubernetes Cluster.

> Please make sure to use Kong Ingress Controller >= 0.6 as the previous
versions contain a [bug](https://github.com/hbagdi/go-kong/issues/6).

## Installation

Please follow the [deployment](../deployment) documentation to install
Kong Ingress Controller onto your Kubernetes cluster.

## Testing connectivity to Kong

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

This is expected since Kong doesn't know how to proxy any requests yet.

## Setup a Sample Service

For the purpose of this guide, we will setup an [httpbin](https://httpbin.org)
service in the cluster and proxy it.

```bash
$ kubectl apply -f https://bit.ly/k8s-httpbin
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
    kubernetes.io/ingress.class: kong
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

Observe the headers and you can see that Kong has proxied the request correctly.

## Setup passive health checking

Now, let's setup passive HTTP health-check for our service.
All health-checking is done at Service-level and not Ingress-level.

Add the following KongIngress resource:

```bash
$ echo "apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
    name: demo-health-checking
upstream:
  healthchecks:
    passive:
      healthy:
        successes: 3
      unhealthy:
        http_failures: 3" | kubectl apply -f -
kongingress.configuration.konghq.com/demo-health-checking created
```

Here, we are configuring Kong to short-circuit requests to a pod
if a pod throws 3 consecutive errors.

Next, associate the KongIngress resource with `httpbin` service:

```bash
$ kubectl patch svc httpbin -p '{"metadata":{"annotations":{"konghq.com/override":"demo-health-checking"}}}'
service/httpbin patched
```

Now, let's send some traffic to test if this works:

Let's send 2 requests that represent a failure from upstream
and then send a request for 200.
Here we are using `/status/500` to simulate a failure from upstream.

```bash
$ curl -i $PROXY_IP/foo/status/500
HTTP/1.1 500 INTERNAL SERVER ERROR
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 05 Aug 2019 22:38:24 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 0
Via: kong/1.2.1

$ curl -i $PROXY_IP/foo/status/500
HTTP/1.1 500 INTERNAL SERVER ERROR
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 05 Aug 2019 22:38:24 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 0
Via: kong/1.2.1

$ curl -i $PROXY_IP/foo/status/200
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 05 Aug 2019 22:38:26 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1

```

Kong has not short-circuited because there were only two failures.
Let's send 3 requests and open the circuit, and then send a normal request.

```bash
$ curl -i $PROXY_IP/foo/status/500
HTTP/1.1 500 INTERNAL SERVER ERROR
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 05 Aug 2019 22:38:24 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 0
Via: kong/1.2.1

$ curl -i $PROXY_IP/foo/status/500
HTTP/1.1 500 INTERNAL SERVER ERROR
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 05 Aug 2019 22:38:24 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 0
Via: kong/1.2.1

$ curl -i $PROXY_IP/foo/status/500
HTTP/1.1 500 INTERNAL SERVER ERROR
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 05 Aug 2019 22:38:24 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 0
Via: kong/1.2.1

curl -i $PROXY_IP/foo/status/200
HTTP/1.1 503 Service Temporarily Unavailable
Date: Mon, 05 Aug 2019 22:41:19 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 58
Server: kong/1.2.1

{"message":"failure to get a peer from the ring-balancer"}

```

As we can see, Kong returns back a 503, representing that the service is
unavailable. Since we have only one pod of httpbin running in our cluster,
and that is throwing errors, Kong will not proxy anymore requests.

Now we have a few options:

- Delete the current httpbin pod; Kong will then proxy requests to the new
  pod that comes in its place.
- Scale the httpbin deployment; Kong will then proxy requests to the new
  pods and leave the short-circuited pod out of the loop.
- Manually change the pod health status in Kong using Kong's Admin API.

These options highlight the fact that once a circuit is opened because of
errors, there is no way for Kong to close the circuit again.

This is a feature which some services might need, where once a pod starts
throwing errors, manual intervention is necessary before that pod can
again handle requests.
To get around this, we can introduce active health-check, where each instance
of Kong actively probes pods to figure out if they are healthy or not.

## Setup active health checking

Let's update our KongIngress resource to use active health-checks:

```bash
$ echo "apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
    name: demo-health-checking
upstream:
  healthchecks:
    active:
      healthy:
        interval: 5
        successes: 3
      http_path: /status/200
      type: http
      unhealthy:
        http_failures: 1
        interval: 5
    passive:
      healthy:
        successes: 3
      unhealthy:
        http_failures: 3" | kubectl apply -f -
kongingress.configuration.konghq.com/demo-health-checking configured
```

Here, we are configuring Kong to actively probe `/status/200` every 5 seconds.
If a pod is unhealthy (from Kong's perspective),
3 successful probes will change the status of the pod to healthy and Kong
will again start to forward requests to that pod.

Now, the requests should flow once again:

```bash
$ curl -i $PROXY_IP/foo/status/200
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 05 Aug 2019 22:38:26 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1

```

Let's trip the circuit again by sending three requests that will return
500s from httpbin:

```bash
$ curl -i $PROXY_IP/foo/status/500
$ curl -i $PROXY_IP/foo/status/500
$ curl -i $PROXY_IP/foo/status/500
```

Now, sending the following request will fail for about 15 seconds,
the duration it will take active healthchecks to re-classify
the httpbin pod as healthy again.

```bash
$ curl -i $PROXY_IP/foo/status/200
HTTP/1.1 503 Service Temporarily Unavailable
Date: Mon, 05 Aug 2019 23:17:47 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 58
Server: kong/1.2.1

{"message":"failure to get a peer from the ring-balancer"}
```

After 15 seconds, you will see:

```bash
$ curl -i $PROXY_IP/foo/status/200
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 05 Aug 2019 22:38:26 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1
```

As we can see, active health-checks automatically marked a pod as healthy
when passive health-checks marked it unhealthy.

## Bonus

Scale the `httpbin` and `ingress-kong` deployments and observe how
multiple pods change the outcome of the above demo. 

Read more about health-checks and ciruit breaker in Kong's
[documentation](https://docs.konghq.com/latest/health-checks-circuit-breakers).
