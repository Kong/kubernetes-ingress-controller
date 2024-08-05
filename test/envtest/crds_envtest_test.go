//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
)

// TestGatewayAPIControllersMayBeDynamicallyStarted ensures that in case of missing CRDs installation in the
// cluster, specific controllers are not started until the CRDs are installed.
func TestGatewayAPIControllersMayBeDynamicallyStarted(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme, WithInstallGatewayCRDs(false))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, loggerHook := RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithGatewayFeatureEnabled,
		WithGatewayAPIControllers(),
		WithPublishService("ns"),
	)

	controllers := []string{
		"Gateway",
		"HTTPRoute",
		"ReferenceGrant",
		"UDPRoute",
		"TCPRoute",
		"TLSRoute",
		"GRPCRoute",
	}

	requireLogForAllControllers := func(expectedLog string) {
		require.Eventually(t, func() bool {
			for _, controller := range controllers {
				if !lo.ContainsBy(loggerHook.All(), func(entry observer.LoggedEntry) bool {
					return strings.Contains(entry.LoggerName, controller) && strings.Contains(entry.Message, expectedLog)
				}) {
					t.Logf("expected log %q not found for %s controller", expectedLog, controller)
					return false
				}
			}
			return true
		}, 30*time.Second, time.Millisecond*500)
	}

	const (
		expectedLogOnStartup      = "Required CustomResourceDefinitions are not installed, setting up a watch for them in case they are installed afterward"
		expectedLogOnCRDInstalled = "All required CustomResourceDefinitions are installed, setting up the controller"
	)

	t.Log("waiting for all controllers to not start due to missing CRDs")
	requireLogForAllControllers(expectedLogOnStartup)

	t.Log("installing missing CRDs")
	installGatewayCRDs(t, scheme, envcfg)

	t.Log("waiting for all controllers to start after CRDs installation")
	requireLogForAllControllers(expectedLogOnCRDInstalled)
}

// TestNoKongCRDsInstalledIsFatal ensures that in case of missing Kong CRDs installation, the manager Run() eventually
// returns an error due to cache synchronisation timeout.
func TestNoKongCRDsInstalledIsFatal(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t)
	envcfg := Setup(t, scheme, WithInstallKongCRDs(false))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := ConfigForEnvConfig(t, envcfg)

	logger := zapr.NewLogger(zap.NewNop())
	ctrl.SetLogger(logger)

	// Reducing the cache sync timeout to speed up the test.
	cfg.CacheSyncTimeout = time.Millisecond * 500
	err := manager.Run(ctx, &cfg, diagnostics.ClientDiagnostic{}, logger)
	require.ErrorContains(t, err, "timed out waiting for cache to be synced")
}

