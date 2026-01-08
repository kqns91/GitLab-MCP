package gitlab

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

func TestErrorCode_Constants(t *testing.T) {
	// Verify error codes exist
	assert.Equal(t, ErrorCode("unauthorized"), ErrCodeUnauthorized)
	assert.Equal(t, ErrorCode("forbidden"), ErrCodeForbidden)
	assert.Equal(t, ErrorCode("not_found"), ErrCodeNotFound)
	assert.Equal(t, ErrorCode("rate_limited"), ErrCodeRateLimited)
	assert.Equal(t, ErrorCode("bad_request"), ErrCodeBadRequest)
	assert.Equal(t, ErrorCode("server_error"), ErrCodeServerError)
	assert.Equal(t, ErrorCode("tool_disabled"), ErrCodeToolDisabled)
}

func TestMCPError_Error(t *testing.T) {
	err := &MCPError{
		Code:    ErrCodeNotFound,
		Message: "Resource not found",
	}

	assert.Equal(t, "[not_found] Resource not found", err.Error())
}

func TestMCPError_IsRetryable(t *testing.T) {
	tests := []struct {
		code      ErrorCode
		retryable bool
	}{
		{ErrCodeRateLimited, true},
		{ErrCodeServerError, true},
		{ErrCodeUnauthorized, false},
		{ErrCodeForbidden, false},
		{ErrCodeNotFound, false},
		{ErrCodeBadRequest, false},
		{ErrCodeToolDisabled, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			err := &MCPError{Code: tt.code}
			assert.Equal(t, tt.retryable, err.IsRetryable())
		})
	}
}

func TestFromGitLabResponse_401(t *testing.T) {
	resp := &gogitlab.Response{
		Response: &http.Response{StatusCode: http.StatusUnauthorized},
	}

	mcpErr := FromGitLabResponse(errors.New("unauthorized"), resp)

	assert.Equal(t, ErrCodeUnauthorized, mcpErr.Code)
	assert.Contains(t, mcpErr.Message, "認証")
}

func TestFromGitLabResponse_403(t *testing.T) {
	resp := &gogitlab.Response{
		Response: &http.Response{StatusCode: http.StatusForbidden},
	}

	mcpErr := FromGitLabResponse(errors.New("forbidden"), resp)

	assert.Equal(t, ErrCodeForbidden, mcpErr.Code)
	assert.Contains(t, mcpErr.Message, "権限")
}

func TestFromGitLabResponse_404(t *testing.T) {
	resp := &gogitlab.Response{
		Response: &http.Response{StatusCode: http.StatusNotFound},
	}

	mcpErr := FromGitLabResponse(errors.New("not found"), resp)

	assert.Equal(t, ErrCodeNotFound, mcpErr.Code)
	assert.Contains(t, mcpErr.Message, "見つかりません")
}

func TestFromGitLabResponse_429(t *testing.T) {
	resp := &gogitlab.Response{
		Response: &http.Response{StatusCode: http.StatusTooManyRequests},
	}

	mcpErr := FromGitLabResponse(errors.New("rate limited"), resp)

	assert.Equal(t, ErrCodeRateLimited, mcpErr.Code)
	assert.Contains(t, mcpErr.Message, "レート制限")
	assert.True(t, mcpErr.IsRetryable())
}

func TestFromGitLabResponse_500(t *testing.T) {
	resp := &gogitlab.Response{
		Response: &http.Response{StatusCode: http.StatusInternalServerError},
	}

	mcpErr := FromGitLabResponse(errors.New("server error"), resp)

	assert.Equal(t, ErrCodeServerError, mcpErr.Code)
	assert.Contains(t, mcpErr.Message, "サーバー")
	assert.True(t, mcpErr.IsRetryable())
}

func TestFromGitLabResponse_NilResponse(t *testing.T) {
	mcpErr := FromGitLabResponse(errors.New("network error"), nil)

	assert.Equal(t, ErrCodeServerError, mcpErr.Code)
	assert.Contains(t, mcpErr.Message, "network error")
}

func TestNewToolDisabledError(t *testing.T) {
	err := NewToolDisabledError("merge_merge_request")

	assert.Equal(t, ErrCodeToolDisabled, err.Code)
	assert.Contains(t, err.Message, "merge_merge_request")
	assert.Contains(t, err.Message, "無効")
}
