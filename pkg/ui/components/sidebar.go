package components

import (
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// SidebarLink denotes a navigational item in the sidebar menu.
type SidebarLink struct {
	Label  string
	Href   string
	Active bool
}

// SidebarProps represents the layout of the primary navigation sidebar.
type SidebarProps struct {
	Title string
	Links []SidebarLink
}

// Sidebar renders the navigation chrome.
func Sidebar(p SidebarProps) g.Node {
	var title g.Node = g.Group(nil)
	if strings.TrimSpace(p.Title) != "" {
		title = html.Div(
			html.Class("sidebar__brand"),
			g.Text(strings.TrimSpace(p.Title)),
		)
	}

	links := make([]g.Node, 0, len(p.Links))
	for _, link := range p.Links {
		className := "sidebar__link"
		if link.Active {
			className += " sidebar__link--active"
		}
		links = append(links, html.Li(
			html.A(
				html.Class(className),
				g.If(strings.TrimSpace(link.Href) != "", html.Href(strings.TrimSpace(link.Href))),
				g.Text(link.Label),
			),
		))
	}

	return html.Div(
		html.Class("sidebar"),
		title,
		html.Nav(
			html.Class("sidebar__nav"),
			html.Ul(
				html.Class("sidebar__list"),
				g.Group(links),
			),
		),
	)
}
