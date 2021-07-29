#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE}")/.."
DIFFROOT="${SCRIPT_ROOT}"

cleanup() {
  git checkout "${DIFFROOT}"
}
trap "cleanup" EXIT SIGINT

if ! git status --porcelain --untracked-files=no "$DIFFROOT" ; then
    echo "error: please run this script on a clean working copy"
    exit 1
fi

cd "${SCRIPT_ROOT}"
make generate

if git diff --quiet "${DIFFROOT}"
then
  echo "${DIFFROOT} up to date."
else
  echo "${DIFFROOT} is out of date. Please run make generate"
  echo "Diff output:"
  git --no-pager diff "${DIFFROOT}"
  exit 1
fi
