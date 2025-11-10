package panel

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
)

//go:embed list.tmpl
var listTemplate string

var listTmpl = templates.New(
	"ui_panel_list",
	"components/ui/panel/list.tmpl",
	listTemplate,
	nil,
)

// ListProps defines the structure of a standard panel list with an empty state.
type ListProps struct {
	Items        []ListItemProps
	ListClass    string
	EmptyMessage string
}

// ListItemProps describes a single list row.
type ListItemProps struct {
	Class          string
	Content        template.HTML
	DisableWrapper bool
}

// HasItems reports whether the list contains any items.
func (p ListProps) HasItems() bool {
	return len(p.Items) > 0
}

// HasEmptyMessage reports whether an empty message is configured.
func (p ListProps) HasEmptyMessage() bool {
	return len(p.EmptyMessage) > 0
}

// RenderContent returns the rendered HTML fragment for a list item.
func (p ListItemProps) RenderContent() template.HTML {
	return p.Content
}

// ShouldWrap reports whether the panel list should wrap the content with an <li>.
func (p ListItemProps) ShouldWrap() bool {
	return !p.DisableWrapper
}

// RenderList executes the list template.
func RenderList(p ListProps) (template.HTML, error) {
	return listTmpl.Render(p)
}

// MustRenderList renders the list, returning an HTML comment when rendering
// fails so parent components can continue rendering gracefully.
func MustRenderList(p ListProps) template.HTML {
	return templates.SafeHTML(RenderList(p))
}
