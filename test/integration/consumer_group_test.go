//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

func TestConsumerGroup(t *testing.T) {
	t.Parallel()

	RunWhenKongVersion(t, fmt.Sprintf(">=%s", versions.ConsumerGroupsVersionCutoff))
	RunWhenKongEnterprise(t)

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	consumerGroupNames := []string{
		"test-consumer-group-1",
		"test-consumer-group-2",
	}
	for _, cgName := range consumerGroupNames {
		t.Logf("configuring consumer group: %q", cgName)
		cg, err := c.ConfigurationV1beta1().KongConsumerGroups(ns.Name).Create(
			ctx,
			&kongv1beta1.KongConsumerGroup{
				ObjectMeta: metav1.ObjectMeta{
					Name: cgName,
					Annotations: map[string]string{
						annotations.IngressClassKey: consts.IngressClass,
					},
				},
			},
			metav1.CreateOptions{},
		)
		require.NoError(t, err)
		cleaner.Add(cg)
	}

	const consumerName = "test-consumer"
	t.Logf("configuring consumer: %q", consumerName)
	consumer, err := c.ConfigurationV1().KongConsumers(ns.Name).Create(
		ctx,
		&kongv1.KongConsumer{
			ObjectMeta: metav1.ObjectMeta{
				Name: consumerName,
				Annotations: map[string]string{
					annotations.IngressClassKey: consts.IngressClass,
				},
			},
			Username:       consumerName,
			ConsumerGroups: consumerGroupNames,
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)
	cleaner.Add(consumer)

	for _, cgName := range consumerGroupNames {
		t.Logf("validating that consumer %q was successfully added to previously configured consumer group: %q", consumerName, cgName)
		require.Eventually(t, func() bool {
			cgPath := fmt.Sprintf("/consumer_groups/%s/consumers/%s", cgName, consumerName)
			var headers map[string]string
			if testenv.DBMode() != testenv.DBModeOff {
				cgPath = fmt.Sprintf(
					"/%s/consumer_groups/%s/consumers/%s", consts.KongTestWorkspace, cgName, consumerName,
				)
				headers = map[string]string{
					"Kong-Admin-Token": consts.KongTestPassword,
				}
			}
			req := helpers.MustHTTPRequest(t, http.MethodGet, proxyAdminURL, cgPath, headers)
			resp, err := helpers.DefaultHTTPClient().Do(req)
			if err != nil {
				t.Logf("WARNING: error while waiting for %s: %v", req.URL, err)
				return false
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Logf("WARNING: error while reading response from %s: %v", req.URL, err)
			}
			switch resp.StatusCode {
			case http.StatusOK:
				return true
			case http.StatusForbidden:
				t.Logf(
					"WARNING: it seems Kong Gateway Enterprise hasn't got a valid license passed - from: %s received: %s with body: %s",
					req.URL, resp.Status, body,
				)
				return false
			default:
				t.Logf("WARNING: from: %s received unexpected: %s with body: %s", req.URL, resp.Status, body)
				return false
			}
		}, ingressWait, waitTick)
	}
}
