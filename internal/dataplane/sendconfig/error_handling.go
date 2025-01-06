package sendconfig

import (
	"fmt"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
)

// resourceErrorsToResourceFailures translates a slice of ResourceError to a slice of failures.ResourceFailure.
func resourceErrorsToResourceFailures(resourceErrors []ResourceError, logger logr.Logger) []failures.ResourceFailure {
	var out []failures.ResourceFailure
	for _, ee := range resourceErrors {
		obj := metav1.PartialObjectMetadata{
			TypeMeta: metav1.TypeMeta{
				Kind:       ee.Kind,
				APIVersion: ee.APIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ee.Namespace,
				Name:      ee.Name,
				UID:       k8stypes.UID(ee.UID),
			},
		}
		for problemSource, problem := range ee.Problems {
			logger.V(logging.DebugLevel).Info("Adding failure", "resource_name", ee.Name, "source", problemSource, "problem", problem)
			resourceFailure, failureCreateErr := failures.NewResourceFailure(
				fmt.Sprintf("invalid %s: %s", problemSource, problem),
				&obj,
			)
			if failureCreateErr != nil {
				logger.Error(failureCreateErr, "Could not create resource failure event")
			} else {
				out = append(out, resourceFailure)
			}
		}
	}

	return out
}
