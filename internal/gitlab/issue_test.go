package gitlab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProjectIssues_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	issues, err := client.ListProjectIssues("test-project", nil)

	require.NoError(t, err)
	assert.Len(t, issues, 2)
	assert.Equal(t, int64(1), issues[0].IID)
	assert.Equal(t, "Bug report", issues[0].Title)
	assert.Equal(t, "opened", issues[0].State)
	assert.Equal(t, int64(2), issues[1].IID)
}

func TestListProjectIssues_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	issues, err := client.ListProjectIssues("unknown-project", nil)

	assert.Nil(t, issues)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}

func TestGetIssue_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	issue, err := client.GetIssue("test-project", 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), issue.IID)
	assert.Equal(t, "Bug report", issue.Title)
	assert.Equal(t, "Something is broken", issue.Description)
	assert.Equal(t, "opened", issue.State)
}

func TestGetIssue_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	issue, err := client.GetIssue("test-project", 999)

	assert.Nil(t, issue)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}

func TestCreateIssue_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/issues", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":      103,
			"iid":     3,
			"title":   "New issue",
			"state":   "opened",
			"web_url": "https://gitlab.example.com/project/-/issues/3",
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	issue, err := client.CreateIssue("test-project", &CreateIssueOptions{
		Title: "New issue",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(3), issue.IID)
	assert.Equal(t, "New issue", issue.Title)
}

func TestUpdateIssue_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	title := "Updated title"
	issue, err := client.UpdateIssue("test-project", 1, &UpdateIssueOptions{
		Title: &title,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), issue.IID)
	assert.Equal(t, "Updated title", issue.Title)
}

func TestDeleteIssue_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/issues/1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	err = client.DeleteIssue("test-project", 1)

	assert.NoError(t, err)
}

func TestDeleteIssue_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	err = client.DeleteIssue("test-project", 999)

	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}

func TestListIssueNotes_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			{
				"id":     11,
				"body":   "Second comment",
				"author": map[string]any{"username": "user2"},
				"system": true,
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	notes, err := client.ListIssueNotes("test-project", 1, nil)

	require.NoError(t, err)
	assert.Len(t, notes, 2)
	assert.Equal(t, "First comment", notes[0].Body)
	assert.Equal(t, false, notes[0].System)
	assert.Equal(t, "Second comment", notes[1].Body)
	assert.Equal(t, true, notes[1].System)
}

func TestCreateIssueNote_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/issues/1/notes", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":     20,
			"body":   "New comment",
			"author": map[string]any{"username": "user1"},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	note, err := client.CreateIssueNote("test-project", 1, "New comment")

	require.NoError(t, err)
	assert.Equal(t, "New comment", note.Body)
}

func TestDeleteIssueNote_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/issues/1/notes/10", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	err = client.DeleteIssueNote("test-project", 1, 10)

	assert.NoError(t, err)
}

func TestListIssueDiscussions_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	discussions, err := client.ListIssueDiscussions("test-project", 1, nil)

	require.NoError(t, err)
	assert.Len(t, discussions, 1)
	assert.Equal(t, "abc123", discussions[0].ID)
	assert.Len(t, discussions[0].Notes, 1)
	assert.Equal(t, "Discussion note", discussions[0].Notes[0].Body)
}

func TestCreateIssueDiscussion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/issues/1/discussions", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": "def456",
			"notes": []map[string]any{
				{
					"id":     40,
					"body":   "New discussion",
					"author": map[string]any{"username": "user1"},
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	discussion, err := client.CreateIssueDiscussion("test-project", 1, "New discussion")

	require.NoError(t, err)
	assert.Equal(t, "def456", discussion.ID)
}

func TestAddIssueDiscussionNote_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/issues/1/discussions/abc123/notes", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":     50,
			"body":   "Reply note",
			"author": map[string]any{"username": "user1"},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	note, err := client.AddIssueDiscussionNote("test-project", 1, "abc123", "Reply note")

	require.NoError(t, err)
	assert.Equal(t, "Reply note", note.Body)
}
