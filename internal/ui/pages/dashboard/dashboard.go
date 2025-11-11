package dashboard

import (
	"bytes"
	"io"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/pkg/ui/components"
	"composer/internal/ui/views"
)

// Props aggregates the components required to render the dashboard page.
type Props struct {
	Sidebar        components.SidebarProps
	WorkflowColumn views.WorkflowColumnProps
	WorkflowModal  views.WorkflowModalProps
	RunColumn      views.RunColumnProps
	RunModal       views.RunModalProps
	TaskColumn     views.WaitingColumnProps
}

// Page returns the root gomponent for the dashboard.
func Page(p Props) g.Node {
	return html.Doctype(
		html.HTML(
			html.Lang("en"),
			head(),
			body(p),
		),
	)
}

// RenderPage writes the rendered dashboard to the supplied writer.
func RenderPage(w io.Writer, props Props) error {
	return Page(props).Render(w)
}

// HTML renders the dashboard and returns the HTML fragment.
func HTML(props Props) (string, error) {
	var buf bytes.Buffer
	if err := RenderPage(&buf, props); err != nil {
		return "", err
	}
	return buf.String(), nil
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

func body(p Props) g.Node {
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
						views.RunColumn(p.RunColumn),
						views.WaitingColumn(p.TaskColumn),
					),
					views.RunModal(p.RunModal),
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
