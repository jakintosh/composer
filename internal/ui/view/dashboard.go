package view

// DashboardTemplate is the page template used when rendering the dashboard.
const DashboardTemplate = "pages/dashboard"

// Dashboard aggregates the components required to render the dashboard.
type Dashboard struct {
	Sidebar        Sidebar
	WorkflowColumn WorkflowColumn
	WorkflowModal  WorkflowModal
	RunColumn      RunColumn
	TaskColumn     WaitingColumn
}
