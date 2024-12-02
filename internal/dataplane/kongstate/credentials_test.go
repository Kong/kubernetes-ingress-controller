package kongstate_test

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
)

func TestKeyAuth_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   kongstate.KeyAuth
		want kongstate.KeyAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes key",
			in: kongstate.KeyAuth{
				KeyAuth: kong.KeyAuth{
					Consumer:  &kong.Consumer{Username: kong.String("foo")},
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Key:       kong.String("3"),
					Tags:      []*string{kong.String("4.1"), kong.String("4.2")},
				},
			},
			want: kongstate.KeyAuth{
				KeyAuth: kong.KeyAuth{
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Key:       kong.String("{vault://52fdfc07-2182-454f-963f-5f0f9a621d72}"),
					Tags:      []*string{kong.String("4.1"), kong.String("4.2")},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy(kongstate.StaticUUIDGenerator{UUID: "52fdfc07-2182-454f-963f-5f0f9a621d72"})
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHMACAuth_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   kongstate.HMACAuth
		want kongstate.HMACAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes secret",
			in: kongstate.HMACAuth{
				HMACAuth: kong.HMACAuth{
					Consumer:  &kong.Consumer{Username: kong.String("foo")},
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Username:  kong.String("3"),
					Secret:    kong.String("4"),
					Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
				},
			},
			want: kongstate.HMACAuth{
				HMACAuth: kong.HMACAuth{
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Username:  kong.String("3"),
					Secret:    kongstate.RedactedString,
					Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJWTAuth_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   kongstate.JWTAuth
		want kongstate.JWTAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes secret",
			in: kongstate.JWTAuth{
				JWTAuth: kong.JWTAuth{
					Consumer:     &kong.Consumer{Username: kong.String("foo")},
					CreatedAt:    kong.Int(1),
					ID:           kong.String("2"),
					Algorithm:    kong.String("3"),
					Key:          kong.String("4"),
					RSAPublicKey: kong.String("5"),
					Secret:       kong.String("6"),
					Tags:         []*string{kong.String("7.1"), kong.String("7.2")},
				},
			},
			want: kongstate.JWTAuth{
				JWTAuth: kong.JWTAuth{
					CreatedAt:    kong.Int(1),
					ID:           kong.String("2"),
					Algorithm:    kong.String("3"),
					Key:          kong.String("4"),
					RSAPublicKey: kong.String("5"),
					Secret:       kongstate.RedactedString,
					Tags:         []*string{kong.String("7.1"), kong.String("7.2")},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBasicAuth_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   kongstate.BasicAuth
		want kongstate.BasicAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes password",
			in: kongstate.BasicAuth{
				BasicAuth: kong.BasicAuth{
					Consumer:  &kong.Consumer{Username: kong.String("foo")},
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Username:  kong.String("3"),
					Password:  kong.String("4"),
					Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
				},
			},
			want: kongstate.BasicAuth{
				BasicAuth: kong.BasicAuth{
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Username:  kong.String("3"),
					Password:  kongstate.RedactedString,
					Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOauth2Credential_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   kongstate.Oauth2Credential
		want kongstate.Oauth2Credential
	}{
		{
			name: "fills all fields but Consumer and sanitizes client secret",
			in: kongstate.Oauth2Credential{
				Oauth2Credential: kong.Oauth2Credential{
					Consumer:     &kong.Consumer{Username: kong.String("foo")},
					CreatedAt:    kong.Int(1),
					ID:           kong.String("2"),
					Name:         kong.String("3"),
					ClientID:     kong.String("4"),
					ClientSecret: kong.String("5"),
					RedirectURIs: []*string{kong.String("6.1"), kong.String("6.2")},
					Tags:         []*string{kong.String("7.1"), kong.String("7.2")},
				},
			},
			want: kongstate.Oauth2Credential{
				Oauth2Credential: kong.Oauth2Credential{
					CreatedAt:    kong.Int(1),
					ID:           kong.String("2"),
					Name:         kong.String("3"),
					ClientID:     kong.String("4"),
					ClientSecret: kongstate.RedactedString,
					RedirectURIs: []*string{kong.String("6.1"), kong.String("6.2")},
					Tags:         []*string{kong.String("7.1"), kong.String("7.2")},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy()
			assert.Equal(t, tt.want, got)
		})
	}
}
