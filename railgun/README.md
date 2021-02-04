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
KONG_EXTERNAL_CONTROLLER=true KONG_CONFIGURATION_NAMESPACE=kube-system go run main.go
```

There's a `make run` to do this as well, but it's presently broken until we fix `go-kong` upstream issues with `make manifest` (see the [TODO List](/TODO)).

Look in the `examples/` directory for `Ingress` resources to deploy for testing.

[kic]:https://github.com/kong/kubernetes-ingress-controller
[kb]:https://github.com/kubernetes-sigs/kubebuilder
[ctrl]:https://github.com/kubernetes-sigs/controller-runtime/releases/tag/v0.7.0
[q3]:https://github.com/ioquake/ioq3
[slack]:https://app.slack.com/client/T0DS5NB27/C011RQPHDC7
