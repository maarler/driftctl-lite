// Package notify provides drift notification hooks.
package notify

import (
	"fmt"
	"io"
	"os"

	"github.com/owner/driftctl-lite/internal/drift"
)

// Channel represents a notification output channel.
type Channel string

const (
	ChannelStdout Channel = "stdout"
	ChannelStderr Channel = "stderr"
)

// Notifier sends drift notifications to a configured channel.
type Notifier struct {
	channel Channel
	w       io.Writer
}

// New creates a Notifier for the given channel.
func New(channel Channel) (*Notifier, error) {
	var w io.Writer
	switch channel {
	case ChannelStdout:
		w = os.Stdout
	case ChannelStderr:
		w = os.Stderr
	default:
		return nil, fmt.Errorf("unsupported notify channel: %s", channel)
	}
	return &Notifier{channel: channel, w: w}, nil
}

// NewWithWriter creates a Notifier writing to w (useful for tests).
func NewWithWriter(channel Channel, w io.Writer) *Notifier {
	return &Notifier{channel: channel, w: w}
}

// Notify writes a summary notification if drift was detected.
func (n *Notifier) Notify(results []drift.Result) error {
	total := len(results)
	if total == 0 {
		_, err := fmt.Fprintln(n.w, "[notify] No drift detected.")
		return err
	}
	_, err := fmt.Fprintf(n.w, "[notify] Drift detected: %d resource(s) out of sync.\n", total)
	return err
}
