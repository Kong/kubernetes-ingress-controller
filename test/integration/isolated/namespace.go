//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

// CreateNSForTest creates a random namespace with the runID as a prefix.
// It is stored in the context so that the deleteNSForTest routine can look it up and delete it.
func CreateNSForTest(ctx context.Context, cfg *envconf.Config, t *testing.T, runID string) (context.Context, error) {
	t.Helper()

	// TODO: We could be tempted to use cfg.Client().Resources() here but when
	// running tests in parallel this causes a data race.
	// Related upstream issue: https://github.com/kubernetes-sigs/e2e-framework/issues/352
	c, err := client.New(cfg.Client().RESTConfig(), client.Options{})
	if err != nil {
		return ctx, err
	}

	t.Logf("Creating namespace for test %v", t.Name())
	nsObj := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "ns-" + runID + "-",
			Labels: map[string]string{
				"kubernetes-ingress-controller.konghq.com/test-name": NameFromT(t),
				"kubernetes-ingress-controller.konghq.com/run-id":    runID,
				"kubernetes-ingress-controller.konghq.com/test-type": "integration",
			},
		},
	}

	if err := c.Create(ctx, &nsObj); err != nil {
		return ctx, err
	}

	t.Logf("Created namespace %s for test %v", nsObj.Name, t.Name())
	ctx = context.WithValue(ctx, getNamespaceKey(t), nsObj.Name)

	return ctx, nil
}

// deleteNSForTest looks up the namespace corresponding to the given test and deletes it.
func deleteNSForTest(ctx context.Context, cfg *envconf.Config, t *testing.T, _ string) (context.Context, error) {
	t.Helper()

	ns := fmt.Sprint(ctx.Value(getNamespaceKey(t)))
	t.Logf("Deleting NS %v for test %v", ns, t.Name())

	nsObj := corev1.Namespace{}
	nsObj.Name = ns
	return ctx, cfg.Client().Resources().Delete(ctx, &nsObj)
}

type (
	NamespaceCtxKey string
)

// getNamespaceKey returns the context key for a given test.
func getNamespaceKey(t *testing.T) NamespaceCtxKey {
	t.Helper()

	// When we pass t.Name() from inside an `assess` step, the name is in the form TestName/Features/Assess
	if strings.Contains(t.Name(), "/") {
		return NamespaceCtxKey(strings.Split(t.Name(), "/")[0])
	}

	// When pass t.Name() from inside a `testenv.BeforeEachTest` function, the name is just TestName
	return NamespaceCtxKey(t.Name())
}

// GetNamespaceForT returns the namespace for a given test.
func GetNamespaceForT(ctx context.Context, t *testing.T) string {
	t.Helper()

	raw := ctx.Value(getNamespaceKey(t))
	str, ok := raw.(string)
	if !ok {
		t.Fatalf("required namespace name to be stored in context but found: %s (of type %T)", raw, raw)
	}

	return str
}
