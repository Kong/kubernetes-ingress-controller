//go:build envtest
// +build envtest

package telemetry_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilversion "k8s.io/apimachinery/pkg/util/version"
	"k8s.io/apimachinery/pkg/version"
	discoveryclient "k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset/scheme"
	"github.com/kong/kubernetes-ingress-controller/v2/test/envtest"
)

func TestTelemetry(t *testing.T) {
	t.Log("configuring TLS listener - server for telemetry data")
	cert, err := generateSelfSignedCert()
	require.NoError(t, err)
	telemetryServer, err := tls.Listen("tcp", "localhost:0", &tls.Config{
		Certificates: []tls.Certificate{cert},
		// The same version as the one used by TLS forwarder in the pkg telemetry.
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
	})
	require.NoError(t, err)
	defer telemetryServer.Close()
	// Run a server that will receive the report, it's expected
	// to be the first connection and the payload.
	reportChan := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleConnectionToTelemetryServer(ctx, t, telemetryServer, reportChan)

	t.Log("configuring envtest and creating K8s objects for telemetry test")
	envcfg := envtest.Setup(t, scheme.Scheme)
	// Let's have long duration due too rate limiter in K8s client.
	cfg := configForEnvTestTelemetry(t, envcfg, telemetryServer.Addr().String(), time.Second)
	c, err := cfg.GetKubeconfig()
	require.NoError(t, err)
	createK8sObjectsForTelemetryTest(ctx, t, c)

	t.Log("starting the controller manager")
	go func() {
		deprecatedLogger, _, err := manager.SetupLoggers(&cfg, io.Discard)
		if !assert.NoError(t, err) {
			return
		}
		err = manager.Run(ctx, &cfg, util.ConfigDumpDiagnostic{}, deprecatedLogger)
		assert.NoError(t, err)
	}()

	dcl, err := discoveryclient.NewDiscoveryClientForConfig(envcfg)
	require.NoError(t, err)
	k8sVersion, err := dcl.ServerVersion()
	require.NoError(t, err)
	t.Log("verifying telemetry report")
	require.Eventually(t, func() bool {
		return verifyTelemetryReport(t, k8sVersion, string(<-reportChan))
	},
		10*time.Second,
		// Tick duration doesn't really matter, because reading from channel is blocking.
		100*time.Millisecond,
	)
}

func configForEnvTestTelemetry(t *testing.T, envcfg *rest.Config, splunkEndpoint string, telemetryPeriod time.Duration) manager.Config {
	t.Helper()

	cfg := manager.Config{}
	cfg.FlagSet() // Just set the defaults.

	// Telemetry is enabled by default so nothing to configure here.

	// Override the APIServer.
	cfg.APIServerHost = envcfg.Host
	cfg.APIServerCertData = envcfg.CertData
	cfg.APIServerKeyData = envcfg.KeyData
	cfg.APIServerCAData = envcfg.CAData
	cfg.KongAdminURLs = []string{envtest.StartAdminAPIServerMock(t).URL}
	cfg.UpdateStatus = false
	cfg.ProxySyncSeconds = 0.1

	// And other settings which are irrelevant.
	cfg.Konnect.ConfigSynchronizationEnabled = false
	cfg.Konnect.LicenseSynchronizationEnabled = false
	cfg.EnableProfiling = false
	cfg.EnableConfigDumps = false

	cfg.FeatureGates = featuregates.GetFeatureGatesDefaults()
	cfg.FeatureGates[featuregates.GatewayFeature] = false

	cfg.SplunkEndpoint = splunkEndpoint
	cfg.SplunkEndpointInsecureSkipVerify = true
	cfg.TelemetryPeriod = telemetryPeriod

	return cfg
}

