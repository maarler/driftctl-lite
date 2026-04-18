// Package audit provides append-only audit logging for driftctl-lite.
//
// Each time a drift detection run completes, an Entry is written as a
// JSON line to the configured audit log file. Entries can be read back
// with ReadAll for reporting or compliance purposes.
//
// Usage:
//
//	logger := audit.NewLogger("/var/log/driftctl/audit.log")
//	if err := logger.Record(stateFile, source, results); err != nil {
//	    log.Fatal(err)
//	}
package audit
