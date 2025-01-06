# Contributing Guidelines

Thank you for showing interest in contributing to Kong Ingress Controller.

Following guide will help you navigate the repository and get your PRs merged in faster.

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

If you wish to contribute code (features or bug fixes) please open pull requests with `main` as the base branch or `release/X.Y.x` if the fix shouldn't be targeted at newer versions.

## Working with Pull Requests

When you finish your work and want your contributions to be included, please submit a pull request
and include necessary information about the pull request in its description following our
[pull request template](https://github.com/Kong/kubernetes-ingress-controller/blob/a5cfca4a97927c2878fce593929a8e48442f3127/.github/PULL_REQUEST_TEMPLATE.md).
A clear title and description is helpful to encourage your pull request to be reviewed and merged.

If your pull request is still in progress, but you want to get early reviews before
it is ready to be merged, please add `do not merge` or `work in progress` labels to the
pull request. If your pull request is not ready to be reviewed yet, please convert it to `Draft`.
`Draft` pull requests will not be reviewed until they are converted to `Open` status.

## Stale issue and pull request policy

To ensure our backlog is organized and up to date, we will close issues and
pull requests that have been inactive awaiting a community response for over 2
weeks. If you wish to reopen a closed issue or PR to continue work, please
leave a comment asking a team member to do so.

## When to close an issue

For a GitHub issue describing a problem/feature request:

- **Duplicates**. if there are other issues in the repository describing the same problem/FR:

    1. find the issue that has the most context (possibly not the first reported)

    1. close all other issues with a comment _Duplicate of #XYZ_

- **Resolution by code**. if resolution involves creating PRs:

    1. ensure that all PRs reference the issue they are solving. Keep in mind that the _fixes_/_resolves_ directive only works for PRs merged to the default branch of the repository.

    1. close the issue as soon as all the PRs have been merged to **`main`**. If it's obvious from PRs that the issue has been resolved, a closing comment on the issue is purely optional.

- **Other resolutions/rejections**. if resolution happens for any other reason (_resolved without code_, _user's question answered_, _won't fix_, _infeasible_, _not useful_, _alternative approach chosen_, _problem will go away in $FUTURE-VERSION_)

    1. close the issue with a comment describing the resolution/reason.

For a closed issue, one can verify which released versions contain the fix/enhancement by navigating into the merge commit of each attached PR, where GitHub lists tags/branches that contain the merge commit.
Thus:

- if the list includes a release tag: the fix/enhancement is included in that release tag.
- if the list includes `release/X.Y.x` but no release tags: the fix/enhancement will come in the nearest patch release.
- if the list includes `main` : the fix/enhancement will come in the next minor release.

# Enhancements

Documenting and communicating the motivation for major enhancements in the Kong Kubernetes Ingress Controller (KIC) is done using an upstream Kubernetes process referred to as [Kubernetes Enhancement Proposals (KEPs)][kep].

[kep]:https://github.com/kubernetes/enhancements

## New Enhancement Proposals

When starting a new enhancement proposal use the upstream [KEP Template][template] file as the starting point for your KEP, and follow the instructions therein.

Initially you can remove a lot of the scaffolding in the template for the first `provisional` iteration and focus on establishing the following sections:

- Summary
- Motivation
- Goals
- User Stories

In general the maintainers here feel establishing these things in a KEP should be done _prior to any technical writeups_ but this is a soft rule.

[template]:https://github.com/kubernetes/enhancements/blob/master/keps/NNNN-kep-template/README.md

## Feature Gates

New features should be added to the [Feature Gates][kic-fg] documentation and `internal/manager/feature_gates.go`.

[kic-fg]:https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md

## Development environment

## Environment

- Golang version matching our [`Dockerfile`](./Dockerfile) installed
- [Kubebuilder][kubebuilder]
- [GNU Make][make]
- [Docker][docker] (for building)
- Access to a Kubernetes cluster (we use [KIND][kind] for development)

[kubebuilder]:https://kubebuilder.io/
[make]:https://www.gnu.org/software/make/
[docker]:https://docs.docker.com/
[kind]:https://github.com/kubernetes-sigs/kind

## Dependencies

The build uses dependencies are managed by [go modules](https://blog.golang.org/using-go-modules)

## Developing

Development of our [Kubernetes Controllers][ctrl] and [APIs][kapi] is managed through the [Kubebuilder SDK][kubebuilder].

Prior to developing we recommend you read through the [Makefile](/Makefile) directives related to generation of API configurations, and run through the [Kubebuilder Quickstart Documentation][kbquick] documentation in order to familiarize yourself with how the command line works, how to add new APIs and controllers, and how to update existing APIs.

Make sure you're [generally familiar with Kubernetes Controllers as a concept, and how to build them][kbctrl].

[ctrl]:https://kubernetes.io/docs/concepts/architecture/controller/
[kapi]:https://kubernetes.io/docs/concepts/overview/kubernetes-api/
[kbquick]:https://kubebuilder.io/quick-start.html
[kbctrl]:https://kubebuilder.io/cronjob-tutorial/controller-overview.html

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
kubectl port-forward deploy/ingress-kong -n kong 8444:8444 2>&1 > /dev/null &
kubectl proxy --port=8002 2>&1 > /dev/null &

export POD_NAME=`kubectl get po -n kong -o json | jq ".items[] | .metadata.name" -r | grep ingress`
export POD_NAMESPACE=kong

go run -tags gcp ./internal/cmd/main.go \
--kubeconfig ~/.kube/config \
--publish-service=kong/kong-proxy \
--apiserver-host=http://localhost:8002 \
--kong-admin-url https://localhost:8444 \
--kong-admin-tls-skip-verify true
```

If you are using Kind we can leverage [extraPortMapping config](https://kind.sigs.k8s.io/docs/user/ingress/)

```shell
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 8000
    hostPort: 8000
    protocol: TCP
  - containerPort: 8443
    hostPort: 8443
    protocol: TCP
EOF

# mapping host ports to a kong ingress container port
kubectl patch -n kong deploy ingress-kong -p '{"spec": {"template": {"spec": {"containers": [{"name": "proxy", "ports": [{"containerPort": 8000, "hostPort": 8000, "name": "proxy", "protocol": "TCP"}, {"containerPort": 8443, "hostPort": 8443, "name": "proxy-ssl", "protocol": "TCP"}]}]}}}}'
```

## Building

Build is performed via [`Makefile`](./Makefile). Depending on your
requirements you can build a raw server binary, a local container image,
or push an image to a remote repository.

### Build a raw server binary

```console
make build
```

### Build a local container image

```console
TAG=DEV REGISTRY=docker.example.com/registry make container
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
docker push docker.example.com/registry/kong-ingress-controller:DEV
```

Note: replace `docker.example.com/registry` with your registry URL.

## Deploying

There are several ways to deploy Kong Ingress Controller onto a cluster.
Please check the [deployment guide](https://docs.konghq.com/kubernetes-ingress-controller/latest/deployment/overview/).

## Linting

You can lint the codebase by running:

```console
make lint
```

passing the above is required by CI.

## Testing

Please see [TESTING.md](./TESTING.md) for a detailed description of testing KIC.

## Releasing

See [RELEASE.md](./RELEASE.md).

## Contributor Badge

If your Pull Request to Kong/kubernetes-ingress-controller was accepted, and it fixes a bug, adds functionality, or makes it significantly easier to use or understand KIC, congratulations! You are eligible to receive a digital Contributor Badge! Go ahead and fill out the [Contributors Submissions form](https://goo.gl/forms/5w6mxLaE4tz2YM0L2).

*Badges expire after 1 year, at which point you may submit a new contribution to renew the badge.* 
