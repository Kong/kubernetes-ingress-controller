package translator

import (
	"sort"
	"testing"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/stretchr/testify/require"
)

func TestMergeCerts(t *testing.T) {
	crt1, key1 := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCommonName("foo.com"))
	crt2, key2 := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCommonName("bar.com"))
	testCases := []struct {
		name         string
		certs        []certWrapper
		mergedCerts  []kongstate.Certificate
		idToMergedID certIDToMergedCertID
	}{
		{
			name: "single certificate",
			certs: []certWrapper{
				{
					identifier: string(crt1) + string(key1),
					cert: kong.Certificate{
						ID:   kong.String("certificate-1"),
						Cert: kong.String(string(crt1)),
						Key:  kong.String(string(key1)),
					},
					snis: []string{"foo.com"},
				},
			},
			mergedCerts: []kongstate.Certificate{
				{
					Certificate: kong.Certificate{
						ID:   kong.String("certificate-1"),
						Cert: kong.String(string(crt1)),
						Key:  kong.String(string(key1)),
						SNIs: kong.StringSlice("foo.com"),
					},
				},
			},
			idToMergedID: certIDToMergedCertID{"certificate-1": "certificate-1"},
		},
		{
			name: "multiple different certifcates",
			certs: []certWrapper{
				{
					identifier: string(crt1) + string(key1),
					cert: kong.Certificate{
						ID:   kong.String("certificate-1"),
						Cert: kong.String(string(crt1)),
						Key:  kong.String(string(key1)),
					},
					snis: []string{"foo.com"},
				},
				{
					identifier: string(crt2) + string(key2),
					cert: kong.Certificate{
						ID:   kong.String("certificate-2"),
						Cert: kong.String(string(crt2)),
						Key:  kong.String(string(key2)),
					},
					snis: []string{"bar.com"},
				},
			},
			mergedCerts: []kongstate.Certificate{
				{
					Certificate: kong.Certificate{
						ID:   kong.String("certificate-1"),
						Cert: kong.String(string(crt1)),
						Key:  kong.String(string(key1)),
						SNIs: kong.StringSlice("foo.com"),
					},
				},
				{
					Certificate: kong.Certificate{
						ID:   kong.String("certificate-2"),
						Cert: kong.String(string(crt2)),
						Key:  kong.String(string(key2)),
						SNIs: kong.StringSlice("bar.com"),
					},
				},
			},
			idToMergedID: certIDToMergedCertID{
				"certificate-1": "certificate-1",
				"certificate-2": "certificate-2",
			},
		},
		{
			name: "multiple certs with same content should be merged",
			certs: []certWrapper{
				{
					identifier: string(crt1) + string(key1),
					cert: kong.Certificate{
						ID:   kong.String("certificate-1"),
						Cert: kong.String(string(crt1)),
						Key:  kong.String(string(key1)),
					},
					snis: []string{"foo.com"},
				},
				{
					identifier: string(crt1) + string(key1),
					cert: kong.Certificate{
						ID:   kong.String("certificate-1-1"),
						Cert: kong.String(string(crt1)),
						Key:  kong.String(string(key1)),
					},
					snis: []string{"baz.com"},
				},
			},
			mergedCerts: []kongstate.Certificate{
				{
					Certificate: kong.Certificate{
						ID:   kong.String("certificate-1"),
						Cert: kong.String(string(crt1)),
						Key:  kong.String(string(key1)),
						// SNIs should be sorted
						SNIs: kong.StringSlice("baz.com", "foo.com"),
					},
				},
			},
			idToMergedID: certIDToMergedCertID{
				"certificate-1":   "certificate-1",
				"certificate-1-1": "certificate-1",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mergedCerts, idToMergedID := mergeCerts(logr.Discard(), tc.certs)
			// sort certs by their IDs to make a stable order of the result merged certs.
			sort.SliceStable(mergedCerts, func(i, j int) bool {
				return *mergedCerts[i].Certificate.ID < *mergedCerts[j].Certificate.ID
			})
			require.Equal(t, tc.mergedCerts, mergedCerts)
			require.Equal(t, tc.idToMergedID, idToMergedID)
		})
	}
}
