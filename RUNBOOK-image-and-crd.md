# Runbook — Rebranded KIC v3.2.4 (image + CRDs)

Build and ship a private Kong Ingress Controller v3.2.4 whose CRDs live under
`acceldata.io` instead of `konghq.com`, so it coexists with an upstream Kong
install in a shared cluster.

## What changed

Both CRD API groups moved (single-repo — types are in `pkg/apis/`):

| Old group | New group |
|---|---|
| `configuration.konghq.com` | `configuration.acceldata.io` |
| `incubator.ingress-controller.konghq.com` | `incubator.ingress-controller.acceldata.io` |

CRDs: `config/crd/bases/<new-group>_*.yaml`; all-in-one manifests: `deploy/single/*.yaml`.
Unchanged: GatewayClass controller name `konghq.com/kic-gateway-controller`,
data-plane annotations `konghq.com/*`, ingress class, webhook names, module path.

## 1. Reproduce the rebrand

```bash
hack/rebrand-core.sh --dry-run
hack/rebrand-core.sh --new-domain acceldata.io --branch
make manifests            # regenerate CRDs / RBAC / webhook / deploy-single
hack/rebrand-core.sh --verify
go build ./...
```

## 2. Build a linux/amd64 image (podman; host is arm64)

```bash
TAG=v3.2.4-acceldata-rebrand
IMAGE=localhost/kong-ingress-acceldata
podman build --platform linux/amd64 --build-arg TAG=${TAG} -t ${IMAGE}:${TAG} .
```

## 3. Load into kind

```bash
CLUSTER=kind-cluster            # your local cluster
podman save ${IMAGE}:${TAG} -o /tmp/kic-rebrand.tar
kind load image-archive /tmp/kic-rebrand.tar --name ${CLUSTER#kind-}
```

> ARCH NOTE: the shippable image is **linux/amd64**. A kind node built on an
> Apple-silicon host is **arm64** and will not run an amd64 image without qemu.
> For a local arm64 kind test, rebuild with `--platform linux/arm64` (native,
> fast) and load that tag instead. The amd64 image remains the delivery artifact.

## 4. Deploy under a DISTINCT release name (avoid ClusterRole clashes)

The in-repo `deploy/single/*` all-in-one manifests are **deprecated stubs** at
v3.2.4 — do not use them. Deploy the rebranded RBAC + a KIC/Kong pod directly.
KIC v3.2.4 requires a reachable Kong admin API at startup, so run a dbless Kong
next to it. Give every **cluster-scoped** object a distinct name so it cannot
collide with an upstream KIC install (CRDs already differ by group; ClusterRole,
ClusterRoleBinding must be renamed):

```bash
# 1. Rebranded RBAC, renamed distinct. Both roles are needed:
#    role.yaml (Kong groups + core) AND crds/role.yaml (apiextensions read).
sed 's/name: kong-ingress$/name: acceldata-kong-ingress/'      config/rbac/role.yaml      | kubectl apply -f -
sed 's/name: kong-ingress-crds/name: acceldata-kong-ingress-crds/' config/rbac/crds/role.yaml | kubectl apply -f -
kubectl create ns kong-acceldata
kubectl create sa acceldata-kong-sa -n kong-acceldata
kubectl create clusterrolebinding acceldata-kong-ingress      --clusterrole=acceldata-kong-ingress      --serviceaccount=kong-acceldata:acceldata-kong-sa
kubectl create clusterrolebinding acceldata-kong-ingress-crds --clusterrole=acceldata-kong-ingress-crds --serviceaccount=kong-acceldata:acceldata-kong-sa

# 2. Install the rebranded CRDs (both dirs):
kubectl apply -f config/crd/bases/
kubectl apply -f config/crd/incubator/incubator.ingress-controller.acceldata.io_kongservicefacades.yaml

# 3. Deploy a pod with two containers — proxy (kong:3.6, dbless) + ingress-controller
#    (your image). KIC args: --kong-admin-url=http://localhost:8001
#    --publish-service=kong-acceldata/acceldata-kong-proxy --anonymous-reports=false
#    (Do NOT install the ValidatingWebhookConfiguration unless you wire its TLS cert;
#     with failurePolicy=Fail it would block Kong CR writes.)
```

A ready-to-apply reference manifest is in the rebrand PR description / kit notes.

## 5. Confirm reconciliation + isolation (verified 2026-07-23 on kind v1.35, arm64)

```bash
# Coexistence — new + pre-existing upstream CRDs both present:
kubectl get crd | grep -E 'kongconsumers\.(configuration\.acceldata\.io|configuration\.konghq\.com)'

# Reconcile a CR under the NEW group -> KIC pushes it into Kong's data plane:
kubectl -n kong-acceldata apply -f - <<'EOF'
apiVersion: configuration.acceldata.io/v1
kind: KongConsumer
metadata:
  name: alice
  namespace: kong-acceldata
  annotations: { kubernetes.io/ingress.class: kong }
username: alice-acceldata
EOF
kubectl -n kong-acceldata get kongconsumer alice -o jsonpath='{.status.conditions[0].reason}'   # Programmed
kubectl -n kong-acceldata port-forward deploy/acceldata-kong 8001:8001 &
curl -s localhost:8001/consumers    # -> alice-acceldata present

# RBAC grants new group, not old:
SA=system:serviceaccount:kong-acceldata:acceldata-kong-sa
kubectl auth can-i list kongplugins.configuration.acceldata.io --as=$SA   # yes
kubectl auth can-i list kongplugins.configuration.konghq.com   --as=$SA   # no
```

GOTCHA: dbless sync is all-or-nothing — a single invalid Kong CR anywhere in the
cluster (e.g. a KongConsumer referencing a missing KongConsumerGroup) blocks the
whole push. Keep the test namespace clean.

## Rollback / cutover

Self-contained commit on `feature/rebrand-v3.2.4`; `git checkout v3.2.4` to abandon.
No CR migration: old-group CRs are a distinct kind with no conversion webhook — a
cutover re-creates Kong CRs under the new `apiVersion`.
