package handlers

import (
	"fmt"
	"net/http"

	"github.com/keircn/kcst/internal/models"
	"github.com/keircn/kcst/internal/templates"
	"github.com/keircn/kcst/internal/upload"
)

type Handler struct {
	templates *templates.Templates
	store     *upload.Store
}

func New(t *templates.Templates, s *upload.Store) *Handler {
	return &Handler{
		templates: t,
		store:     s,
	}
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.home(w, r)
	case http.MethodPost:
		h.upload(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	data := models.PageData{
		Title:   "KCST",
		Message: "Temporary file hosting.",
		BaseURL: baseURL,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.RenderPage(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, err := h.store.Save(file, header)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s\n", filename)
}
