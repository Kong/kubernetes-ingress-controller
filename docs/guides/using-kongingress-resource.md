# Using KongIngress resource

In this guide, we will learn how to use KongIngress resource to control
proxy behavior.

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

This is expected as Kong does not yet know how to proxy the request.

## Install a dummy service

We will start by installing the echo service.

```bash
$ kubectl apply -f https://bit.ly/echo-service
service/echo created
deployment.apps/echo created
```

## Setup Ingress

Let's expose the echo service outside the Kubernetes cluster
by defining an Ingress.

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

Let's test:

```bash
$ curl -i $PROXY_IP/foo
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
	node name:	gke-harry-k8s-dev-default-pool-bb23a167-8pgh
	pod name:	echo-d778ffcd8-n9bss
	pod namespace:	default
	pod IP:	10.60.0.4

Server values:
	server_version=nginx: 1.12.2 - lua: 10010

Request Information:
	client_address=10.60.1.10
	method=GET
	real path=/foo
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://35.233.170.67:8080/foo
<-- clipped -- >
```

## Use KongIngress with Ingress resource

Kong will strip the path defined in the Ingress rule before proxying
the request to the service.
This can be seen in the real path value in the above response.

We can configure Kong to not strip out this path and to only respond to GET requests
for this particular Ingress rule.

To modify these behaviours, let's first create a KongIngress resource
defining the new behaviour:

```bash
$ echo "apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
  name: sample-customization
route:
  methods:
  - GET
  strip_path: true" | kubectl apply -f -
kongingress.configuration.konghq.com/test created
```

Now, let's associate this KongIngress resource with our Ingress resource
using the `konghq.com/override` annotation.

```bash
$ kubectl patch ingress demo -p '{"metadata":{"annotations":{"konghq.com/override":"sample-customization"}}}'
ingress.extensions/demo patched
```

Now, Kong will proxy only GET requests on `/foo/bar` path and
strip away `/foo`:

```bash
$ curl -s $PROXY_IP/foo -X POST
{"message":"no Route matched with those values"}


$ curl -s $PROXY_IP/foo/baz


Hostname: echo-d778ffcd8-vrrtw

Pod Information:
	node name:	gke-harry-k8s-dev-default-pool-bb23a167-8pgh
	pod name:	echo-d778ffcd8-vrrtw
	pod namespace:	default
	pod IP:	10.60.0.9

Server values:
	server_version=nginx: 1.12.2 - lua: 10010

Request Information:
	client_address=10.60.1.10
	method=GET
	real path=/baz
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://35.233.170.67:8080/baz
```

As you can see, the real path value is `/baz`.

## Use KongIngress with Service resource

KongIngress can be used to change load-balancing, health-checking and other
proxy behaviours in Kong.

Next, we are going to tweak two settings:

- Configure Kong to hash the requests based on IP address of the client.
- Configure Kong to proxy all the request on `/foo` to `/bar`.

Let's create a KongIngress resource with these settings:

```bash
$ echo 'apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
  name: demo-customization
upstream:
  hash_on: ip
proxy:
  path: /bar/' | kubectl apply -f -
kongingress.configuration.konghq.com/demo-customization created
```

Now, let's associate this KongIngress resource to the echo service.

```bash
$ kubectl patch service echo -p '{"metadata":{"annotations":{"configuration.konghq.com":"demo-customization"}}}'
service/echo patched
```

Let's test this now:

```bash
$ curl $PROXY_IP/foo/baz
Hostname: echo-d778ffcd8-vrrtw

Pod Information:
	node name:	gke-harry-k8s-dev-default-pool-bb23a167-8pgh
	pod name:	echo-d778ffcd8-vrrtw
	pod namespace:	default
	pod IP:	10.60.0.9

Server values:
	server_version=nginx: 1.12.2 - lua: 10010

Request Information:
	client_address=10.60.1.10
	method=GET
	real path=/bar/baz
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://35.233.170.67:8080/bar/baz

<-- clipped -->
```

Real path received by the upstream service (echo) is now changed to `/bar/baz`.

Also, now all the requests will be sent to the same upstream pod:

```bash
$ curl -s $PROXY_IP/foo | grep "pod IP"
	pod IP:	10.60.0.9
$ curl -s $PROXY_IP/foo | grep "pod IP"
	pod IP:	10.60.0.9
$ curl -s $PROXY_IP/foo | grep "pod IP"
	pod IP:	10.60.0.9
$ curl -s $PROXY_IP/foo | grep "pod IP"
	pod IP:	10.60.0.9
$ curl -s $PROXY_IP/foo | grep "pod IP"
	pod IP:	10.60.0.9
$ curl -s $PROXY_IP/foo | grep "pod IP"
	pod IP:	10.60.0.9
```


You can experiement with various load balancing and healthchecking settings
that KongIngress resource exposes to suit your specific use case.
