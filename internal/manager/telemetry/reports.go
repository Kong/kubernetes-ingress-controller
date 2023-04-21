package telemetry

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"k8s.io/client-go/rest"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/metadata"
)

// GatewayClientsProvider is an interface that provides clients for the currently discovered Gateway instances.
type GatewayClientsProvider interface {
	GatewayClients() []*adminapi.Client
	GatewayClientsCount() int
}

type InstanceIDProvider interface {
	GetID() uuid.UUID
}

// SetupAnonymousReports sets up and starts the anonymous reporting and returns
// a cleanup function and an error.
// The caller is responsible to call the returned function - when the returned
// error is not nil - to stop the reports sending.
func SetupAnonymousReports(
	ctx context.Context,
	kubeCfg *rest.Config,
	clientsProvider GatewayClientsProvider,
	rv ReportValues,
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

	// gather versioning information from the kong client
	kongVersion, ok := root["version"].(string)
	if !ok {
		return nil, fmt.Errorf("malformed Kong version found in Kong client root")
	}
	cfg, ok := root["configuration"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("malformed Kong configuration found in Kong client root")
	}
	kongDB, ok := cfg["database"].(string)
	if !ok {
		return nil, fmt.Errorf("malformed database configuration found in Kong client root")
	}

	fixedPayload := Payload{
		"v":  metadata.Release,
		"kv": kongVersion,
		"db": kongDB,
		"id": instanceIDProvider.GetID(), // universal unique identifier for this system
	}

	tMgr, err := CreateManager(kubeCfg, clientsProvider, fixedPayload, rv)
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
