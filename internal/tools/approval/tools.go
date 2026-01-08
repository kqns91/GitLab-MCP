package approval

import (
	"context"

	"github.com/kqns91/gitlab-mcp/internal/gitlab"
	"github.com/kqns91/gitlab-mcp/internal/registry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ApproveInput は approve_merge_request の入力パラメータ
type ApproveInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
}

// ApproveOutput は approve_merge_request の出力
type ApproveOutput struct {
	Approved        bool  `json:"approved"`
	UserHasApproved bool  `json:"user_has_approved"`
	ApprovalsLeft   int64 `json:"approvals_left"`
}

// UnapproveInput は unapprove_merge_request の入力パラメータ
type UnapproveInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
}

// UnapproveOutput は unapprove_merge_request の出力
type UnapproveOutput struct {
	Success bool `json:"success"`
}

// GetApprovalsInput は get_merge_request_approvals の入力パラメータ
type GetApprovalsInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
}

// Approver は承認者情報
type Approver struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// GetApprovalsOutput は get_merge_request_approvals の出力
type GetApprovalsOutput struct {
	Approved          bool       `json:"approved"`
	ApprovalsRequired int64      `json:"approvals_required"`
	ApprovalsLeft     int64      `json:"approvals_left"`
	UserHasApproved   bool       `json:"user_has_approved"`
	UserCanApprove    bool       `json:"user_can_approve"`
	ApprovedBy        []Approver `json:"approved_by"`
}

// clientHolder holds the GitLab client for handlers
type clientHolder struct {
	client *gitlab.Client
}

var holder *clientHolder

// Register は承認関連ツールを登録する
func Register(reg *registry.Registry, client *gitlab.Client) {
	holder = &clientHolder{client: client}

	registry.RegisterTool(reg, "approve_merge_request",
		"GitLab Merge Request を承認します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ApproveInput) (*mcp.CallToolResult, ApproveOutput, error) {
			return approveHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "unapprove_merge_request",
		"GitLab Merge Request の承認を取り消します",
		func(ctx context.Context, req *mcp.CallToolRequest, input UnapproveInput) (*mcp.CallToolResult, UnapproveOutput, error) {
			return unapproveHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "get_merge_request_approvals",
		"GitLab Merge Request の承認状態を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input GetApprovalsInput) (*mcp.CallToolResult, GetApprovalsOutput, error) {
			return getApprovalsHandler(holder.client, ctx, req, input)
		})
}

func approveHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ApproveInput) (*mcp.CallToolResult, ApproveOutput, error) {
	approvals, err := client.ApproveMergeRequest(input.ProjectID, input.MergeRequestIID)
	if err != nil {
		return nil, ApproveOutput{}, err
	}

	return nil, ApproveOutput{
		Approved:        approvals.Approved,
		UserHasApproved: approvals.UserHasApproved,
		ApprovalsLeft:   approvals.ApprovalsLeft,
	}, nil
}

func unapproveHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input UnapproveInput) (*mcp.CallToolResult, UnapproveOutput, error) {
	err := client.UnapproveMergeRequest(input.ProjectID, input.MergeRequestIID)
	if err != nil {
		return nil, UnapproveOutput{}, err
	}

	return nil, UnapproveOutput{Success: true}, nil
}

func getApprovalsHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input GetApprovalsInput) (*mcp.CallToolResult, GetApprovalsOutput, error) {
	approvals, err := client.GetMergeRequestApprovals(input.ProjectID, input.MergeRequestIID)
	if err != nil {
		return nil, GetApprovalsOutput{}, err
	}

	approvers := make([]Approver, len(approvals.ApprovedBy))
	for i, a := range approvals.ApprovedBy {
		approvers[i] = Approver{
			ID:       int64(a.User.ID),
			Username: a.User.Username,
		}
	}

	return nil, GetApprovalsOutput{
		Approved:          approvals.Approved,
		ApprovalsRequired: approvals.ApprovalsRequired,
		ApprovalsLeft:     approvals.ApprovalsLeft,
		UserHasApproved:   approvals.UserHasApproved,
		UserCanApprove:    approvals.UserCanApprove,
		ApprovedBy:        approvers,
	}, nil
}
