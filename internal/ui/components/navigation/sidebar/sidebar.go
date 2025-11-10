package sidebar

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
)

//go:embed sidebar.tmpl
var sidebarTemplate string

var tmpl = templates.New(
	"navigation_sidebar",
	"components/navigation/sidebar/sidebar.tmpl",
	sidebarTemplate,
	nil,
)

// Link denotes a navigational item in the sidebar menu.
type Link struct {
	Label  string
	Href   string
	Active bool
}

// Props represents the layout of the primary navigation sidebar.
type Props struct {
	Title string
	Links []Link
}

// Render produces the sidebar HTML fragment.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts failures into an HTML comment.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
