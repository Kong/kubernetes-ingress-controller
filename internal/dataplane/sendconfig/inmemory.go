package sendconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

type ConfigService interface {
	ReloadDeclarativeRawConfig(
		ctx context.Context,
		config io.Reader,
		checkHash bool,
		flattenErrors bool,
	) error
}

type ContentToDBLessConfigConverter interface {
	// Convert converts a decK's file.Content to a DBLessConfig.
	// Implementations are allowed to modify the input *file.Content. Make sure it's copied beforehand if needed.
	Convert(content *file.Content) DBLessConfig
}

// UpdateStrategyInMemory implements the UpdateStrategy interface. It updates Kong's data-plane
// configuration using its `POST /config` endpoint that is used by ConfigService.ReloadDeclarativeRawConfig.
type UpdateStrategyInMemory struct {
	configService   ConfigService
	configConverter ContentToDBLessConfigConverter
	logger          logr.Logger
}

func NewUpdateStrategyInMemory(
	configService ConfigService,
	configConverter ContentToDBLessConfigConverter,
	logger logr.Logger,
) UpdateStrategyInMemory {
	return UpdateStrategyInMemory{
		configService:   configService,
		configConverter: configConverter,
		logger:          logger,
	}
}

func (s UpdateStrategyInMemory) Update(ctx context.Context, targetState ContentWithHash) (int, error) {
	dblessConfig := s.configConverter.Convert(targetState.Content)
	config, err := json.Marshal(dblessConfig)
	if err != nil {
		return 0, fmt.Errorf("constructing kong configuration: %w", err)
	}

	if len(targetState.CustomEntities) > 0 {
		unmarshaledConfig := map[string]any{}
		if err := json.Unmarshal(config, &unmarshaledConfig); err != nil {
			return 0, fmt.Errorf("unmarshaling config for adding custom entities: %w", err)
		}
		for entityType, entities := range targetState.CustomEntities {
			unmarshaledConfig[entityType] = entities
			s.logger.V(logging.DebugLevel).Info("Filled custom entities", "entity_type", entityType)
		}
		config, err = json.Marshal(unmarshaledConfig)
		if err != nil {
			return 0, fmt.Errorf("constructing kong configuration again with custom entities: %w", err)
		}
	}

	configSize := len(config)
	if reloadConfigErr := s.configService.ReloadDeclarativeRawConfig(
		ctx,
		bytes.NewReader(config),
		true,
		true,
	); reloadConfigErr != nil {

		// If the returned error is an APIError with a 400 status code, we can try to parse the response body to get the
		// resource errors and produce an UpdateError with them.
		var apiError *kong.APIError
		if errors.As(reloadConfigErr, &apiError) && apiError.Code() == http.StatusBadRequest {
			resourceErrors, parseErr := parseFlatEntityErrors(apiError.Raw(), s.logger)
			if parseErr != nil {
				return 0, fmt.Errorf("failed to parse flat entity errors from error response: %w", parseErr)
			}

			for _, resourceError := range resourceErrors {
				s.logger.V(logging.DebugLevel).Info("Resource error", "resource_error", resourceError)
			}

			return 0, NewUpdateErrorWithResponseBody(
				apiError.Raw(),
				configSize,
				resourceErrorsToResourceFailures(resourceErrors, s.logger),
				reloadConfigErr,
			)
		}
		// ...otherwise, we return the original one.
		return 0, fmt.Errorf("failed to reload declarative configuration: %w", reloadConfigErr)
	}
	return configSize, nil
}

func (s UpdateStrategyInMemory) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDBLess
}

func (s UpdateStrategyInMemory) Type() string {
	return "InMemory"
}

type InMemoryClient interface {
	BaseRootURL() string
	ReloadDeclarativeRawConfig(ctx context.Context, config io.Reader, checkHash bool, flattenErrors bool) ([]byte, error)
}
