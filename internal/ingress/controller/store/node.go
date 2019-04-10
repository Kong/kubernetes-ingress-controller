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

package store

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// NodeLister makes a Store that lists Nodes.
type NodeLister struct {
	cache.Store
}

// ByKey searches for a Node in the local secrets Store
func (nl *NodeLister) ByKey(key string) (*apiv1.Node, error) {
	n, exists, err := nl.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Node %v was not found", key)
	}
	return n.(*apiv1.Node), nil
}

// GetNodes returns all nodes
func (nl *NodeLister) GetNodes() ([]*apiv1.Node, error) {
	var result []*apiv1.Node
	for _, n := range nl.List() {
		node, ok := n.(*apiv1.Node)
		if ok {
			result = append(result, node)
		}
	}
	return result, nil
}
