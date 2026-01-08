package gitlab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddMergeRequestComment_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/notes", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":   1,
			"body": "This is a comment",
			"author": map[string]any{
				"id":       1,
				"username": "testuser",
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	note, err := client.AddMergeRequestComment("test-project", 1, "This is a comment")

	require.NoError(t, err)
	assert.Equal(t, int64(1), note.ID)
	assert.Equal(t, "This is a comment", note.Body)
}

func TestAddMergeRequestComment_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "404 Not found"})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	note, err := client.AddMergeRequestComment("unknown-project", 1, "comment")

	assert.Nil(t, note)
	assert.Error(t, err)
	mcpErr, ok := err.(*MCPError)
	require.True(t, ok)
	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
}

func TestCreateMergeRequestDiscussion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/discussions", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":              "abc123",
			"individual_note": false,
			"notes": []map[string]any{
				{
					"id":         1,
					"body":       "Line comment",
					"resolvable": true,
					"resolved":   false,
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	discussion, err := client.CreateMergeRequestDiscussion("test-project", 1, &CreateDiscussionOptions{
		Body:     "Line comment",
		FilePath: "main.go",
		NewLine:  intPtr(10),
	})

	require.NoError(t, err)
	assert.Equal(t, "abc123", discussion.ID)
	assert.Len(t, discussion.Notes, 1)
	assert.Equal(t, "Line comment", discussion.Notes[0].Body)
}

func TestListMergeRequestDiscussions_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/discussions", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":              "disc1",
				"individual_note": false,
				"notes": []map[string]any{
					{"id": 1, "body": "First comment", "resolved": false},
				},
			},
			{
				"id":              "disc2",
				"individual_note": false,
				"notes": []map[string]any{
					{"id": 2, "body": "Second comment", "resolved": true},
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	discussions, err := client.ListMergeRequestDiscussions("test-project", 1, nil)

	require.NoError(t, err)
	assert.Len(t, discussions, 2)
	assert.Equal(t, "disc1", discussions[0].ID)
	assert.Equal(t, "disc2", discussions[1].ID)
}

func TestResolveMergeRequestDiscussion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/test-project/merge_requests/1/discussions/disc123", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":              "disc123",
			"individual_note": false,
			"notes": []map[string]any{
				{"id": 1, "body": "Comment", "resolved": true},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	discussion, err := client.ResolveMergeRequestDiscussion("test-project", 1, "disc123", true)

	require.NoError(t, err)
	assert.Equal(t, "disc123", discussion.ID)
	assert.True(t, discussion.Notes[0].Resolved)
}

func TestResolveMergeRequestDiscussion_Unresolve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":              "disc123",
			"individual_note": false,
			"notes": []map[string]any{
				{"id": 1, "body": "Comment", "resolved": false},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	require.NoError(t, err)

	discussion, err := client.ResolveMergeRequestDiscussion("test-project", 1, "disc123", false)

	require.NoError(t, err)
	assert.False(t, discussion.Notes[0].Resolved)
}

// Helper functions
func intPtr(i int) *int {
	return &i
}
