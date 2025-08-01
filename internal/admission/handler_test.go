package admission

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/require"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	configurationv1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
)

var (
	secretTypeMeta = metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Secret",
	}
	kongPluginTypeMeta = metav1.TypeMeta{
		APIVersion: configurationv1.GroupVersion.String(),
		Kind:       "KongPlugin",
	}
	kongClusterPluginTypeMeta = metav1.TypeMeta{
		APIVersion: configurationv1.GroupVersion.String(),
		Kind:       "KongClusterPlugin",
	}
)

func TestHandleKongIngress(t *testing.T) {
	tests := []struct {
		name         string
		resource     configurationv1.KongIngress
		wantWarnings []string
	}{
		{
			name: "has proxy",
			resource: configurationv1.KongIngress{
				Proxy: &configurationv1.KongIngressService{},
			},
			wantWarnings: []string{proxyWarning},
		},
		{
			name: "has route",
			resource: configurationv1.KongIngress{
				Route: &configurationv1.KongIngressRoute{},
			},
			wantWarnings: []string{routeWarning},
		},
		{
			name: "has upstream",
			resource: configurationv1.KongIngress{
				Upstream: &configurationv1.KongIngressUpstream{},
			},
			wantWarnings: []string{upstreamWarning},
		},
		{
			name: "has everything",
			resource: configurationv1.KongIngress{
				Proxy:    &configurationv1.KongIngressService{},
				Route:    &configurationv1.KongIngressRoute{},
				Upstream: &configurationv1.KongIngressUpstream{},
			},
			wantWarnings: []string{proxyWarning, routeWarning, upstreamWarning},
		},
	}
	for _, tt := range tests {
		resource := tt.resource
		t.Run(tt.name, func(t *testing.T) {
			validator := KongHTTPValidator{}
			raw, err := json.Marshal(tt.resource)
			require.NoError(t, err)
			request := admissionv1.AdmissionRequest{
				Object: runtime.RawExtension{
					Object: &resource,
					Raw:    raw,
				},
			}

			handler := RequestHandler{
				Validator: validator,
				Logger:    logr.Discard(),
			}

			responseBuilder := NewResponseBuilder(k8stypes.UID(""))

			got, err := handler.handleKongIngress(t.Context(), request, responseBuilder)
			require.NoError(t, err)
			require.True(t, got.Allowed)
			require.Equal(t, tt.wantWarnings, got.Warnings)
		})
	}
}

func TestHandleService(t *testing.T) {
	tests := []struct {
		name         string
		resource     corev1.Service
		validator    KongHTTPValidator
		wantWarnings []string
		wantMessage  string
		isAllowed    bool
	}{
		{
			name: "has kongingress annotation",
			resource: corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.ConfigurationKey: "test",
					},
				},
			},

			wantWarnings: []string{
				fmt.Sprintf(serviceWarning, annotations.AnnotationPrefix+annotations.ConfigurationKey,
					configurationv1beta1.KongUpstreamPolicyAnnotationKey),
			},
			isAllowed: true,
		},
		{
			name: "has upstream policy annotation",
			resource: corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + configurationv1beta1.KongUpstreamPolicyAnnotationKey: "test",
					},
				},
			},
			isAllowed: true,
		},
		{
			name: "passes when many plugins of the same type are attached",
			resource: corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey: "plugin1, plugin2, plugin3",
					},
				},
			},
			validator: KongHTTPValidator{
				ManagerClient: func() client.Client {
					scheme := runtime.NewScheme()
					require.NoError(t, configurationv1.AddToScheme(scheme))
					fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(
						&configurationv1.KongPlugin{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "plugin1",
								Namespace: "default",
							},
							PluginName: "foo",
						},
						&configurationv1.KongPlugin{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "plugin2",
								Namespace: "default",
							},
							PluginName: "bar",
						},
						&configurationv1.KongPlugin{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "plugin3",
								Namespace: "default",
							},
							PluginName: "foo",
						},
					).Build()
					return fakeClient
				}(),
			},
			isAllowed: true,
		},
		{
			name: "pass when many valid plugins are attached",
			resource: corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.PluginsKey:       "plugin1, plugin2, plugin3",
						annotations.AnnotationPrefix + annotations.ConfigurationKey: "test",
					},
				},
			},
			validator: KongHTTPValidator{
				ManagerClient: func() client.Client {
					scheme := runtime.NewScheme()
					require.NoError(t, configurationv1.AddToScheme(scheme))
					fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(
						&configurationv1.KongPlugin{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "plugin1",
								Namespace: "default",
							},
							PluginName: "foo",
						},
						&configurationv1.KongPlugin{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "plugin2",
								Namespace: "default",
							},
							PluginName: "bar",
						},
					).Build()
					return fakeClient
				}(),
			},
			isAllowed: true,
			wantWarnings: []string{
				fmt.Sprintf(serviceWarning, annotations.AnnotationPrefix+annotations.ConfigurationKey,
					configurationv1beta1.KongUpstreamPolicyAnnotationKey),
			},
		},
	}
	for _, tt := range tests {
		resource := tt.resource
		t.Run(tt.name, func(t *testing.T) {
			raw, err := json.Marshal(tt.resource)
			require.NoError(t, err)
			request := admissionv1.AdmissionRequest{
				Object: runtime.RawExtension{
					Object: &resource,
					Raw:    raw,
				},
			}
			handler := RequestHandler{
				Validator: tt.validator,
				Logger:    logr.Discard(),
			}

			responseBuilder := NewResponseBuilder(k8stypes.UID(""))

			got, err := handler.handleService(request, responseBuilder)
			require.NoError(t, err)
			require.Equal(t, tt.isAllowed, got.Allowed)
			require.Equal(t, tt.wantWarnings, got.Warnings)
			require.Equal(t, tt.wantMessage, got.Result.Message)
		})
	}
}

