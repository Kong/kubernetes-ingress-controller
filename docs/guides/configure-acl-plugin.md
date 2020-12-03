# Configuring ACL Plugin

This guide walks through configuring the Kong ACL Plugin. The ACL Plugin
requires the use of at least one Authentication plugin. This example will use
the JWT Auth Plugin

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
$ kubectl apply -f https://bit.ly/k8s-httpbin
service/httpbin created
deployment.apps/httpbin created
```

Create two Ingress rules to proxy the httpbin service we just created:

```bash
$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-get
  annotations:
    konghq.com/strip-path: "false"
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - http:
      paths:
      - path: /get
        backend:
          serviceName: httpbin
          servicePort: 80
' | kubectl apply -f -

$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-post
  annotations:
    konghq.com/strip-path: "false"
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - http:
      paths:
      - path: /post
        backend:
          serviceName: httpbin
          servicePort: 80
' | kubectl apply -f -
```

Test the Ingress rules:

```bash
$ curl -i $PROXY_IP/get
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

$ curl -i --data "foo=bar" -X POST $PROXY_IP/post
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

{
  "args": {},
  "data": "",
  "files": {},
  "form": {
    "foo": "bar"
  },
  "headers": {
    "Accept": "*/*",
    "Connection": "keep-alive",
    "Content-Length": "7",
    "Content-Type": "application/x-www-form-urlencoded",
    "Host": "localhost",
    "User-Agent": "curl/7.54.0",
    "X-Credential-Identifier": "localhost",
    "X-Forwarded-Host": "localhost"
  },
  "json": null,
  "origin": "192.168.0.3",
  "url": "http://some.url/post"
}

```

## Add JWT authentication to the service

With Kong, adding authentication in front of an API is as simple as
enabling a plugin. Let's enable JWT authentication

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: app-jwt
plugin: jwt
" | kubectl apply -f -
```

Now let's associate the plugin to the Ingress rules we created earlier.

```bash
$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-get
  annotations:
    konghq.com/plugins: app-jwt
    konghq.com/strip-path: "false"
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - http:
      paths:
      - path: /get
        backend:
          serviceName: httpbin
          servicePort: 80
' | kubectl apply -f -

$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-post
  annotations:
    konghq.com/plugins: app-jwt
    konghq.com/strip-path: "false"
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - http:
      paths:
      - path: /post
        backend:
          serviceName: httpbin
          servicePort: 80
' | kubectl apply -f -
```

Any requests matching the proxying rules for `demo-get` and `demo` post will
now require a valid JWT and the consumer for the JWT to be associate with the
right ACL.

```bash
$ curl -i $PROXY_IP/get

HTTP/1.1 401 Unauthorized
Date: Mon, 06 Apr 2020 07:27:44 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 50
X-Kong-Response-Latency: 2
Server: kong/2.0.2


{"message":"Unauthorized"}

$ curl -i --data "foo=bar" -X POST $PROXY_IP/post

HTTP/1.1 401 Unauthorized
Date: Mon, 06 Apr 2020 07:27:44 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 50
X-Kong-Response-Latency: 2
Server: kong/2.0.2


{"message":"Unauthorized"}
```

You should get a 401 response telling you that the request is not authorized.

## Provision Consumers

Let's provision 2 KongConsumer resources:

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: admin
  annotations:
    kubernetes.io/ingress.class: kong
username: admin
" | kubectl apply -f -

$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: plain-user
  annotations:
    kubernetes.io/ingress.class: kong
username: plain-user
" | kubectl apply -f -
```

## Secrets

Next, let's provision some Secrets for the KongConsumers to reference. Each
ACL will need its own Secret and each JWT public key will need its own Secret.
The credential type is specified in the `kongCredType` field. In this
case we'll be using `jwt` and `acl`. You can create a secret using any other
method as well.

The JWT signing algorithm is set in the `algorithm` field. The if using a
public key like this example it is stored in the `rsa_pulic_key` field. If you
are using a secret signing key, use the `secret` field. The last field to set
if you are using `RS256` or `ES256` is the `key` field. This should match the
`iss` field in the JWT you will be sending. You can check this value by
decoding your JWT over at [https://jwt.io](https://jwt.io)

Since we are using the Secret resource, Kubernetes will encrypt and store the
JWT signing key and ACL group for us.

### JWT signing key

```bash
# create secret for jwt public key
$ kubectl create secret \
  generic app-admin-jwt  \
  --from-literal=kongCredType=jwt  \
  --from-literal=key="admin-issuer" \
  --from-literal=algorithm=RS256 \
  --from-literal=rsa_public_key="-----BEGIN PUBLIC KEY-----
  MIIBIjA....
  -----END PUBLIC KEY-----"

# create a second secret with a different key
$ kubectl create secret \
  generic app-user-jwt  \
  --from-literal=kongCredType=jwt  \
  --from-literal=key="user-issuer" \
  --from-literal=algorithm=RS256 \
  --from-literal=rsa_public_key="-----BEGIN PUBLIC KEY-----
  qwerlkjqer....
  -----END PUBLIC KEY-----"
