// Package server assembles the HTTP service: routing, static frontend hosting, and rate-limit middleware.
package server

import (
	"encoding/json"
	"io/fs"
	"net/http"

	book "blockchain-cryptography"
	"blockchain-cryptography/internal/handler"
	"blockchain-cryptography/internal/sandbox"
	"blockchain-cryptography/internal/service"
)

// Options holds service configuration.
type Options struct {
	// SandboxBackend sandbox backend: empty = exec (default, local process), "docker" = Docker container.
	SandboxBackend string
}

// Server aggregates the business components and exposes the routes.
type Server struct {
	srv   *service.Service
	run   *handler.RunHandler
	limit *sandbox.RateLimiter
}

// New creates a service instance, loading the embedded content and preparing the sandbox executor.
func New(opts Options) (*Server, error) {
	svc, err := service.New(book.Content, book.Examples)
	if err != nil {
		return nil, err
	}
	return &Server{
		srv:   svc,
		run:   handler.NewRunHandler(sandbox.New(opts.SandboxBackend)),
		limit: sandbox.NewRateLimiter(30, 60), // 30 requests per IP per minute
	}, nil
}

// Router returns the HTTP handler.
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	api := http.NewServeMux()

	api.HandleFunc("GET /api/v1/chapters", s.handleChapters)
	api.HandleFunc("GET /api/v1/chapters/{id}", s.handleChapter)
	api.HandleFunc("GET /api/v1/meta", s.handleMeta)
	api.HandleFunc("GET /api/v1/health", s.handleHealth)
	api.Handle("POST /api/v1/run", s.limit.Wrap(s.run.HandleRun))

	mux.Handle("/api/", api)

	// Frontend static assets (SPA), gracefully degrades when missing
	webFS, err := fs.Sub(book.Web, "web/dist")
	if err != nil {
		panic("web/dist not embedded: " + err.Error())
	}
	mux.Handle("/", http.FileServer(http.FS(webFS)))

	return mux
}

func (s *Server) handleChapters(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.srv.List())
}

func (s *Server) handleMeta(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.srv.Meta())
}

func (s *Server) handleChapter(w http.ResponseWriter, r *http.Request) {
	ch, ok := s.srv.Get(r.PathValue("id"))
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "chapter not found"})
		return
	}
	writeJSON(w, http.StatusOK, ch)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
