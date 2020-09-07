package main

import (
	"testing"
)

func TestFixVersion(t *testing.T) {
	validVersions := map[string]string{
		"0.14.1":                          "0.14.1",
		"0.14.2rc":                        "0.14.2-rc",
		"0.14.2rc1":                       "0.14.2-rc1",
		"0.14.2preview":                   "0.14.2-preview",
		"0.14.2preview1":                  "0.14.2-preview1",
		"0.33-enterprise-edition":         "0.33.0-enterprise",
		"0.33-1-enterprise-edition":       "0.33.1-enterprise",
		"1.3.0.0-enterprise-edition-lite": "1.3.0-0-enterprise-lite",
		"1.3.0.0-enterprise-lite":         "1.3.0-0-enterprise-lite",
	}
	for inputVersion, expectedVersion := range validVersions {
		v, err := getSemVerVer(inputVersion)
		if err != nil {
			t.Errorf("error converting %s: %v", inputVersion, err)
		} else if v.String() != expectedVersion {
			t.Errorf("converting %s, expecting %s, getting %s", inputVersion, expectedVersion, v.String())
		}
	}

	invalidVersions := []string{
		"",
		"0-1-1",
	}
	for _, inputVersion := range invalidVersions {
		_, err := getSemVerVer(inputVersion)
		if err == nil {
			t.Errorf("expecting error converting %s, getting no errors", inputVersion)
		}
	}
}
