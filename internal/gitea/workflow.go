package gitea

import sdk "code.gitea.io/sdk/gitea"

func (c *Client) ListWorkflows() ([]*sdk.ActionWorkflow, error) {
	resp, _, err := c.client.ListRepoActionWorkflows(c.owner, c.repo)
	if err != nil {
		return nil, err
	}
	return resp.Workflows, nil
}

func (c *Client) DispatchWorkflow(workflowID, ref string) error {
	_, _, err := c.client.DispatchRepoActionWorkflow(c.owner, c.repo, workflowID,
		sdk.CreateActionWorkflowDispatchOption{Ref: ref}, false)
	return err
}
