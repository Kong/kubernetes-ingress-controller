# Annotation ingress.class

If you have multiple Ingress controllers in a single cluster,
you can pick one by specifying the `ingress.class` annotation.
Following is an example of
creating an Ingress with an annotation:

```yaml
metadata:
  name: foo
  annotations:
    kubernetes.io/ingress.class: "gce"
```

will target the GCE controller, forcing Kong Ingress Controller to ignore it.

On the other hand, an annotation such as

```yaml
metadata:
  name: foo
  annotations:
    kubernetes.io/ingress.class: "nginx"
```

will target Kong Ingress controller, forcing the GCE controller to ignore it.

__Note__: Deploying multiple ingress controller and not specifying the
annotation will cause both controllers fighting to satisfy the Ingress
and will lead to unknown behaviour.

If you're running multiple ingress controllers, or running on a cloud provider that handles ingress, you need to specify the annotation `kubernetes.io/ingress.class: "nginx"` in all ingresses you would like this controller to claim. This mechanism also provides users the ability to run _multiple_ Kong ingress controllers (e.g. one which serves public traffic, one which serves "internal" traffic).
When using this functionality the option `--ingress-class` should set a value unique for the cluster. Here is a partial example:

```yaml
spec:
  template:
     spec:
       containers:
         - name: kong-ingress-internal-controller
           args:
             - /kong-ingress-controller
             - '--election-id=ingress-controller-leader-internal'
             - '--ingress-class=kong-internal'
```

Not specifying the annotation will lead to multiple ingress controllers claiming the same ingress.
Setting a value which does not match the class of any existing ingress controllers will cause all ingress controllers ignoring the ingress.
