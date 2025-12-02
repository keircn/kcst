package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/keircn/kcst/internal/models"
	"github.com/keircn/kcst/internal/storage"
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
		Title:   "kcst",
		Message: "Temporary file hosting.",
		BaseURL: h.getBaseURL(r),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.RenderPage(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if err := r.ParseMultipartForm(100 << 20); err != nil {
		h.jsonError(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.jsonError(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, meta, err := h.store.Save(file, header)
	if err != nil {
		h.jsonError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	baseURL := h.getBaseURL(r)
	responseMS := time.Since(start).Milliseconds()

	resp := models.UploadResponse{
		Success:      true,
		URL:          fmt.Sprintf("%s/%s", baseURL, filename),
		RawURL:       fmt.Sprintf("%s/%s", baseURL, filename),
		PreviewURL:   fmt.Sprintf("%s/f/%s", baseURL, filename),
		Filename:     filename,
		OriginalName: meta.OriginalName,
		Size:         meta.Size,
		SizeHuman:    formatSize(meta.Size),
		ContentType:  meta.ContentType,
		UploadedAt:   meta.UploadedAt,
		ExpiresAt:    meta.ExpiresAt(),
		RetentionMS:  storage.CalculateTTL(meta.Size).Milliseconds(),
		ResponseMS:   responseMS,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   message,
	})
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

func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/f/")
	if filename == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	filename = filepath.Base(filename)

	_, meta, err := h.store.Get(filename)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	baseURL := h.getBaseURL(r)
	mediaType := getMediaType(meta.ContentType)

	data := models.FilePreviewData{
		Title:        fmt.Sprintf("%s - kcst", meta.OriginalName),
		Description:  fmt.Sprintf("Expires %s", meta.ExpiresAt().Format("Jan 02, 2006")),
		Filename:     filename,
		OriginalName: meta.OriginalName,
		Size:         meta.Size,
		SizeHuman:    formatSize(meta.Size),
		ContentType:  meta.ContentType,
		MediaType:    mediaType,
		RawURL:       fmt.Sprintf("%s/%s", baseURL, filename),
		PreviewURL:   fmt.Sprintf("%s/f/%s", baseURL, filename),
		BaseURL:      baseURL,
		UploadedAt:   meta.UploadedAt,
		ExpiresAt:    meta.ExpiresAt(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.RenderPreview(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getMediaType(contentType string) string {
	ct := strings.ToLower(contentType)
	if strings.HasPrefix(ct, "image/") {
		return "image"
	}
	if strings.HasPrefix(ct, "video/") {
		return "video"
	}
	return ""
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
