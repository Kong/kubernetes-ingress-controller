package translator

import (
	"sort"
	"testing"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
)

func TestGetGatewayCerts(t *testing.T) {
	crt, key := certificate.MustGenerateCertPEMFormat(certificate.WithCommonName("example.com"))

	const (
		gwNS        = "gateway-ns"
		secretNS    = "secret-ns"
		secretName  = "prod-tls"
		gwClassName = "kong"
		listener    = "https"
		secretUID   = "7428fb98-180b-4702-a91f-61351a33c6e4"
	)

	makeSecret := func(ns, name string) *corev1.Secret {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				UID:       k8stypes.UID(secretUID),
				Name:      name,
				Namespace: ns,
			},
			Data: map[string][]byte{
				corev1.TLSCertKey:       crt,
				corev1.TLSPrivateKeyKey: key,
			},
		}
	}

	programmedConditions := func(programmed bool) []metav1.Condition {
		if programmed {
			return []metav1.Condition{{
				Type:               string(gatewayapi.ListenerConditionProgrammed),
				Status:             metav1.ConditionTrue,
				Reason:             string(gatewayapi.ListenerReasonProgrammed),
				ObservedGeneration: 0,
			}}
		}
		return []metav1.Condition{{
			Type:               string(gatewayapi.ListenerConditionProgrammed),
			Status:             metav1.ConditionFalse,
			Reason:             string(gatewayapi.ListenerReasonInvalid),
			ObservedGeneration: 0,
		}}
	}

	makeGWC := func(unmanaged bool) *gatewayapi.GatewayClass {
		gwc := &gatewayapi.GatewayClass{
			ObjectMeta: metav1.ObjectMeta{Name: gwClassName},
			Spec:       gatewayapi.GatewayClassSpec{ControllerName: ""},
		}
		if unmanaged {
			gwc.Annotations = map[string]string{
				annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
			}
		}
		return gwc
	}

	// makeGateway builds a Gateway with a single TLS listener.
	// certNS nil means same-namespace (no cross-namespace ref).
	makeGateway := func(ns string, certNS *string, programmed bool) *gatewayapi.Gateway {
		var refNS *gatewayapi.Namespace
		if certNS != nil {
			n := gatewayapi.Namespace(*certNS)
			refNS = &n
		}
		return &gatewayapi.Gateway{
			ObjectMeta: metav1.ObjectMeta{Name: "gw", Namespace: ns},
			Spec: gatewayapi.GatewaySpec{
				GatewayClassName: gatewayapi.ObjectName(gwClassName),
				Listeners: []gatewayapi.Listener{{
					Name:     gatewayapi.SectionName(listener),
					Port:     443,
					Protocol: gatewayapi.HTTPSProtocolType,
					TLS: &gatewayapi.GatewayTLSConfig{
						CertificateRefs: []gatewayapi.SecretObjectReference{{
							Group:     lo.ToPtr(gatewayapi.Group("")),
							Kind:      lo.ToPtr(gatewayapi.Kind("Secret")),
							Name:      gatewayapi.ObjectName(secretName),
							Namespace: refNS,
						}},
					},
				}},
			},
			Status: gatewayapi.GatewayStatus{
				Listeners: []gatewayapi.ListenerStatus{{
					Name:       gatewayapi.SectionName(listener),
					Conditions: programmedConditions(programmed),
				}},
			},
		}
	}

	makeReferenceGrant := func(targetNS, fromNS string, grantedSecretName *string) *gatewayapi.ReferenceGrant {
		to := gatewayapi.ReferenceGrantTo{
			Group: gatewayapi.Group(""),
			Kind:  gatewayapi.Kind("Secret"),
		}
		if grantedSecretName != nil {
			n := gatewayapi.ObjectName(*grantedSecretName)
			to.Name = &n
		}
		return &gatewayapi.ReferenceGrant{
			ObjectMeta: metav1.ObjectMeta{Name: "grant", Namespace: targetNS},
			Spec: gatewayapi.ReferenceGrantSpec{
				From: []gatewayapi.ReferenceGrantFrom{{
					Group:     gatewayapi.V1Group,
					Kind:      "Gateway",
					Namespace: gatewayapi.Namespace(fromNS),
				}},
				To: []gatewayapi.ReferenceGrantTo{to},
			},
		}
	}

	testCases := []struct {
		name          string
		objects       store.FakeObjects
		wantCertCount int
	}{
		{
			name: "same-namespace cert, managed mode",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:       []*gatewayapi.Gateway{makeGateway(gwNS, nil, true)},
				Secrets:        []*corev1.Secret{makeSecret(gwNS, secretName)},
			},
			wantCertCount: 1,
		},
		{
			name: "cross-namespace cert, valid ReferenceGrant without name restriction",
			objects: store.FakeObjects{
				GatewayClasses:  []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:        []*gatewayapi.Gateway{makeGateway(gwNS, lo.ToPtr(secretNS), true)},
				Secrets:         []*corev1.Secret{makeSecret(secretNS, secretName)},
				ReferenceGrants: []*gatewayapi.ReferenceGrant{makeReferenceGrant(secretNS, gwNS, nil)},
			},
			wantCertCount: 1,
		},
		{
			name: "cross-namespace cert, valid ReferenceGrant with matching secret name",
			objects: store.FakeObjects{
				GatewayClasses:  []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:        []*gatewayapi.Gateway{makeGateway(gwNS, lo.ToPtr(secretNS), true)},
				Secrets:         []*corev1.Secret{makeSecret(secretNS, secretName)},
				ReferenceGrants: []*gatewayapi.ReferenceGrant{makeReferenceGrant(secretNS, gwNS, lo.ToPtr(secretName))},
			},
			wantCertCount: 1,
		},
		{
			name: "cross-namespace cert, no ReferenceGrant — cert blocked",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:       []*gatewayapi.Gateway{makeGateway(gwNS, lo.ToPtr(secretNS), false)},
				Secrets:        []*corev1.Secret{makeSecret(secretNS, secretName)},
			},
			wantCertCount: 0,
		},
		{
			name: "cross-namespace cert, ReferenceGrant from wrong source namespace — cert blocked",
			objects: store.FakeObjects{
				GatewayClasses:  []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:        []*gatewayapi.Gateway{makeGateway(gwNS, lo.ToPtr(secretNS), false)},
				Secrets:         []*corev1.Secret{makeSecret(secretNS, secretName)},
				ReferenceGrants: []*gatewayapi.ReferenceGrant{makeReferenceGrant(secretNS, "other-ns", nil)},
			},
			wantCertCount: 0,
		},
		{
			name: "cross-namespace cert, ReferenceGrant names a different secret — cert blocked",
			objects: store.FakeObjects{
				GatewayClasses:  []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:        []*gatewayapi.Gateway{makeGateway(gwNS, lo.ToPtr(secretNS), false)},
				Secrets:         []*corev1.Secret{makeSecret(secretNS, secretName)},
				ReferenceGrants: []*gatewayapi.ReferenceGrant{makeReferenceGrant(secretNS, gwNS, lo.ToPtr("other-secret"))},
			},
			wantCertCount: 0,
		},
		{
			name: "cross-namespace cert, ReferenceGrant lives in wrong namespace — cert blocked",
			objects: store.FakeObjects{
				GatewayClasses:  []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:        []*gatewayapi.Gateway{makeGateway(gwNS, lo.ToPtr(secretNS), false)},
				Secrets:         []*corev1.Secret{makeSecret(secretNS, secretName)},
				ReferenceGrants: []*gatewayapi.ReferenceGrant{makeReferenceGrant("other-ns", gwNS, nil)},
			},
			wantCertCount: 0,
		},
		{
			name: "unmanaged mode, listener programmed",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{makeGWC(true)},
				Gateways:       []*gatewayapi.Gateway{makeGateway(gwNS, nil, true)},
				Secrets:        []*corev1.Secret{makeSecret(gwNS, secretName)},
			},
			wantCertCount: 1,
		},
		{
			name: "unmanaged mode, listener not programmed — cert blocked",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{makeGWC(true)},
				Gateways:       []*gatewayapi.Gateway{makeGateway(gwNS, nil, false)},
				Secrets:        []*corev1.Secret{makeSecret(gwNS, secretName)},
			},
			wantCertCount: 0,
		},
		{
			// managed mode skips the Programmed check; same-namespace cert needs no ReferenceGrant
			name: "managed mode, listener not programmed, same-namespace cert — cert returned",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:       []*gatewayapi.Gateway{makeGateway(gwNS, nil, false)},
				Secrets:        []*corev1.Secret{makeSecret(gwNS, secretName)},
			},
			wantCertCount: 1,
		},
		{
			// ReferenceGrant check fires even when Programmed=True; cross-namespace still needs a grant
			name: "managed mode, listener programmed, cross-namespace cert, no ReferenceGrant — cert blocked",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:       []*gatewayapi.Gateway{makeGateway(gwNS, lo.ToPtr(secretNS), true)},
				Secrets:        []*corev1.Secret{makeSecret(secretNS, secretName)},
			},
			wantCertCount: 0,
		},
		{
			// primary security regression: managed mode + not-programmed + no ReferenceGrant must not exfiltrate
			name: "managed mode, listener not programmed, cross-namespace cert, no ReferenceGrant — cert blocked",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways:       []*gatewayapi.Gateway{makeGateway(gwNS, lo.ToPtr(secretNS), false)},
				Secrets:        []*corev1.Secret{makeSecret(secretNS, secretName)},
			},
			wantCertCount: 0,
		},
		{
			name: "listener missing status entry — listener skipped",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways: []*gatewayapi.Gateway{{
					ObjectMeta: metav1.ObjectMeta{Name: "gw", Namespace: gwNS},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayapi.ObjectName(gwClassName),
						Listeners: []gatewayapi.Listener{{
							Name:     gatewayapi.SectionName(listener),
							Port:     443,
							Protocol: gatewayapi.HTTPSProtocolType,
							TLS: &gatewayapi.GatewayTLSConfig{
								CertificateRefs: []gatewayapi.SecretObjectReference{{
									Name: gatewayapi.ObjectName(secretName),
								}},
							},
						}},
					},
					// Status.Listeners intentionally empty — no matching status entry
				}},
				Secrets: []*corev1.Secret{makeSecret(gwNS, secretName)},
			},
			wantCertCount: 0,
		},
		{
			name: "listener without TLS — no cert",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{makeGWC(false)},
				Gateways: []*gatewayapi.Gateway{{
					ObjectMeta: metav1.ObjectMeta{Name: "gw", Namespace: gwNS},
					Spec: gatewayapi.GatewaySpec{
						GatewayClassName: gatewayapi.ObjectName(gwClassName),
						Listeners: []gatewayapi.Listener{{
							Name:     gatewayapi.SectionName(listener),
							Port:     80,
							Protocol: gatewayapi.HTTPProtocolType,
						}},
					},
					Status: gatewayapi.GatewayStatus{
						Listeners: []gatewayapi.ListenerStatus{{
							Name:       gatewayapi.SectionName(listener),
							Conditions: programmedConditions(true),
						}},
					},
				}},
			},
			wantCertCount: 0,
		},
		{
			name: "GatewayClass controller name mismatch — gateway skipped",
			objects: store.FakeObjects{
				GatewayClasses: []*gatewayapi.GatewayClass{{
					ObjectMeta: metav1.ObjectMeta{Name: gwClassName},
					Spec:       gatewayapi.GatewayClassSpec{ControllerName: "example.com/other-controller"},
				}},
				Gateways: []*gatewayapi.Gateway{makeGateway(gwNS, nil, true)},
				Secrets:  []*corev1.Secret{makeSecret(gwNS, secretName)},
			},
			wantCertCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeStore, err := store.NewFakeStore(tc.objects)
			require.NoError(t, err)
			p := mustNewTranslator(t, fakeStore)
			certs := p.getGatewayCerts()
			require.Len(t, certs, tc.wantCertCount)
		})
	}
}

func TestMergeCerts(t *testing.T) {
	crt1, key1 := certificate.MustGenerateCertPEMFormat(certificate.WithCommonName("foo.com"))
	crt2, key2 := certificate.MustGenerateCertPEMFormat(certificate.WithCommonName("bar.com"))
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
				return *mergedCerts[i].ID < *mergedCerts[j].ID
			})
			require.Equal(t, tc.mergedCerts, mergedCerts)
			require.Equal(t, tc.idToMergedID, idToMergedID)
		})
	}
}
