package column

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/button"
	appcolumn "composer/internal/ui/components/column"
)

// Task represents a single waiting task awaiting user interaction.
type Task struct {
	Name        string
	Description string
	Prompt      string
}

// Group aggregates pending human tasks for a specific run.
type Group struct {
	RunID          string
	RunDisplayName string
	WorkflowName   string
	TaskCount      int
	Tasks          []Task
}

// Props represents the waiting-task column on the dashboard.
type Props struct {
	Title   string
	Actions []button.Props
	Groups  []Group
}

// Column renders the waiting tasks column, including empty states.
func Column(p Props) g.Node {
	return appcolumn.Section(appcolumn.Props{
		Title:        p.Title,
		Actions:      p.Actions,
		Variant:      appcolumn.VariantMuted,
		ListClass:    "waiting-list",
		EmptyMessage: "No waiting tasks.",
		Items:        waitingItems(p.Groups),
	})
}

func waitingItems(groups []Group) []appcolumn.Item {
	if len(groups) == 0 {
		return nil
	}

	items := make([]appcolumn.Item, 0, len(groups))
	for _, group := range groups {
		group := group
		items = append(items, appcolumn.Item{
			DisableWrapper: true,
			Nodes: []g.Node{html.Li(
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
			)},
		})
	}

	return items
}

func waitingTasks(tasks []Task) []g.Node {
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
