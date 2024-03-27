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
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

func TestConvertGatewayMatchHeadersToKongRouteMatchHeadersVersionBehavior(t *testing.T) {
	type Case struct {
		msg         string
		input       []gatewayv1.HTTPHeaderMatch
		kongVersion func(t *testing.T) semver.Version
		output      map[string][]string
		err         error
	}

	testcases := []Case{
		{
			msg: "regex header matches fail on unsupported versions",
			input: []gatewayv1.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayv1.HeaderMatchRegularExpression),
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
			input: []gatewayv1.HTTPHeaderMatch{{
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
			input: []gatewayv1.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayv1.HeaderMatchRegularExpression),
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
		input  []gatewayv1.HTTPHeaderMatch
		output map[string][]string
		err    error
	}{
		{
			msg: "regex header matches convert correctly",
			input: []gatewayv1.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayv1.HeaderMatchRegularExpression),
				Name:  "Content-Type",
				Value: "^audio/*",
			}},
			output: map[string][]string{
				"Content-Type": {kongHeaderRegexPrefix + "^audio/*"},
			},
		},
		{
			msg: "a single exact header match with no type defaults to exact type and converts properly",
			input: []gatewayv1.HTTPHeaderMatch{{
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
		{
			msg: "a single exact header match with a single value converts properly",
			input: []gatewayv1.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayv1.HeaderMatchExact),
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
		{
			msg: "multiple header matches for the same header are rejected",
			input: []gatewayv1.HTTPHeaderMatch{
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
			input: []gatewayv1.HTTPHeaderMatch{
				{
					Type:  lo.ToPtr(gatewayv1.HeaderMatchExact),
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
	grants := []*gatewayv1beta1.ReferenceGrant{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "fitrat",
			},
			Spec: gatewayv1beta1.ReferenceGrantSpec{
				From: []gatewayv1beta1.ReferenceGrantFrom{
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
				To: []gatewayv1beta1.ReferenceGrantTo{
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("GrantOne"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "cholpon",
			},
			Spec: gatewayv1beta1.ReferenceGrantSpec{
				From: []gatewayv1beta1.ReferenceGrantFrom{
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
				To: []gatewayv1beta1.ReferenceGrantTo{
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("GrantTwo"),
					},
				},
			},
		},
	}
	tests := []struct {
		msg    string
		from   gatewayv1beta1.ReferenceGrantFrom
		result map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo
	}{
		{
			msg: "no matches whatsoever",
			from: gatewayv1beta1.ReferenceGrantFrom{
				Group:     gatewayv1alpha2.Group("invalid.example"),
				Kind:      gatewayv1alpha2.Kind("invalid"),
				Namespace: gatewayv1alpha2.Namespace("invalid"),
			},
			result: map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo{},
		},
		{
			msg: "non-matching namespace",
			from: gatewayv1beta1.ReferenceGrantFrom{
				Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
				Kind:      gatewayv1alpha2.Kind("UDPRoute"),
				Namespace: gatewayv1alpha2.Namespace("niyazi"),
			},
			result: map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo{},
		},
		{
			msg: "non-matching kind",
			from: gatewayv1beta1.ReferenceGrantFrom{
				Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
				Kind:      gatewayv1alpha2.Kind("TLSRoute"),
				Namespace: gatewayv1alpha2.Namespace("behbudiy"),
			},
			result: map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo{},
		},
		{
			msg: "non-matching group",
			from: gatewayv1beta1.ReferenceGrantFrom{
				Group:     gatewayv1alpha2.Group("invalid.example"),
				Kind:      gatewayv1alpha2.Kind("UDPRoute"),
				Namespace: gatewayv1alpha2.Namespace("behbudiy"),
			},
			result: map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo{},
		},
		{
			msg: "single match",
			from: gatewayv1beta1.ReferenceGrantFrom{
				Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
				Kind:      gatewayv1alpha2.Kind("UDPRoute"),
				Namespace: gatewayv1alpha2.Namespace("behbudiy"),
			},
			result: map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo{
				"cholpon": {
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("GrantTwo"),
					},
				},
			},
		},
		{
			msg: "multiple matches",
			from: gatewayv1beta1.ReferenceGrantFrom{
				Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
				Kind:      gatewayv1alpha2.Kind("TCPRoute"),
				Namespace: gatewayv1alpha2.Namespace("qodiriy"),
			},
			result: map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo{
				"cholpon": {
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("GrantTwo"),
					},
				},
				"fitrat": {
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("GrantOne"),
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
	grants := []*gatewayv1beta1.ReferenceGrant{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "fitrat",
			},
			Spec: gatewayv1beta1.ReferenceGrantSpec{
				From: []gatewayv1beta1.ReferenceGrantFrom{
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
				To: []gatewayv1beta1.ReferenceGrantTo{
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("Service"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.NewString(),
				Namespace: "cholpon",
			},
			Spec: gatewayv1beta1.ReferenceGrantSpec{
				From: []gatewayv1beta1.ReferenceGrantFrom{
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
				To: []gatewayv1beta1.ReferenceGrantTo{
					{
						Group: gatewayv1alpha2.Group(""),
						Kind:  gatewayv1alpha2.Kind("Service"),
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
	port := gatewayv1.PortNumber(7777)
	redObjName := gatewayv1.ObjectName("red-service")
	blueObjName := gatewayv1.ObjectName("blue-service")
	cholponNamespace := gatewayv1.Namespace("cholpon")
	serviceKind := gatewayv1.Kind("Service")
	serviceGroup := gatewayv1.Group("")
	tests := []struct {
		msg     string
		route   client.Object
		refs    []gatewayv1.BackendRef
		result  kongstate.Service
		wantErr bool
	}{
		{
			msg: "all backends in route namespace",
			route: &gatewayv1.HTTPRoute{
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
			refs: []gatewayv1.BackendRef{
				{
					BackendObjectReference: gatewayv1.BackendObjectReference{
						Name:  blueObjName,
						Kind:  &serviceKind,
						Port:  &port,
						Group: &serviceGroup,
					},
				},
				{
					BackendObjectReference: gatewayv1.BackendObjectReference{
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
				Parent: &gatewayv1.HTTPRoute{
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
			refs: []gatewayv1.BackendRef{
				{
					BackendObjectReference: gatewayv1.BackendObjectReference{
						Name:      blueObjName,
						Port:      &port,
						Kind:      &serviceKind,
						Namespace: &cholponNamespace,
						Group:     &serviceGroup,
					},
				},
				{
					BackendObjectReference: gatewayv1.BackendObjectReference{
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
			refs: []gatewayv1.BackendRef{
				{
					BackendObjectReference: gatewayv1.BackendObjectReference{
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
				Parent: &gatewayv1alpha2.TCPRoute{
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
			refs: []gatewayv1.BackendRef{
				{
					BackendObjectReference: gatewayv1.BackendObjectReference{
						Name:      blueObjName,
						Port:      &port,
						Kind:      &serviceKind,
						Namespace: &cholponNamespace,
						Group:     &serviceGroup,
					},
				},
				{
					BackendObjectReference: gatewayv1.BackendObjectReference{
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
