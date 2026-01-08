# Technology Stack

## Architecture

Package-based Modular アーキテクチャを採用。機能単位でパッケージを分離し、Go 標準のプロジェクト構造に準拠。

```
MCP Server → Tool Registry → Tools → GitLab Client → GitLab API
```

## Core Technologies

- **Language**: Go 1.25.5
- **MCP SDK**: github.com/modelcontextprotocol/go-sdk v1.2.0+
- **GitLab SDK**: gitlab.com/gitlab-org/api/client-go v1.0+
- **Transport**: stdio

## Key Libraries

- **MCP Go SDK**: MCP プロトコル実装（公式 SDK、Google/Anthropic 共同メンテナンス）
- **GitLab Go SDK**: GitLab API 通信（公式 SDK、GitLab メンテナンス）
- **testify**: テストアサーションとモック

## Development Standards

### Type Safety
- struct タグによる JSON マッピングと jsonschema 生成
- ポインタ型によるオプショナルフィールドの明示

### Code Quality
- `go fmt` による統一フォーマット
- `go vet` による静的解析
- パッケージごとの `*_test.go` によるテスト

### Testing
- `testing` 標準パッケージ + testify
- `httptest` による HTTP モック
- ユニットテスト + 統合テスト

## Development Environment

### Required Tools
- Go 1.25.5
- Git

### Common Commands
```bash
# Build
go build -o gitlab-mcp ./cmd/gitlab-mcp

# Test
go test ./...

# Run
GITLAB_URL=https://gitlab.example.com GITLAB_TOKEN=xxx ./gitlab-mcp
```

## Key Technical Decisions

| Decision | Rationale |
|----------|-----------|
| 公式 MCP Go SDK 採用 | 長期サポートと安定性（Google + Anthropic） |
| 公式 GitLab Go SDK 採用 | 型定義・エラー処理済み、開発効率向上 |
| 環境変数のみの設定 | CLI ツールに最適、CI/CD 親和性 |
| DISABLED_TOOLS 優先 | セキュリティ上、明示的な無効化を優先 |
| PAT 認証のみ | シンプルさ優先、OAuth2 フローは非対応 |

---
_Document standards and patterns, not every dependency_