```

## Assign the credentials

In order to for the ACL and JWT to be validated by Kong, the secrets will need
to be referenced by the KongConsumers we created earlier. Let's update those.

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: admin
  annotations:
    kubernetes.io/ingress.class: kong
username: admin
credentials:
  - app-admin-jwt
" | kubectl apply -f -

$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: plain-user
  annotations:
    kubernetes.io/ingress.class: kong
username: plain-user
credentials:
  - app-user-jwt
" | kubectl apply -f -
```

## Use the credential

Now to use a JWT to pass authentication. Let's store the user and admin jwt's
in some environment variables. `USER_JWT` and `ADMIN_JWT`. If you are using
an identity provider, you should be able to login and get out a JWT from their
API. If you are generating your own, go through the process of generating your
own.

Let's test the get route

```bash
$ curl -i -H "Authorization: Bearer ${USER_JWT}" $PROXY_IP/get

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 947
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 06 Apr 2020 06:45:45 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 7
X-Kong-Proxy-Latency: 2
Via: kong/2.0.2

{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Authorization": "Bearer eyJ...",
    "Connection": "keep-alive",
    "Host": "localhost",
    "User-Agent": "curl/7.54.0",
    "X-Consumer-Id": "393611c3-aea9-510d-9be4-ac429ecc53f4",
    "X-Consumer-Username": "plain-user",
    "X-Credential-Identifier": "localhost",
    "X-Forwarded-Host": "localhost"
  },
  "origin": "192.168.0.3",
  "url": "http://some.url/get"
}



$ curl -i -H "Authorization: Bearer ${ADMIN_JWT}" $PROXY_IP/get

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 947
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 06 Apr 2020 06:45:45 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 7
X-Kong-Proxy-Latency: 2
Via: kong/2.0.2

{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Authorization": "Bearer eyJ...",
    "Connection": "keep-alive",
    "Host": "localhost",
    "User-Agent": "curl/7.54.0",
    "X-Consumer-Id": "a6edc906-2f9f-5fb2-a373-efac406f0ef2",
    "X-Consumer-Username": "admin",
    "X-Credential-Identifier": "localhost",
    "X-Forwarded-Host": "localhost"
  },
  "origin": "192.168.0.3",
  "url": "http://some.url/get"
}

```

Now let's test the post route

```bash
$ curl -i -X POST --data "foo=bar" \
-H "Authorization: Bearer ${USER_JWT}" $PROXY_IP/post

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 947
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 06 Apr 2020 06:45:45 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 7
X-Kong-Proxy-Latency: 2
Via: kong/2.0.2

{
  "args": {},
  "data": "",
  "files": {},
  "form": {
    "foo": "bar"
  },
  "headers": {
    "Accept": "*/*",
    "Authorization": "Bearer eyJ...",
    "Connection": "keep-alive",
    "Content-Length": "7",
    "Content-Type": "application/x-www-form-urlencoded",
    "Host": "localhost",
    "User-Agent": "curl/7.54.0",
    "X-Consumer-Id": "393611c3-aea9-510d-9be4-ac429ecc53f4",
    "X-Consumer-Username": "plain-user",
    "X-Credential-Identifier": "localhost",
    "X-Forwarded-Host": "localhost"
  },
  "json": null,
  "origin": "192.168.0.3",
  "url": "http://some.url/post"
}

$ curl -i -X POST --data "foo=bar" \
-H "Authorization: Bearer ${ADMIN_JWT}" $PROXY_IP/post

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 947
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 06 Apr 2020 06:45:45 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 7
X-Kong-Proxy-Latency: 2
Via: kong/2.0.2

{
  "args": {},
  "data": "",
  "files": {},
  "form": {
    "foo": "bar"
  },
  "headers": {
    "Accept": "*/*",
    "Authorization": "Bearer eyJ...",
    "Connection": "keep-alive",
    "Content-Length": "7",
    "Content-Type": "application/x-www-form-urlencoded",
    "Host": "localhost",
    "User-Agent": "curl/7.54.0",
    "X-Consumer-Id": "393611c3-aea9-510d-9be4-ac429ecc53f4",
    "X-Consumer-Username": "admin",
    "X-Credential-Identifier": "localhost",
    "X-Forwarded-Host": "localhost"
  },
  "json": null,
  "origin": "192.168.0.3",
  "url": "http://some.url/post"
}


```

## Adding ACL's

The JWT plugin doesn't provide the ability to authroize a given issuer to a
given ingress. To do this we need to use the ACL plugin. Let's create an admin
ACL config

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: admin-acl
plugin: acl
config:
  whitelist: ['app-admin']
" | kubectl apply -f -
```

Then let's create a user ACL config. We want our admin to be able to access
the same resources as the user, so let's make sure we include them in the
whitelist.

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: plain-user-acl
plugin: acl
config:
  whitelist: ['app-user','app-admin']
" | kubectl apply -f -
```

