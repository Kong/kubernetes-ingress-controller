# Kong Ingress Controller annotations

Kong Ingress Controller supports the following annotations on various resources:

## Ingress resource

Following annotations are supported on Ingress resources:

| Annotation name | Description |
|-----------------|-------------|
| REQUIRED [`kubernetes.io/ingress.class`](#kubernetesioingressclass) | Restrict the Ingress rules that Kong should satisfy |
| [`konghq.com/plugins`](#konghqcomplugins) | Run plugins for specific Ingress. |
| [`konghq.com/protocols`](#konghqcomprotocols) | Set protocols to handle for each Ingress resource. |
| [`konghq.com/preserve-host`](#konghqcompreserve-host) | Pass the `host` header as is to the upstream service. |
| [`konghq.com/strip-path`](#konghqcomstrip-path) | Strip the path defined in Ingress resource and then forward the request to the upstream service. |
| [`konghq.com/https-redirect-status-code`](#konghqcomhttps-redirect-status-code) | Set the HTTPS redirect status code to use when an HTTP request is recieved. |
| [`konghq.com/regex-priority`](#konghqcomregex-priority) | Set the route's regex priority. |
| [`konghq.com/methods`](#konghqcommethods) | Set methods matched by this Ingress. |
| [`konghq.com/override`](#konghqcomoverride) | Control other routing attributes via `KongIngress` resource. |

`kubernetes.io/ingress.class` is normally required, and its value should match
the value of the `--ingress-class` controller argument ("kong" by default).

Setting the `--process-classless-ingress-v1beta1` controller flag removes that requirement:
when enabled, the controller will process Ingresses with no
`kubernetes.io/ingress.class` annotation. Recommended best practice is to set
the annotation and leave this flag disabled; the flag is intended for
older configurations, as controller versions prior to 0.10 processed classless
Ingress resources by default.

## Service resource

Following annotations are supported on Service resources:

| Annotation name | Description |
|-----------------|-------------|
| [`konghq.com/plugins`](#konghqcomplugins) | Run plugins for a specific Service |
| [`konghq.com/protocol`](#konghqcomprotocol) | Set protocol Kong should use to talk to a Kubernetes service |
| [`konghq.com/path`](#konghqcompath) | HTTP Path that is always prepended to each request that is forwarded to a Kubernetes service |
| [`konghq.com/client-cert`](#konghqcomclient-cert) | Client certificate and key pair Kong should use to authenticate itself to a specific Kubernetes service |
| [`konghq.com/host-header`](#konghqcomhost-header) | Set the value sent in the `Host` header when proxying requests upstream |
| [`konghq.com/override`](#konghqcomoverride) | Fine grained routing and load-balancing |
| [`ingress.kubernetes.io/service-upstream`](#ingresskubernetesioservice-upstream) | Offload load-balancing to kube-proxy or sidecar |

## KongConsumer resource

Following annotaitons are supported on KongConsumer resources:

| Annotation name | Description |
|-----------------|-------------|
| REQUIRED [`kubernetes.io/ingress.class`](#kubernetesioingressclass) | Restrict the KongConsumers that a controller should satisfy |
| [`konghq.com/plugins`](#konghqcomplugins) | Run plugins for a specific consumer |

`kubernetes.io/ingress.class` is normally required, and its value should match
the value of the `--ingress-class` controller argument ("kong" by default).

Setting the `--process-classless-kong-consumer` controller flag removes that requirement:
when enabled, the controller will process KongConsumers with no
`kubernetes.io/ingress.class` annotation. Recommended best practice is to set
the annotation and leave this flag disabled; the flag is primarily intended for
older configurations, as controller versions prior to 0.10 processed classless
KongConsumer resources by default.

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

The following resources _require_ this annotation by default:

- Ingress
- KongConsumer
- TCPIngress
- KongClusterPlugin
- Secret resources with the `ca-cert` label 

You can optionally allow Ingress or KongConsumer resources with no class
annotation (by setting the `--process-classless-ingress-v1beta1` or
`--process-classless-kong-consumer` flags, respectively), though recommended
best practice is to leave these flags disabled: the flags are primarily
intended for compatibility with configuration created before this requirement
was introduced in controller 0.10.

If you allow classless resources, you must take care when using multiple
controller instances in a single cluster: only one controller instance should
enable these flags to avoid different controller instances fighting over
classless resources, which will result in unexpected and unknown behavior.

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

### `konghq.com/regex-priority`

> Available since controller 0.9

Sets the `regex_priority` setting to this value on the Kong route associated
with the Ingress resource. This controls the [matching evaluation
order](https://docs.konghq.com/latest/proxy/#evaluation-order) for regex-based
routes. It accepts any integer value. Routes are evaluated in order of highest
priority to lowest.

Sample usage:

```yaml
konghq.com/regex-priority: "10"
```

Please note the quotes (`"`) around the integer value.

### `konghq.com/methods`

> Available since controller 0.9

Sets the `methods` setting on the Kong route associated with the Ingress
resource. This controls which request methods will match the route. Any
uppercase alpha ASCII string is accepted, though most users will use only
[standard methods](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods).

Sample usage:

```yaml
konghq.com/methods: "GET,POST"
```

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

### `konghq.com/host-header`

> Available since controller 0.9

Sets the `host_header` setting on the Kong upstream created to represent a
Kubernetes Service. By default, Kong upstreams set `Host` to the hostname or IP
address of an individual target (the Pod IP for controller-managed
configuration). This annotation overrides the default behavior and sends
the annotation value as the `Host` header value.

If `konghq.com/preserve-host: true` is present on an Ingress (or
`route.preserve_host: true` is present in a linked KongIngress), it will take
precedence over this annotation, and requests to the application will use the
hostname in the Ingress rule.

Sample usage:

```yaml
konghq.com/host-header: "test.example.com"
```

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
