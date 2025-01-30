package manager

import (
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/config"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// Note: NewConfig is not implemented in the `pkg/manager/config` package to avoid cyclic dependencies:
// for now, it depends on `internal/manager/config` to set defaults using CLI flags parsing.

// NewConfig is used to create a new configuration object with default values.
// Values can be overridden by passing `managercfg.Opt` options.
func NewConfig(opts ...managercfg.Opt) (managercfg.Config, error) {
	cfg, err := NewDefaultConfig()
	if err != nil {
		return managercfg.Config{}, fmt.Errorf("failed to create default manager config: %w", err)
	}

	// Override default values with the provided options.
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg, nil
}

func NewDefaultConfig() (managercfg.Config, error) {
	// Set default values relying on CLI flags parsing.
	cliCfg := config.NewCLIConfig()
	flags := cliCfg.FlagSet()
	if err := flags.Parse([]string{}); err != nil {
		return managercfg.Config{}, fmt.Errorf("failed to parse CLI flags: %w", err)
	}

	// Set default feature gates values as they're not populated by the CLI flags parsing.
	cliCfg.FeatureGates = managercfg.GetFeatureGatesDefaults()

	return *cliCfg.Config, nil
}
