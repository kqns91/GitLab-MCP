# Research & Design Decisions

## Summary
- **Feature**: `gitlab-mr-mcp`
- **Discovery Scope**: New Feature (greenfield)
- **Key Findings**:
  - 公式 MCP Go SDK (`github.com/modelcontextprotocol/go-sdk`) が利用可能（v1.2.0, 2025-12-22 リリース）
  - Go 1.25+ で最新機能（Green Tea GC、container-aware GOMAXPROCS）を活用可能
  - GitLab Go SDK (`gitlab.com/gitlab-org/api/client-go`) v1.0+ を使用（公式 SDK）
  - GitLab REST API v4 は MR 操作に必要な全エンドポイントを提供

## Research Log

### MCP Go SDK 調査
- **Context**: MCP サーバー実装の技術選定（Go 言語）
- **Sources Consulted**:
  - [MCP Go SDK (Official)](https://github.com/modelcontextprotocol/go-sdk) — Google と Anthropic が共同メンテナンス
  - [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) — コミュニティ版、400+ パッケージで使用
- **Findings**:
  - **公式 SDK**: `github.com/modelcontextprotocol/go-sdk` (v1.2.0)
    - MCP spec 2025-11-25 をサポート
    - `mcp.Server` でサーバー構築
    - `mcp.StdioTransport{}` で stdio トランスポート
    - `mcp.AddTool()` でツール登録、struct タグで JSON Schema 生成
  - **コミュニティ SDK**: `github.com/mark3labs/mcp-go`
    - より高レベルな API を提供
    - 400+ パッケージで採用実績あり
- **Implications**:
  - 公式 SDK を採用（長期サポートと安定性を優先）
  - struct タグによるスキーマ定義で型安全性を確保

### Go バージョン調査
- **Context**: 使用する Go バージョンの決定
- **Sources Consulted**:
  - [Go Release History](https://go.dev/doc/devel/release)
  - [Go 1.25 Release Notes](https://go.dev/doc/go1.25)
- **Findings**:
  - **Go 1.25.5** (2025-12-02): 最新安定版
  - **Go 1.26 RC1** (2025-12-16): 2026年2月リリース予定
  - Go 1.25 の主な機能:
    - Green Tea GC（10-40% GC オーバーヘッド削減）
    - container-aware GOMAXPROCS
  - Go 1.24 の主な機能:
    - ジェネリック型エイリアス完全サポート
    - tool directives in go.mod
    - FIPS 140-3 準拠サポート
- **Implications**:
  - Go 1.25.5 を使用（最新安定版）
  - go.mod で `go 1.25` を指定

### GitLab REST API v4 調査
- **Context**: MR 操作に必要な API エンドポイントの特定
- **Sources Consulted**:
  - [GitLab Merge Requests API](https://docs.gitlab.com/api/merge_requests/)
  - [GitLab Discussions API](https://docs.gitlab.com/api/discussions/)
  - [GitLab Merge Request Approvals API](https://docs.gitlab.com/api/merge_request_approvals/)
  - [GitLab Notes API](https://docs.gitlab.com/api/notes/)
- **Findings**:
  - **MR 一覧**: `GET /projects/:id/merge_requests` (state, author_id, assignee_id でフィルタ可能)
  - **MR 詳細**: `GET /projects/:id/merge_requests/:merge_request_iid`
  - **MR 作成**: `POST /projects/:id/merge_requests` (source_branch, target_branch, title 必須)
  - **MR 更新**: `PUT /projects/:id/merge_requests/:merge_request_iid`
  - **MR マージ**: `PUT /projects/:id/merge_requests/:merge_request_iid/merge`
  - **MR 承認**: `POST /projects/:id/merge_requests/:merge_request_iid/approve`
  - **MR 承認取消**: `POST /projects/:id/merge_requests/:merge_request_iid/unapprove`
  - **承認状態**: `GET /projects/:id/merge_requests/:merge_request_iid/approvals`
  - **変更差分**: `GET /projects/:id/merge_requests/:merge_request_iid/diffs`
  - **パイプライン**: `GET /projects/:id/merge_requests/:merge_request_iid/pipelines`
  - **ディスカッション一覧**: `GET /projects/:id/merge_requests/:merge_request_iid/discussions`
  - **スレッド作成**: `POST /projects/:id/merge_requests/:merge_request_iid/discussions`
  - **ディスカッション解決**: `PUT /projects/:id/merge_requests/:merge_request_iid/discussions/:discussion_id`
  - **ノート追加**: `POST /projects/:id/merge_requests/:merge_request_iid/notes`
- **Implications**:
  - すべての要件を満たす API が存在
  - 認証は Private-Token ヘッダーまたは OAuth2 Bearer トークン
  - レート制限あり（429 エラー対応必要）

### 認証方式調査
- **Context**: GitLab API 認証方法の選定
- **Sources Consulted**: GitLab API ドキュメント
- **Findings**:
  - Personal Access Token (PAT): `Private-Token` ヘッダー
  - OAuth2 Access Token: `Authorization: Bearer` ヘッダー
  - Project/Group Access Token: PAT と同様の使用方法
- **Implications**:
  - 環境変数 `GITLAB_TOKEN` で PAT を設定
  - `GITLAB_URL` で GitLab インスタンス URL を設定

### GitLab Go SDK 調査
- **Context**: GitLab API 通信ライブラリの選定
- **Sources Consulted**:
  - [gitlab.com/gitlab-org/api/client-go](https://gitlab.com/gitlab-org/api/client-go) — 公式 GitLab Go SDK
  - [github.com/xanzy/go-gitlab](https://github.com/xanzy/go-gitlab) — 旧リポジトリ（非推奨）
  - [pkg.go.dev/gitlab.com/gitlab-org/api/client-go](https://pkg.go.dev/gitlab.com/gitlab-org/api/client-go) — Go パッケージドキュメント
- **Findings**:
  - **公式 SDK**: `gitlab.com/gitlab-org/api/client-go` (v1.0+)
    - `github.com/xanzy/go-gitlab` から移行した公式 SDK
    - GitLab が公式メンテナンス
    - すべての GitLab API エンドポイントをカバー
    - 型定義、エラーハンドリング、ページネーション組み込み
  - **旧リポジトリ**: `github.com/xanzy/go-gitlab` (v0.115.0)
    - 非推奨、v0.115.0 で deprecation notice
    - 新規プロジェクトでは使用しない
- **Implications**:
  - `gitlab.com/gitlab-org/api/client-go` を採用
  - 自前 HTTP クライアント実装は不要
  - 開発効率とメンテナンス性が向上

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| Layered Architecture | Handler → Service → Client | シンプル、Go らしい | 大規模化時に複雑化 | MCP サーバーの規模に適切 |
| Clean Architecture | Domain を中心にした依存逆転 | テスト容易、依存逆転 | オーバーエンジニアリングのリスク | 将来の拡張性は高い |
| **Package-based Modular** | 機能単位のパッケージ分割 | Go 標準的な構造、テスト容易、適切な複雑度 | パッケージ間依存管理 | **選択**: 各機能を独立パッケージとして実装 |

## Design Decisions

### Decision: 公式 MCP Go SDK の採用
- **Context**: MCP サーバー実装に使用する SDK の選定
- **Alternatives Considered**:
  1. mark3labs/mcp-go — コミュニティ版、高レベル API
  2. 公式 go-sdk — Google/Anthropic 共同メンテナンス
  3. 自前実装 — プロトコル直接実装
- **Selected Approach**: 公式 `github.com/modelcontextprotocol/go-sdk`
- **Rationale**:
  - 長期サポートが期待できる（Google + Anthropic）
  - MCP spec の最新版に追従
  - 標準的な Go パターンに準拠
- **Trade-offs**: コミュニティ版より API が低レベル
- **Follow-up**: SDK のアップデートを定期的に確認

### Decision: Go 1.25.5 を使用
- **Context**: 対応する Go バージョンの決定
- **Alternatives Considered**:
  1. Go 1.21+ — 広い互換性
  2. Go 1.23+ — 最新機能とのバランス
  3. Go 1.25.5 — 最新安定版
- **Selected Approach**: Go 1.25.5（最新安定版）
- **Rationale**:
  - 最新の安定版を使用してプロジェクトを開始
  - Green Tea GC による 10-40% の GC オーバーヘッド削減
  - container-aware GOMAXPROCS の活用
- **Trade-offs**: 古い環境での互換性は考慮しない
- **Follow-up**: go.mod に `go 1.25` を指定

### Decision: Package-based Modular アーキテクチャの採用
- **Context**: 15 の MR 操作ツールを提供する MCP サーバーの構造設計
- **Alternatives Considered**:
  1. 単一パッケージ — シンプルだが保守性低下
  2. Clean Architecture — 過度な抽象化
  3. Package-based Modular — 機能単位のパッケージ分割
- **Selected Approach**: Package-based Modular
  - `internal/gitlab/` 配下に GitLab API クライアント
  - `internal/tools/` 配下にツール群
  - `internal/config/` 配下に設定管理
- **Rationale**: Go の標準的なプロジェクト構造に準拠、各パッケージが独立してテスト可能
- **Trade-offs**: パッケージ間の共通処理は internal/shared に抽出
- **Follow-up**: ツール間で共通するエラーハンドリングを共通化

### Decision: 環境変数ベースの設定管理
- **Context**: GitLab 接続情報とツール有効化設定
- **Alternatives Considered**:
  1. 設定ファイル (JSON/YAML) — 柔軟だが管理複雑
  2. 環境変数のみ — シンプル、CI/CD 親和性高
  3. 環境変数 + 設定ファイル — 柔軟性と簡便性の両立
- **Selected Approach**: 環境変数のみ
  - `GITLAB_URL`: GitLab インスタンス URL
  - `GITLAB_TOKEN`: Personal Access Token
  - `GITLAB_MCP_ENABLED_TOOLS`: 有効ツール（カンマ区切り）
  - `GITLAB_MCP_DISABLED_TOOLS`: 無効ツール（カンマ区切り）
  - `GITLAB_MCP_DEBUG`: デバッグモード
- **Rationale**: MCP サーバーは通常 CLI 経由で起動、環境変数が最適
- **Trade-offs**: 複雑な設定には不向き（現状は不要）
- **Follow-up**: なし

### Decision: ツール有効/無効の優先順位
- **Context**: ENABLED_TOOLS と DISABLED_TOOLS の競合時の動作
- **Alternatives Considered**:
  1. ENABLED_TOOLS 優先 — ホワイトリスト方式
  2. DISABLED_TOOLS 優先 — ブラックリスト方式（セキュリティ重視）
- **Selected Approach**: DISABLED_TOOLS 優先
- **Rationale**: セキュリティ上、明示的な無効化を優先すべき
- **Trade-offs**: ユーザーの意図と異なる場合あり（ドキュメントで明記）
- **Follow-up**: README に動作仕様を明記

### Decision: GitLab Go SDK の採用
- **Context**: GitLab API 通信方法の選定
- **Alternatives Considered**:
  1. net/http 直接使用 — 依存関係最小、完全な制御
  2. gitlab.com/gitlab-org/api/client-go — 公式 SDK、型定義・エラー処理済み
  3. github.com/xanzy/go-gitlab — コミュニティ版（非推奨）
- **Selected Approach**: `gitlab.com/gitlab-org/api/client-go` v1.0+
- **Rationale**:
  - 公式メンテナンスで長期サポート
  - 型定義、エラーハンドリング、ページネーションが組み込み済み
  - MCP ツールの設計に集中できる（API 実装に時間を費やさない）
- **Trade-offs**: 依存関係が増加（許容範囲）
- **Follow-up**: go.mod に追加

## Risks & Mitigations
- **Risk 1**: GitLab API レート制限によるエラー
  - **Mitigation**: SDK のエラーハンドリングを活用、リトライ推奨メッセージ
- **Risk 2**: 認証トークンの漏洩
  - **Mitigation**: ログ出力時のマスキング、環境変数での管理
- **Risk 3**: GitLab Go SDK の破壊的変更
  - **Mitigation**: go.mod でバージョン固定、定期的なアップデート確認
- **Risk 4**: MCP Go SDK の仕様変更
  - **Mitigation**: go.mod でバージョン固定、定期的なアップデート確認

## References
- [MCP Go SDK (Official)](https://github.com/modelcontextprotocol/go-sdk) — 公式 MCP Go SDK
- [GitLab Go SDK (Official)](https://gitlab.com/gitlab-org/api/client-go) — 公式 GitLab Go SDK
- [GitLab Go SDK (pkg.go.dev)](https://pkg.go.dev/gitlab.com/gitlab-org/api/client-go) — Go パッケージドキュメント
- [Go Release History](https://go.dev/doc/devel/release) — Go バージョン履歴
- [GitLab REST API Documentation](https://docs.gitlab.com/api/) — GitLab API 公式ドキュメント
- [GitLab Merge Requests API](https://docs.gitlab.com/api/merge_requests/) — MR 操作 API
- [GitLab Discussions API](https://docs.gitlab.com/api/discussions/) — ディスカッション API
- [GitLab Merge Request Approvals API](https://docs.gitlab.com/api/merge_request_approvals/) — 承認 API
