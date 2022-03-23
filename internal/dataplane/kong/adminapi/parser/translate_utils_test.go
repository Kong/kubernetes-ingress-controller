package parser

import (
	"fmt"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func Test_convertGatewayMatchHeadersToKongRouteMatchHeadersVersionBehavior(t *testing.T) {
	regexType := gatewayv1alpha2.HeaderMatchRegularExpression

	type Case struct {
		msg    string
		input  []gatewayv1alpha2.HTTPHeaderMatch
		output map[string][]string
		err    error
	}

	// util reports Kong version 0.0.0 when not initialized
	belowThresholdTests := []Case{
		{
			msg: "regex header matches fail on unsupported versions",
			input: []gatewayv1alpha2.HTTPHeaderMatch{{
				Type:  &regexType,
				Name:  "Content-Type",
				Value: "^audio/*",
			}},
			output: nil,
			err:    fmt.Errorf("Kong version %s does not support HeaderMatchRegularExpression", semver.MustParse("0.0.0")),
		},
		{
			msg: "a single exact header match succeeds on any version",
			input: []gatewayv1alpha2.HTTPHeaderMatch{{
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
	}

	for _, tt := range belowThresholdTests {
		t.Run(tt.msg, func(t *testing.T) {
			output, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(tt.input)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.output, output)
		})
	}

	util.SetKongVersion(semver.MustParse("2.8.0"))

	aboveThresholdTests := []Case{
		{
			msg: "regex header matches succeed on supported versions",
			input: []gatewayv1alpha2.HTTPHeaderMatch{{
				Type:  &regexType,
				Name:  "Content-Type",
				Value: "^audio/*",
			}},
			output: map[string][]string{
				"Content-Type": {kongHeaderRegexPrefix + "^audio/*"},
			},
		},
	}

	for _, tt := range aboveThresholdTests {
		t.Run(tt.msg, func(t *testing.T) {
			output, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(tt.input)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.output, output)
		})
	}
}

func Test_convertGatewayMatchHeadersToKongRouteMatchHeaders(t *testing.T) {
	regexType := gatewayv1alpha2.HeaderMatchRegularExpression
	exactType := gatewayv1alpha2.HeaderMatchExact
	util.SetKongVersion(semver.MustParse("2.8.0"))

	t.Log("generating several gateway header matches")
	tests := []struct {
		msg    string
		input  []gatewayv1alpha2.HTTPHeaderMatch
		output map[string][]string
		err    error
	}{
		{
			msg: "regex header matches convert correctly",
			input: []gatewayv1alpha2.HTTPHeaderMatch{{
				Type:  &regexType,
				Name:  "Content-Type",
				Value: "^audio/*",
			}},
			output: map[string][]string{
				"Content-Type": {kongHeaderRegexPrefix + "^audio/*"},
			},
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
			msg: "multiple header matches for the same header are rejected",
			input: []gatewayv1alpha2.HTTPHeaderMatch{
				{
					Name:  "Content-Type",
					Value: "audio/vorbis",
				},
				{
					Name:  "Content-Type",
					Value: "audio/flac",
				},
			},
			output: nil,
			err:    fmt.Errorf("multiple header matches for the same header are not allowed: Content-Type"),
		},
		{
			msg: "multiple header matches convert properly",
			input: []gatewayv1alpha2.HTTPHeaderMatch{
				{
					Type:  &exactType,
					Name:  "Content-Type",
					Value: "audio/vorbis",
				},
				{
					Name:  "Content-Length",
					Value: "999999999",
				},
			},
			output: map[string][]string{
				"Content-Type":   {"audio/vorbis"},
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
