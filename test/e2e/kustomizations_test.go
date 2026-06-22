//go:build e2e_tests || istio_tests || performance_tests

package e2e

import (
	"bytes"
	"fmt"
	"io"
	"time"

	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/kubectl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/yaml"
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

	kongRouterFlavorPatch = `- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: KONG_ROUTER_FLAVOR
    value: "%s"`

	kongRouterFlavorPatchDelete = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: proxy-kong
  namespace: kong
spec:
  template:
    spec:
      containers:
      - name: proxy
        env:
        - name: KONG_ROUTER_FLAVOR
          $patch: delete`
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

func patchKongRouterFlavorFn(flavor string) func(io.Reader) (io.Reader, error) {
	kustomization := types.Kustomization{
		Patches: []types.Patch{
			{
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      "proxy-kong",
						Namespace: "kong",
					},
				},
				Patch: kongRouterFlavorPatchDelete,
			},
			{
				Patch: fmt.Sprintf(kongRouterFlavorPatch, flavor),
				Target: &types.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Name:      "proxy-kong",
						Namespace: "kong",
					},
				},
			},
		},
	}

	return func(baseManifestReader io.Reader) (io.Reader, error) {
		return kubectl.GetKustomizedManifest(kustomization, baseManifestReader)
	}
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

const kongLicenseEnvPatch = `- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: KONG_LICENSE_DATA
    valueFrom:
      secretKeyRef:
        key: license
        name: kong-enterprise-license`

// WithLicensePatch injects the enterprise license from the environment into the deployed
// manifest. It adds the KONG_LICENSE_DATA env var (referencing the kong-enterprise-license
// secret) to the proxy container of the given proxy deployment and appends the secret itself.
// proxyDeploymentName must be the name of the deployment that runs the proxy container, which
// differs between manifest variants (proxy-kong for multi-pod, ingress-kong for single-pod);
// use getProxyDeploymentName to derive it from the manifest path.
func WithLicensePatch(proxyDeploymentName string) ManifestPatch {
	return func(r io.Reader) (io.Reader, error) {
		licenseSecret, err := ktfkong.GetLicenseSecretFromEnv()
		if err != nil {
			return nil, err
		}
		// GetLicenseSecretFromEnv returns a secret without TypeMeta/namespace; both are
		// required for it to apply cleanly into the kong namespace via `kubectl apply -f -`.
		licenseSecret.TypeMeta = metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"}
		licenseSecret.Namespace = namespace

		// Add the KONG_LICENSE_DATA env var (referencing the secret) to the proxy
		// container (index 0) of the proxy deployment.
		kustomization := types.Kustomization{
			Patches: []types.Patch{
				{
					Patch: kongLicenseEnvPatch,
					Target: &types.Selector{
						ResId: resid.ResId{
							Gvk:       resid.Gvk{Group: "apps", Version: "v1", Kind: "Deployment"},
							Name:      proxyDeploymentName,
							Namespace: namespace,
						},
					},
				},
			},
		}
		patched, err := kubectl.GetKustomizedManifest(kustomization, r)
		if err != nil {
			return nil, err
		}

		// Append the license secret as an additional document so it is created
		// alongside the deployment.
		secretYAML, err := yaml.Marshal(licenseSecret)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, patched); err != nil {
			return nil, err
		}
		buf.WriteString("\n---\n")
		buf.Write(secretYAML)
		return &buf, nil
	}
}

// konnectLicensingEnvPatch is the patch to add the CONTROLLER_KONNECT_LICENSING_ENABLED env var to the controller container to enable Konnect licensing behavior.
const konnectLicensingEnvPatch = `- op: add
  path: /spec/template/spec/containers/0/env/-
  value:
    name: CONTROLLER_KONNECT_LICENSING_ENABLED
    value: "true"`

// WithKonnectLicensingPatch adds the CONTROLLER_KONNECT_LICENSING_ENABLED env var to the controller container to enable Konnect licensing behavior.
// This is required for tests that involve Konnect licensing with Kong versions that require a license at startup,
// since the presence of this env var changes the licensing activation flow to be compatible with such versions.
func WithKonnectLicensingPatch() ManifestPatch {
	return func(r io.Reader) (io.Reader, error) {
		ingressDeploymentName := "ingress-kong"
		// Add the CONTROLLER_KONNECT_LICENSING_ENABLED env var to the controller container to enable Konnect licensing behavior.
		kustomization := types.Kustomization{
			Patches: []types.Patch{
				{
					Patch: konnectLicensingEnvPatch,
					Target: &types.Selector{
						ResId: resid.ResId{
							Gvk:       resid.Gvk{Group: "apps", Version: "v1", Kind: "Deployment"},
							Name:      ingressDeploymentName,
							Namespace: namespace,
						},
					},
				},
			},
		}
		patched, err := kubectl.GetKustomizedManifest(kustomization, r)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, patched); err != nil {
			return nil, err
		}
		return &buf, nil
	}
}
