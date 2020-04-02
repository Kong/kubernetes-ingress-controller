# Using Kong with Knative

Kong Ingress Controller supports managing ingress traffic for
serverless workloads managed via Knative.

In this guide, we will learn how to use Kong with Knative services and
configure plugins for Knative services.


## Pre-requisite

This guide will be easier to follow if you have access to a Kubernetes
cluster that is running in the cloud rather than Minikube or any other
local environment. The guide requires access to DNS and a public IP
address or DNS name will certainly keep things simpler and easy for you.

## Install Knative

If you don't have knative installed, you need to install Knative:

```
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.13.0/serving-crds.yaml
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.13.0/serving-core.yaml
```

This will install the resources that are required to run Knative.

## Install Kong

Next, we will install Kong Ingress Controller:

```
kubectl apply -f https://bit.ly/k4k8s
```

You can choose to install a different flavor, like using a database,
or using an Enterprise installation instead of Open-Source. You can also
use Helm installation method if that works for you.

Once Kong is installed,
you should note down the IP address or public CNAME of
`kong-proxy` service.

In the current case case,

```shell
$ kubectl get service kong-proxy -n kong
NAME         TYPE           CLUSTER-IP      EXTERNAL-IP    PORT(S)                      AGE
kong-proxy   LoadBalancer   10.63.248.154   35.247.39.83   80:30345/TCP,443:31872/TCP   53m
```

Take a note of the above IP address "`35.247.39.83`". This will be different
for your installation.

## Configure Knative to use Kong for Ingress

### Ingress class

Next, we will configure Knative to use `kong` as the Ingress class:

```
$ kubectl patch configmap/config-network \
  --namespace knative-serving \
    --type merge \
      --patch '{"data":{"ingress.class":"kong"}}'
```

## Setup Knative domain

As the final step, we need to configure Knative's base domain at which
our services will be accessible.

We override the default ConfigMap with the DNS name of `${KONG_IP}.xip.io`.
This will be different for you:

```
$ echo '
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-domain
  namespace: knative-serving
  labels:
    serving.knative.dev/release: v0.13.0
data:
  35.247.39.83.xip.io: ""
' | kubectl apply -f -
configmap/config-domain configured
```

Once this is done, the setup is complete and we can move onto using Knative
and Kong.

## Test connectivity to Kong

Send a request to the above domain that we have configured:

```bash
curl -i http://35.247.39.83.xip.io/
HTTP/1.1 404 Not Found
Date: Wed, 11 Mar 2020 00:18:49 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 48
X-Kong-Response-Latency: 1
Server: kong/1.4.3

{"message":"no Route matched with those values"}
```

The 404 response is expected since we have not configured any services
in Knative yet.

## Install a Knative Service

Let's install our first Knative service:

```
$ echo "
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld-go
  namespace: default
spec:
  template:
    spec:
      containers:
        - image: gcr.io/knative-samples/helloworld-go
          env:
            - name: TARGET
              value: Go Sample v1
" | kubectl apply -f -
```

It can take a couple of minutes for everything to get configured but
eventually, you will see the URL of the Service.
Let's make the call to the URL:

```shell
$ curl -v http://helloworld-go.default.<your-ip>.xip.io
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Content-Length: 20
Connection: keep-alive
Date: Tue, 10 Mar 2020 23:45:14 GMT
X-Kong-Upstream-Latency: 2723
X-Kong-Proxy-Latency: 0
Via: kong/1.4.3

Hello Go Sample v1!
```

The request is served by Knative and from the response HTTP headeres,
we can tell that the request was proxied by Kong.

The first request will also take longer to complete as Knative will spin
up a new Pod to service the request.
We can see how Kong observed this latency and recorded it in the
`X-Kong-Upstream-Latency` header.
If you perform subsequent requests,
they should complete much faster.

## Plugins for knative services

Let's now execute a plugin for our new Knative service.

First, let's create a KongPlugin resource:

```shell
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: add-response-header
config:
  add:
    headers:
    - "demo: injected-by-kong
plugin: response-transformer
" | kubectl apply -f -
kongplugin.configuration.konghq.com/add-response-header created
```

Next, we will update the Knative service created before and add in
annotation in the template:

```shell
$ echo "
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld-go
  namespace: default
spec:
  template:
    metadata:
      annotations:
        konghq.com/plugins: add-response-header
    spec:
      containers:
        - image: gcr.io/knative-samples/helloworld-go
          env:
            - name: TARGET
              value: Go Sample v1
" | kubectl apply -f -
service.serving.knative.dev/helloworld-go configured
```

Please note that the annotation `konghq.com/plugins` is
not added to the Service definition
itself but to the `spec.template.metadata.annotations`.

Let's make the request again:

```shell
$ curl -i http://helloworld-go.default.35.247.39.83.xip.io/
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Content-Length: 20
Connection: keep-alive
Date: Wed, 11 Mar 2020 00:35:07 GMT
demo:  injected-by-kong
X-Kong-Upstream-Latency: 2455
X-Kong-Proxy-Latency: 1
Via: kong/1.4.3

Hello Go Sample v1!
```

As we can see, the response has the `demo` header injected.

This guide demonstrates the power of using Kong and Knative together.
Checkout other plugins and try them out with multiple Knative services.
The possibilities are endless!

