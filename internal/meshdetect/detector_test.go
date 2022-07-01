package meshdetect

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestDetectMeshDeployment(t *testing.T) {
	testScheme := runtime.NewScheme()
	err := corev1.AddToScheme(testScheme)
	require.NoErrorf(t, err, "should add corev1 to scheme successfully")

	b := fake.NewClientBuilder().WithScheme(testScheme)
	b.WithObjects(
		// add services.
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "istio-system",
				Name:      "istiod",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "kuma-sys",
				Name:      "kong-mesh-control-plane",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "traefik-mesh-controller",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "consul-mesh",
				Name:      "consul-server",
			},
		},
	)

	c := b.Build()
	d := &Detector{
		Client:       c,
		PodNamespace: "kong",
		PodName:      "kic-1",
		Logger:       logr.Discard(),
	}

	res := d.DetectMeshDeployment(context.Background())
	expected := map[MeshKind]*DeploymentResults{
		MeshKindIstio: {
			ServiceExists: true,
		},
		MeshKindLinkerd: {
			ServiceExists: false,
		},
		MeshKindKuma: {
			ServiceExists: false,
		},
		MeshKindKongMesh: {
			ServiceExists: true,
		},
		MeshKindConsul: {
			ServiceExists: true,
		},
		MeshKindTraefik: {
			ServiceExists: true,
		},
		MeshKindAWSAppMesh: {
			ServiceExists: false,
		},
	}

	for _, meshKind := range MeshesToDetect {
		require.Equalf(t, expected[meshKind], res[meshKind], "detection result should be the same for mesh %s", meshKind)
	}
}

