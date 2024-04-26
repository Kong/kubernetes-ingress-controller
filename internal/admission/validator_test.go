package admission

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testk8sclient "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	managerscheme "github.com/kong/kubernetes-ingress-controller/v3/internal/manager/scheme"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

type fakePluginSvc struct {
	kong.AbstractPluginService

	err   error
	msg   string
	valid bool
}

func (f *fakePluginSvc) Validate(context.Context, *kong.Plugin) (bool, string, error) {
	return f.valid, f.msg, f.err
}

type fakeConsumersSvc struct {
	kong.AbstractConsumerService
	consumer *kong.Consumer
}

func (f fakeConsumersSvc) Get(context.Context, *string) (*kong.Consumer, error) {
	if f.consumer != nil {
		return f.consumer, nil
	}
	return nil, kong.NewAPIError(http.StatusNotFound, "no consumer found")
}

type fakeServicesProvider struct {
	pluginSvc        kong.AbstractPluginService
	consumerSvc      kong.AbstractConsumerService
	consumerGroupSvc kong.AbstractConsumerGroupService
	infoSvc          kong.AbstractInfoService
	routeSvc         kong.AbstractRouteService
	vaultSvc         kong.AbstractVaultService
}

func (f fakeServicesProvider) GetConsumersService() (kong.AbstractConsumerService, bool) {
	if f.consumerSvc != nil {
		return f.consumerSvc, true
	}
	return nil, false
}

func (f fakeServicesProvider) GetInfoService() (kong.AbstractInfoService, bool) {
	if f.infoSvc != nil {
		return f.infoSvc, true
	}
	return nil, false
}

func (f fakeServicesProvider) GetConsumerGroupsService() (kong.AbstractConsumerGroupService, bool) {
	if f.consumerGroupSvc != nil {
		return f.consumerGroupSvc, true
	}
	return nil, false
}

func (f fakeServicesProvider) GetPluginsService() (kong.AbstractPluginService, bool) {
	if f.pluginSvc != nil {
		return f.pluginSvc, true
	}
	return nil, false
}

func (f fakeServicesProvider) GetRoutesService() (kong.AbstractRouteService, bool) {
	if f.routeSvc != nil {
		return f.routeSvc, true
	}
	return nil, false
}

func (f fakeServicesProvider) GetVaultsService() (kong.AbstractVaultService, bool) {
	if f.vaultSvc != nil {
		return f.vaultSvc, true
	}
	return nil, false
}

