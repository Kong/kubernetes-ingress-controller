package sendconfig_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
)

const validFlattenedErrorsResponse = `{
	"code": 14,
	"name": "invalid declarative configuration",
	"flattened_errors": [
		{
			"entity_name": "ingress.httpbin.httpbin..80",
			"entity_tags": [
				"k8s-name:httpbin",
				"k8s-namespace:default",
				"k8s-kind:Ingress",
				"k8s-uid:7b3f3b3b-0b3b-4b3b-8b3b-3b3b3b3b3b3b",
				"k8s-group:networking.k8s.io",
				"k8s-version:v1"
			],
			"errors": [
				{
					"field": "methods",
					"type": "field",
					"message": "cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'"
				}
			],
			"entity": {
				"regex_priority": 0,
				"preserve_host": true,
				"name": "ingress.httpbin.httpbin..80",
				"protocols": [
					"grpcs"
				],
				"https_redirect_status_code": 426,
				"request_buffering": true,
				"tags": [
					"k8s-name:httpbin",
					"k8s-namespace:default",
					"k8s-kind:Ingress",
					"k8s-uid:7b3f3b3b-0b3b-4b3b-8b3b-3b3b3b3b3b3b",
					"k8s-group:networking.k8s.io",
					"k8s-version:v1"
				],
				"path_handling": "v0",
				"response_buffering": true,
				"methods": [
					"GET"
				],
				"paths": [
					"/bar/",
					"~/bar$"
				]
			},
			"entity_type": "route"
		}
	],
	"message": "declarative config is invalid: {}",
	"fields": {}
}`

type mockConfigService struct {
	responseBody []byte
	err          error
}

func (m *mockConfigService) ReloadDeclarativeRawConfig(context.Context, io.Reader, bool, bool) ([]byte, error) {
	return m.responseBody, m.err
}

type mockConfigConverter struct {
	called bool
}

func (m *mockConfigConverter) Convert(*file.Content) sendconfig.DBLessConfig {
	m.called = true
	return sendconfig.DBLessConfig{}
}

func TestUpdateStrategyInMemory(t *testing.T) {
	testCases := []struct {
		name                      string
		configServiceError        error
		configServiceResponseBody []byte
		expectedError             error
	}{
		{
			name:               "no error returned from config service",
			configServiceError: nil,
			expectedError:      nil,
		},
		{
			name:               "unexpected error returned from config service",
			configServiceError: fmt.Errorf("unexpected error"), // e.g. network error
			expectedError:      fmt.Errorf("failed to reload declarative configuration: %w", fmt.Errorf("unexpected error")),
		},
		{
			name:               "APIError 500 returned from config service",
			configServiceError: kong.NewAPIError(500, "internal error"),
			expectedError:      fmt.Errorf("failed to reload declarative configuration: %w", kong.NewAPIError(500, "internal error")),
		},
		{
			name:               "APIError 400 with no resource failures returned from config service",
			configServiceError: kong.NewAPIError(400, "bad request"),
			expectedError:      sendconfig.NewUpdateError(nil, kong.NewAPIError(400, "bad request")),
		},
		{
			name:                      "APIError 400 with resource failures returned from config service",
			configServiceError:        kong.NewAPIError(400, "bad request"),
			configServiceResponseBody: []byte(validFlattenedErrorsResponse),
			expectedError: sendconfig.NewUpdateErrorWithResponseBody(
				[]byte(validFlattenedErrorsResponse),
				[]failures.ResourceFailure{
					lo.Must(failures.NewResourceFailure("invalid methods: cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'", &metav1.PartialObjectMetadata{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "networking.k8s.io/v1",
							Kind:       "Ingress",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "httpbin",
							Namespace: "default",
							UID:       "7b3f3b3b-0b3b-4b3b-8b3b-3b3b3b3b3b3b",
						},
					})),
				},
				kong.NewAPIError(400, "bad request"),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configService := &mockConfigService{responseBody: tc.configServiceResponseBody, err: tc.configServiceError}
			configConverter := &mockConfigConverter{}
			s := sendconfig.NewUpdateStrategyInMemory(configService, configConverter, logr.Discard())
			err := s.Update(context.Background(), sendconfig.ContentWithHash{})
			require.Equal(t, tc.expectedError, err)
			require.True(t, configConverter.called)
		})
	}
}
