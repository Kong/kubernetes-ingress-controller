#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."
OSS_IMAGES="$REPO_ROOT/config/image/oss/kustomization.yaml"

EXPECTED=$1
KIC_VERSION=$(yq eval-all '.images[] | select(.name=="kic-placeholder") | .newTag' "${OSS_IMAGES}")

if [ "${EXPECTED}" != "${KIC_VERSION}" ]
then
		echo "KIC version in ${OSS_IMAGES} is ${KIC_VERSION}, expected ${EXPECTED}"
		exit 1
fi
