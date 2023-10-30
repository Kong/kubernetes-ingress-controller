package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"k8s.io/client-go/rest"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/utils/kongconfig"
)

// GatewayClientsProvider is an interface that provides clients for the currently discovered Gateway instances.
type GatewayClientsProvider interface {
	GatewayClients() []*adminapi.Client
	GatewayClientsCount() int
}

type InstanceIDProvider interface {
	GetID() uuid.UUID
}

const (
	splunkEndpoint                   = "kong-hf.konghq.com:61833"
	splunkEndpointInsecureSkipVerify = false
	telemetryPeriod                  = time.Hour
)

type ReportConfig struct {
	SplunkEndpoint                   string
	SplunkEndpointInsecureSkipVerify bool
	TelemetryPeriod                  time.Duration
	ReportValues                     ReportValues
}

// SetupAnonymousReports sets up and starts the anonymous reporting and returns
// a cleanup function and an error.
// The caller is responsible to call the returned function - when the returned
// error is not nil - to stop the reports sending.
func SetupAnonymousReports(
	ctx context.Context,
	logger logr.Logger,
	kubeCfg *rest.Config,
	clientsProvider GatewayClientsProvider,
	reportCfg ReportConfig,
	instanceIDProvider InstanceIDProvider,
) (func(), error) {
	// if anonymous reports are enabled this helps provide Kong with insights about usage of the ingress controller
	// which is non-sensitive and predominantly informs us of the controller and cluster versions in use.
	// This data helps inform us what versions, features, e.t.c. end-users are actively using which helps to inform
	// our prioritization of work and we appreciate when our end-users provide them, however if you do feel
	// uncomfortable and would rather turn them off run the controller with the "--anonymous-reports false" flag.

	// This now only uses the first instance for telemetry reporting.
	// That's fine because we allow for now only 1 set of version and db setting
	// throughout all Kong instances that 1 KIC instance configures.
	//
	// When we change that and decide to allow heterogeneous Kong instances to be
	// configured by 1 KIC instance then this will have to change.
	//
	// https://github.com/Kong/kubernetes-ingress-controller/issues/3589
	root, err := clientsProvider.GatewayClients()[0].AdminAPIClient().Root(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Kong root config data: %w", err)
	}

	// Gather versioning information from the kong client
	kongVersion := kong.VersionFromInfo(root)
	if kongVersion == "" {
		return nil, fmt.Errorf("malformed Kong version found in Kong client root")
	}
	kongDB, err := kongconfig.DBModeFromRoot(root)
	if err != nil {
		return nil, err
	}
	routerFlavor, err := kongconfig.RouterFlavorFromRoot(root)
	if err != nil {
		return nil, err
	}

	fixedPayload := Payload{
		"v":  metadata.Release,
		"kv": kongVersion,
		"db": kongDB,
		"rf": routerFlavor,
		"id": instanceIDProvider.GetID(), // universal unique identifier for this system
	}

	// Use defaults when not specified.
	if reportCfg.SplunkEndpoint == "" {
		reportCfg.SplunkEndpoint = splunkEndpoint
	}
	if !reportCfg.SplunkEndpointInsecureSkipVerify {
		reportCfg.SplunkEndpointInsecureSkipVerify = splunkEndpointInsecureSkipVerify
	}
	if reportCfg.TelemetryPeriod == 0 {
		reportCfg.TelemetryPeriod = telemetryPeriod
	}

	tMgr, err := CreateManager(logger, kubeCfg, clientsProvider, fixedPayload, reportCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create anonymous reports manager: %w", err)
	}

	if err := tMgr.Start(); err != nil {
		return nil, fmt.Errorf("anonymous reports failed to start: %w", err)
	}

	if err := tMgr.TriggerExecute(ctx, SignalStart); err != nil {
		return tMgr.Stop, fmt.Errorf("failed to trigger telemetry report during start: %w", err)
	}

	return tMgr.Stop, nil
}
