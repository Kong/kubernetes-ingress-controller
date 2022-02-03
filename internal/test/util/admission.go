package util

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	admregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// -----------------------------------------------------------------------------
// Test Utilities - Admission Webhooks - Public Vars & Consts
// -----------------------------------------------------------------------------

const (
	// XXX (this hack is tracked in https://github.com/Kong/kubernetes-ingress-controller/issues/1613):
	//
	// The test process (`go test github.com/Kong/kubernetes-ingress-controller/test/integration/...`) serves the webhook
	// endpoints to be consumed by the apiserver (so that the tests can apply a ValidatingWebhookConfiguration and test
	// those validations).
	// In order to make that possible, we needed to allow the apiserver (that gets spun up by the test harness) to access
	// the system under test (which runs as a part of the `go test` process).
	// In the constants below, we're making an audacious assumption that the host running the `go test` process is also
	// the Docker host on the default bridge (therefore it can listen on 172.17.0.1), and that the apiserver
	// is running within a context (such as KIND running on that same docker bridge), from which 172.17.0.1 is routable.
	// This works if the test runs against a KIND cluster, and does not work against cloud providers (like GKE).
	AdmissionWebhookListenHost = "172.17.0.1"
	AdmissionWebhookListenPort = 49023

	webhookTimeout = time.Minute * 3
)

// -----------------------------------------------------------------------------
// Test Utilities - Admission Webhooks - Public Functions
// -----------------------------------------------------------------------------

// EnsureAdmissionRegistration ensures that for a provided list of admission
// rules with operations a webhook is registered in the provided test cluster
// and creates a Service where that webhook will be exposed.
//
// A cleanup function is returned which will delete the registered webhook
// when the caller is ready.
func EnsureAdmissionRegistration(ctx context.Context, cluster clusters.Cluster, namespace string, rules ...admregv1.RuleWithOperations) (func() error, error) {
	svcName := "webhook-kong-integration-tests"
	svcCloser, err := ensureWebhookService(ctx, cluster, namespace, svcName)
	if err != nil {
		return nil, err
	}

	fail := admregv1.Fail
	none := admregv1.SideEffectClassNone
	webhook, err := cluster.Client().AdmissionregistrationV1().ValidatingWebhookConfigurations().Create(ctx,
		&admregv1.ValidatingWebhookConfiguration{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "admissionregistration.k8s.io/v1",
				Kind:       "ValidatingWebhookConfiguration",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "kong-integration-tests",
			},
			Webhooks: []admregv1.ValidatingWebhook{
				{
					Name:                    "validations.kong.konghq.com",
					FailurePolicy:           &fail,
					SideEffects:             &none,
					AdmissionReviewVersions: []string{"v1beta1", "v1"},
					Rules:                   rules,
					ClientConfig: admregv1.WebhookClientConfig{
						Service: &admregv1.ServiceReference{
							Namespace: namespace,
							Name:      svcName,
						},
						CABundle: []byte(KongSystemServiceCert),
					},
				},
			},
		}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	closer := func() error {
		if err := cluster.Client().AdmissionregistrationV1().ValidatingWebhookConfigurations().Delete(ctx, webhook.Name, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			return err
		}
		return svcCloser()
	}

	return closer, waitForWebhookService()
}

// -----------------------------------------------------------------------------
// Test Utilities - Admission Webhooks - Private Functions
// -----------------------------------------------------------------------------

func ensureWebhookService(ctx context.Context, cluster clusters.Cluster, namespace, name string) (func() error, error) {
	validationsService, err := cluster.Client().CoreV1().Services(namespace).Create(ctx, &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "default",
					Port:       443,
					TargetPort: intstr.FromInt(AdmissionWebhookListenPort),
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("creating webhook service: %w", err)
	}

	nodeName := "aaaa"
	endpoints, err := cluster.Client().CoreV1().Endpoints(namespace).Create(ctx, &corev1.Endpoints{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Endpoints"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP:       AdmissionWebhookListenHost,
						NodeName: &nodeName,
					},
				},
				Ports: []corev1.EndpointPort{
					{
						Name:     "default",
						Port:     AdmissionWebhookListenPort,
						Protocol: corev1.ProtocolTCP,
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("creating webhook endpoints: %w", err)
	}

	closer := func() error {
		if err := cluster.Client().CoreV1().Services(namespace).Delete(ctx, validationsService.Name, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			return err
		}

		if err := cluster.Client().CoreV1().Endpoints(namespace).Delete(ctx, endpoints.Name, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	return closer, nil
}

func waitForWebhookService() error {
	timeout := time.Now().Add(webhookTimeout)
	consistentReplies := 0
	for timeout.After(time.Now()) {
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", AdmissionWebhookListenHost, AdmissionWebhookListenPort), 1*time.Second)
		if err == nil {
			// if we've gotten at least 3 success replies in a row, consider the webhook service ready
			consistentReplies++
			if consistentReplies >= 3 {
				return nil
			}
			continue // if we haven't quite reached 3 successes in a row, keep trying
		}
		consistentReplies = 0 // any failures reset the consistent replies check
		time.Sleep(time.Second)
	}
	return fmt.Errorf("admission webhook did not become responsive within %s", webhookTimeout)
}
