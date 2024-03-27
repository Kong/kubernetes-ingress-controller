package sendconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/blang/semver/v4"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
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
	log             logrus.FieldLogger
	version         semver.Version
}

func NewUpdateStrategyInMemory(
	configService ConfigService,
	configConverter ContentToDBLessConfigConverter,
	log logrus.FieldLogger,
	version semver.Version,
) UpdateStrategyInMemory {
	return UpdateStrategyInMemory{
		configService:   configService,
		configConverter: configConverter,
		log:             log,
		version:         version,
	}
}

func (s UpdateStrategyInMemory) Update(ctx context.Context, targetState ContentWithHash) (
	err error,
	resourceErrors []ResourceError,
	resourceErrorsParseErr error,
) {
	dblessConfig := s.configConverter.Convert(targetState.Content)
	config, err := json.Marshal(dblessConfig)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err), nil, nil
	}

	flattenErrors := shouldUseFlattenedErrors(s.version)
	if errBody, err := s.configService.ReloadDeclarativeRawConfig(ctx, bytes.NewReader(config), true, flattenErrors); err != nil {
		resourceErrors, parseErr := parseFlatEntityErrors(errBody, s.log)
		return err, resourceErrors, parseErr
	}

	return nil, nil, nil
}

func (s UpdateStrategyInMemory) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDBLess
}

func (s UpdateStrategyInMemory) Type() string {
	return "InMemory"
}

// shouldUseFlattenedErrors verifies whether we should pass flatten errors flag to ReloadDeclarativeRawConfig.
// Kong's API library combines KVs in the request body (the config) and query string (check hash, flattened)
// into a single set of parameters: https://github.com/Kong/go-kong/pull/271#issuecomment-1416212852
// KIC therefore must _not_ request flattened errors on versions that do not support it, as otherwise Kong
// will interpret the query string toggle as part of the config, and will reject it, as "flattened_errors" is
// not a valid config key.
// KIC only sends this query parameter if Kong is 3.2 or higher.
func shouldUseFlattenedErrors(version semver.Version) bool {
	return version.GTE(versions.FlattenedErrorCutoff)
}

type InMemoryClient interface {
	BaseRootURL() string
	ReloadDeclarativeRawConfig(ctx context.Context, config io.Reader, checkHash bool, flattenErrors bool) ([]byte, error)
}
