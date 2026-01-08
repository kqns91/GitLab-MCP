package gitlab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListMergeRequests_Success(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"iid": 1, "title": "MR 1", "state": "opened"},
			{"iid": 2, "title": "MR 2", "state": "merged"},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	mrs, err := client.ListMergeRequests("test-project", nil)

	require.NoError(t, err)
	assert.Len(t, mrs, 2)
	assert.Equal(t, int64(1), mrs[0].IID)
	assert.Equal(t, "MR 1", mrs[0].Title)
}

func TestListMergeRequests_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "opened", r.URL.Query().Get("state"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"iid": 1, "title": "MR 1", "state": "opened"},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	state := "opened"
	mrs, err := client.ListMergeRequests("test-project", &ListMergeRequestsOptions{
		State: &state,
	})

	require.NoError(t, err)
	assert.Len(t, mrs, 1)
}

func TestListMergeRequests_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Project not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	mrs, err := client.ListMergeRequests("unknown-project", nil)

	assert.Nil(t, mrs)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}

func TestGetMergeRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"iid":           1,
			"title":         "Test MR",
			"description":   "Test description",
			"state":         "opened",
			"source_branch": "feature",
			"target_branch": "main",
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	mr, err := client.GetMergeRequest("test-project", 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), mr.IID)
	assert.Equal(t, "Test MR", mr.Title)
	assert.Equal(t, "feature", mr.SourceBranch)
	assert.Equal(t, "main", mr.TargetBranch)
}

func TestGetMergeRequest_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	mr, err := client.GetMergeRequest("test-project", 999)

	assert.Nil(t, mr)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}

func TestCreateMergeRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"iid":           1,
			"title":         "New MR",
			"source_branch": "feature",
			"target_branch": "main",
			"web_url":       "https://gitlab.example.com/project/-/merge_requests/1",
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	mr, err := client.CreateMergeRequest("test-project", &CreateMergeRequestOptions{
		SourceBranch: "feature",
		TargetBranch: "main",
		Title:        "New MR",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), mr.IID)
	assert.Equal(t, "New MR", mr.Title)
	assert.Equal(t, "feature", mr.SourceBranch)
}

func TestUpdateMergeRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"iid":   1,
			"title": "Updated Title",
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	newTitle := "Updated Title"
	mr, err := client.UpdateMergeRequest("test-project", 1, &UpdateMergeRequestOptions{
		Title: &newTitle,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Title", mr.Title)
}

func TestMergeMergeRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/merge", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"iid":   1,
			"state": "merged",
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	squash := true
	removeSourceBranch := true
	mr, err := client.MergeMergeRequest("test-project", 1, &MergeMergeRequestOptions{
		Squash:                   &squash,
		ShouldRemoveSourceBranch: &removeSourceBranch,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), mr.IID)
	assert.Equal(t, "merged", mr.State)
}

func TestMergeMergeRequest_Conflict(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Branch cannot be merged"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	mr, err := client.MergeMergeRequest("test-project", 1, nil)

	assert.Nil(t, mr)
	assert.Error(t, err)
}

func TestGetMergeRequestChanges_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/diffs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{
				"old_path":     "README.md",
				"new_path":     "README.md",
				"diff":         "@@ -1,3 +1,5 @@\n # Project\n+\n+Description",
				"new_file":     false,
				"renamed_file": false,
				"deleted_file": false,
			},
			{
				"old_path":     "",
				"new_path":     "new_file.go",
				"diff":         "@@ -0,0 +1,10 @@\n+package main",
				"new_file":     true,
				"renamed_file": false,
				"deleted_file": false,
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	diffs, err := client.GetMergeRequestChanges("test-project", 1, nil)

	require.NoError(t, err)
	assert.Len(t, diffs, 2)
	assert.Equal(t, "README.md", diffs[0].NewPath)
	assert.False(t, diffs[0].NewFile)
	assert.Equal(t, "new_file.go", diffs[1].NewPath)
	assert.True(t, diffs[1].NewFile)
}
