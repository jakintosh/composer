package column

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
	"composer/pkg/ui/components/columnheader"
)

//go:embed column.tmpl
var columnTemplate string

var tmpl = templates.New(
	"waiting_column",
	"internal/ui/components/waiting_column/column.tmpl",
	columnTemplate,
	nil,
)

// Task represents a single waiting task awaiting user interaction.
type Task struct {
	Name        string
	Description string
	Prompt      string
}

// Group aggregates pending human tasks for a specific run.
type Group struct {
	RunID          string
	RunDisplayName string
	WorkflowName   string
	TaskCount      int
	Tasks          []Task
}

// Props represents the waiting-task column on the dashboard.
type Props struct {
	Header columnheader.Props
	Groups []Group
}

// RenderHeader renders the nested header component.
func (p Props) RenderHeader() template.HTML {
	return templates.SafeHTML(columnheader.Render(p.Header))
}

// Render executes the template for the waiting column component.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts failures into HTML comments.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
