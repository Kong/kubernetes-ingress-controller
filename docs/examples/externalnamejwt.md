# Expose an external application and validate JWT

This example shows how we can expose a service located outside the Kubernetes cluster using an Ingress rule similar to the Kong [Getting Started guide](0).
We then extend the example with a plugin for JWT validation.

Requirements:

- working Kubernetes cluster
- Kong Ingress controller installed. Please check the [deploy guide](1)

1. Create a Kubernetes service. First we need to create a Kubernetes Service [type=ExternalName](2) using the hostname of the application we want to expose:

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

2. Configure a [request-transformer](3) plugin to remove the Host header from the original request. This removes the Host header so when the traffic reaches `mockbin.org` does not contains `foo.bar`:

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

3. Create an Ingress to expose the service in the host `foo.bar`:

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

curl ${PROXY_IP}:${HTTP_PORT} -H "Host:foo.bar"
```
5. Inspect the configuration on kong using curl:

```bash
curl ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/routes/ |jq .
curl ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/services/ |jq .
curl ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/upstreams/ |jq .
curl ${KONG_ADMIN_IP}:${KONG_ADMIN_PORT}/upstreams/default.proxy-to-mockbin.80/targets |jq .
```

## Adding some JWT-security

6. Create a consumer.
```bash
echo "
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: anonymous
username: anonymous
" | kubectl create -f -
```

7. Add the JWT plugin to consumer with credentials for your public key, update the command with your `iss` on the key attribute and add the public key used by your IDP to generate your JWTs
```bash
echo "
apiVersion: configuration.konghq.com/v1
kind: KongCredential
metadata:
  name: credential-jwt
type: jwt
consumerRef: anonymous
config:
  key: https://mydomain/auth/realms/myrealm
  algorithm: RS256
  rsa_public_key: |
      -----BEGIN PUBLIC KEY-----
      MIIBI................DAQAB
      -----END PUBLIC KEY-----
" | kubectl create -f -
```
8. Add JWT-plugin configuration:
```bash
echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: jwt
" | kubectl create -f -
```
9. Update the ingress with a reference to the JWT plugin:
```bash
# First remove the old service (alternatively use kubectl patch)
kubectl delete ingress proxy-from-k8s-to-mockbin

echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: proxy-from-k8s-to-mockbin
  annotations:
    jwt.plugin.konghq.com: jwt
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

10. Service call should return Unauthorized if no Bearer header provided, but if a valid one is provided, i.e. a JWT signed with the private key corresponding to the public key in your consumer credentials from step 7, some data is expected in the response.
```bash
# 401 expected:
curl ${PROXY_IP}:${HTTP_PORT} \
-H "Host:foo.bar" -v
# 200 and data expected in the response:
curl ${PROXY_IP}:${HTTP_PORT} \
-H "Host:foo.bar" \
-H "Authorization: Bearer eyJhbGc......nmXA"
```
