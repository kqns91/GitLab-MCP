# GitLab MCP Server

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![MCP](https://img.shields.io/badge/MCP-Compatible-green.svg)](https://modelcontextprotocol.io/)

A Model Context Protocol (MCP) server that enables AI agents to interact with GitLab Merge Requests. This server exposes GitLab MR operations as MCP tools, allowing seamless integration with Claude Code, Claude Desktop, and other MCP-compatible clients.

[日本語版ドキュメント](docs/README.ja.md)

## Features

- **MR Lifecycle Management**: Create, update, merge, and close merge requests
- **Code Review Support**: Add comments, create discussions, approve/unapprove MRs
- **CI/CD Integration**: View pipeline status and job details
- **Change Analysis**: Get detailed file diffs and changes
- **Flexible Access Control**: Enable/disable tools via environment variables
- **Secure**: Personal Access Token authentication with token masking in logs

## Quick Start

Get up and running in under 5 minutes:

### 1. Install

```bash
# Option A: go install (recommended)
go install github.com/kqns91/gitlab-mcp/cmd/gitlab-mcp@latest

# Option B: Build from source
git clone https://github.com/kqns91/gitlab-mcp.git
cd gitlab-mcp
go build -o gitlab-mcp ./cmd/gitlab-mcp
```

### 2. Create GitLab Token

Create a Personal Access Token with `api` scope:
1. Go to GitLab → Settings → Access Tokens
2. Create a token with `api` scope
3. Copy the token

### 3. Configure MCP Client

Add to your Claude Code or Claude Desktop config (see [Usage with MCP Clients](#usage-with-mcp-clients) for details):

```json
{
  "mcpServers": {
    "gitlab": {
      "command": "gitlab-mcp",
      "env": {
        "GITLAB_URL": "https://gitlab.com",
        "GITLAB_TOKEN": "your-token"
      }
    }
  }
}
```

## Configuration

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GITLAB_URL` | Yes | GitLab instance URL (e.g., `https://gitlab.com`) |
| `GITLAB_TOKEN` | Yes | Personal Access Token with `api` scope |
| `GITLAB_MCP_ENABLED_TOOLS` | No | Comma-separated list of tools to enable (enables all if not set) |
| `GITLAB_MCP_DISABLED_TOOLS` | No | Comma-separated list of tools to disable (takes precedence over enabled) |
| `GITLAB_MCP_DEBUG` | No | Enable debug logging (`true`, `1`, or `yes`) |

### Tool Filtering Examples

```bash
# Enable only read operations
export GITLAB_MCP_ENABLED_TOOLS="list_merge_requests,get_merge_request,get_merge_request_changes"

# Disable destructive operations
export GITLAB_MCP_DISABLED_TOOLS="merge_merge_request,approve_merge_request"

# DISABLED_TOOLS takes precedence over ENABLED_TOOLS
```

## Available Tools

### Merge Request Operations

| Tool | Description |
|------|-------------|
| `list_merge_requests` | List merge requests in a project with filtering options |
| `get_merge_request` | Get detailed information about a specific merge request |
| `create_merge_request` | Create a new merge request |
| `update_merge_request` | Update an existing merge request |
| `merge_merge_request` | Merge a merge request (with squash and delete branch options) |
| `get_merge_request_changes` | Get file changes/diffs in a merge request |

### Discussion & Comments

| Tool | Description |
|------|-------------|
| `add_merge_request_comment` | Add a general comment to a merge request |
| `add_merge_request_discussion` | Create a line-specific discussion on code |
| `list_merge_request_discussions` | List all discussions on a merge request |
| `resolve_discussion` | Resolve or unresolve a discussion |

### Approval

| Tool | Description |
|------|-------------|
| `approve_merge_request` | Approve a merge request |
| `unapprove_merge_request` | Remove approval from a merge request |
| `get_merge_request_approvals` | Get approval status and list of approvers |

### Pipeline & CI/CD

| Tool | Description |
|------|-------------|
| `list_merge_request_pipelines` | List pipelines associated with a merge request |
| `get_pipeline_jobs` | Get jobs in a specific pipeline |

## Usage with MCP Clients

### Claude Code

Add to your Claude Code configuration:

```json
{
  "mcpServers": {
    "gitlab": {
      "command": "/path/to/gitlab-mcp",
      "env": {
        "GITLAB_URL": "https://gitlab.com",
        "GITLAB_TOKEN": "your-token"
      }
    }
  }
}
```

### Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "gitlab": {
      "command": "/path/to/gitlab-mcp",
      "env": {
        "GITLAB_URL": "https://gitlab.com",
        "GITLAB_TOKEN": "your-token"
      }
    }
  }
}
```

**Config file locations:**
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`
- Linux: `~/.config/Claude/claude_desktop_config.json`

## Example Interactions

Once configured, you can ask Claude to:

- "List all open merge requests in project mygroup/myproject"
- "Get the details of MR #42 in myproject"
- "Create a merge request from feature-branch to main with title 'Add new feature'"
- "Add a comment to MR #42 saying 'LGTM!'"
- "Approve MR #42"
- "What's the pipeline status for MR #42?"

## Development

### Prerequisites

- Go 1.23 or later
- Git

### Build

```bash
go build -o gitlab-mcp ./cmd/gitlab-mcp
```

### Test

```bash
go test ./...
```

### Project Structure

```
.
├── cmd/gitlab-mcp/        # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── gitlab/            # GitLab API client wrapper
│   ├── registry/          # MCP tool registry
│   └── tools/             # MCP tool implementations
│       ├── approval/      # Approval tools
│       ├── discussion/    # Discussion tools
│       ├── mergerequest/  # Merge request tools
│       └── pipeline/      # Pipeline tools
└── test/integration/      # Integration tests
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io/) - The protocol specification
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Official Go SDK
- [GitLab Go SDK](https://gitlab.com/gitlab-org/api/client-go) - Official GitLab API client
