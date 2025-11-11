package workflow

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/button"
	"composer/internal/ui/components/card"
	appcolumn "composer/internal/ui/components/column"
	"composer/internal/ui/components/datalist"
)

// Workflow represents a single workflow card in the column.
type Workflow struct {
	DisplayName string
	ID          string
	Description string
	Message     string
	StepNames   []string
}

// ColumnProps describes the workflow column on the dashboard.
type ColumnProps struct {
	Title     string
	Actions   []button.Props
	Workflows []Workflow
}

// Column renders the workflow column, including empty states.
func Column(p ColumnProps) g.Node {
	return appcolumn.Section(appcolumn.Props{
		Title:        p.Title,
		Actions:      p.Actions,
		EmptyMessage: "No workflows available.",
		Items:        workflowItems(p.Workflows),
	})
}

func workflowItems(workflows []Workflow) []appcolumn.Item {
	if len(workflows) == 0 {
		return nil
	}
	items := make([]appcolumn.Item, 0, len(workflows))
	for _, wf := range workflows {
		workflow := wf
		items = append(items, appcolumn.Item{
			DisableWrapper: true,
			Nodes: []g.Node{
				card.Card(card.Props{
					Title:        workflow.DisplayName,
					SummaryItems: workflowSummaryItems(workflow),
					Body:         workflowBody(workflow),
				}),
			},
		})
	}
	return items
}

func workflowSummaryItems(w Workflow) []g.Node {
	id := strings.TrimSpace(w.ID)
	if id == "" {
		return nil
	}

	label := fmt.Sprintf("Start run for workflow %s", strings.TrimSpace(w.DisplayName))
	return []g.Node{
		button.Button(button.Props{
			Label:     "Run",
			Class:     "button--accent button--sm",
			HideIcon:  true,
			AriaLabel: label,
			Data: map[string]string{
				"open-run-modal":   "",
				"workflow-id":      id,
				"workflow-display": w.DisplayName,
			},
		}),
	}
}

func workflowBody(w Workflow) g.Node {
	rows := []g.Node{
		appcolumn.InfoRow("ID:", w.ID),
		appcolumn.InfoRow("Display Name:", w.DisplayName),
	}
	if strings.TrimSpace(w.Description) != "" {
		rows = append(rows, appcolumn.InfoRow("Description:", w.Description))
	}
	if strings.TrimSpace(w.Message) != "" {
		rows = append(rows, appcolumn.InfoRow("Message:", w.Message))
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
	items := make([]datalist.Item, 0, len(stepNames))
	for _, name := range stepNames {
		items = append(items, datalist.Item{
			Primary: name,
		})
	}
	return datalist.List(datalist.Props{Items: items})
}
