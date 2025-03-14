package kongstate

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	configurationv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// Consumer holds a Kong consumer and its plugins and credentials.
type Consumer struct {
	kong.Consumer
	Plugins        []kong.Plugin
	ConsumerGroups []kong.ConsumerGroup

	KeyAuths   []*KeyAuth
	HMACAuths  []*HMACAuth
	JWTAuths   []*JWTAuth
	BasicAuths []*BasicAuth
	ACLGroups  []*ACLGroup

	Oauth2Creds []*Oauth2Credential
	MTLSAuths   []*MTLSAuth

	K8sKongConsumer configurationv1.KongConsumer
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (c *Consumer) SanitizedCopy(uuidGenerator util.UUIDGenerator) Consumer {
	return Consumer{
		Consumer:       c.Consumer,
		Plugins:        c.Plugins,
		ConsumerGroups: c.ConsumerGroups,
		KeyAuths: func() []*KeyAuth {
			if c.KeyAuths == nil {
				return nil
			}
			return lo.Map(c.KeyAuths, func(c *KeyAuth, _ int) *KeyAuth {
				return c.SanitizedCopy(uuidGenerator)
			})
		}(),
		HMACAuths: func() []*HMACAuth {
			if c.HMACAuths == nil {
				return nil
			}
			return lo.Map(c.HMACAuths, func(c *HMACAuth, _ int) *HMACAuth {
				return c.SanitizedCopy()
			})
		}(),
		JWTAuths: func() []*JWTAuth {
			if c.JWTAuths == nil {
				return nil
			}
			return lo.Map(c.JWTAuths, func(c *JWTAuth, _ int) *JWTAuth {
				return c.SanitizedCopy()
			})
		}(),
		BasicAuths: func() []*BasicAuth {
			if c.BasicAuths == nil {
				return nil
			}
			return lo.Map(c.BasicAuths, func(c *BasicAuth, _ int) *BasicAuth {
				return c.SanitizedCopy()
			})
		}(),
		Oauth2Creds: func() []*Oauth2Credential {
			if c.Oauth2Creds == nil {
				return nil
			}
			return lo.Map(c.Oauth2Creds, func(c *Oauth2Credential, _ int) *Oauth2Credential {
				return c.SanitizedCopy()
			})
		}(),
		ACLGroups:       c.ACLGroups,
		MTLSAuths:       c.MTLSAuths,
		K8sKongConsumer: c.K8sKongConsumer,
	}
}

func (c *Consumer) SetCredential(credType string, credConfig interface{}, tags []*string) (interface{}, error) {
	switch credType {
	case "key-auth", "keyauth_credential":
		cred, err := NewKeyAuth(credConfig)
		if err != nil {
			return nil, err
		}
		cred.Tags = tags
		c.KeyAuths = append(c.KeyAuths, cred)
		return cred, nil
	case "basic-auth", "basicauth_credential":
		cred, err := NewBasicAuth(credConfig)
		if err != nil {
			return nil, err
		}
		cred.Tags = tags
		c.BasicAuths = append(c.BasicAuths, cred)
		return cred, nil
	case "hmac-auth", "hmacauth_credential":
		cred, err := NewHMACAuth(credConfig)
		if err != nil {
			return nil, err
		}
		cred.Tags = tags
		c.HMACAuths = append(c.HMACAuths, cred)
		return cred, nil
	case "oauth2":
		cred, err := NewOauth2Credential(credConfig)
		if err != nil {
			return nil, err
		}
		cred.Tags = tags
		c.Oauth2Creds = append(c.Oauth2Creds, cred)
		return cred, nil
	case "jwt", "jwt_secret":
		cred, err := NewJWTAuth(credConfig)
		if err != nil {
			return nil, err
		}
		cred.Tags = tags
		c.JWTAuths = append(c.JWTAuths, cred)
		return cred, nil
	case "acl":
		cred, err := NewACLGroup(credConfig)
		if err != nil {
			return nil, err
		}
		cred.Tags = tags
		c.ACLGroups = append(c.ACLGroups, cred)
		return cred, nil
	case "mtls-auth":
		cred, err := NewMTLSAuth(credConfig)
		if err != nil {
			return nil, err
		}
		cred.Tags = tags
		c.MTLSAuths = append(c.MTLSAuths, cred)
		return cred, nil
	default:
		return nil, fmt.Errorf("invalid credential type: '%v'", credType)
	}
}
