package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kqns91/gitlab-mcp/internal/config"
	"github.com/kqns91/gitlab-mcp/internal/gitlab"
	"github.com/kqns91/gitlab-mcp/internal/registry"
	"github.com/kqns91/gitlab-mcp/internal/tools/approval"
	"github.com/kqns91/gitlab-mcp/internal/tools/discussion"
	"github.com/kqns91/gitlab-mcp/internal/tools/mergerequest"
	"github.com/kqns91/gitlab-mcp/internal/tools/pipeline"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupIntegrationTest は統合テスト用のサーバーとクライアントセッションをセットアップする
func setupIntegrationTest(t *testing.T, cfg *config.Config, gitlabHandler http.HandlerFunc) (*mcp.ClientSession, func()) {
	// Create mock GitLab server
	gitlabServer := httptest.NewServer(gitlabHandler)

	// Override GitLab URL to use mock server
	cfg.GitLabURL = gitlabServer.URL

	// Create GitLab client
	gitlabClient, err := gitlab.NewClient(cfg.GitLabURL, cfg.GitLabToken)
	require.NoError(t, err)

	// Create registry and register all tools
	reg := registry.New(cfg)
	mergerequest.Register(reg, gitlabClient)
	discussion.Register(reg, gitlabClient)
	approval.Register(reg, gitlabClient)
	pipeline.Register(reg, gitlabClient)

	// Create in-memory transports for testing
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	// Start MCP server in background
	mcpServer := reg.Server()
	serverCtx, serverCancel := context.WithCancel(context.Background())
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- mcpServer.Run(serverCtx, serverTransport)
	}()

	// Connect MCP client
	mcpClient := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := mcpClient.Connect(ctx, clientTransport, nil)
	require.NoError(t, err)

	cleanup := func() {
		session.Close()
		serverCancel()
		<-serverDone
		gitlabServer.Close()
	}

	return session, cleanup
}

func TestIntegration_ListTools(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// List available tools
	tools, err := session.ListTools(ctx, nil)
	require.NoError(t, err)

	// Should have all 15 tools registered
	toolNames := make([]string, len(tools.Tools))
	for i, tool := range tools.Tools {
		toolNames[i] = tool.Name
	}

	expectedTools := []string{
		"list_merge_requests",
		"get_merge_request",
		"create_merge_request",
		"update_merge_request",
		"merge_merge_request",
		"get_merge_request_changes",
		"add_merge_request_comment",
		"add_merge_request_discussion",
		"list_merge_request_discussions",
		"resolve_discussion",
		"approve_merge_request",
		"unapprove_merge_request",
		"get_merge_request_approvals",
		"list_merge_request_pipelines",
		"get_pipeline_jobs",
	}

	for _, expected := range expectedTools {
		assert.Contains(t, toolNames, expected, "Tool %s should be registered", expected)
	}
}

func TestIntegration_DisabledToolsExcludedFromList(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:     "https://gitlab.example.com",
		GitLabToken:   "test-token",
		DisabledTools: []string{"merge_merge_request", "approve_merge_request"},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// List available tools
	tools, err := session.ListTools(ctx, nil)
	require.NoError(t, err)

	toolNames := make([]string, len(tools.Tools))
	for i, tool := range tools.Tools {
		toolNames[i] = tool.Name
	}

	// Disabled tools should not be in the list
	assert.NotContains(t, toolNames, "merge_merge_request")
	assert.NotContains(t, toolNames, "approve_merge_request")

	// Other tools should still be available
	assert.Contains(t, toolNames, "list_merge_requests")
	assert.Contains(t, toolNames, "get_merge_request")
}

func TestIntegration_CallTool_ListMergeRequests(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/projects/test-project/merge_requests" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"iid":           1,
					"title":         "Test MR 1",
					"state":         "opened",
					"source_branch": "feature-1",
					"target_branch": "main",
					"web_url":       "https://gitlab.example.com/project/-/merge_requests/1",
					"author":        map[string]any{"id": 10, "username": "user1"},
				},
				{
					"iid":           2,
					"title":         "Test MR 2",
					"state":         "merged",
					"source_branch": "feature-2",
					"target_branch": "main",
					"web_url":       "https://gitlab.example.com/project/-/merge_requests/2",
					"author":        map[string]any{"id": 11, "username": "user2"},
				},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call list_merge_requests tool
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "list_merge_requests",
		Arguments: map[string]any{
			"project_id": "test-project",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.IsError)

	// Verify response contains MR data
	require.Len(t, result.Content, 1)
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, textContent.Text, "Test MR 1")
	assert.Contains(t, textContent.Text, "Test MR 2")
}

