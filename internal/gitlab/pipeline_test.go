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
