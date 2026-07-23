#!/usr/bin/env bash
#
# rebrand-core.sh — GENERIC CRD API-group rebrand engine for kubebuilder /
# controller-runtime operators.
#
# This file is COMMON: it is byte-for-byte identical across every operator you
# rebrand. Never edit it per-operator — everything operator-specific lives in a
# small `rebrand.conf` (see rebrand.conf.example), and anything not configured is
# auto-detected from the repo.
#
# WHY THIS WORKS FOR ANY SUCH OPERATOR
#   A CRD is a cluster-scoped singleton keyed by <plural>.<group>. To make your
#   build coexist with an upstream install in a shared cluster you only need to
#   move its API *group*. kubebuilder/controller-runtime operators encode the
#   group in exactly one source of truth — the `+groupName=` marker and the
#   matching `SchemeGroupVersion` in a groupversion_info.go / *_types.go — and
#   derive everything else (CRD manifest, RBAC, watches, owner-refs) from it.
#   So a fixed-string replace of the group token, plus renaming the generated
#   CRD manifest file, is sufficient AND safe. (Re-verify per operator — see
#   references/verification.md for the assumptions this relies on.)
#
# CONFIG (rebrand.conf, all optional — auto-detected when omitted)
#   REBRAND_PAIRS=("old.group=new.group" ...)   # what to rebrand; auto-detected
#                                                 # from +groupName if left empty,
#                                                 # with NEW = OLD's first label + NEW_DOMAIN
#   NEW_DOMAIN="acceldata.io"                     # used to synthesize NEW when only
#                                                 # the old group is auto-detected
#   CRD_DIRS=("config/crd/bases" ...)            # auto-detected when empty
#   KEEP=("some.domain/anno")                     # informational: things you chose to keep
#
# USAGE
#   rebrand-core.sh [options]
#     --config PATH      Config file (default: ./rebrand.conf or hack/rebrand.conf if present)
#     --dir PATH         Repo root (default: git top-level, else .)
#     --new-domain DOM   Synthesize NEW groups as <first-label-of-old>.<DOM> (e.g. acceldata.io)
#     --pair OLD=NEW     Add an explicit rebrand pair (repeatable; overrides auto-detect)
#     --crd-dir PATH     Add a CRD manifest dir (repeatable; overrides auto-detect)
#     --branch           Create feature/rebrand-<normalized_source_branch>, apply, commit there
#     --branch-name NAME Override the auto branch name (implies --branch)
#     --no-commit        With --branch: branch + apply but skip the commit
#     --dry-run          Show what would change; write nothing
#     --verify           Assert the tree is fully rebranded; change nothing (exit 1 on drift)
#     -h, --help         Show this help
#
# EXAMPLES
#   rebrand-core.sh --new-domain acceldata.io --dry-run     # auto-detect, preview
#   rebrand-core.sh --branch                                 # uses rebrand.conf, branch+commit
#   rebrand-core.sh --pair vault.banzaicloud.com=vault.acceldata.io --branch
#   rebrand-core.sh --verify                                 # CI gate
#
set -euo pipefail

# ------------------------------------------------------------------ defaults / args
CONFIG=""
DIR=""
NEW_DOMAIN="${NEW_DOMAIN:-}"
MODE="apply"        # apply | dry-run | verify
DO_BRANCH=0
BRANCH_NAME=""
DO_COMMIT=1
CLI_PAIRS=()
CLI_CRD_DIRS=()
REBRAND_PAIRS=()    # may be set by config
CRD_DIRS=()         # may be set by config
KEEP=()             # informational only

log()  { printf '  %s\n' "$*"; }
info() { printf '\033[1m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[33mwarning:\033[0m %s\n' "$*" >&2; }
die()  { printf '\033[31merror:\033[0m %s\n' "$*" >&2; exit 1; }
usage() { sed -n '2,/^set -euo/p' "$0" | sed 's/^# \{0,1\}//; s/^#$//' | sed '$d'; }
normalize() { printf '%s' "$1" | tr '[:upper:]' '[:lower:]' | sed -E 's#[^a-z0-9]+#-#g; s#^-+##; s#-+$##'; }

while [[ $# -gt 0 ]]; do
  case "$1" in
    --config)      CONFIG="${2:?}"; shift 2 ;;
    --dir)         DIR="${2:?}"; shift 2 ;;
    --new-domain)  NEW_DOMAIN="${2:?}"; shift 2 ;;
    --pair)        CLI_PAIRS+=("${2:?}"); shift 2 ;;
    --crd-dir)     CLI_CRD_DIRS+=("${2:?}"); shift 2 ;;
    --branch)      DO_BRANCH=1; shift ;;
    --branch-name) BRANCH_NAME="${2:?}"; DO_BRANCH=1; shift 2 ;;
    --no-commit)   DO_COMMIT=0; shift ;;
    --dry-run)     MODE="dry-run"; shift ;;
    --verify)      MODE="verify"; shift ;;
    -h|--help)     usage; exit 0 ;;
    *) die "unknown option: $1 (try --help)" ;;
  esac
