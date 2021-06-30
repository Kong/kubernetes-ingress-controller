package seeder

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	knative "knative.dev/networking/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/proxy"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/config"
)

// -----------------------------------------------------------------------------
// Seeder - Builder
// -----------------------------------------------------------------------------

// Builder is a tool to configure and build Seeder objects.
type Builder struct {
	namespaces       []string
	ingressClassName string
	controllerConfig *config.ControllerConfig

	fieldLogger logrus.FieldLogger

	restCFG *rest.Config
	prx     proxy.Proxy
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

// WithNamespaces configures which Kubernetes namespaces objects should be listed on.
func (b *Builder) WithNamespaces(namespaces ...string) *Builder {
	b.namespaces = append(b.namespaces, namespaces...)
	return b
}

// WithControllerConfig enables filtering out APIs that are disabled
func (b *Builder) WithControllerConfig(controllerConfig *config.ControllerConfig) *Builder {
	b.controllerConfig = controllerConfig
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

	if len(b.namespaces) == 0 {
		b.fieldLogger.Info("no namespace was provided for object seeder: using all namespaces")
		b.namespaces = []string{corev1.NamespaceAll}
	}

	return &Seeder{
		namespaces:       b.namespaces,
		ingressClassName: b.ingressClassName,
		controllerConfig: b.controllerConfig,

		logger: b.fieldLogger,
		prx:    b.prx,

		kc:       kc,
		kongc:    kongc,
		knativec: knativec,
	}, nil
}
