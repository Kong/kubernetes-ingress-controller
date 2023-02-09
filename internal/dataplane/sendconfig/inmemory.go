package sendconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/kong/deck/file"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
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
	configService ConfigService
}

func NewUpdateStrategyInMemory(
	configService ConfigService,
) UpdateStrategyInMemory {
	return UpdateStrategyInMemory{
		configService: configService,
	}
}

func (s UpdateStrategyInMemory) Update(
	ctx context.Context,
	targetState *file.Content,
) error {
	// Kong will error out if this is set
	targetState.Info = nil

	// Kong errors out if `null`s are present in `config` of plugins
	deckgen.CleanUpNullsInPluginConfigs(targetState)

	config, err := json.Marshal(targetState)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err)
	}

	if _, err = s.configService.ReloadDeclarativeRawConfig(ctx, bytes.NewReader(config), true, false); err != nil {
		return err
	}

	return nil
}

func (s UpdateStrategyInMemory) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDBLess
}
