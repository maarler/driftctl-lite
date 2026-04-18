package state_test

import (
	"os"
	"testing"

	"github.com/example/driftctl-lite/internal/state"
)

func writeTempState(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFromFile(t *testing.T) {
	path := writeTempState(t, `{"resources":[{"id":"vpc-1","type":"aws_vpc","attributes":{"cidr":"10.0.0.0/16"}}]}`)
	s, err := state.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(s.Resources))
	}
	if s.Resources[0].ID != "vpc-1" {
		t.Errorf("expected id vpc-1, got %s", s.Resources[0].ID)
	}
}

func TestLoadFromFile_Missing(t *testing.T) {
	_, err := state.LoadFromFile("/nonexistent/path.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestResourceMap(t *testing.T) {
	s := &state.State{
		Resources: []state.Resource{
			{ID: "sg-1", Type: "aws_security_group", Attributes: map[string]string{"name": "default"}},
			{ID: "sg-2", Type: "aws_security_group", Attributes: map[string]string{"name": "web"}},
		},
	}
	m := s.ResourceMap()
	if _, ok := m["sg-1"]; !ok {
		t.Error("expected sg-1 in map")
	}
	if len(m) != 2 {
		t.Errorf("expected map length 2, got %d", len(m))
	}
}
