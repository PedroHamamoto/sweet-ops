package ui

import (
	"html/template"
	"path/filepath"
	"sync"
)

var (
	cache = map[string]*template.Template{}
	mu    sync.RWMutex
)

func load(layout string, page string) *template.Template {
	key := layout + ":" + page

	mu.RLock()
	t, ok := cache[key]
	mu.RUnlock()
	if ok {
		return t
	}

	mu.Lock()
	defer mu.Unlock()

	t = parse(layout, page)
	cache[key] = t
	return t
}

func parse(layout string, page string) *template.Template {
	files := []string{
		"internal/ui/templates/layouts/base.html",
		"internal/ui/templates/pages/" + page + ".html",
	}
	globalPartials, _ := filepath.Glob("internal/ui/templates/partials/*.html")
	files = append(files, globalPartials...)

	pagePartials, _ := filepath.Glob("internal/ui/templates/pages/" + page + "/partials/*.html")
	files = append(files, pagePartials...)

	return template.Must(template.ParseFiles(files...))
}
