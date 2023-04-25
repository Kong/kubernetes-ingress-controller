package deckgen

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
)

// GenerateSHA generates a SHA256 checksum of targetContent, with the purpose
// of change detection.
func GenerateSHA(targetContent *file.Content) ([]byte, error) {
	jsonConfig, err := json.Marshal(targetContent)
	if err != nil {
		return nil, fmt.Errorf("marshaling Kong declarative configuration to JSON: %w", err)
	}

	shaSum := sha256.Sum256(jsonConfig)
	return shaSum[:], nil
}

// CleanUpNullsInPluginConfigs modifies `state` by deleting plugin config map keys that have nil as their value.
func CleanUpNullsInPluginConfigs(state *file.Content) {
	for _, s := range state.Services {
		for _, p := range s.Plugins {
			for k, v := range p.Config {
				if v == nil {
					delete(p.Config, k)
				}
			}
		}
		for _, r := range state.Routes {
			for _, p := range r.Plugins {
				for k, v := range p.Config {
					if v == nil {
						delete(p.Config, k)
					}
				}
			}
		}
	}

	for _, c := range state.Consumers {
		for _, p := range c.Plugins {
			for k, v := range p.Config {
				if v == nil {
					delete(p.Config, k)
				}
			}
		}
	}

	for _, p := range state.Plugins {
		for k, v := range p.Config {
			if v == nil {
				delete(p.Config, k)
			}
		}
	}
}

// GetFCertificateFromKongCert converts a kong.Certificate to a file.FCertificate.
func GetFCertificateFromKongCert(kongCert kong.Certificate) file.FCertificate {
	var res file.FCertificate
	if kongCert.ID != nil {
		res.ID = kong.String(*kongCert.ID)
	}
	if kongCert.Key != nil {
		res.Key = kong.String(*kongCert.Key)
	}
	if kongCert.Cert != nil {
		res.Cert = kong.String(*kongCert.Cert)
	}
	res.SNIs = getSNIs(kongCert.SNIs)
	return res
}

func getSNIs(names []*string) []kong.SNI {
	var snis []kong.SNI
	for _, name := range names {
		snis = append(snis, kong.SNI{
			Name: kong.String(*name),
		})
	}
	return snis
}

// PluginString returns a string representation of a FPlugin suitable as a sorting key.
//
// Deprecated. To be replaced by a predicate that compares two FPlugins.
func PluginString(plugin file.FPlugin) string {
	result := ""
	if plugin.Name != nil {
		result = *plugin.Name
	}
	if plugin.Consumer != nil && plugin.Consumer.ID != nil {
		result += *plugin.Consumer.ID
	}
	if plugin.Route != nil && plugin.Route.ID != nil {
		result += *plugin.Route.ID
	}
	if plugin.Service != nil && plugin.Service.ID != nil {
		result += *plugin.Service.ID
	}
	return result
}
