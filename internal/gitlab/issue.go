package gitlab

import (
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

// ListProjectIssuesOptions はIssue一覧取得のオプション
type ListProjectIssuesOptions struct {
	State      *string
	Labels     []string
	AssigneeID *int
	AuthorID   *int
	Search     *string
	Page       int
	PerPage    int
}

// ListProjectIssues はプロジェクトのIssue一覧を取得する
func (c *Client) ListProjectIssues(projectID string, opts *ListProjectIssuesOptions) ([]*gogitlab.Issue, error) {
	page, perPage := 1, 100
	if opts != nil {
		if opts.Page > 0 {
			page = opts.Page
		}
		if opts.PerPage > 0 {
			perPage = opts.PerPage
		}
	}

	listOpts := &gogitlab.ListProjectIssuesOptions{
		ListOptions: gogitlab.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}

	if opts != nil {
		if opts.State != nil {
			listOpts.State = opts.State
		}
		if len(opts.Labels) > 0 {
			labels := gogitlab.LabelOptions(opts.Labels)
			listOpts.Labels = &labels
		}
		if opts.AssigneeID != nil {
			listOpts.AssigneeID = gogitlab.AssigneeID(*opts.AssigneeID)
		}
		if opts.AuthorID != nil {
			authorID := int64(*opts.AuthorID)
			listOpts.AuthorID = &authorID
		}
		if opts.Search != nil {
			listOpts.Search = opts.Search
		}
	}

	issues, resp, err := c.client.Issues.ListProjectIssues(projectID, listOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return issues, nil
}

// GetIssue はIssueの詳細を取得する
func (c *Client) GetIssue(projectID string, issueIID int) (*gogitlab.Issue, error) {
	issue, resp, err := c.client.Issues.GetIssue(projectID, int64(issueIID))
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return issue, nil
}

// CreateIssueOptions はIssue作成のオプション
type CreateIssueOptions struct {
	Title       string
	Description *string
	Labels      []string
	AssigneeIDs []int
	MilestoneID *int
}

// CreateIssue は新しいIssueを作成する
func (c *Client) CreateIssue(projectID string, opts *CreateIssueOptions) (*gogitlab.Issue, error) {
	createOpts := &gogitlab.CreateIssueOptions{
		Title: &opts.Title,
	}

	if opts.Description != nil {
		createOpts.Description = opts.Description
	}

	if len(opts.AssigneeIDs) > 0 {
		assigneeIDs := make([]int64, len(opts.AssigneeIDs))
		for i, id := range opts.AssigneeIDs {
			assigneeIDs[i] = int64(id)
		}
		createOpts.AssigneeIDs = &assigneeIDs
	}

	if len(opts.Labels) > 0 {
		labels := gogitlab.LabelOptions(opts.Labels)
		createOpts.Labels = &labels
	}

	if opts.MilestoneID != nil {
		milestoneID := int64(*opts.MilestoneID)
		createOpts.MilestoneID = &milestoneID
	}

	issue, resp, err := c.client.Issues.CreateIssue(projectID, createOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return issue, nil
}

// UpdateIssueOptions はIssue更新のオプション
type UpdateIssueOptions struct {
	Title       *string
	Description *string
	StateEvent  *string
	Labels      []string
	AssigneeIDs []int
	MilestoneID *int
}

// UpdateIssue は既存のIssueを更新する
func (c *Client) UpdateIssue(projectID string, issueIID int, opts *UpdateIssueOptions) (*gogitlab.Issue, error) {
	updateOpts := &gogitlab.UpdateIssueOptions{}

	if opts.Title != nil {
		updateOpts.Title = opts.Title
	}
	if opts.Description != nil {
		updateOpts.Description = opts.Description
	}
	if opts.StateEvent != nil {
		updateOpts.StateEvent = opts.StateEvent
	}

	if len(opts.AssigneeIDs) > 0 {
		assigneeIDs := make([]int64, len(opts.AssigneeIDs))
		for i, id := range opts.AssigneeIDs {
			assigneeIDs[i] = int64(id)
		}
		updateOpts.AssigneeIDs = &assigneeIDs
	}

	if len(opts.Labels) > 0 {
		labels := gogitlab.LabelOptions(opts.Labels)
		updateOpts.Labels = &labels
	}

	if opts.MilestoneID != nil {
		milestoneID := int64(*opts.MilestoneID)
		updateOpts.MilestoneID = &milestoneID
	}

	issue, resp, err := c.client.Issues.UpdateIssue(projectID, int64(issueIID), updateOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return issue, nil
}

// DeleteIssue はIssueを削除する
func (c *Client) DeleteIssue(projectID string, issueIID int) error {
	resp, err := c.client.Issues.DeleteIssue(projectID, int64(issueIID))
	if err != nil {
		return FromGitLabResponse(err, resp)
	}
	return nil
}

// ListIssueNotes はIssueのコメント一覧を取得する
func (c *Client) ListIssueNotes(projectID string, issueIID int, pagination *PaginationOptions) ([]*gogitlab.Note, error) {
	page, perPage := 1, 100
	if pagination != nil {
		if pagination.Page > 0 {
			page = pagination.Page
		}
		if pagination.PerPage > 0 {
			perPage = pagination.PerPage
		}
	}

	opts := &gogitlab.ListIssueNotesOptions{
		ListOptions: gogitlab.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}
	notes, resp, err := c.client.Notes.ListIssueNotes(projectID, int64(issueIID), opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return notes, nil
}

// CreateIssueNote はIssueにコメントを追加する
func (c *Client) CreateIssueNote(projectID string, issueIID int, body string) (*gogitlab.Note, error) {
	opts := &gogitlab.CreateIssueNoteOptions{
		Body: &body,
	}

	note, resp, err := c.client.Notes.CreateIssueNote(projectID, int64(issueIID), opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return note, nil
}

// DeleteIssueNote はIssueのコメントを削除する
func (c *Client) DeleteIssueNote(projectID string, issueIID int, noteID int) error {
	resp, err := c.client.Notes.DeleteIssueNote(projectID, int64(issueIID), int64(noteID))
	if err != nil {
		return FromGitLabResponse(err, resp)
	}
	return nil
}

// ListIssueDiscussions はIssueのディスカッション一覧を取得する
func (c *Client) ListIssueDiscussions(projectID string, issueIID int, pagination *PaginationOptions) ([]*gogitlab.Discussion, error) {
	page, perPage := 1, 100
	if pagination != nil {
		if pagination.Page > 0 {
			page = pagination.Page
		}
		if pagination.PerPage > 0 {
			perPage = pagination.PerPage
		}
	}

	opts := &gogitlab.ListIssueDiscussionsOptions{
		ListOptions: gogitlab.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}
	discussions, resp, err := c.client.Discussions.ListIssueDiscussions(projectID, int64(issueIID), opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return discussions, nil
}

// CreateIssueDiscussion はIssueにディスカッションを作成する
func (c *Client) CreateIssueDiscussion(projectID string, issueIID int, body string) (*gogitlab.Discussion, error) {
	opts := &gogitlab.CreateIssueDiscussionOptions{
		Body: &body,
	}

	discussion, resp, err := c.client.Discussions.CreateIssueDiscussion(projectID, int64(issueIID), opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return discussion, nil
}

// AddIssueDiscussionNote はIssueのディスカッションに返信を追加する
func (c *Client) AddIssueDiscussionNote(projectID string, issueIID int, discussionID string, body string) (*gogitlab.Note, error) {
	opts := &gogitlab.AddIssueDiscussionNoteOptions{
		Body: &body,
	}

	note, resp, err := c.client.Discussions.AddIssueDiscussionNote(projectID, int64(issueIID), discussionID, opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return note, nil
}

