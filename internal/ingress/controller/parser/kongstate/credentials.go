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

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (c *KeyAuth) SanitizedCopy() *KeyAuth {
	return &KeyAuth{
		kong.KeyAuth{
			// Consumer field omitted
			CreatedAt: c.CreatedAt,
			ID:        c.ID,
			Key:       redactedString,
			Tags:      c.Tags,
		},
	}
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (c *HMACAuth) SanitizedCopy() *HMACAuth {
	return &HMACAuth{
		kong.HMACAuth{
			// Consumer field omitted
			CreatedAt: c.CreatedAt,
			ID:        c.ID,
			Username:  c.Username,
			Secret:    redactedString,
			Tags:      c.Tags,
		},
	}
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (c *JWTAuth) SanitizedCopy() *JWTAuth {
	return &JWTAuth{
		kong.JWTAuth{
			// Consumer field omitted
			CreatedAt:    c.CreatedAt,
			ID:           c.ID,
			Algorithm:    c.Algorithm,
			Key:          c.Key, // despite field name, "key" is an identifier
			RSAPublicKey: c.RSAPublicKey,
			Secret:       redactedString,
			Tags:         c.Tags,
		},
	}
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (c *BasicAuth) SanitizedCopy() *BasicAuth {
	return &BasicAuth{
		kong.BasicAuth{
			// Consumer field omitted
			CreatedAt: c.CreatedAt,
			ID:        c.ID,
			Username:  c.Username,
			Password:  redactedString,
			Tags:      c.Tags,
		},
	}
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (c *Oauth2Credential) SanitizedCopy() *Oauth2Credential {
	return &Oauth2Credential{
		kong.Oauth2Credential{
			// Consumer field omitted
			CreatedAt:    c.CreatedAt,
			ID:           c.ID,
			Name:         c.Name,
			ClientID:     c.ClientID,
			ClientSecret: redactedString,
			RedirectURIs: c.RedirectURIs,
			Tags:         c.Tags,
		},
	}
}
