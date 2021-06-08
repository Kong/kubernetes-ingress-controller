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
	kongyaml    = "https://bit.ly/k4k8s"
)

func isKnativeReady(ctx context.Context, cluster kind.Cluster) bool {
	timeout := 1
	for timeout < 60 {
		podList, err := cluster.Client().CoreV1().Pods("knative-serving").List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Println("failed retrieving knative pods.", err)
			time.Sleep(1 * time.Second)
			continue
		}

		ready := true
		if len(podList.Items) != 4 {
			fmt.Println("expected 4 pods, found ", len(podList.Items))
			time.Sleep(1 * time.Second)
			continue
		}

		for _, pod := range podList.Items {
			if pod.Status.Phase != v1.PodRunning {
				ready = false
			}
		}
		if ready == true {
			fmt.Println("All knative pods are up and ready.")
			fmt.Println("Covering a window that webhook has been configured itself but not ready to receive traffic yet.")
			time.Sleep(3 * time.Second)
			return true
		}
		if ready == false {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("Failed to bring up knative resources within 60 seconds.")
	return false
}

func TestKnativeIngress(t *testing.T) {
	ctx := context.Background()

	t.Log("Deploying all resources that are required to run knative")
	err := deployManifest(knativeCrds, ctx)
	assert.NoError(t, err)
	err = deployManifest(knativeCore, ctx)
	assert.NoError(t, err)
	knativeReady := isKnativeReady(ctx, cluster)
	assert.Equal(t, knativeReady, true)

	t.Log("Deploying kong ingress.")
	err = deployManifest(kongyaml, ctx)
	assert.NoError(t, err)

	t.Log("Note down the ip address or public CNAME of kong-proxy service.")
	proxy, err := retrieveProxyInfo(ctx)
	assert.NoError(t, err)

	t.Log("Configure Knative NetworkLayer as Kong")
	err = configKnativeNetwork(ctx, cluster)
	assert.NoError(t, err)
	err = configKnativeDomain(ctx, proxy, cluster)
	assert.NoError(t, err)

	t.Log("Install knative service")
	err = installKnativeSrv(ctx)
	assert.NoError(t, err)

	t.Log("Check knative service readiness.")
	err = ensureKnativeSrv(ctx)
	if err != nil {
		t.Fatalf("Knative Service is not ready.")
	}

	t.Log("Test knative service using kong.")
	srvaccessable := accessKnativeSrv(ctx, proxy)
	if srvaccessable == false {
		t.Fatalf("failed to access knative service.")
	}

	t.Log("clean up test deployments.")
	deleteManifest(knativeCrds, ctx)
	deleteManifest(knativeCore, ctx)
	deleteManifest(kongyaml, ctx)
	deleteManifest("helloworldgo.yaml", ctx)
	time.Sleep(5 * time.Second)
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

func retrieveProxyInfo(ctx context.Context) (string, error) {
	cnt := 1
	for cnt < 60 {
		cmd := exec.CommandContext(ctx, "kubectl", "get", "service", "kong-proxy", "--namespace", "kong")
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stdout, stdout.String())
			fmt.Fprintln(os.Stderr, stderr.String())
			return "", err
		}

		if len(stdout.String()) > 0 {
			info := strings.Split(stdout.String(), "\n")
			proxy := strings.Fields(info[1])[3]
			fmt.Println("kong-proxy " + proxy)
			if checkIPAddress(proxy) == true {
				return proxy, nil
			}
			time.Sleep(1 * time.Second)
			continue
		}
	}
	return "", nil
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

func ensureKnativeSrv(ctx context.Context) error {
	cnt := 1
	for cnt < 120 {
		cmd := exec.CommandContext(ctx, "kubectl", "get", "ksvc")
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stdout, stdout.String())
			return err
		}
		if len(stdout.String()) > 0 {
			if strings.Contains(stdout.String(), "True") {
				fmt.Println("knative service has been up.")
				return nil
			}
		}
		time.Sleep(1 * time.Second)
		cnt += 1
	}
	return fmt.Errorf("knative service failed to be up.")
}

func accessKnativeSrv(ctx context.Context, proxy string) bool {
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

	resp, err := client.Do(req)
	fmt.Println("resp {", resp, "}")
	if err != nil {
		fmt.Println("WARNING: error ", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return true
	}
	fmt.Println("knative service query ", resp)
	return false
}
