//go:build envtest

package envtest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	gojson "github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

// configDumpResponse mirrors the diagnostics.configDumpResponse struct, which isn't published.
// It's replicated here since some envtests use the config dump endpoints as a hack to extract the config for
// inspection.

type configDumpResponse struct {
	ConfigHash string       `json:"hash"`
	Config     file.Content `json:"config"`
}

func TestIngressWorksWithServiceBackendsSpecifyingOnlyPortNames(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 10 * time.Second
		tickTime = 10 * time.Millisecond
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Add Gateway API Schema because its controllers are enabled by default.
	scheme := Scheme(t, WithGatewayAPI)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)
	ingressClassName := "kongenvtest"
	deployIngressClass(ctx, t, ingressClassName, ctrlClient)

	diagPort := helpers.GetFreePort(t)
	ns := CreateNamespace(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(ns.Name),
		WithIngressClass(ingressClassName),
		WithProxySyncSeconds(0.01),
		WithDiagnosticsServer(diagPort),
	)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Spec.Template.Spec.Containers[0].Ports[0].Name = "http"
	deployment.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, deployment))

	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.Namespace = ns.Name
	t.Logf("exposing deployment %s via service %s", deployment.Name, service.Name)
	require.NoError(t, ctrlClient.Create(ctx, service))

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-1",
			Namespace: ns.Name,
			Labels: map[string]string{
				"app": "httpbin",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				func() corev1.Container {
					c := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
					c.Ports[0].Name = "http"
					return c
				}(),
			},
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, &pod))

	es := discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      uuid.NewString(),
			Namespace: ns.Name,
			Labels: map[string]string{
				"kubernetes.io/service-name": service.Name,
			},
		},
		AddressType: discoveryv1.AddressTypeIPv4,
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"10.0.0.1"},
				Conditions: discoveryv1.EndpointConditions{
					Ready:       lo.ToPtr(true),
					Terminating: lo.ToPtr(false),
				},
				TargetRef: testPodReference("pod-1", ns.Name),
			},
		},
		Ports: builder.NewEndpointPort(80).WithName("http").IntoSlice(),
	}
	require.NoError(t, ctrlClient.Create(ctx, &es))

	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: ns.Name,
			Annotations: map[string]string{
				"konghq.com/strip-path": "true",
			},
		},
		Spec: netv1.IngressSpec{
			IngressClassName: lo.ToPtr(ingressClassName),
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: lo.ToPtr(netv1.PathTypePrefix),
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: service.Name,
											Port: netv1.ServiceBackendPort{
												Name: service.Spec.Ports[0].Name,
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
	}
	t.Logf("creating ingress %s for service %s", ingress.Name, service.Name)
	require.NoError(t, ctrlClient.Create(ctx, ingress))

	require.Eventually(t, func() bool {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/debug/config/successful", diagPort))
		if err != nil {
			t.Logf("WARNING: error while getting config: %v", err)
			return false
		}
		defer resp.Body.Close()

		var (
			configDump configDumpResponse
			config     file.Content
			buff       bytes.Buffer
		)

		if err := gojson.NewDecoder(io.TeeReader(resp.Body, &buff)).Decode(&configDump); err != nil {
			t.Logf("WARNING: error while decoding config: %+v, response: %s", err, buff.String())
			return false
		}

		config = configDump.Config

		if len(config.Services) != 1 {
			t.Logf("WARNING: expected 1 service in config: %+v", config)
			return false
		}
		if len(config.Services[0].Routes) != 1 {
			t.Logf("WARNING: expected 1 route for service in config: %+v", config)
			return false
		}

		if len(config.Upstreams) != 1 {
			t.Logf("WARNING: expected 1 upstream in config: %+v", config)
			return false
		}

		if len(config.Upstreams[0].Targets) != 1 {
			t.Logf("WARNING: expected 1 target in config: %+v", config)
			return false
		}

		target := config.Upstreams[0].Targets[0].Target.Target
		if target == nil || *target != "10.0.0.1:80" {
			t.Logf("WARNING: expected target to be equal to %s%d: actual %s", es.Endpoints[0].Addresses[0], *es.Ports[0].Port, *target)
			return false
		}

		return true
	}, waitTime, tickTime)
}
