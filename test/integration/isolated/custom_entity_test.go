//go:build integration_tests

package isolated

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-configuration/v2/pkg/clientset"

	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestCustomEntityExample(t *testing.T) {
	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.Kind, testlabels.KindKongCustomEntity).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		Setup(SkipIfEnterpriseNotEnabled).
		Setup(SkipIfDBBacked).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		WithSetup("deploy example manifest", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			manifestPath := examplesManifestPath("kong-custom-entity.yaml")

			b, err := os.ReadFile(manifestPath)
			require.NoError(t, err)
			manifest := string(b)

			ingressClass := GetIngressClassFromCtx(ctx)

			t.Logf("replacing kong ingress class in yaml manifest to %s", ingressClass)
			manifest = strings.ReplaceAll(
				manifest,
				"kubernetes.io/ingress.class: kong",
				fmt.Sprintf("kubernetes.io/ingress.class: %s", ingressClass),
			)
			manifest = strings.ReplaceAll(
				manifest,
				"ingressClassName: kong",
				fmt.Sprintf("ingressClassName: %s", ingressClass),
			)
			manifest = strings.ReplaceAll(
				manifest,
				"controllerName: kong",
				fmt.Sprintf("controllerName: %s", ingressClass),
			)

			t.Logf("applying yaml manifest %s", manifestPath)
			cluster := GetClusterFromCtx(ctx)
			require.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, manifest))

			t.Cleanup(func() {
				t.Logf("deleting yaml manifest %s", manifestPath)
				assert.NoError(t, clusters.DeleteManifestByYAML(ctx, cluster, manifest))
			})

			t.Log("waiting for hasura deployment to be ready")
			helpers.WaitForDeploymentRollout(ctx, t, cluster, "default", "hasura")
			return ctx
		}).
		Assess("degraphql plugin works as expected", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			proxyURL := GetHTTPURLFromCtx(ctx)
			t.Log("Waiting for graphQL service to be available")
			helpers.EventuallyGETPath(t, proxyURL, proxyURL.Host, "/healthz", nil, http.StatusOK, "OK", nil, consts.IngressWait, consts.WaitTick)

			t.Log("injecting data for graphQL service")
			injectDataURL := proxyURL.String() + "/v2/query"
			runSQLCreateTableBody := `{
				"type": "run_sql",
				"args": {
					"sql": "CREATE TABLE contacts(id serial NOT NULL, name text NOT NULL, phone_number text NOT NULL, PRIMARY KEY(id));"
				  }
			}`
			req, err := http.NewRequest(http.MethodPost, injectDataURL, bytes.NewReader([]byte(runSQLCreateTableBody)))
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
			req, err = http.NewRequest(http.MethodPost, injectDataURL, bytes.NewReader([]byte(runSQLInsertRowBody)))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("X-Hasura-Role", "admin")
			resp, err = helpers.DefaultHTTPClient().Do(req)
			require.NoError(t, err)
			resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)

			setMetadataURL := proxyURL.String() + "/v1/metadata"
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
			// The ingress providing graphQL service has a different host, so we need to set the `Host` header.
			helpers.EventuallyGETPath(t, proxyURL, "graphql.service.example", "/contacts", nil, http.StatusOK, `"name":"Alice"`, map[string]string{"Host": "graphql.service.example"}, consts.IngressWait, consts.WaitTick)

			return ctx
		}).
		Assess("another ingress using the same degraphql plugin should also work", func(ctx context.Context, t *testing.T, conf *envconf.Config) context.Context {
			const (
				ingressNamespace = "default"
				serviceName      = "hasura"
				ingressName      = "hasura-ingress-graphql"
				alterServiceName = "hasura-alter"
				alterIngressName = "hasura-ingress-graphql-alter"
			)
			r := conf.Client().Resources()

			t.Log("creating alternative service")
			svc := corev1.Service{}
			require.NoError(t, r.Get(ctx, serviceName, ingressNamespace, &svc))
			alterService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        alterServiceName,
					Namespace:   ingressNamespace,
					Labels:      svc.Labels,
					Annotations: svc.Annotations,
				},
			}
			alterService.Spec = *svc.Spec.DeepCopy()
			alterService.Spec.ClusterIP = ""
			alterService.Spec.ClusterIPs = []string{}
			require.NoError(t, r.Create(ctx, alterService))

			t.Log("creating alternative ingress with the same degraphql plugin attached")
			ingress := netv1.Ingress{}
			require.NoError(t, r.Get(ctx, ingressName, ingressNamespace, &ingress))
			alterIngress := &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:        alterIngressName,
					Namespace:   ingressNamespace,
					Labels:      ingress.Labels,
					Annotations: ingress.Annotations,
				},
			}
			alterIngress.Spec = *ingress.Spec.DeepCopy()
			for i := range alterIngress.Spec.Rules {
				alterIngress.Spec.Rules[i].Host = "alter-graphql.service.example"
				for j := range alterIngress.Spec.Rules[i].HTTP.Paths {
					alterIngress.Spec.Rules[i].HTTP.Paths[j].Backend = netv1.IngressBackend{
						Service: &netv1.IngressServiceBackend{
							Name: alterServiceName,
							Port: netv1.ServiceBackendPort{
								Number: int32(80),
							},
						},
					}
				}
			}
			require.NoError(t, r.Create(ctx, alterIngress))

			t.Log("verifying degraphQL plugin and degraphql_routes entity works")
			proxyURL := GetHTTPURLFromCtx(ctx)
			helpers.EventuallyGETPath(t, proxyURL, "alter-graphql.service.example", "/contacts", nil, http.StatusOK, `"name":"Alice"`, map[string]string{"Host": "graphql.service.example"}, consts.IngressWait, consts.WaitTick)

			return ctx
		}).
		Assess("entities not passing schema validation are excluded from config and produce KongConfigurationTranslationFailed events", func(ctx context.Context, t *testing.T, conf *envconf.Config) context.Context {
			c, err := clientset.NewForConfig(conf.Client().RESTConfig())
			require.NoError(t, err)
			kceClient := c.ConfigurationV1alpha1().KongCustomEntities("default")
			customEntity, err := kceClient.Get(ctx, "degraphql-route-example", metav1.GetOptions{})
			require.NoError(t, err)
			t.Logf("update the custom entity with invalid field, failing validation and thus being excluded from config")
			customEntity.Spec.Fields = apiextensionsv1.JSON{
				Raw: []byte(`{"uri":"/contacts","query":"aaaa","foo":"bar"}`),
			}
			_, err = kceClient.Update(ctx, customEntity, metav1.UpdateOptions{})
			require.NoError(t, err)

			cl, err := client.NewWithWatch(conf.Client().RESTConfig(), client.Options{})
			require.NoError(t, err)

			var events corev1.EventList
			fs := fields.SelectorFromSet(fields.Set{
				"involvedObject.kind":      "KongCustomEntity",
				"involvedObject.name":      "degraphql-route-example",
				"involvedObject.namespace": "default",
				"reason":                   "KongConfigurationTranslationFailed",
			})
			w, err := cl.Watch(ctx, &events, &client.ListOptions{
				Namespace:     "default",
				FieldSelector: fs,
			})
			require.NoError(t, err)
			defer w.Stop()

			timer := time.NewTimer(consts.StatusWait)
			defer timer.Stop()
		forLoop:
			for {
				select {
				case e := <-w.ResultChan():
					event := e.Object.(*corev1.Event)
					// Just perform the assertions and break out of the loop. There shouldn't be more than one event at this point.
					assert.Contains(t, event.Message, "entity default/degraphql-route-example failed validation:")
					assert.Contains(t, event.Message, "schema violation")
					assert.Equal(t, "KongConfigurationTranslationFailed", event.Reason)
					break forLoop

				case <-timer.C:
					t.Logf("timeout waiting for KongConfigurationTranslationFailed Event")
					t.Fail()
					return ctx
				case <-ctx.Done():
					t.Logf("timeout waiting for KongConfigurationTranslationFailed Event")
					t.Fail()
					return ctx
				}
			}

			t.Log("verifying degraphQL plugin was excluded from the configuration as it was failing schema validation")
			proxyURL := GetHTTPURLFromCtx(ctx)
			// The ingress providing graphQL service has a different host, so we need to set the `Host` header.
			helpers.EventuallyGETPath(t, proxyURL, "graphql.service.example", "/contacts", nil, http.StatusNotFound, `{"message":"Not Found"}`, map[string]string{"Host": "graphql.service.example"}, consts.IngressWait, consts.WaitTick)

			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
