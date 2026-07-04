package gitea

import sdk "code.gitea.io/sdk/gitea"

func (c *Client) ListRuns(workflowName string) ([]*sdk.ActionWorkflowRun, error) {
	// SDK has no workflow-name filter; fetch all and filter client-side if needed.
	resp, _, err := c.client.ListRepoActionRuns(c.owner, c.repo, sdk.ListRepoActionRunsOptions{})
	if err != nil {
		return nil, err
	}
	if workflowName == "" {
		return resp.WorkflowRuns, nil
	}
	var filtered []*sdk.ActionWorkflowRun
	for _, r := range resp.WorkflowRuns {
		if r.DisplayTitle == workflowName || r.Path == workflowName {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}
