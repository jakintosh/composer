# Composer

A workflow orchestration tool that executes steps based on declarative TOML definitions.

## Overview

Composer allows you to define workflows as a series of steps with dependencies, then execute those workflows with automatic parallelization and state management. Steps run when their input requirements are met, enabling both sequential and parallel execution patterns.

## Core Concepts

### Workflows
Workflows are defined in TOML files with metadata and a list of steps. Each workflow has:
- **ID**: Derived from the filename (e.g., `example.toml` → `example`)
- **Metadata**: Title, description, and optional message
- **Steps**: An ordered list of step definitions

Example workflow:
```toml
title = "Data Processing Pipeline"
description = "Fetch, process, and validate data"

[[steps]]
name = "fetch-data"
description = "Fetch data from source"
output = "raw-data"

[[steps]]
name = "process-data"
description = "Process the raw data"
inputs = ["raw-data"]
output = "processed-data"

[[steps]]
name = "review-quality"
handler = "human"
description = "Review data quality"
prompt = "Check the processed data for anomalies and approve if it looks good"
inputs = ["processed-data"]
output = "reviewed-data"
```

### Steps
Steps are the individual units of work in a workflow:
- **Name**: Unique identifier for the step
- **Description**: Human-readable description
- **Handler**: Who executes the step - `"tool"` (default, automated) or `"human"` (requires intervention)
- **Prompt**: Instructions for human handlers (optional)
- **Inputs**: List of required outputs from other steps (optional)
- **Output**: Name of the output this step produces

Steps with no inputs can run immediately. Steps with inputs wait until all required outputs are available.

**Handler Types:**
- **tool** (default): Automated steps that execute immediately when dependencies are met
- **human**: Steps requiring human intervention; transition to "ready" status and must be completed via the `do` command

### Runs
A run is an instantiated workflow with state. When you execute a workflow, Composer creates a run directory at `.composer/runs/{run-name}/` (relative to your current directory) that tracks:
- **Workflow name**: Which workflow this run executes
- **Step states**: Status of each step (`pending`, `ready`, `succeeded`, `failed`)
- **Outputs**: List of outputs that have been produced

**Step Statuses:**
- **pending**: Waiting for input dependencies
- **ready**: Human-handler step with dependencies met, awaiting intervention
- **succeeded**: Step completed successfully
- **failed**: Step failed (not yet implemented)

State is persisted as JSON between ticks, allowing you to stop and resume execution.

## Project Structure

```
.
├── bin/                   # Built binary (created by make build)
├── cmd/composer/          # CLI entry point
├── internal/
│   ├── orchestrator/      # Run creation and tick execution
│   └── workflow/          # Workflow loading, state management, paths
└── Makefile               # Build tasks
```

## Building

```bash
make build
```

This creates the `bin/composer` executable.

## Usage

### Create and start a workflow run
```bash
./bin/composer run <workflow-name> <run-name>
```

This loads a workflow, creates a new run with initial state, and executes the first tick.

### Continue execution (tick)
```bash
./bin/composer tick <run-name>
```

Executes one tick: finds all runnable steps (those with satisfied inputs), runs tool steps in parallel, transitions human steps to "ready" status, updates state, and saves.

### List waiting tasks
```bash
./bin/composer tasks <run-name>
```

Lists all tasks with "ready" status that are waiting for human intervention. Each task is shown with an index, description, prompt (if provided), inputs, and output.

### Complete a waiting task
```bash
./bin/composer do <run-name> <task-index>
```

Marks a waiting task as completed, adding its output to the run state. Use the task index from the `tasks` command.

### Example
```bash
# Start a run of the example workflow
./bin/composer run example my-first-run

# Continue execution until complete
./bin/composer tick my-first-run
./bin/composer tick my-first-run
# ... repeat until "Workflow complete!" appears
```

### Example with Human Tasks
```bash
# Run a workflow with human intervention steps
./bin/composer run review-workflow my-review

# Tick until a human task is ready
./bin/composer tick my-review

# List waiting tasks
./bin/composer tasks my-review

# Complete task at index 0
./bin/composer do my-review 0

# Continue workflow
./bin/composer tick my-review
```

## Runtime Directories

When you run Composer, it looks for workflows and stores run state in specific locations:

### Workflow Search Paths
Composer searches for workflow TOML files in this order:
1. `./.composer/workflows/` (current directory)
2. `$XDG_DATA_HOME/composer/workflows/` (or `~/.local/share/composer/workflows/`)
3. `/etc/composer/workflows/` (system-wide)

### Run Storage
Runs are always stored in `./.composer/runs/` relative to the current directory where you execute the `composer` command. Each run gets its own subdirectory containing `state.json`.

## Current Status

This is an early-stage project. Current functionality:
- Loading workflows from TOML files
- Creating runs with initial state
- Tick-based execution with dependency resolution
- Parallel step execution when dependencies allow
- State persistence between ticks
- Human intervention steps that pause workflow execution
- Task listing and completion via CLI

Tool handler steps currently just print their execution (no actual work performed yet). Human handler steps require manual completion via the `do` command.

## Architecture Notes

### Orchestrator (`internal/orchestrator/`)
- **CreateRun**: Initializes a new run with pending steps
- **Tick**: Executes one cycle of the workflow (find runnable steps → run in parallel → save state)
- **findRunnableSteps**: Determines which pending steps have all inputs satisfied

### Workflow Package (`internal/workflow/`)
- **loader.go**: Searches for and loads workflow TOML files
- **schema.go**: Workflow and Step data structures
- **state.go**: RunState management, persistence, and helper methods
- **paths.go**: Path resolution for workflows and runs

### CLI (`cmd/composer/`)
Two commands:
- `run`: Loads workflow, creates run, executes first tick
- `tick`: Loads existing run state, executes one tick
