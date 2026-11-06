// Package handler 提供 HTTP 处理器。
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"blockchain-cryptography/internal/sandbox"
)

// RunRequest 运行接口入参。
type RunRequest struct {
	Code      string `json:"code"`
	TimeoutMS int64  `json:"timeout_ms"` // 可选，服务端有硬上限
}

// RunHandler 处理 POST /api/v1/run。
type RunHandler struct {
	runner sandbox.Runner
}

// NewRunHandler 创建运行 handler。
func NewRunHandler(runner sandbox.Runner) *RunHandler {
	return &RunHandler{runner: runner}
}

// HandleRun 编译并运行用户代码。
func (h *RunHandler) HandleRun(w http.ResponseWriter, r *http.Request) {
	var req RunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Code == "" {
		writeErr(w, http.StatusBadRequest, "code is required")
		return
	}

	timeout := time.Duration(req.TimeoutMS) * time.Millisecond
	res, err := h.runner.Run(req.Code, timeout)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
