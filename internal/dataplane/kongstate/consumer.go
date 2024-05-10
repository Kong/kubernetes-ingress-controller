package kongstate

import (
	"fmt"

	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
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

	K8sKongConsumer kongv1.KongConsumer
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (c *Consumer) SanitizedCopy(uuidGenerator util.UUIDGenerator) *Consumer {
	return &Consumer{
		Consumer: c.Consumer,
		Plugins:  c.Plugins,
		KeyAuths: func() (res []*KeyAuth) {
			for _, v := range c.KeyAuths {
				res = append(res, v.SanitizedCopy(uuidGenerator))
			}
			return
		}(),
		HMACAuths: func() (res []*HMACAuth) {
			for _, v := range c.HMACAuths {
				res = append(res, v.SanitizedCopy())
			}
			return
		}(),
		JWTAuths: func() (res []*JWTAuth) {
			for _, v := range c.JWTAuths {
				res = append(res, v.SanitizedCopy())
			}
			return
		}(),
		BasicAuths: func() (res []*BasicAuth) {
			for _, v := range c.BasicAuths {
				res = append(res, v.SanitizedCopy())
			}
			return
		}(),
		Oauth2Creds: func() (res []*Oauth2Credential) {
			for _, v := range c.Oauth2Creds {
				res = append(res, v.SanitizedCopy())
			}
			return
		}(),
		ACLGroups:       c.ACLGroups,
		MTLSAuths:       c.MTLSAuths,
		K8sKongConsumer: c.K8sKongConsumer,
	}
}

func (c *Consumer) SetCredential(credType string, credConfig interface{}, tags []*string) error {
	switch credType {
	case "key-auth", "keyauth_credential":
		cred, err := NewKeyAuth(credConfig)
		if err != nil {
			return err
		}
		cred.Tags = tags
		c.KeyAuths = append(c.KeyAuths, cred)
	case "basic-auth", "basicauth_credential":
		cred, err := NewBasicAuth(credConfig)
		if err != nil {
			return err
		}
		cred.Tags = tags
		c.BasicAuths = append(c.BasicAuths, cred)
	case "hmac-auth", "hmacauth_credential":
		cred, err := NewHMACAuth(credConfig)
		if err != nil {
			return err
		}
		cred.Tags = tags
		c.HMACAuths = append(c.HMACAuths, cred)
	case "oauth2":
		cred, err := NewOauth2Credential(credConfig)
		if err != nil {
			return err
		}
		cred.Tags = tags
		c.Oauth2Creds = append(c.Oauth2Creds, cred)
	case "jwt", "jwt_secret":
		cred, err := NewJWTAuth(credConfig)
		if err != nil {
			return err
		}
		cred.Tags = tags
		c.JWTAuths = append(c.JWTAuths, cred)
	case "acl":
		cred, err := NewACLGroup(credConfig)
		if err != nil {
			return err
		}
		cred.Tags = tags
		c.ACLGroups = append(c.ACLGroups, cred)
	case "mtls-auth":
		cred, err := NewMTLSAuth(credConfig)
		if err != nil {
			return err
		}
		cred.Tags = tags
		c.MTLSAuths = append(c.MTLSAuths, cred)
	default:
		return fmt.Errorf("invalid credential type: '%v'", credType)
	}
	return nil
}
