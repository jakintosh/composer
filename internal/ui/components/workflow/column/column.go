package column

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/components/ui/columnheader"
	"composer/internal/ui/templates"
)

//go:embed column.tmpl
var columnTemplate string

var tmpl = templates.New(
	"workflow_column",
	"components/workflow/column/column.tmpl",
	columnTemplate,
	nil,
)

// Workflow represents a single workflow card in the column.
type Workflow struct {
	DisplayName string
	ID          string
	Description string
	Message     string
	StepNames   []string
}

// Props describes the workflow column on the dashboard.
type Props struct {
	Header    columnheader.Props
	Workflows []Workflow
}

// RenderHeader renders the nested column header.
func (p Props) RenderHeader() template.HTML {
	return templates.SafeHTML(columnheader.Render(p.Header))
}

// Render executes the component template.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts failures into HTML comments.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
