package kong

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// -----------------------------------------------------------------------------
// Public Consts & Vars
// -----------------------------------------------------------------------------

var (
	// ProxyPodTimeout is the default amount of time to wait for a Kong proxy deployment to be ready.
	ProxyPodTimeout = time.Minute * 3
)

// -----------------------------------------------------------------------------
// Public Functions - Gateway Management
// -----------------------------------------------------------------------------

// SimpleProxySetup deploys a single (bare) Kong Gateway configured with an exposed Admin API.
// This does not provide mTLS and is only intended for local testing purposes.
func SimpleProxySetup(kc *kubernetes.Clientset, namespace string) error {
	// ensure that the namespace provided exists
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	if _, err := kc.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	// create the deployment for the proxy
	_, err := kc.AppsV1().Deployments(namespace).Create(context.Background(), defaultProxyDeployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// expose the proxy via ClusterIP
	_, err = kc.CoreV1().Services(namespace).Create(context.Background(), defaultProxyService, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// wait for the proxy to be ready before continuing
	ok := false
	timeout := time.Now().Add(ProxyPodTimeout)
	for time.Now().Before(timeout) {
		deployment, err := kc.AppsV1().Deployments(namespace).Get(context.Background(), "kong", metav1.GetOptions{})
		if err != nil {
			return err
		}
		if deployment.Status.ReadyReplicas > 0 {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("kong proxy was not ready after %s\n", ProxyPodTimeout)
	}

	return nil
}

// DeployControllers deploys the Kong Kubernetes Ingress Controller (KIC) and other relevant
// Ingress controllers to the provided cluster given a *kubernetes.Clientset for it.
// FIXME: this is a total hack for now
func DeployControllers(kc *kubernetes.Clientset, containerImage, namespace string) (context.CancelFunc, error) {
	// ensure the controller namespace is created
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	if _, err := kc.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	// FIXME: temp logging file
	tmpfile, err := ioutil.TempFile(os.TempDir(), "kong-integration-tests-")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "go", "run", "../../main.go")
	cmd.Stdout = tmpfile
	cmd.Stderr = tmpfile
	fmt.Fprintf(os.Stdout, "INFO: tempfile for controller logs: %s\n", tmpfile.Name())

	go func() {
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}()

	return cancel, nil
}

// -----------------------------------------------------------------------------
// Private Consts & Vars
// -----------------------------------------------------------------------------

var (
	defaultProxyDeployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "kong"},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app.kubernetes.io/name": "kong"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app.kubernetes.io/name": "kong"}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "proxy",
						Image:           "kong:2.2",
						ImagePullPolicy: "IfNotPresent",
						Env: []corev1.EnvVar{
							{Name: "KONG_ADMIN_LISTEN", Value: "0.0.0.0:8001"},
							{Name: "KONG_DATABASE", Value: "off"},
							{Name: "KONG_PROXY_LISTEN", Value: "0.0.0.0:8000, 0.0.0.0:8443 http2 ssl"},
							{Name: "KONG_STATUS_LISTEN", Value: "0.0.0.0:8100"},
						},
						Ports: []corev1.ContainerPort{
							{Name: "admin", ContainerPort: 8001, Protocol: corev1.ProtocolTCP},
							{Name: "proxy", ContainerPort: 8000, Protocol: corev1.ProtocolTCP},
							{Name: "proxy-tls", ContainerPort: 8443, Protocol: corev1.ProtocolTCP},
							{Name: "status", ContainerPort: 8100, Protocol: corev1.ProtocolTCP},
						}}}}}}}
	defaultProxyService = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "kong"},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app.kubernetes.io/name": "kong",
			},
			Ports: []corev1.ServicePort{
				{Name: "admin", Port: 8001, TargetPort: intstr.FromInt(8001), Protocol: corev1.ProtocolTCP},
				{Name: "proxy", Port: 8000, TargetPort: intstr.FromInt(8000), Protocol: corev1.ProtocolTCP},
				{Name: "proxy-tls", Port: 8443, TargetPort: intstr.FromInt(8443), Protocol: corev1.ProtocolTCP},
				{Name: "status", Port: 8100, TargetPort: intstr.FromInt(8100), Protocol: corev1.ProtocolTCP},
			}}}
)
