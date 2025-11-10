package view

// WorkflowColumn describes the workflow column panel on the dashboard.
type WorkflowColumn struct {
	Header    ColumnHeader
	Workflows []Workflow
}

// Workflow captures the essential display information for a workflow card.
type Workflow struct {
	DisplayName string
	ID          string
	Description string
	Message     string
	StepNames   []string
}

// WorkflowModal holds the controls displayed in the workflow creation modal.
type WorkflowModal struct {
	AddStepButton Button
}
