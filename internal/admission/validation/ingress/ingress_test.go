package ingress

import (
	"context"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

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
				context.Background(),
				mockRoutesValidator{},
				translator.FeatureFlags{},
				tt.ingress,
				logger,
				fakestore,
				fake.NewFakeClient(),
			)
			assert.Equal(t, tt.valid, valid, tt.msg)
			assert.Equal(t, tt.validationMsg, validMsg, tt.msg)
			assert.Equal(t, tt.err, err, tt.msg)
		})
	}
}

type mockRoutesValidator struct{}

func (mockRoutesValidator) Validate(_ context.Context, _ *kong.Route) (bool, string, error) {
	return true, "", nil
}
