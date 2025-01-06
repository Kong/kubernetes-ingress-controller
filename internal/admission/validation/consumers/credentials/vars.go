package credentials

import "k8s.io/apimachinery/pkg/util/sets"

// -----------------------------------------------------------------------------
// Validation - Vars
// -----------------------------------------------------------------------------

// SupportedTypes indicates all the Kong credential types which are supported for KongConsumer credentials.
var SupportedTypes = sets.NewString(
	"basic-auth",
	"hmac-auth",
	"jwt",
	"key-auth",
	"oauth2",
	"acl",
	"mtls-auth",
)

var (
	KeyAuthFields    = []string{"key"}
	BasicAuthFields  = []string{"username", "password"}
	HMACAuthFields   = []string{"username", "secret"}
	JWTAuthFields    = []string{"algorithm", "rsa_public_key", "key", "secret"}
	MTLsAuthFields   = []string{"subject_name"}
	OAUTH2AuthFields = []string{"name", "client_id", "client_secret", "redirect_uris"}
	ACLAuthFields    = []string{"group"}
)

var CredTypeToFields = map[string][]string{
	"key-auth":             KeyAuthFields,
	"keyauth_credential":   KeyAuthFields,
	"basic-auth":           BasicAuthFields,
	"basicauth_credential": BasicAuthFields,
	"hmac-auth":            HMACAuthFields,
	"hmacauth_credential":  HMACAuthFields,
	"jwt":                  JWTAuthFields,
	"jwt_secret":           JWTAuthFields,
	"oauth2":               OAUTH2AuthFields,
	"acl":                  ACLAuthFields,
	"mtls-auth":            MTLsAuthFields,
}
