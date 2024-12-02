package app

import (
	"html/template"
	"net/http"

	"github.com/javernus/quote-unquote/internal/handler"
)

func (a *App) loadRoutes(tmpl *template.Template) {
	quotebook := handler.New(a.logger, a.db, tmpl)

	files := http.FileServer(http.Dir("./static"))

	a.router.Handle("GET /static/", http.StripPrefix("/static", files))

	a.router.Handle("GET /{$}", http.HandlerFunc(quotebook.Home))

	a.router.Handle("POST /{$}", http.HandlerFunc(quotebook.Create))
}
