package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

func Test_convertGatewayMatchHeadersToKongRouteMatchHeaders(t *testing.T) {
	regexType := gatewayv1alpha2.HeaderMatchRegularExpression
	exactType := gatewayv1alpha2.HeaderMatchExact

	t.Log("generating several gateway header matches")
	tests := []struct {
		msg    string
		input  []gatewayv1alpha2.HTTPHeaderMatch
		output map[string][]string
		err    error
	}{
		{
			msg: "regex header matches will fail due to lack of support",
			input: []gatewayv1alpha2.HTTPHeaderMatch{{
				Type:  &regexType,
				Name:  "Content-Type",
				Value: "^audio/*",
			}},
			output: nil,
			err:    fmt.Errorf("regular expression header matches are not yet supported"),
		},
		{
			msg: "a single exact header match with no type defaults to exact type and converts properly",
			input: []gatewayv1alpha2.HTTPHeaderMatch{{
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
		{
			msg: "a single exact header match with a single value converts properly",
			input: []gatewayv1alpha2.HTTPHeaderMatch{{
				Type:  &exactType,
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
		{
			msg: "a single exact header match with multiple values converts properly",
			input: []gatewayv1alpha2.HTTPHeaderMatch{{
				Type:  &exactType,
				Name:  "Content-Type",
				Value: "audio/vorbis,audio/mpeg",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis", "audio/mpeg"},
			},
		},
		{
			msg: "multiple header matches with a mixture of value counts convert properly",
			input: []gatewayv1alpha2.HTTPHeaderMatch{
				{
					Type:  &exactType,
					Name:  "Content-Type",
					Value: "audio/vorbis,audio/mpeg",
				},
				{
					Name:  "Content-Length",
					Value: "999999999",
				},
			},
			output: map[string][]string{
				"Content-Type":   {"audio/vorbis", "audio/mpeg"},
				"Content-Length": {"999999999"},
			},
		},
		{
			msg:    "an empty list of headers will produce no converted headers",
			output: map[string][]string{},
		},
	}

	t.Log("verifying header match conversions")
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			output, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(tt.input)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.output, output)
		})
	}
}
