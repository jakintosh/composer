package collapsible

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
)

//go:embed collapsible.tmpl
var collapsibleTemplate string

var tmpl = templates.New(
	"ui_card_collapsible",
	"components/ui/card/collapsible/collapsible.tmpl",
	collapsibleTemplate,
	nil,
)

// Props describes a collapsible card with a summary row and body content.
type Props struct {
	Title        string
	SummaryItems []template.HTML
	Body         template.HTML
}

// HasSummaryItems reports whether summary content was provided.
func (p Props) HasSummaryItems() bool {
	return len(p.SummaryItems) > 0
}

// RenderBody returns the rendered HTML for the details region.
func (p Props) RenderBody() template.HTML {
	return p.Body
}

// Render executes the collapsible card template.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender renders the collapsible card, returning an HTML comment when
// rendering fails.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
