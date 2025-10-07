package workflow

// Workflow represents a simple throwaway workflow definition for testing.
// This schema will be replaced with the actual workflow structure later.
type Workflow struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Message     string `toml:"message"`
}
