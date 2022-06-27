package parser

import (
	"fmt"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
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
			result := isRefAllowedByPolicy(tt.ref.Namespace, tt.ref.Name, tt.ref.Group, tt.ref.Kind, fakeMap)
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

func Test_generateKongServiceFromBackendRef(t *testing.T) {
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
						Kind:  gatewayv1alpha2.Kind("Service"),
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
						Kind:  gatewayv1alpha2.Kind("Service"),
					},
				},
			},
		},
	}
	fakestore, err := store.NewFakeStore(store.FakeObjects{ReferencePolicies: policies})
	assert.Nil(t, err)
	p := NewParser(logrus.New(), fakestore)
	// empty since we always want to actually generate a service for tests
	// static values for the basic string format inputs since nothing interesting happens with them
	rules := ingressRules{ServiceNameToServices: map[string]kongstate.Service{}}
	ruleNumber := 999
	protocol := "example"
	port := gatewayv1alpha2.PortNumber(7777)
	redObjName := gatewayv1alpha2.ObjectName("red-service")
	blueObjName := gatewayv1alpha2.ObjectName("blue-service")
	cholponNamespace := gatewayv1alpha2.Namespace("cholpon")
	serviceKind := gatewayv1alpha2.Kind("Service")
	serviceGroup := gatewayv1alpha2.Group("")
	tests := []struct {
		msg     string
		route   client.Object
		refs    []gatewayv1alpha2.BackendRef
		result  kongstate.Service
		wantErr bool
	}{
		{
			msg:     "empty backend list",
			route:   &gatewayv1alpha2.HTTPRoute{},
			refs:    []gatewayv1alpha2.BackendRef{},
			result:  kongstate.Service{},
			wantErr: true,
		},
		{
			msg: "all backends in route namespace",
			route: &gatewayv1alpha2.HTTPRoute{
				// normally the k8s api call populates TypeMeta properly, but we have no such luxuries here
				TypeMeta: metav1.TypeMeta{
					Kind:       "HTTPRoute",
					APIVersion: "gateway.networking.k8s.io/v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tong-sirlari",
					Namespace: "cholpon",
				},
			},
			refs: []gatewayv1alpha2.BackendRef{
				{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name:  blueObjName,
						Kind:  &serviceKind,
						Port:  &port,
						Group: &serviceGroup,
					},
				},
				{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name:  redObjName,
						Kind:  &serviceKind,
						Port:  &port,
						Group: &serviceGroup,
					},
				},
			},
			result: kongstate.Service{
				Service: kong.Service{
					Name:           kong.String("httproute.cholpon.tong-sirlari.999"),
					Host:           kong.String("httproute.cholpon.tong-sirlari.999"),
					Protocol:       kong.String(protocol),
					ConnectTimeout: kong.Int(DefaultServiceTimeout),
					ReadTimeout:    kong.Int(DefaultServiceTimeout),
					WriteTimeout:   kong.Int(DefaultServiceTimeout),
					Retries:        kong.Int(DefaultRetries),
				},
				Namespace: "cholpon",
				Backends: []kongstate.ServiceBackend{
					{
						Name: string(blueObjName),
						PortDef: kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: int32(port),
						},
					},
					{
						Name: string(redObjName),
						PortDef: kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: int32(port),
						},
					},
				},
				Parent: &gatewayv1alpha2.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tong-sirlari",
						Namespace: "cholpon",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "HTTPRoute",
						APIVersion: "gateway.networking.k8s.io/v1alpha2",
					},
				},
			},
			wantErr: false,
		},
		{
			msg: "same and different ns backend",
			route: &gatewayv1alpha2.UDPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "UDPRoute",
					APIVersion: "gateway.networking.k8s.io/v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "padarkush",
					Namespace: "behbudiy",
				},
			},
			refs: []gatewayv1alpha2.BackendRef{
				{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name:      blueObjName,
						Port:      &port,
						Kind:      &serviceKind,
						Namespace: &cholponNamespace,
						Group:     &serviceGroup,
					},
				},
				{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name:  redObjName,
						Port:  &port,
						Kind:  &serviceKind,
						Group: &serviceGroup,
					},
				},
			},
			result: kongstate.Service{
				Service: kong.Service{
					Name:           kong.String("udproute.behbudiy.padarkush.999"),
					Host:           kong.String("udproute.behbudiy.padarkush.999"),
					Protocol:       kong.String(protocol),
					ConnectTimeout: kong.Int(DefaultServiceTimeout),
					ReadTimeout:    kong.Int(DefaultServiceTimeout),
					WriteTimeout:   kong.Int(DefaultServiceTimeout),
					Retries:        kong.Int(DefaultRetries),
				},
				Namespace: "behbudiy",
				Backends: []kongstate.ServiceBackend{
					{
						Name:      string(blueObjName),
						Namespace: "cholpon",
						PortDef: kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: int32(port),
						},
					},
					{
						Name: string(redObjName),
						PortDef: kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: int32(port),
						},
					},
				},
				Parent: &gatewayv1alpha2.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "padarkush",
						Namespace: "behbudiy",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "UDPRoute",
						APIVersion: "gateway.networking.k8s.io/v1alpha2",
					},
				},
			},
			wantErr: false,
		},
		{
			msg: "only not permitted remote ns",
			route: &gatewayv1alpha2.TCPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "TCPRoute",
					APIVersion: "gateway.networking.k8s.io/v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kitab-ul-atfol",
					Namespace: "behbudiy",
				},
			},
			refs: []gatewayv1alpha2.BackendRef{
				{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name:      blueObjName,
						Port:      &port,
						Kind:      &serviceKind,
						Namespace: &cholponNamespace,
						Group:     &serviceGroup,
					},
				},
			},
			result:  kongstate.Service{},
			wantErr: true,
		},
		{
			msg: "same and different ns backend",
			route: &gatewayv1alpha2.TCPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "TCPRoute",
					APIVersion: "gateway.networking.k8s.io/v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "muntaxabi-jugrofiyai-umumiy",
					Namespace: "behbudiy",
				},
			},
			refs: []gatewayv1alpha2.BackendRef{
				{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name:      blueObjName,
						Port:      &port,
						Kind:      &serviceKind,
						Namespace: &cholponNamespace,
						Group:     &serviceGroup,
					},
				},
				{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name:  redObjName,
						Port:  &port,
						Kind:  &serviceKind,
						Group: &serviceGroup,
					},
				},
			},
			result: kongstate.Service{
				Service: kong.Service{
					Name:           kong.String("tcproute.behbudiy.muntaxabi-jugrofiyai-umumiy.999"),
					Host:           kong.String("tcproute.behbudiy.muntaxabi-jugrofiyai-umumiy.999"),
					Protocol:       kong.String(protocol),
					ConnectTimeout: kong.Int(DefaultServiceTimeout),
					ReadTimeout:    kong.Int(DefaultServiceTimeout),
					WriteTimeout:   kong.Int(DefaultServiceTimeout),
					Retries:        kong.Int(DefaultRetries),
				},
				Namespace: "behbudiy",
				Backends: []kongstate.ServiceBackend{
					{
						Name: string(redObjName),
						PortDef: kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: int32(port),
						},
					},
				},
				Parent: &gatewayv1alpha2.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "muntaxabi-jugrofiyai-umumiy",
						Namespace: "behbudiy",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "TCPRoute",
						APIVersion: "gateway.networking.k8s.io/v1alpha2",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			result, err := p.generateKongServiceFromBackendRef(&rules, tt.route, ruleNumber, protocol, tt.refs...)
			assert.Equal(t, tt.result, result)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
