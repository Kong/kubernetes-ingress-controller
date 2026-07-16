package translator

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
)

func TestGetCACerts(t *testing.T) {
	// validCACert is a self-signed certificate with the CA basic constraint set, so it can be translated.
	validCACert, _ := certificate.MustGenerateCertPEMFormat(
		certificate.WithCATrue(),
		certificate.WithCommonName("ca.example.com"),
	)
	// nonCACert is a valid X.509 certificate but lacks the CA basic constraint, so it fails translation.
	nonCACert, _ := certificate.MustGenerateCertPEMFormat(
		certificate.WithCommonName("leaf.example.com"),
	)

	const dupID = "8214a145-a328-4c56-ab72-2973a56d4eba"

	// makeCACertSecret builds a Secret labelled as a CA cert (so ListCACerts picks it up) with the given data.
	makeCACertSecret := func(namespace, name string, creationTime time.Time, data map[string][]byte) *corev1.Secret {
		return &corev1.Secret{
			// TypeMeta must be set so causing objects have a non-empty GVK; otherwise the
			// failures collector silently drops the registered translation failures.
			TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
			ObjectMeta: metav1.ObjectMeta{
				Name:              name,
				Namespace:         namespace,
				CreationTimestamp: metav1.NewTime(creationTime),
				Labels: map[string]string{
					"konghq.com/ca-cert": "true",
				},
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Data: data,
		}
	}

	// defaultSecret is a convenience wrapper for single-secret tests that don't exercise selection.
	defaultSecret := func(name string, data map[string][]byte) *corev1.Secret {
		return makeCACertSecret("default", name, time.Time{}, data)
	}

	t0 := time.Unix(0, 0)
	t1 := time.Unix(1, 0)
	t2 := time.Unix(2, 0)

	testCases := []struct {
		name                string
		secrets             []*corev1.Secret
		wantCertCount       int
		wantFailureMessages []string
	}{
		{
			name: "valid CA cert is translated",
			secrets: []*corev1.Secret{
				defaultSecret("valid", map[string][]byte{
					"id":   []byte("8214a145-a328-4c56-ab72-2973a56d4eba"),
					"cert": validCACert,
				}),
			},
			wantCertCount:       1,
			wantFailureMessages: nil,
		},
		{
			name: "secret without 'id' field cannot be translated",
			secrets: []*corev1.Secret{
				defaultSecret("no-id", map[string][]byte{
					"cert": validCACert,
				}),
			},
			wantCertCount:       0,
			wantFailureMessages: []string{"missing 'id' field"},
		},
		{
			name: "secret without 'cert' field cannot be translated",
			secrets: []*corev1.Secret{
				defaultSecret("no-cert", map[string][]byte{
					"id": []byte("8214a145-a328-4c56-ab72-2973a56d4eba"),
				}),
			},
			wantCertCount:       0,
			wantFailureMessages: []string{`neither "cert" nor "ca.crt"`},
		},
		{
			// With equal timestamps and namespace, "dup-1" < "dup-2" so "dup-1" is preserved.
			name: "secrets with duplicate 'id': lexicographically smallest name wins",
			secrets: []*corev1.Secret{
				makeCACertSecret("default", "dup-1", t0, map[string][]byte{
					"id":   []byte(dupID),
					"cert": validCACert,
				}),
				makeCACertSecret("default", "dup-2", t0, map[string][]byte{
					"id":   []byte(dupID),
					"cert": validCACert,
				}),
			},
			wantCertCount:       1,
			wantFailureMessages: []string{"duplicate 'id'"},
		},
		{
			// The secret with the earliest creation timestamp is preserved; the later one gets a failure.
			name: "secrets with duplicate 'id': earliest creation time wins",
			secrets: []*corev1.Secret{
				makeCACertSecret("default", "later", t2, map[string][]byte{
					"id":   []byte(dupID),
					"cert": validCACert,
				}),
				makeCACertSecret("default", "earlier", t1, map[string][]byte{
					"id":   []byte(dupID),
					"cert": validCACert,
				}),
			},
			wantCertCount:       1,
			wantFailureMessages: []string{"duplicate 'id'"},
		},
		{
			// With equal timestamps, the secret in the lexicographically smaller namespace wins.
			name: "secrets with duplicate 'id': smallest namespace wins on timestamp tie",
			secrets: []*corev1.Secret{
				makeCACertSecret("ns-b", "same-name", t0, map[string][]byte{
					"id":   []byte(dupID),
					"cert": validCACert,
				}),
				makeCACertSecret("ns-a", "same-name", t0, map[string][]byte{
					"id":   []byte(dupID),
					"cert": validCACert,
				}),
			},
			wantCertCount:       1,
			wantFailureMessages: []string{"duplicate 'id'"},
		},
		{
			// Three objects with the same ID: only the overall winner is translated; the other two get failures.
			name: "secrets with duplicate 'id': three duplicates, only winner translated",
			secrets: []*corev1.Secret{
				makeCACertSecret("default", "c", t2, map[string][]byte{
					"id":   []byte(dupID),
					"cert": validCACert,
				}),
				makeCACertSecret("default", "a", t1, map[string][]byte{
					"id":   []byte(dupID),
					"cert": validCACert,
				}),
				makeCACertSecret("default", "b", t1, map[string][]byte{
					"id":   []byte(dupID),
					"cert": validCACert,
				}),
			},
			wantCertCount:       1,
			wantFailureMessages: []string{"duplicate 'id'", "duplicate 'id'"},
		},
		{
			name: "secret with invalid CA certificate cannot be translated",
			secrets: []*corev1.Secret{
				defaultSecret("not-a-ca", map[string][]byte{
					"id":   []byte("8214a145-a328-4c56-ab72-2973a56d4eba"),
					"cert": nonCACert,
				}),
			},
			wantCertCount:       0,
			wantFailureMessages: []string{"invalid secret CA certificate"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeStore, err := store.NewFakeStore(store.FakeObjects{Secrets: tc.secrets})
			require.NoError(t, err)
			p := mustNewTranslator(t, fakeStore)

			certs := p.getCACerts()
			require.Len(t, certs, tc.wantCertCount)

			translationFailures := p.popTranslationFailures()
			require.Len(t, translationFailures, len(tc.wantFailureMessages))
			// Failures may be reported in any order, so assert that every expected message
			// has a matching failure rather than pairing them by index.
			for _, want := range tc.wantFailureMessages {
				require.Truef(t, func() bool {
					for _, f := range translationFailures {
						if strings.Contains(f.Message(), want) {
							return true
						}
					}
					return false
				}(), "expected a translation failure containing %q, got %v", want, translationFailures)
			}
		})
	}
}
