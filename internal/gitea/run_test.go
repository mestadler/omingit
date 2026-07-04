package gitea

import (
	"testing"

	sdk "code.gitea.io/sdk/gitea"
)

func TestClient_ListRuns(t *testing.T) {
	t.Run("method signature", func(t *testing.T) {
		var fn func(*Client, string) ([]*sdk.ActionWorkflowRun, error) = (*Client).ListRuns
		_ = fn
	})
}
