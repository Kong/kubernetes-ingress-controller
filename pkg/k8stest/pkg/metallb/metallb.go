package metallb

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"os/exec"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ktfdocker "github.com/kong/kubernetes-ingress-controller/pkg/k8stest/pkg/docker"
	ktfkind "github.com/kong/kubernetes-ingress-controller/pkg/k8stest/pkg/kind"
	ktfnet "github.com/kong/kubernetes-ingress-controller/pkg/k8stest/pkg/networking"
)

// -----------------------------------------------------------------------------
// Public Functions - Metallb Management
// -----------------------------------------------------------------------------

// DeployMetallbForKindCluster deploys Metallb to the given Kind cluster using the Docker network provided for LoadBalancer IPs.
func DeployMetallbForKindCluster(kindClusterName, dockerNetwork string) error {
	// grab a kubernetes client for the cluster
	kc, err := ktfkind.ClientForKindCluster(kindClusterName)
	if err != nil {
		return err
	}

	// ensure the namespace for metallb is created
	ns := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "metallb-system"}}
	if _, err := kc.CoreV1().Namespaces().Create(context.Background(), &ns, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	// get an IP range for the docker container network to use for MetalLB
	network, err := ktfdocker.GetDockerContainerIPNetwork(ktfkind.GetKindDockerContainerID(kindClusterName), dockerNetwork)
	if err != nil {
		return err
	}
	ipStart, ipEnd := GetIPRangeForMetallb(*network)

	// deploy the metallb configuration
	cfgMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config",
			Namespace: "metallb-system",
		},
		Data: map[string]string{
			"config": GetMetallbYAMLCfg(ipStart, ipEnd),
		},
	}
	if _, err := kc.CoreV1().ConfigMaps(ns.Name).Create(context.Background(), cfgMap, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	// generate and deploy a metallb memberlist secret
	secretKey := make([]byte, 128)
	if _, err := rand.Read(secretKey); err != nil {
		return err
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "memberlist",
			Namespace: ns.Name,
		},
		StringData: map[string]string{
			"secretkey": base64.StdEncoding.EncodeToString(secretKey),
		},
	}
	if _, err := kc.CoreV1().Secrets(ns.Name).Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	// create the metallb deployment and related resources
	return metallbDeployHack(kindClusterName)
}

// -----------------------------------------------------------------------------
// Public Functions - Helper
// -----------------------------------------------------------------------------

// GetIPRangeForMetallb provides a range of IP addresses to use for MetalLB given an IPv4 Network
// FIXME: just choosing specific default IPs for now, need to check range validity and dynamically assign IPs.
func GetIPRangeForMetallb(network net.IPNet) (startIP, endIP net.IP) {
	startIP = ktfnet.ConvertUint32ToIPv4(ktfnet.ConvertIPv4ToUint32(network.IP) | ktfnet.ConvertIPv4ToUint32(defaultStartIP))
	endIP = ktfnet.ConvertUint32ToIPv4(ktfnet.ConvertIPv4ToUint32(network.IP) | ktfnet.ConvertIPv4ToUint32(defaultEndIP))
	return
}

func GetMetallbYAMLCfg(ipStart, ipEnd net.IP) string {
	return fmt.Sprintf(`
address-pools:
- name: default
  protocol: layer2
  addresses:
  - %s
`, ktfnet.GetIPRangeStr(ipStart, ipEnd))
}

// -----------------------------------------------------------------------------
// Private Consts & Vars
// -----------------------------------------------------------------------------

var (
	defaultStartIP = net.ParseIP("0.0.0.240")
	defaultEndIP   = net.ParseIP("0.0.0.250")
	metalManifest  = "https://raw.githubusercontent.com/metallb/metallb/v0.9.5/manifests/metallb.yaml"
)

// -----------------------------------------------------------------------------
// Private Functions
// -----------------------------------------------------------------------------

// FIXME: needs to be replaced with non-kubectl, just used this originally for speed.
func metallbDeployHack(clusterName string) error {
	deployArgs := []string{
		"--context", fmt.Sprintf("kind-%s", clusterName),
		"apply", "-f", metalManifest,
	}
	cmd := exec.Command("kubectl", deployArgs...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
