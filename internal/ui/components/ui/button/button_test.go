package button

import (
	"testing"

	"gotest.tools/v3/golden"
)

func TestRenderButton(t *testing.T) {
	cases := []struct {
		name  string
		props Props
	}{
		{
			name: "default",
			props: Props{
				Label: "Add",
			},
		},
		{
			name: "custom",
			props: Props{
				ID:        "add-step",
				Class:     "button--accent button--sm",
				Title:     "Create",
				AriaLabel: "Create item",
				Label:     "Create",
				Type:      "submit",
				IconSize:  20,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			html, err := Render(tt.props)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}

			golden.Assert(t, string(html), tt.name+".golden")
		})
	}
}
