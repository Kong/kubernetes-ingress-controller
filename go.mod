module github.com/kong/kubernetes-ingress-controller

go 1.15

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/docker/docker v20.10.3+incompatible
	github.com/eapache/channels v1.1.0
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.3.0
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-memdb v1.2.0 // indirect
	github.com/hashicorp/go-uuid v1.0.1
	github.com/kong/deck v1.2.1
	github.com/kong/go-kong v0.15.0
	github.com/lithammer/dedent v1.1.0
	github.com/mitchellh/mapstructure v1.3.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.2.1
	github.com/tidwall/match v1.0.1 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/pool.v3 v3.1.1
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/klog v1.0.0
	knative.dev/networking v0.0.0-20201028144035-3287613a3b41
	knative.dev/pkg v0.0.0-20201026165741-2f75016c1368
	sigs.k8s.io/controller-runtime v0.7.0
)

replace k8s.io/client-go => k8s.io/client-go v0.20.2
