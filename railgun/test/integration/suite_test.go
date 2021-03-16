//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kong/kubernetes-ingress-controller/railgun/controllers"
	"github.com/kong/kubernetes-testing-framework/pkg/runbooks"
)

var (
	// ClusterName indicates the name of the Kind test cluster setup for this test suite.
	ClusterName = uuid.New().String()

	// kc is a kubernetes clientset for the default Kind cluster created for this test suite.
	kc *kubernetes.Clientset

	// ProxyReadyChannel is the channel that indicates when the Kong proxy is ready to use.
	ProxyReadyChannel = make(chan *url.URL)

	// ProxyErrorCh indicates if the Proxy provisioning has failed on the cluster.
	ProxyErrorCh = make(chan error)

	// IngressTimeout is the maximum amount of time that the tests should wait for an Ingress record to be provisioned and the backend accessible.
	IngressTimeout = time.Minute * 5

	// IngressTimeoutTick is the time to wait between Ingress resource timeout checks
	IngressTimeoutTick = time.Second * 1
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	var cleanup func()
	kc, cleanup, err = runbooks.CreateKindClusterWithKongProxy(ctx, ClusterName, ProxyReadyChannel, ProxyErrorCh)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(10)
	}
	defer cleanup()

	// deploy the Kong Kubernetes Ingress Controller (KIC) to the cluster
	// TODO - need to fix the context handling here
	if err := deployControllers(ctx, kc, os.Getenv("KONG_CONTROLLER_TEST_IMAGE"), controllers.DefaultNamespace); err != nil {
		cleanup()
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(17)
	}

	// run the tests
	code := m.Run()

	// cleanup
	cleanup()
	os.Exit(code)
}

// FIXME: this is a total hack for now
func deployControllers(ctx context.Context, kc *kubernetes.Clientset, containerImage, namespace string) error {
	// ensure the controller namespace is created
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	if _, err := kc.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	// FIXME: temp logging file
	tmpfile, err := ioutil.TempFile(os.TempDir(), "kong-integration-tests-")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "INFO: tempfile for controller logs: %s\n", tmpfile.Name())

	go func() {
		u, err := proxyURL()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}

		stderr := new(bytes.Buffer)
		cmd := exec.CommandContext(ctx, "go", "run", "../../main.go", "--kong-url", fmt.Sprintf("http://%s:8001", u.Hostname()))
		cmd.Stdout = tmpfile
		cmd.Stderr = io.MultiWriter(stderr, tmpfile)

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}()

	return nil
}

var prx *url.URL
var lock = sync.Mutex{}

func proxyURL() (*url.URL, error) {
	lock.Lock()
	defer lock.Unlock()

	if prx == nil {
		prx = <-ProxyReadyChannel
	}

	return prx, nil
}
