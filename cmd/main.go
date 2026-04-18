package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/live"
	"github.com/driftctl-lite/internal/output"
	"github.com/driftctl-lite/internal/state"
)

func main() {
	declaredPath := flag.String("state", "state.json", "path to declared state file")
	livePath := flag.String("live", "live.json", "path to live state file")
	format := flag.String("format", "text", "output format: text or json")
	flag.Parse()

	declared, err := state.LoadFromFile(*declaredPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading declared state: %v\n", err)
		os.Exit(1)
	}

	fetcher := live.NewFetcher(live.SourceFile, *livePath)
	liveResources, err := fetcher.Fetch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching live state: %v\n", err)
		os.Exit(1)
	}

	result := drift.Detect(declared.ResourceMap(), liveResources)

	reporter := output.NewReporter(*format, os.Stdout)
	if err := reporter.Report(result); err != nil {
		fmt.Fprintf(os.Stderr, "error reporting results: %v\n", err)
		os.Exit(1)
	}

	if result.HasDrift() {
		os.Exit(2)
	}
}
