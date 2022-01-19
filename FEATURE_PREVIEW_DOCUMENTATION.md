# Feature Preview Documentation

This document partners with the [feature gates documentation][fg] and includes early draft documentation for features that are currently considered [alpha maturity][fg-alpha].

As features mature to `BETA` or `GA` state this documentation will be grown and eventually will graduate to the [main Kong Documentation][kong-docs].

**WARNING**: the documentation here should be considered experimental and for development use. DO NOT USE FOR PRODUCTION.

[fg]:/FEATURE_GATES.md
[fg-alpha]:/FEATURE_GATES.md#feature-gates-for-alpha-or-beta-features
[kong-docs]:https://github.com/Kong/docs.konghq.com

## Feature Preview: Gateway APIs

[Gateway APIs][gwapis] is an upcoming upstream API for Kubernetes which is intended to functionally replace the [Ingress API][ingress].

In the previous `Ingress` model the actual "gateway" (the underlying proxy server) was implied by configuring an "ingress class" on the `Ingress` resource to identify the otherwise deployed gateway that should serve the traffic defined by the object. With the advent of Gateway APIs the `Gateway` is an actual resource in the cluster in addition to referencing an underlying proxy server.

Gateway APIs natively supports more options for ingress network traffic than was previously possible with `Ingress`, including (but not limited to):

- TCP routing
- UDP routing
- TLS routing
- Extended featuresets for HTTP routing (over `Ingress`)

For a more complete overview of Gateway APIs please see the [Gateway APIs documentation][gw-docs].

[gwapis]:https://github.com/kubernetes-sigs/gateway-api
[ingress]:https://kubernetes.io/docs/concepts/services-networking/ingress/
[gw-docs]:https://gateway-api.sigs.k8s.io/

### API Overview

The following Gateway APIs are available in the current build of the Kong Ingress Controller (KIC):

- [GatewayClass][gwc-api]
- [Gateway][gw-api]
- [HTTPRoute][httproute-api]

Note that some notable features of these APIs are currently unsupported by this implementation:

- `Gateway` resources are not yet capable of being provisioned by the controller yet, they must currently reference an existing Kong Gateway
- `HTTPRoute` does not currently support regex based path matching
- `HTTPRoute` does not currently support regex based header matching
- `HTTPRoute` does not currently support query param based matching

**NOTE**: only the APIs which are supported by this implementation are currently documented here, for a full upstream list see [the api reference][gw-ref]

[gwc-api]:https://gateway-api.sigs.k8s.io/v1alpha2/api-types/gatewayclass/
[gw-api]:https://gateway-api.sigs.k8s.io/v1alpha2/api-types/gateway/
[httproute-api]:https://gateway-api.sigs.k8s.io/v1alpha2/api-types/httproute/
[gw-ref]:https://gateway-api.sigs.k8s.io/v1alpha2/api-types/gatewayclass/

### Example - Gateway & HTTPRoute

Current functionality supports routing traffic with the `HTTPRoute` API.

In order to route HTTP traffic you'll need to:

- have an existing Kubernetes cluster with the KIC deployed and the Gateway feature enabled
- create a `GatewayClass` which refers to your KIC controller
- create a `Gateway` which links to the `GatewayClass`
- create an `HTTPRoute` which is bound to `Listeners` on the `Gateway`

In the following sections we'll walk through each of these steps.

