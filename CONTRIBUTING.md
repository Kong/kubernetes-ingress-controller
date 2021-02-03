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

If you wish to contribute code (features or bug fixes) please keep in mind the following:

- bug fix pull requests should be opened against `main` as the base branch
- feature pull requests should be opened with `next` as the base branch

## Stale issue and pull request policy

To ensure our backlog is organized and up to date, we will close issues and
pull requests that have been inactive awaiting a community response for over 2
weeks. If you wish to reopen a closed issue or PR to continue work, please
leave a comment asking a team member to do so.

## Development environment

## Environment

- Golang version matching our [`Dockerfile`](./Dockerfile) installed
- Access to a k8s cluster, you can use Minikube or GKE
- make
- Docker (for building)

## Dependencies

The build uses dependencies are managed by [go modules](https://blog.golang.org/using-go-modules)

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
--kong-admin-url http://localhost:8001
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
$ TAG=DEV REGISTRY=docker.example.com/registry make container
```

Note: this will use the Docker daemon
running on your system.
If you're developing using minikube, you
should execute the following to use the
Docker daemon running inside the Minikube VM:

```console
eval $(minikube docker-env)
```

This will allow you to publish images to
Minikube VM, allowing you to reference them
in your Deployment specs.

### Push the container image to a remote repository

```console
$ docker push docker.example.com/registry/kong-ingress-controller:DEV
```

Note: replace `docker.example.com/registry` with your registry URL.

## Deploying

There are several ways to deploy Kong Ingress Controller onto a cluster.
Please check the [deployment guide](https://docs.konghq.com/kubernetes-ingress-controller/latest/deployment/overview/).

## Testing

To run unit-tests, just run

```console
$ cd $GOPATH/src/github.com/kong/kubernetes-ingress-controller
$ make test
```

To run integration tests, see the [integration test readme](test/integration/README.md).

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
