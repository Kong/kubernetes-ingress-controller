# Kong Ingress on Google Kubernetes Engine (GKE)

## Requirements

1. A fully functional GKE cluster.
   The easiest way to do this is to do it via the web UI:
   Go to Google Cloud's console > Kubernetes Engine > Cluster >
   Create a new cluster.
   This documentation has been tested on a zonal cluster in
   europe-west-4a using 1.10.5-gke.4 as Master version.
   The default pool has been assigned 2 nodes of kind 1VCPU
   with 3.75GB memory (default setting).
   The OS used is COS (Container Optimized OS) and the auto-scaling
   has been enabled. Default settings are being used except for
   `HTTP load balancing` which has been disabled (you probably wanna use
   Kong features for this). For more information on GKE clusters,
   refer to
   [the GKE documentation](https://cloud.google.com/kubernetes-engine/docs/)
1. If you wish to use a static IP for Kong, you have to reserve a static IP
   address (in Google Cloud's console > VPC network >
   External IP addresses). For information,
   you must create a regional IP
   global is not supported as `loadBalancerIP` yet)
1. Basic understanding of Kubernetes
1. A working `kubectl`  linked to the GKE Kubernetes
   cluster we will work on. For information, you can associate a new `kubectl`
   context by using:

   ```bash
   gcloud container clusters get-credentials <my-cluster-name> --zone <my-zone> --project <my-project-id>
    ```

## Deploy Kong Ingress Controller

## Downloads basic resources

It's recommended to keep your Kuberenetes configuration versioned,
so we will first download basic resources:

```bash

wget https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/namespace.yaml
wget https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/custom-types.yaml
wget https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/postgres.yaml
wget https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/rbac.yaml
wget https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/ingress-controller.yaml
wget https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/kong.yaml
wget https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/provider/cloud/gke-kong-proxy-loadbalancer-service.yaml

```

This command will create  followings files:
- `namespace.yaml` : The definition of the Kong's namespace
- `custom-types.yaml` : Contains Custom Types Definition
- `postgres.yaml` : Contains PostgreSQL deployment.
   You are free to use a DaaS like CloudSQL if needed.
   If you do so, this file should not be used and you have to configure
   `kong.yaml` and `ingress-controller.yaml` accordingly.
- `rbac.yaml` : Contains Service-Account and Roles definitions used by Kong.
- `ingress-controller.yaml` : Deployment file of the Kong Ingress Controller.

    Note: Ingress Controller is deployed as `NodePort` exposing ports
    internally into node's private network,
    you might edit the Service `kong-ingress-controller`
    to use `ClusterIP` Type as follows:

    ```yaml

    apiVersion: v1
    kind: Service
    metadata:
      name: kong-ingress-controller
      namespace: kong

      type: ClusterIP
      ports:
      - name: kong-admin
        port: 8001
        targetPort: 8001
        protocol: TCP
      selector:
        app: ingress-kong

    ```

- `kong.yaml` :  Deployment file of Kong
- `gke-kong-proxy-loadbalancer-service.yaml` : The Service associated to Kong

You can now edit these files to fill your need
(for example, you can edit them to use Google's CloudSQL
for the PostgreSQL database, or edit `KONG_ADMIN_LISTEN`
if you wish to access Kong admin API).
If you want to use a static IP, add the IP value as
`loadBalancerIP` in the Service `kong-proxy` in the file
gke-kong-proxy-loadbalancer-service.yaml. For example:

```yaml

apiVersion: v1
kind: Service
metadata:
  name: kong-proxy
  namespace: kong
spec:
  externalTrafficPolicy: Local
  type: LoadBalancer
  ports:
  - name: kong-proxy
    port: 80
    targetPort: 8000
    protocol: TCP
  - name: kong-proxy-ssl
    port: 443
    targetPort: 8443
    protocol: TCP
  loadBalancerIP: <my-reserved-ip-address>
  selector:
    app: kong

```

## Update User Permissions

> [Because of the way Kubernetes Engine checks permissions
when you create a Role or ClusterRole, you must
first create a RoleBinding that grants you all of
the permissions included in the role you want to create.
An example workaround is to create a RoleBinding that
gives your Google identity a cluster-admin role
before attempting to create additional Role or
ClusterRole permissions.
This is a known issue in RBAC in Kubernetes and
Kubernetes Engine versions 1.6 and
later.](https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control).

A fast workaround:

