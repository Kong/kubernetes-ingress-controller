package util

import (
	"fmt"
	"os"

	"github.com/bluele/gcache"
)

var goCache gcache.Cache

func InitCache() error {
	if goCache != nil {
		msg := "memory cache already initialized."
		fmt.Println(msg)
		return nil
	}

	goCache = gcache.New(1000).
		LRU().
		Build()

	if goCache == nil {
		msg := "failed to initialize memory cache"
		fmt.Println(msg)
		return fmt.Errorf(msg)
	}

	if os.Getenv("TEST_MODE") == "true" {
		fmt.Printf("test mode : every restart removes all key-value pairs from cache.")
		goCache.Purge()
	}

	return nil
}

func SetValue(key, value interface{}) error {
	if goCache == nil {
		msg := "memory cache is not connected"
		fmt.Println(msg)
		return fmt.Errorf(msg)
	}

	if err := goCache.Set(key, value); err != nil {
		errMsg := fmt.Sprintf("failed to cache key %s value %s", key, value)
		fmt.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	okMsg := fmt.Sprintf("sucessfully cached key %s value %d", key, value)
	fmt.Println(okMsg)
	return nil
}

func GetValue(key interface{}) (interface{}, error) {
	if goCache == nil {
		msg := "memory cache is not connected"
		fmt.Println(msg)
		return nil, fmt.Errorf(msg)
	}
	value, err := goCache.Get(key)
	if err != nil {
		errMsg := fmt.Sprintf("failed to get cache key %s err %v", key, err)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	okMsg := fmt.Sprintf("successfully get cache key %s value %d", key, value)
	fmt.Println(okMsg)
	return value, nil
}
