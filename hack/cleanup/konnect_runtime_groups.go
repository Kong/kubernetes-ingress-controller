package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/roles"
	rg "github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/runtimegroups"
)

const (
	konnectRuntimeGroupsBaseURL     = "https://us.kic.api.konghq.tech/v2"
	konnectRuntimeGroupsLimit       = 100
	konnectRolesBaseURL             = "https://global.api.konghq.tech/v2"
	createdInTestsRuntimeGroupLabel = "created_in_tests"
	timeUntilRuntimeGroupOrphaned   = time.Hour
)

// cleanupKonnectRuntimeGroups deletes orphaned runtime groups created by the tests and their roles.
func cleanupKonnectRuntimeGroups(ctx context.Context, log logr.Logger) error {
	rgClient, err := rg.NewClientWithResponses(konnectRuntimeGroupsBaseURL, rg.WithRequestEditorFn(
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+konnectAccessToken)
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to create runtime groups client: %w", err)
	}

	orphanedRGs, err := findOrphanedRuntimeGroups(ctx, log, rgClient)
	if err != nil {
		return fmt.Errorf("failed to find orphaned runtime groups: %w", err)
	}
	if err := deleteRuntimeGroups(ctx, log, orphanedRGs, rgClient); err != nil {
		return fmt.Errorf("failed to delete runtime groups: %w", err)
	}

	// We have to manually delete roles created for the runtime group because Konnect doesn't do it automatically.
	// If we don't do it, we will eventually hit a problem with Konnect APIs answering our requests with 504s
	// because of a performance issue when there's too many roles for the account
	// (see https://konghq.atlassian.net/browse/TPS-1319).
	//
	// We can drop this once the automated cleanup is implemented on Konnect side:
	// https://konghq.atlassian.net/browse/TPS-1453.
	rolesClient := roles.NewClient(&http.Client{}, konnectRolesBaseURL, konnectAccessToken)
	rolesToDelete, err := findOrphanedRolesToDelete(ctx, log, orphanedRGs, rolesClient)
	if err != nil {
		return fmt.Errorf("failed to list runtime group roles to delete: %w", err)
	}
	if err := deleteRoles(ctx, log, rolesToDelete, rolesClient); err != nil {
		return fmt.Errorf("failed to delete runtime group roles: %w", err)
	}

	return nil
}

// findOrphanedRuntimeGroups finds runtime groups that were created by the tests and are older than timeUntilRuntimeGroupOrphaned.
func findOrphanedRuntimeGroups(ctx context.Context, log logr.Logger, c *rg.ClientWithResponses) ([]types.UUID, error) {
	response, err := c.ListRuntimeGroupsWithResponse(ctx, &rg.ListRuntimeGroupsParams{
		PageSize: lo.ToPtr(konnectRuntimeGroupsLimit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list runtime groups: %w", err)
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to list runtime groups, status: %s, body: %s", response.Status(), string(response.Body))
	}
	if response.JSON200 == nil || response.JSON200.Data == nil {
		return nil, fmt.Errorf("no data in the response, status: %s, body: %s", response.Status(), string(response.Body))
	}

	var orphanedRuntimeGroups []types.UUID
	for _, runtimeGroup := range *response.JSON200.Data {
		if runtimeGroup.Labels == nil || (*runtimeGroup.Labels)[createdInTestsRuntimeGroupLabel] != "true" {
			log.Info("runtime group was not created by the tests, skipping", "name", *runtimeGroup.Name)
			continue
		}
		if runtimeGroup.CreatedAt == nil {
			log.Info("runtime group has no creation timestamp, skipping", "name", *runtimeGroup.Name)
			continue
		}
		orphanedAfter := (*runtimeGroup.CreatedAt).Add(timeUntilRuntimeGroupOrphaned)
		if !time.Now().After(orphanedAfter) {
			log.Info("runtime group is not old enough to be considered orphaned, skipping", "name", *runtimeGroup.Name, "created_at", *runtimeGroup.CreatedAt)
			continue
		}
		orphanedRuntimeGroups = append(orphanedRuntimeGroups, *runtimeGroup.Id)
	}
	return orphanedRuntimeGroups, nil
}

// deleteRuntimeGroups deletes runtime groups by their IDs.
func deleteRuntimeGroups(ctx context.Context, log logr.Logger, rgsIDs []types.UUID, c *rg.ClientWithResponses) error {
	if len(rgsIDs) < 1 {
		log.Info("no runtime groups to clean up")
		return nil
	}

	var errs []error
	for _, rgID := range rgsIDs {
		log.Info("deleting runtime group", "name", rgID)
		if _, err := c.DeleteRuntimeGroupWithResponse(ctx, rgID); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete runtime group %s: %w", rgID, err))
		}
	}
	return errors.Join(errs...)
}

// findOrphanedRolesToDelete gets a list of roles that belong to the orphaned runtime groups.
func findOrphanedRolesToDelete(ctx context.Context, log logr.Logger, orphanedRGsIDs []types.UUID, rolesClient *roles.Client) ([]string, error) {
	if len(orphanedRGsIDs) < 1 {
		log.Info("no runtime groups to clean up, skipping listing roles")
		return nil, nil
	}

	existingRoles, err := rolesClient.ListRuntimeGroupsRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list runtime group roles: %w", err)
	}

	var rolesIDsToDelete []string
	for _, role := range existingRoles {
		belongsToOrphanedRuntimeGroup := lo.ContainsBy(orphanedRGsIDs, func(rgID types.UUID) bool {
			return rgID.String() == role.EntityID
		})
		if !belongsToOrphanedRuntimeGroup {
			log.Info("role is not assigned to an orphaned runtime group, skipping", "id", role.ID)
			continue
		}
		rolesIDsToDelete = append(rolesIDsToDelete, role.ID)
	}
	return rolesIDsToDelete, nil
}

// deleteRoles deletes roles by their IDs.
func deleteRoles(ctx context.Context, log logr.Logger, rolesIDsToDelete []string, rolesClient *roles.Client) error {
	if len(rolesIDsToDelete) == 0 {
		log.Info("no roles to delete")
		return nil
	}

	var errs []error
	for _, roleID := range rolesIDsToDelete {
		log.Info("deleting role", "id", roleID)
		if err := rolesClient.DeleteRole(ctx, roleID); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete role %s: %w", roleID, err))
		}
	}

	return errors.Join(errs...)
}
