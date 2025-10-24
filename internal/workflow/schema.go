package workflow

// Step represents a single step in a workflow
type Step struct {
	Name        string   `toml:"name" json:"name"`
	Description string   `toml:"description" json:"description"`
	Handler     string   `toml:"handler" json:"handler"` // "tool" (default), "human"
	Prompt      string   `toml:"prompt" json:"prompt"`   // Instructions for cognitive handlers
	Content     string   `toml:"content" json:"content"` // Inline content for steps with no inputs
	Inputs      []string `toml:"inputs" json:"inputs"`
	Output      string   `toml:"output" json:"output"`
}

// Workflow represents a workflow definition
type Workflow struct {
	// ID is the workflow identifier derived from the filename (not stored in TOML)
	ID          string `toml:"-" json:"id"`
	DisplayName string `toml:"display_name" json:"display_name"`
	Description string `toml:"description" json:"description"`
	Message     string `toml:"message" json:"message"`
	Steps       []Step `toml:"steps" json:"steps"`
}