func TestKongHTTPValidator_ValidatePlugin(t *testing.T) {
	store, _ := store.NewFakeStore(store.FakeObjects{
		Secrets: []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "",
					Name:      "conf-secret",
				},
				Data: map[string][]byte{
					"valid-conf":   []byte(`{"foo":"bar"}`),
					"invalid-conf": []byte(`{"foo":"baz}`),
				},
			},
		},
	})
	type args struct {
		plugin          kongv1.KongPlugin
		overrideSecrets []*corev1.Secret
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
				plugin: kongv1.KongPlugin{PluginName: "foo"},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:      "plugin is not valid",
			PluginSvc: &fakePluginSvc{valid: false, msg: "now where could my pipe be"},
			args: args{
				plugin: kongv1.KongPlugin{PluginName: "foo"},
			},
			wantOK:      false,
			wantMessage: fmt.Sprintf(ErrTextPluginConfigViolatesSchema, "now where could my pipe be"),
			wantErr:     false,
		},
		{
			name:      "plugin has invalid configuration",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: kongv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{{}`),
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigInvalid,
			wantErr:     false,
		},
		{
			name:      "plugin has valid configPatches",
			PluginSvc: &fakePluginSvc{valid: true},
			args: args{
				plugin: kongv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"k1":"v1"}`),
					},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/foo",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Secret: "conf-secret",
									Key:    "valid-conf",
								},
							},
						},
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:      "plugin has invalid configPatches",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: kongv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"k1":"v1"}`),
					},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/foo",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Secret: "conf-secret",
									Key:    "invalid-conf",
								},
							},
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigInvalid,
			wantErr:     false,
		},
		{
			name:      "plugin ConfigFrom references non-existent Secret",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: kongv1.KongPlugin{
					PluginName: "key-auth",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "key-auth-config",
							Secret: "conf-secret",
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginSecretConfigUnretrievable,
			wantErr:     false,
		},
		{
			name:      "failed to retrieve validation info",
			PluginSvc: &fakePluginSvc{valid: false, err: fmt.Errorf("everything broke")},
			args: args{
				plugin: kongv1.KongPlugin{PluginName: "foo"},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigValidationFailed,
			wantErr:     true,
		},
		{
			name:      "validate from override secret which generates valid configuration",
			PluginSvc: &fakePluginSvc{valid: true},
			args: args{
				plugin: kongv1.KongPlugin{
					PluginName: "key-auth",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Key:    "valid-conf",
							Secret: "another-conf-secret",
						},
					},
				},
				overrideSecrets: []*corev1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "",
							Name:      "another-conf-secret",
						},
						Data: map[string][]byte{
							"valid-conf":   []byte(`{"foo":"bar"}`),
							"invalid-conf": []byte(`{"foo":"baz}`),
						},
					},
				},
			},
			wantOK: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := KongHTTPValidator{
				SecretGetter: store,
				AdminAPIServicesProvider: fakeServicesProvider{
					pluginSvc: tt.PluginSvc,
				},
				ingressClassMatcher: fakeClassMatcher,
			}
			gotOK, gotMessage, err := validator.ValidatePlugin(context.Background(), tt.args.plugin, tt.args.overrideSecrets)
			assert.Equalf(t, tt.wantOK, gotOK,
				"KongHTTPValidator.ValidatePlugin() want OK: %v, got OK: %v",
				tt.wantOK, gotOK,
			)
			if tt.wantMessage != "" {
				assert.Containsf(t, gotMessage, tt.wantMessage,
					"KongHTTPValidator.ValidatePlugin() want message: %v, got message: %v",
					tt.wantMessage, gotMessage,
				)
			}
			assert.Equalf(t, tt.wantErr, err != nil,
				"KongHTTPValidator.ValidatePlugin() wantErr %v, got error %v",
				tt.wantErr, err,
			)
		})
	}
}

func TestKongHTTPValidator_ValidateClusterPlugin(t *testing.T) {
	store, _ := store.NewFakeStore(store.FakeObjects{
		Secrets: []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "conf-secret",
				},
				Data: map[string][]byte{
					"valid-conf":   []byte(`{"foo":"bar"}`),
					"invalid-conf": []byte(`{"foo":"baz}`),
				},
			},
		},
	})
	type args struct {
		plugin          kongv1.KongClusterPlugin
		overrideSecrets []*corev1.Secret
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
				plugin: kongv1.KongClusterPlugin{PluginName: "foo"},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:      "plugin is not valid",
			PluginSvc: &fakePluginSvc{valid: false, msg: "now where could my pipe be"},
			args: args{
				plugin: kongv1.KongClusterPlugin{PluginName: "foo"},
			},
			wantOK:      false,
			wantMessage: fmt.Sprintf(ErrTextPluginConfigViolatesSchema, "now where could my pipe be"),
			wantErr:     false,
		},
		{
			name:      "plugin has invalid configuration",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{{}`),
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigInvalid,
			wantErr:     false,
		},
		{
			name:      "plugin has valid configPatches",
			PluginSvc: &fakePluginSvc{valid: true},
			args: args{
				plugin: kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"k1":"v1"}`),
					},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/foo",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Secret:    "conf-secret",
									Key:       "valid-conf",
								},
							},
						},
					},
				},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:      "plugin has invalid configPatches",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"k1":"v1"}`),
					},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/foo",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: "default",
									Secret:    "conf-secret",
									Key:       "invalid-conf",
								},
							},
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigInvalid,
			wantErr:     false,
		},
		{
			name:      "plugin ConfigFrom references non-existent Secret",
			PluginSvc: &fakePluginSvc{},
			args: args{
				plugin: kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Key:       "key-auth-config",
							Secret:    "conf-secret",
							Namespace: "default",
						},
					},
				},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginSecretConfigUnretrievable,
			wantErr:     false,
		},
		{
			name:      "failed to retrieve validation info",
			PluginSvc: &fakePluginSvc{valid: false, err: fmt.Errorf("everything broke")},
			args: args{
				plugin: kongv1.KongClusterPlugin{PluginName: "foo"},
			},
			wantOK:      false,
			wantMessage: ErrTextPluginConfigValidationFailed,
			wantErr:     true,
		},
		{
			name:      "validate from override secret which generates valid configuration",
			PluginSvc: &fakePluginSvc{valid: true},
			args: args{
				plugin: kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Namespace: "default",
							Key:       "valid-conf",
							Secret:    "another-conf-secret",
						},
					},
				},
				overrideSecrets: []*corev1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "another-conf-secret",
						},
						Data: map[string][]byte{
							"valid-conf":   []byte(`{"foo":"bar"}`),
							"invalid-conf": []byte(`{"foo":"baz}`),
						},
					},
				},
			},
			wantOK: true,
		},
		{
			name:      "no gateway was available at the time of validation",
			PluginSvc: nil, // no plugin service is available as there's no gateways
			args: args{
				plugin: kongv1.KongClusterPlugin{PluginName: "foo"},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := KongHTTPValidator{
				SecretGetter: store,
				AdminAPIServicesProvider: fakeServicesProvider{
					pluginSvc: tt.PluginSvc,
				},
				ingressClassMatcher: fakeClassMatcher,
			}

			gotOK, gotMessage, err := validator.ValidateClusterPlugin(context.Background(), tt.args.plugin, tt.args.overrideSecrets)
			assert.Equalf(t, tt.wantOK, gotOK,
				"KongHTTPValidator.ValidateClusterPlugin() want OK: %v, got OK: %v",
				tt.wantOK, gotOK,
			)
			if tt.wantMessage != "" {
				assert.Containsf(t, gotMessage, tt.wantMessage,
					"KongHTTPValidator.ValidateClusterPlugin() want message: %v, got message: %v",
					tt.wantMessage, gotMessage,
				)
			}
			assert.Equalf(t, tt.wantErr, err != nil,
				"KongHTTPValidator.ValidateClusterPlugin() wantErr %v, got error %v",
				tt.wantErr, err,
			)
		})
	}
}

