package kongstate

import "github.com/kong/go-kong/kong"

var redactedString = kong.String("REDACTED")

// KeyAuth represents a key-auth credential.
type KeyAuth struct {
	kong.KeyAuth
}

// HMACAuth represents a HMAC credential.
type HMACAuth struct {
	kong.HMACAuth
}

// JWTAuth represents a JWT credential.
type JWTAuth struct {
	kong.JWTAuth
}

// BasicAuth represents a basic authentication credential.
type BasicAuth struct {
	kong.BasicAuth
}

// ACLGroup represents an ACL associated with a consumer. Due to ACL implementation in Kong being similar to
// credentials, ACLs are treated as credentials, too.
type ACLGroup struct {
	kong.ACLGroup
}

// Oauth2Credential represents an OAuth2 client configuration including credentials.
type Oauth2Credential struct {
	kong.Oauth2Credential
}
