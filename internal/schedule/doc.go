// Package schedule provides a simple interval-based scheduler for running
// periodic drift detection jobs. A Job encapsulates an interval and a handler
// function; the Scheduler ticks on that interval until the context is
// cancelled or the handler returns an error.
package schedule
