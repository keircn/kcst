package server

import (
	"net/http"

	"github.com/keircn/kcst/internal/handlers"
	"github.com/keircn/kcst/internal/templates"
)

type Server struct {
	addr   string
	mux    *http.ServeMux
	server *http.Server
}

func New(addr string) *Server {
	mux := http.NewServeMux()

	tmpl := templates.New()
	h := handlers.New(tmpl)

	mux.HandleFunc("/", h.Home)

	return &Server{
		addr: addr,
		mux:  mux,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}
