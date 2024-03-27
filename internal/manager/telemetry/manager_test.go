package telemetry

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/kubernetes-telemetry/pkg/forwarders"
	"github.com/kong/kubernetes-telemetry/pkg/serializers"
	"github.com/kong/kubernetes-telemetry/pkg/telemetry"
	"github.com/kong/kubernetes-telemetry/pkg/types"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	testdynclient "k8s.io/client-go/dynamic/fake"
	testk8sclient "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	// openShiftVersionPodNamespace is a namespace expected to contain pods whose environment includes OpenShift version
	// information.
	openShiftVersionPodNamespace = "openshift-apiserver-operator"
	// openShiftVersionPodApp is a value for the "app" label to select pods whose environment includes OpenShift version
	// information.
	openShiftVersionPodApp = "openshift-apiserver-operator"
)

type mockGatewaysCounter int

func (m mockGatewaysCounter) GatewayClientsCount() int {
	return int(m)
}

func TestCreateManager(t *testing.T) {
	var (
		payload = types.ProviderReport{
			"db": "off",
			"kv": "3.1.1",
		}
		featureGates = map[string]bool{
			"gateway": true,
			"knative": false,
		}
		publishService = k8stypes.NamespacedName{
			Namespace: "kong",
			Name:      "kong-proxy",
		}
		pod = k8stypes.NamespacedName{
			Namespace: "kong",
			Name:      "kong-ingress-controller",
		}
	)
	t.Setenv("POD_NAME", pod.Name)
	t.Setenv("POD_NAMESPACE", pod.Namespace)

	hostname, err := os.Hostname()
	require.NoError(t, err)

	scheme := prepareScheme(t)
	objs := prepareObjects(pod)

	dyn := testdynclient.NewSimpleDynamicClient(scheme, objs...)
	ctrlClient := fakeclient.NewClientBuilder().
		WithScheme(scheme).
		// We need this for mesh detection which lists services.
		WithIndex(&corev1.Service{}, "metadata.name", func(o client.Object) []string {
			return []string{o.GetName()}
		}).
		WithRuntimeObjects(objs...).
		Build()

	k8sclient := testk8sclient.NewSimpleClientset(objs...)

	d, ok := k8sclient.Discovery().(*fakediscovery.FakeDiscovery)
	require.True(t, ok)
	d.FakedServerVersion = &version.Info{
		Major:        "1",
		Minor:        "24",
		GitVersion:   "v1.24.5",
		GitCommit:    "cc6a1b4915a99f49f5510ef0667f94b9ca832a8a",
		GitTreeState: "clean",
		BuildDate:    "2022-06-09T18:24:04Z",
		GoVersion:    "go1.16.15",
		Compiler:     "gc",
		Platform:     "linux/amd64",
	}

	reportValues := ReportValues{
		FeatureGates:                   featureGates,
		MeshDetection:                  true,
		PublishServiceNN:               publishService,
		KonnectSyncEnabled:             true,
		GatewayServiceDiscoveryEnabled: true,
	}

	runManagerTest(
		t,
		k8sclient,
		dyn,
		ctrlClient,
		mockGatewaysCounter(5),
		payload,
		reportValues,
		func(t *testing.T, actualReport string) {
			require.Equal(t,
				fmt.Sprintf(
					"<14>"+
						"signal=test-signal;"+
						"db=off;"+
						"feature-gateway-service-discovery=true;"+
						"feature-gateway=true;"+
						"feature-knative=false;"+
						"feature-konnect-sync=true;"+
						"hn=%s;"+
						"kv=3.1.1;"+
						"uptime=0;"+
						"discovered_gateways_count=5;"+
						"k8s_arch=linux/amd64;"+
						"k8s_provider=UNKNOWN;"+
						"k8sv=v1.24.5;"+
						"k8sv_semver=v1.24.5;"+
						"openshift_version=4.13.0;"+
						"k8s_nodes_count=4;"+
						"k8s_pods_count=11;"+
						"k8s_services_count=18;"+
						"kinm=c3,l2,l3,l4;"+
						"mdep=i3,k3,km3,l3,t3;"+
						"mdist=all18,c1,i2,k1,km1,l3,t1;"+
						"\n",
					hostname,
				),
				actualReport,
			)
		},
	)
}

