package ui

import (
	"sort"
	"strings"

	"composer/internal/orchestrator"
	sidebar "composer/internal/ui/components/navigation/sidebar"
	runcolumn "composer/internal/ui/components/run/column"
	"composer/internal/ui/components/ui/button"
	"composer/internal/ui/components/ui/columnheader"
	waitingcolumn "composer/internal/ui/components/waiting/column"
	workflowcolumn "composer/internal/ui/components/workflow/column"
	dashboardpage "composer/internal/ui/pages/dashboard"
	"composer/internal/workflow"
)

func buildDashboardModel(
	workflows []workflow.Workflow,
	runs []workflow.RunState,
	waitingTasks map[string][]orchestrator.WaitingTask,
) dashboardpage.Props {
	workflowVMs := make([]workflowcolumn.Workflow, 0, len(workflows))
	for _, wf := range workflows {
		workflowVMs = append(workflowVMs, workflowcolumn.Workflow{
			DisplayName: strings.TrimSpace(wf.DisplayName),
			ID:          strings.TrimSpace(wf.ID),
			Description: strings.TrimSpace(wf.Description),
			Message:     strings.TrimSpace(wf.Message),
			StepNames:   collectStepNames(wf.Steps),
		})
	}
	sort.Slice(workflowVMs, func(i, j int) bool {
		return workflowVMs[i].DisplayName < workflowVMs[j].DisplayName
	})

	runVMs := make([]runcolumn.Run, 0, len(runs))
	for _, runState := range runs {
		status := summarizeRunState(runState)
		stepNames := sortedStepNames(runState.StepStates)

		displayName := strings.TrimSpace(runState.Name)
		runID := strings.TrimSpace(runState.ID)

		steps := make([]runcolumn.Step, 0, len(stepNames))
		for _, name := range stepNames {
			stepState := runState.StepStates[name]
			steps = append(steps, runcolumn.Step{
				Name:        name,
				Status:      string(stepState.Status),
				StatusClass: stateClassForStatus(stepState.Status),
			})
		}

		runVMs = append(runVMs, runcolumn.Run{
			DisplayName:  displayName,
			ID:           runID,
			StateLabel:   status.Label,
			StateClass:   status.Class,
			WorkflowName: strings.TrimSpace(runState.WorkflowName),
			Steps:        steps,
		})
	}
	sort.Slice(runVMs, func(i, j int) bool {
		return runVMs[i].DisplayName < runVMs[j].DisplayName
	})

	taskGroupVMs := make([]waitingcolumn.Group, 0, len(waitingTasks))
	for _, runState := range runs {
		runID := strings.TrimSpace(runState.ID)
		displayName := strings.TrimSpace(runState.Name)
		if displayName == "" {
			displayName = runID
		}

		tasks := waitingTasks[runID]
		if len(tasks) == 0 {
			continue
		}

		taskVMs := make([]waitingcolumn.Task, 0, len(tasks))
		for _, task := range tasks {
			taskVMs = append(taskVMs, waitingcolumn.Task{
				Name:        strings.TrimSpace(task.Name),
				Description: strings.TrimSpace(task.Description),
				Prompt:      strings.TrimSpace(task.Prompt),
			})
		}

		sort.Slice(taskVMs, func(i, j int) bool {
			return taskVMs[i].Name < taskVMs[j].Name
		})

		taskGroupVMs = append(taskGroupVMs, waitingcolumn.Group{
			RunID:          runID,
			RunDisplayName: displayName,
			WorkflowName:   strings.TrimSpace(runState.WorkflowName),
			TaskCount:      len(taskVMs),
			Tasks:          taskVMs,
		})
	}
	sort.Slice(taskGroupVMs, func(i, j int) bool {
		return taskGroupVMs[i].RunDisplayName < taskGroupVMs[j].RunDisplayName
	})

	return dashboardpage.Props{
		Sidebar: sidebar.Props{
			Title: "Composer",
			Links: []sidebar.Link{
				{
					Label:  "Dashboard",
					Href:   "/",
					Active: true,
				},
			},
		},
		WorkflowColumn: workflowcolumn.Props{
			Header: columnheader.Props{
				Title: "Workflows",
				Actions: []button.Props{
					{
						ID:        "open-workflow-modal",
						Class:     "button--accent button--icon",
						Title:     "Create workflow",
						AriaLabel: "Create workflow",
						Type:      "button",
						IconSize:  16,
					},
				},
			},
			Workflows: workflowVMs,
		},
		WorkflowModal: dashboardpage.DefaultWorkflowModal(),
		RunColumn: runcolumn.Props{
			Header: columnheader.Props{Title: "Runs"},
			Runs:   runVMs,
		},
		RunModal: dashboardpage.DefaultRunModal(),
		TaskColumn: waitingcolumn.Props{
			Header: columnheader.Props{Title: "Tasks"},
			Groups: taskGroupVMs,
		},
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

func stateClassForStatus(status workflow.StepStatus) string {
	switch status {
	case workflow.StatusFailed:
		return "status-badge--failed"
	case workflow.StatusSucceeded:
		return "status-badge--succeeded"
	case workflow.StatusReady:
		return "status-badge--ready"
	case workflow.StatusPending:
		return "status-badge--pending"
	default:
		return "status-badge--unknown"
	}
}

type runStatus struct {
	Label string
	Class string
}

func summarizeRunState(rs workflow.RunState) runStatus {
	if len(rs.StepStates) == 0 {
		return runStatus{Label: "pending", Class: "status-badge--pending"}
	}

	allSucceeded := len(rs.StepStates) > 0
	hasReady := false
	hasPending := false

	for _, step := range rs.StepStates {
		switch step.Status {
		case workflow.StatusFailed:
			return runStatus{Label: "failed", Class: "status-badge--failed"}
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
		return runStatus{Label: "succeeded", Class: "status-badge--succeeded"}
	}
	if hasReady {
		return runStatus{Label: "ready", Class: "status-badge--ready"}
	}
	if hasPending {
		return runStatus{Label: "pending", Class: "status-badge--pending"}
	}
	return runStatus{Label: "unknown", Class: "status-badge--unknown"}
}
