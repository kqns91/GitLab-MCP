# Project Structure

## Organization Philosophy

Go 標準のプロジェクトレイアウトに準拠した Package-based Modular 構造。`cmd/` にエントリーポイント、`internal/` にビジネスロジックを配置。

## Directory Patterns

### Entry Point
**Location**: `cmd/gitlab-mcp/`
**Purpose**: アプリケーションのエントリーポイント
**Example**: `main.go` - サーバー起動と依存関係の初期化

### Configuration
**Location**: `internal/config/`
**Purpose**: 環境変数からの設定読み込みと管理
**Example**: `config.go` - Config 構造体と Load 関数

### API Client
**Location**: `internal/gitlab/`
**Purpose**: GitLab Go SDK のラッパーとエラー変換
**Contents**: SDK クライアント初期化、サービスアクセサ、エラー型変換

### Tool Registry
**Location**: `internal/registry/`
**Purpose**: MCP ツールの登録と有効/無効管理

### Tools
**Location**: `internal/tools/{domain}/`
**Purpose**: 機能ドメイン別の MCP ツール実装
**Domains**: `mergerequest/`, `discussion/`, `approval/`, `pipeline/`

### Documentation
**Location**: `docs/`
**Purpose**: 追加ドキュメント（日本語版など）

## Naming Conventions

- **Files**: snake_case (`config.go`, `client_test.go`)
- **Packages**: 小文字単語 (`config`, `gitlab`, `mergerequest`)
- **Types**: PascalCase (`MergeRequest`, `Config`)
- **Functions**: PascalCase for exported, camelCase for unexported

## Import Organization

```go
import (
    // Standard library
    "context"
    "net/http"

    // External packages
    "github.com/modelcontextprotocol/go-sdk/mcp"

    // Internal packages
    "gitlab-mcp/internal/config"
    "gitlab-mcp/internal/gitlab"
)
```

## Code Organization Principles

- **internal/ の活用**: 外部からのインポートを防止
- **ドメイン分離**: tools/ 配下で機能ドメインごとにパッケージ分割
- **依存方向**: Tools → GitLab Client → External API（一方向）
- **テストの配置**: 実装ファイルと同じディレクトリに `*_test.go`

---
_Document patterns, not file trees. New files following patterns shouldn't require updates_
