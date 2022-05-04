package parser

import (
	"fmt"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func Test_isRefAllowedByPolicy(t *testing.T) {
	fitrat := gatewayv1alpha2.Namespace("fitrat")
	cholpon := gatewayv1alpha2.Namespace("cholpon")
	behbudiy := gatewayv1alpha2.Namespace("behbudiy")

	group := gatewayv1alpha2.Group("fake.example.com")
	kind := gatewayv1alpha2.Kind("fakeKind")
	badKind := gatewayv1alpha2.Kind("badFakeKind")
	cholponName := gatewayv1alpha2.ObjectName("cholpon")

	fakeMap := map[gatewayv1alpha2.Namespace][]gatewayv1alpha2.ReferencePolicyTo{
		fitrat:   {{Group: group, Kind: kind}, {Group: gatewayv1alpha2.Group("extra.example"), Kind: badKind}},
		cholpon:  {{Group: group, Kind: kind, Name: &cholponName}},
		behbudiy: {},
	}
	tests := []struct {
		msg    string
		ref    gatewayv1alpha2.BackendRef
		result bool
	}{
		{
			msg:    "empty",
			ref:    gatewayv1alpha2.BackendRef{},
			result: true,
		},
		{
			msg: "no namespace",
			ref: gatewayv1alpha2.BackendRef{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name: gatewayv1alpha2.ObjectName("foo"),
				},
			},
			result: true,
		},
		{
			msg: "valid namespace+group+kind",
			ref: gatewayv1alpha2.BackendRef{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name:      gatewayv1alpha2.ObjectName("foo"),
					Group:     &group,
					Kind:      &kind,
					Namespace: &fitrat,
				},
			},
			result: true,
		},
		{
			msg: "valid namespace+group+kind+name",
			ref: gatewayv1alpha2.BackendRef{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name:      cholponName,
					Group:     &group,
					Kind:      &kind,
					Namespace: &cholpon,
				},
			},
			result: true,
		},
		{
			msg: "invalid namespace+group+kind",
			ref: gatewayv1alpha2.BackendRef{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name:      gatewayv1alpha2.ObjectName("foo"),
					Group:     &group,
					Kind:      &badKind,
					Namespace: &fitrat,
				},
			},
			result: false,
		},
		{
			msg: "invalid namespace+group+kind+name",
			ref: gatewayv1alpha2.BackendRef{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name:      gatewayv1alpha2.ObjectName("sadness"),
					Group:     &group,
					Kind:      &kind,
					Namespace: &cholpon,
				},
			},
			result: false,
		},
		{
			msg: "no policies in target namespace",
			ref: gatewayv1alpha2.BackendRef{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name:      gatewayv1alpha2.ObjectName("foo"),
					Group:     &group,
					Kind:      &kind,
					Namespace: &behbudiy,
				},
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			result := isRefAllowedByPolicy(tt.ref, fakeMap)
			assert.Equal(t, tt.result, result)
		})
	}
}

