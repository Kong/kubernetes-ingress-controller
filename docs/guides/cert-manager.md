# Using cert-manager for automated TLS certificate

This guide will walk through steps to setup Kong Ingress Controller with
cert-manager to automate certificate management using Let's Encrypt.
Any ACME-based CA can be used in-place of Let's Encrypt as well.

## Before you begin

You will need the following:

- Kubernetes cluster that can provision an IP address that is routable from
  the Internet. If you don't have one, you can use GKE or any managed k8s
  cloud offering.
- A domain name for which you control the DNS records.
  This is necessary so that
  Let's Encrypt can verify the ownership of the domain and issue a certificate.
  In the current guide, we use `yolo42.com`, please replace this with a domain
  you control.

This tutorial was written using Google Kubernetes Engine.

## Setup Kong Ingress Controller

Execute the following to install the Ingress Controller:

```bash
$ kubectl create -f https://bit.ly/k4k8s
namespace/kong created
customresourcedefinition.apiextensions.k8s.io/kongplugins.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongconsumers.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongcredentials.configuration.konghq.com created
customresourcedefinition.apiextensions.k8s.io/kongingresses.configuration.konghq.com created
serviceaccount/kong-serviceaccount created
clusterrole.rbac.authorization.k8s.io/kong-ingress-clusterrole created
clusterrolebinding.rbac.authorization.k8s.io/kong-ingress-clusterrole-nisa-binding created
configmap/kong-server-blocks created
service/kong-proxy created
service/kong-validation-webhook created
deployment.extensions/kong created
```

## Setup cert-manager

Please follow cert-manager's [documentation](https://docs.cert-manager.io/en/latest/getting-started/install/kubernetes.html)
on how to install cert-manager onto your cluster.

Once installed, verify all the components are running using:

```bash
kubectl get all -n cert-manager
NAME                                           READY   STATUS    RESTARTS   AGE
pod/cert-manager-86478c5ff-mkhb9               1/1     Running   0          23m
pod/cert-manager-cainjector-65dbccb8b6-6dnjl   1/1     Running   0          23m
pod/cert-manager-webhook-78f9d55fdf-5wcnp      1/1     Running   0          23m

NAME                           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/cert-manager-webhook   ClusterIP   10.63.240.251   <none>        443/TCP   23m

NAME                                      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/cert-manager              1         1         1            1           23m
deployment.apps/cert-manager-cainjector   1         1         1            1           23m
deployment.apps/cert-manager-webhook      1         1         1            1           23m

NAME                                                 DESIRED   CURRENT   READY   AGE
replicaset.apps/cert-manager-86478c5ff               1         1         1       23m
replicaset.apps/cert-manager-cainjector-65dbccb8b6   1         1         1       23m
replicaset.apps/cert-manager-webhook-78f9d55fdf      1         1         1       23m
```

## Setup your application

Any HTTP-based application can be used, for the purpose of the demo, install
the following echo server:

```bash
$ kubectl apply -f https://bit.ly/echo-service
service/echo created
deployment.apps/echo created
```

## Setup DNS

Get the IP address of the load balancer for Kong:

```bash
$ kubectl get service -n kong kong-proxy
NAME         TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                      AGE
kong-proxy   LoadBalancer   10.63.250.199   35.233.170.67   80:31929/TCP,443:31408/TCP   58d
```

To get only the IP address:

```bash
$ kubectl get -o jsonpath="{.status.loadBalancer.ingress[0].ip}" service -n kong kong-proxy
35.233.170.67
```

Please note that the IP address in your case will be different.

Next, setup a DNS records to resolve `proxy.yolo42.com` to the
above IP address:

```bash
$ dig +short proxy.yolo42.com
35.233.170.67
```

Next, setup a CNAME DNS record to resolve `demo.yolo42.com` to
`proxy.yolo42.com`.

```bash
$ dig +short demo.yolo2.com
proxy.yolo42.com.
35.233.170.67
```

## Expose your application to the Internet

Setup an Ingress rule to expose the application:

```bash
$ echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-yolo42-com
  annotations:
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - host: demo.yolo42.com
    http:
      paths:
      - path: /
        backend:
          serviceName: echo
          servicePort: 80
" | kubectl apply -f -
ingress.extensions/demo-yolo42-com created
```

Access your application:

```bash
$ curl -I demo.yolo42.com
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Fri, 21 Jun 2019 21:14:45 GMT
Server: echoserver
X-Kong-Upstream-Latency: 1
X-Kong-Proxy-Latency: 1
Via: kong/1.1.2
```

## Request TLS Certificate from Let's Encrypt

First, setup a ClusterIssuer for cert-manager

```bash
$ echo "apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
  namespace: cert-manager
spec:
  acme:
    email: user@example.com #please change this
    privateKeySecretRef:
      name: letsencrypt-prod
    server: https://acme-v02.api.letsencrypt.org/directory
    solvers:
    - http01:
        ingress: {}" | kubectl apply -f -
clusterissuer.cert-manager.io/letsencrypt-prod configured
```

*Note*: If you run into issues configuring this,
be sure that the group (`cert-manager.io`) and
version (`v1alpha2`) match those in the output of
`kubectl describe crd clusterissuer`.
This directs cert-manager which CA authority to use to issue the certificate.

Next, update your Ingress resource to provision a certificate and then use it:

```bash
$ echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-yolo42-com
  annotations:
    kubernetes.io/tls-acme: "true"
    cert-manager.io/cluster-issuer: letsencrypt-prod
    kubernetes.io/ingress.class: kong
spec:
  tls:
  - secretName: demo-yolo42-com
    hosts:
    - demo.yolo42.com
  rules:
  - host: demo.yolo42.com
    http:
      paths:
      - path: /
        backend:
          serviceName: echo
          servicePort: 80
' | kubectl apply -f -
ingress.extensions/demo-yolo42-com configured
```

