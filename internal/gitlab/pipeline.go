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
func (c *Client) ListPipelineJobs(projectID string, pipelineID int, pagination *PaginationOptions) ([]*gogitlab.Job, error) {
	page, perPage := 1, 100
	if pagination != nil {
		if pagination.Page > 0 {
			page = pagination.Page
		}
		if pagination.PerPage > 0 {
			perPage = pagination.PerPage
		}
	}

	opts := &gogitlab.ListJobsOptions{
		ListOptions: gogitlab.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}
	jobs, resp, err := c.client.Jobs.ListPipelineJobs(projectID, int64(pipelineID), opts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return jobs, nil
}