This documentation assumes that you're running the _latest_ copy of the KIC on your development/testing cluster deployed via [Helm][helm]. If that's true then you can skip the next section, otherwise move on to the [local development cluster setup section](/#local-development-cluster-setup-optional).

[helm]:https://helm.sh

#### Local Development Cluster Setup (Optional)

If you're starting from scratch and just trying to test Gateway support quickly you can use [Kubernetes In Docker (KIND)][kind] to quickly spin up a local Kubernetes cluster and try things out. Make sure you follow the [installation instructions for Kind on your system][kind-install] before proceeding to the next steps. Also ensure that you have the latest version of [kubectl][kubectl] and [helm][helm] installed in order to deploy the KIC to the cluster.

Create the cluster by running the following:

```console
$ kind create cluster --name kic-gateway-testing
$ kubectl cluster-info --context kind-kic-gateway-testing
```

Once that's completed ensure that you have the Kong charts repo registered with Helm:

```console
$ helm repo add kong https://charts.konghq.com
$ helm repo update
```

Start a configuration file to manage your Kong deployment named `values.yaml` and make sure you keep it:

```yaml
proxy:
  type: NodePort
  http:
    nodePort: 32080
```

Deploy KIC to the cluster with Helm:

```console
$ helm install --create-namespace --namespace kong -f values.yaml kong kong/kong
```

Once the deployments are `Ready` you'll need to capture the IP address and port of the container running the Kong Gateway:

```console
$ export PROXY_ENDPOINT="$(docker inspect kic-gateway-testing-control-plane -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'):32080"
```

If things have worked correctly, the contents of that variable should be something like this:

```console
$ echo $PROXY_ENDPOINT
172.18.0.2:32080
```

And you should be able to communicate with the Kong Gateway using it:

```console
$ curl http://${PROXY_ENDPOINT}/
{"message":"no Route matched with those values"}
```

The above 404 is expected when there's nothing yet configured.

[kind]:https://kind.sigs.k8s.io
[kind-install]:https://kind.sigs.k8s.io/docs/user/quick-start/#installation
[kubectl]:https://kubernetes.io/docs/reference/kubectl/overview/
[helm]:https://helm.sh/

#### Enabling Gateway Feature

By default `Gateway` support in the KIC is disabled by way of [feature gates][feature-gates].

To enable the `Gateway` feature gate you'll need the Gateway APIs CRDs installed on the cluster and the `--feature-gates=Gateway=true` argument needs to be added to the KIC's `ingress-controller` container.

First you'll need to deploy the Gateway CRDs:

```console
$ kubectl kustomize https://github.com/kubernetes-sigs/gateway-api.git/config/crd?ref=master | kubectl apply -f -
```

Now to enable the feature, configure this as an argument for the ingress controller in your Helm `values.yaml`:

```yaml
ingressController:
  args:
    - --feature-gates=Gateway=true
```

Complete the enablement by upgrading the Helm release with the new options:

```console
$ helm upgrade --namespace kong -f values.yaml kong kong/kong
```

If everything is working the pod logs should show entries like the following:

```console
time="2022-01-11T21:04:30Z" level=info msg="found configuration option for gated feature" enabled=true feature=Gateway logger=setup
time="2022-01-11T21:04:30Z" level=info msg="Starting Controller" logger=controllers.Gateway.controller.gateway-controller
```

#### Deploying a Gateway resource

In order for `HTTPRoute` objects to route traffic we will need to create a `Gateway` object.

Note that currently we only support a deployment mode called "unmanaged gateway mode": in this mode the `Gateway` resource is a reference to the Kubernetes `Service` where an existing Gateway is running. In this mode the gateway controller will automatically manage the specification of the `Gateway` resource and keep it updated according to that reference `Service`. Later iterations of gateway support may add managed gateways where the act of creating the resource results in underlying deployments and pods being provisioned according to the spec.

Deploy a `GatewayClass` which will link the KIC's gateway controller to `Gateway` objects:

```yaml
kind: GatewayClass
apiVersion: gateway.networking.k8s.io/v1alpha2
metadata:
  name: kong
spec:
  controllerName: konghq.com/kic-gateway-controller
```

Next create the `Gateway` resource linked to that `GatewayClass`:

```yaml
kind: Gateway
apiVersion: gateway.networking.k8s.io/v1alpha2
metadata:
  annotations:
    konghq.com/gateway-unmanaged: "true"
  name: kong
spec:
  gatewayClassName: kong
  listeners:
  - name: http
    protocol: HTTP
    port: 80
```

Given some time you should see the Gateway become `READY`, e.g.:

```console
$ kubectl get gateways
NAME   CLASS   ADDRESS        READY   AGE
kong   kong    10.96.220.25   True    10s
```

Once the `Gateway` is ready it can start accepting `HTTPRoutes`.

#### Routing HTTP Traffic

This example will create a basic `HTTPRoute` to allow http traffic to a simple webserver.

First we'll create a `Deployment` with our backend webserver:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpbin
  labels:
    app: httpbin
spec:
  selector:
    matchLabels:
      app: httpbin
  template:
    metadata:
      labels:
        app: httpbin
    spec:
      containers:
      - name: httpbin
        image: kennethreitz/httpbin
        ports:
        - containerPort: 80
```

Expose the `Deployment` with a `Service`:

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: httpbin
  name: httpbin
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: httpbin
  type: ClusterIP
```

And ultimately route traffic to the service via `HTTPRoute`:

```yaml
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: HTTPRoute
metadata:
  name: httpbin
spec:
  parentRefs:
  - name: kong
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /httpbin
    backendRefs:
    - name: httpbin
      port: 80
```

Once the pods for the `httpbin` deployment are ready, you should be able to start accessing the webserver via Kong:

```console
$ curl -w 'STATUS: %{http_code}\n' http://${PROXY_ENDPOINT}/httpbin/status/200
STATUS: 200
```

Where `${PROXY_ENDPOINT}` is the IP address and port of your KIC proxy service (e.g. `172.18.0.2:32080` if you used the local cluster setup option above).

[feature-gates]:https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
