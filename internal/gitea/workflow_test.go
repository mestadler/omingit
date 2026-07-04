package gitea

import (
	"testing"

	sdk "code.gitea.io/sdk/gitea"
)

func TestClient_ListWorkflows(t *testing.T) {
	t.Run("method signature", func(t *testing.T) {
		var fn func(*Client) ([]*sdk.ActionWorkflow, error) = (*Client).ListWorkflows
		_ = fn
	})
}

func TestClient_DispatchWorkflow(t *testing.T) {
	t.Run("method signature", func(t *testing.T) {
		var fn func(*Client, string, string) error = (*Client).DispatchWorkflow
		_ = fn
	})
}
