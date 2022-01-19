package admission

import (
	"context"
	"fmt"
	"testing"

	"github.com/kong/go-kong/kong"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

type fakePluginSvc struct {
	kong.AbstractPluginService

	err   error
	msg   string
	valid bool
}

func (f *fakePluginSvc) Validate(ctx context.Context, plugin *kong.Plugin) (bool, string, error) {
	return f.valid, f.msg, f.err
}

func TestKongHTTPValidator_ValidatePlugin(t *testing.T) {
	store, _ := store.NewFakeStore(store.FakeObjects{})
	type args struct {
		plugin configurationv1.KongPlugin
	}
	tests := []struct {
		name        string
		PluginSvc   kong.AbstractPluginService
		args        args
		wantOK      bool
		wantMessage string
		wantErr     bool
	}{
		{
			name:      "plugin is valid",
			PluginSvc: &fakePluginSvc{valid: true},
			args: args{
				plugin: configurationv1.KongPlugin{PluginName: "foo"},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:      "plugin is not valid",
			PluginSvc: &fakePluginSvc{valid: false, msg: "now where could my pipe be"},
			args: args{
				plugin: configurationv1.KongPlugin{PluginName: "foo"},
			},
			wantOK:      false,
			wantMessage: fmt.Sprintf(ErrTextPluginConfigViolatesSchema, "now where could my pipe be"),
			wantErr:     false,
		},
		{
			name:      "plugin lacks plugin name",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: configurationv1.KongPlugin{},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginNameEmpty,
			wantErr:     false,
		},
		{
			name:      "plugin has invalid configuration",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: configurationv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{{}`),
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigInvalid,
			wantErr:     true,
		},
		{
			name:      "plugin has both Config and ConfigFrom",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: configurationv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"key_names": "whatever"}`),
					},
					ConfigFrom: &configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "key-auth-config",
							Secret: "conf-secret",
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginUsesBothConfigTypes,
			wantErr:     false,
		},
		{
			name:      "plugin ConfigFrom references non-existent Secret",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: configurationv1.KongPlugin{
					PluginName: "key-auth",
					ConfigFrom: &configurationv1.ConfigSource{
						SecretValue: configurationv1.SecretValueFromSource{
							Key:    "key-auth-config",
							Secret: "conf-secret",
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginSecretConfigUnretrievable,
			wantErr:     true,
		},
		{
			name:      "failed to retrieve validation info",
			PluginSvc: &fakePluginSvc{valid: false, err: fmt.Errorf("everything broke")},
			args: args{
				plugin: configurationv1.KongPlugin{PluginName: "foo"},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigValidationFailed,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := KongHTTPValidator{
				SecretGetter:        store,
				PluginSvc:           tt.PluginSvc,
				ingressClassMatcher: fakeClassMatcher,
			}
			got, got1, err := validator.ValidatePlugin(context.Background(), tt.args.plugin)
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

func TestKongHTTPValidator_ValidateClusterPlugin(t *testing.T) {
	store, _ := store.NewFakeStore(store.FakeObjects{})
	type args struct {
		plugin configurationv1.KongClusterPlugin
	}
	tests := []struct {
		name        string
		PluginSvc   kong.AbstractPluginService
		args        args
		wantOK      bool
		wantMessage string
		wantErr     bool
	}{
		{
			name:      "plugin is valid",
			PluginSvc: &fakePluginSvc{valid: true},
			args: args{
				plugin: configurationv1.KongClusterPlugin{PluginName: "foo"},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:      "plugin is not valid",
			PluginSvc: &fakePluginSvc{valid: false, msg: "now where could my pipe be"},
			args: args{
				plugin: configurationv1.KongClusterPlugin{PluginName: "foo"},
			},
			wantOK:      false,
			wantMessage: fmt.Sprintf(ErrTextPluginConfigViolatesSchema, "now where could my pipe be"),
			wantErr:     false,
		},
		{
			name:      "plugin lacks plugin name",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: configurationv1.KongClusterPlugin{},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginNameEmpty,
			wantErr:     false,
		},
		{
			name:      "plugin has invalid configuration",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: configurationv1.KongClusterPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{{}`),
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigInvalid,
			wantErr:     true,
		},
		{
			name:      "plugin has both Config and ConfigFrom",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: configurationv1.KongClusterPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"key_names": "whatever"}`),
					},
					ConfigFrom: &configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "key-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginUsesBothConfigTypes,
			wantErr:     false,
		},
		{
			name:      "plugin ConfigFrom references non-existent Secret",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: configurationv1.KongClusterPlugin{
					PluginName: "key-auth",
					ConfigFrom: &configurationv1.NamespacedConfigSource{
						SecretValue: configurationv1.NamespacedSecretValueFromSource{
							Key:       "key-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginSecretConfigUnretrievable,
			wantErr:     true,
		},
		{
			name:      "failed to retrieve validation info",
			PluginSvc: &fakePluginSvc{valid: false, err: fmt.Errorf("everything broke")},
			args: args{
				plugin: configurationv1.KongClusterPlugin{PluginName: "foo"},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigValidationFailed,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := KongHTTPValidator{
				SecretGetter:        store,
				PluginSvc:           tt.PluginSvc,
				ingressClassMatcher: fakeClassMatcher,
			}
			got, got1, err := validator.ValidateClusterPlugin(context.Background(), tt.args.plugin)
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

func fakeClassMatcher(*metav1.ObjectMeta, annotations.ClassMatching) bool { return true }
