package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"
)

func TestExtractEntityFieldDefinitions(t *testing.T) {
	testCases := []struct {
		name           string
		schema         kong.Schema
		expectedFields map[string]EntityField
	}{
		{
			name: "absent fields should have a nil value",
			schema: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"foo": map[string]interface{}{
							"type":     "string",
							"required": true,
						},
					},
					map[string]interface{}{
						"bar": map[string]interface{}{
							"type":      "foreign",
							"required":  true,
							"reference": "service",
						},
					},
				},
			},
			expectedFields: map[string]EntityField{
				"foo": {
					Name:     "foo",
					Type:     EntityFieldTypeString,
					Required: true,
					Auto:     false,
					UUID:     false,
				},
				"bar": {
					Name:      "bar",
					Type:      EntityFieldTypeForeign,
					Required:  true,
					Reference: "service",
				},
			},
		},
		{
			name: "irrelevant fields should be safely ignored",
			schema: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"protocol": map[string]interface{}{
							"type":     "string",
							"required": true,
							"default":  "http",
							"one_of":   []string{"http", "https"},
						},
						"port": map[string]interface{}{
							"type":     "integer",
							"required": true,
							"default":  80,
							"min":      1,
							"max":      65535,
						},
					},
				},
				"checks": "some_check",
			},
			expectedFields: map[string]EntityField{
				"protocol": {
					Name:     "protocol",
					Type:     EntityFieldTypeString,
					Required: true,
					Default:  "http",
				},
				"port": {
					Name:     "port",
					Type:     EntityFieldTypeInteger,
					Required: true,
					Default:  int(80),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := ExtractEntityFieldDefinitions(tc.schema).Fields
			for fieldName, expectedField := range tc.expectedFields {
				actualField, ok := fields[fieldName]
				require.Truef(t, ok, "field %s should exist", fieldName)
				require.Equalf(t, expectedField, actualField, "field %s should be same as expected", fieldName)
			}
		})
	}
}
