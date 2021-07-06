//+build performance_tests

package performance

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"bytes"
	"context"
	"fmt"
	"testing"

	types "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-testing-framework/pkg/kind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knativenetworkingversioned "knative.dev/networking/pkg/client/clientset/versioned"
	"knative.dev/pkg/apis"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	knativeversioned "knative.dev/serving/pkg/client/clientset/versioned"
)

const (
	knativeCrds = "https://github.com/knative/serving/releases/download/v0.13.0/serving-crds.yaml"
	knativeCore = "https://github.com/knative/serving/releases/download/v0.13.0/serving-core.yaml"
)

func TestPerfKnativePerformance(t *testing.T) {

	proxy := KongInfo.ProxyURL.Hostname()
	assert.NotEmpty(t, proxy)
	t.Logf("proxy url %s", proxy)

	ctx := context.Background()

	t.Log("Deploying all resources that are required to run knative")
	require.NoError(t, deployManifest(knativeCrds, ctx, t))
	require.NoError(t, deployManifest(knativeCore, ctx, t))
	require.True(t, isKnativeReady(ctx, cluster, t), true)

	t.Log("Configure Knative NetworkLayer as Kong")
	require.NoError(t, perfconfigKnativeNetwork(ctx, cluster, t))
	require.NoError(t, perfconfigKnativeDomain(ctx, proxy, cluster, t))

	cnt := 1
	cost := 0
	for cnt <= max_ingress {
		namespace := fmt.Sprintf("knative-%d", cnt)
		nsName := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err := cluster.Client().CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Logf("Install knative service into namespace %s", namespace)
		require.Eventually(t, func() bool {
			err := perfInstallKnativeSrv(ctx, t, namespace)
			if err != nil {
				t.Logf("checking knativing webhook readiness.")
				return false
			}
			return true
		}, 30*time.Second, 2*time.Second, true)

		t.Logf("Test knative service using kong within namespace %s.", namespace)
		require.True(t, perfaccessKnativeSrv(ctx, proxy, t, namespace, &cost), true)
		cnt++
	}
	t.Logf("knative ingress cost %d", cost/max_ingress)

	// cleanup clean ingress also to continueally repeat
}

func deployManifest(yml string, ctx context.Context, t *testing.T) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", yml)
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stdout, stdout.String())
		return err
	}
	t.Logf("successfully deploy manifest " + yml)
	return nil
}

func perfcheckIPAddress(ip string, t *testing.T) bool {
	if net.ParseIP(ip) == nil {
		t.Logf("IP Address: %s - Invalid\n", ip)
		return false
	} else {
		t.Logf("IP Address: %s - Valid\n", ip)
		return true
	}
}

func perfconfigKnativeNetwork(ctx context.Context, cluster kind.Cluster, t *testing.T) error {
	payloadBytes := []byte(fmt.Sprintf("{\"data\": {\"ingress.class\": \"%s\"}}", "kong"))
	_, err := cluster.Client().CoreV1().ConfigMaps("knative-serving").Patch(ctx, "config-network", types.MergePatchType, payloadBytes, metav1.PatchOptions{})
	if err != nil {
		t.Logf("failed updating config map %v", err)
		return err
	}

	t.Logf("successfully configured knative network.")
	return nil
}

func perfInstallKnativeSrv(ctx context.Context, t *testing.T, namespace string) error {
	tobeDeployedService := &knservingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "helloworld-go",
			Namespace: namespace,
		},
		Spec: knservingv1.ServiceSpec{
			ConfigurationSpec: knservingv1.ConfigurationSpec{
				Template: knservingv1.RevisionTemplateSpec{
					Spec: knservingv1.RevisionSpec{
						PodSpec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Image: "gcr.io/knative-samples/helloworld-go",
									Env: []corev1.EnvVar{
										{
											Name:  "TARGET",
											Value: "Go Sample v1",
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
	knativeCli, err := knativeversioned.NewForConfig(cluster.Config())
	_, err = knativeCli.ServingV1().Services(namespace).Create(ctx, tobeDeployedService, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create knative service. %v", err)
	}
	t.Logf("successfully installed knative service.")
	return nil
}

func perfconfigKnativeDomain(ctx context.Context, proxy string, cluster kind.Cluster, t *testing.T) error {
	configMapData := make(map[string]string, 0)
	configMapData[proxy] = ""
	labels := make(map[string]string)
	labels["serving.knative.dev/release"] = "v0.13.0"
	configMap := v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-domain",
			Namespace: "knative-serving",
			Labels:    labels,
		},
		Data: configMapData,
	}
	_, err := cluster.Client().CoreV1().ConfigMaps("knative-serving").Update(ctx, &configMap, metav1.UpdateOptions{})
	if err != nil {
		t.Logf("failed updating config map %v", err)
		return err
	}
	t.Logf("successfully update knative config domain.")
	return nil
}

func perfaccessKnativeSrv(ctx context.Context, proxy string, t *testing.T, namesapce string, cost *int) bool {
	knativeCli, err := knativenetworkingversioned.NewForConfig(cluster.Config())
	if err != nil {
		return false
	}
	ingCli := knativeCli.NetworkingV1alpha1().Ingresses(namesapce)
	s := time.Now().Second()
	assert.Eventually(t, func() bool {
		curIng, err := ingCli.Get(ctx, "helloworld-go", metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		conds := curIng.Status.Status.GetConditions()
		for _, cond := range conds {
			if cond.Type == apis.ConditionReady && cond.Status == v1.ConditionTrue {
				e := time.Now().Second()
				t.Logf("knative ingress status is ready. cost %d seconds ", (e - s))
				*cost += (e - s)
				return true
			}
		}
		return false
	}, 120*time.Second, 1*time.Second, true)

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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Logf("failed generating httpquerst err %v", err)
	}
	req.Header.Set("Host", "helloworld-go."+namesapce+"."+proxy)
	req.Host = "helloworld-go." + namesapce + "." + proxy

	return assert.Eventually(t, func() bool {
		resp, err := client.Do(req)
		t.Logf("resp <%v>", resp.StatusCode)
		if err != nil {
			return false
		}
		if resp.StatusCode == http.StatusOK {
			t.Logf("service is successfully accessed through kong.")
			return true
		}
		return false

	}, 120*time.Second, 1*time.Second)
}

func isKnativeReady(ctx context.Context, cluster kind.Cluster, t *testing.T) bool {
	return assert.Eventually(t, func() bool {
		podList, err := cluster.Client().CoreV1().Pods("knative-serving").List(ctx, metav1.ListOptions{})
		if err != nil {
			t.Logf("failed retrieving knative pods. %v", err)
			return false
		}

		if len(podList.Items) != 4 {
			t.Logf("expected 4 pods, found %d", len(podList.Items))
			return false
		}

		for _, pod := range podList.Items {
			if pod.Status.Phase != v1.PodRunning {
				return false
			}
		}

		t.Logf("All knative pods are up and ready.")
		return true

	}, 60*time.Second, 1*time.Second, true)
}
