#!/usr/bin/env bash
#
# rebrand-test.sh — GENERIC local-Kubernetes correctness suite for a CRD rebrand.
#
# COMMON file — identical across operators. It proves the CRD-level guarantees
# that make a shared-cluster install safe, using only kubectl + a reachable
# cluster (kind/minikube/etc.). No image build or operator deploy required — that
# operator-specific "does it reconcile" step is driven by the skill separately.
#
# WHAT IT CHECKS (per old=new group pair from rebrand.conf)
#   T1 CRD installs and reaches Established
#   T2 Coexistence — the new CRD and the (upstream) old CRD can both exist
#   T3 Discovery — `kubectl api-resources` lists the new group distinctly
#   T4 Admission — a sample CR under the new apiVersion is accepted & stored
#   T5 Isolation — that CR is NOT visible under the old group
#   T6 Schema parity — new CRD == old CRD after group normalization (best-effort)
#
# USAGE
#   rebrand-test.sh [options]
#     --config PATH     rebrand.conf (default: ./rebrand.conf or hack/rebrand.conf)
#     --dir PATH        Repo root (default: git top-level, else .)
#     --namespace NS    Test namespace (default: rebrand-test)
#     --install-old     Also install the OLD CRD (from git history) to prove coexistence
#     --sample PATH     Sample CR manifest to use for admission (else auto-discovered)
#     --cleanup         Remove everything this script created, then exit
#     -h, --help        Show help
#
set -euo pipefail

CONFIG=""; DIR=""; NS="rebrand-test"; INSTALL_OLD=0; SAMPLE=""; CLEANUP=0
REBRAND_PAIRS=(); CRD_DIRS=()

log()  { printf '  %s\n' "$*"; }
info() { printf '\033[1m==>\033[0m %s\n' "$*"; }
pass() { printf '  \033[32m✓\033[0m %s\n' "$*"; }
failx(){ printf '  \033[31m✗\033[0m %s\n' "$*"; FAILED=$((FAILED+1)); }
skip() { printf '  \033[33m∼\033[0m %s (skipped: %s)\n' "$1" "$2"; }
warn() { printf '\033[33mwarning:\033[0m %s\n' "$*" >&2; }
die()  { printf '\033[31merror:\033[0m %s\n' "$*" >&2; exit 1; }
usage(){ sed -n '2,/^set -euo/p' "$0" | sed 's/^# \{0,1\}//; s/^#$//' | sed '$d'; }
FAILED=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --config)    CONFIG="${2:?}"; shift 2 ;;
    --dir)       DIR="${2:?}"; shift 2 ;;
    --namespace) NS="${2:?}"; shift 2 ;;
    --install-old) INSTALL_OLD=1; shift ;;
    --sample)    SAMPLE="${2:?}"; shift 2 ;;
    --cleanup)   CLEANUP=1; shift ;;
    -h|--help)   usage; exit 0 ;;
    *) die "unknown option: $1" ;;
  esac
done

command -v kubectl >/dev/null 2>&1 || die "kubectl is required"
[[ -z "$DIR" ]] && DIR="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$DIR"

# ---- config -------------------------------------------------------------------
if [[ -z "$CONFIG" ]]; then
  for c in rebrand.conf hack/rebrand.conf .rebrand.conf; do [[ -f "$c" ]] && { CONFIG="$c"; break; }; done
fi
[[ -n "$CONFIG" && -f "$CONFIG" ]] && { source "$CONFIG"; info "loaded config: $CONFIG"; }

