## Running Kong Ingress Controller with Istio

In this guide, you will:
* Install Istio v1.6.7 and Kong in your cluster.
* Deploy an example Istio-enabled application (_bookinfo_).
* Deploy an `Ingress` customized with a `KongPlugin` for the example application.
* Make several requests to the sample application via Kong and Istio.
* See the performance metrics of the sample application, provided by Istio.

### Prerequisites

* A Kubernetes v1.15-1.18 cluster which can pull container images from public registries. Examples:
    * A managed Kubernetes cluster (AWS EKS, Google Cloud GKE, Azure AKS),
    * Minikube,
    * `microk8s` with the `dns` addon enabled,
* `kubectl` with admin access to the cluster.

### Download Istio

Download the Istio bundle at version 1.6.7:

```console
$ curl -L https://istio.io/downloadIstio | env ISTIO_VERSION=1.6.7 sh -
...
...
Istio 1.6.7 Download Complete!                                                                                                 
                                                               
Istio has been successfully downloaded into the istio-1.6.7 folder on your system.                                                                                                                                                                            
...
...
```

### Install Istio Operator

Invoke `istioctl` to deploy the Istio Operator to the Kubernetes cluster:

```console
$ ./istio-1.6.7/bin/istioctl operator init
Using operator Deployment image: docker.io/istio/operator:1.6.7
✔ Istio operator installed                                                                                                                                                                                                                                    
✔ Installation complete
```

### Deploy Istio using Operator

Deploy Istio using Istio Operator:

```console
$ kubectl create namespace istio-system
namespace/istio-system created
```
```console
$ kubectl apply -f - <<EOF
  apiVersion: install.istio.io/v1alpha1
  kind: IstioOperator
  metadata:
    namespace: istio-system
    name: example-istiocontrolplane
  spec:
    profile: demo
EOF
istiooperator.install.istio.io/example-istiocontrolplane created
```
```console
$ kubectl describe istiooperator -n istio-system
...
...
Status:
  Status:  RECONCILING
...
...
```

Wait until the `kubectl describe istiooperator` command returns `Status: HEALTHY`.

### Deploy Kong Ingress Controller in an Istio-enabled namespace

```console
$ kubectl create namespace kong-istio
namespace/kong-istio created
```
```console
$ kubectl label namespace kong-istio istio-injection=enabled
namespace/kong-istio labeled
```
```console
$ helm install -n kong-istio example-kong kong/kong --set ingressController.installCRDs=false
...
NAME: example-kong
LAST DEPLOYED: Mon Aug 10 15:14:44 2020
NAMESPACE: kong-istio
STATUS: deployed
...
```

_Optional:_ Run `kubectl describe pod -n kong-istio -l app.kubernetes.io/instance=example-kong` to see that the Istio sidecar (`istio-proxy`) is running alongside Kong Ingress Controller.

### Deploy bookinfo in an Istio-enabled namespace

Deploy the sample _bookinfo_ app from the Istio bundle:

```console
$ kubectl create namespace my-istio-app
namespace/my-istio-app created
```
```console
$ kubectl label namespace my-istio-app istio-injection=enabled
namespace/my-istio-app labeled
kubectl apply -n my-istio-app -f istio-1.6.7/samples/bookinfo/platform/kube/bookinfo.yaml
```
Wait until the application is up:
```console
$ kubectl wait --for=condition=Available deployment productpage -n my-istio-app --timeout=240s
```
### Deploy ingress

Define a `KongPlugin` rate-limiting access to 100 requests per minute. Define an `Ingress` telling Kong to proxy traffic
to a service belonging to the sample application:

```console
$ kubectl apply -f - <<EOF
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: rate-limit
  namespace: my-istio-app
plugin: rate-limiting
config:
  minute: 30
  policy: local
EOF
```

```console
$ kubectl apply -f - <<EOF
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: productpage
  namespace: my-istio-app
  annotations:
    konghq.com/plugins: rate-limit
spec:
  rules:
  - http:
      paths:
      - path: /
        backend:
          serviceName: productpage
          servicePort: 9080
```

### Make some requests to the sample application

Connect to the sample application served via Kong and Istio.

Note that `8080:80` means that `kubectl` will open the `tcp/8080` port on the local system and forward all requests to
Kong's port `80`.

```console
$ # Keep the command below running in the background
$ kubectl port-forward service/example-kong-kong-proxy 8080:80 -n kong-istio
Forwarding from 127.0.0.1:8080 -> 8000
Forwarding from [::1]:8080 -> 8000
...
```

Navigate your web browser to `http://localhost:8080/` You should be able to see a bookstore web application. Click
through any available links several times. As you hit 30 requests per minute (for example, by holding down the "Refresh"
key combination, e.g. `<Ctrl-R>` or `<Command-R>`), you should obtain a `Kong Error - API rate limit exceeded` response.

### See the connection graph in Kiali

Connect to Kiali (the Istio dashboard):

```console
$ # Keep the command below running in the background
$ kubectl port-forward service/kiali 20001:20001 -n istio-system
Forwarding from 127.0.0.1:20001 -> 20001
Forwarding from [::1]:20001 -> 20001
...
```

* Navigate your web browser to `http://localhost:20001/`.
* Log in using the default credentials (`admin`/`admin`).
* Choose _Workloads_ from the menu on the left.
* Select `my-istio-app` in the _Namespace_ drop-down menu.
* Click the _productpage-v1_ service name.
* Click the three dots button in the top-right corner of _Graph Overview_ and click _Show full graph_.
* Select `kong-istio` alongside `my-istio-app` in the _Namespace_ diagram.
* Observe a connection graph spanning from `example-kong-kong-proxy` through `productpage-v1` to the other sample
application services such as `ratings-v1` and `details-v1`.

### See the metrics in Grafana

Connect to Grafana (a dashboard frontend for Prometheus which has been deployed with Istio):

```console
$ # Keep the command below running in the background
$ kubectl port-forward service/grafana 3000:3000 -n istio-system
Forwarding from 127.0.0.1:3000 -> 3000
Forwarding from [::1]:3000 -> 3000
...
```

* Navigate your web browser to `http://localhost:3000/`.
* Expand the dashboard selection drop-down menu from the top of the screen. Expand the `istio` directory and choose the
_Istio Workload Dashboard_ from the list.
* Choose _Namespace: my-istio-app_ and _Workload: productpage-v1_ from the drop-downs.
* Choose a timespan in the top-right of the page to include the time when you made requests to the sample application (e.g. _Last 1 hour_).
* Observe the incoming and outgoing request graphs reflecting actual requests from Kong to `productpage-v1`, and from `productpage-v1` to its backends.

Note that the requests from the web browser to Kong are not reflected in inbound stats of `example-kong-kong-proxy`
because we've issued these requests by `kubectl port-forward`, thus bypassing the Istio proxy sidecar in Kong.
