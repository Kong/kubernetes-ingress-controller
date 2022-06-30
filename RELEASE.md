# Release Process

To release, [create a new issue](https://github.com/Kong/kubernetes-ingress-controller/issues/new/choose) from the "Release" [template](https://github.com/Kong/kubernetes-ingress-controller/blob/main/.github/workflows/release.yaml).

Fill out the issue title and release type, create the issue, and proceed through the release steps, marking them done as you go.

# Release Troubleshooting

## Manual Docker image build

If the "Build and push development images" Github action is not appropriate for your release, or is not operating properly, you can build and push Docker images manually:

- Check out your release tag.
- Run `make container`. Note that you can set the `TAG` environment variable if you need to override the current tag in Makefile.
- Add additional tags for your container (e.g. `docker tag kong/kubernetes-ingress-controller:1.2.0-alpine kong/kubernetes-ingress-controller:1.2.0; docker tag kong/kubernetes-ingress-controller:1.2.0-alpine kong/kubernetes-ingress-controller:1.2`)
- Create a temporary token for the `kongbot` user (see 1Password) and log in using it.
- Push each of your tags (e.g. `docker push kong/kubernetes-ingress-controller:1.2.0-alpine`)

## GKE test failures

If GKE test clusters are not successfully starting, you can review their Pod logs from the Kubernetes Engine section of https://console.cloud.google.com/

You can run GKE tests locally by [creating a service account and token](https://cloud.google.com/docs/authentication/getting-started) and running, for example:

```
KUBERNETES_MAJOR_VERSION=1 KUBERNETES_MINOR_VERSION=21 GOOGLE_APPLICATION_CREDENTIALS=`cat /tmp/credentials.json` GOOGLE_PROJECT='<project name>' GOOGLE_LOCATION=us-central1 hack/e2e/dlv-tests.sh
```

You may wish to run a modified version of the script to start it with dlv and/or run a single test only. Spawning clusters is also fairly slow, so you can remove the `trap cleanup EXIT SIGINT SIGQUIT` and change `CLUSTER_NAME="e2e-$(uuidgen)"` to a static value to reuse the same cluster for multiple runs. Remember to run the cleanup function after to discard the cluster.
