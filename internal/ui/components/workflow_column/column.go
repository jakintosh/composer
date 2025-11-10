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
)

// Workflow represents a single workflow card in the column.
type Workflow struct {
	DisplayName string
	ID          string
	Description string
	Message     string
	StepNames   []string
}

// Props describes the workflow column on the dashboard.
type Props struct {
	Header    columnheader.Props
	Workflows []Workflow
}

// Column renders the workflow column, including empty states.
func Column(p Props) g.Node {
	return panel.Section(panel.SectionProps{
		Header: p.Header,
		Body: panel.List(panel.ListProps{
			Items:        workflowItems(p.Workflows),
			EmptyMessage: "No workflows available.",
		}),
	})
}

func workflowItems(workflows []Workflow) []panel.ListItemProps {
	if len(workflows) == 0 {
		return nil
	}
	items := make([]panel.ListItemProps, 0, len(workflows))
	for _, wf := range workflows {
		workflow := wf
		items = append(items, panel.ListItemProps{
			DisableWrapper: true,
			Content: card.Card(card.Props{
				Title:        workflow.DisplayName,
				SummaryItems: workflowSummaryItems(workflow),
				Body:         workflowBody(workflow),
			}),
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
		infoRow("ID:", w.ID),
		infoRow("Display Name:", w.DisplayName),
	}
	if strings.TrimSpace(w.Description) != "" {
		rows = append(rows, infoRow("Description:", w.Description))
	}
	if strings.TrimSpace(w.Message) != "" {
		rows = append(rows, infoRow("Message:", w.Message))
	}
	if steps := workflowSteps(w.StepNames); steps != nil {
		rows = append(rows, html.H3(g.Text("Steps")), steps)
	}
	return g.Group(rows)
}

func infoRow(label, value string) g.Node {
	return html.P(
		html.Strong(g.Text(label+" ")),
		g.Text(value),
	)
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
