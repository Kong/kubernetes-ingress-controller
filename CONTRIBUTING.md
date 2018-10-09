# Contributing Guidelines

Thank you for showing interest in contributing to
Kong Ingress Controller.

Following guide will help you navigate
the repository and get your PRs
merged in faster.

## Finding work

If you're new to the project and want to help, but don't know where to start,
look for "Help wanted" or "Good first issue" labels in our
[issue tracker](https://github.com/Kong/kubernetes-ingress-controller/issues).
Alternatively, read our documentation and fix or
improve any issues that you see. We really value documentation contributions
since that makes life easier for a lot of people.

All of the following tasks are noble and worthy contributions that you can
make without coding:

- Reporting a bug
- Helping other members of the community on the
  [support channels](https://github.com/Kong/kubernetes-ingress-controller#seeking-help)
- Fixing a typo in the code
- Fixing a typo in the documentation
- Providing your feedback on the proposed features and designs
- Reviewing Pull Requests

If you wish to contribute code (features or bug fixes), see the [Submitting a
patch](#submitting-a-patch) section.

## Development environement

## Environment

- Golang version >= 1.10 installed
- Access to a k8s cluster, you can use Minikube or GKE
- Install dep for dependency management
- make
- Docker (for building)

## Dependencies

The build uses dependencies in the `vendor` directory, which
must be installed before building a binary/image. Occasionally, you
might need to update the dependencies.

Check the version of `dep` you are using and make sure it is up to date.
If you have an older version of `dep`, you can update it as follows:

```console
$ go get -u github.com/golang/dep
```

This will automatically save the dependencies to the `vendor/` directory.

```console
$ cd $GOPATH/src/github.com/kong/ingress-controller
$ dep ensure -v -vendor-only
```

## Running in dev mode

You can run the ingress controller without building a Docker
Image and installing it onto your docker container.

Following is a helpful shell script that you could
use to run the Ingress Controller without building
the Ingress Controller:

```shell
#!/bin/bash
pkill -f kubectl
# setup proxies
kubectl port-forward svc/kong-proxy -n kong 8443:443 2>&1 > /dev/null &
kubectl port-forward svc/kong-proxy -n kong 8000:80 2>&1 > /dev/null &
kubectl port-forward svc/kong-ingress-controller -n kong 8001:8001 2>&1 > /dev/null &
kubectl proxy --port=8002 2>&1 > /dev/null &

export POD_NAME=`kubectl get po -n kong -o json | jq ".items[] | .metadata.name" -r | grep ingress`
export POD_NAMESPACE=kong

go run -tags gcp ./cli/ingress-controller/ \
--default-backend-service kong/kong-proxy \
--kubeconfig ~/.kube/config \
--publish-service=kong/kong-proxy \
--apiserver-host=http://localhost:8002 \
--kong-url http://localhost:8001
```

## Building

Build is performed via Makefile. Depending on your
requirements you can build a raw server binary, a local container image,
or push an image to a remote repository.

### Build a raw server binary

```console
$ make build
```

### Build a local container image

```console
$ TAG=DEV make docker-build
```

Note: this will use the Docker daemon
running on your system.
If you're developing using minikube, you
should exectue the following to use the
Docker daemon running inside the Minikube VM:

```console
eval $(minikube docker-env)
```

This will allow you to publish images to
Minikube VM, allowing you to reference them
in your Deployment specs.

### Push the container image to a remote repository

```console
$ TAG=DEV REGISTRY=$USER/kong-ingress-controller make docker-push
```

## Deploying

There are several ways to deploy Kong Ingress Controller onto a cluster.
Please check the [deployment guide](/deploy/README.md)

## Testing

To run unit-tests, just run

```console
$ cd $GOPATH/src/github.com/kong/kubernetes-ingress-controller
$ make test
```

## Releasing

Makefile will produce a release binary, as shown above. To publish this
to a wider Kubernetes user base, push the image to a container registry.
Our images are hosted on
[Bintray](https://bintray.com/kong/kubernetes-ingress-controller).

An example release might look like:

```shell
$ export TAG=42
$ make release
```

Please follow these guidelines to cut a release:

- Update the [release](https://help.github.com/articles/creating-releases/)
  page with a link to changelog.
- Cut a release branch, if appropriate.
  All major feature work is done in HEAD. Specific bug fixes are
  cherry-picked into a release branch.
- If you're not confident about the stability of the code,
  [tag](https://help.github.com/articles/working-with-tags/) it as alpha or beta.
  Typically, a release branch should have stable code.
