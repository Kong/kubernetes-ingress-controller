package kongstate_test

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: Implement this test for all cred types
func TestCredentialsConflictsDetector(t *testing.T) {
	someSecret := func(name string) *corev1.Secret {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "default",
			},
		}
	}
	someConsumer := func(name string) *kongv1.KongConsumer {
		return &kongv1.KongConsumer{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
		}
	}
	expectedConflict := func(msg string, secret *corev1.Secret, consumer *kongv1.KongConsumer) kongstate.CredentialConflict {
		return kongstate.CredentialConflict{
			Message: msg,
			Credential: kongstate.CredentialWithConsumer{
				CredentialSecret: secret,
				Consumer:         consumer,
			},
		}
	}

	type testCredentialDetails struct {
		credential any
		secret     *corev1.Secret
		consumer   *kongv1.KongConsumer
	}
	testCases := []struct {
		name              string
		credentials       []testCredentialDetails
		expectedConflicts []kongstate.CredentialConflict
	}{
		{
			name: "key-auth - no conflicts",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.KeyAuth{
						KeyAuth: kong.KeyAuth{
							Key: kong.String("key1"),
						},
					},
					secret:   someSecret("key1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.KeyAuth{
						KeyAuth: kong.KeyAuth{
							Key: kong.String("key2"),
						},
					},
					secret:   someSecret("key2"),
					consumer: someConsumer("consumer2"),
				},
			},
			expectedConflicts: nil,
		},
		{
			name: "2 conflict and 1 no conflict",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.KeyAuth{
						KeyAuth: kong.KeyAuth{
							Key: kong.String("key1"),
						},
					},
					secret:   someSecret("key1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.KeyAuth{
						KeyAuth: kong.KeyAuth{
							Key: kong.String("key1"),
						},
					},
					secret:   someSecret("key1"),
					consumer: someConsumer("consumer2"),
				},
				{
					credential: &kongstate.KeyAuth{
						KeyAuth: kong.KeyAuth{
							Key: kong.String("key2"),
						},
					},
				},
			},
			expectedConflicts: []kongstate.CredentialConflict{
				expectedConflict(`conflict detected in "key-auth on 'key'" index`, someSecret("key1"), someConsumer("consumer1")),
				expectedConflict(`conflict detected in "key-auth on 'key'" index`, someSecret("key1"), someConsumer("consumer2")),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conflictsDetector := kongstate.NewCredentialConflictsDetector()
			for _, cred := range tc.credentials {
				conflictsDetector.RegisterForConflictDetection(cred.credential, cred.secret, cred.consumer)
			}
			conflicts := conflictsDetector.DetectConflicts()
			require.Equal(t, tc.expectedConflicts, conflicts)
		})
	}

}
