#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE}")/.."

if ! git status --porcelain --untracked-files=no "$SCRIPT_ROOT" ; then
    echo "error: repository not clean (changed files found)"
    exit 1
fi
