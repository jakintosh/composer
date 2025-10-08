package workflow

// Step represents a single step in a workflow
type Step struct {
	Name        string   `toml:"name"`
	Description string   `toml:"description"`
	Inputs      []string `toml:"inputs"`
	Output      string   `toml:"output"`
}

// Workflow represents a workflow definition
type Workflow struct {
	// ID is the workflow identifier derived from the filename (not stored in TOML)
	ID          string `toml:"-"`
	Title       string `toml:"title"`
	Description string `toml:"description"`
	Message     string `toml:"message"`
	Steps       []Step `toml:"steps"`
}
