package gitlab

import (
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

// ListMergeRequestPipelines はMRに関連するパイプライン一覧を取得する
func (c *Client) ListMergeRequestPipelines(projectID string, mrIID int) ([]*gogitlab.PipelineInfo, error) {
	pipelines, resp, err := c.client.MergeRequests.ListMergeRequestPipelines(projectID, int64(mrIID))
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return pipelines, nil
}

// ListPipelineJobs はパイプラインのジョブ一覧を取得する
func (c *Client) ListPipelineJobs(projectID string, pipelineID int) ([]*gogitlab.Job, error) {
	jobs, resp, err := c.client.Jobs.ListPipelineJobs(projectID, int64(pipelineID), nil)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return jobs, nil
}
