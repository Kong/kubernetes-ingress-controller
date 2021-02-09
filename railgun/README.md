# Kubernetes Ingress Controller (KIC) Prototype

This repo holds a prototype implementation of the [Kong Kubernetes Ingress Controller (KIC)][kic] updated and built for [Kubebuilder][kb] and [Controller Runtime v0.7.x+][ctrl].

The naming `railgun` is a codename used arbitrarily based on a [video game reference][q3].

If you're confused what you're looking at or have any questions, please reach out to @kong/team-k8s on our [Slack Channel][slack]!

## Usage

To run the dev environment currently, you need a Kubernetes cluster (such as [kind](https://github.com/kubernetes-sigs/kind)):

```shell
$ kind create cluster --name kong-test
kubectl cluster-info --context kind-kong-test
```

Deploy some kong proxies which the controllers will test against with:

```shell
$ make deploy.test.proxy
```

This will deploy the helm chart with the KIC disabled, the Admin API enabled over `ClusterIP`, and 3 proxy replicas.

Now you can run the controller:

```shell
KONG_EXTERNAL_CONTROLLER=true KONG_CONFIGURATION_NAMESPACE=kong-system go run main.go
```

There's a `make run` to do this as well, but it's presently broken until we fix `go-kong` upstream issues with `make manifest` (see the [TODO List](/TODO)).

Look in the `examples/` directory for `Ingress` resources to deploy for testing.

### Integration Demo

```shell
kubectl create namespace kong-system
docker run -d --rm --name kong-dbless -e KONG_ADMIN_LISTEN="0.0.0.0:8001" -e KONG_DATABASE=off kong:2.2
env KONG_CONFIGURATION_NAMESPACE=kong-system ./bin/manager --kong-url=http://172.17.0.2:8001
kubectl create secret -n kong-system generic kong-config --from-literal=a=b
```

# Tests

Right now the integration tests in `test/integration` rely on a private repo, make sure you run the following before running tests:

```shell
export GOPRIVATE="github.com/kong/kubernetes-testing-framework,${GOPRIVATE}"
git config --global url."git@github.com:kong/kubernetes-testing-framework".insteadOf https://github.com/kong/kubernetes-testing-framework

[kic]:https://github.com/kong/kubernetes-ingress-controller
[kb]:https://github.com/kubernetes-sigs/kubebuilder
[ctrl]:https://github.com/kubernetes-sigs/controller-runtime/releases/tag/v0.7.0
[q3]:https://github.com/ioquake/ioq3
[slack]:https://app.slack.com/client/T0DS5NB27/C011RQPHDC7
