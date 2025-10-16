package ui

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"composer/internal/workflow"
)

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

var dashboardTemplate = template.Must(template.New("dashboard").Parse(
	`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>Composer Workflow Dashboard</title>
	<style>
		body {
			margin: 0;
			background: #0f0f0f;
			color: #ffffff;
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
			line-height: 1.5;
		}
		main {
			margin: 1.5rem auto;
			max-width: 960px;
		}
		h1, h2 {
			margin-top: 0;
		}
		.columns {
			display: flex;
			gap: 1.5rem;
			align-items: flex-start;
			flex-wrap: wrap;
		}
		.column {
			background: #1a1a1a;
			flex: 1 1 280px;
			padding: 0.75rem 1rem;
			border-radius: 6px;
			box-sizing: border-box;
		}
		.column ul {
			margin: 0;
			padding: 0;
		}
		.state {
			font-weight: 600;
		}
		.state-ready {
			color: #6cb4ff;
		}
		.state-succeeded {
			color: #65d57c;
		}
		.state-failed {
			color: #ff6b6b;
		}
		.state-pending {
			color: #ffd966;
		}
		.state-unknown {
			color: #bbbbbb;
		}
		.item-list {
			list-style: none;
		}
		.list-item {
			margin: 0.5rem 0;
			background: #232323;
			border-radius: 6px;
			border: 1px solid #2d2d2d;
			overflow: hidden;
		}
		.list-item details {
			display: block;
		}
		.list-item summary {
			cursor: pointer;
			padding: 0.6rem 0.8rem;
			margin: 0;
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 0.75rem;
			font-weight: 600;
		}
		.list-item summary::-webkit-details-marker {
			display: none;
		}
		.item-title {
			flex: 1;
		}
		.item-details {
			padding: 0.6rem 0.8rem 0.8rem;
			font-size: 0.95rem;
			color: #dddddd;
		}
		.item-details p {
			margin: 0.2rem 0;
		}
		.item-details h3 {
			margin: 0.6rem 0 0.3rem;
			font-size: 1rem;
		}
		.detail-list {
			list-style: none;
			padding: 0;
			margin: 0.2rem 0 0;
		}
		.detail-list li {
			margin: 0.2rem 0;
		}
		.detail-list .state {
			font-weight: 500;
		}
	</style>
</head>
<body>
	<main>
		<h1>Workflow Dashboard</h1>
		<div class="columns">
			<section class="column">
				<h2>Workflows</h2>
				{{if .Workflows}}
				<ul class="item-list">
					{{range .Workflows}}
					<li class="list-item">
						<details>
							<summary>
								<span class="item-title">{{.DisplayName}}</span>
							</summary>
							<div class="item-details">
								<p><strong>ID:</strong> {{.Workflow.ID}}</p>
								{{if .Workflow.Title}}<p><strong>Title:</strong> {{.Workflow.Title}}</p>{{end}}
								{{if .Workflow.Description}}<p><strong>Description:</strong> {{.Workflow.Description}}</p>{{end}}
								{{if .Workflow.Message}}<p><strong>Message:</strong> {{.Workflow.Message}}</p>{{end}}
								{{if .StepNames}}
								<h3>Steps</h3>
								<ul class="detail-list">
									{{range .StepNames}}
									<li>{{.}}</li>
									{{end}}
								</ul>
								{{end}}
							</div>
						</details>
					</li>
					{{end}}
				</ul>
				{{else}}
				<p>No workflows available.</p>
				{{end}}
			</section>
			<section class="column">
				<h2>Runs</h2>
				{{if .Runs}}
				<ul class="item-list">
					{{range .Runs}}
					<li class="list-item">
						<details>
							<summary>
								<span class="item-title">{{.Name}}</span>
								<span class="item-state state {{.StateClass}}">{{.StateLabel}}</span>
							</summary>
							<div class="item-details">
								<p><strong>Workflow:</strong> {{.WorkflowName}}</p>
								{{if .Steps}}
								<h3>Steps</h3>
								<ul class="detail-list">
									{{range .Steps}}
									<li>
										<span>{{.Name}}</span>
										<span class="state {{.StatusClass}}">— {{.Status}}</span>
									</li>
									{{end}}
								</ul>
								{{end}}
							</div>
						</details>
					</li>
					{{end}}
				</ul>
				{{else}}
				<p>No runs found.</p>
				{{end}}
			</section>
		</div>
	</main>
</body>
</html>
`))

// BuildRouter creates and configures the UI router.
func BuildRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleDashboard)
	return mux
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	workflows, err := workflow.ListWorkflows()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load workflows: %v", err), http.StatusInternalServerError)
		return
	}

	runs, err := workflow.ListRuns()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load runs: %v", err), http.StatusInternalServerError)
		return
	}

	data := buildDashboardData(workflows, runs)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := dashboardTemplate.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("failed to render dashboard: %v", err), http.StatusInternalServerError)
		return
	}
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
