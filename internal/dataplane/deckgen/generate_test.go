package deckgen

import (
	"context"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
)

func TestToDeckContent(t *testing.T) {
	testCases := []struct {
		name     string
		params   GenerateDeckContentParams
		input    *kongstate.KongState
		expected *file.Content
	}{
		{
			name:   "empty",
			params: GenerateDeckContentParams{},
			input:  &kongstate.KongState{},
			expected: &file.Content{
				FormatVersion: versions.DeckFileFormatVersion,
			},
		},
		{
			name: "empty, generate stub entity",
			params: GenerateDeckContentParams{
				AppendStubEntityWhenConfigEmpty: true,
			},
			input: &kongstate.KongState{},
			expected: &file.Content{
				FormatVersion: versions.DeckFileFormatVersion,
				Upstreams: []file.FUpstream{
					{
						Upstream: kong.Upstream{
							Name: lo.ToPtr(StubUpstreamName),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToDeckContent(context.Background(), zapr.NewLogger(zap.NewNop()), tc.input, tc.params)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestFillPlugin(t *testing.T) {
	testCases := []struct {
		name          string
		plugin        *file.FPlugin
		schemas       PluginSchemaStore
		expected      *file.FPlugin
		expectedError error
	}{
		{
			name: "Required field provided for plugin",
			plugin: &file.FPlugin{
				Plugin: kong.Plugin{
					Name: lo.ToPtr("plugin"),
					Config: kong.Configuration{
						"endpoint": "https://example.com",
					},
				},
			},
			schemas: &mockPluginSchemaStore{
				map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"protocols": map[string]interface{}{
								"elements": map[string]interface{}{
									"type": "string",
									"one_of": []interface{}{
										"grpc",
										"grpcs",
										"http",
										"https",
									},
								},
								"description": "A set of strings representing HTTP protocols.",
								"type":        "set",
								"default": []interface{}{
									"grpc",
									"grpcs",
									"http",
									"https",
								},
								"required": true,
							},
						},
						map[string]interface{}{
							"config": map[string]interface{}{
								"type": "record",
								"fields": []interface{}{
									map[string]interface{}{
										"endpoint": map[string]interface{}{
											"type":          "string",
											"required":      true,
											"description":   "A string representing a URL, such as https://example.com/path/to/resource?q=search.",
											"referenceable": true,
										},
									},
								},
							},
						},
					},
				},
			},
			expected: &file.FPlugin{
				Plugin: kong.Plugin{
					Name: lo.ToPtr("plugin"),
					Protocols: []*string{
						lo.ToPtr("grpc"),
						lo.ToPtr("grpcs"),
						lo.ToPtr("http"),
						lo.ToPtr("https"),
					},
					Enabled: lo.ToPtr(true),
					Config: kong.Configuration{
						"endpoint": "https://example.com",
					},
				},
			},
		},
		{
			name: "Required field not provided for plugin gets filled in with nil",
			plugin: &file.FPlugin{
				Plugin: kong.Plugin{
					Name: lo.ToPtr("plugin"),
				},
			},
			schemas: &mockPluginSchemaStore{
				map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"protocols": map[string]interface{}{
								"elements": map[string]interface{}{
									"type": "string",
									"one_of": []interface{}{
										"grpc",
										"grpcs",
										"http",
										"https",
									},
								},
								"description": "A set of strings representing HTTP protocols.",
								"type":        "set",
								"default": []interface{}{
									"grpc",
									"grpcs",
									"http",
									"https",
								},
								"required": true,
							},
						},
						map[string]interface{}{
							"config": map[string]interface{}{
								"type": "record",
								"fields": []interface{}{
									map[string]interface{}{
										"endpoint": map[string]interface{}{
											"type":          "string",
											"required":      true,
											"description":   "A string representing a URL, such as https://example.com/path/to/resource?q=search.",
											"referenceable": true,
										},
									},
								},
							},
						},
					},
				},
			},
			expected: &file.FPlugin{
				Plugin: kong.Plugin{
					Name: lo.ToPtr("plugin"),
					Protocols: []*string{
						lo.ToPtr("grpc"),
						lo.ToPtr("grpcs"),
						lo.ToPtr("http"),
						lo.ToPtr("https"),
					},
					Enabled: lo.ToPtr(true),
					Config: kong.Configuration{
						"endpoint": nil,
					},
				},
			},
		},
		{
			// NOTE: This would fail for go-kong v0.52.0 and older.
			name: "OpenTelemetry plugin for Gateway 3.7.x",
			plugin: &file.FPlugin{
				Plugin: kong.Plugin{
					Name: lo.ToPtr("opentelemetry"),
				},
			},
			schemas: &mockPluginSchemaStore{
				schema: map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"protocols": map[string]interface{}{
								"elements": map[string]interface{}{
									"type": "string",
									"one_of": []interface{}{
										"grpc",
										"grpcs",
										"http",
										"https",
									},
								},
								"description": "A set of strings representing HTTP protocols.",
								"type":        "set",
								"default": []interface{}{
									"grpc",
									"grpcs",
									"http",
									"https",
								},
								"required": true,
							},
						},
						map[string]interface{}{
							"config": map[string]interface{}{
								"type": "record",
								"fields": []interface{}{
									map[string]interface{}{
										"endpoint": map[string]interface{}{
											"type":          "string",
											"required":      true,
											"description":   "A string representing a URL, such as https://example.com/path/to/resource?q=search.",
											"referenceable": true,
										},
									},
									map[string]interface{}{
										"headers": map[string]interface{}{
											"description": "The custom headers to be added in the HTTP request sent to the OTLP server. This setting is useful for adding the authentication headers (token) for the APM backend.",
											"type":        "map",
											"values": map[string]interface{}{
												"type":          "string",
												"referenceable": true,
											},
											"keys": map[string]interface{}{
												"type":        "string",
												"description": "A string representing an HTTP header name.",
											},
										},
									},
									map[string]interface{}{
										"resource_attributes": map[string]interface{}{
											"type": "map",
											"keys": map[string]interface{}{
												"type":     "string",
												"required": true,
											},
											"values": map[string]interface{}{
												"type":     "string",
												"required": true,
											},
											"description": "Attributes to add to the OpenTelemetry resource object, following the spec for Semantic Attributes. \nThe following attributes are automatically added:\n- \"service.name\": The name of the service (default:'kong').\n-'service.version': The version of Kong Gateway.\n-'service.instance.id': The node ID of Kong Gateway.\n\nYou can use this property to override default attribute values. For example, to override the default for'service.name', you can specify'{ \"service.name\": \"my-service\" }'.",
										},
									},
									map[string]interface{}{
										"queue": map[string]interface{}{
											"type": "record",
											"fields": []interface{}{
												map[string]interface{}{
													"max_batch_size": map[string]interface{}{
														"type": "integer",
														"between": []interface{}{
															1,
															1000000,
														},
														"default":     1,
														"description": "Maximum number of entries that can be processed at a time.",
													},
												},
												map[string]interface{}{
													"max_coalescing_delay": map[string]interface{}{
														"type": "number",
														"between": []interface{}{
															0,
															3600,
														},
														"default":     1,
														"description": "Maximum number of (fractional) seconds to elapse after the first entry was queued before the queue starts calling the handler.",
													},
												},
												map[string]interface{}{
													"max_entries": map[string]interface{}{
														"type": "integer",
														"between": []interface{}{
															1,
															1000000,
														},
														"default":     10000,
														"description": "Maximum number of entries that can be waiting on the queue.",
													},
												},
												map[string]interface{}{
													"max_bytes": map[string]interface{}{
														"type":        "integer",
														"description": "Maximum number of bytes that can be waiting on a queue, requires string content.",
													},
												},
												map[string]interface{}{
													"max_retry_time": map[string]interface{}{
														"type":        "number",
														"default":     60,
														"description": "Time in seconds before the queue gives up calling a failed handler for a batch.",
													},
												},
												map[string]interface{}{
													"initial_retry_delay": map[string]interface{}{
														"type": "number",
														"between": []interface{}{
															0.001,
															1000000,
														},
														"default":     0.01,
														"description": "Time in seconds before the initial retry is made for a failing batch.",
													},
												},
												map[string]interface{}{
													"max_retry_delay": map[string]interface{}{
														"type": "number",
														"between": []interface{}{
															0.001,
															1000000,
														},
														"default":     60,
														"description": "Maximum time in seconds between retries, caps exponential backoff.",
													},
												},
											},
											"default": map[string]interface{}{
												"max_batch_size": 200,
											},
											"required": true,
										},
									},
									map[string]interface{}{
										"batch_span_count": map[string]interface{}{
											"description": "The number of spans to be sent in a single batch.",
											"type":        "integer",
											"deprecation": map[string]interface{}{
												"old_default":        200,
												"removal_in_version": "4.0",
												"message":            "opentelemetry: config.batch_span_count is deprecated, please use config.queue.max_batch_size instead",
											},
										},
									},
									map[string]interface{}{
										"batch_flush_delay": map[string]interface{}{
											"description": "The delay, in seconds, between two consecutive batches.",
											"type":        "integer",
											"deprecation": map[string]interface{}{
												"old_default":        3,
												"removal_in_version": "4.0",
												"message":            "opentelemetry: config.batch_flush_delay is deprecated, please use config.queue.max_coalescing_delay instead",
											},
										},
									},
									map[string]interface{}{
										"connect_timeout": map[string]interface{}{
											"type":        "integer",
											"description": "An integer representing a timeout in milliseconds. Must be between 0 and 2^31-2.",
											"default":     1000,
											"between": []interface{}{
												0,
												2147483646,
											},
										},
									},
									map[string]interface{}{
										"send_timeout": map[string]interface{}{
											"type":        "integer",
											"description": "An integer representing a timeout in milliseconds. Must be between 0 and 2^31-2.",
											"default":     5000,
											"between": []interface{}{
												0,
												2147483646,
											},
										},
									},
									map[string]interface{}{
										"read_timeout": map[string]interface{}{
											"type":        "integer",
											"description": "An integer representing a timeout in milliseconds. Must be between 0 and 2^31-2.",
											"default":     5000,
											"between": []interface{}{
												0,
												2147483646,
											},
										},
									},
									map[string]interface{}{
										"http_response_header_for_traceid": map[string]interface{}{
											"description": "Specifies a custom header for the'trace_id'. If set, the plugin sets the corresponding header in the response.",
											"type":        "string",
										},
									},
									map[string]interface{}{
										"header_type": map[string]interface{}{
											"deprecation": map[string]interface{}{
												"old_default":        "preserve",
												"removal_in_version": "4.0",
												"message":            "opentelemetry: config.header_type is deprecated, please use config.propagation options instead",
											},
											"one_of": []interface{}{
												"preserve",
												"ignore",
												"b3",
												"b3-single",
												"w3c",
												"jaeger",
												"ot",
												"aws",
												"gcp",
												"datadog",
											},
											"type":        "string",
											"description": "All HTTP requests going through the plugin are tagged with a tracing HTTP request. This property codifies what kind of tracing header the plugin expects on incoming requests.",
											"default":     "preserve",
											"required":    false,
										},
									},
									map[string]interface{}{
										"sampling_rate": map[string]interface{}{
											"between": []interface{}{
												0,
												1,
											},
											"description": "Tracing sampling rate for configuring the probability-based sampler. When set, this value supersedes the global'tracing_sampling_rate' setting from kong.conf.",
											"type":        "number",
											"required":    false,
										},
									},
									map[string]interface{}{
										"propagation": map[string]interface{}{
											"type": "record",
											"fields": []interface{}{
												map[string]interface{}{
													"extract": map[string]interface{}{
														"description": "Header formats used to extract tracing context from incoming requests. If multiple values are specified, the first one found will be used for extraction. If left empty, Kong will not extract any tracing context information from incoming requests and generate a trace with no parent and a new trace ID.",
														"type":        "array",
														"elements": map[string]interface{}{
															"type": "string",
															"one_of": []interface{}{
																"ot",
																"w3c",
																"datadog",
																"b3",
																"gcp",
																"jaeger",
																"aws",
															},
														},
													},
												},
												map[string]interface{}{
													"clear": map[string]interface{}{
														"description": "Header names to clear after context extraction. This allows to extract the context from a certain header and then remove it from the request, useful when extraction and injection are performed on different header formats and the original header should not be sent to the upstream. If left empty, no headers are cleared.",
														"type":        "array",
														"elements": map[string]interface{}{
															"type": "string",
														},
													},
												},
												map[string]interface{}{
													"inject": map[string]interface{}{
														"description": "Header formats used to inject tracing context. The value 'preserve' will use the same header format as the incoming request. If multiple values are specified, all of them will be used during injection. If left empty, Kong will not inject any tracing context information in outgoing requests.",
														"type":        "array",
														"elements": map[string]interface{}{
															"type": "string",
															"one_of": []interface{}{
																"preserve",
																"ot",
																"w3c",
																"datadog",
																"b3",
																"gcp",
																"b3-single",
																"jaeger",
																"aws",
															},
														},
													},
												},
												map[string]interface{}{
													"default_format": map[string]interface{}{
														"description": "The default header format to use when extractors did not match any format in the incoming headers and'inject' is configured with the value:'preserve'. This can happen when no tracing header was found in the request, or the incoming tracing header formats were not included in'extract'.",
														"one_of": []interface{}{
															"ot",
															"w3c",
															"datadog",
															"b3",
															"gcp",
															"b3-single",
															"jaeger",
															"aws",
														},
														"type":     "string",
														"required": true,
													},
												},
											},
											"default": map[string]interface{}{
												"default_format": "w3c",
											},
											"required": true,
										},
									},
								},
								"required": true,
							},
						},
					},
					"entity_checks": []interface{}{},
				},
			},
			expected: &file.FPlugin{
				Plugin: kong.Plugin{
					Name: lo.ToPtr("opentelemetry"),
					Protocols: []*string{
						lo.ToPtr("grpc"),
						lo.ToPtr("grpcs"),
						lo.ToPtr("http"),
						lo.ToPtr("https"),
					},
					Enabled: lo.ToPtr(true),
					Config: kong.Configuration{
						"endpoint":                         nil,
						"batch_flush_delay":                nil,
						"batch_span_count":                 nil,
						"connect_timeout":                  float64(1000),
						"header_type":                      "preserve",
						"headers":                          nil,
						"http_response_header_for_traceid": nil,
						"propagation": map[string]interface{}{
							"clear":          nil,
							"default_format": "w3c",
							"extract":        nil,
							"inject":         nil,
						},
						"queue": map[string]interface{}{
							"initial_retry_delay":  float64(0.01),
							"max_batch_size":       float64(200),
							"max_bytes":            nil,
							"max_coalescing_delay": float64(1),
							"max_entries":          float64(10000),
							"max_retry_delay":      float64(60),
							"max_retry_time":       float64(60),
						},
						"read_timeout":        float64(5000),
						"resource_attributes": nil,
						"sampling_rate":       nil,
						"send_timeout":        float64(5000),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plugin := tc.plugin.DeepCopy()
			err := fillPlugin(context.Background(), plugin, tc.schemas)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, plugin)
			}
		})
	}
}

type mockPluginSchemaStore struct {
	schema map[string]interface{}
}

func (m *mockPluginSchemaStore) Schema(_ context.Context, _ string) (map[string]interface{}, error) {
	return m.schema, nil
}
