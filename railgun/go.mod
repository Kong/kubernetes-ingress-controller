module github.com/kong/railgun

go 1.15

require (
	github.com/go-logr/logr v0.3.0
	github.com/google/uuid v1.2.0
	github.com/kong/go-kong v0.15.0
	github.com/kong/kubernetes-ingress-controller v1.1.2-0.20210205121726-5872d8431e8c
	github.com/kong/kubernetes-testing-framework v0.0.0-20210209161849-62d8d8602a0e
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/networking v0.0.0-20201028144035-3287613a3b41
	sigs.k8s.io/controller-runtime v0.7.0
)

replace k8s.io/client-go => k8s.io/client-go v0.20.2
