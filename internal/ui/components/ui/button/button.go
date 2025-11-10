package button

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
)

//go:embed button.tmpl
var buttonTemplate string

var tmpl = templates.New(
	"ui_button",
	"components/ui/button/button.tmpl",
	buttonTemplate,
	nil,
)

// Props describes the configuration for rendering a button with optional icon.
type Props struct {
	ID        string
	Class     string
	Title     string
	AriaLabel string
	Label     string
	Type      string
	IconSize  int
}

// Render executes the template with the provided props and returns safe HTML.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender renders the button, returning a safe placeholder comment if an
// error occurs. Useful when composing components in templates.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
