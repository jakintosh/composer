package statusbadge

import (
	"testing"

	"gotest.tools/v3/golden"
)

func TestRenderStatusBadgeDefault(t *testing.T) {
	html, err := Render(Props{Label: "pending"})
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	golden.Assert(t, string(html), "status_badge_default.golden")
}

func TestRenderStatusBadgeVariant(t *testing.T) {
	html, err := Render(Props{Label: "ready", Variant: "status-badge--ready"})
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	golden.Assert(t, string(html), "status_badge_variant.golden")
}