func Test_getPermittedForReferencePolicyFrom(t *testing.T) {
	policies := []*gatewayv1alpha2.ReferencePolicy{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:        uuid.NewString(),
				Annotations: map[string]string{},
				Namespace:   "fitrat",
			},
			Spec: gatewayv1alpha2.ReferencePolicySpec{
				From: []gatewayv1alpha2.ReferencePolicyFrom{
					{
						Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
						Kind:      gatewayv1alpha2.Kind("TCPRoute"),
						Namespace: gatewayv1alpha2.Namespace("garbage"),
					},
					{
						Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
						Kind:      gatewayv1alpha2.Kind("TCPRoute"),
						Namespace: gatewayv1alpha2.Namespace("behbudiy"),
					},
					{
						Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
						Kind:      gatewayv1alpha2.Kind("TCPRoute"),
						Namespace: gatewayv1alpha2.Namespace("qodiriy"),
					},
				},
				To: []gatewayv1alpha2.ReferencePolicyTo{
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("PolicyOne"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:        uuid.NewString(),
				Annotations: map[string]string{},
				Namespace:   "cholpon",
			},
			Spec: gatewayv1alpha2.ReferencePolicySpec{
				From: []gatewayv1alpha2.ReferencePolicyFrom{
					{
						Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
						Kind:      gatewayv1alpha2.Kind("UDPRoute"),
						Namespace: gatewayv1alpha2.Namespace("behbudiy"),
					},
					{
						Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
						Kind:      gatewayv1alpha2.Kind("TCPRoute"),
						Namespace: gatewayv1alpha2.Namespace("qodiriy"),
					},
				},
				To: []gatewayv1alpha2.ReferencePolicyTo{
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("PolicyTwo"),
					},
				},
			},
		},
	}
	tests := []struct {
		msg    string
		from   gatewayv1alpha2.ReferencePolicyFrom
		result map[gatewayv1alpha2.Namespace][]gatewayv1alpha2.ReferencePolicyTo
	}{
		{
			msg: "no matches whatsoever",
			from: gatewayv1alpha2.ReferencePolicyFrom{
				Group:     gatewayv1alpha2.Group("invalid.example"),
				Kind:      gatewayv1alpha2.Kind("invalid"),
				Namespace: gatewayv1alpha2.Namespace("invalid"),
			},
			result: map[gatewayv1alpha2.Namespace][]gatewayv1alpha2.ReferencePolicyTo{},
		},
		{
			msg: "non-matching namespace",
			from: gatewayv1alpha2.ReferencePolicyFrom{
				Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
				Kind:      gatewayv1alpha2.Kind("UDPRoute"),
				Namespace: gatewayv1alpha2.Namespace("niyazi"),
			},
			result: map[gatewayv1alpha2.Namespace][]gatewayv1alpha2.ReferencePolicyTo{},
		},
		{
			msg: "non-matching kind",
			from: gatewayv1alpha2.ReferencePolicyFrom{
				Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
				Kind:      gatewayv1alpha2.Kind("TLSRoute"),
				Namespace: gatewayv1alpha2.Namespace("behbudiy"),
			},
			result: map[gatewayv1alpha2.Namespace][]gatewayv1alpha2.ReferencePolicyTo{},
		},
		{
			msg: "non-matching group",
			from: gatewayv1alpha2.ReferencePolicyFrom{
				Group:     gatewayv1alpha2.Group("invalid.example"),
				Kind:      gatewayv1alpha2.Kind("UDPRoute"),
				Namespace: gatewayv1alpha2.Namespace("behbudiy"),
			},
			result: map[gatewayv1alpha2.Namespace][]gatewayv1alpha2.ReferencePolicyTo{},
		},
		{
			msg: "single match",
			from: gatewayv1alpha2.ReferencePolicyFrom{
				Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
				Kind:      gatewayv1alpha2.Kind("UDPRoute"),
				Namespace: gatewayv1alpha2.Namespace("behbudiy"),
			},
			result: map[gatewayv1alpha2.Namespace][]gatewayv1alpha2.ReferencePolicyTo{
				"cholpon": {
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("PolicyTwo"),
					},
				},
			},
		},
		{
			msg: "multiple matches",
			from: gatewayv1alpha2.ReferencePolicyFrom{
				Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
				Kind:      gatewayv1alpha2.Kind("TCPRoute"),
				Namespace: gatewayv1alpha2.Namespace("qodiriy"),
			},
			result: map[gatewayv1alpha2.Namespace][]gatewayv1alpha2.ReferencePolicyTo{
				"cholpon": {
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("PolicyTwo"),
					},
				},
				"fitrat": {
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("PolicyOne"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			result := getPermittedForReferencePolicyFrom(tt.from, policies)
			assert.Equal(t, tt.result, result)
		})
	}
}
