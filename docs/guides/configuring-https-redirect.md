# Configuring https redirect

This guide walks through how to configure Kong Ingress Controller to
redirect HTTP request to HTTPS so that all communication
from the external world to your APIs and microservices is encrypted.

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
$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
  annotations:
    konghq.com/strip-path: "true"
spec:
  rules:
  - http:
      paths:
      - path: /foo
        backend:
          serviceName: httpbin
          servicePort: 80
' | kubectl apply -f -
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

## Setup HTTPS redirect

Next, we will create a KongIngress resource which will enforce a policy
on Kong to accept only HTTPS requests for the above Ingress rule and
send back a redirect if the request matches the Ingress rule:

```bash
$ echo "apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
    name: https-only
route:
  protocols:
  - http
  https_redirect_status_code: 302
        " | kubectl apply -f -
kongingress.configuration.konghq.com/https-only created
```

Next, we need to associate the KongIngress resource with the Ingress resource
we created before:

```bash
$ kubectl patch ingress demo -p '{"metadata":{"annotations":{"konghq.com/override":"https-only"}}}'
ingress.extensions/demo patched
```

## Test it

Now, making a plain-text HTTP request to Kong will result in a redirect
being issued from Kong:

```bash
$ curl $PROXY_IP/foo/headers -I
HTTP/1.1 302 Moved Temporarily
Date: Tue, 06 Aug 2019 18:04:38 GMT
Content-Type: text/html
Content-Length: 167
Connection: keep-alive
Location: https://35.197.125.63/foo/headers
Server: kong/1.2.1
```

The `Location` header will contain the URL you need to use for an HTTPS
request. Please note that this URL will be different depending on your
installation method. You can also grab the IP address of the load balancer
fronting Kong and send a HTTPS request to test it.

Let's test it:

```bash
$ curl -k https://35.197.125.63/foo/headers
{
  "headers": {
    "Accept": "*/*",
    "Connection": "keep-alive",
    "Host": "35.197.125.63",
    "User-Agent": "curl/7.54.0",
    "X-Forwarded-Host": "35.197.125.63"
  }
}
```

We can see that Kong correctly serves the request only on HTTPS protocol
and redirects the user if plaint-text HTTP protocol is used.
We had to use `-k` flag in cURL to skip certificate validation as the
certificate served by Kong is a self-signed one.
If you are serving this traffic via a domain that you control and have
configured TLS properties for it, then the flag won't
be necessary.

If you have a domain that you control but don't have TLS/SSL certificates
for it, please check out out
[Using cert-manager with Kong](cert-manager.md) guide which can get TLS
certificates setup for you automatically. And it's free, thanks to
Let's Encrypt!
