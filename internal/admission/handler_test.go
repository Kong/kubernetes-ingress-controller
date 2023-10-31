package admission

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

func TestHandleKongIngress(t *testing.T) {
	tests := []struct {
		name         string
		resource     kongv1.KongIngress
		wantWarnings []string
	}{
		{
			name: "has proxy",
			resource: kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{},
			},
			wantWarnings: []string{proxyWarning},
		},
		{
			name: "has route",
			resource: kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{},
			},
			wantWarnings: []string{routeWarning},
		},
		{
			name: "has upstream",
			resource: kongv1.KongIngress{
				Upstream: &kongv1.KongIngressUpstream{},
			},
			wantWarnings: []string{upstreamWarning},
		},
		{
			name: "has everything",
			resource: kongv1.KongIngress{
				Proxy:    &kongv1.KongIngressService{},
				Route:    &kongv1.KongIngressRoute{},
				Upstream: &kongv1.KongIngressUpstream{},
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

			got, err := handler.handleKongIngress(context.Background(), request, responseBuilder)
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
		wantWarnings []string
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
					kongv1beta1.KongUpstreamPolicyAnnotationKey),
			},
		},
		{
			name: "has upstream policy annotation",
			resource: corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + kongv1beta1.KongUpstreamPolicyAnnotationKey: "test",
					},
				},
			},
			wantWarnings: nil,
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

			got, err := handler.handleService(context.Background(), request, responseBuilder)
			require.NoError(t, err)
			require.True(t, got.Allowed)
			require.Equal(t, tt.wantWarnings, got.Warnings)
		})
	}
}