func TestKongHTTPValidator_ValidateConsumer(t *testing.T) {
	t.Run("passes with and without consumers service available", func(t *testing.T) {
		s, _ := store.NewFakeStore(store.FakeObjects{})
		validator := KongHTTPValidator{
			SecretGetter: s,
			AdminAPIServicesProvider: fakeServicesProvider{
				consumerSvc: fakeConsumersSvc{consumer: nil},
			},
			ingressClassMatcher: fakeClassMatcher,
		}

		valid, errText, err := validator.ValidateConsumer(context.Background(), kongv1.KongConsumer{
			Username: "username",
		})
		require.NoError(t, err)
		require.True(t, valid)
		require.Empty(t, errText)

		// make services unavailable
		validator.AdminAPIServicesProvider = fakeServicesProvider{}

		valid, errText, err = validator.ValidateConsumer(context.Background(), kongv1.KongConsumer{
			Username: "username",
		})
		require.NoError(t, err)
		require.True(t, valid)
		require.Empty(t, errText)
	})

	t.Run("fails when services available and consumer exists", func(t *testing.T) {
		s, _ := store.NewFakeStore(store.FakeObjects{})
		validator := KongHTTPValidator{
			SecretGetter: s,
			AdminAPIServicesProvider: fakeServicesProvider{
				consumerSvc: fakeConsumersSvc{consumer: &kong.Consumer{Username: lo.ToPtr("username")}},
			},
			ingressClassMatcher: fakeClassMatcher,
		}

		valid, errText, err := validator.ValidateConsumer(context.Background(), kongv1.KongConsumer{
			Username: "username",
		})
		require.NoError(t, err)
		require.False(t, valid)
		require.Equal(t, ErrTextConsumerExists, errText)
	})
}

type fakeConsumerGroupSvc struct {
	kong.AbstractConsumerGroupService
	err error
}

func (f fakeConsumerGroupSvc) List(context.Context, *kong.ListOpt) ([]*kong.ConsumerGroup, *kong.ListOpt, error) {
	if f.err != nil {
		return []*kong.ConsumerGroup{}, &kong.ListOpt{}, f.err
	}
	return []*kong.ConsumerGroup{}, &kong.ListOpt{}, nil
}

type fakeInfoSvc struct {
	kong.AbstractInfoService
	version string
}

func (f fakeInfoSvc) Get(context.Context) (*kong.Info, error) {
	if f.version != "" {
		return &kong.Info{Version: f.version}, nil
	}
	return nil, kong.NewAPIError(http.StatusInternalServerError, "bogus fake info")
}

