package parser

import (
	"fmt"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

func TestConvertGatewayMatchHeadersToKongRouteMatchHeadersVersionBehavior(t *testing.T) {
	type Case struct {
		msg         string
		input       []gatewayapi.HTTPHeaderMatch
		kongVersion func(t *testing.T) semver.Version
		output      map[string][]string
		err         error
	}

	testcases := []Case{
		{
			msg: "regex header matches fail on unsupported versions",
			input: []gatewayapi.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayapi.HeaderMatchRegularExpression),
				Name:  "Content-Type",
				Value: "^audio/*",
			}},
			kongVersion: func(t *testing.T) semver.Version {
				kongVersion, err := semver.Parse("2.5.0")
				require.NoError(t, err)
				return kongVersion
			},
			output: nil,
			err:    fmt.Errorf("Kong version %s does not support HeaderMatchRegularExpression", semver.MustParse("2.5.0")),
		},
		{
			msg: "a single exact header match succeeds on any version",
			input: []gatewayapi.HTTPHeaderMatch{{
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			kongVersion: func(t *testing.T) semver.Version {
				kongVersion, err := semver.Parse("2.5.0")
				require.NoError(t, err)
				return kongVersion
			},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
		{
			msg: "regex header matches succeed on supported versions",
			input: []gatewayapi.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayapi.HeaderMatchRegularExpression),
				Name:  "Content-Type",
				Value: "^audio/*",
			}},
			kongVersion: func(t *testing.T) semver.Version {
				kongVersion, err := semver.Parse("2.8.0")
				require.NoError(t, err)
				return kongVersion
			},
			output: map[string][]string{
				"Content-Type": {kongHeaderRegexPrefix + "^audio/*"},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.msg, func(t *testing.T) {
			output, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(tt.input, tt.kongVersion(t))
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.output, output)
		})
	}
}

