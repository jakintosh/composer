package ui

import (
	"sort"
	"strings"

	"composer/internal/orchestrator"
	"composer/internal/workflow"
)

type dashboardViewModel struct {
	Sidebar        sidebarViewModel
	WorkflowColumn workflowColumnViewModel
	WorkflowModal  workflowModalViewModel
	RunColumn      runColumnViewModel
	TaskColumn     waitingTaskColumnViewModel
}

type workflowColumnViewModel struct {
	Header    columnHeaderViewModel
	Workflows []workflowViewModel
}

type workflowModalViewModel struct {
	AddStepButton uiButtonViewModel
}

type uiButtonViewModel struct {
	ID        string
	Class     string
	Title     string
	AriaLabel string
	Label     string
	Type      string
	IconSize  int
}

type columnHeaderViewModel struct {
	Title   string
	Actions []uiButtonViewModel
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

type runColumnViewModel struct {
	Header columnHeaderViewModel
	Runs   []runViewModel
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

type waitingTaskColumnViewModel struct {
	Header columnHeaderViewModel
	Groups []waitingTaskGroupViewModel
}

type waitingTaskGroupViewModel struct {
	RunName      string
	WorkflowName string
	TaskCount    int
	Tasks        []waitingTaskViewModel
}

type waitingTaskViewModel struct {
	Name        string
	Description string
	Prompt      string
}

type sidebarViewModel struct {
	Title string
	Links []sidebarLinkViewModel
}

type sidebarLinkViewModel struct {
	Label  string
	Href   string
	Active bool
}

func buildDashboardViewModel(
	workflows []workflow.Workflow,
	runs []workflow.RunState,
	waitingTasks map[string][]orchestrator.WaitingTask,
) dashboardViewModel {
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

	taskGroupVMs := make([]waitingTaskGroupViewModel, 0, len(waitingTasks))
	for _, run := range runs {
		tasks := waitingTasks[run.RunName]
		if len(tasks) == 0 {
			continue
		}

		taskVMs := make([]waitingTaskViewModel, 0, len(tasks))
		for _, task := range tasks {
			taskVMs = append(taskVMs, waitingTaskViewModel{
				Name:        strings.TrimSpace(task.Name),
				Description: strings.TrimSpace(task.Description),
				Prompt:      strings.TrimSpace(task.Prompt),
			})
		}

		sort.Slice(taskVMs, func(i, j int) bool { return taskVMs[i].Name < taskVMs[j].Name })

		taskGroupVMs = append(taskGroupVMs, waitingTaskGroupViewModel{
			RunName:      strings.TrimSpace(run.RunName),
			WorkflowName: strings.TrimSpace(run.WorkflowName),
			TaskCount:    len(taskVMs),
			Tasks:        taskVMs,
		})
	}
	sort.Slice(taskGroupVMs, func(i, j int) bool { return taskGroupVMs[i].RunName < taskGroupVMs[j].RunName })

	return dashboardViewModel{
		Sidebar: sidebarViewModel{
			Title: "Composer",
			Links: []sidebarLinkViewModel{
				{
					Label:  "Dashboard",
					Href:   "/",
					Active: true,
				},
			},
		},
		WorkflowColumn: workflowColumnViewModel{
			Header: columnHeaderViewModel{
				Title: "Workflows",
				Actions: []uiButtonViewModel{
					{
						ID:        "open-workflow-modal",
						Class:     "primary-action",
						Title:     "Create workflow",
						AriaLabel: "Create workflow",
						Type:      "button",
						IconSize:  16,
					},
				},
			},
			Workflows: workflowVMs,
		},
		WorkflowModal: workflowModalViewModel{
			AddStepButton: uiButtonViewModel{
				ID:       "add-workflow-step",
				Class:    "add-step-button",
				Label:    "Add Step",
				Type:     "button",
				IconSize: 16,
			},
		},
		RunColumn: runColumnViewModel{
			Header: columnHeaderViewModel{Title: "Runs"},
			Runs:   runVMs,
		},
		TaskColumn: waitingTaskColumnViewModel{
			Header: columnHeaderViewModel{Title: "Tasks"},
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
