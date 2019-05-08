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

package main

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/eapache/channels"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
)

func TestCreateApiserverClient(t *testing.T) {
	t.Skip("Skipping TestCreateApiserverClient.")
	home := os.Getenv("HOME")
	kubeConfigFile := fmt.Sprintf("%v/.kube/config", home)

	_, kubeClient, err := createApiserverClient("", kubeConfigFile)
	if err != nil {
		t.Fatalf("unexpected error creating api server client: %v", err)
	}
	if kubeClient == nil {
		t.Fatalf("expected a kubernetes client but none returned")
	}

	_, _, err = createApiserverClient("", "")
	if err == nil {
		t.Fatalf("expected an error creating api server client without an api server URL or kubeconfig file")
	}
}

func TestHandleSigterm(t *testing.T) {
	t.Skip("Skipping TestHandleSigterm.")
	home := os.Getenv("HOME")
	kubeConfigFile := fmt.Sprintf("%v/.kube/config", home)

	kubeConf, kubeClient, err := createApiserverClient("", kubeConfigFile)
	if err != nil {
		t.Fatalf("unexpected error creating api server client: %v", err)
	}

	resetForTesting(func() { t.Fatal("bad parse") })

	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "test")
	defer os.Setenv("POD_NAME", "")
	defer os.Setenv("POD_NAMESPACE", "")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "--default-backend-service", "ingress-nginx/default-backend-http"}

	_, conf, err := parseFlags()
	if err != nil {
		t.Errorf("unexpected error creating Kong controller: %v", err)
	}
	conf.KubeClient = kubeClient
	conf.KubeConf = kubeConf

	conf.Kong = controller.Kong{}

	kong, err := controller.NewKongController(conf,
		channels.NewRingChannel(1024),
		store.New(
			store.CacheStores{},
			annotations.IngressClassValidatorFuncFromObjectMeta(conf.IngressClass),
		),
	)

	go handleSigterm(kong, make(chan struct{}), func(code int) {
		if code != 1 {
			t.Errorf("expected exit code 1 but %v received", code)
		}

		return
	})

	time.Sleep(1 * time.Second)

	t.Logf("sending SIGTERM to process PID %v", syscall.Getpid())
	err = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	if err != nil {
		t.Errorf("unexpected error sending SIGTERM signal")
	}
}
