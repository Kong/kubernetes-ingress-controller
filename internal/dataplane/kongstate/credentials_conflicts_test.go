package kongstate_test

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
)

func TestCredentialsConflictsDetector(t *testing.T) {
	someSecret := func(name string) *corev1.Secret {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
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
			name: "key auth - 2 conflict and 1 no conflict",
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
		{
			name: "basic auth - no conflicts",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.BasicAuth{
						BasicAuth: kong.BasicAuth{
							Username: kong.String("user1"),
						},
					},
					secret:   someSecret("basic1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.BasicAuth{
						BasicAuth: kong.BasicAuth{
							Username: kong.String("user2"),
						},
					},
					secret:   someSecret("basic2"),
					consumer: someConsumer("consumer2"),
				},
			},
			expectedConflicts: nil,
		},
		{
			name: "basic auth - 2 conflict and 1 no conflict",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.BasicAuth{
						BasicAuth: kong.BasicAuth{
							Username: kong.String("user1"),
						},
					},
					secret:   someSecret("basic1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.BasicAuth{
						BasicAuth: kong.BasicAuth{
							Username: kong.String("user1"),
						},
					},
					secret:   someSecret("basic1"),
					consumer: someConsumer("consumer2"),
				},
				{
					credential: &kongstate.BasicAuth{
						BasicAuth: kong.BasicAuth{
							Username: kong.String("user2"),
						},
					},
					secret:   someSecret("basic2"),
					consumer: someConsumer("consumer3"),
				},
			},
			expectedConflicts: []kongstate.CredentialConflict{
				expectedConflict(`conflict detected in "basic-auth on 'username'" index`, someSecret("basic1"), someConsumer("consumer1")),
				expectedConflict(`conflict detected in "basic-auth on 'username'" index`, someSecret("basic1"), someConsumer("consumer2")),
			},
		},
		{
			name: "jwt auth - no conflicts",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.JWTAuth{
						JWTAuth: kong.JWTAuth{
							Key: kong.String("key1"),
						},
					},
					secret:   someSecret("jwt1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.JWTAuth{
						JWTAuth: kong.JWTAuth{
							Key: kong.String("key2"),
						},
					},
					secret:   someSecret("jwt2"),
					consumer: someConsumer("consumer2"),
				},
			},
			expectedConflicts: nil,
		},
		{
			name: "jwt auth - 2 conflict and 1 no conflict",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.JWTAuth{
						JWTAuth: kong.JWTAuth{
							Key: kong.String("key1"),
						},
					},
					secret:   someSecret("jwt1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.JWTAuth{
						JWTAuth: kong.JWTAuth{
							Key: kong.String("key1"),
						},
					},
					secret:   someSecret("jwt1"),
					consumer: someConsumer("consumer2"),
				},
				{
					credential: &kongstate.JWTAuth{
						JWTAuth: kong.JWTAuth{
							Key: kong.String("key2"),
						},
					},
					secret:   someSecret("jwt2"),
					consumer: someConsumer("consumer3"),
				},
			},
			expectedConflicts: []kongstate.CredentialConflict{
				expectedConflict(`conflict detected in "jwt-auth on 'key'" index`, someSecret("jwt1"), someConsumer("consumer1")),
				expectedConflict(`conflict detected in "jwt-auth on 'key'" index`, someSecret("jwt1"), someConsumer("consumer2")),
			},
		},
		{
			name: "hmac auth - no conflicts",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.HMACAuth{
						HMACAuth: kong.HMACAuth{
							Username: kong.String("user1"),
						},
					},
					secret:   someSecret("hmac1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.HMACAuth{
						HMACAuth: kong.HMACAuth{
							Username: kong.String("user2"),
						},
					},
					secret:   someSecret("hmac2"),
					consumer: someConsumer("consumer2"),
				},
			},
			expectedConflicts: nil,
		},
		{
			name: "hmac auth - 2 conflict and 1 no conflict",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.HMACAuth{
						HMACAuth: kong.HMACAuth{
							Username: kong.String("user1"),
						},
					},
					secret:   someSecret("hmac1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.HMACAuth{
						HMACAuth: kong.HMACAuth{
							Username: kong.String("user1"),
						},
					},
					secret:   someSecret("hmac1"),
					consumer: someConsumer("consumer2"),
				},
				{
					credential: kongstate.HMACAuth{
						HMACAuth: kong.HMACAuth{
							Username: kong.String("user2"),
						},
					},
					secret:   someSecret("hmac2"),
					consumer: someConsumer("consumer3"),
				},
			},
			expectedConflicts: []kongstate.CredentialConflict{
				expectedConflict(`conflict detected in "hmac-auth on 'username'" index`, someSecret("hmac1"), someConsumer("consumer1")),
				expectedConflict(`conflict detected in "hmac-auth on 'username'" index`, someSecret("hmac1"), someConsumer("consumer2")),
			},
		},
		{
			name: "oauth2 credential - no conflicts",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.Oauth2Credential{
						Oauth2Credential: kong.Oauth2Credential{
							ClientID: kong.String("client1"),
						},
					},
					secret:   someSecret("oauth2-1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.Oauth2Credential{
						Oauth2Credential: kong.Oauth2Credential{
							ClientID: kong.String("client2"),
						},
					},
					secret:   someSecret("oauth2-2"),
					consumer: someConsumer("consumer2"),
				},
			},
			expectedConflicts: nil,
		},
		{
			name: "oauth2 credential - 2 conflict and 1 no conflict",
			credentials: []testCredentialDetails{
				{
					credential: &kongstate.Oauth2Credential{
						Oauth2Credential: kong.Oauth2Credential{
							ClientID: kong.String("client1"),
						},
					},
					secret:   someSecret("oauth2-1"),
					consumer: someConsumer("consumer1"),
				},
				{
					credential: &kongstate.Oauth2Credential{
						Oauth2Credential: kong.Oauth2Credential{
							ClientID: kong.String("client1"),
						},
					},
					secret:   someSecret("oauth2-1"),
					consumer: someConsumer("consumer2"),
				},
				{
					credential: kongstate.Oauth2Credential{
						Oauth2Credential: kong.Oauth2Credential{
							ClientID: kong.String("client2"),
						},
					},
					secret:   someSecret("oauth2-2"),
					consumer: someConsumer("consumer3"),
				},
			},
			expectedConflicts: []kongstate.CredentialConflict{
				expectedConflict(`conflict detected in "oauth2-credentials on 'client_id'" index`, someSecret("oauth2-1"), someConsumer("consumer1")),
				expectedConflict(`conflict detected in "oauth2-credentials on 'client_id'" index`, someSecret("oauth2-1"), someConsumer("consumer2")),
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