func TestCRDValidations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme)
	ctrlClient := NewControllerClient(t, scheme, envcfg)

	testCases := []struct {
		name     string
		scenario func(ctx context.Context, t *testing.T, ns string)
	}{
		{
			name: "invalid TCPIngress service name",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyTCPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Backend.ServiceName = ""
				})

				require.ErrorContains(t, err, "serviceName")
			},
		},
		{
			name: "invalid TCPIngress service port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyTCPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Backend.ServicePort = 0
				})

				require.ErrorContains(t, err, "servicePort")
			},
		},
		{
			name: "invalid TCPIngress rule port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyTCPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Port = 0
				})

				require.ErrorContains(t, err, "spec.rules[0].port")
			},
		},
		{
			name: "invalid UDPIngress service name",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyUDPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Backend.ServiceName = ""
				})

				require.ErrorContains(t, err, "serviceName")
			},
		},
		{
			name: "invalid UDPIngress service port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyUDPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Backend.ServicePort = 0
				})

				require.ErrorContains(t, err, "servicePort")
			},
		},
		{
			name: "invalid UDPIngress rule port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyUDPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Port = 0
				})

				require.ErrorContains(t, err, "spec.rules[0].port")
			},
		},
		{
			name: "KongUpstreamPolicy - only one of spec.hashOn.(cookie|header|uriCapture|queryArg) can be set",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				for i, invalidHashOn := range generateInvalidHashOns() {
					invalidHashOn := invalidHashOn
					t.Run(fmt.Sprintf("invalidHashOn[%d]", i), func(t *testing.T) {
						err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
							HashOn: &invalidHashOn,
						})
						require.ErrorContains(t, err, "Only one of spec.hashOn.(input|cookie|header|uriCapture|queryArg) can be set.")
					})
				}
			},
		},
		{
			name: "KongUpstreamPolicy - only one of spec.hashOnFallback.(header|uriCapture|queryArg) can be set",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				invalidHashOns := lo.Reject(generateInvalidHashOns(), func(hashOn kongv1beta1.KongUpstreamHash, _ int) bool {
					// Filter out Cookie which is not allowed in spec.hashOnFallback.
					return hashOn.Cookie != nil
				})
				for i, invalidHashOn := range invalidHashOns {
					invalidHashOn := invalidHashOn
					t.Run(fmt.Sprintf("invalidHashOn[%d]", i), func(t *testing.T) {
						err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
							HashOnFallback: &invalidHashOn,
						})
						require.ErrorContains(t, err, "Only one of spec.hashOnFallback.(input|header|uriCapture|queryArg) can be set.")
					})
				}
			},
		},
		{
			name: "KongUpstreamPolicy - spec.hashOn.cookie and spec.hashOn.cookiePath are set",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					Algorithm: lo.ToPtr("consistent-hashing"),
					HashOn: &kongv1beta1.KongUpstreamHash{
						Cookie:     lo.ToPtr("cookie-name"),
						CookiePath: lo.ToPtr("/"),
					},
				})
				require.NoError(t, err)
			},
		},
		{
			name: "KongUpstreamPolicy - spec.hashOn.cookie is set, spec.hashOn.cookiePath is required",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					Algorithm: lo.ToPtr("consistent-hashing"),
					HashOn: &kongv1beta1.KongUpstreamHash{
						Cookie: lo.ToPtr("cookie-name"),
					},
				})
				require.ErrorContains(t, err, "When spec.hashOn.cookie is set, spec.hashOn.cookiePath is required.")
			},
		},
		{
			name: "KongUpstreamPolicy - spec.hashOn.cookiePath is set, spec.hashOn.cookie is required",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					Algorithm: lo.ToPtr("consistent-hashing"),
					HashOn: &kongv1beta1.KongUpstreamHash{
						CookiePath: lo.ToPtr("/"),
					},
				})
				require.ErrorContains(t, err, "When spec.hashOn.cookiePath is set, spec.hashOn.cookie is required.")
			},
		},
		{
			name: "KongUpstreamPolicy - spec.hashOnFallback.cookie must not be set",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					HashOnFallback: &kongv1beta1.KongUpstreamHash{
						CookiePath: lo.ToPtr("/"),
					},
				})
				require.ErrorContains(t, err, "spec.hashOnFallback.cookiePath must not be set.")
			},
		},
		{
			name: "KongUpstreamPolicy - spec.hashOnFallback.cookiePath must not be set",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					HashOnFallback: &kongv1beta1.KongUpstreamHash{
						CookiePath: lo.ToPtr("/"),
					},
				})
				require.ErrorContains(t, err, "spec.hashOnFallback.cookiePath must not be set.")
			},
		},
		{
			name: "KongUpstreamPolicy - healthchecks.active.healthy.httpStatuses contains invalid HTTP status code",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					Healthchecks: &kongv1beta1.KongUpstreamHealthcheck{
						Active: &kongv1beta1.KongUpstreamActiveHealthcheck{
							Healthy: &kongv1beta1.KongUpstreamHealthcheckHealthy{
								HTTPStatuses: []kongv1beta1.HTTPStatus{600},
							},
						},
					},
				})
				require.ErrorContains(t, err, "should be less than or equal to 599")
			},
		},
		{
			name: "KongUpstreamPolicy - healthchecks.active.unhealthy.httpStatuses contains invalid HTTP status code",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					Healthchecks: &kongv1beta1.KongUpstreamHealthcheck{
						Active: &kongv1beta1.KongUpstreamActiveHealthcheck{
							Unhealthy: &kongv1beta1.KongUpstreamHealthcheckUnhealthy{
								HTTPStatuses: []kongv1beta1.HTTPStatus{99},
							},
						},
					},
				})
				require.ErrorContains(t, err, "should be greater than or equal to 100")
			},
		},
		{
			name: "KongUpstreamPolicy - healthchecks.passive.healthy.interval must not be set",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					Healthchecks: &kongv1beta1.KongUpstreamHealthcheck{
						Passive: &kongv1beta1.KongUpstreamPassiveHealthcheck{
							Healthy: &kongv1beta1.KongUpstreamHealthcheckHealthy{
								Interval: lo.ToPtr(10),
							},
						},
					},
				})
				require.ErrorContains(t, err, "spec.healthchecks.passive.healthy.interval must not be set.")
			},
		},
		{
			name: "KongUpstreamPolicy - healthchecks.passive.unhealthy.interval must not be set",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					Healthchecks: &kongv1beta1.KongUpstreamHealthcheck{
						Passive: &kongv1beta1.KongUpstreamPassiveHealthcheck{
							Unhealthy: &kongv1beta1.KongUpstreamHealthcheckUnhealthy{
								Interval: lo.ToPtr(10),
							},
						},
					},
				})
				require.ErrorContains(t, err, "spec.healthchecks.passive.unhealthy.interval must not be set.")
			},
		},
		{
			name: "KongUpstreamPolicy - hashOn can only be set when algorithm is set to consistent-hashing",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					HashOn: &kongv1beta1.KongUpstreamHash{
						Header: lo.ToPtr("header-name"), // Could be any of the hashOn fields.
					},
				})
				require.ErrorContains(t, err, `spec.algorithm must be set to "consistent-hashing" when spec.hashOn is set.`)
			},
		},
		{
			name: "KongUpstreamPolicy - hashOnFallback can only be set when algorithm is set to consistent-hashing",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					HashOnFallback: &kongv1beta1.KongUpstreamHash{
						Header: lo.ToPtr("header-name"), // Could be any of the hashOn fields.
					},
				})
				require.ErrorContains(t, err, `spec.algorithm must be set to "consistent-hashing" when spec.hashOnFallback is set.`)
			},
		},
		{
			name: "KongUpstreamPolicy - hashOn(Fallback).input enum is validated",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				validValues := []string{"ip", "consumer", "path"}
				for _, validValue := range validValues {
					t.Run(fmt.Sprintf("valid-value[%s]", validValue), func(t *testing.T) {
						err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
							Algorithm: lo.ToPtr("consistent-hashing"),
							HashOn: &kongv1beta1.KongUpstreamHash{
								Input: lo.ToPtr(kongv1beta1.HashInput(validValue)),
							},
							HashOnFallback: &kongv1beta1.KongUpstreamHash{
								Input: lo.ToPtr(kongv1beta1.HashInput(validValue)),
							},
						})
						require.NoError(t, err)
					})
				}

				t.Run("invalid value", func(t *testing.T) {
					err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
						Algorithm: lo.ToPtr("consistent-hashing"),
						HashOn: &kongv1beta1.KongUpstreamHash{
							Input: lo.ToPtr(kongv1beta1.HashInput("unknown-input")),
						},
						HashOnFallback: &kongv1beta1.KongUpstreamHash{
							Input: lo.ToPtr(kongv1beta1.HashInput("unknown-input-fallback")),
						},
					})
					require.ErrorContains(t, err, `spec.hashOn.input: Unsupported value: "unknown-input": supported values: "ip", "consumer", "path"`)
					require.ErrorContains(t, err, `spec.hashOnFallback.input: Unsupported value: "unknown-input-fallback": supported values: "ip", "consumer", "path"`)
				})
			},
		},
		{
			name: "KongUpstreamPolicy - algorithm enum is validated",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				validValues := []string{"consistent-hashing", "round-robin", "least-connections", "latency"}
				for _, validValue := range validValues {
					t.Run(fmt.Sprintf("valid-value[%s]", validValue), func(t *testing.T) {
						err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
							Algorithm: lo.ToPtr(validValue),
						})
						require.NoError(t, err)
					})
				}

				t.Run("invalid value", func(t *testing.T) {
					err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
						Algorithm: lo.ToPtr("unknown-algorithm"),
					})
					require.ErrorContains(t, err, `spec.algorithm: Unsupported value: "unknown-algorithm": supported values: "round-robin", "consistent-hashing", "least-connections", "latency"`)
				})
			},
		},
		{
			name: "KongUpstreamPolicy - healthcheck.(active|passive).type enum is validated",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				validValues := []string{"http", "https", "tcp", "grpc", "grpcs"}
				for _, validValue := range validValues {
					t.Run(fmt.Sprintf("valid-value[%s]", validValue), func(t *testing.T) {
						err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
							Healthchecks: &kongv1beta1.KongUpstreamHealthcheck{
								Active: &kongv1beta1.KongUpstreamActiveHealthcheck{
									Type: lo.ToPtr(validValue),
								},
								Passive: &kongv1beta1.KongUpstreamPassiveHealthcheck{
									Type: lo.ToPtr(validValue),
								},
							},
						})
						require.NoError(t, err)
					})
				}
				t.Run("invalid value", func(t *testing.T) {
					err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
						Healthchecks: &kongv1beta1.KongUpstreamHealthcheck{
							Active: &kongv1beta1.KongUpstreamActiveHealthcheck{
								Type: lo.ToPtr("unknown-type-active"),
							},
							Passive: &kongv1beta1.KongUpstreamPassiveHealthcheck{
								Type: lo.ToPtr("unknown-type-passive"),
							},
						},
					})
					require.ErrorContains(t, err, `spec.healthchecks.active.type: Unsupported value: "unknown-type-active": supported values: "http", "https", "tcp", "grpc", "grpcs"`)
					require.ErrorContains(t, err, `spec.healthchecks.passive.type: Unsupported value: "unknown-type-passive": supported values: "http", "https", "tcp", "grpc", "grpcs"`)
				})
			},
		},
		{
			name: "KongUpstreamPolicy - hashOnFallback must not be set when spec.hasOn.cookie is set",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongUpstreamPolicy(ctx, ctrlClient, ns, kongv1beta1.KongUpstreamPolicySpec{
					Algorithm: lo.ToPtr("consistent-hashing"),
					HashOn: &kongv1beta1.KongUpstreamHash{
						Cookie:     lo.ToPtr("cookie-name"),
						CookiePath: lo.ToPtr("/"),
					},
					HashOnFallback: &kongv1beta1.KongUpstreamHash{
						Header: lo.ToPtr("header-name"),
					},
				})
				require.ErrorContains(t, err, `spec.hashOnFallback must not be set when spec.hashOn.cookie is set.`)
			},
		},
		{
			name: "KongIngress - proxy is not allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongIngress(ctx, ctrlClient, ns, &kongv1.KongIngress{
					Proxy: &kongv1.KongIngressService{
						Retries: lo.ToPtr(5),
					},
				})
				require.ErrorContains(t, err, "'proxy' field is no longer supported, use Service's annotations instead")
			},
		},
		{
			name: "KongIngress - route is not allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongIngress(ctx, ctrlClient, ns, &kongv1.KongIngress{
					Route: &kongv1.KongIngressRoute{
						PreserveHost: lo.ToPtr(true),
					},
				})
				require.ErrorContains(t, err, "'route' field is no longer supported, use Ingress' annotations instead")
			},
		},
		{
			name: "KongPlugin - no config is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongPlugin(ctx, ctrlClient, ns, &kongv1.KongPlugin{
					PluginName: "key-auth",
				})
				require.NoError(t, err)
			},
		},
		{
			name: "KongPlugin - with config is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongPlugin(ctx, ctrlClient, ns, &kongv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"key_names":["apikey"]}`),
					},
				})
				require.NoError(t, err)
			},
		},
		{
			name: "KongPlugin - with configFrom is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongPlugin(ctx, ctrlClient, ns, &kongv1.KongPlugin{
					PluginName: "key-auth",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Secret: "secret-name",
							Key:    "key-name",
						},
					},
				})
				require.NoError(t, err)
			},
		},
		{
			name: "KongPlugin - with configFrom and config is rejected",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongPlugin(ctx, ctrlClient, ns, &kongv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"key_names":["apikey"]}`),
					},
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Secret: "secret-name",
							Key:    "key-name",
						},
					},
				})
				assert.Error(t, err)
				assert.ErrorContains(t, err, "Using both config and configFrom fields is not allowed.")
			},
		},
		{
			name: "KongPlugin - change plugin field is not allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				plugin := &kongv1.KongPlugin{
					PluginName: "key-auth",
				}
				err := createKongPlugin(ctx, ctrlClient, ns, plugin)
				require.NoError(t, err)
				plugin.PluginName = "cors"
				err = updateKongPlugin(ctx, ctrlClient, ns, plugin)
				assert.Error(t, err)
				assert.ErrorContains(t, err, "The plugin field is immutable")
			},
		},
		{
			name: "KongPlugin - using configPatches is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				plugin := &kongv1.KongPlugin{
					PluginName: "key-auth",
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/key_names",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Secret: "secret-name",
									Key:    "key-name",
								},
							},
						},
					},
				}
				require.NoError(t, createKongPlugin(ctx, ctrlClient, ns, plugin))
			},
		},
		{
			name: "KongPlugin - using config and configPatches is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				plugin := &kongv1.KongPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"key_name":"apikey"}`),
					},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/key_names",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Secret: "secret-name",
									Key:    "key-name",
								},
							},
						},
					},
				}
				require.NoError(t, createKongPlugin(ctx, ctrlClient, ns, plugin))
			},
		},
		{
			name: "KongPlugin - using configFrom and configPatches is not allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				plugin := &kongv1.KongPlugin{
					PluginName: "key-auth",
					ConfigFrom: &kongv1.ConfigSource{
						SecretValue: kongv1.SecretValueFromSource{
							Secret: "secret-name",
							Key:    "key-name",
						},
					},
					ConfigPatches: []kongv1.ConfigPatch{
						{
							Path: "/key_names",
							ValueFrom: kongv1.ConfigSource{
								SecretValue: kongv1.SecretValueFromSource{
									Secret: "secret-name",
									Key:    "key-name",
								},
							},
						},
					},
				}
				assert.ErrorContains(t, createKongPlugin(ctx, ctrlClient, ns, plugin),
					"Using both configFrom and configPatches fields is not allowed.")
			},
		},
		{
			name: "KongClusterPlugin - no config is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongClusterPlugin(ctx, ctrlClient, ns, &kongv1.KongClusterPlugin{
					PluginName: "key-auth",
				})
				require.NoError(t, err)
			},
		},
		{
			name: "KongClusterPlugin - with config is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongClusterPlugin(ctx, ctrlClient, ns, &kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"key_names":["apikey"]}`),
					},
				})
				require.NoError(t, err)
			},
		},
		{
			name: "KongClusterPlugin - with configFrom is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongClusterPlugin(ctx, ctrlClient, ns, &kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Secret:    "secret-name",
							Key:       "key-name",
							Namespace: "ns",
						},
					},
				})
				require.NoError(t, err)
			},
		},
		{
			name: "KongClusterPlugin - with configFrom and config is rejected",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createKongClusterPlugin(ctx, ctrlClient, ns, &kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"key_names":["apikey"]}`),
					},
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Secret:    "secret-name",
							Key:       "key-name",
							Namespace: "ns",
						},
					},
				})
				assert.Error(t, err)
				assert.ErrorContains(t, err, "Using both config and configFrom fields is not allowed.")
			},
		},
		{
			name: "KongClusterPlugin - change plugin field is not allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				plugin := &kongv1.KongClusterPlugin{
					PluginName: "key-auth",
				}
				err := createKongClusterPlugin(ctx, ctrlClient, ns, plugin)
				require.NoError(t, err)
				plugin.PluginName = "cors"
				err = updateKongClusterPlugin(ctx, ctrlClient, ns, plugin)
				assert.Error(t, err)
				assert.ErrorContains(t, err, "The plugin field is immutable")
			},
		},
		{
			name: "KongClusterPlugin - using configPatches is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				plugin := &kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/key_names",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: ns,
									Secret:    "secret-name",
									Key:       "key-name",
								},
							},
						},
					},
				}
				require.NoError(t, createKongClusterPlugin(ctx, ctrlClient, ns, plugin))
			},
		},
		{
			name: "KongClusterPlugin - using config and configPatches is allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				plugin := &kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"key_name":"apikey"}`),
					},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/key_names",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: ns,
									Secret:    "secret-name",
									Key:       "key-name",
								},
							},
						},
					},
				}
				require.NoError(t, createKongClusterPlugin(ctx, ctrlClient, ns, plugin))
			},
		},
		{
			name: "KongClusterPlugin - using configFrom and configPatches is not allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				plugin := &kongv1.KongClusterPlugin{
					PluginName: "key-auth",
					ConfigFrom: &kongv1.NamespacedConfigSource{
						SecretValue: kongv1.NamespacedSecretValueFromSource{
							Secret: "secret-name",
							Key:    "key-name",
						},
					},
					ConfigPatches: []kongv1.NamespacedConfigPatch{
						{
							Path: "/key_names",
							ValueFrom: kongv1.NamespacedConfigSource{
								SecretValue: kongv1.NamespacedSecretValueFromSource{
									Namespace: ns,
									Secret:    "secret-name",
									Key:       "key-name",
								},
							},
						},
					},
				}
				assert.ErrorContains(t, createKongClusterPlugin(ctx, ctrlClient, ns, plugin),
					"Using both configFrom and configPatches fields is not allowed.")
			},
		},
		{
			name: "KongVault - changing spec.prefix is not allowed",
			scenario: func(ctx context.Context, t *testing.T, _ string) {
				vault := &kongv1alpha1.KongVault{
					Spec: kongv1alpha1.KongVaultSpec{
						Backend: "env",
						Prefix:  "env-0",
					},
				}
				require.NoError(t, createKongVault(ctx, ctrlClient, vault))
				vault.Spec.Prefix = "env-1"
				assert.ErrorContains(t, updateKongVault(ctx, ctrlClient, vault),
					"The spec.prefix field is immutable")
			},
		},
		{
			name: "KongConsumer - duplicate credentials are not allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				consumer := &kongv1.KongConsumer{
					Username:    "u1",
					Credentials: []string{"c1", "c1", "c2"},
				}
				assert.ErrorContains(t, createKongConsumer(ctx, ctrlClient, ns, consumer), `Duplicate value: "c1"`)
			},
		},
		{
			name: "KongConsumer - unique credentials are allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				consumer := &kongv1.KongConsumer{
					Username:    "u1",
					Credentials: []string{"c1", "c2"},
				}
				assert.NoError(t, createKongConsumer(ctx, ctrlClient, ns, consumer))
			},
		},
		{
			name: "KongConsumer - duplicate consumer groups are not allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				consumer := &kongv1.KongConsumer{
					Username:       "u1",
					ConsumerGroups: []string{"cg1", "cg1", "cg2"},
				}
				assert.ErrorContains(t, createKongConsumer(ctx, ctrlClient, ns, consumer), `Duplicate value: "cg1"`)
			},
		},
		{
			name: "KongConsumer - unique consumer groups are allowed",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				consumer := &kongv1.KongConsumer{
					Username:       "u1",
					ConsumerGroups: []string{"cg1", "cg2"},
				}
				assert.NoError(t, createKongConsumer(ctx, ctrlClient, ns, consumer))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ns := CreateNamespace(ctx, t, ctrlClient)
			tc.scenario(ctx, t, ns.Name)
		})
	}
}

