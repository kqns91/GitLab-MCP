package gitlab

import (
	"errors"

	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

// Client は GitLab API クライアントのラッパー
type Client struct {
	client *gogitlab.Client
}

// NewClient は新しい GitLab クライアントを作成する
func NewClient(baseURL, token string) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("GitLab URL is required")
	}
	if token == "" {
		return nil, errors.New("GitLab token is required")
	}

	client, err := gogitlab.NewClient(token, gogitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

// MergeRequests returns the MergeRequestsService
func (c *Client) MergeRequests() gogitlab.MergeRequestsServiceInterface {
	return c.client.MergeRequests
}

// Discussions returns the DiscussionsService
func (c *Client) Discussions() gogitlab.DiscussionsServiceInterface {
	return c.client.Discussions
}

// MergeRequestApprovals returns the MergeRequestApprovalsService
func (c *Client) MergeRequestApprovals() gogitlab.MergeRequestApprovalsServiceInterface {
	return c.client.MergeRequestApprovals
}

// Pipelines returns the PipelinesService
func (c *Client) Pipelines() gogitlab.PipelinesServiceInterface {
	return c.client.Pipelines
}

// Jobs returns the JobsService
func (c *Client) Jobs() gogitlab.JobsServiceInterface {
	return c.client.Jobs
}

// Notes returns the NotesService
func (c *Client) Notes() gogitlab.NotesServiceInterface {
	return c.client.Notes
}
