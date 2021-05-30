# KIC Release process

## Git workflow and image builds

- [ ] Merge main into next (ignore if patch release). Use a merge commit.
- [ ] Check out a new branch from next (major/minor release) or main (patch release), e.g. `release/1.2.0`.
- [ ] Update manifest versions under `deploy/base` with the new version and run `hack/build-single-manifests.sh` (ignore for pre-releases).
- [ ] Bump the Makefile version (see [example](https://github.com/Kong/kubernetes-ingress-controller/pull/851/commits/b874b36bfdb0d7c6a13cc35ed06f666d135b04a4)
- [ ] Make any final adjustments to CHANGELOG.md. Double-check that dates are correct, that link anchors point to the correct header, and that you've included a link to the Github compare link at the end.
- [ ] Open a PR from your branch to next (major/minor) or main (patch).
- [ ] Once the PR is merged, tag your release (e.g. `git fetch; git tag origin/next 1.2.0; git push origin --tags`.
- [ ] Wait for CI to build images and push them to Docker Hub.
- [ ] Open a PR from next to main (ignore if patch release).
- [ ] [Draft a new release](https://github.com/Kong/kubernetes-ingress-controller/releases), using a title and body similar to previous releases. Use your existing tag.
- [ ] Create a new minor version branch and push it, e.g. `git checkout 1.2.0; git checkout -b 1.2.x; git push 1.2.x` (ignore if patch or pre-release).

## Documentation

For major/minor releases only.

- [ ] Create a new branch in the [documentation site repo](https://github.com/Kong/docs.konghq.com).
- [ ] Copy `app/kubernetes-ingress-controller/OLD_VERSION` to `app/kubernetes-ingress-controller/NEW_VERSION`.
- [ ] Update articles in the new version as needed.
- [ ] Update `references/version-compatibility.md` to include the new version.
- [ ] Copy `app/_data/docs_nav_kic_OLDVERSION.yml` to `app/_data/docs_nav_kic_NEWVERSION.yml`. Add entries for any new articles.
- [ ] Add a section to `app/_data/kong_versions.yml` for your version.
- [ ] Open a PR from your branch.

## Manual Docker image build

If the "Build and push development images" Github action is not appropriate for your release, or is not operating properly, you can build and push Docker images manually.

- [ ] Check out your release tag.
- [ ] Run `make container` (legacy) or `make railgun-container` (Railgun/2.x). Note that you can set the `TAG` environment variable if you need to override the current tag in Makefile.
- [ ] Add additional tags for your container (e.g. `docker tag kong/kubernetes-ingress-controller:1.2.0-alpine kong/kubernetes-ingress-controller:1.2.0; docker tag kong/kubernetes-ingress-controller:1.2.0-alpine kong/kubernetes-ingress-controller:1.2`)
- [ ] Create a temporary token for the `kongbot` user (see 1Password) and log in using it.
- [ ] Push each of your tags (e.g. `docker push kong/kubernetes-ingress-controller:1.2.0-alpine`)
