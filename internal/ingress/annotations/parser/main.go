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

package parser

import (
	"fmt"
	"strconv"

	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress/errors"
)

var (
	// AnnotationsPrefix defines the common prefix used in the nginx ingress controller
	AnnotationsPrefix = "nginx.ingress.kubernetes.io"
)

// IngressAnnotation has a method to parse annotations located in Ingress
type IngressAnnotation interface {
	Parse(ing *extensions.Ingress) (interface{}, error)
}

type ingAnnotations map[string]string

func (a ingAnnotations) parseBool(name string) (bool, error) {
	val, ok := a[name]
	if ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false, errors.NewInvalidAnnotationContent(name, val)
		}
		return b, nil
	}
	return false, errors.ErrMissingAnnotations
}

func (a ingAnnotations) parseString(name string) (string, error) {
	val, ok := a[name]
	if ok {
		return val, nil
	}
	return "", errors.ErrMissingAnnotations
}

func (a ingAnnotations) parseInt(name string) (int, error) {
	val, ok := a[name]
	if ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, errors.NewInvalidAnnotationContent(name, val)
		}
		return i, nil
	}
	return 0, errors.ErrMissingAnnotations
}

func checkAnnotation(name string, objectMeta *metav1.ObjectMeta) error {
	if objectMeta == nil || len(objectMeta.GetAnnotations()) == 0 {
		return errors.ErrMissingAnnotations
	}
	if name == "" {
		return errors.ErrInvalidAnnotationName
	}

	return nil
}

// GetBoolAnnotation extracts a boolean from an Ingress annotation
func GetBoolAnnotation(name string, ing *extensions.Ingress) (bool, error) {
	v := GetAnnotationWithPrefix(name)
	if ing == nil {
		return false, errors.ErrMissingAnnotations
	}

	err := checkAnnotation(v, &ing.ObjectMeta)

	if err != nil {
		return false, err
	}

	return ingAnnotations(ing.GetAnnotations()).parseBool(v)
}

// GetStringAnnotation extracts a string from an Ingress annotation
func GetStringAnnotation(name string, objectMeta *metav1.ObjectMeta) (string, error) {
	v := GetAnnotationWithPrefix(name)
	if objectMeta == nil {
		return "", errors.ErrMissingAnnotations
	}

	err := checkAnnotation(v, objectMeta)
	if err != nil {
		return "", err
	}
	return ingAnnotations(objectMeta.GetAnnotations()).parseString(v)
}

// GetIntAnnotation extracts an int from an Ingress annotation
func GetIntAnnotation(name string, ing *extensions.Ingress) (int, error) {
	v := GetAnnotationWithPrefix(name)
	if ing == nil {
		return 0, errors.ErrMissingAnnotations
	}

	err := checkAnnotation(v, &ing.ObjectMeta)
	if err != nil {
		return 0, err
	}
	return ingAnnotations(ing.GetAnnotations()).parseInt(v)
}

// GetAnnotationWithPrefix returns the prefix of ingress annotations
func GetAnnotationWithPrefix(suffix string) string {
	return fmt.Sprintf("%v/%v", AnnotationsPrefix, suffix)
}
