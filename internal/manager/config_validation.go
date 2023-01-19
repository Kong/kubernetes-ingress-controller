package manager

import (
	"errors"
	"strings"

	"k8s.io/apimachinery/pkg/types"
)

// This file contains a set of constructors that are used to validate and set validated values in Config.
// They're meant to be used together with ValidatedValue[T] type.

func createNamespacedName(s string) (types.NamespacedName, error) {
	parts := strings.SplitN(s, "/", 3)
	if len(parts) != 2 {
		return types.NamespacedName{}, errors.New("the expected format is namespace/name")
	}
	return types.NamespacedName{
		Namespace: parts[0],
		Name:      parts[1],
	}, nil
}

func createGatewayAPIControllerName(s string) (string, error) {
	if !isControllerNameValid(s) {
		return "", errors.New("the expected format is example.com/controller-name")
	}
	return s, nil
}
