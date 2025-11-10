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
	"run_column",
	"components/run/column/column.tmpl",
	columnTemplate,
	nil,
)

// Step represents a single workflow step state within a run.
type Step struct {
	Name        string
	Status      string
	StatusClass string
}

// Run summarizes a workflow run and its current state.
type Run struct {
	DisplayName  string
	ID           string
	StateLabel   string
	StateClass   string
	WorkflowName string
	Steps        []Step
}

// Props describes the runs column rendered on the dashboard.
type Props struct {
	Header columnheader.Props
	Runs   []Run
}

// RenderHeader renders the nested header component.
func (p Props) RenderHeader() template.HTML {
	return templates.SafeHTML(columnheader.Render(p.Header))
}

// Render executes the template for the run column component.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts failures into HTML comments.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
