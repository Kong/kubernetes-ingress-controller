package kongstate

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	configurationv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
)

// CredentialConflictsDetector registers all credentials and detects conflicts globally using indices.
type CredentialConflictsDetector struct {
	keyAuthsByKey               credentialIndex
	hmacAuthsByUsername         credentialIndex
	jwtAuthsByByKey             credentialIndex
	basicAuthsByUsername        credentialIndex
	oauth2CredentialsByClientID credentialIndex
}

func NewCredentialConflictsDetector() *CredentialConflictsDetector {
	return &CredentialConflictsDetector{
		keyAuthsByKey:               newCredentialIndex("key-auth on 'key'"),
		hmacAuthsByUsername:         newCredentialIndex("hmac-auth on 'username'"),
		jwtAuthsByByKey:             newCredentialIndex("jwt-auth on 'key'"),
		basicAuthsByUsername:        newCredentialIndex("basic-auth on 'username'"),
		oauth2CredentialsByClientID: newCredentialIndex("oauth2-credentials on 'client_id'"),
	}
}

func (d *CredentialConflictsDetector) RegisterForConflictDetection(
	cred any,
	credSource *corev1.Secret,
	consumerRef *configurationv1.KongConsumer,
) {
	credWithParent := CredentialWithConsumer{
		CredentialSecret: credSource,
		Consumer:         consumerRef,
	}
	switch c := cred.(type) {
	case *KeyAuth:
		d.keyAuthsByKey.add(*c.Key, credWithParent)
	case *BasicAuth:
		d.basicAuthsByUsername.add(*c.Username, credWithParent)
	case *HMACAuth:
		d.hmacAuthsByUsername.add(*c.Username, credWithParent)
	case *Oauth2Credential:
		d.oauth2CredentialsByClientID.add(*c.ClientID, credWithParent)
	case *JWTAuth:
		d.jwtAuthsByByKey.add(*c.Key, credWithParent)
	case *ACLGroup:
		// ACLs do not have any unique field to index on.
	case *MTLSAuth:
		// MTLSAuths do not have any unique field to index on.
	}
}

// DetectConflicts returns all conflicts detected.
func (d *CredentialConflictsDetector) DetectConflicts() []CredentialConflict {
	var conflicts []CredentialConflict
	for _, index := range []credentialIndex{
		d.keyAuthsByKey,
		d.hmacAuthsByUsername,
		d.jwtAuthsByByKey,
		d.basicAuthsByUsername,
		d.oauth2CredentialsByClientID,
	} {
		for _, creds := range index.index {
			// If there are more than one credential with the same key, it's a conflict.
			if len(creds) > 1 {
				for _, cred := range creds {
					conflicts = append(conflicts, CredentialConflict{
						Message:    fmt.Sprintf("conflict detected in %q index", index.name),
						Credential: cred,
					})
				}
			}
		}
	}
	return conflicts
}

// CredentialConflict represents a conflict in credentials.
type CredentialConflict struct {
	// Message is the human-readable message for the conflict.
	Message string

	// Credential is the conflicting credential and its associated consumer.
	Credential CredentialWithConsumer
}

// CredentialWithConsumer represents a credential and its associated consumer.
type CredentialWithConsumer struct {
	CredentialSecret *corev1.Secret
	Consumer         *configurationv1.KongConsumer
}

// credentialIndex is an index for credentials.
type credentialIndex struct {
	// name is the field name to index on.
	name string

	// index is the actual index.
	index map[string][]CredentialWithConsumer
}

func newCredentialIndex(name string) credentialIndex {
	return credentialIndex{
		name:  name,
		index: make(map[string][]CredentialWithConsumer),
	}
}

// add adds a credential to the index.
func (c credentialIndex) add(key string, cred CredentialWithConsumer) {
	c.index[key] = append(c.index[key], cred)
}
