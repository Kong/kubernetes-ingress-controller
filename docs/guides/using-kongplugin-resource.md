# Using KongPlugin resource

In this guide, we will learn how to use KongPlugin resource to configure
plugins in Kong to modify requests for a specific request path.
The guide will cover configuring a plugin for a specific service, a set of Ingress rules
and for a specific user of the API.

## Installation

Please follow the [deployment](../deployment) documentation to install
Kong Ingress Controller onto your Kubernetes cluster.

## Testing connectivity to Kong

This guide assumes that `PROXY_IP` environment variable is
set to contain the IP address or URL pointing to Kong.
If you've not done so, please follow one of the
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

## Installing sample services

We will start by installing two services,
an echo service and an httpbin service.

```bash
$ kubectl apply -f https://bit.ly/k8s-httpbin
service/httpbin created
deployment.apps/httpbin created
```

```bash
$ kubectl apply -f https://bit.ly/echo-service
service/echo created
deployment.apps/echo created
```

## Setup Ingress rules

Let's expose these services outside the Kubernetes cluster
by defining Ingress rules.

```bash
echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
spec:
  rules:
  - http:
      paths:
      - path: /foo
        backend:
          serviceName: httpbin
          servicePort: 80
      - path: /bar
        backend:
          serviceName: echo
          servicePort: 80
" | kubectl apply -f -
ingress.extensions/demo created
```

Let's test these endpoints:

```bash
$ curl -i $PROXY_IP/foo/status/200
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Wed, 17 Jul 2019 21:38:00 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1

$ curl -i $PROXY_IP/bar
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Transfer-Encoding: chunked
Connection: keep-alive
Date: Wed, 17 Jul 2019 21:38:17 GMT
Server: echoserver
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1



Hostname: echo-d778ffcd8-n9bss

Pod Information:
    node name:  gke-harry-k8s-dev-default-pool-bb23a167-8pgh
    pod name:  echo-d778ffcd8-n9bss
    pod namespace:  default
    pod IP:  10.60.0.4
<-- clipped -- >
```

Let's add another Ingress resource which proxies requests to `/baz` to httpbin
service:

```bash
echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-2
spec:
  rules:
  - http:
      paths:
      - path: /baz
        backend:
          serviceName: httpbin
          servicePort: 80
" | kubectl apply -f -
ingress.extensions/demo-2 created
```

We will use this path later.

## Configuring plugins on Ingress resource

Next, we will configure two plugins on the Ingress resource.

First, we will create a KongPlugin resource:

```bash
$ echo '
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: add-response-header
config:
  add:
    headers:
    - "demo: injected-by-kong"
plugin: response-transformer
' | kubectl apply -f -
kongplugin.configuration.konghq.com/add-response-header created
```

Next, we will associate it with our Ingress rules:

```bash
$ kubectl patch ingress demo -p '{"metadata":{"annotations":{"plugins.konghq.com":"add-response-header"}}}'
ingress.extensions/demo patched
```

Here, we are asking Kong Ingress Controller to execute the response-transformer
plugin whenever a request matching the Ingress rule is processed.

Let's test it out:

```bash
curl -i $PROXY_IP/foo/status/200
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 9593
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Wed, 17 Jul 2019 21:54:31 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
demo:  injected-by-kong
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1

$ curl -I $PROXY_IP/bar
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Wed, 17 Jul 2019 21:54:39 GMT
Server: echoserver
demo:  injected-by-kong
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1
```

As can be seen in the output, the `demo` header is injected by Kong when
the request matches the Ingress rules defined in the `demo` Ingress resource.

If we send a request to `/baz`, then we can see that the header is not injected
by Kong:

```bash
curl -I $PROXY_IP/baz
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 9593
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Wed, 17 Jul 2019 21:56:20 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 3
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1
```

Here, we have successfully setup a plugin which is executed only when a
request matches a specific `Ingress` rule.

## Configuring plugins on Service resource

Next, we will see how we can configure Kong to execute plugins for requests
which are sent to a specific service.

Let's add a `KongPlugin` resource for authentication on the httpbin service:

```bash
$ echo "apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: httpbin-auth
plugin: key-auth
" | kubectl apply -f -

kongplugin.configuration.konghq.com/httpbin-auth created
```

Next, we will associate this plugin to the httpbin service running in our
cluster:

```bash
kubectl patch service httpbin -p '{"metadata":{"annotations":{"plugins.konghq.com":"httpbin-auth"}}}'
service/httpbin patched
```

Now, any request sent to the service will require authentication,
no matter which `Ingress` rule it matched:

```bash
$ curl -I $PROXY_IP/baz
HTTP/1.1 401 Unauthorized
Date: Wed, 17 Jul 2019 22:09:04 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
WWW-Authenticate: Key realm="kong"
Content-Length: 41
Server: kong/1.2.1

$ curl -I $PROXY_IP/foo
HTTP/1.1 401 Unauthorized
Date: Wed, 17 Jul 2019 22:12:13 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
WWW-Authenticate: Key realm="kong"
Content-Length: 41
demo:  injected-by-kong
Server: kong/1.2.1
```

You can also see how the `demo` header was injected as the request also
matched one of the rules defined in the `demo` `Ingress` resource.

## Configure consumer and credential

Follow the [Using Consumers and Credentials](using-consumer-credential-resource.md)
guide to provision a user and an apikey.
Once you have it, please continue:

Use the API key to pass authentication:

```bash
$ curl -I $PROXY_IP/baz -H 'apikey: sooper-secret-key'
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 9593
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Wed, 17 Jul 2019 22:16:35 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1

$ curl -I $PROXY_IP/foo -H 'apikey: sooper-secret-key'
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 9593
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Wed, 17 Jul 2019 22:15:34 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
demo:  injected-by-kong
X-Kong-Upstream-Latency: 3
X-Kong-Proxy-Latency: 0
Via: kong/1.2.1
```

## Configure a global plugin

Now, we will protect our Kubernetes cluster.
For this, we will be configuring a rate-limiting plugin, which
will throttle requests coming from the same client.

Let's create the `KongPlugin` resource:

```bash
echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: global-rate-limit
  labels:
    global: \"true\"
config:
  minute: 5
  limit_by: consumer
  policy: local
plugin: rate-limiting
" | kubectl apply -f -
kongplugin.configuration.konghq.com/global-rate-limit created
```

With this plugin (please note the `global` label), every request through
Kong Ingress Controller will be rate-limited:

```bash
$ curl -I $PROXY_IP/foo -H 'apikey: sooper-secret-key'
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 9593
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Wed, 17 Jul 2019 22:34:10 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-RateLimit-Limit-minute: 5
X-RateLimit-Remaining-minute: 4
demo:  injected-by-kong
X-Kong-Upstream-Latency: 3
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1

$ curl -I $PROXY_IP/bar
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Wed, 17 Jul 2019 22:34:14 GMT
Server: echoserver
X-RateLimit-Limit-minute: 5
X-RateLimit-Remaining-minute: 4
demo:  injected-by-kong
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1
```

## Configure a plugin for a specific consumer

Now, let's say we would like to give a specific consumer a higher rate-limit.

For this, we can create a `KongPlugin` resource and then associate it with
a specific consumer.

First, create the `KongPlugin` resource:

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: harry-rate-limit
config:
  minute: 10
  limit_by: consumer
  policy: local
plugin: rate-limiting
" | kubectl apply -f -
kongplugin.configuration.konghq.com/harry-rate-limit created
```

Next, associate this with the consumer:

```bash
echo "apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: harry
  annotations:
    plugins.konghq.com: harry-rate-limit
username: harry" | kubectl apply -f -
kongconsumer.configuration.konghq.com/harry configured
```

Note the annotation being added to the `KongConsumer` resource.

Now, if the request is made as the `harry` consumer, the client
will be rate-limited differently:

```bash
$ curl -I $PROXY_IP/foo -H 'apikey: sooper-secret-key'
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 9593
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Wed, 17 Jul 2019 22:34:10 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-RateLimit-Limit-minute: 10
X-RateLimit-Remaining-minute: 9
demo:  injected-by-kong
X-Kong-Upstream-Latency: 3
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1

# a regular unauthenticated request
$ curl -I $PROXY_IP/bar
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Wed, 17 Jul 2019 22:34:14 GMT
Server: echoserver
X-RateLimit-Limit-minute: 5
X-RateLimit-Remaining-minute: 4
demo:  injected-by-kong
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1
```

This guide demonstrates how you can use Kong Ingress Controller to
impose restrictions and transformations
on various levels using Kubernetes style APIs.
