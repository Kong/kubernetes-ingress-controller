package configfetcher

import (
	"github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
)

// -----------------------------------------------------------------------------
// KongRawState to KongState conversion functions
// -----------------------------------------------------------------------------

// KongRawStateToKongState converts a Deck kongRawState to a KIC KongState.
func KongRawStateToKongState(rawstate *utils.KongRawState) *kongstate.KongState {
	kongState := &kongstate.KongState{}
	if rawstate == nil {
		return kongState
	}

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

	for _, cg := range rawstate.ConsumerGroups {
		kongState.ConsumerGroups = append(kongState.ConsumerGroups, kongstate.ConsumerGroup{ConsumerGroup: sanitizeConsumerGroup(*cg.ConsumerGroup)})
	}

	targets := make(map[string][]*kong.Target)
	for _, u := range rawstate.Targets {
		if u.Upstream != nil && u.Upstream.ID != nil {
			targets[*u.Upstream.ID] = append(targets[*u.Upstream.ID], u)
		}
	}

	for i, s := range rawstate.Services {
		kongState.Services = append(kongState.Services, kongstate.Service{
			Service: sanitizeKongService(*s),
			Routes:  []kongstate.Route{},
			Plugins: []kong.Plugin{},
		})
		for j, r := range routes[*s.ID] {
			kongState.Services[i].Routes = append(kongState.Services[i].Routes, kongstate.Route{
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
		newUpstream := kongstate.Upstream{
			Upstream: *u,
		}
		if u.ID != nil {
			newUpstream.Targets = rawTargetsToTargets(targets[*u.ID])
		}
		kongState.Upstreams = append(kongState.Upstreams, sanitizeUpstream(newUpstream))
	}

	kongState.Vaults = rawVaultsToVaults(rawstate.Vaults)

	kongState.CACertificates = rawCACertificatesToCACertificates(rawstate.CACertificates)
	kongState.Certificates = rawCertificatesToCertificates(rawstate.Certificates)

	// Build lookup maps from consumer ID to credential slices to avoid O(n²) nested loops
	// when associating credentials with consumers (critical for large consumer counts).
	keyAuthsByConsumerID := make(map[string][]*kong.KeyAuth, len(rawstate.KeyAuths))
	for _, keyAuth := range rawstate.KeyAuths {
		if keyAuth.Consumer != nil && keyAuth.Consumer.ID != nil {
			consumerID := *keyAuth.Consumer.ID
			sanitizeAuth(keyAuth)
			keyAuthsByConsumerID[consumerID] = append(keyAuthsByConsumerID[consumerID], keyAuth)
		}
	}
	hmacAuthsByConsumerID := make(map[string][]*kong.HMACAuth, len(rawstate.HMACAuths))
	for _, hmacAuth := range rawstate.HMACAuths {
		if hmacAuth.Consumer != nil && hmacAuth.Consumer.ID != nil {
			consumerID := *hmacAuth.Consumer.ID
			sanitizeAuth(hmacAuth)
			hmacAuthsByConsumerID[consumerID] = append(hmacAuthsByConsumerID[consumerID], hmacAuth)
		}
	}
	jwtAuthsByConsumerID := make(map[string][]*kong.JWTAuth, len(rawstate.JWTAuths))
	for _, jwtAuth := range rawstate.JWTAuths {
		if jwtAuth.Consumer != nil && jwtAuth.Consumer.ID != nil {
			consumerID := *jwtAuth.Consumer.ID
			sanitizeAuth(jwtAuth)
			jwtAuthsByConsumerID[consumerID] = append(jwtAuthsByConsumerID[consumerID], jwtAuth)
		}
	}
	basicAuthsByConsumerID := make(map[string][]*kong.BasicAuth, len(rawstate.BasicAuths))
	for _, basicAuth := range rawstate.BasicAuths {
		if basicAuth.Consumer != nil && basicAuth.Consumer.ID != nil {
			consumerID := *basicAuth.Consumer.ID
			sanitizeAuth(basicAuth)
			basicAuthsByConsumerID[consumerID] = append(basicAuthsByConsumerID[consumerID], basicAuth)
		}
	}
	aclGroupsByConsumerID := make(map[string][]*kong.ACLGroup, len(rawstate.ACLGroups))
	for _, aclGroup := range rawstate.ACLGroups {
		if aclGroup.Consumer != nil && aclGroup.Consumer.ID != nil {
			consumerID := *aclGroup.Consumer.ID
			sanitizeAuth(aclGroup)
			aclGroupsByConsumerID[consumerID] = append(aclGroupsByConsumerID[consumerID], aclGroup)
		}
	}
	oauth2CredsByConsumerID := make(map[string][]*kong.Oauth2Credential, len(rawstate.Oauth2Creds))
	for _, oauth2Cred := range rawstate.Oauth2Creds {
		if oauth2Cred.Consumer != nil && oauth2Cred.Consumer.ID != nil {
			consumerID := *oauth2Cred.Consumer.ID
			sanitizeAuth(oauth2Cred)
			oauth2CredsByConsumerID[consumerID] = append(oauth2CredsByConsumerID[consumerID], oauth2Cred)
		}
	}
	mTLSAuthsByConsumerID := make(map[string][]*kong.MTLSAuth, len(rawstate.MTLSAuths))
	for _, mTLSAuth := range rawstate.MTLSAuths {
		if mTLSAuth.Consumer != nil && mTLSAuth.Consumer.ID != nil {
			consumerID := *mTLSAuth.Consumer.ID
			sanitizeAuth(mTLSAuth)
			mTLSAuthsByConsumerID[consumerID] = append(mTLSAuthsByConsumerID[consumerID], mTLSAuth)
		}
	}

	for i, consumer := range rawstate.Consumers {
		kongState.Consumers = append(kongState.Consumers, kongstate.Consumer{
			Consumer: sanitizeConsumer(*consumer),
		})
		if consumer.ID == nil {
			continue
		}
		consumerID := *consumer.ID
		for _, keyAuth := range keyAuthsByConsumerID[consumerID] {
			kongState.Consumers[i].KeyAuths = append(kongState.Consumers[i].KeyAuths,
				&kongstate.KeyAuth{KeyAuth: *keyAuth},
			)
		}
		for _, hmacAuth := range hmacAuthsByConsumerID[consumerID] {
			kongState.Consumers[i].HMACAuths = append(kongState.Consumers[i].HMACAuths,
				&kongstate.HMACAuth{HMACAuth: *hmacAuth},
			)
		}
		for _, jwtAuth := range jwtAuthsByConsumerID[consumerID] {
			kongState.Consumers[i].JWTAuths = append(kongState.Consumers[i].JWTAuths,
				&kongstate.JWTAuth{JWTAuth: *jwtAuth},
			)
		}
		for _, basicAuth := range basicAuthsByConsumerID[consumerID] {
			kongState.Consumers[i].BasicAuths = append(kongState.Consumers[i].BasicAuths,
				&kongstate.BasicAuth{BasicAuth: *basicAuth},
			)
		}
		for _, aclGroup := range aclGroupsByConsumerID[consumerID] {
			kongState.Consumers[i].ACLGroups = append(kongState.Consumers[i].ACLGroups,
				&kongstate.ACLGroup{ACLGroup: *aclGroup},
			)
		}
		for _, oauth2Cred := range oauth2CredsByConsumerID[consumerID] {
			kongState.Consumers[i].Oauth2Creds = append(kongState.Consumers[i].Oauth2Creds,
				&kongstate.Oauth2Credential{Oauth2Credential: *oauth2Cred},
			)
		}
		for _, mTLSAuth := range mTLSAuthsByConsumerID[consumerID] {
			kongState.Consumers[i].MTLSAuths = append(kongState.Consumers[i].MTLSAuths,
				&kongstate.MTLSAuth{MTLSAuth: *mTLSAuth},
			)
		}
	}

	for _, entity := range rawstate.CustomEntities {
		entityType := entity.Type()
		obj := entity.Object()
		ksEntity := kongstate.CustomEntity{
			Object: obj,
		}
		kongState.AddCustomEntity(string(entityType), kongstate.EntitySchema{}, ksEntity)
	}

	return kongState
}

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

func rawTargetsToTargets(targets []*kong.Target) []kongstate.Target {
	if len(targets) == 0 {
		return nil
	}
	ts := []kongstate.Target{}

	for _, t := range targets {
		ts = append(ts, kongstate.Target{Target: *t})
	}
	return ts
}

func rawCertificatesToCertificates(certificates []*kong.Certificate) []kongstate.Certificate {
	if len(certificates) == 0 {
		return nil
	}
	certs := []kongstate.Certificate{}

	for _, c := range certificates {
		certs = append(certs, kongstate.Certificate{
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

func rawVaultsToVaults(rawVaults []*kong.Vault) []kongstate.Vault {
	if len(rawVaults) == 0 {
		return nil
	}
	vaults := []kongstate.Vault{}

	for _, v := range rawVaults {
		vaults = append(vaults, kongstate.Vault{
			Vault: sanitizeVault(*v),
		})
	}
	return vaults
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

func sanitizeUpstream(upstream kongstate.Upstream) kongstate.Upstream {
	upstream.CreatedAt = nil
	upstream.ID = nil
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

func sanitizeVault(v kong.Vault) kong.Vault {
	v.ID = nil
	v.CreatedAt = nil
	return v
}

func sanitizeConsumer(consumer kong.Consumer) kong.Consumer {
	consumer.ID = nil
	consumer.CreatedAt = nil
	return consumer
}

func sanitizeConsumerGroup(consumerGroup kong.ConsumerGroup) kong.ConsumerGroup {
	consumerGroup.ID = nil
	consumerGroup.CreatedAt = nil
	return consumerGroup
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
	switch a := (any)(auth).(type) {
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
