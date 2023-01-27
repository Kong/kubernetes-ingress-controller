package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestKonnectRuntimeGroupConfigPush(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("building test cluster and environment")
	builder, err := getEnvironmentBuilder(ctx)
	require.NoError(t, err)
	env, err := builder.Build(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		helpers.TeardownCluster(ctx, t, env.Cluster())
	})

	// TODO: create runtime group in konnect

	// TODO: generate client certificate and POST it as dp cert to runtime group

	// TODO: create

	// TODO: deploy config/konnect with

	t.Log("deploying kong components")
	const konnectDeploymentPath = "../../deploy/single/all-in-one-dbless.yaml"
	manifest, err := getTestManifest(t, konnectDeploymentPath)
	require.NoError(t, err)
	deployment := deployKong(ctx, t, env, manifest)
	require.NotNil(t, deployment)
}
