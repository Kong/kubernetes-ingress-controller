package configfetcher_test

import (
	"reflect"
	"testing"

	"github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/configfetcher"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
)

func TestKongRawStateToKongState(t *testing.T) {
	for _, tt := range []struct {
		name              string
		kongRawState      utils.KongRawState
		expectedKongState *kongstate.KongState
	}{
		{
			name: "sanitizes all services, routes, and upstreams and create a KongState out of a KongRawState",
			kongRawState: utils.KongRawState{
				Services: []*kong.Service{
					{
						Name:      kong.String("service"),
						ID:        kong.String("service"),
						CreatedAt: kong.Int(100),
					},
				},
				Routes: []*kong.Route{
					{
						Name:      kong.String("route"),
						ID:        kong.String("route"),
						CreatedAt: kong.Int(101),
						Service: &kong.Service{
							ID: kong.String("service"),
						},
					},
				},
				Upstreams: []*kong.Upstream{
					{
						Name: kong.String("upstream"),
						ID:   kong.String("upstream"),
					},
				},
				Targets: []*kong.Target{
					{
						ID:        kong.String("target"),
						CreatedAt: kong.Float64(102),
						Weight:    kong.Int(999),
						Upstream: &kong.Upstream{
							ID: kong.String("upstream"),
						},
					},
				},
				Plugins: []*kong.Plugin{
					{
						Name: kong.String("plugin1"),
						ID:   kong.String("plugin1"),
						Service: &kong.Service{
							ID: kong.String("service"),
						},
					},
					{
						Name: kong.String("plugin2"),
						ID:   kong.String("plugin2"),
						Route: &kong.Route{
							ID: kong.String("route"),
						},
					},
				},
				Certificates: []*kong.Certificate{
					{
						ID:   kong.String("certificate"),
						Cert: kong.String("cert"),
					},
				},
				CACertificates: []*kong.CACertificate{
					{
						ID:   kong.String("CACertificate"),
						Cert: kong.String("cert"),
					},
				},
				Consumers: []*kong.Consumer{
					{
						ID:       kong.String("consumer"),
						CustomID: kong.String("customID"),
					},
				},
				KeyAuths: []*kong.KeyAuth{
					{
						ID:  kong.String("keyAuth"),
						Key: kong.String("key"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer"),
						},
					},
				},
				HMACAuths: []*kong.HMACAuth{
					{
						ID: kong.String("hmacAuth"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer"),
						},
						Username: kong.String("username"),
					},
				},
				JWTAuths: []*kong.JWTAuth{
					{
						ID: kong.String("jwtAuth"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer"),
						},
						Key: kong.String("key"),
					},
				},
				BasicAuths: []*kong.BasicAuth{
					{
						ID: kong.String("basicAuth"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer"),
						},
						Username: kong.String("username"),
					},
				},
				ACLGroups: []*kong.ACLGroup{
					{
						ID: kong.String("basicAuth"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer"),
						},
						Group: kong.String("group"),
					},
				},
				Oauth2Creds: []*kong.Oauth2Credential{
					{
						ID: kong.String("basicAuth"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer"),
						},
						Name: kong.String("name"),
					},
				},
				MTLSAuths: []*kong.MTLSAuth{
					{
						ID: kong.String("basicAuth"),
						Consumer: &kong.Consumer{
							ID: kong.String("consumer"),
						},
						SubjectName: kong.String("subjectName"),
					},
				},
			},
			expectedKongState: &kongstate.KongState{
				Services: []kongstate.Service{
					{
						Service: kong.Service{
							Name: kong.String("service"),
						},
						Plugins: []kong.Plugin{
							{
								Name: kong.String("plugin1"),
							},
						},
						Routes: []kongstate.Route{
							{
								Route: kong.Route{
									Name: kong.String("route"),
								},
								Plugins: []kong.Plugin{
									{
										Name: kong.String("plugin2"),
									},
								},
							},
						},
					},
				},
				Upstreams: []kongstate.Upstream{
					{
						Upstream: kong.Upstream{
							Name: kong.String("upstream"),
						},
						Targets: []kongstate.Target{
							{
								Target: kong.Target{
									Weight: kong.Int(999),
								},
							},
						},
					},
				},
				Certificates: []kongstate.Certificate{
					{
						Certificate: kong.Certificate{
							Cert: kong.String("cert"),
						},
					},
				},
				CACertificates: []kong.CACertificate{
					{
						Cert: kong.String("cert"),
					},
				},
				Consumers: []kongstate.Consumer{
					{
						Consumer: kong.Consumer{
							CustomID: kong.String("customID"),
						},
						KeyAuths: []*kongstate.KeyAuth{
							{
								KeyAuth: kong.KeyAuth{
									Key: kong.String("key"),
								},
							},
						},
						HMACAuths: []*kongstate.HMACAuth{
							{
								HMACAuth: kong.HMACAuth{
									Username: kong.String("username"),
								},
							},
						},
						JWTAuths: []*kongstate.JWTAuth{
							{
								JWTAuth: kong.JWTAuth{
									Key: kong.String("key"),
								},
							},
						},
						BasicAuths: []*kongstate.BasicAuth{
							{
								BasicAuth: kong.BasicAuth{
									Username: kong.String("username"),
								},
							},
						},
						ACLGroups: []*kongstate.ACLGroup{
							{
								ACLGroup: kong.ACLGroup{
									Group: kong.String("group"),
								},
							},
						},
						Oauth2Creds: []*kongstate.Oauth2Credential{
							{
								Oauth2Credential: kong.Oauth2Credential{
									Name: kong.String("name"),
								},
							},
						},
						MTLSAuths: []*kongstate.MTLSAuth{
							{
								MTLSAuth: kong.MTLSAuth{
									SubjectName: kong.String("subjectName"),
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			state := configfetcher.KongRawStateToKongState(&tt.kongRawState)
			require.Equal(t, tt.expectedKongState, state)
		})
	}
}

func TestKongStateToKongRawState_Ensure(t *testing.T) {
	kongRawStateFieldsKICDoesntSupport := []string{
		// These are fields that KIC explicitly doesn't support.
		"SNIs",
		"ConsumerGroups",
		"CustomEntities",
		"Vaults",
		"RBACRoles",
		"RBACEndpointPermissions",
	}
	allKongRawStateFields := func() []string {
		var fields []string
		typ := reflect.ValueOf(utils.KongRawState{}).Type()
		for i := 0; i < typ.NumField(); i++ {
			fields = append(fields, typ.Field(i).Name)
		}
		return fields
	}()

	testedFields := []string{
		"Services",
		"Routes",
		"Upstreams",
		"Targets",
		"Plugins",
		"Certificates",
		"CACertificates",
		"Consumers",
		"KeyAuths",
		"HMACAuths",
		"JWTAuths",
		"BasicAuths",
		"ACLGroups",
		"Oauth2Creds",
		"MTLSAuths",
	}

	// Meta test - ensure we have testcases covering all fields in KongRawState.
	for _, field := range allKongRawStateFields {
		if lo.Contains(kongRawStateFieldsKICDoesntSupport, field) {
			t.Logf("skipping field %s - explicitly unsupported", field)
			continue
		}
		assert.True(t, lo.Contains(testedFields, field), "field %s unsupported", field)
	}
}
