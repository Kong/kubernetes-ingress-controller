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

package class

import (
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// IngressKey picks a specific "class" for the Ingress.
	// The controller only processes Ingresses with this annotation either
	// unset, or set to either the configured value or the empty string.
	IngressKey = "kubernetes.io/ingress.class"
)

var (
	// DefaultClass defines the default class used in the nginx ingres controller
	DefaultClass = "nginx"

	// IngressClass sets the runtime ingress class to use
	// An empty string means accept all ingresses without
	// annotation and the ones configured with class nginx
	IngressClass = "nginx"
)

// IsValid returns true if the given KongConsumer either doesn't specify
// the ingress.class annotation, or it's set to the configured in the
// ingress controller.
func IsValid(objectMeta *metav1.ObjectMeta) bool {
	ingress, ok := objectMeta.GetAnnotations()[IngressKey]
	if !ok {
		glog.V(3).Infof("annotation %v is not present in custom resources %v/%v", IngressKey, objectMeta.Namespace, objectMeta.Name)
	}

	// we have 2 valid combinations
	// 1 - ingress with default class | blank annotation on ingress
	// 2 - ingress with specific class | same annotation on ingress
	//
	// and 2 invalid combinations
	// 3 - ingress with default class | fixed annotation on ingress
	// 4 - ingress with specific class | different annotation on ingress
	if ingress == "" && IngressClass == DefaultClass {
		return true
	}

	return ingress == IngressClass
}
