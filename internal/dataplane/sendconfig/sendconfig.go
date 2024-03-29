package sendconfig

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// -----------------------------------------------------------------------------
// Sendconfig - Public Functions
// -----------------------------------------------------------------------------

type UpdateStrategyResolver interface {
	ResolveUpdateStrategy(client UpdateClient) UpdateStrategy
}

type AdminAPIClient interface {
	AdminAPIClient() *kong.Client
	LastConfigSHA() []byte
	SetLastConfigSHA([]byte)
	BaseRootURL() string
	PluginSchemaStore() *util.PluginSchemaStore

	IsKonnect() bool
	KonnectControlPlane() string
}

// Note that UpdateError is used for both DB-less and DB-backed update strategies, but does not have parity between the
// two. Pending a refactor of database reconciler (https://github.com/Kong/go-database-reconciler/issues/22), KIC does
// not have access to per-resource errors or any original response bodies in database mode. Future DB mode work would
// need to find some way to reconcile the single /config raw body and many per-resource responses from DB endpoints.

// UpdateError wraps several pieces of error information relevant to a failed Kong update attempt.
type UpdateError struct {
	// RawBody is the original Kong HTTP error response body from a failed update.
	RawBody []byte
	// ResourceFailures are per-resource failures from a Kong configuration update attempt.
	ResourceFailures []failures.ResourceFailure
	// Err is an overall description of the update failure.
	Err error
}

// Error implements the Error interface. It returns the string value of the Err field.
func (e UpdateError) Error() string {
	return fmt.Sprintf("%s", e.Err)
}

// PerformUpdate writes `targetContent` to Kong Admin API specified by `kongConfig`.
func PerformUpdate(
	ctx context.Context,
	logger logr.Logger,
	client AdminAPIClient,
	config Config,
	targetContent *file.Content,
	promMetrics *metrics.CtrlFuncMetrics,
	updateStrategyResolver UpdateStrategyResolver,
	configChangeDetector ConfigurationChangeDetector,
) ([]byte, UpdateError) {
	oldSHA := client.LastConfigSHA()
	newSHA, err := deckgen.GenerateSHA(targetContent)
	if err != nil {
		return oldSHA, UpdateError{ResourceFailures: []failures.ResourceFailure{}, Err: err}
	}

	// disable optimization if reverse sync is enabled
	if !config.EnableReverseSync {
		configurationChanged, err := configChangeDetector.HasConfigurationChanged(ctx, oldSHA, newSHA, targetContent, client, client.AdminAPIClient())
		if err != nil {
			return nil, UpdateError{Err: err}
		}
		if !configurationChanged {
			if client.IsKonnect() {
				logger.V(util.DebugLevel).Info("No configuration change, skipping sync to Konnect")
			} else {
				logger.V(util.DebugLevel).Info("No configuration change, skipping sync to Kong")
			}
			return oldSHA, UpdateError{}
		}
	}

	updateStrategy := updateStrategyResolver.ResolveUpdateStrategy(client)
	logger = logger.WithValues("update_strategy", updateStrategy.Type())
	timeStart := time.Now()
	err, resourceErrors, rawErrBody, resourceErrorsParseErr := updateStrategy.Update(ctx, ContentWithHash{
		Content: targetContent,
		Hash:    newSHA,
	})
	duration := time.Since(timeStart)

	metricsProtocol := updateStrategy.MetricsProtocol()
	if err != nil {
		// Not pushing metrics in case it's an update skip due to a backoff.
		if errors.As(err, &UpdateSkippedDueToBackoffStrategyError{}) {
			return nil, UpdateError{RawBody: rawErrBody, Err: err}
		}

		resourceFailures := resourceErrorsToResourceFailures(resourceErrors, resourceErrorsParseErr, logger)
		promMetrics.RecordPushFailure(metricsProtocol, duration, client.BaseRootURL(), len(resourceFailures), err)
		return nil, UpdateError{ResourceFailures: resourceFailures, RawBody: rawErrBody, Err: err}
	}

	promMetrics.RecordPushSuccess(metricsProtocol, duration, client.BaseRootURL())

	if client.IsKonnect() {
		logger.V(util.InfoLevel).Info("Successfully synced configuration to Konnect")
	} else {
		logger.V(util.InfoLevel).Info("Successfully synced configuration to Kong")
	}

	return newSHA, UpdateError{}
}

// -----------------------------------------------------------------------------
// Sendconfig - Private Functions
// -----------------------------------------------------------------------------

// resourceErrorsToResourceFailures translates a slice of ResourceError to a slice of failures.ResourceFailure.
// In case of parseErr being not nil, it just returns a nil slice.
func resourceErrorsToResourceFailures(resourceErrors []ResourceError, parseErr error, logger logr.Logger) []failures.ResourceFailure {
	if parseErr != nil {
		logger.Error(parseErr, "Failed parsing resource errors")
		return nil
	}

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
			logger.V(util.DebugLevel).Info("Adding failure", "resource_name", ee.Name, "source", problemSource, "problem", problem)
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
