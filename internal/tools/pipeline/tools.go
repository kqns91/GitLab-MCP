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
	jobs, err := client.ListPipelineJobs(input.ProjectID, input.PipelineID)
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
