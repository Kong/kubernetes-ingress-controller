//go:build e2e_tests || istio_tests

package e2e

import (
	"fmt"
	"io"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/kubectl"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

const (
	initRetryPatch = `- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: CONTROLLER_KONG_ADMIN_INIT_RETRY_DELAY
    value: "%s"
- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: CONTROLLER_KONG_ADMIN_INIT_RETRIES
    value: "%d"`

	livenessProbePatch = `- op: replace
  path: /spec/template/spec/containers/0/livenessProbe/initialDelaySeconds
  value: %[1]d
- op: replace
  path: /spec/template/spec/containers/0/livenessProbe/timeoutSeconds
  value: %[2]d
- op: replace
  path: /spec/template/spec/containers/0/livenessProbe/failureThreshold
  value: %[3]d`
)

// patchControllerImage replaces the kong/kubernetes-ingress-controller image with the provided image and tag,
// and returns the modified manifest.
func patchControllerImage(baseManifestReader io.Reader, image, tag string) (io.Reader, error) {
	kustomization := types.Kustomization{
		Images: []types.Image{
			{
				Name:    "kong/kubernetes-ingress-controller",
				NewName: image,
				NewTag:  tag,
			},
		},
	}
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

// patchKongImage replaces the kong and kong/kong-gateway images in a manifest with the provided image and tag,
// and returns the modified manifest.
func patchKongImage(baseManifestReader io.Reader, image, tag string) (io.Reader, error) {
	kustomization := types.Kustomization{
		Images: []types.Image{
			{
				Name:    "kong/kong-gateway",
				NewName: image,
				NewTag:  tag,
			},
			{
				Name:    "kong",
				NewName: image,
				NewTag:  tag,
			},
		},
	}
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

// patchControllerStartTimeout adds or updates the controller container CONTROLLER_KONG_ADMIN_INIT_RETRIES and
// CONTROLLER_KONG_ADMIN_INIT_RETRY_DELAY environment variables with the provided values, and returns the modified
// manifest.
func patchControllerStartTimeout(baseManifestReader io.Reader, tries int, delay time.Duration) (io.Reader, error) {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Patch: fmt.Sprintf(initRetryPatch, delay.String(), tries),
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      controllerDeploymentName,
						Namespace: "kong",
					},
				},
			},
		},
	}
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

// patchLivenessProbes patches the given deployment's liveness probe, replacing the initial delay, period, and failure
// threshold.
func patchLivenessProbes(baseManifestReader io.Reader, deployment k8stypes.NamespacedName, failure int, initial, period time.Duration) (io.Reader, error) {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Patch: fmt.Sprintf(livenessProbePatch, int(initial.Seconds()), int(period.Seconds()), failure),
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      deployment.Name,
						Namespace: deployment.Namespace,
					},
				},
			},
		},
	}
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

// formatKustomizePatchFirstContainerReadinessProbePath returns the patch
// to update the readiness probe path of the first container in the deployment to the probePath.
func formatKustomizePatchFirstContainerReadinessProbePath(probePath string) string {
	const readinessProbePathPatchFormat = "- op: replace\n" +
		"  path: /spec/template/spec/containers/0/readinessProbe/httpGet/path\n" +
		`  value: "%s"`

	return fmt.Sprintf(readinessProbePathPatchFormat, probePath)
}

// patchReadinessProbePath patches the given deployment's path of readiness probe.
// Kong gateway supports endpoint `/status/ready` since 3.4, and prior versions uses `/status`,
// so we need to change the path of readiness probe of Kong gateway deployment when its version is < 3.4.
func patchReadinessProbePath(baseManifestReader io.Reader, deployment k8stypes.NamespacedName, probePath string) (io.Reader, error) {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Patch: formatKustomizePatchFirstContainerReadinessProbePath(probePath),
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      deployment.Name,
						Namespace: deployment.Namespace,
					},
				},
			},
		},
	}
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}

const kustomizePatchKongAdminAPIListenFormat = `apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: %s
  name: %s
spec:
  template:
    spec:
      containers:
      - name: %s
        env:
        - name: KONG_ADMIN_LISTEN
          value: %s
`

func patchKongAdminAPIListen(baseManifestReader io.Reader, deployment k8stypes.NamespacedName, adminAPIListenConfig string) (io.Reader, error) {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Patch: fmt.Sprintf(kustomizePatchKongAdminAPIListenFormat,
					deployment.Namespace, deployment.Name, proxyContainerName, adminAPIListenConfig,
				),
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      deployment.Name,
						Namespace: deployment.Namespace,
					},
				},
			},
		},
	}
	return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
}
