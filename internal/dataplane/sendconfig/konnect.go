package sendconfig

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kong/deck/diff"
	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
)

type KonnectConfig struct {
	Token        string
	RuntimeGroup string
	Address      string
}

const defaultKonnectAPIAddress = "https://api.konghq.com"

func NewKongClientForKonnect(konnectConfig KonnectConfig) (*kong.Client, error) {
	httpClient := deckutils.HTTPClient()
	if konnectConfig.Address == "" {
		konnectConfig.Address = defaultKonnectAPIAddress
	}
	if konnectConfig.Token == "" {
		return nil, errors.New("empty konnect token provided")
	}

	// TODO: when https://github.com/Kong/koko-private/pull/1384 is done we should switch to using KIC-specific
	// endpoints and pass client cert in headers.
	headers := []string{"Authorization:Bearer " + konnectConfig.Token}
	return deckutils.GetKongClient(deckutils.KongClientConfig{
		Address:    konnectConfig.Address + "/konnect-api/api/runtime_groups/" + konnectConfig.RuntimeGroup,
		HTTPClient: httpClient,
		Debug:      false,
		Headers:    headers,
		Retryable:  true,
	})
}

func syncWithKonnect(
	ctx context.Context,
	targetContent *file.Content,
	kongConfig *Kong,
	skipCACertificates bool,
) error {
	address := os.Getenv("KONG_KONNECT_ADDRESS")
	if address == "" {
		address = defaultKonnectAPIAddress
	}
	rg := os.Getenv("KONG_KONNECT_RG")
	c, err := NewKongClientForKonnect(KonnectConfig{
		Token:        os.Getenv("KONG_KONNECT_TOKEN"),
		Address:      address,
		RuntimeGroup: rg,
	})
	if err != nil {
		return fmt.Errorf("failed to create kong client for konnect: %w", err)
	}

	dumpConfig := dump.Config{
		SkipCACerts:         skipCACertificates,
		KonnectRuntimeGroup: rg,
	}

	cs, err := currentState(ctx, c, dumpConfig)
	if err != nil {
		return fmt.Errorf("could not build current state: %w", err)
	}

	ts, err := targetState(ctx, targetContent, cs, kongConfig.Version, c, dumpConfig)
	if err != nil {
		return fmt.Errorf("could not build target state: %w", err)
	}

	syncer, err := diff.NewSyncer(diff.SyncerOpts{
		CurrentState:    cs,
		TargetState:     ts,
		KongClient:      c,
		SilenceWarnings: false,
	})
	if err != nil {
		return fmt.Errorf("creating a new syncer for konnect: %w", err)
	}

	_, errs := syncer.Solve(ctx, kongConfig.Concurrency, false)
	if errs != nil {
		return deckutils.ErrArray{Errors: errs}
	}

	return nil
}
