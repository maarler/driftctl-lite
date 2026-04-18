package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds CLI configuration options.
type Config struct {
	StateFile  string `json:"state_file"`
	LiveSource string `json:"live_source"`
	OutputFmt  string `json:"output_format"`
	FilterType string `json:"filter_type"`
	OnlyDrift  bool   `json:"only_drift"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		StateFile:  "state.json",
		LiveSource: "live.json",
		OutputFmt:  "text",
		FilterType: "",
		OnlyDrift:  false,
	}
}

// LoadFromFile reads a JSON config file and merges it into defaults.
func LoadFromFile(path string) (*Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return cfg, nil
}

// Validate checks that required fields are non-empty.
func (c *Config) Validate() error {
	if c.StateFile == "" {
		return fmt.Errorf("state_file must not be empty")
	}
	if c.LiveSource == "" {
		return fmt.Errorf("live_source must not be empty")
	}
	if c.OutputFmt != "text" && c.OutputFmt != "json" {
		return fmt.Errorf("output_format must be 'text' or 'json', got %q", c.OutputFmt)
	}
	return nil
}
