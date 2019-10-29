#### Development Environment
If you want to develop the ingress controller locally, take the following steps.

1. This guide assumes you are running in GKE.
2. From this directory, apply the dbless config to get kong running in your k8s cluster:

    `k apply -f dev-config.yaml`

3. Run Kong Ingress Controller locally:

    `bash start.sh kong`

#### Optional: GRPC

1. To test running grpc, you can create a demo grpc ingress rule:

```
echo "apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
spec:
  rules:
  - http:
      paths:
      - path: /
        backend:
          serviceName: grpc
          servicePort: 9001" | k apply -f -
```

```
echo "apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
    name: grpc-only
proxy:
    protocol: grpcs
route:
  protocols:
    - grpc
    - grpcs" | k apply -f -
```

2. Update your ingress with `kubectl patch ingress demo -p '{"metadata":{"annotations":{"configuration.konghq.com":"grpc-only"}}}'`

3. You should be able to run a request over grpcs (`grpcurl -v -insecure $PROXY_IP:443 hello.HelloService.SayHello`) or grpc (`grpcurl -v -plaintext $PROXY_IP:9080`).
