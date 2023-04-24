package parser

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclientfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset/fake"
)

func TestGetPluginsAssociatedWithCACertSecret(t *testing.T) {
	kongPluginWithSecret := func(name, secretID string) kongv1.KongPlugin {
		return kongv1.KongPlugin{
			TypeMeta: metav1.TypeMeta{
				Kind:       "KongPlugin",
				APIVersion: "configuration.konghq.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				// https://github.com/kubernetes-sigs/controller-runtime/blob/22718275bffe3185276dc835d610c658f06dac07/pkg/client/fake/client.go#L247-L250
				ResourceVersion: "999",
				Name:            name,
			},
			Config: v1.JSON{
				Raw: []byte(fmt.Sprintf(`{"ca_certificates":["%s"]}`, secretID)),
			},
		}
	}
	kongClusterPluginWithSecret := func(name, secretID string) kongv1.KongClusterPlugin {
		return kongv1.KongClusterPlugin{
			TypeMeta: metav1.TypeMeta{
				Kind:       "KongClusterPlugin",
				APIVersion: "configuration.konghq.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				// https://github.com/kubernetes-sigs/controller-runtime/blob/22718275bffe3185276dc835d610c658f06dac07/pkg/client/fake/client.go#L247-L250
				ResourceVersion: "999",
				Name:            name,
				Annotations:     map[string]string{annotations.IngressClassKey: annotations.DefaultIngressClass},
			},
			Config: v1.JSON{
				Raw: []byte(fmt.Sprintf(`{"ca_certificates":["%s"]}`, secretID)),
			},
		}
	}

	const (
		secretID        = "8a3753e0-093b-43d9-9d39-27985c987d92"
		anotherSecretID = "99fa09c7-f849-4449-891e-19b9a0015763"
	)
	var (
		associatedPlugin           = kongPluginWithSecret("associated_plugin", secretID)
		nonAssociatedPlugin        = kongPluginWithSecret("non_associated_plugin", anotherSecretID)
		associatedClusterPlugin    = kongClusterPluginWithSecret("associated_cluster_plugin", secretID)
		nonAssociatedClusterPlugin = kongClusterPluginWithSecret("non_associated_cluster_plugin", anotherSecretID)
	)

	require.NoError(t, fake.AddToScheme(scheme.Scheme))
	require.NoError(t, gatewayv1alpha2.AddToScheme(scheme.Scheme))
	require.NoError(t, gatewayv1beta1.AddToScheme(scheme.Scheme))
	require.NoError(t, knative.AddToScheme(scheme.Scheme))

	cl := ctrlclientfake.NewClientBuilder().
		WithLists(
			&kongv1.KongPluginList{
				Items: []kongv1.KongPlugin{
					associatedPlugin,
					nonAssociatedPlugin,
				},
			},
			&kongv1.KongClusterPluginList{
				Items: []kongv1.KongClusterPlugin{
					associatedClusterPlugin,
					nonAssociatedClusterPlugin,
				},
			},
		).
		Build()
	storer := store.New(cl, annotations.DefaultIngressClass, logrus.New())

	gotPlugins := getPluginsAssociatedWithCACertSecret(context.TODO(), secretID, storer)
	expectedPlugins := []client.Object{&associatedPlugin, &associatedClusterPlugin}
	require.Len(t, gotPlugins, 2)
	require.Equal(t, expectedPlugins[0], gotPlugins[0])
	require.Equal(t, expectedPlugins[1], gotPlugins[1])
}
