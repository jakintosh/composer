package button

import (
	"sort"
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Props describes the configuration for rendering a call-to-action button with
// an optional label and icon.
type Props struct {
	ID        string
	Class     string
	Title     string
	AriaLabel string
	Label     string
	Type      string
	IconSize  int
	HideIcon  bool
	Data      map[string]string
}

// Button renders the configured button as a gomponent tree.
func Button(p Props) g.Node {
	size := p.IconSize
	if size == 0 {
		size = 16
	}

	buttonType := p.Type
	if buttonType == "" {
		buttonType = "button"
	}

	className := "button"
	if extra := strings.TrimSpace(p.Class); extra != "" {
		className += " " + extra
	}

	nodes := []g.Node{
		html.Type(buttonType),
		html.Class(className),
	}
	if p.ID != "" {
		nodes = append(nodes, html.ID(p.ID))
	}
	if p.Title != "" {
		nodes = append(nodes, html.Title(p.Title))
	}
	if p.AriaLabel != "" {
		nodes = append(nodes, html.Aria("label", p.AriaLabel))
	}
	if len(p.Data) > 0 {
		keys := make([]string, 0, len(p.Data))
		for key := range p.Data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			name := strings.TrimSpace(key)
			if name == "" {
				continue
			}
			nodes = append(nodes, g.Attr("data-"+name, p.Data[key]))
		}
	}

	if !p.HideIcon {
		nodes = append(nodes, plusIcon(size))
	}
	if p.Label != "" {
		nodes = append(nodes, html.Span(g.Text(p.Label)))
	}

	return html.Button(nodes...)
}

func plusIcon(size int) g.Node {
	dimension := strconv.Itoa(size)
	return html.SVG(
		g.Attr("aria-hidden", "true"),
		g.Attr("focusable", "false"),
		html.Width(dimension),
		html.Height(dimension),
		g.Attr("viewBox", "0 0 16 16"),
		g.El("path",
			g.Attr("d", "M8 3v10M3 8h10"),
			g.Attr("stroke", "currentColor"),
			g.Attr("stroke-width", "2"),
			g.Attr("stroke-linecap", "round"),
		),
	)
}
