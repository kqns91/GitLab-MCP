package pipeline

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

func TestListMergeRequestPipelinesTool(t *testing.T) {
	t.Run("returns pipelines list successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/pipelines", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id":         100,
					"status":     "success",
					"ref":        "feature-branch",
					"sha":        "abc123",
					"web_url":    "https://gitlab.example.com/project/-/pipelines/100",
					"created_at": "2024-01-01T10:00:00Z",
				},
				{
					"id":         99,
					"status":     "failed",
					"ref":        "feature-branch",
					"sha":        "def456",
					"web_url":    "https://gitlab.example.com/project/-/pipelines/99",
					"created_at": "2024-01-01T09:00:00Z",
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("list_merge_request_pipelines"))

		input := ListPipelinesInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, output, err := listPipelinesHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.Pipelines, 2)
		assert.Equal(t, int64(100), output.Pipelines[0].ID)
		assert.Equal(t, "success", output.Pipelines[0].Status)
		assert.Equal(t, int64(99), output.Pipelines[1].ID)
		assert.Equal(t, "failed", output.Pipelines[1].Status)
	})

	t.Run("returns error for non-existent MR", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := ListPipelinesInput{
			ProjectID:       "test-project",
			MergeRequestIID: 999,
		}

		ctx := context.Background()
		_, _, err := listPipelinesHandler(client, ctx, nil, input)

		assert.Error(t, err)
	})
}

func TestGetPipelineJobsTool(t *testing.T) {
	t.Run("returns jobs list successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/pipelines/100/jobs", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id":     1,
					"name":   "build",
					"stage":  "build",
					"status": "success",
				},
				{
					"id":     2,
					"name":   "test",
					"stage":  "test",
					"status": "success",
				},
				{
					"id":     3,
					"name":   "deploy",
					"stage":  "deploy",
					"status": "running",
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("get_pipeline_jobs"))

		input := GetJobsInput{
			ProjectID:  "test-project",
			PipelineID: 100,
		}

		ctx := context.Background()
		_, output, err := getJobsHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.Jobs, 3)
		assert.Equal(t, "build", output.Jobs[0].Name)
		assert.Equal(t, "build", output.Jobs[0].Stage)
		assert.Equal(t, "success", output.Jobs[0].Status)
		assert.Equal(t, "test", output.Jobs[1].Name)
		assert.Equal(t, "deploy", output.Jobs[2].Name)
		assert.Equal(t, "running", output.Jobs[2].Status)
	})

	t.Run("returns error for non-existent pipeline", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "404 Pipeline not found"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := GetJobsInput{
			ProjectID:  "test-project",
			PipelineID: 999,
		}

		ctx := context.Background()
		_, _, err := getJobsHandler(client, ctx, nil, input)

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
			DisabledTools: []string{"list_merge_request_pipelines"},
		}

		client, err := gitlab.NewClient(server.URL, "test-token")
		require.NoError(t, err)

		reg := registry.New(cfg)
		Register(reg, client)

		assert.True(t, reg.IsRegistered("list_merge_request_pipelines"))
		assert.False(t, reg.IsToolEnabled("list_merge_request_pipelines"))
	})
}
