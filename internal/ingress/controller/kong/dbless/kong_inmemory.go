package dbless

import (
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

// Service is a Kong Service, it's plugins and routes associated with it.
type Service struct {
	kong.Service
	Plugins []kong.Plugin `json:"plugins"`
	Routes  []Route       `json:"routes"`
}

// Route is a Kong Route and the plugins associated with it.
type Route struct {
	kong.Route
	Plugins []kong.Plugin `json:"plugins"`
}

// Upstream is a Kong Upstream and it's targets.
type Upstream struct {
	kong.Upstream
	Targets []kong.Target `json:"targets"`
}

// Certificate is a Kong Certificate and it's associated SNIs.
type Certificate struct {
	// Duplicated to avoid the problem of Certificate struct having an
	// SNI as well outer layer.

	ID        *string    `json:"id,omitempty" yaml:"id,omitempty"`
	Cert      *string    `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key       *string    `json:"key,omitempty" yaml:"key,omitempty"`
	CreatedAt *int64     `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	SNIs      []kong.SNI `json:"snis"`
}

// Consumer is a Kong consumer, and plugins and credentials associated with it.
type Consumer struct {
	kong.Consumer
	Plugins     []kong.Plugin            `json:"plugins"`
	KeyAuths    []*kong.KeyAuth          `json:"keyauth_credentials,omitempty"`
	HMACAuths   []*kong.HMACAuth         `json:"hmacauth_credentials,omitempty"`
	JWTAuths    []*kong.JWTAuth          `json:"jwt_secrets,omitempty"`
	BasicAuths  []*kong.BasicAuth        `json:"basicauth_credentials,omitempty"`
	ACLGroups   []*kong.ACLGroup         `json:"acls,omitempty"`
	Oauth2Creds []*kong.Oauth2Credential `json:"oauth2_credentials,omitempty"`
}

// KongDeclarativeConfig holds Kong's configuration and can be marshalled
// into Kong's native declarative configuration.
type KongDeclarativeConfig struct {
	FormatVersion string        `json:"_format_version"`
	Services      []Service     `json:"services"`
	Upstreams     []Upstream    `json:"upstreams"`
	Certificates  []Certificate `json:"certificates"`
	Plugins       []kong.Plugin `json:"plugins"`
	Consumers     []Consumer    `json:"consumers"`
}

// KongNativeState takes in a parser state and spits out Kong's native
// declarative configuration format.
func KongNativeState(k8sState *parser.KongState) *KongDeclarativeConfig {
	var result KongDeclarativeConfig
	result.FormatVersion = "1.1"
	if k8sState == nil {
		return &result
	}
	for _, s := range k8sState.Services {
		service := Service{Service: s.Service}

		for _, p := range s.Plugins {
			service.Plugins = append(service.Plugins, *p.DeepCopy())
		}

		for _, r := range s.Routes {
			route := Route{Route: r.Route}

			for _, p := range r.Plugins {
				route.Plugins = append(route.Plugins, *p.DeepCopy())
			}
			service.Routes = append(service.Routes, route)
		}
		result.Services = append(result.Services, service)
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
		consumer := Consumer{Consumer: c.Consumer,
			Plugins:    c.Plugins,
			KeyAuths:   c.KeyAuths,
			HMACAuths:  c.HMACAuths,
			ACLGroups:  c.ACLGroups,
			BasicAuths: c.BasicAuths,
			JWTAuths:   c.JWTAuths,
		}
		result.Consumers = append(result.Consumers, consumer)
	}
	return &result
}
