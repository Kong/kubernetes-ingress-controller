# Using OIDC plugin

Kong Enterprise's OIDC plugin can authenticate requests using OpenID Connect protocol.
This guide shows a basic example of how to setup the OIDC plugin using
the Ingress Controller.

> Note: This works only with Enterprise version of Kong.

## Installation

Please follow the [deployment](../deployment/k4k8s-enterprise.md) documentation
to install enterprise version of Kong Ingress Controller.

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
$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
spec:
  rules:
  - host: 192.0.2.8.xip.io
    http:
      paths:
      - path: /
        backend:
          serviceName: httpbin
          servicePort: 80
' | kubectl apply -f -
ingress.extensions/demo created
```

We are using `192.0.2.8.xip.io` as our host, you can use any domain name
of your choice. A domain name is a prerequisite for this guide.
For demo purpose, we are using [xip.io](http://xip.io)
service to avoid setting up a DNS record.

Test the Ingress rule:

```bash
$ curl -i 192.0.2.8.xip.io/status/200
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

Next, open a browser and browse to `http://192.0.2.8.xip.io`.
You should see landing page same as httpbin.org.

## Setup OIDC plugin

Now we are going to protect our dummy service with OpenID Connect
protocol using Google as our identity provider.

First, setup an Oauth 2.0 application in
[Google](https://developers.google.com/identity/protocols/oauth2/openid-connect).

Once you have setup your application in Google, use the client ID and client
secret and create a KongPlugin resource in Kubernetes:

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: oidc-auth
config:
  issuer:
  - https://accounts.google.com/.well-known/openid-configuration
  client_id:
  - <client-id>
  client_secret:
  - <client-secret>
  redirect_uri:
  - http://192.0.2.8.xip.io
plugin: openid-connect
" | kubectl apply -f -
kongplugin.configuration.konghq.com/global-rate-limit created
```

The `redirect_uri` parameter must be a URI that matches the Ingress rule we
created earlier. You must also [add it to your Google OIDC
configuration](https://developers.google.com/identity/protocols/oauth2/openid-connect#setredirecturi)

Next, enable the plugin on our Ingress:

```bash
$ kubectl patch ing demo -p '{"metadata":{"annotations":{"konghq.com/plugin":"oidc-auth"}}}'
ingress.extensions/demo patched
```
## Test

Now, if you visit the host you have set up in your Ingress resource,
Kong should redirect you to Google to verify your identity.
Once you identify yourself, you should be able to browse our dummy service
once again.

This basic configuration permits any user with a valid Google account to access
the dummy service.
For setting up more complicated authentication and authorization flows,
please read
[plugin docs](https://docs.konghq.com/enterprise/1.5.x/plugins/oidc-google/).
