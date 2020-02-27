# Configuring ACL Plugin

This guide walks through configuring the Kong ACL Plugin. The ACL Plugin requires the use of at least one Authentication plugin. This example will use the JWT Auth Plugin

## Secrets

In order to configure a consumer for the JWT plugin, a public key needs to be created in a `secret`

```bash
# create secret for jwt public key
kubectl create secret \
  generic app-admin-jwt  \
  --from-literal=kongCredType=jwt  \
  --from-literal=algorithm=RS256 \
  --from-literal=rsa_public_key="-----BEGIN PUBLIC KEY-----\nMIIBIjA....-----END PUBLIC KEY-----"

# create a second secret with a different key
kubectl create secret \
  generic app-user-jwt  \
  --from-literal=kongCredType=jwt  \
  --from-literal=algorithm=RS256 \
  --from-literal=rsa_public_key="-----BEGIN PUBLIC KEY-----\nMIIBIjA....-----END PUBLIC KEY-----"
```

In order to assign a consumer to an ACL group, a `secret` needs to be created with the group name

```bash
# create secrets for acl group
kubectl create secret \
  generic app-admin-acl  \
  --from-literal=kongCredType=acl  \
  --from-literal=group=app-admin

kubectl create secret \
  generic app-user-acl  \
  --from-literal=kongCredType=acl  \
  --from-literal=group=app-user
```

## Consumers

To illustrate having more than one acl group, the next step will be to create two consumers

```yaml
---
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: admin
username: admin
credentials:
  - app-admin-jwt
  - app-admin-acl
---
apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: plain-user
username: plain-user
credentials:
  - app-user-jwt
  - app-user-acl
---

```

## Plugins

To enable JWT authentication, the jwt plugin needs to be configured

```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: app-jwt
plugin: jwt
```

To enable ACL's, the acl plugin needs to be configured

```yaml
---
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: admin-acl
plugin: acl
config:
  whitelist: ['app-admin']
---
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: plain-user-acl
plugin: acl
config:
  whitelist: ['app-user']
---

```

## Ingress Annotation

Last step. Annotate the ingress resources with the jwt and acl plugins

```yaml
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
  annotations:
    plugins.konghq.com: app-jwt,admin-acl
spec:
  rules:
    - http:
        paths:
          - path: /console
            backend:
              serviceName: my-service
              servicePort: 80
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
  annotations:
    plugins.konghq.com: app-jwt,user-acl
spec:
  rules:
    - http:
        paths:
          - path: /
            backend:
              serviceName: my-service
              servicePort: 80
---

```
