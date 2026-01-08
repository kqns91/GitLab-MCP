package gitlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_Success(t *testing.T) {
	client, err := NewClient("https://gitlab.example.com", "test-token")

	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewClient_EmptyURL(t *testing.T) {
	client, err := NewClient("", "test-token")

	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "URL")
}

func TestNewClient_EmptyToken(t *testing.T) {
	client, err := NewClient("https://gitlab.example.com", "")

	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token")
}

func TestClient_ServiceAccessors(t *testing.T) {
	client, err := NewClient("https://gitlab.example.com", "test-token")
	require.NoError(t, err)

	// Verify all service accessors return non-nil
	assert.NotNil(t, client.MergeRequests())
	assert.NotNil(t, client.Discussions())
	assert.NotNil(t, client.MergeRequestApprovals())
	assert.NotNil(t, client.Pipelines())
	assert.NotNil(t, client.Jobs())
	assert.NotNil(t, client.Notes())
}
