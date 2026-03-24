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

func TestListProjectPipelinesTool(t *testing.T) {
	t.Run("returns project pipelines successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/pipelines", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id":      200,
					"status":  "success",
					"ref":     "main",
					"sha":     "abc123",
					"web_url": "https://gitlab.example.com/project/-/pipelines/200",
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("list_project_pipelines"))

		input := ListProjectPipelinesInput{
			ProjectID: "test-project",
		}

		ctx := context.Background()
		_, output, err := listProjectPipelinesHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.Pipelines, 1)
		assert.Equal(t, int64(200), output.Pipelines[0].ID)
		assert.Equal(t, "success", output.Pipelines[0].Status)
	})
}

func TestGetPipelineTool(t *testing.T) {
	t.Run("returns pipeline detail successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/pipelines/200", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":       200,
				"status":   "success",
				"ref":      "main",
				"sha":      "abc123",
				"web_url":  "https://gitlab.example.com/project/-/pipelines/200",
				"source":   "push",
				"duration": 120,
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := GetPipelineInput{
			ProjectID:  "test-project",
			PipelineID: 200,
		}

		ctx := context.Background()
		_, output, err := getPipelineHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(200), output.ID)
		assert.Equal(t, "success", output.Status)
		assert.Equal(t, "push", output.Source)
		assert.Equal(t, int64(120), output.Duration)
	})
}

func TestCreatePipelineTool(t *testing.T) {
	t.Run("creates pipeline successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/pipeline", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":      300,
				"status":  "pending",
				"web_url": "https://gitlab.example.com/project/-/pipelines/300",
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := CreatePipelineInput{
			ProjectID: "test-project",
			Ref:       "main",
		}

		ctx := context.Background()
		_, output, err := createPipelineHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(300), output.ID)
		assert.Equal(t, "pending", output.Status)
	})
}

func TestRetryPipelineTool(t *testing.T) {
	t.Run("retries pipeline successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/pipelines/200/retry", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":      200,
				"status":  "pending",
				"web_url": "https://gitlab.example.com/project/-/pipelines/200",
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := RetryPipelineInput{
			ProjectID:  "test-project",
			PipelineID: 200,
		}

		ctx := context.Background()
		_, output, err := retryPipelineHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(200), output.ID)
		assert.Equal(t, "pending", output.Status)
	})
}

func TestCancelPipelineTool(t *testing.T) {
	t.Run("cancels pipeline successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/pipelines/200/cancel", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":      200,
				"status":  "canceled",
				"web_url": "https://gitlab.example.com/project/-/pipelines/200",
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := CancelPipelineInput{
			ProjectID:  "test-project",
			PipelineID: 200,
		}

		ctx := context.Background()
		_, output, err := cancelPipelineHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(200), output.ID)
		assert.Equal(t, "canceled", output.Status)
	})
}

func TestGetPipelineJobTool(t *testing.T) {
	t.Run("returns job detail successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/jobs/10", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":       10,
				"name":     "build",
				"stage":    "build",
				"status":   "success",
				"ref":      "main",
				"web_url":  "https://gitlab.example.com/project/-/jobs/10",
				"duration": 30.5,
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := GetPipelineJobInput{
			ProjectID: "test-project",
			JobID:     10,
		}

		ctx := context.Background()
		_, output, err := getPipelineJobHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(10), output.ID)
		assert.Equal(t, "build", output.Name)
		assert.Equal(t, "success", output.Status)
	})
}

func TestGetJobLogTool(t *testing.T) {
	t.Run("returns job log successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/jobs/10/trace", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("Running build...\nBuild succeeded!"))
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := GetJobLogInput{
			ProjectID: "test-project",
			JobID:     10,
		}

		ctx := context.Background()
		_, output, err := getJobLogHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Contains(t, output.Log, "Running build...")
		assert.Contains(t, output.Log, "Build succeeded!")
	})
}

func TestRetryPipelineJobTool(t *testing.T) {
	t.Run("retries job successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/jobs/10/retry", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":      11,
				"name":    "build",
				"status":  "pending",
				"web_url": "https://gitlab.example.com/project/-/jobs/11",
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := RetryPipelineJobInput{
			ProjectID: "test-project",
			JobID:     10,
		}

		ctx := context.Background()
		_, output, err := retryPipelineJobHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(11), output.ID)
		assert.Equal(t, "pending", output.Status)
	})
}

func TestNewToolRegistration(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	_, reg, cleanup := setupTestServer(t, handler)
	defer cleanup()

	newTools := []string{
		"list_project_pipelines",
		"get_pipeline",
		"create_pipeline",
		"retry_pipeline",
		"cancel_pipeline",
		"get_pipeline_job",
		"get_job_log",
		"retry_pipeline_job",
	}

	for _, tool := range newTools {
		assert.True(t, reg.IsRegistered(tool), "tool %s should be registered", tool)
	}
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
