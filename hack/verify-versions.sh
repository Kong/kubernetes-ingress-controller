#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."
MANIFEST_ROOT="$REPO_ROOT/deploy/single"
FAILED=0

MANIFESTS="$MANIFEST_ROOT/*"
for MANIFEST in ${MANIFESTS}
do
	CONTAINERS=$(yq eval-all ".spec.template.spec.containers[].image" "${MANIFEST}" -N)
	INITCONTAINERS=$(yq eval-all ".spec.template.spec.initContainers[].image" "${MANIFEST}" -N)
	KONGS=$(printf "%s\n%s" "$CONTAINERS" "$INITCONTAINERS" | sort | uniq | grep -oP "kong(\/kong-gateway)?:[\d\.]+")
	if [ "$(echo "${KONGS}" | wc -l)" -gt 1 ]
	then
		echo "multiple Kong images in $MANIFEST, verify image consistency in source:"
		echo "$KONGS"
		FAILED=$((FAILED+1))
	fi
	if [ $FAILED -gt 0 ]
	then
		exit 1
	fi
done
