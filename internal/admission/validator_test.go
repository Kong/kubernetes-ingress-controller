package admission

import (
	"testing"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestKongHTTPValidator_ValidateCredential(t *testing.T) {
	type args struct {
		secret corev1.Secret
	}
	tests := []struct {
		name        string
		args        args
		wantOK      bool
		wantMessage string
		wantErr     bool
	}{
		{
			name: "valid key-auth credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"key":          []byte("foo"),
						"kongCredType": []byte("key-auth"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "valid jwt credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"algorithm":      []byte("foo-algorithm"),
						"key":            []byte("foo-key"),
						"secret":         []byte("foo-secret"),
						"rsa_public_key": []byte("my-key"),
						"kongCredType":   []byte("key-auth"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "valid jwt credential without rsa_public_key",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"algorithm":    []byte("foo-algorithm"),
						"key":          []byte("foo-key"),
						"secret":       []byte("foo-secret"),
						"kongCredType": []byte("key-auth"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "valid keyauth_credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"key":          []byte("foo"),
						"kongCredType": []byte("keyauth_credential"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "invalid key-auth credential",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"key-wrong":    []byte("foo"),
						"kongCredType": []byte("key-auth"),
					},
				},
			},
			wantOK:      false,
			wantMessage: "missing required field(s): key",
			wantErr:     false,
		},
		{
			name: "invalid credential type",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"kongCredType": []byte("foo"),
					},
				},
			},
			wantOK:      false,
			wantMessage: "invalid credential type: foo",
			wantErr:     false,
		},
		{
			name: "non-kong secrets are passed",
			args: args{
				secret: corev1.Secret{
					Data: map[string][]byte{
						"key": []byte("foo"),
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := KongHTTPValidator{}
			got, got1, err := validator.ValidateCredential(tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("KongHTTPValidator.ValidateCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantOK {
				t.Errorf("KongHTTPValidator.ValidateCredential() got = %v, want %v", got, tt.wantOK)
			}
			if got1 != tt.wantMessage {
				t.Errorf("KongHTTPValidator.ValidateCredential() got1 = %v, want %v", got1, tt.wantMessage)
			}
		})
	}
}

func TestKongHTTPValidator_ValidatePlugin(t *testing.T) {
	store, _ := store.NewFakeStore(store.FakeObjects{})
	type args struct {
		plugin configurationv1.KongPlugin
	}
	tests := []struct {
		name        string
		args        args
		wantOK      bool
		wantMessage string
		wantErr     bool
	}{
		{
			name: "plugin lacks plugin name",
			args: args{
				plugin: configurationv1.KongPlugin{},
			},
			wantOK:      false,
			wantMessage: "plugin name cannot be empty",
			wantErr:     false,
		},
		{
			name: "plugin has invalid configuration",
			args: args{
				plugin: configurationv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{{}`),
					},
				},
			},
			wantOK:      false,
			wantMessage: "could not unmarshal plugin configuration",
			wantErr:     true,
		},
		{
			name: "plugin has both Config and ConfigFrom",
			args: args{
				plugin: configurationv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"key_names": "whatever"}`),
					},
					ConfigFrom: configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "key-auth-config",
							Secret: "conf-secret",
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: "plugin cannot use both Config and ConfigFrom",
			wantErr:     false,
		},
		{
			name: "plugin ConfigFrom references non-existent Secret",
			args: args{
				plugin: configurationv1.KongPlugin{
					PluginName: "key-auth",
					ConfigFrom: configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "key-auth-config",
							Secret: "conf-secret",
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: "could not load secret plugin configuration",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := KongHTTPValidator{
				Store: store,
			}
			got, got1, err := validator.ValidatePlugin(tt.args.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("KongHTTPValidator.ValidatePlugin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantOK {
				t.Errorf("KongHTTPValidator.ValidatePlugin() got = %v, want %v", got, tt.wantOK)
			}
			if got1 != tt.wantMessage {
				t.Errorf("KongHTTPValidator.ValidatePlugin() got1 = %v, want %v", got1, tt.wantMessage)
			}
		})
	}
}
