package alert_test

import (
	"bytes"
	"testing"

	"driftctl-lite/internal/alert"
	"driftctl-lite/internal/drift"
)

func makeResult(id, typ string, status drift.Status) drift.Result {
	return drift.Result{ID: id, Type: typ, Status: status}
}

func TestEvaluate_AllInSync(t *testing.T) {
	results := []drift.Result{
		makeResult("r1", "aws_s3_bucket", drift.StatusInSync),
		makeResult("r2", "aws_s3_bucket", drift.StatusInSync),
	}
	a := alert.Evaluate(results, alert.DefaultThresholds())
	if a.Level != alert.LevelOK {
		t.Errorf("expected OK, got %s", a.Level)
	}
	if a.Drifted != 0 {
		t.Errorf("expected 0 drifted, got %d", a.Drifted)
	}
}

func TestEvaluate_Warning(t *testing.T) {
	results := []drift.Result{
		makeResult("r1", "aws_s3_bucket", drift.StatusMissing),
		makeResult("r2", "aws_s3_bucket", drift.StatusInSync),
	}
	t2 := alert.Thresholds{Warning: 1, Critical: 5}
	a := alert.Evaluate(results, t2)
	if a.Level != alert.LevelWarning {
		t.Errorf("expected WARNING, got %s", a.Level)
	}
	if a.Drifted != 1 {
		t.Errorf("expected 1 drifted, got %d", a.Drifted)
	}
}

func TestEvaluate_Critical(t *testing.T) {
	results := []drift.Result{
		makeResult("r1", "aws_s3_bucket", drift.StatusMissing),
		makeResult("r2", "aws_s3_bucket", drift.StatusExtra),
		makeResult("r3", "aws_s3_bucket", drift.StatusModified),
		makeResult("r4", "aws_s3_bucket", drift.StatusMissing),
		makeResult("r5", "aws_s3_bucket", drift.StatusExtra),
	}
	t2 := alert.Thresholds{Warning: 1, Critical: 5}
	a := alert.Evaluate(results, t2)
	if a.Level != alert.LevelCritical {
		t.Errorf("expected CRITICAL, got %s", a.Level)
	}
	if a.Drifted != 5 {
		t.Errorf("expected 5 drifted, got %d", a.Drifted)
	}
}

func TestEvaluate_Empty(t *testing.T) {
	a := alert.Evaluate(nil, alert.DefaultThresholds())
	if a.Level != alert.LevelOK {
		t.Errorf("expected OK for empty results, got %s", a.Level)
	}
}

func TestFprint_Output(t *testing.T) {
	a := alert.Alert{Level: alert.LevelWarning, Message: "1 drifted resource(s) exceed warning threshold (1)", Drifted: 1}
	var buf bytes.Buffer
	alert.Fprint(&buf, a)
	got := buf.String()
	if got == "" {
		t.Error("expected non-empty output from Fprint")
	}
	expected := "[WARNING] 1 drifted resource(s) exceed warning threshold (1)\n"
	if got != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", got, expected)
	}
}
