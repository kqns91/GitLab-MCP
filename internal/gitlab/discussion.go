package gitlab

import (
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

// AddMergeRequestComment はMRに一般コメントを追加する
func (c *Client) AddMergeRequestComment(projectID string, mrIID int, body string) (*gogitlab.Note, error) {
	opts := &gogitlab.CreateMergeRequestNoteOptions{
		Body: &body,
	}

	note, resp, err := c.client.Notes.CreateMergeRequestNote(projectID, int64(mrIID), opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return note, nil
}

// CreateDiscussionOptions はディスカッション作成のオプション
type CreateDiscussionOptions struct {
	Body     string
	FilePath string
	OldLine  *int
	NewLine  *int
	BaseSHA  string
	HeadSHA  string
	StartSHA string
}

// CreateMergeRequestDiscussion は行コメント（ディスカッション）を作成する
func (c *Client) CreateMergeRequestDiscussion(projectID string, mrIID int, opts *CreateDiscussionOptions) (*gogitlab.Discussion, error) {
	createOpts := &gogitlab.CreateMergeRequestDiscussionOptions{
		Body: &opts.Body,
	}

	// 位置情報が指定されている場合
	if opts.FilePath != "" {
		positionType := "text"
		position := &gogitlab.PositionOptions{
			PositionType: &positionType,
			NewPath:      &opts.FilePath,
			OldPath:      &opts.FilePath,
		}

		if opts.NewLine != nil {
			newLine := int64(*opts.NewLine)
			position.NewLine = &newLine
		}
		if opts.OldLine != nil {
			oldLine := int64(*opts.OldLine)
			position.OldLine = &oldLine
		}
		if opts.BaseSHA != "" {
			position.BaseSHA = &opts.BaseSHA
		}
		if opts.HeadSHA != "" {
			position.HeadSHA = &opts.HeadSHA
		}
		if opts.StartSHA != "" {
			position.StartSHA = &opts.StartSHA
		}

		createOpts.Position = position
	}

	discussion, resp, err := c.client.Discussions.CreateMergeRequestDiscussion(projectID, int64(mrIID), createOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return discussion, nil
}

// PaginationOptions はページネーションのオプション
type PaginationOptions struct {
	Page    int
	PerPage int
}

// ListMergeRequestDiscussions はMRのディスカッション一覧を取得する
func (c *Client) ListMergeRequestDiscussions(projectID string, mrIID int, pagination *PaginationOptions) ([]*gogitlab.Discussion, error) {
	page, perPage := 1, 100
	if pagination != nil {
		if pagination.Page > 0 {
			page = pagination.Page
		}
		if pagination.PerPage > 0 {
			perPage = pagination.PerPage
		}
	}

	opts := &gogitlab.ListMergeRequestDiscussionsOptions{
		ListOptions: gogitlab.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}
	discussions, resp, err := c.client.Discussions.ListMergeRequestDiscussions(projectID, int64(mrIID), opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return discussions, nil
}

// ResolveMergeRequestDiscussion はディスカッションの解決状態を変更する
func (c *Client) ResolveMergeRequestDiscussion(projectID string, mrIID int, discussionID string, resolved bool) (*gogitlab.Discussion, error) {
	opts := &gogitlab.ResolveMergeRequestDiscussionOptions{
		Resolved: &resolved,
	}

	discussion, resp, err := c.client.Discussions.ResolveMergeRequestDiscussion(projectID, int64(mrIID), discussionID, opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return discussion, nil
}

// DeleteMergeRequestNote はMRのコメント（ノート）を削除する
func (c *Client) DeleteMergeRequestNote(projectID string, mrIID int, noteID int) error {
	resp, err := c.client.Notes.DeleteMergeRequestNote(projectID, int64(mrIID), int64(noteID))
	if err != nil {
		return FromGitLabResponse(err, resp)
	}
	return nil
}
