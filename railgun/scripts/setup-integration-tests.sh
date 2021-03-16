#!/bin/bash
# TODO: (shane) this is all going to be replaced by a container image soon

set -euox pipefail

KIND_VERSION="v0.10.0"

# ensure docker command is accessible
if ! command -v docker &> /dev/null
then
    echo "ERROR: docker command not found"
    exit 10
fi

# ensure docker is functional
docker info

# ensure kind command is accessible
if ! command -v kind &> /dev/null
then
    go get -v sigs.k8s.io/kind@${KIND_VERSION}
fi

# ensure kind is functional
kind version
