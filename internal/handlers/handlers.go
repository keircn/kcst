package handlers

import (
	"net/http"

	"github.com/keircn/kcst/internal/models"
	"github.com/keircn/kcst/internal/templates"
)

type Handler struct {
	templates *templates.Templates
}

func New(t *templates.Templates) *Handler {
	return &Handler{
		templates: t,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		Title:   "skibidi",
		Message: "sigma sigma",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.RenderPage(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
