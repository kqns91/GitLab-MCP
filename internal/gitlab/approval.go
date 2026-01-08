package gitlab

import (
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

// ApproveMergeRequest はMRを承認する
func (c *Client) ApproveMergeRequest(projectID string, mrIID int) (*gogitlab.MergeRequestApprovals, error) {
	approvals, resp, err := c.client.MergeRequestApprovals.ApproveMergeRequest(projectID, int64(mrIID), nil)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return approvals, nil
}

// UnapproveMergeRequest はMRの承認を取り消す
func (c *Client) UnapproveMergeRequest(projectID string, mrIID int) error {
	resp, err := c.client.MergeRequestApprovals.UnapproveMergeRequest(projectID, int64(mrIID))
	if err != nil {
		return FromGitLabResponse(err, resp)
	}
	return nil
}

// GetMergeRequestApprovals はMRの承認状態を取得する
func (c *Client) GetMergeRequestApprovals(projectID string, mrIID int) (*gogitlab.MergeRequestApprovals, error) {
	approvals, resp, err := c.client.MergeRequestApprovals.GetConfiguration(projectID, int64(mrIID))
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return approvals, nil
}
