package admission

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
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

func TestKongHTTPValidator_ValidateConsumer(t *testing.T) {
	basicConsumer := func() configurationv1.KongConsumer {
		return configurationv1.KongConsumer{
			TypeMeta: metav1.TypeMeta{
				Kind:       "KongConsumer",
				APIVersion: configurationv1.GroupVersion.Version,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "consumer",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Username: "username",
		}
	}
	validSecret := func() *corev1.Secret {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"kongCredType": []byte("key-auth"),
				"key":          []byte("secret-key"),
			},
		}
	}

	tests := []struct {
		name                           string
		secrets                        []*corev1.Secret
		consumers                      []*configurationv1.KongConsumer
		modifyBasicConsumer            func(c *configurationv1.KongConsumer)
		consumerAlreadyExistsInGateway bool
		expectError                    bool
		expectOK                       bool
		expectedMessage                string
	}{
		{
			name: "consumer refers a non-existent secret",
			modifyBasicConsumer: func(c *configurationv1.KongConsumer) {
				c.Credentials = []string{"non-existing-secret"}
			},
			expectOK:        false,
			expectError:     true,
			expectedMessage: "could not retrieve secrets from the kubernetes API",
		},
		{
			name: "consumer refers a valid secret",
			modifyBasicConsumer: func(c *configurationv1.KongConsumer) {
				c.Credentials = []string{"secret"}
			},
			secrets:  []*corev1.Secret{validSecret()},
			expectOK: true,
		},
		{
			name: "consumer refers a secret with no kongCredType",
			modifyBasicConsumer: func(c *configurationv1.KongConsumer) {
				c.Credentials = []string{"secret"}
			},
			secrets: []*corev1.Secret{
				func() *corev1.Secret {
					s := validSecret()
					delete(s.Data, "kongCredType")
					return s
				}(),
			},
			expectError:     true,
			expectedMessage: "consumer credential failed validation",
		},
		{
			name: "consumer refers a secret that is already referred by another consumer",
			modifyBasicConsumer: func(c *configurationv1.KongConsumer) {
				c.Credentials = []string{"secret"}
			},
			secrets: []*corev1.Secret{validSecret()},
			consumers: []*configurationv1.KongConsumer{
				func() *configurationv1.KongConsumer {
					c := basicConsumer()
					c.Name = "consumer-2"
					c.Credentials = []string{"secret"}
					return &c
				}(),
			},
			expectError:     true,
			expectedMessage: "consumer credential violated unique key constraint",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			s, _ := store.NewFakeStore(store.FakeObjects{
				Secrets:       tt.secrets,
				KongConsumers: tt.consumers,
			})
			validator := KongHTTPValidator{
				// For the sake of tests we use the same store for both SecretGetter and Store as it does not matter here
				// if secrets are managed by us or not yet.
				Store:               s,
				SecretGetter:        s,
				ConsumerSvc:         &fakeConsumerService{consumerAlreadyExists: tt.consumerAlreadyExistsInGateway},
				ingressClassMatcher: fakeClassMatcher,
			}

			toValidate := basicConsumer()
			tt.modifyBasicConsumer(&toValidate)
			ok, message, err := validator.ValidateConsumer(context.Background(), toValidate)
			assert.Equal(t, tt.expectOK, ok)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedMessage, message)
		})
	}
}

func fakeClassMatcher(*metav1.ObjectMeta, string, annotations.ClassMatching) bool { return true }

type fakeConsumerService struct {
	kong.AbstractConsumerService
	consumerAlreadyExists bool
}

func (f *fakeConsumerService) Get(context.Context, *string) (*kong.Consumer, error) {
	if f.consumerAlreadyExists {
		return &kong.Consumer{}, nil
	}
	return nil, kong.NewAPIError(http.StatusNotFound, "consumer not found")
}
