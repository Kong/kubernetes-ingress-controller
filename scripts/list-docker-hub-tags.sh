#!/usr/bin/env bash

# ------------------------------------------------------------------------------
# Lists tags of a Docker Hub repository that were pushed within a time window
# [since, until) and whose names match a given regex.
#
# Matching tags are printed to stdout, one per line, de-duplicated and sorted
# by version. Diagnostic messages are printed to stderr.
# ------------------------------------------------------------------------------

set -euo pipefail

usage() {
  cat >&2 <<EOF
Usage: $(basename "$0") --repository <repo> --regex <regex> [--name-filter <filter>] <window>

  --repository <repo>       Docker Hub repository, e.g. kong/kong-gateway-dev
  --regex <regex>           Regex (jq syntax) the full tag name must match.
  --name-filter <filter>    Optional substring passed to Docker Hub to narrow down results
                            server side (Docker Hub only supports substring filtering).
                            Recommended for repositories with many tags, as paginating
                            through all of them may get rate limited by Docker Hub.

Exactly one of the following time windows must be given:
  --since <time> [--until <time>]
                            Tags pushed in [since, until). Times are RFC3339 UTC timestamps,
                            e.g. 2026-10-14T06:00:00Z. Without --until there is no upper bound.
  --lookback-hours <hours>  Tags pushed within the given number of hours before now.
  --daily-window <HH:MM>    Tags pushed in the 24 hours before the latest HH:MM (UTC) that
                            is not after now. Consecutive runs scheduled daily at HH:MM cover
                            adjacent windows, even if a run starts late.
EOF
  exit 1
}

REPOSITORY=""
REGEX=""
NAME_FILTER=""
SINCE=""
UNTIL=""
LOOKBACK_HOURS=""
DAILY_WINDOW=""

while [[ $# -gt 0 ]]; do
  [[ $# -ge 2 ]] || usage
  case "$1" in
    --repository) REPOSITORY="$2" ;;
    --regex) REGEX="$2" ;;
    --name-filter) NAME_FILTER="$2" ;;
    --since) SINCE="$2" ;;
    --until) UNTIL="$2" ;;
    --lookback-hours) LOOKBACK_HOURS="$2" ;;
    --daily-window) DAILY_WINDOW="$2" ;;
    *) usage ;;
  esac
  shift 2
done

[[ -n "${REPOSITORY}" && -n "${REGEX}" ]] || usage

windows=0
[[ -n "${SINCE}" ]] && windows=$((windows + 1))
[[ -n "${LOOKBACK_HOURS}" ]] && windows=$((windows + 1))
[[ -n "${DAILY_WINDOW}" ]] && windows=$((windows + 1))
if [[ "${windows}" -ne 1 || ( -n "${UNTIL}" && -z "${SINCE}" ) ]]; then
  usage
fi

if [[ -n "${LOOKBACK_HOURS}" ]]; then
  [[ "${LOOKBACK_HOURS}" =~ ^[0-9]+$ ]] || usage
  SINCE=$(jq -rn --argjson h "${LOOKBACK_HOURS}" 'now - $h * 3600 | floor | todate')
elif [[ -n "${DAILY_WINDOW}" ]]; then
  [[ "${DAILY_WINDOW}" =~ ^([01][0-9]|2[0-3]):[0-5][0-9]$ ]] || usage
  read -r SINCE UNTIL < <(jq -rn --arg anchor "${DAILY_WINDOW}" '
    (now | floor) as $now
    | ($now | strftime("%Y-%m-%dT") + $anchor + ":00Z" | fromdate) as $today
    | (if $today > $now then $today - 86400 else $today end) as $until
    | "\($until - 86400 | todate) \($until | todate)"
  ')
fi

url="https://hub.docker.com/v2/repositories/${REPOSITORY}/tags?page_size=100&ordering=last_updated"
if [[ -n "${NAME_FILTER}" ]]; then
  url+="&name=$(jq -rn --arg v "${NAME_FILTER}" '$v | @uri')"
fi

echo "Listing tags of ${REPOSITORY} pushed in [${SINCE}, ${UNTIL:-now}) matching ${REGEX}" >&2

# Tags are sorted newest first, so stop paginating once we reach tags older than SINCE.
tags=""
while [[ -n "${url}" ]]; do
  page=$(curl -fsSL --retry 3 "${url}")
  tags+=$(jq -r --arg since "${SINCE}" --arg until "${UNTIL}" --arg regex "${REGEX}" '
    .results[]
    | (.tag_last_pushed // .last_updated) as $pushed
    | select($pushed >= $since and ($until == "" or $pushed < $until))
    | .name
    | select(test($regex))
  ' <<< "${page}")$'\n'
  oldest=$(jq -r '.results[-1] | (.tag_last_pushed // .last_updated) // ""' <<< "${page}")
  if [[ -z "${oldest}" || "${oldest}" < "${SINCE}" ]]; then
    break
  fi
  url=$(jq -r '.next // ""' <<< "${page}")
done

echo "${tags}" | sed '/^$/d' | sort -uV
