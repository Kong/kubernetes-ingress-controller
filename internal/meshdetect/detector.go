package meshdetect

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Detector provides methods to detect the following:
//
//   - whether a service mesh is deployed to the cluster
//   - whether KIC is injected with service mesh
//   - count of services injected within mesh networks
type Detector struct {
	// Client is the kubernetes client to read kubernetes services.
	Client client.Client

	// PodNamespace is the namespace of the pod in which the mesh detector is running.
	PodNamespace string

	// PodName is the name of pod  in which the mesh detector is running.
	PodName string

	// PublishServiceName is the name of the Kubernetes service used for
	// ingress traffic to the Kong Gateway.
	PublishServiceName string

	Logger logr.Logger
}

// NewDetectorByConfig creates a new Detector provided a Kubernetes
// config for the relevant cluster and the name of the Kubernetes service
// for ingress traffic to the Kong Gateway.
func NewDetectorByConfig(kubeCfg *rest.Config,
	podNamespace, podName, publishServiceName string,
	logger logr.Logger) (*Detector, error) {
	c, err := client.New(kubeCfg, client.Options{})
	if err != nil {
		return nil, err
	}

	return &Detector{
		Client:             c,
		PodNamespace:       podNamespace,
		PodName:            podName,
		PublishServiceName: publishServiceName,
		Logger:             logger,
	}, nil
}

// DetectMeshDeployment detects which kinds of mesh networks are deployed.
func (d *Detector) DetectMeshDeployment(ctx context.Context) map[MeshKind]*DeploymentResults {
	deploymentResults := map[MeshKind]*DeploymentResults{}

	for _, meshKind := range MeshesToDetect {
		deploymentResult := &DeploymentResults{}
		if d.detectMeshDeploymentByService(ctx, meshKind) {
			deploymentResult.ServiceExists = true
		}
		deploymentResults[meshKind] = deploymentResult
	}

	return deploymentResults
}

// detectMeshDeploymentByService finds the service for each mesh in the cluster.
func (d *Detector) detectMeshDeploymentByService(ctx context.Context, meshKind MeshKind) bool {
	serviceName := meshServiceName[meshKind]

	svcList := &corev1.ServiceList{}
	err := d.Client.List(ctx, svcList, &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{"metadata.name": serviceName}),
	})
	if err != nil {
		return false
	}

	for _, svc := range svcList.Items {
		if svc.Name == serviceName {
			return true
		}
	}

	return false
}

// DetectRunUnder detects whether KIC pod is running under each kind of service mesh.
// in this function, we want to detect whether the pod running this detector, which is
// also the KIC pod, is running under a certain kind of service mesh.
// for example, if the pod is injected with istio sidecar container and init container,
// we report that the detector is running under istio mesh.
func (d *Detector) DetectRunUnder(ctx context.Context) map[MeshKind]*RunUnderResults {
	runUnderResults := map[MeshKind]*RunUnderResults{}
	// get the pod itself.
	pod := &corev1.Pod{}
	err := d.Client.Get(ctx, client.ObjectKey{Namespace: d.PodNamespace, Name: d.PodName}, pod)
	if err != nil {
		return runUnderResults
	}

	publishService := &corev1.Service{}
	if d.PublishServiceName != "" {
		parts := strings.Split(d.PublishServiceName, "/")
		objKey := client.ObjectKey{}
		// format <namespace>/<name>
		if len(parts) == 2 {
			objKey.Namespace = parts[0]
			objKey.Name = parts[1]
		} else {
			d.Logger.Info("publish service " + d.PublishServiceName + " is not in <namespace>/<name> format")
		}
		// only try to get service if the namespace and name are correctly filled
		if objKey.Name != "" && objKey.Namespace != "" {
			err := d.Client.Get(ctx, objKey, publishService)
			if err != nil {
				d.Logger.Info(
					"failed to get service to publish gateway configuration:"+err.Error(),
					"namespace", objKey.Namespace, "name", objKey.Name)
			}
		}
	}

	for _, meshKind := range MeshesToDetect {
		runUnderResults[meshKind] = &RunUnderResults{}

		// detect if service for kong-gateway has annotations(only for traefik)
		if publishService != nil && isServiceInjected(meshKind, publishService) {
			runUnderResults[meshKind].PodOrServiceAnnotation = true
		}

		// detect if pod has annotations.
		podAnnotations := meshPodAnnotations[meshKind]
		if podAnnotations != nil && podAnnotations.Matches(labels.Set(pod.Annotations)) {
			runUnderResults[meshKind].PodOrServiceAnnotation = true
		}

		// detect if pod has a sidecar.
		runUnderResults[meshKind].SidecarContainerInjected = isPodSidecarInjected(meshKind, pod)

		// detect if pod has a init container.
		runUnderResults[meshKind].InitContainerInjected = isPodInitContainerInjected(meshKind, pod)

	}

	return runUnderResults
}

