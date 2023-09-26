package sendconfig

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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
	KonnectRuntimeGroup() string
}

// PerformUpdate writes `targetContent` to Kong Admin API specified by `kongConfig`.
func PerformUpdate(
	ctx context.Context,
	log logrus.FieldLogger,
	client AdminAPIClient,
	config Config,
	targetContent *file.Content,
	promMetrics *metrics.CtrlFuncMetrics,
	updateStrategyResolver UpdateStrategyResolver,
	configChangeDetector ConfigurationChangeDetector,
) ([]byte, []FlatEntityError, error) {
	oldSHA := client.LastConfigSHA()
	newSHA, err := deckgen.GenerateSHA(targetContent)
	if err != nil {
		return oldSHA, nil, err
	}

	// disable optimization if reverse sync is enabled
	if !config.EnableReverseSync {
		configurationChanged, err := configChangeDetector.HasConfigurationChanged(ctx, oldSHA, newSHA, targetContent, client, client.AdminAPIClient())
		if err != nil {
			return nil, nil, err
		}
		if !configurationChanged {
			if client.IsKonnect() {
				log.Debug("no configuration change, skipping sync to Konnect")
			} else {
				log.Debug("no configuration change, skipping sync to Kong")
			}
			return oldSHA, nil, nil
		}
	}

	updateStrategy := updateStrategyResolver.ResolveUpdateStrategy(client)
	log = log.WithField("update_strategy", updateStrategy.Type())
	timeStart := time.Now()
	err, entityErrors := updateStrategy.Update(ctx, ContentWithHash{
		Content: targetContent,
		Hash:    newSHA,
	})
	duration := time.Since(timeStart)

	metricsProtocol := updateStrategy.MetricsProtocol()
	if err != nil || len(entityErrors) > 0 {
		// Not pushing metrics in case it's an update skip due to a backoff.
		if errors.As(err, &UpdateSkippedDueToBackoffStrategyError{}) {
			return nil, nil, err
		}

		resourceErrors := ResourceErrorsFromEntityErrors(entityErrors, log)
		resourceFailures := ResourceErrorsToResourceFailures(resourceErrors, log)
		promMetrics.RecordPushFailure(metricsProtocol, duration, client.BaseRootURL(), len(resourceFailures), err)
		return nil, entityErrors, err
	}

	promMetrics.RecordPushSuccess(metricsProtocol, duration, client.BaseRootURL())

	if client.IsKonnect() {
		log.Info("successfully synced configuration to Konnect")
	} else {
		log.Info("successfully synced configuration to Kong")
	}

	return newSHA, nil, nil
}

// ResourceErrorsToResourceFailures translates a slice of ResourceError to a slice of failures.ResourceFailure.
// In case of parseErr being not nil, it just returns a nil slice.
func ResourceErrorsToResourceFailures(resourceErrors []ResourceError, log logrus.FieldLogger) []failures.ResourceFailure {
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
		for field, problem := range ee.Problems {
			log.Debug(fmt.Sprintf("adding failure for %s: %s = %s", ee.Name, field, problem))
			resourceFailure, failureCreateErr := failures.NewResourceFailure(
				fmt.Sprintf("invalid %s: %s", field, problem),
				&obj,
			)
			if failureCreateErr != nil {
				log.WithError(failureCreateErr).Error("could create resource failure event")
			} else {
				out = append(out, resourceFailure)
			}
		}
	}

	return out
}
