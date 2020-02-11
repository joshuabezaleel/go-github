// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Step represents a single task from a sequence of tasks of a job.
type Step struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	Number      int64     `json:"number"`
	StartedAt   Timestamp `json:"started_at"`
	CompletedAt Timestamp `json:"completed_at"`
}

// Job represents a repository action workflow job.
type Job struct {
	ID          int64     `json:"id"`
	RunID       int64     `json:"run_id"`
	RunURL      string    `json:"run_url"`
	NodeID      string    `json:"node_id"`
	HeadSHA     string    `json:"head_sha"`
	URL         string    `json:"url"`
	HTMLURL     string    `json:"html_url"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	StartedAt   Timestamp `json:"started_at"`
	CompletedAt Timestamp `json:"completed_at"`
	Name        string    `json:"name"`
	Steps       []*Step   `json:"steps"`
	CheckRunURL string    `json:"check_run_url"`
}

// Jobs represents a slice of repository action workflow job.
type Jobs struct {
	TotalCount int    `json:"total_count"`
	Jobs       []*Job `json:"jobs"`
}

// ListWorkflowJobs lists all jobs for a workflow run.
//
// GitHub API docs: https://developer.github.com/v3/actions/workflow_jobs/#list-jobs-for-a-workflow-run
func (s *ActionsService) ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *ListOptions) (*Jobs, *Response, error) {
	u := fmt.Sprintf("repos/%s/%s/actions/runs/%v/jobs", owner, repo, runID)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	jobs := new(Jobs)
	resp, err := s.client.Do(ctx, req, &jobs)
	if err != nil {
		return nil, resp, err
	}

	return jobs, resp, nil
}

// GetWorkflowJobByID gets a specific job in a workflow run by ID.
//
// GitHub API docs: https://developer.github.com/v3/actions/workflow_jobs/#list-jobs-for-a-workflow-run
func (s *ActionsService) GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*Job, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/jobs/%v", owner, repo, jobID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	job := new(Job)
	resp, err := s.client.Do(ctx, req, job)
	if err != nil {
		return nil, resp, err
	}

	return job, resp, nil
}

// GetWorkflowJobLogs gets a redirect URL to a download a plain text file of logs for a workflow job.
//
// GitHub API docs: https://developer.github.com/v3/actions/workflow_jobs/#list-workflow-job-logs
func (s *ActionsService) GetWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64, followRedirects bool) (*url.URL, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/jobs/%v/logs", owner, repo, jobID)

	resp, err := s.getWorkflowJobLogsFromURL(ctx, u, followRedirects)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusFound {
		return nil, newResponse(resp), fmt.Errorf("unexpected status code: %s", resp.Status)
	}
	parsedURL, err := url.Parse(resp.Header.Get("Location"))
	return parsedURL, newResponse(resp), err
}

func (s *ActionsService) getWorkflowJobLogsFromURL(ctx context.Context, u string, followRedirects bool) (*http.Response, error) {
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	// Use http.DefaultTransport if no custom Transport is configured
	req = withContext(ctx, req)
	if s.client.client.Transport == nil {
		resp, err = http.DefaultTransport.RoundTrip(req)
	} else {
		resp, err = s.client.client.Transport.RoundTrip(req)
	}
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// If redirect response is returned, follow it
	if followRedirects && resp.StatusCode == http.StatusMovedPermanently {
		u = resp.Header.Get("Location")
		resp, err = s.getWorkflowJobLogsFromURL(ctx, u, false)
	}
	return resp, err

}
