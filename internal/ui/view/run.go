package view

// RunColumn describes the runs column rendered on the dashboard.
type RunColumn struct {
	Header ColumnHeader
	Runs   []Run
}

// Run summarizes a workflow run and its current state.
type Run struct {
	DisplayName  string
	ID           string
	StateLabel   string
	StateClass   string
	WorkflowName string
	Steps        []RunStep
}

// RunStep represents a single workflow step state within a run.
type RunStep struct {
	Name        string
	Status      string
	StatusClass string
}
