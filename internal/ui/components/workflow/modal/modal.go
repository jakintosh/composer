package modal

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/components/ui/button"
	"composer/internal/ui/templates"
)

//go:embed modal.tmpl
var modalTemplate string

var tmpl = templates.New(
	"workflow_modal",
	"components/workflow/modal/modal.tmpl",
	modalTemplate,
	nil,
)

// Props holds the controls displayed in the workflow creation modal.
type Props struct {
	AddStepButton button.Props
}

// RenderAddStepButton renders the nested add-step button component.
func (p Props) RenderAddStepButton() template.HTML {
	return button.MustRender(p.AddStepButton)
}

// Render executes the modal template.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts failures into HTML comments.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
