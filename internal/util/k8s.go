/*
Copyright 2015 The Kubernetes Authors.

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

package util

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

// ParseNameNS parses a string searching a namespace and name.
func ParseNameNS(input string) (string, string, error) {
	nsName := strings.Split(input, "/")
	if len(nsName) != 2 {
		return "", "", fmt.Errorf("invalid format (namespace/name) found in '%v'", input)
	}

	return nsName[0], nsName[1], nil
}

// GetNodeIPOrName returns the IP address or the name of a node in the cluster.
func GetNodeIPOrName(ctx context.Context, kubeClient clientset.Interface, name string) string {
	node, err := kubeClient.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return ""
	}

	ip := ""

	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeExternalIP {
			if address.Address != "" {
				ip = address.Address
				break
			}
		}
	}

	// Report the external IP address of the node
	if ip != "" {
		return ip
	}

	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			if address.Address != "" {
				ip = address.Address
				break
			}
		}
	}

	return ip
}

// GetPodNN returns NamespacedName of pod that this process is running in.
func GetPodNN() (k8stypes.NamespacedName, error) {
	nn := k8stypes.NamespacedName{
		Namespace: os.Getenv("POD_NAMESPACE"),
		Name:      os.Getenv("POD_NAME"),
	}
	if nn.Name == "" || nn.Namespace == "" {
		return k8stypes.NamespacedName{},
			fmt.Errorf("unable to get POD information (missing POD_NAME or POD_NAMESPACE environment variable")
	}

	return nn, nil
}

// PodInfo contains runtime information about the pod running the Ingres controller.
type PodInfo struct {
	Name      string
	Namespace string
	NodeIP    string
	// Labels selectors of the running pod
	// This is used to search for other Ingress controller pods
	Labels map[string]string
}

// GetPodDetails returns runtime information about the pod:
// name, namespace and IP of the node where it is running.
func GetPodDetails(ctx context.Context, kubeClient clientset.Interface) (*PodInfo, error) {
	nn, err := GetPodNN()
	if err != nil {
		return nil, err
	}

	pod, err := kubeClient.CoreV1().Pods(nn.Namespace).Get(ctx, nn.Name, metav1.GetOptions{})
	if pod == nil {
		return nil, fmt.Errorf("unable to get POD information: %w", err)
	}

	return &PodInfo{
		Name:      nn.Name,
		Namespace: nn.Namespace,
		NodeIP:    GetNodeIPOrName(ctx, kubeClient, pod.Spec.NodeName),
		Labels:    pod.GetLabels(),
	}, nil
}

// map of all the supported Group/Kinds for the backend. At the moment, only
// core services are supported, but to provide support to other kinds, it is
// enough to add entries to this map.
var backendRefSupportedGroupKinds = map[string]struct{}{
	"core/Service": {},
}

// IsBackendRefGroupKindSupported checks if the GroupKind of the object used as
// BackendRef for the HTTPRoute is supported.
func IsBackendRefGroupKindSupported(gatewayAPIGroup *gatewayv1.Group, gatewayAPIKind *gatewayv1.Kind) bool {
	if gatewayAPIKind == nil {
		return false
	}

	group := "core"
	if gatewayAPIGroup != nil && *gatewayAPIGroup != "" {
		group = string(*gatewayAPIGroup)
	}

	_, ok := backendRefSupportedGroupKinds[fmt.Sprintf("%s/%s", group, *gatewayAPIKind)]
	return ok
}

const (
	K8sNamespaceTagPrefix = "k8s-namespace:"
	K8sNameTagPrefix      = "k8s-name:"
	K8sUIDTagPrefix       = "k8s-uid:"
	K8sKindTagPrefix      = "k8s-kind:"
	K8sGroupTagPrefix     = "k8s-group:"
	K8sVersionTagPrefix   = "k8s-version:"
)

// GenerateTagsForObject returns a subset of an object's metadata as a slice of prefixed string pointers.
func GenerateTagsForObject(obj client.Object) []*string {
	if obj == nil {
		// this should never happen in practice, but it happen in some unit tests
		// in those cases, the nil object has no tags
		return nil
	}
	gvk := obj.GetObjectKind().GroupVersionKind()
	tags := []string{}
	if obj.GetName() != "" {
		tags = append(tags, K8sNameTagPrefix+obj.GetName())
	}
	if obj.GetNamespace() != "" {
		tags = append(tags, K8sNamespaceTagPrefix+obj.GetNamespace())
	}
	if gvk.Kind != "" {
		tags = append(tags, K8sKindTagPrefix+gvk.Kind)
	}
	if string(obj.GetUID()) != "" {
		tags = append(tags, K8sUIDTagPrefix+string(obj.GetUID()))
	}
	if gvk.Group != "" {
		tags = append(tags, K8sGroupTagPrefix+gvk.Group)
	}
	if gvk.Version != "" {
		tags = append(tags, K8sVersionTagPrefix+gvk.Version)
	}

	tags = append(tags,
		lo.Uniq(
			annotations.ExtractUserTags(obj.GetAnnotations()),
		)...,
	)
	return kong.StringSlice(tags...)
}
