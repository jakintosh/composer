package modal

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
	"composer/pkg/ui/components/button"
	modalshell "composer/pkg/ui/components/modalshell"
)

//go:embed modal.tmpl
var modalTemplate string

var tmpl = templates.New(
	"workflow_modal",
	"internal/ui/components/workflow_modal/modal.tmpl",
	modalTemplate,
	nil,
)

//go:embed body.tmpl
var bodyTemplate string

var bodyTmpl = templates.New(
	"workflow_modal_body",
	"internal/ui/components/workflow_modal/body.tmpl",
	bodyTemplate,
	nil,
)

// Props holds the controls displayed in the workflow creation modal.
type Props struct {
	AddStepButton button.Props
}

// RenderShell renders the shared modal chrome with the workflow form body.
func (p Props) RenderShell() template.HTML {
	return modalshell.MustRender(modalshell.Props{
		ID:         "workflow-modal",
		Title:      "Create Workflow",
		CloseLabel: "Close create workflow form",
		Body:       p.renderBody(),
	})
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

func (p Props) renderBody() template.HTML {
	return templates.SafeHTML(bodyTmpl.Render(p))
}
