package sendconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

type ConfigService interface {
	ReloadDeclarativeRawConfig(
		ctx context.Context,
		config io.Reader,
		checkHash bool,
		flattenErrors bool,
	) ([]byte, error)
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

func (s UpdateStrategyInMemory) Update(ctx context.Context, targetState ContentWithHash) (
	err error,
	resourceErrors []ResourceError,
	rawErrBody []byte,
	resourceErrorsParseErr error,
) {
	dblessConfig := s.configConverter.Convert(targetState.Content)
	config, err := json.Marshal(dblessConfig)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err), nil, nil, nil
	}

	if errBody, err := s.configService.ReloadDeclarativeRawConfig(ctx, bytes.NewReader(config), true, true); err != nil {
		resourceErrors, parseErr := parseFlatEntityErrors(errBody, s.logger)
		return err, resourceErrors, errBody, parseErr
	}

	return nil, nil, nil, nil
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
