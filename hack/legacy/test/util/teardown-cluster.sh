#!/bin/bash

kubectl config set-context --current --namespace kong

echo -e "\n--------------------- kubectl get all -A ---------------------\n"
kubectl get all -A
echo -e "\n--------------------------------------------------------------\n"

echo -e "\n---------- kubectl describe deployment ingress-kong ----------\n"
kubectl describe deployment ingress-kong
echo -e "\n--------------------------------------------------------------\n"

for POD in $(kubectl get pods -o=go-template='{{range .items}}{{.metadata.name}}{{end}}')
do
    echo -e "\n----------------- kubectl describe pod ${POD} ----------------\n"
    kubectl describe pod ${POD}
    echo -e "\n------------ kubectl logs ${POD} ingress-controller ----------\n"
    kubectl logs ${POD} ingress-controller
    echo -e "\n------------------ kubectl logs ${POD} proxy -----------------\n"
    kubectl logs ${POD} proxy
    echo -e "\n--------------------------------------------------------------\n"
done

"$KIND_BINARY" delete cluster "--name=$CLUSTER_NAME"
docker rm "$REGISTRY_NAME" -f

