#!/bin/bash

# This script adapts auto-generated api-reference.md to the requirements of
# docs.konghq.com:
#   - adds a title section
#   - turns vale linter off for the whole document

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")/../.."
CRD_REF_DOC="${SCRIPT_ROOT}/docs/api-reference.md"
POST_PROCESSED_DOC="$1"

echo "---
title: Custom Resource Definitions API Reference
---
<!-- vale off -->
" > "${POST_PROCESSED_DOC}"
cat "${CRD_REF_DOC}" >> "${POST_PROCESSED_DOC}"
echo "<!-- vale on -->" >> "${POST_PROCESSED_DOC}"
