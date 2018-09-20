# Custom Resource Definitions

Kong relies on several [Custom Resource Definitions][0] to declare
additional information to Ingress rules and synchronize configuration with the Kong admin API.

The custom resource names are:

- **KongConsumer**: These resources map to [Consumer][1] entities in Kong.
- **KongCredential**: These resources map to credentials (key-auth, basic-auth, etc) that belong to consumers.
- **KongPlugin**: These resources belong to [Plugin][2] entities in Kong.

### KongPlugin

This resource allows the configuration of Kong plugins in the same way we [add plugins using the admin API][4]

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: <object name>
  namespace: <object namespace>
  labels:
    global: "true" # optional, please note the quotes around true
consumerRef: <optional, name of an existing consumer> # optional
disabled: <boolean>  # optional
config:
    key: value
plugin: <name-of-plugin>
```

- `consumerRef`, an optional field, implies the plugin will be used for a particular consumer only.
  The value of the field must reference an existing consumer in the same namespace.
  If specified, the plugin will execute for the specific consumer only.
- `config` contains a list of `key` and `value` required to configure the plugin.
  All configuration values specific to the type of plugin go in here.
  Please read the documentation of the plugin being configured to set values
  in here.
- `disabled` if set to true, disables the plugin in Kong (but not delete it).
- `plugin` field determines the name of the plugin in Kong.
  This field was introduced in Kong Ingress Controller 0.2.0.
- Setting a label `global` to `"true"` will result in the plugin being
  applied globally in Kong, meaning it will be executed for every
  request that is proxied via Kong.

**Important:** validation of the configuration fields is left to the user.
Setting invalid fields will result in errors in the Ingress Controller.
This behavior is set to improve in future.

The plugins can be associated with Ingress resources using `plugins.konghq.com` annotation.

*Example:*

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: http-svc-consumer-ratelimiting
  namespace: default
config:
  key: value
plugin: my-plugin
```

### KongConsumer

*Definition:*

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: <object name>
  namespace: <object namespace>
username: <user name>
custom_id: <custom ID>
```

This resource allows configuring Consumers in Kong

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: consumer-team-x
username: team-X
custom_id: my_team_x # optional and not recommended, please use `username`

---

apiVersion: configuration.konghq.com/v1
kind: KongCredential
metadata:
  name: credential-team-x
consumerRef: consumer-team-x
type: key-auth
config:
  key: 62eb165c070a41d5c1b58d9d3d725ca1

---

apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: http-svc-consumer-ratelimiting
consumerRef: consumer-team-x
config:
  hour: 1000
  limit_by: ip
  second: 100
plugin: rate-limiting
```

### KongIngress

Ingress resource spec in Kubernetes can define routing policies based on HTTP Host header and paths.
While this is sufficient in most cases, sometimes, users may want more control over routing at the Ingress level.
`KongIngress` works in conjunction with existing Ingress resource and extends it. It is not meant as a replacement to the
`Ingress` resource in Kubernetes.
Using `KongIngress`, all properties of [Upstream][5], [Service][6] and [Route][7] entitise in Kong related to an Ingress resource
can be modified.

Once a `KongIngress` resource is created, it can be associated with `Ingress` resource in two ways:
- Create a `KongIngress` object in the same namespace as that of the Ingress rule using the same name.
  This avoids a need of additional annotation in Ingress oresource.
  On the other hand, this approach requires a `KongIngress` resource per Ingress, which becomes hard to maintain with multiple Ingresses.

- Create an `KongIngress` resource and then using the annotation `configuration.konghq.com: <KongIngress-resource-name>`,
  associate it with one or more Ingress resources. This approach allows you to reuse the same `KongIngress`.

*Note:* Is not required to define the complete object, we can define the `upstream`, `proxy` or `route` sections

Following is a complete spec example:

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

[0]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/
[1]: https://getkong.org/docs/0.14.x/admin-api/#consumer-object
[2]: https://getkong.org/docs/0.14.x/admin-api/#plugin-object
[3]: https://kubernetes.io/docs/reference/kubectl/overview/
[4]: https://getkong.org/docs/0.14.x/admin-api/#add-plugin
[5]: https://getkong.org/docs/0.14.x/admin-api/#upstream-objects
[6]: https://getkong.org/docs/0.14.x/admin-api/#service-object
[7]: https://getkong.org/docs/0.14.x/admin-api/#route-object
