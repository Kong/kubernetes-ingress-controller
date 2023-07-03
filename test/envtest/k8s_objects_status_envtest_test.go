//go:build envtest

package envtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func TestKubernetesObjectsStatus_NotBlocked(t *testing.T) {
	envcfg := Setup(t, scheme.Scheme)
	cfg := ConfigForEnvConfig(t, envcfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var client ctrlclient.Client
	{
		var err error
		client, err = ctrlclient.New(envcfg, ctrlclient.Options{
			Scheme: scheme.Scheme,
		})
		require.NoError(t, err)
	}

	ns := CreateNamespace(ctx, t, client)
	publishSvc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "publish-svc",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{Port: 80}},
		},
	}
	require.NoError(t, client.Create(ctx, &publishSvc))
	t.Cleanup(func() { _ = client.Delete(ctx, &publishSvc) })

	cfg.UpdateStatus = true
	cfg.UpdateStatusQueueBufferSize = 0
	cfg.PublishStatusAddress = []string{"127.0.0.1"}
	cfg.PublishService = mo.Some(k8stypes.NamespacedName{
		Namespace: publishSvc.Namespace,
		Name:      publishSvc.Name,
	})
	cfg.FeatureGates[featuregates.GatewayFeature] = true
	cfg.FeatureGates[featuregates.GatewayAlphaFeature] = true

	go func(ctx context.Context) {
		err := manager.Run(ctx, &cfg, util.ConfigDumpDiagnostic{}, logrus.New())
		require.NoError(t, err)
	}(ctx)

	gwc := gatewayv1beta1.GatewayClass{
		Spec: gatewayv1beta1.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				"konghq.com/gatewayclass-unmanaged": "placeholder",
			},
		},
	}

	require.NoError(t, client.Create(ctx, &gwc))
	t.Cleanup(func() { _ = client.Delete(ctx, &gwc) })

	gw := gatewayv1beta1.Gateway{
		Spec: gatewayv1beta1.GatewaySpec{
			GatewayClassName: gatewayv1beta1.ObjectName(gwc.Name),
			Listeners: []gatewayv1beta1.Listener{
				{
					Name:     gatewayv1beta1.SectionName("http"),
					Port:     gatewayv1beta1.PortNumber(80),
					Protocol: gatewayv1beta1.HTTPProtocolType,
					AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
						Namespaces: &gatewayv1beta1.RouteNamespaces{
							From: lo.ToPtr(gatewayv1beta1.NamespacesFromAll),
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
	require.NoError(t, client.Create(ctx, &gw))

	route := gatewayv1beta1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: "v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      uuid.NewString(),
		},
		Spec: gatewayv1beta1.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
				ParentRefs: []gatewayv1beta1.ParentReference{{
					Name:      gatewayv1beta1.ObjectName(gw.Name),
					Namespace: lo.ToPtr(gatewayv1beta1.Namespace(ns.Name)),
				}},
			},
			Rules: []gatewayv1beta1.HTTPRouteRule{{
				BackendRefs: builder.NewHTTPBackendRef("backend-1").WithNamespace(ns.Name).ToSlice(),
			}},
		},
	}
	require.NoError(t, client.Create(ctx, &route))

	time.Sleep(time.Second * 20)
}
