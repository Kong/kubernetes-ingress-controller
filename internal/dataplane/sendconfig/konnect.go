package sendconfig

//
// func syncWithKonnect(
// 	ctx context.Context,
// 	targetContent *file.Content,
// 	kongConfig *Kong,
// 	skipCACertificates bool,
// ) error {
// 	address := os.Getenv("KONG_KONNECT_ADDRESS")
// 	if address == "" {
// 		address = defaultKonnectAPIAddress
// 	}
// 	rg := os.Getenv("KONG_KONNECT_RG")
// 	c, err := NewKongClientForKonnect(KonnectConfig{
// 		Token:        os.Getenv("KONG_KONNECT_TOKEN"),
// 		Address:      address,
// 		RuntimeGroup: rg,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create kong client for konnect: %w", err)
// 	}
//
// 	dumpConfig := dump.Config{
// 		SkipCACerts:         skipCACertificates,
// 		KonnectRuntimeGroup: rg,
// 	}
//
// 	cs, err := currentState(ctx, c, dumpConfig)
// 	if err != nil {
// 		return fmt.Errorf("could not build current state: %w", err)
// 	}
//
// 	ts, err := targetState(ctx, targetContent, cs, kongConfig.Version, c, dumpConfig)
// 	if err != nil {
// 		return fmt.Errorf("could not build target state: %w", err)
// 	}
//
// 	syncer, err := diff.NewSyncer(diff.SyncerOpts{
// 		CurrentState:    cs,
// 		TargetState:     ts,
// 		KongClient:      c,
// 		SilenceWarnings: false,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("creating a new syncer for konnect: %w", err)
// 	}
//
// 	_, errs := syncer.Solve(ctx, kongConfig.Concurrency, false)
// 	if errs != nil {
// 		return deckutils.ErrArray{Errors: errs}
// 	}
//
// 	return nil
// }
