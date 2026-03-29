package releaseinstall

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type InstallStatus struct {
	Mode             string         `json:"mode"`
	ExecutablePath   string         `json:"executable_path"`
	ReleaseRoot      string         `json:"release_root,omitempty"`
	ManifestPath     string         `json:"manifest_path,omitempty"`
	ManifestReadable bool           `json:"manifest_readable"`
	Manifest         map[string]any `json:"manifest,omitempty"`
}

func Detect() InstallStatus {
	exePath, err := os.Executable()
	if err != nil {
		return InstallStatus{Mode: "unknown"}
	}
	exePath, _ = filepath.EvalSymlinks(exePath)

	status := InstallStatus{
		Mode:           detectMode(exePath),
		ExecutablePath: exePath,
	}
	if status.Mode == "release-installed" {
		releaseRoot := filepath.Dir(filepath.Dir(exePath))
		manifestPath := filepath.Join(releaseRoot, "manifest.json")
		status.ReleaseRoot = releaseRoot
		status.ManifestPath = manifestPath
		manifest, ok := readManifest(manifestPath)
		status.ManifestReadable = ok
		if ok {
			status.Manifest = manifest
		}
	}
	return status
}

func detectMode(exePath string) string {
	clean := filepath.ToSlash(exePath)
	switch {
	case strings.Contains(clean, "/dist/releases/") && strings.HasSuffix(clean, "/bin/aih"):
		return "release-installed"
	case strings.Contains(clean, "/go-build") || strings.Contains(clean, "/tmp/go-build") || strings.Contains(clean, "/T/go-build"):
		return "dev-go-run"
	default:
		return "custom-binary"
	}
}

func readManifest(path string) (map[string]any, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, false
	}
	return manifest, true
}
