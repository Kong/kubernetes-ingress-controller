#!/bin/bash

SCRIPT_DIR="$(dirname "$BASH_SOURCE")"
export KUBECONFIG=$PWD/kubeconfig-test-cluster
export CLUSTER_NAME="test-cluster"
export REGISTRY_NAME="test-local-registry"

export KIND_BINARY=./kind
export KIND_URL=https://github.com/kubernetes-sigs/kind/releases/download/v0.9.0/kind-linux-amd64

wget "$KIND_URL" -O "$KIND_BINARY" || exit 1
chmod +x "$KIND_BINARY" || exit 1

"$SCRIPT_DIR/util/create-cluster.sh" || exit 1
if [ -z "$SKIP_TEARDOWN" ]; then
	trap "$SCRIPT_DIR/util/teardown-cluster.sh" EXIT
fi
"$SCRIPT_DIR/util/install-sut.sh" || exit 1
"$SCRIPT_DIR/util/run-all-tests.sh" || exit 1
