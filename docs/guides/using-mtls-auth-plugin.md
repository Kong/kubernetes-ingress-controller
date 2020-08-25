# Using mtls-auth plugin

This guide walks through how to configure Kong Ingress Controller to
verify client certificates using CA certificates and
[mtls-auth](https://docs.konghq.com/hub/kong-inc/mtls-auth/) plugin
for HTTPS requests.

> Note: You need an Enterprise license to use this feature.

## Installation

Please follow the [deployment](../deployment) documentation to install
Kong for Kubernetes Enterprise on your Kubernetes cluster.

## Testing Connectivity to Kong

This guide assumes that the `PROXY_IP` environment variable is
set to contain the IP address or URL pointing to Kong.
Please follow one of the
[deployment guides](../deployment/k4k8s-enterprise.md) to configure
this environment variable.

If everything is set up correctly, making a request to Kong should return
HTTP 404 Not Found.

```bash
$ curl -i $PROXY_IP
HTTP/1.1 404 Not Found
Date: Fri, 21 Jun 2019 17:01:07 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 48
Server: kong/1.2.1

{"message":"no Route matched with those values"}
```

This is expected as Kong does not yet know how to proxy the request.

## Provision a CA certificate in Kong

CA certificates in Kong are provisioned by create a `Secret` resource in
Kubernetes.

The secret resource should have a few properties:
- It should have the `konghq.com/ca-cert: "true"` label.
- It should have a `cert` data property which contains a valid CA certificate
  in PEM format.
- It should have an `id` data property which contains a random UUID.

Note that a self-signed CA certificate is being used for the purpose of this
guide. You should use your own CA certificate that is backed by
your PKI infrastructure.

```bash
$ echo "apiVersion: v1
kind: Secret
metadata:
  name: my-ca-cert
  labels:
    konghq.com/ca-cert: 'true'
type: Opaque
stringData:
  cert: |
    -----BEGIN CERTIFICATE-----
    MIICwTCCAamgAwIBAgIUHGUzUWvHJHrREvIZIcORiFUvze4wDQYJKoZIhvcNAQEL
    BQAwEDEOMAwGA1UEAwwFSGVsbG8wHhcNMjAwNTA4MjExODA1WhcNMjAwNjA3MjEx
    ODA1WjAQMQ4wDAYDVQQDDAVIZWxsbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
    AQoCggEBANCMMBngjuTvqts8ZXtZhqdr181QH/NmytW1KlyqZd6ppXUer+i0OWhP
    1nAyHsBPJljKAFLd8l1EioPFkN78/wJFDJrHOtfniIQPVLdS2cnNQ72dLyQH6smH
    JQDV8ePBQ2GdRP6s61+Da8eoaW6nSLtmEUhxvyteboqwmi2CtUtAfuiU1m5sOdpS
    z+L4D08CE+SFIT4MGD3gxNdg7lccWCHIfk54VRSdGDKEVwed8OQvxD0TdpHY+ym5
    nJ4JSkhiS9XIodnxR3AZ6rIPRqk+MQ4LGTjX2EbM0/Yg4qvnZ7m4fcpK2goDZIVL
    EF8F+ka1RaAYWTsXI1BAkJbb3kdo/yUCAwEAAaMTMBEwDwYDVR0TBAgwBgEB/wIB
    ADANBgkqhkiG9w0BAQsFAAOCAQEAVvB/PeVZpeQ7q2IQQQpADtTd8+22Ma3jNZQD
    EkWGZEQLkRws4EJNCCIvkApzpx1GqRcLLL9lbV+iCSiIdlR5W9HtK07VZ318gpsG
    aTMNrP9/2XWTBzdHWaeZKmRKB04H4z7V2Dl58D+wxjdqNWsMIHeqqPNKGamk/q8k
    YFNqNwisRxMhU6qPOpOj5Swl2jLTuVMAeGWBWmPGU2MUoaJb8sc2Vix9KXcyDZIr
    eidkzkqSrjNzI0yJ2gdCDRS4/Rw9iV3B3SRMs0mJMLBDrsowhNfLAd8I3NHzLwps
    dZFcvZcT/p717K3hlFVdjGnKIgKcG7aYji/XRR87HKnc+cJMCw==
    -----END CERTIFICATE-----
  id: cce8c384-721f-4f58-85dd-50834e3e733a" | kubectl create -f -
secret/my-ca-cert created
```

Please note the ID, you can use this ID one or use a different one but
the ID is important in the next step when we create the plugin.
Each CA certificate that you create needs a unique ID.
Any random UUID will suffice here and it doesn't have an security
implication.

You can use [uuidgen](https://linux.die.net/man/1/uuidgen) (Linux, OS X) or
[New-Guid](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/new-guid)
(Windows) to generate an ID.

For example:
```bash
$ uuidgen
907821fc-cd09-4186-afb5-0b06530f2524
```

## Configure mtls-auth plugin

Next, we are going to create an `mtls-auth` KongPlugin resource which references
CA certificate provisioned in the last step:

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: mtls-auth
config:
  ca_certificates:
  - cce8c384-721f-4f58-85dd-50834e3e733a
  skip_consumer_lookup: true
  revocation_check_mode: SKIP
plugin: mtls-auth
" | kubectl apply -f -
kongplugin.configuration.konghq.com/mtls-auth created
```

## Install a dummy service

Let's deploy an echo service which we wish to protect
using TLS client certificate authentication.

```bash
$ kubectl apply -f https://bit.ly/echo-service
service/echo created
deployment.apps/echo created
```

You can deploy a different service or skip this step if you already
have a service deployed in Kubernetes.

## Set up Ingress

Let's expose the echo service outside the Kubernetes cluster
by defining an Ingress.

```bash
$ echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
  annotations:
    konghq.com/plugins: mtls-auth
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - http:
      paths:
      - path: /foo
        backend:
          serviceName: echo
          servicePort: 80
" | kubectl apply -f -
ingress.extensions/demo created
```

## Test the endpoint

Now, let's test to see if Kong is asking for client certificate
or not when we make the request:

```
$ curl -k https://$PROXY_IP
HTTP/2 401
date: Mon, 11 May 2020 18:15:05 GMT
content-type: application/json; charset=utf-8
content-length: 50
x-kong-response-latency: 0
server: kong/2.0.4.0-enterprise-k8s

{"message":"No required TLS certificate was sent"}
```

As we can see, Kong is restricting the request because it doesn't
have the necessary authentication information.

Two things to note here:
- `-k` is used because Kong is set up to serve a self-signed certificate
  by default. For full mutual authentication in production use cases,
  you must configure Kong to serve a certificate that is signed by a trusted CA.
- For some deployments `$PROXY_IP` might contain a port that points to
  `http` port of Kong. In others, it might happen that it contains a DNS name
  instead of an IP address. If needed, please update the
  command to send an `https` request to the `https` port of Kong or
  the load balancer in front of it.


## Provisioning credential

Next, in order to authenticate against Kong, create the client
certificate and private key with the following content:

```bash
$ cat client.crt
-----BEGIN CERTIFICATE-----
MIIEFTCCAv0CAWUwDQYJKoZIhvcNAQELBQAwEDEOMAwGA1UEAwwFSGVsbG8wHhcN
MjAwNTA4MjE0OTE1WhcNMjEwNTA4MjE0OTE1WjCBkDELMAkGA1UEBhMCQVUxEzAR
BgNVBAgMClNvbWUtU3RhdGUxDTALBgNVBAcMBHNvbWUxETAPBgNVBAoMCHNvbWUg
b3JnMRAwDgYDVQQLDAdvcmd1bml0MRswGQYDVQQDDBJleGFtcGxlLmtvbmdocS5j
b20xGzAZBgkqhkiG9w0BCQEWDGZvb0Bzb21lLmNvbTCCAiIwDQYJKoZIhvcNAQEB
BQADggIPADCCAgoCggIBAM/y80ppzwGYS7zl+A6fx4Xkjwja+ZUK/AoBDazS3TkR
W1tDFZ71koLd60qK2W1d9Wh0/F3iNTcobVefr02mEcLtl+d4zUug+W7RsK/8JSCM
MIDVDYzlTWdd7RJzV1c/0NFZyTRkEVSjGn6eQoC/1aviftiNyfqWtuIDQ5ctSBt8
2fyvDwu/tBR5VyKu7CLnjZ/ffjNT8WDfbO704XeBBId0+L8i8J7ddYlRhZufdjEw
hKx2Su8PZ9RnJYShTBOpD0xdveh16eb7dpCZiPnp1/MOCyIyo1Iwu570VoMde9SW
sPFLdUMiCXw+A4Gp/e9Am+D/98PiL4JChKsiowbzpDfMrVQH4Sblpcgn/Pp+u1be
2Kl/7wqr3TA+w/unLnBnB859v3wDhSW4hhKASoFwyX3VfJ43AkmWFUBX/bpDvHto
rFw+MvbSLsS3QD5KlZmega1pNZtin5KV8H/oJI/CjEc9HHwd27alW9VkUu0WrH0j
c98wLHB/9xXLjunabxSmd+wv25SgYNqpsRNOLgcJraJbaRh4XkbDyuvjF2bRJVP4
pIjntxQHS/oDFFFK3wc7fp/rTAl0PJ7tytYj4urg45N3ts7unwnB8WmKzD9Avcwe
8Kst12cEibS8X2sg8wOqgB0yarC17mBEqONK7Fw4VH+VzZYw0KGF5DWjeSXj/XsD
AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAEvTMHe27npmyJUBxQeHcNFniMJUWZf0
i9EGd+XlF+m/l3rh1/mCecV7s32QTZEiFHv4UJPYASbgtx7+mEZuq7dVsxIUICWs
gyRkwvKjMqK2tR5IRkquhK5PuDS0QC3M/ZsDwnTgaezFrplFYf80z1kAAkm/c7eh
ZEjI6+1vuaS+HX1w2unk42PiAEB6oKFi3b8xl4TC6acYfMYiC3cOa/d3ZKHhqXhT
wM0VtDe0Qn1kExe+19XJG5cROelxmMXBm1+/c2KUw1yK8up6kJlEsmd8JLw/wMUp
xcJUKIH1qGBlRlFTYbVell+dB7IkHhadrnw27Z47uHobB/lzN69r63c=
-----END CERTIFICATE-----
```

```bash
$ cat client.pem
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
```

Now, use the key and certificate to authenticate against Kong and use the
service:

```bash
$ curl --key client.key --cert client.crt  https://$PROXY_IP/foo -k -I
HTTP/2 200
content-type: text/plain; charset=UTF-8
date: Mon, 11 May 2020 18:27:22 GMT
server: echoserver
x-kong-upstream-latency: 1
x-kong-proxy-latency: 1
via: kong/2.0.4.0-enterprise-k8s
```

## Conclusion

This guide demonstrates how to implement client TLS authentication
using Kong.
You are free to use other features that mtls-auth plugin in Kong to
achieve more complicated use-cases.
