//go:build integration_tests

package integration

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

const examplesDIR = "../../examples"

func TestTLSRouteExample(t *testing.T) {
	t.Log("locking Gateway TLS ports")
	tlsMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	})

	var (
		tlsrouteExampleManifests = fmt.Sprintf("%s/gateway-tlsroute.yaml", examplesDIR)
		ctx                      = context.Background()
	)
	_, cleaner := helpers.Setup(ctx, t, env)

	t.Logf("applying yaml manifest %s", tlsrouteExampleManifests)
	b, err := os.ReadFile(tlsrouteExampleManifests)
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyManifestByYAML(ctx, env.Cluster(), string(b)))
	cleaner.AddManifest(string(b))
	// Copy pasted cert from gateway-tlsroute.yaml to use it in certPool for checking validity.
	const caForCert = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVGekNDQXYrZ0F3SUJBZ0lVZTFvWnRWQVBOM1V2bXRkSHo5OFpYcDd2a3Znd0RRWUpLb1pJaHZjTkFRRUwKQlFBd2dZZ3hDekFKQmdOVkJBWVRBbFZUTVJNd0VRWURWUVFJREFwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSApEQTFUWVc0Z1JuSmhibU5wYzJOdk1SSXdFQVlEVlFRS0RBbExiMjVuSUVsdVl5NHhHREFXQmdOVkJBc01EMVJsCllXMGdTM1ZpWlhKdVpYUmxjekVlTUJ3R0ExVUVBd3dWZEd4emNtOTFkR1V1YTI5dVp5NWxlR0Z0Y0d4bE1DQVgKRFRJME1EY3dOVEUwTlRjek5sb1lEekl4TWpRd05qRXhNVFExTnpNMldqQ0JpREVMTUFrR0ExVUVCaE1DVlZNeApFekFSQmdOVkJBZ01Da05oYkdsbWIzSnVhV0V4RmpBVUJnTlZCQWNNRFZOaGJpQkdjbUZ1WTJselkyOHhFakFRCkJnTlZCQW9NQ1V0dmJtY2dTVzVqTGpFWU1CWUdBMVVFQ3d3UFZHVmhiU0JMZFdKbGNtNWxkR1Z6TVI0d0hBWUQKVlFRRERCVjBiSE55YjNWMFpTNXJiMjVuTG1WNFlXMXdiR1V3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQgpEd0F3Z2dFS0FvSUJBUURleklnV1E4L1crYUx4VERPanZvOGY5OVpGYVBoOWQrTkFYQ1NmSTFLcS9TdVdrMnJyClcxKyt2QzA4MERWTnc4dmx3Z0VRUUlhV3c2bVh6V2hQNmppdHIvemxSVGg4TWFwSFMvTXhXbjN0WnFKZ3ZVdVoKMkFnMTBXUE14UHV4UUlaU2FucU95M0RNeDJDcGlMQ1c0SVBERlRhQm5XT1hOeFg4bEMvQit6QlZYYzBIYVdUUwpqUFViUUZONGVGcEFtcHlxak1Dak53Y1VSd3BBVSs0cXpDeVZ2ZU5VU0RLWHpoN04rUDlPRkFiVjNqL0IyOXpqCk9sVFZKNTUvZ2VUeGJqZVZCa0ZDZXAvQkh4UEY4MnhtWUJOQnJ2WVU5dFkyc0JCZmh6OGFUNFJaMmx5NXJxVnYKRnZ4TDF1R3ZmU29CeUdoVTVFWDg0NmZVYm5uc2xJSDdBKzdOQWdNQkFBR2pkVEJ6TUIwR0ExVWREZ1FXQkJTdwp2c2VNR08wN1JXMXpxVWNsOFZEeXY2M25HakFmQmdOVkhTTUVHREFXZ0JTd3ZzZU1HTzA3UlcxenFVY2w4VkR5CnY2M25HakFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQ0FHQTFVZEVRUVpNQmVDRlhSc2MzSnZkWFJsTG10dmJtY3UKWlhoaGJYQnNaVEFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBa3JJQWp6Ulhadm51bjAvMGc2NHNnK0laaDZXLwo0UWFsYklwWlh0Z2JrKzE3d1FVN0hrTFNaS2tmL2IwZHlTR1RpLzBDMUdYK1lxRTlKb2NtWjdjUWI2Y0RheDl3CjBUeFZ1NVRxYWVTQVlLZktGZlExcjB5VzhjYVc3TFhoYVZ3Lzh2YTVQMWRzdnkwNUs0K3dydCtBK1NudFRxL2EKYzV2T0ZQcmk3ZlBMWEZ5SVE5eXhVZXdSYnphdUY2SEE5eDY4bWt3WUVvSTUxMnM1SjBtUVByUGhPN1VhRnFwVApQL3NqdXdHNk1qR0t1MzU2VXJKTGlFV1NkZmtiTkU0bGFLa3Z2U0paMjh3MXVGemsrOUl5QWpiRzQ3RXl3dTVICjJ0ZzZzV3VTa29xT2hwc3JTMjBteFlkVHU0dVpWZXdlR1FjbmZXUGNjNlRQME8vQzZCZGpRVFB3Nmc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
	certPool := x509.NewCertPool()
	decodedCert, err := base64.StdEncoding.DecodeString(caForCert)
	assert.NoError(t, err)
	assert.True(t, certPool.AppendCertsFromPEM(decodedCert))

	t.Log("verifying that TLSRoute becomes routable")
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		err := tlsEchoResponds(
			proxyTLSURL, "tlsroute-example-manifest", "tlsroute.kong.example", certPool, true,
		)
		assert.NoError(c, err)
	}, ingressWait, waitTick)
}
