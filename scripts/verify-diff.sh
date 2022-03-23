#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE}")/.."

if git diff --quiet "${SCRIPT_ROOT}"
then
  echo "${SCRIPT_ROOT} up to date."
else
  echo "${SCRIPT_ROOT} appears to be out of date (make sure you've run 'make manifests' and 'make generate')"
  echo "Diff output:"
  git --no-pager diff "${SCRIPT_ROOT}"
  exit 1
fi
