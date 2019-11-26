# Expose an external application
This example shows how we can expose a service located outside the Kubernetes cluster using an Ingress.

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

## Create a Kubernetes service

First we need to create a Kubernetes Service [type=ExternalName][0] using the hostname of the application we want to expose.

```bash
echo "
kind: Service
apiVersion: v1
metadata:
  name: proxy-to-httpbin
spec:
  ports:
  - protocol: TCP
    port: 80
  type: ExternalName
  externalName: httpbin.org
" | kubectl create -f -
```

## Configure a [request-transformer][1] plugin to remove the Host header from the original request

This removes the Host header so when the traffic reaches `httpbin.org` does not contain `foo.bar`

```bash
echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: transform-request-to-httpbin
config:
  remove:
    headers: host
plugin: request-transformer
" | kubectl create -f -
```

## Create an Ingress to expose the service in the host `foo.bar`

```bash
echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: proxy-from-k8s-to-httpbin
  annotations:
    plugins.konghq.com: transform-request-to-httpbin
spec:
  rules:
  - host: foo.bar
    http:
      paths:
      - path: /
        backend:
          serviceName: proxy-to-httpbin
          servicePort: 80
" | kubectl create -f -
```

## Test the service

```bash
export KONG_ADMIN_PORT=$(minikube service -n kong kong-ingress-controller --url --format "{{ .Port }}")
export KONG_ADMIN_IP=$(minikube service   -n kong kong-ingress-controller --url --format "{{ .IP }}")
export PROXY_IP=$(minikube   service -n kong kong-proxy --url --format "{{ .IP }}" | head -1)
export HTTP_PORT=$(minikube  service -n kong kong-proxy --url --format "{{ .Port }}" | head -1)
export HTTPS_PORT=$(minikube service -n kong kong-proxy --url --format "{{ .Port }}" | tail -1)

http ${PROXY_IP}:${HTTP_PORT} Host:foo.bar
```

## View the Kong configuration

```bash
http ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/routes/
http ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/services/
http ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/upstreams/
http ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/upstreams/default.proxy-to-httpbin.80/targets
```


[0]: https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors
[1]: https://getkong.org/plugins/request-transformer/
