#!/bin/bash

HAS_ERR=false
while read -r line ; do
  if [ $(gh issue view --json=state "${line##*/}" | jq '.state' | sed 's/\"//g') = "CLOSED" ]; then
    echo "closed issue reference: $line"
    HAS_ERR=true
  fi
done < <(grep -roiIE "//( )*TODO(:){0,1}( )*(https://){0,1}github.com/kong/kubernetes-ingress-controller/issues/[0-9]+" .)
if $HAS_ERR; then
  exit 1
fi
