package column

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/button"
	"composer/internal/ui/components/card"
	"composer/internal/ui/components/columnheader"
	"composer/internal/ui/components/datalist"
	"composer/internal/ui/components/panel"
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

// Props describes the runs column rendered on the dashboard.
type Props struct {
	Header columnheader.Props
	Runs   []Run
}

// Column renders the runs column including an empty state.
func Column(p Props) g.Node {
	return panel.Section(panel.SectionProps{
		Header: p.Header,
		Body: panel.List(panel.ListProps{
			Items:        runItems(p.Runs),
			EmptyMessage: "No runs found.",
		}),
	})
}

func runItems(runs []Run) []panel.ListItemProps {
	if len(runs) == 0 {
		return nil
	}
	items := make([]panel.ListItemProps, 0, len(runs))
	for _, run := range runs {
		current := run
		items = append(items, panel.ListItemProps{
			DisableWrapper: true,
			Content: card.Card(card.Props{
				Title:        current.DisplayName,
				SummaryItems: runSummaryItems(current),
				Body:         runBody(current),
			}),
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
		infoRow("Run ID:", run.ID),
		infoRow("Workflow:", run.WorkflowName),
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

func infoRow(label, value string) g.Node {
	return html.P(
		html.Strong(g.Text(label+" ")),
		g.Text(value),
	)
}
