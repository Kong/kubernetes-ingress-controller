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
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// resetForTesting clears all flag state and sets the usage function as directed.
// After calling resetForTesting, parse errors in flag handling will not
// exit the program.
// Extracted from https://github.com/golang/go/blob/master/src/flag/export_test.go
func resetForTesting(usage func()) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.Usage = usage
}

func TestDefaults(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{}
	assert := assert.New(t)

	conf, err := parseFlags()

	expectedConf := cliConfig{
		AdmissionWebhookListen:   ":8080",
		AdmissionWebhookCertPath: "/admission-webhook/tls.crt",
		AdmissionWebhookKeyPath:  "/admission-webhook/tls.key",

		KongAdminURL:           "http://localhost:8001",
		KongWorkspace:          "",
		KongAdminFilterTags:    []string{"managed-by-ingress-controller"},
		KongAdminHeaders:       []string{},
		KongAdminTLSSkipVerify: false,
		KongAdminTLSServerName: "",
		KongAdminCACertPath:    "",

		WatchNamespace: "",
		IngressClass:   "kong",
		ElectionID:     "ingress-controller-leader",

		PublishService:         "",
		PublishStatusAddress:   "",
		UpdateStatus:           true,
		UpdateStatusOnShutdown: true,

		SyncPeriod:    600 * time.Second,
		SyncRateLimit: 0.3,

		APIServerHost:      "",
		KubeConfigFilePath: "",

		EnableProfiling: true,

		ShowVersion:      false,
		AnonymousReports: true,
	}
	assert.Equal(expectedConf, conf)
	assert.Nil(err, "unexpected error parsing default flags")
}

func TestOverrideViaCLIFlags(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	assert := assert.New(t)

	os.Args = []string{
		"cmd",
		"--admission-webhook-listen", ":8081",
		"--admission-webhook-cert-file", "/cert-file",
		"--admission-webhook-key-file", "/key-file",

		"--kong-url", "https://kong.example.com",
		"--kong-workspace", "yolo",
		"--kong-admin-filter-tag", "foo-tag",
		"--admin-header", "foo:bar",
		"--admin-tls-skip-verify",
		"--admin-tls-server-name", "kong-admin.example.com",
		"--admin-ca-cert-file", "/path/to/ca-cert",

		"--watch-namespace", "foons",
		"--ingress-class", "kong-internal",
		"--election-id", "new-election-id",

		"--publish-service", "published-kong-proxy",
		"--publish-status-address", "some-custom-address",
		"--update-status=false",
		"--update-status-on-shutdown=false",

		"--sync-period", "10s",
		"--sync-rate-limit", "0.9",

		"--apiserver-host", "kube-apiserver.internal",
		"--kubeconfig", "/path/to/kubeconfig",

		"--profiling=false",
		"--version",
		"--anonymous-reports=false",
	}
	conf, err := parseFlags()

	expectedConf := cliConfig{
		AdmissionWebhookListen:   ":8081",
		AdmissionWebhookCertPath: "/cert-file",
		AdmissionWebhookKeyPath:  "/key-file",

		KongAdminURL:           "https://kong.example.com",
		KongWorkspace:          "yolo",
		KongAdminFilterTags:    []string{"foo-tag"},
		KongAdminHeaders:       []string{"foo:bar"},
		KongAdminTLSSkipVerify: true,
		KongAdminTLSServerName: "kong-admin.example.com",
		KongAdminCACertPath:    "/path/to/ca-cert",

		WatchNamespace: "foons",
		IngressClass:   "kong-internal",
		ElectionID:     "new-election-id",

		PublishService:         "published-kong-proxy",
		PublishStatusAddress:   "some-custom-address",
		UpdateStatus:           false,
		UpdateStatusOnShutdown: false,

		SyncPeriod:    10 * time.Second,
		SyncRateLimit: 0.9,

		APIServerHost:      "kube-apiserver.internal",
		KubeConfigFilePath: "/path/to/kubeconfig",

		EnableProfiling:  false,
		ShowVersion:      true,
		AnonymousReports: false,
	}
	assert.Equal(expectedConf, conf)
	assert.Nil(err, "unexpected error parsing default flags")
}

