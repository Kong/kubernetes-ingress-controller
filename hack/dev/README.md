# Development Environment

This guide shows how to setup a development environment for Kong Ingress
Controller.

## Pre-requisite

- A Kubernetes cluster: You can use Minikube, kind or use a managed cluster
  service like GKE.
- In your workstation, you need following installed:
  - Go (v1.13+)
  - kubectl
  - kustomize
  - Docker (for build purposes only)

## Overview

In order to develop the Ingress Controller, we are going to set up the following
architecture:

- Run Kong inside a Kubernetes cluster
- Proxy Kubernetes API-server to port 8002 on the workstation
- Proxy Kong's Admin API and Proxy port on the workstation
- Run the Ingress Controller locally (`go run`)


## Setup resources in Kubernetes

### Common resources

First, we will install the common required resources:

```
kustomize build ./hack/dev/common | kubectl apply -f -
```

### Install Kong

To install Kong, we will use a deploy script.
This script is capable of install Kong in DB or DB-less mode and install
various versions of Kong as well.

To install DB-less version:

```
./hack/dev/deploy.sh dbless 2.0
```

To install Kong running without a database:

```
./hack/dev/deploy.sh db 2.0
```

Substitute `2.0` with another version to install a different version of Kong.
The possible versions are listed in `hack/dev/db/` or
`hack/dev/dbless/` directories.


### Run the Ingress Controller


Once installed, you can now run the Ingress Controller locally:

```
make run
```

Please ensure that the current context in `kubectl` is set appropriately.

Please also note that the RBAC profile being used is of the current user
of `kubectl`, likely your RBAC profile.
    `kubectl apply -f hack/dev-env/dev-config.yaml`


## Running multiple version of Kong

You can deploy multiple version of Kong at the same time in the same
k8s cluster using the deploy script described above.

You can override the version and mode (db/dbless) to use using RUN_VERSION
and DB environment variables to the `make run` command:

To run the controller against 1.4 DB version of Kong:
```
DB=true RUN_VERSION=14 make run
```

Please note that you will have to ensure that Kong 1.4 with DB is deployed
already using the above deploy script.

This setup is designed to keep the code, test, repeat cycle fast and lean.
If you run into issues, please feel free to open Github issues.
