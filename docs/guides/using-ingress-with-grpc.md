## Installation

Please follow the [deployment](../deployment) documentation to install
Kong Ingress Controller onto your Kubernetes cluster.

## Pre-requisite

To make `gRPC` requests, you need a client which can invoke gRPC requests.
In this guide, we use
[`grpcurl`](https://github.com/fullstorydev/grpcurl#installation).
Please ensure that you have that installed in on your local system.

## Testing connectivity to Kong

This guide assumes that `PROXY_IP` environment variable is
set to contain the IP address or URL pointing to Kong.
If you've not done so, please follow one of the
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

#### Running GRPC

1. Add a grpc deployment and service

```bash
$ kubectl apply -f https://bit.ly/grpcbin-service
service/grpcbin created
deployment.apps/grpcbin created
```
1. Create a demo grpc ingress rule:

```bash
$ echo "apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
spec:
  rules:
  - http:
      paths:
      - path: /
        backend:
          serviceName: grpcbin
          servicePort: 9001" | kubectl apply -f -
ingress.extensions/demo created
```
1. Next, we need to update the Ingress rule to specify gRPC as the protocol.
By default, all routes are assumed to be either HTTP or HTTPS. This annotation
informs Kong that this route is a gRPC(s) route and not a plain HTTP route:

```
$ kubectl patch ingress demo -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc,grpcs"}}}'
```

1. Next, we also update the upstream protocol to be `grpcs`.
Similar to routes, Kong assumes that services are HTTP-based by default.
With this annotation, we configure Kong to use gRPCs protocol when it
talks to the upstream service:

```
$ kubectl patch svc grpcbin -p '{"metadata":{"annotations":{"konghq.com/protocol":"grpcs"}}}'
```

1. You should be able to run a request over `gRPC`:

```
$ grpcurl -v -d '{"greeting": "Kong Hello world!"}' -insecure $PROXY_IP:443 hello.HelloService.SayHello
```
