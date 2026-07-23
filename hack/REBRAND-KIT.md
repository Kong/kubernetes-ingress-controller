# CRD Rebrand Kit — kubernetes-ingress-controller @ v3.2.4

Move Kong's CRD API groups off `konghq.com` onto the private `acceldata.io`
domain so this KIC build coexists with an upstream Kong install in a shared
cluster. Behavior is unchanged — only the group moves.

## Why single-repo (important)

KIC **v3.2.4 predates the `kubernetes-configuration` module split**. The CRD Go
types live **in this repo** under `pkg/apis/{configuration,incubator}/**`, so the
rebrand is entirely local — **there is no `go.mod` repoint** and the
`kong-kubernetes-configuration` fork is not involved at this version. (Newer KIC
— v3.3+ — imports that module and needs the two-repo strategy instead.)

The group is defined in one place per version: the `+groupName=` marker + the
matching `GroupVersion`/`SchemeGroupVersion` in `pkg/apis/**/groupversion_info.go`.
Everything else (generated CRD YAML, RBAC, webhook rules, clientsets, owner-refs)
derives from it.

## Groups moved (both at this version)

| Old | New |
|---|---|
| `configuration.konghq.com` | `configuration.acceldata.io` |
| `incubator.ingress-controller.konghq.com` | `incubator.ingress-controller.acceldata.io` |

## Files

- `hack/rebrand-core.sh` — rebrand engine (dry-run / apply+branch / verify). Common, never edit per-repo.
- `hack/rebrand-test.sh` — local-Kubernetes CRD correctness suite (T1–T6).
- `hack/rebrand.conf` — this repo's two group pairs + CRD dir.

## Run it

```bash
hack/rebrand-core.sh --dry-run                              # preview every edit + CRD rename
hack/rebrand-core.sh --new-domain acceldata.io --branch     # apply on feature/rebrand-<branch>, commit
make manifests                                              # regenerate CRDs/RBAC/webhook from rebranded types
make generate.clientsets 2>/dev/null || true                # regenerate typed clients
git diff --stat                                             # expect only fold/format noise, if any
hack/rebrand-core.sh --verify                               # CI gate: no old group token remains
go build ./...                                              # must pass
```

## Not touched (by design)

`konghq.com/*` data-plane annotations, `konghq.com/kic-gateway-controller`
(GatewayClass controller identity), `*.validation.ingress-controller.konghq.com`
webhook names, `*.konghq.com` leader-election id, `kubernetes.io/ingress.class`,
and the Go module path. These are routing/identity, not CRD groups.

The group-as-key strings in `internal/dataplane/fallback/cache_to_graph_test.go`
(e.g. `configuration.konghq.com/KongPlugin:...`) ARE rebranded on purpose — they
are runtime-derived GroupKind cache keys that change with the group; leaving them
would break the tests.
