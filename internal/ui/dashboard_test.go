package ui

import (
	"testing"

	"composer/internal/orchestrator"
	"composer/internal/workflow"
)

func TestBuildDashboardModel(t *testing.T) {
	workflows := []workflow.Workflow{
		{
			DisplayName: "Beta Flow",
			ID:          "wf-beta",
			Description: "B",
			Message:     "Beta",
			Steps: []workflow.Step{
				{Name: " beta-step "},
			},
		},
		{
			DisplayName: "Alpha Flow",
			ID:          "wf-alpha",
			Description: "A",
			Message:     "Alpha",
			Steps: []workflow.Step{
				{Name: "first"},
				{Name: "Second"},
			},
		},
	}

	runs := []workflow.RunState{
		{
			ID:           "run-b",
			Name:         "Run B",
			WorkflowName: "Beta Flow",
			StepStates: map[string]workflow.StepState{
				"beta-step": {Status: workflow.StatusSucceeded},
			},
		},
		{
			ID:           "run-a",
			Name:         "Run A",
			WorkflowName: "Alpha Flow",
			StepStates: map[string]workflow.StepState{
				"first":  {Status: workflow.StatusPending},
				"second": {Status: workflow.StatusReady},
			},
		},
	}

	waiting := map[string][]orchestrator.WaitingTask{
		"run-a": {
			{
				Name:        "Review",
				Description: "Check",
				Prompt:      "Do it",
			},
			{
				Name:        "Approve",
				Description: "OK",
			},
		},
	}

	model := buildDashboardModel(workflows, runs, waiting)

	if model.Sidebar.Title != "Composer" {
		t.Fatalf("Sidebar title = %q, want %q", model.Sidebar.Title, "Composer")
	}

	if got := model.WorkflowColumn.Workflows[0].DisplayName; got != "Alpha Flow" {
		t.Fatalf("first workflow = %q, want Alpha Flow", got)
	}

	if got := len(model.WorkflowColumn.Workflows[0].StepNames); got != 2 {
		t.Fatalf("Alpha Flow step count = %d, want 2", got)
	}

	if got := model.RunColumn.Runs[0].DisplayName; got != "Run A" {
		t.Fatalf("run order incorrect, first run = %q want Run A", got)
	}

	if got := model.RunColumn.Runs[0].Steps[0].Name; got != "first" {
		t.Fatalf("run step ordering incorrect, first step = %q want first", got)
	}

	if got := len(model.TaskColumn.Groups); got != 1 {
		t.Fatalf("task groups = %d, want 1", got)
	}

	if got := model.TaskColumn.Groups[0].TaskCount; got != 2 {
		t.Fatalf("task count = %d, want 2", got)
	}
}

func TestSummarizeRunState(t *testing.T) {
	tests := []struct {
		name     string
		state    workflow.RunState
		expected runStatus
	}{
		{
			name:     "no steps defaults to pending",
			state:    workflow.RunState{StepStates: map[string]workflow.StepState{}},
			expected: runStatus{Label: "pending", Class: "status-badge--pending"},
		},
		{
			name: "failed overrides",
			state: workflow.RunState{
				StepStates: map[string]workflow.StepState{
					"a": {Status: workflow.StatusFailed},
					"b": {Status: workflow.StatusSucceeded},
				},
			},
			expected: runStatus{Label: "failed", Class: "status-badge--failed"},
		},
		{
			name: "all succeeded",
			state: workflow.RunState{
				StepStates: map[string]workflow.StepState{
					"a": {Status: workflow.StatusSucceeded},
				},
			},
			expected: runStatus{Label: "succeeded", Class: "status-badge--succeeded"},
		},
		{
			name: "ready takes precedence over pending",
			state: workflow.RunState{
				StepStates: map[string]workflow.StepState{
					"a": {Status: workflow.StatusReady},
					"b": {Status: workflow.StatusPending},
				},
			},
			expected: runStatus{Label: "ready", Class: "status-badge--ready"},
		},
		{
			name: "pending fallback",
			state: workflow.RunState{
				StepStates: map[string]workflow.StepState{
					"a": {Status: workflow.StatusPending},
				},
			},
			expected: runStatus{Label: "pending", Class: "status-badge--pending"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := summarizeRunState(tt.state); got != tt.expected {
				t.Fatalf("summarizeRunState() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestStateClassForStatus(t *testing.T) {
	cases := map[workflow.StepStatus]string{
		workflow.StatusFailed:    "status-badge--failed",
		workflow.StatusSucceeded: "status-badge--succeeded",
		workflow.StatusReady:     "status-badge--ready",
		workflow.StatusPending:   "status-badge--pending",
		workflow.StepStatus("x"): "status-badge--unknown",
	}

	for status, want := range cases {
		if got := stateClassForStatus(status); got != want {
			t.Fatalf("stateClassForStatus(%q) = %q, want %q", status, got, want)
		}
	}
}
