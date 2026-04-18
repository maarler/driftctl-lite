// Package diff computes field-level differences between declared infrastructure
// attributes and their live counterparts.
//
// It is used by the explain package to provide granular drift details, showing
// exactly which fields changed, were added, or were removed for a given resource.
//
// Usage:
//
//	diffs := diff.Compute(declaredAttrs, liveAttrs)
//	for _, d := range diffs {
//		fmt.Println(d)
//	}
package diff
