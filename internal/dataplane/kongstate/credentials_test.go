package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func TestKeyAuth_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   KeyAuth
		want KeyAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes key",
			in: KeyAuth{
				KeyAuth: kong.KeyAuth{
					Consumer:  &kong.Consumer{Username: kong.String("foo")},
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Key:       kong.String("3"),
					Tags:      []*string{kong.String("4.1"), kong.String("4.2")},
				},
			},
			want: KeyAuth{
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
			got := *tt.in.SanitizedCopy(mocks.StaticUUIDGenerator{UUID: "52fdfc07-2182-454f-963f-5f0f9a621d72"})
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
			in: HMACAuth{
				HMACAuth: kong.HMACAuth{
					Consumer:  &kong.Consumer{Username: kong.String("foo")},
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Username:  kong.String("3"),
					Secret:    kong.String("4"),
					Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
				},
			},
			want: HMACAuth{
				HMACAuth: kong.HMACAuth{
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Username:  kong.String("3"),
					Secret:    redactedString,
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
		in   JWTAuth
		want JWTAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes secret",
			in: JWTAuth{
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
			want: JWTAuth{
				JWTAuth: kong.JWTAuth{
					CreatedAt:    kong.Int(1),
					ID:           kong.String("2"),
					Algorithm:    kong.String("3"),
					Key:          kong.String("4"),
					RSAPublicKey: kong.String("5"),
					Secret:       redactedString,
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
		in   BasicAuth
		want BasicAuth
	}{
		{
			name: "fills all fields but Consumer and sanitizes password",
			in: BasicAuth{
				BasicAuth: kong.BasicAuth{
					Consumer:  &kong.Consumer{Username: kong.String("foo")},
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Username:  kong.String("3"),
					Password:  kong.String("4"),
					Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
				},
			},
			want: BasicAuth{
				BasicAuth: kong.BasicAuth{
					CreatedAt: kong.Int(1),
					ID:        kong.String("2"),
					Username:  kong.String("3"),
					Password:  redactedString,
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
		in   Oauth2Credential
		want Oauth2Credential
	}{
		{
			name: "fills all fields but Consumer and sanitizes client secret",
			in: Oauth2Credential{
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
			want: Oauth2Credential{
				Oauth2Credential: kong.Oauth2Credential{
					CreatedAt:    kong.Int(1),
					ID:           kong.String("2"),
					Name:         kong.String("3"),
					ClientID:     kong.String("4"),
					ClientSecret: redactedString,
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

func TestGetKeyAuthsConflictingOnKey(t *testing.T) {
	testCases := []struct {
		name              string
		keyAuths          []*KeyAuth
		expectedConflicts int
	}{
		{
			name: "no conflict",
			keyAuths: []*KeyAuth{
				{
					KeyAuth: kong.KeyAuth{
						Key: kong.String("key1"),
					},
				},
				{
					KeyAuth: kong.KeyAuth{
						Key: kong.String("key2"),
					},
				},
			},
			expectedConflicts: 0,
		},
		{
			name: "2 conflict and 1 no conflict",
			keyAuths: []*KeyAuth{
				{
					KeyAuth: kong.KeyAuth{
						Key: kong.String("key1"),
					},
				},
				{
					KeyAuth: kong.KeyAuth{
						Key: kong.String("key1"),
					},
				},
				{
					KeyAuth: kong.KeyAuth{
						Key: kong.String("key2"),
					},
				},
			},
			expectedConflicts: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conflictingKeyAuths := getKeyAuthsConflictingOnKey(tc.keyAuths)
			assert.Len(t, conflictingKeyAuths, tc.expectedConflicts)
		})
	}
}

func TestGetHMACAuthsConflictingOnUsername(t *testing.T) {
	testCases := []struct {
		name              string
		hmacAuths         []*HMACAuth
		expectedConflicts int
	}{
		{
			name: "no conflict",
			hmacAuths: []*HMACAuth{
				{
					HMACAuth: kong.HMACAuth{
						Username: kong.String("user1"),
						Secret:   kong.String("whatever"),
					},
				},
				{
					HMACAuth: kong.HMACAuth{
						Username: kong.String("user2"),
						Secret:   kong.String("whatever"),
					},
				},
			},
			expectedConflicts: 0,
		},
		{
			name: "2 conflict and 1 no conflict",
			hmacAuths: []*HMACAuth{
				{
					HMACAuth: kong.HMACAuth{
						Username: kong.String("user1"),
						Secret:   kong.String("wonderful"),
					},
				},
				{
					HMACAuth: kong.HMACAuth{
						Username: kong.String("user1"),
						Secret:   kong.String("terrible"),
					},
				},
				{
					HMACAuth: kong.HMACAuth{
						Username: kong.String("user2"),
						Secret:   kong.String("whatever"),
					},
				},
			},
			expectedConflicts: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conflictingHMACAuths := getHMACAuthsConflictingOnUsername(tc.hmacAuths)
			assert.Len(t, conflictingHMACAuths, tc.expectedConflicts)
		})
	}
}

func TestGetJWTAuthsConflictingOnKey(t *testing.T) {
	testCases := []struct {
		name              string
		jwtAuths          []*JWTAuth
		expectedConflicts int
	}{
		{
			name: "no conflict",
			jwtAuths: []*JWTAuth{
				{
					JWTAuth: kong.JWTAuth{
						Key:          kong.String("key1"),
						Algorithm:    kong.String("HS256"),
						RSAPublicKey: kong.String("----- BEGIN PUBLIC KEY"),
					},
				},
				{
					JWTAuth: kong.JWTAuth{
						Key:          kong.String("key2"),
						Algorithm:    kong.String("HS256"),
						RSAPublicKey: kong.String("----- BEGIN PUBLIC KEY"),
					},
				},
			},
			expectedConflicts: 0,
		},
		{
			name: "2 conflict and 1 no conflict",
			jwtAuths: []*JWTAuth{
				{
					JWTAuth: kong.JWTAuth{
						Key:          kong.String("key1"),
						Algorithm:    kong.String("HS256"),
						RSAPublicKey: kong.String("----- BEGIN PUBLIC KEY"),
					},
				},
				{
					JWTAuth: kong.JWTAuth{
						Key:          kong.String("key1"),
						Algorithm:    kong.String("HS256"),
						RSAPublicKey: kong.String("----- BEGIN PUBLIC KEY"),
					},
				},
				{
					JWTAuth: kong.JWTAuth{
						Key:          kong.String("key2"),
						Algorithm:    kong.String("HS256"),
						RSAPublicKey: kong.String("----- BEGIN PUBLIC KEY"),
					},
				},
			},
			expectedConflicts: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conflictingJWTAuths := getJWTAuthsConflictingOnKey(tc.jwtAuths)
			assert.Len(t, conflictingJWTAuths, tc.expectedConflicts)
		})
	}
}

func TestGetBasicAuthsConflictingOnUsername(t *testing.T) {
	testCases := []struct {
		name              string
		basicAuths        []*BasicAuth
		expectedConflicts int
	}{
		{
			name: "no conflict",
			basicAuths: []*BasicAuth{
				{
					BasicAuth: kong.BasicAuth{
						Username: kong.String("user1"),
						Password: kong.String("123456"),
					},
				},
				{
					BasicAuth: kong.BasicAuth{
						Username: kong.String("user2"),
						Password: kong.String("234567"),
					},
				},
			},
			expectedConflicts: 0,
		},
		{
			name: "2 conflict and 1 no conflict",
			basicAuths: []*BasicAuth{
				{
					BasicAuth: kong.BasicAuth{
						Username: kong.String("user1"),
						Password: kong.String("123456"),
					},
				},
				{
					BasicAuth: kong.BasicAuth{
						Username: kong.String("user1"),
						Password: kong.String("098765"),
					},
				},
				{
					BasicAuth: kong.BasicAuth{
						Username: kong.String("user2"),
						Password: kong.String("234567"),
					},
				},
			},
			expectedConflicts: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conflictingBasicAuths := getBasicAuthsConflictingOnUsername(tc.basicAuths)
			assert.Len(t, conflictingBasicAuths, tc.expectedConflicts)
		})
	}
}

func TestGetOAuth2CredentialsConflictingOnClientID(t *testing.T) {
	testCases := []struct {
		name              string
		oauth2Creds       []*Oauth2Credential
		expectedConflicts int
	}{
		{
			name: "no conflict",
			oauth2Creds: []*Oauth2Credential{
				{
					Oauth2Credential: kong.Oauth2Credential{
						ClientID:     kong.String("client-1"),
						ClientSecret: kong.String("client-secret-1"),
					},
				},
				{
					Oauth2Credential: kong.Oauth2Credential{
						ClientID:     kong.String("client-2"),
						ClientSecret: kong.String("client-secret-2"),
					},
				},
			},
			expectedConflicts: 0,
		},
		{
			name: "2 conflict and 1 no conflict",
			oauth2Creds: []*Oauth2Credential{
				{
					Oauth2Credential: kong.Oauth2Credential{
						ClientID:     kong.String("client-1"),
						ClientSecret: kong.String("client-secret-1"),
					},
				},
				{
					Oauth2Credential: kong.Oauth2Credential{
						ClientID:     kong.String("client-1"),
						ClientSecret: kong.String("client-secret-1"),
					},
				},
				{
					Oauth2Credential: kong.Oauth2Credential{
						ClientID:     kong.String("client-2"),
						ClientSecret: kong.String("client-secret-2"),
					},
				},
			},
			expectedConflicts: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conflictingOauth2Creds := getOAuth2CredentialsConflictingOnClientID(tc.oauth2Creds)
			assert.Len(t, conflictingOauth2Creds, tc.expectedConflicts)
		})
	}
}