done

command -v perl >/dev/null 2>&1 || die "perl is required but not found"

# ------------------------------------------------------------------ locate repo
[[ -z "$DIR" ]] && DIR="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
[[ -d "$DIR" ]] || die "directory not found: $DIR"
cd "$DIR"
IS_GIT=0; git rev-parse --is-inside-work-tree >/dev/null 2>&1 && IS_GIT=1

# ------------------------------------------------------------------ load config
if [[ -z "$CONFIG" ]]; then
  for c in rebrand.conf hack/rebrand.conf .rebrand.conf; do [[ -f "$c" ]] && { CONFIG="$c"; break; }; done
fi
if [[ -n "$CONFIG" ]]; then
  [[ -f "$CONFIG" ]] || die "config not found: $CONFIG"
  # shellcheck disable=SC1090
  source "$CONFIG"
  info "loaded config: $CONFIG"
fi
# CLI overrides win over config.
[[ ${#CLI_PAIRS[@]}    -gt 0 ]] && REBRAND_PAIRS=("${CLI_PAIRS[@]}")
[[ ${#CLI_CRD_DIRS[@]} -gt 0 ]] && CRD_DIRS=("${CLI_CRD_DIRS[@]}")

# ------------------------------------------------------------------ sanity: is this an operator repo?
looks_like_operator() {
  [[ -f PROJECT ]] && return 0
  grep -rql -- '+groupName=' --include='*.go' . 2>/dev/null && return 0
  for d in config/crd deploy/crd; do [[ -d "$d" ]] && return 0; done
  return 1
}
looks_like_operator || die "'$DIR' does not look like a kubebuilder/controller-runtime operator (no PROJECT, +groupName marker, or */crd dir)"

# ------------------------------------------------------------------ auto-detect CRD dirs
if [[ ${#CRD_DIRS[@]} -eq 0 ]]; then
  for d in config/crd/bases deploy/crd/bases config/crd deploy/crd; do
    [[ -d "$d" ]] && CRD_DIRS+=("$d")
  done
fi

# ------------------------------------------------------------------ auto-detect rebrand pairs
detect_old_groups() {
  # from +groupName markers (authoritative source of truth)
  grep -rhoE '\+groupName=[A-Za-z0-9.-]+' --include='*.go' . 2>/dev/null \
    | sed 's/^+groupName=//' | sort -u
}
if [[ ${#REBRAND_PAIRS[@]} -eq 0 ]]; then
  [[ -n "$NEW_DOMAIN" ]] || die "no REBRAND_PAIRS configured and no --new-domain given.
Either create a rebrand.conf with REBRAND_PAIRS=(\"old=new\"), pass --pair old=new,
or pass --new-domain <domain> to auto-synthesize NEW from the detected group(s)."
  olds="$(detect_old_groups)"
  [[ -n "$olds" ]] || die "could not auto-detect any API group (+groupName markers not found); use --pair or rebrand.conf"
  while IFS= read -r og; do
    [[ -n "$og" ]] || continue
    first_label="${og%%.*}"                 # e.g. vault.banzaicloud.com -> vault
    REBRAND_PAIRS+=("${og}=${first_label}.${NEW_DOMAIN}")
  done <<< "$olds"
fi
[[ ${#REBRAND_PAIRS[@]} -gt 0 ]] || die "no rebrand pairs resolved"

# Parse pairs into OLD/NEW arrays; validate.
OLDS=(); NEWS=()
dns_re='^[a-z0-9]([-a-z0-9.]*[a-z0-9])?$'
for p in "${REBRAND_PAIRS[@]}"; do
  [[ "$p" == *=* ]] || die "invalid rebrand pair '$p' (expected old.group=new.group)"
  o="${p%%=*}"; n="${p#*=}"
  [[ "$o" =~ $dns_re ]] || die "invalid old group '$o' (must be a DNS subdomain)"
  [[ "$n" =~ $dns_re ]] || die "invalid new group '$n' (must be a DNS subdomain)"
  [[ "$o" != "$n" ]] || die "old and new group are identical: '$o'"
  OLDS+=("$o"); NEWS+=("$n")
done

info "operator CRD rebrand"
log  "repo      : $DIR"
log  "CRD dirs  : ${CRD_DIRS[*]:-(none found)}"
for i in "${!OLDS[@]}"; do log  "group     : ${OLDS[$i]}  ->  ${NEWS[$i]}"; done
[[ ${#KEEP[@]} -gt 0 ]] && log "kept      : ${KEEP[*]}"
log  "mode      : $MODE"
echo

# ------------------------------------------------------------------ helpers
# List text files containing a literal token (git-aware, binary/self-safe).
# The kit's own scripts, config, and docs legitimately contain the OLD group as
# examples/explanation — they must never be rewritten (that would corrupt them).
SELF_EXCLUDES=(
  hack/rebrand-core.sh hack/rebrand-test.sh hack/rebrand-pipeline.sh
  hack/rebrand.conf hack/rebrand.conf.example
  hack/REBRAND-KIT.md hack/REBRAND.md RUNBOOK-image-and-crd.md
  rebrand-core.sh rebrand-test.sh rebrand.conf rebrand.conf.example
)
list_files_with() {
  local token="$1"
  if [[ "$IS_GIT" == "1" ]]; then
    local excl=() p; for p in "${SELF_EXCLUDES[@]}"; do excl+=(":(exclude)$p"); done
    git grep -I -l -F -e "$token" -- . "${excl[@]}" 2>/dev/null || true
  else
    local ge=() p; for p in "${SELF_EXCLUDES[@]}"; do ge+=(--exclude="$(basename "$p")"); done
    grep -rIlF --exclude-dir=.git "${ge[@]}" -e "$token" . 2>/dev/null | sed 's#^\./##' || true
  fi
}
replace_in_files() {  # from to files...
  local from="$1" to="$2"; shift 2
  [[ $# -gt 0 ]] || return 0
  REBRAND_FROM="$from" REBRAND_TO="$to" perl -i -pe 's/\Q$ENV{REBRAND_FROM}\E/$ENV{REBRAND_TO}/g' "$@"
}
read_into() {  # array_name < stream
  local __n="$1" __l; eval "$__n=()"
  while IFS= read -r __l; do [[ -n "$__l" ]] && eval "$__n+=(\"\$__l\")"; done
}
# All CRD manifest files named <old>_*.yaml under the CRD dirs.
crd_files_for() {  # old_group
  local og="$1" d
  for d in ${CRD_DIRS[@]+"${CRD_DIRS[@]}"}; do
    [[ -d "$d" ]] || continue
    find "$d" -maxdepth 1 -type f -name "${og}_*.yaml" 2>/dev/null || true
  done
}

# ------------------------------------------------------------------ verify mode
if [[ "$MODE" == "verify" ]]; then
  rc=0
  for i in "${!OLDS[@]}"; do
    og="${OLDS[$i]}"; ng="${NEWS[$i]}"
    leftover="$(list_files_with "$og")"
    if [[ -n "$leftover" ]]; then warn "files still referencing '$og':"; printf '    %s\n' $leftover >&2; rc=1; fi
    named="$(crd_files_for "$og")"
    if [[ -n "$named" ]]; then warn "CRD manifest(s) still named for '$og':"; printf '    %s\n' $named >&2; rc=1; fi
    if ! grep -rql -F -- "+groupName=$ng" --include='*.go' . 2>/dev/null; then
      warn "no Go source of truth contains '+groupName=$ng'"; rc=1
    fi
  done
  [[ $rc -eq 0 ]] && { info "verify OK — fully rebranded"; exit 0; } || die "verify FAILED"
fi

# ------------------------------------------------------------------ plan
TOTAL_FILES=0
if [[ "$MODE" == "dry-run" ]]; then
  for i in "${!OLDS[@]}"; do
    og="${OLDS[$i]}"; ng="${NEWS[$i]}"
    files=(); read_into files < <(list_files_with "$og")
    info "'$og' -> '$ng' : ${#files[@]} file(s)"
    for f in ${files[@]+"${files[@]}"}; do
      printf '    %-58s (%s occ.)\n' "$f" "$(grep -c -F "$og" "$f" 2>/dev/null || echo 0)"
    done
    renames=(); read_into renames < <(crd_files_for "$og")
    for f in ${renames[@]+"${renames[@]}"}; do
      b="$(basename "$f")"; printf '    rename: %s -> %s\n' "$f" "$(dirname "$f")/${b/${og}_/${ng}_}"
    done
    TOTAL_FILES=$((TOTAL_FILES + ${#files[@]}))
  done
  info "dry-run complete — nothing written ($TOTAL_FILES file-changes across ${#OLDS[@]} group(s))"
  exit 0
fi

# ------------------------------------------------------------------ branch (apply only)
TARGET_BRANCH=""; SRC_BRANCH=""
if [[ "$DO_BRANCH" == "1" ]]; then
  [[ "$IS_GIT" == "1" ]] || die "--branch requires a git repository"
  SRC_BRANCH="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || true)"
  if [[ -n "$BRANCH_NAME" ]]; then
    TARGET_BRANCH="$BRANCH_NAME"
  else
    [[ -n "$SRC_BRANCH" && "$SRC_BRANCH" != "HEAD" ]] || die "--branch: detached HEAD; use --branch-name NAME"
    TARGET_BRANCH="feature/rebrand-$(normalize "$SRC_BRANCH")"
  fi
  if [[ "$SRC_BRANCH" == "$TARGET_BRANCH" ]]; then
    info "already on rebrand branch '$TARGET_BRANCH'"
  elif git show-ref --verify --quiet "refs/heads/$TARGET_BRANCH"; then
    die "branch '$TARGET_BRANCH' already exists — delete it or pass --branch-name"
  else
    info "creating rebrand branch '$TARGET_BRANCH' from '$SRC_BRANCH'"
    git checkout -b "$TARGET_BRANCH" >/dev/null
  fi
fi

# ------------------------------------------------------------------ persist config
# If no config existed, record the resolved pairs so --verify / rebrand-test / CI
# work later without re-passing --pair/--new-domain. Written into the commit.
if [[ -z "$CONFIG" ]]; then
  conf_out="rebrand.conf"; [[ -d hack ]] && conf_out="hack/rebrand.conf"
  {
    echo "# Generated by rebrand-core.sh — records this rebrand for --verify / rebrand-test / CI."
    printf 'REBRAND_PAIRS=('
    for i in "${!OLDS[@]}"; do printf '"%s=%s" ' "${OLDS[$i]}" "${NEWS[$i]}"; done
    printf ')\n'
    if [[ ${#CRD_DIRS[@]} -gt 0 ]]; then
      printf 'CRD_DIRS=('; for d in "${CRD_DIRS[@]}"; do printf '"%s" ' "$d"; done; printf ')\n'
    fi
  } > "$conf_out"
  CONFIG="$conf_out"
  info "wrote $conf_out (records the rebrand for verify/test/CI)"
fi

# ------------------------------------------------------------------ apply
for i in "${!OLDS[@]}"; do
  og="${OLDS[$i]}"; ng="${NEWS[$i]}"
  files=(); read_into files < <(list_files_with "$og")
  if [[ ${#files[@]} -gt 0 ]]; then
    info "rewriting '$og' -> '$ng' in ${#files[@]} file(s)"
    replace_in_files "$og" "$ng" ${files[@]+"${files[@]}"}
  fi
  # rename CRD manifests <og>_*.yaml -> <ng>_*.yaml
  renames=(); read_into renames < <(crd_files_for "$og")
  for f in ${renames[@]+"${renames[@]}"}; do
    d="$(dirname "$f")"; b="$(basename "$f")"; nf="$d/${b/${og}_/${ng}_}"
    info "renaming CRD manifest: $b -> $(basename "$nf")"
    if [[ "$IS_GIT" == "1" ]] && git ls-files --error-unmatch "$f" >/dev/null 2>&1; then
      git mv -f "$f" "$nf"
    else
      mv -f "$f" "$nf"
    fi
  done
done

# ------------------------------------------------------------------ post-check
echo
for i in "${!OLDS[@]}"; do
  og="${OLDS[$i]}"
  left="$(list_files_with "$og" | wc -l | tr -d ' ')"
  [[ "$left" == "0" ]] || { warn "$left file(s) still reference '$og':"; list_files_with "$og" | sed 's/^/    /' >&2; die "rebrand incomplete"; }
  named="$(crd_files_for "$og")"
  [[ -z "$named" ]] || { warn "CRD manifest(s) still named for '$og':"; printf '    %s\n' $named >&2; die "rebrand incomplete (file not renamed)"; }
done
info "rebrand complete"

# ------------------------------------------------------------------ commit (apply + --branch)
if [[ "$DO_BRANCH" == "1" && "$DO_COMMIT" == "1" ]]; then
  info "committing on '$TARGET_BRANCH'"
  git add -A -- . ':(exclude)*.swp' ':(exclude)**/*.swp'
  commit_id=()
  git config user.email >/dev/null 2>&1 || commit_id=(-c user.name="rebrand" -c user.email="rebrand@localhost")
  msg="chore: rebrand CRD group(s) so the CRD does not clash with upstream in a shared cluster"$'\n'
  for i in "${!OLDS[@]}"; do msg+=$'\n'"  ${OLDS[$i]} -> ${NEWS[$i]}"; done
  msg+=$'\n\n'"Applied by rebrand-core.sh from source branch '${SRC_BRANCH:-?}'."$'\n\n'"Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
  git ${commit_id[@]+"${commit_id[@]}"} commit -q -m "$msg"
  info "committed: $(git log --oneline -1)"
fi

cat <<EOF

Next steps:
  * Regenerate derived artifacts to prove no drift (optional):
      make generate    (or: make manifests)   # git diff should be empty
  * Test in a local cluster:
      hack/rebrand-test.sh
  * CI gate:
      hack/rebrand-core.sh --verify
EOF
