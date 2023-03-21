package telemetry

import (
	"fmt"
	"strings"
	"time"

	"github.com/bombsimon/logrusr/v2"
	"github.com/kong/kubernetes-telemetry/pkg/forwarders"
	"github.com/kong/kubernetes-telemetry/pkg/provider"
	"github.com/kong/kubernetes-telemetry/pkg/serializers"
	"github.com/kong/kubernetes-telemetry/pkg/telemetry"
	"github.com/kong/kubernetes-telemetry/pkg/types"
	"github.com/sirupsen/logrus"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/telemetry/workflows"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	splunkEndpoint  = "kong-hf.konghq.com:61833"
	telemetryPeriod = time.Hour

	prefix      = "kic"
	SignalStart = prefix + "-start"
	SignalPing  = prefix + "-ping"
)

type Payload = types.ProviderReport

type ReportValues struct {
	FeatureGates                   map[string]bool
	MeshDetection                  bool
	PublishServiceNN               apitypes.NamespacedName
	KonnectSyncEnabled             bool
	GatewayServiceDiscoveryEnabled bool
}

// CreateManager creates telemetry manager using the provided rest.Config.
func CreateManager(restConfig *rest.Config, gatewaysCounter workflows.DiscoveredGatewaysCounter, fixedPayload Payload, rv ReportValues) (telemetry.Manager, error) {
	logger := logrusr.New(logrus.New())

	k, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create client-go kubernetes client: %w", err)
	}
	cl, err := client.New(restConfig, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create controller-runtime kubernetes client: %w", err)
	}
	dyn := dynamic.New(k.Discovery().RESTClient())

	m, err := createManager(k, dyn, cl, gatewaysCounter, fixedPayload, rv,
		telemetry.OptManagerPeriod(telemetryPeriod),
		telemetry.OptManagerLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	tf, err := forwarders.NewTLSForwarder(splunkEndpoint, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry TLSForwarder: %w", err)
	}
	serializer := serializers.NewSemicolonDelimited()
	consumer := telemetry.NewConsumer(serializer, tf)
	if err := m.AddConsumer(consumer); err != nil {
		return nil, fmt.Errorf("failed to add TLSforwarder: %w", err)
	}

	return m, nil
}

func createManager(
	k kubernetes.Interface,
	dyn dynamic.Interface,
	cl client.Client,
	gatewaysCounter workflows.DiscoveredGatewaysCounter,
	fixedPayload Payload,
	rv ReportValues,
	opts ...telemetry.OptManager,
) (telemetry.Manager, error) {
	m, err := telemetry.NewManager(
		SignalPing,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry manager: %w", err)
	}

	// Add identify cluster workflow
	{
		w, err := telemetry.NewIdentifyPlatformWorkflow(k)
		if err != nil {
			return nil, fmt.Errorf("failed to create identify platform workflow: %w", err)
		}
		m.AddWorkflow(w)
	}

	// Add cluster state workflow
	{
		w, err := telemetry.NewClusterStateWorkflow(dyn, cl.RESTMapper())
		if err != nil {
			return nil, fmt.Errorf("failed to create cluster state workflow: %w", err)
		}

		m.AddWorkflow(w)
	}

	// Add mesh detect workflow
	{
		if rv.MeshDetection {
			podNN, err := util.GetPodNN()
			// Don't fail if an err is no nil, just don't include mesh detection workflow.
			// We could probably add conditions around this, so that only the
			// part responsible for detecting the mesh that current pod is running
			// gets disabled.
			if err == nil {
				w, err := telemetry.NewMeshDetectWorkflow(cl, podNN, rv.PublishServiceNN)
				if err != nil {
					return nil, fmt.Errorf("failed to create mesh detect workflow: %w", err)
				}

				m.AddWorkflow(w)
			}
		}
	}

	// Add state workflow
	{
		w, err := telemetry.NewStateWorkflow()
		if err != nil {
			return nil, fmt.Errorf("failed to create state workflow: %w", err)
		}

		{
			p, err := provider.NewFixedValueProvider("payload", fixedPayload)
			if err != nil {
				return nil, fmt.Errorf("failed to create fixed value provider: %w", err)
			}
			w.AddProvider(p)
		}
		{
			p, err := provider.NewFixedValueProvider("feature-gates", featureGatesToTelemetryPayload(rv.FeatureGates))
			if err != nil {
				return nil, fmt.Errorf("failed to create fixed value provider: %w", err)
			}
			w.AddProvider(p)
		}
		{
			p, err := provider.NewFixedValueProvider("feature-flags", types.ProviderReport{
				"feature-konnect-sync":              rv.KonnectSyncEnabled,
				"feature-gateway-service-discovery": rv.GatewayServiceDiscoveryEnabled,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create fixed value provider: %w", err)
			}
			w.AddProvider(p)
		}

		m.AddWorkflow(w)
	}

	if rv.GatewayServiceDiscoveryEnabled {
		w, err := workflows.NewGatewayDiscoveryWorkflow(gatewaysCounter)
		if err != nil {
			return nil, fmt.Errorf("failed to create gateway discovery workflow: %w", err)
		}
		m.AddWorkflow(w)
	}

	return m, nil
}

func featureGatesToTelemetryPayload(featureGates map[string]bool) types.ProviderReport {
	report := make(types.ProviderReport)
	for k, v := range featureGates {
		key := fmt.Sprintf("feature-%s", strings.ToLower(k))
		report[types.ProviderReportKey(key)] = v
	}
	return report
}