```yaml

echo -n "
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: cluster-admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: User
  name: <the current user using kubectl ; usually the google account>
  namespace: kube-system" | kubectl apply -f -

```

## Deploy Kong

You can now deploy Kong:

```bash

kubectl apply -f namespace.yaml &&
kubectl apply -f custom-types.yaml &&
kubectl apply -f postgres.yaml &&
kubectl apply -f rbac.yaml &&
kubectl apply -f ingress-controller.yaml &&
kubectl apply -f kong.yaml  &&
kubectl apply -f gke-kong-proxy-loadbalancer-service.yaml

```

Should display:

```bash

namespace "kong" created
customresourcedefinition "kongplugins.configuration.konghq.com" configured
customresourcedefinition "kongconsumers.configuration.konghq.com" configured
customresourcedefinition "kongcredentials.configuration.konghq.com" configured
customresourcedefinition "kongingresses.configuration.konghq.com" configured
service "postgres" created
statefulset "postgres" created
serviceaccount "kong-serviceaccount" created
clusterrole "kong-ingress-clusterrole" configured
role "kong-ingress-role" created
rolebinding "kong-ingress-role-nisa-binding" created
clusterrolebinding "kong-ingress-clusterrole-nisa-binding" configured
service "kong-ingress-controller" created
deployment "kong-ingress-controller" created
deployment "kong" created
service "kong-proxy" created

```

You can now retrieve the associated IP for the Service `kong-proxy`
(or you can use  directly your static IP if you used one):

`kubectl get services -n kong` should display :

```bash

NAME                      TYPE           CLUSTER-IP      EXTERNAL-IP    PORT(S)
kong-ingress-controller   ClusterIP      10.42.42.1   <none>         8001/TCP
kong-proxy                LoadBalancer   10.42.42.2   203.0.113.1   80:30095/TCP,443:31166/TCP
postgres                  ClusterIP      10.42.42.3   <none>         5432/TCP

```

Now,

```bash

curl 203.0.113.1

```

Should display:

```bash

{"message":"no route and no API found with those values"}

```

> Note: It may take a while for Google to actually associate the
IP address to the `kong-proxy` Service.

## Test your deployment

- Deploy a dummy application :
  
  ```bash

  kubectl create namespace dummy && curl https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/master/deploy/manifests/dummy-application.yaml -n dummy

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
      - host: dummy.kong.example
        http:
          paths:
            - path: "/"
              backend:
                serviceName: http-svc
                servicePort: http" | kubectl apply -f -

  ```

- Edit your /etc/hosts and add:

  ```text

  203.0.113.1 dummy.kong.example

  ```

Now, access to dummy.kong.example should display some information.

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

### Setup TLS (HTTPS)

You need to set your API with HTTPS in order
to expose your service securely. In this section,
I will explain how to secure it with [Let’s Encrypt](https://letsencrypt.org/).

#### Register your domain

First of all, you must register your domain with any domain
registration services such as [Google Domains](https://domains.google/).

#### Follow instructions of Let’s Encrypt on GKE

[Let’s Encrypt on GKE](https://github.com/ahmetb/gke-letsencrypt)
is a tutorial for installing `cert-manager` to get HTTPS certificates
from Let’s Encrypt.
There is an important thing you need to configure,
if you want to accomplish correctly. You should apply
[KongIngress](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/custom-types.md#kongingress)
and set `preserve_host` configuration `true` at the
[4th step](https://github.com/ahmetb/gke-letsencrypt/blob/master/40-deploy-an-app.md)
so that you could keep hostname in request headers.

[cert-manager](https://github.com/jetstack/cert-manager) checks equality
of hostname and domain name when it creates HTTPS certificates.
However, Kong removes hostname as default.
I recommend you to create a `KongIngress` spec file to avoid
the following error:

```bash

[dummy.kong.example] Invalid host 'xxx.xxx.xxx.xxx'

```

These are examples of `KongIngress` and `Ingress` spec.

```sh

echo -n "
apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
  name: sample-kong-ingress
  namespace: kong
proxy:
  path: /
route:
  protocols:
  - https
  - http
  strip_path: false
  preserve_host: true" | kubectl apply -f -

```

```sh

echo -n "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: dummy
  namespace:  dummy
  annotations:
    kubernetes.io/ingress.class: "nginx"
    configuration.konghq.com: sample-kong-ingress
spec:
  rules:
    - host: dummy.kong.example
      http:
        paths:
          - path: "/"
            backend:
              serviceName: http-svc
              servicePort: http" | kubectl apply -f -

```