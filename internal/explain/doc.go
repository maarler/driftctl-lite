// Package explain translates raw drift detection results into human-readable
// explanations, providing actionable guidance for each type of drift detected.
//
// Usage:
//
//	results := drift.Detect(state, live)
//	explanations := explain.Explain(results)
//	for _, e := range explanations {
//		fmt.Println(e.Summary)
//		for _, d := range e.Details {
//			fmt.Println(d)
//		}
//	}
package explain
