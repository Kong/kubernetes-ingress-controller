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
) ([]byte, []failures.ResourceFailure, []byte, error) {
	oldSHA := client.LastConfigSHA()
	newSHA, err := deckgen.GenerateSHA(targetContent)
	if err != nil {
		return oldSHA, []failures.ResourceFailure{}, nil, err
	}

	// disable optimization if reverse sync is enabled
	if !config.EnableReverseSync {
		configurationChanged, err := configChangeDetector.HasConfigurationChanged(ctx, oldSHA, newSHA, targetContent, client, client.AdminAPIClient())
		if err != nil {
			return nil, []failures.ResourceFailure{}, nil, err
		}
		if !configurationChanged {
			if client.IsKonnect() {
				logger.V(util.DebugLevel).Info("No configuration change, skipping sync to Konnect")
			} else {
				logger.V(util.DebugLevel).Info("No configuration change, skipping sync to Kong")
			}
			return oldSHA, []failures.ResourceFailure{}, nil, nil
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
			return nil, []failures.ResourceFailure{}, rawErrBody, err
		}

		resourceFailures := resourceErrorsToResourceFailures(resourceErrors, resourceErrorsParseErr, logger)
		promMetrics.RecordPushFailure(metricsProtocol, duration, client.BaseRootURL(), len(resourceFailures), err)
		return nil, resourceFailures, rawErrBody, err
	}

	promMetrics.RecordPushSuccess(metricsProtocol, duration, client.BaseRootURL())

	if client.IsKonnect() {
		logger.V(util.InfoLevel).Info("Successfully synced configuration to Konnect")
	} else {
		logger.V(util.InfoLevel).Info("Successfully synced configuration to Kong")
	}

	return newSHA, nil, rawErrBody, nil
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