Things to note here:

- The annotation `kubernetes.io/tls-acme`  is set to `true`, informing
  cert-manager that it should provision a certificate for hosts in this
  Ingress using ACME protocol.
- `certmanager.k8s.io/cluster-issuer` is set to `letsencrypt-prod`, directing
  cert-manager to use Let's Encrypt's production server to provision a TLS
  certificate.
- `tls` section of the Ingress directs Kong Ingress Controller to use the
  secret `demo-yolo42-com` to encrypt the traffic for `demo.yolo42.com`.
  This secret will be created by cert-manager.

Once you update the Ingress resource, cert-manager will start provisioning
the certificate and in sometime the certificate will be available for use.

You can track the progress of certificate issuance:

```bash
$ kubectl describe certificate demo-example-com
Name:         demo-example-com
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  certmanager.k8s.io/v1alpha1
Kind:         Certificate
Metadata:
  Creation Timestamp:  2019-06-21T20:41:54Z
  Generation:          1
  Owner References:
    API Version:           extensions/v1beta1
    Block Owner Deletion:  true
    Controller:            true
    Kind:                  Ingress
    Name:                  demo-example-com
    UID:                   261d15d3-9464-11e9-9965-42010a8a01ad
  Resource Version:        19561898
  Self Link:               /apis/certmanager.k8s.io/v1alpha1/namespaces/default/certificates/demo-example-com
  UID:                     014d3f1d-9465-11e9-9965-42010a8a01ad
Spec:
  Acme:
    Config:
      Domains:
        demo.yolo42.com
      Http 01:
  Dns Names:
    demo.yolo42.com
  Issuer Ref:
    Kind:       ClusterIssuer
    Name:       letsencrypt-prod
  Secret Name:  demo-example-com
Status:
  Conditions:
    Last Transition Time:  2019-06-21T20:42:20Z
    Message:               Certificate is up to date and has not expired
    Reason:                Ready
    Status:                True
    Type:                  Ready
  Not After:               2019-09-19T19:42:19Z
Events:
  Type    Reason              Age   From          Message
  ----    ------              ----  ----          -------
  Normal  Generated           53m   cert-manager  Generated new private key
  Normal  GenerateSelfSigned  53m   cert-manager  Generated temporary self signed certificate
  Normal  OrderCreated        53m   cert-manager  Created Order resource "demo-example-com-3811625818"
  Normal  OrderComplete       53m   cert-manager  Order "demo-example-com-3811625818" completed successfully
  Normal  CertIssued          53m   cert-manager  Certificate issued successfully
```

## Test HTTPS

Once all is in place, you can use HTTPS:

```bash
$ curl -v https://demo.yolo42.com
* Rebuilt URL to: https://demo.yolo42.com/
*   Trying 35.233.170.67...
* TCP_NODELAY set
* Connected to demo.yolo42.com (35.233.170.67) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* Cipher selection: ALL:!EXPORT:!EXPORT40:!EXPORT56:!aNULL:!LOW:!RC4:@STRENGTH
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/cert.pem
  CApath: none
* TLSv1.2 (OUT), TLS handshake, Client hello (1):
* TLSv1.2 (IN), TLS handshake, Server hello (2):
* TLSv1.2 (IN), TLS handshake, Certificate (11):
* TLSv1.2 (IN), TLS handshake, Server key exchange (12):
* TLSv1.2 (IN), TLS handshake, Server finished (14):
* TLSv1.2 (OUT), TLS handshake, Client key exchange (16):
* TLSv1.2 (OUT), TLS change cipher, Client hello (1):
* TLSv1.2 (OUT), TLS handshake, Finished (20):
* TLSv1.2 (IN), TLS change cipher, Client hello (1):
* TLSv1.2 (IN), TLS handshake, Finished (20):
* SSL connection using TLSv1.2 / ECDHE-RSA-AES256-GCM-SHA384
* ALPN, server accepted to use http/1.1
* Server certificate:
*  subject: CN=demo.yolo42.com
*  start date: Jun 21 19:42:19 2019 GMT
*  expire date: Sep 19 19:42:19 2019 GMT
*  subjectAltName: host "demo.yolo42.com" matched cert's "demo.yolo42.com"
*  issuer: C=US; O=Let's Encrypt; CN=Let's Encrypt Authority X3
*  SSL certificate verify ok.
> GET / HTTP/1.1
> Host: demo.yolo42.com
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=UTF-8
< Transfer-Encoding: chunked
< Connection: keep-alive
< Date: Fri, 21 Jun 2019 21:37:43 GMT
< Server: echoserver
< X-Kong-Upstream-Latency: 1
< X-Kong-Proxy-Latency: 1
< Via: kong/1.1.2
<


Hostname: echo-d778ffcd8-52ddj

Pod Information:
  node name: gke-harry-k8s-dev-default-pool-bb23a167-9w4t
  pod name: echo-d778ffcd8-52ddj
  pod namespace: default
  pod IP:10.60.2.246

Server values:
  server_version=nginx: 1.12.2 - lua: 10010

Request Information:
  client_address=10.60.2.239
  method=GET
  real path=/
  query=
  request_version=1.1
  request_scheme=http
  request_uri=http://demo.yolo42.com:8080/

Request Headers:
  accept=*/*
  connection=keep-alive
  host=demo.yolo42.com
  user-agent=curl/7.54.0
  x-forwarded-for=10.138.0.6
  x-forwarded-host=demo.yolo42.com
  x-forwarded-port=8443
  x-forwarded-proto=https
  x-real-ip=10.138.0.6

Request Body:
  -no body in request-
```

Et voilà ! You've secured your API with HTTPS
with Kong Ingress Controller and cert-manager.
