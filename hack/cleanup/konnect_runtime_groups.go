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

	cp "github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/controlplanes"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/roles"
)

const (
	konnectControlPlanesBaseURL     = "https://us.kic.api.konghq.tech/v2"
	konnectControlPlanesLimit       = 100
	konnectRolesBaseURL             = "https://global.api.konghq.tech/v2"
	createdInTestsControlPlaneLabel = "created_in_tests"
	timeUntilControlPlaneOrphaned   = time.Hour
)

// cleanupKonnectControlPlanes deletes orphaned control planes created by the tests and their roles.
func cleanupKonnectControlPlanes(ctx context.Context, log logr.Logger) error {
	cpClient, err := cp.NewClientWithResponses(konnectControlPlanesBaseURL, cp.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+konnectAccessToken)
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to create control planes client: %w", err)
	}

	orphanedCPs, err := findOrphanedControlPlanes(ctx, log, cpClient)
	if err != nil {
		return fmt.Errorf("failed to find orphaned control planes: %w", err)
	}
	if err := deleteControlPlanes(ctx, log, orphanedCPs, cpClient); err != nil {
		return fmt.Errorf("failed to delete control planes: %w", err)
	}

	// We have to manually delete roles created for the control plane because Konnect doesn't do it automatically.
	// If we don't do it, we will eventually hit a problem with Konnect APIs answering our requests with 504s
	// because of a performance issue when there's too many roles for the account
	// (see https://konghq.atlassian.net/browse/TPS-1319).
	//
	// We can drop this once the automated cleanup is implemented on Konnect side:
	// https://konghq.atlassian.net/browse/TPS-1453.
	rolesClient := roles.NewClient(&http.Client{}, konnectRolesBaseURL, konnectAccessToken)
	rolesToDelete, err := findOrphanedRolesToDelete(ctx, log, orphanedCPs, rolesClient)
	if err != nil {
		return fmt.Errorf("failed to list control plane roles to delete: %w", err)
	}
	if err := deleteRoles(ctx, log, rolesToDelete, rolesClient); err != nil {
		return fmt.Errorf("failed to delete control plane roles: %w", err)
	}

	return nil
}

// findOrphanedControlPlanes finds control planes that were created by the tests and are older than timeUntilControlPlaneOrphaned.
func findOrphanedControlPlanes(ctx context.Context, log logr.Logger, c *cp.ClientWithResponses) ([]types.UUID, error) {
	response, err := c.ListControlPlanesWithResponse(ctx, &cp.ListControlPlanesParams{
		PageSize: lo.ToPtr(konnectControlPlanesLimit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list control planes: %w", err)
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to list control planes, status: %s, body: %s", response.Status(), string(response.Body))
	}
	if response.JSON200 == nil || response.JSON200.Data == nil {
		return nil, fmt.Errorf("no data in the response, status: %s, body: %s", response.Status(), string(response.Body))
	}

	var orphanedControlPlanes []types.UUID
	for _, ControlPlane := range *response.JSON200.Data {
		if ControlPlane.Labels == nil || (*ControlPlane.Labels)[createdInTestsControlPlaneLabel] != "true" {
			log.Info("Control plane was not created by the tests, skipping", "name", *ControlPlane.Name)
			continue
		}
		if ControlPlane.CreatedAt == nil {
			log.Info("Control plane has no creation timestamp, skipping", "name", *ControlPlane.Name)
			continue
		}
		orphanedAfter := ControlPlane.CreatedAt.Add(timeUntilControlPlaneOrphaned)
		if !time.Now().After(orphanedAfter) {
			log.Info("Control plane is not old enough to be considered orphaned, skipping", "name", *ControlPlane.Name, "created_at", *ControlPlane.CreatedAt)
			continue
		}
		orphanedControlPlanes = append(orphanedControlPlanes, *ControlPlane.Id)
	}
	return orphanedControlPlanes, nil
}

// deleteControlPlanes deletes control planes by their IDs.
func deleteControlPlanes(ctx context.Context, log logr.Logger, cpsIDs []types.UUID, c *cp.ClientWithResponses) error {
	if len(cpsIDs) < 1 {
		log.Info("No control planes to clean up")
		return nil
	}

	var errs []error
	for _, cpID := range cpsIDs {
		log.Info("Deleting control plane", "name", cpID)
		if _, err := c.DeleteControlPlaneWithResponse(ctx, cpID); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete control plane %s: %w", cpID, err))
		}
	}
	return errors.Join(errs...)
}

// findOrphanedRolesToDelete gets a list of roles that belong to the orphaned control planes.
func findOrphanedRolesToDelete(ctx context.Context, log logr.Logger, orphanedCPsIDs []types.UUID, rolesClient *roles.Client) ([]string, error) {
	if len(orphanedCPsIDs) < 1 {
		log.Info("No control planes to clean up, skipping listing roles")
		return nil, nil
	}

	existingRoles, err := rolesClient.ListControlPlanesRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list control plane roles: %w", err)
	}

	var rolesIDsToDelete []string
	for _, role := range existingRoles {
		belongsToOrphanedControlPlane := lo.ContainsBy(orphanedCPsIDs, func(cpID types.UUID) bool {
			return cpID.String() == role.EntityID
		})
		if !belongsToOrphanedControlPlane {
			log.Info("Role is not assigned to an orphaned control plane, skipping", "id", role.ID)
			continue
		}
		rolesIDsToDelete = append(rolesIDsToDelete, role.ID)
	}
	return rolesIDsToDelete, nil
}

// deleteRoles deletes roles by their IDs.
func deleteRoles(ctx context.Context, log logr.Logger, rolesIDsToDelete []string, rolesClient *roles.Client) error {
	if len(rolesIDsToDelete) == 0 {
		log.Info("No roles to delete")
		return nil
	}

	var errs []error
	for _, roleID := range rolesIDsToDelete {
		log.Info("Deleting role", "id", roleID)
		if err := rolesClient.DeleteRole(ctx, roleID); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete role %s: %w", roleID, err))
		}
	}

	return errors.Join(errs...)
}
