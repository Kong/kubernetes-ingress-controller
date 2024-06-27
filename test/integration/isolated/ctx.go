//go:build integration_tests

package isolated

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
)

// SetInCtxForT sets a value in the context that can be later retrieved with GetFromCtxForT.
// It has 1 type parameter which is the type of the value to set in context.
//
// Due to the usage of type parameter this allows to get only 1 value from the context.
// If users want to be able to store more than 1 value of the same type then they
// should use e.g. a slice or a struct containing those.
func SetInCtxForT[R any](ctx context.Context, t *testing.T, r R) context.Context {
	t.Helper()

	return context.WithValue(ctx, getCtxKey[R](t), r)
}

// GetFromCtxForT gets a value from the context that was previously set with SetInCtxForT.
// It has 1 type parameter which is the type of the value to get.
// It will fail the test if the value is not found or if the type does not match.
//
// Due to the usage of type parameter this allows to get only 1 value from the context.
// If users want to be able to store more than 1 value of the same type then they
// should use e.g. a slice or a struct containing those.
func GetFromCtxForT[R any](ctx context.Context, t *testing.T) R {
	t.Helper()

	raw := ctx.Value(getCtxKey[R](t))
	result, ok := raw.(R)
	if !ok {
		var r R
		t.Fatalf("required %T to be stored in context but found: %s (of type %T)", r, raw, raw)
	}
	return result
}

type ctxKey[R any] string

func getCtxKey[R any](t *testing.T) ctxKey[R] {
	t.Helper()

	// When we pass t.Name() from inside an `assess` step, the name is in the form TestName/Features/Assess.
	if strings.Contains(t.Name(), "/") {
		return ctxKey[R](strings.Split(t.Name(), "/")[0])
	}

	// When we pass t.Name() from inside a `testenv.BeforeEachTest` function, the name is just TestName.
	return ctxKey[R](t.Name())
}

func setInCtx[KeyT comparable, R any](ctx context.Context, key KeyT, r R) context.Context {
	return context.WithValue(ctx, key, r)
}

type _cluster struct{}

// SetClusterInCtx sets the cluster in the context.
func SetClusterInCtx(ctx context.Context, c clusters.Cluster) context.Context {
	return setInCtx(ctx, _cluster{}, c)
}

// GetClusterFromCtx gets the cluster from the context.
func GetClusterFromCtx(ctx context.Context) clusters.Cluster {
	r := ctx.Value(_cluster{})
	if r == nil {
		return nil
	}
	return r.(clusters.Cluster)
}

type _runID struct{}

// SetRunIDInCtx sets the runID in the context.
func SetRunIDInCtx(ctx context.Context, runID string) context.Context {
	return setInCtx(ctx, _runID{}, runID)
}

// GetRunIDFromCtx gets the runID from the context.
func GetRunIDFromCtx(ctx context.Context) string {
	r := ctx.Value(_runID{})
	if r == nil {
		return ""
	}
	return r.(string)
}

type _udpURL struct{}

// SetUDPURLInCtx sets the UDP URL in the context.
func SetUDPURLInCtx(ctx context.Context, url string) context.Context {
	return setInCtx(ctx, _udpURL{}, url)
}

// GetUDPURLFromCtx gets the UDP URL from the context.
func GetUDPURLFromCtx(ctx context.Context) string {
	return ctx.Value(_udpURL{}).(string)
}

type _tlsURL struct{}

// SetTLSURLInCtx sets the TLS URL in the context.
func SetTLSURLInCtx(ctx context.Context, url string) context.Context {
	return setInCtx(ctx, _tlsURL{}, url)
}

// GetTLSURLFromCtx gets the TLS URL from the context.
func GetTLSURLFromCtx(ctx context.Context) string {
	return ctx.Value(_tlsURL{}).(string)
}

type _tcpURL struct{}

// SetTCPURLInCtx sets the TCP URL in the context.
func SetTCPURLInCtx(ctx context.Context, url string) context.Context {
	return setInCtx(ctx, _tcpURL{}, url)
}

// GetTCPURLFromCtx gets the TCP URL from the context.
func GetTCPURLFromCtx(ctx context.Context) string {
	return ctx.Value(_tcpURL{}).(string)
}

type _proxyHTTPURL struct{}

// SetHTTPURLInCtx sets the proxy URL in the context.
func SetHTTPURLInCtx(ctx context.Context, url *url.URL) context.Context {
	return setInCtx(ctx, _proxyHTTPURL{}, url)
}

// GetHTTPURLFromCtx gets the proxy URL from the context.
func GetHTTPURLFromCtx(ctx context.Context) *url.URL {
	u := ctx.Value(_proxyHTTPURL{})
	if u == nil {
		return nil
	}
	return u.(*url.URL)
}

type _proxyHTTPSURL struct{}

// SetHTTPSURLInCtx sets the proxy URL in the context.
func SetHTTPSURLInCtx(ctx context.Context, url *url.URL) context.Context {
	return setInCtx(ctx, _proxyHTTPSURL{}, url)
}

// GetHTTPSURLFromCtx gets the proxy URL from the context.
func GetHTTPSURLFromCtx(ctx context.Context) *url.URL {
	u := ctx.Value(_proxyHTTPSURL{})
	if u == nil {
		return nil
	}
	return u.(*url.URL)
}

type _adminURL struct{}

// SetAdminURLInCtx sets the admin URL in the context.
func SetAdminURLInCtx(ctx context.Context, url *url.URL) context.Context {
	return setInCtx(ctx, _adminURL{}, url)
}

// GetAdminURLFromCtx gets the admin URL from the context.
func GetAdminURLFromCtx(ctx context.Context) *url.URL {
	u := ctx.Value(_adminURL{})
	if u == nil {
		return nil
	}
	return u.(*url.URL)
}

type _diagURL struct{}

// SetDiagURLInCtx sets the diag URL in the context.
func SetDiagURLInCtx(ctx context.Context, url *url.URL) context.Context {
	return setInCtx(ctx, _diagURL{}, url)
}

// GetDiagURLFromCtx gets the diag URL from the context.
func GetDiagURLFromCtx(ctx context.Context) *url.URL {
	u := ctx.Value(_diagURL{})
	if u == nil {
		return nil
	}
	return u.(*url.URL)
}

type _ingressClass struct{}

// GetIngressClassFromCtx gets the Ingress Class from the context.
func GetIngressClassFromCtx(ctx context.Context) string {
	r := ctx.Value(_ingressClass{})
	if r == nil {
		return ""
	}
	return r.(string)
}
