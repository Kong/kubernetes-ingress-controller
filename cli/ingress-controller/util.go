package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver"
	"github.com/hbagdi/go-kong/kong"
)

func getSemVerVer(v string) (semver.Version, error) {
	// fix enterprise edition semver adding patch number
	// fix enterprise edition version with dash
	// fix bad version formats like 0.13.0preview1
	re := regexp.MustCompile(`(\d+\.\d+)(?:[\.-](\d+))?(?:\-?(.+)$|$)`)
	m := re.FindStringSubmatch(v)
	if len(m) != 4 {
		return semver.Version{}, fmt.Errorf("Unknown Kong version : '%v'", v)
	}
	if m[2] == "" {
		m[2] = "0"
	}
	if m[3] != "" {
		m[3] = "-" + strings.Replace(m[3], "enterprise-edition", "enterprise", 1)
		m[3] = strings.Replace(m[3], ".", "", -1)
	}
	v = fmt.Sprintf("%s.%s%s", m[1], m[2], m[3])
	return semver.Make(v)
}

func ensureWorkspace(client *kong.Client, workspace string) error {
	req, err := client.NewRequest("GET", "/workspaces/"+workspace, nil, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(context.TODO(), req, nil)
	if err != nil {
		if kong.IsNotFoundErr(err) {
			if err := createWorkspace(client, workspace); err != nil {
				return fmt.Errorf("creating workspace '%v': %w", workspace, err)
			}
			return nil
		}
		return fmt.Errorf("looking up workspace '%v': %w", workspace, err)
	}
	return nil
}

func createWorkspace(client *kong.Client, workspace string) error {
	body := map[string]string{"name": workspace}
	req, err := client.NewRequest("POST", "/workspaces", nil, body)
	if err != nil {
		return err
	}
	_, err = client.Do(context.TODO(), req, nil)
	return err
}
