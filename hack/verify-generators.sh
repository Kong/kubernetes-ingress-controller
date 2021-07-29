#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

DIFFROOT="$(dirname "${BASH_SOURCE}")/.."

if ! git status --porcelain --untracked-files=no "$DIFFROOT" ; then
    echo "error: please run this script on a clean working copy"
    exit 1
fi

cleanup() {
  git checkout "${DIFFROOT}"
}
trap "cleanup" EXIT SIGINT

cd "${DIFFROOT}"
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
