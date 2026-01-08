package gitlab

import (
	"fmt"
	"net/http"

	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

// ErrorCode はエラーコード
type ErrorCode string

const (
	ErrCodeUnauthorized ErrorCode = "unauthorized"
	ErrCodeForbidden    ErrorCode = "forbidden"
	ErrCodeNotFound     ErrorCode = "not_found"
	ErrCodeRateLimited  ErrorCode = "rate_limited"
	ErrCodeBadRequest   ErrorCode = "bad_request"
	ErrCodeServerError  ErrorCode = "server_error"
	ErrCodeToolDisabled ErrorCode = "tool_disabled"
)

// MCPError は MCP 互換エラー
type MCPError struct {
	Code    ErrorCode
	Message string
}

// Error implements the error interface
func (e *MCPError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// IsRetryable はリトライ可能なエラーかどうかを返す
func (e *MCPError) IsRetryable() bool {
	return e.Code == ErrCodeRateLimited || e.Code == ErrCodeServerError
}

// FromGitLabResponse は GitLab SDK レスポンスから MCPError を作成する
func FromGitLabResponse(err error, resp *gogitlab.Response) *MCPError {
	if resp == nil {
		return &MCPError{
			Code:    ErrCodeServerError,
			Message: fmt.Sprintf("GitLab API エラー: %v", err),
		}
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &MCPError{
			Code:    ErrCodeUnauthorized,
			Message: "認証トークンが無効または期限切れです",
		}
	case http.StatusForbidden:
		return &MCPError{
			Code:    ErrCodeForbidden,
			Message: "この操作を実行する権限がありません",
		}
	case http.StatusNotFound:
		return &MCPError{
			Code:    ErrCodeNotFound,
			Message: "指定されたリソースが見つかりません",
		}
	case http.StatusTooManyRequests:
		return &MCPError{
			Code:    ErrCodeRateLimited,
			Message: "API レート制限に達しました。しばらく待ってから再試行してください",
		}
	case http.StatusBadRequest:
		return &MCPError{
			Code:    ErrCodeBadRequest,
			Message: fmt.Sprintf("リクエストが無効です: %v", err),
		}
	default:
		if resp.StatusCode >= 500 {
			return &MCPError{
				Code:    ErrCodeServerError,
				Message: "GitLab サーバーでエラーが発生しました。しばらく待ってから再試行してください",
			}
		}
		return &MCPError{
			Code:    ErrCodeServerError,
			Message: fmt.Sprintf("予期しないエラーが発生しました: %v", err),
		}
	}
}

// NewToolDisabledError はツール無効化エラーを作成する
func NewToolDisabledError(toolName string) *MCPError {
	return &MCPError{
		Code:    ErrCodeToolDisabled,
		Message: fmt.Sprintf("ツール '%s' は無効化されています", toolName),
	}
}
