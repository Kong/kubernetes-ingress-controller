package sendconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/kong/deck/file"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
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

// UpdateStrategyInMemory implements the UpdateStrategy interface. It updates Kong's data-plane
// configuration using its `POST /config` endpoint that is used by ConfigService.ReloadDeclarativeRawConfig.
type UpdateStrategyInMemory struct {
	configService               ConfigService
	preserveNullsinPluginConfig bool
	log                         logrus.FieldLogger
}

func NewUpdateStrategyInMemory(
	configService ConfigService,
	preserveNullsinPluginConfig bool,
	log logrus.FieldLogger,
) UpdateStrategyInMemory {
	return UpdateStrategyInMemory{
		configService:               configService,
		preserveNullsinPluginConfig: preserveNullsinPluginConfig,
		log:                         log,
	}
}

func (s UpdateStrategyInMemory) Update(ctx context.Context, targetState *file.Content) (
	err error,
	resourceErrors []ResourceError,
	resourceErrorsParseErr error,
) {
	// Kong will error out if this is set
	targetState.Info = nil

	if !s.preserveNullsinPluginConfig {
		deckgen.CleanUpNullsInPluginConfigs(targetState)
	}

	config, err := json.Marshal(targetState)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err), nil, nil
	}

	flattenErrors := shouldUseFlattenedErrors()
	if errBody, err := s.configService.ReloadDeclarativeRawConfig(ctx, bytes.NewReader(config), true, flattenErrors); err != nil {
		resourceErrors, parseErr := parseFlatEntityErrors(errBody, s.log)
		return err, resourceErrors, parseErr
	}

	return nil, nil, nil
}

func (s UpdateStrategyInMemory) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDBLess
}

// shouldUseFlattenedErrors verifies whether we should pass flatten errors flag to ReloadDeclarativeRawConfig.
// Kong's API library combines KVs in the request body (the config) and query string (check hash, flattened)
// into a single set of parameters: https://github.com/Kong/go-kong/pull/271#issuecomment-1416212852
// KIC therefore must _not_ request flattened errors on versions that do not support it, as otherwise Kong
// will interpret the query string toggle as part of the config, and will reject it, as "flattened_errors" is
// not a valid config key. KIC only sends this query parameter if Kong is 3.2 or higher.
func shouldUseFlattenedErrors() bool {
	return !versions.GetKongVersion().MajorMinorOnly().LTE(versions.FlattenedErrorCutoff)
}

type InMemoryClient interface {
	BaseRootURL() string
	ReloadDeclarativeRawConfig(ctx context.Context, config io.Reader, checkHash bool, flattenErrors bool) ([]byte, error)
}