func TestOverrideViaEnvVars(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{}
	assert := assert.New(t)

	envs := map[string]string{
		"CONTROLLER_ADMISSION_WEBHOOK_LISTEN":    ":9001",
		"CONTROLLER_ADMISSION_WEBHOOK_CERT_FILE": "/new-cert-path",
		"CONTROLLER_ADMISSION_WEBHOOK_KEY_FILE":  "/new-key-path",
		"CONTROLLER_ANONYMOUS_REPORTS":           "false",
	}
	for k, v := range envs {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	conf, err := parseFlags()

	expectedConf := cliConfig{
		AdmissionWebhookListen:   ":9001",
		AdmissionWebhookCertPath: "/new-cert-path",
		AdmissionWebhookKeyPath:  "/new-key-path",

		KongAdminFilterTags:    []string{"managed-by-ingress-controller"},
		KongAdminURL:           "http://localhost:8001",
		KongWorkspace:          "",
		KongAdminHeaders:       []string{},
		KongAdminTLSSkipVerify: false,
		KongAdminTLSServerName: "",
		KongAdminCACertPath:    "",

		WatchNamespace: "",
		IngressClass:   "kong",
		ElectionID:     "ingress-controller-leader",

		PublishService:         "",
		PublishStatusAddress:   "",
		UpdateStatus:           true,
		UpdateStatusOnShutdown: true,

		SyncPeriod:    600 * time.Second,
		SyncRateLimit: 0.3,

		APIServerHost:      "",
		KubeConfigFilePath: "",

		EnableProfiling: true,

		ShowVersion:      false,
		AnonymousReports: false,
	}
	assert.Equal(expectedConf, conf)
	assert.Nil(err, "unexpected error parsing default flags")
}

func TestDeprecatedFlags(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	assert := assert.New(t)

	os.Args = []string{
		"cmd",
		"--kong-url", "https://kong.example.com",
		"--kong-workspace", "yolo",
		"--admin-header", "foo:bar",
		"--admin-tls-skip-verify",
		"--admin-tls-server-name", "kong-admin.example.com",
		"--admin-ca-cert-file", "/path/to/ca-cert",
	}
	conf, err := parseFlags()

	expectedConf := cliConfig{
		KongAdminURL:           "https://kong.example.com",
		KongWorkspace:          "yolo",
		KongAdminFilterTags:    []string{"managed-by-ingress-controller"},
		KongAdminHeaders:       []string{"foo:bar"},
		KongAdminTLSSkipVerify: true,
		KongAdminTLSServerName: "kong-admin.example.com",
		KongAdminCACertPath:    "/path/to/ca-cert",

		AdmissionWebhookListen:   ":8080",
		AdmissionWebhookCertPath: "/admission-webhook/tls.crt",
		AdmissionWebhookKeyPath:  "/admission-webhook/tls.key",

		WatchNamespace: "",
		IngressClass:   "kong",
		ElectionID:     "ingress-controller-leader",

		PublishService:         "",
		PublishStatusAddress:   "",
		UpdateStatus:           true,
		UpdateStatusOnShutdown: true,

		SyncPeriod:    600 * time.Second,
		SyncRateLimit: 0.3,

		APIServerHost:      "",
		KubeConfigFilePath: "",

		EnableProfiling: true,

		ShowVersion:      false,
		AnonymousReports: true,
	}
	assert.Equal(expectedConf, conf)
	assert.Nil(err, "unexpected error parsing default flags")
}

func TestDeprecatedFlagPrecedences(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	assert := assert.New(t)

	os.Args = []string{
		"cmd",
		"--kong-url", "https://kong.example.com",
		"--kong-admin-url", "http://kong.yolo42.com",
		"--kong-workspace", "yolo",
		"--admin-header", "foo:bar",
		"--kong-admin-header", "fuu:baz",
		"--kong-admin-tls-skip-verify",
		"--admin-tls-server-name", "kong-admin.example.com",
		"--kong-admin-tls-server-name", "kong-admin-new.example.com",
		"--admin-ca-cert-file", "/path/to/ca-cert",
		"--kong-admin-ca-cert-file", "/path/to/new/ca-cert",
	}
	conf, err := parseFlags()

	expectedConf := cliConfig{
		KongAdminURL:           "http://kong.yolo42.com",
		KongWorkspace:          "yolo",
		KongAdminFilterTags:    []string{"managed-by-ingress-controller"},
		KongAdminHeaders:       []string{"fuu:baz"},
		KongAdminTLSSkipVerify: true,
		KongAdminTLSServerName: "kong-admin-new.example.com",
		KongAdminCACertPath:    "/path/to/new/ca-cert",

		AdmissionWebhookListen:   ":8080",
		AdmissionWebhookCertPath: "/admission-webhook/tls.crt",
		AdmissionWebhookKeyPath:  "/admission-webhook/tls.key",

		WatchNamespace: "",
		IngressClass:   "kong",
		ElectionID:     "ingress-controller-leader",

		PublishService:         "",
		PublishStatusAddress:   "",
		UpdateStatus:           true,
		UpdateStatusOnShutdown: true,

		SyncPeriod:    600 * time.Second,
		SyncRateLimit: 0.3,

		APIServerHost:      "",
		KubeConfigFilePath: "",

		EnableProfiling: true,

		ShowVersion:      false,
		AnonymousReports: true,
	}
	assert.Equal(expectedConf, conf)
	assert.Nil(err, "unexpected error parsing default flags")
}

func TestKongAdminHeaders(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	assert := assert.New(t)

	os.Args = []string{
		"cmd",
		"--kong-admin-header", "key0:value0",
		"--kong-admin-header", "key1:value1",
	}
	conf, err := parseFlags()
	assert.Equal([]string{"key0:value0", "key1:value1"}, conf.KongAdminHeaders)

	assert.Nil(err, "unexpected error parsing default flags")
}

func TestKongAdminHeadersEnvVar(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	assert := assert.New(t)

	k := "CONTROLLER_KONG_ADMIN_HEADER"
	v := "key0:value0 key1:value1"
	os.Setenv(k, v)
	defer os.Unsetenv(k)
	conf, err := parseFlags()
	assert.Equal([]string{"key0:value0", "key1:value1"}, conf.KongAdminHeaders)

	assert.Nil(err, "unexpected error parsing default flags")
}

func TestKongFilterTags(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	assert := assert.New(t)

	// comma-separated
	os.Args = []string{
		"cmd",
		"--kong-admin-filter-tag", "foo,bar",
	}
	conf, err := parseFlags()
	assert.Equal([]string{"foo", "bar"}, conf.KongAdminFilterTags)

	assert.Nil(err, "unexpected error parsing default flags")

	resetForTesting(func() { t.Fatal("bad parse") })

	// differnt flags
	os.Args = []string{
		"cmd",
		"--kong-admin-filter-tag", "foo",
		"--kong-admin-filter-tag", "bar",
	}
	conf, err = parseFlags()
	assert.Equal([]string{"foo", "bar"}, conf.KongAdminFilterTags)

	assert.Nil(err, "unexpected error parsing default flags")
}

func TestKongAdminFilterTagEnvVar(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	assert := assert.New(t)

	k := "CONTROLLER_KONG_ADMIN_FILTER_TAG"
	v := "tag1 tag2"
	os.Setenv(k, v)
	defer os.Unsetenv(k)

	conf, err := parseFlags()
	assert.Equal([]string{"tag1", "tag2"},
		conf.KongAdminFilterTags)

	assert.Nil(err, "unexpected error parsing default flags")
}
