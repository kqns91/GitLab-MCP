# Implementation Plan

## Task 1: プロジェクト基盤と設定管理

- [ ] 1.1 Go モジュールの初期化と依存関係の設定
  - Go モジュールを作成し、MCP Go SDK と必要な依存関係を追加する
  - プロジェクトのディレクトリ構造を作成する
  - _Requirements: 11.1_

- [ ] 1.2 環境変数から GitLab 接続情報を読み込む設定管理を実装
  - GitLab URL と認証トークンを環境変数から読み込む
  - 必須項目が未設定の場合にエラーを返す
  - 認証トークンがログに出力されないことを保証する
  - デバッグモードの有効/無効を環境変数で制御する
  - _Requirements: 1.1, 1.4, 12.4_

- [ ] 1.3 ツール有効化/無効化の設定ロジックを実装
  - ENABLED_TOOLS 環境変数のパース（カンマ区切り）
  - DISABLED_TOOLS 環境変数のパース（カンマ区切り）
  - DISABLED_TOOLS を ENABLED_TOOLS より優先する判定ロジック
  - ツールが有効かどうかを判定するメソッド
  - _Requirements: 13.1, 13.2, 13.3, 13.6_

## Task 2: GitLab API クライアント基盤

- [ ] 2.1 GitLab Go SDK クライアントラッパーを実装
  - gitlab.com/gitlab-org/api/client-go を使用してクライアントを初期化
  - 環境変数から URL とトークンを受け取りクライアントを構成
  - SDK の各サービス（MergeRequests, Discussions, Pipelines など）へのアクセサを提供
  - _Requirements: 1.2, 1.3, 2.5, 3.2, 3.3, 6.5, 9.2, 10.2_

- [ ] 2.2 エラー型とエラーハンドリングを実装
  - エラーコード（unauthorized, forbidden, not_found, rate_limited など）を定義
  - SDK エラー（*gitlab.ErrorResponse）から内部エラー型への変換
  - レート制限エラー（429）の検出とリトライ可能フラグ
  - ネットワークエラーの検出
  - ツール無効化エラーの定義
  - _Requirements: 1.2, 1.3, 12.1, 12.2, 12.3, 13.4_

## Task 3: ツールレジストリ

- [ ] 3.1 ツールの登録と管理を行うレジストリを実装
  - ツールの登録機能
  - 設定に基づくツールの有効/無効判定
  - 無効化されたツールをツール一覧から除外
  - 無効化されたツール呼び出し時のエラー返却
  - _Requirements: 11.3, 11.4, 13.4, 13.5_

## Task 4: MergeRequest API クライアントメソッド

- [ ] 4.1 (P) MR 一覧取得と詳細取得のラッパーメソッドを実装
  - SDK の MergeRequestsService.ListProjectMergeRequests を呼び出す
  - SDK の ListMergeRequestsOptions でフィルタリング（state, author_id, assignee_id）
  - SDK の MergeRequestsService.GetMergeRequest で詳細情報を取得
  - SDK エラーから内部エラーへの変換
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 3.1, 3.4_

- [ ] 4.2 (P) MR 作成と更新のラッパーメソッドを実装
  - SDK の MergeRequestsService.CreateMergeRequest を呼び出す
  - SDK の CreateMergeRequestOptions でオプション指定（アサイン先、レビュアー、ラベル）
  - SDK の MergeRequestsService.UpdateMergeRequest で既存 MR を更新
  - SDK エラーから内部エラーへの変換
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 5.1, 5.2, 5.3, 5.4_

- [ ] 4.3 (P) MR マージと変更差分取得のラッパーメソッドを実装
  - SDK の MergeRequestsService.AcceptMergeRequest を呼び出す（squash、削除オプション付き）
  - マージ不可状態（conflicts, pipeline failure）のエラー変換
  - SDK の MergeRequestsService.GetMergeRequestDiffVersions で変更差分を取得
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 9.1, 9.2, 9.3_

## Task 5: Discussion API クライアントメソッド

- [ ] 5.1 (P) ディスカッション操作のラッパーメソッドを実装
  - SDK の NotesService.CreateMergeRequestNote でコメントを追加
  - SDK の DiscussionsService.CreateMergeRequestDiscussion で位置情報付きディスカッションを作成
  - SDK の DiscussionsService.ListMergeRequestDiscussions でディスカッション一覧を取得
  - SDK の DiscussionsService.ResolveMergeRequestDiscussion で解決状態を変更
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

## Task 6: Approval API クライアントメソッド

- [ ] 6.1 (P) 承認操作のラッパーメソッドを実装
  - SDK の MergeRequestApprovalsService.ApproveMergeRequest で承認
  - SDK の MergeRequestApprovalsService.UnapproveMergeRequest で承認取消
  - SDK の MergeRequestApprovalsService.GetConfiguration で承認状態・承認者を取得
  - SDK エラーから承認権限エラーへの変換
  - _Requirements: 7.1, 7.2, 7.3, 7.4_

## Task 7: Pipeline API クライアントメソッド

- [ ] 7.1 (P) パイプライン情報取得のラッパーメソッドを実装
  - SDK の MergeRequestsService.ListMergeRequestPipelines で MR 関連パイプラインを取得
  - SDK の JobsService.ListPipelineJobs でジョブ一覧を取得
  - _Requirements: 10.1, 10.2, 10.3_

## Task 8: MCP ツール実装 - MergeRequest

