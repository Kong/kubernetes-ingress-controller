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
POST_PROCESSED_DOC="${1}"

# Add a title and turn the vale linter off
echo "---
title: Custom Resource (CRD) API Reference

description: |
  See the generated CRD structure containing all possible resource fields and descriptions for each property.

content_type: reference
layout: reference
tags:
  - crd
search_aliases: 
  - kic CRD
products:
  - kic
breadcrumbs:
  - /kubernetes-ingress-controller/
works_on:
  - on-prem
  - konnect
related_resources:
  - text: Gateway API
    url: /kubernetes-ingress-controller/gateway-api/
---

<!-- vale off -->
" > "${POST_PROCESSED_DOC}"

# Add the generated doc content
cat "${CRD_REF_DOC}" >> "${POST_PROCESSED_DOC}"

# Turn the linter back on
echo "<!-- vale on -->" >> "${POST_PROCESSED_DOC}"

SED=sed
if [[ $(uname -s) == "Darwin" ]]; then
  if gsed --version 2>&1 >/dev/null ; then
    SED=gsed
  else
    echo "GNU sed is required on macOS. You can install it via Homebrew with 'brew install gnu-sed'."
    exit 1
  fi
fi

# Replace all description placeholders with proper include directives
${SED} -i \
  's/<!-- \(.*\) description placeholder -->/{% include md\/kic\/crd-ref\/\1_description.md kong_version=page.kong_version %}/' \
  "${POST_PROCESSED_DOC}"
