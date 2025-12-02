package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/keircn/kcst/internal/models"
	"github.com/keircn/kcst/internal/templates"
	"github.com/keircn/kcst/internal/upload"
)

type Handler struct {
	templates *templates.Templates
	store     *upload.Store
	baseURL   string
}

func New(t *templates.Templates, s *upload.Store, baseURL string) *Handler {
	return &Handler{
		templates: t,
		store:     s,
		baseURL:   strings.TrimSuffix(baseURL, "/"),
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

func (h *Handler) getBaseURL(r *http.Request) string {
	if h.baseURL != "" {
		return h.baseURL
	}
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		Title:   "KCST",
		Message: "Temporary file hosting.",
		BaseURL: h.getBaseURL(r),
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
	fmt.Fprintf(w, "%s/%s\n", h.getBaseURL(r), filename)
}

func (h *Handler) ServeFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/")
	if filename == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	filename = filepath.Base(filename)

	file, meta, err := h.store.Get(filename)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", meta.Size))
	http.ServeContent(w, r, filename, meta.UploadedAt, file)
}
