/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package status

import (
	"context"
	"os"
	"testing"
	"time"

	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	"github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress/task"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
)

func buildLoadBalancerIngressByIP() []apiv1.LoadBalancerIngress {
	return []apiv1.LoadBalancerIngress{
		{
			IP:       "10.0.0.1",
			Hostname: "foo1",
		},
		{
			IP:       "10.0.0.2",
			Hostname: "foo2",
		},
		{
			IP:       "10.0.0.3",
			Hostname: "",
		},
		{
			IP:       "",
			Hostname: "foo4",
		},
	}
}

func buildSimpleClientSet(extraObjects ...runtime.Object) *testclient.Clientset {
	objects := []runtime.Object{
		&apiv1.PodList{Items: []apiv1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo1",
					Namespace: apiv1.NamespaceDefault,
					Labels: map[string]string{
						"lable_sig": "foo_pod",
					},
				},
				Spec: apiv1.PodSpec{
					NodeName: "foo_node_2",
				},
				Status: apiv1.PodStatus{
					Phase: apiv1.PodRunning,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo1-unknown",
					Namespace: apiv1.NamespaceDefault,
				},
				Spec: apiv1.PodSpec{
					NodeName: "foo_node_1",
				},
				Status: apiv1.PodStatus{
					Phase: apiv1.PodUnknown,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo2",
					Namespace: apiv1.NamespaceDefault,
					Labels: map[string]string{
						"lable_sig": "foo_no",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo3",
					Namespace: metav1.NamespaceSystem,
					Labels: map[string]string{
						"lable_sig": "foo_pod",
					},
				},
				Spec: apiv1.PodSpec{
					NodeName: "foo_node_2",
				},
				Status: apiv1.PodStatus{
					Phase: apiv1.PodRunning,
				},
			},
		}},
		&apiv1.ServiceList{Items: []apiv1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: apiv1.NamespaceDefault,
				},
				Status: apiv1.ServiceStatus{
					LoadBalancer: apiv1.LoadBalancerStatus{
						Ingress: buildLoadBalancerIngressByIP(),
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo_non_exist",
					Namespace: apiv1.NamespaceDefault,
				},
			},
		}},
		&apiv1.NodeList{Items: []apiv1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo_node_1",
				},
				Status: apiv1.NodeStatus{
					Addresses: []apiv1.NodeAddress{
						{
							Type:    apiv1.NodeInternalIP,
							Address: "10.0.0.1",
						}, {
							Type:    apiv1.NodeExternalIP,
							Address: "10.0.0.2",
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo_node_2",
				},
				Status: apiv1.NodeStatus{
					Addresses: []apiv1.NodeAddress{
						{
							Type:    apiv1.NodeInternalIP,
							Address: "11.0.0.1",
						},
						{
							Type:    apiv1.NodeExternalIP,
							Address: "11.0.0.2",
						},
					},
				},
			},
		}},
		&apiv1.EndpointsList{Items: []apiv1.Endpoints{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ingress-controller-leader",
					Namespace: apiv1.NamespaceDefault,
					SelfLink:  "/api/v1/namespaces/default/endpoints/ingress-controller-leader",
				},
			}}},
	}

	return testclient.NewSimpleClientset(append(objects, extraObjects...)...)
}

func fakeSynFn(interface{}) error {
	return nil
}

