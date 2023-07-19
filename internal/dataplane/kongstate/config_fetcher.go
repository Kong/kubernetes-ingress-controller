package kongstate

import (
	"context"

	"github.com/kong/deck/dump"
	"github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
)

type KongLastGoodConfigFetcher interface {
	// GetKongRawState fetches the configuration loaded in a Kong node.
	GetKongRawState(ctx context.Context, client *kong.Client) (*utils.KongRawState, error)
	// GetKongStatus fetches the status of a Kong node.
	GetKongStatus(ctx context.Context, client *kong.Client) (*kong.Status, error)
}

type DefaultKongLastGoodConfigFetcher struct {
	config dump.Config
}

func (g *DefaultKongLastGoodConfigFetcher) GetKongRawState(ctx context.Context, client *kong.Client) (*utils.KongRawState, error) {
	return dump.Get(ctx, client, g.config)
}

func (g *DefaultKongLastGoodConfigFetcher) GetKongStatus(ctx context.Context, client *kong.Client) (*kong.Status, error) {
	return client.Status(ctx)
}

func NewDefaultKongLastGoodConfigFetcher() *DefaultKongLastGoodConfigFetcher {
	return &DefaultKongLastGoodConfigFetcher{
		config: dump.Config{},
	}
}

// KongRawStateToKongState converts a Deck kongRawState to a KIC KongState.
func KongRawStateToKongState(rawstate *utils.KongRawState) *KongState {
	kongState := &KongState{}

	routes := make(map[string][]*kong.Route)
	for _, r := range rawstate.Routes {
		if r.Service != nil && r.Service.ID != nil {
			routes[*r.Service.ID] = append(routes[*r.Service.ID], r)
		}
	}

	pluginsByService := make(map[string][]*kong.Plugin)
	pluginsByRoute := make(map[string][]*kong.Plugin)
	for _, p := range rawstate.Plugins {
		if p.Service != nil && p.Service.ID != nil {
			pluginsByService[*p.Service.ID] = append(pluginsByService[*p.Service.ID], p)
		}
		if p.Route != nil && p.Route.ID != nil {
			pluginsByRoute[*p.Route.ID] = append(pluginsByRoute[*p.Route.ID], p)
		}
	}

	targets := make(map[string][]*kong.Target)
	for _, u := range rawstate.Targets {
		if u.Upstream != nil && u.Upstream.ID != nil {
			targets[*u.Upstream.ID] = append(targets[*u.Upstream.ID], u)
		}
	}

	for i, s := range rawstate.Services {
		kongState.Services = append(kongState.Services, Service{
			Service: sanitizeKongService(*s),
			Routes:  []Route{},
			Plugins: []kong.Plugin{},
		})
		for j, r := range routes[*s.ID] {
			kongState.Services[i].Routes = append(kongState.Services[i].Routes, Route{
				Route:   sanitizeKongRoute(*r),
				Plugins: []kong.Plugin{},
			})
			if r.ID != nil {
				kongState.Services[i].Routes[j].Plugins = rawPluginsToPlugins(pluginsByRoute[*r.ID])
			}
		}
		kongState.Services[i].Plugins = rawPluginsToPlugins(pluginsByService[*s.ID])
	}

	for _, u := range rawstate.Upstreams {
		newUpstream := Upstream{
			Upstream: *u,
		}
		if u.ID != nil {
			newUpstream.Targets = rawTargetsToTargets(targets[*u.ID])
		}
		kongState.Upstreams = append(kongState.Upstreams, sanitizeUpstream(newUpstream))
	}

	kongState.CACertificates = rawCACertificatesToCACertificates(rawstate.CACertificates)
	kongState.Certificates = rawCertificatesToCertificates(rawstate.Certificates)

	for i, consumer := range rawstate.Consumers {
		kongState.Consumers = append(kongState.Consumers, Consumer{
			Consumer: sanitizeConsumer(*consumer),
		})
		for _, keyAuth := range rawstate.KeyAuths {
			if keyAuth.Consumer != nil {
				if *keyAuth.Consumer.ID == *consumer.ID {
					sanitizeAuth(keyAuth)
					kongState.Consumers[i].KeyAuths = append(kongState.Consumers[i].KeyAuths,
						&KeyAuth{
							KeyAuth: *keyAuth,
						},
					)
				}
			}
		}
		for _, hmacAuth := range rawstate.HMACAuths {
			if hmacAuth.Consumer != nil {
				if *hmacAuth.Consumer.ID == *consumer.ID {
					sanitizeAuth(hmacAuth)
					kongState.Consumers[i].HMACAuths = append(kongState.Consumers[i].HMACAuths,
						&HMACAuth{
							HMACAuth: *hmacAuth,
						},
					)
				}
			}
		}
		for _, jwtAuth := range rawstate.JWTAuths {
			if jwtAuth.Consumer != nil {
				if *jwtAuth.Consumer.ID == *consumer.ID {
					sanitizeAuth(jwtAuth)
					kongState.Consumers[i].JWTAuths = append(kongState.Consumers[i].JWTAuths,
						&JWTAuth{
							JWTAuth: *jwtAuth,
						},
					)
				}
			}
		}
		for _, basicAuth := range rawstate.BasicAuths {
			if basicAuth.Consumer != nil {
				if *basicAuth.Consumer.ID == *consumer.ID {
					sanitizeAuth(basicAuth)
					kongState.Consumers[i].BasicAuths = append(kongState.Consumers[i].BasicAuths,
						&BasicAuth{
							BasicAuth: *basicAuth,
						},
					)
				}
			}
		}
		for _, aclGroup := range rawstate.ACLGroups {
			if aclGroup.Consumer != nil {
				if *aclGroup.Consumer.ID == *consumer.ID {
					sanitizeAuth(aclGroup)
					kongState.Consumers[i].ACLGroups = append(kongState.Consumers[i].ACLGroups,
						&ACLGroup{
							ACLGroup: *aclGroup,
						},
					)
				}
			}
		}
		for _, oauth2Cred := range rawstate.Oauth2Creds {
			if oauth2Cred.Consumer != nil {
				if *oauth2Cred.Consumer.ID == *consumer.ID {
					sanitizeAuth(oauth2Cred)
					kongState.Consumers[i].Oauth2Creds = append(kongState.Consumers[i].Oauth2Creds,
						&Oauth2Credential{
							Oauth2Credential: *oauth2Cred,
						},
					)
				}
			}
		}
		for _, mTLSAuth := range rawstate.MTLSAuths {
			if mTLSAuth.Consumer != nil {
				if *mTLSAuth.Consumer.ID == *consumer.ID {
					sanitizeAuth(mTLSAuth)
					kongState.Consumers[i].MTLSAuths = append(kongState.Consumers[i].MTLSAuths,
						&MTLSAuth{
							MTLSAuth: *mTLSAuth,
						},
					)
				}
			}
		}
	}

	return kongState
}

