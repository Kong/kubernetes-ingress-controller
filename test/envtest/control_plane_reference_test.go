//go:build envtest

package envtest

import (
	"context"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/conditions"
)

// TestControlPlaneReferenceHandling tests ControlPlaneReference handling in controllers supporting it.
// It expects that if an object has a ControlPlaneReference set, it should only be programmed if the reference
// is set to 'kic'.
func TestControlPlaneReferenceHandling(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const ingressClassName = "kongenvtest"
	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)
	deployIngressClass(ctx, t, ingressClassName, ctrlClient)
	logger := zapr.NewLogger(zap.NewNop())
	ctrl.SetLogger(logger)
	ns := CreateNamespace(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithUpdateStatus(),
		WithIngressClass(ingressClassName),
		WithPublishService(ns.Name),
		WithProxySyncSeconds(0.10),
	)

	var (
		kicCPRef = &kongv1alpha1.ControlPlaneRef{
			Type: kongv1alpha1.ControlPlaneRefKIC,
		}
		konnectCPRef = &kongv1alpha1.ControlPlaneRef{
			Type:      kongv1alpha1.ControlPlaneRefKonnectID,
			KonnectID: lo.ToPtr("konnect-id"),
		}

		validConsumer = func() *kongv1.KongConsumer {
			return &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "consumer-",
					Namespace:    ns.Name,
					Annotations: map[string]string{
						annotations.IngressClassKey: ingressClassName,
					},
				},
				Username: "consumer",
			}
		}
		validConsumerGroup = func() *kongv1beta1.KongConsumerGroup {
			return &kongv1beta1.KongConsumerGroup{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "consumer-group-",
					Namespace:    ns.Name,
					Annotations: map[string]string{
						annotations.IngressClassKey: ingressClassName,
					},
				},
				Spec: kongv1beta1.KongConsumerGroupSpec{
					Name: "consumer-group",
				},
			}
		}
		validVault = func() *kongv1alpha1.KongVault {
			return &kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "vault-",
					Namespace:    ns.Name,
					Annotations: map[string]string{
						annotations.IngressClassKey: ingressClassName,
					},
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend: "env",
					Prefix:  "env",
				},
			}
		}
	)

	testCases := []struct {
		name   string
		object interface {
			client.Object
			GetConditions() []metav1.Condition
			SetControlPlaneRef(*kongv1alpha1.ControlPlaneRef)
		}
		controlPlaneRef      *kongv1alpha1.ControlPlaneRef
		expectToBeProgrammed bool
	}{
		{
			name:                 "KongConsumer - without ControlPlaneRef",
			object:               validConsumer(),
			expectToBeProgrammed: true,
		},
		{
			name:                 "KongConsumer - with ControlPlaneRef == kic",
			object:               validConsumer(),
			controlPlaneRef:      kicCPRef,
			expectToBeProgrammed: true,
		},
		{
			name:                 "KongConsumer - with ControlPlaneRef != kic",
			object:               validConsumer(),
			controlPlaneRef:      konnectCPRef,
			expectToBeProgrammed: false,
		},
		{
			name:                 "KongConsumerGroup - without ControlPlaneRef",
			object:               validConsumerGroup(),
			expectToBeProgrammed: true,
		},
		{
			name:                 "KongConsumerGroup - with ControlPlaneRef == kic",
			object:               validConsumerGroup(),
			controlPlaneRef:      kicCPRef,
			expectToBeProgrammed: true,
		},
		{
			name:                 "KongConsumerGroup - with ControlPlaneRef != kic",
			object:               validConsumerGroup(),
			controlPlaneRef:      konnectCPRef,
			expectToBeProgrammed: false,
		},
		{
			name:                 "KongVault - without ControlPlaneRef",
			object:               validVault(),
			expectToBeProgrammed: true,
		},
		{
			name:                 "KongVault - with ControlPlaneRef == kic",
			object:               validVault(),
			controlPlaneRef:      kicCPRef,
			expectToBeProgrammed: true,
		},
		{
			name:                 "KongVault - with ControlPlaneRef != kic",
			object:               validVault(),
			controlPlaneRef:      konnectCPRef,
			expectToBeProgrammed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.controlPlaneRef != nil {
				tc.object.SetControlPlaneRef(tc.controlPlaneRef)
			}
			require.NoError(t, ctrlClient.Create(ctx, tc.object))

			if tc.expectToBeProgrammed {
				require.EventuallyWithT(t, func(t *assert.CollectT) {
					if !assert.NoError(t, ctrlClient.Get(ctx, client.ObjectKeyFromObject(tc.object), tc.object)) {
						return
					}
					assert.Equal(t, tc.expectToBeProgrammed, conditions.Contain(
						tc.object.GetConditions(),
						conditions.WithType(string(kongv1.ConditionProgrammed)),
						conditions.WithStatus(metav1.ConditionTrue),
					))
				}, waitTime, tickDuration, "expected object to be programmed")
			} else {
				// We'll wait for `waitTime` to ensure the object does not get programmed. We need a following boolean
				// to make sure the object was fetched successfully at least once.
				var wasObjectSuccessfullyFetched bool
				require.Never(t, func() bool {
					err := ctrlClient.Get(ctx, client.ObjectKeyFromObject(tc.object), tc.object)
					if err != nil {
						t.Logf("Error fetching object: %v", err)
						return false // Most likely that would is NotFound error. We want to keep waiting in any case.
					}
					wasObjectSuccessfullyFetched = true
					return conditions.Contain(
						tc.object.GetConditions(),
						conditions.WithType(string(kongv1.ConditionProgrammed)),
						conditions.WithStatus(metav1.ConditionTrue),
					)
				}, waitTime, tickDuration, "expected object not to be programmed")
				assert.True(t, wasObjectSuccessfullyFetched, "expected object to be fetched at least once")
			}
		})
	}
}
