package admission_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission"
)

func TestResponseBuilder(t *testing.T) {
	someUID := k8stypes.UID("f0c595d0-09f5-4dca-95b7-80869e250068")
	testCases := []struct {
		name             string
		modifyBuilderFn  func(b *admission.ResponseBuilder)
		expectedResponse *admissionv1.AdmissionResponse
	}{
		{
			name: "allowed",
			modifyBuilderFn: func(b *admission.ResponseBuilder) {
				b.Allowed(true)
			},
			expectedResponse: &admissionv1.AdmissionResponse{
				UID:     someUID,
				Allowed: true,
				Result:  &metav1.Status{},
			},
		},
		{
			name: "disallowed",
			modifyBuilderFn: func(b *admission.ResponseBuilder) {
				b.Allowed(false)
			},
			expectedResponse: &admissionv1.AdmissionResponse{
				UID:     someUID,
				Allowed: false,
				Result: &metav1.Status{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "disallowed with a message",
			modifyBuilderFn: func(b *admission.ResponseBuilder) {
				b.Allowed(false).WithMessage("something is invalid")
			},
			expectedResponse: &admissionv1.AdmissionResponse{
				UID:     someUID,
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusBadRequest,
					Message: "something is invalid",
				},
			},
		},
		{
			name: "allowed with a warning",
			modifyBuilderFn: func(b *admission.ResponseBuilder) {
				b.Allowed(true).WithWarning("you've been warned")
			},
			expectedResponse: &admissionv1.AdmissionResponse{
				UID:      someUID,
				Allowed:  true,
				Result:   &metav1.Status{},
				Warnings: []string{"you've been warned"},
			},
		},
		{
			name: "allowed with multiple warnings",
			modifyBuilderFn: func(b *admission.ResponseBuilder) {
				b.Allowed(true).
					WithWarning("you've been warned once").
					WithWarning("you've been warned twice")
			},
			expectedResponse: &admissionv1.AdmissionResponse{
				UID:     someUID,
				Allowed: true,
				Result:  &metav1.Status{},
				Warnings: []string{
					"you've been warned once",
					"you've been warned twice",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := admission.NewResponseBuilder(someUID)

			tc.modifyBuilderFn(builder)

			builtResponse := builder.Build()
			require.Equal(t, tc.expectedResponse, builtResponse)
		})
	}
}