func TestKongHTTPValidator_ValidateConsumerGroup(t *testing.T) {
	store, err := store.NewFakeStore(store.FakeObjects{})
	require.NoError(t, err)
	type args struct {
		cg kongv1beta1.KongConsumerGroup
	}
	tests := []struct {
		name             string
		ConsumerGroupSvc kong.AbstractConsumerGroupService
		InfoSvc          kong.AbstractInfoService
		args             args
		wantOK           bool
		wantMessage      string
		wantErr          bool
	}{
		{
			name:             "Enterprise version",
			ConsumerGroupSvc: &fakeConsumerGroupSvc{err: nil},
			InfoSvc:          &fakeInfoSvc{version: "3.4.1.0"},
			args: args{
				cg: kongv1beta1.KongConsumerGroup{},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:             "OSS version",
			ConsumerGroupSvc: &fakeConsumerGroupSvc{err: nil},
			InfoSvc:          &fakeInfoSvc{version: "3.4.1"},
			args: args{
				cg: kongv1beta1.KongConsumerGroup{},
			},
			wantOK:      false,
			wantMessage: ErrTextConsumerGroupUnsupported,
			wantErr:     false,
		},
		{
			name:             "Enterprise version, unlicensed",
			ConsumerGroupSvc: &fakeConsumerGroupSvc{err: kong.NewAPIError(http.StatusForbidden, "no license")},
			InfoSvc:          &fakeInfoSvc{version: "3.4.1.0"},
			args: args{
				cg: kongv1beta1.KongConsumerGroup{},
			},
			wantOK:      false,
			wantMessage: ErrTextConsumerGroupUnlicensed,
			wantErr:     false,
		},
		{
			name:             "Enterprise version, API somehow missing",
			ConsumerGroupSvc: &fakeConsumerGroupSvc{err: kong.NewAPIError(http.StatusNotFound, "well, this is awkward")},
			InfoSvc:          &fakeInfoSvc{version: "3.4.1.0"},
			args: args{
				cg: kongv1beta1.KongConsumerGroup{},
			},
			wantOK:      false,
			wantMessage: ErrTextConsumerGroupUnsupported,
			wantErr:     false,
		},
		{
			name:             "invalid semver with consumer groups support",
			ConsumerGroupSvc: &fakeConsumerGroupSvc{err: nil},
			InfoSvc:          &fakeInfoSvc{version: "a.4.0.0"},
			args: args{
				cg: kongv1beta1.KongConsumerGroup{},
			},
			wantOK:      true,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:             "invalid semver with no consumer groups support",
			ConsumerGroupSvc: &fakeConsumerGroupSvc{err: kong.NewAPIError(http.StatusNotFound, "ConsumerGroups API not found")},
			InfoSvc:          &fakeInfoSvc{version: "a.4.0.0"},
			args: args{
				cg: kongv1beta1.KongConsumerGroup{},
			},
			wantOK:      false,
			wantMessage: ErrTextConsumerGroupUnsupported,
			wantErr:     false,
		},
		{
			name:             "Enterprise version, API returning unexpected error",
			ConsumerGroupSvc: &fakeConsumerGroupSvc{err: kong.NewAPIError(http.StatusTeapot, "I'm a teapot")},
			InfoSvc:          &fakeInfoSvc{version: "3.4.1.0"},
			args: args{
				cg: kongv1beta1.KongConsumerGroup{},
			},
			wantOK:      false,
			wantMessage: fmt.Sprintf("%s: %s", ErrTextConsumerGroupUnexpected, `HTTP status 418 (message: "I'm a teapot")`),
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := KongHTTPValidator{
				SecretGetter: store,
				AdminAPIServicesProvider: fakeServicesProvider{
					infoSvc:          tt.InfoSvc,
					consumerGroupSvc: tt.ConsumerGroupSvc,
				},
				ingressClassMatcher: fakeClassMatcher,
				Logger:              zapr.NewLogger(zap.NewNop()),
			}
			got, gotMsg, err := validator.ValidateConsumerGroup(context.Background(), tt.args.cg)
			if (err != nil) != tt.wantErr {
				t.Errorf("KongHTTPValidator.ValidateConsumerGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantOK {
				t.Errorf("KongHTTPValidator.ValidateConsumerGroups() got = %v, want %v", got, tt.wantOK)
			}
			if gotMsg != tt.wantMessage {
				t.Errorf("KongHTTPValidator.ValidateConsumerGroups() gotMsg = %v, want %v", gotMsg, tt.wantMessage)
			}
		})
	}
}

func fakeClassMatcher(*metav1.ObjectMeta, string, annotations.ClassMatching) bool { return true }

func TestKongHTTPValidator_ValidateCredential(t *testing.T) {
	testCases := []struct {
		name            string
		consumers       []kongv1.KongConsumer
		secret          corev1.Secret
		wantOK          bool
		wantMessage     string
		wantErrContains string
	}{
		{
			name: "labeled valid key-auth credential with no consumers gets accepted",
			secret: corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"konghq.com/credential": "key-auth",
					},
				},
				Data: map[string][]byte{
					"key": []byte("my-key"),
				},
			},
			wantOK: true,
		},
		{
			name: "labeled invalid key-auth credential with no consumers gets rejected",
			secret: corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"konghq.com/credential": "key-auth",
					},
				},
				Data: map[string][]byte{
					// missing key
				},
			},
			wantOK:      false,
			wantMessage: fmt.Sprintf("%s: %s", ErrTextConsumerCredentialValidationFailed, "missing required field(s): key"),
		},
		{
			name: "valid key-auth credential with no consumers gets accepted",
			secret: corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "credential-0",
					Labels: map[string]string{
						"konghq.com/credential": "key-auth",
					},
				},
				Data: map[string][]byte{
					"key": []byte("my-key"),
				},
			},
			wantOK: true,
		},
		{
			name: "valid key-auth credential using only konghq.com/credential with a consumer gets accepted",
			consumers: []kongv1.KongConsumer{
				{
					Username: "username",
					Credentials: []string{
						"username-key-auth-1",
					},
				},
			},
			secret: corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "username-key-auth-1",
					Labels: map[string]string{
						"konghq.com/credential": "key-auth",
					},
				},
				Data: map[string][]byte{
					"key": []byte("my-key"),
				},
			},
			wantOK: true,
		},
		{
			name: "invalid key-auth credential with no consumers gets rejected",
			secret: corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "credential-0",
					Labels: map[string]string{
						"konghq.com/credential": "key-auth",
					},
				},
				Data: map[string][]byte{
					// missing key
				},
			},
			wantOK:      false,
			wantMessage: fmt.Sprintf("%s: %s", ErrTextConsumerCredentialValidationFailed, "missing required field(s): key"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, testk8sclient.AddToScheme(scheme))
			require.NoError(t, kongv1.AddToScheme(scheme))
			b := fake.NewClientBuilder().WithScheme(scheme)

			validator := KongHTTPValidator{
				ManagerClient: b.Build(),
				ConsumerGetter: fakeConsumerGetter{
					consumers: tc.consumers,
				},
				AdminAPIServicesProvider: fakeServicesProvider{},
				ingressClassMatcher:      fakeClassMatcher,
				Logger:                   logr.Discard(),
			}

			ok, msg := validator.ValidateCredential(context.Background(), tc.secret)
			assert.Equal(t, tc.wantOK, ok)
			assert.Equal(t, tc.wantMessage, msg)
		})
	}
}

