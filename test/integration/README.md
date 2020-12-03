## E2E test

The end-to-end (E2E) integration test performs the following:

1. Start up an ephemeral Kubernetes cluster (KIND)
1. Install a container registry (localhost:5000)
1. Push the Kong Ingress Controller image to localhost:5000/kic:local
1. Run Kong Ingress Controller in the local Kubernetes cluster
1. Start up example services (echo and httpbin)
1. Run test cases (`/test/integration/cases/*`):
    1. Apply manifests from `.yaml` files present in the test case directory
		1. Wait for KIC to apply the configuration
		1. Run ./verify.sh to verify assertions
		1. Delete manifests created in the apply step above.

### How to run the test

```bash
make integration-test
```

This command builds a KIC Docker image and runs the test against that image.

It is possible to run the test against any prebuilt KIC image, skipping the build step:

```bash
env KIC_IMAGE=some-kic-image:tag ./test/integration/test.sh
```

### Troubleshooting

If you want to troubleshoot a specific test case, here's how to do it:

#### Run tests with `SKIP_TEARDOWN=yes`
By passing `SKIP_TEARDOWN=yes` to the test you can inspect the test environment after failure, and run certain test cases manually:

```bash
make SKIP_TEARDOWN=yes integration-test
# or
env KIC_IMAGE=some-kic-image:tag SKIP_TEARDOWN=yes ./test/integration/test.sh
```

#### Access the test cluster with `kubectl`

After the test invocation command with `SKIP_TEARDOWN=yes` set terminates, KIND will continue running in the cluster. You can access it using `./kubeconfig-test-cluster` as kubeconfig:

```bash
kubectl --kubeconfig=./kubeconfig-test-cluster get pods --all-namespaces
```

#### Run a test case manually

You can run a test case manually:
```bash
kubectl --kubeconfig=./kubeconfig-test-cluster port-forward -n kong svc/kong-proxy "27080:80" "27443:443"

# in a separate terminal window:
kubectl --kubeconfig=./kubeconfig-test-cluster apply -f ./test/integration/cases/01-https
env SUT_HTTP_HOST=127.0.0.1:27080 SUT_HTTPS_HOST=127.0.0.1:27443 ./test/integration/cases/01-https/verify.sh
kubectl --kubeconfig=./kubeconfig-test-cluster delete -f ./test/integration/cases/01-https
```

Run service-dual-stack test case use different command manually:
```bash
# get ingress-kong pod name
INGRESS_KONG_POD_ID=`kubectl --kubeconfig=./kubeconfig-test-cluster -n kong get pods -l app=ingress-kong \
  -o jsonpath='{.items[*].metadata.name}' | head -n 1`
  
kubectl --kubeconfig=./kubeconfig-test-cluster port-forward -n kong pod/"$INGRESS_KONG_POD_ID" "28444:8444"
# in a separate terminal window:
kubectl --kubeconfig=./kubeconfig-test-cluster apply -f ./test/integration/cases/07-service-dual-stack
env SUT_ADMIN_API_HOST=127.0.0.1:28444 ./test/integration/cases/07-service-dual-stack/verify.sh
kubectl --kubeconfig=./kubeconfig-test-cluster delete -f ./test/integration/cases/07-service-dual-stack
```

#### Manually tear down the test cluster

At the end of the debugging session, you can tear down the environment like this:
```bash
docker rm -f test-cluster-control-plane test-local-registry
```
