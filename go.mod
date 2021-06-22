module github.com/kong/kubernetes-ingress-controller

go 1.16

require (
	github.com/blang/semver/v4 v4.0.0
	github.com/bombsimon/logrusr v1.1.0
	github.com/docker/docker v20.10.5+incompatible // indirect
	github.com/eapache/channels v1.1.0
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/fatih/color v1.12.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kong/deck v1.7.0
	github.com/kong/go-kong v0.19.0
	github.com/kong/kubernetes-testing-framework v0.0.12
	github.com/lithammer/dedent v1.1.0
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.22.0
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.8.0
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210603081109-ebe580a85c40 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/pool.v3 v3.1.1
	gopkg.in/ini.v1 v1.62.0 // indirect
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20210527160623-6fdb442a123b // indirect
	knative.dev/networking v0.0.0-20210512050647-ace2d3306f0b
	knative.dev/pkg v0.0.0-20210510175900-4564797bf3b7
	knative.dev/serving v0.23.1
	sigs.k8s.io/controller-runtime v0.9.0-beta.6
	sigs.k8s.io/yaml v1.2.0
)