func TestConvertGatewayMatchHeadersToKongRouteMatchHeaders(t *testing.T) {
	kongVersion, err := semver.Parse("2.8.0")
	require.NoError(t, err)

	t.Log("generating several gateway header matches")
	tests := []struct {
		msg    string
		input  []gatewayapi.HTTPHeaderMatch
		output map[string][]string
		err    error
	}{
		{
			msg: "regex header matches convert correctly",
			input: []gatewayapi.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayapi.HeaderMatchRegularExpression),
				Name:  "Content-Type",
				Value: "^audio/*",
			}},
			output: map[string][]string{
				"Content-Type": {kongHeaderRegexPrefix + "^audio/*"},
			},
		},
		{
			msg: "a single exact header match with no type defaults to exact type and converts properly",
			input: []gatewayapi.HTTPHeaderMatch{{
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
		{
			msg: "a single exact header match with a single value converts properly",
			input: []gatewayapi.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayapi.HeaderMatchExact),
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
		{
			msg: "multiple header matches for the same header are rejected",
			input: []gatewayapi.HTTPHeaderMatch{
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
			input: []gatewayapi.HTTPHeaderMatch{
				{
					Type:  lo.ToPtr(gatewayapi.HeaderMatchExact),
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
			output, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(tt.input, kongVersion)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.output, output)
		})
	}
}

func TestGetPermittedForReferenceGrantFrom(t *testing.T) {
	grants := []*gatewayapi.ReferenceGrant{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "fitrat",
			},
			Spec: gatewayapi.ReferenceGrantSpec{
				From: []gatewayapi.ReferenceGrantFrom{
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("TCPRoute"),
						Namespace: gatewayapi.Namespace("garbage"),
					},
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("TCPRoute"),
						Namespace: gatewayapi.Namespace("behbudiy"),
					},
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("TCPRoute"),
						Namespace: gatewayapi.Namespace("qodiriy"),
					},
				},
				To: []gatewayapi.ReferenceGrantTo{
					{
						Group: gatewayapi.Group(""),
						Kind:  gatewayapi.Kind("GrantOne"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "cholpon",
			},
			Spec: gatewayapi.ReferenceGrantSpec{
				From: []gatewayapi.ReferenceGrantFrom{
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("UDPRoute"),
						Namespace: gatewayapi.Namespace("behbudiy"),
					},
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("TCPRoute"),
						Namespace: gatewayapi.Namespace("qodiriy"),
					},
				},
				To: []gatewayapi.ReferenceGrantTo{
					{
						Group: gatewayapi.Group(""),
						Kind:  gatewayapi.Kind("GrantTwo"),
					},
				},
			},
		},
	}
	tests := []struct {
		msg    string
		from   gatewayapi.ReferenceGrantFrom
		result map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo
	}{
		{
			msg: "no matches whatsoever",
			from: gatewayapi.ReferenceGrantFrom{
				Group:     gatewayapi.Group("invalid.example"),
				Kind:      gatewayapi.Kind("invalid"),
				Namespace: gatewayapi.Namespace("invalid"),
			},
			result: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{},
		},
		{
			msg: "non-matching namespace",
			from: gatewayapi.ReferenceGrantFrom{
				Group:     gatewayapi.Group("gateway.networking.k8s.io"),
				Kind:      gatewayapi.Kind("UDPRoute"),
				Namespace: gatewayapi.Namespace("niyazi"),
			},
			result: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{},
		},
		{
			msg: "non-matching kind",
			from: gatewayapi.ReferenceGrantFrom{
				Group:     gatewayapi.Group("gateway.networking.k8s.io"),
				Kind:      gatewayapi.Kind("TLSRoute"),
				Namespace: gatewayapi.Namespace("behbudiy"),
			},
			result: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{},
		},
		{
			msg: "non-matching group",
			from: gatewayapi.ReferenceGrantFrom{
				Group:     gatewayapi.Group("invalid.example"),
				Kind:      gatewayapi.Kind("UDPRoute"),
				Namespace: gatewayapi.Namespace("behbudiy"),
			},
			result: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{},
		},
		{
			msg: "single match",
			from: gatewayapi.ReferenceGrantFrom{
				Group:     gatewayapi.Group("gateway.networking.k8s.io"),
				Kind:      gatewayapi.Kind("UDPRoute"),
				Namespace: gatewayapi.Namespace("behbudiy"),
			},
			result: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{
				"cholpon": {
					{
						Group: gatewayapi.Group(""),
						Kind:  gatewayapi.Kind("GrantTwo"),
					},
				},
			},
		},
		{
			msg: "multiple matches",
			from: gatewayapi.ReferenceGrantFrom{
				Group:     gatewayapi.Group("gateway.networking.k8s.io"),
				Kind:      gatewayapi.Kind("TCPRoute"),
				Namespace: gatewayapi.Namespace("qodiriy"),
			},
			result: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{
				"cholpon": {
					{
						Group: gatewayapi.Group(""),
						Kind:  gatewayapi.Kind("GrantTwo"),
					},
				},
				"fitrat": {
					{
						Group: gatewayapi.Group(""),
						Kind:  gatewayapi.Kind("GrantOne"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			result := getPermittedForReferenceGrantFrom(tt.from, grants)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestGenerateKongServiceFromBackendRef(t *testing.T) {
	grants := []*gatewayapi.ReferenceGrant{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "fitrat",
			},
			Spec: gatewayapi.ReferenceGrantSpec{
				From: []gatewayapi.ReferenceGrantFrom{
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("TCPRoute"),
						Namespace: gatewayapi.Namespace("garbage"),
					},
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("TCPRoute"),
						Namespace: gatewayapi.Namespace("behbudiy"),
					},
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("TCPRoute"),
						Namespace: gatewayapi.Namespace("qodiriy"),
					},
				},
				To: []gatewayapi.ReferenceGrantTo{
					{
						Group: gatewayapi.Group(""),
						Kind:  gatewayapi.Kind("Service"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "cholpon",
			},
			Spec: gatewayapi.ReferenceGrantSpec{
				From: []gatewayapi.ReferenceGrantFrom{
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("UDPRoute"),
						Namespace: gatewayapi.Namespace("behbudiy"),
					},
					{
						Group:     gatewayapi.Group("gateway.networking.k8s.io"),
						Kind:      gatewayapi.Kind("TCPRoute"),
						Namespace: gatewayapi.Namespace("qodiriy"),
					},
				},
				To: []gatewayapi.ReferenceGrantTo{
					{
						Group: gatewayapi.Group(""),
						Kind:  gatewayapi.Kind("Service"),
					},
				},
			},
		},
	}
	fakestore, err := store.NewFakeStore(store.FakeObjects{ReferenceGrants: grants})
	assert.Nil(t, err)
	p := mustNewParser(t, fakestore)
	// empty since we always want to actually generate a service for tests
	// static values for the basic string format inputs since nothing interesting happens with them
	rules := ingressRules{ServiceNameToServices: map[string]kongstate.Service{}}
	ruleNumber := 999
	protocol := "example"
	port := gatewayapi.PortNumber(7777)
	redObjName := gatewayapi.ObjectName("red-service")
	blueObjName := gatewayapi.ObjectName("blue-service")
	cholponNamespace := gatewayapi.Namespace("cholpon")
	serviceKind := gatewayapi.Kind("Service")
	serviceGroup := gatewayapi.Group("")
	tests := []struct {
		msg     string
		route   client.Object
		refs    []gatewayapi.BackendRef
		result  kongstate.Service
		wantErr bool
	}{
		{
			msg: "all backends in route namespace",
			route: &gatewayapi.HTTPRoute{
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
			refs: []gatewayapi.BackendRef{
				{
					BackendObjectReference: gatewayapi.BackendObjectReference{
						Name:  blueObjName,
						Kind:  &serviceKind,
						Port:  &port,
						Group: &serviceGroup,
					},
				},
				{
					BackendObjectReference: gatewayapi.BackendObjectReference{
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
				Parent: &gatewayapi.HTTPRoute{
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
			route: &gatewayapi.UDPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "UDPRoute",
					APIVersion: "gateway.networking.k8s.io/v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "padarkush",
					Namespace: "behbudiy",
				},
			},
			refs: []gatewayapi.BackendRef{
				{
					BackendObjectReference: gatewayapi.BackendObjectReference{
						Name:      blueObjName,
						Port:      &port,
						Kind:      &serviceKind,
						Namespace: &cholponNamespace,
						Group:     &serviceGroup,
					},
				},
				{
					BackendObjectReference: gatewayapi.BackendObjectReference{
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
				Parent: &gatewayapi.UDPRoute{
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
			route: &gatewayapi.TCPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "TCPRoute",
					APIVersion: "gateway.networking.k8s.io/v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kitab-ul-atfol",
					Namespace: "behbudiy",
				},
			},
			refs: []gatewayapi.BackendRef{
				{
					BackendObjectReference: gatewayapi.BackendObjectReference{
						Name:      blueObjName,
						Port:      &port,
						Kind:      &serviceKind,
						Namespace: &cholponNamespace,
						Group:     &serviceGroup,
					},
				},
			},
			result: kongstate.Service{
				Service: kong.Service{
					Name:           kong.String("tcproute.behbudiy.kitab-ul-atfol.999"),
					Host:           kong.String("tcproute.behbudiy.kitab-ul-atfol.999"),
					Protocol:       kong.String(protocol),
					ConnectTimeout: kong.Int(DefaultServiceTimeout),
					ReadTimeout:    kong.Int(DefaultServiceTimeout),
					WriteTimeout:   kong.Int(DefaultServiceTimeout),
					Retries:        kong.Int(DefaultRetries),
				},
				Namespace: "behbudiy",
				Backends:  []kongstate.ServiceBackend{},
				Plugins: []kong.Plugin{
					{
						Name: kong.String("request-termination"),
						Config: kong.Configuration{
							"status_code": 500,
							"message":     "no existing backendRef provided",
						},
					},
				},
				Parent: &gatewayapi.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kitab-ul-atfol",
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
		{
			msg: "same and different ns backend",
			route: &gatewayapi.TCPRoute{
				TypeMeta: metav1.TypeMeta{
					Kind:       "TCPRoute",
					APIVersion: "gateway.networking.k8s.io/v1alpha2",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "muntaxabi-jugrofiyai-umumiy",
					Namespace: "behbudiy",
				},
			},
			refs: []gatewayapi.BackendRef{
				{
					BackendObjectReference: gatewayapi.BackendObjectReference{
						Name:      blueObjName,
						Port:      &port,
						Kind:      &serviceKind,
						Namespace: &cholponNamespace,
						Group:     &serviceGroup,
					},
				},
				{
					BackendObjectReference: gatewayapi.BackendObjectReference{
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
				Parent: &gatewayapi.TCPRoute{
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
			result, err := generateKongServiceFromBackendRefWithRuleNumber(p.logger, p.storer, &rules, tt.route, ruleNumber, protocol, tt.refs...)
			assert.Equal(t, tt.result, result)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
