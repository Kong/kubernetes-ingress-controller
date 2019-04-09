package dbless

import (
	"encoding/json"

	"github.com/hbagdi/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser"
)

var credNameToDAOName = map[string]string{
	"key-auth":   "keyauth_credentials",
	"hmac-auth":  "hmacauth_credentials",
	"basic-auth": "basicauth_credentials",
	"jwt":        "jwt_secrets",
	"acl":        "acls",
}

type Service struct {
	kong.Service
	Plugins []kong.Plugin `json:"plugins"`
	Routes  []Route       `json:"routes"`
}

type Route struct {
	kong.Route
	Plugins []kong.Plugin `json:"plugins"`
}

type Upstream struct {
	kong.Upstream
	Targets []kong.Target `json:"targets"`
}

type Certificate struct {
	// Duplicated to avoid the problem of Certificate struct having an
	// SNI as well outer layer.

	ID        *string    `json:"id,omitempty" yaml:"id,omitempty"`
	Cert      *string    `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key       *string    `json:"key,omitempty" yaml:"key,omitempty"`
	CreatedAt *int64     `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	SNIs      []kong.SNI `json:"snis"`
}

type Consumer struct {
	kong.Consumer
	Plugins     []kong.Plugin `json:"plugins"`
	Credentials map[string][]map[string]interface{}
}

func (c Consumer) MarshalJSON() ([]byte, error) {
	res := map[string]interface{}{}

	consumerJSON, err := json.Marshal(&c.Consumer)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(consumerJSON, &res)
	if err != nil {
		return nil, err
	}

	res["plugins"] = c.Plugins
	for credType, creds := range c.Credentials {
		if DAOName, ok := credNameToDAOName[credType]; ok {
			credType = DAOName
		}
		res[credType] = creds
	}
	return json.Marshal(&res)
}

type KongDeclarativeConfig struct {
	FormatVersion string        `json:"_format_version"`
	Services      []Service     `json:"services"`
	Routes        []Route       `json:"routes"`
	Upstreams     []Upstream    `json:"upstreams"`
	Certificates  []Certificate `json:"certificates"`
	Plugins       []kong.Plugin `json:"plugins"`
	Consumers     []Consumer    `json:"consumers"`
}

func KongNativeState(k8sState *parser.KongState) *KongDeclarativeConfig {
	var result KongDeclarativeConfig
	result.FormatVersion = "1.1"
	for _, s := range k8sState.Services {
		service := Service{Service: s.Service}

		for _, p := range s.Plugins {
			service.Plugins = append(service.Plugins, *p.DeepCopy())
		}
		result.Services = append(result.Services, service)

		for _, r := range s.Routes {
			route := Route{Route: r.Route}

			for _, p := range r.Plugins {
				route.Plugins = append(route.Plugins, *p.DeepCopy())
			}
			result.Routes = append(result.Routes, route)
		}
	}

	for _, plugin := range k8sState.GlobalPlugins {
		result.Plugins = append(result.Plugins, *plugin.DeepCopy())
	}

	for _, u := range k8sState.Upstreams {
		upstream := Upstream{Upstream: u.Upstream}
		for _, t := range u.Targets {
			upstream.Targets = append(upstream.Targets, *t.DeepCopy())
		}
		result.Upstreams = append(result.Upstreams, upstream)
	}

	for _, c := range k8sState.Certificates {
		cert := Certificate{
			Key:  c.Key,
			Cert: c.Cert,
		}
		for _, sni := range c.SNIs {
			cert.SNIs = append(cert.SNIs, kong.SNI{Name: kong.String(*sni)})
		}
		result.Certificates = append(result.Certificates, cert)
	}

	for _, c := range k8sState.Consumers {
		c := Consumer{Consumer: c.Consumer,
			Plugins:     c.Plugins,
			Credentials: c.Credentials,
		}
		result.Consumers = append(result.Consumers, c)
	}
	return &result
}
