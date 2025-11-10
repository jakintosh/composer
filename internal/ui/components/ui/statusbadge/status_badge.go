package statusbadge

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
)

//go:embed status_badge.tmpl
var statusBadgeTemplate string

var tmpl = templates.New(
	"ui_status_badge",
	"components/ui/statusbadge/status_badge.tmpl",
	statusBadgeTemplate,
	nil,
)

// Props defines the content and variant for a status badge.
type Props struct {
	Label   string
	Variant string
}

// Render executes the status badge template.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender renders the status badge, returning an HTML comment on failure to
// keep parent components resilient.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
