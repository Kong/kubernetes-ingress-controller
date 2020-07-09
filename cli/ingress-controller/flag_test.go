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

type TLSPair struct {
	Key, Cert string
}

var (
	tlsPairs = []TLSPair{
		{
			Cert: `-----BEGIN CERTIFICATE-----
MIIC2DCCAcACCQC32eFOsWpKojANBgkqhkiG9w0BAQsFADAuMRcwFQYDVQQDDA5z
ZWN1cmUtZm9vLWJhcjETMBEGA1UECgwKa29uZ2hxLm9yZzAeFw0xODEyMTgyMTI4
MDBaFw0xOTEyMTgyMTI4MDBaMC4xFzAVBgNVBAMMDnNlY3VyZS1mb28tYmFyMRMw
EQYDVQQKDAprb25naHEub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAqhl/HSwV6PbMv+cMFU9X+HuM7QbNNPh39GKa4pkxzFgiAnuuJ4jw9V/bzsEy
S+ZIyjzo+QKB1LzmgdcX4vkdI22BjxUd9HPHdZxtv3XilbNmSk9UOl2Hh1fORJoS
7YH+VbvVwiz5lo7qKRepbg/jcKkbs6AUE0YWFygtDLTvhP2qkphQkxZ0m8qroW91
CWgI73Ar6U2W/YQBRI3+LwtsKo0p2ASDijvqxElQBgBIiyGIr0RZc5pkCJ1eQdDB
2F6XaMfpeEyBj0MxypNL4S9HHfchOt55J1KOzYnUPkQnSoxp6oEjef4Q/ZCj5BRL
EGZnTb3tbwzHZCxGtgl9KqO9pQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAKQ5BX
kkBL+alERL31hsOgWgRiUMw+sPDtRS96ozUlPtVvAg9XFdpY4ldtWkxFcmBnhKzp
UewjrHkf9rR16NISwUTjlGIwaJu/ACQrY15v+r301Crq2DV+GjiUJFVuT495dp/l
0LZbt2Sh/uD+r3UNTcJpJ7jb1V0UP7FWXFj8oafsoFSgmxAPjpKQySTC54JK4AYb
QSnWu1nQLyohnrB9qLZhe2+jOQZnkKuCcWJQ5njvU6SxT3SOKE5XaOZCezEQ6IVL
U47YCCXsq+7wKWXBhKl4H2Ztk6x3HOC56l0noXWezsMfrou/kjwGuuViGnrjqelS
WQ7uVeNCUBY+l+qY
-----END CERTIFICATE-----`,
			Key: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCqGX8dLBXo9sy/
5wwVT1f4e4ztBs00+Hf0YprimTHMWCICe64niPD1X9vOwTJL5kjKPOj5AoHUvOaB
1xfi+R0jbYGPFR30c8d1nG2/deKVs2ZKT1Q6XYeHV85EmhLtgf5Vu9XCLPmWjuop
F6luD+NwqRuzoBQTRhYXKC0MtO+E/aqSmFCTFnSbyquhb3UJaAjvcCvpTZb9hAFE
jf4vC2wqjSnYBIOKO+rESVAGAEiLIYivRFlzmmQInV5B0MHYXpdox+l4TIGPQzHK
k0vhL0cd9yE63nknUo7NidQ+RCdKjGnqgSN5/hD9kKPkFEsQZmdNve1vDMdkLEa2
CX0qo72lAgMBAAECggEADxMTYNJ3Xp4Ap0EioQDXGv5YDul7ZiZe+xmCAHLzJtjo
qq+rT3WjZRuJr1kPzAosiT+8pdTDDMdw5jDZvRO2sV0TDksgzHk2RAYI897OpdWw
SwWcwU9oo2X0sb+1zbang5GR8BNsSxt/RQUDzu05itJx0gltvgeIDaVR2L5wO6ja
USa8OVuj/92XtIIve9OtyK9jAzgR6LQOTFrCCEv89/vmy5Bykv4Uz8s8swZmTs3v
XJmAmruHGuSLMfXk8lBRp/gVyNTi3uMsdph5AJbVKnra5TZLguEozZKbLdNUYk0p
+aAc7rxDcH2sPqa/7DwRvei9dvd5oB3VJlxGVgC8AQKBgQDfznRSSKAD15hoSDzt
cKNyhLgWAL+MD0jhHKUy3x+Z9OCvf0DVnmru5HfQKq5UfT0t8VTRPGKmOtAMD4cf
LYjIurvMvpVzQGSJfhtHQuULZTh3dfsM7xivMqSV+9txklMAakM7vGQlOQxhrScM
21Mp5LWDU6+e2pFCrQPop0IPkQKBgQDCkVE+dou2yFuJx3uytCH1yKPSy9tkdhQH
dGF12B5dq8MZZozAz5P9YN/COa9WjsNKDqWbEgLEksEQUq4t8SBjHnSV/D3x7rEF
qgwii0GETYxax6gms8nueIqWZQf+0NbX7Gc5mTqeVb7v3TrhsKr0VNMFRXXQwP2E
M/pxJq8q1QKBgQC3rH7oXLP+Ez0AMHDYSL3LKULOw/RvpMeh/9lQA6+ysTaIsP3r
kuSdhCEUVULXEiVYhBug0FcBp3jAvSmem8cLPb0Mjkim2mzoLfeDJ1JEZODPoaLU
fZEbj4tlj9oLvhOiXpMo/jaOGeCgdPN8aK86zXlt+wtBao0WVFnF4SalEQKBgQC1
uLfi2SGgs/0a8B/ORoO5ZY3s4c2lRMtsMvyb7iBeaIAuByPLKZUVABe89deXxnsL
fiaacPX41wBO2IoqCp2vNdC6DP9mKQNZQPtYgCvPAAbo+rVIgH9HpXn7AZ24FyGy
RfAbUcv3+in9KelGxZTF4zu8HqXtNXMSuOFeMT1FiQKBgF0R+IFDGHhD4nudAQvo
hncXsgyzK6QUzak6HmFji/CMZ6EU9q6A67JkiEWrYoKqIAKZ2Og8+Eucr/rDdGWc
kqlmLPBJAJeUsP/9KidBjTE5mIbn/2n089VPMBvnlt2xIcuB6+zrf2NjvlcZEyKS
Gn+T2uCyOP4a1DTUoPyoNJXo
-----END PRIVATE KEY-----`,
		},
	}
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
		AdmissionWebhookListen:   "off",
		AdmissionWebhookCertPath: "/admission-webhook/tls.crt",
		AdmissionWebhookKeyPath:  "/admission-webhook/tls.key",

		KongAdminURL:           "http://localhost:8001",
		KongAdminConcurrency:   10,
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

		LogLevel:  "info",
		LogFormat: "text",

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
		"--kong-admin-concurrency", "1",
		"--kong-workspace", "yolo",
		"--kong-admin-filter-tag", "foo-tag",
		"--admin-header", "foo:bar",
		"--kong-admin-token", "my-token",
		"--admin-tls-skip-verify",
		"--admin-tls-server-name", "kong-admin.example.com",
		"--admin-ca-cert-file", "/path/to/ca-cert",

		"--kong-custom-entities-secret", "foons/foosecretname",

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

		"--log-format", "json",

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
		KongAdminConcurrency:   1,
		KongWorkspace:          "yolo",
		KongAdminFilterTags:    []string{"foo-tag"},
		KongAdminHeaders:       []string{"foo:bar", "kong-admin-token:my-token"},
		KongAdminTLSSkipVerify: true,
		KongAdminTLSServerName: "kong-admin.example.com",
		KongAdminCACertPath:    "/path/to/ca-cert",

		KongCustomEntitiesSecret: "foons/foosecretname",

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

		LogLevel:  "info",
		LogFormat: "json",

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
		"CONTROLLER_KONG_ADMIN_CONCURRENCY":      "100",
		"CONTROLLER_KONG_ADMIN_TOKEN":            "my-secret-token",

		"CONTROLLER_LOG_LEVEL": "panic",

		"CONTROLLER_KONG_CUSTOM_ENTITIES_SECRET": "foons/barsecretname",
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
		KongAdminConcurrency:   100,
		KongWorkspace:          "",
		KongAdminHeaders:       []string{"kong-admin-token:my-secret-token"},
		KongAdminTLSSkipVerify: false,
		KongAdminTLSServerName: "",
		KongAdminCACertPath:    "",

		KongCustomEntitiesSecret: "foons/barsecretname",

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

		LogLevel:  "panic",
		LogFormat: "text",

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
		KongAdminConcurrency:   10,
		KongAdminFilterTags:    []string{"managed-by-ingress-controller"},
		KongAdminHeaders:       []string{"foo:bar"},
		KongAdminTLSSkipVerify: true,
		KongAdminTLSServerName: "kong-admin.example.com",
		KongAdminCACertPath:    "/path/to/ca-cert",

		AdmissionWebhookListen:   "off",
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

		LogLevel:  "info",
		LogFormat: "text",

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
		"--admission-webhook-listen", ":8080",
	}
	conf, err := parseFlags()

	expectedConf := cliConfig{
		KongAdminURL:           "http://kong.yolo42.com",
		KongWorkspace:          "yolo",
		KongAdminConcurrency:   10,
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

		LogLevel:  "info",
		LogFormat: "text",

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

	// different flags
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

// test the certificate environment variables
// these are mutually exclusive with their _FILE partners
// and aren't tested in the regular override test as such
func TestEnvironmentCertificates(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	assert := assert.New(t)

	envs := map[string]string{
		"CONTROLLER_ADMISSION_WEBHOOK_LISTEN": ":9001",
		"CONTROLLER_ADMISSION_WEBHOOK_CERT":   tlsPairs[0].Cert,
		"CONTROLLER_ADMISSION_WEBHOOK_KEY":    tlsPairs[0].Key,
		"CONTROLLER_KONG_ADMIN_CA_CERT":       tlsPairs[0].Cert,
	}
	for k, v := range envs {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	conf, err := parseFlags()

	expected := cliConfig{
		AdmissionWebhookListen:   ":9001",
		AdmissionWebhookCertPath: "/admission-webhook/tls.crt",
		AdmissionWebhookKeyPath:  "/admission-webhook/tls.key",
		AdmissionWebhookCert:     tlsPairs[0].Cert,
		AdmissionWebhookKey:      tlsPairs[0].Key,

		KongAdminCACert: tlsPairs[0].Cert,
	}
	assert.Nil(err, "unexpected error supplying certificates via environment")
	assert.Equal(expected.AdmissionWebhookCert, conf.AdmissionWebhookCert)
	assert.Equal(expected.AdmissionWebhookKey, conf.AdmissionWebhookKey)
	assert.Equal(expected.KongAdminCACert, conf.KongAdminCACert)
}
