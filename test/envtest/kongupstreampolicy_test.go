//go:build envtest

package envtest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/go-logr/zapr"
	gojson "github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

func TestKongUpstreamPolicyWithoutGatewayAPICRDs(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme, WithInstallGatewayCRDs(false))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrlClient := NewControllerClient(t, scheme, envcfg)
	ingressClassName := "kongenvtest"
	deployIngressClass(ctx, t, ingressClassName, ctrlClient)

	logger := zapr.NewLogger(zap.NewNop())
	ctrl.SetLogger(logger)

	diagPort := helpers.GetFreePort(t)
	ns := CreateNamespace(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(ns.Name),
		WithIngressClass(ingressClassName),
		WithGatewayFeatureEnabled,
		WithKongServiceFacadeFeatureEnabled(),
		WithGatewayAPIControllers(),
		WithProxySyncSeconds(0.10),
		WithDiagnosticsServer(diagPort),
	)

	t.Log("verfying that KongUpstreamPolicy works without gateway APIs")

	t.Log("creating a KongUpstreamPolicy")
	const KongUpstreamPolicyName = "test-upstream-policy"
	kup := &kongv1beta1.KongUpstreamPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KongUpstreamPolicyName,
			Namespace: ns.Name,
		},
		Spec: kongv1beta1.KongUpstreamPolicySpec{
			Algorithm: lo.ToPtr("round-robin"),
			Slots:     lo.ToPtr(32),
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, kup))

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Spec.Template.Spec.Containers[0].Ports[0].Name = "http"
	deployment.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, deployment))

	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.Namespace = ns.Name
	service.Annotations = map[string]string{
		kongv1beta1.KongUpstreamPolicyAnnotationKey: KongUpstreamPolicyName,
	}
	t.Logf("exposing deployment %s via service %s", deployment.Name, service.Name)
	require.NoError(t, ctrlClient.Create(ctx, service))

	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: ns.Name,
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

	t.Logf("verify that upstream policy is configured in Kong gateway correctly")
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

		if len(config.Upstreams) != 1 {
			t.Logf("WARNING: expected 1 upstream in config: %+v", config)
			return false
		}
		upstream := config.Upstreams[0]
		return upstream.Algorithm != nil && *upstream.Algorithm == "round-robin" && upstream.Slots != nil && *upstream.Slots == 32
	}, waitTime, tickTime)

	t.Logf("verify that ancestor status of KongUpstreamPolicy is updated correctly")
	err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
		Namespace: ns.Name,
		Name:      KongUpstreamPolicyName,
	}, kup)
	require.NoError(t, err)
	require.Len(t, kup.Status.Ancestors, 1)
	require.Equal(t, "Service", string(*kup.Status.Ancestors[0].AncestorRef.Kind))
	require.Equal(t, ns.Name, string(*kup.Status.Ancestors[0].AncestorRef.Namespace))
	require.Equal(t, service.Name, string(kup.Status.Ancestors[0].AncestorRef.Name))
}

