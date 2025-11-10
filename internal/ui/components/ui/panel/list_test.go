package panel

import (
	"html/template"
	"testing"

	"gotest.tools/v3/golden"
)

func TestRenderListWithItems(t *testing.T) {
	props := ListProps{
		ListClass:    "custom",
		EmptyMessage: "No items",
		Items: []ListItemProps{
			{
				Class:   "first",
				Content: template.HTML("<span>Alpha</span>"),
			},
			{
				Class:   "second",
				Content: template.HTML("<span>Beta</span>"),
			},
		},
	}

	html, err := RenderList(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "list_with_items.golden")
}

func TestRenderListEmpty(t *testing.T) {
	props := ListProps{
		EmptyMessage: "Nothing here",
	}

	html, err := RenderList(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "list_empty.golden")
}

func TestRenderListWithoutWrapper(t *testing.T) {
	props := ListProps{
		Items: []ListItemProps{
			{
				DisableWrapper: true,
				Content:        template.HTML("<div class=\"raw-item\">Custom</div>"),
			},
		},
	}

	html, err := RenderList(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "list_without_wrapper.golden")
}
