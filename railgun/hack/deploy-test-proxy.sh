#!/bin/bash

set -euox pipefail

helm repo add kong https://charts.konghq.com
helm repo update
helm install kong-test-proxy kong/kong \
    --namespace kong-system \
    --set ingressController.enabled=false \
    --set replicaCount=3 \
    --set admin.enabled=true \
    --set admin.type=ClusterIP

PROXY_PODS="$(kubectl -n kong-system get pods -o=go-template='{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | perl -ne 'print if m{^kong-test-proxy-kong-}')"

for POD in ${PROXY_PODS}
do
    kubectl -n kong-system label --overwrite=true pods ${POD} konghq.com/proxy-instance=true
done

kubectl -n kong-system get pods --selector konghq.com/proxy-instance
