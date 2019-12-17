# Provisioning Consumers and Credentials

This guide walks through how to use the KongConsumer custom
resource and use Secret resources to associate credentials with those
consumers.

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
          serviceName: httpbin
          servicePort: 80
" | kubectl apply -f -
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

## Add authentication to the service

With Kong, adding authentication in front of an API is as simple as
enabling a plugin.

Let's add a KongPlugin resource to protect the API:

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: httpbin-auth
plugin: key-auth
" | kubectl apply -f -
kongplugin.configuration.konghq.com/httpbin-auth created
```

Now, associate this plugin with the previous Ingress rule we created
using the `plugins.konghq.com` annotation:

```bash
$ echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
  annotations:
    plugins.konghq.com: httpbin-auth
spec:
  rules:
  - http:
      paths:
      - path: /foo
        backend:
          serviceName: httpbin
          servicePort: 80
" | kubectl apply -f -
```

Any request matching the proxying rules defined in the `demo` ingress will
now require a valid API key:

```bash
$ curl -i $PROXY_IP/foo/status/200
HTTP/1.1 401 Unauthorized
Date: Wed, 17 Jul 2019 19:30:33 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
WWW-Authenticate: Key realm="kong"
Content-Length: 41
Server: kong/1.2.1

{"message":"No API key found in request"}
```

As you can see above, Kong returns back a `401 Unauthorized` because
we didn't provide an API key.

## Provision a Consumer

Let's create a KongConsumer resource:

```bash
$ echo "apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: harry
username: harry" | kubectl apply -f -
kongconsumer.configuration.konghq.com/harry created
```

Now, let's provision an API-key associated with
this consumer so that we can pass the authentication imposed by Kong:

Next, we will create a [Secret](https://kubernetes.io/docs/concepts/configuration/secret/)
resource with an API-key inside it:

```bash
kubectl create secret generic harry-apikey  \
  --from-literal=kongCredType=key-auth  \
  --from-literal=key=my-sooper-secret-key
```

The type of credential is specified via `kongCredType`.
You can create the Secret using any other method as well.

Since we are using the Secret resource,
Kubernetes will encrypt and store this API-key for us.

Next, we will associate this API-key with the consumer we created previously.

Please note that we are not re-creating the KongConsumer resource but
only updating it to add the `credentials` array:

```bash
$ echo "apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: harry
username: harry
credentials:
- harry-apikey" | kubectl apply -f -
kongconsumer.configuration.konghq.com/harry configured
```

## Use the credential

Now, use the credential to pass authentication:

```bash
$ curl -i -H 'apikey: my-sooper-secret-key' $PROXY_IP/foo/status/200
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Wed, 17 Jul 2019 19:34:44 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 3
X-Kong-Proxy-Latency: 1
Via: kong/1.2.1
```

In this guide, we learned how to leverage an authentication plugin in Kong
and provision credentials. This enables you to offload authentication into
your Ingress layer and keeps the application logic simple.

All other authentication plugins bundled with Kong work in this
way and can be used to quickly add an authentication layer on top of
your microservices.
