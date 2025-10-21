package ui

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"maps"
	"os"
	"sort"
	"strings"
	"sync"
)

//go:embed templates/*.tmpl templates/components/**/*.tmpl templates/pages/*.tmpl
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

type Renderer struct {
	mu        sync.RWMutex
	templates *template.Template
	funcs     template.FuncMap
	dev       bool
}

func NewRenderer(funcs template.FuncMap) (*Renderer, error) {
	if funcs == nil {
		funcs = template.FuncMap{}
	}

	parsed, err := parseTemplates(templateFS, funcs)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		templates: parsed,
		funcs:     copyFuncMap(funcs),
		dev:       devMode(),
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
		parsed, err := parseTemplates(os.DirFS("internal/ui"), r.funcs)
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

func devMode() bool {
	return os.Getenv("DEV") == "1"
}

func staticFileSystem() (fs.FS, error) {
	if devMode() {
		return os.DirFS("internal/ui/static"), nil
	}
	return fs.Sub(staticFS, "static")
}
