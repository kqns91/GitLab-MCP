package registry

import (
	"context"

	"github.com/kqns91/gitlab-mcp/internal/config"
	"github.com/kqns91/gitlab-mcp/internal/gitlab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Registry はツールの登録と管理を行う
type Registry struct {
	server          *mcp.Server
	config          *config.Config
	registeredTools map[string]bool
}

// New は新しい Registry を作成する
func New(cfg *config.Config) *Registry {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "gitlab-mcp",
			Version: "1.0.0",
		},
		nil,
	)

	return &Registry{
		server:          server,
		config:          cfg,
		registeredTools: make(map[string]bool),
	}
}

// Server は MCP サーバーを返す
func (r *Registry) Server() *mcp.Server {
	return r.server
}

// IsRegistered はツールが登録されているかを返す
func (r *Registry) IsRegistered(toolName string) bool {
	return r.registeredTools[toolName]
}

// IsToolEnabled はツールが有効かどうかを返す
func (r *Registry) IsToolEnabled(toolName string) bool {
	return r.config.IsToolEnabled(toolName)
}

// CheckToolEnabled はツールが有効でない場合にエラーを返す
func (r *Registry) CheckToolEnabled(toolName string) error {
	if !r.IsToolEnabled(toolName) {
		return gitlab.NewToolDisabledError(toolName)
	}
	return nil
}

// GetEnabledTools は有効なツールの名前リストを返す
func (r *Registry) GetEnabledTools() []string {
	var enabled []string
	for name := range r.registeredTools {
		if r.IsToolEnabled(name) {
			enabled = append(enabled, name)
		}
	}
	return enabled
}

// ToolHandlerFor is a type alias for MCP tool handlers
type ToolHandlerFor[In, Out any] func(ctx context.Context, req *mcp.CallToolRequest, input In) (*mcp.CallToolResult, Out, error)

// RegisterTool は新しいツールを登録する
// ツールが無効化されている場合でも登録はするが、呼び出し時にチェックされる
func RegisterTool[In, Out any](r *Registry, name, description string, handler ToolHandlerFor[In, Out]) {
	r.registeredTools[name] = true

	// Wrap handler to check if tool is enabled
	wrappedHandler := func(ctx context.Context, req *mcp.CallToolRequest, input In) (*mcp.CallToolResult, Out, error) {
		if err := r.CheckToolEnabled(name); err != nil {
			var zero Out
			return nil, zero, err
		}
		return handler(ctx, req, input)
	}

	// Only add to server if enabled (to exclude from tools/list)
	if r.config.IsToolEnabled(name) {
		mcp.AddTool(r.server, &mcp.Tool{
			Name:        name,
			Description: description,
		}, wrappedHandler)
	}
}
