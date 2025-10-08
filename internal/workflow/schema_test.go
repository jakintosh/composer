package workflow

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestStepUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		toml     string
		expected Step
	}{
		{
			name: "step with no inputs",
			toml: `
name = "init"
description = "Initialize workflow"
output = "initialized"
`,
			expected: Step{
				Name:        "init",
				Description: "Initialize workflow",
				Inputs:      nil,
				Output:      "initialized",
			},
		},
		{
			name: "step with single input",
			toml: `
name = "process"
description = "Process data"
inputs = ["initialized"]
output = "processed"
`,
			expected: Step{
				Name:        "process",
				Description: "Process data",
				Inputs:      []string{"initialized"},
				Output:      "processed",
			},
		},
		{
			name: "step with multiple inputs",
			toml: `
name = "combine"
description = "Combine results"
inputs = ["processed", "validated"]
output = "combined"
`,
			expected: Step{
				Name:        "combine",
				Description: "Combine results",
				Inputs:      []string{"processed", "validated"},
				Output:      "combined",
			},
		},
		{
			name: "step with inline content",
			toml: `
name = "init"
description = "Initialize with content"
content = "This is initial content"
output = "initialized"
`,
			expected: Step{
				Name:        "init",
				Description: "Initialize with content",
				Content:     "This is initial content",
				Output:      "initialized",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var step Step
			err := toml.Unmarshal([]byte(tt.toml), &step)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			if step.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", step.Name, tt.expected.Name)
			}
			if step.Description != tt.expected.Description {
				t.Errorf("Description = %v, want %v", step.Description, tt.expected.Description)
			}
			if step.Content != tt.expected.Content {
				t.Errorf("Content = %v, want %v", step.Content, tt.expected.Content)
			}
			if step.Output != tt.expected.Output {
				t.Errorf("Output = %v, want %v", step.Output, tt.expected.Output)
			}

			if len(step.Inputs) != len(tt.expected.Inputs) {
				t.Errorf("Inputs length = %v, want %v", len(step.Inputs), len(tt.expected.Inputs))
			} else {
				for i, input := range step.Inputs {
					if input != tt.expected.Inputs[i] {
						t.Errorf("Inputs[%d] = %v, want %v", i, input, tt.expected.Inputs[i])
					}
				}
			}
		})
	}
}

func TestWorkflowWithSteps(t *testing.T) {
	tomlData := `
title = "Test Workflow"
description = "A test workflow"
message = "Testing steps"

[[steps]]
name = "start"
description = "Start step"
output = "started"

[[steps]]
name = "middle"
description = "Middle step"
inputs = ["started"]
output = "processed"

[[steps]]
name = "end"
description = "End step"
inputs = ["processed"]
output = "finished"
`

	var workflow Workflow
	err := toml.Unmarshal([]byte(tomlData), &workflow)
	if err != nil {
		t.Fatalf("failed to unmarshal workflow: %v", err)
	}

	if workflow.Title != "Test Workflow" {
		t.Errorf("Title = %v, want Test Workflow", workflow.Title)
	}

	if len(workflow.Steps) != 3 {
		t.Fatalf("Steps length = %v, want 3", len(workflow.Steps))
	}

	// Verify first step
	if workflow.Steps[0].Name != "start" {
		t.Errorf("Steps[0].Name = %v, want start", workflow.Steps[0].Name)
	}
	if workflow.Steps[0].Output != "started" {
		t.Errorf("Steps[0].Output = %v, want started", workflow.Steps[0].Output)
	}
	if len(workflow.Steps[0].Inputs) != 0 {
		t.Errorf("Steps[0].Inputs length = %v, want 0", len(workflow.Steps[0].Inputs))
	}

	// Verify middle step
	if workflow.Steps[1].Name != "middle" {
		t.Errorf("Steps[1].Name = %v, want middle", workflow.Steps[1].Name)
	}
	if len(workflow.Steps[1].Inputs) != 1 || workflow.Steps[1].Inputs[0] != "started" {
		t.Errorf("Steps[1].Inputs = %v, want [started]", workflow.Steps[1].Inputs)
	}

	// Verify end step
	if workflow.Steps[2].Name != "end" {
		t.Errorf("Steps[2].Name = %v, want end", workflow.Steps[2].Name)
	}
	if len(workflow.Steps[2].Inputs) != 1 || workflow.Steps[2].Inputs[0] != "processed" {
		t.Errorf("Steps[2].Inputs = %v, want [processed]", workflow.Steps[2].Inputs)
	}
}
