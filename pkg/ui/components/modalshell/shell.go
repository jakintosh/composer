package shell

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
)

//go:embed shell.tmpl
var shellTemplate string

var tmpl = templates.New(
	"ui_modal_shell",
	"pkg/ui/components/modalshell/shell.tmpl",
	shellTemplate,
	nil,
)

// Props defines the shared modal chrome rendered around modal-specific bodies.
type Props struct {
	ID         string
	Title      string
	TitleID    string
	CloseLabel string
	Body       template.HTML
}

// DialogLabelID returns the aria-labelledby target for the modal dialog.
func (p Props) DialogLabelID() string {
	if p.TitleID != "" {
		return p.TitleID
	}
	if p.ID != "" {
		return p.ID + "-title"
	}
	return "modal-title"
}

// CloseLabelText returns the accessible label for the close button.
func (p Props) CloseLabelText() string {
	if p.CloseLabel != "" {
		return p.CloseLabel
	}
	return "Close dialog"
}

// RenderBody returns the rendered modal body fragment.
func (p Props) RenderBody() template.HTML {
	return p.Body
}

// Render executes the modal shell template.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender renders the modal shell, emitting an HTML comment on failure to
// keep parent components resilient.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
