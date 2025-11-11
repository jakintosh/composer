package run

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/button"
	"composer/internal/ui/components/modal"
)

// ModalProps contains the data required to render the run modal.
type ModalProps struct{}

// Modal renders the run modal shell and form.
func Modal(ModalProps) g.Node {
	return modal.Shell(modal.Props{
		ID:         "run-modal",
		Title:      "Start Run",
		CloseLabel: "Close start run form",
		Body:       runForm(),
	})
}

func runForm() g.Node {
	return g.Group([]g.Node{
		html.Div(
			html.ID("run-form-error"),
			html.Class("alert alert--error"),
			g.Attr("role", "alert"),
		),
		html.Form(
			html.ID("run-form"),
			inputField("run-workflow-name", "Workflow Name",
				html.Input(
					html.ID("run-workflow-name"),
					html.Name("run-workflow-name"),
					html.Type("text"),
					html.ReadOnly(),
					html.Aria("readonly", "true"),
				),
			),
			inputField("run-workflow-id", "Workflow ID",
				html.Input(
					html.ID("run-workflow-id"),
					html.Name("run-workflow-id"),
					html.Type("text"),
					html.ReadOnly(),
					html.Aria("readonly", "true"),
				),
			),
			formField(
				html.Label(
					html.For("run-name"),
					g.Text("Display Name"),
				),
				html.Input(
					html.ID("run-name"),
					html.Name("run-name"),
					html.Type("text"),
					html.Placeholder("My Example Run"),
					html.Required(),
				),
				html.P(
					html.Class("form__hint"),
					g.Text("Shown in the dashboard. Choose something descriptive."),
				),
				html.P(
					html.Class("form__hint"),
					g.Text("Run ID (used for CLI commands): "),
					html.Code(
						html.ID("run-id-preview"),
						g.Text("—"),
					),
				),
			),
			html.Input(
				html.ID("run-id"),
				html.Name("run-id"),
				html.Type("hidden"),
			),
			html.Div(
				html.Class("form__actions"),
				button.Button(button.Props{
					Label:    "Cancel",
					Class:    "button--ghost",
					HideIcon: true,
					Data: map[string]string{
						"close-modal": "",
					},
				}),
				button.Button(button.Props{
					ID:       "run-submit",
					Label:    "Start Run",
					Class:    "button--accent",
					HideIcon: true,
					Type:     "submit",
				}),
			),
		),
	})
}

func inputField(id, label string, control g.Node) g.Node {
	return formField(
		html.Label(
			html.For(id),
			g.Text(label),
		),
		control,
	)
}

func formField(children ...g.Node) g.Node {
	return html.Div(
		html.Class("form__field"),
		g.Group(children),
	)
}
