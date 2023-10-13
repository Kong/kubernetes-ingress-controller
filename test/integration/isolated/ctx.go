//go:build integration_tests

package isolated

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
)

type CtxKey[R any] string

func getCtxKey[R any](t *testing.T) CtxKey[R] {
	t.Helper()

	// When we pass t.Name() from inside an `assess` step, the name is in the form TestName/Features/Assess
	if strings.Contains(t.Name(), "/") {
		return CtxKey[R](strings.Split(t.Name(), "/")[0])
	}

	// When pass t.Name() from inside a `testenv.BeforeEachTest` function, the name is just TestName
	return CtxKey[R](t.Name())
}

func SetInCtxForT[R any](ctx context.Context, t *testing.T, r R) context.Context {
	t.Helper()

	return context.WithValue(ctx, getCtxKey[R](t), r)
}

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

func setInCtx[KeyT comparable, R any](ctx context.Context, key KeyT, r R) context.Context {
	return context.WithValue(ctx, key, r)
}

// -----
// RunID
// -----

type _runID struct{}

func SetRunIDInCtx(ctx context.Context, runID string) context.Context {
	return setInCtx(ctx, _runID{}, runID)
}

func GetRunIDFromCtx(ctx context.Context) string {
	r := ctx.Value(_runID{})
	if r == nil {
		return ""
	}
	return r.(string)
}

// -----
// Generic
// -----

type _cluster struct{}

func SetClusterInCtx(ctx context.Context, c clusters.Cluster) context.Context {
	return setInCtx(ctx, _cluster{}, c)
}

func GetClusterFromCtx(ctx context.Context) clusters.Cluster {
	r := ctx.Value(_cluster{})
	if r == nil {
		return nil
	}
	return r.(clusters.Cluster)
}

// -----
// UDPURL
// -----

type _udpURL struct{}

func SetUDPURLInCtx(ctx context.Context, url *url.URL) context.Context {
	return setInCtx(ctx, _udpURL{}, url)
}

func GetUDPURLFromCtx(ctx context.Context) *url.URL {
	u := ctx.Value(_udpURL{})
	if u == nil {
		return nil
	}
	return u.(*url.URL)
}

// -----
// AdminURL
// -----

type _adminURL struct{}

func SetAdminURLInCtx(ctx context.Context, url *url.URL) context.Context {
	return setInCtx(ctx, _adminURL{}, url)
}

func GetAdminURLFromCtx(ctx context.Context) *url.URL {
	u := ctx.Value(_adminURL{})
	if u == nil {
		return nil
	}
	return u.(*url.URL)
}