- [ ] 8.1 list_merge_requests ツールを実装
  - MCP ツールとして MR 一覧取得を公開
  - 入力パラメータの JSON Schema を定義
  - フィルタリングオプション（state, author_id, assignee_id）をサポート
  - 各 MR の基本情報を構造化して返却
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 11.4_

- [ ] 8.2 get_merge_request ツールを実装
  - MCP ツールとして MR 詳細取得を公開
  - 変更差分を含めるオプションをサポート
  - MR 不在時の適切なエラーメッセージ
  - _Requirements: 3.1, 3.2, 3.4, 11.4_

- [ ] 8.3 create_merge_request ツールを実装
  - MCP ツールとして MR 作成を公開
  - 必須パラメータ（ソースブランチ、ターゲットブランチ、タイトル、説明）の検証
  - オプションパラメータ（アサイン先、レビュアー、ラベル）のサポート
  - 作成成功時に IID と URL を返却
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 11.4_

- [ ] 8.4 update_merge_request ツールを実装
  - MCP ツールとして MR 更新を公開
  - 更新可能フィールド（タイトル、説明、アサイン先、レビュアー、ラベル、ターゲットブランチ）
  - 更新後の MR 情報を返却
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 11.4_

- [ ] 8.5 merge_merge_request ツールを実装
  - MCP ツールとして MR マージを公開
  - squash オプションとソースブランチ削除オプション
  - マージ不可時のエラーメッセージ、パイプライン失敗時の警告
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 11.4_

- [ ] 8.6 get_merge_request_changes ツールを実装
  - MCP ツールとして変更差分取得を公開
  - ファイルごとの変更タイプ、パス、diff 内容を返却
  - _Requirements: 9.1, 9.2, 9.3, 11.4_

## Task 9: MCP ツール実装 - Discussion

- [ ] 9.1 add_merge_request_comment ツールを実装
  - MCP ツールとして一般コメント追加を公開
  - コメント本文を MR に追加
  - _Requirements: 6.1, 11.4_

- [ ] 9.2 add_merge_request_discussion ツールを実装
  - MCP ツールとして行コメント（ディスカッション）作成を公開
  - ファイルパスと行番号を指定した位置情報付きコメント
  - _Requirements: 6.2, 11.4_

- [ ] 9.3 list_merge_request_discussions ツールを実装
  - MCP ツールとしてディスカッション一覧取得を公開
  - 作成者、内容、解決状態、関連ファイル・行情報を含める
  - _Requirements: 6.3, 6.5, 11.4_

- [ ] 9.4 resolve_discussion ツールを実装
  - MCP ツールとしてディスカッション解決を公開
  - ディスカッションを解決済み/未解決に設定
  - _Requirements: 6.4, 11.4_

## Task 10: MCP ツール実装 - Approval

- [ ] 10.1 approve_merge_request ツールを実装
  - MCP ツールとして MR 承認を公開
  - 承認権限エラーの適切なハンドリング
  - _Requirements: 7.1, 7.4, 11.4_

- [ ] 10.2 unapprove_merge_request ツールを実装
  - MCP ツールとして MR 承認取消を公開
  - _Requirements: 7.2, 11.4_

- [ ] 10.3 get_merge_request_approvals ツールを実装
  - MCP ツールとして承認状態取得を公開
  - 承認状態と承認者一覧を返却
  - _Requirements: 7.3, 11.4_

## Task 11: MCP ツール実装 - Pipeline

- [ ] 11.1 list_merge_request_pipelines ツールを実装
  - MCP ツールとして MR パイプライン一覧取得を公開
  - 各パイプラインの状態、日時、URL を返却
  - _Requirements: 10.1, 10.2, 11.4_

- [ ] 11.2 get_pipeline_jobs ツールを実装
  - MCP ツールとしてパイプラインジョブ一覧取得を公開
  - 各ジョブの名前、ステージ、状態を返却
  - _Requirements: 10.3, 11.4_

## Task 12: MCP サーバー統合

- [ ] 12.1 MCP サーバーのエントリーポイントを実装
  - 設定の読み込みと検証
  - GitLab クライアントの初期化
  - ツールレジストリへの全ツール登録
  - stdio トランスポートでの MCP サーバー起動
  - MCP 仕様に準拠したエラーレスポンス
  - _Requirements: 1.1, 11.1, 11.2, 11.3, 11.5_

## Task 13: ユニットテスト

- [ ] 13.1 (P) 設定管理のテストを実装
  - 環境変数パースのテスト
  - ツール有効化/無効化ロジックのテスト
  - DISABLED_TOOLS 優先ルールのテスト
  - _Requirements: 1.1, 13.1, 13.2, 13.3_

- [ ] 13.2 (P) GitLab API クライアントのテストを実装
  - SDK クライアントラッパーのテスト（httptest でモック GitLab API を構築）
  - SDK エラーから内部エラーへの変換テスト
  - エラーハンドリングのテスト（401, 403, 404, 429, 5xx）
  - _Requirements: 12.1, 12.2, 12.3_

- [ ] 13.3 (P) ツールレジストリのテストを実装
  - ツール登録のテスト
  - 有効/無効判定のテスト
  - 無効化ツールの一覧除外テスト
  - _Requirements: 13.4, 13.5_

## Task 14: 統合テスト

- [ ] 14.1 MCP サーバー全体の統合テストを実装
  - モック GitLab API を使用したエンドツーエンドフロー
  - ツール呼び出しとレスポンス検証
  - エラーケース（認証失敗、リソース不在、レート制限）の動作確認
  - ツール有効化/無効化の動作確認
  - _Requirements: 11.1, 11.2, 11.3, 11.5, 12.1, 12.2, 13.4_
