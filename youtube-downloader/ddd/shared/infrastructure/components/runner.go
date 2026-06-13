// Package components holds shared, cross-cutting infrastructure helpers with no
// external dependencies. Runner locates and runs the bundled external binaries
// (yt-dlp, ffmpeg, ffprobe): the binaries live in a bin/ directory next to the
// project or the compiled executable; nothing is installed system-wide and no
// Python is needed.
package components

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// FindBinary locates a bundled binary by name (without extension), searching
// bin/ directories near the working directory and the executable, walking a few
// levels up, and finally the system PATH.
func FindBinary(name string) (string, error) {
	if runtime.GOOS == "windows" && filepath.Ext(name) == "" {
		name += ".exe"
	}

	var roots []string
	if wd, err := os.Getwd(); err == nil {
		roots = append(roots, wd)
	}
	if exe, err := os.Executable(); err == nil {
		roots = append(roots, filepath.Dir(exe))
	}

	seen := map[string]bool{}
	for _, root := range roots {
		dir := root
		for i := 0; i < 5; i++ {
			for _, cand := range []string{
				filepath.Join(dir, "bin", name),
				filepath.Join(dir, name),
			} {
				if seen[cand] {
					continue
				}
				seen[cand] = true
				if fi, err := os.Stat(cand); err == nil && !fi.IsDir() {
					return cand, nil
				}
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("binary %q not found in bin/ or PATH", name)
}

// Runner executes a single resolved binary.
type Runner struct {
	Path string
}

// NewRunner resolves a binary by name and returns a Runner for it.
func NewRunner(name string) (*Runner, error) {
	p, err := FindBinary(name)
	if err != nil {
		return nil, err
	}
	return &Runner{Path: p}, nil
}

// Run executes the binary with args and returns the combined stdout+stderr.
func (r *Runner) Run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, r.Path, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}
