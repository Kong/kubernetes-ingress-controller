package parser

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
	KeyAuths   []*kong.KeyAuth
	HMACAuths  []*kong.HMACAuth
	JWTAuths   []*kong.JWTAuth
	BasicAuths []*kong.BasicAuth
	ACLGroups  []*kong.ACLGroup

	Oauth2Creds []*kong.Oauth2Credential

	k8sKongConsumer configurationv1.KongConsumer
}

func (c *Consumer) setCredential(log logrus.FieldLogger, credType string, credConfig interface{}) error {
	switch credType {
	case "key-auth", "keyauth_credential":
		var cred kong.KeyAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode key-auth credential: %w", err)

		}
		c.KeyAuths = append(c.KeyAuths, &cred)
	case "basic-auth", "basicauth_credential":
		var cred kong.BasicAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode basic-auth credential: %w", err)
		}
		c.BasicAuths = append(c.BasicAuths, &cred)
	case "hmac-auth", "hmacauth_credential":
		var cred kong.HMACAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode hmac-auth credential: %w", err)
		}
		c.HMACAuths = append(c.HMACAuths, &cred)
	case "oauth2":
		var cred kong.Oauth2Credential
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode oauth2 credential: %w", err)
		}
		c.Oauth2Creds = append(c.Oauth2Creds, &cred)
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
		c.JWTAuths = append(c.JWTAuths, &cred)
	case "acl":
		var cred kong.ACLGroup
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			log.Errorf("failed to process ACL group: %v", err)
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
