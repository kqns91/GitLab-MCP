package mergerequest

import (
	"context"

	"github.com/kqns91/gitlab-mcp/internal/gitlab"
	"github.com/kqns91/gitlab-mcp/internal/registry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListMergeRequestsInput は list_merge_requests の入力パラメータ
type ListMergeRequestsInput struct {
	ProjectID  string  `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	State      *string `json:"state,omitempty" jsonschema:"enum:opened,enum:closed,enum:merged,enum:all,description:MR state filter"`
	AuthorID   *int    `json:"author_id,omitempty" jsonschema:"description:Author user ID filter"`
	AssigneeID *int    `json:"assignee_id,omitempty" jsonschema:"description:Assignee user ID filter"`
}

// MergeRequestSummary はMR一覧の各項目
type MergeRequestSummary struct {
	IID          int64  `json:"iid"`
	Title        string `json:"title"`
	State        string `json:"state"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	WebURL       string `json:"web_url"`
	AuthorName   string `json:"author_name,omitempty"`
}

// ListMergeRequestsOutput は list_merge_requests の出力
type ListMergeRequestsOutput struct {
	MergeRequests []MergeRequestSummary `json:"merge_requests"`
}

// GetMergeRequestInput は get_merge_request の入力パラメータ
type GetMergeRequestInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
}

// MergeRequestDetail はMR詳細情報
type MergeRequestDetail struct {
	IID          int64  `json:"iid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	State        string `json:"state"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	WebURL       string `json:"web_url"`
	AuthorName   string `json:"author_name,omitempty"`
}

// GetMergeRequestOutput は get_merge_request の出力
type GetMergeRequestOutput = MergeRequestDetail

// CreateMergeRequestInput は create_merge_request の入力パラメータ
type CreateMergeRequestInput struct {
	ProjectID    string   `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	SourceBranch string   `json:"source_branch" jsonschema:"description:Source branch name"`
	TargetBranch string   `json:"target_branch" jsonschema:"description:Target branch name"`
	Title        string   `json:"title" jsonschema:"description:Merge request title"`
	Description  *string  `json:"description,omitempty" jsonschema:"description:Merge request description"`
	AssigneeIDs  []int    `json:"assignee_ids,omitempty" jsonschema:"description:Assignee user IDs"`
	ReviewerIDs  []int    `json:"reviewer_ids,omitempty" jsonschema:"description:Reviewer user IDs"`
	Labels       []string `json:"labels,omitempty" jsonschema:"description:Labels to add"`
}

// CreateMergeRequestOutput は create_merge_request の出力
type CreateMergeRequestOutput struct {
	IID    int64  `json:"iid"`
	Title  string `json:"title"`
	WebURL string `json:"web_url"`
}

// UpdateMergeRequestInput は update_merge_request の入力パラメータ
type UpdateMergeRequestInput struct {
	ProjectID       string   `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int      `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
	Title           *string  `json:"title,omitempty" jsonschema:"description:New title"`
	Description     *string  `json:"description,omitempty" jsonschema:"description:New description"`
	AssigneeIDs     []int    `json:"assignee_ids,omitempty" jsonschema:"description:New assignee user IDs"`
	ReviewerIDs     []int    `json:"reviewer_ids,omitempty" jsonschema:"description:New reviewer user IDs"`
	Labels          []string `json:"labels,omitempty" jsonschema:"description:New labels"`
	TargetBranch    *string  `json:"target_branch,omitempty" jsonschema:"description:New target branch"`
}

// UpdateMergeRequestOutput は update_merge_request の出力
type UpdateMergeRequestOutput struct {
	IID    int64  `json:"iid"`
	Title  string `json:"title"`
	WebURL string `json:"web_url"`
}

// MergeMergeRequestInput は merge_merge_request の入力パラメータ
type MergeMergeRequestInput struct {
	ProjectID                string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID          int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
	Squash                   *bool  `json:"squash,omitempty" jsonschema:"description:Squash commits when merging"`
	ShouldRemoveSourceBranch *bool  `json:"should_remove_source_branch,omitempty" jsonschema:"description:Remove source branch after merge"`
}

// MergeMergeRequestOutput は merge_merge_request の出力
type MergeMergeRequestOutput struct {
	IID    int64  `json:"iid"`
	State  string `json:"state"`
	WebURL string `json:"web_url"`
}

// GetMergeRequestChangesInput は get_merge_request_changes の入力パラメータ
type GetMergeRequestChangesInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
}

// ChangeInfo は変更ファイルの情報
type ChangeInfo struct {
	OldPath     string `json:"old_path"`
	NewPath     string `json:"new_path"`
	Diff        string `json:"diff"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
}