type fakeConsumerGetter struct {
	consumers []kongv1.KongConsumer
}

func (f fakeConsumerGetter) ListAllConsumers(context.Context) ([]kongv1.KongConsumer, error) {
	return f.consumers, nil
}

func TestValidator_ValidateIngress(t *testing.T) {
	const testSvcFacadeName = "svc-facade"
	s := lo.Must(managerscheme.Get())
	b := fake.NewClientBuilder().WithScheme(s)

	testCases := []struct {
		name                          string
		storerObjects                 store.FakeObjects
		kongRouteValidationShouldFail bool
		translatorFeatures            translator.FeatureFlags
		ingress                       *netv1.Ingress
		wantOK                        bool
		wantMessage                   string
	}{
		{
			name: "not matching ingress class is always ok",
			ingress: builder.NewIngress("ingress", "not-kong").
				WithNamespace("default").
				WithRules(
					newHTTPIngressRule(netv1.IngressBackend{
						Service: &netv1.IngressServiceBackend{
							Name: "svc",
							Port: netv1.ServiceBackendPort{
								Number: 8080,
							},
						},
					}),
				).
				Build(),
			kongRouteValidationShouldFail: true, // Despite the route validation failing, the ingress class is not kong, so it's ok.
			storerObjects:                 store.FakeObjects{},
			wantOK:                        true,
		},
		{
			name: "valid with Service backend",
			ingress: builder.NewIngress("ingress", "not-kong").
				WithNamespace("default").
				WithRules(
					newHTTPIngressRule(netv1.IngressBackend{
						Service: &netv1.IngressServiceBackend{
							Name: "svc",
							Port: netv1.ServiceBackendPort{
								Number: 8080,
							},
						},
					}),
				).
				Build(),
			storerObjects: store.FakeObjects{},
			wantOK:        true,
		},
		{
			name: "invalid with Service backend",
			ingress: builder.NewIngress("ingress", "kong").
				WithNamespace("default").
				WithRules(
					newHTTPIngressRule(netv1.IngressBackend{
						Service: &netv1.IngressServiceBackend{
							Name: "svc",
							Port: netv1.ServiceBackendPort{
								Number: 8080,
							},
						},
					}),
				).
				Build(),
			kongRouteValidationShouldFail: true,
			storerObjects:                 store.FakeObjects{},
			wantOK:                        false,
			wantMessage:                   "Ingress failed schema validation: something is wrong with the route",
		},
		{
			name: "valid with KongServiceFacade backend",
			ingress: builder.NewIngress("ingress", "not-kong").
				WithNamespace("default").
				WithRules(
					newHTTPIngressRule(netv1.IngressBackend{
						Resource: &corev1.TypedLocalObjectReference{
							APIGroup: lo.ToPtr(incubatorv1alpha1.SchemeGroupVersion.Group),
							Kind:     incubatorv1alpha1.KongServiceFacadeKind,
							Name:     testSvcFacadeName,
						},
					}),
				).
				Build(),
			translatorFeatures: translator.FeatureFlags{
				KongServiceFacade: true,
			},
			storerObjects: store.FakeObjects{
				KongServiceFacades: []*incubatorv1alpha1.KongServiceFacade{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      testSvcFacadeName,
							Namespace: "default",
						},
						Spec: incubatorv1alpha1.KongServiceFacadeSpec{
							Backend: incubatorv1alpha1.KongServiceFacadeBackend{
								Name: "svc",
								Port: 8080,
							},
						},
					},
				},
			},
			wantOK: true,
		},
		{
			name: "invalid with KongServiceFacade backend",
			translatorFeatures: translator.FeatureFlags{
				KongServiceFacade: true,
			},
			ingress: builder.NewIngress("ingress", "kong").
				WithNamespace("default").
				WithRules(
					newHTTPIngressRule(netv1.IngressBackend{
						Resource: &corev1.TypedLocalObjectReference{
							APIGroup: lo.ToPtr(incubatorv1alpha1.SchemeGroupVersion.Group),
							Kind:     incubatorv1alpha1.KongServiceFacadeKind,
							Name:     testSvcFacadeName,
						},
					}),
				).
				Build(),
			storerObjects: store.FakeObjects{}, // No KongServiceFacade will be found resulting in an error.
			wantOK:        false,
			wantMessage:   `Ingress failed schema validation: failed to get backend for ingress path "/": failed to get KongServiceFacade "svc-facade": KongServiceFacade default/svc-facade not found`,
		},
		{
			name: "invalid with KongServiceFacade backend with feature flag off is ok",
			translatorFeatures: translator.FeatureFlags{
				KongServiceFacade: false,
			},
			ingress: builder.NewIngress("ingress", "not-kong").
				WithNamespace("default").
				WithRules(
					newHTTPIngressRule(netv1.IngressBackend{
						Resource: &corev1.TypedLocalObjectReference{
							APIGroup: lo.ToPtr(incubatorv1alpha1.SchemeGroupVersion.Group),
							Kind:     incubatorv1alpha1.KongServiceFacadeKind,
							Name:     testSvcFacadeName,
						},
					}),
				).
				Build(),
			storerObjects: store.FakeObjects{}, // No KongServiceFacade found would result in an error, but the feature flag is off.
			wantOK:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storer := lo.Must(store.NewFakeStore(tc.storerObjects))
			validator := KongHTTPValidator{
				ManagerClient: b.Build(),
				Storer:        storer,
				AdminAPIServicesProvider: fakeServicesProvider{
					routeSvc: &fakeRouteSvc{
						shouldFail: tc.kongRouteValidationShouldFail,
					},
				},
				TranslatorFeatures: tc.translatorFeatures,
				ingressClassMatcher: func(*metav1.ObjectMeta, string, annotations.ClassMatching) bool {
					return false // Always return false, we'll use Spec.IngressClassName matcher.
				},
				ingressV1ClassMatcher: func(ingress *netv1.Ingress, _ annotations.ClassMatching) bool {
					return *ingress.Spec.IngressClassName == annotations.DefaultIngressClass
				},
				Logger: logr.Discard(),
			}
			ok, msg, err := validator.ValidateIngress(context.Background(), *tc.ingress)
			require.NoError(t, err)
			assert.Equal(t, tc.wantOK, ok)
			assert.Equal(t, tc.wantMessage, msg)
		})
	}
}

