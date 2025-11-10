package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"composer/internal/ui/templates"
)

//go:embed static/*
var staticFS embed.FS

type Mode string

const (
	ModeProduction  Mode = "production"
	ModeDevelopment Mode = "development"
)

type Server struct {
	static fs.FS
}

func newServer(mode Mode) (*Server, error) {
	templates.Configure(uiSourceRoot())

	provider, err := newProvider(mode)
	if err != nil {
		return nil, err
	}

	return &Server{
		static: provider.staticFS,
	}, nil
}

// Init configures the UI server for the supplied mode.
func Init(mode Mode) (*Server, error) {
	return newServer(mode)
}

type provider struct {
	staticFS fs.FS
}

func newProvider(mode Mode) (*provider, error) {
	switch mode {
	case ModeProduction:
		static, err := fs.Sub(staticFS, "static")
		if err != nil {
			return nil, fmt.Errorf("ui: load embedded static assets: %w", err)
		}
		return &provider{staticFS: static}, nil
	case ModeDevelopment:
		root := uiSourceRoot()
		staticDir := filepath.Join(root, "static")
		if _, err := os.Stat(staticDir); err != nil {
			return nil, fmt.Errorf("ui: static directory not found at %q: %w", staticDir, err)
		}
		return &provider{staticFS: os.DirFS(staticDir)}, nil
	default:
		return nil, fmt.Errorf("ui: unsupported mode %q", mode)
	}
}

func uiSourceRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "internal/ui"
	}
	return filepath.Dir(filename)
}
