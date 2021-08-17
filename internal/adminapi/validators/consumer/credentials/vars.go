package credentials

import "k8s.io/apimachinery/pkg/util/sets"

// -----------------------------------------------------------------------------
// Validation - Vars
// -----------------------------------------------------------------------------

// TypeKey indicates the key in a consumer secret which identifies the type
// of credential that is being provided for the consumer.
const TypeKey = "kongCredType"

// SupportedCreds indicates all the "kongCredType"s which are supported for KongConsumer credentials.
var SupportedTypes = sets.NewString(
	"basic-auth",
	"hmac-auth",
	"jwt",
	"key-auth",
	"oauth2",
	"acl",
)
