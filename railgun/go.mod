module github.com/shaneutt/railgun

go 1.15

require (
	github.com/go-logr/logr v0.3.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/kong/go-kong v0.15.0
	github.com/kong/kubernetes-ingress-controller v1.1.2-0.20210203125350-f9d44c654a00
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	golang.org/x/sys v0.0.0-20200928205150-006507a75852 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	sigs.k8s.io/controller-runtime v0.7.0
)

replace k8s.io/client-go => k8s.io/client-go v0.19.0
