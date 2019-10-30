# Validating Admission Controller

Kong Ingress Controller ships with an Admission Controller for KongPlugin
and KongConsumer resources in the `configuration.konghq.com` API group.

The Admission Controller needs a TLS certificate and key pair which
you need to generate as part of the deployment.

Following guide walks through a setup of how to create the required key-pair
and enable the admission controller.

Please note that this requires Kong Ingress Controller >= 0.6.

## Create a certificate for the admission controller

Kuberentes API-server makes an HTTPS call to the Admission Controller to verify
if the custom resource is valid or not. For this to work, Kubernetes API-server
needs to trust the CA certificate that is used to sign Admission Controller's
TLS certificate.

This can be accomplished either using a self-signed certificate or using
Kubernetes CA. Follow one of the steps below and then go to
[Create the secret](#create-the-secret) step below.

Please note the `CN` field of the x509 certificate takes the form
`<validation-service-name>.<ingress-controller-namespace>.svc`, which
in the default case is `kong-validation-webhook.kong.svc`.

### Using self-signed certificate

Use openssl to generate a self-signed certificate:

```bash
$ openssl req -x509 -newkey rsa:2048 -keyout tls.key -out tls.crt -days 365  \
    -nodes -subj "/CN=kong-validation-webhook.kong.svc"
Generating a 2048 bit RSA private key
..........................................................+++
.............+++
writing new private key to 'key.pem'
```

### Using in-built Kubernetes CA

Kubernetes comes with an in-built CA which can be used to provision
a certificate for the Admission Controller.
Please refer to the
[this guide](https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/)
on how to generate a certificate using the in-built CA.

### Create the secret

Next, create a Kubernetes secret object based on the key and certificate that
was generatd in the previous steps.
Here, we assume that the PEM-encoded certificate is stored in a file named
`tls.crt` and private key is stored in `tls.key`.

```bash
$ kubectl create secret tls kong-validation-webhook -n kong \
    --key tls.key --cert tls.crt
secret/kong-validation-webhook created
```

## Update the deployment

Once the secret is created, update the Ingress Controller deployment:

Execute the following command to patch the Kong Ingress Controller deployment
to mount the certificate and key pair:

```bash
$ kubectl patch deploy -n kong ingress-kong
    -p '{"spec":{"template":{"spec":{"containers":[{"name":"ingress-controller","volumeMounts":[{"name":"validation-webhook","mountPath":"/admission-webhook"}]}],"volumes":[{"secret":{"secretName":"kong-validation-webhook"},"name":"validation-webhook"}]}}}}'
deployment.extensions/ingress-kong patched
```

## Enable the validating admission

If you are using Kubernetes CA to generate the certificate, you don't need
to supply a CA certificate (in the `caBunde` param)
as part of the Validation Webhook configuration
as the API-server already trusts the internal CA.

```bash
$ echo "apiVersion: admissionregistration.k8s.io/v1beta1
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
  failurePolicy: Fail
  clientConfig:
    service:
      namespace: kong
      name: kong-validation-webhook
    caBundle: $(cat tls.crt  | base64) " | kubectl apply -f -
```

## Verify if it works

### Verify duplicate KongConsumers

Create a KongConsumer with username as `harry`:

```bash
$ echo "apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: harry
username: harry" | kubectl apply -f -
kongconsumer.configuration.konghq.com/harry created
```

Now, create another KongConsumer with the same username:

```bash
$ echo "apiVersion: configuration.konghq.com/v1
kind: KongConsumer
metadata:
  name: harry2
username: harry" | kubectl apply -f -
Error from server: error when creating "STDIN": admission webhook "validations.kong.konghq.com" denied the request: consumer already exists
```

The validation webhook rejected the KongConsumer resource as there already
exists a consumer in Kong with the same username.

### Verify incorrect KongPlugins

Try to create the folowing KongPlugin resource.
The `foo` config property does not exist in the configuration definition and
hence the Admission Controller returns back an error.
If you remove the `foo: bar` configuration line, the plugin will be
created succesfully.

```bash
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: request-id
config:
  foo: bar
  header_name: my-request-id
plugin: correlation-id
" | kubectl apply -f -
Error from server: error when creating "STDIN": admission webhook "validations.kong.konghq.com" denied the request: 400 Bad Request {"fields":{"config":{"foo":"unknown field"}},"name":"schema violation","code":2,"message":"schema violation (config.foo: unknown field)"}
```
