package utils

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// GetKubeconfig returns a Kubernetes REST config object based on the configuration.
func GetKubeconfig(c config.Config) (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags(c.APIServerHost, c.KubeconfigPath)
	if err != nil {
		return nil, err
	}

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

	config.UserAgent = metadata.UserAgent()

	return config, err
}

// GetKubeClient returns a Kubernetes client based on the configuration.
func GetKubeClient(c config.Config) (client.Client, error) {
	conf, err := GetKubeconfig(c)
	if err != nil {
		return nil, err
	}
	return client.New(conf, client.Options{})
}
