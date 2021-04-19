---
status: implemented
---

# Kubernetes Ingress Controller (KIC) Minimal Definition of Done

## Summary

Create a minimal **Definition of Done** for issues. While there are many areas a definition of done shall cover, start small - with satisfactory and sufficient conditions to close an issue.

## Motivation

- For a closed issue, the following questions don't have obvious answers:
    - What was the resolution (fixed/rejected-infeasible/rejected-unnecessary/duplicate/not understood/stale)?
    - Has the fix/enhancement made it to `main`?
    - Has the fix/enhancement made it to a public release?
    - As a maintainer, should I close an issue:
        - as soon as all the PRs get merged to `next`?
        - as soon as all the PRs get merged to `main`?
        - as soon as a {RC, alpha, beta, GA} release includes the fix/enhancement?

### Goals

Applies to the following GitHub repos:
- kubernetes-ingress-controller
- charts
- deck
- kong-operator
- go-kong



- Closed GitHub issues have resolutions that are easy to understand.
- Maintainers know under what circumstances they shall close issues.
- For a GitHub issue, it's easy to see what released version includes the fix, if any.

### Non-Goals

- Decide when a PR is good to be merged to `next` or `main`. This is a separate discussion to be had.
- Decide on a quality bar for new code (this belongs to a Definition of Done, but is to be discussed separately to keep scope under control).

## Proposal

### User Stories

#### Story 1

As a {contributor|maintainer|user} of the repositories (or their artifacts), when I find a GitHub issue describing a problem, I understand what the status of the fix/enhancement is. Especially if the issue is closed, I can understand whether the issue has been fixed or rejected, and which released version includes the fix.

#### Story 2

As a maintainer of the repositories (or contributor), assuming that I have solved the problem or contributed the fix/enhancement, I know under what circumstances I shall close the issue, and what information to provide.

## Design Details

Add the following guideline to `CONTRIBUTING.MD`:

For a GitHub issue describing a problem/feature request:

- **Duplicates**. if there are other issues in the repository describing the same problem/FR:

    1. find the issue that has the most context (possibly not the first reported)

    1. close all other issues with a comment _Duplicate of #XYZ_

- **Resolution by code**. if resolution involves creating PRs:

    1. ensure that all PRs reference the issue they are solving. Keep in mind that the _fixes_/_resolves_ directive only works for PRs merged to the default branch of the repository.

    1. close the issue as soon as all the PRs have been merged to **`main` or `next`**. If it's obvious from PRs that the issue has been resolved, a closing comment on the issue is purely optional.

- **Other resolutions/rejections**. if resolution happens for any other reason (_resolved without code_, _user's question answered_, _won't fix_, _infeasible_, _not useful_, _alternative approach chosen_, _problem will go away in $FUTURE-VERSION_)

    1. close the issue with a comment describing the resolution/reason.

For a closed issue, one can verify which released versions contain the fix/enhancement by navigating into the merge commit of each attached PR, where GitHub lists tags/branches that contain the merge commit.
Thus:
- if the list includes a release tag: the fix/enhancement is included in that release tag.
- if the list includes `next` but no release tags: the fix/enhancement will come in the nearest minor release.
- if the list includes `main` but no release tags: the fix/enhancement will come in the nearest patch release.

## Alternatives

Closing an issue as "resolved-fixed" only when all the code is released in some major/minor/patch release.

Pros:

- fewer hops for a user to verify whether the fix/enhancement has been released

Cons:

- ties issue resolution to releases
- makes planning/acceptance more complicated because there is uncertainty whether a "closed" issue is actually "fully complete"
- makes completion of milestones unclear (if a milestone is to be released at once, it would be visible as 0% completed until a release is made, and then suddenly jump to 100% completion)
- requires additional maintainer/community coordination to retroactively close issues post-release
