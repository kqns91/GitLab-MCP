package issue

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

func TestListIssuesTool(t *testing.T) {
	t.Run("returns issues list successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/issues", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id":      101,
					"iid":     1,
					"title":   "Bug report",
					"state":   "opened",
					"web_url": "https://gitlab.example.com/project/-/issues/1",
					"author":  map[string]any{"username": "user1"},
					"labels":  []string{"bug"},
				},
				{
					"id":      102,
					"iid":     2,
					"title":   "Feature request",
					"state":   "opened",
					"web_url": "https://gitlab.example.com/project/-/issues/2",
					"author":  map[string]any{"username": "user2"},
					"labels":  []string{"enhancement"},
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("list_issues"))

		input := ListIssuesInput{
			ProjectID: "test-project",
		}

		ctx := context.Background()
		_, output, err := listIssuesHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.Issues, 2)
		assert.Equal(t, int64(1), output.Issues[0].IID)
		assert.Equal(t, "Bug report", output.Issues[0].Title)
		assert.Equal(t, "opened", output.Issues[0].State)
		assert.Equal(t, "user1", output.Issues[0].AuthorName)
		assert.Equal(t, int64(2), output.Issues[1].IID)
	})

	t.Run("returns error for non-existent project", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := ListIssuesInput{
			ProjectID: "unknown-project",
		}

		ctx := context.Background()
		_, _, err := listIssuesHandler(client, ctx, nil, input)

		assert.Error(t, err)
	})
}

func TestGetIssueTool(t *testing.T) {
	t.Run("returns issue detail successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/issues/1", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":          101,
				"iid":         1,
				"title":       "Bug report",
				"description": "Something is broken",
				"state":       "opened",
				"web_url":     "https://gitlab.example.com/project/-/issues/1",
				"author":      map[string]any{"username": "user1"},
				"labels":      []string{"bug"},
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := GetIssueInput{
			ProjectID: "test-project",
			IssueIID:  1,
		}

		ctx := context.Background()
		_, output, err := getIssueHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(1), output.IID)
		assert.Equal(t, "Bug report", output.Title)
		assert.Equal(t, "Something is broken", output.Description)
		assert.Equal(t, "user1", output.AuthorName)
	})
}

func TestCreateIssueTool(t *testing.T) {
	t.Run("creates issue successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/issues", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":      103,
				"iid":     3,
				"title":   "New issue",
				"web_url": "https://gitlab.example.com/project/-/issues/3",
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := CreateIssueInput{
			ProjectID: "test-project",
			Title:     "New issue",
		}

		ctx := context.Background()
		_, output, err := createIssueHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(3), output.IID)
		assert.Equal(t, "New issue", output.Title)
	})
}

func TestUpdateIssueTool(t *testing.T) {
	t.Run("updates issue successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/issues/1", r.URL.Path)
			assert.Equal(t, "PUT", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":      101,
				"iid":     1,
				"title":   "Updated title",
				"state":   "opened",
				"web_url": "https://gitlab.example.com/project/-/issues/1",
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		title := "Updated title"
		input := UpdateIssueInput{
			ProjectID: "test-project",
			IssueIID:  1,
			Title:     &title,
		}

		ctx := context.Background()
		_, output, err := updateIssueHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(1), output.IID)
		assert.Equal(t, "Updated title", output.Title)
	})
}

func TestDeleteIssueTool(t *testing.T) {
	t.Run("deletes issue successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/issues/1", r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := DeleteIssueInput{
			ProjectID: "test-project",
			IssueIID:  1,
		}

		ctx := context.Background()
		_, output, err := deleteIssueHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.True(t, output.Success)
	})
}

func TestListIssueNotesTool(t *testing.T) {
	t.Run("returns notes list successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/issues/1/notes", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id":     10,
					"body":   "First comment",
					"author": map[string]any{"username": "user1"},
					"system": false,
				},
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := ListIssueNotesInput{
			ProjectID: "test-project",
			IssueIID:  1,
		}

		ctx := context.Background()
		_, output, err := listIssueNotesHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.Notes, 1)
		assert.Equal(t, "First comment", output.Notes[0].Body)
		assert.Equal(t, "user1", output.Notes[0].AuthorName)
		assert.Equal(t, false, output.Notes[0].System)
	})
}

func TestCreateIssueNoteTool(t *testing.T) {
	t.Run("creates note successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/issues/1/notes", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":     20,
				"body":   "New comment",
				"author": map[string]any{"username": "user1"},
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := CreateIssueNoteInput{
			ProjectID: "test-project",
			IssueIID:  1,
			Body:      "New comment",
		}

		ctx := context.Background()
		_, output, err := createIssueNoteHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(20), output.ID)
		assert.Equal(t, "New comment", output.Body)
	})
}

func TestListIssueDiscussionsTool(t *testing.T) {
	t.Run("returns discussions list successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/issues/1/discussions", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id": "abc123",
					"notes": []map[string]any{
						{
							"id":     30,
							"body":   "Discussion note",
							"author": map[string]any{"username": "user1"},
						},
					},
				},
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := ListIssueDiscussionsInput{
			ProjectID: "test-project",
			IssueIID:  1,
		}

		ctx := context.Background()
		_, output, err := listIssueDiscussionsHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.Discussions, 1)
		assert.Equal(t, "abc123", output.Discussions[0].ID)
		assert.Len(t, output.Discussions[0].Notes, 1)
		assert.Equal(t, "Discussion note", output.Discussions[0].Notes[0].Body)
	})
}

func TestReplyToIssueDiscussionTool(t *testing.T) {
	t.Run("replies to discussion successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/issues/1/discussions/abc123/notes", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":     50,
				"body":   "Reply",
				"author": map[string]any{"username": "user1"},
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := ReplyToIssueDiscussionInput{
			ProjectID:    "test-project",
			IssueIID:     1,
			DiscussionID: "abc123",
			Body:         "Reply",
		}

		ctx := context.Background()
		_, output, err := replyToIssueDiscussionHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(50), output.ID)
		assert.Equal(t, "Reply", output.Body)
	})
}

func TestToolRegistration(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	_, reg, cleanup := setupTestServer(t, handler)
	defer cleanup()

	expectedTools := []string{
		"list_issues",
		"get_issue",
		"create_issue",
		"update_issue",
		"delete_issue",
		"list_issue_notes",
		"create_issue_note",
		"delete_issue_note",
		"list_issue_discussions",
		"create_issue_discussion",
		"reply_to_issue_discussion",
	}

	for _, tool := range expectedTools {
		assert.True(t, reg.IsRegistered(tool), "tool %s should be registered", tool)
	}
}

func TestToolDisabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		GitLabURL:     server.URL,
		GitLabToken:   "test-token",
		DisabledTools: []string{"list_issues"},
	}

	client, err := gitlab.NewClient(server.URL, "test-token")
	require.NoError(t, err)

	reg := registry.New(cfg)
	Register(reg, client)

	assert.True(t, reg.IsRegistered("list_issues"))
	assert.False(t, reg.IsToolEnabled("list_issues"))
}
