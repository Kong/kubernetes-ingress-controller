## Packages
- [configuration.konghq.com/v1](#configurationkonghqcomv1)


## configuration.konghq.com/v1

Package v1 contains API Schema definitions for the konghq.com v1 API group

### Resource Types
- [KongClusterPlugin](#kongclusterplugin)
- [KongClusterPluginList](#kongclusterpluginlist)
- [KongConsumer](#kongconsumer)
- [KongConsumerList](#kongconsumerlist)
- [KongIngress](#kongingress)
- [KongIngressList](#kongingresslist)
- [KongPlugin](#kongplugin)
- [KongPluginList](#kongpluginlist)



#### ConfigSource



ConfigSource is a wrapper around SecretValueFromSource

_Appears in:_
- [KongPlugin](#kongplugin)

| Field | Description |
| --- | --- |
| `secretKeyRef` _[SecretValueFromSource](#secretvaluefromsource)_ |  |


#### KongClusterPlugin



KongClusterPlugin is the Schema for the kongclusterplugins API

_Appears in:_
- [KongClusterPluginList](#kongclusterpluginlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `configuration.konghq.com/v1`
| `kind` _string_ | `KongClusterPlugin`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `consumerRef` _string_ | ConsumerRef is a reference to a particular consumer |
| `disabled` _boolean_ | Disabled set if the plugin is disabled or not |
| `config` _[JSON](#json)_ | Config contains the plugin configuration. |
| `configFrom` _[NamespacedConfigSource](#namespacedconfigsource)_ | ConfigFrom references a secret containing the plugin configuration. |
| `plugin` _string_ | PluginName is the name of the plugin to which to apply the config |
| `run_on` _string_ | RunOn configures the plugin to run on the first or the second or both nodes in case of a service mesh deployment. |
| `protocols` _[KongProtocol](#kongprotocol) array_ | Protocols configures plugin to run on requests received on specific protocols. |
| `ordering` _PluginOrdering_ | Ordering overrides the normal plugin execution order |


#### KongClusterPluginList



KongClusterPluginList contains a list of KongClusterPlugin



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `configuration.konghq.com/v1`
| `kind` _string_ | `KongClusterPluginList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KongClusterPlugin](#kongclusterplugin) array_ |  |


#### KongConsumer



KongConsumer is the Schema for the kongconsumers API

_Appears in:_
- [KongConsumerList](#kongconsumerlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `configuration.konghq.com/v1`
| `kind` _string_ | `KongConsumer`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `username` _string_ | Username unique username of the consumer. |
| `custom_id` _string_ | CustomID existing unique ID for the consumer - useful for mapping Kong with users in your existing database |
| `credentials` _string array_ | Credentials are references to secrets containing a credential to be provisioned in Kong. |


#### KongConsumerList



KongConsumerList contains a list of KongConsumer



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `configuration.konghq.com/v1`
| `kind` _string_ | `KongConsumerList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KongConsumer](#kongconsumer) array_ |  |


#### KongIngress



KongIngress is the Schema for the kongingresses API

_Appears in:_
- [KongIngressList](#kongingresslist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `configuration.konghq.com/v1`
| `kind` _string_ | `KongIngress`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `upstream` _[KongIngressUpstream](#kongingressupstream)_ | Upstream represents a virtual hostname and can be used to loadbalance incoming requests over multiple targets (e.g. Kubernetes `Services` can be a target, OR `Endpoints` can be targets). |
| `proxy` _[KongIngressService](#kongingressservice)_ | Proxy defines additional connection options for the routes to be configured in the Kong Gateway, e.g. `connection_timeout`, `retries`, e.t.c. |
| `route` _[KongIngressRoute](#kongingressroute)_ | Route define rules to match client requests. Each Route is associated with a Service, and a Service may have multiple Routes associated to it. |


#### KongIngressList



KongIngressList contains a list of KongIngress



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `configuration.konghq.com/v1`
| `kind` _string_ | `KongIngressList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KongIngress](#kongingress) array_ |  |


#### KongIngressRoute



KongIngressRoute contains KongIngress route configuration

_Appears in:_
- [KongIngress](#kongingress)

| Field | Description |
| --- | --- |
| `methods` _string array_ | Methods is a list of HTTP methods that match this Route. |
| `headers` _object (keys:string, values:string array)_ | Headers contains one or more lists of values indexed by header name that will cause this Route to match if present in the request. The Host header cannot be used with this attribute. |
| `protocols` _[KongProtocol](#kongprotocol) array_ | Protocols is an array of the protocols this Route should allow. |
| `regex_priority` _integer_ | RegexPriority is a number used to choose which route resolves a given request when several routes match it using regexes simultaneously. |
| `strip_path` _boolean_ | StripPath sets When matching a Route via one of the paths strip the matching prefix from the upstream request URL. |
| `preserve_host` _boolean_ | PreserveHost sets When matching a Route via one of the hosts domain names, use the request Host header in the upstream request headers. If set to false, the upstream Host header will be that of the Service’s host. |
| `https_redirect_status_code` _integer_ | HTTPSRedirectStatusCode is the status code Kong responds with when all properties of a Route match except the protocol. |
| `path_handling` _string_ | PathHandling controls how the Service path, Route path and requested path are combined when sending a request to the upstream. |
| `snis` _string array_ | SNIs is a list of SNIs that match this Route when using stream routing. |
| `request_buffering` _boolean_ | RequestBuffering sets whether to enable request body buffering or not. |
| `response_buffering` _boolean_ | ResponseBuffering sets whether to enable response body buffering or not. |


#### KongIngressService



KongIngressService contains KongIngress service configuration.

_Appears in:_
- [KongIngress](#kongingress)

| Field | Description |
| --- | --- |
| `protocol` _string_ | The protocol used to communicate with the upstream. |
| `path` _string_ | The path to be used in requests to the upstream server.(optional) |
| `retries` _integer_ | The number of retries to execute upon failure to proxy. |
| `connect_timeout` _integer_ | The timeout in milliseconds for establishing a connection to the upstream server. |
| `read_timeout` _integer_ | The timeout in milliseconds between two successive read operations for transmitting a request to the upstream server. |
| `write_timeout` _integer_ | The timeout in milliseconds between two successive write operations for transmitting a request to the upstream server. |


#### KongIngressUpstream



KongIngressUpstream contains KongIngress upstream configuration

_Appears in:_
- [KongIngress](#kongingress)

| Field | Description |
| --- | --- |
| `host_header` _string_ | HostHeader is The hostname to be used as Host header when proxying requests through Kong. |
| `algorithm` _string_ | Algorithm is the load balancing algorithm to use. |
| `slots` _integer_ | Slots is the number of slots in the load balancer algorithm. |
| `healthchecks` _Healthcheck_ | Healthchecks defines the health check configurations in Kong. |
| `hash_on` _string_ | HashOn defines what to use as hashing input. Accepted values are: "none", "consumer", "ip", "header", "cookie", "path", "query_arg", "uri_capture". |
| `hash_fallback` _string_ | HashFallback defines What to use as hashing input if the primary hash_on does not return a hash. Accepted values are: "none", "consumer", "ip", "header", "cookie". |
| `hash_on_header` _string_ | HashOnHeader defines the header name to take the value from as hash input. Only required when "hash_on" is set to "header". |
| `hash_fallback_header` _string_ | HashFallbackHeader is the header name to take the value from as hash input. Only required when "hash_fallback" is set to "header". |
| `hash_on_cookie` _string_ | The cookie name to take the value from as hash input. Only required when "hash_on" or "hash_fallback" is set to "cookie". |
| `hash_on_cookie_path` _string_ | The cookie path to set in the response headers. Only required when "hash_on" or "hash_fallback" is set to "cookie". |
| `hash_on_query_arg` _string_ | HashOnQueryArg is the query string parameter whose value is the hash input when "hash_on" is set to "query_arg". |
| `hash_fallback_query_arg` _string_ | HashFallbackQueryArg is the "hash_fallback" version of HashOnQueryArg. |
| `hash_on_uri_capture` _string_ | HashOnURICapture is the name of the capture group whose value is the hash input when "hash_on" is set to "uri_capture". |
| `hash_fallback_uri_capture` _string_ | HashFallbackURICapture is the "hash_fallback" version of HashOnURICapture. |


#### KongPlugin



KongPlugin is the Schema for the kongplugins API

_Appears in:_
- [KongPluginList](#kongpluginlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `configuration.konghq.com/v1`
| `kind` _string_ | `KongPlugin`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `consumerRef` _string_ | ConsumerRef is a reference to a particular consumer |
| `disabled` _boolean_ | Disabled set if the plugin is disabled or not |
| `config` _[JSON](#json)_ | Config contains the plugin configuration. |
| `configFrom` _[ConfigSource](#configsource)_ | ConfigFrom references a secret containing the plugin configuration. |
| `plugin` _string_ | PluginName is the name of the plugin to which to apply the config |
| `run_on` _string_ | RunOn configures the plugin to run on the first or the second or both nodes in case of a service mesh deployment. |
| `protocols` _[KongProtocol](#kongprotocol) array_ | Protocols configures plugin to run on requests received on specific protocols. |
| `ordering` _[PluginOrdering](#pluginordering)_ | Ordering overrides the normal plugin execution order |


#### KongPluginList



KongPluginList contains a list of KongPlugin



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `configuration.konghq.com/v1`
| `kind` _string_ | `KongPluginList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KongPlugin](#kongplugin) array_ |  |


#### KongProtocol

_Underlying type:_ `string`



_Appears in:_
- [KongClusterPlugin](#kongclusterplugin)
- [KongIngressRoute](#kongingressroute)
- [KongPlugin](#kongplugin)



#### NamespacedConfigSource



NamespacedConfigSource is a wrapper around NamespacedSecretValueFromSource

_Appears in:_
- [KongClusterPlugin](#kongclusterplugin)

| Field | Description |
| --- | --- |
| `secretKeyRef` _[NamespacedSecretValueFromSource](#namespacedsecretvaluefromsource)_ |  |


#### NamespacedSecretValueFromSource



NamespacedSecretValueFromSource represents the source of a secret value specifying the secret namespace

_Appears in:_
- [NamespacedConfigSource](#namespacedconfigsource)

| Field | Description |
| --- | --- |
| `namespace` _string_ | The namespace containing the secret |
| `name` _string_ | the secret containing the key |
| `key` _string_ | the key containing the value |


#### SecretValueFromSource



SecretValueFromSource represents the source of a secret value

_Appears in:_
- [ConfigSource](#configsource)

| Field | Description |
| --- | --- |
| `name` _string_ | the secret containing the key |
| `key` _string_ | the key containing the value |


