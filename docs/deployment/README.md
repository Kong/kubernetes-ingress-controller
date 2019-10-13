# Deploying Kong Ingress Controller

Kong ingress controller can be installed on a local, managed
or any Kubernetes cluster which supports a service of type `LoadBalancer`.

Here are some guides to get you started:

1. [Using minikube][0]:

   This guide helps you get Kong Ingress Controller on a local
   Kubernetes cluster.

2. [Elastic Kubernetes Service][1]:

   [EKS](https://aws.amazon.com/eks/) is a managed Kubnernetes cluster
   offering by Amazon Web Services. This guide is a walkthrough to set up
   Kong Ingress Controller on EKS.

3. [Google Kubernetes Engine(GKE)][2]:

   [GKE](https://cloud.google.com/kubernetes-engine/)
   is a managed Kubernetes cluster offering by Google.
   This guide is a walkthrough to set up Kong Ingress Controller on GKE.
   If you've access to GKE, please use this guide.

4. [Azure Kubernetes Service(AKS))][3]:

   [AKS](https://azure.microsoft.com/en-us/services/kubernetes-service/)
   is a managed Kubernetes cluster offering by Microsoft Azure.
   This guide is a walkthrough to set up Kong Ingress
   Controller on AKS.

Once you've Kong Ingress Controlled installed, please follow our
[getting started](../guides/getting-started.md) tutorial to learn
about how to use the Ingress Controller.

## Deploying Admission Controller

Kong Ingress Controller also ships with a Validating
Admission Controller that
can be enabled to verify KongConsumer and KongPlugin resources as they
are created.
Please follow the [admission-webhook](admission-webhook.md) deployment
guide to set it up.

[0]: minikube.md
[1]: eks.md
[2]: gke.md
[3]: aks.md


## Deploying Kong Enterprise Ingress Controller

Please follow [this guide](kong-enterprise.md) to deploy Kong Enterprise Ingress Controller
if you have purchased or are trying out Kong Enterprise.
