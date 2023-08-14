package consts

// KongHelmChartVersion is the version of the Kong Helm chart to use in tests.
// TODO: use chart 2.25.0 as a workaround before charts deal with invalid semver in image tag properly
// https://github.com/Kong/charts/issues/856
const KongHelmChartVersion = "2.25.0"