func buildIngressesV1beta1() []networkingv1beta1.Ingress {
	return []networkingv1beta1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo_ingress_1",
				Namespace: apiv1.NamespaceDefault,
			},
			Status: networkingv1beta1.IngressStatus{
				LoadBalancer: apiv1.LoadBalancerStatus{
					Ingress: []apiv1.LoadBalancerIngress{
						{
							IP:       "10.0.0.1",
							Hostname: "foo1",
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo_ingress_different_class",
				Namespace: metav1.NamespaceDefault,
				Annotations: map[string]string{
					"kubernetes.io/ingress.class": "no-nginx",
				},
			},
			Status: networkingv1beta1.IngressStatus{
				LoadBalancer: apiv1.LoadBalancerStatus{
					Ingress: []apiv1.LoadBalancerIngress{
						{
							IP:       "0.0.0.0",
							Hostname: "foo.bar.com",
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo_ingress_2",
				Namespace: apiv1.NamespaceDefault,
			},
			Status: networkingv1beta1.IngressStatus{
				LoadBalancer: apiv1.LoadBalancerStatus{
					Ingress: []apiv1.LoadBalancerIngress{},
				},
			},
		},
	}
}

func buildIngressesV1() []networkingv1.Ingress {
	return []networkingv1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo_ingress_1",
				Namespace: apiv1.NamespaceDefault,
			},
			Status: networkingv1.IngressStatus{
				LoadBalancer: apiv1.LoadBalancerStatus{
					Ingress: []apiv1.LoadBalancerIngress{
						{
							IP:       "10.0.0.1",
							Hostname: "foo1",
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo_ingress_different_class",
				Namespace: metav1.NamespaceDefault,
				Annotations: map[string]string{
					"kubernetes.io/ingress.class": "no-nginx",
				},
			},
			Status: networkingv1.IngressStatus{
				LoadBalancer: apiv1.LoadBalancerStatus{
					Ingress: []apiv1.LoadBalancerIngress{
						{
							IP:       "0.0.0.0",
							Hostname: "foo.bar.com",
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo_ingress_2",
				Namespace: apiv1.NamespaceDefault,
			},
			Status: networkingv1.IngressStatus{
				LoadBalancer: apiv1.LoadBalancerStatus{
					Ingress: []apiv1.LoadBalancerIngress{},
				},
			},
		},
	}
}

var sampleIngressesV1beta1 = []*networkingv1beta1.Ingress{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo_ingress_non_01",
			Namespace: apiv1.NamespaceDefault,
		},
	},

	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo_ingress_1",
			Namespace: apiv1.NamespaceDefault,
		},
		Status: networkingv1beta1.IngressStatus{
			LoadBalancer: apiv1.LoadBalancerStatus{
				Ingress: buildLoadBalancerIngressByIP(),
			},
		},
	},
}

var sampleIngressesV1 = []*networkingv1.Ingress{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo_ingress_non_01",
			Namespace: apiv1.NamespaceDefault,
		},
	},

	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo_ingress_1",
			Namespace: apiv1.NamespaceDefault,
		},
		Status: networkingv1.IngressStatus{
			LoadBalancer: apiv1.LoadBalancerStatus{
				Ingress: buildLoadBalancerIngressByIP(),
			},
		},
	},
}

type testIngressLister struct {
	ingressesV1beta1 []*networkingv1beta1.Ingress
	ingressesV1      []*networkingv1.Ingress
}

func (til *testIngressLister) ListIngressesV1beta1() []*networkingv1beta1.Ingress {
	return til.ingressesV1beta1
}

func (til *testIngressLister) ListIngressesV1() []*networkingv1.Ingress {
	return til.ingressesV1
}

func (til *testIngressLister) ListTCPIngresses() ([]*configurationv1beta1.TCPIngress, error) {
	return nil, nil
}

func (til *testIngressLister) ListKnativeIngresses() ([]*knative.Ingress, error) {
	return nil, nil
}

func buildStatusSync() statusSync {
	return statusSync{
		pod: &util.PodInfo{
			Name:      "foo_base_pod",
			Namespace: apiv1.NamespaceDefault,
			Labels: map[string]string{
				"lable_sig": "foo_pod",
			},
		},
		syncQueue: task.NewTaskQueue(fakeSynFn, logrus.New()),
		Config: Config{
			CoreClient:     buildSimpleClientSet(&networkingv1beta1.IngressList{Items: buildIngressesV1beta1()}),
			PublishService: apiv1.NamespaceDefault + "/" + "foo",
			IngressLister:  &testIngressLister{ingressesV1beta1: sampleIngressesV1beta1},
			IngressAPI:     util.ExtensionsV1beta1,
		},
	}
}

