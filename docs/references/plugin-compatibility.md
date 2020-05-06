# Plugin Compatibility

DB-less mode is the preferred choice for controller-managed Kong and Kong
Enterprise clusters. However, not all plugins are available in DB-less mode.
Review the table below to check if a plugin you wish to use requires a
database.

Note that some DB-less compatible plugins [have some limitations or require
non-default configuration for
compatibility](https://docs.konghq.com/latest/db-less-and-declarative-config/#plugin-compatibility).

## Kong

|Plugin                          |  Kong                |  Kong (DB-less)    |
|--------------------------------|----------------------|--------------------|
|                                |  :white_check_mark:  |  :white_check_mark:|
|acl                             |  :white_check_mark:  |  :white_check_mark:|
|aws-lambda                      |  :white_check_mark:  |  :white_check_mark:|
|azure-functions                 |  :white_check_mark:  |  :white_check_mark:|
|basic-auth                      |  :white_check_mark:  |  :white_check_mark:|
|bot-detection                   |  :white_check_mark:  |  :white_check_mark:|
|correlation-id                  |  :white_check_mark:  |  :white_check_mark:|
|cors                            |  :white_check_mark:  |  :white_check_mark:|
|datadog                         |  :white_check_mark:  |  :white_check_mark:|
|file-log                        |  :white_check_mark:  |  :white_check_mark:|
|hmac-auth                       |  :white_check_mark:  |  :white_check_mark:|
|http-log                        |  :white_check_mark:  |  :white_check_mark:|
|ip-restriction                  |  :white_check_mark:  |  :white_check_mark:|
|jwt                             |  :white_check_mark:  |  :white_check_mark:|
|key-auth                        |  :white_check_mark:  |  :white_check_mark:|
|kubernetes-sidecar-injector     |  :white_check_mark:  |  :white_check_mark:|
|oauth2                          |  :white_check_mark:  |  :x:               |
|prometheus                      |  :white_check_mark:  |  :white_check_mark:|
|rate-limiting                   |  :white_check_mark:  |  :white_check_mark:|
|request-termination             |  :white_check_mark:  |  :white_check_mark:|
|request-transformer             |  :white_check_mark:  |  :white_check_mark:|
|response-ratelimiting           |  :white_check_mark:  |  :white_check_mark:|
|response-transformer            |  :white_check_mark:  |  :white_check_mark:|
|syslog                          |  :white_check_mark:  |  :white_check_mark:|
|tcp-log                         |  :white_check_mark:  |  :white_check_mark:|
|udp-log                         |  :white_check_mark:  |  :white_check_mark:|
|zipkin                          |  :white_check_mark:  |  :white_check_mark:|
|application-registration        |  :x:                 |  :x:               |
|canary release                  |  :x:                 |  :x:               |
|collector                       |  :x:                 |  :x:               |
|degraphql                       |  :x:                 |  :x:               |
|exit-transformer                |  :x:                 |  :x:               |
|forward-proxy                   |  :x:                 |  :x:               |
|graphql-proxy-cache-advanced    |  :x:                 |  :x:               |
|graphql-rate-limiting-advanced  |  :x:                 |  :x:               |
|jwt-signer                      |  :x:                 |  :x:               |
|kafka-log                       |  :x:                 |  :x:               |
|kafka-upstream                  |  :x:                 |  :x:               |
|key-auth-enc                    |  :x:                 |  :x:               |
|ldap-auth-advanced              |  :x:                 |  :x:               |
|mtls-auth                       |  :x:                 |  :x:               |
|oauth2-introspection            |  :x:                 |  :x:               |
|openid-connect                  |  :x:                 |  :x:               |
|proxy-cache-advanced            |  :x:                 |  :x:               |
|rate-limiting-advanced          |  :x:                 |  :x:               |
|request-transformer-advanced    |  :x:                 |  :x:               |
|request-validator               |  :x:                 |  :x:               |
|response-transformer-advanced   |  :x:                 |  :x:               |
|route-transformer-advanced      |  :x:                 |  :x:               |
|statsd-advanced                 |  :x:                 |  :x:               |
|vault-auth                      |  :x:                 |  :x:               |

## Kong Enterprise

There are [two distributions of Kong Enterprise](https://github.com/Kong/kubernetes-ingress-controller/tree/master/docs/deployment#overview),
`kong-enterprise-k8s` (used with Kong for Kubernetes Enterprise) and
`kong-enterprise-edition` (used with Kong for Kubernetes with Kong Enterprise).

|Plugin                          |  Kong for Kubernetes Enterprise  |  Kong for Kubernetes with Kong Enterprise  |
|--------------------------------|----------------------------------|--------------------------------------------|
|                                |  :white_check_mark:              |  :white_check_mark:                        |
|acl                             |  :white_check_mark:              |  :white_check_mark:                        |
|aws-lambda                      |  :white_check_mark:              |  :white_check_mark:                        |
|azure-functions                 |  :white_check_mark:              |  :white_check_mark:                        |
|basic-auth                      |  :white_check_mark:              |  :white_check_mark:                        |
|bot-detection                   |  :white_check_mark:              |  :white_check_mark:                        |
|correlation-id                  |  :white_check_mark:              |  :white_check_mark:                        |
|cors                            |  :white_check_mark:              |  :white_check_mark:                        |
|datadog                         |  :white_check_mark:              |  :white_check_mark:                        |
|file-log                        |  :white_check_mark:              |  :white_check_mark:                        |
|hmac-auth                       |  :white_check_mark:              |  :white_check_mark:                        |
|http-log                        |  :white_check_mark:              |  :white_check_mark:                        |
|ip-restriction                  |  :white_check_mark:              |  :white_check_mark:                        |
|jwt                             |  :white_check_mark:              |  :white_check_mark:                        |
|key-auth                        |  :white_check_mark:              |  :white_check_mark:                        |
|kubernetes-sidecar-injector     |  :white_check_mark:              |  :white_check_mark:                        |
|oauth2                          |  :white_check_mark:              |  :x:                                       |
|prometheus                      |  :white_check_mark:              |  :white_check_mark:                        |
|rate-limiting                   |  :white_check_mark:              |  :white_check_mark:                        |
|request-termination             |  :white_check_mark:              |  :white_check_mark:                        |
|request-transformer             |  :white_check_mark:              |  :white_check_mark:                        |
|response-ratelimiting           |  :white_check_mark:              |  :white_check_mark:                        |
|response-transformer            |  :white_check_mark:              |  :white_check_mark:                        |
|syslog                          |  :white_check_mark:              |  :white_check_mark:                        |
|tcp-log                         |  :white_check_mark:              |  :white_check_mark:                        |
|udp-log                         |  :white_check_mark:              |  :white_check_mark:                        |
|zipkin                          |  :white_check_mark:              |  :white_check_mark:                        |
|application-registration        |  :white_check_mark:              |  :x:<sup>1</sup>                           |
|canary release                  |  :white_check_mark:              |  :white_check_mark:                        |
|collector                       |  :white_check_mark:              |  :white_check_mark:                        |
|degraphql                       |  :white_check_mark:              |  :white_check_mark:                        |
|exit-transformer                |  :white_check_mark:              |  :x:                                       |
|forward-proxy                   |  :white_check_mark:              |  :white_check_mark:                        |
|graphql-proxy-cache-advanced    |  :white_check_mark:              |  :white_check_mark:                        |
|graphql-rate-limiting-advanced  |  :white_check_mark:              |  :white_check_mark:                        |
|jwt-signer                      |  :white_check_mark:              |  :white_check_mark:                        |
|kafka-log                       |  :white_check_mark:              |  :white_check_mark:                        |
|kafka-upstream                  |  :white_check_mark:              |  :white_check_mark:                        |
|key-auth-enc                    |  :white_check_mark:              |  :x:                                       |
|ldap-auth-advanced              |  :white_check_mark:              |  :white_check_mark:                        |
|mtls-auth                       |  :white_check_mark:              |  :white_check_mark:                        |
|oauth2-introspection            |  :white_check_mark:              |  :white_check_mark:                        |
|openid-connect                  |  :white_check_mark:              |  :white_check_mark:                        |
|proxy-cache-advanced            |  :white_check_mark:              |  :white_check_mark:                        |
|rate-limiting-advanced          |  :white_check_mark:              |  :white_check_mark:                        |
|request-transformer-advanced    |  :white_check_mark:              |  :x:<sup>2</sup>                           |
|request-validator               |  :white_check_mark:              |  :white_check_mark:                        |
|response-transformer-advanced   |  :white_check_mark:              |  :white_check_mark:                        |
|route-transformer-advanced      |  :white_check_mark:              |  :x:                                       |
|statsd-advanced                 |  :white_check_mark:              |  :x:<sup>3</sup>                           |
|vault-auth                      |  :white_check_mark:              |  :white_check_mark:                        |

<sup>1</sup> Only used with Dev Portal

<sup>2</sup> request-transformer now has feature parity with
  request-transformer-advanced. request-transformer-advanced remains only for
  compatibility with existing configurations.

<sup>3</sup> Only used with Vitals