func createK8sObjectsForTelemetryTest(ctx context.Context, t *testing.T, cfg *rest.Config) {
	t.Helper()
	cl, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	gcl, err := gatewayclient.NewForConfig(cfg)
	require.NoError(t, err)

	const additionalNamespace = "test-ns-1"
	_, err = cl.CoreV1().Namespaces().Create(
		ctx,
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: additionalNamespace},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)

	for i := 0; i < 2; i++ {
		_, err = gcl.GatewayV1beta1().GatewayClasses().Create(
			ctx,
			&gatewayv1beta1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
				Spec: gatewayv1beta1.GatewayClassSpec{
					ControllerName: "test.com/gateway-controller",
				},
			},
			metav1.CreateOptions{},
		)
		require.NoError(t, err)

		_, err = cl.CoreV1().Nodes().Create(
			ctx,
			&corev1.Node{
				ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
			},
			metav1.CreateOptions{},
		)
		require.NoError(t, err)

		for _, namespace := range []string{metav1.NamespaceDefault, additionalNamespace} {
			_, err := cl.CoreV1().Services(namespace).Create(
				ctx,
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Port: 443,
							},
						},
					},
				},
				metav1.CreateOptions{},
			)
			require.NoError(t, err)

			_, err = cl.CoreV1().Pods(namespace).Create(
				ctx,
				&corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{Name: "test", Image: "test"},
						},
					},
				},
				metav1.CreateOptions{},
			)
			require.NoError(t, err)

			_, err = gcl.GatewayV1beta1().Gateways(namespace).Create(
				ctx,
				&gatewayv1beta1.Gateway{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
					Spec: gatewayv1beta1.GatewaySpec{
						GatewayClassName: gatewayv1beta1.ObjectName("test"),
						Listeners: []gatewayv1beta1.Listener{
							{
								Name:     "test",
								Port:     443,
								Protocol: "HTTP",
							},
						},
					},
				},
				metav1.CreateOptions{},
			)
			require.NoError(t, err)

			_, err = gcl.GatewayV1beta1().HTTPRoutes(namespace).Create(
				ctx,
				&gatewayv1beta1.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
					Spec:       gatewayv1beta1.HTTPRouteSpec{},
				},
				metav1.CreateOptions{},
			)
			require.NoError(t, err)

			_, err = gcl.GatewayV1alpha2().GRPCRoutes(namespace).Create(
				ctx,
				&gatewayv1alpha2.GRPCRoute{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
					Spec:       gatewayv1alpha2.GRPCRouteSpec{},
				},
				metav1.CreateOptions{},
			)
			require.NoError(t, err)

			_, err = gcl.GatewayV1alpha2().TCPRoutes(namespace).Create(
				ctx,
				&gatewayv1alpha2.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
					Spec: gatewayv1alpha2.TCPRouteSpec{
						Rules: []gatewayv1alpha2.TCPRouteRule{{}},
					},
				},
				metav1.CreateOptions{},
			)
			require.NoError(t, err)

			_, err = gcl.GatewayV1alpha2().UDPRoutes(namespace).Create(
				ctx,
				&gatewayv1alpha2.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
					Spec: gatewayv1alpha2.UDPRouteSpec{
						Rules: []gatewayv1alpha2.UDPRouteRule{{}},
					},
				},
				metav1.CreateOptions{},
			)
			require.NoError(t, err)

			_, err = gcl.GatewayV1alpha2().TLSRoutes(namespace).Create(
				ctx,
				&gatewayv1alpha2.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
					Spec: gatewayv1alpha2.TLSRouteSpec{
						Rules: []gatewayv1alpha2.TLSRouteRule{{}},
					},
				},
				metav1.CreateOptions{},
			)
			require.NoError(t, err)

			_, err = gcl.GatewayV1beta1().ReferenceGrants(namespace).Create(
				ctx,
				&gatewayv1beta1.ReferenceGrant{
					ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)},
					Spec: gatewayv1beta1.ReferenceGrantSpec{
						From: []gatewayv1beta1.ReferenceGrantFrom{
							{
								Kind:      "test",
								Namespace: metav1.NamespaceDefault,
							},
						},
						To: []gatewayv1beta1.ReferenceGrantTo{
							{
								Kind: "test",
							},
						},
					},
				},
				metav1.CreateOptions{},
			)
			require.NoError(t, err)
		}
	}
}

