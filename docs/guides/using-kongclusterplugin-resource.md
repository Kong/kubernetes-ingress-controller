# Using KongPlugin resource

In this guide, we will learn how to use KongClusterPlugin resource to configure
plugins in Kong.
The guide will cover configuring a plugin for services across different
namespaces.

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
an echo service and an httpbin service in their corresponding namespaces.

```bash
$ kubectl create namespace httpbin
namespace/httpbin created
$ kubectl apply -n httpbin -f https://bit.ly/k8s-httpbin
service/httpbin created
deployment.apps/httpbin created
```

```bash
$ kubectl create namespace echo
namespace/echo created
$ kubectl apply -n echo -f https://bit.ly/echo-service
service/echo created
deployment.apps/echo created
```

## Setup Ingress rules

Let's expose these services outside the Kubernetes cluster
by defining Ingress rules.

```bash
$ echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: httpbin-app
  namespace: httpbin
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
" | kubectl apply -f -
ingress.extensions/demo created

$ echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: echo-app
  namespace: echo
  annotations:
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - http:
      paths:
      - path: /bar
        backend:
          serviceName: echo
          servicePort: 80
" | kubectl apply -f -
ingress.extensions/demo created
```

Let's test these endpoints:

```bash
# access httpbin service
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

# access echo service
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

## Create KongClusterPlugin resource

```bash
$ echo '
apiVersion: configuration.konghq.com/v1
kind: KongClusterPlugin
metadata:
  name: add-response-header
  annotations:
    kubernetes.io/ingress.class: kong
config:
  add:
    headers:
    - "demo: injected-by-kong"
plugin: response-transformer
' | kubectl apply -f -
kongclusterplugin.configuration.konghq.com/add-response-header created
```

Note how the resource is created at cluster-level and not in any specific
namespace:

```bash
$ kubectl get kongclusterplugins
NAME                  PLUGIN-TYPE            AGE
add-response-header   response-transformer   4s
```

If you send requests to `PROXY_IP` now, you will see that the header is not
injected in the responses. The reason being that we have created a
resource but we have not told Kong when to execute the plugin.

## Configuring plugins on Ingress resources

We will associate the `KongClusterPlugin` resource with the two Ingress resources
that we previously created:

```bash
$ kubectl patch ingress -n httpbin httpbin-app -p '{"metadata":{"annotations":{"plugins.konghq.com":"add-response-header"}}}'
ingress.extensions/httpbin-app patched

$ kubectl patch ingress -n echo echo-app -p '{"metadata":{"annotations":{"plugins.konghq.com":"add-response-header"}}}'
ingress.extensions/echo-app patched
```

Here, we are asking Kong Ingress Controller to execute the response-transformer
plugin whenever a request matching any of the above two Ingress rules is
processed.

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
the request matches the Ingress rules defined in our two Ingress rules.

## Updating plugin configuration

Now, let's update the plugin configuration to change the header value from
`injected-by-kong` to `injected-by-kong-for-kubernetes`:

```bash
$ echo '
apiVersion: configuration.konghq.com/v1
kind: KongClusterPlugin
metadata:
  name: add-response-header
  annotations:
    kubernetes.io/ingress.class: kong
config:
  add:
    headers:
    - "demo: injected-by-kong-for-kubernetes"
plugin: response-transformer
' | kubectl apply -f -
kongclusterplugin.configuration.konghq.com/add-response-header configured
```

If you repeat the requests from the last step, you will see Kong
now responds with updated header value.

This guides demonstrates how plugin configuration can be shared across
services running in different namespaces.
This can prove to be useful if the persona controlling the plugin
configuration is different from service owners that are responsible for the
Service and Ingress resources in Kubernetes.

