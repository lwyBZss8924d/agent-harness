package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstalledRuntimeDirFromExecutable(t *testing.T) {
	tmp := t.TempDir()
	releaseRoot := filepath.Join(tmp, "dist", "releases", "0.0.1-dev.1")
	binDir := filepath.Join(releaseRoot, "bin")
	runtimeDir := filepath.Join(releaseRoot, "runtime")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("MkdirAll bin: %v", err)
	}
	if err := os.MkdirAll(runtimeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll runtime: %v", err)
	}
	exe := filepath.Join(binDir, "aih")
	if err := os.WriteFile(exe, []byte(""), 0o755); err != nil {
		t.Fatalf("WriteFile exe: %v", err)
	}

	restoreExe := executablePath
	restoreEval := evalSymlinks
	executablePath = func() (string, error) { return exe, nil }
	evalSymlinks = func(path string) (string, error) { return path, nil }
	defer func() {
		executablePath = restoreExe
		evalSymlinks = restoreEval
	}()

	got := InstalledRuntimeDir()
	if got != runtimeDir {
		t.Fatalf("InstalledRuntimeDir() = %q, want %q", got, runtimeDir)
	}
}

func TestInstalledRuntimeDirEmptyWhenRuntimeMissing(t *testing.T) {
	tmp := t.TempDir()
	binDir := filepath.Join(tmp, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("MkdirAll bin: %v", err)
	}
	exe := filepath.Join(binDir, "aih")
	if err := os.WriteFile(exe, []byte(""), 0o755); err != nil {
		t.Fatalf("WriteFile exe: %v", err)
	}

	restoreExe := executablePath
	restoreEval := evalSymlinks
	executablePath = func() (string, error) { return exe, nil }
	evalSymlinks = func(path string) (string, error) { return path, nil }
	defer func() {
		executablePath = restoreExe
		evalSymlinks = restoreEval
	}()

	if got := InstalledRuntimeDir(); got != "" {
		t.Fatalf("InstalledRuntimeDir() = %q, want empty", got)
	}
}
