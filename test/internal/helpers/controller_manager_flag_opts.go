package helpers

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type ControllerManagerOpt func([]string) []string

// ControllerManagerOptAdditionalWatchNamespace adds the provided namespace to
// controller manager's watch namespaces if it's not there yet.
func ControllerManagerOptAdditionalWatchNamespace(ns string) ControllerManagerOpt {
	return func(args []string) []string {
		// Check if watch namespace is set at all.
		wn, idx, ok := lo.FindIndexOf(args, func(arg string) bool {
			return strings.HasPrefix(arg, "--watch-namespace=")
		})
		// If it isn't then append new watch namespace flag with the provided namespace.
		if !ok {
			args = append(args, fmt.Sprintf("--watch-namespace=%s", ns))
			return args
		}

		// If it is, then check the existing value (split by a comma) if it's in the list.
		v := strings.TrimPrefix(wn, "--watch-namespace=")
		namespaces := strings.Split(v, ",")
		if !lo.Contains(namespaces, ns) {
			// Replace the existing value with the new one.
			args[idx] = fmt.Sprintf("%s,%s", wn, ns)
		}

		return args
	}
}

// ControllerManagerOptFlagUseLastValidConfigForFallback sets --use-last-valid-config-for-fallback
// controller manager flag.
func ControllerManagerOptFlagUseLastValidConfigForFallback() ControllerManagerOpt {
	return func(args []string) []string {
		return append(args, "--use-last-valid-config-for-fallback")
	}
}
