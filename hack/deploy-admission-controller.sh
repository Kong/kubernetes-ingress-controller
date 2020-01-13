#!/bin/bash

# create a self-signed certificate
openssl req -x509 -newkey rsa:2048 -keyout tls.key -out tls.crt -days 365  \
  -nodes -subj "/CN=kong-validation-webhook.kong.svc"
# create a secret out of this self-signed cert-key pair
kubectl create secret tls kong-validation-webhook -n kong \
      --key tls.key --cert tls.crt
# enable the Admission Webhook Server server
kubectl patch deploy -n kong ingress-kong \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"ingress-controller","env":[{"name":"CONTROLLER_ADMISSION_WEBHOOK_LISTEN","value":":8080"}],"volumeMounts":[{"name":"validation-webhook","mountPath":"/admission-webhook"}]}],"volumes":[{"secret":{"secretName":"kong-validation-webhook"},"name":"validation-webhook"}]}}}}'
# configure k8s apiserver to send validations to the webhook
echo "apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: kong-validations
webhooks:
- name: validations.kong.konghq.com
  rules:
  - apiGroups:
    - configuration.konghq.com
    apiVersions:
    - '*'
    operations:
    - CREATE
    - UPDATE
    resources:
    - kongconsumers
    - kongplugins
  - apiGroups:
    - ''
    apiVersions:
    - 'v1'
    operations:
    - CREATE
    - UPDATE
    resources:
    - secrets
  failurePolicy: Fail
  clientConfig:
    service:
      namespace: kong
      name: kong-validation-webhook
    caBundle: $(cat tls.crt  | base64 -w 0) " | kubectl apply -f -

