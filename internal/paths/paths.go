package paths

import (
	"os"
	"path/filepath"
)

const (
	DefaultSourceRepo  = "dev-space/aih-toolkit"
	LegacyRuntimeRoot  = ".agents/tools"
	AgentsRoot         = ".agents"
	GeneratedFactsRoot = ".agents/os-dev-environment/generated"
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

func AgentsDir() string {
	return filepath.Join(HomeDir(), AgentsRoot)
}

func LegacyRuntimeDir() string {
	return filepath.Join(HomeDir(), LegacyRuntimeRoot)
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
