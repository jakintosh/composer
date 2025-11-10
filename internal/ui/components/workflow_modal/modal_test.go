package modal

import (
	"testing"

	"composer/internal/ui/components/button"
	"composer/internal/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderWorkflowModal(t *testing.T) {
	props := Props{AddStepButton: button.Props{Label: "Add Step"}}

	html := testutil.Render(t, Modal(props))
	golden.Assert(t, html, "modal.golden")
}
