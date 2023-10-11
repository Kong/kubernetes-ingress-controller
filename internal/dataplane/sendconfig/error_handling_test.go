package sendconfig

import (
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestParseFlatEntityErrors(t *testing.T) {
	logger := zapr.NewLogger(zap.NewNop())
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
					Name:       "httpbin",
					Namespace:  "67338dc2-31fd-47b6-85a9-9c11d347d090",
					Kind:       "Ingress",
					APIVersion: "networking.k8s.io/v1",
					UID:        "ea569579-f7e9-4d4e-973b-b207bfb848d8",
					Problems: map[string]string{
						"methods": "cannot set methods when protocols is grpc or grpcs",
					},
				},
				{
					Name:       "httpbin",
					Namespace:  "67338dc2-31fd-47b6-85a9-9c11d347d090",
					Kind:       "Service",
					APIVersion: "v1",
					UID:        "e7e5c93e-4d56-4cc3-8f4f-ff1fcbe95eb2",
					Problems: map[string]string{
						"service:67338dc2-31fd-47b6-85a9-9c11d347d090.httpbin.httpbin.80": "failed conditional validation given value of field protocol",
						"path": "value must be null",
					},
				},
			},
			wantErr: false,
			body: []byte(`{
  "name": "invalid declarative configuration",
  "fields": {
    "services": [
      {
        "@entity": [
          "failed conditional validation given value of field protocol"
        ],
        "path": "value must be null"
      }
    ]
  },
  "flattened_errors": [
    {
      "entity_type": "route",
      "entity_name": "67338dc2-31fd-47b6-85a9-9c11d347d090.httpbin.httpbin..80",
      "entity": {
        "paths": [
          "/bar/",
          "~/bar$"
        ],
        "methods": [
          "GET"
        ],
        "response_buffering": true,
        "tags": [
          "k8s-name:httpbin",
          "k8s-namespace:67338dc2-31fd-47b6-85a9-9c11d347d090",
          "k8s-kind:Ingress",
          "k8s-uid:ea569579-f7e9-4d4e-973b-b207bfb848d8",
          "k8s-group:networking.k8s.io",
          "k8s-version:v1"
        ],
        "name": "67338dc2-31fd-47b6-85a9-9c11d347d090.httpbin.httpbin..80",
        "request_buffering": true,
        "preserve_host": true,
        "https_redirect_status_code": 426,
        "path_handling": "v0",
        "protocols": [
          "grpcs"
        ],
        "regex_priority": 0
      },
      "errors": [
        {
          "message": "cannot set methods when protocols is grpc or grpcs",
          "type": "field",
          "field": "methods"
        }
      ],
      "entity_tags": [
        "k8s-name:httpbin",
        "k8s-namespace:67338dc2-31fd-47b6-85a9-9c11d347d090",
        "k8s-kind:Ingress",
        "k8s-uid:ea569579-f7e9-4d4e-973b-b207bfb848d8",
        "k8s-group:networking.k8s.io",
        "k8s-version:v1"
      ]
    },
    {
      "entity_type": "service",
      "entity_name": "67338dc2-31fd-47b6-85a9-9c11d347d090.httpbin.httpbin.80",
      "entity": {
        "protocol": "tcp",
        "tags": [
          "k8s-name:httpbin",
          "k8s-namespace:67338dc2-31fd-47b6-85a9-9c11d347d090",
          "k8s-kind:Service",
          "k8s-uid:e7e5c93e-4d56-4cc3-8f4f-ff1fcbe95eb2",
          "k8s-group:",
          "k8s-version:v1"
        ],
        "retries": 5,
        "connect_timeout": 60000,
        "path": "/aitmatov",
        "name": "67338dc2-31fd-47b6-85a9-9c11d347d090.httpbin.httpbin.80",
        "read_timeout": 60000,
        "port": 80,
        "host": "httpbin.67338dc2-31fd-47b6-85a9-9c11d347d090.80.svc",
        "write_timeout": 60000
      },
      "errors": [
        {
          "message": "failed conditional validation given value of field protocol",
          "type": "entity"
        },
        {
          "message": "value must be null",
          "type": "field",
          "field": "path"
        }
      ],
      "entity_tags": [
        "k8s-name:httpbin",
        "k8s-namespace:67338dc2-31fd-47b6-85a9-9c11d347d090",
        "k8s-kind:Service",
        "k8s-uid:e7e5c93e-4d56-4cc3-8f4f-ff1fcbe95eb2",
        "k8s-group:",
        "k8s-version:v1"
      ]
    }
  ],
  "message": "declarative config is invalid: {services={{[\"@entity\"]={\"failed conditional validation given value of field protocol\"},path=\"value must be null\"}}}",
  "code": 14
}`),
		},
		{
			name:    "empty response body",
			body:    nil,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFlatEntityErrors(tt.body, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFlatEntityErrors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
