package mgrutils

import (
	"context"

	credsvalidation "github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi/validators/consumer/credentials"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GlobalValidationIndexForCredentials builds an index of all consumer credentials
// using a given controller-runtime client. This provides an index based on
// ALL namespaces in the cluster. This can be very expensive with high numbers
// of consumer credentials, particularly if the client you provide is not cached
func GlobalValidationIndexForCredentials(ctx context.Context, managerClient client.Client, consumers []*kongv1.KongConsumer) (credsvalidation.Index, error) {
	// pull the reference secrets for credentials from each consumer in the list
	index := make(credsvalidation.Index)
	for _, consumer := range consumers {
		for _, secretName := range consumer.Credentials {
			// grab a copy of the credential secret
			secret := &corev1.Secret{}
			if err := managerClient.Get(ctx, client.ObjectKey{
				Namespace: consumer.Namespace,
				Name:      secretName,
			}, secret); err != nil {
				if errors.IsNotFound(err) { // ignore missing secrets
					continue
				}
				return nil, err
			}

			// add the credential secret to the index
			if err := index.ValidateCredentialsForUniqueKeyConstraints(consumer.Name, secret); err != nil {
				return nil, err
			}
		}
	}

	return index, nil
}
