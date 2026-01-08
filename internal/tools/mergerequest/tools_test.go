package mergerequest

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

func TestListMergeRequestsTool(t *testing.T) {
	t.Run("returns MR list successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id":            1,
					"iid":           1,
					"title":         "Test MR 1",
					"state":         "opened",
					"source_branch": "feature-1",
					"target_branch": "main",
					"web_url":       "https://gitlab.example.com/project/-/merge_requests/1",
				},
				{
					"id":            2,
					"iid":           2,
					"title":         "Test MR 2",
					"state":         "merged",
					"source_branch": "feature-2",
					"target_branch": "main",
					"web_url":       "https://gitlab.example.com/project/-/merge_requests/2",
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		// Verify tool is registered
		assert.True(t, reg.IsRegistered("list_merge_requests"))

		// Test handler directly
		input := ListMergeRequestsInput{
			ProjectID: "test-project",
		}

		ctx := context.Background()
		_, output, err := listMergeRequestsHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.MergeRequests, 2)
		assert.Equal(t, int64(1), output.MergeRequests[0].IID)
		assert.Equal(t, "Test MR 1", output.MergeRequests[0].Title)
		assert.Equal(t, "opened", output.MergeRequests[0].State)
	})

	t.Run("filters by state", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "opened", r.URL.Query().Get("state"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id":    1,
					"iid":   1,
					"title": "Open MR",
					"state": "opened",
				},
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		state := "opened"
		input := ListMergeRequestsInput{
			ProjectID: "test-project",
			State:     &state,
		}

		ctx := context.Background()
		_, output, err := listMergeRequestsHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.MergeRequests, 1)
	})
}

func TestGetMergeRequestTool(t *testing.T) {
	t.Run("returns MR details successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":            1,
				"iid":           1,
				"title":         "Test MR",
				"description":   "Test description",
				"state":         "opened",
				"source_branch": "feature",
				"target_branch": "main",
				"web_url":       "https://gitlab.example.com/project/-/merge_requests/1",
				"author": map[string]any{
					"id":       10,
					"username": "testuser",
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("get_merge_request"))

		input := GetMergeRequestInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, output, err := getMergeRequestHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(1), output.IID)
		assert.Equal(t, "Test MR", output.Title)
		assert.Equal(t, "Test description", output.Description)
	})

	t.Run("returns error for non-existent MR", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := GetMergeRequestInput{
			ProjectID:       "test-project",
			MergeRequestIID: 999,
		}

		ctx := context.Background()
		_, _, err := getMergeRequestHandler(client, ctx, nil, input)

		assert.Error(t, err)
	})
}

func TestCreateMergeRequestTool(t *testing.T) {
	t.Run("creates MR successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && r.URL.Path == "/api/v4/projects/test-project/merge_requests" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]any{
					"id":            1,
					"iid":           1,
					"title":         "New Feature",
					"source_branch": "feature",
					"target_branch": "main",
					"web_url":       "https://gitlab.example.com/project/-/merge_requests/1",
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("create_merge_request"))

		desc := "Feature description"
		input := CreateMergeRequestInput{
			ProjectID:    "test-project",
			SourceBranch: "feature",
			TargetBranch: "main",
			Title:        "New Feature",
			Description:  &desc,
		}

		ctx := context.Background()
		_, output, err := createMergeRequestHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(1), output.IID)
		assert.Equal(t, "New Feature", output.Title)
		assert.Contains(t, output.WebURL, "merge_requests/1")
	})
}

func TestUpdateMergeRequestTool(t *testing.T) {
	t.Run("updates MR successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PUT" && r.URL.Path == "/api/v4/projects/test-project/merge_requests/1" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{
					"id":      1,
					"iid":     1,
					"title":   "Updated Title",
					"web_url": "https://gitlab.example.com/project/-/merge_requests/1",
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("update_merge_request"))

		newTitle := "Updated Title"
		input := UpdateMergeRequestInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			Title:           &newTitle,
		}

		ctx := context.Background()
		_, output, err := updateMergeRequestHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, "Updated Title", output.Title)
	})
}

func TestMergeMergeRequestTool(t *testing.T) {
	t.Run("merges MR successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "PUT" && r.URL.Path == "/api/v4/projects/test-project/merge_requests/1/merge" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{
					"id":      1,
					"iid":     1,
					"state":   "merged",
					"web_url": "https://gitlab.example.com/project/-/merge_requests/1",
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("merge_merge_request"))

		squash := true
		input := MergeMergeRequestInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			Squash:          &squash,
		}

		ctx := context.Background()
		_, output, err := mergeMergeRequestHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, "merged", output.State)
	})

	t.Run("returns error when merge is not allowed", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"message": "Method Not Allowed"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := MergeMergeRequestInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, _, err := mergeMergeRequestHandler(client, ctx, nil, input)

		assert.Error(t, err)
	})
}

func TestGetMergeRequestChangesTool(t *testing.T) {
	t.Run("returns changes successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v4/projects/test-project/merge_requests/1/diffs" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]map[string]any{
					{
						"old_path": "file.go",
						"new_path": "file.go",
						"diff":     "@@ -1,3 +1,4 @@\n+new line",
						"new_file": false,
					},
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("get_merge_request_changes"))

		input := GetMergeRequestChangesInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, output, err := getMergeRequestChangesHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.Changes, 1)
		assert.Equal(t, "file.go", output.Changes[0].NewPath)
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
			DisabledTools: []string{"list_merge_requests"},
		}

		client, err := gitlab.NewClient(server.URL, "test-token")
		require.NoError(t, err)

		reg := registry.New(cfg)
		Register(reg, client)

		// Tool is registered internally for tracking
		assert.True(t, reg.IsRegistered("list_merge_requests"))
		// But it's not enabled
		assert.False(t, reg.IsToolEnabled("list_merge_requests"))
	})
}
