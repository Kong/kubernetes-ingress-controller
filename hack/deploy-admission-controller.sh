#!/usr/bin/env bash

set -o nounset
set -o pipefail

if [[ ${#} -eq 1 && -n "${1}" ]]; then
  export KUBECONFIG="${1}"
fi

BASE64_OPTIONS=""
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  BASE64_OPTIONS="-w0"
fi

# create a self-signed certificate
TMPDIR="$(mktemp -d )"
openssl req -x509 -newkey rsa:2048 -keyout "${TMPDIR}/tls.key" -out "${TMPDIR}/tls.crt" -days 365  \
    -nodes -subj "/CN=kong-validation-webhook.kong.svc" \
    -extensions EXT -config <( \
   printf "[dn]\nCN=kong-validation-webhook.kong.svc\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:kong-validation-webhook.kong.svc\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
# create a secret out of this self-signed cert-key pair
kubectl create secret tls kong-validation-webhook -n kong \
      --key "${TMPDIR}/tls.key" --cert "${TMPDIR}/tls.crt"
# enable the Admission Webhook Server server
kubectl patch -n kong deploy/ingress-kong \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"ingress-controller","env":[{"name":"CONTROLLER_ADMISSION_WEBHOOK_LISTEN","value":":8080"}],"volumeMounts":[{"name":"validation-webhook","mountPath":"/admission-webhook"}]}],"volumes":[{"secret":{"secretName":"kong-validation-webhook"},"name":"validation-webhook"}]}}}}'

readonly CABUNDLE=$(base64 ${BASE64_OPTIONS:+${BASE64_OPTIONS}} < ${TMPDIR}/tls.crt)

# configure k8s apiserver to send validations to the webhook
(
cat << EOF
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: kong-validations
webhooks:
- name: validations.kong.konghq.com
  objectSelector:
    matchExpressions:
    - key: owner
      operator: NotIn
      values:
      - helm
  failurePolicy: Ignore
  sideEffects: None
  admissionReviewVersions: ["v1", "v1beta1"]
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
    - kongconsumergroups
    - kongplugins
    - kongclusterplugins
    - kongingresses
    - kongvaults
  - apiGroups:
    - ''
    apiVersions:
    - 'v1'
    operations:
    - CREATE
    - UPDATE
    resources:
    - secrets
    - services
  - apiGroups:
    - networking.k8s.io
    apiVersions:
      - 'v1'
    operations:
    - CREATE
    - UPDATE
    resources:
    - ingresses
  - apiGroups:
    - gateway.networking.k8s.io
    apiVersions:
    - 'v1alpha2'
    - 'v1beta1'
    - 'v1'
    operations:
    - CREATE
    - UPDATE
    resources:
    - gateways
    - httproutes
  clientConfig:
    service:
      namespace: kong
      name: kong-validation-webhook
    caBundle: ${CABUNDLE}
EOF
) | kubectl apply -f -
