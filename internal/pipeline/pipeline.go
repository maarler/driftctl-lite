// Package pipeline provides a composable processing pipeline for drift results.
// It allows chaining multiple transformation steps over a slice of DriftResult.
package pipeline

import "github.com/driftctl-lite/internal/drift"

// Stage is a function that transforms a slice of DriftResult.
type Stage func([]drift.Result) []drift.Result

// Pipeline holds an ordered list of processing stages.
type Pipeline struct {
	stages []Stage
}

// New creates an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Add appends a Stage to the pipeline.
func (p *Pipeline) Add(s Stage) *Pipeline {
	p.stages = append(p.stages, s)
	return p
}

// Run executes all stages in order, passing results through each one.
// If the input slice is nil, an empty slice is returned.
func (p *Pipeline) Run(results []drift.Result) []drift.Result {
	if results == nil {
		results = []drift.Result{}
	}
	for _, stage := range p.stages {
		results = stage(results)
	}
	return results
}

// Len returns the number of stages registered in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.stages)
}
