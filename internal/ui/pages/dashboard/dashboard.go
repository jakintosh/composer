package dashboard

import (
	"bytes"
	"io"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/button"
	runcomponent "composer/internal/ui/components/run"
	sidebar "composer/internal/ui/components/sidebar"
	waitingcolumn "composer/internal/ui/components/waiting_column"
	workflowcomponent "composer/internal/ui/components/workflow"
)

// Props aggregates the components required to render the dashboard page.
type Props struct {
	Sidebar        sidebar.Props
	WorkflowColumn workflowcomponent.ColumnProps
	WorkflowModal  workflowcomponent.ModalProps
	RunColumn      runcomponent.ColumnProps
	RunModal       runcomponent.ModalProps
	TaskColumn     waitingcolumn.Props
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
func DefaultWorkflowModal() workflowcomponent.ModalProps {
	return workflowcomponent.ModalProps{
		AddStepButton: buttonProps(),
	}
}

// DefaultRunModal returns an empty run modal props value for convenience.
func DefaultRunModal() runcomponent.ModalProps {
	return runcomponent.ModalProps{}
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
				sidebar.Sidebar(p.Sidebar),
			),
			html.Main(
				html.Class("ui-shell__main"),
				html.Div(
					html.Class("ui-shell__content"),
					html.H1(g.Text("Workflow Dashboard")),
					html.Div(
						html.Class("panel-grid"),
						workflowcomponent.Column(p.WorkflowColumn),
						runcomponent.Column(p.RunColumn),
						waitingcolumn.Column(p.TaskColumn),
					),
					runcomponent.Modal(p.RunModal),
					workflowcomponent.Modal(p.WorkflowModal),
				),
			),
		),
		html.Script(
			html.Src("/static/dashboard.js"),
			html.Defer(),
		),
	)
}

func buttonProps() button.Props {
	return button.Props{
		ID:       "add-workflow-step",
		Class:    "button--outline button--sm",
		Label:    "Add Step",
		Type:     "button",
		IconSize: 16,
	}
}
