package util

import (
	"os"
	"time"
)

func ControllersCacheSyncTimeout() time.Duration {
	if v := os.Getenv("TEST_KONG_CONTROLLERS_CACHE_SYNC_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}

	defaultCacheSyncTimeout := 2 * time.Minute
	return defaultCacheSyncTimeout
}
