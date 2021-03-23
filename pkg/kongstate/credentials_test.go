package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
)

func TestKeyAuth_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   KeyAuth
		want KeyAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes key",
			in: KeyAuth{kong.KeyAuth{
				Consumer:  &kong.Consumer{Username: kong.String("foo")},
				CreatedAt: kong.Int(1),
				ID:        kong.String("2"),
				Key:       kong.String("3"),
				Tags:      []*string{kong.String("4.1"), kong.String("4.2")},
			}},
			want: KeyAuth{kong.KeyAuth{
				CreatedAt: kong.Int(1),
				ID:        kong.String("2"),
				Key:       redactedString,
				Tags:      []*string{kong.String("4.1"), kong.String("4.2")},
			}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHMACAuth_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   HMACAuth
		want HMACAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes secret",
			in: HMACAuth{kong.HMACAuth{
				Consumer:  &kong.Consumer{Username: kong.String("foo")},
				CreatedAt: kong.Int(1),
				ID:        kong.String("2"),
				Username:  kong.String("3"),
				Secret:    kong.String("4"),
				Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
			}},
			want: HMACAuth{kong.HMACAuth{
				CreatedAt: kong.Int(1),
				ID:        kong.String("2"),
				Username:  kong.String("3"),
				Secret:    redactedString,
				Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
			}},
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
		in   JWTAuth
		want JWTAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes secret",
			in: JWTAuth{kong.JWTAuth{
				Consumer:     &kong.Consumer{Username: kong.String("foo")},
				CreatedAt:    kong.Int(1),
				ID:           kong.String("2"),
				Algorithm:    kong.String("3"),
				Key:          kong.String("4"),
				RSAPublicKey: kong.String("5"),
				Secret:       kong.String("6"),
				Tags:         []*string{kong.String("7.1"), kong.String("7.2")},
			}},
			want: JWTAuth{kong.JWTAuth{
				CreatedAt:    kong.Int(1),
				ID:           kong.String("2"),
				Algorithm:    kong.String("3"),
				Key:          kong.String("4"),
				RSAPublicKey: kong.String("5"),
				Secret:       redactedString,
				Tags:         []*string{kong.String("7.1"), kong.String("7.2")},
			}},
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
		in   BasicAuth
		want BasicAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes password",
			in: BasicAuth{kong.BasicAuth{
				Consumer:  &kong.Consumer{Username: kong.String("foo")},
				CreatedAt: kong.Int(1),
				ID:        kong.String("2"),
				Username:  kong.String("3"),
				Password:  kong.String("4"),
				Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
			}},
			want: BasicAuth{kong.BasicAuth{
				CreatedAt: kong.Int(1),
				ID:        kong.String("2"),
				Username:  kong.String("3"),
				Password:  redactedString,
				Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
			}},
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
		in   Oauth2Credential
		want Oauth2Credential
	}{
		{
			name: "fills all fields but Consumer and sanitizes client secret",
			in: Oauth2Credential{kong.Oauth2Credential{
				Consumer:     &kong.Consumer{Username: kong.String("foo")},
				CreatedAt:    kong.Int(1),
				ID:           kong.String("2"),
				Name:         kong.String("3"),
				ClientID:     kong.String("4"),
				ClientSecret: kong.String("5"),
				RedirectURIs: []*string{kong.String("6.1"), kong.String("6.2")},
				Tags:         []*string{kong.String("7.1"), kong.String("7.2")},
			}},
			want: Oauth2Credential{kong.Oauth2Credential{
				CreatedAt:    kong.Int(1),
				ID:           kong.String("2"),
				Name:         kong.String("3"),
				ClientID:     kong.String("4"),
				ClientSecret: redactedString,
				RedirectURIs: []*string{kong.String("6.1"), kong.String("6.2")},
				Tags:         []*string{kong.String("7.1"), kong.String("7.2")},
			}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy()
			assert.Equal(t, tt.want, got)
		})
	}
}