func TestStatusActionsV1beta1(t *testing.T) {
	ctx := context.Background()
	// make sure election can be created
	os.Setenv("POD_NAME", "foo1")
	os.Setenv("POD_NAMESPACE", apiv1.NamespaceDefault)
	c := Config{
		CoreClient:             buildSimpleClientSet(&networkingv1beta1.IngressList{Items: buildIngressesV1beta1()}),
		PublishService:         apiv1.NamespaceDefault + "/" + "foo",
		IngressLister:          &testIngressLister{ingressesV1beta1: sampleIngressesV1beta1},
		UpdateStatusOnShutdown: true,
		IngressAPI:             util.NetworkingV1beta1,
		Logger:                 logrus.New(),
	}
	// create object
	fkSync, err := NewStatusSyncer(ctx, c)
	if fkSync == nil {
		t.Fatalf("expected a valid Sync")
	}

	fk := fkSync.(statusSync)

	// start it and wait for the election and syn actions
	go fk.Run()
	//  wait for the election
	time.Sleep(100 * time.Millisecond)
	// execute sync
	err = fk.sync("just-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newIPs := []apiv1.LoadBalancerIngress{{
		IP: "11.0.0.2",
	}}
	fooIngress1, err1 := fk.CoreClient.NetworkingV1beta1().Ingresses(
		apiv1.NamespaceDefault).Get(ctx, "foo_ingress_1", metav1.GetOptions{})
	if err1 != nil {
		t.Fatalf("unexpected error")
	}
	fooIngress1CurIPs := fooIngress1.Status.LoadBalancer.Ingress
	if !ingressSliceEqual(fooIngress1CurIPs, newIPs) {
		t.Fatalf("returned %v but expected %v", fooIngress1CurIPs, newIPs)
	}

	time.Sleep(1 * time.Second)

	// execute shutdown
	fk.Shutdown(true)
	// ingress should be empty
	newIPs2 := []apiv1.LoadBalancerIngress{}
	fooIngress2, err2 := fk.CoreClient.NetworkingV1beta1().Ingresses(
		apiv1.NamespaceDefault).Get(ctx, "foo_ingress_1", metav1.GetOptions{})
	if err2 != nil {
		t.Fatalf("unexpected error")
	}
	fooIngress2CurIPs := fooIngress2.Status.LoadBalancer.Ingress
	if !ingressSliceEqual(fooIngress2CurIPs, newIPs2) {
		t.Fatalf("returned %v but expected %v", fooIngress2CurIPs, newIPs2)
	}

	oic, err := fk.CoreClient.NetworkingV1beta1().Ingresses(
		metav1.NamespaceDefault).Get(ctx, "foo_ingress_different_class", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if oic.Status.LoadBalancer.Ingress[0].IP != "0.0.0.0" && oic.Status.LoadBalancer.Ingress[0].Hostname != "foo.bar.com" {
		t.Fatalf("invalid ingress status for rule with different class")
	}
}

func TestStatusActionsV1(t *testing.T) {
	ctx := context.Background()
	// make sure election can be created
	os.Setenv("POD_NAME", "foo1")
	os.Setenv("POD_NAMESPACE", apiv1.NamespaceDefault)
	c := Config{
		CoreClient:             buildSimpleClientSet(&networkingv1.IngressList{Items: buildIngressesV1()}),
		PublishService:         apiv1.NamespaceDefault + "/" + "foo",
		IngressLister:          &testIngressLister{ingressesV1: sampleIngressesV1},
		UpdateStatusOnShutdown: true,
		IngressAPI:             util.NetworkingV1,
		Logger:                 logrus.New(),
	}
	// create object
	fkSync, err := NewStatusSyncer(ctx, c)
	if fkSync == nil {
		t.Fatalf("expected a valid Sync")
	}

	fk := fkSync.(statusSync)

	// start it and wait for the election and syn actions
	go fk.Run()
	//  wait for the election
	time.Sleep(100 * time.Millisecond)
	// execute sync
	err = fk.sync("just-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newIPs := []apiv1.LoadBalancerIngress{{
		IP: "11.0.0.2",
	}}
	fooIngress1, err1 := fk.CoreClient.NetworkingV1().Ingresses(
		apiv1.NamespaceDefault).Get(ctx, "foo_ingress_1", metav1.GetOptions{})
	if err1 != nil {
		t.Fatalf("unexpected error")
	}
	fooIngress1CurIPs := fooIngress1.Status.LoadBalancer.Ingress
	if !ingressSliceEqual(fooIngress1CurIPs, newIPs) {
		t.Fatalf("returned %v but expected %v", fooIngress1CurIPs, newIPs)
	}

	time.Sleep(1 * time.Second)

	// execute shutdown
	fk.Shutdown(true)
	// ingress should be empty
	newIPs2 := []apiv1.LoadBalancerIngress{}
	fooIngress2, err2 := fk.CoreClient.NetworkingV1().Ingresses(
		apiv1.NamespaceDefault).Get(ctx, "foo_ingress_1", metav1.GetOptions{})
	if err2 != nil {
		t.Fatalf("unexpected error")
	}
	fooIngress2CurIPs := fooIngress2.Status.LoadBalancer.Ingress
	if !ingressSliceEqual(fooIngress2CurIPs, newIPs2) {
		t.Fatalf("returned %v but expected %v", fooIngress2CurIPs, newIPs2)
	}

	oic, err := fk.CoreClient.NetworkingV1().Ingresses(
		metav1.NamespaceDefault).Get(ctx, "foo_ingress_different_class", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if oic.Status.LoadBalancer.Ingress[0].IP != "0.0.0.0" && oic.Status.LoadBalancer.Ingress[0].Hostname != "foo.bar.com" {
		t.Fatalf("invalid ingress status for rule with different class")
	}
}

