#!/bin/bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

export JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'

echo "deploying Kong Ingress controller"
cat deploy/single/all-in-one-postgres.yaml | kubectl create -f -

echo "updating image..."
kubectl set image \
    deployments \
    --namespace kong \
	--selector app=ingress-kong \
    ingress-controller=kong-docker-kubernetes-ingress-controller.bintray.io/kong-ingress-controller:test

sleep 5


function waitForPod() {
    until kubectl get pods -n kong -l app="$1" -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True";
    do
        sleep 1;
    done
}

export -f waitForPod

echo "waiting Postgres pod..."
timeout 30s bash -c waitForPod postgres
echo "waiting Kong pod..."
timeout 30s bash -c waitForPod kong
echo "waiting Kong ingress pod..."
timeout 3m  bash -c waitForPod ingress-kong

if kubectl get pods -n kong -l app=ingress-kong -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True";
then
    echo "Kubernetes deployments started"
else
    echo "Kubernetes deployments with issues:"
    kubectl get pods -n kong

    echo "Reason:"
    kubectl describe pods -n kong
    kubectl logs -n kong -l app=ingress-kong ingress-controller
    exit 1
fi
