package utils

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/metadata"
)

// GetKubeconfig returns a Kubernetes REST config object based on the configuration.
func GetKubeconfig(c managercfg.Config) (*rest.Config, error) {
	var (
		config *rest.Config
		err    error
	)
	switch c.KubeRestConfig {
	case nil:
		// If no kubeconfig path or REST config is provided, use the in-cluster config.
		config, err = clientcmd.BuildConfigFromFlags(c.APIServerHost, c.KubeconfigPath)
		if err != nil {
			return nil, err
		}
	default:
		// If a REST config is provided, use it directly.
		config = c.KubeRestConfig
	}

	// Set the user agent so it's possible to identify the controller in the API server logs.
	config.UserAgent = metadata.UserAgent()

	// Configure K8s client rate-limiting.
	config.QPS = float32(c.APIServerQPS)
	config.Burst = c.APIServerBurst

	if c.APIServerCertData != nil {
		config.CertData = c.APIServerCertData
	}
	if c.APIServerCAData != nil {
		config.CAData = c.APIServerCAData
	}
	if c.APIServerKeyData != nil {
		config.KeyData = c.APIServerKeyData
	}
	if c.Impersonate != "" {
		config.Impersonate.UserName = c.Impersonate
	}

	return config, err
}

// GetKubeClient returns a Kubernetes client based on the configuration.
func GetKubeClient(c managercfg.Config) (client.Client, error) {
	conf, err := GetKubeconfig(c)
	if err != nil {
		return nil, err
	}
	return client.New(conf, client.Options{})
}
