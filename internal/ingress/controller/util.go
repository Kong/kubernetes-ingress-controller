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

package controller

import (
	"strings"

	"github.com/hbagdi/go-kong/kong"
)

func isEmpty(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

func toStringPtrArray(array []string) []*string {
	var result []*string
	for _, element := range array {
		e := element
		result = append(result, &e)
	}
	return result
}

func toStringArray(array []*string) []string {
	var result []string
	for _, element := range array {
		e := *element
		result = append(result, e)
	}
	return result
}

// TODO refactor this away
func compareRoute(r1, r2 *kong.Route) bool {
	if r1 == r2 {
		return true
	}
	if r1 == nil || r2 == nil {
		return false
	}

	if len(r1.Hosts) != len(r2.Hosts) {
		return false
	}

	for _, r1b := range r1.Hosts {
		found := false
		for _, r2b := range r2.Hosts {
			if *r1b == *r2b {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(r1.Paths) != len(r2.Paths) {
		return false
	}

	for _, r1b := range r1.Paths {
		found := false
		for _, r2b := range r2.Paths {
			if *r1b == *r2b {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if r1.Service != nil && r2.Service != nil {
		if *r1.Service.ID != *r2.Service.ID {
			return false
		}

	}

	return true
}
