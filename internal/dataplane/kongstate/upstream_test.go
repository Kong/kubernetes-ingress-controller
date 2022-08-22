package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

func TestOverrideUpstream(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inUpstream     Upstream
		inKongIngresss *configurationv1.KongIngress
		outUpstream    Upstream
		svc            *corev1.Service
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
				Upstream: &configurationv1.KongIngressUpstream{
					HashOn:                 kong.String("HashOn"),
					HashOnCookie:           kong.String("HashOnCookie"),
					HashOnCookiePath:       kong.String("HashOnCookiePath"),
					HashOnHeader:           kong.String("HashOnHeader"),
					HashFallback:           kong.String("HashFallback"),
					HashFallbackHeader:     kong.String("HashFallbackHeader"),
					HashOnQueryArg:         kong.String("HashOnQueryArg"),
					HashFallbackQueryArg:   kong.String("HashFallbackQueryArg"),
					HashOnURICapture:       kong.String("HashOnURICapture"),
					HashFallbackURICapture: kong.String("HashFallbackURICapture"),
					Slots:                  kong.Int(42),
				},
			},
			outUpstream: Upstream{
				Upstream: kong.Upstream{
					Name:                   kong.String("foo.com"),
					HashOn:                 kong.String("HashOn"),
					HashOnCookie:           kong.String("HashOnCookie"),
					HashOnCookiePath:       kong.String("HashOnCookiePath"),
					HashOnHeader:           kong.String("HashOnHeader"),
					HashFallback:           kong.String("HashFallback"),
					HashFallbackHeader:     kong.String("HashFallbackHeader"),
					HostHeader:             kong.String("foo.com"),
					HashOnQueryArg:         kong.String("HashOnQueryArg"),
					HashFallbackQueryArg:   kong.String("HashFallbackQueryArg"),
					HashOnURICapture:       kong.String("HashOnURICapture"),
					HashFallbackURICapture: kong.String("HashFallbackURICapture"),
					Slots:                  kong.Int(42),
				},
			},
			svc: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"konghq.com/host-header": "foo.com",
					},
				},
			},
		},
	}

	for _, testcase := range testTable {
		testcase.inUpstream.override(testcase.inKongIngresss, testcase.svc)
		assert.Equal(testcase.inUpstream, testcase.outUpstream)
	}

	assert.NotPanics(func() {
		var nilUpstream *Upstream
		nilUpstream.override(nil, nil)
	})
}
