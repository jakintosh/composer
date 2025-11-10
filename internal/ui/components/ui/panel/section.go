package panel

import (
	_ "embed"
	"html/template"

	"composer/internal/ui/components/ui/columnheader"
	"composer/internal/ui/templates"
)

//go:embed section.tmpl
var sectionTemplate string

var sectionTmpl = templates.New(
	"ui_panel_section",
	"components/ui/panel/section.tmpl",
	sectionTemplate,
	nil,
)

// Variant controls the visual style of a panel container.
type Variant string

const (
	// VariantDefault renders the standard panel surface.
	VariantDefault Variant = ""
	// VariantMuted renders the low-contrast panel surface.
	VariantMuted Variant = "panel--muted"
)

// SectionProps describes a panel section with a header and optional body.
type SectionProps struct {
	Header  columnheader.Props
	Variant Variant
	Body    template.HTML
}

// VariantClass exposes the BEM modifier to append to the panel container.
func (p SectionProps) VariantClass() string {
	return string(p.Variant)
}

// HasBody reports whether a body fragment was provided.
func (p SectionProps) HasBody() bool {
	return len(p.Body) > 0
}

// RenderHeader renders the reusable column header.
func (p SectionProps) RenderHeader() template.HTML {
	return templates.SafeHTML(columnheader.Render(p.Header))
}

// RenderBody returns the already-rendered body fragment.
func (p SectionProps) RenderBody() template.HTML {
	return p.Body
}

// RenderSection executes the panel section template.
func RenderSection(p SectionProps) (template.HTML, error) {
	return sectionTmpl.Render(p)
}

// MustRenderSection renders the panel section, returning an HTML comment if
// rendering fails so parent components remain resilient.
func MustRenderSection(p SectionProps) template.HTML {
	return templates.SafeHTML(RenderSection(p))
}
