# Ingress v1 and v1beta1 differences

## Table of Contents

- [**Introduction**](#introduction)
- [**Paths**](#paths)
- [**Ingress Class**](#ingress-class)
- [**Hostnames**](#hostnames)
- [**Backend Types**](#backend-types)

## Introduction

The Kubernetes [Ingress resource](https://kubernetes.io/docs/concepts/services-networking/ingress/)
has evolved. TODO Expand this later!

## Paths

Both Ingress v1beta1 and v1 HTTP rules require a path, which represents a [URI
path][uri-rfc-paths]. Although v1beta1 had specified that paths were [POSIX
regular expressions][posix-regex] and enforced this, but in practice most
controllers used other other implementations that did not match the
specification. v1 seeks to reduce confusion by introducing several i[path
types][path-types] and lifting restrictions on regular expression grammars used
by controllers.

### v1beta1

The controller passes paths directly to Kong and relies on its [path handling
logic][kong-path]. The Kong proxy treats paths as a prefix unless they include
characters [not allowed in RFC 3986 paths][uri-rfc-paths], in which case the
proxy assumes they are a regular expression, and does not treat slashes as
special characters. For example, the prefix `/foo` can match any of the
following:

```
/foo
/foo/
/foobar
/foo/bar
```

### v1

Although v1 Ingresses provide path types with more clearly-defined logic, the
controller must still create Kong routes and work within the Kong proxy's
routing logic. As such, the controller translates Ingress rule paths to create
Kong routes that match the specification:

#### Exact

To handle exact matches, the controller creates a Kong route with a regular
expression that matches the rule path only, e.g. an exact rule for `/foo` in an
Ingress translates to a Kong route with a `/foo$` regular expression path.

#### Prefix

To handle prefix matches, the controller creates a Kong route with two path
criteria, e.g. `/foo` will create a route with a `/foo$` regular expression and
`/foo/` plain path.

#### ImplementationSpecific

The controller leaves ImplementationSpecific path rules entirely up to the Kong
router. It creates a route with the exact same path string as the Ingress rule.

## Ingress class

[Ingress class][ingress-class] indicates which resources an ingress controller
should process. It provides a means to separate out configuration intended for
other controllers or other instances of the Kong Ingress Controller.

In v1beta1, ingress class was handled informally using
`kubernetes.io/ingress.class` [annotations][deprecated-annotation]. v1
introduces a new [IngressClass resource][ingress-class-api] which provides
richer information about the controller. v1 Ingresses are bound to a class via
their `ingressClassName` field.

For example, consider this v1beta1 Ingress:

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: example-ingress
  annotations:
    kubernetes.io/ingress.class: "kong"
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /test
        backend:
          serviceName: echo
          servicePort: 80
```

Its ingress class annotation is set to `kong`, and ingress controllers set to
process `kong` class Ingresses will process it.

In v1, the equivalent configuration declares a `kong` IngressClass resource
whose `metadata.name` field indicates the class name. The `ingressClassName`
value of the Ingress object must match the value of the `name` field in the
IngressClass metadata:

```
apiVersion: networking.k8s.io/v1
kind: IngressClass
metadata:
  name: kong
spec:
  controller: konghq.com/kong
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
spec:
  ingressClassName: kong
  rules:
  - http:
      paths:
      - path: /testpath
        pathType: Prefix
        backend:
          service:
            name: test
            port:
              number: 80
```

## Hostnames

Ingress v1 formally codifies support for [wildcard hostnames][wildcard-hostnames].
v1beta1 Ingresses did not reject wildcard hostnames, however, and Kong had
[existing support for them][kong-wildcard-hostnames].

As such, while the v1beta1 specification did not officially support wildcard
hostnames, you can use wildcard hostnames with either version. Setting a
hostname like `*.example.com` will match requests for both `foo.example.com`
and `bar.example.com` with either v1 or v1beta1 Ingresses.

## Backend types

Ingress v1 introduces support for backends other than Kubernetes Services, via
[resource backends][resource-backends].

Kong has support for routes without a service in some cases, such as routes
that use the [AWS Lambda plugin][lambda-plugin]. However, it does not yet offer
Kubernetes resources for them. As such, you should instead create a placeholder
Kubernetes Service for them, using an [ExternalName Service][external-name]
with an RFC 2606 invalid hostname, e.g. `kong.invalid`. You can use these
placeholder services with either v1 or v1beta1 Ingresses.

[ingress-class]: ./ingress-classes.md
[uri-rfc-paths]: https://tools.ietf.org/html/rfc3986#section-3.3
[posix-regex]: https://www.boost.org/doc/libs/1_38_0/libs/regex/doc/html/boost_regex/syntax/basic_extended.html
[path-types]: https://kubernetes.io/docs/concepts/services-networking/ingress/#path-types
[kong-paths]: https://docs.konghq.com/latest/proxy/#request-path
[wildcard-hostnames]: https://kubernetes.io/docs/concepts/services-networking/ingress/#hostname-wildcards
[kong-wildcard-hostnames]: https://docs.konghq.com/latest/proxy/#using-wildcard-hostnames
[resource-backends]: https://kubernetes.io/docs/concepts/services-networking/ingress/#resource-backend
[lambda-plugin]: https://docs.konghq.com/hub/kong-inc/aws-lambda/
[external-name]: https://kubernetes.io/docs/concepts/services-networking/service/#externalname
[deprecated-annotation]: https://kubernetes.io/docs/concepts/services-networking/ingress/#deprecated-annotation
[ingress-class-api]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#ingressclass-v1-networking-k8s-io
