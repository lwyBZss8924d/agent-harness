package paths

import (
	"os"
	"path/filepath"
)

const (
	DefaultSourceRepo  = "dev-space/aih-toolkit"
	AgentsRoot         = ".agents"
	GeneratedFactsRoot = ".agents/state/facts"
)

var (
	executablePath = os.Executable
	evalSymlinks   = filepath.EvalSymlinks
)

func HomeDir() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	return ""
}

func SourceRepoRoot() string {
	if value := os.Getenv("AIH_SOURCE_ROOT"); value != "" {
		return value
	}
	return filepath.Join(HomeDir(), DefaultSourceRepo)
}

func InstalledRuntimeDir() string {
	exePath, err := executablePath()
	if err != nil {
		return ""
	}
	if resolved, err := evalSymlinks(exePath); err == nil {
		exePath = resolved
	}
	exePath = filepath.Clean(exePath)
	if filepath.Base(exePath) != "aih" {
		return ""
	}
	binDir := filepath.Dir(exePath)
	if filepath.Base(binDir) != "bin" {
		return ""
	}
	runtimeDir := filepath.Join(filepath.Dir(binDir), "runtime")
	if _, err := os.Stat(runtimeDir); err == nil {
		return runtimeDir
	}
	return ""
}

func AgentsDir() string {
	return filepath.Join(HomeDir(), AgentsRoot)
}

func LegacyRuntimeDir() string {
	candidate := filepath.Join(HomeDir(), AgentsRoot, "tools")
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return ""
}

func GeneratedFactsDir() string {
	return filepath.Join(HomeDir(), GeneratedFactsRoot)
}

func FactsJSONPath() string {
	return filepath.Join(GeneratedFactsDir(), "host-facts.json")
}

func FactsMarkdownPath() string {
	return filepath.Join(GeneratedFactsDir(), "host-facts.md")
}
