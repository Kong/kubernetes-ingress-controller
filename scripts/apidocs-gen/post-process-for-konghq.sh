#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")/../.."
CRD_REF_DOC="${SCRIPT_ROOT}/docs/api-reference.md"
TEMPORARY_POST_PROCESSED_DOC="$1"

echo "---
title: Custom Resource Definitions API Reference
---
<!-- vale off -->
" > "${TEMPORARY_POST_PROCESSED_DOC}"
cat "${CRD_REF_DOC}" >> "${TEMPORARY_POST_PROCESSED_DOC}"
echo "<!-- vale on -->" >> "${TEMPORARY_POST_PROCESSED_DOC}"
