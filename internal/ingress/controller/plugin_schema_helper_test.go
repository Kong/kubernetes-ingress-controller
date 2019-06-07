package controller

import (
	"encoding/json"
	"testing"

	"github.com/hbagdi/go-kong/kong"
	"github.com/stretchr/testify/assert"
)

var (
	KeyAuthSchema = `{
		"fields": [
			{
				"key_names": {
					"default": [
						"apikey"
					],
					"elements": {
						"type": "string"
					},
					"required": true,
					"type": "array"
				}
			},
			{
				"hide_credentials": {
					"default": false,
					"type": "boolean"
				}
			},
			{
				"anonymous": {
					"legacy": true,
					"type": "string",
					"uuid": true
				}
			},
			{
				"key_in_body": {
					"default": false,
					"type": "boolean"
				}
			},
			{
				"run_on_preflight": {
					"default": true,
					"type": "boolean"
				}
			}
		]
	}`
	KeyAuthDefaultConfig = `{
		"anonymous": null,
		"hide_credentials": false,
		"key_in_body": false,
		"key_names": [
			"apikey"
		],
		"run_on_preflight": true
	}`
	StatsDSchema = `{
		"fields": [
			{
				"host": {
					"default": "localhost",
					"type": "string"
				}
			},
			{
				"port": {
					"between": [
						0,
						65535
					],
					"default": 8125,
					"type": "integer"
				}
			},
			{
				"prefix": {
					"default": "kong",
					"type": "string"
				}
			},
			{
				"metrics": {
					"default": [
						{
							"name": "request_count",
							"sample_rate": 1,
							"stat_type": "counter"
						},
						{
							"name": "latency",
							"stat_type": "timer"
						},
						{
							"name": "request_size",
							"stat_type": "timer"
						},
						{
							"name": "status_count",
							"sample_rate": 1,
							"stat_type": "counter"
						},
						{
							"name": "response_size",
							"stat_type": "timer"
						},
						{
							"consumer_identifier": "custom_id",
							"name": "unique_users",
							"stat_type": "set"
						},
						{
							"consumer_identifier": "custom_id",
							"name": "request_per_user",
							"sample_rate": 1,
							"stat_type": "counter"
						},
						{
							"name": "upstream_latency",
							"stat_type": "timer"
						},
						{
							"name": "kong_latency",
							"stat_type": "timer"
						},
						{
							"consumer_identifier": "custom_id",
							"name": "status_count_per_user",
							"sample_rate": 1,
							"stat_type": "counter"
						}
					],
					"elements": {
						"entity_checks": [
							{
								"conditional": {
									"if_field": "name",
									"if_match": {
										"eq": "unique_users"
									},
									"then_field": "stat_type",
									"then_match": {
										"eq": "set"
									}
								}
							},
							{
								"conditional": {
									"if_field": "stat_type",
									"if_match": {
										"one_of": [
											"counter",
											"gauge"
										]
									},
									"then_field": "sample_rate",
									"then_match": {
										"required": true
									}
								}
							},
							{
								"conditional": {
									"if_field": "name",
									"if_match": {
										"one_of": [
											"status_count_per_user",
											"request_per_user",
											"unique_users"
										]
									},
									"then_field": "consumer_identifier",
									"then_match": {
										"required": true
									}
								}
							},
							{
								"conditional": {
									"if_field": "name",
									"if_match": {
										"one_of": [
											"status_count",
											"status_count_per_user",
											"request_per_user"
										]
									},
									"then_field": "stat_type",
									"then_match": {
										"eq": "counter"
									}
								}
							}
						],
						"fields": [
							{
								"name": {
									"one_of": [
										"kong_latency",
										"latency",
										"request_count",
										"request_per_user",
										"request_size",
										"response_size",
										"status_count",
										"status_count_per_user",
										"unique_users",
										"upstream_latency"
									],
									"required": true,
									"type": "string"
								}
							},
							{
								"stat_type": {
									"one_of": [
										"counter",
										"gauge",
										"histogram",
										"meter",
										"set",
										"timer"
									],
									"required": true,
									"type": "string"
								}
							},
							{
								"sample_rate": {
									"gt": 0,
									"type": "number"
								}
							},
							{
								"consumer_identifier": {
									"one_of": [
										"consumer_id",
										"custom_id",
										"username"
									],
									"type": "string"
								}
							}
						],
						"type": "record"
					},
					"type": "array"
				}
			}
		]
	}`
	StatsDDefaultConfig = `{
        "host": "localhost",
        "metrics": [
            {
                "name": "request_count",
                "sample_rate": 1,
                "stat_type": "counter"
            },
            {
                "name": "latency",
                "stat_type": "timer"
            },
            {
                "name": "request_size",
                "stat_type": "timer"
            },
            {
                "name": "status_count",
                "sample_rate": 1,
                "stat_type": "counter"
            },
            {
                "name": "response_size",
                "stat_type": "timer"
            },
            {
                "consumer_identifier": "custom_id",
                "name": "unique_users",
                "stat_type": "set"
            },
            {
                "consumer_identifier": "custom_id",
                "name": "request_per_user",
                "sample_rate": 1,
                "stat_type": "counter"
            },
            {
                "name": "upstream_latency",
                "stat_type": "timer"
            },
            {
                "name": "kong_latency",
                "stat_type": "timer"
            },
            {
                "consumer_identifier": "custom_id",
                "name": "status_count_per_user",
                "sample_rate": 1,
                "stat_type": "counter"
            }
        ],
        "port": 8125,
        "prefix": "kong"
	}`
	RequestTransformerSchema = `{
		"fields": [
			{
				"http_method": {
					"match": "^%u+$",
					"type": "string"
				}
			},
			{
				"remove": {
					"fields": [
						{
							"body": {
								"default": [],
								"elements": {
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"headers": {
								"default": [],
								"elements": {
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"querystring": {
								"default": [],
								"elements": {
									"type": "string"
								},
								"type": "array"
							}
						}
					],
					"required": true,
					"type": "record"
				}
			},
			{
				"rename": {
					"fields": [
						{
							"body": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"headers": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"querystring": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						}
					],
					"required": true,
					"type": "record"
				}
			},
			{
				"replace": {
					"fields": [
						{
							"body": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"headers": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"querystring": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						}
					],
					"required": true,
					"type": "record"
				}
			},
			{
				"add": {
					"fields": [
						{
							"body": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"headers": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"querystring": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						}
					],
					"required": true,
					"type": "record"
				}
			},
			{
				"append": {
					"fields": [
						{
							"body": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"headers": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						},
						{
							"querystring": {
								"default": [],
								"elements": {
									"match": "^[^:]+:.*$",
									"type": "string"
								},
								"type": "array"
							}
						}
					],
					"required": true,
					"type": "record"
				}
			}
		]
	}`
	RequestTransformerConfig = `{
		"add": {
			"body": [],
			"headers": [],
			"querystring": []
		},
		"append": {
			"body": [],
			"headers": [],
			"querystring": []
		},
		"http_method": null,
		"remove": {
			"body": [],
			"headers": [],
			"querystring": []
		},
		"rename": {
			"body": [],
			"headers": [],
			"querystring": []
		},
		"replace": {
			"body": [],
			"headers": [],
			"querystring": []
		}
	}`
	RequestTransformerNonEmptyConfig = `{
		"remove": {
			"headers": [ "cookie" ],
			"body": [ "foo" ]
		},
		"add": {
			"body": [ "bar" ]
		}
	}`
	RequestTransformerNonEmptyFilledConfig = `{
		"add": {
			"body": [ "bar" ],
			"headers": [],
			"querystring": []
		},
		"append": {
			"body": [],
			"headers": [],
			"querystring": []
		},
		"http_method": null,
		"remove": {
			"body": [ "foo" ],
			"headers": [ "cookie" ],
			"querystring": []
		},
		"rename": {
			"body": [],
			"headers": [],
			"querystring": []
		},
		"replace": {
			"body": [],
			"headers": [],
			"querystring": []
		}
	}`
)

