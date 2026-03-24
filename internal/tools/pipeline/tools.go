package pipeline

import (
	"context"

	"github.com/kqns91/gitlab-mcp/internal/gitlab"
	"github.com/kqns91/gitlab-mcp/internal/registry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListPipelinesInput は list_merge_request_pipelines の入力パラメータ
type ListPipelinesInput struct {
	ProjectID       string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	MergeRequestIID int    `json:"merge_request_iid" jsonschema:"description:Merge Request IID"`
}

// PipelineInfo はパイプライン情報
type PipelineInfo struct {
	ID        int64  `json:"id"`
	Status    string `json:"status"`
	Ref       string `json:"ref"`
	SHA       string `json:"sha"`
	WebURL    string `json:"web_url"`
	CreatedAt string `json:"created_at,omitempty"`
}

// ListPipelinesOutput は list_merge_request_pipelines の出力
type ListPipelinesOutput struct {
	Pipelines []PipelineInfo `json:"pipelines"`
}

// GetJobsInput は get_pipeline_jobs の入力パラメータ
type GetJobsInput struct {
	ProjectID  string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	PipelineID int    `json:"pipeline_id" jsonschema:"description:Pipeline ID"`
	Page       int    `json:"page,omitempty" jsonschema:"description:Page number (default: 1)"`
	PerPage    int    `json:"per_page,omitempty" jsonschema:"description:Number of items per page (default: 100, max: 100)"`
}

// JobInfo はジョブ情報
type JobInfo struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Stage  string `json:"stage"`
	Status string `json:"status"`
}

// GetJobsOutput は get_pipeline_jobs の出力
type GetJobsOutput struct {
	Jobs []JobInfo `json:"jobs"`
}

// ListProjectPipelinesInput は list_project_pipelines の入力パラメータ
type ListProjectPipelinesInput struct {
	ProjectID string  `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	Status    *string `json:"status,omitempty" jsonschema:"enum:created,enum:waiting_for_resource,enum:preparing,enum:pending,enum:running,enum:success,enum:failed,enum:canceled,enum:skipped,enum:manual,enum:scheduled,description:Pipeline status filter"`
	Ref       *string `json:"ref,omitempty" jsonschema:"description:Branch or tag name filter"`
	SHA       *string `json:"sha,omitempty" jsonschema:"description:Commit SHA filter"`
	Source    *string `json:"source,omitempty" jsonschema:"description:Pipeline source filter (e.g. push, web, trigger)"`
	Page      int     `json:"page,omitempty" jsonschema:"description:Page number (default: 1)"`
	PerPage   int     `json:"per_page,omitempty" jsonschema:"description:Number of items per page (default: 100, max: 100)"`
}

// ListProjectPipelinesOutput は list_project_pipelines の出力
type ListProjectPipelinesOutput struct {
	Pipelines []PipelineInfo `json:"pipelines"`
}

// GetPipelineInput は get_pipeline の入力パラメータ
type GetPipelineInput struct {
	ProjectID  string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	PipelineID int    `json:"pipeline_id" jsonschema:"description:Pipeline ID"`
}

// PipelineDetail はパイプライン詳細情報
type PipelineDetail struct {
	ID         int64  `json:"id"`
	Status     string `json:"status"`
	Ref        string `json:"ref"`
	SHA        string `json:"sha"`
	WebURL     string `json:"web_url"`
	Source     string `json:"source"`
	Duration   int64  `json:"duration,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	StartedAt  string `json:"started_at,omitempty"`
	FinishedAt string `json:"finished_at,omitempty"`
}

// GetPipelineOutput は get_pipeline の出力
type GetPipelineOutput = PipelineDetail

