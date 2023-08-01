//go:build integration_tests

package integration

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

const webhookName = "kong-validations-consumer-group"

func TestValidationWebhookConsumerGroupKongOSS(t *testing.T) {
	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("webhook tests are only available on KIND clusters currently")
	}
	RunWhenKongVersion(t, fmt.Sprintf(">=%s", versions.ConsumerGroupsVersionCutoff))
	RunWhenKongOSS(t)

	t.Parallel()

	ctx := context.Background()
	ns, kongClient := setupWebhook(ctx, t, env, webhookName)

	t.Log("verifying validation for consumer groups for Kong OSS (rejecting creation)")
	require.ErrorContains(
		t,
		createRandomConsumerGroup(ctx, ns, kongClient),
		"consumer groups are not supported in Kong OSS (only in Enterprise)",
	)
}

func TestValidationWebhookConsumerGroupKongEnterprise(t *testing.T) {
	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("webhook tests are only available on KIND clusters currently")
	}
	RunWhenKongVersion(t, fmt.Sprintf(">=%s", versions.ConsumerGroupsVersionCutoff))
	RunWhenKongEnterprise(t)

	t.Parallel()

	ctx := context.Background()
	ns, kongClient := setupWebhook(ctx, t, env, webhookName)

	t.Log("verifying validation for consumer groups for Kong Enterprise with valid license")
	require.NoError(t, createRandomConsumerGroup(ctx, ns, kongClient))

	t.Log("make license invalid for Kong Enterprise")
	k8sClient := env.Cluster().Client()
	// Fill in the license secret with an invalid license.
	_, err := k8sClient.CoreV1().Secrets(consts.ControllerNamespace).Patch(
		ctx,
		"kong-enterprise-license",
		k8stypes.StrategicMergePatchType,
		[]byte(fmt.Sprintf(`{"data": {"license": "%s"}}`, base64.StdEncoding.EncodeToString([]byte("invalid")))),
		metav1.PatchOptions{},
	)
	require.NoError(t, err)
	// License secret is passed as env var to the Kong Gateway, so it has to be restarted.
	restartAnnotation := []byte(fmt.Sprintf(
		`{"spec": {"template": {"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`, time.Now().Format(time.Stamp),
	))
	deploymentClient := k8sClient.AppsV1().Deployments(consts.ControllerNamespace)
	const gateway = "ingress-controller-kong"
	deploymentPatch, err := deploymentClient.Patch(ctx, gateway, k8stypes.StrategicMergePatchType, restartAnnotation, metav1.PatchOptions{})
	require.NoError(t, err)

	require.Eventuallyf(t, func() bool {
		current, err := deploymentClient.Get(ctx, gateway, metav1.GetOptions{})
		if err != nil {
			t.Log("WARNING: unexpected error getting deployment: ", err)
			return false
		}
		return deploymentComplete(deploymentPatch, &current.Status)
	}, ingressWait, waitTick, "deployment %q (Kong Gateway) should be ready after restart", gateway)
	helpers.EventuallyGETPath(t, proxyAdminURL, "/", 200, "", map[string]string{
		"Kong-Admin-Token": consts.KongTestPassword,
	}, ingressWait, waitTick)
	time.Sleep(20 * time.Second) // WHY!!!???? (I want to get rid of this sleep)
	t.Log("verifying validation for consumer groups for Kong Enterprise with invalid license (rejecting creation)")
	require.ErrorContains(
		t,
		createRandomConsumerGroup(ctx, ns, kongClient),
		"consumer groups are not supported in Kong Enterprise runs without a license",
	)
}

func setupWebhook(ctx context.Context, t *testing.T, env environments.Environment, name string) (string, *clientset.Clientset) {
	ns := helpers.Namespace(ctx, t, env)
	closer, err := ensureAdmissionRegistration(
		ctx,
		name,
		[]admregv1.RuleWithOperations{
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"configuration.konghq.com"},
					APIVersions: []string{"v1beta1"},
					Resources:   []string{"kongconsumergroups"},
				},
				Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
			},
		},
	)
	assert.NoError(t, err, "creating webhook config")
	t.Cleanup(func() {
		assert.NoError(t, closer())
	})

	err = waitForWebhookServiceConnective(ctx, name)
	require.NoError(t, err)

	kongClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	return ns.Name, kongClient
}

func createRandomConsumerGroup(ctx context.Context, namespace string, kongClient *clientset.Clientset) error {
	_, err := kongClient.ConfigurationV1beta1().KongConsumerGroups(namespace).Create(ctx, &kongv1beta1.KongConsumerGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
	}, metav1.CreateOptions{})
	return err
}

// deploymentComplete considers a deployment to be complete once all of its desired replicas
// are updated and available, and no old pods are running.
// src: https://github.com/kubernetes/kubernetes/blob/513da69f76f64e5292ee661e033fb9f33ec89161/pkg/controller/deployment/util/deployment_util.go#L706-L713
func deploymentComplete(deployment *appsv1.Deployment, newStatus *appsv1.DeploymentStatus) bool {
	return newStatus.UpdatedReplicas == *(deployment.Spec.Replicas) &&
		newStatus.Replicas == *(deployment.Spec.Replicas) &&
		newStatus.AvailableReplicas == *(deployment.Spec.Replicas) &&
		newStatus.ObservedGeneration >= deployment.Generation
}
