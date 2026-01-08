# Requirements Document

## Introduction
本ドキュメントは、GitLab の Merge Request (MR) 操作を AI エージェントが実行できるようにするための MCP (Model Context Protocol) サーバーの要件を定義します。このサーバーは、GitLab REST API と連携し、MR の作成、レビュー、承認、マージなどの操作を MCP ツールとして公開します。

## Requirements

### Requirement 1: GitLab 認証・接続管理
**Objective:** As a AI エージェント, I want GitLab インスタンスに安全に接続する, so that API を介して MR 操作を実行できる

#### Acceptance Criteria
1. When MCP サーバーが起動する, the GitLab MCP Server shall 環境変数または設定ファイルから GitLab の URL と認証トークンを読み込む
2. When 認証トークンが無効または期限切れである, the GitLab MCP Server shall 明確なエラーメッセージを返す
3. When GitLab インスタンスに接続できない, the GitLab MCP Server shall 接続エラーの詳細を含むエラーメッセージを返す
4. The GitLab MCP Server shall 認証トークンをログや出力に含めない

### Requirement 2: Merge Request 一覧取得
**Objective:** As a AI エージェント, I want プロジェクトの MR 一覧を取得する, so that レビュー対象や作業対象の MR を把握できる

#### Acceptance Criteria
1. When `list_merge_requests` ツールが呼び出される, the GitLab MCP Server shall 指定されたプロジェクトの MR 一覧を返す
2. When state パラメータが指定される, the GitLab MCP Server shall 指定された状態 (opened, closed, merged, all) の MR のみをフィルタリングして返す
3. When author_id パラメータが指定される, the GitLab MCP Server shall 指定した作成者の MR のみを返す
4. When assignee_id パラメータが指定される, the GitLab MCP Server shall 指定したアサイン先の MR のみを返す
5. The GitLab MCP Server shall 各 MR について、タイトル、説明、状態、作成者、アサイン先、ソースブランチ、ターゲットブランチの情報を含める

### Requirement 3: Merge Request 詳細取得
**Objective:** As a AI エージェント, I want 特定の MR の詳細情報を取得する, so that MR の内容を理解してレビューや操作ができる

#### Acceptance Criteria
1. When `get_merge_request` ツールがプロジェクト ID と MR IID で呼び出される, the GitLab MCP Server shall 該当 MR の詳細情報を返す
2. The GitLab MCP Server shall MR の変更差分 (diff) 情報を取得できるオプションを提供する
3. The GitLab MCP Server shall MR に関連するパイプライン情報を取得できるオプションを提供する
4. If 指定された MR が存在しない, then the GitLab MCP Server shall 「MR が見つからない」というエラーを返す

### Requirement 4: Merge Request 作成
**Objective:** As a AI エージェント, I want 新しい MR を作成する, so that コード変更をレビュー可能な形で提出できる

#### Acceptance Criteria
1. When `create_merge_request` ツールが必要なパラメータで呼び出される, the GitLab MCP Server shall 新しい MR を作成する
2. The GitLab MCP Server shall ソースブランチ、ターゲットブランチ、タイトル、説明を必須パラメータとして受け付ける
3. The GitLab MCP Server shall アサイン先、レビュアー、ラベルをオプションパラメータとして受け付ける
4. If ソースブランチが存在しない, then the GitLab MCP Server shall 「ソースブランチが見つからない」というエラーを返す
5. When MR が正常に作成される, the GitLab MCP Server shall 作成された MR の IID と URL を返す

### Requirement 5: Merge Request 更新
**Objective:** As a AI エージェント, I want 既存の MR を更新する, so that タイトル、説明、アサイン先などを変更できる

#### Acceptance Criteria
1. When `update_merge_request` ツールがプロジェクト ID、MR IID、更新内容で呼び出される, the GitLab MCP Server shall MR を更新する
2. The GitLab MCP Server shall タイトル、説明、アサイン先、レビュアー、ラベル、ターゲットブランチの更新をサポートする
3. If 指定された MR が存在しない, then the GitLab MCP Server shall 「MR が見つからない」というエラーを返す
4. When 更新が正常に完了する, the GitLab MCP Server shall 更新後の MR 情報を返す

### Requirement 6: Merge Request コメント・ディスカッション
**Objective:** As a AI エージェント, I want MR にコメントを追加したりディスカッションを管理する, so that コードレビューのフィードバックを提供できる

#### Acceptance Criteria
1. When `add_merge_request_comment` ツールが呼び出される, the GitLab MCP Server shall MR に一般的なコメントを追加する
2. When `add_merge_request_discussion` ツールが呼び出される, the GitLab MCP Server shall 特定のファイルの特定の行に対するコメント（ディスカッション）を作成する
3. When `list_merge_request_discussions` ツールが呼び出される, the GitLab MCP Server shall MR のディスカッション一覧を返す
4. When `resolve_discussion` ツールが呼び出される, the GitLab MCP Server shall 指定されたディスカッションを解決済みとしてマークする
5. The GitLab MCP Server shall 各ディスカッションについて、作成者、内容、解決状態、関連するファイル・行情報を含める
6. When `delete_merge_request_comment` ツールが呼び出される, the GitLab MCP Server shall 指定されたノート ID のコメントを削除する
7. If 削除対象のコメントが存在しない, then the GitLab MCP Server shall 「コメントが見つからない」というエラーを返す
8. If 削除権限がない, then the GitLab MCP Server shall 「削除権限がありません」というエラーを返す

### Requirement 7: Merge Request 承認
**Objective:** As a AI エージェント, I want MR を承認する, so that コードレビュー完了を示せる

