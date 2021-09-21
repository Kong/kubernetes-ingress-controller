//go:build integration_tests
// +build integration_tests

package integration

import (
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestEventHooks(t *testing.T) {
	if enterpriseEnablement != "on" {
		t.Logf("enterprise is not enabled. Skip Event Hooks Test.")
		t.Skip()
	}
	// configure event hooks for consumers using kongIngress
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("deploying a kong ingress to test webhooks.")
	testName := "event-hooks"

	webhook := "webhook"
	crud := "crud"
	consumers := "consumers"
	toBeCreatedWebhook := kongv1.KongIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},

		EventHooks: kong.EventHooks{
			Config: map[string]interface{}{
				"url":        "https://webhook.site/ec707ef0-ab91-4693-8dd2-114471ff6f90",
				"ssl_verify": false,
				"secret":     "",
			},
			Handler: &webhook,
			Source:  &crud,
			Event:   &consumers,
		},
	}

	_, err := c.ConfigurationV1().KongIngresses(ns.Name).Create(ctx, toBeCreatedWebhook, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that KongIngress %s is cleaned up", king.Name)
		if err := c.ConfigurationV1().KongIngresses(ns.Name).Delete(ctx, webhook.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	// create consumer through kIC
	consumer := &kongv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: ns.Name,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		"username":  "my-username",
		"custom_id": "my-custom-id",
		"tags":      {"user-level", "low-priority"},
	}

	c.ConfigurationV1().KongConsumer(testName).Create(ctx, consumer, metav1.CreateOptions{})
	assert.NoError(t, err)
	fmt.Printf("successfully created consumer.")

	// verify the event-hooks actions
	cmd := exec.CommandContext(ctx, "curl", "-i", "-X", "Get", "--url", "https://webhook.site/ec707ef0-ab91-4693-8dd2-114471ff6f90")
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		exitOnErr(fmt.Errorf("%s: %w", stderr.String(), err))
	}
}
