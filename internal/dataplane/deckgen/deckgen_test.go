package deckgen_test

import (
	"testing"

	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
)

func TestGetFCertificateFromKongCert(t *testing.T) {
const (
	certID = "c6ac927c-4f5a-4e88-8b5d-c7b01d0f43af"
	certPEM = "-----BEGIN CERTIFICATE-----\nfake\n-----END CERTIFICATE-----"
	keyPEM = "-----BEGIN PRIVATE KEY-----\nfake\n-----END PRIVATE KEY-----"
	sniName = "example.com"
	tag1 = "k8s-name:sooper-secret"
	tag2 = "k8s-namespace:bar-namespace"
)

	testCases := []struct {
		name     string
		input    kong.Certificate
		wantTags []*string
	}{
		{
			name: "copies tags",
			input: kong.Certificate{
				ID:   kong.String(certID),
				Cert: kong.String(certPEM),
				Key:  kong.String(keyPEM),
				SNIs: []*string{kong.String(sniName)},
				Tags: []*string{kong.String(tag1), kong.String(tag2)},
			},
			wantTags: []*string{kong.String(tag1), kong.String(tag2)},
		},
		{
			name: "nil tags",
			input: kong.Certificate{
				ID:   kong.String(certID),
				Cert: kong.String(certPEM),
				Key:  kong.String(keyPEM),
				SNIs: []*string{kong.String(sniName)},
				Tags: nil,
			},
			wantTags: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := deckgen.GetFCertificateFromKongCert(tc.input)

			require.Equal(t, tc.input.ID, got.ID)
			require.Equal(t, tc.input.Cert, got.Cert)
			require.Equal(t, tc.input.Key, got.Key)
			require.Equal(t, tc.wantTags, got.Tags)

			require.Len(t, got.SNIs, len(tc.input.SNIs))
			for i, sni := range got.SNIs {
				require.Equal(t, tc.input.SNIs[i], sni.Name)
				require.NotNil(t, sni.Certificate)
				require.Equal(t, tc.input.ID, sni.Certificate.ID)
			}
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
