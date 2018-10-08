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

package store

import (
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eapache/channels"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	extensions "k8s.io/api/extensions/v1beta1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"github.com/kong/kubernetes-ingress-controller/test/e2e/framework"
)

func TestStore(t *testing.T) {
	t.Skip("Skipping tests in TestStore.")
	// TODO: find a way to avoid the need to use a real api server
	home := os.Getenv("HOME")
	kubeConfigFile := fmt.Sprintf("%v/.kube/config", home)
	kubeContext := ""

	kubeConfig, err := framework.LoadConfig(kubeConfigFile, kubeContext)
	if err != nil {
		t.Errorf("unexpected error loading kubeconfig file: %v", err)
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		t.Errorf("unexpected error creating ingress client: %v", err)
	}

	t.Run("should return an error searching for non existing objects", func(t *testing.T) {
		ns := createNamespace(clientSet, t)
		defer deleteNamespace(ns, clientSet, t)

		stopCh := make(chan struct{})
		updateCh := channels.NewRingChannel(1024)

		go func(ch *channels.RingChannel) {
			for {
				<-ch.Out()
			}
		}(updateCh)

		storer := New(ns.Name,
			fmt.Sprintf("%v/config", ns.Name),
			fmt.Sprintf("%v/tcp", ns.Name),
			fmt.Sprintf("%v/udp", ns.Name),
			"",
			10*time.Minute,
			clientSet,
			kubeConfig,
			updateCh)

		storer.Run(stopCh)

		key := fmt.Sprintf("%v/anything", ns.Name)
		ing, err := storer.GetIngress(key)
		if err == nil {
			t.Errorf("expected an error but none returned")
		}
		if ing != nil {
			t.Errorf("expected an Ingres but none returned")
		}

		s, err := storer.GetSecret(key)
		if err == nil {
			t.Errorf("expected an error but none returned")
		}
		if s != nil {
			t.Errorf("expected an Ingres but none returned")
		}

		svc, err := storer.GetService(key)
		if err == nil {
			t.Errorf("expected an error but none returned")
		}
		if svc != nil {
			t.Errorf("expected an Ingres but none returned")
		}

		updateCh.Close()
		close(stopCh)
	})

	t.Run("should return ingress one event for add, update and delete", func(t *testing.T) {
		ns := createNamespace(clientSet, t)
		defer deleteNamespace(ns, clientSet, t)

		stopCh := make(chan struct{})
		updateCh := channels.NewRingChannel(1024)

		var add uint64
		var upd uint64
		var del uint64

		go func(ch *channels.RingChannel) {
			for {
				evt, ok := <-ch.Out()
				if !ok {
					return
				}

				e := evt.(Event)
				if e.Obj == nil {
					continue
				}
				if _, ok := e.Obj.(*extensions.Ingress); !ok {
					t.Errorf("expected an Ingress type but %T returned", e.Obj)
				}
				switch e.Type {
				case CreateEvent:
					atomic.AddUint64(&add, 1)
				case UpdateEvent:
					atomic.AddUint64(&upd, 1)
				case DeleteEvent:
					atomic.AddUint64(&del, 1)
				}
			}
		}(updateCh)

		storer := New(ns.Name,
			fmt.Sprintf("%v/config", ns.Name),
			fmt.Sprintf("%v/tcp", ns.Name),
			fmt.Sprintf("%v/udp", ns.Name),
			"",
			10*time.Minute,
			clientSet,
			kubeConfig,
			updateCh)

		storer.Run(stopCh)

		ing, err := ensureIngress(&v1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dummy",
				Namespace: ns.Name,
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "dummy",
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "/",
										Backend: v1beta1.IngressBackend{
											ServiceName: "http-svc",
											ServicePort: intstr.FromInt(80),
										},
									},
								},
							},
						},
					},
				},
			},
		}, clientSet)
		if err != nil {
			t.Errorf("unexpected error creating ingress: %v", err)
		}

		// create an invalid ingress (different class)
		_, err = ensureIngress(&v1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "custom-class",
				Namespace: ns.Name,
				Annotations: map[string]string{
					"kubernetes.io/ingress.class": "something",
				},
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "dummy",
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "/",
										Backend: v1beta1.IngressBackend{
											ServiceName: "http-svc",
											ServicePort: intstr.FromInt(80),
										},
									},
								},
							},
						},
					},
				},
			},
		}, clientSet)
		if err != nil {
			t.Errorf("unexpected error creating ingress: %v", err)
		}

		ni := ing.DeepCopy()
		ni.Spec.Rules[0].Host = "update-dummy"
		_, err = ensureIngress(ni, clientSet)
		if err != nil {
			t.Errorf("unexpected error creating ingress: %v", err)
		}

		err = clientSet.ExtensionsV1beta1().
			Ingresses(ni.Namespace).
			Delete(ni.Name, &metav1.DeleteOptions{})
		if err != nil {
			t.Errorf("unexpected error creating ingress: %v", err)
		}

		framework.WaitForNoIngressInNamespace(clientSet, ni.Namespace, ni.Name)

		if atomic.LoadUint64(&add) != 1 {
			t.Errorf("expected 1 event of type Create but %v occurred", add)
		}
		if atomic.LoadUint64(&upd) != 1 {
			t.Errorf("expected 1 event of type Update but %v occurred", upd)
		}
		if atomic.LoadUint64(&del) != 1 {
			t.Errorf("expected 1 event of type Delete but %v occurred", del)
		}

		updateCh.Close()
		close(stopCh)
	})

	t.Run("should not receive events from secret not referenced from ingress", func(t *testing.T) {
		ns := createNamespace(clientSet, t)
		defer deleteNamespace(ns, clientSet, t)

		stopCh := make(chan struct{})
		updateCh := channels.NewRingChannel(1024)

		var add uint64
		var upd uint64
		var del uint64

		go func(ch *channels.RingChannel) {
			for {
				evt, ok := <-ch.Out()
				if !ok {
					return
				}

				e := evt.(Event)
				if e.Obj == nil {
					continue
				}
				switch e.Type {
				case CreateEvent:
					atomic.AddUint64(&add, 1)
				case UpdateEvent:
					atomic.AddUint64(&upd, 1)
				case DeleteEvent:
					atomic.AddUint64(&del, 1)
				}
			}
		}(updateCh)

		storer := New(ns.Name,
			fmt.Sprintf("%v/config", ns.Name),
			fmt.Sprintf("%v/tcp", ns.Name),
			fmt.Sprintf("%v/udp", ns.Name),
			"",
			10*time.Minute,
			clientSet,
			kubeConfig,
			updateCh)

		storer.Run(stopCh)

		secretName := "no-referenced"
		_, _, _, err = framework.CreateIngressTLSSecret(clientSet, []string{"foo"}, secretName, ns.Name)
		if err != nil {
			t.Errorf("unexpected error creating secret: %v", err)
		}

		err = framework.WaitForSecretInNamespace(clientSet, ns.Name, secretName)
		if err != nil {
			t.Errorf("unexpected error waiting for secret: %v", err)
		}

		if atomic.LoadUint64(&add) != 0 {
			t.Errorf("expected 0 events of type Create but %v occurred", add)
		}
		if atomic.LoadUint64(&upd) != 0 {
			t.Errorf("expected 0 events of type Update but %v occurred", upd)
		}
		if atomic.LoadUint64(&del) != 0 {
			t.Errorf("expected 0 events of type Delete but %v occurred", del)
		}

		err = clientSet.CoreV1().Secrets(ns.Name).Delete(secretName, &metav1.DeleteOptions{})
		if err != nil {
			t.Errorf("unexpected error deleting secret: %v", err)
		}

		time.Sleep(1 * time.Second)

		if atomic.LoadUint64(&add) != 0 {
			t.Errorf("expected 0 events of type Create but %v occurred", add)
		}
		if atomic.LoadUint64(&upd) != 0 {
			t.Errorf("expected 0 events of type Update but %v occurred", upd)
		}
		if atomic.LoadUint64(&del) != 1 {
			t.Errorf("expected 1 events of type Delete but %v occurred", del)
		}

		updateCh.Close()
		close(stopCh)
	})

	t.Run("should create an ingress with a secret it doesn't exists", func(t *testing.T) {
		ns := createNamespace(clientSet, t)
		defer deleteNamespace(ns, clientSet, t)

		stopCh := make(chan struct{})
		updateCh := channels.NewRingChannel(1024)

		var add uint64
		var upd uint64
		var del uint64

		go func(ch *channels.RingChannel) {
			for {
				evt, ok := <-ch.Out()
				if !ok {
					return
				}

				e := evt.(Event)
				if e.Obj == nil {
					continue
				}
				switch e.Type {
				case CreateEvent:
					atomic.AddUint64(&add, 1)
				case UpdateEvent:
					atomic.AddUint64(&upd, 1)
				case DeleteEvent:
					atomic.AddUint64(&del, 1)
				}
			}
		}(updateCh)

		storer := New(ns.Name,
			fmt.Sprintf("%v/config", ns.Name),
			fmt.Sprintf("%v/tcp", ns.Name),
			fmt.Sprintf("%v/udp", ns.Name),
			"",
			10*time.Minute,
			clientSet,
			kubeConfig,
			updateCh)

		storer.Run(stopCh)

		name := "ingress-with-secret"
		secretHosts := []string{name}

		_, err := ensureIngress(&v1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns.Name,
			},
			Spec: v1beta1.IngressSpec{
				TLS: []v1beta1.IngressTLS{
					{
						Hosts:      secretHosts,
						SecretName: name,
					},
				},
				Rules: []v1beta1.IngressRule{
					{
						Host: name,
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "/",
										Backend: v1beta1.IngressBackend{
											ServiceName: "http-svc",
											ServicePort: intstr.FromInt(80),
										},
									},
								},
							},
						},
					},
				},
			},
		}, clientSet)
		if err != nil {
			t.Errorf("unexpected error creating ingress: %v", err)
		}

		err = framework.WaitForIngressInNamespace(clientSet, ns.Name, name)
		if err != nil {
			t.Errorf("unexpected error waiting for secret: %v", err)
		}

		if atomic.LoadUint64(&add) != 1 {
			t.Errorf("expected 1 events of type Create but %v occurred", add)
		}
		if atomic.LoadUint64(&upd) != 0 {
			t.Errorf("expected 0 events of type Update but %v occurred", upd)
		}
		if atomic.LoadUint64(&del) != 0 {
			t.Errorf("expected 0 events of type Delete but %v occurred", del)
		}

		_, _, _, err = framework.CreateIngressTLSSecret(clientSet, secretHosts, name, ns.Name)
		if err != nil {
			t.Errorf("unexpected error creating secret: %v", err)
		}

		// skip closing updateCh.input channel to skip syncing with informer-started go routines
		// updateCh.Close()
		close(stopCh)
	})

	// test add ingress with secret it doesn't exists and then add secret
	// check secret is generated on fs
	// check ocsp
	// check invalid secret (missing crt)
	// check invalid secret (missing key)
	// check invalid secret (missing ca)
}

func createNamespace(clientSet *kubernetes.Clientset, t *testing.T) *apiv1.Namespace {
	t.Log("creating temporal namespace")
	ns, err := framework.CreateKubeNamespace("store-test", clientSet)
	if err != nil {
		t.Errorf("unexpected error creating ingress client: %v", err)
	}
	t.Logf("temporal namespace %v created", ns.Name)

	return ns
}

func deleteNamespace(ns *apiv1.Namespace, clientSet *kubernetes.Clientset, t *testing.T) {
	t.Logf("deleting temporal namespace %v created", ns.Name)
	err := framework.DeleteKubeNamespace(clientSet, ns.Name)
	if err != nil {
		t.Errorf("unexpected error creating ingress client: %v", err)
	}
	t.Logf("temporal namespace %v deleted", ns.Name)
}

func ensureIngress(ingress *extensions.Ingress, clientSet *kubernetes.Clientset) (*extensions.Ingress, error) {
	s, err := clientSet.ExtensionsV1beta1().Ingresses(ingress.Namespace).Update(ingress)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return clientSet.ExtensionsV1beta1().Ingresses(ingress.Namespace).Create(ingress)
		}
		return nil, err
	}
	return s, nil
}
