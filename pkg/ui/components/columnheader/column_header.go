package columnheader

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
	"composer/pkg/ui/components/button"
)

//go:embed column_header.tmpl
var columnHeaderTemplate string

var tmpl = templates.New(
	"ui_column_header",
	"pkg/ui/components/columnheader/column_header.tmpl",
	columnHeaderTemplate,
	template.FuncMap{
		"renderButton": button.MustRender,
	},
)

// Props defines the structure of a reusable column header.
type Props struct {
	Title   string
	Actions []button.Props
}

// Render executes the template for the column header component.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts errors into HTML comments so parent
// components can embed the header without manual error handling.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
