//go:build envtest

package envtest

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	commonv1alpha1 "github.com/kong/kubernetes-configuration/api/common/v1alpha1"
	configurationv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	configurationv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	configurationv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/conditions"
)

// TestControlPlaneReferenceHandling tests ControlPlaneReference handling in controllers supporting it.
// It expects that if an object has a ControlPlaneReference set, it should only be programmed if the reference
// is set to 'kic'.
func TestControlPlaneReferenceHandling(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	const ingressClassName = "kongenvtest"
	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)
	deployIngressClass(ctx, t, ingressClassName, ctrlClient)
	ns := CreateNamespace(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithUpdateStatus(),
		WithIngressClass(ingressClassName),
		WithPublishService(ns.Name),
		WithProxySyncSeconds(0.10),
	)

	var (
		kicCPRef = &commonv1alpha1.ControlPlaneRef{
			Type: commonv1alpha1.ControlPlaneRefKIC,
		}
		konnectCPRef = &commonv1alpha1.ControlPlaneRef{
			Type:      commonv1alpha1.ControlPlaneRefKonnectID,
			KonnectID: lo.ToPtr(commonv1alpha1.KonnectIDType("konnect-id")),
		}

		validConsumer = func() *configurationv1.KongConsumer {
			return &configurationv1.KongConsumer{
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
		validConsumerGroup = func() *configurationv1beta1.KongConsumerGroup {
			return &configurationv1beta1.KongConsumerGroup{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "consumer-group-",
					Namespace:    ns.Name,
					Annotations: map[string]string{
						annotations.IngressClassKey: ingressClassName,
					},
				},
				Spec: configurationv1beta1.KongConsumerGroupSpec{
					Name: "consumer-group",
				},
			}
		}
		validVault = func() *configurationv1alpha1.KongVault {
			return &configurationv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "vault-",
					Namespace:    ns.Name,
					Annotations: map[string]string{
						annotations.IngressClassKey: ingressClassName,
					},
				},
				Spec: configurationv1alpha1.KongVaultSpec{
					Backend: "env",
					// Prefix has to be unique for each Vault object as it's validated by KIC in translation.
					Prefix: "prefix-" + lo.RandomString(8, lo.LowerCaseLettersCharset),
				},
			}
		}
	)

	testCases := []struct {
		name   string
		object interface {
			client.Object
			GetConditions() []metav1.Condition
			SetControlPlaneRef(*commonv1alpha1.ControlPlaneRef)
		}
		controlPlaneRef                 *commonv1alpha1.ControlPlaneRef
		expectedErrorOnCreationContains string
	}{
		{
			name:   "KongConsumer - without ControlPlaneRef",
			object: validConsumer(),
		},
		{
			name:            "KongConsumer - with ControlPlaneRef == kic",
			object:          validConsumer(),
			controlPlaneRef: kicCPRef,
		},
		{
			name:                            "KongConsumer - with ControlPlaneRef != kic",
			object:                          validConsumer(),
			controlPlaneRef:                 konnectCPRef,
			expectedErrorOnCreationContains: "spec.controlPlaneRef: Invalid value: \"object\": 'konnectID' type is not supported",
		},
		{
			name:   "KongConsumerGroup - without ControlPlaneRef",
			object: validConsumerGroup(),
		},
		{
			name:            "KongConsumerGroup - with ControlPlaneRef == kic",
			object:          validConsumerGroup(),
			controlPlaneRef: kicCPRef,
		},
		{
			name:                            "KongConsumerGroup - with ControlPlaneRef != kic",
			object:                          validConsumerGroup(),
			controlPlaneRef:                 konnectCPRef,
			expectedErrorOnCreationContains: "spec.controlPlaneRef: Invalid value: \"object\": 'konnectID' type is not supported",
		},
		{
			name:   "KongVault - without ControlPlaneRef",
			object: validVault(),
		},
		{
			name:            "KongVault - with ControlPlaneRef == kic",
			object:          validVault(),
			controlPlaneRef: kicCPRef,
		},
		{
			name:                            "KongVault - with ControlPlaneRef != kic",
			object:                          validVault(),
			controlPlaneRef:                 konnectCPRef,
			expectedErrorOnCreationContains: "spec.controlPlaneRef: Invalid value: \"object\": 'konnectID' type is not supported",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.controlPlaneRef != nil {
				tc.object.SetControlPlaneRef(tc.controlPlaneRef)
			}
			err := ctrlClient.Create(ctx, tc.object)
			if tc.expectedErrorOnCreationContains != "" {
				require.ErrorContains(
					t,
					err,
					tc.expectedErrorOnCreationContains,
				)
				return
			}
			require.NoError(t, err)

			require.EventuallyWithT(t, func(t *assert.CollectT) {
				if !assert.NoError(t, ctrlClient.Get(ctx, client.ObjectKeyFromObject(tc.object), tc.object)) {
					return
				}
				assert.True(t, conditions.Contain(
					tc.object.GetConditions(),
					conditions.WithType(string(configurationv1.ConditionProgrammed)),
					conditions.WithStatus(metav1.ConditionTrue),
				))
			}, waitTime, tickDuration, "expected object to be programmed")
		})
	}
}
