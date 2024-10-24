package subtranslator

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
)

func TestApplyExpressionToL4KongRoute(t *testing.T) {
	testCases := []struct {
		name    string
		route   kong.Route
		subExpr string
	}{
		{
			name:    "destination port",
			subExpr: "net.dst.port == 1234",
			route: kong.Route{
				Destinations: []*kong.CIDRPort{
					{
						Port: lo.ToPtr(1234),
					},
				},
				Protocols: []*string{
					lo.ToPtr("tcp"),
				},
			},
		},
		{
			name:    "multiple destination ports",
			subExpr: "(net.dst.port == 1234) || (net.dst.port == 5678)",
			route: kong.Route{
				Destinations: []*kong.CIDRPort{
					{
						Port: lo.ToPtr(1234),
					},
					{
						Port: lo.ToPtr(5678),
					},
				},
				Protocols: []*string{
					lo.ToPtr("tcp"),
				},
			},
		},
		{
			name:    "SNI host",
			subExpr: "tls.sni == \"example.com\"",
			route: kong.Route{
				SNIs: []*string{
					lo.ToPtr("example.com"),
				},
				Protocols: []*string{
					lo.ToPtr("tcp"),
				},
			},
		},
		{
			name:    "multiple SNI hosts",
			subExpr: "(tls.sni == \"example.com\") || (tls.sni == \"example.net\")",
			route: kong.Route{
				SNIs: []*string{
					lo.ToPtr("example.com"),
					lo.ToPtr("example.net"),
				},
				Protocols: []*string{
					lo.ToPtr("tcp"),
				},
			},
		},
		{
			name:    "SNI host and multiple destination ports",
			subExpr: "(tls.sni == \"example.com\") && ((net.dst.port == 1234) || (net.dst.port == 5678))",
			route: kong.Route{
				Destinations: []*kong.CIDRPort{
					{
						Port: lo.ToPtr(1234),
					},
					{
						Port: lo.ToPtr(5678),
					},
				},
				SNIs: []*string{
					lo.ToPtr("example.com"),
				},
				Protocols: []*string{
					lo.ToPtr("tcp"),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wrapped := kongstate.Route{Route: tc.route}
			ApplyExpressionToL4KongRoute(&wrapped)
			require.Contains(t, *wrapped.Route.Expression, tc.subExpr)
		})
	}
}
