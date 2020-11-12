#!/bin/bash

"$KIND_BINARY" delete cluster "--name=$CLUSTER_NAME"
docker rm "$REGISTRY_NAME" -f

