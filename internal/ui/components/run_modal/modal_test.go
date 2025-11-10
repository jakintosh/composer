package modal

import (
	"testing"

	"gotest.tools/v3/golden"
)

func TestRenderRunModal(t *testing.T) {
	html, err := Render(Props{})
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "modal.golden")
}
