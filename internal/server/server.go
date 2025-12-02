package server

import (
	"net/http"
	"strings"

	"github.com/keircn/kcst/internal/config"
	"github.com/keircn/kcst/internal/handlers"
	"github.com/keircn/kcst/internal/storage"
	"github.com/keircn/kcst/internal/templates"
	"github.com/keircn/kcst/internal/upload"
)

type Server struct {
	addr        string
	mux         *http.ServeMux
	server      *http.Server
	db          *storage.DB
	store       *upload.Store
	stopCleanup chan struct{}
}

func New(cfg *config.Config) (*Server, error) {
	mux := http.NewServeMux()

	storage.SetRetention(storage.RetentionConfig{
		MinTTL:      cfg.Retention.MinTTL,
		MaxTTL:      cfg.Retention.MaxTTL,
		MaxFileSize: cfg.Retention.MaxFileSize,
	})

	db, err := storage.Open(cfg.Storage.DBPath)
	if err != nil {
		return nil, err
	}

	tmpl := templates.New()
	store, err := upload.NewStore(cfg.Storage.UploadDir, db, cfg.Retention.CleanupInterval)
	if err != nil {
		db.Close()
		return nil, err
	}
	h := handlers.New(tmpl, store, cfg.Server.BaseURL)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			h.Root(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/f/") {
			h.Preview(w, r)
		} else {
			h.ServeFile(w, r)
		}
	})

	return &Server{
		addr:        cfg.Server.Address,
		mux:         mux,
		db:          db,
		store:       store,
		stopCleanup: make(chan struct{}),
		server: &http.Server{
			Addr:    cfg.Server.Address,
			Handler: mux,
		},
	}, nil
}

func (s *Server) Run() error {
	s.store.StartCleanupRoutine(s.stopCleanup)
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	close(s.stopCleanup)
	return s.db.Close()
}
