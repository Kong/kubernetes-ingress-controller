module github.com/kong/kubernetes-ingress-controller

go 1.16

require (
	cloud.google.com/go v0.90.0
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/blang/semver/v4 v4.0.0
	github.com/bombsimon/logrusr v1.1.0
	github.com/fatih/color v1.12.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/kong/deck v1.7.0
	github.com/kong/go-kong v0.20.0
	github.com/kong/kubernetes-testing-framework v0.6.1
	github.com/lithammer/dedent v1.1.0
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.26.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.8.1
	google.golang.org/api v0.52.0
	google.golang.org/genproto v0.0.0-20210728212813-7823e685a01f
	k8s.io/api v0.22.0
	k8s.io/apiextensions-apiserver v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.22.0
	knative.dev/networking v0.0.0-20210622182128-53f45d6d2cfa
	knative.dev/pkg v0.0.0-20210622173328-dd0db4b05c80
	knative.dev/serving v0.24.0
	sigs.k8s.io/controller-runtime v0.9.5
	sigs.k8s.io/yaml v1.2.0
)