func TestDetectRunUnder(t *testing.T) {
	b := fake.NewClientBuilder()
	b.WithObjects(
		// add KIC pod.
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "kong",
				Name:      "kong-ingress-controller",
				Annotations: map[string]string{
					"sidecar.istio.io/status":  "injected",
					"linkerd.io/proxy-version": "1.0.0",
					"kuma.io/sidecar-injected": "true",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					// KIC container
					{
						Name:  "ingress-controller",
						Image: "kong/kubernetes-ingress-controller:2.4",
					},
					// sidecars
					{Name: "istio-proxy"},
					{Name: "kuma-sidecar"},
					{
						Name:  "envoy",
						Image: "public.ecr.aws/appmesh/aws-appmesh-envoy:v1.22.2.0-prod",
					},
				},
				// init containers.
				InitContainers: []corev1.Container{
					{Name: "istio-init"},
				},
			},
		},
		// add kong-proxy service.
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "kong",
				Name:      "kong-proxy",
				Annotations: map[string]string{
					"mesh.traefik.io/traffic-type": "HTTP",
				},
			},
		},
		// add another namespace.
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "kong-2",
				Annotations: map[string]string{
					"linkerd.io/inject": "enabled",
				},
				Labels: map[string]string{
					"appmesh.k8s.aws/sidecarInjectorWebhook": "enabled",
				},
			},
		},
		// add another KIC pod and kong-proxy service.
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "kong-2",
				Name:      "kong-ingress-controller",
				Annotations: map[string]string{
					"linkerd.io/proxy-version":                   "1.0.0",
					"consul.hashicorp.com/connect-inject-status": "injected",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					// KIC container
					{
						Name:  "ingress-controller",
						Image: "kong/kubernetes-ingress-controller:2.4",
					},
					// sidecars
					{Name: "linkerd-proxy"},
					{Name: "envoy-sidecar"},
				},
				// init containers.
				InitContainers: []corev1.Container{
					{Name: "linkerd-init"},
					{Name: "consul-connect-inject-init"},
				},
			},
		},
		// add a KIC pod without a publishing service, and no injection.
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "kong-3",
				Name:      "kong-ingress-controller",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					// KIC container
					{
						Name:  "ingress-controller",
						Image: "kong/kubernetes-ingress-controller:2.4",
					},
				},
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "kong-2",
				Name:      "kong-proxy",
				Annotations: map[string]string{
					"mesh.traefik.io/traffic-type": "UDP",
				},
			},
		},
	)

	testCases := []struct {
		caseName        string
		podNamespace    string
		podName         string
		expectedResults map[MeshKind]*RunUnderResults
	}{
		{
			caseName:     "injected-istio,kuma,traefik,aws;annotation-linkerd",
			podNamespace: "kong",
			podName:      "kong-ingress-controller",
			expectedResults: map[MeshKind]*RunUnderResults{
				MeshKindIstio: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: true,
					InitContainerInjected:    true,
				},
				MeshKindLinkerd: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: false,
					InitContainerInjected:    false,
				},
				MeshKindKuma: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: true,
					InitContainerInjected:    false,
				},
				MeshKindKongMesh: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: true,
					InitContainerInjected:    false,
				},
				MeshKindConsul: {
					// all false
				},
				MeshKindTraefik: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: false,
					InitContainerInjected:    false,
				},
				MeshKindAWSAppMesh: {
					PodOrServiceAnnotation:   false,
					SidecarContainerInjected: true,
					InitContainerInjected:    false,
				},
			},
		},
		{
			caseName:     "injected-linkerd,consul",
			podNamespace: "kong-2",
			podName:      "kong-ingress-controller",
			expectedResults: map[MeshKind]*RunUnderResults{
				MeshKindIstio: {
					// all false
				},
				MeshKindLinkerd: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: true,
					InitContainerInjected:    true,
				},
				MeshKindKuma: {
					// all false
				},
				MeshKindKongMesh: {
					// all false
				},
				MeshKindConsul: {
					PodOrServiceAnnotation:   true,
					SidecarContainerInjected: true,
					InitContainerInjected:    true,
				},
				MeshKindTraefik: {
					// all false
				},
				MeshKindAWSAppMesh: {
					PodOrServiceAnnotation:   false,
					SidecarContainerInjected: false,
					InitContainerInjected:    false,
				},
			},
		},
		{
			caseName:     "nothing injected",
			podNamespace: "kong-3",
			podName:      "kong-ingress-controller",
			expectedResults: map[MeshKind]*RunUnderResults{
				// all mesh kinds -> all false
				MeshKindIstio:      {},
				MeshKindLinkerd:    {},
				MeshKindKuma:       {},
				MeshKindKongMesh:   {},
				MeshKindConsul:     {},
				MeshKindTraefik:    {},
				MeshKindAWSAppMesh: {},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			c := b.Build()
			d := &Detector{
				Client:       c,
				PodNamespace: tc.podNamespace,
				PodName:      tc.podName,
				Logger:       logr.Discard(),
			}
			d.PublishServiceName = tc.podNamespace + "/kong-proxy"
			res := d.DetectRunUnder(context.Background())
			for _, meshKind := range MeshesToDetect {
				require.Equalf(t, tc.expectedResults[meshKind], res[meshKind],
					"test case %s: detection result should be same for mesh %s", tc.caseName, meshKind)
			}
		})

	}
}

func TestDetectServiceDistribution(t *testing.T) {
	b := fake.NewClientBuilder()
	// add services/endpoints/pods.
	b.WithObjects(
		// services.
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
						{},
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
	)

	c := b.Build()
	d := &Detector{
		Client: c,
		Logger: logr.Discard(),
	}

	expectedTotal := 6
	expected := map[MeshKind]int{
		MeshKindIstio:    2,
		MeshKindLinkerd:  1,
		MeshKindKuma:     1,
		MeshKindKongMesh: 1,
		MeshKindTraefik:  1,
	}

	res, err := d.DetectServiceDistribution(context.Background())
	require.NoErrorf(t, err, "should not return error in detecting service distribution")
	require.Equalf(t, expectedTotal, res.TotalServices, "total number of services should be the same")
	for _, meshKind := range MeshesToDetect {
		require.Equalf(t, expected[meshKind], res.MeshDistribution[meshKind],
			"service within mesh %s should be same as expected", meshKind)
	}
}
