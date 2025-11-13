package pages

import (
	"composer/internal/ui/views"
	"composer/pkg/ui/components"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// DashboardProps aggregates the components required to render the dashboard page.
type DashboardProps struct {
	Sidebar        components.SidebarProps
	WorkflowColumn views.WorkflowColumnProps
	WorkflowModal  views.WorkflowModalProps
	RunColumn      views.RunColumnProps
	RunModal       views.RunModalProps
	TaskColumn     views.WaitingColumnProps
}

// Dashboard returns the root gomponent for the dashboard.
func Dashboard(p DashboardProps) g.Node {
	return html.Doctype(
		html.HTML(
			html.Lang("en"),
			head(),
			body(p),
		),
	)
}

// DefaultWorkflowModal returns the default modal configuration shared by the
// dashboard builder.
func DefaultWorkflowModal() views.WorkflowModalProps {
	return views.WorkflowModalProps{
		AddStepButton: buttonProps(),
	}
}

// DefaultRunModal returns an empty run modal props value for convenience.
func DefaultRunModal() views.RunModalProps {
	return views.RunModalProps{}
}

func head() g.Node {
	return html.Head(
		html.Meta(html.Charset("utf-8")),
		html.Meta(
			html.Name("viewport"),
			html.Content("width=device-width, initial-scale=1"),
		),
		html.TitleEl(g.Text("Composer Workflow Dashboard")),
		html.Link(
			html.Rel("stylesheet"),
			html.Href("/static/app.css"),
		),
	)
}

func body(p DashboardProps) g.Node {
	return html.Body(
		html.Div(
			html.Class("ui-shell"),
			html.Aside(
				html.Class("ui-shell__sidebar"),
				components.Sidebar(p.Sidebar),
			),
			html.Main(
				html.Class("ui-shell__main"),
				html.Div(
					html.Class("ui-shell__content"),
					html.H1(g.Text("Workflow Dashboard")),
					html.Div(
						html.Class("panel-grid"),
						views.WorkflowColumn(p.WorkflowColumn),
						p.RunColumn.Render(),
						views.WaitingColumn(p.TaskColumn),
					),
					p.RunModal.Render(),
					views.WorkflowModal(p.WorkflowModal),
				),
			),
		),
		html.Script(
			html.Src("/static/dashboard.js"),
			html.Defer(),
		),
	)
}

func buttonProps() components.ButtonProps {
	return components.ButtonProps{
		ID:       "add-workflow-step",
		Class:    "button--outline button--sm",
		Label:    "Add Step",
		Type:     "button",
		IconSize: 16,
	}
}
