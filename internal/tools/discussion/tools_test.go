package discussion

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

func TestAddMergeRequestCommentTool(t *testing.T) {
	t.Run("adds comment successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/notes", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"id":   100,
				"body": "This is a test comment",
				"author": map[string]any{
					"id":       1,
					"username": "testuser",
				},
				"created_at": "2024-01-01T00:00:00Z",
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("add_merge_request_comment"))

		input := AddCommentInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			Body:            "This is a test comment",
		}

		ctx := context.Background()
		_, output, err := addCommentHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(100), output.ID)
		assert.Equal(t, "This is a test comment", output.Body)
	})

	t.Run("returns error for non-existent MR", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := AddCommentInput{
			ProjectID:       "test-project",
			MergeRequestIID: 999,
			Body:            "Comment",
		}

		ctx := context.Background()
		_, _, err := addCommentHandler(client, ctx, nil, input)

		assert.Error(t, err)
	})
}

func TestAddMergeRequestDiscussionTool(t *testing.T) {
	t.Run("creates discussion successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/discussions", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"id": "abc123",
				"notes": []map[string]any{
					{
						"id":   100,
						"body": "Review comment on line",
						"author": map[string]any{
							"id":       1,
							"username": "testuser",
						},
					},
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("add_merge_request_discussion"))

		newLine := 10
		input := AddDiscussionInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			Body:            "Review comment on line",
			Position: &DiffPosition{
				BaseSHA:  "abc",
				StartSHA: "def",
				HeadSHA:  "ghi",
				OldPath:  "file.go",
				NewPath:  "file.go",
				NewLine:  &newLine,
			},
		}

		ctx := context.Background()
		_, output, err := addDiscussionHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, "abc123", output.ID)
	})

	t.Run("creates discussion without position", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"id": "xyz789",
				"notes": []map[string]any{
					{
						"id":   101,
						"body": "General discussion",
					},
				},
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := AddDiscussionInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			Body:            "General discussion",
		}

		ctx := context.Background()
		_, output, err := addDiscussionHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, "xyz789", output.ID)
	})
}

func TestListMergeRequestDiscussionsTool(t *testing.T) {
	t.Run("returns discussions list successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/discussions", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id": "disc1",
					"notes": []map[string]any{
						{
							"id":   100,
							"body": "First discussion",
							"author": map[string]any{
								"id":       1,
								"username": "user1",
							},
							"resolvable": true,
							"resolved":   false,
						},
					},
				},
				{
					"id": "disc2",
					"notes": []map[string]any{
						{
							"id":   101,
							"body": "Second discussion",
							"author": map[string]any{
								"id":       2,
								"username": "user2",
							},
							"resolvable": true,
							"resolved":   true,
						},
					},
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("list_merge_request_discussions"))

		input := ListDiscussionsInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
		}

		ctx := context.Background()
		_, output, err := listDiscussionsHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Len(t, output.Discussions, 2)
		assert.Equal(t, "disc1", output.Discussions[0].ID)
		assert.Equal(t, "disc2", output.Discussions[1].ID)
	})
}

func TestResolveDiscussionTool(t *testing.T) {
	t.Run("resolves discussion successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/discussions/disc123", r.URL.Path)
			assert.Equal(t, "PUT", r.Method)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id": "disc123",
				"notes": []map[string]any{
					{
						"id":         100,
						"body":       "Discussion content",
						"resolvable": true,
						"resolved":   true,
					},
				},
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("resolve_discussion"))

		input := ResolveDiscussionInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			DiscussionID:    "disc123",
			Resolved:        true,
		}

		ctx := context.Background()
		_, output, err := resolveDiscussionHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, "disc123", output.ID)
		assert.True(t, output.Resolved)
	})

	t.Run("unresolves discussion successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id": "disc123",
				"notes": []map[string]any{
					{
						"id":         100,
						"resolvable": true,
						"resolved":   false,
					},
				},
			})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := ResolveDiscussionInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			DiscussionID:    "disc123",
			Resolved:        false,
		}

		ctx := context.Background()
		_, output, err := resolveDiscussionHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.False(t, output.Resolved)
	})
}

func TestDeleteMergeRequestCommentTool(t *testing.T) {
	t.Run("deletes comment successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/notes/123", r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("delete_merge_request_comment"))

		input := DeleteCommentInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			NoteID:          123,
		}

		ctx := context.Background()
		_, output, err := deleteCommentHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.True(t, output.Success)
		assert.Equal(t, "Comment deleted successfully", output.Message)
	})

	t.Run("returns error for non-existent comment", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := DeleteCommentInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			NoteID:          999,
		}

		ctx := context.Background()
		_, _, err := deleteCommentHandler(client, ctx, nil, input)

		assert.Error(t, err)
	})

	t.Run("returns error for forbidden access", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "403 Forbidden"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := DeleteCommentInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			NoteID:          123,
		}

		ctx := context.Background()
		_, _, err := deleteCommentHandler(client, ctx, nil, input)

		assert.Error(t, err)
	})
}

func TestReplyToMergeRequestCommentTool(t *testing.T) {
	t.Run("replies to discussion successfully", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/discussions/disc123/notes", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"id":   200,
				"body": "This is a reply to the discussion",
				"author": map[string]any{
					"id":       1,
					"username": "testuser",
				},
				"created_at": "2024-01-01T12:00:00Z",
			})
		}

		client, reg, cleanup := setupTestServer(t, handler)
		defer cleanup()

		assert.True(t, reg.IsRegistered("reply_to_merge_request_comment"))

		input := ReplyToCommentInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			DiscussionID:    "disc123",
			Body:            "This is a reply to the discussion",
		}

		ctx := context.Background()
		_, output, err := replyToCommentHandler(client, ctx, nil, input)

		require.NoError(t, err)
		assert.Equal(t, int64(200), output.ID)
		assert.Equal(t, "This is a reply to the discussion", output.Body)
		assert.Equal(t, "testuser", output.AuthorName)
	})

	t.Run("returns error for non-existent discussion", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
		}

		client, _, cleanup := setupTestServer(t, handler)
		defer cleanup()

		input := ReplyToCommentInput{
			ProjectID:       "test-project",
			MergeRequestIID: 1,
			DiscussionID:    "unknown-disc",
			Body:            "Reply",
		}

		ctx := context.Background()
		_, _, err := replyToCommentHandler(client, ctx, nil, input)

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
			DisabledTools: []string{"add_merge_request_comment"},
		}

		client, err := gitlab.NewClient(server.URL, "test-token")
		require.NoError(t, err)

		reg := registry.New(cfg)
		Register(reg, client)

		assert.True(t, reg.IsRegistered("add_merge_request_comment"))
		assert.False(t, reg.IsToolEnabled("add_merge_request_comment"))
	})
}
