package gitlab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApproveMergeRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
						"id":       1,
						"username": "testuser",
					},
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	approval, err := client.ApproveMergeRequest("test-project", 1)

	require.NoError(t, err)
	assert.True(t, approval.Approved)
	assert.True(t, approval.UserHasApproved)
	assert.Equal(t, int64(0), approval.ApprovalsLeft)
}

func TestApproveMergeRequest_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "You cannot approve this merge request"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	approval, err := client.ApproveMergeRequest("test-project", 1)

	assert.Nil(t, approval)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeForbidden, mcpErr.Code)
}

func TestUnapproveMergeRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/unapprove", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	err = client.UnapproveMergeRequest("test-project", 1)

	assert.NoError(t, err)
}

func TestUnapproveMergeRequest_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "You cannot unapprove this merge request"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	err = client.UnapproveMergeRequest("test-project", 1)

	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeForbidden, mcpErr.Code)
}

func TestGetMergeRequestApprovals_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
						"id":       1,
						"username": "approver1",
					},
				},
				{
					"user": map[string]any{
						"id":       2,
						"username": "approver2",
					},
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	approvals, err := client.GetMergeRequestApprovals("test-project", 1)

	require.NoError(t, err)
	assert.True(t, approvals.Approved)
	assert.Equal(t, int64(2), approvals.ApprovalsRequired)
	assert.Equal(t, int64(0), approvals.ApprovalsLeft)
	assert.Len(t, approvals.ApprovedBy, 2)
}

func TestGetMergeRequestApprovals_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	approvals, err := client.GetMergeRequestApprovals("unknown-project", 999)

	assert.Nil(t, approvals)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}
