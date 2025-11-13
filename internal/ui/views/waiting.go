package views

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/pkg/ui/components"
)

// WaitingTask represents a single waiting task awaiting user interaction.
type WaitingTask struct {
	Name        string
	Description string
	Prompt      string
}

// WaitingGroup aggregates pending human tasks for a specific run.
type WaitingGroup struct {
	RunID          string
	RunDisplayName string
	WorkflowName   string
	TaskCount      int
	Tasks          []WaitingTask
}

// WaitingColumnProps represents the waiting-task column on the dashboard.
type WaitingColumnProps struct {
	Title   string
	Actions []components.ButtonProps
	Groups  []WaitingGroup
}

// WaitingColumn renders the waiting tasks column, including empty states.
func WaitingColumn(p WaitingColumnProps) g.Node {
	return components.ColumnSection(components.ColumnProps{
		Title:        p.Title,
		Actions:      p.Actions,
		Variant:      components.ColumnVariantMuted,
		ListClass:    "waiting-list",
		EmptyMessage: "No waiting tasks.",
		Items:        waitingItems(p.Groups),
	})
}

func waitingItems(groups []WaitingGroup) []components.ColumnItem {
	if len(groups) == 0 {
		return nil
	}

	items := make([]components.ColumnItem, 0, len(groups))
	for _, current := range groups {
		group := current
		items = append(items, components.ColumnItem{
			DisableWrapper: true,
			Node: html.Li(
				html.Div(
					html.Class("waiting-group__header"),
					html.Span(g.Text(group.RunDisplayName)),
					html.Span(
						html.Class("waiting-group__divider"),
						html.Aria("hidden", "true"),
					),
				),
				html.Ul(
					html.Class("waiting-group__tasks"),
					g.Group(waitingTasks(group.Tasks)),
				),
			),
		})
	}

	return items
}

func waitingTasks(tasks []WaitingTask) []g.Node {
	nodes := make([]g.Node, 0, len(tasks))
	for _, task := range tasks {
		nodes = append(nodes, html.Li(
			html.Class("card card--compact waiting-task"),
			html.Div(
				html.Class("waiting-task__name"),
				g.Text(task.Name),
			),
			g.If(task.Description != "", html.Div(
				html.Class("waiting-task__description"),
				g.Text(task.Description),
			)),
		))
	}
	return nodes
}
