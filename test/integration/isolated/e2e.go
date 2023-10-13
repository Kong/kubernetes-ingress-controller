//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"strings"
	"testing"

	ktfkongaddon "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func NameFromT(t *testing.T) string {
	t.Helper()

	name := strings.ToLower(t.Name())
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, "/", ".")

	return name
}

// CreateNSForTest creates a random namespace with the runID as a prefix. It is stored in the context
// so that the deleteNSForTest routine can look it up and delete it.
func createNSForTest(ctx context.Context, cfg *envconf.Config, t *testing.T, runID string) (context.Context, error) {
	t.Helper()

	ns := RandomName("ns-"+runID, 10)
	ctx = context.WithValue(ctx, GetNamespaceKey(t), ns)

	t.Logf("Creating NS %v for test %v", ns, t.Name())
	nsObj := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
			Labels: map[string]string{
				"kubernetes-ingress-controller.konghq.com/test-name": NameFromT(t),
				"kubernetes-ingress-controller.konghq.com/run-id":    runID,
				"kubernetes-ingress-controller.konghq.com/test-type": "integration",
			},
		},
	}
	c, err := client.New(cfg.Client().RESTConfig(), client.Options{})
	if err != nil {
		return ctx, err
	}
	return ctx, c.Create(ctx, &nsObj)

	// TODO: cfg.Client().Resources() causes a data race when run in parallel
	//
	// WARNING: DATA RACE
	// Write at 0x00c0010a8f20 by goroutine 168:
	// sigs.k8s.io/e2e-framework/klient/k8s/resources.(*Resources).WithNamespace()
	// 	/home/runner/go/pkg/mod/sigs.k8s.io/e2e-framework@v0.3.0/klient/k8s/resources/resources.go:85 +0x117
	// sigs.k8s.io/e2e-framework/klient.(*client).Resources()
	// 	/home/runner/go/pkg/mod/sigs.k8s.io/e2e-framework@v0.3.0/klient/client.go:80 +0xe9
	// github.com/kong/kubernetes-ingress-controller/v2/test/integration/isolated.createNSForTest()
	// 	/home/runner/work/kubernetes-ingress-controller/kubernetes-ingress-controller/test/integration/isolated/e2e.go:46 +0x437
	// github.com/kong/kubernetes-ingress-controller/v2/test/integration/isolated.TestUDPRoute.featureSetup.func18()
	// 	/home/runner/work/kubernetes-ingress-controller/kubernetes-ingress-controller/test/integration/isolated/suite_test.go:247 +0x15a
	// sigs.k8s.io/e2e-framework/pkg/env.(*testEnv).executeSteps()
	// 	/home/runner/go/pkg/mod/sigs.k8s.io/e2e-framework@v0.3.0/pkg/env/env.go:428 +0x12a
	// sigs.k8s.io/e2e-framework/pkg/env.(*testEnv).processTestFeature.(*testEnv).execFeature.func1()
	// 	/home/runner/go/pkg/mod/sigs.k8s.io/e2e-framework@v0.3.0/pkg/env/env.go:447 +0x1bd
	// testing.tRunner()
	// 	/opt/hostedtoolcache/go/1.21.1/x64/src/testing/testing.go:1595 +0x238
	// testing.(*T).Run.func1()
	// 	/opt/hostedtoolcache/go/1.21.1/x64/src/testing/testing.go:1648 +0x44

	// Previous write at 0x00c0010a8f20 by goroutine 166:
	// sigs.k8s.io/e2e-framework/klient/k8s/resources.(*Resources).WithNamespace()
	//	/home/runner/go/pkg/mod/sigs.k8s.io/e2e-framework@v0.3.0/klient/k8s/resources/resources.go:85 +0x117
	// sigs.k8s.io/e2e-framework/klient.(*client).Resources()
	//	/home/runner/go/pkg/mod/sigs.k8s.io/e2e-framework@v0.3.0/klient/client.go:80 +0xe9
	// github.com/kong/kubernetes-ingress-controller/v2/test/integration/isolated.createNSForTest()
	//	/home/runner/work/kubernetes-ingress-controller/kubernetes-ingress-controller/test/integration/isolated/e2e.go:46 +0x437
	// github.com/kong/kubernetes-ingress-controller/v2/test/integration/isolated.TestUDPRoute.featureSetup.func15()
	//	/home/runner/work/kubernetes-ingress-controller/kubernetes-ingress-controller/test/integration/isolated/suite_test.go:247 +0x15a
	// sigs.k8s.io/e2e-framework/pkg/env.(*testEnv).executeSteps()
	//	/home/runner/go/pkg/mod/sigs.k8s.io/e2e-framework@v0.3.0/pkg/env/env.go:428 +0x12a
	// sigs.k8s.io/e2e-framework/pkg/env.(*testEnv).processTestFeature.(*testEnv).execFeature.func1()
	//	/home/runner/go/pkg/mod/sigs.k8s.io/e2e-framework@v0.3.0/pkg/env/env.go:447 +0x1bd
	// testing.tRunner()
	//	/opt/hostedtoolcache/go/1.21.1/x64/src/testing/testing.go:1595 +0x238
	// testing.(*T).Run.func1()
	//	/opt/hostedtoolcache/go/1.21.1/x64/src/testing/testing.go:1648 +0x44
}

// DeleteNSForTest looks up the namespace corresponding to the given test and deletes it.
func deleteNSForTest(ctx context.Context, cfg *envconf.Config, t *testing.T, _ string) (context.Context, error) {
	t.Helper()

	ns := fmt.Sprint(ctx.Value(GetNamespaceKey(t)))
	t.Logf("Deleting NS %v for test %v", ns, t.Name())

	nsObj := corev1.Namespace{}
	nsObj.Name = ns
	return ctx, cfg.Client().Resources().Delete(ctx, &nsObj)
}

type (
	NamespaceCtxKey string
	KongAddonCtxKey string
)

// GetNamespaceKey returns the context key for a given test.
func GetNamespaceKey(t *testing.T) NamespaceCtxKey {
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

	raw := ctx.Value(GetNamespaceKey(t))
	str, ok := raw.(string)
	if !ok {
		t.Fatalf("required namespace name to be stored in context but found: %s (of type %T)", raw, raw)
	}

	return str
}

func GetKongAddonCtxKey(t *testing.T) KongAddonCtxKey {
	t.Helper()

	// When we pass t.Name() from inside an `assess` step, the name is in the form TestName/Features/Assess
	if strings.Contains(t.Name(), "/") {
		return KongAddonCtxKey(strings.Split(t.Name(), "/")[0])
	}

	// When pass t.Name() from inside a `testenv.BeforeEachTest` function, the name is just TestName
	return KongAddonCtxKey(t.Name())
}

func SetKongAddonForT(ctx context.Context, t *testing.T, addon *ktfkongaddon.Addon) context.Context {
	t.Helper()

	return context.WithValue(ctx, GetKongAddonCtxKey(t), addon)
}

func GetKongAddonForT(ctx context.Context, t *testing.T) *ktfkongaddon.Addon {
	t.Helper()

	raw := ctx.Value(GetKongAddonCtxKey(t))
	addon, ok := raw.(*ktfkongaddon.Addon)
	if !ok {
		t.Fatalf("required kong addon to be stored in context but found: %s (of type %T)", raw, raw)
	}

	return addon
}