func newHTTPIngressRule(backend netv1.IngressBackend) netv1.IngressRule {
	return netv1.IngressRule{
		IngressRuleValue: netv1.IngressRuleValue{
			HTTP: &netv1.HTTPIngressRuleValue{
				Paths: []netv1.HTTPIngressPath{
					{
						Path:     "/",
						PathType: lo.ToPtr(netv1.PathTypeImplementationSpecific),
						Backend:  backend,
					},
				},
			},
		},
	}
}

type fakeRouteSvc struct {
	kong.AbstractRouteService
	shouldFail bool
}

func (f *fakeRouteSvc) Validate(context.Context, *kong.Route) (bool, string, error) {
	if f.shouldFail {
		return false, "something is wrong with the route", nil
	}
	return true, "", nil
}

func TestValidator_ValidateVault(t *testing.T) {
	testCases := []struct {
		name            string
		kongVault       kongv1alpha1.KongVault
		storerObjects   store.FakeObjects
		validateSvcFail bool
		expectedOK      bool
		expectedMessage string
		expectedError   string
	}{
		{
			name: "valid vault",
			kongVault: kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vault-1",
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend: "env",
					Prefix:  "env-1",
				},
			},
			expectedOK: true,
		},
		{
			name: "vault with invalid(malformed) configuration",
			kongVault: kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vault-1",
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend: "env",
					Prefix:  "env-1",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{{}`),
					},
				},
			},
			expectedOK:      false,
			expectedMessage: "failed to unmarshal vault configuration",
		},
		{
			name: "vault with duplicate prefix",
			kongVault: kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vault-1",
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend: "env",
					Prefix:  "env-dupe",
				},
			},
			storerObjects: store.FakeObjects{
				KongVaults: []*kongv1alpha1.KongVault{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "vault-0",
							Annotations: map[string]string{
								annotations.IngressClassKey: annotations.DefaultIngressClass,
							},
						},
						Spec: kongv1alpha1.KongVaultSpec{
							Backend: "env",
							Prefix:  "env-dupe",
						},
					},
				},
			},
			validateSvcFail: false,
			expectedOK:      false,
			expectedMessage: "spec.prefix \"env-dupe\" is duplicate with existing KongVault",
		},
		{
			name: "vault with failure in validating service",
			kongVault: kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vault-1",
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend: "env",
					Prefix:  "env-1",
				},
			},
			validateSvcFail: true,
			expectedOK:      false,
			expectedMessage: "something is wrong with the vault",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			storer := lo.Must(store.NewFakeStore(tc.storerObjects))
			validator := KongHTTPValidator{
				AdminAPIServicesProvider: fakeServicesProvider{
					vaultSvc: &fakeVaultSvc{
						shouldFail: tc.validateSvcFail,
					},
				},
				Storer:              storer,
				ingressClassMatcher: fakeClassMatcher,
				Logger:              logr.Discard(),
			}
			ok, msg, err := validator.ValidateVault(context.Background(), tc.kongVault)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOK, ok)
			assert.Contains(t, msg, tc.expectedMessage)
		})
	}
}

type fakeVaultSvc struct {
	kong.AbstractVaultService
	shouldFail bool
}

func (s fakeVaultSvc) Validate(context.Context, *kong.Vault) (bool, string, error) {
	if s.shouldFail {
		return false, "something is wrong with the vault", nil
	}
	return true, "", nil
}
