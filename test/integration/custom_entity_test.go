//go:build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/go-kong/kong/custom"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestCustomEntity(t *testing.T) {
	RunWhenKongDBMode(t, dpconf.DBModeOff, "Custom entities are only enabled with dbless mode")
	RunWhenKongEnterprise(t)
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)
	// "sessions" is used by "session" plugin to store existing sessions.
	// The sessions are usually created after successful authentication by auth plugins (like `key-auth`) but not directly created by users.
	// Here we use the entity to test the reconciling and translating functions.
	customEntitySession := &kongv1alpha1.KongCustomEntity{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "session-1",
		},
		Spec: kongv1alpha1.KongCustomEntitySpec{
			ControllerName: consts.IngressClass,
			EntityType:     "sessions",
			Fields: apiextensionsv1.JSON{
				Raw: []byte(`{"session_id":"test-session","data":"foobar"}`),
			},
		},
	}

	c, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	customEntityClient := c.ConfigurationV1alpha1().KongCustomEntities(ns.Name)
	t.Logf("creating a KongCustomEntity in namespace %s to test reconciliation", ns.Name)
	_, err = customEntityClient.Create(ctx, customEntitySession, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(customEntitySession)

	t.Log("waiting for the KongCustomEntity's Programmed condition to be set to True")
	require.Eventually(t, func() bool {
		e, err := customEntityClient.Get(ctx, customEntitySession.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return lo.ContainsBy(e.Status.Conditions, func(condition metav1.Condition) bool {
			return condition.Type == "Programmed" && condition.Status == metav1.ConditionTrue
		})
	}, ingressWait, waitTick, "Programmed condition should become True")

	t.Logf("verifying that the entity is translated and applied on Kong gateway")
	kongClient, err := adminapi.NewKongAPIClient(proxyAdminURL.String(), nil)
	require.NoError(t, err)
	err = kongClient.Registry.Register(custom.Type("sessions"), &custom.EntityCRUDDefinition{
		Name:     "sessions",
		CRUDPath: "/sessions",
		// Although the actual primary key is "id", we cannot know the "id" of the "session" entity.
		// So we use the alternative identifier "session_id" as the primary key here.
		PrimaryKey: "session_id",
	})
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		obj := custom.NewEntityObject("sessions")
		obj.SetObject(custom.Object{
			"session_id": "test-session",
		})

		entity, err := kongClient.CustomEntities.Get(ctx, obj)
		require.NoError(t, err)
		gotObj := entity.Object()
		sessionID, ok := gotObj["session_id"].(string)
		if !ok || sessionID != "test-session" {
			return false
		}
		sessionData, ok := gotObj["data"]
		if !ok || sessionData != "foobar" {
			return false
		}
		return true
	}, ingressWait, waitTick)

	// Test degraphql plugin and degraphql_route custom entity with a graphQL service.
	t.Log("deploying a container providing graphQL services")
	hasuraContainer := generators.NewContainer("hasura", test.HasuraGraphQLEngineImage, test.HasuraGraphQLEnginePort)
	hasuraContainer.Env = []corev1.EnvVar{
		{
			Name:  "HASURA_GRAPHQL_DATABASE_URL",
			Value: "postgres://user:password@localhost:5432/hasura_data",
		},
		{
			Name:  "HASURA_GRAPHQL_ENABLE_CONSOLE",
			Value: "true",
		},
		{
			Name:  "HASURA_GRAPHQL_DEV_MODE",
			Value: "true",
		},
	}
	deployment := generators.NewDeploymentForContainer(hasuraContainer)
	postgresContainer := generators.NewContainer("postgres", test.PostgresImage, test.PostgresPort)
	postgresContainer.Env = []corev1.EnvVar{
		{
			Name:  "POSTGRES_USER",
			Value: "user",
		},
		{
			Name:  "POSTGRES_PASSWORD",
			Value: "password",
		},
		{
			Name:  "POSTGRES_DB",
			Value: "hasura_data",
		},
	}
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, postgresContainer)
	_, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating a degraphql plugin")
	degraphqlPlugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "degraphql-test",
		},
		PluginName: "degraphql",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"graphql_server_path":"/v1/graphql"}`),
		},
	}
	_, err = c.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, degraphqlPlugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(degraphqlPlugin)

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, consts.IngressClass)
	ingress := generators.NewIngressForService("/", map[string]string{
		"konghq.com/strip-path": "true",
		"konghq.com/plugins":    degraphqlPlugin.Name,
	}, service)
	ingress.Spec.IngressClassName = kong.String(consts.IngressClass)
	require.NoError(t, clusters.DeployIngress(ctx, env.Cluster(), ns.Name, ingress))
	cleaner.Add(ingress)

	t.Logf("creating a degraphql_route custom entity")
	customEntityDegraqhQLRoute := &kongv1alpha1.KongCustomEntity{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "degraphql-route-test",
		},
		Spec: kongv1alpha1.KongCustomEntitySpec{
			EntityType:     "degraphql_routes",
			ControllerName: consts.IngressClass,
			Fields: apiextensionsv1.JSON{
				Raw: []byte(`{"uri":"/contacts","query":"query { contacts { name } }"}`),
			},
			ParentRef: &kongv1alpha1.ObjectReference{
				Group: &kongv1.GroupVersion.Group,
				Kind:  lo.ToPtr("KongPlugin"),
				Name:  degraphqlPlugin.Name,
			},
		},
	}
	_, err = customEntityClient.Create(ctx, customEntityDegraqhQLRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(customEntityDegraqhQLRoute)

	t.Logf("verifying that the degraphql_routes entity is translated and applied on Kong gateway")
	require.Eventually(t, func() bool {
		kongServices, err := kongClient.Services.ListAll(ctx)
		require.NoError(t, err)
		if len(kongServices) != 1 {
			return false
		}
		svc := kongServices[0]

		degrapqhQLRoutes, err := kongClient.DegraphqlRoutes.ListAll(ctx, svc.ID)
		require.NoError(t, err)
		if len(degrapqhQLRoutes) != 1 {
			return false
		}
		return degrapqhQLRoutes[0].URI != nil && *degrapqhQLRoutes[0].URI == "/contacts"
	}, ingressWait, waitTick)

	t.Log("waiting for the graqhQL service to get a LoadBalacer IP")
	graqhQLServiceIP := ""
	require.Eventually(t, func() bool {
		service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Get(ctx, service.Name, metav1.GetOptions{})
		require.NoError(t, err)
		if len(service.Status.LoadBalancer.Ingress) == 0 {
			return false
		}
		lbIngressStatus := service.Status.LoadBalancer.Ingress[0]
		if lbIngressStatus.IP != "" {
			graqhQLServiceIP = lbIngressStatus.IP
			return true
		}
		return false
	}, ingressWait, waitTick)
	graphQLServiceURLBase := fmt.Sprintf("http://%s:%d", graqhQLServiceIP, test.HasuraGraphQLEnginePort)

	t.Log("waiting for graqhQL service ready")
	helpers.EventuallyGETPath(
		t, helpers.MustParseURL(t, graphQLServiceURLBase),
		graqhQLServiceIP, "/healthz", http.StatusOK, "OK", nil, ingressWait, waitTick,
	)

	t.Log("configuring data for graqhQL service")
	queryURL := graphQLServiceURLBase + "/v2/query"
	runSQLCreateTableBody := `{
		"type": "run_sql",
		"args": {
			"sql": "CREATE TABLE contacts(id serial NOT NULL, name text NOT NULL, phone_number text NOT NULL, PRIMARY KEY(id));"
		  }
	}`
	req, err := http.NewRequest(http.MethodPost, queryURL, bytes.NewReader([]byte(runSQLCreateTableBody)))
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Hasura-Role", "admin")
	resp, err := helpers.DefaultHTTPClient().Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	runSQLInsertRowBody := `{
		"type": "run_sql",
		"args": {
			"sql": "INSERT INTO contacts (name, phone_number) VALUES ('Alice','0123456789');"
		  }
	}`
	req, err = http.NewRequest(http.MethodPost, queryURL, bytes.NewReader([]byte(runSQLInsertRowBody)))
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Hasura-Role", "admin")
	resp, err = helpers.DefaultHTTPClient().Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	setMetadataURL := graphQLServiceURLBase + "/v1/metadata"
	trackTableBody := `{
		"type": "pg_track_table",
		"args": {
		  "schema": "public",
		  "name": "contacts"
		}
	}`
	req, err = http.NewRequest(http.MethodPost, setMetadataURL, bytes.NewReader([]byte(trackTableBody)))
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Hasura-Role", "admin")
	resp, err = helpers.DefaultHTTPClient().Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	t.Log("verifying degraphQL plugin and degraphql_routes entity works")
	helpers.EventuallyGETPath(t, proxyHTTPURL, proxyHTTPURL.Host, "/contacts", http.StatusOK, `"name":"Alice"`, nil, ingressWait, waitTick)
}