func TestCreateManager_GatewayDiscoverySpecifics(t *testing.T) {
	testCases := []struct {
		name                           string
		gatewayServiceDiscoveryEnabled bool
		expectReportToContain          []string
		expectReportToNotContain       []string
	}{
		{
			name:                           "gateway service discovery disabled",
			gatewayServiceDiscoveryEnabled: false,
			expectReportToContain: []string{
				"feature-gateway-service-discovery=false",
			},
			expectReportToNotContain: []string{
				"discovered_gateways_count=",
			},
		},
		{
			name:                           "gateway service discovery enabled",
			gatewayServiceDiscoveryEnabled: true,
			expectReportToContain: []string{
				"feature-gateway-service-discovery=true",
				"discovered_gateways_count=5",
			},
		},
	}

	scheme := prepareScheme(t)
	dyn := testdynclient.NewSimpleDynamicClient(scheme)
	ctrlClient := fakeclient.NewClientBuilder().Build()
	k8sclient := testk8sclient.NewSimpleClientset()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			runManagerTest(
				t,
				k8sclient,
				dyn,
				ctrlClient,
				mockGatewaysCounter(5),
				Payload{},
				ReportValues{
					GatewayServiceDiscoveryEnabled: tc.gatewayServiceDiscoveryEnabled,
				},
				func(t *testing.T, actualReport string) {
					for _, expected := range tc.expectReportToContain {
						require.Contains(t, actualReport, expected)
					}
					for _, expected := range tc.expectReportToNotContain {
						require.NotContains(t, actualReport, expected)
					}
				})
		})
	}
}

// runManagerTest is a helper function that creates a manager with the dependencies provided as arguments, and
// calls the testFn with the actual report string it receives after triggering `test-signal` execution.
func runManagerTest(
	t *testing.T,

	// Following arguments map with the arguments of the createManager function, but instead of interfaces,
	// concrete test types are used.
	k8sclient *testk8sclient.Clientset,
	dyn *testdynclient.FakeDynamicClient,
	ctrlClient client.Client,
	gatewaysCounter mockGatewaysCounter,
	payload Payload,
	reportValues ReportValues,

	// testFn is a function that will be called with the actual report string.
	testFn func(t *testing.T, actualReport string),
) {
	ctx := context.Background()
	mgr, err := createManager(
		k8sclient,
		dyn,
		ctrlClient,
		gatewaysCounter,
		payload,
		reportValues,
		telemetry.OptManagerPeriod(time.Hour),
		telemetry.OptManagerLogger(logr.Discard()),
	)
	require.NoError(t, err)
	require.NotNil(t, mgr)

	ch := make(chan []byte)
	consumer := telemetry.NewConsumer(
		serializers.NewSemicolonDelimited(),
		forwarders.NewChannelForwarder(ch),
	)
	require.NoError(t, mgr.AddConsumer(consumer))

	require.NoError(t, mgr.Start())
	defer mgr.Stop()
	require.NoError(t, mgr.TriggerExecute(ctx, "test-signal"))
	select {
	case b := <-ch:
		testFn(t, string(b))
	case <-time.After(time.Second):
		t.Fatal("we should get a report but we didn't")
	}
}

func prepareScheme(t *testing.T) *runtime.Scheme {
	scheme := runtime.NewScheme()
	require.NoError(t, testk8sclient.AddToScheme(scheme))
	// Note: this has no effect on the object listing because pluralising gateways
	// does not work.
	// Ref: https://github.com/kubernetes/kubernetes/pull/110053
	require.NoError(t, gatewayv1.Install(scheme))
	return scheme
}

