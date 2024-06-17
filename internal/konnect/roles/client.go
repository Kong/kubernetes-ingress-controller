package roles

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/useragent"
)

const (
	konnectUsersMeURL            = "/users/me"
	konnectUsersAssignedRolesURL = "/users/%s/assigned-roles?filter%%5Bentity_type_name%%5D=Runtime+Groups"
	konnectAssignedRoleURL       = "/users/%s/assigned-roles/%s"
)

type Client struct {
	httpClient          *http.Client
	personalAccessToken string
	currentUserID       string
	baseURL             string
}

type Role struct {
	// ID is the role ID.
	ID string

	// EntityID is the ID of the entity the role is assigned to (e.g. Control Plane).
	EntityID string
}

func NewRole(id, entityID string) (Role, error) {
	if id == "" {
		return Role{}, fmt.Errorf("role ID is required")
	}
	if entityID == "" {
		return Role{}, fmt.Errorf("entity ID is required")
	}
	return Role{
		ID:       id,
		EntityID: entityID,
	}, nil
}

func NewClient(httpClient *http.Client, baseURL string, personalAccessToken string) *Client {
	httpClient.Transport = useragent.NewTransport(httpClient.Transport)
	return &Client{
		baseURL:             baseURL,
		httpClient:          httpClient,
		personalAccessToken: personalAccessToken,
	}
}

// ListControlPlanesRoles lists all roles assigned to the current user for Control Planes.
func (c *Client) ListControlPlanesRoles(ctx context.Context) ([]Role, error) {
	currentUserID, err := c.getCurrentUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user ID: %w", err)
	}

	listRolesURL := fmt.Sprintf(konnectUsersAssignedRolesURL, currentUserID)
	req, err := c.newRequestWithAuth(ctx, http.MethodGet, listRolesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list roles, status: %s", resp.Status)
	}

	var rolesResponse struct {
		Data []struct {
			ID       string `json:"id"`
			EntityID string `json:"entity_id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rolesResponse); err != nil {
		return nil, fmt.Errorf("failed to decode roles response: %w", err)
	}

	roles := make([]Role, 0, len(rolesResponse.Data))
	for _, role := range rolesResponse.Data {
		r, err := NewRole(role.ID, role.EntityID)
		if err != nil {
			return nil, fmt.Errorf("failed to create role: %w", err)
		}
		roles = append(roles, r)
	}

	return roles, nil
}

// DeleteRole deletes a role assigned to the current user.
func (c *Client) DeleteRole(ctx context.Context, roleID string) error {
	currentUserID, err := c.getCurrentUserID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current user ID: %w", err)
	}

	deleteRoleURL := fmt.Sprintf(konnectAssignedRoleURL, currentUserID, roleID)
	req, err := c.newRequestWithAuth(ctx, http.MethodDelete, deleteRoleURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete role %s: %w", roleID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete role %s, status: %s", roleID, resp.Status)
	}

	return nil
}

func (c *Client) getCurrentUserID(ctx context.Context) (string, error) {
	// It's already cached, no need to make a request.
	if c.currentUserID != "" {
		return c.currentUserID, nil
	}

	meRequest, err := c.newRequestWithAuth(ctx, http.MethodGet, konnectUsersMeURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	meResponse, err := c.httpClient.Do(meRequest)
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	defer meResponse.Body.Close()
	if meResponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get current user, status: %s", meResponse.Status)
	}

	var meResponseData struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(meResponse.Body).Decode(&meResponseData); err != nil {
		return "", fmt.Errorf("failed to decode current user response: %w", err)
	}

	if meResponseData.ID == "" {
		return "", fmt.Errorf("failed to get current user, empty id")
	}

	c.currentUserID = meResponseData.ID
	return meResponseData.ID, nil
}

func (c *Client) newRequestWithAuth(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.personalAccessToken)
	return req, nil
}
