package kongconfig

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

// ValidateRoots checks if all provided kong roots are the same given that we
// only care about the fact that the following fields are the same:
// - database setting
// - kong version.
func ValidateRoots(roots []Root, skipCACerts bool) (string, kong.Version, error) {
	if err := multierr.Combine(lo.Map(roots, validateRootFunc(skipCACerts))...); err != nil {
		return "", kong.Version{}, fmt.Errorf("failed to validate kong Roots: %w", err)
	}

	uniqs := lo.UniqBy(roots, getRootKeyFunc(skipCACerts))
	if len(uniqs) != 1 {
		return "", kong.Version{},
			fmt.Errorf("there should only be one dbmode:version combination across configured kong instances while there are (%d): %v", len(uniqs), uniqs)
	}

	dbMode, err := DBModeFromRoot(uniqs[0])
	if err != nil {
		return "", kong.Version{}, err
	}

	kongVersion, err := KongVersionFromRoot(uniqs[0])
	if err != nil {
		return "", kong.Version{}, err
	}

	return dbMode, kongVersion, nil
}

func DBModeFromRoot(r Root) (string, error) {
	rootConfig, ok := r["configuration"].(map[string]any)
	if !ok {
		return "", fmt.Errorf(
			"invalid root configuration, expected a map[string]any got %T",
			r["configuration"],
		)
	}

	dbMode, ok := rootConfig["database"].(string)
	if !ok {
		return "", fmt.Errorf(
			"invalid database configuration, expected a string got %t",
			rootConfig["database"],
		)
	}
	return dbMode, nil
}

func KongVersionFromRoot(r Root) (kong.Version, error) {
	v := kong.VersionFromInfo(r)
	kv, err := kong.ParseSemanticVersion(v)
	if err != nil {
		return kong.Version{}, fmt.Errorf("could not parse Kong version (%s): %w", v, err)
	}
	return kv, nil
}

// Root represents Kong Gateway configuration root.
type Root map[string]any

func (kr Root) Validate(skipCACerts bool) error {
	dbMode, err := DBModeFromRoot(kr)
	if err != nil {
		return err
	}

	if err := validateDBMode(dbMode, skipCACerts); err != nil {
		return err
	}

	if _, err = KongVersionFromRoot(kr); err != nil {
		return err
	}

	return nil
}

func (kr Root) Key(skipCACerts bool) string {
	dbMode, err := DBModeFromRoot(kr)
	if err != nil {
		return ""
	}

	if err := validateDBMode(dbMode, skipCACerts); err != nil {
		return ""
	}

	v := kong.VersionFromInfo(kr)
	_, err = kong.ParseSemanticVersion(v)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", dbMode, v)
}

func validateRootFunc(skipCACerts bool) func(Root, int) error {
	return func(r Root, _ int) error {
		return r.Validate(skipCACerts)
	}
}

// getRootKeyFunc generates a key for mapping a kong root to a string.
// It assumes that the provided configuration root has already been verified with a validation
// function return by validateRootFunc.
func getRootKeyFunc(skipCACerts bool) func(Root) string {
	return func(r Root) string {
		return r.Key(skipCACerts)
	}
}

// validateDBMode validates the provided dbMode string.
func validateDBMode(dbMode string, skipCACerts bool) error {
	switch dbMode {
	case "off", "":
		if skipCACerts {
			return fmt.Errorf("--skip-ca-certificates is not available for use with DB-less Kong instances")
		}
	case "postgres":
		return nil
	case "cassandra":
		return fmt.Errorf("Cassandra-backed deployments of Kong managed by the ingress controller are no longer supported; you must migrate to a Postgres-backed or DB-less deployment")
	default:
		return fmt.Errorf("%s is not a supported database backend", dbMode)
	}
	return nil
}

// GetRoots fetches all the configuration roots using the provided kong clients.
func GetRoots(
	ctx context.Context,
	setupLog logr.Logger,
	retries uint,
	retryDelay time.Duration,
	kongClients []*adminapi.Client,
) ([]Root, error) {
	var (
		roots []Root
		lock  sync.Mutex
	)

	eg, ctx := errgroup.WithContext(ctx)

	for _, client := range kongClients {
		client := client
		eg.Go(func() error {
			return retry.Do(
				func() error {
					root, err := client.AdminAPIClient().Root(ctx)
					// Abort if the provided context has been cancelled.
					if errors.Is(err, context.Canceled) {
						return retry.Unrecoverable(err)
					}
					if err != nil {
						return err
					}

					lock.Lock()
					roots = append(roots, root)
					lock.Unlock()
					return nil
				},
				retry.Attempts(retries),
				retry.Delay(retryDelay),
				retry.DelayType(retry.FixedDelay),
				retry.LastErrorOnly(true),
				retry.OnRetry(func(n uint, err error) {
					setupLog.Info("Retrying kong admin api client call after error",
						"retries", fmt.Sprintf("%d/%d", n, retries),
						"error", err.Error(),
					)
				}),
			)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return roots, nil
}
