package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"

	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func int64Ptr(i int64) *int64 {
	return &i
}

func TestConsumer_SanitizedCopy(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   Consumer
		want Consumer
	}{
		{
			name: "sanitizes all credentials and copies all other fields",
			in: Consumer{
				Consumer: kong.Consumer{
					ID:        kong.String("1"),
					CustomID:  kong.String("2"),
					Username:  kong.String("3"),
					CreatedAt: int64Ptr(4),
					Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
				},
				Plugins:    []kong.Plugin{{ID: kong.String("1")}},
				KeyAuths:   []*KeyAuth{{kong.KeyAuth{ID: kong.String("1"), Key: kong.String("secret")}}},
				HMACAuths:  []*HMACAuth{{kong.HMACAuth{ID: kong.String("1"), Secret: kong.String("secret")}}},
				JWTAuths:   []*JWTAuth{{kong.JWTAuth{ID: kong.String("1"), Secret: kong.String("secret")}}},
				BasicAuths: []*BasicAuth{{kong.BasicAuth{ID: kong.String("1"), Password: kong.String("secret")}}},
				ACLGroups:  []*ACLGroup{{kong.ACLGroup{ID: kong.String("1")}}},
				Oauth2Creds: []*Oauth2Credential{
					{kong.Oauth2Credential{ID: kong.String("1"), ClientSecret: kong.String("secret")}},
				},
				MTLSAuths:       []*MTLSAuth{{kong.MTLSAuth{ID: kong.String("1"), SubjectName: kong.String("foo@example.com")}}},
				K8sKongConsumer: kongv1.KongConsumer{Username: "foo"},
			},
			want: Consumer{
				Consumer: kong.Consumer{
					ID:        kong.String("1"),
					CustomID:  kong.String("2"),
					Username:  kong.String("3"),
					CreatedAt: int64Ptr(4),
					Tags:      []*string{kong.String("5.1"), kong.String("5.2")},
				},
				Plugins:    []kong.Plugin{{ID: kong.String("1")}},
				KeyAuths:   []*KeyAuth{{kong.KeyAuth{ID: kong.String("1"), Key: kong.String("{vault://52fdfc07-2182-454f-963f-5f0f9a621d72}")}}},
				HMACAuths:  []*HMACAuth{{kong.HMACAuth{ID: kong.String("1"), Secret: redactedString}}},
				JWTAuths:   []*JWTAuth{{kong.JWTAuth{ID: kong.String("1"), Secret: redactedString}}},
				BasicAuths: []*BasicAuth{{kong.BasicAuth{ID: kong.String("1"), Password: redactedString}}},
				ACLGroups:  []*ACLGroup{{kong.ACLGroup{ID: kong.String("1")}}},
				Oauth2Creds: []*Oauth2Credential{
					{kong.Oauth2Credential{ID: kong.String("1"), ClientSecret: redactedString}},
				},
				MTLSAuths:       []*MTLSAuth{{kong.MTLSAuth{ID: kong.String("1"), SubjectName: kong.String("foo@example.com")}}},
				K8sKongConsumer: kongv1.KongConsumer{Username: "foo"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := *tt.in.SanitizedCopy(mocks.StaticUUIDGenerator{UUID: "52fdfc07-2182-454f-963f-5f0f9a621d72"})
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConsumer_SetCredential(t *testing.T) {
	username := "example"
	type args struct {
		credType   string
		consumer   *Consumer
		credConfig interface{}
	}
	type Case struct {
		name    string
		args    args
		result  *Consumer
		wantErr bool
	}

	tests := []Case{
		{
			name: "invalid cred type errors",
			args: args{
				credType:   "invalid-type",
				consumer:   &Consumer{},
				credConfig: nil,
			},
			result:  &Consumer{},
			wantErr: true,
		},
		{
			name: "key-auth",
			args: args{
				credType:   "key-auth",
				consumer:   &Consumer{},
				credConfig: map[string]string{"key": "foo"},
			},
			result: &Consumer{
				KeyAuths: []*KeyAuth{
					{kong.KeyAuth{
						Key:  kong.String("foo"),
						Tags: []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "key-auth without key",
			args: args{
				credType:   "key-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "key-auth with invalid key type",
			args: args{
				credType:   "key-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]interface{}{"key": true},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "keyauth_credential",
			args: args{
				credType:   "keyauth_credential",
				consumer:   &Consumer{},
				credConfig: map[string]string{"key": "foo"},
			},
			result: &Consumer{
				KeyAuths: []*KeyAuth{
					{kong.KeyAuth{
						Key:  kong.String("foo"),
						Tags: []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "basic-auth",
			args: args{
				credType: "basic-auth",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"username": "foo",
					"password": "bar",
				},
			},
			result: &Consumer{
				BasicAuths: []*BasicAuth{
					{kong.BasicAuth{
						Username: kong.String("foo"),
						Password: kong.String("bar"),
						Tags:     []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "basic-auth without username",
			args: args{
				credType:   "basic-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "basic-auth with invalid username type",
			args: args{
				credType:   "basic-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]interface{}{"username": true},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "basicauth_credential",
			args: args{
				credType: "basicauth_credential",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"username": "foo",
					"password": "bar",
				},
			},
			result: &Consumer{
				BasicAuths: []*BasicAuth{
					{kong.BasicAuth{
						Username: kong.String("foo"),
						Password: kong.String("bar"),
						Tags:     []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "hmac-auth",
			args: args{
				credType: "hmac-auth",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"username": "foo",
					"secret":   "bar",
				},
			},
			result: &Consumer{
				HMACAuths: []*HMACAuth{
					{kong.HMACAuth{
						Username: kong.String("foo"),
						Secret:   kong.String("bar"),
						Tags:     []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "hmac-auth without username",
			args: args{
				credType:   "hmac-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "hmac-auth with invalid username type",
			args: args{
				credType:   "hmac-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]interface{}{"username": true},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "hmacauth_credential",
			args: args{
				credType: "hmacauth_credential",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"username": "foo",
					"secret":   "bar",
				},
			},
			result: &Consumer{
				HMACAuths: []*HMACAuth{
					{kong.HMACAuth{
						Username: kong.String("foo"),
						Secret:   kong.String("bar"),
						Tags:     []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "oauth2",
			args: args{
				credType: "oauth2",
				consumer: &Consumer{},
				credConfig: map[string]interface{}{
					"name":          "foo",
					"client_id":     "bar",
					"client_secret": "baz",
					"redirect_uris": []string{"example.com"},
				},
			},
			result: &Consumer{
				Oauth2Creds: []*Oauth2Credential{
					{kong.Oauth2Credential{
						Name:         kong.String("foo"),
						ClientID:     kong.String("bar"),
						ClientSecret: kong.String("baz"),
						RedirectURIs: kong.StringSlice("example.com"),
						Tags:         []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "oauth2 without name",
			args: args{
				credType: "oauth2",
				consumer: &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]interface{}{
					"client_id": "bar",
				},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "oauth2 without client_id",
			args: args{
				credType: "oauth2",
				consumer: &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]interface{}{
					"name": "bar",
				},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "oauth2 with invalid client_id type",
			args: args{
				credType:   "oauth2",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]interface{}{"client_id": true},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "jwt",
			args: args{
				credType: "jwt",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"key":            "foo",
					"rsa_public_key": "bar",
					"secret":         "baz",
				},
			},
			result: &Consumer{
				JWTAuths: []*JWTAuth{
					{kong.JWTAuth{
						Key:          kong.String("foo"),
						RSAPublicKey: kong.String("bar"),
						Secret:       kong.String("baz"),
						// set by default
						Algorithm: kong.String("HS256"),
						Tags:      []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "jwt without key",
			args: args{
				credType:   "jwt",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "jwt with invald key type",
			args: args{
				credType:   "jwt",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]interface{}{"key": true},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "jwt_secret",
			args: args{
				credType: "jwt_secret",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"key":            "foo",
					"rsa_public_key": "bar",
					"secret":         "baz",
				},
			},
			result: &Consumer{
				JWTAuths: []*JWTAuth{
					{kong.JWTAuth{
						Key:          kong.String("foo"),
						RSAPublicKey: kong.String("bar"),
						Secret:       kong.String("baz"),
						// set by default
						Algorithm: kong.String("HS256"),
						Tags:      []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "acl",
			args: args{
				credType:   "acl",
				consumer:   &Consumer{},
				credConfig: map[string]string{"group": "group-foo"},
			},
			result: &Consumer{
				ACLGroups: []*ACLGroup{
					{kong.ACLGroup{
						Group: kong.String("group-foo"),
						Tags:  []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "acl without group",
			args: args{
				credType:   "acl",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "acl with invalid group type",
			args: args{
				credType:   "acl",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]interface{}{"group": true},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.consumer.SetCredential(tt.args.credType, tt.args.credConfig, []*string{})
			if (err != nil) != tt.wantErr {
				t.Errorf("processCredential() error = %v, wantErr %v",
					err, tt.wantErr)
			}
			assert.Equal(t, tt.args.consumer, tt.result)
		})
	}

	mtlsSupportedTests := []Case{
		{
			name: "mtls-auth",
			args: args{
				credType:   "mtls-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]string{"subject_name": "foo@example.com"},
			},
			result: &Consumer{
				Consumer: kong.Consumer{Username: &username, Tags: []*string{}},
				MTLSAuths: []*MTLSAuth{
					{kong.MTLSAuth{
						SubjectName: kong.String("foo@example.com"),
						Tags:        []*string{},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "mtls-auth without subject_name",
			args: args{
				credType:   "mtls-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
		{
			name: "mtls-auth with invalid subject_name type",
			args: args{
				credType:   "mtls-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
				credConfig: map[string]interface{}{"subject_name": true},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username, Tags: []*string{}}},
			wantErr: true,
		},
	}

	for _, tt := range mtlsSupportedTests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.consumer.SetCredential(tt.args.credType, tt.args.credConfig, []*string{})
			if (err != nil) != tt.wantErr {
				t.Errorf("processCredential() error = %v, wantErr %v",
					err, tt.wantErr)
			}
			assert.Equal(t, tt.args.consumer, tt.result)
		})
	}
}
