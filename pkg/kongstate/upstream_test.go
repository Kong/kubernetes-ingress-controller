package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/stretchr/testify/assert"
)

func TestOverrideUpstream(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inUpstream     Upstream
		inKongIngresss configurationv1.KongIngress
		outUpstream    Upstream
	}{
		{
			Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			configurationv1.KongIngress{
				Upstream: &kong.Upstream{},
			},
			Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
		},
		{
			Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			configurationv1.KongIngress{
				Upstream: &kong.Upstream{
					Name:               kong.String("wrong.com"),
					HashOn:             kong.String("HashOn"),
					HashOnCookie:       kong.String("HashOnCookie"),
					HashOnCookiePath:   kong.String("HashOnCookiePath"),
					HashOnHeader:       kong.String("HashOnHeader"),
					HashFallback:       kong.String("HashFallback"),
					HashFallbackHeader: kong.String("HashFallbackHeader"),
					Slots:              kong.Int(42),
				},
			},
			Upstream{
				Upstream: kong.Upstream{
					Name:               kong.String("foo.com"),
					HashOn:             kong.String("HashOn"),
					HashOnCookie:       kong.String("HashOnCookie"),
					HashOnCookiePath:   kong.String("HashOnCookiePath"),
					HashOnHeader:       kong.String("HashOnHeader"),
					HashFallback:       kong.String("HashFallback"),
					HashFallbackHeader: kong.String("HashFallbackHeader"),
					Slots:              kong.Int(42),
				},
			},
		},
	}

	for _, testcase := range testTable {
		testcase.inUpstream.override(&testcase.inKongIngresss, make(map[string]string))
		assert.Equal(testcase.inUpstream, testcase.outUpstream)
	}

	assert.NotPanics(func() {
		var nilUpstream *Upstream
		nilUpstream.override(nil, make(map[string]string))
	})
}
