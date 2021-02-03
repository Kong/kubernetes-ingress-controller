package kongstate

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
)

// Consumer holds a Kong consumer and its plugins and credentials.
type Consumer struct {
	kong.Consumer
	Plugins    []kong.Plugin
	KeyAuths   []*KeyAuth
	HMACAuths  []*HMACAuth
	JWTAuths   []*JWTAuth
	BasicAuths []*BasicAuth
	ACLGroups  []*ACLGroup

	Oauth2Creds []*Oauth2Credential

	K8sKongConsumer configurationv1.KongConsumer
}

// SanitizedCopy returns a shallow copy with sensitive values redacted best-effort.
func (c *Consumer) SanitizedCopy() *Consumer {
	return &Consumer{
		Consumer: c.Consumer,
		Plugins:  c.Plugins,
		KeyAuths: func() (res []*KeyAuth) {
			for _, v := range c.KeyAuths {
				res = append(res, v.SanitizedCopy())
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
		K8sKongConsumer: c.K8sKongConsumer,
	}
}

func (c *Consumer) SetCredential(credType string, credConfig interface{}) error {
	switch credType {
	case "key-auth", "keyauth_credential":
		cred, err := NewKeyAuth(credConfig)
		if err != nil {
			return err
		}
		c.KeyAuths = append(c.KeyAuths, cred)
	case "basic-auth", "basicauth_credential":
		cred, err := NewBasicAuth(credConfig)
		if err != nil {
			return err
		}
		c.BasicAuths = append(c.BasicAuths, cred)
	case "hmac-auth", "hmacauth_credential":
		cred, err := NewHMACAuth(credConfig)
		if err != nil {
			return err
		}
		c.HMACAuths = append(c.HMACAuths, cred)
	case "oauth2":
		cred, err := NewOauth2Credential(credConfig)
		if err != nil {
			return err
		}
		c.Oauth2Creds = append(c.Oauth2Creds, cred)
	case "jwt", "jwt_secret":
		cred, err := NewJWTAuth(credConfig)
		if err != nil {
			return err
		}
		c.JWTAuths = append(c.JWTAuths, cred)
	case "acl":
		cred, err := NewACLGroup(credConfig)
		if err != nil {
			return err
		}
		c.ACLGroups = append(c.ACLGroups, cred)
	default:
		return fmt.Errorf("invalid credential type: '%v'", credType)
	}
	return nil
}
