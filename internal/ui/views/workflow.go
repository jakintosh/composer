package views

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/pkg/ui/components"
)

// WorkflowView represents a single workflow card in the column.
type WorkflowView struct {
	DisplayName string
	ID          string
	Description string
	Message     string
	StepNames   []string
}

// WorkflowColumnProps describes the workflow column on the dashboard.
type WorkflowColumnProps struct {
	Title     string
	Actions   []components.ButtonProps
	Workflows []WorkflowView
}

// WorkflowColumn renders the workflow column, including empty states.
func WorkflowColumn(p WorkflowColumnProps) g.Node {
	return components.ColumnSection(components.ColumnProps{
		Title:        p.Title,
		Actions:      p.Actions,
		EmptyMessage: "No workflows available.",
		Items:        workflowItems(p.Workflows),
	})
}

func workflowItems(workflows []WorkflowView) []components.ColumnItem {
	if len(workflows) == 0 {
		return nil
	}
	items := make([]components.ColumnItem, 0, len(workflows))
	for _, wf := range workflows {
		workflow := wf
		items = append(items, components.ColumnItem{
			DisableWrapper: true,
			Node: components.Card(components.CardProps{
				Title:  workflow.DisplayName,
				Body:   workflowBody(workflow),
				Action: workflowActionButton(workflow),
			}),
		})
	}
	return items
}

func workflowActionButton(w WorkflowView) *components.ButtonProps {
	id := strings.TrimSpace(w.ID)
	if id == "" {
		return nil
	}

	label := fmt.Sprintf("Start run for workflow %s", strings.TrimSpace(w.DisplayName))
	return &components.ButtonProps{
		Label:     "Run",
		Class:     "button--accent button--sm",
		HideIcon:  true,
		AriaLabel: label,
		Data: map[string]string{
			"open-run-modal":   "",
			"workflow-id":      id,
			"workflow-display": w.DisplayName,
		},
	}
}

func workflowBody(w WorkflowView) g.Node {
	rows := []g.Node{
		components.ColumnInfoRow("ID:", w.ID),
		components.ColumnInfoRow("Display Name:", w.DisplayName),
	}
	if strings.TrimSpace(w.Description) != "" {
		rows = append(rows, components.ColumnInfoRow("Description:", w.Description))
	}
	if strings.TrimSpace(w.Message) != "" {
		rows = append(rows, components.ColumnInfoRow("Message:", w.Message))
	}
	if steps := workflowSteps(w.StepNames); steps != nil {
		rows = append(rows, html.H3(g.Text("Steps")), steps)
	}
	return g.Group(rows)
}

func workflowSteps(stepNames []string) g.Node {
	if len(stepNames) == 0 {
		return nil
	}
	items := make([]components.DataListItem, 0, len(stepNames))
	for _, name := range stepNames {
		items = append(items, components.DataListItem{
			Primary: name,
		})
	}
	return components.DataList(components.DataListProps{Items: items})
}

// WorkflowModalProps holds the controls displayed in the workflow creation modal.
type WorkflowModalProps struct {
	AddStepButton components.ButtonProps
}

// WorkflowModal renders the workflow modal shell, form, and step template.
func WorkflowModal(p WorkflowModalProps) g.Node {
	return g.Group([]g.Node{
		components.ModalShell(components.ModalProps{
			ID:         "workflow-modal",
			Title:      "Create Workflow",
			CloseLabel: "Close create workflow form",
			Body:       workflowForm(p),
		}),
		workflowStepTemplate(),
	})
}

func workflowForm(p WorkflowModalProps) g.Node {
	return g.Group([]g.Node{
		html.Div(
			html.ID("workflow-form-error"),
			html.Class("alert alert--error"),
			g.Attr("role", "alert"),
		),
		html.Form(
			html.ID("workflow-form"),
			workflowInputField("workflow-id", "Workflow ID",
				html.Input(
					html.ID("workflow-id"),
					html.Name("workflow-id"),
					html.Type("text"),
					html.Placeholder("my-workflow"),
					html.Required(),
				),
			),
			workflowInputField("workflow-title", "Title",
				html.Input(
					html.ID("workflow-title"),
					html.Name("workflow-title"),
					html.Type("text"),
					html.Placeholder("Human-friendly workflow title"),
					html.Required(),
				),
			),
			workflowInputField("workflow-description", "Description",
				html.Textarea(
					html.ID("workflow-description"),
					html.Name("workflow-description"),
					html.Placeholder("Explain what this workflow accomplishes"),
				),
			),
			workflowInputField("workflow-message", "Message",
				html.Textarea(
					html.ID("workflow-message"),
					html.Name("workflow-message"),
					html.Placeholder("Optional run message shown to operators"),
				),
			),
			workflowFormField(
				html.Div(
					html.Class("workflow-steps__header"),
					html.H3(g.Text("Steps")),
					components.Button(p.AddStepButton),
				),
				html.Div(
					html.ID("workflow-steps"),
				),
			),
			html.Div(
				html.Class("form__actions"),
				components.Button(components.ButtonProps{
					Label:    "Cancel",
					Class:    "button--ghost",
					HideIcon: true,
					Data: map[string]string{
						"close-modal": "",
					},
				}),
				components.Button(components.ButtonProps{
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
				components.Button(components.ButtonProps{
					Label:    "Remove",
					Class:    "button--text button--danger remove-step",
					HideIcon: true,
				}),
			),
			workflowFormField(
				html.Label(g.Text("Step Name")),
				html.Input(
					html.Type("text"),
					html.Name("step-name"),
					html.Placeholder("identify-step"),
					html.Required(),
				),
			),
			workflowFormField(
				html.Label(g.Text("Description")),
				html.Textarea(
					html.Name("step-description"),
					html.Placeholder("Optional description"),
				),
			),
			workflowFormField(
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
			workflowFormField(
				html.Label(g.Text("Prompt")),
				html.Textarea(
					html.Name("step-prompt"),
					html.Placeholder("Guidance for human or cognitive handlers"),
				),
			),
			workflowFormField(
				html.Label(g.Text("Content")),
				html.Textarea(
					html.Name("step-content"),
					html.Placeholder("Optional inline content"),
				),
			),
			workflowFormField(
				html.Label(g.Text("Inputs (one per line or comma separated)")),
				html.Textarea(
					html.Name("step-inputs"),
					html.Placeholder("input-a\ninput-b"),
				),
			),
			workflowFormField(
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

func workflowInputField(id, label string, control g.Node) g.Node {
	return workflowFormField(
		html.Label(
			html.For(id),
			g.Text(label),
		),
		control,
	)
}

func workflowFormField(children ...g.Node) g.Node {
	return html.Div(
		html.Class("form__field"),
		g.Group(children),
	)
}