func TestCallback(t *testing.T) {
	buildStatusSync()
}

func TestKeyfunc(t *testing.T) {
	fk := buildStatusSync()

	i := "foo_base_pod"
	r, err := fk.keyfunc(i)

	if err != nil {
		t.Fatalf("unexpected error")
	}
	if r != i {
		t.Errorf("returned %v but expected %v", r, i)
	}
}

func TestRunningAddresessWithPublishService(t *testing.T) {
	ctx := context.Background()
	fk := buildStatusSync()

	r, _ := fk.runningAddresses(ctx)
	if r == nil {
		t.Fatalf("returned nil but expected valid []string")
	}
	rl := len(r)
	if len(r) != 1 {
		t.Errorf("returned %v but expected %v", rl, 1)
	}
}

func TestRunningAddresessWithPods(t *testing.T) {
	ctx := context.Background()
	fk := buildStatusSync()
	fk.PublishService = apiv1.NamespaceDefault + "/" + "foo"

	r, err := fk.runningAddresses(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatalf("returned nil but expected valid []string")
	}
	rl := len(r)
	if len(r) != 1 {
		t.Fatalf("returned %v but expected %v", rl, 1)
	}
	rv := r[0]
	if rv != "11.0.0.2" {
		t.Errorf("returned %v but expected %v", rv, "11.0.0.2")
	}
}

func TestRunningAddresessWithPublishStatusAddress(t *testing.T) {
	ctx := context.Background()
	fk := buildStatusSync()
	fk.PublishService = ""
	fk.PublishStatusAddress = "127.0.0.1"

	r, _ := fk.runningAddresses(ctx)
	if r == nil {
		t.Fatalf("returned nil but expected valid []string")
	}
	rl := len(r)
	if len(r) != 1 {
		t.Errorf("returned %v but expected %v", rl, 1)
	}
	rv := r[0]
	if rv != "127.0.0.1" {
		t.Errorf("returned %v but expected %v", rv, "127.0.0.1")
	}
}

func TestSliceToStatus(t *testing.T) {
	fkEndpoints := []string{
		"10.0.0.1",
		"2001:db8::68",
		"opensource-k8s-ingress",
	}

	r := sliceToStatus(fkEndpoints)

	if r == nil {
		t.Fatalf("returned nil but expected a valid []apiv1.LoadBalancerIngress")
	}
	rl := len(r)
	if rl != 3 {
		t.Fatalf("returned %v but expected %v", rl, 3)
	}
	re1 := r[0]
	if re1.Hostname != "opensource-k8s-ingress" {
		t.Fatalf("returned %v but expected %v", re1, apiv1.LoadBalancerIngress{Hostname: "opensource-k8s-ingress"})
	}
	re2 := r[1]
	if re2.IP != "10.0.0.1" {
		t.Fatalf("returned %v but expected %v", re2, apiv1.LoadBalancerIngress{IP: "10.0.0.1"})
	}
	re3 := r[2]
	if re3.IP != "2001:db8::68" {
		t.Fatalf("returned %v but expected %v", re3, apiv1.LoadBalancerIngress{IP: "2001:db8::68"})
	}
}

func TestIngressSliceEqual(t *testing.T) {
	fk1 := buildLoadBalancerIngressByIP()
	fk2 := append(buildLoadBalancerIngressByIP(), apiv1.LoadBalancerIngress{
		IP:       "10.0.0.5",
		Hostname: "foo5",
	})
	fk3 := buildLoadBalancerIngressByIP()
	fk3[0].Hostname = "foo_no_01"
	fk4 := buildLoadBalancerIngressByIP()
	fk4[2].IP = "11.0.0.3"

	fooTests := []struct {
		lhs []apiv1.LoadBalancerIngress
		rhs []apiv1.LoadBalancerIngress
		er  bool
	}{
		{fk1, fk1, true},
		{fk2, fk1, false},
		{fk3, fk1, false},
		{fk4, fk1, false},
		{fk1, nil, false},
		{nil, nil, true},
		{[]apiv1.LoadBalancerIngress{}, []apiv1.LoadBalancerIngress{}, true},
	}

	for _, fooTest := range fooTests {
		r := ingressSliceEqual(fooTest.lhs, fooTest.rhs)
		if r != fooTest.er {
			t.Errorf("returned %v but expected %v", r, fooTest.er)
		}
	}
}
