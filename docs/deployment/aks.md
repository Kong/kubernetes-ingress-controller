# Kong Ingress on Azure Kubernetes Service (AKS)

## Requirements

1. A fully functional AKS cluster.
   Please follow Azure's Guide to
   [setup an AKS cluster](https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough).
1. Basic understanding of Kubernetes
1. A working `kubectl`  linked to the AKS Kubernetes
   cluster we will work on. The above AKS setup guide will help
   you set this up.

## Deploy Kong Ingress Controller

Deploy Kong as Ingress controller:

```shell

$ kubectl apply -f https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/single/all-in-one-postgres.yaml

namespace/kong created
customresourcedefinition.apiextensions.k8s.io/kongplugins.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongconsumers.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongcredentials.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongingresses.configuration.konghq.com created
service/postgres created
statefulset.apps/postgres created
serviceaccount/kong-serviceaccount created
clusterrole.rbac.authorization.k8s.io/kong-ingress-clusterrole created
role.rbac.authorization.k8s.io/kong-ingress-role created
rolebinding.rbac.authorization.k8s.io/kong-ingress-role-nisa-binding created
clusterrolebinding.rbac.authorization.k8s.io/kong-ingress-clusterrole-nisa-binding created
service/kong-ingress-controller created
deployment.extensions/kong-ingress-controller created
service/kong-proxy created
deployment.extensions/kong created

```

It will take a few minutes for all containers to start and report
healthy status.

You can now retrieve the associated IP for the Service `kong-proxy`

```bash

$ kubectl get services -n kong

NAME                      TYPE           CLUSTER-IP      EXTERNAL-IP    PORT(S)
kong-ingress-controller   ClusterIP      10.42.42.1   <none>         8001/TCP
kong-proxy                LoadBalancer   10.42.42.2   203.0.113.42   80:30095/TCP,443:31166/TCP
postgres                  ClusterIP      10.42.42.3   <none>         5432/TCP

```

Now,

```bash

curl 203.0.113.42

```

Should display:

```bash

{"message":"no route and no API found with those values"}

```

> Note: It may take a while to actually associate the
IP address to the `kong-proxy` Service so please be patient.

## Test your deployment

- Deploy a dummy application :

  ```shell

  $ kubectl create namespace dummy
  namespace/dummy created
  $ kubectl apply -n dummy -f https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/dummy-application.yaml
  deployment.extensions/http-svc created
  service/http-svc created

  ```

- Add an Ingress:

  ```yaml

  echo -n "
  apiVersion: extensions/v1beta1
  kind: Ingress
  metadata:
    name: dummy
    namespace:  dummy
    annotations:
      kubernetes.io/ingress.class: "nginx"
  spec:
    rules:
      - host:
        http:
          paths:
            - path: "/"
              backend:
                serviceName: http-svc
                servicePort: http" | kubectl apply -f -

  ```

- Edit your /etc/hosts and add:

  ```text

  203.0.113.42 dummy.kong.example

  ```

Now, access to dummy.kong.example should display some informations.

## Bonus: Expose the Kong admin API

If you want to expose the Kong admin API,
you must configure Kong correctly via `KONG_ADMIN_LISTEN` and add an Ingress:

```yaml

echo -n "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: kong-admin
  namespace:  kong
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
  rules:
    - host: dummy.kong.example
      http:
        paths:
          - path: "/"
            backend:
              serviceName: kong-ingress-controller
              servicePort: 8001" | kubectl apply -f -

```

Do keep in mind that anyone can configure Kong at this point.
You can use one of the Kong's authentication plugin to protect
the Admin API on the Internet.