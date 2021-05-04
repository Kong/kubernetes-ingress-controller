package mgrutils

import (
	"context"
	"fmt"
	"github.com/kong/go-kong/kong"
)

// TODO these functions have been copied from 1.x cli/ingress-controller/util.go
// https://github.com/Kong/go-kong/pull/56 provides go-kong replacements for them,
// but is not yet released

func EnsureWorkspace(ctx context.Context, client *kong.Client, workspace string) error {
	req, err := client.NewRequest("GET", "/workspaces/"+workspace, nil, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(ctx, req, nil)
	if err != nil {
		if kong.IsNotFoundErr(err) {
			if err := createWorkspace(ctx, client, workspace); err != nil {
				return fmt.Errorf("creating workspace '%v': %w", workspace, err)
			}
			return nil
		}
		return fmt.Errorf("looking up workspace '%v': %w", workspace, err)
	}
	return nil
}

func createWorkspace(ctx context.Context, client *kong.Client, workspace string) error {
	body := map[string]string{"name": workspace}
	req, err := client.NewRequest("POST", "/workspaces", nil, body)
	if err != nil {
		return err
	}
	_, err = client.Do(ctx, req, nil)
	return err
}
