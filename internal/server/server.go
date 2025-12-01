package server

import (
	"net/http"

	"github.com/keircn/kcst/internal/handlers"
	"github.com/keircn/kcst/internal/templates"
	"github.com/keircn/kcst/internal/upload"
)

type Server struct {
	addr   string
	mux    *http.ServeMux
	server *http.Server
}

func New(addr, uploadDir string) (*Server, error) {
	mux := http.NewServeMux()
	tmpl := templates.New()
	store, err := upload.NewStore(uploadDir)
	if err != nil {
		return nil, err
	}
	h := handlers.New(tmpl, store)
	mux.HandleFunc("/", h.Root)

	return &Server{
		addr: addr,
		mux:  mux,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}, nil
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}
