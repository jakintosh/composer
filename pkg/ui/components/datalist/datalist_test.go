package datalist

import (
	"html/template"
	"testing"

	"gotest.tools/v3/golden"
)

func TestRenderDataList(t *testing.T) {
	props := Props{
		Items: []Item{
			{Primary: "Alpha"},
			{
				Primary:   "Beta",
				Secondary: template.HTML("<span>Ready</span>"),
			},
		},
	}

	html, err := Render(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "datalist.golden")
}
