# Introduction

This directory contains three scripts that allow us to install the same application used in the [minikube guide](0) and create 500 Ingress rules.
During the creation, we also create a new Kubernetes service per Ingress to make sure we have multiple Services, Upstreams and Targets in Kong.
By doing this we make sure the pagination used by the controller is working properly and get a more realistic scenario than just creating an Ingress

### How to run the test:

1. Download the repository
2. execute the script `setup.sh`
3. execute the script `up.sh`
4. Check the ingress controller and Kong admin API logs
5. Wait until the process ends
6. Optionally we can scale the http-svc deployment replica count to see how the upstream targets are updated running:

    `kubectl scale deployment --namespace batch-demo http-svc --replicas=4`

7. Scaling the replica count to one should remove three targets

    `kubectl scale deployment --namespace batch-demo http-svc --replicas=1`

8. To remove all the created resources execute the script `down.sh`

[0]: ../../deploy/minikube.md
