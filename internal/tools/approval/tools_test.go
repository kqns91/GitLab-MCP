package approval

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kqns91/gitlab-mcp/internal/config"
	"github.com/kqns91/gitlab-mcp/internal/gitlab"
	"github.com/kqns91/gitlab-mcp/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*gitlab.Client, *registry.Registry, func()) {
	server := httptest.NewServer(handler)

	cfg := &config.Config{
		GitLabURL:   server.URL,
		GitLabToken: "test-token",
	}

	client, err := gitlab.NewClient(server.URL, "test-token")
	require.NoError(t, err)

	reg := registry.New(cfg)
	Register(reg, client)

	return client, reg, server.Close
}

func TestApproveMergeRequestTool(t *testing.T) {
	t.Run("approves MR successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/approve", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":                1,
				"iid":               1,
				"approved":          true,
				"approvals_left":    0,
				"user_has_approved": true,
				"approved_by": []map[string]any{
					{
						"user": map[string]any{
							"id":       10,
							"username": "approver",
						},
					},
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("approve_merge_request"))

		input := ApproveInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, output, err := approveHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.True(t, output.Approved)
		assert.True(t, output.UserHasApproved)
		assert.Equal(t, int64(0), output.ApprovalsLeft)
	})

	t.Run("returns error when approval forbidden", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "You cannot approve this merge request"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := ApproveInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, _, err := approveHandler(client, ctx, nil, input)

		assert.Error(t, err)
		mcpErr, ok := err.(*gitlab.MCPError)
		require.True(t, ok)
		assert.Equal(t, gitlab.ErrCodeForbidden, mcpErr.Code)
	})
}

func TestUnapproveMergeRequestTool(t *testing.T) {
	t.Run("unapproves MR successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/unapprove", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.WriteHeader(http.StatusOK)
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("unapprove_merge_request"))

		input := UnapproveInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, output, err := unapproveHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.True(t, output.Success)
	})

	t.Run("returns error when unapproval forbidden", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "You cannot unapprove this merge request"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := UnapproveInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, _, err := unapproveHandler(client, ctx, nil, input)

		assert.Error(t, err)
	})
}

func TestGetMergeRequestApprovalsTool(t *testing.T) {
	t.Run("returns approvals successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/approvals", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":                 1,
				"iid":                1,
				"approved":           true,
				"approvals_required": 2,
				"approvals_left":     0,
				"user_has_approved":  false,
				"user_can_approve":   true,
				"approved_by": []map[string]any{
					{
						"user": map[string]any{
							"id":       10,
							"username": "approver1",
						},
					},
					{
						"user": map[string]any{
							"id":       11,
							"username": "approver2",
						},
					},
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("get_merge_request_approvals"))

		input := GetApprovalsInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, output, err := getApprovalsHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.True(t, output.Approved)
		assert.Equal(t, int64(2), output.ApprovalsRequired)
		assert.Equal(t, int64(0), output.ApprovalsLeft)
		assert.Len(t, output.ApprovedBy, 2)
		assert.Equal(t, "approver1", output.ApprovedBy[0].Username)
		assert.Equal(t, "approver2", output.ApprovedBy[1].Username)
	})

	t.Run("returns error for non-existent MR", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := GetApprovalsInput{
			ProjectID:       "test-project",
			MergeRequestIID: 999,
		}

		ctx := context.Background()
		_, _, err := getApprovalsHandler(client, ctx, nil, input)

		assert.Error(t, err)
	})
}

func TestToolDisabled(t *testing.T) {
	t.Run("disabled tool is not registered in server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		cfg := &config.Config{
			GitLabURL:     server.URL,
			GitLabToken:   "test-token",
			DisabledTools: []string{"approve_merge_request"},
		}

		client, err := gitlab.NewClient(server.URL, "test-token")
		require.NoError(t, err)

		reg := registry.New(cfg)
		Register(reg, client)

		assert.True(t, reg.IsRegistered("approve_merge_request"))
		assert.False(t, reg.IsToolEnabled("approve_merge_request"))
	})
}
