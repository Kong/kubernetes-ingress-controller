# Using mTLS with Kong

This guide walks through on how to setup Kong to perform mutual-TLS
authentication with an upstream service.

> Please note that this guide walks through mTLS configuration between
Kong and a Service and not Kong and a client or consumer.

## What is mTLS?

Mutual authentication refers to two-way authencation, where the client and
server, both can authenticate themselves to the other party.

With mutual TLS authentication, client and server both present TLS
certificates to the other party (and can prove their identity using their
private key) during the TLS handshake. They can verify the other's
certificate using the their trusted CAs.

## mTLS with Kong

Kong 1.3 and above support mutual TLS authentication between Kong and the
upstream service.

Let's take a look at how one can configure it.

## Configure Kong to verify upstream server certificate

Kong, by default, does not verify the certificate presented by the upstream
service.

To enforce certificate verification, you need to configure the following
environment variables on Kong's container in your deployment:

```
KONG_NGINX_PROXY_PROXY_SSL_VERIFY="on"
KONG_NGINX_PROXY_PROXY_SSL_VERIFY_DEPTH="on"
KONG_NGINX_PROXY_PROXY_SSL_TRUSTED_CERTIFICATE="/path/to/ca_certs.pem"
```

These basically translate to
[NGINX directives](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
to configure NGINX to verify certificates.

Please make sure that the trusted certificates are correctly
mounted into Kong's container and the path to certificate is correctly
reflected in the above environment variable.

## Configure Kong to present its certificate to the upstream server

In the above section, we achieved one side of mutual authentication,
where Kong has been configured to verify the identity of the upstream server.

In this section, we will configure Kong to present its identity to the
upstream server.

To configure this, you have two options, depending on your use-case.
If you would like Kong to present its client certificate to each and every
service that it talks to, you can configure the client certificate
at the global level using Nginx directives.
If you would like to configure a different certificate for
each service that Kong talks to or want to configure Kong to present a
client certificate only to a subset of all services that it is configured to
communicate with, then you can configure that using an annotation on
the Kubernetes Service resource.

### Global Nginx directive

You need to configure two Nginx directives for this purpose:
- [`proxy_ssl_certificate`](https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_ssl_certificate)
- [`proxy_ssl_certificate_key`](https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_ssl_certificate_key)

You can mount the certificate and key pair using secrets into the Kong pod
and then set the following two environment variables to set the above two
directives:

```
KONG_NGINX_PROXY_PROXY_SSL_CERTIFICATE="/path/to/client_cert.pem"
KONG_NGINX_PROXY_PROXY_SSL_CERTIFICATE_KEY="/path/to/key.pem"
```

Once configured, Kong will present its client certificate to every upstream
server that it talks to.

### Per service annotation

To configure a different client certificate for each service or only for a
subset of services, you can do so using the
[`configuration.konghq.com/client-cert`](../references/annotations.md#configurationkonghqcom/client-cert)
annotation.

To use the annotation, you first need to create a TLS secret with the 
client certificate and key in Kubernetes.
The secret should be created in the same namespace as your Kubernetes
Service to which Kong should authenticate itself.

Once the secret is in place, add the follow annotation on the service:

```
configuration.konghq.com/client-cert: <name-of-secret>
```

Kong will then use the TLS key-pair to authenticate itself against that service.