func createFaultyTCPIngress(ctx context.Context, t *testing.T, envcfg *rest.Config, ns string, modifier func(*kongv1beta1.TCPIngress)) error {
	ingress := validTCPIngress()
	modifier(ingress)

	gatewayClient, err := clientset.NewForConfig(envcfg)
	require.NoError(t, err)

	c := gatewayClient.ConfigurationV1beta1().TCPIngresses(ns)
	ingress, err = c.Create(ctx, ingress, metav1.CreateOptions{})
	if !assert.Error(t, err) {
		t.Cleanup(func() { _ = c.Delete(ctx, ingress.Name, metav1.DeleteOptions{}) })
	}
	return err
}

func validTCPIngress() *kongv1beta1.TCPIngress {
	return &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Port: 80,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: "service-name",
						ServicePort: 80,
					},
				},
			},
		},
	}
}

func createFaultyUDPIngress(ctx context.Context, t *testing.T, envcfg *rest.Config, ns string, modifier func(ingress *kongv1beta1.UDPIngress)) error {
	ingress := validUDPIngress()
	modifier(ingress)

	gatewayClient, err := clientset.NewForConfig(envcfg)
	require.NoError(t, err)

	c := gatewayClient.ConfigurationV1beta1().UDPIngresses(ns)
	ingress, err = c.Create(ctx, ingress, metav1.CreateOptions{})
	if !assert.Error(t, err) {
		t.Cleanup(func() { _ = c.Delete(ctx, ingress.Name, metav1.DeleteOptions{}) })
	}
	return err
}

