# Product Overview

GitLab MCP Server は、GitLab の Merge Request (MR) 操作を AI エージェントが実行できるようにする MCP (Model Context Protocol) サーバーです。

## Core Capabilities

- **MR ライフサイクル管理**: MR の作成、更新、マージ、クローズの一連の操作
- **コードレビュー支援**: ディスカッション、コメント、承認操作による AI アシストレビュー
- **CI/CD 状態確認**: パイプラインとジョブの状態をリアルタイムで取得
- **変更差分分析**: MR に含まれるファイル変更と diff の詳細取得
- **柔軟なアクセス制御**: 環境変数によるツールの有効化/無効化

## Target Use Cases

- AI エージェント（Claude Code, Claude Desktop など）による MR 操作の自動化
- コードレビューワークフローの効率化
- MR 作成からマージまでの一連のプロセスをプログラマティックに実行

## Value Proposition

- **MCP 準拠**: 標準プロトコルにより任意の MCP クライアントから利用可能
- **セキュア**: Personal Access Token による認証、トークンのログ出力防止
- **運用柔軟性**: ツール単位での有効化/無効化によるセキュリティポリシー対応
- **Go 実装**: 型安全で高性能、シングルバイナリでの配布

---
_Focus on patterns and purpose, not exhaustive feature lists_
