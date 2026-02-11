package ui

import (
	"net/http"
)

func Render(w http.ResponseWriter, r *http.Request, page string, data any) {
	layout := r.Context().Value(layoutKey).(string)

	t := load(layout, page)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, layout, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
