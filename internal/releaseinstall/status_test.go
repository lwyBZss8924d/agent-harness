package releaseinstall

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectMode(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"/tmp/project/dist/releases/0.0.1/bin/aih", "release-installed"},
		{"/private/var/folders/x/T/go-build123/b001/exe/aih", "dev-go-run"},
		{"/tmp/custom/aih", "custom-binary"},
	}
	for _, tc := range cases {
		if got := detectMode(tc.path); got != tc.want {
			t.Fatalf("detectMode(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestReadManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(path, []byte(`{"version":"0.0.1-dev.2"}`), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	manifest, ok := readManifest(path)
	if !ok {
		t.Fatal("readManifest returned !ok")
	}
	if got := manifest["version"]; got != "0.0.1-dev.2" {
		t.Fatalf("manifest version = %#v", got)
	}
}
