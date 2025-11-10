package column

import (
	_ "embed"
	"html/template"

	cardcollapsible "composer/internal/ui/components/ui/card/collapsible"
	"composer/internal/ui/components/ui/columnheader"
	"composer/internal/ui/components/ui/datalist"
	"composer/internal/ui/components/ui/panel"
	"composer/internal/ui/components/ui/statusbadge"
	"composer/internal/ui/templates"
)

//go:embed column.tmpl
var columnTemplate string

var tmpl = templates.New(
	"run_column",
	"components/run/column/column.tmpl",
	columnTemplate,
	nil,
)

//go:embed body.tmpl
var bodyTemplate string

var bodyTmpl = templates.New(
	"run_column_body",
	"components/run/column/body.tmpl",
	bodyTemplate,
	nil,
)

//go:embed tick_button.tmpl
var tickButtonTemplate string

var tickButtonTmpl = templates.New(
	"run_column_tick_button",
	"components/run/column/tick_button.tmpl",
	tickButtonTemplate,
	nil,
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

// RenderSection renders the panel section that wraps the run list.
func (p Props) RenderSection() template.HTML {
	section := panel.SectionProps{
		Header: p.Header,
		Body: panel.MustRenderList(panel.ListProps{
			Items:        p.listItems(),
			EmptyMessage: "No runs found.",
		}),
	}
	return panel.MustRenderSection(section)
}

// Render executes the template for the run column component.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts failures into HTML comments.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}

func (p Props) listItems() []panel.ListItemProps {
	if len(p.Runs) == 0 {
		return nil
	}
	items := make([]panel.ListItemProps, 0, len(p.Runs))
	for _, run := range p.Runs {
		items = append(items, panel.ListItemProps{
			DisableWrapper: true,
			Content: cardcollapsible.MustRender(cardcollapsible.Props{
				Title:        run.DisplayName,
				SummaryItems: renderRunSummaryItems(run),
				Body:         renderRunBody(run),
			}),
		})
	}
	return items
}

func renderRunSummaryItems(run Run) []template.HTML {
	items := []template.HTML{
		statusbadge.MustRender(statusbadge.Props{
			Label:   run.StateLabel,
			Variant: run.StateClass,
		}),
	}
	if run.ID != "" {
		items = append(items, templates.SafeHTML(tickButtonTmpl.Render(run)))
	}
	return items
}

func renderRunBody(run Run) template.HTML {
	var steps template.HTML
	if len(run.Steps) > 0 {
		items := make([]datalist.Item, 0, len(run.Steps))
		for _, step := range run.Steps {
			items = append(items, datalist.Item{
				Primary: step.Name,
				Secondary: statusbadge.MustRender(statusbadge.Props{
					Label:   step.Status,
					Variant: step.StatusClass,
				}),
			})
		}
		steps = datalist.MustRender(datalist.Props{Items: items})
	}

	props := runBodyProps{
		Run:   run,
		Steps: steps,
	}
	return templates.SafeHTML(bodyTmpl.Render(props))
}

type runBodyProps struct {
	Run
	Steps template.HTML
}

func (p runBodyProps) HasSteps() bool {
	return len(p.Steps) > 0
}

func (p runBodyProps) RenderSteps() template.HTML {
	return p.Steps
}
