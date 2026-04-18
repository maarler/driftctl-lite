// Package notify implements lightweight drift notification hooks for
// driftctl-lite. It supports writing human-readable drift alerts to
// configured output channels (stdout, stderr) after a drift detection
// run completes.
//
// Usage:
//
//	n, err := notify.New(notify.ChannelStdout)
//	if err != nil { ... }
//	n.Notify(results)
package notify
