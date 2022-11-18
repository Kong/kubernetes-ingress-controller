package sendconfig

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"

	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/stretchr/testify/require"
)

func TestParseFlatEntityErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    []byte
		want    []ResourceError
		wantErr bool
	}{
		{
			name: "a route nested under a service, with two and one errors, respectively",
			want: []ResourceError{
				{
					Name:       "scallion",
					Namespace:  "default",
					Kind:       "Ingress",
					APIVersion: "networking.k8s.io/v1",
					UID:        "d7300db1-14eb-5a09-b594-2db904ed8eca",
					Problems: map[string]string{
						"methods": "cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'",
					},
				},
				{
					Name:       "turnip",
					Namespace:  "default",
					Kind:       "Ingress",
					APIVersion: "networking.k8s.io/v1",
					UID:        "d7300db1-14eb-5a09-b594-2db904ed8eca",
					Problems: map[string]string{
						"methods[0]": "expected a string",
						"methods[1]": "expected a string",
					},
				},
				{
					Name:       "radish",
					Namespace:  "default",
					Kind:       "Service",
					UID:        "b8aa692c-6d8d-580e-a767-a7dbc1f58344",
					APIVersion: "v1",
					Problems: map[string]string{
						"read_timeout":  "expected an integer",
						"write_timeout": "expected an integer",
					},
				},
			},
			wantErr: false,
			body: []byte(`{
    "code": 14,
    "fields": {
        "flattened": [
            {
                "entity_id": "d7300db1-14eb-5a09-b594-2db904ed8eca",
                "entity_name": "default.demo.00",
                "entity_tags": [
                    "k8s-name:scallion",
                    "k8s-namespace:default",
                    "k8s-kind:Ingress",
                    "k8s-uid:d7300db1-14eb-5a09-b594-2db904ed8eca",
                    "k8s-version:v1",
                    "k8s-group:networking.k8s.io"
                ],
                "errors": [
                    {
                        "field": "methods",
                        "message": "cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'"
                    }
                ]
            },
            {
                "entity_id": "d7300db1-14eb-5a09-b594-2db904ed8eca",
                "entity_name": "default.demo.01",
                "entity_tags": [
                    "k8s-name:turnip",
                    "k8s-namespace:default",
                    "k8s-kind:Ingress",
                    "k8s-uid:d7300db1-14eb-5a09-b594-2db904ed8eca",
                    "k8s-version:v1",
                    "k8s-group:networking.k8s.io"
                ],
                "errors": [
                    {
                        "field": "methods",
                        "messages": [
                            "expected a string",
                            "expected a string"
                        ]
                    }
                ]
            },
            {
                "entity_id": "b8aa692c-6d8d-580e-a767-a7dbc1f58344",
                "entity_name": "default.echo.pnum-80",
                "entity_tags": [
                    "k8s-name:radish",
                    "k8s-namespace:default",
                    "k8s-kind:Service",
                    "k8s-uid:b8aa692c-6d8d-580e-a767-a7dbc1f58344",
                    "k8s-version:v1"
                ],
                "errors": [
                    {
                        "field": "write_timeout",
                        "message": "expected an integer"
                    },
                    {
                        "field": "read_timeout",
                        "message": "expected an integer"
                    }
                ]
            }
        ],
        "services": [
            {
                "read_timeout": "expected an integer",
                "routes": [
                    {
                        "methods": "cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'"
                    },
                    {
                        "methods": [
                            "expected a string",
                            "expected a string"
                        ]
                    }
                ],
                "write_timeout": "expected an integer"
            }
        ]
    },
    "message": "declarative config is invalid: {services={{read_timeout=\"expected an integer\",routes={{methods=\"cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs'\"},{methods={\"expected a string\",\"expected a string\"}}},write_timeout=\"expected an integer\"}}}",
    "name": "invalid declarative configuration"
}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFlatEntityErrors(tt.body, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFlatEntityErrors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, got, tt.want)
		})
	}
}

func TestPushFailureReason(t *testing.T) {
	apiConflictErr := kong.NewAPIError(http.StatusConflict, "conflict api error", []byte{})
	networkErr := net.UnknownNetworkError("network error")
	genericError := errors.New("generic error")

	testCases := []struct {
		name           string
		err            error
		expectedReason string
	}{
		{
			name:           "generic_error",
			err:            genericError,
			expectedReason: metrics.FailureReasonOther,
		},
		{
			name:           "api_conflict_error",
			err:            apiConflictErr,
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "api_conflict_error_wrapped",
			err:            fmt.Errorf("wrapped conflict api err: %w", apiConflictErr),
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "deck_config_conflict_error_empty",
			err:            deckConfigConflictError{},
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "deck_config_conflict_error_with_generic_error",
			err:            deckConfigConflictError{genericError},
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "deck_err_array_with_api_conflict_error",
			err:            deckutils.ErrArray{Errors: []error{apiConflictErr}},
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "wrapped_deck_err_array_with_api_conflict_error",
			err:            fmt.Errorf("wrapped: %w", deckutils.ErrArray{Errors: []error{apiConflictErr}}),
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "deck_err_array_with_generic_error",
			err:            deckutils.ErrArray{Errors: []error{genericError}},
			expectedReason: metrics.FailureReasonOther,
		},
		{
			name:           "deck_err_array_empty",
			err:            deckutils.ErrArray{Errors: []error{genericError}},
			expectedReason: metrics.FailureReasonOther,
		},
		{
			name:           "network_error",
			err:            networkErr,
			expectedReason: metrics.FailureReasonNetwork,
		},
		{
			name:           "network_error_wrapped_in_deck_config_conflict_error",
			err:            deckConfigConflictError{networkErr},
			expectedReason: metrics.FailureReasonNetwork,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			reason := pushFailureReason(tc.err)
			require.Equal(t, tc.expectedReason, reason)
		})
	}
}
