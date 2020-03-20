# Getting started with Kong Ingress Controller

## Installation

Please follow the [deployment](../deployment) documentation to install
Kong Ingress Controller onto your Kubernetes cluster.

## Testing connectivity to Kong

This guide assumes that `PROXY_IP` environment variable is
set to contain the IP address or URL pointing to Kong.
If you've not done so, please follow one of the
[deployment guides](../deployment) to configure this environment variable.

If everything is setup correctly, making a request to Kong should return back
a HTTP 404 Not Found.

```bash
$ curl -i $PROXY_IP
HTTP/1.1 404 Not Found
Date: Fri, 21 Jun 2019 17:01:07 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 48
Server: kong/1.1.2

{"message":"no Route matched with those values"}
```

This is expected since Kong doesn't know how to proxy the request yet.

## Setup an echo-server

Setup an echo-server application to demonstrate how
to use Kong Ingress Controller:

```bash
$ curl -sL bit.ly/echo-service | kubectl apply -f -
service/echo created
deployment.apps/echo created
```

This application just returns information about the
pod and details from the HTTP request.

## Basic proxy

Create an Ingress rule to proxy the echo-server created previously:

```bash
$ echo "
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
          serviceName: echo
          servicePort: 80
" | kubectl apply -f -
ingress.extensions/demo created
```

Test the Ingress rule:

```bash
$ curl -i $PROXY_IP/foo
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Transfer-Encoding: chunked
Connection: keep-alive
Date: Fri, 21 Jun 2019 17:12:49 GMT
Server: echoserver
X-Kong-Upstream-Latency: 0
X-Kong-Proxy-Latency: 1
Via: kong/1.1.2



Hostname: echo-758859bbfb-txt52

Pod Information:
        node name:      minikube
        pod name:       echo-758859bbfb-txt52
        pod namespace:  default
        pod IP: 172.17.0.14
<-- clipped -->
```

If everything is deployed correctly, you should see the above response.
This verifies that Kong can correctly route traffic to an application running
inside Kubernetes.

## Using plugins in Kong

Setup a KongPlugin resource:

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: request-id
config:
  header_name: my-request-id
plugin: correlation-id
" | kubectl apply -f -
kongplugin.configuration.konghq.com/request-id created
```

Create a new Ingress resource which uses this plugin:

```bash
$ echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-example-com
  annotations:
    plugins.konghq.com: request-id
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /bar
        backend:
          serviceName: echo
          servicePort: 80
" | kubectl apply -f -
ingress.extensions/demo-example-com created
```

The above resource directs Kong to execute the request-id plugin whenever
a request is proxied matching any rule defined in the resource.

Send a request to Kong:

```bash
$ curl -i -H "Host: example.com" $PROXY_IP/bar/sample
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Transfer-Encoding: chunked
Connection: keep-alive
Date: Fri, 21 Jun 2019 18:09:02 GMT
Server: echoserver
X-Kong-Upstream-Latency: 1
X-Kong-Proxy-Latency: 1
Via: kong/1.1.2



Hostname: echo-758859bbfb-cnfmx

Pod Information:
        node name:      minikube
        pod name:       echo-758859bbfb-cnfmx
        pod namespace:  default
        pod IP: 172.17.0.9

Server values:
        server_version=nginx: 1.12.2 - lua: 10010

Request Information:
        client_address=172.17.0.2
        method=GET
        real path=/sample
        query=
        request_version=1.1
        request_scheme=http
        request_uri=http://example.com:8080/sample

Request Headers:
        accept=*/*
        connection=keep-alive
        host=example.com
        my-request-id=7250803a-a85a-48da-94be-1aa342ca276f#6
        user-agent=curl/7.54.0
        x-forwarded-for=172.17.0.1
        x-forwarded-host=example.com
        x-forwarded-port=8000
        x-forwarded-proto=http
        x-real-ip=172.17.0.1

Request Body:
        -no body in request-
```

The `my-request-id` can be seen in the request received by echo-server.
It is injected by Kong as the request matches one
of the Ingress rules defined in `demo-example-com` resource.

## Using plugins on Services

Kong Ingress allows plugins to be executed on a service level, meaning
Kong will execute a plugin whenever a request is sent to a specific k8s service,
no matter which Ingress path it came from.

Create a KongPlugin resource:

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: rl-by-ip
config:
  minute: 5
  limit_by: ip
  policy: local
plugin: rate-limiting
" | kubectl apply -f -
kongplugin.configuration.konghq.com/rl-by-ip created
```

Next, apply the `plugins.konghq.com` annotation on the Kubernetes Service
that needs rate-limiting:

```bash
kubectl patch svc echo \
  -p '{"metadata":{"annotations":{"plugins.konghq.com": "rl-by-ip\n"}}}'
```

Now, any request sent to this service will be protected by a rate-limit
enforced by Kong:

```bash
$ curl -I $PROXY_IP/foo
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Fri, 21 Jun 2019 18:25:49 GMT
Server: echoserver
X-RateLimit-Limit-minute: 5
X-RateLimit-Remaining-minute: 2
X-Kong-Upstream-Latency: 0
X-Kong-Proxy-Latency: 4
Via: kong/1.1.2

$ curl -I -H "Host: example.com" $PROXY_IP/bar/sample
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Fri, 21 Jun 2019 18:28:30 GMT
Server: echoserver
X-RateLimit-Limit-minute: 5
X-RateLimit-Remaining-minute: 4
X-Kong-Upstream-Latency: 1
X-Kong-Proxy-Latency: 2
Via: kong/1.1.2
```

## Result

This guide sets up the following configuration:

```text
HTTP requests with /foo -> Kong enforces rate-limit -> echo server

HTTP requests with /bar -> Kong enforces rate-limit +   -> echo-server
   on example.com          injects my-request-id header
```