func validUDPIngress() *kongv1beta1.UDPIngress {
	return &kongv1beta1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.UDPIngressSpec{
			Rules: []kongv1beta1.UDPIngressRule{
				{
					Port: 80,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: "service-name",
						ServicePort: 80,
					},
				},
			},
		},
	}
}

func createKongUpstreamPolicy(ctx context.Context, client client.Client, ns string, spec kongv1beta1.KongUpstreamPolicySpec) error {
	return client.Create(ctx, &kongv1beta1.KongUpstreamPolicy{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-",
			Namespace:    ns,
		},
		Spec: spec,
	})
}

// generateInvalidHashOns generates a list of KongUpstreamHash objects with all possible invalid fields pairs.
func generateInvalidHashOns() []kongv1beta1.KongUpstreamHash {
	fieldSetFns := []func(h *kongv1beta1.KongUpstreamHash){
		func(h *kongv1beta1.KongUpstreamHash) {
			h.Input = lo.ToPtr(kongv1beta1.HashInput("consumer"))
		},
		func(h *kongv1beta1.KongUpstreamHash) {
			h.Cookie = lo.ToPtr("cookie-name")
			h.CookiePath = lo.ToPtr("/")
		},
		func(h *kongv1beta1.KongUpstreamHash) {
			h.Header = lo.ToPtr("header-name")
		},
		func(h *kongv1beta1.KongUpstreamHash) {
			h.URICapture = lo.ToPtr("uri-capture")
		},
		func(h *kongv1beta1.KongUpstreamHash) {
			h.QueryArg = lo.ToPtr("query-arg")
		},
	}

	var invalidHashOns []kongv1beta1.KongUpstreamHash
	for outerIdx, fieldSetFn := range fieldSetFns {
		hashOn := kongv1beta1.KongUpstreamHash{}
		fieldSetFn(&hashOn)

		for innerIdx, innerFieldSetFn := range fieldSetFns {
			if outerIdx == innerIdx {
				continue
			}
			invalidHashOn := hashOn.DeepCopy()
			innerFieldSetFn(invalidHashOn)
			invalidHashOns = append(invalidHashOns, *invalidHashOn)
		}
	}

	optStr := func(s *string) string {
		if s == nil {
			return "<nil>"
		}
		return *s
	}
	return lo.UniqBy(invalidHashOns, func(h kongv1beta1.KongUpstreamHash) string {
		return fmt.Sprintf("%s.%s.%s.%s", optStr(h.Cookie), optStr(h.Header), optStr(h.URICapture), optStr(h.QueryArg))
	})
}

