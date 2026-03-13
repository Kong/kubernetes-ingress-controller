package deckgen_test

import (
	"testing"

	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
)

func TestGetFCertificateFromKongCert(t *testing.T) {
	testCases := []struct {
		name     string
		inmemory bool
		cert     kong.Certificate
		want     file.FCertificate
	}{
		{
			name:     "empty certificate",
			inmemory: false,
			cert:     kong.Certificate{},
			want: file.FCertificate{
				SNIs: []kong.SNI{},
			},
		},
		{
			name:     "all fields set, inmemory=true, SNIs have no certificate ref",
			inmemory: true,
			cert: kong.Certificate{
				ID:   kong.String("cert-id"),
				Key:  kong.String("cert-key"),
				Cert: kong.String("cert-pem"),
				SNIs: []*string{kong.String("example.com"), kong.String("other.com")},
			},
			want: file.FCertificate{
				ID:   kong.String("cert-id"),
				Key:  kong.String("cert-key"),
				Cert: kong.String("cert-pem"),
				SNIs: []kong.SNI{
					{Name: kong.String("example.com")},
					{Name: kong.String("other.com")},
				},
			},
		},
		{
			name:     "all fields set, inmemory=false, SNIs have certificate ref",
			inmemory: false,
			cert: kong.Certificate{
				ID:   kong.String("cert-id"),
				Key:  kong.String("cert-key"),
				Cert: kong.String("cert-pem"),
				SNIs: []*string{kong.String("example.com")},
			},
			want: file.FCertificate{
				ID:   kong.String("cert-id"),
				Key:  kong.String("cert-key"),
				Cert: kong.String("cert-pem"),
				SNIs: []kong.SNI{
					{
						Name:        kong.String("example.com"),
						Certificate: &kong.Certificate{ID: kong.String("cert-id")},
					},
				},
			},
		},
		{
			name:     "nil ID, inmemory=false, SNIs have no certificate ref",
			inmemory: false,
			cert: kong.Certificate{
				Key:  kong.String("cert-key"),
				Cert: kong.String("cert-pem"),
				SNIs: []*string{kong.String("example.com")},
			},
			want: file.FCertificate{
				Key:  kong.String("cert-key"),
				Cert: kong.String("cert-pem"),
				SNIs: []kong.SNI{
					{Name: kong.String("example.com")},
				},
			},
		},
		{
			name:     "no SNIs",
			inmemory: false,
			cert: kong.Certificate{
				ID:   kong.String("cert-id"),
				Key:  kong.String("cert-key"),
				Cert: kong.String("cert-pem"),
			},
			want: file.FCertificate{
				ID:   kong.String("cert-id"),
				Key:  kong.String("cert-key"),
				Cert: kong.String("cert-pem"),
				SNIs: []kong.SNI{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := deckgen.GetFCertificateFromKongCert(tc.inmemory, tc.cert)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestIsContentEmpty(t *testing.T) {
	testCases := []struct {
		name    string
		content *file.Content
		want    bool
	}{
		{
			name: "non-empty content",
			content: &file.Content{
				Upstreams: []file.FUpstream{
					{
						Upstream: kong.Upstream{
							Name: kong.String("test"),
						},
					},
				},
			},
			want: false,
		},
		{
			name:    "empty content",
			content: &file.Content{},
			want:    true,
		},
		{
			name: "empty with version and info",
			content: &file.Content{
				FormatVersion: "1.1",
				Info: &file.Info{
					SelectorTags: []string{"tag1", "tag2"},
				},
			},
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := deckgen.IsContentEmpty(tc.content)
			require.Equal(t, tc.want, got)
		})
	}
}
