package templates

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

// DevMode indicates whether templates should be reloaded from disk on every
// render. It is controlled via the APP_ENV environment variable.
var DevMode = os.Getenv("APP_ENV") == "development"

var projectRoot atomic.Pointer[string]

func init() {
	cwd := "."
	projectRoot.Store(&cwd)
}

// Configure overrides the root directory used to locate template files in
// development mode. Passing an empty path leaves the existing value unchanged.
func Configure(root string) {
	if strings.TrimSpace(root) == "" {
		return
	}
	path := root
	projectRoot.Store(&path)
}

func currentRoot() string {
	if v := projectRoot.Load(); v != nil {
		return *v
	}
	return "."
}

// Manager handles parsing templates for a specific component.
type Manager struct {
	name     string
	relPath  string
	embedded string
	funcs    template.FuncMap

	once   sync.Once
	cached *template.Template
	err    error
}

// New constructs a Manager for a template with the provided metadata.
func New(name, relPath, embedded string, funcs template.FuncMap) *Manager {
	var fn template.FuncMap
	if len(funcs) > 0 {
		fn = make(template.FuncMap, len(funcs))
		for k, v := range funcs {
			fn[k] = v
		}
	}
	return &Manager{
		name:     name,
		relPath:  relPath,
		embedded: embedded,
		funcs:    fn,
	}
}

func (m *Manager) parse(content string) (*template.Template, error) {
	t := template.New(m.name)
	if len(m.funcs) > 0 {
		t = t.Funcs(m.funcs)
	}
	return t.Parse(content)
}

func (m *Manager) getTemplate() (*template.Template, error) {
	if DevMode {
		path := filepath.Join(currentRoot(), m.relPath)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("templates: read %s: %w", path, err)
		}
		return m.parse(string(data))
	}

	m.once.Do(func() {
		m.cached, m.err = m.parse(m.embedded)
	})
	if m.err != nil {
		return nil, m.err
	}
	return m.cached, nil
}

// MustGet returns the parsed template or panics if parsing fails.
func (m *Manager) MustGet() *template.Template {
	tmpl, err := m.getTemplate()
	if err != nil {
		panic(err)
	}
	return tmpl
}

// Execute renders the template into the provided writer.
func (m *Manager) Execute(w io.Writer, data any) error {
	tmpl, err := m.getTemplate()
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// Render returns the rendered template output as template.HTML.
func (m *Manager) Render(data any) (template.HTML, error) {
	tmpl, err := m.getTemplate()
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// SafeHTML returns html if err is nil, otherwise an HTML comment describing the
// failure. This helps parent components recover gracefully when child renders
// fail.
func SafeHTML(html template.HTML, err error) template.HTML {
	if err != nil {
		return template.HTML("<!-- render error: " + template.HTMLEscapeString(err.Error()) + " -->")
	}
	return html
}
