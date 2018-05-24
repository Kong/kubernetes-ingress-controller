# Expose an external application

This example shows how we can expose a service located outside the Kubernetes cluster using an Ingress rule similar to the Kong [Getting Started guide][0]

Requirements:

- working Kubernetes cluster
- Kong Ingress controller installed. Please check the [deploy guide][1]

1. Create a Kubernetes service

First we need to create a Kubernetes Service [type=ExternalName][2] using the hostname of the application we want to expose

```bash
echo "
kind: Service
apiVersion: v1
metadata:
  name: proxy-to-mockbin
spec:
  ports:
  - protocol: TCP
    port: 80
  type: ExternalName
  externalName: mockbin.org
" | kubectl create -f -
```

2. Configure a [request-transformer][3] plugin to remove the Host header from the original request.

This removes the Host header so when the traffic reaches `mockbin.org` does not contains `foo.bar`

```bash
echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: transform-request-to-mockbin
config:
  remove:
    headers: host
" | kubectl create -f -
```

3. Create an Ingress to expose the service in the host `foo.bar`

```bash
echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: proxy-from-k8s-to-mockbin
  annotations:
    request-transformer.plugin.konghq.com: |
      transform-request-to-mockbin
spec:
  rules:
  - host: foo.bar
    http:
      paths:
      - path: /
        backend:
          serviceName: proxy-to-mockbin
          servicePort: 80
" | kubectl create -f -
```

4. Now we can test the service running:

```bash
export KONG_ADMIN_PORT=$(minikube service -n kong kong-ingress-controller --url --format "{{ .Port }}")
export KONG_ADMIN_IP=$(minikube service   -n kong kong-ingress-controller --url --format "{{ .IP }}")
export PROXY_IP=$(minikube   service -n kong kong-proxy --url --format "{{ .IP }}" | head -1)
export HTTP_PORT=$(minikube  service -n kong kong-proxy --url --format "{{ .Port }}" | head -1)
export HTTPS_PORT=$(minikube service -n kong kong-proxy --url --format "{{ .Port }}" | tail -1)

http ${PROXY_IP}:${HTTP_PORT} Host:foo.bar

```

5. What is configured in Kong?

```bash
http ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/routes/
http ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/services/
http ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/upstreams/
http ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/upstreams/default.proxy-to-mockbin.80/targets
```

[0]: https://getkong.org/docs/0.13.x/getting-started/configuring-a-service/
[1]: https://github.com/Kong/kubernetes-ingress-controller/tree/master/deploy
[2]: https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors
[3]: https://getkong.org/plugins/request-transformer/
