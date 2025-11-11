package run

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/button"
	"composer/internal/ui/components/card"
	appcolumn "composer/internal/ui/components/column"
	"composer/internal/ui/components/datalist"
	"composer/internal/ui/components/statusbadge"
)

// Step represents a single workflow step state within a run.
type Step struct {
	Name        string
	Status      string
	StatusClass string
}

// Run summarizes a workflow run and its current state.
type Run struct {
	DisplayName  string
	ID           string
	StateLabel   string
	StateClass   string
	WorkflowName string
	Steps        []Step
}

// ColumnProps describes the runs column rendered on the dashboard.
type ColumnProps struct {
	Title   string
	Actions []button.Props
	Runs    []Run
}

// Column renders the runs column including an empty state.
func Column(p ColumnProps) g.Node {
	return appcolumn.Section(appcolumn.Props{
		Title:        p.Title,
		Actions:      p.Actions,
		EmptyMessage: "No runs found.",
		Items:        runItems(p.Runs),
	})
}

func runItems(runs []Run) []appcolumn.Item {
	if len(runs) == 0 {
		return nil
	}
	items := make([]appcolumn.Item, 0, len(runs))
	for _, run := range runs {
		current := run
		items = append(items, appcolumn.Item{
			DisableWrapper: true,
			Nodes: []g.Node{
				card.Card(card.Props{
					Title:        current.DisplayName,
					SummaryItems: runSummaryItems(current),
					Body:         runBody(current),
				}),
			},
		})
	}
	return items
}

func runSummaryItems(run Run) []g.Node {
	items := []g.Node{
		statusbadge.Badge(statusbadge.Props{
			Label:   run.StateLabel,
			Variant: run.StateClass,
		}),
	}

	if id := strings.TrimSpace(run.ID); id != "" {
		items = append(items, button.Button(button.Props{
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

func runBody(run Run) g.Node {
	rows := []g.Node{
		appcolumn.InfoRow("Run ID:", run.ID),
		appcolumn.InfoRow("Workflow:", run.WorkflowName),
	}
	if steps := runSteps(run.Steps); steps != nil {
		rows = append(rows, html.H3(g.Text("Steps")), steps)
	}
	return g.Group(rows)
}

func runSteps(steps []Step) g.Node {
	if len(steps) == 0 {
		return nil
	}
	items := make([]datalist.Item, 0, len(steps))
	for _, step := range steps {
		items = append(items, datalist.Item{
			Primary: step.Name,
			Secondary: statusbadge.Badge(statusbadge.Props{
				Label:   step.Status,
				Variant: step.StatusClass,
			}),
		})
	}
	return datalist.List(datalist.Props{Items: items})
}
