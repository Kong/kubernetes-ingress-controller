#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")/..

DIFFROOT="${SCRIPT_ROOT}/deploy/single/"

cleanup() {
  git checkout "${DIFFROOT}"
}
trap "cleanup" EXIT SIGINT

cleanup

"${SCRIPT_ROOT}/hack/build-single-manifests.sh"
echo "diffing ${DIFFROOT} against freshly generated single manifests"
ret=0
git diff --quiet "${DIFFROOT}" || ret=$?
if [[ $ret -eq 0 ]]
then
  echo "${DIFFROOT} up to date."
else
  echo "${DIFFROOT} is out of date. Please run hack/build-single-manifests.sh"
  git checkout "${DIFFROOT}"
  exit 1
fi
