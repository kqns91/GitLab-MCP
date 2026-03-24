package issue

import (
	"context"

	"github.com/kqns91/gitlab-mcp/internal/gitlab"
	"github.com/kqns91/gitlab-mcp/internal/registry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListIssuesInput は list_issues の入力パラメータ
type ListIssuesInput struct {
	ProjectID  string   `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	State      *string  `json:"state,omitempty" jsonschema:"enum:opened,enum:closed,enum:all,description:Issue state filter"`
	Labels     []string `json:"labels,omitempty" jsonschema:"description:Label names filter"`
	AssigneeID *int     `json:"assignee_id,omitempty" jsonschema:"description:Assignee user ID filter"`
	AuthorID   *int     `json:"author_id,omitempty" jsonschema:"description:Author user ID filter"`
	Search     *string  `json:"search,omitempty" jsonschema:"description:Search query"`
	Page       int      `json:"page,omitempty" jsonschema:"description:Page number (default: 1)"`
	PerPage    int      `json:"per_page,omitempty" jsonschema:"description:Number of items per page (default: 100, max: 100)"`
}

// IssueSummary はIssue一覧の各項目
type IssueSummary struct {
	IID        int64    `json:"iid"`
	Title      string   `json:"title"`
	State      string   `json:"state"`
	WebURL     string   `json:"web_url"`
	AuthorName string   `json:"author_name,omitempty"`
	Labels     []string `json:"labels,omitempty"`
}

// ListIssuesOutput は list_issues の出力
type ListIssuesOutput struct {
	Issues []IssueSummary `json:"issues"`
}

// GetIssueInput は get_issue の入力パラメータ
type GetIssueInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	IssueIID  int    `json:"issue_iid" jsonschema:"description:Issue IID"`
}

// IssueDetail はIssue詳細情報
type IssueDetail struct {
	IID         int64    `json:"iid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	State       string   `json:"state"`
	WebURL      string   `json:"web_url"`
	AuthorName  string   `json:"author_name,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

// GetIssueOutput は get_issue の出力
type GetIssueOutput = IssueDetail

// CreateIssueInput は create_issue の入力パラメータ
type CreateIssueInput struct {
	ProjectID   string   `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	Title       string   `json:"title" jsonschema:"description:Issue title"`
	Description *string  `json:"description,omitempty" jsonschema:"description:Issue description"`
	Labels      []string `json:"labels,omitempty" jsonschema:"description:Labels to add"`
	AssigneeIDs []int    `json:"assignee_ids,omitempty" jsonschema:"description:Assignee user IDs"`
	MilestoneID *int     `json:"milestone_id,omitempty" jsonschema:"description:Milestone ID"`
}

// CreateIssueOutput は create_issue の出力
type CreateIssueOutput struct {
	IID    int64  `json:"iid"`
	Title  string `json:"title"`
	WebURL string `json:"web_url"`
}

// UpdateIssueInput は update_issue の入力パラメータ
type UpdateIssueInput struct {
	ProjectID   string   `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	IssueIID    int      `json:"issue_iid" jsonschema:"description:Issue IID"`
	Title       *string  `json:"title,omitempty" jsonschema:"description:New title"`
	Description *string  `json:"description,omitempty" jsonschema:"description:New description"`
	StateEvent  *string  `json:"state_event,omitempty" jsonschema:"enum:close,enum:reopen,description:State event"`
	Labels      []string `json:"labels,omitempty" jsonschema:"description:New labels"`
	AssigneeIDs []int    `json:"assignee_ids,omitempty" jsonschema:"description:New assignee user IDs"`
	MilestoneID *int     `json:"milestone_id,omitempty" jsonschema:"description:New milestone ID"`
}

// UpdateIssueOutput は update_issue の出力
type UpdateIssueOutput struct {
	IID    int64  `json:"iid"`
	Title  string `json:"title"`
	State  string `json:"state"`
	WebURL string `json:"web_url"`
}

// DeleteIssueInput は delete_issue の入力パラメータ
type DeleteIssueInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	IssueIID  int    `json:"issue_iid" jsonschema:"description:Issue IID"`
}

// DeleteIssueOutput は delete_issue の出力
type DeleteIssueOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ListIssueNotesInput は list_issue_notes の入力パラメータ
type ListIssueNotesInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	IssueIID  int    `json:"issue_iid" jsonschema:"description:Issue IID"`
	Page      int    `json:"page,omitempty" jsonschema:"description:Page number (default: 1)"`
	PerPage   int    `json:"per_page,omitempty" jsonschema:"description:Number of items per page (default: 100, max: 100)"`
}

// NoteInfo はノート情報
type NoteInfo struct {
	ID         int64  `json:"id"`
	Body       string `json:"body"`
	AuthorName string `json:"author_name,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	System     bool   `json:"system"`
}

// ListIssueNotesOutput は list_issue_notes の出力
type ListIssueNotesOutput struct {
	Notes []NoteInfo `json:"notes"`
}

// CreateIssueNoteInput は create_issue_note の入力パラメータ
type CreateIssueNoteInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	IssueIID  int    `json:"issue_iid" jsonschema:"description:Issue IID"`
	Body      string `json:"body" jsonschema:"description:Comment body text"`
}

