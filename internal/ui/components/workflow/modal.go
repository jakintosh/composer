package workflow

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/button"
	"composer/internal/ui/components/modal"
)

// ModalProps holds the controls displayed in the workflow creation modal.
type ModalProps struct {
	AddStepButton button.Props
}

// Modal renders the workflow modal shell, form, and step template.
func Modal(p ModalProps) g.Node {
	return g.Group([]g.Node{
		modal.Shell(modal.Props{
			ID:         "workflow-modal",
			Title:      "Create Workflow",
			CloseLabel: "Close create workflow form",
			Body:       workflowForm(p),
		}),
		workflowStepTemplate(),
	})
}

func workflowForm(p ModalProps) g.Node {
	return g.Group([]g.Node{
		html.Div(
			html.ID("workflow-form-error"),
			html.Class("alert alert--error"),
			g.Attr("role", "alert"),
		),
		html.Form(
			html.ID("workflow-form"),
			inputField("workflow-id", "Workflow ID",
				html.Input(
					html.ID("workflow-id"),
					html.Name("workflow-id"),
					html.Type("text"),
					html.Placeholder("my-workflow"),
					html.Required(),
				),
			),
			inputField("workflow-title", "Title",
				html.Input(
					html.ID("workflow-title"),
					html.Name("workflow-title"),
					html.Type("text"),
					html.Placeholder("Human-friendly workflow title"),
					html.Required(),
				),
			),
			inputField("workflow-description", "Description",
				html.Textarea(
					html.ID("workflow-description"),
					html.Name("workflow-description"),
					html.Placeholder("Explain what this workflow accomplishes"),
				),
			),
			inputField("workflow-message", "Message",
				html.Textarea(
					html.ID("workflow-message"),
					html.Name("workflow-message"),
					html.Placeholder("Optional run message shown to operators"),
				),
			),
			formField(
				html.Div(
					html.Class("workflow-steps__header"),
					html.H3(g.Text("Steps")),
					button.Button(p.AddStepButton),
				),
				html.Div(
					html.ID("workflow-steps"),
				),
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
					ID:       "workflow-submit",
					Label:    "Save Workflow",
					Class:    "button--accent",
					HideIcon: true,
					Type:     "submit",
				}),
			),
		),
	})
}

func workflowStepTemplate() g.Node {
	return html.Template(
		html.ID("workflow-step-template"),
		html.Article(
			html.Class("card card--form workflow-step"),
			html.Header(
				html.Class("workflow-step__header"),
				html.H4(
					g.Text("Step "),
					html.Span(
						html.Class("step-number"),
					),
				),
				button.Button(button.Props{
					Label:    "Remove",
					Class:    "button--text button--danger remove-step",
					HideIcon: true,
				}),
			),
			formField(
				html.Label(g.Text("Step Name")),
				html.Input(
					html.Type("text"),
					html.Name("step-name"),
					html.Placeholder("identify-step"),
					html.Required(),
				),
			),
			formField(
				html.Label(g.Text("Description")),
				html.Textarea(
					html.Name("step-description"),
					html.Placeholder("Optional description"),
				),
			),
			formField(
				html.Label(g.Text("Handler")),
				html.Select(
					html.Name("step-handler"),
					html.Option(
						html.Value("tool"),
						html.Selected(),
						g.Text("tool"),
					),
					html.Option(
						html.Value("human"),
						g.Text("human"),
					),
				),
			),
			formField(
				html.Label(g.Text("Prompt")),
				html.Textarea(
					html.Name("step-prompt"),
					html.Placeholder("Guidance for human or cognitive handlers"),
				),
			),
			formField(
				html.Label(g.Text("Content")),
				html.Textarea(
					html.Name("step-content"),
					html.Placeholder("Optional inline content"),
				),
			),
			formField(
				html.Label(g.Text("Inputs (one per line or comma separated)")),
				html.Textarea(
					html.Name("step-inputs"),
					html.Placeholder("input-a\ninput-b"),
				),
			),
			formField(
				html.Label(g.Text("Output")),
				html.Input(
					html.Type("text"),
					html.Name("step-output"),
					html.Placeholder("result-key"),
				),
			),
		),
	)
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
