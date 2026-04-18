package config

import (
	"flag"
	"fmt"
)

// Flags holds parsed CLI flag values.
type Flags struct {
	ConfigFile string
	StateFile  string
	LiveSource string
	OutputFmt  string
	FilterType string
	OnlyDrift  bool
}

// ParseFlags parses command-line flags and returns a Flags struct.
func ParseFlags(args []string) (*Flags, error) {
	fs := flag.NewFlagSet("driftctl-lite", flag.ContinueOnError)
	f := &Flags{}
	fs.StringVar(&f.ConfigFile, "config", ".driftctl.json", "path to config file")
	fs.StringVar(&f.StateFile, "state", "", "override state file path")
	fs.StringVar(&f.LiveSource, "live", "", "override live source path")
	fs.StringVar(&f.OutputFmt, "output", "", "output format: text or json")
	fs.StringVar(&f.FilterType, "type", "", "filter by resource type")
	fs.BoolVar(&f.OnlyDrift, "only-drift", false, "show only drifted resources")
	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}
	return f, nil
}

// Merge applies non-zero flag values on top of a Config.
func Merge(cfg *Config, f *Flags) *Config {
	out := *cfg
	if f.StateFile != "" {
		out.StateFile = f.StateFile
	}
	if f.LiveSource != "" {
		out.LiveSource = f.LiveSource
	}
	if f.OutputFmt != "" {
		out.OutputFmt = f.OutputFmt
	}
	if f.FilterType != "" {
		out.FilterType = f.FilterType
	}
	if f.OnlyDrift {
		out.OnlyDrift = true
	}
	return &out
}