// CreatePipelineInput は create_pipeline の入力パラメータ
type CreatePipelineInput struct {
	ProjectID string                  `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	Ref       string                  `json:"ref" jsonschema:"description:Branch or tag name to create the pipeline for"`
	Variables []PipelineVariableInput `json:"variables,omitempty" jsonschema:"description:Pipeline variables"`
}

// PipelineVariableInput はパイプライン変数の入力
type PipelineVariableInput struct {
	Key   string `json:"key" jsonschema:"description:Variable key"`
	Value string `json:"value" jsonschema:"description:Variable value"`
}

// CreatePipelineOutput は create_pipeline の出力
type CreatePipelineOutput struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
	WebURL string `json:"web_url"`
}

// RetryPipelineInput は retry_pipeline の入力パラメータ
type RetryPipelineInput struct {
	ProjectID  string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	PipelineID int    `json:"pipeline_id" jsonschema:"description:Pipeline ID"`
}

// RetryPipelineOutput は retry_pipeline の出力
type RetryPipelineOutput struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
	WebURL string `json:"web_url"`
}

// CancelPipelineInput は cancel_pipeline の入力パラメータ
type CancelPipelineInput struct {
	ProjectID  string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	PipelineID int    `json:"pipeline_id" jsonschema:"description:Pipeline ID"`
}

// CancelPipelineOutput は cancel_pipeline の出力
type CancelPipelineOutput struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
	WebURL string `json:"web_url"`
}

// GetPipelineJobInput は get_pipeline_job の入力パラメータ
type GetPipelineJobInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	JobID     int    `json:"job_id" jsonschema:"description:Job ID"`
}

// JobDetail はジョブ詳細情報
type JobDetail struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Stage         string  `json:"stage"`
	Status        string  `json:"status"`
	Ref           string  `json:"ref"`
	WebURL        string  `json:"web_url"`
	Duration      float64 `json:"duration,omitempty"`
	FailureReason string  `json:"failure_reason,omitempty"`
	CreatedAt     string  `json:"created_at,omitempty"`
	StartedAt     string  `json:"started_at,omitempty"`
	FinishedAt    string  `json:"finished_at,omitempty"`
}

// GetPipelineJobOutput は get_pipeline_job の出力
type GetPipelineJobOutput = JobDetail

// GetJobLogInput は get_job_log の入力パラメータ
type GetJobLogInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	JobID     int    `json:"job_id" jsonschema:"description:Job ID"`
}

// GetJobLogOutput は get_job_log の出力
type GetJobLogOutput struct {
	Log string `json:"log"`
}

// RetryPipelineJobInput は retry_pipeline_job の入力パラメータ
type RetryPipelineJobInput struct {
	ProjectID string `json:"project_id" jsonschema:"description:Project ID or URL-encoded path"`
	JobID     int    `json:"job_id" jsonschema:"description:Job ID"`
}

// RetryPipelineJobOutput は retry_pipeline_job の出力
type RetryPipelineJobOutput struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	WebURL string `json:"web_url"`
}

// clientHolder holds the GitLab client for handlers
type clientHolder struct {
	client *gitlab.Client
}

var holder *clientHolder

// Register はパイプライン関連ツールを登録する
func Register(reg *registry.Registry, client *gitlab.Client) {
	holder = &clientHolder{client: client}

	registry.RegisterTool(reg, "list_merge_request_pipelines",
		"GitLab Merge Request に関連するパイプライン一覧を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ListPipelinesInput) (*mcp.CallToolResult, ListPipelinesOutput, error) {
			return listPipelinesHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "get_pipeline_jobs",
		"GitLab パイプラインのジョブ一覧を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input GetJobsInput) (*mcp.CallToolResult, GetJobsOutput, error) {
			return getJobsHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "list_project_pipelines",
		"GitLab プロジェクトのパイプライン一覧を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input ListProjectPipelinesInput) (*mcp.CallToolResult, ListProjectPipelinesOutput, error) {
			return listProjectPipelinesHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "get_pipeline",
		"GitLab パイプラインの詳細情報を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input GetPipelineInput) (*mcp.CallToolResult, GetPipelineOutput, error) {
			return getPipelineHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "create_pipeline",
		"GitLab で新しいパイプラインを作成します",
		func(ctx context.Context, req *mcp.CallToolRequest, input CreatePipelineInput) (*mcp.CallToolResult, CreatePipelineOutput, error) {
			return createPipelineHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "retry_pipeline",
		"GitLab パイプラインの失敗したジョブを再試行します",
		func(ctx context.Context, req *mcp.CallToolRequest, input RetryPipelineInput) (*mcp.CallToolResult, RetryPipelineOutput, error) {
			return retryPipelineHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "cancel_pipeline",
		"GitLab パイプラインをキャンセルします",
		func(ctx context.Context, req *mcp.CallToolRequest, input CancelPipelineInput) (*mcp.CallToolResult, CancelPipelineOutput, error) {
			return cancelPipelineHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "get_pipeline_job",
		"GitLab ジョブの詳細情報を取得します",
		func(ctx context.Context, req *mcp.CallToolRequest, input GetPipelineJobInput) (*mcp.CallToolResult, GetPipelineJobOutput, error) {
			return getPipelineJobHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "get_job_log",
		"GitLab ジョブのログを取得します（最大100KB）",
		func(ctx context.Context, req *mcp.CallToolRequest, input GetJobLogInput) (*mcp.CallToolResult, GetJobLogOutput, error) {
			return getJobLogHandler(holder.client, ctx, req, input)
		})

	registry.RegisterTool(reg, "retry_pipeline_job",
		"GitLab ジョブを再試行します",
		func(ctx context.Context, req *mcp.CallToolRequest, input RetryPipelineJobInput) (*mcp.CallToolResult, RetryPipelineJobOutput, error) {
			return retryPipelineJobHandler(holder.client, ctx, req, input)
		})
}

func listPipelinesHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ListPipelinesInput) (*mcp.CallToolResult, ListPipelinesOutput, error) {
	pipelines, err := client.ListMergeRequestPipelines(input.ProjectID, input.MergeRequestIID)
	if err != nil {
		return nil, ListPipelinesOutput{}, err
	}

	infos := make([]PipelineInfo, len(pipelines))
	for i, p := range pipelines {
		createdAt := ""
		if p.CreatedAt != nil {
			createdAt = p.CreatedAt.String()
		}
		infos[i] = PipelineInfo{
			ID:        p.ID,
			Status:    p.Status,
			Ref:       p.Ref,
			SHA:       p.SHA,
			WebURL:    p.WebURL,
			CreatedAt: createdAt,
		}
	}

	return nil, ListPipelinesOutput{Pipelines: infos}, nil
}

func getJobsHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input GetJobsInput) (*mcp.CallToolResult, GetJobsOutput, error) {
	var pagination *gitlab.PaginationOptions
	if input.Page > 0 || input.PerPage > 0 {
		pagination = &gitlab.PaginationOptions{
			Page:    input.Page,
			PerPage: input.PerPage,
		}
	}

	jobs, err := client.ListPipelineJobs(input.ProjectID, input.PipelineID, pagination)
	if err != nil {
		return nil, GetJobsOutput{}, err
	}

	infos := make([]JobInfo, len(jobs))
	for i, j := range jobs {
		infos[i] = JobInfo{
			ID:     int64(j.ID),
			Name:   j.Name,
			Stage:  j.Stage,
			Status: j.Status,
		}
	}

	return nil, GetJobsOutput{Jobs: infos}, nil
}

func listProjectPipelinesHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input ListProjectPipelinesInput) (*mcp.CallToolResult, ListProjectPipelinesOutput, error) {
	opts := &gitlab.ListProjectPipelinesOptions{
		Status:  input.Status,
		Ref:     input.Ref,
		SHA:     input.SHA,
		Source:  input.Source,
		Page:    input.Page,
		PerPage: input.PerPage,
	}

	pipelines, err := client.ListProjectPipelines(input.ProjectID, opts)
	if err != nil {
		return nil, ListProjectPipelinesOutput{}, err
	}

	infos := make([]PipelineInfo, len(pipelines))
	for i, p := range pipelines {
		createdAt := ""
		if p.CreatedAt != nil {
			createdAt = p.CreatedAt.String()
		}
		infos[i] = PipelineInfo{
			ID:        p.ID,
			Status:    p.Status,
			Ref:       p.Ref,
			SHA:       p.SHA,
			WebURL:    p.WebURL,
			CreatedAt: createdAt,
		}
	}

	return nil, ListProjectPipelinesOutput{Pipelines: infos}, nil
}

func getPipelineHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input GetPipelineInput) (*mcp.CallToolResult, GetPipelineOutput, error) {
	p, err := client.GetPipeline(input.ProjectID, input.PipelineID)
	if err != nil {
		return nil, GetPipelineOutput{}, err
	}

	createdAt := ""
	if p.CreatedAt != nil {
		createdAt = p.CreatedAt.String()
	}
	startedAt := ""
	if p.StartedAt != nil {
		startedAt = p.StartedAt.String()
	}
	finishedAt := ""
	if p.FinishedAt != nil {
		finishedAt = p.FinishedAt.String()
	}

	return nil, GetPipelineOutput{
		ID:         p.ID,
		Status:     p.Status,
		Ref:        p.Ref,
		SHA:        p.SHA,
		WebURL:     p.WebURL,
		Source:     string(p.Source),
		Duration:   p.Duration,
		CreatedAt:  createdAt,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}, nil
}

func createPipelineHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input CreatePipelineInput) (*mcp.CallToolResult, CreatePipelineOutput, error) {
	vars := make([]gitlab.PipelineVariable, len(input.Variables))
	for i, v := range input.Variables {
		vars[i] = gitlab.PipelineVariable{
			Key:   v.Key,
			Value: v.Value,
		}
	}

	opts := &gitlab.CreatePipelineOptions{
		Ref:       input.Ref,
		Variables: vars,
	}

	p, err := client.CreatePipeline(input.ProjectID, opts)
	if err != nil {
		return nil, CreatePipelineOutput{}, err
	}

	return nil, CreatePipelineOutput{
		ID:     p.ID,
		Status: p.Status,
		WebURL: p.WebURL,
	}, nil
}

func retryPipelineHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input RetryPipelineInput) (*mcp.CallToolResult, RetryPipelineOutput, error) {
	p, err := client.RetryPipeline(input.ProjectID, input.PipelineID)
	if err != nil {
		return nil, RetryPipelineOutput{}, err
	}

	return nil, RetryPipelineOutput{
		ID:     p.ID,
		Status: p.Status,
		WebURL: p.WebURL,
	}, nil
}

func cancelPipelineHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input CancelPipelineInput) (*mcp.CallToolResult, CancelPipelineOutput, error) {
	p, err := client.CancelPipeline(input.ProjectID, input.PipelineID)
	if err != nil {
		return nil, CancelPipelineOutput{}, err
	}

	return nil, CancelPipelineOutput{
		ID:     p.ID,
		Status: p.Status,
		WebURL: p.WebURL,
	}, nil
}

func getPipelineJobHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input GetPipelineJobInput) (*mcp.CallToolResult, GetPipelineJobOutput, error) {
	j, err := client.GetJob(input.ProjectID, input.JobID)
	if err != nil {
		return nil, GetPipelineJobOutput{}, err
	}

	createdAt := ""
	if j.CreatedAt != nil {
		createdAt = j.CreatedAt.String()
	}
	startedAt := ""
	if j.StartedAt != nil {
		startedAt = j.StartedAt.String()
	}
	finishedAt := ""
	if j.FinishedAt != nil {
		finishedAt = j.FinishedAt.String()
	}

	return nil, GetPipelineJobOutput{
		ID:            int64(j.ID),
		Name:          j.Name,
		Stage:         j.Stage,
		Status:        j.Status,
		Ref:           j.Ref,
		WebURL:        j.WebURL,
		Duration:      j.Duration,
		FailureReason: j.FailureReason,
		CreatedAt:     createdAt,
		StartedAt:     startedAt,
		FinishedAt:    finishedAt,
	}, nil
}

func getJobLogHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input GetJobLogInput) (*mcp.CallToolResult, GetJobLogOutput, error) {
	log, err := client.GetJobTrace(input.ProjectID, input.JobID)
	if err != nil {
		return nil, GetJobLogOutput{}, err
	}

	return nil, GetJobLogOutput{Log: log}, nil
}

func retryPipelineJobHandler(client *gitlab.Client, ctx context.Context, req *mcp.CallToolRequest, input RetryPipelineJobInput) (*mcp.CallToolResult, RetryPipelineJobOutput, error) {
	j, err := client.RetryJob(input.ProjectID, input.JobID)
	if err != nil {
		return nil, RetryPipelineJobOutput{}, err
	}

	return nil, RetryPipelineJobOutput{
		ID:     int64(j.ID),
		Name:   j.Name,
		Status: j.Status,
		WebURL: j.WebURL,
	}, nil
}
