package aih

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestRunReleaseJSON(t *testing.T) {
	app := New()
	var out bytes.Buffer
	var errOut bytes.Buffer
	app.Stdout = &out
	app.Stderr = &errOut

	code := app.Run([]string{"release", "--json"})
	if code != 0 {
		t.Fatalf("Run returned %d stderr=%s", code, errOut.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload["release_target"] != "0.0.1" {
		t.Fatalf("release_target = %#v", payload["release_target"])
	}
}

func TestRunFactsPathJSON(t *testing.T) {
	app := New()
	var out bytes.Buffer
	var errOut bytes.Buffer
	app.Stdout = &out
	app.Stderr = &errOut

	code := app.Run([]string{"facts", "path", "--json"})
	if code != 0 {
		t.Fatalf("Run returned %d stderr=%s", code, errOut.String())
	}
	var payload map[string]string
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload["facts_json"] == "" || payload["facts_md"] == "" {
		t.Fatalf("payload = %#v", payload)
	}
}