if [[ ${#CRD_DIRS[@]} -eq 0 ]]; then
  for d in config/crd/bases deploy/crd/bases config/crd deploy/crd; do [[ -d "$d" ]] && CRD_DIRS+=("$d"); done
fi

# Resolve OLD/NEW group pairs. Prefer config; else derive NEW from installed CRD
# files and leave OLD empty (subset of tests).
OLDS=(); NEWS=()
if [[ ${#REBRAND_PAIRS[@]} -gt 0 ]]; then
  for p in "${REBRAND_PAIRS[@]}"; do OLDS+=("${p%%=*}"); NEWS+=("${p#*=}"); done
else
  warn "no rebrand.conf pairs found — deriving NEW groups from CRD manifests (T2/T5/T6 need old group; will be limited)"
  for d in ${CRD_DIRS[@]+"${CRD_DIRS[@]}"}; do
    for f in "$d"/*.yaml; do
      [[ -f "$f" ]] || continue
      g="$(grep -m1 -E '^  group: ' "$f" | awk '{print $2}')"
      [[ -n "$g" ]] && { OLDS+=(""); NEWS+=("$g"); }
    done
  done
fi
[[ ${#NEWS[@]} -gt 0 ]] || die "could not determine any CRD group to test"

# Kubernetes context guard (avoid accidentally hitting a prod cluster).
CTX="$(kubectl config current-context 2>/dev/null || echo '?')"
info "kube-context: $CTX   namespace: $NS"
case "$CTX" in
  kind-*|minikube|docker-desktop|k3d-*|rancher-desktop|orbstack) : ;;
  *) warn "context '$CTX' does not look like a local cluster — Ctrl-C now if this is not intended"; sleep 3 2>/dev/null || true ;;
esac

# Find the current (rebranded) CRD manifest file for a given new group.
crd_file_for_new() {  # new_group
  local ng="$1" d
  for d in ${CRD_DIRS[@]+"${CRD_DIRS[@]}"}; do
    find "$d" -maxdepth 1 -type f -name "${ng}_*.yaml" 2>/dev/null
  done | head -1
}

# Most recent commit:path whose tree still contained the OLD-group CRD manifest.
# Walks recent history (robust to git recording the rebrand as a rename, which
# --diff-filter=D would miss). Prints "<commit>:<path>" or nothing.
old_crd_ref() {  # old_group
  local og="$1" c p
  for c in $(git rev-list -n 25 HEAD 2>/dev/null); do
    p="$(git ls-tree -r --name-only "$c" 2>/dev/null | grep -m1 "${og}_[^/]*\.yaml$" || true)"
    [ -n "$p" ] && { printf '%s:%s' "$c" "$p"; return 0; }
  done
  return 1
}

# ---- cleanup ------------------------------------------------------------------
do_cleanup() {
  info "cleanup"
  kubectl delete ns "$NS" --ignore-not-found --wait=false >/dev/null 2>&1 || true
  for i in "${!NEWS[@]}"; do
    ng="${NEWS[$i]}"
    for crd in $(kubectl get crd -o name 2>/dev/null | grep "\.${ng}$" || true); do
      kubectl delete "$crd" --ignore-not-found >/dev/null 2>&1 || true
      log "deleted $crd"
    done
  done
  info "cleanup done (old/upstream CRDs left untouched)"
}
if [[ "$CLEANUP" == "1" ]]; then do_cleanup; exit 0; fi

# Ensure a usable namespace — wait out any lingering Terminating state first.
if kubectl get ns "$NS" -o jsonpath='{.status.phase}' 2>/dev/null | grep -q Terminating; then
  kubectl wait --for=delete ns/"$NS" --timeout=60s >/dev/null 2>&1 || true
fi
kubectl get ns "$NS" >/dev/null 2>&1 || kubectl create ns "$NS" >/dev/null
kubectl label ns "$NS" purpose=rebrand-test --overwrite >/dev/null 2>&1 || true

# ==============================================================================
for i in "${!NEWS[@]}"; do
  NG="${NEWS[$i]}"; OG="${OLDS[$i]}"
  echo; info "group: ${OG:+$OG -> }$NG"
  CRDF="$(crd_file_for_new "$NG")"
  [[ -n "$CRDF" ]] || { failx "no rebranded CRD manifest found for group '$NG' under ${CRD_DIRS[*]}"; continue; }

  # T1 install + Established
  kubectl apply -f "$CRDF" >/dev/null 2>&1 || { failx "T1 apply CRD ($CRDF)"; continue; }
  crd_name="$(grep -m1 -E '^  name: ' "$CRDF" | awk '{print $2}')"
  ok=0; for _ in 1 2 3 4 5 6 7 8 9 10; do
    [[ "$(kubectl get crd "$crd_name" -o jsonpath='{.status.conditions[?(@.type=="Established")].status}' 2>/dev/null)" == "True" ]] && { ok=1; break; }
    sleep 1 2>/dev/null || true
  done
  [[ "$ok" == "1" ]] && pass "T1 CRD Established: $crd_name" || failx "T1 CRD not Established: $crd_name"

  # T3 discovery
  if kubectl api-resources --api-group="$NG" 2>/dev/null | grep -q .; then
    pass "T3 discovery lists group '$NG'"
  else failx "T3 group '$NG' not in api-resources"; fi

  # T2 coexistence (needs old group)
  if [[ -n "$OG" ]]; then
    old_crd="${crd_name/$NG/$OG}"
    if [[ "$INSTALL_OLD" == "1" ]] && ! kubectl get crd "$old_crd" >/dev/null 2>&1; then
      # Recreate the pre-rebrand CRD from git history to demonstrate coexistence.
      ref="$(old_crd_ref "$OG" || true)"
      [[ -n "$ref" ]] && git show "$ref" 2>/dev/null | kubectl apply -f - >/dev/null 2>&1 || true
    fi
    if kubectl get crd "$old_crd" >/dev/null 2>&1 && kubectl get crd "$crd_name" >/dev/null 2>&1; then
      pass "T2 coexistence: both '$crd_name' and '$old_crd' present"
    else
      skip "T2 coexistence" "old CRD '$old_crd' not in cluster (use --install-old)"
    fi
  else
    skip "T2 coexistence" "old group unknown (no rebrand.conf)"
  fi

  # T4/T5 admission + isolation (needs a sample CR)
  s="$SAMPLE"
  if [[ -z "$s" ]]; then
    for cand in config/samples deploy/examples examples; do
      [[ -d "$cand" ]] || continue
      s="$(grep -rl -m1 "apiVersion:.*${NG}/" "$cand" 2>/dev/null | head -1)"; [[ -n "$s" ]] && break
    done
  fi
  if [[ -n "$s" && -f "$s" ]]; then
    if kubectl -n "$NS" apply -f "$s" >/dev/null 2>&1; then
      kind_plural="$(kubectl get crd "$crd_name" -o jsonpath='{.spec.names.plural}' 2>/dev/null)"
      if [[ -n "$(kubectl -n "$NS" get "${kind_plural}.${NG}" -o name 2>/dev/null)" ]]; then
        pass "T4 admission: sample CR stored under '$NG' ($(basename "$s"))"
      else failx "T4 sample CR not found under '$NG'"; fi
      if [[ -n "$OG" ]] && kubectl get crd "${crd_name/$NG/$OG}" >/dev/null 2>&1; then
        if [[ -z "$(kubectl -n "$NS" get "${kind_plural}.${OG}" -o name 2>/dev/null)" ]]; then
          pass "T5 isolation: CR invisible under old group '$OG'"
        else failx "T5 CR leaked into old group '$OG'"; fi
      else skip "T5 isolation" "old CRD not present"; fi
    else failx "T4 sample CR rejected by apiserver ($s)"; fi
  else
    skip "T4/T5 admission" "no sample CR found (pass --sample PATH)"
  fi

  # T6 schema parity (best-effort, from git)
  if [[ -n "$OG" ]]; then
    ref="$(old_crd_ref "$OG" || true)"
    if [[ -n "$ref" ]]; then
      tmp_old="$(mktemp)"; tmp_new="$(mktemp)"
      git show "$ref" 2>/dev/null | sed "s/$OG/GROUP/g" > "$tmp_old" || true
      sed "s/$NG/GROUP/g" "$CRDF" > "$tmp_new"
      if [[ -s "$tmp_old" ]] && diff -q "$tmp_old" "$tmp_new" >/dev/null 2>&1; then
        pass "T6 schema parity: identical after group normalization"
      elif [[ -s "$tmp_old" ]]; then failx "T6 schema differs beyond the group token"; else skip "T6 parity" "old manifest unreadable"; fi
      rm -f "$tmp_old" "$tmp_new"
    else skip "T6 parity" "old manifest not found in git history"; fi
  else skip "T6 parity" "old group unknown"; fi
done

echo
if [[ "$FAILED" -eq 0 ]]; then
  info "ALL CHECKS PASSED"
  echo "  (operator reconciliation is a separate, operator-specific step — see the skill)"
  echo "  cleanup with: $0 --cleanup --namespace $NS"
else
  die "$FAILED check(s) failed"
fi
