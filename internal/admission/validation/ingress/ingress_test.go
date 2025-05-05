package ingress

import (
	"context"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestValidateIngress(t *testing.T) {
	for _, tt := range []struct {
		msg           string
		ingress       *netv1.Ingress
		valid         bool
		validationMsg string
		err           error
	}{
		{
			msg: "invalid protocols",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.ProtocolsKey: "ohno",
					},
				},
			},
			valid:         false,
			validationMsg: "Ingress has invalid Kong annotations: invalid konghq.com/protocols value: ohno",
		},
		{
			msg: "invalid protocol combination: mutally exclusive protocols",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing",
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.ProtocolsKey: "http, tcp",
					},
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path:     "/",
											PathType: lo.ToPtr(netv1.PathTypePrefix),
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "httpbin",
													Port: netv1.ServiceBackendPort{Number: int32(80)},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			valid:         false,
			validationMsg: "Ingress failed schema validation: Invalid protocols: http and tcp are mutally exclusive",
		},
	} {
		t.Run(tt.msg, func(t *testing.T) {
			logger := zapr.NewLogger(zap.NewNop())
			fakestore, err := store.NewFakeStore(store.FakeObjects{
				IngressesV1: []*netv1.Ingress{
					tt.ingress,
				},
			})
			require.NoError(t, err)
			valid, validMsg, err := ValidateIngress(
				t.Context(),
				mockRoutesValidator{},
				translator.FeatureFlags{},
				tt.ingress,
				logger,
				fakestore,
			)
			assert.Equal(t, tt.valid, valid, tt.msg)
			assert.Equal(t, tt.validationMsg, validMsg, tt.msg)
			assert.Equal(t, tt.err, err, tt.msg)
		})
	}
}

type mockRoutesValidator struct{}

func (mockRoutesValidator) Validate(_ context.Context, r *kong.Route) (bool, string, error) {
	protocols := lo.FromSlicePtr(r.Protocols)
	if lo.Contains(protocols, "http") && lo.Contains(protocols, "tcp") {
		return false, "Invalid protocols: http and tcp are mutally exclusive", nil
	}
	return true, "", nil
}
