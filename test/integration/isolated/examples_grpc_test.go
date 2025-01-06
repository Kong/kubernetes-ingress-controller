//go:build integration_tests

package isolated

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestGRPCRouteExample(t *testing.T) {
	testGRPC := func(ctx context.Context, t *testing.T, manifestName string, gatewayPort int, hostname string, certPool *x509.CertPool) {
		t.Helper()
		cluster := GetClusterFromCtx(ctx)
		proxyURL := GetHTTPURLFromCtx(ctx)
		manifestPath := examplesManifestPath(manifestName)
		t.Logf("applying yaml manifest %s", manifestPath)
		b, err := os.ReadFile(manifestPath)
		assert.NoError(t, err)
		manifest := string(b)
		assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, manifest))

		t.Log("verifying that GRPCRoute becomes routable")
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			err := grpcEchoResponds(
				ctx, fmt.Sprintf("%s:%d", proxyURL.Hostname(), gatewayPort), hostname, "kong", certPool,
			)
			assert.NoError(c, err)
		}, consts.IngressWait, consts.WaitTick)

		t.Logf("deleting yaml manifest %s", manifestPath)
		assert.NoError(t, clusters.DeleteManifestByYAML(ctx, cluster, manifest))
	}

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindGRPCRoute).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
			withKongProxyEnvVars(map[string]string{
				"PROXY_LISTEN": `0.0.0.0:8000 http2\, 0.0.0.0:8443 http2 ssl`,
			}),
		)).
		Assess("deploying to cluster works and deployed GRPC via HTTP responds", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			testGRPC(ctx, t, "gateway-grpcroute-via-http.yaml", ktfkong.DefaultProxyHTTPPort, "example-grpc-via-http.com", nil)
			return ctx
		}).
		Assess("deploying to cluster works and deployed GRPC via HTTPS responds", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			// Copy pasted cert from gateway-grpcroute-via-https.yaml to use it in certPool for checking validity.
			const caForCert = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUQ5ekNDQXQrZ0F3SUJBZ0lVTWQrODVFTE9BT2hzN3FmclRhUi9yclh1UEFjd0RRWUpLb1pJaHZjTkFRRUwKQlFBd2ZqRUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdNQ2tOaGJHbG1iM0p1YVdFeEZqQVVCZ05WQkFjTQpEVk5oYmlCR2NtRnVZMmx6WTI4eEVqQVFCZ05WQkFvTUNVdHZibWNnU1c1akxqRVlNQllHQTFVRUN3d1BWR1ZoCmJTQkxkV0psY201bGRHVnpNUlF3RWdZRFZRUUREQXRsZUdGdGNHeGxMbU52YlRBZ0Z3MHlOREEzTURVeE16VTMKTWpKYUdBOHlNVEkwTURZeE1URXpOVGN5TWxvd2ZqRUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdNQ2tOaApiR2xtYjNKdWFXRXhGakFVQmdOVkJBY01EVk5oYmlCR2NtRnVZMmx6WTI4eEVqQVFCZ05WQkFvTUNVdHZibWNnClNXNWpMakVZTUJZR0ExVUVDd3dQVkdWaGJTQkxkV0psY201bGRHVnpNUlF3RWdZRFZRUUREQXRsZUdGdGNHeGwKTG1OdmJUQ0NBU0l3RFFZSktvWklodmNOQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQU1FLzVxeFRoMFphVmkwdgpUY3dZdlZ0RjlXVjhLb3paTFN2SHJkQm96dUdYWHFpZmp0NzRzSkFwdVIrT3Nsam8rTkp1Wm0xRVl4NnRuV0dkCnBOandSa2VHdWdCVFFqa1Q5NmJ5V3dwZ0c0K21QL01RU3pDcjE2T3BhNTFOckVHV3lUYzJ6K1B3TlF6SzJ1SWUKbVlxaHJRa2xFUG1WemRRTXZoeWV4dkpoY0p0RWZ3MUgrUFlNNVN1cmwyUDJFNXhwZGpRTXJxZytjTVNqeSs2TApKb283VXFSZU9hNHBtM0Z5Sm9NQTdXUU9GUWt3U2dsV21QUFJ4RDJaaG9FQnl3YmlwODdHYkRIbkFVOUhhaG5SCiswL0FrdldlaXM2YW9GajV4bWJxWmxQek1YUVdXZ1dvM2ZzQVY4U0lKV3o3UVVGM21WUWtGU24welRBc2RrREoKZ21GUDNmTUNBd0VBQWFOck1Ha3dIUVlEVlIwT0JCWUVGS25Ia21YaEZibzZ4L21JZXdmb3dWTkZnVEFLTUI4RwpBMVVkSXdRWU1CYUFGS25Ia21YaEZibzZ4L21JZXdmb3dWTkZnVEFLTUE4R0ExVWRFd0VCL3dRRk1BTUJBZjh3CkZnWURWUjBSQkE4d0RZSUxaWGhoYlhCc1pTNWpiMjB3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUdDRlRzZHkKVmVCMGxUOE9xcXNCZzUvYUFyeGttSDJ3UCtTYndTbXlGRG0xS1pycmwyQWwxUHozMDA3aWVHTE9KRkp3LzNkZgpmek56MjFSNytmZThjNG51eWNJR3Yvc2ZwbEtVbWNRRm9kWXdqUkRON3UvOVoycVYzMjNFSldIaXExVnF1VXpqCnpKWXBDWWRXWlRraVlwMmNWdUxzZlBLbFI2VVB3Z3JoSU94MEVSL0wveTFIVzU2NnNHZ3lDU1k2V1crQ0UzS28KZEVVZUdjYzM4NlR0WGNQWVRTa2tGejdWdkM3QVNrQjdtWmlSV080RlFRWXUrelRoYm8vVXlhYXJSa3lzb0xGbgpnb2lVaVBhaFhWVTl6L3Mva0ppZWNTZ0t3UU5tQkhsR084RXhjODYyekZyRVFHUFpXTFhZWTRlQXVjV2VLQUpZCjgzbWZYQ3I0c3V1aU1ldz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
			certPool := x509.NewCertPool()
			decodedCert, err := base64.StdEncoding.DecodeString(caForCert)
			assert.NoError(t, err)
			assert.True(t, certPool.AppendCertsFromPEM(decodedCert))

			testGRPC(ctx, t, "gateway-grpcroute-via-https.yaml", ktfkong.DefaultProxyTLSServicePort, "example.com", certPool)
			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
