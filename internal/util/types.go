/*
Copyright 2018 The Kubernetes Authors.

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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
)

// Endpoint describes a kubernetes endpoint, same as a target in Kong.
type Endpoint struct {
	// Address IP address of the endpoint
	Address string `json:"address"`
	// Port number of the TCP port
	Port string `json:"port"`
	// Terminating indicates if the endpoint is in terminating state
	Terminating bool `json:"terminating,omitempty"`
}

// TypeMeta is stripped after unmarshaling into Go struct due to the issue described in
// https://github.com/kubernetes/kubernetes/issues/3030, but we need it for various purposes.

// This function is adopted from https://github.com/kubernetes/cli-runtime/blob/v0.28.2/pkg/printers/typesetter.go#L39-L72
// It has been modified to remove the io.Writer output and to use the scheme package directly instead of the Typer helper.

// PopulateTypeMeta adds GVK information to a runtime.Object that may not have it available in the object TypeMeta.
func PopulateTypeMeta(obj runtime.Object, s *runtime.Scheme) error {
	gvks, _, err := s.ObjectKinds(obj)
	if err != nil {
		return fmt.Errorf("missing apiVersion or kind and cannot assign it; %w", err)
	}

	for _, gvk := range gvks {
		if len(gvk.Kind) == 0 {
			continue
		}
		if len(gvk.Version) == 0 || gvk.Version == runtime.APIVersionInternal {
			continue
		}
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		break
	}

	return nil
}
