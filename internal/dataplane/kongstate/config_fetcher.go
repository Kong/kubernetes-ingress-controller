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
		})
		for _, r := range routes[*s.ID] {
			kongState.Services[i].Routes = append(kongState.Services[i].Routes, Route{
				Route: sanitizeKongRoute(*r),
			})
		}
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

	return kongState
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
			Certificate: *c,
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
		certs = append(certs, *c)
	}
	return certs
}

func sanitizeKongService(service kong.Service) kong.Service {
	service.CreatedAt = nil
	service.ID = nil
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