func isServiceInjected(meshKind MeshKind, svc *corev1.Service) bool {
	if svc == nil {
		return false
	}
	if svc.Annotations == nil {
		return false
	}

	svcAnnotations := meshServiceAnnotations[meshKind]
	if svcAnnotations == nil {
		return false
	}
	if svcAnnotations.Matches(labels.Set(svc.Annotations)) {
		return true
	}

	return false
}

func isPodSidecarInjected(meshKind MeshKind, pod *corev1.Pod) bool {
	sidecarName := meshSidecarContainerName[meshKind]
	if sidecarName == "" {
		return false
	}
	for _, container := range pod.Spec.Containers {
		if container.Name == sidecarName {

			switch meshKind { // nolint:exhaustive
			case MeshKindAWSAppMesh:
				// special judgement for AWS appmesh here:
				// AWS appmesh uses `envoy` as sidecar name, which is a very common name.
				// We do a further check on the container and treat as really injected
				// if the container uses image `aws-appmesh-envoy:*`.
				if strings.Contains(container.Image, awsAppMeshEnvoyImageName) {
					return true
				}
			default:
				// for meshes other than AWS app mesh, directly return true (pod injected)
				// when the container with the sidecar name is found.
				return true
			}
		}
	}
	return false
}

func isPodInitContainerInjected(meshKind MeshKind, pod *corev1.Pod) bool {
	initContainerName := meshInitContainerName[meshKind]
	if initContainerName == "" {
		return false
	}
	for _, initContainer := range pod.Spec.InitContainers {
		if initContainer.Name == initContainerName {
			return true
		}
	}

	return false
}

// DetectServiceDistribution detects how many services are running under each mesh.
func (d *Detector) DetectServiceDistribution(ctx context.Context) (*ServiceDistributionResults, error) {
	ret := &ServiceDistributionResults{
		MeshDistribution: map[MeshKind]int{},
	}

	// list all services
	serviceList := &corev1.ServiceList{}
	err := d.Client.List(ctx, serviceList)
	if err != nil {
		d.Logger.Error(err, "failed to list all services in cluster")
		return nil, err
	}

	ret.TotalServices = len(serviceList.Items)

	for i := range serviceList.Items {
		svc := &serviceList.Items[i]
		endpoints := &corev1.Endpoints{}
		err := d.Client.Get(ctx, client.ObjectKey{Namespace: svc.Namespace, Name: svc.Name}, endpoints)
		if err != nil {
			continue
		}

		// injected is set to true if the service(pod) is injected by mesh.
		injected := map[MeshKind]bool{}

		// detect if service has annotations to indicate that the service is injected
		// (only for traefik)

		for meshKind := range meshServiceAnnotations {
			injected[meshKind] = isServiceInjected(meshKind, svc)
		}

		for _, subset := range endpoints.Subsets {
			for _, address := range subset.Addresses {
				// skip if the target endpoint address is not a pod.
				if address.TargetRef == nil {
					continue
				}
				if address.TargetRef.Kind != "Pod" {
					continue
				}

				// if one of the pod is injected, we consider this service as running under the mesh.
				pod := &corev1.Pod{}
				err := d.Client.Get(ctx,
					client.ObjectKey{Namespace: address.TargetRef.Namespace, Name: address.TargetRef.Name},
					pod)
				if err != nil {
					continue
				}

				for _, meshKind := range MeshesToDetect {
					// set injected to true if one of pods in service is injected with sidecar and init container.
					injected[meshKind] = injected[meshKind] ||
						(isPodSidecarInjected(meshKind, pod) || isPodInitContainerInjected(meshKind, pod))
				}
			}
		}

		for meshKind := range injected {
			if injected[meshKind] {
				ret.MeshDistribution[meshKind]++
			}
		}
	}

	return ret, nil
}