#### Acceptance Criteria
1. When `approve_merge_request` ツールが呼び出される, the GitLab MCP Server shall 指定された MR を承認する
2. When `unapprove_merge_request` ツールが呼び出される, the GitLab MCP Server shall 指定された MR の承認を取り消す
3. When `get_merge_request_approvals` ツールが呼び出される, the GitLab MCP Server shall MR の承認状態と承認者一覧を返す
4. If 承認権限がない, then the GitLab MCP Server shall 「承認権限がありません」というエラーを返す

### Requirement 8: Merge Request マージ
**Objective:** As a AI エージェント, I want MR をマージする, so that レビュー済みのコード変更をメインブランチに統合できる

#### Acceptance Criteria
1. When `merge_merge_request` ツールが呼び出される, the GitLab MCP Server shall 指定された MR をマージする
2. The GitLab MCP Server shall squash オプションをサポートする
3. The GitLab MCP Server shall ソースブランチ削除オプションをサポートする
4. If MR がマージできない状態である, then the GitLab MCP Server shall マージできない理由を含むエラーを返す
5. If パイプラインが失敗している, then the GitLab MCP Server shall 「パイプラインが失敗しています」という警告を返す

### Requirement 9: Merge Request 変更差分取得
**Objective:** As a AI エージェント, I want MR の変更差分を詳細に取得する, so that コード変更内容を正確にレビューできる

#### Acceptance Criteria
1. When `get_merge_request_changes` ツールが呼び出される, the GitLab MCP Server shall MR に含まれるファイル変更の一覧を返す
2. The GitLab MCP Server shall 各変更ファイルについて、旧パス、新パス、変更タイプ（追加、削除、変更、名前変更）、diff 内容を含める
3. When ファイル数が多い場合, the GitLab MCP Server shall ページネーションをサポートする

### Requirement 10: Merge Request パイプライン情報
**Objective:** As a AI エージェント, I want MR のパイプライン状態を確認する, so that CI/CD の結果に基づいて判断できる

#### Acceptance Criteria
1. When `list_merge_request_pipelines` ツールが呼び出される, the GitLab MCP Server shall MR に関連するパイプライン一覧を返す
2. The GitLab MCP Server shall 各パイプラインについて、状態、作成日時、完了日時、URL を含める
3. When `get_pipeline_jobs` ツールが呼び出される, the GitLab MCP Server shall パイプラインのジョブ一覧と各ジョブの状態を返す

### Requirement 11: MCP サーバー基盤
**Objective:** As a 開発者, I want 標準的な MCP プロトコルに準拠したサーバーを構築する, so that 任意の MCP クライアントから利用できる

#### Acceptance Criteria
1. The GitLab MCP Server shall MCP (Model Context Protocol) 仕様に準拠する
2. The GitLab MCP Server shall stdio トランスポートをサポートする
3. The GitLab MCP Server shall 利用可能なツールの一覧と説明を提供する
4. The GitLab MCP Server shall 各ツールの入力スキーマを JSON Schema 形式で提供する
5. When ツール呼び出しでエラーが発生する, the GitLab MCP Server shall MCP 仕様に準拠したエラーレスポンスを返す

### Requirement 12: エラーハンドリングとロギング
**Objective:** As a 開発者, I want 適切なエラーハンドリングとロギングを行う, so that 問題の診断とデバッグが容易になる

#### Acceptance Criteria
1. The GitLab MCP Server shall すべての GitLab API エラーを適切にハンドリングする
2. The GitLab MCP Server shall レート制限エラー (429) を検出し、適切なメッセージを返す
3. The GitLab MCP Server shall ネットワークエラーを検出し、リトライ可能かどうかを示す
4. While デバッグモードが有効である, the GitLab MCP Server shall 詳細なログを stderr に出力する

### Requirement 13: ツール有効化・無効化設定
**Objective:** As a 運用者, I want 環境変数でツールの有効・無効を制御する, so that セキュリティポリシーや運用要件に応じて利用可能なツールを制限できる

#### Acceptance Criteria
1. When `GITLAB_MCP_ENABLED_TOOLS` 環境変数が設定される, the GitLab MCP Server shall 指定されたツールのみを有効化する
2. When `GITLAB_MCP_DISABLED_TOOLS` 環境変数が設定される, the GitLab MCP Server shall 指定されたツールを無効化する
3. If 両方の環境変数が設定される, then the GitLab MCP Server shall `DISABLED_TOOLS` を優先し、`ENABLED_TOOLS` に含まれていても無効化する
4. When 無効化されたツールが呼び出される, the GitLab MCP Server shall 「このツールは無効化されています」というエラーを返す
5. The GitLab MCP Server shall ツール一覧取得時に無効化されたツールを含めない
6. The GitLab MCP Server shall 環境変数の値をカンマ区切りのツール名リストとして解釈する

### Requirement 14: ドキュメント
**Objective:** As a OSS 利用者, I want 包括的なドキュメントを参照する, so that セットアップと利用方法を理解できる

#### Acceptance Criteria
1. The GitLab MCP Server shall README.md に以下の内容を英語で記載する: プロジェクト概要、機能一覧、インストール方法、設定方法、使用例、ライセンス
2. The GitLab MCP Server shall README.md に Quick Start セクションを含め、5 分以内にセットアップできる手順を提供する
3. The GitLab MCP Server shall README.md に利用可能な全ツールの一覧と各ツールの説明を含める
4. The GitLab MCP Server shall README.md に環境変数の完全なリファレンスを含める
5. The GitLab MCP Server shall README.md に Claude Code / Claude Desktop での設定例を含める
6. The GitLab MCP Server shall README.md に Contributing ガイドラインを含める
7. The GitLab MCP Server shall docs/README.ja.md に日本語版ドキュメントを提供する
8. The GitLab MCP Server shall README.md に badges（Go version, License, Release）を含める
