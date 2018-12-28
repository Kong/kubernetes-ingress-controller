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

	consumerv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/consumer/v1"
	credentialv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/credential/v1"
	pluginv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/plugin/v1"
	extensions "k8s.io/api/extensions/v1beta1"

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

func checkAnnotation(name string, ing *extensions.Ingress) error {
	if ing == nil || len(ing.GetAnnotations()) == 0 {
		return errors.ErrMissingAnnotations
	}
	if name == "" {
		return errors.ErrInvalidAnnotationName
	}

	return nil
}

func checkAnnotationPlugin(name string, plugin *pluginv1.KongPlugin) error {
	if plugin == nil || len(plugin.GetAnnotations()) == 0 {
		return errors.ErrMissingAnnotations
	}
	if name == "" {
		return errors.ErrInvalidAnnotationName
	}

	return nil
}

func checkAnnotationCredential(name string, credential *credentialv1.KongCredential) error {
	if credential == nil || len(credential.GetAnnotations()) == 0 {
		return errors.ErrMissingAnnotations
	}
	if name == "" {
		return errors.ErrInvalidAnnotationName
	}

	return nil
}

func checkAnnotationConsumer(name string, consumer *consumerv1.KongConsumer) error {
	if consumer == nil || len(consumer.GetAnnotations()) == 0 {
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
	err := checkAnnotation(v, ing)
	if err != nil {
		return false, err
	}
	return ingAnnotations(ing.GetAnnotations()).parseBool(v)
}

// GetStringAnnotation extracts a string from an Ingress annotation
func GetStringAnnotation(name string, ing *extensions.Ingress) (string, error) {
	v := GetAnnotationWithPrefix(name)
	err := checkAnnotation(v, ing)
	if err != nil {
		return "", err
	}
	return ingAnnotations(ing.GetAnnotations()).parseString(v)
}

// GetStringAnnotationPlugin extracts a string from an Ingress annotation
func GetStringAnnotationPlugin(name string, plugin *pluginv1.KongPlugin) (string, error) {
	v := GetAnnotationWithPrefix(name)
	err := checkAnnotationPlugin(v, plugin)
	if err != nil {
		return "", err
	}
	return ingAnnotations(plugin.GetAnnotations()).parseString(v)
}

// GetStringAnnotationCredential extracts a string from an Ingress annotation
func GetStringAnnotationCredential(name string, credential *credentialv1.KongCredential) (string, error) {
	v := GetAnnotationWithPrefix(name)
	err := checkAnnotationCredential(v, credential)
	if err != nil {
		return "", err
	}
	return ingAnnotations(credential.GetAnnotations()).parseString(v)
}

// GetStringAnnotationConsumer extracts a string from an Ingress annotation
func GetStringAnnotationConsumer(name string, consumer *consumerv1.KongConsumer) (string, error) {
	v := GetAnnotationWithPrefix(name)
	err := checkAnnotationConsumer(v, consumer)
	if err != nil {
		return "", err
	}
	return ingAnnotations(consumer.GetAnnotations()).parseString(v)
}

// GetIntAnnotation extracts an int from an Ingress annotation
func GetIntAnnotation(name string, ing *extensions.Ingress) (int, error) {
	v := GetAnnotationWithPrefix(name)
	err := checkAnnotation(v, ing)
	if err != nil {
		return 0, err
	}
	return ingAnnotations(ing.GetAnnotations()).parseInt(v)
}

// GetAnnotationWithPrefix returns the prefix of ingress annotations
func GetAnnotationWithPrefix(suffix string) string {
	return fmt.Sprintf("%v/%v", AnnotationsPrefix, suffix)
}
