//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"net"
	"net/url"

	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
)

func createResolver(proxyUDPURL *url.URL) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: consts.WaitTick,
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:%d", proxyUDPURL.Hostname(), ktfkong.DefaultUDPServicePort))
		},
	}
}

func not(fn func() bool) func() bool {
	return func() bool {
		return !fn()
	}
}

func urlResolvesSuccessfullyFn(ctx context.Context, proxyUDPURL *url.URL) func() bool {
	return func() bool {
		resolver := createResolver(proxyUDPURL)
		_, err := resolver.LookupHost(ctx, "kernel.org")
		return err == nil
	}
}
