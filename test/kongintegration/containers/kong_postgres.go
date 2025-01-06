package containers

import (
	"context"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	// postgresImage is the default image used for the Postgres container.
	postgresImage = "postgres:16-alpine"

	// postgresUser is the default user for the Postgres container.
	postgresUser = "postgres"

	// postgresPassword is the default password for the Postgres container.
	postgresPassword = "pass"

	// postgresDatabase is the default database for the Postgres container.
	postgresDatabase = "kong"

	// postgres is the default port for the Postgres container.
	postgresPort = 5432

	// postgresContainerNetworkAlias is the hostname alias for the Postgres container.
	postgresContainerNetworkAlias = "db"
)

type Postgres struct {
	container testcontainers.Container
}

func NewPostgres(ctx context.Context, t *testing.T, net *testcontainers.DockerNetwork) *Postgres {
	postgresC, err := postgres.Run(ctx,
		postgresImage,
		network.WithNetwork([]string{postgresContainerNetworkAlias}, net),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDatabase,
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections"),
		),
	)
	require.NoError(t, err)
	t.Logf("Postgres container ID: %s", postgresC.GetContainerID())

	t.Cleanup(func() { //nolint:contextcheck
		// If the container is already terminated, we don't need to terminate it again.
		if postgresC.IsRunning() {
			assert.NoError(t, postgresC.Terminate(context.Background()))
		}
	})

	runKongDBMigrations(ctx, t, net.Name)

	return &Postgres{
		container: postgresC,
	}
}

// runKongDBMigrations runs the Kong migrations bootstrap command in a container against a Postgres database container.
func runKongDBMigrations(ctx context.Context, t *testing.T, networkName string) {
	// Run Kong migrations bootstrap command in a container.
	kongMigrationsC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: kongImageUnderTest(),
			Env: map[string]string{
				"KONG_DATABASE":    "postgres",
				"KONG_PG_HOST":     postgresContainerNetworkAlias,
				"KONG_PG_PORT":     strconv.Itoa(postgresPort),
				"KONG_PG_USER":     postgresUser,
				"KONG_PG_PASSWORD": postgresPassword,
				"KONG_PG_DATABASE": postgresDatabase,
			},
			Cmd: []string{
				"kong", "migrations", "bootstrap",
				"--yes", "--force",
				"--db-timeout", "30",
			},
			Networks: []string{networkName},
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Logf("Kong migrations container ID: %s", kongMigrationsC.GetContainerID())

	// Wait for migrations to finish successfully (status == "exited" and exit code == 0).
	const (
		timeout = 30 * time.Second
		period  = 1 * time.Second
	)
	timer := time.After(timeout)
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-timer:
			assert.Fail(t, "Kong migrations bootstrap timed out")
			return
		default:
		}

		state, err := kongMigrationsC.State(ctx)
		require.NoError(t, err)
		if state.Status == "exited" {
			if !assert.Equal(t, 0, state.ExitCode, "Kong migrations bootstrap failed") {
				logs, err := kongMigrationsC.Logs(ctx)
				require.NoError(t, err)

				logsB, err := io.ReadAll(logs)
				require.NoError(t, err)

				t.Logf("Kong migrations bootstrap logs: %s", string(logsB))
			}
			return
		}

		t.Logf("Waiting for Kong migrations to finish, current state: %s", state.Status)
	}
}
