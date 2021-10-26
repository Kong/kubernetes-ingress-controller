package sendconfig

import (
	"context"
	"time"

	"github.com/kong/deck/file"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Sendconfig - Workflow Functions
// -----------------------------------------------------------------------------

// UpdateKongAdminSimple is a helper function for the most common usage of PerformUpdate() with only minimal
// upfront configuration required. This function is specialized and highly opinionated.
//
// If you're implementation needs to expand on the configuration and usage of the following inner components:
//
//   - store.Storer
//   - kongstate.Kong
//   - deckgen.ToDeckContent()
//   - sendconfig.PerformUpdate()
//
// Or any other encapsulated components this function makes all of that opaque to the caller.
// Treat this function as a very specific "workflow" to update the Kong Admin API,
// and use it as a reference to implement the workflow you need.
func UpdateKongAdminSimple(ctx context.Context,
	lastConfigSHA []byte,
	cache *store.CacheStores,
	ingressClassName string,
	deprecatedLogger logrus.FieldLogger,
	kongConfig Kong,
	enableReverseSync bool,
	diagnostic util.ConfigDumpDiagnostic,
	proxyRequestTimeout time.Duration,
	promMetrics *metrics.CtrlFuncMetrics,
) ([]byte, error) {
	// build the kongstate object from the Kubernetes objects in the storer
	storer := store.New(*cache, ingressClassName, false, false, false, deprecatedLogger)
	kongstate, err := parser.Build(deprecatedLogger, storer)
	if err != nil {
		promMetrics.TranslationCount.With(prometheus.Labels{
			metrics.SuccessKey: metrics.SuccessFalse,
		}).Inc()
		return nil, err
	}
	promMetrics.TranslationCount.With(prometheus.Labels{
		metrics.SuccessKey: metrics.SuccessTrue,
	}).Inc()
	var diagnosticConfig *file.Content

	// generate the deck configuration to be applied to the admin API
	targetConfig := deckgen.ToDeckContent(ctx,
		deprecatedLogger, kongstate,
		kongConfig.PluginSchemaStore, kongConfig.FilterTags)

	// generate diagnostic configuration if enabled
	// "diagnostic" will be empty if --dump-config is not set
	if diagnostic != (util.ConfigDumpDiagnostic{}) {
		if !diagnostic.DumpsIncludeSensitive {
			redactedConfig := deckgen.ToDeckContent(ctx,
				deprecatedLogger, kongstate.SanitizedCopy(),
				kongConfig.PluginSchemaStore, kongConfig.FilterTags)
			diagnosticConfig = redactedConfig
		} else {
			diagnosticConfig = targetConfig
		}
	}

	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(ctx, proxyRequestTimeout)
	defer cancel()

	configSHA, err := PerformUpdate(timedCtx,
		deprecatedLogger, &kongConfig,
		kongConfig.InMemory, enableReverseSync,
		targetConfig, kongConfig.FilterTags, nil, lastConfigSHA, false, promMetrics,
	)
	if err != nil {
		if diagnostic != (util.ConfigDumpDiagnostic{}) {
			select {
			case diagnostic.Configs <- util.ConfigDump{Failed: true, Config: *diagnosticConfig}:
				deprecatedLogger.Debug("shipping config to diagnostic server")
			default:
				deprecatedLogger.Error("config diagnostic buffer full, dropping diagnostic config")
			}
		}
		return nil, err
	}
	if diagnostic != (util.ConfigDumpDiagnostic{}) {
		select {
		case diagnostic.Configs <- util.ConfigDump{Failed: false, Config: *diagnosticConfig}:
			deprecatedLogger.Debug("shipping config to diagnostic server")
		default:
			deprecatedLogger.Error("config diagnostic buffer full, dropping diagnostic config")
		}
	}

	return configSHA, nil
}
