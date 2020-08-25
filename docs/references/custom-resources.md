# Custom Resource Definitions

The Ingress Controller can configure Kong specific features
using several [Custom Resource Definitions(CRDs)][k8s-crd].

Following CRDs enables users to declaratively configure all aspects of Kong:

- [**KongPlugin**](#kongplugin): This resource corresponds to
  the [Plugin][kong-plugin] entity in Kong.
- [**KongIngress**](#kongingress): This resource provides fine-grained control
  over all aspects of proxy behaviour like routing, load-balancing,
  and health checking. It serves as an "extension" to the Ingress resources
  in Kubernetes.
- [**KongConsumer**](#kongconsumer):
  This resource maps to the [Consumer][kong-consumer] entity in Kong.
- [**TCPIngress**](#tcpingress):
  This resource can configure TCP-based routing in Kong for non-HTTP
  services running inside Kubernetes.
- [**KongCredential (Deprecated)**](#kongcredential-deprecated):
  This resource maps to
  a credential (key-auth, basic-auth, jwt, hmac-auth) that is associated with
  a specific KongConsumer.

## KongPlugin

This resource provides an API to configure plugins inside Kong using
Kubernetes-style resources.

Please see the [concept](../concepts/custom-resources.md#KongPlugin)
document for how the resource should be used.

The following snippet shows the properties available in KongPlugin resource:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: <object name>
  namespace: <object namespace>
disabled: <boolean>  # optionally disable the plugin in Kong
config:              # configuration for the plugin
    key: value
configFrom:
    secretKeyRef:
       name: <Secret name>
       key: <Secret key>
plugin: <name-of-plugin> # like key-auth, rate-limiting etc
```

- `config` contains a list of `key` and `value`
  required to configure the plugin.
  All configuration values specific to the type of plugin go in here.
  Please read the documentation of the plugin being configured to set values
  in here. For any plugin in Kong, anything that goes in the `config` JSON
  key in the Admin API request, goes into the  `config` YAML key in this resource.
  Please use a valid JSON to YAML convertor and place the content under the
  `config` key in the YAML above.
- `configFrom` contains a reference to a Secret and key, where the key contains
  a complete JSON or YAML configuration. This should be used when the plugin
  configuration contains sensitive information, such as AWS credentials in the
  Lambda plugin or the client secret in the OIDC plugin. Only one of `config`
  or `configFrom` may be used in a KongPlugin, not both at once.
- `plugin` field determines the name of the plugin in Kong.
  This field was introduced in Kong Ingress Controller 0.2.0.

**Please note:** validation of the configuration fields is left to the user
by default. It is advised to setup and use the admission validating controller
to catch user errors.

The plugins can be associated with Ingress
or Service object in Kubernetes using `plugins.konghq.com` annotation.

### Examples

#### Applying a plugin to a service

Given the following plugin:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: request-id
config:
  header_name: my-request-id
  echo_downstream: true
plugin: correlation-id
```

It can be applied to a service by annotating like:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
  labels:
     app: myapp-service
  annotations:
     plugins.konghq.com: request-id
spec:
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
    name: myapp-service
  selector:
    app: myapp-service
```

#### Applying a plugin to an ingress

The KongPlugin above can be applied to a specific ingress (route or routes):

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-example-com
  annotations:
    plugins.konghq.com: request-id
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /bar
        backend:
          serviceName: echo
          servicePort: 80
```

A plugin can also be applied to a specific KongConsumer by adding
`plugins.konghq.com` annotation to the KongConsumer resource.

Please follow the
[Using the KongPlugin resource](../guides/using-kongplugin-resource.md)
guide for details on how to use this resource.

#### Applying a plugin with a secret configuration

The plugin above can be modified to store its configuration in a secret:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: request-id
configFrom:
  secretKeyRef:
    name: plugin-conf-secret
    key: request-id
plugin: correlation-id
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: plugin-conf-secret
stringData:
  request-id: |
    header_name: my-request-id
    echo_downstream: true
type: Opaque
```

## KongClusterPlugin

A `KongClusterPlugin` is same as `KongPlugin` resource. The only differences
are that it is a Kubernetes cluster-level resource instead of a namespaced
resource, and can be applied as a global plugin using labels.

Please consult the [KongPlugin](#kongplugin) section for details.

*Example:*

KongClusterPlugin example:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongClusterPlugin
metadata:
  name: request-id
  labels:
    global: "true"   # optional, if set, then the plugin will be executed
                     # for every request that Kong proxies
                     # please note the quotes around true
config:
  header_name: my-request-id
configFrom:
    secretKeyRef:
       name: <Secret name>
       key: <Secret key>
       namespace: <Secret namespace>
plugin: correlation-id
```

As with KongPlugin, only one of `config` or `configFrom` can be used.

Setting the label `global` to `"true"` will apply the plugin globally in Kong,
meaning it will be executed for every request that is proxied via Kong.

## KongIngress

Ingress resource spec in Kubernetes can define routing policies
based on HTTP Host header and paths.
While this is sufficient in most cases,
sometimes, users may want more control over routing at the Ingress level.
`KongIngress` serves as an "extension" to Ingress resource.
It is not meant as a replacement to the
`Ingress` resource in Kubernetes.

Please read the [concept](../concepts/custom-resources.md#kongingress)
document for why this resource exists and how it relates to the existing
Ingress resource.

Using `KongIngress`, all properties of [Upstream][kong-upstream],
[Service][kong-service] and
[Route][kong-route] entities in Kong related to an Ingress resource
can be modified.

Once a `KongIngress` resource is created, it needs to be associated with
an Ingress or Service resource using the following annotation:

```yaml
configuration.konghq.com: kong-ingress-resource-name
```

Specifically,

- To override any properties related to health-checking, load-balancing,
  or details specific to a service, add the annotation to the Kubernetes
  Service that is being exposed via the Ingress API.
- To override routing configuration (like protocol or method based routing),
  add the annotation to the Ingress resource.

Please follow the
[Using the KongIngress resource](../guides/using-kongingress-resource.md)
guide for details on how to use this resource.

For reference, the following is a complete spec for KongIngress:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
  name: configuration-demo
upstream:
  slots: 10
  hash_on: none
  hash_fallback: none
  healthchecks:
    threshold: 25
    active:
      concurrency: 10
      healthy:
        http_statuses:
        - 200
        - 302
        interval: 0
        successes: 0
      http_path: "/"
      timeout: 1
      unhealthy:
        http_failures: 0
        http_statuses:
        - 429
        interval: 0
        tcp_failures: 0
        timeouts: 0
    passive:
      healthy:
        http_statuses:
        - 200
        successes: 0
      unhealthy:
        http_failures: 0
        http_statuses:
        - 429
        - 503
        tcp_failures: 0
        timeouts: 0
proxy:
  protocol: http
  path: /
  connect_timeout: 10000
  retries: 10
  read_timeout: 10000
  write_timeout: 10000
route:
  methods:
  - POST
  - GET
  regex_priority: 0
  strip_path: false
  preserve_host: true
  protocols:
  - http
  - https
```

## TCPIngress

The Ingress resource in Kubernetes is HTTP-only.
This custom resource is modeled similar to the Ingress resource but for
TCP and TLS SNI based routing purposes:

```yaml
apiVersion: configuration.konghq.com/v1beta1
kind: TCPIngress
metadata:
  name: <object name>
  namespace: <object namespace>
spec:
  rules:
  - host: <SNI, optional>
    port: <port on which to expose this service, required>
    backend:
      serviceName: <name of the kubernetes service, required>
      servicePort: <port number to forward on the service, required>
```

If `host` is not specified, then port-based TCP routing is performed. Kong
doesn't care about the content of TCP stream in this case.

If `host` is specified, then Kong expects the TCP stream to be TLS-encrypted
and Kong will terminate the TLS session based on the SNI.
Also note that, the port in this case should be configured with `ssl` parameter
in Kong.

## KongConsumer

This custom resource configures a consumer in Kong:

The following snippet shows the field available in the resource:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: <object name>
  namespace: <object namespace>
  annotations:
    kubernetes.io/ingress.class: kong
username: <user name>
custom_id: <custom ID>
```

An example:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: consumer-team-x
  annotations:
    kubernetes.io/ingress.class: kong
username: team-X
```

When this resource is created, a corresponding consumer entity will be
created in Kong.

Consumers' `username` and `custom_id` values must be unique across the Kong
cluster. While KongConsumers exist in a specific Kubernetes namespace,
KongConsumers from all namespaces are combined into a single Kong
configuration, and no KongConsumers with the same `kubernetes.io/ingress.class`
may share the same `username` or `custom_id` value.

## KongCredential (Deprecated)

This custom resource can be used to configure a consumer specific
entities in Kong.
The resource reference the KongConsumer resource via the `consumerRef` key.

The validation of the config object is left up to the user.

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongCredential
metadata:
  name: credential-team-x
consumerRef: consumer-team-x
type: key-auth
config:
  key: 62eb165c070a41d5c1b58d9d3d725ca1
```

The following credential types can be provisioned using the KongCredential
resource:

- `key-auth` for [Key authentication](https://docs.konghq.com/plugins/key-authentication/)
- `basic-auth` for [Basic authenticaiton](https://docs.konghq.com/plugins/basic-authentication/)
- `hmac-auth` for [HMAC authentication](http://docs.konghq.com/plugins/hmac-authentication/)
- `jwt` for [JWT based authentication](http://docs.konghq.com/plugins/jwt/)
- `oauth2` for [Oauth2 Client credentials](https://docs.konghq.com/hub/kong-inc/oauth2/)
- `acl` for [ACL group associations](https://docs.konghq.com/hub/kong-inc/acl/)

Please ensure that all fields related to the credential in Kong
are present in the definition of KongCredential's `config` section.

Please refer to the
[using the Kong Consumer and Credential resource](../guides/using-consumer-credential-resource.md)
guide for details on how to use this resource.

[k8s-crd]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/
[kong-consumer]: https://getkong.org/docs/latest/admin-api/#consumer-object
[kong-plugin]: https://getkong.org/docs/latest/admin-api/#plugin-object
[kong-upstream]: https://getkong.org/docs/latest/admin-api/#upstream-objects
[kong-service]: https://getkong.org/docs/latest/admin-api/#service-object
[kong-route]: https://getkong.org/docs/latest/admin-api/#route-object
