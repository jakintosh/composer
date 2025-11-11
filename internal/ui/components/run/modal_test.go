package run

import (
	"testing"

	"composer/internal/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderRunModal(t *testing.T) {
	html := testutil.Render(t, Modal(ModalProps{}))
	golden.Assert(t, html, "modal.golden")
}
