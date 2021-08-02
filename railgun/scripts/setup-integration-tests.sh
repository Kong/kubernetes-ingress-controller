#!/bin/bash

set -euox pipefail

KIND_VERSION="v0.11.1"

# ensure docker command is accessible
if ! command -v docker &> /dev/null
then
    echo "ERROR: docker command not found"
    exit 10
fi

# ensure docker is functional
docker info 1>/dev/null

# ensure kind command is accessible
if ! command -v kind &> /dev/null
then
    go get -v sigs.k8s.io/kind@${KIND_VERSION}
fi

# ensure kind is functional
kind version
