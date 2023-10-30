package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOverrideUpstream(t *testing.T) {
	testTable := []struct {
		inUpstream  Upstream
		outUpstream Upstream
		svc         *corev1.Service
	}{
		{
			inUpstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
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
			svc: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"konghq.com/host-header": "foo.com",
					},
				},
			},
			outUpstream: Upstream{
				Upstream: kong.Upstream{
					Name:       kong.String("foo.com"),
					HostHeader: kong.String("foo.com"),
				},
			},
		},
	}

	for _, testcase := range testTable {
		testcase.inUpstream.override(testcase.svc)
		assert.Equal(t, testcase.inUpstream, testcase.outUpstream)
	}

	assert.NotPanics(t, func() {
		var nilUpstream *Upstream
		nilUpstream.override(nil)
	})
}