Next let's create the secrets that will define the ACL groups.

```bash
# create secrets for acl groups
$ kubectl create secret \
  generic app-admin-acl \
  --from-literal=kongCredType=acl  \
  --from-literal=group=app-admin

$ kubectl create secret \
  generic app-user-acl \
  --from-literal=kongCredType=acl  \
  --from-literal=group=app-user
```

After we create the secrets, the consumers need to be updated to reference the
ACL credentials

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: admin
  annotations:
    kubernetes.io/ingress.class: kong
username: admin
credentials:
  - app-admin-jwt
  - app-admin-acl
" | kubectl apply -f -

$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: plain-user
  annotations:
    kubernetes.io/ingress.class: kong
username: plain-user
credentials:
  - app-user-jwt
  - app-user-acl
" | kubectl apply -f -
```

The last thing to configure is the ingress to use the new plguins. Note, if you
set more than one ACL plugin, the last one supplied will be the only one
evaluated.

```bash
$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-get
  annotations:
    konghq.com/plugins: app-jwt,plain-user-acl
    konghq.com/strip-path: "false"
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - http:
      paths:
      - path: /get
        backend:
          serviceName: httpbin
          servicePort: 80
' | kubectl apply -f -

$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-post
  annotations:
    konghq.com/plugins: app-jwt,admin-acl
    konghq.com/strip-path: "false"
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - http:
      paths:
      - path: /post
        backend:
          serviceName: httpbin
          servicePort: 80
' | kubectl apply -f -
```

Now let's test it.

```bash
$ curl -i -H "Authorization: Bearer ${USER_JWT}" $PROXY_IP/get

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 947
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 06 Apr 2020 06:45:45 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 7
X-Kong-Proxy-Latency: 2
Via: kong/2.0.2

{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Authorization": "Bearer eyJ...",
    "Connection": "keep-alive",
    "Host": "localhost",
    "User-Agent": "curl/7.54.0",
    "X-Consumer-Groups": "app-user",
    "X-Consumer-Id": "393611c3-aea9-510d-9be4-ac429ecc53f4",
    "X-Consumer-Username": "plain-user",
    "X-Credential-Identifier": "localhost",
    "X-Forwarded-Host": "localhost"
  },
  "origin": "192.168.0.3",
  "url": "http://some.url/get"
}



$ curl -i -H "Authorization: Bearer ${ADMIN_JWT}" $PROXY_IP/get

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 947
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 06 Apr 2020 06:45:45 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 7
X-Kong-Proxy-Latency: 2
Via: kong/2.0.2

{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Authorization": "Bearer eyJ...",
    "Connection": "keep-alive",
    "Host": "localhost",
    "User-Agent": "curl/7.54.0",
    "X-Consumer-Groups": "app-admin",
    "X-Consumer-Id": "a6edc906-2f9f-5fb2-a373-efac406f0ef2",
    "X-Consumer-Username": "admin",
    "X-Credential-Identifier": "localhost",
    "X-Forwarded-Host": "localhost"
  },
  "origin": "192.168.0.3",
  "url": "http://some.url/get"
}

```

Now let's test the post route

```bash
$ curl -i -X POST --data "foo=bar" \
-H "Authorization: Bearer ${USER_JWT}" $PROXY_IP/post
HTTP/1.1 403 Forbidden
Date: Mon, 06 Apr 2020 07:11:59 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 45
X-Kong-Response-Latency: 1
Server: kong/2.0.2

{"message":"You cannot consume this service"}
```

The `plain-user` user is not in the `admin-acl` whitelist, and is therefore
unauthorized to access the resource

```bash
$ curl -i -X POST --data "foo=bar" \
-H "Authorization: Bearer ${ADMIN_JWT}" $PROXY_IP/post

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 1156
Connection: keep-alive
Server: gunicorn/19.9.0
Date: Mon, 06 Apr 2020 07:20:35 GMT
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
X-Kong-Upstream-Latency: 4
X-Kong-Proxy-Latency: 4
Via: kong/2.0.2

{
  "args": {},
  "data": "",
  "files": {},
  "form": {
    "foo": "bar"
  },
  "headers": {
    "Accept": "*/*",
    "Authorization": "Bearer eyJ...",
    "Connection": "keep-alive",
    "Content-Length": "7",
    "Content-Type": "application/x-www-form-urlencoded",
    "Host": "localhost",
    "User-Agent": "curl/7.54.0",
    "X-Consumer-Groups": "app-admin",
    "X-Consumer-Id": "393611c3-aea9-510d-9be4-ac429ecc53f4",
    "X-Consumer-Username": "admin",
    "X-Credential-Identifier": "localhost",
    "X-Forwarded-Host": "localhost"
  },
  "json": null,
  "origin": "192.168.0.3",
  "url": "http://some.url/post"
}
```
