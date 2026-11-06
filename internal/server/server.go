// Package server 组装 HTTP 服务：路由、静态前端托管、限流中间件。
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

// Options 服务配置。
type Options struct {
	// SandboxBackend 沙箱后端：空 = exec(默认，本地进程)，"docker" = Docker 容器。
	SandboxBackend string
}

// Server 聚合各业务组件并暴露路由。
type Server struct {
	srv   *service.Service
	run   *handler.RunHandler
	limit *sandbox.RateLimiter
}

// New 创建服务实例，加载内嵌教材并准备沙箱执行器。
func New(opts Options) (*Server, error) {
	svc, err := service.New(book.Content, book.Examples)
	if err != nil {
		return nil, err
	}
	return &Server{
		srv:   svc,
		run:   handler.NewRunHandler(sandbox.New(opts.SandboxBackend)),
		limit: sandbox.NewRateLimiter(30, 60), // 每个 IP 每分钟 30 次
	}, nil
}

// Router 返回 HTTP 处理器。
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	api := http.NewServeMux()

	api.HandleFunc("GET /api/v1/chapters", s.handleChapters)
	api.HandleFunc("GET /api/v1/chapters/{id}", s.handleChapter)
	api.HandleFunc("GET /api/v1/meta", s.handleMeta)
	api.HandleFunc("GET /api/v1/health", s.handleHealth)
	api.Handle("POST /api/v1/run", s.limit.Wrap(s.run.HandleRun))

	mux.Handle("/api/", api)

	// 前端静态资源（SPA），缺失时优雅降级
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
