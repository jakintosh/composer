package ui

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

//go:embed templates/*.tmpl templates/components/**/*.tmpl templates/pages/*.tmpl
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

type Mode string

const (
	ModeProduction  Mode = "production"
	ModeDevelopment Mode = "development"
)

type Server struct {
	renderer *Renderer
	static   fs.FS
}

// Renderer returns the template renderer backing this UI server.
func (s *Server) Renderer() *Renderer {
	return s.renderer
}

func newServer(mode Mode, funcs template.FuncMap) (*Server, error) {
	provider, err := newProvider(mode)
	if err != nil {
		return nil, err
	}

	renderer, err := newRenderer(provider.templateFS, funcs, provider.dev)
	if err != nil {
		return nil, err
	}

	return &Server{
		renderer: renderer,
		static:   provider.staticFS,
	}, nil
}

// Init configures the UI server for the supplied mode.
func Init(mode Mode) (*Server, error) {
	return newServer(mode, nil)
}

type provider struct {
	templateFS fs.FS
	staticFS   fs.FS
	dev        bool
}

func newProvider(mode Mode) (*provider, error) {
	switch mode {
	case ModeProduction:
		static, err := fs.Sub(staticFS, "static")
		if err != nil {
			return nil, fmt.Errorf("ui: load embedded static assets: %w", err)
		}
		return &provider{
			templateFS: templateFS,
			staticFS:   static,
			dev:        false,
		}, nil
	case ModeDevelopment:
		root := uiSourceRoot()
		templateDir := filepath.Join(root, "templates")
		if _, err := os.Stat(templateDir); err != nil {
			return nil, fmt.Errorf("ui: templates directory not found at %q: %w", templateDir, err)
		}
		staticDir := filepath.Join(root, "static")
		if _, err := os.Stat(staticDir); err != nil {
			return nil, fmt.Errorf("ui: static directory not found at %q: %w", staticDir, err)
		}
		return &provider{
			templateFS: os.DirFS(root),
			staticFS:   os.DirFS(staticDir),
			dev:        true,
		}, nil
	default:
		return nil, fmt.Errorf("ui: unsupported mode %q", mode)
	}
}

type Renderer struct {
	mu         sync.RWMutex
	templates  *template.Template
	funcs      template.FuncMap
	templateFS fs.FS
	dev        bool
}

func newRenderer(templateSource fs.FS, funcs template.FuncMap, dev bool) (*Renderer, error) {
	if funcs == nil {
		funcs = template.FuncMap{}
	}

	parsed, err := parseTemplates(templateSource, funcs)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		templates:  parsed,
		funcs:      copyFuncMap(funcs),
		templateFS: templateSource,
		dev:        dev,
	}, nil
}

func (r *Renderer) Page(w io.Writer, name string, data any) error {
	tpl, err := r.currentTemplates()
	if err != nil {
		return err
	}
	return tpl.ExecuteTemplate(w, name, data)
}

func (r *Renderer) currentTemplates() (*template.Template, error) {
	if r.dev {
		parsed, err := parseTemplates(r.templateFS, r.funcs)
		if err != nil {
			return nil, err
		}
		r.mu.Lock()
		r.templates = parsed
		r.mu.Unlock()
		return parsed, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.templates, nil
}

func parseTemplates(fsys fs.FS, funcs template.FuncMap) (*template.Template, error) {
	files := make([]string, 0)
	if err := fs.WalkDir(fsys, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".tmpl") {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no template files found")
	}

	sort.Strings(files)

	root := template.New("_root")
	if len(funcs) > 0 {
		root = root.Funcs(copyFuncMap(funcs))
	}
	return root.ParseFS(fsys, files...)
}

func copyFuncMap(in template.FuncMap) template.FuncMap {
	out := make(template.FuncMap, len(in))
	maps.Copy(out, in)
	return out
}

func uiSourceRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "internal/ui"
	}
	return filepath.Dir(filename)
}
