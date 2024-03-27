package deckgen_test

import (
	"encoding/json"
	"testing"

	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
)

func BenchmarkDeckgenGenerateSHA(b *testing.B) {
	var targetContent file.Content
	require.NoError(b, json.Unmarshal([]byte(configJSON), &targetContent))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bb, err := deckgen.GenerateSHA(&targetContent)
		require.NoError(b, err)
		_ = bb
	}
}

const configJSON = `{
	"_format_version": "3.0",
	"_info": {
		"select_tags": [
			"managed-by-ingress-controller"
		],
		"defaults": {}
	},
	"services": [
		{
			"connect_timeout": 60000,
			"host": "httproute.default.httproute-testing.0",
			"name": "httproute.default.httproute-testing.0",
			"protocol": "http",
			"read_timeout": 60000,
			"retries": 5,
			"write_timeout": 60000,
			"tags": [
				"k8s-name:httproute-testing",
				"k8s-namespace:default",
				"k8s-kind:HTTPRoute",
				"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
				"k8s-group:gateway.networking.k8s.io",
				"k8s-version:v1beta1"
			],
			"routes": [
				{
					"name": "httproute.default.httproute-testing.0.0",
					"paths": [
						"/httproute-testing"
					],
					"path_handling": "v0",
					"preserve_host": true,
					"protocols": [
						"http",
						"https"
					],
					"strip_path": true,
					"tags": [
						"k8s-name:httproute-testing",
						"k8s-namespace:default",
						"k8s-kind:HTTPRoute",
						"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
						"k8s-group:gateway.networking.k8s.io",
						"k8s-version:v1beta1"
					],
					"https_redirect_status_code": 426
				}
			]
		},
		{
			"connect_timeout": 60000,
			"host": "httproute.default.httproute-testing.1",
			"name": "httproute.default.httproute-testing.1",
			"protocol": "http",
			"read_timeout": 60000,
			"retries": 5,
			"write_timeout": 60000,
			"tags": [
				"k8s-name:httproute-testing",
				"k8s-namespace:default",
				"k8s-kind:HTTPRoute",
				"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
				"k8s-group:gateway.networking.k8s.io",
				"k8s-version:v1beta1"
			],
			"routes": [
				{
					"name": "httproute.default.httproute-testing.1.0",
					"paths": [
						"/httproute-testing"
					],
					"path_handling": "v0",
					"preserve_host": true,
					"protocols": [
						"http",
						"https"
					],
					"strip_path": true,
					"tags": [
						"k8s-name:httproute-testing",
						"k8s-namespace:default",
						"k8s-kind:HTTPRoute",
						"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
						"k8s-group:gateway.networking.k8s.io",
						"k8s-version:v1beta1"
					],
					"https_redirect_status_code": 426
				}
			]
		},
		{
			"connect_timeout": 60000,
			"host": "httproute.default.httproute-testing.2",
			"name": "httproute.default.httproute-testing.2",
			"protocol": "http",
			"read_timeout": 60000,
			"retries": 5,
			"write_timeout": 60000,
			"tags": [
				"k8s-name:httproute-testing",
				"k8s-namespace:default",
				"k8s-kind:HTTPRoute",
				"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
				"k8s-group:gateway.networking.k8s.io",
				"k8s-version:v1beta1"
			],
			"routes": [
				{
					"name": "httproute.default.httproute-testing.2.0",
					"paths": [
						"/httproute-testing"
					],
					"path_handling": "v0",
					"preserve_host": true,
					"protocols": [
						"http",
						"https"
					],
					"strip_path": true,
					"tags": [
						"k8s-name:httproute-testing",
						"k8s-namespace:default",
						"k8s-kind:HTTPRoute",
						"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
						"k8s-group:gateway.networking.k8s.io",
						"k8s-version:v1beta1"
					],
					"https_redirect_status_code": 426
				}
			]
		},
		{
			"connect_timeout": 60000,
			"host": "httproute.default.httproute-testing.3",
			"name": "httproute.default.httproute-testing.3",
			"protocol": "http",
			"read_timeout": 60000,
			"retries": 5,
			"write_timeout": 60000,
			"tags": [
				"k8s-name:httproute-testing",
				"k8s-namespace:default",
				"k8s-kind:HTTPRoute",
				"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
				"k8s-group:gateway.networking.k8s.io",
				"k8s-version:v1beta1"
			],
			"routes": [
				{
					"name": "httproute.default.httproute-testing.3.0",
					"paths": [
						"/httproute-testing"
					],
					"path_handling": "v0",
					"preserve_host": true,
					"protocols": [
						"http",
						"https"
					],
					"strip_path": true,
					"tags": [
						"k8s-name:httproute-testing",
						"k8s-namespace:default",
						"k8s-kind:HTTPRoute",
						"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
						"k8s-group:gateway.networking.k8s.io",
						"k8s-version:v1beta1"
					],
					"https_redirect_status_code": 426
				}
			]
		},
		{
			"connect_timeout": 60000,
			"host": "httproute.default.httproute-testing.4",
			"name": "httproute.default.httproute-testing.4",
			"protocol": "http",
			"read_timeout": 60000,
			"retries": 5,
			"write_timeout": 60000,
			"tags": [
				"k8s-name:httproute-testing",
				"k8s-namespace:default",
				"k8s-kind:HTTPRoute",
				"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
				"k8s-group:gateway.networking.k8s.io",
				"k8s-version:v1beta1"
			],
			"routes": [
				{
					"name": "httproute.default.httproute-testing.4.0",
					"paths": [
						"/httproute-testing"
					],
					"path_handling": "v0",
					"preserve_host": true,
					"protocols": [
						"http",
						"https"
					],
					"strip_path": true,
					"tags": [
						"k8s-name:httproute-testing",
						"k8s-namespace:default",
						"k8s-kind:HTTPRoute",
						"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
						"k8s-group:gateway.networking.k8s.io",
						"k8s-version:v1beta1"
					],
					"https_redirect_status_code": 426
				}
			]
		}
	],
	"upstreams": [
		{
			"name": "httproute.default.httproute-testing.0",
			"algorithm": "round-robin",
			"tags": [
				"k8s-name:httproute-testing",
				"k8s-namespace:default",
				"k8s-kind:HTTPRoute",
				"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
				"k8s-group:gateway.networking.k8s.io",
				"k8s-version:v1beta1"
			],
			"targets": [
				{
					"target": "10.244.0.11:80",
					"weight": 75
				}
			]
		},
		{
			"name": "httproute.default.httproute-testing.1",
			"algorithm": "round-robin",
			"tags": [
				"k8s-name:httproute-testing",
				"k8s-namespace:default",
				"k8s-kind:HTTPRoute",
				"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
				"k8s-group:gateway.networking.k8s.io",
				"k8s-version:v1beta1"
			],
			"targets": [
				{
					"target": "10.244.0.11:80",
					"weight": 75
				}
			]
		},
		{
			"name": "httproute.default.httproute-testing.2",
			"algorithm": "round-robin",
			"tags": [
				"k8s-name:httproute-testing",
				"k8s-namespace:default",
				"k8s-kind:HTTPRoute",
				"k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
				"k8s-group:gateway.networking.k8s.io",
				"k8s-version:v1beta1"
			],
			"targets": [
				{
					"target": "10.244.0.11:80",
					"weight": 75
				}
			]
		}
	]
}`

func TestIsContentEmpty(t *testing.T) {
	testCases := []struct {
		name    string
		content *file.Content
		want    bool
	}{
		{
			name: "non-empty content",
			content: &file.Content{
				Upstreams: []file.FUpstream{
					{
						Upstream: kong.Upstream{
							Name: kong.String("test"),
						},
					},
				},
			},
			want: false,
		},
		{
			name:    "empty content",
			content: &file.Content{},
			want:    true,
		},
		{
			name: "empty with version and info",
			content: &file.Content{
				FormatVersion: "1.1",
				Info: &file.Info{
					SelectorTags: []string{"tag1", "tag2"},
				},
			},
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := deckgen.IsContentEmpty(tc.content)
			require.Equal(t, tc.want, got)
		})
	}
}
