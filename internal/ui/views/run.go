package views

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/pkg/ui/components"
)

// RunStep represents a single workflow step state within a run.
type RunStep struct {
	Name        string
	Status      string
	StatusClass string
}

// RunView summarizes a workflow run and its current state.
type RunView struct {
	DisplayName  string
	ID           string
	StateLabel   string
	StateClass   string
	WorkflowName string
	Steps        []RunStep
}

// RunColumnProps describes the runs column rendered on the dashboard.
type RunColumnProps struct {
	Title   string
	Actions []components.ButtonProps
	Runs    []RunView
}

// RunColumn renders the runs column including an empty state.
func RunColumn(p RunColumnProps) g.Node {
	return components.ColumnSection(components.ColumnProps{
		Title:        p.Title,
		Actions:      p.Actions,
		EmptyMessage: "No runs found.",
		Items:        runItems(p.Runs),
	})
}

func runItems(runs []RunView) []components.ColumnItem {
	if len(runs) == 0 {
		return nil
	}
	items := make([]components.ColumnItem, 0, len(runs))
	for _, current := range runs {
		run := current
		items = append(items, components.ColumnItem{
			DisableWrapper: true,
			Nodes: []g.Node{
				components.Card(components.CardProps{
					Title:        run.DisplayName,
					SummaryItems: runSummaryItems(run),
					Body:         runBody(run),
				}),
			},
		})
	}
	return items
}

func runSummaryItems(run RunView) []g.Node {
	items := []g.Node{
		components.StatusBadge(components.StatusBadgeProps{
			Label:   run.StateLabel,
			Variant: run.StateClass,
		}),
	}

	if id := strings.TrimSpace(run.ID); id != "" {
		items = append(items, components.Button(components.ButtonProps{
			Label:     "Tick",
			Class:     "button--primary button--sm run-tick-button",
			HideIcon:  true,
			AriaLabel: fmt.Sprintf("Run tick for %s", run.DisplayName),
			Data: map[string]string{
				"run-id":      id,
				"run-display": run.DisplayName,
			},
		}))
	}
	return items
}

func runBody(run RunView) g.Node {
	rows := []g.Node{
		components.ColumnInfoRow("Run ID:", run.ID),
		components.ColumnInfoRow("Workflow:", run.WorkflowName),
	}
	if steps := runSteps(run.Steps); steps != nil {
		rows = append(rows, html.H3(g.Text("Steps")), steps)
	}
	return g.Group(rows)
}

func runSteps(steps []RunStep) g.Node {
	if len(steps) == 0 {
		return nil
	}
	items := make([]components.DataListItem, 0, len(steps))
	for _, step := range steps {
		items = append(items, components.DataListItem{
			Primary: step.Name,
			Secondary: components.StatusBadge(components.StatusBadgeProps{
				Label:   step.Status,
				Variant: step.StatusClass,
			}),
		})
	}
	return components.DataList(components.DataListProps{Items: items})
}

// RunModalProps contains the data required to render the run modal.
type RunModalProps struct{}

// RunModal renders the run modal shell and form.
func RunModal(RunModalProps) g.Node {
	return components.ModalShell(components.ModalProps{
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
			runInputField("run-workflow-name", "Workflow Name",
				html.Input(
					html.ID("run-workflow-name"),
					html.Name("run-workflow-name"),
					html.Type("text"),
					html.ReadOnly(),
					html.Aria("readonly", "true"),
				),
			),
			runInputField("run-workflow-id", "Workflow ID",
				html.Input(
					html.ID("run-workflow-id"),
					html.Name("run-workflow-id"),
					html.Type("text"),
					html.ReadOnly(),
					html.Aria("readonly", "true"),
				),
			),
			runFormField(
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
				components.Button(components.ButtonProps{
					Label:    "Cancel",
					Class:    "button--ghost",
					HideIcon: true,
					Data: map[string]string{
						"close-modal": "",
					},
				}),
				components.Button(components.ButtonProps{
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

func runInputField(id, label string, control g.Node) g.Node {
	return runFormField(
		html.Label(
			html.For(id),
			g.Text(label),
		),
		control,
	)
}

func runFormField(children ...g.Node) g.Node {
	return html.Div(
		html.Class("form__field"),
		g.Group(children),
	)
}
