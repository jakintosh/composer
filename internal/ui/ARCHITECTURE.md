# Go UI Component Architecture

This document outlines the architecture for building reusable, testable, and composable server-side rendered UI components in Go.

## Core Principles

1. Component-Based: Each UI element (Button, Card, Page) is a self-contained package.
2. Props-Driven: Components are pure functions of their input properties (Props struct).
3. Standard Library: Leverages html/template, go:embed, and standard Go testing.
4. Golden File Testing: Ensures visual consistency and catches HTML regressions without brittle DOM assertions.

## Component Structure

Each component resides in its own package and consists of three key files:

### 1. The Template (component.tmpl)

Standard html/template file defining the raw HTML structure. It uses data from the Props struct.
```html
<!-- internal/ui/components/ui/button/button.tmpl -->
<button class="button {{ .Variant }}">
  {{.Label}}
</button>
```

### 2. The Go Logic (component.go)
Defines the Props struct, embeds the template, and provides Render helpers. We use the
central `templates.Manager` to parse embedded templates once in production and reload
from disk in development.

```go
// internal/ui/components/ui/button/button.go
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

type Props struct {
    Label string
}

func Render(p Props) (template.HTML, error) {
    return tmpl.Render(p)
}

func MustRender(p Props) template.HTML {
    return templates.SafeHTML(Render(p))
}
```

### 3. The Test (component_test.go)

Uses golden file testing to snapshot the rendered HTML.

```go
func TestButton(t *testing.T) {
    html, err := Render(Props{Label: "Save"})
    if err != nil {
        t.Fatalf("render error: %v", err)
    }

    golden.Assert(t, string(html), filepath.Join("fixtures", "save.golden"))
}
```

## Composition Pattern

Complex components build upon simpler ones by including their Props structs and exposing helper methods for rendering.

### Parent Component (usercard.go):
```go
type UserCardProps struct {
    UserName    string
    ButtonProps button.Props // Embed child props
}

// Helper method called by parent template: {{ .RenderButton }}
func (p UserCardProps) RenderButton() template.HTML {
    return p.ButtonProps.Render()
}
```

### Parent Template (usercard.tmpl):
```html
<div class="card">
    <h2>{{ .UserName }}</h2>
    <!-- Renders child component safely -->
    {{ .RenderButton }}
</div>
```

This pattern ensures encapsulation while allowing parent components full control over their children's data and rendering context.

## Template Manager

`internal/ui/templates.Manager` centralizes template loading:

- When `APP_ENV=development`, templates are read from disk on each render using the
  configured project root (`templates.Configure(uiSourceRoot())`).
- In other environments, go:embedded template strings are parsed once and cached with
  `sync.Once` for maximum performance.

Components simply register themselves with `templates.New`, avoiding any per-component
environment branching.
