package seeder

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	knative "knative.dev/networking/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/proxy"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
)

// -----------------------------------------------------------------------------
// Seeder - Builder
// -----------------------------------------------------------------------------

// Builder is a tool to configure and build Seeder objects.
type Builder struct {
	fieldLogger      logrus.FieldLogger
	restCFG          *rest.Config
	prx              proxy.Proxy
	ingressClassName string
}

// NewBuilder produces a new *Builder object to build Seeders.
func NewBuilder(restCFG *rest.Config, prx proxy.Proxy) *Builder {
	return &Builder{
		restCFG:          restCFG,
		prx:              prx,
		fieldLogger:      logrus.New(),
		ingressClassName: annotations.DefaultIngressClass,
	}
}

// WithFieldLogger allows the caller to provide a custom logger for the Seeder.
func (b *Builder) WithFieldLogger(fieldLogger logrus.FieldLogger) *Builder {
	b.fieldLogger = fieldLogger
	return b
}

// WithIngressClass overrides the default ingress class which is used to identify supported objects.
func (b *Builder) WithIngressClass(ingressClassName string) *Builder {
	b.ingressClassName = ingressClassName
	return b
}

// Build generates the Seeder object based on the current configuration.
func (b *Builder) Build() (*Seeder, error) {
	kc, err := kubernetes.NewForConfig(b.restCFG)
	if err != nil {
		return nil, err
	}

	kongc, err := clientset.NewForConfig(b.restCFG)
	if err != nil {
		return nil, err
	}

	knativec, err := knative.NewForConfig(b.restCFG)
	if err != nil {
		return nil, err
	}

	return &Seeder{
		ingressClassName: b.ingressClassName,

		logger: b.fieldLogger,
		prx:    b.prx,

		kc:       kc,
		kongc:    kongc,
		knativec: knativec,
	}, nil
}
