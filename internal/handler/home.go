package handler

import (
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"strings"

	goaway "github.com/TwiN/go-away"
	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/x-way/crawlerdetect"

	"github.com/javernus/quote-unquote/internal/quote"
	"github.com/javernus/quote-unquote/internal/repository"
)

type Quotebook struct {
	logger *slog.Logger
	tmpl   *template.Template
	repo   *repository.Queries
}

func New(
	logger *slog.Logger, db *pgxpool.Pool, tmpl *template.Template,
) *Quotebook {
	return &Quotebook{
		tmpl:   tmpl,
		repo:   repository.New(db),
		logger: logger,
	}
}

type indexPage struct {
	Quotes []repository.Quote
	Total  int64
}

type errorPage struct {
	ErrorMessage string
}

func (h *Quotebook) Home(w http.ResponseWriter, r *http.Request) {
	quotes, err := h.repo.FindAll(r.Context(), 200)
	if err != nil {
		h.logger.Error("failed to find quotes", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	count, err := h.repo.Count(r.Context())
	if err != nil {
		h.logger.Error("failed to get count", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "index.html", indexPage{
		Quotes: quotes,
		Total:  count,
	})
}

func (h *Quotebook) Create(w http.ResponseWriter, r *http.Request) {
	// if crawlerdetect.IsCrawler(r.Header.Get("User-Agent")) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	//
	if err := r.ParseForm(); err != nil {
		h.logger.Error("failed to parse form", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg, ok := r.Form["message"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message := strings.Join(msg, " ")

	if strings.TrimSpace(message) == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "error.html", errorPage{
			ErrorMessage: "Blank messages don't count",
		})

		return
	}

	splits := strings.Split(r.RemoteAddr, ":")
	ipStr := strings.Trim(strings.Join(splits[:len(splits)-1], ":"), "[]")
	ip := net.ParseIP(ipStr)

	if goaway.IsProfane(message) {
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "error.html", errorPage{
			ErrorMessage: fmt.Sprintf(
				"Please don't use profanity. Your IP has been tracked %s",
				ipStr,
			),
		})
		return
	}

	psn, ok := r.Form["person"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	person := strings.Join(psn, " ")

	if strings.TrimSpace(person) == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "error.html", errorPage{
			ErrorMessage: "Blank persons don't count",
		})

		return
	}

	if goaway.IsProfane(person) {
		w.WriteHeader(http.StatusBadRequest)
		h.tmpl.ExecuteTemplate(w, "error.html", errorPage{
			ErrorMessage: fmt.Sprintf(
				"Please don't use profanity. Your IP has been tracked %s",
				ipStr,
			),
		})
		return
	}

	quote, err := quote.NewQuote(message, person, ip)
	if err != nil {
		h.logger.Error("failed to create quote", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.repo.Insert(r.Context(), repository.InsertParams{
		ID:        quote.ID,
		Message:   quote.Message,
		Person:    quote.Person,
		CreatedAt: quote.CreatedAt,
		Ip:        quote.IP,
	})
	if err != nil {
		h.logger.Error("failed to insert quote", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
