#!/bin/bash

# TODO: for now our e2e tests are effectively just our integration tests run
#       against a conformant, production grade cluster. In the future we will
#       add a dedicated e2e test suite.
#
#       See: https://github.com/Kong/kubernetes-ingress-controller/issues/1605

set -euo pipefail

WORKDIR="$(dirname "${BASH_SOURCE}")/../.."
cd "${WORKDIR}"

CLUSTER_NAME="e2e-$(uuidgen)"
KUBERNETES_CLUSTER_NAME="${CLUSTER_NAME}" go run hack/e2e/main.go

function cleanup() {
    KUBERNETES_CLUSTER_NAME="${CLUSTER_NAME}" go run hack/e2e/main.go cleanup
}
trap cleanup EXIT SIGINT SIGQUIT

GOFLAGS="-tags=integration_tests" KONG_TEST_CLUSTER="gke:${CLUSTER_NAME}" go test -count 1 -timeout 45m -v ./test/integration/...
