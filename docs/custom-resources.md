# Custom Resource Definitions

The Ingress Controller can configure Kong specific features
using several [Custom Resource Definitions(CRDs)][k8s-crd].

Following CRDs enables users to declaratively configure all aspects of Kong:

- [**KongPlugin**](#kongplugin): These resources correspond to [Plugin][kong-plugin]
  entities in Kong.
- [**KongIngress**](#kongingress): These resources can control routing, load-balancing,
  health checking properties in Kong.  
  It works with the Ingress resources in Kubernetes.
- [**KongConsumer**](#kongconsumer):
  These resources map to [Consumer][kong-consumer] entities in Kong.
- [**KongCredential**](#kongcredential): These resources map to
  credentials (key-auth, basic-auth, etc) that belong to consumers.

## KongPlugin

This resource allows the configuration of
Kong plugins in the same way we add plugins using the
[admin API][kong-add-plugin]

Following is an example no how to define a `KongPlugin` resource:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: <object name>
  namespace: <object namespace>
  labels:
    global: "true" # optional, please note the quotes around true
                  # configures the plugin Globally in Kong
consumerRef: <name of an existing consumer> # optional
                                            # applies the plugin
                                            # in on specific route and consumer
disabled: <boolean>  # optionally disable the plugin in Kong
config:
    key: value
plugin: <name-of-plugin> # like key-auth, rate-limiting etc
```

- `consumerRef`, an optional field,
  implies the plugin will be used for a particular consumer only.
  The value of the field must reference an existing consumer
  in the same namespace.
  If specified, the plugin will execute for the specific consumer only.
- `config` contains a list of `key` and `value`
  required to configure the plugin.
  All configuration values specific to the type of plugin go in here.
  Please read the documentation of the plugin being configured to set values
  in here.
- `disabled` if set to true, disables the plugin in Kong (but not delete it).
- `plugin` field determines the name of the plugin in Kong.
  This field was introduced in Kong Ingress Controller 0.2.0.
- Setting a label `global` to `"true"` will result in the plugin being
  applied globally in Kong, meaning it will be executed for every
  request that is proxied via Kong.

**Please note:** validation of the configuration fields is left to the user.
Setting invalid fields will result in errors in the Ingress Controller.
This behavior is set to improve in the future.

The plugins can be associated with Ingress
or Service object in Kubernetes using `plugins.konghq.com` annotation.

*Example:*

Given the following plugin:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: http-svc-consumer-ratelimiting
  namespace: default #this should match the namespace of the route or service you're adding it too.
config:
  key: value
plugin: my-plugin
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
     plugins.konghq.com: http-svc-consumer-ratelimiting

spec:
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
    name: myapp-service
  selector:
    app: myapp-service
```

It can be applied to a specific ingress (route or routes) like:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
   name: myapp-ingress
   annotations:
      plugins.konghq.com: http-svc-consumer-ratelimiting
spec:
   rules:
     - host: my.host.com
       http:
         paths:
           - path: /myendpoint
             backend:
               serviceName: myapp-service
               servicePort: 80
```

## KongIngress

Ingress resource spec in Kubernetes can define routing policies
based on HTTP Host header and paths.  
While this is sufficient in most cases,
sometimes, users may want more control over routing at the Ingress level.
`KongIngress` works in conjunction with existing Ingress resource
and extends it. It is not meant as a replacement to the
`Ingress` resource in Kubernetes.
Using `KongIngress`, all properties of [Upstream][kong-upstream],
[Service][kong-service] and
[Route][kong-route] entities in Kong related to an Ingress resource
can be modified.

Once a `KongIngress` resource is created, it can be associated with
`Ingress` resource in two ways:

- Create a `KongIngress` object in the same namespace as that of the
  Ingress rule using the same name.
  This avoids a need of additional annotation in Ingress resource.
  On the other hand, this approach requires a `KongIngress`
  resource per Ingress, which becomes hard to maintain with multiple Ingresses.

- Create a `KongIngress` resource and then using the annotation
  `configuration.konghq.com: <KongIngress-resource-name>`,
  associate it with one or more Ingress resources.  
  This approach allows you to reuse the same `KongIngress`.

*Note:* Is not required to define the complete object,
one can define only one of the `upstream`, `proxy` or `route` sections

Following is a complete spec for KongIngress:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
  name: configuration-demo
upstream:
  hash_on: none
  hash_fallback: none
  healthchecks:
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
    slots: 10
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

## KongConsumer

This custom resource configures consumers in Kong:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: <object name>
  namespace: <object namespace>
username: <user name>
custom_id: <custom ID>
```

An example:

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: consumer-team-x
username: team-X
custom_id: my_team_x # optional and not recommended, please use `username`
```

## KongCredential

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

[k8s-crd]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/
[kong-consumer]: https://getkong.org/docs/latest/admin-api/#consumer-object
[kong-plugin]: https://getkong.org/docs/latest/admin-api/#plugin-object
[kubectl-doc]: https://kubernetes.io/docs/reference/kubectl/overview/
[kong-add-plugin]: https://getkong.org/docs/latest/admin-api/#add-plugin
[kong-upstream]: https://getkong.org/docs/latest/admin-api/#upstream-objects
[kong-service]: https://getkong.org/docs/latest/admin-api/#service-object
[kong-route]: https://getkong.org/docs/latest/admin-api/#route-object
