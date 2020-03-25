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

	"github.com/eapache/channels"
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

	_, kubeClient, err := createApiserverClient("", kubeConfigFile)
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

	conf, err := parseFlags()
	if err != nil {
		t.Errorf("unexpected error creating Kong controller: %v", err)
	}

	kong, err := controller.NewKongController(
		&controller.Configuration{
			KubeClient: kubeClient,
		},
		channels.NewRingChannel(1024),
		store.New(store.CacheStores{}, conf.IngressClass),
	)

	exitCh := make(chan int, 1)
	go handleSigterm(kong, make(chan struct{}), exitCh)

	t.Logf("sending SIGTERM to process PID %v", syscall.Getpid())
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM); err != nil {
		t.Errorf("unexpected error sending SIGTERM signal")
	}

	// Allow test to time out if no value becomes avaialble soon enough.
	if code := <-exitCh; code != 1 {
		t.Errorf("expected exit code 1 but %v received", code)
	}
}
