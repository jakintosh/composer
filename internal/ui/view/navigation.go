package view

// Sidebar represents the layout of the primary navigation sidebar.
type Sidebar struct {
	Title string
	Links []SidebarLink
}

// SidebarLink denotes a navigational item in the sidebar menu.
type SidebarLink struct {
	Label  string
	Href   string
	Active bool
}
