//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	knativenetworkingversioned "knative.dev/networking/pkg/client/clientset/versioned"
	"knative.dev/pkg/apis"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	knativeversioned "knative.dev/serving/pkg/client/clientset/versioned"
)

const (
	// knativeWaitTime indicates how long to wait for knative components to be up and running
	// on the cluster. The current value is based on deployment times seen in a GKE environment.
	knativeWaitTime = time.Minute * 2

	// knativeNamespace is the testing namespace where Knative components will be deployed
	knativeNamespace = "knative-serving"

	knativeCrds = "https://github.com/knative/serving/releases/download/v0.13.0/serving-crds.yaml"
	knativeCore = "https://github.com/knative/serving/releases/download/v0.13.0/serving-core.yaml"
)

func TestKnativeIngress(t *testing.T) {
	cluster := env.Cluster()
	proxy := proxyURL.Hostname()
	assert.NotEmpty(t, proxy)
	t.Logf("proxy url %s", proxy)

	ctx := context.Background()

	t.Log("Deploying all resources that are required to run knative")
	require.NoError(t, deployManifest(ctx, knativeCrds, t))
	require.NoError(t, deployManifest(ctx, knativeCore, t))
	require.True(t, isKnativeReady(ctx, cluster, t), true)

	t.Log("Configure Knative NetworkLayer as Kong")
	require.NoError(t, configKnativeNetwork(ctx, cluster, t))
	require.NoError(t, configKnativeDomain(ctx, proxy, cluster, t))

	t.Log("Install knative service")
	require.Eventually(t, func() bool {
		err := installKnativeSrv(ctx, t)
		if err != nil {
			t.Log("checking knativing webhook readiness.")
			return false
		}
		return true
	}, 30*time.Second, 2*time.Second, true)

	t.Log("Test knative service using kong.")
	require.True(t, accessKnativeSrv(ctx, proxy, t), true)
}

// -----------------------------------------------------------------------------
// Knative Deployment Functions
// -----------------------------------------------------------------------------

// TODO: in future iterations Knative components will become deployable as an "addon" for test clusters
//       in our testing framework, and then we can remove this deployment logic and just have the tests.
//       See: https://github.com/Kong/kubernetes-testing-framework/issues/75

func deployManifest(ctx context.Context, yml string, t *testing.T) error {
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

func configKnativeNetwork(ctx context.Context, cluster clusters.Cluster, t *testing.T) error {
	payloadBytes := []byte(fmt.Sprintf("{\"data\": {\"ingress.class\": \"%s\"}}", ingressClass))
	_, err := cluster.Client().CoreV1().ConfigMaps(knativeNamespace).Patch(ctx, "config-network", types.MergePatchType, payloadBytes, metav1.PatchOptions{})
	if err != nil {
		t.Logf("failed updating config map %v", err)
		return err
	}

	t.Log("successfully configured knative network.")
	return nil
}

func installKnativeSrv(ctx context.Context, t *testing.T) error {
	tobeDeployedService := &knservingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "helloworld-go",
			Namespace: "default",
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
	knativeCli, err := knativeversioned.NewForConfig(env.Cluster().Config())
	if err != nil {
		return fmt.Errorf("failed to create knative service. %v", err)
	}

	_, err = knativeCli.ServingV1().Services("default").Create(ctx, tobeDeployedService, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create knative service. %v", err)
	}

	t.Log("successfully installed knative service.")
	return nil
}

func configKnativeDomain(ctx context.Context, proxy string, cluster clusters.Cluster, t *testing.T) error {
	configMapData := make(map[string]string)
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
			Namespace: knativeNamespace,
			Labels:    labels,
		},
		Data: configMapData,
	}
	_, err := cluster.Client().CoreV1().ConfigMaps(knativeNamespace).Update(ctx, &configMap, metav1.UpdateOptions{})
	if err != nil {
		t.Logf("failed updating config map %v", err)
		return err
	}
	t.Log("successfully update knative config domain.")
	return nil
}

func accessKnativeSrv(ctx context.Context, proxy string, t *testing.T) bool {
	knativeCli, err := knativenetworkingversioned.NewForConfig(env.Cluster().Config())
	if err != nil {
		return false
	}
	ingCli := knativeCli.NetworkingV1alpha1().Ingresses("default")
	assert.Eventually(t, func() bool {
		curIng, err := ingCli.Get(ctx, "helloworld-go", metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		conds := curIng.Status.Status.GetConditions()
		for _, cond := range conds {
			if cond.Type == apis.ConditionReady && cond.Status == v1.ConditionTrue {
				t.Log("knative ingress status is ready.")
				return true
			}
		}
		return false
	}, 120*time.Second, 1*time.Second, true)

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
	req.Header.Set("Host", "helloworld-go.default."+proxy)
	req.Host = "helloworld-go.default." + proxy

	return assert.Eventually(t, func() bool {
		resp, err := client.Do(req)
		t.Logf("resp <%v> ", resp.StatusCode)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return false
			}
			bodyString := string(bodyBytes)
			t.Logf(bodyString)
			t.Log("service is successfully accessed through kong.")
			return true
		}
		return false

	}, 120*time.Second, 1*time.Second)
}

func isKnativeReady(ctx context.Context, cluster clusters.Cluster, t *testing.T) bool {
	// the deployment manifests for knative include some CPU and Memory limits which
	// are good for production, but mostly just problematic when running simple tests
	// where these components are going to be brought up and torn down quickly.
	// we tear out these requirements ad as long as the pods start we will likely have
	// all the CPU and memory we need to complete the tests (whereafter we will tear
	// all of the knative components down anyhow).
	deploymentList, err := cluster.Client().AppsV1().Deployments(knativeNamespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err)
	for _, deployment := range deploymentList.Items {
		require.Eventually(t, func() bool {
			for i := 0; i < len(deployment.Spec.Template.Spec.Containers); i++ {
				deployment.Spec.Template.Spec.Containers[i].Resources = corev1.ResourceRequirements{}
			}
			_, err = cluster.Client().AppsV1().Deployments(knativeNamespace).Update(ctx, &deployment, metav1.UpdateOptions{})
			return err == nil
		}, knativeWaitTime, waitTick)
	}

	return assert.Eventually(t, func() bool {
		podList, err := cluster.Client().CoreV1().Pods(knativeNamespace).List(ctx, metav1.ListOptions{})
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

		t.Log("All knative pods are up and ready.")
		return true

	}, knativeWaitTime, waitTick, true)
}
