package column

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
	cardcollapsible "composer/pkg/ui/components/cardcollapsible"
	"composer/pkg/ui/components/columnheader"
	"composer/pkg/ui/components/datalist"
	"composer/pkg/ui/components/panel"
)

//go:embed column.tmpl
var columnTemplate string

var tmpl = templates.New(
	"workflow_column",
	"internal/ui/components/workflow_column/column.tmpl",
	columnTemplate,
	nil,
)

//go:embed body.tmpl
var bodyTemplate string

var bodyTmpl = templates.New(
	"workflow_column_body",
	"internal/ui/components/workflow_column/body.tmpl",
	bodyTemplate,
	nil,
)

//go:embed run_button.tmpl
var runButtonTemplate string

var runButtonTmpl = templates.New(
	"workflow_column_run_button",
	"internal/ui/components/workflow_column/run_button.tmpl",
	runButtonTemplate,
	nil,
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

// RenderSection renders the shared panel section with the workflow list.
func (p Props) RenderSection() template.HTML {
	section := panel.SectionProps{
		Header: p.Header,
		Body: panel.MustRenderList(panel.ListProps{
			Items:        p.listItems(),
			EmptyMessage: "No workflows available.",
		}),
	}
	return panel.MustRenderSection(section)
}

// Render executes the component template.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts failures into HTML comments.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}

func (p Props) listItems() []panel.ListItemProps {
	if len(p.Workflows) == 0 {
		return nil
	}
	items := make([]panel.ListItemProps, 0, len(p.Workflows))
	for _, wf := range p.Workflows {
		items = append(items, panel.ListItemProps{
			DisableWrapper: true,
			Content: cardcollapsible.MustRender(cardcollapsible.Props{
				Title:        wf.DisplayName,
				SummaryItems: renderWorkflowSummaryItems(wf),
				Body:         renderWorkflowBody(wf),
			}),
		})
	}
	return items
}

func renderWorkflowSummaryItems(w Workflow) []template.HTML {
	if w.ID == "" {
		return nil
	}
	return []template.HTML{
		templates.SafeHTML(runButtonTmpl.Render(w)),
	}
}

func renderWorkflowBody(w Workflow) template.HTML {
	var steps template.HTML
	if len(w.StepNames) > 0 {
		items := make([]datalist.Item, 0, len(w.StepNames))
		for _, name := range w.StepNames {
			items = append(items, datalist.Item{Primary: name})
		}
		steps = datalist.MustRender(datalist.Props{Items: items})
	}

	props := workflowBodyProps{
		Workflow: w,
		Steps:    steps,
	}
	return templates.SafeHTML(bodyTmpl.Render(props))
}

type workflowBodyProps struct {
	Workflow
	Steps template.HTML
}

func (p workflowBodyProps) HasSteps() bool {
	return len(p.Steps) > 0
}

func (p workflowBodyProps) RenderSteps() template.HTML {
	return p.Steps
}