func prepareObjects(pod k8stypes.NamespacedName) []runtime.Object {
	setRandomUUIDName := func(o client.Object) runtime.Object {
		o.SetName(uuid.NewString())
		return o
	}

	return []runtime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: pod.Namespace,
				Name:      pod.Name,
				Annotations: map[string]string{
					"linkerd.io/proxy-version": "1.0.0",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "ingress-controller",
						Image: "kong/kubernetes-ingress-controller:2.8",
					},
					{
						Name:  "kong",
						Image: "kong/kong:3.1.1",
					},
					// sidecars
					{Name: "linkerd-proxy"},
					{Name: "envoy-sidecar"},
				},
				// init containers.
				InitContainers: []corev1.Container{
					{Name: "linkerd-init"},
				},
			},
		},
		// service with no endpoints.
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: pod.Namespace,
				Name:      "kong-proxy",
			},
		},
		// Service with multiple EndpointSlices.
		&discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: pod.Namespace,
				Name:      "kong-proxy-1",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "kong-proxy",
				},
			},
			Endpoints: []discoveryv1.Endpoint{
				{
					TargetRef: &corev1.ObjectReference{
						Kind:      "Pod",
						Namespace: pod.Namespace,
						Name:      pod.Name,
					},
				},
			},
		},

		setRandomUUIDName(&corev1.Node{}),
		setRandomUUIDName(&corev1.Node{}),
		setRandomUUIDName(&corev1.Node{}),
		setRandomUUIDName(&corev1.Node{}),

		setRandomUUIDName(&corev1.Service{}),
		setRandomUUIDName(&corev1.Service{}),
		setRandomUUIDName(&corev1.Service{}),
		setRandomUUIDName(&corev1.Service{}),
		setRandomUUIDName(&corev1.Service{}),

		setRandomUUIDName(&corev1.Pod{}),
		setRandomUUIDName(&corev1.Pod{}),
		setRandomUUIDName(&corev1.Pod{}),

		// services.
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "kuma-control-plane",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "kong-mesh-control-plane",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "linkerd-proxy-injector",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "traefik-mesh-controller",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "istiod",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "service1",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "service2",
			},
		},
		&corev1.Service{
			// service with no available endpoints.
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "service3",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service1",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service2",
				Annotations: map[string]string{
					"mesh.traefik.io/traffic-type": "TCP",
				},
			},
		},
		// Service with no EndpointSlices.
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service3",
			},
		},
		// Service with multiple EndpointSlices.
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service4",
			},
		},
		// EndpointSlices for Pods.
		&discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "service1-1",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "service1",
				},
			},
			Endpoints: []discoveryv1.Endpoint{
				{
					TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns1", Name: "pod1"},
				},
			},
		},
		&discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "service2-1",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "service2",
				},
			},
			Endpoints: []discoveryv1.Endpoint{
				{
					TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns1", Name: "pod2"},
				},
			},
		},
		&discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "service3-1",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "service3",
				},
			},
			// EndpointSlice with no endpoints.
		},
		&discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service1-1",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "service1",
				},
			},
			Endpoints: []discoveryv1.Endpoint{
				{
					TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns2", Name: "pod1"},
				},
			},
		},
		&discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service2-1",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "service2",
				},
			},
			Endpoints: []discoveryv1.Endpoint{
				{
					TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns2", Name: "pod2"},
				},
				{},
			},
		},
		// Two EndpointSlices for the same service.
		&discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service3-1",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "service3",
				},
			},
			Endpoints: []discoveryv1.Endpoint{
				{
					TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns2", Name: "pod3-1"},
				},
			},
		},
		&discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service3-2",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "service3",
				},
			},
			Endpoints: []discoveryv1.Endpoint{
				{
					TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns2", Name: "pod3-2"},
				},
			},
		},
		// Pods referenced by EndpointSlices.
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "pod1",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "worker"},
					{Name: "istio-proxy"},
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "pod2",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "worker"},
					{Name: "kuma-sidecar"},
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "pod1",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "worker"},
					{Name: "istio-proxy"},
					{Name: "linkerd-proxy"},
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "pod2",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "worker"},
				},
			},
		},
		// One Pod has service mesh sidecar, the other doesn't.
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "pod3-1",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "worker"},
					{Name: "linkerd-proxy"},
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "pod3-2",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "worker"},
				},
			},
		},

		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: openShiftVersionPodNamespace,
			},
			Spec: corev1.NamespaceSpec{},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: openShiftVersionPodNamespace,
				Name:      openShiftVersionPodApp + "-85c4c6dbb7-zbrkm",
				Labels: map[string]string{
					"app": openShiftVersionPodApp,
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "worker",
						Env: []corev1.EnvVar{
							{
								// this is hardcoded here. the upstream const in telemetry lives inside an internal package and
								// doesn't make sense to stick elsewhere. this is, however, only in a test, and it will be obvious if
								// we break it. not ideal, but only requires some additional busywork to fix if we change upstream.
								Name:  "OPERATOR_IMAGE_VERSION",
								Value: "4.13.0",
							},
						},
					},
				},
			},
		},
	}
}
