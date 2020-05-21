# Rewriting paths
This guide demonstrates host and path rewrites using Ingress and Service configuration.

## Installation

Please follow the [deployment](../deployment) documentation to install
Kong Ingress Controller on your Kubernetes cluster.

## Testing Connectivity to Kong

This guide assumes that the `PROXY_IP` environment variable is
set to contain the IP address or URL pointing to Kong.
Please follow one of the
[deployment guides](../deployment) to configure this environment variable.

If everything is setup correctly, making a request to Kong should return
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

## Create a Kubernetes service

First, create a Kubernetes Service:

```bash
echo "
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  ports:
    - name: http
      protocol: TCP
      port: 80
      targetPort: 9376
    - name: https
      protocol: TCP
      port: 443
      targetPort: 9377
" | kubectl create -f -
```

This Service will create a Kong service and upstream that uses the upstream IPs
(Pod IPs) for its `Host` header and appends request paths starting at `/`.

## Create an Ingress to expose the service at the path `/myapp` on `example.com`

```bash
echo '
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: my-app
spec:
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /myapp
        backend:
          serviceName: my-service
          servicePort: 80
' | kubectl create -f -
```

This Ingress will create a Kong route attached to the service we created above.
It will preserve its path but honor the service's hostname, so this request:

```
$ curl -svX GET http://myapp.example.com/myapp/foo
GET /myapp/foo HTTP/1.1
Host: myapp.example.com
User-Agent: curl/7.70.0
Accept: */*
```
will appear upstream as:

```
GET /myapp/foo HTTP/1.1
Host: 10.16.4.8
User-Agent: curl/7.70.0
Accept: */*
```

Note that this default behavior uses `strip_path=false` on the route. This
differs from Kong's standard default to conform with expected ingress
controller behavior.

## Rewriting the host

There are two options to override the default `Host` header behavior:

- Add the [`konghq.com/preserve-host` annotation][0] to your Ingress, which
  sends the route/Ingress hostname:
  ```bash
  $ kubectl patch ingress my-app -p '{"metadata":{"annotations":{"konghq.com/preserve-host":"true"}}}'
  ```
  The request upstream will now look like:
  ```
  GET /myapp/foo HTTP/1.1
  Host: myapp.example.com
  User-Agent: curl/7.70.0
  Accept: */*
  ```
- Add the [`konghq.com/host-header` annotation][1] to your Service, which sets
  the `Host` header directly:
  ```bash
  $ kubectl patch service my-service -p '{"metadata":{"annotations":{"konghq.com/host-header":"internal.myapp.example.com"}}}'
  ```
  The request upstream will now look like:
  ```
  GET /myapp/foo HTTP/1.1
  Host: internal.myapp.example.com
  User-Agent: curl/7.70.0
  Accept: */*
  ```

The `preserve-host` annotation takes precedence, so if you add both annotations
above, the upstream host header would be `myapp.example.com`.

## Rewriting the path

There are two options to rewrite the default path handling behavior:

- Add the [`konghq.com/strip-path` annotation][2] to your Ingress, which strips
  the path component of the route/Ingress, leaving the remainder of the path at
  the root:
  ```bash
  $ kubectl patch ingress my-app -p '{"metadata":{"annotations":{"konghq.com/strip-path":"true"}}}'
  ```
  The request upstream will now look like:
  ```
  GET /foo HTTP/1.1
  Host: myapp.example.com
  User-Agent: curl/7.70.0
  Accept: */*
  ```
- Add the [`konghq.com/path` annotation][3] to your Service, which prepends
  that value to the upstream path:
  ```bash
  $ kubectl patch service my-service -p '{"metadata":{"annotations":{"konghq.com/path":"/api"}}}'
  ```
  The request upstream will now look like:
  ```
  GET /api/myapp/foo HTTP/1.1
  Host: myapp.example.com
  User-Agent: curl/7.70.0
  Accept: */*
  ```
`strip-path` and `path` can be combined together, with the `path` component
coming first. Adding both annotations above will send requests for `/api/foo`.

[0]: https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/references/annotations.md#konghqcompreserve-host
[1]: https://github.com/Kong/kubernetes-ingress-controller/blob/next/docs/references/annotations.md#konghqcomhost-header
[2]: https://github.com/Kong/kubernetes-ingress-controller/blob/next/docs/references/annotations.md#konghqcomstrip-path
[3]: https://github.com/Kong/kubernetes-ingress-controller/blob/next/docs/references/annotations.md#konghqcompath
