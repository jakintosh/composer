package view

// Button describes the structure of a reusable UI button.
type Button struct {
	ID        string
	Class     string
	Title     string
	AriaLabel string
	Label     string
	Type      string
	IconSize  int
}

// ColumnHeader defines the metadata rendered at the top of a column panel.
type ColumnHeader struct {
	Title   string
	Actions []Button
}
