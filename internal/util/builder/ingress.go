package builder

import (
	"strings"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
)

type IngressBuilder struct {
	ingress netv1.Ingress
}

// NewIngress builds an Ingress object with the given name and class, when "" is passed as class parameter
// the field .Spec.IngressClassName is not set.
func NewIngress(name string, class string) *IngressBuilder {
	var classToSet *string
	if class != "" {
		classToSet = &class
	}
	return &IngressBuilder{
		ingress: netv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Annotations: make(map[string]string),
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "Ingress",
				APIVersion: "networking.k8s.io/v1",
			},
			Spec: netv1.IngressSpec{
				IngressClassName: classToSet,
			},
		},
	}
}

func (b *IngressBuilder) Build() *netv1.Ingress {
	return &b.ingress
}

func (b *IngressBuilder) WithLegacyClassAnnotation(class string) *IngressBuilder {
	b.ingress.Annotations[annotations.IngressClassKey] = class
	return b
}

func (b *IngressBuilder) WithRules(rules ...netv1.IngressRule) *IngressBuilder {
	b.ingress.Spec.Rules = append(b.ingress.Spec.Rules, rules...)
	return b
}

func (b *IngressBuilder) WithNamespace(namespace string) *IngressBuilder {
	b.ingress.ObjectMeta.Namespace = namespace
	return b
}

func (b *IngressBuilder) WithAnnotations(annotations map[string]string) *IngressBuilder {
	if b.ingress.Annotations == nil {
		b.ingress.Annotations = annotations
		return b
	}
	for k, v := range annotations {
		b.ingress.Annotations[k] = v
	}
	return b
}

func (b *IngressBuilder) WithKongPlugins(names ...string) *IngressBuilder {
	return b.WithAnnotations(map[string]string{
		annotations.AnnotationPrefix + annotations.PluginsKey: strings.Join(names, ","),
	})
}

func (b *IngressBuilder) WithDefaultBackend(backend *netv1.IngressBackend) *IngressBuilder {
	b.ingress.Spec.DefaultBackend = backend
	return b
}
