# Testing guidelines

Following guide will help you decide what to test and how to go about it
when changing code in this repository.

## Testing levels

In KIC, we use several levels of testing:

- [unit tests](#unit-tests)
- [integration tests](#integration-tests)
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
or using e.g. kubernetes related fakes like for instance
[`TestAddressesFromEndpointSlice`][test_AddressesFromEndpointSlice] does.

#### How to run

Unit tests can be either run via appropriate `go test ...` invocation or using any of
the provided Makefile targets:

- `test.unit` run unit tests with standard verbose output on stderr
- `test.unit.pretty` run unit tests with opinionated output format, one line per package

Unit tests are typically located in the same package as the code they test
and they do not use any [build tags][go_build_tag].

> **NOTE**: Some of these tests that do require [`envtest`][envtest] setup and
> use a special [build tags][go_build_tag]: `envtest`. This is in place so that
> simply running `go test ...` without any special build tags would pass without
> requiring any manual setup from the developer.

[test_GetAdminAPIsForService]: https://github.com/Kong/kubernetes-ingress-controller/blob/753e91f73dea5e51a3610d50c8a5928da79baa0f/internal/adminapi/endpoints_envtest_test.go#L28
[test_HTTPRouteReconcilerProperlyReactsToReferenceGrant]: https://github.com/Kong/kubernetes-ingress-controller/blob/ccafa7ca9da7fb52ba959c2ebbc0974e22497b5b/internal/controllers/gateway/httproute_controller_envtest_test.go#L37
[test_AddressesFromEndpointSlice]: https://github.com/Kong/kubernetes-ingress-controller/blob/753e91f73dea5e51a3610d50c8a5928da79baa0f/internal/adminapi/endpoints_test.go#L23-L24
[issue4099]: https://github.com/Kong/kubernetes-ingress-controller/issues/4099
[envtest]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest

### Integration tests

KIC's integration tests rely on its controller manager being run in the same
process as the tests via [`manager.Run()`][manager_run].

These tests rely on a [test suite setup][integration_test_suite]
using [`TestMain()`][pkggodev_testmain].

Said setup will either create a new cluster or use an existing one to run the
tests against.

Most of the setup - and cleanup after the tests have run - is being done using
[ktf][ktf].

Currently, the cluster is being shared across all the tests that are a part
of the test suite so special care needs to be taken in order to clean up after the tests
that run.

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

[ktf]: https://github.com/Kong/kubernetes-testing-framework
[pkggodev_testmain]: https://pkg.go.dev/testing#hdr-Main
[integration_test_suite]: https://github.com/Kong/kubernetes-ingress-controller/blob/61e06ee64ff913aa9952816121125fca7ed59ba5/test/integration/suite_test.go#L36
[manager_run]: https://github.com/Kong/kubernetes-ingress-controller/blob/5abc699aeee552945a76c82e3f7abb3e1b2fabf1/internal/cmd/rootcmd/run.go#L14-L22

### End-to-End (E2E) tests

End to end tests in KIC are used to test features or use cases which rely on KIC
being deployed in cluster.

For instance:

- deploying all in one manifests tested in [all in one tests][test_e2e_all_in_one]
- upgrade scenarios tested in [upgrade tests][test_e2e_upgrade]

These tests deploy KIC and Kong Gateway in cluster using the requested image(s)
which could be customized via dedicated environment variables like

- [`TEST_KONG_CONTROLLER_IMAGE_OVERRIDE`][env_var_controller_image_override]
- [`TEST_KONG_IMAGE_OVERRIDE`][env_var_kong_image_override]

#### How to run

E2E tests are located under `tests/e2e/` and use `e2e_tests` [build tag][go_build_tag].

You can run them using one of the dedicated Makefile targets:

- `test.e2e` run all e2e tests.

[test_e2e_all_in_one]: https://github.com/Kong/kubernetes-ingress-controller/blob/3d45c822bdb907caba568f86062af83406785fc5/test/e2e/all_in_one_test.go
[test_e2e_upgrade]: https://github.com/Kong/kubernetes-ingress-controller/blob/43e797f7394c5f0a9394c6f158f5efff5e2321ec/test/e2e/upgrade_test.go
[env_var_controller_image_override]: https://github.com/Kong/kubernetes-ingress-controller/blob/50a4c3f1e57c56950808b90bcb0b57fefc2f3d7c/test/e2e/environment.go#L12-L13
[env_var_kong_image_override]: https://github.com/Kong/kubernetes-ingress-controller/blob/50a4c3f1e57c56950808b90bcb0b57fefc2f3d7c/test/e2e/environment.go#L19-L20
[go_build_tag]: https://pkg.go.dev/go/build#hdr-Build_Constraints
