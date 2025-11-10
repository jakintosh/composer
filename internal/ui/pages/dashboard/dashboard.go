package dashboard

import (
	_ "embed"
	"html/template"
	"io"

	"composer/internal/ui/components/navigation/sidebar"
	runcolumn "composer/internal/ui/components/run/column"
	runmodal "composer/internal/ui/components/run/modal"
	"composer/internal/ui/components/ui/button"
	waitingcolumn "composer/internal/ui/components/waiting/column"
	workflowcolumn "composer/internal/ui/components/workflow/column"
	workflowmodal "composer/internal/ui/components/workflow/modal"
	"composer/internal/ui/templates"
)

//go:embed dashboard.tmpl
var pageTemplate string

var tmpl = templates.New(
	"page_dashboard",
	"pages/dashboard/dashboard.tmpl",
	pageTemplate,
	nil,
)

// Props aggregates the components required to render the dashboard page.
type Props struct {
	Sidebar        sidebar.Props
	WorkflowColumn workflowcolumn.Props
	WorkflowModal  workflowmodal.Props
	RunColumn      runcolumn.Props
	RunModal       runmodal.Props
	TaskColumn     waitingcolumn.Props
}

// RenderSidebar renders the sidebar component for template composition.
func (p Props) RenderSidebar() template.HTML {
	return sidebar.MustRender(p.Sidebar)
}

// RenderWorkflowColumn renders the workflow column component.
func (p Props) RenderWorkflowColumn() template.HTML {
	return workflowcolumn.MustRender(p.WorkflowColumn)
}

// RenderRunColumn renders the run column component.
func (p Props) RenderRunColumn() template.HTML {
	return runcolumn.MustRender(p.RunColumn)
}

// RenderTaskColumn renders the waiting-task column component.
func (p Props) RenderTaskColumn() template.HTML {
	return waitingcolumn.MustRender(p.TaskColumn)
}

// RenderWorkflowModal renders the workflow modal component.
func (p Props) RenderWorkflowModal() template.HTML {
	return workflowmodal.MustRender(p.WorkflowModal)
}

// RenderRunModal renders the run modal component.
func (p Props) RenderRunModal() template.HTML {
	return runmodal.MustRender(p.RunModal)
}

// RenderPage writes the rendered dashboard to the supplied writer.
func RenderPage(w io.Writer, props Props) error {
	return tmpl.Execute(w, props)
}

// HTML renders the dashboard and returns the HTML fragment.
func HTML(props Props) (template.HTML, error) {
	return tmpl.Render(props)
}

// DefaultWorkflowModal returns the default modal configuration shared by the
// dashboard builder.
func DefaultWorkflowModal() workflowmodal.Props {
	return workflowmodal.Props{
		AddStepButton: button.Props{
			ID:       "add-workflow-step",
			Class:    "button--outline button--sm",
			Label:    "Add Step",
			Type:     "button",
			IconSize: 16,
		},
	}
}

// DefaultRunModal returns an empty run modal props value for convenience.
func DefaultRunModal() runmodal.Props {
	return runmodal.Props{}
}