func TestKongUpstreamPolicyWithHTTPRoute(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithKong, WithGatewayAPI)
	envcfg := Setup(t, scheme)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrlClient := NewControllerClient(t, scheme, envcfg)
	ingressClassName := "kongenvtest"
	deployIngressClass(ctx, t, ingressClassName, ctrlClient)

	logger := zapr.NewLogger(zap.NewNop())
	ctrl.SetLogger(logger)

	diagPort := helpers.GetFreePort(t)
	ns := CreateNamespace(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(ns.Name),
		WithIngressClass(ingressClassName),
		WithGatewayFeatureEnabled,
		WithGatewayAPIControllers(),
		WithProxySyncSeconds(0.10),
		WithDiagnosticsServer(diagPort),
	)

	t.Log("verfying that KongUpstreamPolicy works without gateway APIs")

	t.Log("creating a KongUpstreamPolicy")
	const KongUpstreamPolicyName = "test-upstream-policy"
	kup := &kongv1beta1.KongUpstreamPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KongUpstreamPolicyName,
			Namespace: ns.Name,
		},
		Spec: kongv1beta1.KongUpstreamPolicySpec{
			Algorithm: lo.ToPtr("round-robin"),
			Slots:     lo.ToPtr(32),
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, kup))

	gwc := gatewayapi.GatewayClass{
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				"konghq.com/gatewayclass-unmanaged": "placeholder",
			},
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, &gwc))
	t.Cleanup(func() { _ = ctrlClient.Delete(ctx, &gwc) })

	gw := gatewayapi.Gateway{
		Spec: gatewayapi.GatewaySpec{
			GatewayClassName: gatewayapi.ObjectName(gwc.Name),
			Listeners: []gatewayapi.Listener{
				{
					Name:     gatewayapi.SectionName("http"),
					Port:     gatewayapi.PortNumber(80),
					Protocol: gatewayapi.HTTPProtocolType,
					AllowedRoutes: &gatewayapi.AllowedRoutes{
						Namespaces: &gatewayapi.RouteNamespaces{
							From: lo.ToPtr(gatewayapi.NamespacesFromAll),
						},
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      uuid.NewString(),
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, &gw))

	gwOld := gw.DeepCopy()
	gw.Status = gatewayapi.GatewayStatus{
		Addresses: []gatewayapi.GatewayStatusAddress{
			{
				Type:  lo.ToPtr(gatewayapi.IPAddressType),
				Value: "10.0.0.1",
			},
		},
		Conditions: []metav1.Condition{
			{
				Type:               "Programmed",
				Status:             metav1.ConditionTrue,
				Reason:             "Programmed",
				LastTransitionTime: metav1.Now(),
				ObservedGeneration: gw.Generation,
			},
			{
				Type:               "Accepted",
				Status:             metav1.ConditionTrue,
				Reason:             "Accepted",
				LastTransitionTime: metav1.Now(),
				ObservedGeneration: gw.Generation,
			},
		},
		Listeners: []gatewayapi.ListenerStatus{
			{
				Name: "http",
				Conditions: []metav1.Condition{
					{
						Type:               "Accepted",
						Status:             metav1.ConditionTrue,
						Reason:             "Accepted",
						LastTransitionTime: metav1.Now(),
					},
					{
						Type:               "Programmed",
						Status:             metav1.ConditionTrue,
						Reason:             "Programmed",
						LastTransitionTime: metav1.Now(),
					},
				},
				SupportedKinds: []gatewayapi.RouteGroupKind{
					{
						Group: lo.ToPtr(gatewayapi.Group(gatewayv1.GroupVersion.Group)),
						Kind:  "HTTPRoute",
					},
				},
			},
		},
	}
	require.NoError(t, ctrlClient.Status().Patch(ctx, &gw, ctrlclient.MergeFrom(gwOld)))

	t.Log("deploying a minimal HTTP container deployment to test HTTPRoutes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Spec.Template.Spec.Containers[0].Ports[0].Name = "http"
	deployment.Namespace = ns.Name
	require.NoError(t, ctrlClient.Create(ctx, deployment))

	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.Namespace = ns.Name
	service.Annotations = map[string]string{
		kongv1beta1.KongUpstreamPolicyAnnotationKey: KongUpstreamPolicyName,
	}
	t.Logf("exposing deployment %s via service %s", deployment.Name, service.Name)
	require.NoError(t, ctrlClient.Create(ctx, service))

	route := gatewayapi.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: "v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      uuid.NewString(),
		},
		Spec: gatewayapi.HTTPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name:      gatewayapi.ObjectName(gw.Name),
					Namespace: lo.ToPtr(gatewayapi.Namespace(ns.Name)),
				}},
			},
			Rules: []gatewayapi.HTTPRouteRule{{
				BackendRefs: builder.NewHTTPBackendRef(service.Name).WithNamespace(ns.Name).ToSlice(),
			}},
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, &route))

	t.Logf("verifying that the Service as backend of HTTPRoute is added to ancestor status of KongUpstreamPolicy")
	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
			Namespace: ns.Name,
			Name:      KongUpstreamPolicyName,
		}, kup)
		require.NoError(t, err)
		return lo.ContainsBy(kup.Status.Ancestors, func(ancestorStatus gatewayapi.PolicyAncestorStatus) bool {
			return ancestorStatus.AncestorRef.Kind != nil && string(*ancestorStatus.AncestorRef.Kind) == "Service" &&
				string(ancestorStatus.AncestorRef.Name) == service.Name
		})
	}, waitTime, tickTime)
}

