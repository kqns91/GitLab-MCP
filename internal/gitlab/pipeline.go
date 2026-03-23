package gitlab

import (
	"io"

	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

// maxJobLogSize はジョブログの最大サイズ（100KB）
const maxJobLogSize = 100 * 1024

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

// ListProjectPipelinesOptions はパイプライン一覧取得のオプション
type ListProjectPipelinesOptions struct {
	Status  *string
	Ref     *string
	SHA     *string
	Source  *string
	Page    int
	PerPage int
}

// ListProjectPipelines はプロジェクトのパイプライン一覧を取得する
func (c *Client) ListProjectPipelines(projectID string, opts *ListProjectPipelinesOptions) ([]*gogitlab.PipelineInfo, error) {
	page, perPage := 1, 100
	if opts != nil {
		if opts.Page > 0 {
			page = opts.Page
		}
		if opts.PerPage > 0 {
			perPage = opts.PerPage
		}
	}

	listOpts := &gogitlab.ListProjectPipelinesOptions{
		ListOptions: gogitlab.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}

	if opts != nil {
		if opts.Status != nil {
			status := gogitlab.BuildStateValue(*opts.Status)
			listOpts.Status = &status
		}
		if opts.Ref != nil {
			listOpts.Ref = opts.Ref
		}
		if opts.SHA != nil {
			listOpts.SHA = opts.SHA
		}
		if opts.Source != nil {
			listOpts.Source = opts.Source
		}
	}

	pipelines, resp, err := c.client.Pipelines.ListProjectPipelines(projectID, listOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return pipelines, nil
}

// GetPipeline はパイプラインの詳細を取得する
func (c *Client) GetPipeline(projectID string, pipelineID int) (*gogitlab.Pipeline, error) {
	pipeline, resp, err := c.client.Pipelines.GetPipeline(projectID, int64(pipelineID))
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return pipeline, nil
}

// CreatePipelineOptions はパイプライン作成のオプション
type CreatePipelineOptions struct {
	Ref       string
	Variables []PipelineVariable
}

// PipelineVariable はパイプライン変数
type PipelineVariable struct {
	Key   string
	Value string
}

// CreatePipeline は新しいパイプラインを作成する
func (c *Client) CreatePipeline(projectID string, opts *CreatePipelineOptions) (*gogitlab.Pipeline, error) {
	createOpts := &gogitlab.CreatePipelineOptions{
		Ref: &opts.Ref,
	}

	if len(opts.Variables) > 0 {
		vars := make([]*gogitlab.PipelineVariableOptions, len(opts.Variables))
		for i, v := range opts.Variables {
			key := v.Key
			value := v.Value
			vars[i] = &gogitlab.PipelineVariableOptions{
				Key:   &key,
				Value: &value,
			}
		}
		createOpts.Variables = &vars
	}

	pipeline, resp, err := c.client.Pipelines.CreatePipeline(projectID, createOpts)
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return pipeline, nil
}

// RetryPipeline はパイプラインの失敗ジョブを再試行する
func (c *Client) RetryPipeline(projectID string, pipelineID int) (*gogitlab.Pipeline, error) {
	pipeline, resp, err := c.client.Pipelines.RetryPipelineBuild(projectID, int64(pipelineID))
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return pipeline, nil
}

// CancelPipeline はパイプラインをキャンセルする
func (c *Client) CancelPipeline(projectID string, pipelineID int) (*gogitlab.Pipeline, error) {
	pipeline, resp, err := c.client.Pipelines.CancelPipelineBuild(projectID, int64(pipelineID))
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return pipeline, nil
}

// GetJob はジョブの詳細を取得する
func (c *Client) GetJob(projectID string, jobID int) (*gogitlab.Job, error) {
	job, resp, err := c.client.Jobs.GetJob(projectID, int64(jobID))
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return job, nil
}

// GetJobTrace はジョブのログを取得する
func (c *Client) GetJobTrace(projectID string, jobID int) (string, error) {
	reader, resp, err := c.client.Jobs.GetTraceFile(projectID, int64(jobID))
	if err != nil {
		return "", FromGitLabResponse(err, resp)
	}

	// ログサイズを制限する
	limitedReader := io.LimitReader(reader, maxJobLogSize)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", &MCPError{
			Code:    ErrCodeServerError,
			Message: "ジョブログの読み取りに失敗しました",
		}
	}

	return string(data), nil
}

// RetryJob はジョブを再試行する
func (c *Client) RetryJob(projectID string, jobID int) (*gogitlab.Job, error) {
	job, resp, err := c.client.Jobs.RetryJob(projectID, int64(jobID))
	if err != nil {
		return nil, FromGitLabResponse(err, resp)
	}
	return job, nil
}
