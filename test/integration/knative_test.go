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

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/knative"
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
)

func TestKnativeIngress(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("generating a knative clientset")
	knativec, err := knativeversioned.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Logf("configure knative network for ingress class %s", ingressClass)
	payloadBytes := []byte(fmt.Sprintf("{\"data\": {\"ingress.class\": \"%s\"}}", ingressClass))
	_, err = env.Cluster().Client().CoreV1().ConfigMaps(knative.DefaultNamespace).Patch(ctx, "config-network", types.MergePatchType, payloadBytes, metav1.PatchOptions{})
	require.NoError(t, err)
	require.NoError(t, configKnativeDomain(ctx, proxyURL.Hostname(), knative.DefaultNamespace, env.Cluster()))

	t.Log("deploying a native service to test routing")
	var service *knservingv1.Service
	require.Eventually(t, func() bool {
		service, err = installKnativeSrv(ctx, ns.Name, knativec)
		return err == nil
	}, knativeWaitTime, waitTick, true)

	defer func() {
		t.Log("cleaning up knative services used for testing")
		assert.NoError(t, knativec.ServingV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Log("Test knative service using kong.")
	require.True(t, accessKnativeSrv(ctx, proxyURL.Hostname(), ns.Name, t), true)
}

// -----------------------------------------------------------------------------
// Knative Deployment Functions
// -----------------------------------------------------------------------------

func installKnativeSrv(ctx context.Context, nsn string, knativec *knativeversioned.Clientset) (*knservingv1.Service, error) {
	// generate the knative service resource
	tobeDeployedService := &knservingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "helloworld-go",
			Namespace: nsn,
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
										}}}}}}}}}}

	// deploy the new service to the cluster
	service, err := knativec.ServingV1().Services(nsn).Create(ctx, tobeDeployedService, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create knative service. %w", err)
	}

	return service, nil
}

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
	req.Header.Set("Host", fmt.Sprintf("helloworld-go.%s.%s", nsn, proxy))
	req.Host = fmt.Sprintf("helloworld-go.%s.%s", nsn, proxy)

	return assert.Eventually(t, func() bool {
		resp, err := client.Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return false
			}
			bodyString := string(bodyBytes)
			t.Logf(bodyString)
			t.Log("service is successfully accessed through kong.")
			return true
		}
		return false
	}, knativeWaitTime, waitTick)
}
