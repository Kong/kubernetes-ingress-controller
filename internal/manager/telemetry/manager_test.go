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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	testdynclient "k8s.io/client-go/dynamic/fake"
	testk8sclient "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

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
		ctx            = context.Background()
		meshDetection  = true
		publishService = apitypes.NamespacedName{
			Namespace: "kong",
			Name:      "kong-proxy",
		}
		pod = apitypes.NamespacedName{
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

	mgr, err := createManager(ctx, k8sclient, dyn, ctrlClient, payload, featureGates, meshDetection, publishService,
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
		require.Equal(t,
			fmt.Sprintf(
				"<14>"+
					"signal=test-signal;"+
					"db=off;"+
					"feature-gateway=true;"+
					"feature-knative=false;"+
					"hn=%s;"+
					"kv=3.1.1;"+
					"uptime=0;"+
					"k8s_arch=linux/amd64;"+
					"k8s_provider=UNKNOWN;"+
					"k8sv=v1.24.5;"+
					"k8sv_semver=v1.24.5;"+
					"k8s_nodes_count=4;"+
					"k8s_pods_count=8;"+
					"k8s_services_count=17;"+
					"kinm=c3,l2,l3,l4;"+
					"mdep=i3,k3,km3,l3,t3;"+
					"mdist=all17,c1,i2,k1,km1,l2,t1;"+
					"\n",
				hostname),
			string(b),
		)
	case <-time.After(time.Second):
		t.Fatal("we should get a report but we didn't")
	}
}

func prepareScheme(t *testing.T) *runtime.Scheme {
	scheme := runtime.NewScheme()
	require.NoError(t, testk8sclient.AddToScheme(scheme))
	// This doesn't work :(
	// https://kubernetes.slack.com/archives/C0EG7JC6T/p1676125039126589
	require.NoError(t, gatewayv1beta1.Install(scheme))
	return scheme
}

func prepareObjects(pod apitypes.NamespacedName) []runtime.Object {
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
		// endpoints.
		&corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: pod.Namespace,
				Name:      "kong-proxy",
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							TargetRef: &corev1.ObjectReference{
								Kind:      "Pod",
								Namespace: pod.Namespace,
								Name:      pod.Name,
							},
						},
					},
				},
			},
		},
		// &gatewayv1beta1.Gateway{},

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
		// service with no endpoints.
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service3",
			},
		},
		// endpoints.
		&corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "service1",
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns1", Name: "pod1"},
						},
					},
				},
			},
		},
		&corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "service2",
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns1", Name: "pod2"},
						},
					},
				},
			},
		},
		&corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns1",
				Name:      "service3",
			},
			// endpoints with no subsets.
		},
		&corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service1",
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns2", Name: "pod1"},
						},
					},
				},
			},
		},
		&corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns2",
				Name:      "service2",
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							TargetRef: &corev1.ObjectReference{Kind: "Pod", Namespace: "ns2", Name: "pod2"},
						},
					},
				},
			},
		},
		// pods.
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
	}
}
