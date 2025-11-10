# Go UI Component Architecture (gomponents)

This document outlines a "Go-first" architecture for building UI components using the `gomponents` library (`maragu.dev/gomponents`). This approach eliminates templates in favor of pure, type-safe Go code for defining, composing, and rendering HTML.

## Core Principles

1.  **Code, Not Templates:** All HTML elements, attributes, and logic are defined as Go functions and structs.
2.  **Type-Safe Composition:** The Go compiler validates the component tree. Composition is achieved by calling other component functions.
3.  **Functions as Components:** A component is a simple Go function that accepts properties (as arguments or a props struct) and returns a `gomponents.Node`.
4.  **Golden File Testing:** Component output is tested by rendering the node to a string and snapshot-testing it against a "golden" file.

## Component Structure

A component is typically a single `.go` file containing its logic and test file. There are no `.tmpl` files.

### 1. The Component (`button.go`)

The component is defined as a function. It uses `gomponents` primitives (`g`), attributes (`ga`), and flow control (`g.If`) to build a node.

```go
// mylib/button/button.go
package button

import (
    "github.com/maragu/gomponents"
    "github.com/maragu/gomponents/el" // g
    "github.com/maragu/gomponents/attr" // ga
)

// Props can be a struct or simple function arguments
type Props struct {
    Label    string
    Primary  bool
    Disabled bool
}

// The component is just a function.
func Button(p Props) g.Node {
    return el.Button(
        // Use standard Go logic
        g.If(p.Primary, ga.Class("btn btn-primary")),
        g.If(!p.Primary, ga.Class("btn")),
        g.If(p.Disabled, ga.Disabled()),

        // Children are just other nodes
        g.Text(p.Label),
    )
}
```

### 2\. The Test (`button_test.go`)

Testing remains almost identical to the template-based approach. We render the component to a string and use `gotest.tools/golden` to snapshot it.

```go
// mylib/button/button_test.go
package button

import (
    "bytes"
    "testing"

    "gotest.tools/v3/golden"
)

func TestButton(t *testing.T) {
    tests := []struct {
        name  string
        props Props
    }{
        {"primary", Props{Label: "Save", Primary: true}},
        {"disabled", Props{Label: "Save", Disabled: true}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 1. Create the component node
            component := Button(tt.props)
            
            // 2. Render it to a buffer
            var buf bytes.Buffer
            err := component.Render(&buf)
            if err != nil {
                t.Fatal(err)
            }
            
            // 3. Assert against the golden file
            golden.Assert(t, buf.String(), "fixtures/"+tt.name+".golden")
        })
    }
}
```

## Composition Pattern

Composition is achieved by calling component functions within other component functions. This is checked by the compiler.

```go
// mylib/usercard/usercard.go
package usercard

import (
    "mylib/button" // Import the button component
    g "[github.com/maragu/gomponents/el](https://github.com/maragu/gomponents/el)"
)

type Props struct {
    UserName    string
    ButtonProps button.Props // Embed child props
}

func UserCard(p Props) g.Node {
    return g.Div(
        ga.Class("card"),
        g.H2(g.Text(p.UserName)),
        
        // Just call the function. No magic.
        button.Button(p.ButtonProps), 
    )
}
```