// CreateIssueNoteOutput は create_issue_note の出力
type CreateIssueNoteOutput struct {
	ID         int64  `json:"id"`
	Body       string `json:"body"`
	AuthorName string `json:"author_name,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

// DeleteIssueNoteInput は delete_issue_note の入力パラメータ
type DeleteIssueNoteInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	IssueIID  int    `json:"issue_iid" jsonschema:"description:Issue IID"`
	NoteID    int    `json:"note_id" jsonschema:"description:Note ID to delete"`
}

// DeleteIssueNoteOutput は delete_issue_note の出力
type DeleteIssueNoteOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ListIssueDiscussionsInput は list_issue_discussions の入力パラメータ
type ListIssueDiscussionsInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	IssueIID  int    `json:"issue_iid" jsonschema:"description:Issue IID"`
	Page      int    `json:"page,omitempty" jsonschema:"description:Page number (default: 1)"`
	PerPage   int    `json:"per_page,omitempty" jsonschema:"description:Number of items per page (default: 100, max: 100)"`
}

// DiscussionNote はディスカッション内のノート情報
type DiscussionNote struct {
	ID         int64  `json:"id"`
	Body       string `json:"body"`
	AuthorName string `json:"author_name,omitempty"`
}

// DiscussionSummary はディスカッションのサマリー情報
type DiscussionSummary struct {
	ID    string           `json:"id"`
	Notes []DiscussionNote `json:"notes"`
}

// ListIssueDiscussionsOutput は list_issue_discussions の出力
type ListIssueDiscussionsOutput struct {
	Discussions []DiscussionSummary `json:"discussions"`
}

// CreateIssueDiscussionInput は create_issue_discussion の入力パラメータ
type CreateIssueDiscussionInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	IssueIID  int    `json:"issue_iid" jsonschema:"description:Issue IID"`
	Body      string `json:"body" jsonschema:"description:Discussion body text"`
}

// CreateIssueDiscussionOutput は create_issue_discussion の出力
type CreateIssueDiscussionOutput struct {
	ID string `json:"id"`
}

// ReplyToIssueDiscussionInput は reply_to_issue_discussion の入力パラメータ
type ReplyToIssueDiscussionInput struct {
	ProjectID    string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	IssueIID     int    `json:"issue_iid" jsonschema:"description:Issue IID"`
	DiscussionID string `json:"discussion_id" jsonschema:"description:Discussion ID to reply to"`
	Body         string `json:"body" jsonschema:"description:Reply body text"`
}

// ReplyToIssueDiscussionOutput は reply_to_issue_discussion の出力
type ReplyToIssueDiscussionOutput struct {
	ID         int64  `json:"id"`
	Body       string `json:"body"`
	AuthorName string `json:"author_name,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

// clientHolder holds the GitLab client for handlers
type clientHolder struct {
	client *gitlab.Client
}

var holder *clientHolder

// Register はIssue関連ツールを登録する
func Register(reg *registry.Registry, client *gitlab.Client) {
	holder = &clientHolder{client: client}

	registry.RegisterTool(reg, "list_issues",
		"GitLab プロジェクトの Issue 一覧を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ListIssuesInput) (*mcp.CallToolResult, ListIssuesOutput, error) {
			return listIssuesHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "get_issue",
		"GitLab Issue の詳細情報を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input GetIssueInput) (*mcp.CallToolResult, GetIssueOutput, error) {
			return getIssueHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "create_issue",
		"GitLab に新しい Issue を作成します",
		func(ctx context.Context, req *mcp.CallToolRequest, input CreateIssueInput) (*mcp.CallToolResult, CreateIssueOutput, error) {
			return createIssueHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "update_issue",
		"GitLab Issue を更新します",
		func(ctx context.Context, req *mcp.CallToolRequest, input UpdateIssueInput) (*mcp.CallToolResult, UpdateIssueOutput, error) {
			return updateIssueHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "delete_issue",
		"GitLab Issue を削除します",
		func(ctx context.Context, req *mcp.CallToolRequest, input DeleteIssueInput) (*mcp.CallToolResult, DeleteIssueOutput, error) {
			return deleteIssueHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "list_issue_notes",
		"GitLab Issue のコメント一覧を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ListIssueNotesInput) (*mcp.CallToolResult, ListIssueNotesOutput, error) {
			return listIssueNotesHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "create_issue_note",
		"GitLab Issue にコメントを追加します",
		func(ctx context.Context, req *mcp.CallToolRequest, input CreateIssueNoteInput) (*mcp.CallToolResult, CreateIssueNoteOutput, error) {
			return createIssueNoteHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "delete_issue_note",
		"GitLab Issue のコメントを削除します",
		func(ctx context.Context, req *mcp.CallToolRequest, input DeleteIssueNoteInput) (*mcp.CallToolResult, DeleteIssueNoteOutput, error) {
			return deleteIssueNoteHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "list_issue_discussions",
		"GitLab Issue のディスカッション一覧を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ListIssueDiscussionsInput) (*mcp.CallToolResult, ListIssueDiscussionsOutput, error) {
			return listIssueDiscussionsHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "create_issue_discussion",
		"GitLab Issue にディスカッションを作成します",
		func(ctx context.Context, req *mcp.CallToolRequest, input CreateIssueDiscussionInput) (*mcp.CallToolResult, CreateIssueDiscussionOutput, error) {
			return createIssueDiscussionHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "reply_to_issue_discussion",
		"GitLab Issue のディスカッションに返信を追加します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ReplyToIssueDiscussionInput) (*mcp.CallToolResult, ReplyToIssueDiscussionOutput, error) {
			return replyToIssueDiscussionHandler(holder.client, ctx, req, input)
		})
}

func listIssuesHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ListIssuesInput) (*mcp.CallToolResult, ListIssuesOutput, error) {
	opts := &gitlab.ListProjectIssuesOptions{
		State:      input.State,
		Labels:     input.Labels,
		AssigneeID: input.AssigneeID,
		AuthorID:   input.AuthorID,
		Search:     input.Search,
		Page:       input.Page,
		PerPage:    input.PerPage,
	}

	issues, err := client.ListProjectIssues(input.ProjectID, opts)
	if err != nil {
		return nil, ListIssuesOutput{}, err
	}

	summaries := make([]IssueSummary, len(issues))
	for i, issue := range issues {
		authorName := ""
		if issue.Author != nil {
			authorName = issue.Author.Username
		}
		summaries[i] = IssueSummary{
			IID:        issue.IID,
			Title:      issue.Title,
			State:      issue.State,
			WebURL:     issue.WebURL,
			AuthorName: authorName,
			Labels:     issue.Labels,
		}
	}

	return nil, ListIssuesOutput{Issues: summaries}, nil
}

func getIssueHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input GetIssueInput) (*mcp.CallToolResult, GetIssueOutput, error) {
	issue, err := client.GetIssue(input.ProjectID, input.IssueIID)
	if err != nil {
		return nil, GetIssueOutput{}, err
	}

	authorName := ""
	if issue.Author != nil {
		authorName = issue.Author.Username
	}

	createdAt := ""
	if issue.CreatedAt != nil {
		createdAt = issue.CreatedAt.String()
	}
	updatedAt := ""
	if issue.UpdatedAt != nil {
		updatedAt = issue.UpdatedAt.String()
	}

	return nil, GetIssueOutput{
		IID:         issue.IID,
		Title:       issue.Title,
		Description: issue.Description,
		State:       issue.State,
		WebURL:      issue.WebURL,
		AuthorName:  authorName,
		Labels:      issue.Labels,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func createIssueHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input CreateIssueInput) (*mcp.CallToolResult, CreateIssueOutput, error) {
	opts := &gitlab.CreateIssueOptions{
		Title:       input.Title,
		Description: input.Description,
		Labels:      input.Labels,
		AssigneeIDs: input.AssigneeIDs,
		MilestoneID: input.MilestoneID,
	}

	issue, err := client.CreateIssue(input.ProjectID, opts)
	if err != nil {
		return nil, CreateIssueOutput{}, err
	}

	return nil, CreateIssueOutput{
		IID:    issue.IID,
		Title:  issue.Title,
		WebURL: issue.WebURL,
	}, nil
}

func updateIssueHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input UpdateIssueInput) (*mcp.CallToolResult, UpdateIssueOutput, error) {
	opts := &gitlab.UpdateIssueOptions{
		Title:       input.Title,
		Description: input.Description,
		StateEvent:  input.StateEvent,
		Labels:      input.Labels,
		AssigneeIDs: input.AssigneeIDs,
		MilestoneID: input.MilestoneID,
	}

	issue, err := client.UpdateIssue(input.ProjectID, input.IssueIID, opts)
	if err != nil {
		return nil, UpdateIssueOutput{}, err
	}

	return nil, UpdateIssueOutput{
		IID:    issue.IID,
		Title:  issue.Title,
		State:  issue.State,
		WebURL: issue.WebURL,
	}, nil
}

func deleteIssueHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input DeleteIssueInput) (*mcp.CallToolResult, DeleteIssueOutput, error) {
	err := client.DeleteIssue(input.ProjectID, input.IssueIID)
	if err != nil {
		return nil, DeleteIssueOutput{}, err
	}

	return nil, DeleteIssueOutput{
		Success: true,
		Message: "Issue deleted successfully",
	}, nil
}

func listIssueNotesHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ListIssueNotesInput) (*mcp.CallToolResult, ListIssueNotesOutput, error) {
	var pagination *gitlab.PaginationOptions
	if input.Page > 0 || input.PerPage > 0 {
		pagination = &gitlab.PaginationOptions{
			Page:    input.Page,
			PerPage: input.PerPage,
		}
	}

	notes, err := client.ListIssueNotes(input.ProjectID, input.IssueIID, pagination)
	if err != nil {
		return nil, ListIssueNotesOutput{}, err
	}

	infos := make([]NoteInfo, len(notes))
	for i, n := range notes {
		authorName := ""
		if n.Author.Username != "" {
			authorName = n.Author.Username
		}
		createdAt := ""
		if n.CreatedAt != nil {
			createdAt = n.CreatedAt.String()
		}
		infos[i] = NoteInfo{
			ID:         int64(n.ID),
			Body:       n.Body,
			AuthorName: authorName,
			CreatedAt:  createdAt,
			System:     n.System,
		}
	}

	return nil, ListIssueNotesOutput{Notes: infos}, nil
}

func createIssueNoteHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input CreateIssueNoteInput) (*mcp.CallToolResult, CreateIssueNoteOutput, error) {
	note, err := client.CreateIssueNote(input.ProjectID, input.IssueIID, input.Body)
	if err != nil {
		return nil, CreateIssueNoteOutput{}, err
	}

	authorName := note.Author.Username
	createdAt := ""
	if note.CreatedAt != nil {
		createdAt = note.CreatedAt.String()
	}

	return nil, CreateIssueNoteOutput{
		ID:         int64(note.ID),
		Body:       note.Body,
		AuthorName: authorName,
		CreatedAt:  createdAt,
	}, nil
}

func deleteIssueNoteHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input DeleteIssueNoteInput) (*mcp.CallToolResult, DeleteIssueNoteOutput, error) {
	err := client.DeleteIssueNote(input.ProjectID, input.IssueIID, input.NoteID)
	if err != nil {
		return nil, DeleteIssueNoteOutput{}, err
	}

	return nil, DeleteIssueNoteOutput{
		Success: true,
		Message: "Note deleted successfully",
	}, nil
}

func listIssueDiscussionsHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ListIssueDiscussionsInput) (*mcp.CallToolResult, ListIssueDiscussionsOutput, error) {
	var pagination *gitlab.PaginationOptions
	if input.Page > 0 || input.PerPage > 0 {
		pagination = &gitlab.PaginationOptions{
			Page:    input.Page,
			PerPage: input.PerPage,
		}
	}

	discussions, err := client.ListIssueDiscussions(input.ProjectID, input.IssueIID, pagination)
	if err != nil {
		return nil, ListIssueDiscussionsOutput{}, err
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
			}
		}
		summaries[i] = DiscussionSummary{
			ID:    d.ID,
			Notes: notes,
		}
	}

	return nil, ListIssueDiscussionsOutput{Discussions: summaries}, nil
}

func createIssueDiscussionHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input CreateIssueDiscussionInput) (*mcp.CallToolResult, CreateIssueDiscussionOutput, error) {
	discussion, err := client.CreateIssueDiscussion(input.ProjectID, input.IssueIID, input.Body)
	if err != nil {
		return nil, CreateIssueDiscussionOutput{}, err
	}

	return nil, CreateIssueDiscussionOutput{
		ID: discussion.ID,
	}, nil
}

func replyToIssueDiscussionHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ReplyToIssueDiscussionInput) (*mcp.CallToolResult, ReplyToIssueDiscussionOutput, error) {
	note, err := client.AddIssueDiscussionNote(input.ProjectID, input.IssueIID, input.DiscussionID, input.Body)
	if err != nil {
		return nil, ReplyToIssueDiscussionOutput{}, err
	}

	authorName := note.Author.Username
	createdAt := ""
	if note.CreatedAt != nil {
		createdAt = note.CreatedAt.String()
	}

	return nil, ReplyToIssueDiscussionOutput{
		ID:         int64(note.ID),
		Body:       note.Body,
		AuthorName: authorName,
		CreatedAt:  createdAt,
	}, nil
}
