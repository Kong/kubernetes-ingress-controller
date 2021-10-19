package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

func TestOverrideUpstream(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inUpstream     Upstream
		inKongIngresss *configurationv1.KongIngress
		outUpstream    Upstream
		annotations    map[string]string
	}{
		{
			inUpstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			inKongIngresss: nil,
			outUpstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
		},
		{
			inUpstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			inKongIngresss: &configurationv1.KongIngress{
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
			outUpstream: Upstream{
				Upstream: kong.Upstream{
					Name:               kong.String("foo.com"),
					HashOn:             kong.String("HashOn"),
					HashOnCookie:       kong.String("HashOnCookie"),
					HashOnCookiePath:   kong.String("HashOnCookiePath"),
					HashOnHeader:       kong.String("HashOnHeader"),
					HashFallback:       kong.String("HashFallback"),
					HashFallbackHeader: kong.String("HashFallbackHeader"),
					HostHeader:         kong.String("foo.com"),
					Slots:              kong.Int(42),
				},
			},
			annotations: map[string]string{
				"konghq.com/host-header": "foo.com",
			},
		},
	}

	for _, testcase := range testTable {
		testcase.inUpstream.override(testcase.inKongIngresss, testcase.annotations)
		assert.Equal(testcase.inUpstream, testcase.outUpstream)
	}

	assert.NotPanics(func() {
		var nilUpstream *Upstream
		nilUpstream.override(nil, make(map[string]string))
	})
}
