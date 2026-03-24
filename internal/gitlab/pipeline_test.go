package gitlab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListMergeRequestPipelines_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/pipelines", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":      100,
				"status":  "success",
				"ref":     "feature-branch",
				"sha":     "abc123",
				"web_url": "https://gitlab.example.com/project/-/pipelines/100",
			},
			{
				"id":      99,
				"status":  "failed",
				"ref":     "feature-branch",
				"sha":     "def456",
				"web_url": "https://gitlab.example.com/project/-/pipelines/99",
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	pipelines, err := client.ListMergeRequestPipelines("test-project", 1)

	require.NoError(t, err)
	assert.Len(t, pipelines, 2)
	assert.Equal(t, int64(100), pipelines[0].ID)
	assert.Equal(t, "success", pipelines[0].Status)
	assert.Equal(t, int64(99), pipelines[1].ID)
	assert.Equal(t, "failed", pipelines[1].Status)
}

func TestListMergeRequestPipelines_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	pipelines, err := client.ListMergeRequestPipelines("unknown-project", 999)

	assert.Nil(t, pipelines)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}

func TestListPipelineJobs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	jobs, err := client.ListPipelineJobs("test-project", 100, nil)

	require.NoError(t, err)
	assert.Len(t, jobs, 3)
	assert.Equal(t, "build", jobs[0].Name)
	assert.Equal(t, "build", jobs[0].Stage)
	assert.Equal(t, "success", jobs[0].Status)
	assert.Equal(t, "test", jobs[1].Name)
	assert.Equal(t, "deploy", jobs[2].Name)
	assert.Equal(t, "running", jobs[2].Status)
}

func TestListPipelineJobs_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Pipeline not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	jobs, err := client.ListPipelineJobs("test-project", 999, nil)

	assert.Nil(t, jobs)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}

func TestListProjectPipelines_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	pipelines, err := client.ListProjectPipelines("test-project", nil)

	require.NoError(t, err)
	assert.Len(t, pipelines, 1)
	assert.Equal(t, int64(200), pipelines[0].ID)
	assert.Equal(t, "success", pipelines[0].Status)
}

func TestGetPipeline_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	pipeline, err := client.GetPipeline("test-project", 200)

	require.NoError(t, err)
	assert.Equal(t, int64(200), pipeline.ID)
	assert.Equal(t, "success", pipeline.Status)
	assert.Equal(t, "main", pipeline.Ref)
	assert.Equal(t, int64(120), pipeline.Duration)
}

func TestGetPipeline_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	pipeline, err := client.GetPipeline("test-project", 999)

	assert.Nil(t, pipeline)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}

func TestCreatePipeline_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/pipeline", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      300,
			"status":  "pending",
			"ref":     "main",
			"web_url": "https://gitlab.example.com/project/-/pipelines/300",
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	pipeline, err := client.CreatePipeline("test-project", &CreatePipelineOptions{
		Ref: "main",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(300), pipeline.ID)
	assert.Equal(t, "pending", pipeline.Status)
}

func TestRetryPipeline_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/pipelines/200/retry", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      200,
			"status":  "pending",
			"web_url": "https://gitlab.example.com/project/-/pipelines/200",
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	pipeline, err := client.RetryPipeline("test-project", 200)

	require.NoError(t, err)
	assert.Equal(t, int64(200), pipeline.ID)
	assert.Equal(t, "pending", pipeline.Status)
}

func TestCancelPipeline_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/pipelines/200/cancel", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      200,
			"status":  "canceled",
			"web_url": "https://gitlab.example.com/project/-/pipelines/200",
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	pipeline, err := client.CancelPipeline("test-project", 200)

	require.NoError(t, err)
	assert.Equal(t, int64(200), pipeline.ID)
	assert.Equal(t, "canceled", pipeline.Status)
}

func TestGetJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	job, err := client.GetJob("test-project", 10)

	require.NoError(t, err)
	assert.Equal(t, int64(10), job.ID)
	assert.Equal(t, "build", job.Name)
	assert.Equal(t, "success", job.Status)
}

func TestGetJobTrace_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/jobs/10/trace", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Running build...\nBuild succeeded!"))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	trace, err := client.GetJobTrace("test-project", 10)

	require.NoError(t, err)
	assert.Contains(t, trace, "Running build...")
	assert.Contains(t, trace, "Build succeeded!")
}

func TestRetryJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/jobs/10/retry", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      11,
			"name":    "build",
			"status":  "pending",
			"web_url": "https://gitlab.example.com/project/-/jobs/11",
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	job, err := client.RetryJob("test-project", 10)

	require.NoError(t, err)
	assert.Equal(t, int64(11), job.ID)
	assert.Equal(t, "pending", job.Status)
}
