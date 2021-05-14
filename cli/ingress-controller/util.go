package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"k8s.io/client-go/tools/cache"
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

func ensureWorkspace(ctx context.Context, client *kong.Client, workspaceName string) error {
	exists, err := client.Workspaces.Exists(ctx, workspaceName)
	if err != nil {
		return fmt.Errorf("looking up workspace '%v': %w", workspaceName, err)
	}
	if !exists {
		if err := createWorkspace(ctx, client, workspaceName); err != nil {
			return fmt.Errorf("creating workspace '%v': %w", workspaceName, err)
		}
	}
	return nil
}

func createWorkspace(ctx context.Context, client *kong.Client, workspaceName string) error {
	workspace := &kong.Workspace{
		Name: kong.String(workspaceName),
	}
	_, err := client.Workspaces.Create(ctx, workspace)
	return err
}

func newEmptyStore() cache.Store {
	return cache.NewStore(func(interface{}) (string, error) { return "", errors.New("this store cannot add elements") })
}
