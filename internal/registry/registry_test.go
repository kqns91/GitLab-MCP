package registry

import (
	"context"
	"testing"

	"github.com/kqns91/gitlab-mcp/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// DummyInput is test input type
type DummyInput struct {
	Name string `json:"name"`
}

// DummyOutput is test output type
type DummyOutput struct {
	Result string `json:"result"`
}

func dummyHandler(ctx context.Context, req *mcp.CallToolRequest, input DummyInput) (*mcp.CallToolResult, DummyOutput, error) {
	return nil, DummyOutput{Result: "ok"}, nil
}

func TestNewRegistry(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	reg := New(cfg)

	assert.NotNil(t, reg)
	assert.NotNil(t, reg.Server())
}

func TestRegistry_RegisterTool(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	reg := New(cfg)

	// Register a tool
	RegisterTool(reg, "test_tool", "Test tool description", dummyHandler)

	// Tool should be registered
	assert.True(t, reg.IsRegistered("test_tool"))
	assert.False(t, reg.IsRegistered("unknown_tool"))
}

func TestRegistry_IsToolEnabled_NoRestrictions(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "test-token",
	}

	reg := New(cfg)
	RegisterTool(reg, "test_tool", "Test tool", dummyHandler)

	// Without restrictions, tool should be enabled
	assert.True(t, reg.IsToolEnabled("test_tool"))
}

func TestRegistry_IsToolEnabled_Disabled(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:     "https://gitlab.example.com",
		GitLabToken:   "test-token",
		DisabledTools: []string{"test_tool"},
	}

	reg := New(cfg)
	RegisterTool(reg, "test_tool", "Test tool", dummyHandler)

	// Tool should be disabled
	assert.False(t, reg.IsToolEnabled("test_tool"))
}

func TestRegistry_IsToolEnabled_OnlyEnabled(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:    "https://gitlab.example.com",
		GitLabToken:  "test-token",
		EnabledTools: []string{"allowed_tool"},
	}

	reg := New(cfg)
	RegisterTool(reg, "test_tool", "Test tool", dummyHandler)
	RegisterTool(reg, "allowed_tool", "Allowed tool", dummyHandler)

	// Only allowed_tool should be enabled
	assert.False(t, reg.IsToolEnabled("test_tool"))
	assert.True(t, reg.IsToolEnabled("allowed_tool"))
}

func TestRegistry_GetEnabledTools(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:     "https://gitlab.example.com",
		GitLabToken:   "test-token",
		DisabledTools: []string{"disabled_tool"},
	}

	reg := New(cfg)
	RegisterTool(reg, "enabled_tool", "Enabled tool", dummyHandler)
	RegisterTool(reg, "disabled_tool", "Disabled tool", dummyHandler)

	enabledTools := reg.GetEnabledTools()

	assert.Contains(t, enabledTools, "enabled_tool")
	assert.NotContains(t, enabledTools, "disabled_tool")
}

func TestRegistry_CheckToolEnabled_Error(t *testing.T) {
	cfg := &config.Config{
		GitLabURL:     "https://gitlab.example.com",
		GitLabToken:   "test-token",
		DisabledTools: []string{"disabled_tool"},
	}

	reg := New(cfg)
	RegisterTool(reg, "disabled_tool", "Disabled tool", dummyHandler)

	// Should return error for disabled tool
	err := reg.CheckToolEnabled("disabled_tool")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "disabled_tool")
	assert.Contains(t, err.Error(), "無効")

	// Should return nil for enabled tool
	RegisterTool(reg, "enabled_tool", "Enabled tool", dummyHandler)
	err = reg.CheckToolEnabled("enabled_tool")
	assert.NoError(t, err)
}
