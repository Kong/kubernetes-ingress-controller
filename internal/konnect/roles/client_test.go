package roles_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/roles"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
)

const (
	currentUserID = "515d12f2-aab2-42b3-a093-6ad793c0c7ab"
	testToken     = "test-token"
)

func newMockRolesServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	requireToken := func(r *http.Request) {
		token := r.Header.Get("Authorization")
		require.Equal(t, "Bearer "+testToken, token)
	}
	requireUserAgent := func(r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		require.Equal(t, metadata.UserAgent(), userAgent)
	}

	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		requireUserAgent(r)
		requireToken(r)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
		{
		   "id": "515d12f2-aab2-42b3-a093-6ad793c0c7ab",
		   "email": "team-k8s+konnect-testing-2@konghq.com",
		   "full_name": "Kubernetes Team",
		   "preferred_name": "",
		   "active": true,
		   "created_at": "2023-06-26T12:16:08Z",
		   "updated_at": "2023-06-26T12:16:28Z"
		}`))
	})

	mux.HandleFunc("/users/"+currentUserID+"/assigned-roles/", func(w http.ResponseWriter, r *http.Request) {
		requireUserAgent(r)
		requireToken(r)

		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`
		{
		   "meta": {
		       "page": {
		           "number": 1,
		           "size": 10,
		           "total": 2
		       }
		   },
		   "data": [
		       {
		           "id": "24ac168d-4ffb-46ec-8dd6-5a26b5ec6f0b",
		           "role_name": "Admin",
		           "entity_region": "us",
		           "entity_type_name": "Control Planes",
		           "entity_id": "e3f155ec-1786-4017-98d0-b0a0f5e179c3"
		       },
		       {
		           "id": "7edaf68b-8f07-4827-b540-fce06e45429e",
		           "role_name": "Admin",
		           "entity_region": "us",
		           "entity_type_name": "Control Planes",
		           "entity_id": "c486f518-9fc8-461f-af0a-2bc85b70e492"
		       }
		   ]
		}`))
		case http.MethodDelete:
			requireUserAgent(r)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	return httptest.NewServer(mux)
}

func TestRolesClient(t *testing.T) {
	ctx := context.Background()
	server := newMockRolesServer(t)
	c := roles.NewClient(&http.Client{}, server.URL, testToken)

	rgRoles, err := c.ListControlPlanesRoles(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(rgRoles))
	role := rgRoles[0]
	require.Equal(t, "24ac168d-4ffb-46ec-8dd6-5a26b5ec6f0b", role.ID)
	require.Equal(t, "e3f155ec-1786-4017-98d0-b0a0f5e179c3", role.EntityID)

	err = c.DeleteRole(ctx, "24ac168d-4ffb-46ec-8dd6-5a26b5ec6f0b")
	require.NoError(t, err)
}
