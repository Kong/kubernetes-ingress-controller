module github.com/kong/kubernetes-ingress-controller

go 1.16

require (
	github.com/blang/semver/v4 v4.0.0
	github.com/bombsimon/logrusr v1.1.0
	github.com/docker/docker v20.10.5+incompatible // indirect
	github.com/eapache/channels v1.1.0
	github.com/fatih/color v1.12.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.4.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/kong/deck v1.7.0
	github.com/kong/go-kong v0.19.0
	github.com/kong/kubernetes-testing-framework v0.0.11
	github.com/lithammer/dedent v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.0
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.8.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/pool.v3 v3.1.1
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/klog v1.0.0
	knative.dev/networking v0.0.0-20210216014426-94bfc013982b
	knative.dev/pkg v0.0.0-20210216013737-584933f8280b
	sigs.k8s.io/controller-runtime v0.9.0
	sigs.k8s.io/yaml v1.2.0
)
