package discussion

import (
	"context"

	"github.com/kqns91/gitlab-mcp/internal/gitlab"
	"github.com/kqns91/gitlab-mcp/internal/registry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AddCommentInput は add_merge_request_comment の入力パラメータ
type AddCommentInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
	Body            string `json:"body" jsonschema:"description:Comment body text"`
}

// AddCommentOutput は add_merge_request_comment の出力
type AddCommentOutput struct {
	ID         int64  `json:"id"`
	Body       string `json:"body"`
	AuthorName string `json:"author_name,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

// DiffPosition は差分内の位置を指定する
type DiffPosition struct {
	BaseSHA  string `json:"base_sha" jsonschema:"description:Base commit SHA"`
	StartSHA string `json:"start_sha" jsonschema:"description:Start commit SHA"`
	HeadSHA  string `json:"head_sha" jsonschema:"description:Head commit SHA"`
	OldPath  string `json:"old_path" jsonschema:"description:Old file path"`
	NewPath  string `json:"new_path" jsonschema:"description:New file path"`
	OldLine  *int   `json:"old_line,omitempty" jsonschema:"description:Line number in old file"`
	NewLine  *int   `json:"new_line,omitempty" jsonschema:"description:Line number in new file"`
}

// AddDiscussionInput は add_merge_request_discussion の入力パラメータ
type AddDiscussionInput struct {
	ProjectID       string        `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int           `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
	Body            string        `json:"body" jsonschema:"description:Discussion body text"`
	Position        *DiffPosition `json:"position,omitempty" jsonschema:"description:Position for line comment"`
}

// AddDiscussionOutput は add_merge_request_discussion の出力
type AddDiscussionOutput struct {
	ID string `json:"id"`
}

// ListDiscussionsInput は list_merge_request_discussions の入力パラメータ
type ListDiscussionsInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
	Page            int    `json:"page,omitempty" jsonschema:"description:Page number (default: 1)"`
	PerPage         int    `json:"per_page,omitempty" jsonschema:"description:Number of items per page (default: 100, max: 100)"`
}

// DiscussionNote はディスカッション内のノート情報
type DiscussionNote struct {
	ID         int64  `json:"id"`
	Body       string `json:"body"`
	AuthorName string `json:"author_name,omitempty"`
	Resolvable bool   `json:"resolvable"`
	Resolved   bool   `json:"resolved"`
}

// DiscussionSummary はディスカッションのサマリー情報
type DiscussionSummary struct {
	ID    string           `json:"id"`
	Notes []DiscussionNote `json:"notes"`
}

// ListDiscussionsOutput は list_merge_request_discussions の出力
type ListDiscussionsOutput struct {
	Discussions []DiscussionSummary `json:"discussions"`
}

// ResolveDiscussionInput は resolve_discussion の入力パラメータ
type ResolveDiscussionInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
	DiscussionID    string `json:"discussion_id" jsonschema:"description:Discussion ID"`
	Resolved        bool   `json:"resolved" jsonschema:"description:Set to true to resolve or false to unresolve"`
}

// ResolveDiscussionOutput は resolve_discussion の出力
type ResolveDiscussionOutput struct {
	ID       string `json:"id"`
	Resolved bool   `json:"resolved"`
}

// DeleteCommentInput は delete_merge_request_comment の入力パラメータ
type DeleteCommentInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
	NoteID          int    `json:"note_id" jsonschema:"description:Note ID to delete"`
}

// DeleteCommentOutput は delete_merge_request_comment の出力
type DeleteCommentOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// clientHolder holds the GitLab client for handlers
type clientHolder struct {
	client *gitlab.Client
}

var holder *clientHolder

// Register はディスカッション関連ツールを登録する
func Register(reg *registry.Registry, client *gitlab.Client) {
	holder = &clientHolder{client: client}

	registry.RegisterTool(reg, "add_merge_request_comment",
		"GitLab Merge Request に一般コメントを追加します",
		func(ctx context.Context, req *mcp.CallToolRequest, input AddCommentInput) (*mcp.CallToolResult, AddCommentOutput, error) {
			return addCommentHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "add_merge_request_discussion",
		"GitLab Merge Request に行コメント（ディスカッション）を作成します",
		func(ctx context.Context, req *mcp.CallToolRequest, input AddDiscussionInput) (*mcp.CallToolResult, AddDiscussionOutput, error) {
			return addDiscussionHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "list_merge_request_discussions",
		"GitLab Merge Request のディスカッション一覧を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ListDiscussionsInput) (*mcp.CallToolResult, ListDiscussionsOutput, error) {
			return listDiscussionsHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "resolve_discussion",
		"GitLab Merge Request のディスカッションを解決済み/未解決に設定します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ResolveDiscussionInput) (*mcp.CallToolResult, ResolveDiscussionOutput, error) {
			return resolveDiscussionHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "delete_merge_request_comment",
		"GitLab Merge Request のコメントを削除します",
		func(ctx context.Context, req *mcp.CallToolRequest, input DeleteCommentInput) (*mcp.CallToolResult, DeleteCommentOutput, error) {
			return deleteCommentHandler(holder.client, ctx, req, input)
		})
}

func addCommentHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input AddCommentInput) (*mcp.CallToolResult, AddCommentOutput, error) {
	note, err := client.AddMergeRequestComment(input.ProjectID, input.MergeRequestIID, input.Body)
	if err != nil {
		return nil, AddCommentOutput{}, err
	}

	authorName := note.Author.Username

	createdAt := ""
	if note.CreatedAt != nil {
		createdAt = note.CreatedAt.String()
	}

	return nil, AddCommentOutput{
		ID:         int64(note.ID),
		Body:       note.Body,
		AuthorName: authorName,
		CreatedAt:  createdAt,
	}, nil
}

func addDiscussionHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input AddDiscussionInput) (*mcp.CallToolResult, AddDiscussionOutput, error) {
	opts := &gitlab.CreateDiscussionOptions{
		Body: input.Body,
	}

	if input.Position != nil {
		opts.FilePath = input.Position.NewPath
		opts.OldLine = input.Position.OldLine
		opts.NewLine = input.Position.NewLine
		opts.BaseSHA = input.Position.BaseSHA
		opts.HeadSHA = input.Position.HeadSHA
		opts.StartSHA = input.Position.StartSHA
	}

	discussion, err := client.CreateMergeRequestDiscussion(input.ProjectID, input.MergeRequestIID, opts)
	if err != nil {
		return nil, AddDiscussionOutput{}, err
	}

	return nil, AddDiscussionOutput{
		ID: discussion.ID,
	}, nil
}

func listDiscussionsHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ListDiscussionsInput) (*mcp.CallToolResult, ListDiscussionsOutput, error) {
	var pagination *gitlab.PaginationOptions
	if input.Page > 0 || input.PerPage > 0 {
		pagination = &gitlab.PaginationOptions{
			Page:    input.Page,
			PerPage: input.PerPage,
		}
	}

	discussions, err := client.ListMergeRequestDiscussions(input.ProjectID, input.MergeRequestIID, pagination)
	if err != nil {
		return nil, ListDiscussionsOutput{}, err
	}

	summaries := make([]DiscussionSummary, len(discussions))
	for i, d := range discussions {
		notes := make([]DiscussionNote, len(d.Notes))
		for j, n := range d.Notes {
			authorName := n.Author.Username
			notes[j] = DiscussionNote{
				ID:         int64(n.ID),
				Body:       n.Body,
				AuthorName: authorName,
				Resolvable: n.Resolvable,
				Resolved:   n.Resolved,
			}
		}
		summaries[i] = DiscussionSummary{
			ID:    d.ID,
			Notes: notes,
		}
	}

	return nil, ListDiscussionsOutput{Discussions: summaries}, nil
}

func resolveDiscussionHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ResolveDiscussionInput) (*mcp.CallToolResult, ResolveDiscussionOutput, error) {
	discussion, err := client.ResolveMergeRequestDiscussion(input.ProjectID, input.MergeRequestIID, input.DiscussionID, input.Resolved)
	if err != nil {
		return nil, ResolveDiscussionOutput{}, err
	}

	// Determine resolved state from first note
	resolved := false
	if len(discussion.Notes) > 0 {
		resolved = discussion.Notes[0].Resolved
	}

	return nil, ResolveDiscussionOutput{
		ID:       discussion.ID,
		Resolved: resolved,
	}, nil
}

func deleteCommentHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input DeleteCommentInput) (*mcp.CallToolResult, DeleteCommentOutput, error) {
	err := client.DeleteMergeRequestNote(input.ProjectID, input.MergeRequestIID, input.NoteID)
	if err != nil {
		return nil, DeleteCommentOutput{}, err
	}

	return nil, DeleteCommentOutput{
		Success: true,
		Message: "Comment deleted successfully",
	}, nil
}
