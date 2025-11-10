package datalist

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/templates"
)

//go:embed datalist.tmpl
var datalistTemplate string

var tmpl = templates.New(
	"ui_data_list",
	"components/ui/datalist/datalist.tmpl",
	datalistTemplate,
	nil,
)

// Item represents a row in the data list.
type Item struct {
	Primary   string
	Secondary template.HTML
}

// HasSecondary reports whether the row has a secondary value.
func (i Item) HasSecondary() bool {
	return len(i.Secondary) > 0
}

// RenderSecondary returns the rendered HTML for the secondary value.
func (i Item) RenderSecondary() template.HTML {
	return i.Secondary
}

// Props describes a collection of rows rendered in a data list.
type Props struct {
	Items []Item
}

// HasItems reports whether the list contains any rows.
func (p Props) HasItems() bool {
	return len(p.Items) > 0
}

// Render executes the data list template.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender renders the data list, returning an HTML comment on failure.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
