//+build integration_tests

package integration

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
	"strings"
	"testing"

	types "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-testing-framework/pkg/kind"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	knativeCrds = "https://github.com/knative/serving/releases/download/v0.13.0/serving-crds.yaml"
	knativeCore = "https://github.com/knative/serving/releases/download/v0.13.0/serving-core.yaml"
)

func isKnativeReady(ctx context.Context, cluster kind.Cluster, t *testing.T) bool {
	return assert.Eventually(t, func() bool {
		podList, err := cluster.Client().CoreV1().Pods("knative-serving").List(ctx, metav1.ListOptions{})
		if err != nil {
			t.Logf("failed retrieving knative pods. %v", err)
			return false
		}

		if len(podList.Items) != 4 {
			fmt.Println("expected 4 pods, found ", len(podList.Items))
			return false
		}

		for _, pod := range podList.Items {
			if pod.Status.Phase != v1.PodRunning {
				return false
			}
		}

		fmt.Println("All knative pods are up and ready.")
		fmt.Println("Covering a window that webhook has been configured itself but not ready to receive traffic yet.")
		time.Sleep(3 * time.Second)
		return true

	}, 60*time.Second, 1*time.Second, true)
}

func TestKnativeIngress(t *testing.T) {
	_ = proxyReady()
	ctx := context.Background()

	t.Log("Deploying all resources that are required to run knative")
	err := deployManifest(knativeCrds, ctx)
	assert.NoError(t, err)
	err = deployManifest(knativeCore, ctx)
	assert.NoError(t, err)
	knativeReady := isKnativeReady(ctx, cluster, t)
	assert.Equal(t, knativeReady, true)

	t.Log("Note down the ip address or public CNAME of kong-proxy service.")
	proxy := retrieveProxyInfo(ctx, t)
	assert.NotEmpty(t, proxy)
	if err != nil {
		t.Fatalf("kong-proxy service ip/public name is not ready.")
	}

	t.Log("Configure Knative NetworkLayer as Kong")
	err = configKnativeNetwork(ctx, cluster)
	assert.NoError(t, err)
	err = configKnativeDomain(ctx, proxy, cluster)
	assert.NoError(t, err)

	t.Log("Install knative service")
	err = installKnativeSrv(ctx)
	assert.NoError(t, err)

	t.Log("Test knative service using kong.")
	srvaccessable := accessKnativeSrv(ctx, proxy, t)
	if srvaccessable == false {
		t.Fatalf("failed to access knative service.")
	}
}

func deleteManifest(yml string, ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "delete", "-f", yml)
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stdout, stdout.String())
		return err
	}
	fmt.Println("successfully delete manifest " + yml)
	return nil
}

func deployManifest(yml string, ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", yml)
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stdout, stdout.String())
		return err
	}
	fmt.Println("successfully deploy manifest " + yml)
	return nil
}

func checkIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		fmt.Printf("IP Address: %s - Invalid\n", ip)
		return false
	} else {
		fmt.Printf("IP Address: %s - Valid\n", ip)
		return true
	}
}

func retrieveProxyInfo(ctx context.Context, t assert.TestingT) string {
	var proxy string
	assert.Eventually(t, func() bool {
		cmd := exec.CommandContext(ctx, "kubectl", "get", "service", "ingress-controller-kong-proxy", "--namespace", "kong-system")
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stdout, stdout.String())
			fmt.Fprintln(os.Stderr, stderr.String())
			return false
		}

		if len(stdout.String()) > 0 {
			info := strings.Split(stdout.String(), "\n")
			proxy = strings.Fields(info[1])[3]
			fmt.Println("kong-proxy " + proxy)
			if checkIPAddress(proxy) == true {
				return true
			}
		}
		return false
	}, 60*time.Second, 1*time.Second, true)

	return proxy
}

func configKnativeNetwork(ctx context.Context, cluster kind.Cluster) error {
	payloadBytes := []byte("{\"data\": {\"ingress.class\": \"kong\"}}")
	_, err := cluster.Client().CoreV1().ConfigMaps("knative-serving").Patch(ctx, "config-network", types.MergePatchType, payloadBytes, metav1.PatchOptions{})
	if err != nil {
		fmt.Println("failed updating config map ", err)
		return err
	}

	fmt.Println("successfully configured knative network.")
	return nil
}

func installKnativeSrv(ctx context.Context) error {
	cmd := exec.CommandContext(ctx,
		"kubectl", "apply", "-f", "helloworldgo.yaml")
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stdout, stdout.String())
		fmt.Fprintln(os.Stderr, stderr.String())
		return err
	}
	fmt.Println("successfully installed knative service.")
	return nil
}

func configKnativeDomain(ctx context.Context, proxy string, cluster kind.Cluster) error {
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
		fmt.Println("failed updating config map ", err)
		return err
	}
	fmt.Println("successfully update knative config domain.")
	return nil
}

func accessKnativeSrv(ctx context.Context, proxy string, t *testing.T) bool {
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
		fmt.Println("failed generating httpquerst err ", err)
	}
	req.Header.Set("Host", "helloworld-go.default."+proxy)
	req.Host = "helloworld-go.default." + proxy

	return assert.Eventually(t, func() bool {
		resp, err := client.Do(req)
		fmt.Println("resp {", resp.StatusCode, "}")
		if err != nil {
			return false
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println("service is successfully accessed through kong.")
			return true
		}
		return false

	}, 120*time.Second, 1*time.Second)
}