func TestFillNil(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(fill(nil, nil))
}

func TestFillKeyAuth(t *testing.T) {
	assert := assert.New(t)

	var schema map[string]interface{}
	err := json.Unmarshal([]byte(KeyAuthSchema), &schema)
	assert.Nil(err)

	config := make(kong.Configuration)

	def := make(kong.Configuration)
	err = json.Unmarshal([]byte(KeyAuthDefaultConfig), &def)
	assert.Nil(err)

	res, err := fill(schema, config)
	assert.Equal(def, res)
}

func TestFillStatsD(t *testing.T) {
	assert := assert.New(t)

	var schema map[string]interface{}
	err := json.Unmarshal([]byte(StatsDSchema), &schema)
	assert.Nil(err)

	config := make(kong.Configuration)

	def := make(kong.Configuration)
	err = json.Unmarshal([]byte(StatsDDefaultConfig), &def)
	assert.Nil(err)

	res, err := fill(schema, config)
	assert.Equal(def, res)
}

func TestKeyAuthSetKeys(t *testing.T) {
	assert := assert.New(t)

	var schema map[string]interface{}
	err := json.Unmarshal([]byte(KeyAuthSchema), &schema)
	assert.Nil(err)

	config := make(kong.Configuration)
	config["key_in_body"] = true

	def := make(kong.Configuration)
	err = json.Unmarshal([]byte(KeyAuthDefaultConfig), &def)
	assert.Nil(err)

	res, err := fill(schema, config)
	assert.NotEqual(def, res)
	assert.Equal(true, res["key_in_body"])
}

func TestFillReqeustTransformer(t *testing.T) {
	assert := assert.New(t)
	var schema map[string]interface{}
	err := json.Unmarshal([]byte(RequestTransformerSchema), &schema)
	assert.Nil(err)

	config := make(kong.Configuration)
	def := make(kong.Configuration)
	err = json.Unmarshal([]byte(RequestTransformerConfig), &def)
	assert.Nil(err)

	res, err := fill(schema, config)
	assert.Equal(def, res)
}

func TestFillReqeustTransformerNestedConfig(t *testing.T) {
	assert := assert.New(t)
	var schema map[string]interface{}
	err := json.Unmarshal([]byte(RequestTransformerSchema), &schema)
	assert.Nil(err)

	config := make(kong.Configuration)
	err = json.Unmarshal([]byte(RequestTransformerNonEmptyConfig), &config)
	assert.Nil(err)
	want := make(kong.Configuration)
	err = json.Unmarshal([]byte(RequestTransformerNonEmptyFilledConfig), &want)
	res, err := fill(schema, config)
	assert.Equal(want, res)
	assert.Nil(err)
}
