package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	// Setup
	os.Setenv("GITLAB_URL", "https://gitlab.example.com")
	os.Setenv("GITLAB_TOKEN", "test-token")
	defer func() {
		os.Unsetenv("GITLAB_URL")
		os.Unsetenv("GITLAB_TOKEN")
	}()

	// Execute
	cfg, err := Load()

	// Verify
	require.NoError(t, err)
	assert.Equal(t, "https://gitlab.example.com", cfg.GitLabURL)
	assert.Equal(t, "test-token", cfg.GitLabToken)
	assert.False(t, cfg.Debug)
}

func TestLoad_MissingURL(t *testing.T) {
	// Setup
	os.Unsetenv("GITLAB_URL")
	os.Setenv("GITLAB_TOKEN", "test-token")
	defer os.Unsetenv("GITLAB_TOKEN")

	// Execute
	cfg, err := Load()

	// Verify
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GITLAB_URL")
}

func TestLoad_MissingToken(t *testing.T) {
	// Setup
	os.Setenv("GITLAB_URL", "https://gitlab.example.com")
	os.Unsetenv("GITLAB_TOKEN")
	defer os.Unsetenv("GITLAB_URL")

	// Execute
	cfg, err := Load()

	// Verify
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GITLAB_TOKEN")
}

func TestLoad_DebugMode(t *testing.T) {
	// Setup
	os.Setenv("GITLAB_URL", "https://gitlab.example.com")
	os.Setenv("GITLAB_TOKEN", "test-token")
	os.Setenv("GITLAB_MCP_DEBUG", "true")
	defer func() {
		os.Unsetenv("GITLAB_URL")
		os.Unsetenv("GITLAB_TOKEN")
		os.Unsetenv("GITLAB_MCP_DEBUG")
	}()

	// Execute
	cfg, err := Load()

	// Verify
	require.NoError(t, err)
	assert.True(t, cfg.Debug)
}

func TestLoad_EnabledTools(t *testing.T) {
	// Setup
	os.Setenv("GITLAB_URL", "https://gitlab.example.com")
	os.Setenv("GITLAB_TOKEN", "test-token")
	os.Setenv("GITLAB_MCP_ENABLED_TOOLS", "list_merge_requests,get_merge_request")
	defer func() {
		os.Unsetenv("GITLAB_URL")
		os.Unsetenv("GITLAB_TOKEN")
		os.Unsetenv("GITLAB_MCP_ENABLED_TOOLS")
	}()

	// Execute
	cfg, err := Load()

	// Verify
	require.NoError(t, err)
	assert.Equal(t, []string{"list_merge_requests", "get_merge_request"}, cfg.EnabledTools)
}

func TestLoad_DisabledTools(t *testing.T) {
	// Setup
	os.Setenv("GITLAB_URL", "https://gitlab.example.com")
	os.Setenv("GITLAB_TOKEN", "test-token")
	os.Setenv("GITLAB_MCP_DISABLED_TOOLS", "merge_merge_request,approve_merge_request")
	defer func() {
		os.Unsetenv("GITLAB_URL")
		os.Unsetenv("GITLAB_TOKEN")
		os.Unsetenv("GITLAB_MCP_DISABLED_TOOLS")
	}()

	// Execute
	cfg, err := Load()

	// Verify
	require.NoError(t, err)
	assert.Equal(t, []string{"merge_merge_request", "approve_merge_request"}, cfg.DisabledTools)
}

func TestConfig_String_MasksToken(t *testing.T) {
	cfg := &Config{
		GitLabURL:   "https://gitlab.example.com",
		GitLabToken: "super-secret-token",
		Debug:       true,
	}

	// Token should not appear in string representation
	str := cfg.String()
	assert.NotContains(t, str, "super-secret-token")
	assert.Contains(t, str, "***")
}

func TestIsToolEnabled_NoRestrictions(t *testing.T) {
	cfg := &Config{
		GitLabURL:     "https://gitlab.example.com",
		GitLabToken:   "test-token",
		EnabledTools:  nil,
		DisabledTools: nil,
	}

	// All tools should be enabled when no restrictions
	assert.True(t, cfg.IsToolEnabled("list_merge_requests"))
	assert.True(t, cfg.IsToolEnabled("merge_merge_request"))
	assert.True(t, cfg.IsToolEnabled("any_tool"))
}

func TestIsToolEnabled_OnlyEnabledTools(t *testing.T) {
	cfg := &Config{
		GitLabURL:    "https://gitlab.example.com",
		GitLabToken:  "test-token",
		EnabledTools: []string{"list_merge_requests", "get_merge_request"},
	}

	// Only listed tools should be enabled
	assert.True(t, cfg.IsToolEnabled("list_merge_requests"))
	assert.True(t, cfg.IsToolEnabled("get_merge_request"))
	assert.False(t, cfg.IsToolEnabled("merge_merge_request"))
	assert.False(t, cfg.IsToolEnabled("unknown_tool"))
}

func TestIsToolEnabled_OnlyDisabledTools(t *testing.T) {
	cfg := &Config{
		GitLabURL:     "https://gitlab.example.com",
		GitLabToken:   "test-token",
		DisabledTools: []string{"merge_merge_request", "approve_merge_request"},
	}

	// All tools except disabled ones should be enabled
	assert.True(t, cfg.IsToolEnabled("list_merge_requests"))
	assert.True(t, cfg.IsToolEnabled("get_merge_request"))
	assert.False(t, cfg.IsToolEnabled("merge_merge_request"))
	assert.False(t, cfg.IsToolEnabled("approve_merge_request"))
}

func TestIsToolEnabled_DisabledTakesPrecedence(t *testing.T) {
	cfg := &Config{
		GitLabURL:     "https://gitlab.example.com",
		GitLabToken:   "test-token",
		EnabledTools:  []string{"list_merge_requests", "merge_merge_request"},
		DisabledTools: []string{"merge_merge_request"},
	}

	// DISABLED_TOOLS should take precedence over ENABLED_TOOLS
	assert.True(t, cfg.IsToolEnabled("list_merge_requests"))
	assert.False(t, cfg.IsToolEnabled("merge_merge_request")) // disabled even though in enabled list
	assert.False(t, cfg.IsToolEnabled("get_merge_request"))   // not in enabled list
}
