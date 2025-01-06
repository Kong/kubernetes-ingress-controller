# Testing guidelines

Following guide will help you decide what to test and how to go about it
when changing code in this repository.

## Testing levels

In KIC, we use several levels of testing:

- [unit tests](#unit-tests)
- [`envtest` based tests](#envtest-based-tests)
- [integration tests](#integration-tests)
- [isolated integration tests](#isolated-integration-tests)
- [Kong integration tests](#kong-integration-tests)
- [end to end (E2E) tests](#end-to-end-e2e-tests)

### Unit tests

Unit tests verify the functionality of one or more `func` or `struct`s
in an isolated environment.

These tests are independent of the environment they are run in and require no
setup or teardown code.

> **NOTE**: Some of these tests do require [`envtest`][envtest] setup to be present e.g.
> [this one][test_GetAdminAPIsForService] or [this one][test_HTTPRouteReconcilerProperlyReactsToReferenceGrant].
> This setup is handled automatically via `setup-envtest` Makefile target which is a dependency
> of mentioned above unit test targets.
>
> As part of [#4099][issue4099] this organization might change and tests that do
> require this setup might end up in a separate place.

These tests can either be written in a way that do not require any mocks or fakes,
or using e.g. Kubernetes related fakes like for instance
[`TestAddressesFromEndpointSlice`][test_AddressesFromEndpointSlice] does.

#### How to run

Unit tests can be either run via appropriate `go test ...` invocation or using any of
the provided Makefile targets:

- `test.unit` run unit tests with standard verbose output on stderr
- `test.unit.pretty` run unit tests with opinionated output format, one line per package

Unit tests are typically located in the same package as the code they test
and they do not use any [build tags][go_build_tag].

[test_GetAdminAPIsForService]: https://github.com/Kong/kubernetes-ingress-controller/blob/753e91f73dea5e51a3610d50c8a5928da79baa0f/internal/adminapi/endpoints_envtest_test.go#L28
[test_HTTPRouteReconcilerProperlyReactsToReferenceGrant]: https://github.com/Kong/kubernetes-ingress-controller/blob/ccafa7ca9da7fb52ba959c2ebbc0974e22497b5b/internal/controllers/gateway/httproute_controller_envtest_test.go#L37
[test_AddressesFromEndpointSlice]: https://github.com/Kong/kubernetes-ingress-controller/blob/753e91f73dea5e51a3610d50c8a5928da79baa0f/internal/adminapi/endpoints_test.go#L23-L24
[issue4099]: https://github.com/Kong/kubernetes-ingress-controller/issues/4099

### `envtest` based tests

[`envtest`][envtest] based tests use a middle ground approach, somewhere between
unit and integration tests.
They run tests within an isolated, limited environment which allows some basic
Kubernetes cluster actions like CRUD operations on Kubernetes entities and running
controllers against this environment.

Ensuring that the environment is ready for tests is done by:

- [`setup-envtest`][setup_envtest] which ensures that api-server and etcd binaries
  are present on the local system
- `Makefile` invocation of `setup-envtest use` which does what's described above

There already exist several helpers for working with `envtest` environment which
can be found at [`test/envtest/`][envtest_helpers].

> **NOTE**: Currently, these tests that do require [`envtest`][envtest] setup and
> use a special [build tags][go_build_tag]: `envtest`. This is in place so that
> simply running `go test ...` without any special build tags would pass without
> requiring any manual setup from the developer (no failure due to missing `envtest`
> related binaries or environment variables).

Test definitions can be found both in dedicated directory - `test/envtest` - and
throughout the codebase next to code that they are testing.
This approach however may change over time.
Suggestion for new `envtest` based tests is to be placed in `test/envtest` and to
test code as a black box, i.e. using only exported methods (which will be enforced
via running them from this dedicated package).

[envtest]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest
[setup_envtest]: https://github.com/kubernetes-sigs/controller-runtime/tree/main/tools/setup-envtest
[envtest_helpers]: https://github.com/Kong/kubernetes-ingress-controller/tree/611f3c6334424a700f9a00f2801c3cfa2b352d81/test/envtest

#### How to run

[`envtest`][envtest] based tests can be run via `test.envtest` `Makefile` target.

One can also use `./bin/setup-envtest use -p env` to obtain the asset dir environment
variable:

```
./bin/setup-envtest use -p env
export KUBEBUILDER_ASSETS='/Users/username/Library/Application Support/io.kubebuilder.envtest/k8s/1.27.1-darwin-arm64'
```

and use that to manually run the tests:

```
$ eval $(./bin/setup-envtest use -p env)
$ go test -v -count 1 -tags envtest ./pkg/to/test
...
```

### Integration tests

KIC's integration tests rely on its controller manager being run in the same
process as the tests via [`manager.Run()`][manager_run].

These tests rely on a [test suite setup][integration_test_suite]
using [`TestMain()`][pkggodev_testmain].

Said setup will either create a new [kind](https://kind.sigs.k8s.io) cluster or use an existing one to run the
tests against.

Most of the setup - and cleanup after the tests have run - is being done using
[ktf][ktf].

Currently, the cluster is being shared across all the tests that are a part
of the test suite so special care needs to be taken in order to clean up after the tests
that run.
The typical approach is to run them in a dedicated, disposable Kubernetes namespace created just for the purposes of that test.

#### How to run

Integration tests are located under `tests/integration/` and use `integration_tests`
[build tag][go_build_tag].

You can run them using one of the dedicated Makefile targets:

- `test.integration` run all integration tests with standard verbose output on stderr.
  This will run tests for dbless, postgres, enterprise and non enterprise.

- `test.integration.dbless` run all dbless integration tests with standard
  verbose output on stderr. The output will also include controllers' logs.

- `test.integration.postgres` run all postgres integration tests with standard
  verbose output on stderr. The output will also include controllers' logs.

Through `GOTESTFLAGS` you can specify custom flags that will be passed to `go test`.
This can allow you to run a subset of all the tests for faster feedback times, e.g.:

```
make test.integration.dbless GOTESTFLAGS="-count 1 -run TestUDPRouteEssentials"
```

[ktf]: https://github.com/Kong/kubernetes-testing-framework
[pkggodev_testmain]: https://pkg.go.dev/testing#hdr-Main
[integration_test_suite]: https://github.com/Kong/kubernetes-ingress-controller/blob/61e06ee64ff913aa9952816121125fca7ed59ba5/test/integration/suite_test.go#L36
[manager_run]: https://github.com/Kong/kubernetes-ingress-controller/blob/5abc699aeee552945a76c82e3f7abb3e1b2fabf1/internal/cmd/rootcmd/run.go#L14-L22

### Isolated integration tests

Similarly to KIC's integration tests, isolated integration tests rely on its
controller manager being run in the same process as the tests via [`manager.Run()`][manager_run].

These tests rely on [kubernetes-sigs/e2e-framework][github-e2e-framework] for setup,
teardown, tests filtering etc.

Said setup will either create a new [kind](https://kind.sigs.k8s.io) cluster or use an existing one to run the
tests against.

Most of the setup - and cleanup after the tests have run - is being done using [ktf][ktf].

Currently, the cluster is being shared across all the tests that are a part
of the test suite so special care needs to be taken in order to clean up after the tests
that run.
The typical approach is to run them in a dedicated, disposable Kubernetes namespace created just for the purposes of that test.

#### Difference between isolated and regular integration tests

There are a couple of key differences between isolated and regular integration tests:

- In isolated tests each test's [`Feature`][pkggodev-e2e-framework-feature] gets
  its own Kong deployment through the means of [ktf][ktf]'s Kong addon.
- In isolated tests each test's [`Feature`][pkggodev-e2e-framework-feature] gets
  its own controller manager instance started which is configured against the above
  mentioned kong instance.
  Said controller manager instance will be configured to only watch `Feature`'s
  dedicated namespace (more info below) via `--watch-namespace`.
  You can add more namespaces for the controller manager to watch via
  `ControllerManagerOptAdditionalWatchNamespace` `featureSetup` option:

  ```go
  ...
  WithSetup("deploy kong addon into cluster", featureSetup(
    helpers.ControllerManagerOptAdditionalWatchNamespace("my-additional-namespace"),
  )).
  ...
  ```

- In regular integration tests each test gets its own Kubernetes namespace through
  manual creation via helper functions like:

  ```go
  ns, cleaner := helpers.Setup(ctx, t, env)
  ```

  Which will automatically remove the namspace when the test is finished.

  In isolated tests each test's [`Feature`][pkggodev-e2e-framework-feature] will
  get its own Kubernetes namespace which you can get via:

  ```go
  namespace := GetNamespaceForT(ctx, t)
  ```

#### Room for improvement

- Eventually the whole integration suite could be migrated to use this setup.
  When that happens, we can add logic into setup which would make tests that don't 
  need this level of separation to reuse a common installation of Kong (e.g. to
  its default - `kong` - namespace).

  This way we'll have the best of both worlds:

  - seprate tests where it's needed
  - shared, when it's not and where speed is the priority

#### How to run

Tests are located under `tests/integration/isolated` and use `integration_tests`
[build tag][go_build_tag].

You can run them using one of the dedicated Makefile targets:

- `test.integration.isolated.dbless` run all dbless isolated integration tests with standard
  verbose output on stderr. The output will also include controllers' logs.

Through `GOTESTFLAGS` you can specify custom flags that will be passed to `go test`.
This can allow you to run a subset of all the tests for faster feedback times, e.g.:

```
make test.integration.isolated.dbless GOTESTFLAGS="-count 1 -run TestUDPRouteEssentials"
```

You can also specify e2e-framework's flags e.g. to filter tests via [labels][github-e2e-framework-labels].

```
make test.integration.isolated.dbless E2E_FRAMEWORK_FLAGS="-labels=kind=UDPRoute,example=true"
```

[github-e2e-framework]: https://github.com/kubernetes-sigs/e2e-framework
[github-e2e-framework-labels]: https://github.com/kubernetes-sigs/e2e-framework/tree/main/examples/skip_flags#use-labels-in-your-tests
[pkggodev-e2e-framework-feature]: https://pkg.go.dev/sigs.k8s.io/e2e-framework/pkg/features
[ktf]: https://github.com/Kong/kubernetes-testing-framework
[pkggodev_testmain]: https://pkg.go.dev/testing#hdr-Main
[integration_test_suite]: https://github.com/Kong/kubernetes-ingress-controller/blob/61e06ee64ff913aa9952816121125fca7ed59ba5/test/integration/suite_test.go#L36
[manager_run]: https://github.com/Kong/kubernetes-ingress-controller/blob/5abc699aeee552945a76c82e3f7abb3e1b2fabf1/internal/cmd/rootcmd/run.go#L14-L22

### Kong integration tests

Tests located under `test/kongintegration/` allow verifying individual components of KIC
against a running Kong Gateway instance.

Every test in this package requires a Kong instance to be running, but do not require
a Kubernetes cluster nor a full KIC deployment. Each test is responsible for setting up
its own Kong instance and cleaning it up after the test is finished. [testcontainers-go]
is used for that purpose.

Examples of tests that could fit into this category are:

- Ensuring that a component responsible for configuring Kong Gateway works as expected
- Ensuring that configuration generated by a translator is accepted by Kong Gateway

#### How to run

The only requirement for running Kong integration tests is to have a docker daemon running.

All tests can be run using the following Makefile targets:

```shell
# Verbose output
make test.kongintegration

# Pretty output
make test.kongintegration.pretty
```

By default, `kong:latest` image is used for running Kong Gateway. This can be changed
by setting the `TEST_KONG_IMAGE` and `TEST_KONG_TAG` environment variables.

#### Potential issues

##### testcontainer's Garbage Collector

In case your environment doesn't allow running containers in privileged mode, you may encounter
issues with testcontainer's [garbage collector][testcontainers-gc] (Ryuk) which requires it. In
Makefile targets we disable the garbage collector by setting `TESTCONTAINERS_RYUK_DISABLED` environment
variable to `true`.

If you encounter issues with this while running tests directly (e.g. when debugging
in an IDE), you can fix it by:

- setting `TESTCONTAINERS_RYUK_DISABLED` environment variable to `true` in your IDE configuration,
- adding `ryuk.disabled=true` to `.testcontainers.properties` file in your home directory (see
  [Configuration locations][testcontainers-config] for exact locations depending on your OS).

[testcontainers-go]: https://golang.testcontainers.org/
[testcontainers-config]: https://golang.testcontainers.org/features/configuration/#configuration-locations
[testcontainers-gc]: https://golang.testcontainers.org/features/garbage_collector/

### End-to-End (E2E) tests

End to end tests in KIC are used to test features or use cases which rely on KIC
being deployed in-cluster.

For instance:

- deploying all in one manifests tested in [all in one tests][test_e2e_all_in_one]
- upgrade scenarios tested in [upgrade tests][test_e2e_upgrade]

These tests deploy KIC and Kong Gateway in a cluster using the requested image(s)
which could be customized via dedicated environment variables like:

- [`TEST_KONG_CONTROLLER_IMAGE_OVERRIDE`][env_var_controller_image_override]
- [`TEST_KONG_IMAGE_OVERRIDE`][env_var_kong_image_override]

On CI, those are being run both in kind and GKE environments.

#### How to run

E2E tests are located under `tests/e2e/` and use `e2e_tests` [build tag][go_build_tag].

You can run them using one of the dedicated Makefile targets:

- `test.e2e` run all e2e tests.

Through `GOTESTFLAGS` you can specify custom flags that will be passed to `go test`.

`E2E_TEST_RUN` is also available to specify the name of the test to run.
This is being used in CI where Github Actions' matrix specifies tests to run
in each workflow.

Exemplar local invocation can look like this:

```
make test.e2e GOTESTFLAGS="-count 1" E2E_TEST_RUN=TestDeployAllInOneDBLESS
```

##### Konnect tests

Some of the E2E tests are Konnect-specific and are being run only when the `TEST_KONG_KONNECT_ACCESS_TOKEN`
environment variable is set. You can create the access token using the [Konnect UI][konnect_tokens].
The tests run against a `.tech` Konnect environment therefore the token you obtain must be created for an account
in that environment.

Exemplary local invocation for running a Konnect test can look like this:

```
make test.e2e GOTESTFLAGS="-count 1" E2E_TEST_RUN=TestKonnectConfigPush TEST_KONG_KONNECT_ACCESS_TOKEN=<token>
```

[test_e2e_all_in_one]: https://github.com/Kong/kubernetes-ingress-controller/blob/3d45c822bdb907caba568f86062af83406785fc5/test/e2e/all_in_one_test.go
[test_e2e_upgrade]: https://github.com/Kong/kubernetes-ingress-controller/blob/43e797f7394c5f0a9394c6f158f5efff5e2321ec/test/e2e/upgrade_test.go
[env_var_controller_image_override]: https://github.com/Kong/kubernetes-ingress-controller/blob/50a4c3f1e57c56950808b90bcb0b57fefc2f3d7c/test/e2e/environment.go#L12-L13
[env_var_kong_image_override]: https://github.com/Kong/kubernetes-ingress-controller/blob/50a4c3f1e57c56950808b90bcb0b57fefc2f3d7c/test/e2e/environment.go#L19-L20
[go_build_tag]: https://pkg.go.dev/go/build#hdr-Build_Constraints
[konnect_tokens]: https://cloud.konghq.tech/global/account/tokens
