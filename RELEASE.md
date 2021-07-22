# Release Process

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) `v20.10.x`
- [GNU Make](https://www.gnu.org/software/make/) `v4.x`
- [Kustomize](https://github.com/kubernetes-sigs/kustomize) `v1.3.x`

## Release Branch

**For all releases**

For this step we're going to start with the `next` branch and merge in `main` to create our release branch (e.g. `release/X.Y.Z`) which will later be submitted as a pull request back to `main`.

- [ ] ensure that you have up to date copy of `main` and `next`: `git fetch --all`
- [ ] merge main into next: `git checkout next && git merge main`
- [ ] create the release branch for the version (e.g. `release/1.3.1`): `git branch -m release/x.y.z`
- [ ] Make any final adjustments to CHANGELOG.md. Double-check that dates are correct, that link anchors point to the correct header, and that you've included a link to the Github compare link at the end.
- [ ] update manifest files: `./hack/build-single-manifests.sh`
- [ ] update the `TAG` variable in the `Makefile` to the new version release and commit the change
- [ ] push the branch up to the remote: `git push --set-upstream origin release/x.y.z`

## Release Pull Request

**For all releases**

- [ ] Open a PR from your branch to `next` (major/minor) or `main` (patch)
- [ ] Once the PR is merged, tag your release: `git fetch --all && git tag origin/next 1.3.1 && git push origin --tags`
- [ ] Wait for CI to build images and push them to Docker Hub
- [ ] Open a PR from next to main (ignore if patch release)

## Github Release

**For all releases**

- [ ] verify that CI is passing for `main` first: if there are CI errors on main they must be investigated and fixed
- [ ] draft a new [release](https://github.com/Kong/kubernetes-ingress-controller/releases), using a title and body similar to previous releases. Use your existing tag.
- [ ] for new `major` version releases create a new branch (e.g. `1.3.x`) from the release tag and push it
- [ ] for `minor` and `patch` version releases rebase the release tag onto the release branch: `git checkout 1.3.x && git rebase 1.3.1 && git push`

## Documentation

**For major/minor releases only**

- [ ] Create a new branch in the [documentation site repo](https://github.com/Kong/docs.konghq.com).
- [ ] Copy `app/kubernetes-ingress-controller/OLD_VERSION` to `app/kubernetes-ingress-controller/NEW_VERSION`.
- [ ] Update articles in the new version as needed.
- [ ] Update `references/version-compatibility.md` to include the new version.
- [ ] Copy `app/_data/docs_nav_kic_OLDVERSION.yml` to `app/_data/docs_nav_kic_NEWVERSION.yml`. Add entries for any new articles.
- [ ] Add a section to `app/_data/kong_versions.yml` for your version.
- [ ] Open a PR from your branch.

# Release Troubleshooting

## Manual Docker image build

If the "Build and push development images" Github action is not appropriate for your release, or is not operating properly, you can build and push Docker images manually:

- [ ] Check out your release tag.
- [ ] Run `make container` (legacy) or `make railgun-container` (Railgun/2.x). Note that you can set the `TAG` environment variable if you need to override the current tag in Makefile.
- [ ] Add additional tags for your container (e.g. `docker tag kong/kubernetes-ingress-controller:1.2.0-alpine kong/kubernetes-ingress-controller:1.2.0; docker tag kong/kubernetes-ingress-controller:1.2.0-alpine kong/kubernetes-ingress-controller:1.2`)
- [ ] Create a temporary token for the `kongbot` user (see 1Password) and log in using it.
- [ ] Push each of your tags (e.g. `docker push kong/kubernetes-ingress-controller:1.2.0-alpine`)
