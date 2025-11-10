package modal

import (
	"testing"

	"composer/pkg/ui/components/button"
	"gotest.tools/v3/golden"
)

func TestRenderWorkflowModal(t *testing.T) {
	props := Props{AddStepButton: button.Props{Label: "Add Step"}}

	html, err := Render(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "modal.golden")
}
