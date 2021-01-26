package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConsumer_SetCredential(t *testing.T) {
	username := "example"
	type args struct {
		credType   string
		consumer   *Consumer
		credConfig interface{}
	}
	tests := []struct {
		name    string
		args    args
		result  *Consumer
		wantErr bool
	}{
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
						Key: kong.String("foo"),
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "key-auth without key",
			args: args{
				credType:   "key-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username}},
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
						Key: kong.String("foo"),
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
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "basic-auth without username",
			args: args{
				credType:   "basic-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username}},
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
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "hmac-auth without username",
			args: args{
				credType:   "hmac-auth",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username}},
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
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "oauth2 without client_id",
			args: args{
				credType:   "oauth2",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username}},
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
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "jwt without key",
			args: args{
				credType:   "jwt",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username}},
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
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "acl without group",
			args: args{
				credType:   "acl",
				consumer:   &Consumer{Consumer: kong.Consumer{Username: &username}},
				credConfig: map[string]string{},
			},
			result:  &Consumer{Consumer: kong.Consumer{Username: &username}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.consumer.SetCredential(logrus.New(), tt.args.credType,
				tt.args.credConfig); (err != nil) != tt.wantErr {
				t.Errorf("processCredential() error = %v, wantErr %v",
					err, tt.wantErr)
			}
			assert.Equal(t, tt.result, tt.args.consumer)
		})
	}
}
