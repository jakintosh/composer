package view

// WaitingColumn represents the waiting-task column on the dashboard.
type WaitingColumn struct {
	Header ColumnHeader
	Groups []WaitingGroup
}

// WaitingGroup aggregates pending human tasks for a specific run.
type WaitingGroup struct {
	RunID          string
	RunDisplayName string
	WorkflowName   string
	TaskCount      int
	Tasks          []WaitingTask
}

// WaitingTask represents a single waiting task awaiting user interaction.
type WaitingTask struct {
	Name        string
	Description string
	Prompt      string
}
