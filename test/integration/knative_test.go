//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/knative"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	knativenetworkingversioned "knative.dev/networking/pkg/client/clientset/versioned"
	"knative.dev/pkg/apis"
)

const (
	// knativeWaitTime indicates how long to wait for knative components to be up and running
	// on the cluster. The current value is based on deployment times seen in a GKE environment.
	knativeWaitTime = time.Minute * 2
)

// knativeMinKubernetesVersion indicates the minimum Kubernetes version
// required in order to successfully run Knative tests.
var knativeMinKubernetesVersion = semver.MustParse("1.21.0")

func TestKnativeIngress(t *testing.T) {
	if clusterVersion.LT(knativeMinKubernetesVersion) {
		t.Skip("knative tests can't be run on cluster versions prior to 1.21")
	}

	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("generating a knative clientset")
	dynamicClient, err := dynamic.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	knativeGVR := schema.GroupVersionResource{
		Group:    "serving.knative.dev",
		Version:  "v1",
		Resource: "services",
	}
	knativeClient := dynamicClient.Resource(knativeGVR).Namespace(ns.Name)

	t.Logf("configure knative network for ingress class %s", ingressClass)
	payloadBytes := []byte(fmt.Sprintf("{\"data\": {\"ingress-class\": \"%s\"}}", ingressClass))
	_, err = env.Cluster().Client().CoreV1().ConfigMaps(knative.DefaultNamespace).Patch(ctx, "config-network", types.MergePatchType, payloadBytes, metav1.PatchOptions{})
	require.NoError(t, err)
	require.NoError(t, configKnativeDomain(ctx, proxyURL.Hostname(), knative.DefaultNamespace, env.Cluster()))

	t.Log("deploying a native service to test routing")
	service := &unstructured.Unstructured{}
	service.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "serving.knative.dev/v1",
		"kind":       "Service",
		"metadata": map[string]interface{}{
			"name":      "helloworld-go",
			"namespace": ns.Name,
		},
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"image": "gcr.io/knative-samples/helloworld-go",
							"env": []map[string]interface{}{
								{
									"name":  "TARGET",
									"value": "Go Sample v1",
								},
							},
						},
					},
				},
			},
		},
	})
	require.Eventually(t, func() bool {
		_, err = knativeClient.Create(ctx, service, metav1.CreateOptions{})
		return err == nil
	}, knativeWaitTime, waitTick, true)

	defer func() {
		t.Log("cleaning up knative services used for testing")
		assert.NoError(t, knativeClient.Delete(ctx, "helloworld-go", metav1.DeleteOptions{}))
	}()

	t.Log("Test knative service using kong.")
	require.True(t, accessKnativeSrv(ctx, proxyURL.Hostname(), ns.Name, t), true)
}

// -----------------------------------------------------------------------------
// Knative Deployment Functions
// -----------------------------------------------------------------------------

func configKnativeDomain(ctx context.Context, proxy, nsn string, cluster clusters.Cluster) error {
	// generate the new config-domain configmap
	configMap := v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-domain",
			Namespace: nsn,
			Labels: map[string]string{
				"serving.knative.dev/release": "v0.18.0",
			},
		},
		Data: map[string]string{
			proxy: "",
		},
	}

	// update the config-domain configmap with the new values
	_, err := cluster.Client().CoreV1().ConfigMaps(nsn).Update(ctx, &configMap, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func accessKnativeSrv(ctx context.Context, proxy, nsn string, t *testing.T) bool {
	knativec, err := knativenetworkingversioned.NewForConfig(env.Cluster().Config())
	if err != nil {
		return false
	}
	ingCli := knativec.NetworkingV1alpha1().Ingresses(nsn)
	assert.Eventually(t, func() bool {
		curIng, err := ingCli.Get(ctx, "helloworld-go", metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting knative ingress: %v", err)
			return false
		}
		if curIng == nil {
			t.Log("getting knative ingress: got nil, want non-nil")
			return false
		}

		conds := curIng.Status.Status.GetConditions()
		for _, cond := range conds {
			if cond.Type == apis.ConditionReady && cond.Status == v1.ConditionTrue {
				t.Logf("knative ingress %s/%s status is ready.", curIng.Namespace, curIng.Name)
				return true
			}
		}
		t.Logf("knative ingress %s/%s not ready yet", curIng.Namespace, curIng.Name)
		return false
	}, statusWait, waitTick, true)

	url := "http://" + proxy
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   time.Second * 60,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Logf("failed generating httpquerst err %v", err)
	}
	req.Header.Set("Host", fmt.Sprintf("helloworld-go.%s.%s", nsn, proxy))
	req.Host = fmt.Sprintf("helloworld-go.%s.%s", nsn, proxy)

	return assert.Eventually(t, func() bool {
		resp, err := client.Do(req)
		if err != nil {
			t.Logf("error requesting %q: %v", req.URL, err)
			return false
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Logf("error reading response body (url: %s): %v", req.URL, err)
			return false
		}

		if resp.StatusCode == http.StatusOK {
			bodyString := string(bodyBytes)
			t.Logf(bodyString)
			t.Log("service is successfully accessed through kong.")
			return true
		}

		t.Logf("expected HTTP 200 but got %d, with body: %q", resp.StatusCode, bodyBytes)
		return false
	}, knativeWaitTime, waitTick)
}
