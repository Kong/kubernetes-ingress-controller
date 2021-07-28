#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE}")/.."

DIFFROOT="${SCRIPT_ROOT}/deploy/single-v2/"

cleanup() {
  git checkout "${DIFFROOT}"
}
trap "cleanup" EXIT SIGINT

if ! git status --porcelain --untracked-files=no "$DIFFROOT" ; then
    echo "error: please run this script on a clean working copy"
    exit 1
fi


"${SCRIPT_ROOT}/hack/deploy/build-single-manifests.sh"

if git diff --quiet "${DIFFROOT}"
then
  echo "${DIFFROOT} up to date."
else
  echo "${DIFFROOT} is out of date. Please run hack/build-single-manifests.sh"
  echo "Diff output:"
  git --no-pager diff "${DIFFROOT}"
  exit 1
fi
