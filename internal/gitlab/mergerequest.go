package gitlab

import (
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

// ListMergeRequestsOptions はMR一覧取得のオプション
type ListMergeRequestsOptions struct {
	State      *string
	AuthorID   *int
	AssigneeID *int
	Page       int
	PerPage    int
}

// ListMergeRequests はプロジェクトのMR一覧を取得する
func (c *Client) ListMergeRequests(projectID string, opts *ListMergeRequestsOptions) ([]*gogitlab.BasicMergeRequest, error) {
	page, perPage := 1, 100
	if opts != nil {
		if opts.Page > 0 {
			page = opts.Page
		}
		if opts.PerPage > 0 {
			perPage = opts.PerPage
		}
	}

	listOpts := &gogitlab.ListProjectMergeRequestsOptions{
		ListOptions: gogitlab.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}

	if opts != nil {
		if opts.State != nil {
			listOpts.State = opts.State
		}
		if opts.AuthorID != nil {
			authorID := int64(*opts.AuthorID)
			listOpts.AuthorID = &authorID
		}
		if opts.AssigneeID != nil {
			listOpts.AssigneeID = gogitlab.AssigneeID(*opts.AssigneeID)
		}
	}

	mrs, resp, err := c.client.MergeRequests.ListProjectMergeRequests(projectID, listOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return mrs, nil
}

// GetMergeRequest はMRの詳細を取得する
func (c *Client) GetMergeRequest(projectID string, mrIID int) (*gogitlab.MergeRequest, error) {
	mr, resp, err := c.client.MergeRequests.GetMergeRequest(projectID, int64(mrIID), nil)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return mr, nil
}

// CreateMergeRequestOptions はMR作成のオプション
type CreateMergeRequestOptions struct {
	SourceBranch string
	TargetBranch string
	Title        string
	Description  *string
	AssigneeIDs  []int
	ReviewerIDs  []int
	Labels       []string
}

// CreateMergeRequest は新しいMRを作成する
func (c *Client) CreateMergeRequest(projectID string, opts *CreateMergeRequestOptions) (*gogitlab.MergeRequest, error) {
	createOpts := &gogitlab.CreateMergeRequestOptions{
		SourceBranch: &opts.SourceBranch,
		TargetBranch: &opts.TargetBranch,
		Title:        &opts.Title,
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

	if len(opts.ReviewerIDs) > 0 {
		reviewerIDs := make([]int64, len(opts.ReviewerIDs))
		for i, id := range opts.ReviewerIDs {
			reviewerIDs[i] = int64(id)
		}
		createOpts.ReviewerIDs = &reviewerIDs
	}

	if len(opts.Labels) > 0 {
		labels := gogitlab.LabelOptions(opts.Labels)
		createOpts.Labels = &labels
	}

	mr, resp, err := c.client.MergeRequests.CreateMergeRequest(projectID, createOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return mr, nil
}

// UpdateMergeRequestOptions はMR更新のオプション
type UpdateMergeRequestOptions struct {
	Title        *string
	Description  *string
	AssigneeIDs  []int
	ReviewerIDs  []int
	Labels       []string
	TargetBranch *string
}

// UpdateMergeRequest は既存のMRを更新する
func (c *Client) UpdateMergeRequest(projectID string, mrIID int, opts *UpdateMergeRequestOptions) (*gogitlab.MergeRequest, error) {
	updateOpts := &gogitlab.UpdateMergeRequestOptions{}

	if opts.Title != nil {
		updateOpts.Title = opts.Title
	}
	if opts.Description != nil {
		updateOpts.Description = opts.Description
	}
	if opts.TargetBranch != nil {
		updateOpts.TargetBranch = opts.TargetBranch
	}

	if len(opts.AssigneeIDs) > 0 {
		assigneeIDs := make([]int64, len(opts.AssigneeIDs))
		for i, id := range opts.AssigneeIDs {
			assigneeIDs[i] = int64(id)
		}
		updateOpts.AssigneeIDs = &assigneeIDs
	}

	if len(opts.ReviewerIDs) > 0 {
		reviewerIDs := make([]int64, len(opts.ReviewerIDs))
		for i, id := range opts.ReviewerIDs {
			reviewerIDs[i] = int64(id)
		}
		updateOpts.ReviewerIDs = &reviewerIDs
	}

	if len(opts.Labels) > 0 {
		labels := gogitlab.LabelOptions(opts.Labels)
		updateOpts.Labels = &labels
	}

	mr, resp, err := c.client.MergeRequests.UpdateMergeRequest(projectID, int64(mrIID), updateOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return mr, nil
}

// MergeMergeRequestOptions はMRマージのオプション
type MergeMergeRequestOptions struct {
	Squash                   *bool
	ShouldRemoveSourceBranch *bool
	MergeCommitMessage       *string
	SquashCommitMessage      *string
}

// MergeMergeRequest はMRをマージする
func (c *Client) MergeMergeRequest(projectID string, mrIID int, opts *MergeMergeRequestOptions) (*gogitlab.MergeRequest, error) {
	mergeOpts := &gogitlab.AcceptMergeRequestOptions{}

	if opts != nil {
		if opts.Squash != nil {
			mergeOpts.Squash = opts.Squash
		}
		if opts.ShouldRemoveSourceBranch != nil {
			mergeOpts.ShouldRemoveSourceBranch = opts.ShouldRemoveSourceBranch
		}
		if opts.MergeCommitMessage != nil {
			mergeOpts.MergeCommitMessage = opts.MergeCommitMessage
		}
		if opts.SquashCommitMessage != nil {
			mergeOpts.SquashCommitMessage = opts.SquashCommitMessage
		}
	}

	mr, resp, err := c.client.MergeRequests.AcceptMergeRequest(projectID, int64(mrIID), mergeOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return mr, nil
}

// GetMergeRequestChanges はMRの変更差分を取得する
func (c *Client) GetMergeRequestChanges(projectID string, mrIID int, pagination *PaginationOptions) ([]*gogitlab.MergeRequestDiff, error) {
	page, perPage := 1, 100
	if pagination != nil {
		if pagination.Page > 0 {
			page = pagination.Page
		}
		if pagination.PerPage > 0 {
			perPage = pagination.PerPage
		}
	}

	opts := &gogitlab.ListMergeRequestDiffsOptions{
		ListOptions: gogitlab.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}
	diffs, resp, err := c.client.MergeRequests.ListMergeRequestDiffs(projectID, int64(mrIID), opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return diffs, nil
}