func createKongIngress(ctx context.Context, client client.Client, ns string, ingress *kongv1.KongIngress) error {
	ingress.GenerateName = "test-"
	ingress.Namespace = ns
	return client.Create(ctx, ingress)
}

func createKongPlugin(ctx context.Context, client client.Client, ns string, plugin *kongv1.KongPlugin) error {
	plugin.GenerateName = "test-"
	plugin.Namespace = ns
	return client.Create(ctx, plugin)
}

func createKongClusterPlugin(ctx context.Context, client client.Client, ns string, plugin *kongv1.KongClusterPlugin) error {
	plugin.GenerateName = "test-"
	plugin.Namespace = ns
	return client.Create(ctx, plugin)
}

func updateKongPlugin(ctx context.Context, client client.Client, ns string, plugin *kongv1.KongPlugin) error {
	plugin.GenerateName = "test-"
	plugin.Namespace = ns
	return client.Update(ctx, plugin)
}

func updateKongClusterPlugin(ctx context.Context, client client.Client, ns string, plugin *kongv1.KongClusterPlugin) error {
	plugin.GenerateName = "test-"
	plugin.Namespace = ns
	return client.Update(ctx, plugin)
}

func createKongVault(ctx context.Context, c client.Client, vault *kongv1alpha1.KongVault) error {
	vault.GenerateName = "test-"
	return c.Create(ctx, vault)
}

func updateKongVault(ctx context.Context, c client.Client, vault *kongv1alpha1.KongVault) error {
	vault.GenerateName = "test-"
	return c.Update(ctx, vault)
}

func createKongConsumer(ctx context.Context, c client.Client, ns string, consumer *kongv1.KongConsumer) error {
	consumer.GenerateName = "test-"
	consumer.Namespace = ns
	return c.Create(ctx, consumer)
}