// -----------------------------------------------------------------------------
// KongRawState to KongState conversion functions
// -----------------------------------------------------------------------------

func rawPluginsToPlugins(plugins []*kong.Plugin) []kong.Plugin {
	if len(plugins) == 0 {
		return nil
	}
	ps := []kong.Plugin{}

	for _, p := range plugins {
		ps = append(ps, sanitizePlugin(*p))
	}
	return ps
}

func rawTargetsToTargets(targets []*kong.Target) []Target {
	if len(targets) == 0 {
		return nil
	}
	ts := []Target{}

	for _, t := range targets {
		ts = append(ts, Target{Target: *t})
	}
	return ts
}

func rawCertificatesToCertificates(certificates []*kong.Certificate) []Certificate {
	if len(certificates) == 0 {
		return nil
	}
	certs := []Certificate{}

	for _, c := range certificates {
		certs = append(certs, Certificate{
			Certificate: sanitizeCertificate(*c),
		})
	}
	return certs
}

func rawCACertificatesToCACertificates(caCertificates []*kong.CACertificate) []kong.CACertificate {
	if len(caCertificates) == 0 {
		return nil
	}
	certs := []kong.CACertificate{}

	for _, c := range caCertificates {
		certs = append(certs, sanitizeCACertificate(*c))
	}
	return certs
}

// -----------------------------------------------------------------------------
// Sanitization functions
// -----------------------------------------------------------------------------

func sanitizeKongService(service kong.Service) kong.Service {
	service.ID = nil
	service.CreatedAt = nil
	service.UpdatedAt = nil
	return service
}

func sanitizeKongRoute(route kong.Route) kong.Route {
	route.CreatedAt = nil
	route.ID = nil
	route.UpdatedAt = nil
	route.Service = nil
	return route
}

func sanitizeUpstream(upstream Upstream) Upstream {
	upstream.Upstream.CreatedAt = nil
	upstream.Upstream.ID = nil
	for i := range upstream.Targets {
		upstream.Targets[i].CreatedAt = nil
		upstream.Targets[i].ID = nil
		upstream.Targets[i].Upstream = nil
	}
	return upstream
}

func sanitizePlugin(plugin kong.Plugin) kong.Plugin {
	plugin.ID = nil
	plugin.CreatedAt = nil
	plugin.Service = nil
	plugin.Route = nil
	return plugin
}

func sanitizeCertificate(certificate kong.Certificate) kong.Certificate {
	certificate.ID = nil
	certificate.CreatedAt = nil
	return certificate
}

func sanitizeCACertificate(caCertificate kong.CACertificate) kong.CACertificate {
	caCertificate.ID = nil
	caCertificate.CreatedAt = nil
	return caCertificate
}

func sanitizeConsumer(consumer kong.Consumer) kong.Consumer {
	consumer.ID = nil
	consumer.CreatedAt = nil
	return consumer
}

type authT interface {
	*kong.KeyAuth |
		*kong.HMACAuth |
		*kong.JWTAuth |
		*kong.BasicAuth |
		*kong.ACLGroup |
		*kong.Oauth2Credential |
		*kong.MTLSAuth
}

func sanitizeAuth[t authT](auth t) {
	switch a := (interface{})(auth).(type) {
	case *kong.KeyAuth:
		a.ID = nil
		a.CreatedAt = nil
		a.Consumer = nil
	case *kong.HMACAuth:
		a.ID = nil
		a.CreatedAt = nil
		a.Consumer = nil
	case *kong.JWTAuth:
		a.ID = nil
		a.CreatedAt = nil
		a.Consumer = nil
	case *kong.BasicAuth:
		a.ID = nil
		a.CreatedAt = nil
		a.Consumer = nil
	case *kong.ACLGroup:
		a.ID = nil
		a.CreatedAt = nil
		a.Consumer = nil
	case *kong.Oauth2Credential:
		a.ID = nil
		a.CreatedAt = nil
		a.Consumer = nil
	case *kong.MTLSAuth:
		a.ID = nil
		a.CreatedAt = nil
		a.Consumer = nil
		if a.CACertificate != nil {
			a.CACertificate.ID = nil
			a.CACertificate.CreatedAt = nil
		}
	}
}
