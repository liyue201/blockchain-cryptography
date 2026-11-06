package sandbox

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// dockerRunner 通过 Docker 在隔离容器中编译运行用户代码（生产加固路径）。
type dockerRunner struct{}

// Run 将代码挂载进一次性容器执行：--network none、内存/CPU 限额、去能力。
func (d *dockerRunner) Run(code string, timeout time.Duration) (*Result, error) {
	if len(code) > maxCodeLen {
		return nil, errors.New("code too long")
	}
	if timeout <= 0 || timeout > maxTimeout {
		timeout = maxTimeout
	}

	dir, err := os.MkdirTemp("", "bkcc-docker-sandbox-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o600); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), buildLimit)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"--network", "none",
		"--memory", "128m",
		"--cpus", "0.5",
		"--pids-limit", "128",
		"--cap-drop", "ALL",
		"--read-only",
		"-v", dir+":/app:ro",
		"-w", "/app",
		"golang:1.25-alpine",
		"go", "run", "/app/main.go",
	)

	start := time.Now()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	elapsed := time.Since(start).Milliseconds()

	res := &Result{Stdout: trim(stdout.String()), Stderr: trim(stderr.String()), ElapsedMS: elapsed}
	if runErr != nil {
		exitCode := -1
		var ee *exec.ExitError
		if errors.As(runErr, &ee) {
			exitCode = ee.ExitCode()
		} else if ctx.Err() != nil {
			exitCode = 124
		}
		res.ExitCode = exitCode
	}
	return res, nil
}
