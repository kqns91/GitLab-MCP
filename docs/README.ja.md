# GitLab MCP Server

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)
[![MCP](https://img.shields.io/badge/MCP-Compatible-green.svg)](https://modelcontextprotocol.io/)

AI エージェントが GitLab の Merge Request を操作できるようにする Model Context Protocol (MCP) サーバーです。GitLab の MR 操作を MCP ツールとして公開し、Claude Code、Claude Desktop、その他の MCP 対応クライアントとシームレスに連携できます。

[English Documentation](../README.md)

## 機能

- **MR ライフサイクル管理**: Merge Request の作成、更新、マージ、クローズ
- **コードレビュー支援**: コメント追加、ディスカッション作成、承認/承認取消
- **CI/CD 連携**: パイプライン状態とジョブ詳細の確認
- **変更分析**: ファイル差分と変更内容の詳細取得
- **柔軟なアクセス制御**: 環境変数によるツールの有効化/無効化
- **セキュア**: Personal Access Token 認証、ログへのトークン出力防止

## クイックスタート

5分以内でセットアップできます：

### 1. インストール

```bash
# 方法 A: go install（推奨）
go install github.com/kqns91/gitlab-mcp/cmd/gitlab-mcp@latest

# 方法 B: ソースからビルド
git clone https://github.com/kqns91/gitlab-mcp.git
cd gitlab-mcp
go build -o gitlab-mcp ./cmd/gitlab-mcp
```

### 2. GitLab トークンの作成

`api` スコープを持つ Personal Access Token を作成します：
1. GitLab → 設定 → アクセストークン に移動
2. `api` スコープを持つトークンを作成
3. トークンをコピー

### 3. MCP クライアントの設定

Claude Code または Claude Desktop の設定に追加（詳細は [MCP クライアントでの使用方法](#mcp-クライアントでの使用方法) を参照）：

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

## 設定

### 環境変数

| 変数名 | 必須 | 説明 |
|--------|------|------|
| `GITLAB_URL` | はい | GitLab インスタンス URL（例: `https://gitlab.com`） |
| `GITLAB_TOKEN` | はい | `api` スコープを持つ Personal Access Token |
| `GITLAB_MCP_ENABLED_TOOLS` | いいえ | 有効にするツールのカンマ区切りリスト（未設定時は全て有効） |
| `GITLAB_MCP_DISABLED_TOOLS` | いいえ | 無効にするツールのカンマ区切りリスト（ENABLED_TOOLS より優先） |
| `GITLAB_MCP_DEBUG` | いいえ | デバッグログを有効化（`true`、`1`、または `yes`） |

### ツールフィルタリング例

```bash
# 読み取り操作のみを有効化
export GITLAB_MCP_ENABLED_TOOLS="list_merge_requests,get_merge_request,get_merge_request_changes"

# 破壊的操作を無効化
export GITLAB_MCP_DISABLED_TOOLS="merge_merge_request,approve_merge_request"

# DISABLED_TOOLS は ENABLED_TOOLS より優先される
```

## 利用可能なツール

### Merge Request 操作

| ツール | 説明 |
|--------|------|
| `list_merge_requests` | プロジェクトの Merge Request 一覧を取得（フィルタリング対応） |
| `get_merge_request` | 特定の Merge Request の詳細情報を取得 |
| `create_merge_request` | 新しい Merge Request を作成 |
| `update_merge_request` | 既存の Merge Request を更新 |
| `merge_merge_request` | Merge Request をマージ（squash、ブランチ削除オプション対応） |
| `get_merge_request_changes` | Merge Request のファイル変更/差分を取得 |

### ディスカッション・コメント

| ツール | 説明 |
|--------|------|
| `add_merge_request_comment` | Merge Request に一般コメントを追加 |
| `add_merge_request_discussion` | コードの特定行にディスカッションを作成 |
| `list_merge_request_discussions` | Merge Request の全ディスカッションを一覧取得 |
| `resolve_discussion` | ディスカッションを解決済み/未解決に設定 |

### 承認

| ツール | 説明 |
|--------|------|
| `approve_merge_request` | Merge Request を承認 |
| `unapprove_merge_request` | Merge Request の承認を取り消し |
| `get_merge_request_approvals` | 承認状態と承認者一覧を取得 |

### パイプライン・CI/CD

| ツール | 説明 |
|--------|------|
| `list_merge_request_pipelines` | Merge Request に関連するパイプライン一覧を取得 |
| `get_pipeline_jobs` | 特定のパイプラインのジョブ一覧を取得 |

## MCP クライアントでの使用方法

### Claude Code

Claude Code の設定に追加：

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

`claude_desktop_config.json` に追加：

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

**設定ファイルの場所:**
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`
- Linux: `~/.config/Claude/claude_desktop_config.json`

## 使用例

設定後、Claude に以下のような依頼ができます：

- 「mygroup/myproject の全ての open な MR を一覧表示して」
- 「myproject の MR #42 の詳細を取得して」
- 「feature-branch から main への MR を 'Add new feature' というタイトルで作成して」
- 「MR #42 に 'LGTM!' とコメントして」
- 「MR #42 を承認して」
- 「MR #42 のパイプライン状態を確認して」

## 開発

### 前提条件

- Go 1.23 以降
- Git

### ビルド

```bash
go build -o gitlab-mcp ./cmd/gitlab-mcp
```

### テスト

```bash
go test ./...
```

### プロジェクト構造

```
.
├── cmd/gitlab-mcp/        # アプリケーションエントリーポイント
├── internal/
│   ├── config/            # 設定管理
│   ├── gitlab/            # GitLab API クライアントラッパー
│   ├── registry/          # MCP ツールレジストリ
│   └── tools/             # MCP ツール実装
│       ├── approval/      # 承認ツール
│       ├── discussion/    # ディスカッションツール
│       ├── mergerequest/  # Merge Request ツール
│       └── pipeline/      # パイプラインツール
└── test/integration/      # 統合テスト
```

## ライセンス

このプロジェクトは MIT ライセンスの下でライセンスされています - 詳細は [LICENSE](../LICENSE) ファイルを参照してください。

## 謝辞

- [Model Context Protocol](https://modelcontextprotocol.io/) - プロトコル仕様
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - 公式 Go SDK
- [GitLab Go SDK](https://gitlab.com/gitlab-org/api/client-go) - 公式 GitLab API クライアント
