# Release Process

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) `v20.10.x`
- [GNU Make](https://www.gnu.org/software/make/) `v4.x`
- [Kustomize](https://github.com/kubernetes-sigs/kustomize) `v1.3.x`

## Github Workflow Test Matrix Checkup

**For all releases**

We maintain some integration tests with 3rd party components which we need to manually verify and update before cutting any release.

- [ ] check the testing workflow (`.github/workflows/test.yaml`) and ensure that all matrix versions are up to date for various component releases. If there have been any new releases (major, minor or patch) of those components since the latest version seen in that configuration make sure the new versions get added before proceeding with the release.
  - [ ] Kubernetes
  - [ ] Istio

An issue exists to automate the above actions: https://github.com/Kong/kubernetes-ingress-controller/issues/1886

## Release Testing

**For all releases**

We currently provide a suite of tests with an extensive matrix of component versions (see `.github/workflows/release-test.yaml`).

Prior to any release open the [workflow page](https://github.com/Kong/kubernetes-ingress-controller/actions/workflows/release-test.yaml) and trigger a test run against `main`.

**Wait for tests to pass before continuing with any release**, if any problems are found hold on the release until patches are provided and then run the tests again.

## Release Branch

**For all releases**

For this step we're going to start with the `main` branch to create our release branch (e.g. `release/X.Y.Z`) which will later be submitted as a pull request back to `main`.

- [ ] ensure that you have up to date copy of `main`: `git fetch --all`
- [ ] create the release branch for the version (e.g. `release/1.3.1`): `git branch -m release/x.y.z`
- [ ] Make any final adjustments to CHANGELOG.md. Double-check that dates are correct, that link anchors point to the correct header, and that you've included a link to the Github compare link at the end.
- [ ] update the `TAG` variable in the `Makefile` to the new version release and commit the change
- [ ] ensure base manifest versions use the new version and update manifest files: `make manifests`
- [ ] ensure that the Kubernetes versions provisioned in the cloud (GKE, etc.) as part of the release CI pipeline are up to date
  - [ ] remove any versions that are no longer supported by the cloud provider, or the release pipeline will fail
- [ ] push the branch up to the remote: `git push --set-upstream origin release/x.y.z`

## Release Pull Request

**For all releases**

- [ ] Open a PR from your branch to `main`
- [ ] Once the PR is merged, tag your release: `git fetch --all && git tag origin/main 1.3.1 && git push origin --tags`
- [ ] Wait for CI to build images and push them to Docker Hub

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
- [ ] Update `references/version-compatibility.md` to include the new versions (make sure you capture any new Kubernetes/Istio versions that have been tested)
- [ ] Copy `app/_data/docs_nav_kic_OLDVERSION.yml` to `app/_data/docs_nav_kic_NEWVERSION.yml`. Add entries for any new articles.
- [ ] Add a section to `app/_data/kong_versions.yml` for your version.
- [ ] Open a PR from your branch.

# Release Troubleshooting

## Manual Docker image build

If the "Build and push development images" Github action is not appropriate for your release, or is not operating properly, you can build and push Docker images manually:

- [ ] Check out your release tag.
- [ ] Run `make container`. Note that you can set the `TAG` environment variable if you need to override the current tag in Makefile.
- [ ] Add additional tags for your container (e.g. `docker tag kong/kubernetes-ingress-controller:1.2.0-alpine kong/kubernetes-ingress-controller:1.2.0; docker tag kong/kubernetes-ingress-controller:1.2.0-alpine kong/kubernetes-ingress-controller:1.2`)
- [ ] Create a temporary token for the `kongbot` user (see 1Password) and log in using it.
- [ ] Push each of your tags (e.g. `docker push kong/kubernetes-ingress-controller:1.2.0-alpine`)
