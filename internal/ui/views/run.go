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
func (p RunColumnProps) Render() g.Node {

	// render the runs into column items
	numRuns := len(p.Runs)
	var items []components.ColumnItem
	if 0 < numRuns {
		items = make([]components.ColumnItem, numRuns)
		for i, run := range p.Runs {
			node := run.Render()
			item := components.ColumnItem{
				DisableWrapper: true,
				Node:           node,
			}
			items[i] = item
		}
	}

	// build and render the column
	props := components.ColumnProps{
		Title:        p.Title,
		Actions:      p.Actions,
		EmptyMessage: "No runs found.",
		Items:        items,
	}
	return components.ColumnSection(props)
}

// Render returns the run view card node.
func (run RunView) Render() g.Node {

	var status *components.StatusBadgeProps
	if label := strings.TrimSpace(run.StateLabel); label != "" {
		status = &components.StatusBadgeProps{
			Label:   run.StateLabel,
			Variant: run.StateClass,
		}
	}

	var action *components.ButtonProps
	if id := strings.TrimSpace(run.ID); id != "" {
		action = &components.ButtonProps{
			Label:     "Tick",
			Class:     "button--primary button--sm run-tick-button",
			HideIcon:  true,
			AriaLabel: fmt.Sprintf("Run tick for %s", run.DisplayName),
			Data: map[string]string{
				"run-id":      id,
				"run-display": run.DisplayName,
			},
		}
	}

	// render the body
	bodyNodes := []g.Node{
		components.ColumnInfoRow("Run ID:", run.ID),
		components.ColumnInfoRow("Workflow:", run.WorkflowName),
	}

	// render steps in the body
	numSteps := len(run.Steps)
	if numSteps > 0 {
		bodyNodes = append(bodyNodes, html.H3(g.Text("Steps")))

		// transform runstep into datalistitem
		items := make([]components.DataListItem, numSteps)
		for i, step := range run.Steps {

			// build and render status badge
			props := components.StatusBadgeProps{
				Label:   step.Status,
				Variant: step.StatusClass,
			}
			badge := components.StatusBadge(props)

			// build datalistitem
			items[i] = components.DataListItem{
				Primary:   step.Name,
				Secondary: badge,
			}
		}

		// build and render data list
		props := components.DataListProps{
			Items: items,
		}
		list := components.DataList(props)
		bodyNodes = append(bodyNodes, list)
	}

	// build and render card
	card := components.CardProps{
		Title:  run.DisplayName,
		Status: status,
		Action: action,
		Body:   g.Group(bodyNodes),
	}
	return components.Card(card)
}

// RunModalProps contains the data required to render the run modal.
type RunModalProps struct{}

// RunModal renders the run modal shell and form.
func (RunModalProps) Render() g.Node {
	workflowNameField := FormFieldProps{
		Label: "Workflow Name",
		FormInputProps: FormInputProps{
			ID:       "run-workflow-name",
			Name:     "run-workflow-name",
			Type:     "text",
			ReadOnly: true,
		},
	}

	workflowIdField := FormFieldProps{
		Label: "Workflow ID",
		FormInputProps: FormInputProps{
			ID:       "run-workflow-id",
			Name:     "run-workflow-id",
			Type:     "text",
			ReadOnly: true,
		},
	}

	runNameField := FormFieldProps{
		Label: "Display Name",
		FormInputProps: FormInputProps{
			ID:          "run-name",
			Name:        "run-name",
			Type:        "text",
			Placeholder: "My Example Run",
			Required:    true,
		},
		Hint: "Shown in the dashboard. Choose something descriptive.",
	}

	runIdPreviewField := FormFieldProps{
		Label: "Run ID",
		FormInputProps: FormInputProps{
			ID:       "run-id-preview",
			Name:     "run-id-preview",
			Type:     "text",
			ReadOnly: true,
		},
		Hint: "Used as internal unique id.",
	}

	runIdHiddenField := FormInputProps{
		ID:   "run-id",
		Name: "run-id",
		Type: "hidden",
	}

	cancelAction := components.ButtonProps{
		Label:    "Cancel",
		Class:    "button--ghost",
		HideIcon: true,
		Data: map[string]string{
			"close-modal": "",
		},
	}

	submitAction := components.ButtonProps{
		ID:       "run-submit",
		Label:    "Start Run",
		Class:    "button--accent",
		HideIcon: true,
		Type:     "submit",
	}

	form := FormProps{
		ID:      "run-form",
		ErrorID: "run-form-error",
		Fields: []FormFieldProps{
			workflowNameField,
			workflowIdField,
			runNameField,
			runIdPreviewField,
		},
		Hidden: []FormInputProps{
			runIdHiddenField,
		},
		Actions: []components.ButtonProps{
			cancelAction,
			submitAction,
		},
	}

	props := components.ModalProps{
		ID:         "run-modal",
		Title:      "Start Run",
		CloseLabel: "Close start run form",
		Body:       form.Render(),
	}
	return components.ModalShell(props)
}

// FormProps represents the run creation form data.
type FormProps struct {
	ID      string
	ErrorID string
	Fields  []FormFieldProps
	Hidden  []FormInputProps
	Actions []components.ButtonProps
}

// FormFieldProps describes a visible form field (label + control + optional hints).
type FormFieldProps struct {
	Label string
	Hint  string
	FormInputProps
}

// FormInputProps provides the data necessary to render a single input element.
type FormInputProps struct {
	ID          string
	Name        string
	Type        string
	Placeholder string
	ReadOnly    bool
	Required    bool
}

func (form FormProps) Render() g.Node {
	return g.Group{
		// error
		html.Div(
			html.ID(form.ErrorID),
			html.Class("alert alert--error"),
			g.Attr("role", "alert"),
		),
		// form
		html.Form(
			html.ID(form.ID),
			g.Map(form.Fields, func(f FormFieldProps) g.Node { return f.Render() }),
			g.Map(form.Hidden, func(i FormInputProps) g.Node { return i.Render() }),
			html.Div(
				html.Class("form__actions"),
				g.Map(form.Actions, func(props components.ButtonProps) g.Node {
					return components.Button(props)
				}),
			),
		),
	}
}

func (field FormFieldProps) Render() g.Node {
	return html.Div(
		html.Class("form__field"),
		html.Label(
			html.For(field.ID),
			g.Text(field.Label),
		),
		field.FormInputProps.Render(),
		g.If(field.Hint != "", html.P(
			html.Class("form__hint"),
			g.Text(field.Hint),
		)),
	)
}

func (input FormInputProps) Render() g.Node {

	inputType := strings.TrimSpace(input.Type)
	if inputType == "" {
		inputType = "text"
	}

	return html.Input(
		g.If(input.ID != "", html.ID(input.ID)),
		g.If(input.Name != "", html.Name(input.Name)),
		html.Type(inputType),
		g.If(input.Placeholder != "", html.Placeholder(input.Placeholder)),
		g.If(input.ReadOnly, html.ReadOnly()),
		g.If(input.ReadOnly, html.Aria("readonly", "true")),
		g.If(input.Required, html.Required()),
	)
}
