//go:build integration_tests && knative
// +build integration_tests,knative

package integration

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/knative"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	knativenetworkingversioned "knative.dev/networking/pkg/client/clientset/versioned"
	"knative.dev/pkg/apis"
	kservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	kservingclientsetv "knative.dev/serving/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

const (
	// knativeWaitTime indicates how long to wait for knative components to be up and running
	// on the cluster. The current value is based on deployment times seen in a GKE environment.
	knativeWaitTime = time.Minute * 2
)

func TestKnativeIngress(t *testing.T) {
	skipTestForExpressionRouter(t)
	ctx := context.Background()

	if clusterVersion.LT(knativeMinKubernetesVersion) {
		t.Skip("knative tests can't be run on cluster versions prior to", knativeMinKubernetesVersion)
	}

	t.Parallel()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("creating a knative client")
	kservingClient, err := kservingclientsetv.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Logf("configure knative network for ingress class %s", consts.IngressClass)
	payloadBytes := []byte(fmt.Sprintf("{\"data\": {\"ingress-class\": \"%s\"}}", consts.IngressClass))
	_, err = env.Cluster().Client().CoreV1().ConfigMaps(knative.DefaultNamespace).Patch(ctx, "config-network", types.MergePatchType, payloadBytes, metav1.PatchOptions{})
	require.NoError(t, err)
	require.NoError(t, configKnativeDomain(ctx, proxyURL.Hostname(), knative.DefaultNamespace, env.Cluster()))

	service := &kservingv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "httpbin",
			Namespace: ns.Name,
		},
		Spec: kservingv1.ServiceSpec{
			ConfigurationSpec: kservingv1.ConfigurationSpec{
				Template: kservingv1.RevisionTemplateSpec{
					Spec: kservingv1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Image: "kong/httpbin:0.1.0",
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	cleaner.Add(service)

	t.Logf("deploying knative service %s to test routing", service.GetName())
	require.Eventually(t, func() bool {
		_, err = kservingClient.ServingV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			t.Logf("failed creating knative service: %v", err)
			return false
		}
		return true
	}, knativeWaitTime, waitTick)

	t.Log("test knative service using kong")
	require.True(t, accessKnativeSrv(t, ctx, proxyURL.Hostname(), ns.Name, service.GetName()))
}

// -----------------------------------------------------------------------------
// Knative Deployment Functions
// -----------------------------------------------------------------------------

func configKnativeDomain(ctx context.Context, proxy, nsn string, cluster clusters.Cluster) error {
	// generate the new config-domain configmap
	configMap := corev1.ConfigMap{
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

func accessKnativeSrv(t *testing.T, ctx context.Context, proxy, nsn, serviceName string) bool {
	knativec, err := knativenetworkingversioned.NewForConfig(env.Cluster().Config())
	if err != nil {
		return false
	}
	ingCli := knativec.NetworkingV1alpha1().Ingresses(nsn)
	assert.Eventually(t, func() bool {
		curIng, err := ingCli.Get(ctx, serviceName, metav1.GetOptions{})
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
			if cond.Type == apis.ConditionReady && cond.Status == corev1.ConditionTrue {
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
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   time.Second * 60,
	}

	host := fmt.Sprintf("%s.%s.%s", serviceName, nsn, proxy)
	req := helpers.MustHTTPRequest(t, "GET", proxyURL, url, nil)
	req.Header.Set("Host", host)
	req.Host = host
	req.URL.Path = "/status/200"

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