// GetMergeRequestChangesOutput は get_merge_request_changes の出力
type GetMergeRequestChangesOutput struct {
	Changes []ChangeInfo `json:"changes"`
}

// clientHolder holds the GitLab client for handlers
type clientHolder struct {
	client *gitlab.Client
}

var holder *clientHolder

// Register は MR 関連ツールを登録する
func Register(reg *registry.Registry, client *gitlab.Client) {
	holder = &clientHolder{client: client}

	registry.RegisterTool(reg, "list_merge_requests",
		"GitLab プロジェクトの Merge Request 一覧を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ListMergeRequestsInput) (*mcp.CallToolResult, ListMergeRequestsOutput, error) {
			return listMergeRequestsHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "get_merge_request",
		"GitLab Merge Request の詳細情報を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input GetMergeRequestInput) (*mcp.CallToolResult, GetMergeRequestOutput, error) {
			return getMergeRequestHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "create_merge_request",
		"GitLab に新しい Merge Request を作成します",
		func(ctx context.Context, req *mcp.CallToolRequest, input CreateMergeRequestInput) (*mcp.CallToolResult, CreateMergeRequestOutput, error) {
			return createMergeRequestHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "update_merge_request",
		"GitLab Merge Request を更新します",
		func(ctx context.Context, req *mcp.CallToolRequest, input UpdateMergeRequestInput) (*mcp.CallToolResult, UpdateMergeRequestOutput, error) {
			return updateMergeRequestHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "merge_merge_request",
		"GitLab Merge Request をマージします",
		func(ctx context.Context, req *mcp.CallToolRequest, input MergeMergeRequestInput) (*mcp.CallToolResult, MergeMergeRequestOutput, error) {
			return mergeMergeRequestHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "get_merge_request_changes",
		"GitLab Merge Request の変更差分を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input GetMergeRequestChangesInput) (*mcp.CallToolResult, GetMergeRequestChangesOutput, error) {
			return getMergeRequestChangesHandler(holder.client, ctx, req, input)
		})
}

func listMergeRequestsHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ListMergeRequestsInput) (*mcp.CallToolResult, ListMergeRequestsOutput, error) {
	opts := &gitlab.ListMergeRequestsOptions{
		State:      input.State,
		AuthorID:   input.AuthorID,
		AssigneeID: input.AssigneeID,
	}

	mrs, err := client.ListMergeRequests(input.ProjectID, opts)
	if err != nil {
		return nil, ListMergeRequestsOutput{}, err
	}

	summaries := make([]MergeRequestSummary, len(mrs))
	for i, mr := range mrs {
		authorName := ""
		if mr.Author != nil {
			authorName = mr.Author.Username
		}
		summaries[i] = MergeRequestSummary{
			IID:          mr.IID,
			Title:        mr.Title,
			State:        mr.State,
			SourceBranch: mr.SourceBranch,
			TargetBranch: mr.TargetBranch,
			WebURL:       mr.WebURL,
			AuthorName:   authorName,
		}
	}

	return nil, ListMergeRequestsOutput{MergeRequests: summaries}, nil
}

func getMergeRequestHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input GetMergeRequestInput) (*mcp.CallToolResult, GetMergeRequestOutput, error) {
	mr, err := client.GetMergeRequest(input.ProjectID, input.MergeRequestIID)
	if err != nil {
		return nil, GetMergeRequestOutput{}, err
	}

	authorName := ""
	if mr.Author != nil {
		authorName = mr.Author.Username
	}

	return nil, GetMergeRequestOutput{
		IID:          mr.IID,
		Title:        mr.Title,
		Description:  mr.Description,
		State:        mr.State,
		SourceBranch: mr.SourceBranch,
		TargetBranch: mr.TargetBranch,
		WebURL:       mr.WebURL,
		AuthorName:   authorName,
	}, nil
}

func createMergeRequestHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input CreateMergeRequestInput) (*mcp.CallToolResult, CreateMergeRequestOutput, error) {
	opts := &gitlab.CreateMergeRequestOptions{
		SourceBranch: input.SourceBranch,
		TargetBranch: input.TargetBranch,
		Title:        input.Title,
		Description:  input.Description,
		AssigneeIDs:  input.AssigneeIDs,
		ReviewerIDs:  input.ReviewerIDs,
		Labels:       input.Labels,
	}

	mr, err := client.CreateMergeRequest(input.ProjectID, opts)
	if err != nil {
		return nil, CreateMergeRequestOutput{}, err
	}

	return nil, CreateMergeRequestOutput{
		IID:    mr.IID,
		Title:  mr.Title,
		WebURL: mr.WebURL,
	}, nil
}

func updateMergeRequestHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input UpdateMergeRequestInput) (*mcp.CallToolResult, UpdateMergeRequestOutput, error) {
	opts := &gitlab.UpdateMergeRequestOptions{
		Title:        input.Title,
		Description:  input.Description,
		AssigneeIDs:  input.AssigneeIDs,
		ReviewerIDs:  input.ReviewerIDs,
		Labels:       input.Labels,
		TargetBranch: input.TargetBranch,
	}

	mr, err := client.UpdateMergeRequest(input.ProjectID, input.MergeRequestIID, opts)
	if err != nil {
		return nil, UpdateMergeRequestOutput{}, err
	}

	return nil, UpdateMergeRequestOutput{
		IID:    mr.IID,
		Title:  mr.Title,
		WebURL: mr.WebURL,
	}, nil
}

func mergeMergeRequestHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input MergeMergeRequestInput) (*mcp.CallToolResult, MergeMergeRequestOutput, error) {
	opts := &gitlab.MergeMergeRequestOptions{
		Squash:                   input.Squash,
		ShouldRemoveSourceBranch: input.ShouldRemoveSourceBranch,
	}

	mr, err := client.MergeMergeRequest(input.ProjectID, input.MergeRequestIID, opts)
	if err != nil {
		return nil, MergeMergeRequestOutput{}, err
	}

	return nil, MergeMergeRequestOutput{
		IID:    mr.IID,
		State:  mr.State,
		WebURL: mr.WebURL,
	}, nil
}

func getMergeRequestChangesHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input GetMergeRequestChangesInput) (*mcp.CallToolResult, GetMergeRequestChangesOutput, error) {
	diffs, err := client.GetMergeRequestChanges(input.ProjectID, input.MergeRequestIID)
	if err != nil {
		return nil, GetMergeRequestChangesOutput{}, err
	}

	changes := make([]ChangeInfo, len(diffs))
	for i, diff := range diffs {
		changes[i] = ChangeInfo{
			OldPath:     diff.OldPath,
			NewPath:     diff.NewPath,
			Diff:        diff.Diff,
			NewFile:     diff.NewFile,
			RenamedFile: diff.RenamedFile,
			DeletedFile: diff.DeletedFile,
		}
	}

	return nil, GetMergeRequestChangesOutput{Changes: changes}, nil
}