func handleConnectionToTelemetryServer(ctx context.Context, t *testing.T, listener net.Listener, reportChan chan<- []byte) {
	// Any function return indicates that either the
	// report was sent or there was nothing to send.
	defer close(reportChan)
	conn, err := listener.Accept()
	if !assert.NoError(t, err) {
		return
	}
	defer conn.Close()
	for {
		report := make([]byte, 2048) // Report is much shorter.
		n, err := conn.Read(report)
		if !assert.NoError(t, err) {
			return
		}
		select {
		case reportChan <- report[:n]:
		case <-ctx.Done():
			return
		}
	}
}

func verifyTelemetryReport(t *testing.T, k8sVersion *version.Info, report string) bool {
	t.Helper()
	hostname, err := os.Hostname()
	if err != nil {
		t.Logf("failed to get hostname: %s", err)
		return false
	}
	semver, err := utilversion.ParseGeneric(k8sVersion.GitVersion)
	if err != nil {
		t.Logf("failed to parse k8s version: %s", err)
		return false
	}

	// Report contains stanza like:
	// id=57a7a76c-25d0-4394-ab9a-954f7190e39a;
	// uptime=9;
	// that are not stable across runs, so we need to remove them.
	for _, s := range []string{"id", "uptime"} {
		report, err = removeStanzaFromReport(report, s)
		if err != nil {
			t.Logf("failed to remove stanza %q from report: %s", s, err)
			return false
		}
	}

	expectedReport := fmt.Sprintf(
		"<14>"+
			"signal=kic-ping;"+
			"db=off;"+
			"feature-combinedroutes=true;"+
			"feature-combinedservices=true;"+
			"feature-expressionroutes=false;"+
			"feature-fillids=false;"+
			"feature-gateway-service-discovery=false;"+
			"feature-gateway=false;"+
			"feature-gatewayalpha=false;"+
			"feature-knative=false;"+
			"feature-konnect-sync=false;"+
			"hn=%s;"+
			"kv=3.3.0;"+
			"v=NOT_SET;"+
			"k8s_arch=%s;"+
			"k8s_provider=UNKNOWN;"+
			"k8sv=%s;"+
			"k8sv_semver=%s;"+
			"k8s_gatewayclasses_count=2;"+
			"k8s_gateways_count=4;"+
			"k8s_grpcroutes_count=4;"+
			"k8s_httproutes_count=4;"+
			"k8s_nodes_count=2;"+
			"k8s_pods_count=4;"+
			"k8s_referencegrants_count=4;"+
			"k8s_services_count=5;"+
			"k8s_tcproutes_count=4;"+
			"k8s_tlsroutes_count=4;"+
			"k8s_udproutes_count=4;"+
			"\n",
		hostname,
		k8sVersion.Platform,
		k8sVersion.GitVersion,
		"v"+semver.String(),
	)
	return report == expectedReport
}

// removeStanzaFromReport removes stanza from report. Report contains stanzas like:
// id=57a7a76c-25d0-4394-ab9a-954f7190e39a;uptime=9;k8s_services_count=5; etc.
// Pass e.g. "uptime" to remove the whole uptime=9; from the report.
func removeStanzaFromReport(report string, stanza string) (string, error) {
	const idStanzaEnd = ";"
	stanza = stanza + "="
	start := strings.Index(report, stanza)
	if start == -1 {
		return "", errors.New("stanza not found in report")
	}
	end := strings.Index(report[start:], idStanzaEnd)
	if end == -1 {
		return "", errors.New("stanza end not found in report")
	}
	end += start
	return report[:start] + report[end+1:], nil
}

func generateSelfSignedCert() (tls.Certificate, error) {
	// Generate a new RSA private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create a self-signed X.509 certificate.
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Kong HQ"},
			Country:       []string{"US"},
			Province:      []string{"California"},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"150 Spear Street, Suite 1600"},
			PostalCode:    []string{"94105"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Create a tls.Certificate from the generated private key and certificate.
	certificate := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privateKey,
	}

	return certificate, nil
}
