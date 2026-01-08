package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kqns91/gitlab-mcp/internal/config"
	"github.com/kqns91/gitlab-mcp/internal/gitlab"
	"github.com/kqns91/gitlab-mcp/internal/registry"
	"github.com/kqns91/gitlab-mcp/internal/tools/approval"
	"github.com/kqns91/gitlab-mcp/internal/tools/discussion"
	"github.com/kqns91/gitlab-mcp/internal/tools/mergerequest"
	"github.com/kqns91/gitlab-mcp/internal/tools/pipeline"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if cfg.Debug {
		log.Printf("Configuration loaded: %s", cfg)
	}

	// Initialize GitLab client
	client, err := gitlab.NewClient(cfg.GitLabURL, cfg.GitLabToken)
	if err != nil {
		return fmt.Errorf("failed to create GitLab client: %w", err)
	}

	if cfg.Debug {
		log.Printf("GitLab client initialized for %s", cfg.GitLabURL)
	}

	// Create registry and register all tools
	reg := registry.New(cfg)
	registerAllTools(reg, client)

	if cfg.Debug {
		enabled := reg.GetEnabledTools()
		log.Printf("Registered %d tools: %v", len(enabled), enabled)
	}

	// Start MCP server with stdio transport
	server := reg.Server()

	if cfg.Debug {
		log.Printf("Starting MCP server with stdio transport")
	}

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// registerAllTools は全てのツールをレジストリに登録する
func registerAllTools(reg *registry.Registry, client *gitlab.Client) {
	mergerequest.Register(reg, client)
	discussion.Register(reg, client)
	approval.Register(reg, client)
	pipeline.Register(reg, client)
}

// Version information (set via ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	// Print version information if requested
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("gitlab-mcp %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}
}
