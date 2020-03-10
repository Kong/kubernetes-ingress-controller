# Kong Ingress Controller annotations

Kong Ingress Controller supports the following annotations on various resources:

## Ingress resource

Following annotations are supported on Ingress resources:

| Annotation name | Description |
|-----------------|-------------|
| [`kubernetes.io/ingress.class`](#kubernetesioingressclass) | Restrict the Ingress rules that Kong should satisfy |
| [`konghq.com/plugins`](#konghqcom/plugins) | Run plugins for specific Ingress. |
| [`konghq.com/protocols`](#konghqcom/protocols) | Set protocols to handle for each Ingress resource. |
| [`konghq.com/preserve-host`](#konghqcom/preserve-host) | Pass the `host` header as is to the upstream service. |
| [`konghq.com/strip-path`](#konghqcom/strip-path) | Strip the path defined in Ingress resource and then forward the request to the upstream service. |
| [`konghq.com/https-redirect-status-code`](#konghqcom/https-redirect-status-code) | Set the HTTPS redirect status code to use when an HTTP request is recieved. |
| [`konghq.com/override`](#konghqcom/override) | Control other routing attributes via `KongIngress` resource. |
| DEPRECATED [`plugins.konghq.com`](#pluginskonghqcom) | Please use [`konghq.com/plugins`](#konghqcom/plugins) |
| DEPRECATED [`configuration.konghq.com`](#configurationkonghqcom) | Please use [`konghq.com/override`](#konghqcomoverride) |
| DEPRECATED [`configuration.konghq.com/protocols`](#configurationkonghqcom/protocols) | Please use [`konghq.com/protocols`](#konghqcom/protocols) |

## Service resource

Following annotations are supported on Service resources:

| Annotation name | Description |
|-----------------|-------------|
| [`konghq.com/plugins`](#konghqcom/plugins) | Run plugins for a specific Service |
| [`konghq.com/protocol`](#konghqcom/protocol) | Set protocol Kong should use to talk to a Kubernetes service |
| [`konghq.com/path`](#konghqcom/path) | HTTP Path that is always prepended to each request that is forwarded to a Kubernetes service |
| [`konghq.com/client-cert`](#konghqcom/client-cert) | Client certificate and key pair Kong should use to authenticate itself to a specific Kubernetes service |
| [`konghq.com/override`](#konghqcomoverride) | Fine grained routing and load-balancing |
| [`ingress.kubernetes.io/service-upstream`](#ingresskubernetesioservice-upstream) | Offload load-balancing to kube-proxy or sidecar |
| DEPRECATED [`plugins.konghq.com`](#pluginskonghqcom) | Please use [`konghq.com/plugins`](#konghqcom/plugins) |
| DEPRECATED [`configuration.konghq.com`](#configurationkonghqcom) | Please use [`konghq.com/override`](#konghqcomoverride) |
| DEPRECATED [`configuration.konghq.com/protocol`](#configurationkonghqcom/protocol) | Please use [`konghq.com/protocol`](#konghqcom/protocol) |
| DEPRECATED [`configuration.konghq.com/client-cert`](#configurationkonghqcom/client-cert) | Please use [`konghq.com/client-cert`](#konghqcom/client-cert) |

## KongConsumer resource

Following annotaitons are supported on KongConsumer resources:

| Annotation name | Description |
|-----------------|-------------|
| [`kubernetes.io/ingress.class`](#kubernetesioingressclass) | Restrict the KongConsumers that a controller should satisfy |
| [`konghq.com/plugins`](#konghqcom/plugins) | Run plugins for a specific consumer |
| DEPRECATED [`plugins.konghq.com`](#pluginskonghqcom) | Please use [`konghq.com/plugins`](#konghqcom/plugins) |

## Annotations

### `kubernetes.io/ingress.class`

If you have multiple Ingress controllers in a single cluster,
you can pick one by specifying the `ingress.class` annotation.
Following is an example of
creating an Ingress with an annotation:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: test-1
  annotations:
    kubernetes.io/ingress.class: "gce"
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /test1
        backend:
          serviceName: echo
          servicePort: 80
```

This will target the GCE controller, forcing Kong Ingress Controller to ignore it.

On the other hand, an annotation such as

```yaml
metadata:
  name: test-1
  annotations:
    kubernetes.io/ingress.class: "kong"
```

will target Kong Ingress controller, forcing the GCE controller
to ignore it.

> Deploying multiple ingress controller and not specifying the
annotation will cause both controllers fighting to satisfy the Ingress
and will lead to unknown behaviour.

The ingress class used by Kong Ingress Controller to filter Ingress
resources can be changed using the `CONTROLLER_INGRESS_CLASS`
environment variable.

```yaml
spec:
  template:
     spec:
       containers:
         - name: kong-ingress-internal-controller
           env:
           - name: CONTROLLER_INGRESS_CLASS
             value: kong-internal
```

#### Multiple unrelated Kong Ingress Controllers

In some deployments, one might use multiple Kong Ingress Controller
in the same Kubernetes cluster
(e.g. one which serves public traffic, one which serves "internal" traffic).
For such deployments, please ensure that in addition to different
`ingress-class`, the `--election-id` is also different.

In such deployments, `kubernetes.io/ingress.class` annotation can be used on the
following custom resources as well:

- KongPlugin: To configure (global) plugins only in one of the Kong clusters.
- KongConsumer: To create different consumers in different Kong clusters.
- KongCredential: To create associated credentials for consumers.

### `konghq.com/plugins`

> Available since controller 0.8

Kong's power comes from its plugin architecture, where plugins can modify
the request and response or impose certain policies on the requests as they
are proxied to your service.

With Kong Ingress Controller, plugins can be configured by creating `KongPlugin`
Custom Resources and then associating them with an Ingress, Service,
KongConsumer or a combination of those.

Following is an example of how to use the annotation:

```yaml
konghq.com/plugins: high-rate-limit, docs-site-cors
```

Here, `high-rate-limit` and `docs-site-cors`
are the names of the KongPlugin resources which
should be to be applied to the Ingress rules defined in the
Ingress resource on which the annotation is being applied.

This annotation can also be applied to a Service resource in Kubernetes, which
will result in the plugin being executed at Service-level in Kong,
meaning the plugin will be
executed for every request that is proxied, no matter which Route it came from.

This annotation can also be applied to a KongConsumer resource,
which results in plugin being executed whenever the specific consumer
is accessing any of the defined APIs.

Finally, this annotation can also be applied on a combination of the
following resources:
- **Ingress and KongConsumer**  
  If an Ingress resource and a KongConsumer resource share a plugin in the
  `konghq.com/plugins` annotation then the plugin will be created for the
  combination of those to resources in Kong.
- **Service and KongConsumer**  
  Same as the above case, if you would like to give a specific consumer or
  client of your service some special treatment, you can do so by applying
  the same annotation to both of the resources.

Please follow the
[Using the KongPlugin resource](../guides/using-kongplugin-resource.md)
guide for details on how this annotation can be used.


### `konghq.com/path`

> Available since controller 0.8

This annotation can be used on a Service resource only.
This annotation can be used to prepend an HTTP path of a request,
before the request is forwarded.

For example, if the annotation `konghq.com/path: "/baz"` is applied to a
Kubernetes Service `billings`, then any request that is routed to the
`billings` service will be prepended with `/baz` HTTP path. If the
request contains `/foo/something` as the path, then the service will
receive an HTTP request with path set as `/baz/foo/something`.

### `konghq.com/strip-path`

> Available since controller 0.8

This annotation can be applied to an Ingress resource and can take two values:
- `"true"`: If set to true, the part of the path specified in the Ingress rule
  will be stripped out before the request is sent to the service.
  For example, if the Ingress rule has a path of `/foo` and the HTTP request
  that matches the Ingress rule has the path `/foo/bar/something`, then
  the request sent to the Kubernetes service will have the path
  `/bar/something`.
- `"false"`: If set to false, no path manipulation is performed.

All other values are ignored.
Please note the quotes (`"`) around the boolean value.

Sample usage:

```yaml
konghq.com/strip-path: "true"
```

### `konghq.com/preserve-host`

> Available since controller 0.8

This annotation can be applied to an Ingress resource and can take two values:
- `"true"`: If set to true, the `host` header of the request will be sent
  as is to the Service in Kubernetes.
- `"false"`: If set to false, the `host` header of the request is not preserved.

Please note the quotes (`"`) around the boolean value.

Sample usage:

```yaml
konghq.com/preserve-host: "true"
```

### `konghq.com/https-redirect-status-code`

> Available since controller 0.8

By default, Kong sends HTTP Status Code 426 for requests
that need to be redirected to HTTPS.
This can be changed using this annotations.
Acceptable values are:
- 301
- 302
- 307
- 308
- 426

Any other value will be ignored.

Sample usage:

```yaml
konghq.com/https-redirect-status-code: "301"
```

Please note the quotes (`"`) around the integer value.

### `konghq.com/override`

> Available since controller 0.8

This annotation can associate a KongIngress resource with
an Ingress or a Service resource.
It serves as a way to bridge the gap between a sparse Ingress API in Kubernetes
with fine-grained controlled using the properties of Service, Route
and Upstream entities in Kong.

Please follow the
[Using the KongIngress resource](../guides/using-kongingress-resource.md)
guide for details on how to use this annotation.

### `konghq.com/protocol`

> Available since controller 0.8

This annotation can be set on a Kubernetes Service resource and indicates
the protocol that should be used by Kong to communicate with the service.
In other words, the protocol is used for communication between a
[Kong Service](https://docs.konghq.com/latest/admin-api/#service-object) and 
a Kubernetes Service, internally in the Kubernetes cluster.

Accepted values are:
- `http`
- `https`
- `grpc`
- `grpcs`
- `tcp`
- `tls`

### `konghq.com/protocols`

> Available since controller 0.8

This annotation sets the list of acceptable protocols for the all the rules
defined in the Ingress resource.
The protocols are used for communication between the
Kong and the external client/user of the Service.

You usually want to set this annotation for the following two use-cases:
- You want to redirect HTTP traffic to HTTPS, in which case you will use
  `konghq.com/protocols: "https"`
- You want to define gRPC routing, in which case you should use
  `konghq.com/protocols: "grpc,grpcs"`

### `konghq.com/client-cert`

> Available since controller 0.8

This annotation sets the certificate and key-pair Kong should use to
authenticate itself against the upstream service, if the upstream service
is performing mutual-TLS (mTLS) authentication.

The value of this annotation should be the name of the Kubernetes TLS Secret
resource which contains the TLS cert and key pair.

Under the hood, the controller creates a Certificate in Kong and then
sets the
[`service.client_certificate`](https://docs.konghq.com/latest/admin-api/#service-object)
for the service.

### `ingress.kubernetes.io/service-upstream`

By default, Kong Ingress Controller distributes traffic amongst all the Pods
of a Kubernetes `Service` by forwarding the requests directly to
Pod IP addresses. One can choose the load-balancing strategy to use
by specifying a KongIngress resource.

However, in some use-cases, the load-balancing should be left up
to `kube-proxy`, or a sidecar component in the case of Service Mesh deployments.

Setting this annotation to a Service resource in Kubernetes will configure
Kong Ingress Controller to directly forward
the traffic outbound for this Service
to the IP address of the service (usually the ClusterIP).

`kube-proxy` can then decide how it wants to handle the request and route the
traffic accordingly. If a sidecar intercepts the traffic from the controller,
it can also route traffic as it sees fit in this case.

Following is an example snippet you can use to configure this annotation
on a `Service` resource in Kubernetes, (please note the quotes around `true`):

```yaml
annotations:
  ingress.kubernetes.io/service-upstream: "true"
```

You need Kong Ingress Controller >= 0.6 for this annotation.

### `plugins.konghq.com`

> DEPRECATED in Controller 0.8

Please instead use [`konghq.com/plugins`](#konghqcom/plugins).

### `configuration.konghq.com`

> DEPRECATED in Controller 0.8

Please instead use [`konghq.com/override`](#konghqcomoverride).

### `configuration.konghq.com/protocol`

> DEPRECATED in Controller 0.8

Please instead use [`konghq.com/protocol`](#konghqcom/protocol).

### `configuration.konghq.com/protocols`

> DEPRECATED in Controller 0.8

Please instead use [`konghq.com/protocols`](#konghqcom/protocols).

### `configuration.konghq.com/client-cert`

> DEPRECATED in Controller 0.8

Please instead use [`konghq.com/client-cert`](#konghqcom/client-cert).
