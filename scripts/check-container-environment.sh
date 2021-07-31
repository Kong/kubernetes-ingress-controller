#!/bin/bash

# check-container-environment.sh
# ---
#
# This script is used to validate that the current environment has containers
# available which can serve our testing suites, and will download a copy of
# Kubernetes In Docker (KIND) if needed so that we can deploy a Kubernetes
# cluster inside a container for testing purposes. This script is generally only
# run inside CI.

set -euo pipefail

KIND_VERSION="v0.11.1"

# ensure docker command is accessible
if ! command -v docker &> /dev/null
then
    echo "ERROR: docker command not found"
    exit 10
fi

# ensure kind command is accessible
if ! command -v kind &> /dev/null
then
    go get -v sigs.k8s.io/kind@${KIND_VERSION}
fi

# ensure kubectl command is accessible
if ! command -v kubectl &> /dev/null
then
    echo "ERROR: kubectl command not found"
    exit 10
fi

# ensure helm command is accessible
if ! command -v helm &> /dev/null
then
    echo "ERROR: helm command not found"
    exit 10
fi

DOCKER_VERSION="$(docker -v)"
KIND_VERSION="$(kind version)"
KUBECTL_VERSION="$(kubectl version --short --client)"
HELM_VERSION="$(helm version --short)"

echo "INFO: container environment ready
INFO: DOCKER=(${DOCKER_VERSION})
INFO: KIND=(${KIND_VERSION})
INFO: KUBECTL=(${KUBECTL_VERSION})
INFO: HELM=(${HELM_VERSION})"
