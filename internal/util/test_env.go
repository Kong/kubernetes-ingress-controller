package util

import (
	"os"
	"time"
)

// ControllersCacheSyncTimeout evaluates every controller's `controller.Opts.CacheSyncTimeout`.
// By default, it's set to 2 minutes. Tweaking it is possible by setting TEST_KONG_CONTROLLERS_CACHE_SYNC_TIMEOUT
// environment variable and is meant to be used in tests only.
func ControllersCacheSyncTimeout() time.Duration {
	if v := os.Getenv("TEST_KONG_CONTROLLERS_CACHE_SYNC_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}

	defaultCacheSyncTimeout := 2 * time.Minute
	return defaultCacheSyncTimeout
}
