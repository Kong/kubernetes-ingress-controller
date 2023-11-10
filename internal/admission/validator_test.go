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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testk8sclient "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
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

func TestKongHTTPValidator_ValidatePlugin(t *testing.T) {
	store, _ := store.NewFakeStore(store.FakeObjects{})
	type args struct {
		plugin kongv1.KongPlugin
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
			wantErr:     true,
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
			wantErr:     true,
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
		plugin kongv1.KongClusterPlugin
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
			wantErr:     true,
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
			wantErr:     true,
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
				Data: map[string][]byte{
					"kongCredType": []byte("key-auth"),
					"key":          []byte("my-key"),
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
				Data: map[string][]byte{
					"kongCredType": []byte("key-auth"),
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
