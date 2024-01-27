package linter

import (
	"testing"
)

func TestNewLinter(t *testing.T) {
	workdir := "/path/to/workdir"
	output := OutputJson
	linter := New(workdir, output)

	if linter.Workdir != workdir {
		t.Errorf("Expected workdir to be %v, got %v", workdir, linter.Workdir)
	}

	if linter.Output != output {
		t.Errorf("Expected output to be the same instance as the input")
	}

	if len(linter.Rules) != len(ruleFuncs) {
		t.Fatalf("Expected linter to have %d rules, got %d", len(ruleFuncs), len(linter.Rules))
	}
}
