package testutil

import (
	"bytes"
	"testing"

	g "maragu.dev/gomponents"
)

// Render renders the gomponent node and fails the test on error.
func Render(t testing.TB, node g.Node) string {
	t.Helper()

	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}