func TestKongUpstreamPolicyNotReferencedInReconciledIngress(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme, WithInstallGatewayCRDs(false))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrlClient := NewControllerClient(t, scheme, envcfg)
	ingressClassName := "kongenvtest"
	deployIngressClass(ctx, t, ingressClassName, ctrlClient)
	alterIngressClassName := "kongenvtest-alter"
	deployIngressClass(ctx, t, alterIngressClassName, ctrlClient)

	logger := zapr.NewLogger(zap.NewNop())
	ctrl.SetLogger(logger)

	ns := CreateNamespace(ctx, t, ctrlClient)
	RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(ns.Name),
		WithIngressClass(ingressClassName),
		WithProxySyncSeconds(0.10),
	)

	t.Log("creating a KongUpstreamPolicy")
	const KongUpstreamPolicyName = "test-upstream-policy"
	kup := &kongv1beta1.KongUpstreamPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KongUpstreamPolicyName,
			Namespace: ns.Name,
		},
		Spec: kongv1beta1.KongUpstreamPolicySpec{
			Algorithm: lo.ToPtr("round-robin"),
			Slots:     lo.ToPtr(32),
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, kup))

	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Spec.Template.Spec.Containers[0].Ports[0].Name = "http"
	deployment.Namespace = ns.Name
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
	service.Namespace = ns.Name
	service.Annotations = map[string]string{
		kongv1beta1.KongUpstreamPolicyAnnotationKey: KongUpstreamPolicyName,
	}
	t.Logf("Creating a Service %s with the KongUpstreamPolicy and used as backend of Ingress", service.Name)
	require.NoError(t, ctrlClient.Create(ctx, service))

	// KongUpstreamPolicy should not be reconciled when no reconciled Ingresses/HTTPRoutes reference it.
	alterIngress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name + "-alter",
			Namespace: ns.Name,
		},
		Spec: netv1.IngressSpec{
			IngressClassName: lo.ToPtr(alterIngressClassName),
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
	t.Logf("creating ingress %s with another ingress class for service %s", alterIngress.Name, service.Name)
	require.NoError(t, ctrlClient.Create(ctx, alterIngress))

	t.Logf("verify that ancestor status of KongUpstreamPolicy is not updated when it is not referenced by reconciled Ingress")
	require.Never(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
			Namespace: ns.Name,
			Name:      KongUpstreamPolicyName,
		}, kup)
		require.NoError(t, err)
		return len(kup.Status.Ancestors) != 0
	}, waitTime, tickTime)

	// KongUpstreamPolicy should get reconciled when a reconciled Ingress reference it.
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: ns.Name,
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

	t.Logf("verifying that ancestor status of KongUpstreamPolicy is updated when it is references by reconciled Ingress")
	require.Eventually(t, func() bool {
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{
			Namespace: ns.Name,
			Name:      KongUpstreamPolicyName,
		}, kup)
		require.NoError(t, err)
		if len(kup.Status.Ancestors) != 1 {
			return false
		}
		ancestorRef := kup.Status.Ancestors[0].AncestorRef
		return string(*ancestorRef.Kind) == "Service" && string(*ancestorRef.Namespace) == ns.Name && string(ancestorRef.Name) == service.Name
	}, waitTime, tickTime)
}
