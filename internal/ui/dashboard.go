package ui

import (
	"html/template"
	"io"
	"sort"
	"strings"

	"composer/internal/workflow"
)

var dashboardTemplate = template.Must(template.New("dashboard").Parse(
	strings.Join([]string{
		workflowColumnTemplate,
		runColumnTemplate,
		workflowModalTemplate,
		dashboardPageTemplate,
	}, "\n"),
))

type dashboardData struct {
	Workflows []workflowItem
	Runs      []runItem
}

type workflowItem struct {
	DisplayName string
	Workflow    workflow.Workflow
	StepNames   []string
}

type runItem struct {
	Name         string
	StateLabel   string
	StateClass   string
	WorkflowName string
	Steps        []runStep
}

type runStatus struct {
	Label string
	Class string
}

type runStep struct {
	Name        string
	Status      string
	StatusClass string
}

func renderDashboard(w io.Writer, data dashboardData) error {
	return dashboardTemplate.ExecuteTemplate(w, "dashboard", data)
}

func buildDashboardData(workflows []workflow.Workflow, runs []workflow.RunState) dashboardData {
	items := make([]workflowItem, 0, len(workflows))
	for _, wf := range workflows {
		stepNames := make([]string, 0, len(wf.Steps))
		for _, step := range wf.Steps {
			name := strings.TrimSpace(step.Name)
			if name != "" {
				stepNames = append(stepNames, name)
			}
		}

		items = append(items, workflowItem{
			DisplayName: workflowDisplayName(wf),
			Workflow:    wf,
			StepNames:   stepNames,
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].DisplayName < items[j].DisplayName })

	runItems := make([]runItem, 0, len(runs))
	for _, run := range runs {
		status := summarizeRunState(run)

		stepNames := make([]string, 0, len(run.StepStates))
		for name := range run.StepStates {
			stepNames = append(stepNames, name)
		}
		sort.Strings(stepNames)

		steps := make([]runStep, 0, len(stepNames))
		for _, name := range stepNames {
			stepState := run.StepStates[name]
			steps = append(steps, runStep{
				Name:        name,
				Status:      string(stepState.Status),
				StatusClass: stateClassForStatus(stepState.Status),
			})
		}

		runItems = append(runItems, runItem{
			Name:         run.RunName,
			StateLabel:   status.Label,
			StateClass:   status.Class,
			WorkflowName: run.WorkflowName,
			Steps:        steps,
		})
	}
	sort.Slice(runItems, func(i, j int) bool { return runItems[i].Name < runItems[j].Name })

	return dashboardData{
		Workflows: items,
		Runs:      runItems,
	}
}

func workflowDisplayName(wf workflow.Workflow) string {
	if title := strings.TrimSpace(wf.Title); title != "" {
		return title
	}
	if strings.TrimSpace(wf.ID) != "" {
		return wf.ID
	}
	return "Untitled workflow"
}

func stateClassForStatus(status workflow.StepStatus) string {
	switch status {
	case workflow.StatusFailed:
		return "state-failed"
	case workflow.StatusSucceeded:
		return "state-succeeded"
	case workflow.StatusReady:
		return "state-ready"
	case workflow.StatusPending:
		return "state-pending"
	default:
		return "state-unknown"
	}
}

func summarizeRunState(rs workflow.RunState) runStatus {
	if len(rs.StepStates) == 0 {
		return runStatus{Label: "pending", Class: "state-pending"}
	}

	allSucceeded := len(rs.StepStates) > 0
	hasReady := false
	hasPending := false

	for _, step := range rs.StepStates {
		switch step.Status {
		case workflow.StatusFailed:
			return runStatus{Label: "failed", Class: "state-failed"}
		case workflow.StatusSucceeded:
			// keep allSucceeded true unless another status changes it
		case workflow.StatusReady:
			hasReady = true
			allSucceeded = false
		case workflow.StatusPending:
			hasPending = true
			allSucceeded = false
		default:
			allSucceeded = false
		}
	}

	if allSucceeded {
		return runStatus{Label: "succeeded", Class: "state-succeeded"}
	}
	if hasReady {
		return runStatus{Label: "ready", Class: "state-ready"}
	}
	if hasPending {
		return runStatus{Label: "pending", Class: "state-pending"}
	}
	return runStatus{Label: "unknown", Class: "state-unknown"}
}
