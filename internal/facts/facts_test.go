package facts

import "testing"

func TestPathAnomalies(t *testing.T) {
	issues := pathAnomalies([]string{"/definitely/missing/path", "~/bad"})
	if len(issues) != 2 {
		t.Fatalf("len(issues) = %d", len(issues))
	}
}