func TestIntegration_CallTool_GetMergeRequest(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/projects/test-project/merge_requests/1" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"iid":           1,
				"title":         "Test MR",
				"description":   "Test description",
				"state":         "opened",
				"source_branch": "feature",
				"target_branch": "main",
				"web_url":       "https://gitlab.example.com/project/-/merge_requests/1",
				"author":        map[string]any{"id": 10, "username": "testuser"},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call get_merge_request tool
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_merge_request",
		Arguments: map[string]any{
			"project_id":        "test-project",
			"merge_request_iid": 1,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.IsError)

	// Verify response
	require.Len(t, result.Content, 1)
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, textContent.Text, "Test MR")
	assert.Contains(t, textContent.Text, "Test description")
}

func TestIntegration_CallTool_NotFound(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call tool for non-existent MR
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_merge_request",
		Arguments: map[string]any{
			"project_id":        "test-project",
			"merge_request_iid": 999,
		},
	})

	// Should return error in result
	require.NoError(t, err) // No transport error
	require.NotNil(t, result)
	assert.True(t, result.IsError, "Result should indicate error")
}

func TestIntegration_CallTool_Unauthorized(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "invalid-token",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "401 Unauthorized"})
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call tool with invalid token
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "list_merge_requests",
		Arguments: map[string]any{
			"project_id": "test-project",
		},
	})

	require.NoError(t, err) // No transport error
	require.NotNil(t, result)
	assert.True(t, result.IsError, "Result should indicate authentication error")
}

func TestIntegration_CallTool_RateLimited(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{"message": "429 Too Many Requests"})
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call tool when rate limited
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "list_merge_requests",
		Arguments: map[string]any{
			"project_id": "test-project",
		},
	})

	require.NoError(t, err) // No transport error
	require.NotNil(t, result)
	assert.True(t, result.IsError, "Result should indicate rate limit error")
}

func TestIntegration_ApprovalWorkflow(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v4/projects/test-project/merge_requests/1/approvals" && r.Method == "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"approved":           false,
				"approvals_required": 2,
				"approvals_left":     2,
				"user_has_approved":  false,
				"user_can_approve":   true,
				"approved_by":        []any{},
			})
		case r.URL.Path == "/api/v4/projects/test-project/merge_requests/1/approve" && r.Method == "POST":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"approved":          false,
				"approvals_left":    1,
				"user_has_approved": true,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get approvals status
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_merge_request_approvals",
		Arguments: map[string]any{
			"project_id":        "test-project",
			"merge_request_iid": 1,
		},
	})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	// Approve MR
	result, err = session.CallTool(ctx, &mcp.CallToolParams{
		Name: "approve_merge_request",
		Arguments: map[string]any{
			"project_id":        "test-project",
			"merge_request_iid": 1,
		},
	})
	require.NoError(t, err)
	assert.False(t, result.IsError)
}

func TestIntegration_PipelineWorkflow(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v4/projects/test-project/merge_requests/1/pipelines" && r.Method == "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id":         100,
					"status":     "success",
					"ref":        "feature",
					"sha":        "abc123",
					"web_url":    "https://gitlab.example.com/project/-/pipelines/100",
					"created_at": "2024-01-01T10:00:00Z",
				},
			})
		case r.URL.Path == "/api/v4/projects/test-project/pipelines/100/jobs" && r.Method == "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{"id": 1, "name": "build", "stage": "build", "status": "success"},
				{"id": 2, "name": "test", "stage": "test", "status": "success"},
				{"id": 3, "name": "deploy", "stage": "deploy", "status": "success"},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// List pipelines
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "list_merge_request_pipelines",
		Arguments: map[string]any{
			"project_id":        "test-project",
			"merge_request_iid": 1,
		},
	})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	// Get pipeline jobs
	result, err = session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_pipeline_jobs",
		Arguments: map[string]any{
			"project_id":  "test-project",
			"pipeline_id": 100,
		},
	})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	// Verify response contains job data
	require.Len(t, result.Content, 1)
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, textContent.Text, "build")
	assert.Contains(t, textContent.Text, "test")
	assert.Contains(t, textContent.Text, "deploy")
}

func TestIntegration_EnabledToolsOnly(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:    "https://gitlab.example.com",
		GitLabToken:  "test-token",
		EnabledTools: []string{"list_merge_requests", "get_merge_request"},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	session, cleanup := setupIntegrationTest(t, cfg, handler)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// List available tools
	tools, err := session.ListTools(ctx, nil)
	require.NoError(t, err)

	toolNames := make([]string, len(tools.Tools))
	for i, tool := range tools.Tools {
		toolNames[i] = tool.Name
	}

	// Only enabled tools should be in the list
	assert.Len(t, toolNames, 2)
	assert.Contains(t, toolNames, "list_merge_requests")
	assert.Contains(t, toolNames, "get_merge_request")
	assert.NotContains(t, toolNames, "create_merge_request")
	assert.NotContains(t, toolNames, "approve_merge_request")
}
