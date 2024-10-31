//go:build envtest

package envtest

import (
	"context"
	"strings"
	"testing"
	"time"

	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
)

func TestKongStateFillConsumersAndCredentialsFailure(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 10 * time.Second
		tickTime = 100 * time.Millisecond
	)

	scheme := Scheme(t, WithKong)
	cfg := Setup(t, scheme)
	client := NewControllerClient(t, scheme, cfg)

	// We use a deferred cancel to stop the manager and not wait for its timeout.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ns := CreateNamespace(ctx, t, client)

	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "key-auth-cred",
				Namespace: ns.Name,
				Labels: map[string]string{
					labels.CredentialTypeLabel: "key-auth",
				},
			},
			Data: map[string][]byte{
				"key": []byte("whatever"),
				"ttl": []byte("1024"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "empty-cred",
				Namespace: ns.Name,
				Labels: map[string]string{
					labels.CredentialTypeLabel: "key-auth",
				},
			},
			Data: map[string][]byte{},
		},
	}
	for _, secret := range secrets {
		require.NoError(t, client.Create(ctx, secret))
	}

	kongConsumers := []*kongv1.KongConsumer{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "consumer-key-auth-cred",
				Namespace:   ns.Name,
				Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
			},
			Username: "foo",
			Credentials: []string{
				"key-auth-cred",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "consumer-empty-cred",
				Namespace:   ns.Name,
				Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
			},
			CustomID: "bar",
			Credentials: []string{
				"empty-cred",
			},
		},
	}
	for _, kongConsumer := range kongConsumers {
		require.NoError(t, client.Create(ctx, kongConsumer))
	}

	// These KongConsumers should fail admission via the CRD Validation Expressions.
	brokenKongConsumers := []*kongv1.KongConsumer{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "consumer-no-username-and-no-custom-id",
				Namespace:   ns.Name,
				Annotations: map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
			},
			Credentials: []string{
				"key-auth-cred",
			},
		},
	}
	for _, brokenKongConsumer := range brokenKongConsumers {
		require.Error(t, client.Create(ctx, brokenKongConsumer))
	}

	// KongConsumer name -> event message
	kongConsumerTranslationFailureMessages := map[string]string{
		"consumer-empty-cred": `credential "empty-cred" failure: failed to provision credential: key-auth is invalid: no key`,
	}

	RunManager(ctx, t, cfg, AdminAPIOptFns(), WithProxySyncSeconds(0.5))

	require.Eventually(t, func() bool {
		events := &corev1.EventList{}
		err := client.List(ctx, events, &ctrlclient.ListOptions{
			Namespace: ns.Name,
		})
		if err != nil {
			t.Logf("Failed to list events in namespace %s: %v", ns.Name, err)
			return false
		}

		for name, msg := range kongConsumerTranslationFailureMessages {
			// find the translation failure event attached to each expected KongConumser.
			_, found := lo.Find(events.Items, func(e corev1.Event) bool {
				return e.InvolvedObject.Kind == "KongConsumer" && e.InvolvedObject.Name == name &&
					e.Reason == "KongConfigurationTranslationFailed" &&
					strings.Contains(e.Message, msg)
			})
			if !found {
				return false
			}
		}
		return true
	}, waitTime, tickTime)
}
