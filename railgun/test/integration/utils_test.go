//+build integration_tests

package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/railgun/controllers"
)

var (
	l = sync.RWMutex{}
	u *url.URL
)

// proxyURL is a threadsafe way to wait for the proxy to be ready
// and then receive the URL where it can be reached.
func proxyURL() *url.URL {
	l.Lock()
	defer l.Unlock()

	if u == nil {
		u = <-proxyReady
	}

	return u
}

// FIXME - move this into KTF
const proxyDeploymentName = "ingress-controller-kong"

func updateProxyListeners(ctx context.Context, name, kongStreamListen string, containerPorts ...corev1.ContainerPort) (svc *corev1.Service, cleanup func() error, err error) {
	// gather the proxy container as it will need to be specially configured to serve TCP
	var proxy *appsv1.Deployment
	proxy, err = cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Get(ctx, proxyDeploymentName, metav1.GetOptions{})
	if err != nil {
		return
	}
	if count := len(proxy.Spec.Template.Spec.Containers); count != 1 {
		err = fmt.Errorf("expected 1 container for proxy deployment, found %d", count)
		return
	}
	container := proxy.Spec.Template.Spec.Containers[0].DeepCopy()

	// override the KONG_STREAM_LISTEN env var in the proxy container
	var originalVal *corev1.EnvVar
	originalVal, err = overrideEnvVar(container, "KONG_STREAM_LISTEN", kongStreamListen)
	if err != nil {
		return
	}

	// make sure we clean up after ourselves
	cleanup = func() error {
		// remove any created Service for the proxy deployment
		if err := cluster.Client().CoreV1().Services(controllers.DefaultNamespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) { // if the service is not found, that's not a problem nothing to do.
				return err
			}
		}

		// retrieve the current proxy
		proxy, err := cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Get(ctx, "ingress-controller-kong", metav1.GetOptions{})
		if err != nil {
			return err
		}

		// update the KONG_STREAM_LISTEN environment configuration back to its previous value
		container := proxy.Spec.Template.Spec.Containers[0].DeepCopy()
		_, err = overrideEnvVar(container, "KONG_STREAM_LISTEN", originalVal.Value)
		if err != nil {
			return err
		}

		// remove the container ports that were added
		newPorts := make([]corev1.ContainerPort, 0, len(container.Ports))
		for _, port := range container.Ports {
			includePort := true
			for _, configuredPort := range containerPorts {
				if port.Name == configuredPort.Name {
					includePort = false
					break
				}
			}

			if includePort {
				newPorts = append(newPorts, port)
			}
		}
		container.Ports = newPorts

		// revert the corev1.Container to its state prior to the test
		proxy.Spec.Template.Spec.Containers[0] = *container
		_, err = cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Update(ctx, proxy, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		// ensure that the proxy deployment is ready before we proceed
		// FIXME - make this generic k8s test functionality (see below)
		ready := false
		timeout := time.Now().Add(proxyUpdateWait)
		for timeout.After(time.Now()) {
			d, err := cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Get(ctx, proxy.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			if d.Status.ReadyReplicas == d.Status.Replicas && d.Status.AvailableReplicas == d.Status.Replicas && d.Status.UnavailableReplicas < 1 {
				ready = true
				break
			}

			time.Sleep(waitTick)
		}

		if ready {
			return nil
		}

		return fmt.Errorf("proxy did not become ready after %s", proxyUpdateWait)
	}

	// add the provided container ports to the pod
	container.Ports = append(container.Ports, containerPorts...)
	proxy.Spec.Template.Spec.Containers[0] = *container

	// update the deployment with the new container configurations
	proxy, err = cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Update(ctx, proxy, metav1.UpdateOptions{})
	if err != nil {
		return
	}

	// ensure that the proxy deployment is ready before we proceed
	// FIXME - make this generic k8s test functionality
	ready := false
	timeout := time.Now().Add(proxyUpdateWait)
	for timeout.After(time.Now()) {
		var d *appsv1.Deployment
		d, err = cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Get(ctx, proxy.Name, metav1.GetOptions{})
		if err != nil {
			return
		}
		if d.Status.ReadyReplicas == d.Status.Replicas && d.Status.AvailableReplicas == d.Status.Replicas && d.Status.UnavailableReplicas < 1 {
			ready = true
			break
		}

		time.Sleep(waitTick)
	}

	// if the proxy deployment is ready, expose the new container ports via a LoadBalancer Service
	if ready {
		// configure the provided container ports for a LB service
		svc = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: corev1.ServiceSpec{
				Type:     corev1.ServiceTypeLoadBalancer,
				Selector: proxy.Spec.Selector.MatchLabels,
			},
		}
		for _, containerPort := range containerPorts {
			servicePort := corev1.ServicePort{
				Protocol:   containerPort.Protocol,
				Port:       containerPort.ContainerPort,
				TargetPort: intstr.FromInt(int(containerPort.ContainerPort)),
			}
			svc.Spec.Ports = append(svc.Spec.Ports, servicePort)
		}
		svc, err = cluster.Client().CoreV1().Services(controllers.DefaultNamespace).Create(ctx, svc, metav1.CreateOptions{})
		if err != nil {
			return
		}

		// wait for the LB service to be provisioned
		provisioned := false
		timeout := time.Now().Add(serviceWait)
		for timeout.After(time.Now()) {
			svc, err = cluster.Client().CoreV1().Services(controllers.DefaultNamespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return
			}
			if len(svc.Status.LoadBalancer.Ingress) > 0 {
				ing := svc.Status.LoadBalancer.Ingress[0]
				if ip := ing.IP; ip != "" {
					resp, err := http.Get(fmt.Sprintf("http://%s:8001/services", ip))
					if err != nil {
						continue
					}
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						provisioned = true
						break
					}
				}
			}
			time.Sleep(waitTick)
		}

		if !provisioned {
			err = fmt.Errorf("load balancer service for deployment %s did not provision successfully within %s", name, serviceWait)
		}
	} else {
		err = fmt.Errorf("deployment not ready after %s", timeout)
	}

	return
}

func overrideEnvVar(container *corev1.Container, key, val string) (original *corev1.EnvVar, err error) {
	newEnv := make([]corev1.EnvVar, 0, len(container.Env))
	for _, envvar := range container.Env {
		// override any existing value with our custom configuration
		if envvar.Name == key {
			// save the original configuration so we can put it back after we finish testing
			original = envvar.DeepCopy()
			envvar.Value = val
		}
		newEnv = append(newEnv, envvar)
	}

	if original == nil {
		err = fmt.Errorf("could not override env var: %s was not present on container %s", key, container.Name)
	}

	container.Env = newEnv
	return
}
