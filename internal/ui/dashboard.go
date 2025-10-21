package ui

import (
	"sort"
	"strings"

	"composer/internal/workflow"
)

type dashboardViewModel struct {
	Workflows []workflowViewModel
	Runs      []runViewModel
}

type workflowViewModel struct {
	DisplayName string
	ID          string
	Title       string
	Description string
	Message     string
	StepNames   []string
}

type runViewModel struct {
	Name         string
	StateLabel   string
	StateClass   string
	WorkflowName string
	Steps        []runStepViewModel
}

type runStepViewModel struct {
	Name        string
	Status      string
	StatusClass string
}

type runStatus struct {
	Label string
	Class string
}

func buildDashboardViewModel(workflows []workflow.Workflow, runs []workflow.RunState) dashboardViewModel {
	workflowVMs := make([]workflowViewModel, 0, len(workflows))
	for _, wf := range workflows {
		workflowVMs = append(workflowVMs, workflowViewModel{
			DisplayName: workflowDisplayName(wf),
			ID:          strings.TrimSpace(wf.ID),
			Title:       strings.TrimSpace(wf.Title),
			Description: strings.TrimSpace(wf.Description),
			Message:     strings.TrimSpace(wf.Message),
			StepNames:   collectStepNames(wf.Steps),
		})
	}
	sort.Slice(workflowVMs, func(i, j int) bool { return workflowVMs[i].DisplayName < workflowVMs[j].DisplayName })

	runVMs := make([]runViewModel, 0, len(runs))
	for _, run := range runs {
		status := summarizeRunState(run)
		stepNames := sortedStepNames(run.StepStates)

		steps := make([]runStepViewModel, 0, len(stepNames))
		for _, name := range stepNames {
			stepState := run.StepStates[name]
			steps = append(steps, runStepViewModel{
				Name:        name,
				Status:      string(stepState.Status),
				StatusClass: stateClassForStatus(stepState.Status),
			})
		}

		runVMs = append(runVMs, runViewModel{
			Name:         strings.TrimSpace(run.RunName),
			StateLabel:   status.Label,
			StateClass:   status.Class,
			WorkflowName: strings.TrimSpace(run.WorkflowName),
			Steps:        steps,
		})
	}
	sort.Slice(runVMs, func(i, j int) bool { return runVMs[i].Name < runVMs[j].Name })

	return dashboardViewModel{
		Workflows: workflowVMs,
		Runs:      runVMs,
	}
}

func collectStepNames(steps []workflow.Step) []string {
	names := make([]string, 0, len(steps))
	for _, step := range steps {
		if name := strings.TrimSpace(step.Name); name != "" {
			names = append(names, name)
		}
	}
	return names
}

func sortedStepNames(stepStates map[string]workflow.StepState) []string {
	names := make([]string, 0, len(stepStates))
	for name := range stepStates {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
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
			// still successful unless other statuses contradict
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
