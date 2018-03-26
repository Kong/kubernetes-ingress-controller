package admin

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/blang/semver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type KongInterface interface {
	RESTClient() rest.Interface

	RouteGetter
	ServiceGetter
	UpstreamGetter
	TargetGetter
	SNIGetter
	CertificateGetter
	PluginGetter
	CredentialGetter
}

type RestClient struct {
	restClient rest.Interface
}

func (c *RestClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}

func (c *RestClient) Routes() RouteInterface {
	return &routeAPI{
		client: &apiClient{c.RESTClient(),
			&metav1.APIResource{
				Name:       "routes",
				Namespaced: false,
			},
		},
	}
}
func (c *RestClient) Services() ServiceInterface {
	return &serviceAPI{
		client: &apiClient{c.RESTClient(),
			&metav1.APIResource{
				Name:       "services",
				Namespaced: false,
			},
		},
	}
}

func (c *RestClient) Upstreams() UpstreamInterface {
	return &upstreamAPI{
		client: &apiClient{c.RESTClient(),
			&metav1.APIResource{
				Name:       "upstreams",
				Namespaced: false,
			},
		},
	}
}

func (c *RestClient) Targets() TargetInterface {
	return &targetAPI{
		client: &apiClient{c.RESTClient(),
			&metav1.APIResource{
				Name:       "targets",
				Namespaced: false,
			},
		},
	}
}

func (c *RestClient) SNIs() SNIInterface {
	return &sniAPI{
		client: &apiClient{c.RESTClient(),
			&metav1.APIResource{
				Name:       "snis",
				Namespaced: false,
			},
		},
	}
}

func (c *RestClient) Certificates() CertificateInterface {
	return &certificateAPI{
		client: &apiClient{c.RESTClient(),
			&metav1.APIResource{
				Name:       "certificates",
				Namespaced: false,
			},
		},
	}
}

func (c *RestClient) Plugins() PluginInterface {
	return &pluginAPI{
		client: &apiClient{c.RESTClient(),
			&metav1.APIResource{
				Name:       "plugins",
				Namespaced: false,
			},
		},
	}
}

func (c *RestClient) Consumers() ConsumerInterface {
	return &consumerAPI{
		client: &apiClient{c.RESTClient(),
			&metav1.APIResource{
				Name:       "consumers",
				Namespaced: false,
			},
		},
	}
}

func (c *RestClient) Credentials() CredentialInterface {
	return &credentialAPI{c.RESTClient()}
}

func (c *RestClient) GetVersion() (semver.Version, error) {
	var info map[string]interface{}
	data, err := c.RESTClient().Get().RequestURI("/").DoRaw()
	if err != nil {
		return semver.Version{}, err
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return semver.Version{}, err
	}

	if version, ok := info["version"]; ok {
		v := version.(string)
		// fix bad version formats like 0.13.0preview1
		re := regexp.MustCompile(`(.*\d)(preview.*|rc.*)`)
		if re.MatchString(v) {
			v = re.ReplaceAllString(v, "$1-$2")
		}

		return semver.Make(v)
	}

	return semver.Version{}, fmt.Errorf("Unknown Kong version")
}

func NewRESTClient(c *rest.Config) (*RestClient, error) {
	c.ContentConfig = dynamic.ContentConfig()
	cl, err := rest.UnversionedRESTClientFor(c)
	if err != nil {
		return nil, err
	}

	return &RestClient{
		restClient: cl,
	}, nil
}
