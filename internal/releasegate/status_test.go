package releasegate

import "testing"

func TestEvaluateReleaseStatus(t *testing.T) {
	status := Evaluate("0.0.1-dev")
	if status.ReleaseTarget != "0.0.1" {
		t.Fatalf("ReleaseTarget = %q", status.ReleaseTarget)
	}
	if status.CurrentVersion != "0.0.1-dev" {
		t.Fatalf("CurrentVersion = %q", status.CurrentVersion)
	}
	if status.TotalGateCount == 0 {
		t.Fatal("expected gates")
	}
	if status.CompletedGateCount != status.TotalGateCount {
		t.Fatalf("completed %d total %d", status.CompletedGateCount, status.TotalGateCount)
	}
	if status.OverallStatus != StatusComplete {
		t.Fatalf("OverallStatus = %q", status.OverallStatus)
	}
	if len(status.FutureDomains) == 0 {
		t.Fatal("expected future domains")
	}
}
