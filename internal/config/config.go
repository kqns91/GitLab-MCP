package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Config はアプリケーション設定を保持する
type Config struct {
	GitLabURL     string
	GitLabToken   string
	EnabledTools  []string // nil = all enabled
	DisabledTools []string
	Debug         bool
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	gitlabURL := os.Getenv("GITLAB_URL")
	if gitlabURL == "" {
		return nil, errors.New("GITLAB_URL environment variable is required")
	}

	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		return nil, errors.New("GITLAB_TOKEN environment variable is required")
	}

	cfg := &Config{
		GitLabURL:   gitlabURL,
		GitLabToken: gitlabToken,
		Debug:       parseDebug(os.Getenv("GITLAB_MCP_DEBUG")),
	}

	if enabledTools := os.Getenv("GITLAB_MCP_ENABLED_TOOLS"); enabledTools != "" {
		cfg.EnabledTools = parseToolList(enabledTools)
	}

	if disabledTools := os.Getenv("GITLAB_MCP_DISABLED_TOOLS"); disabledTools != "" {
		cfg.DisabledTools = parseToolList(disabledTools)
	}

	return cfg, nil
}

// parseDebug は debug フラグをパースする
func parseDebug(value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	return v == "true" || v == "1" || v == "yes"
}

// parseToolList はカンマ区切りのツール名リストをパースする
func parseToolList(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	tools := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			tools = append(tools, trimmed)
		}
	}
	return tools
}

// String は設定の文字列表現を返す（トークンはマスキング）
func (c *Config) String() string {
	maskedToken := "***"
	if len(c.GitLabToken) > 4 {
		maskedToken = c.GitLabToken[:2] + "***" + c.GitLabToken[len(c.GitLabToken)-2:]
	}
	return fmt.Sprintf("Config{GitLabURL: %q, GitLabToken: %q, EnabledTools: %v, DisabledTools: %v, Debug: %v}",
		c.GitLabURL, maskedToken, c.EnabledTools, c.DisabledTools, c.Debug)
}

// IsToolEnabled はツールが有効かどうかを判定する
// DISABLED_TOOLS が ENABLED_TOOLS より優先される
func (c *Config) IsToolEnabled(toolName string) bool {
	// Check disabled list first (takes precedence)
	if slices.Contains(c.DisabledTools, toolName) {
		return false
	}

	// If enabled list is specified, tool must be in it
	if len(c.EnabledTools) > 0 {
		return slices.Contains(c.EnabledTools, toolName)
	}

	// No restrictions - tool is enabled
	return true
}