func TestHandleSecret(t *testing.T) {
	testCases := []struct {
		name             string
		secret           *corev1.Secret
		referrers        []client.Object
		validatorOK      bool
		validatorMessage string
		validatorError   error
		expectAllowed    bool
		expectStatusCode int32
		expectMessage    string
		expectError      bool
	}{
		{
			name: "secret used as a credential and passes the validation of credential",
			secret: &corev1.Secret{
				TypeMeta: secretTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "credential-0",
					Labels: map[string]string{
						"konghq.com/credential": "basic-auth",
					},
				},
				Data: map[string][]byte{
					"username": []byte("user"),
					"password": []byte("password"),
				},
			},
			validatorOK:   true,
			expectAllowed: true,
		},
		{
			name: "secret used as credential and fails the validation of credential",
			secret: &corev1.Secret{
				TypeMeta: secretTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "credential-1",
					Labels: map[string]string{
						"konghq.com/credential": "basic-auth",
					},
				},
				Data: map[string][]byte{
					"username": []byte("user"),
					"password": []byte("password"),
				},
			},
			validatorOK:      false,
			validatorMessage: "invalid credential",
			expectAllowed:    false,
			expectStatusCode: http.StatusBadRequest,
			expectMessage:    "invalid credential",
		},
		{
			name: "secret with not supported type of credential is ignored",
			secret: &corev1.Secret{
				TypeMeta: secretTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "credential-0",
					Labels: map[string]string{
						"konghq.com/credential": "whatever-credential",
					},
				},
				Data: map[string][]byte{
					"username": []byte("user"),
					"password": []byte("password"),
				},
			},
			validatorOK:   true,
			expectAllowed: true,
		},
		{
			name: "secret used as KongPlugin config and KongClusterPlugin and passes validation of both CRDs",
			secret: &corev1.Secret{
				TypeMeta: secretTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "plugin-conf",
					Labels: map[string]string{
						labels.ValidateLabel: "plugin",
					},
				},
			},
			referrers: []client.Object{
				&configurationv1.KongPlugin{
					TypeMeta: kongPluginTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "plugin-0",
					},
					PluginName: "test-plugin",
				},
				&configurationv1.KongClusterPlugin{
					TypeMeta: kongClusterPluginTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster-plugin-0",
						Labels: map[string]string{
							labels.ValidateLabel: "plugin",
						},
					},
					PluginName: "test-plugin",
				},
			},
			validatorOK:   true,
			expectAllowed: true,
		},
		{
			name: "secret used as KongPlugin config and fails validation of KongPlugin",
			secret: &corev1.Secret{
				TypeMeta: secretTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "plugin-conf",
					Labels: map[string]string{
						labels.ValidateLabel: "plugin",
					},
				},
			},
			referrers: []client.Object{
				&configurationv1.KongPlugin{
					TypeMeta: kongPluginTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "plugin-0",
					},
					PluginName: "test-plugin",
				},
			},
			validatorOK:      false,
			validatorMessage: "invalid KongPlugin",
			expectAllowed:    false,
			expectStatusCode: http.StatusBadRequest,
			expectMessage:    "Change on secret will generate invalid configuration for KongPlugin default/plugin-0: invalid KongPlugin",
		},
		{
			name: "secret used as KongClusterPlugin config and fails validation of KongClusterPlugin",
			secret: &corev1.Secret{
				TypeMeta: secretTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "plugin-conf",
					Labels: map[string]string{
						labels.ValidateLabel: "plugin",
					},
				},
			},
			referrers: []client.Object{
				&configurationv1.KongClusterPlugin{
					TypeMeta: kongClusterPluginTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster-plugin-0",
					},
					PluginName: "test-plugin",
				},
			},
			validatorOK:      false,
			validatorMessage: "invalid KongClusterPlugin",
			expectAllowed:    false,
			expectStatusCode: http.StatusBadRequest,
			expectMessage:    "Change on secret will generate invalid configuration for KongClusterPlugin cluster-plugin-0: invalid KongClusterPlugin",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator := KongFakeValidator{
				Result:  tc.validatorOK,
				Message: tc.validatorMessage,
				Error:   tc.validatorError,
			}
			raw, err := json.Marshal(tc.secret)
			require.NoError(t, err)
			request := admissionv1.AdmissionRequest{
				Object: runtime.RawExtension{
					Object: tc.secret,
					Raw:    raw,
				},
				Operation: admissionv1.Update,
			}

			logger := testr.NewWithOptions(t, testr.Options{Verbosity: logging.DebugLevel})
			referenceIndexer := ctrlref.NewCacheIndexers(logger)

			handler := RequestHandler{
				Validator:         validator,
				Logger:            logger,
				ReferenceIndexers: referenceIndexer,
			}
			for _, obj := range tc.referrers {
				err := handler.ReferenceIndexers.SetObjectReference(obj, tc.secret)
				require.NoError(t, err)
			}

			responseBuilder := NewResponseBuilder(k8stypes.UID(""))
			got, err := handler.handleSecret(t.Context(), request, responseBuilder)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectAllowed, got.Allowed, "should return expected result of allowed")
			require.Equal(t, tc.expectStatusCode, got.Result.Code)
			if len(tc.expectMessage) > 0 {
				require.Contains(t, got.Result.Message, tc.expectMessage)
			}
		})
	}
}
