package notify_test

import (
	"bytes"
	"testing"

	"github.com/owner/driftctl-lite/internal/drift"
	"github.com/owner/driftctl-lite/internal/notify"
)

func makeResults(n int) []drift.Result {
	results := make([]drift.Result, n)
	for i := range results {
		results[i] = drift.Result{ResourceID: fmt.Sprintf("res-%d", i), Status: drift.StatusModified}
	}
	return results
}

func TestNotify_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewWithWriter(notify.ChannelStdout, &buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if got != "[notify] No drift detected.\n" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestNotify_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewWithWriter(notify.ChannelStderr, &buf)
	results := []drift.Result{
		{ResourceID: "vpc-1", Status: drift.StatusModified},
		{ResourceID: "sg-2", Status: drift.StatusMissing},
	}
	if err := n.Notify(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	want := "[notify] Drift detected: 2 resource(s) out of sync.\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNew_UnsupportedChannel(t *testing.T) {
	_, err := notify.New("slack")
	if err == nil {
		t.Fatal("expected error for unsupported channel")
	}
}

func TestNew_ValidChannels(t *testing.T) {
	for _, ch := range []notify.Channel{notify.ChannelStdout, notify.ChannelStderr} {
		_, err := notify.New(ch)
		if err != nil {
			t.Errorf("unexpected error for channel %s: %v", ch, err)
		}
	}
}
