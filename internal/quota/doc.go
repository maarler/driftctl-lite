// Package quota provides per-resource-type drift quota enforcement for
// driftctl-lite. It allows operators to cap the number of drifted resources
// of a given type that are surfaced in a single scan, either by flagging
// violations or by dropping the excess results entirely.
package quota
