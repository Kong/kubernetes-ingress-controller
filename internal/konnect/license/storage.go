package license

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/license"
)

const (
	// licenseResourceNamePrefix is the prefix of the secret name storing the konnect license.
	licenseResourceNamePrefix = "konnect-license-"
	// LabelKeyManagedBy is the label key to mark that the secret is managed by KIC.
	LabelKeyManagedBy = "managed-by"
)

// Storer is the interface to store license fetched from Konnect and load license from storage if failed to fetch.
type Storer interface {
	Store(context.Context, license.KonnectLicense) error
	Load(context.Context) (license.KonnectLicense, error)
}

type SecretLicenseStore struct {
	cl             client.Client
	namespace      string
	controlPlaneID string
}

var _ Storer = &SecretLicenseStore{}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;update

// NewSecretLicenseStore creates a storage to store the fetched license to a secret.
func NewSecretLicenseStore(cl client.Client, namespace, controlPlaneID string) *SecretLicenseStore {
	return &SecretLicenseStore{
		cl:             cl,
		namespace:      namespace,
		controlPlaneID: controlPlaneID,
	}
}

// Store stores license to the secret `konnect-license-<cpid>`.
func (s *SecretLicenseStore) Store(ctx context.Context, l license.KonnectLicense) error {
	secret := &corev1.Secret{}
	err := s.cl.Get(ctx, k8stypes.NamespacedName{
		Namespace: s.namespace,
		Name:      licenseResourceNamePrefix + s.controlPlaneID,
	}, secret)
	if err != nil {
		return err
	}

	// Add label to mark that the secret is managed by KIC.
	if secret.Labels == nil {
		secret.Labels = map[string]string{}
	}
	secret.Labels[LabelKeyManagedBy] = "konghq.com/ingress-controller"

	secret.StringData = map[string]string{
		"payload":    l.Payload,
		"updated_at": strconv.FormatInt(l.UpdatedAt.Unix(), 10),
		"id":         l.ID,
	}
	return s.cl.Update(ctx, secret)
}

// Load loads the license from the secret from secret `konnect-license-<cpid>`.
func (s *SecretLicenseStore) Load(
	ctx context.Context,
) (license.KonnectLicense, error) {
	secret := &corev1.Secret{}
	err := s.cl.Get(ctx, k8stypes.NamespacedName{
		Namespace: s.namespace,
		Name:      licenseResourceNamePrefix + s.controlPlaneID,
	}, secret)
	if err != nil {
		return license.KonnectLicense{}, err
	}

	if (!lo.HasKey(secret.Data, "payload")) || (!lo.HasKey(secret.Data, "updated_at")) || (!lo.HasKey(secret.Data, "id")) {
		return license.KonnectLicense{}, fmt.Errorf("missing required key in secret %s", secret.Name)
	}

	decodedPayload, err := base64.StdEncoding.DecodeString(string(secret.Data["payload"]))
	if err != nil {
		return license.KonnectLicense{}, fmt.Errorf("failed to decode payload of license stored in secret %s: %w", secret.Name, err)
	}
	decodedID, err := base64.StdEncoding.DecodeString(string(secret.Data["id"]))
	if err != nil {
		return license.KonnectLicense{}, fmt.Errorf("failed to decode id of license stored in secret %s: %w", secret.Name, err)
	}
	decodedUpdateAt, err := base64.StdEncoding.DecodeString(string(secret.Data["updated_at"]))
	if err != nil {
		return license.KonnectLicense{}, fmt.Errorf("failed to decode updated_at of license stored in secret %s: %w", secret.Name, err)
	}
	updateAt, err := strconv.ParseInt(string(decodedUpdateAt), 10, 64)
	if err != nil {
		return license.KonnectLicense{}, fmt.Errorf("failed to parse updated_at as timestamp of license stored in secret %s: %w", secret.Name, err)
	}
	return license.KonnectLicense{
		Payload:   string(decodedPayload),
		UpdatedAt: time.Unix(updateAt, 0),
		ID:        string(decodedID),
	}, nil
}
