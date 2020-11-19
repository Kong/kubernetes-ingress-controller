# TCPIngress with Kong

This guide walks through using TCPIngress Custom Resource
resource to expose TCP-based services running in Kubernetes to the out
side world.

## Overview

TCP-based Ingress means that Kong simply forwards the TCP stream to a Pod
of a Service that's running inside Kubernetes. Kong will not perform any
sort of transformations.

There are two modes avaialble:
- **Port based routing**: In this mode, Kong simply proxies all traffic it
  receives on a specific port to the Kubernetes Service. TCP connections are
  load balanced across all the available pods of the Service.
- **SNI based routing**: In this mode, Kong accepts a TLS-encrypted stream
  at the specified port and can route traffic to different services based on
  the `SNI` present in the TLS handshake. Kong will also terminate the TLS
  handshake and forward the TCP stream to the Kubernetes Service.

## Installation

Please follow the [deployment](../deployment) documentation to install
Kong Ingress Controller on your Kubernetes cluster.

> **Note**: This feature works with Kong versions 2.0.4 and above.

> **Note**: This feature is available in Controller versions 0.8 and above.

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

## Configure Kong for new ports

First, we will configure Kong's Deployment and Service to expose two new ports
9000 and 9443. Port 9443 expects a TLS connection from the client.

```shell
$ kubectl patch deploy -n kong ingress-kong --patch '{
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "name": "proxy",
            "env": [
              {
                "name": "KONG_STREAM_LISTEN",
                "value": "0.0.0.0:9000, 0.0.0.0:9443 ssl"
              }
            ],
            "ports": [
              {
                "containerPort": 9000,
                "name": "stream9000",
                "protocol": "TCP"
              },
              {
                "containerPort": 9443,
                "name": "stream9443",
                "protocol": "TCP"
              }
            ]
          }
        ]
      }
    }
  }
}'
deployment.extensions/ingress-kong patched
```

```shell
$ kubectl patch service -n kong kong-proxy --patch '{
  "spec": {
    "ports": [
      {
        "name": "stream9000",
        "port": 9000,
        "protocol": "TCP",
        "targetPort": 9000
      },
      {
        "name": "stream9443",
        "port": 9443,
        "protocol": "TCP",
        "targetPort": 9443
      }
    ]
  }
}'
service/kong-proxy patched
```

You are free to choose other ports as well.

## Install TCP echo service

Next, we will install a dummy TCP service.
If you already have a TCP-based service running in your cluster,
you can use that as well.

```shell
$ kubectl apply -f https://bit.ly/tcp-echo
deployment.apps/tcp-echo created
service/tcp-echo created
```

Now, we have a TCP echo service running in Kubernetes.
We will now expose this on plain-text and a TLS based port.

## TCP port based routing

To expose our service to the outside world, create the following
`TCPIngress` resource:

```shell
$ echo "apiVersion: configuration.konghq.com/v1beta1
kind: TCPIngress
metadata:
  name: echo-plaintext
  annotations:
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - port: 9000
    backend:
      serviceName: tcp-echo
      servicePort: 2701
" | kubectl apply -f -
tcpingress.configuration.konghq.com/echo-plaintext created
```

Here we are instructing Kong to forward all traffic it receives on port
9000 to `tcp-echo` service on port 2701.

Once created, we can see the IP address at which this is available:

```shell
$ kubectl get tcpingress
NAME             ADDRESS        AGE
echo-plaintext   <PROXY_IP>   3m18s
```

Lets connect to this service using `telnet`:

```shell
$ telnet $PROXY_IP 9000
Trying 35.247.39.83...
Connected to 35.247.39.83.
Escape character is '^]'.
Welcome, you are connected to node gke-harry-k8s-dev-pool-1-e9ebab5e-c4gw.
Running on Pod tcp-echo-844545646c-gvmkd.
In namespace default.
With IP address 10.60.1.17.
This text will be echoed back.
This text will be echoed back.
^]
telnet> Connection closed.
```

We can see here that the `tcp-echo` service is now available outside the
Kubernetes cluster via Kong.

## TLS SNI based routing

Next, we will demonstrate how Kong can help expose the `tcp-echo` service
in a secure manner to the outside world.

Create the following TCPIngress resource:

```
$ echo "apiVersion: configuration.konghq.com/v1beta1
kind: TCPIngress
metadata:
  name: echo-tls
  annotations:
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - host: example.com
    port: 9443
    backend:
      serviceName: tcp-echo
      servicePort: 2701
" | kubectl apply -f -
tcpingress.configuration.konghq.com/echo-tls created
```

Now, we can access the `tcp-echo` service on port 9443, on SNI `example.com`.

You should setup a DNS record for a Domain that you control
to point to PROXY_IP and then access
the service via that for production usage.

In our contrived demo example, we can connect to the service via TLS
using `openssl`'s `s_client` command:

```shell
$ openssl s_client -connect $PROXY_IP:9443 -servername example.com -quiet
openssl s_client -connect 35.247.39.83:9443 -servername foo.com -quiet
depth=0 C = US, ST = California, L = San Francisco, O = Kong, OU = IT Department, CN = localhost
verify error:num=18:self signed certificate
verify return:1
depth=0 C = US, ST = California, L = San Francisco, O = Kong, OU = IT Department, CN = localhost
verify return:1
Welcome, you are connected to node gke-harry-k8s-dev-pool-1-e9ebab5e-c4gw.
Running on Pod tcp-echo-844545646c-gvmkd.
In namespace default.
With IP address 10.60.1.17.
This text will be echoed back.
This text will be echoed back.
^C
```

Since Kong is not configured with a TLS cert-key pair for `example.com`, Kong
is returning a self-signed default certificate, which is not trusted.
You can also see that the echo service is running as expected.

## Bonus

Scale the `tcp-echo` Deployment to have multiple replicas and observe how
Kong load-balances the TCP-connections between pods.

## Conclusion

In this guide, we see how to use Kong's TCP routing capabilities using
TCPIngress Custom Resource. This can be very useful if you have services
running inside Kubernetes that have custom protocols instead of the more
popular HTTP or gRPC protocols.
