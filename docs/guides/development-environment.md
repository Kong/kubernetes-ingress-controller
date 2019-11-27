#### Development Environment
If you want to develop the ingress controller locally, take the following steps.

1. This guide assumes you are running in GKE.
2. From this directory, apply the dbless config below to get kong running in your k8s cluster:

    `kubectl apply -f dev-config.yaml`

3. Run Kong Ingress Controller locally:

    `bash start.sh kong`

#### Optional: GRPC

To test running grpc:

1. Add a grpc deployment and service

`kubectl apply -f deploy/manifests/grpc.yaml`

2. Create a demo grpc ingress rule:

`kubectl apply -f  sample-grpc-ingress.yaml`

3. Update your ingress with `kubectl patch ingress demo -p '{"metadata":{"annotations":{"configuration.konghq.com/protocols":"grpc,grpcs"}}}'`

4. Update your grpc service with `kubectl patch svc grpc -p '{"metadata":{"annotations":{"configuration.konghq.com/protocol":"grpcs"}}}'`

5. You should be able to run a request over grpcs via `grpcurl -v -d '{"greeting": "Kong 1.3!"}' -H 'kong-debug: 1' -insecure $PROXY_IP:443 hello.HelloService.SayHello`.
