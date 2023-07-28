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

	const consumerGroupName = "test-consumer-group"
	t.Logf("configuring consumer group: %q", consumerGroupName)
	cg, err := c.ConfigurationV1beta1().KongConsumerGroups(ns.Name).Create(
		ctx,
		&kongv1beta1.KongConsumerGroup{
			ObjectMeta: metav1.ObjectMeta{
				Name:      consumerGroupName,
				Namespace: ns.Name,
				Annotations: map[string]string{
					annotations.IngressClassKey: consts.IngressClass,
				},
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)
	cleaner.Add(cg)

	t.Logf("validating that consumer group: %q was successfully configured", consumerGroupName)
	require.Eventually(t, func() bool {
		cgPath := fmt.Sprintf("/consumer_groups/%s", consumerGroupName)
		var headers map[string]string
		if testenv.DBMode() != testenv.DBModeOff {
			cgPath = fmt.Sprintf("/%s/consumer_groups/%s", consts.KongTestWorkspace, consumerGroupName)
			headers = map[string]string{
				"Kong-Admin-Token": consts.KongTestPassword,
			}
		}
		req := helpers.MustHTTPRequest(t, http.MethodGet, proxyAdminURL, cgPath, headers)
		resp, err := helpers.DefaultHTTPClient().Do(req)
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", resp.Request.URL, err)
			return false
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Logf("WARNING: error while reading response from %s: %v", resp.Request.URL, err)
		}
		switch resp.StatusCode {
		case http.StatusOK:
			return true
		case http.StatusForbidden:
			t.Logf(
				"WARNING: it seems Kong Gateway Enterprise hasn't got a valid license passed - from: %s received: %s with body: %s",
				resp.Request.URL, resp.Status, body,
			)
			return false
		default:
			t.Logf("WARNING: from: %s received unexpected: %s with body: %s", resp.Request.URL, resp.Status, body)
			return false
		}
	}, ingressWait, waitTick)
}
