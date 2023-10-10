package admission

import (
	"context"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// KongValidator validates Kong entities.
type KongValidator interface {
	ValidateConsumer(ctx context.Context, consumer kongv1.KongConsumer) (bool, string, error)
	ValidateConsumerGroup(ctx context.Context, consumerGroup kongv1beta1.KongConsumerGroup) (bool, string, error)
	ValidatePlugin(ctx context.Context, plugin kongv1.KongPlugin) (bool, string, error)
	ValidateClusterPlugin(ctx context.Context, plugin kongv1.KongClusterPlugin) (bool, string, error)
	ValidateCredential(ctx context.Context, secret corev1.Secret) (bool, string, error)
	ValidateGateway(ctx context.Context, gateway gatewayapi.Gateway) (bool, string, error)
	ValidateHTTPRoute(ctx context.Context, httproute gatewayapi.HTTPRoute) (bool, string, error)
	ValidateIngress(ctx context.Context, ingress netv1.Ingress) (bool, string, error)
}

// AdminAPIServicesProvider provides KongHTTPValidator with Kong Admin API services that are needed to perform
// validation against entities stored by the Gateway.
type AdminAPIServicesProvider interface {
	GetConsumersService() (kong.AbstractConsumerService, bool)
	GetPluginsService() (kong.AbstractPluginService, bool)
	GetConsumerGroupsService() (kong.AbstractConsumerGroupService, bool)
	GetInfoService() (kong.AbstractInfoService, bool)
	GetRoutesService() (kong.AbstractRouteService, bool)
}

// KongHTTPValidator implements KongValidator interface to validate Kong
// entities using the Admin API of Kong.
type KongHTTPValidator struct {
	Logger                   logrus.FieldLogger
	SecretGetter             kongstate.SecretGetter
	ManagerClient            client.Client
	AdminAPIServicesProvider AdminAPIServicesProvider
	ParserFeatures           parser.FeatureFlags
	KongVersion              semver.Version

	ingressClassMatcher   func(*metav1.ObjectMeta, string, annotations.ClassMatching) bool
	ingressV1ClassMatcher func(*netv1.Ingress, annotations.ClassMatching) bool
}

// NewKongHTTPValidator provides a new KongHTTPValidator object provided a
// controller-runtime client which will be used to retrieve reference objects
// such as consumer credentials secrets. If you do not pass a cached client
// here, the performance of this validator can get very poor at high scales.
func NewKongHTTPValidator(
	logger logrus.FieldLogger,
	managerClient client.Client,
	ingressClass string,
	servicesProvider AdminAPIServicesProvider,
	parserFeatures parser.FeatureFlags,
	kongVersion semver.Version,
) KongHTTPValidator {
	return KongHTTPValidator{
		Logger:                   logger,
		SecretGetter:             &managerClientSecretGetter{managerClient: managerClient},
		ManagerClient:            managerClient,
		AdminAPIServicesProvider: servicesProvider,
		ParserFeatures:           parserFeatures,
		KongVersion:              kongVersion,

		ingressClassMatcher:   annotations.IngressClassValidatorFuncFromObjectMeta(ingressClass),
		ingressV1ClassMatcher: annotations.IngressClassValidatorFuncFromV1Ingress(ingressClass),
	}
}

// -----------------------------------------------------------------------------
// Private - Manager Client Secret Getter
// -----------------------------------------------------------------------------

type managerClientSecretGetter struct {
	managerClient client.Client
}

func (m *managerClientSecretGetter) GetSecret(namespace, name string) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	return secret, m.managerClient.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, secret)
}
