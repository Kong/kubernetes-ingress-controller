## Installation

Please follow the [deployment](../deployment) documentation to install
Kong Ingress Controller onto your Kubernetes cluster.

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

`kubectl apply -f https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/sample-apps/grpc.yaml`

2. Create a demo grpc ingress rule:

`kubectl apply -f  https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/sample-apps/sample-grpc-ingress.yaml`

3. Update your ingress with `kubectl patch ingress demo -p '{"metadata":{"annotations":{"configuration.konghq.com/protocols":"grpc,grpcs"}}}'`

4. Update your grpc service with `kubectl patch svc grpc -p '{"metadata":{"annotations":{"configuration.konghq.com/protocol":"grpcs"}}}'`

5. You should be able to run a request over grpcs via `grpcurl -v -d '{"greeting": "Kong 1.3!"}' -H 'kong-debug: 1' -insecure $PROXY_IP:443 hello.HelloService.SayHello`.
