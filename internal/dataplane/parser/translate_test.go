package parser

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
)

func TestPathsFromK8s(t *testing.T) {
	for _, tt := range []struct {
		name         string
		path         string
		wantPrefix   []*string
		wantExact    []*string
		wantImplSpec []*string
	}{
		{
			name:         "empty",
			wantPrefix:   kong.StringSlice("/"),
			wantExact:    kong.StringSlice("/$"),
			wantImplSpec: kong.StringSlice("/"),
		},
		{
			name:         "root",
			path:         "/",
			wantPrefix:   kong.StringSlice("/"),
			wantExact:    kong.StringSlice("/$"),
			wantImplSpec: kong.StringSlice("/"),
		},
		{
			name:         "one segment, no trailing slash",
			path:         "/foo",
			wantPrefix:   kong.StringSlice("/foo$", "/foo/"),
			wantExact:    kong.StringSlice("/foo$"),
			wantImplSpec: kong.StringSlice("/foo"),
		},
		{
			name:         "one segment, has trailing slash",
			path:         "/foo/",
			wantPrefix:   kong.StringSlice("/foo$", "/foo/"),
			wantExact:    kong.StringSlice("/foo/$"),
			wantImplSpec: kong.StringSlice("/foo/"),
		},
		{
			name:         "two segments, no trailing slash",
			path:         "/foo/bar",
			wantPrefix:   kong.StringSlice("/foo/bar$", "/foo/bar/"),
			wantExact:    kong.StringSlice("/foo/bar$"),
			wantImplSpec: kong.StringSlice("/foo/bar"),
		},
		{
			name:         "two segments, has trailing slash",
			path:         "/foo/bar/",
			wantPrefix:   kong.StringSlice("/foo/bar$", "/foo/bar/"),
			wantExact:    kong.StringSlice("/foo/bar/$"),
			wantImplSpec: kong.StringSlice("/foo/bar/"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			{
				gotPrefix, gotErr := pathsFromK8s(tt.path, networkingv1.PathTypePrefix)
				require.NoError(t, gotErr)
				require.Equal(t, tt.wantPrefix, gotPrefix, "prefix match")
			}
			{
				gotExact, gotErr := pathsFromK8s(tt.path, networkingv1.PathTypeExact)
				require.NoError(t, gotErr)
				require.Equal(t, tt.wantExact, gotExact, "exact match")
			}
			{
				gotImplSpec, gotErr := pathsFromK8s(tt.path, networkingv1.PathTypeImplementationSpecific)
				require.NoError(t, gotErr)
				require.Equal(t, tt.wantImplSpec, gotImplSpec, "implementation specific match")
			}
		})
	}
}
