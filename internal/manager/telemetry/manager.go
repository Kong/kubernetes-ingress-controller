package telemetry

import (
	"context"
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

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	splunkEndpoint  = "kong-hf.konghq.com:61833"
	telemetryPeriod = time.Second * 3

	prefix      = "kic"
	SignalStart = prefix + "-start"
	SignalPing  = prefix + "-ping"
)

type Payload = types.ProviderReport

// CreateManager creates telemetry manager using the provided rest.Config.
func CreateManager(
	ctx context.Context,
	restConfig *rest.Config,
	fixedPayload Payload,
	featureGates map[string]bool,
	meshDetection bool,
	publishServiceNN apitypes.NamespacedName,
) (telemetry.Manager, error) {
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

	m, err := createManager(ctx, k, dyn, cl, fixedPayload, featureGates, meshDetection, publishServiceNN,
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
	ctx context.Context,
	k kubernetes.Interface,
	dyn dynamic.Interface,
	cl client.Client,
	fixedPayload Payload,
	featureGates map[string]bool,
	meshDetection bool,
	publishServiceNN apitypes.NamespacedName,
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
		if meshDetection {
			podInfo, err := util.GetPodDetails(ctx, k)
			if err != nil {
				// return nil, fmt.Errorf("failed to get pod details: %w", err)
			} else {
				podNN := apitypes.NamespacedName{
					Namespace: podInfo.Namespace,
					Name:      podInfo.Name,
				}

				w, err := telemetry.NewMeshDetectWorkflow(cl, podNN, publishServiceNN)
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
			p, err := provider.NewFixedValueProvider("feature-gates", featureGatesToTelemetryPayload(featureGates))
			if err != nil {
				return nil, fmt.Errorf("failed to create fixed value provider: %w", err)
			}
			w.AddProvider(p)
		}

		m.AddWorkflow(w)
	}

	return m, nil
}

// feature-gateway=false;feature-combinedroutes=false.
func featureGatesToTelemetryPayload(featureGates map[string]bool) types.ProviderReport {
	report := make(types.ProviderReport)
	for k, v := range featureGates {
		key := fmt.Sprintf("feature-%s", strings.ToLower(k))
		report[types.ProviderReportKey(key)] = v
	}
	return report
}

// signal=kic-ping
// feature-combinedroutes=true
// feature-gateway=true
// feature-gatewayalpha=true
// feature-knative=false
//
// hn=P-Maek-MBP
// uptime=3
// v=NOT_SET
// k8s_arch=linux/arm64
// k8s_provider=kind
// k8sv=v1.26.0
// k8sv_semver=v1.26.0
// k8s_gateways_count=0
// k8s_nodes_count=1
// k8s_pods_count=14
// k8s_services_count=6;

// signal=kic-ping
// uptime=3600
// v=2.8.0
// k8sv=v1.24.8
// kv=2.8.3
// db=off
// id=4ef7d6df-553e-4c02-be63-e9e4f3a83014
// hn=ingress-kong-8479884f54-pmfvw
//
// feature-knative=false
// feature-gateway=true
// feature-gatewayalpha=false
// feature-combinedroutes=true;;;mdist="all113"
