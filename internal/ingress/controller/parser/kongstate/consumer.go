package kongstate

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// Consumer holds a Kong consumer and its plugins and credentials.
type Consumer struct {
	kong.Consumer
	Plugins    []kong.Plugin
	KeyAuths   map[string]*kong.KeyAuth
	HMACAuths  map[string]*kong.HMACAuth
	JWTAuths   map[string]*kong.JWTAuth
	BasicAuths map[string]*kong.BasicAuth
	ACLGroups  []*kong.ACLGroup

	Oauth2Creds map[string]*kong.Oauth2Credential

	K8sKongConsumer configurationv1.KongConsumer
}

// NewConsumer initializes an empty Consumer object.
func NewConsumer() Consumer {
	return Consumer{}.initEmpty()
}

func (c Consumer) initEmpty() Consumer {
	if c.KeyAuths == nil {
		c.KeyAuths = map[string]*kong.KeyAuth{}
	}
	if c.HMACAuths == nil {
		c.HMACAuths = map[string]*kong.HMACAuth{}
	}
	if c.JWTAuths == nil {
		c.JWTAuths = map[string]*kong.JWTAuth{}
	}
	if c.BasicAuths == nil {
		c.BasicAuths = map[string]*kong.BasicAuth{}
	}
	if c.Oauth2Creds == nil {
		c.Oauth2Creds = map[string]*kong.Oauth2Credential{}
	}
	return c
}

func (c *Consumer) SetCredential(log logrus.FieldLogger, credType string, credConfig interface{}) error {
	switch credType {
	case "key-auth", "keyauth_credential":
		var cred kong.KeyAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode key-auth credential: %w", err)

		}
		// TODO we perform these validity checks here because passing credentials without these fields will panic deck
		// later on. Ideally this should not be handled in the controller, but we cannot currently handle it elsewhere
		// (i.e. in deck or go-kong) without entering a sync failure loop that cannot actually report the problem
		// piece of configuration. if we can address those limitations, we should remove these checks.
		// See https://github.com/Kong/deck/pull/223 and https://github.com/Kong/kubernetes-ingress-controller/issues/532
		// for more discussion.
		if cred.Key == nil {
			return fmt.Errorf("key-auth for consumer %s is invalid: no key", *c.Username)
		}
		if _, ok := c.KeyAuths[*cred.Key]; ok {
			return fmt.Errorf("key-auth for consumer %s: duplicate key", *c.Username)
		}
		c.KeyAuths[*cred.Key] = &cred
	case "basic-auth", "basicauth_credential":
		var cred kong.BasicAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode basic-auth credential: %w", err)
		}
		if cred.Username == nil {
			return fmt.Errorf("basic-auth for consumer %s is invalid: no username", *c.Username)
		}
		if _, ok := c.BasicAuths[*cred.Username]; ok {
			return fmt.Errorf("basic-auth for consumer %s: duplicate username %q", *c.Username, *cred.Username)
		}
		c.BasicAuths[*cred.Username] = &cred
	case "hmac-auth", "hmacauth_credential":
		var cred kong.HMACAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode hmac-auth credential: %w", err)
		}
		if cred.Username == nil {
			return fmt.Errorf("hmac-auth for consumer %s is invalid: no username", *c.Username)
		}
		if _, ok := c.HMACAuths[*cred.Username]; ok {
			return fmt.Errorf("hmac-auth for consumer %s: duplicate username %q", *c.Username, *cred.Username)
		}
		c.HMACAuths[*cred.Username] = &cred
	case "oauth2":
		var cred kong.Oauth2Credential
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode oauth2 credential: %w", err)
		}
		if cred.ClientID == nil {
			return fmt.Errorf("oauth2 for consumer %s is invalid: no client_id", *c.Username)
		}
		if _, ok := c.Oauth2Creds[*cred.ClientID]; ok {
			return fmt.Errorf("oauth2 for consumer %s: duplicate client ID %q", *c.Username, *cred.ClientID)
		}
		c.Oauth2Creds[*cred.ClientID] = &cred
	case "jwt", "jwt_secret":
		var cred kong.JWTAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			log.Errorf("failed to process JWT credential: %v", err)
		}
		// This is treated specially because only this
		// field might be omitted by user under the expectation
		// that Kong will insert the default.
		// If we don't set it, decK will detect a diff and PUT this
		// credential everytime it performs a sync operation, which
		// leads to unnecessary cache invalidations in Kong.
		if cred.Algorithm == nil || *cred.Algorithm == "" {
			cred.Algorithm = kong.String("HS256")
		}
		if cred.Key == nil {
			return fmt.Errorf("jwt-auth for consumer %s is invalid: no key", *c.Username)
		}
		if _, ok := c.JWTAuths[*cred.Key]; ok {
			return fmt.Errorf("jwt-auth for consumer %s: duplicate key", *c.Username)
		}
		c.JWTAuths[*cred.Key] = &cred
	case "acl":
		var cred kong.ACLGroup
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			log.Errorf("failed to process ACL group: %v", err)
		}
		if cred.Group == nil {
			return fmt.Errorf("acl for consumer %s is invalid: no group", *c.Username)
		}
		c.ACLGroups = append(c.ACLGroups, &cred)
	default:
		return fmt.Errorf("invalid credential type: '%v'", credType)
	}
	return nil
}

func decodeCredential(credConfig interface{},
	credStructPointer interface{}) error {
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{TagName: "json",
			Result: credStructPointer,
		})
	if err != nil {
		return fmt.Errorf("failed to create a decoder: %w", err)
	}
	err = decoder.Decode(credConfig)
	if err != nil {
		return fmt.Errorf("failed to decode credential: %w", err)
	}
	return nil
}
