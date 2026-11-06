// Package sandbox 负责把用户提交的 Go 代码编译并运行。
// 默认使用执行器进程（exec，适合本地/内网演示），可选 Docker 容器后端加固。
package sandbox

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// 运行硬限制（服务端兜底，不信任前端参数）。
const (
	maxTimeout = 5 * time.Second // 单次运行超时，达到即强杀
	buildLimit = 12 * time.Second
	maxOutput  = 1 << 20 // 输出截断上限 1MB
	maxCodeLen = 1 << 15 // 源码长度上限 32KB
)

// Result 返回一次运行的输出与统计。
type Result struct {
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	ExitCode  int    `json:"exit_code"`
	ElapsedMS int64  `json:"elapsed_ms"`
}

// Runner 抽象不同的沙箱后端。
type Runner interface {
	Run(code string, timeout time.Duration) (*Result, error)
}

// New 按后端名创建 Runner："docker" 走容器，否则走本地进程执行。
func New(backend string) Runner {
	if strings.EqualFold(backend, "docker") {
		return &dockerRunner{}
	}
	return &execRunner{}
}

// execRunner 本地进程执行：写入临时 main.go → go build → 限时运行。
type execRunner struct{}

func (e *execRunner) Run(code string, timeout time.Duration) (*Result, error) {
	if len(code) > maxCodeLen {
		return nil, errors.New("code too long")
	}
	if timeout <= 0 || timeout > maxTimeout {
		timeout = maxTimeout
	}

	dir, err := os.MkdirTemp("", "bkcc-sandbox-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	mainPath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainPath, []byte(code), 0o600); err != nil {
		return nil, err
	}
	binPath := filepath.Join(dir, "app")

	buildCtx, cancelBuild := context.WithTimeout(context.Background(), buildLimit)
	defer cancelBuild()
	build := exec.CommandContext(buildCtx, "go", "build", "-o", binPath, mainPath)
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		return &Result{Stderr: trim(string(out)), ExitCode: 2}, nil // 编译错误
	}

	start := time.Now()
	runCtx, cancelRun := context.WithTimeout(context.Background(), timeout)
	defer cancelRun()
	run := exec.CommandContext(runCtx, binPath)
	run.Dir = dir
	var stdout, stderr bytes.Buffer
	run.Stdout = &stdout
	run.Stderr = &stderr
	runErr := run.Run()
	elapsed := time.Since(start).Milliseconds()

	res := &Result{Stdout: trim(stdout.String()), Stderr: trim(stderr.String()), ElapsedMS: elapsed}
	if runErr != nil {
		exitCode := -1
		var ee *exec.ExitError
		if errors.As(runErr, &ee) {
			exitCode = ee.ExitCode()
		} else if runCtx.Err() != nil {
			exitCode = 124 // 超时
		}
		res.ExitCode = exitCode
	}
	return res, nil
}

// trim 截断过大输出，避免拖垮内存。
func trim(s string) string {
	if len(s) > maxOutput {
		return s[:maxOutput] + "\n...[truncated]"
	}
	return s
}
